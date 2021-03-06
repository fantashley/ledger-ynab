package migrate

import (
	"errors"

	"github.com/fantashley/ledger-ynab/pkg/ledger"
	"github.com/fantashley/ledger-ynab/pkg/ynab"
)

func Convert(tx ledger.Transaction) (ynab.Transaction, error) {
	if len(tx.Payers) > 1 || len(tx.Categories) > 1 {
		return ynab.Transaction{}, errors.New("multiple payers and categories not supported")
	}

	return ynab.Transaction{
		Date:      tx.Date.Format(ynab.DateFormat),
		PayeeName: tx.Payee,
		Memo:      tx.Comment,
		Amount:    int64(tx.Categories[0].Amount.CentsTotal()) * 10,
	}, nil
}
