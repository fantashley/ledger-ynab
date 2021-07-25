package ynab

const DateFormat = "2006-01-02"

type Transaction struct {
	ID                    string        `json:"id,omitempty"`
	Date                  string        `json:"date,omitempty"`
	Amount                int64         `json:"amount,omitempty"`
	Memo                  string        `json:"memo,omitempty"`
	Cleared               string        `json:"cleared,omitempty"`
	Approved              bool          `json:"approved,omitempty"`
	FlagColor             string        `json:"flag_color,omitempty"`
	AccountID             string        `json:"account_id,omitempty"`
	PayeeID               string        `json:"payee_id,omitempty"`
	CategoryID            string        `json:"category_id,omitempty"`
	TransferAccountID     string        `json:"transfer_account_id,omitempty"`
	TransferTransactionID string        `json:"transfer_transaction_id,omitempty"`
	MatchedTransactionID  string        `json:"matched_transaction_id,omitempty"`
	ImportID              string        `json:"import_id,omitempty"`
	Deleted               bool          `json:"deleted,omitempty"`
	AccountName           string        `json:"account_name,omitempty"`
	PayeeName             string        `json:"payee_name,omitempty"`
	CategoryName          string        `json:"category_name,omitempty"`
	Subtransactions       []Transaction `json:"subtransactions,omitempty"`
}
