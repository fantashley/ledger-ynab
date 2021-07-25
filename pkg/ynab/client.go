package ynab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
)

const (
	APIPath      = "/v1"
	jsonMimeType = "application/json"
)

type Client struct {
	config     Config
	httpClient *http.Client
	apiURL     *url.URL
}

func NewClient(httpClient *http.Client, config Config) (*Client, error) {
	config = config.FillDefaults()

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	url, err := url.Parse(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
		apiURL:     url,
	}, nil
}

func (c *Client) CreateTransactions(ctx context.Context, transactions []Transaction) error {
	type RequestBody struct {
		Transactions []Transaction `json:"transactions"`
	}

	request := RequestBody{
		Transactions: transactions,
	}

	// var buf bytes.Buffer
	// if err := json.NewEncoder(&buf).Encode(request); err != nil {
	// 	return fmt.Errorf("failed to encode request: %w", err)
	// }

	// body := bytes.NewReader(buf.Bytes())

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Debugf("Create tx request body: %v", string(reqBytes))

	url := *c.apiURL
	url.Path = path.Join(url.Path, APIPath, "budgets", c.config.BudgetID, "transactions")
	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewReader(reqBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.config.PersonalAccessToken)
	req.Header.Add("Content-Type", jsonMimeType)
	log.Debugf("Request: %+v", req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create transaction request failed: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		var statusErr error
		statusErr = multierror.Append(statusErr, fmt.Errorf("received status code of %d", resp.StatusCode))

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return multierror.Append(statusErr, fmt.Errorf("error reading response body: %w", err))
		}
		defer resp.Body.Close()

		return multierror.Append(statusErr, fmt.Errorf("response body: %v", string(respBody)))
	}

	return nil
}

type ListAccountTransactionsOptions struct {
	SinceDate             time.Time
	Type                  string
	LastKnowledgeOfServer int64
}

func (o *ListAccountTransactionsOptions) Values() url.Values {
	vals := url.Values{}

	if !o.SinceDate.IsZero() {
		vals.Set("since_date", o.SinceDate.Format(DateFormat))
	}
	if o.Type != "" {
		vals.Set("type", o.Type)
	}
	if o.LastKnowledgeOfServer != 0 {
		vals.Set("last_knowledge_of_server", strconv.FormatInt(o.LastKnowledgeOfServer, 10))
	}

	return vals
}

func (c *Client) ListAccountTransactions(ctx context.Context, accountID string, options *ListAccountTransactionsOptions) ([]Transaction, error) {
	url := *c.apiURL

	if options != nil {
		url.RawQuery = options.Values().Encode()
	}

	url.Path = path.Join(url.Path, APIPath, "budgets", c.config.BudgetID, "accounts", accountID, "transactions")
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request for listing account transactions: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.config.PersonalAccessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error listing account transactions: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list account transactions returned code %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	defer resp.Body.Close()

	return response.Data.Transactions, nil
}
