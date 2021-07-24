// +build integration_test

package ynab_test

import (
	"context"
	"os"
	"testing"

	"github.com/fantashley/ledger-ynab/pkg/ynab"
)

func TestListTransactions(t *testing.T) {
	var (
		accessToken = os.Getenv("YNAB_API_TOKEN")
		budgetID    = os.Getenv("YNAB_BUDGET_ID")
		accountID   = os.Getenv("YNAB_ACCOUNT_ID")
	)

	config := ynab.Config{
		PersonalAccessToken: accessToken,
		BudgetID:            budgetID,
	}

	client, err := ynab.NewClient(nil, config)
	if err != nil {
		t.Fatalf("Failed to create ynab client: %v", err)
	}

	transactions, err := client.ListAccountTransactions(context.Background(), accountID, nil)
	if err != nil {
		t.Fatalf("Error listing account transactions: %v", err)
	}

	t.Logf("Found %d transactions", len(transactions))
}
