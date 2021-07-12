package ledger

import (
	"errors"
	"fmt"
	"time"

	"github.com/fantashley/ledger-ynab/pkg/usd"
	"github.com/hashicorp/go-multierror"
)

type Ledger struct {
	Transactions []Transaction
}

type Transaction struct {
	Date       time.Time
	Payee      string
	Comment    string
	Categories []Account
	Payers     []Account
}

type Account struct {
	Name   string
	Amount usd.USD
}

func New(filepath string) (Ledger, error) {
	txs, err := loadTransactions(filepath)
	if err != nil {
		return Ledger{}, fmt.Errorf("failed to parse ledger file: %w", err)
	}

	return Ledger{
		Transactions: txs,
	}, nil
}

func (t Transaction) Validate() error {
	var err error
	if t.Date.IsZero() {
		err = multierror.Append(err, errors.New("missing field 'Date'"))
	}
	if t.Payee == "" {
		err = multierror.Append(err, errors.New("missing field 'Payee'"))
	}
	if len(t.Categories)+len(t.Payers) == 0 {
		err = multierror.Append(err, errors.New("missing accounts"))
	}

	var sum usd.USD
	for _, account := range append(t.Categories, t.Payers...) {
		sum = sum.Add(account.Amount)
	}

	if sum.CentsTotal() != 0 {
		err = multierror.Append(err, fmt.Errorf("transaction amounts should balance to $0, instead got total of $%d.%2d", sum.Dollars, sum.Cents))
	}

	return err
}

func (t Transaction) Empty() bool {
	switch {
	case t.Date.IsZero() &&
		t.Payee == "" &&
		t.Comment == "" &&
		len(t.Categories)+len(t.Payers) == 0:
		return true
	}

	return false
}
