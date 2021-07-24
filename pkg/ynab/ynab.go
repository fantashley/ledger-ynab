package ynab

const DateFormat = "2006-01-02"

type Transaction struct {
	ID                    string        `json:"id"`
	Date                  string        `json:"date"`
	Amount                int64         `json:"amount"`
	Memo                  string        `json:"memo"`
	Cleared               string        `json:"cleared"`
	Approved              bool          `json:"approved"`
	FlagColor             string        `json:"flag_color"`
	AccountID             string        `json:"account_id"`
	PayeeID               string        `json:"payee_id"`
	CategoryID            string        `json:"category_id"`
	TransferAccountID     string        `json:"transfer_account_id"`
	TransferTransactionID string        `json:"transfer_transaction_id"`
	MatchedTransactionID  string        `json:"matched_transaction_id"`
	ImportID              string        `json:"import_id"`
	Deleted               bool          `json:"deleted"`
	AccountName           string        `json:"account_name"`
	PayeeName             string        `json:"payee_name"`
	CategoryName          string        `json:"category_name"`
	Subtransactions       []Transaction `json:"subtransactions"`
}
