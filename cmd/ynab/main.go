package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fantashley/ledger-ynab/internal/migrate"
	"github.com/fantashley/ledger-ynab/pkg/ledger"
	"github.com/fantashley/ledger-ynab/pkg/ynab"
	"github.com/peterbourgon/ff/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		fs          = flag.NewFlagSet("ynab", flag.ExitOnError)
		ledgerFile  = fs.String("ledger", "ledger.dat", "Path to ledger file")
		accessToken = fs.String("access-token", "", "YNAB personal access token")
		budget      = fs.String("budget-id", "default", "YNAB budget ID")
		accountID   = fs.String("account-id", "", "Account to create transactions in")
		accountName = fs.String("account-name", "", "Name of the YNAB account to create transactions in")
		debug       = fs.Bool("debug", false, "Set logging to debug level")
		_           = fs.String("config", "", "Path to config file")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("YNAB"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
	); err != nil {
		log.Panicf("Error parsing flags: %v", err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	ledger, err := ledger.New(*ledgerFile)
	if err != nil {
		log.Panicf("Failed to create ledger: %v", err)
	}

	ynabClient, err := ynab.NewClient(http.DefaultClient, ynab.Config{
		PersonalAccessToken: *accessToken,
		BudgetID:            *budget,
	})
	if err != nil {
		log.Fatalf("Error creating YNAB client: %v", err)
	}

	ctx := context.Background()
	ynabTxs, err := ynabClient.ListAccountTransactions(ctx, *accountID, &ynab.ListAccountTransactionsOptions{
		SinceDate: time.Now().AddDate(0, -1, 0),
	})
	if err != nil || len(ynabTxs) == 0 {
		log.Fatalf("Error listing account transactions: %v", err)
	}

	lastTx := ynabTxs[len(ynabTxs)-1]
	lastTxTime, err := time.Parse(ynab.DateFormat, lastTx.Date)
	if err != nil {
		log.Fatalf("Failed to parse time of most recent transaction: %v", err)
	}

	for i := len(ledger.Transactions) - 1; i >= 0 && ledger.Transactions[i].Date.After(lastTxTime); i-- {
		currentTx := ledger.Transactions[i]
		if len(currentTx.Payers) > 1 || len(currentTx.Categories) > 1 {
			log.Errorf("Transaction %+v not supported", currentTx)
			continue
		}

		ynabTx, err := migrate.Convert(currentTx)
		if err != nil {
			log.Errorf("Failed to convert transaction %+v: %v", currentTx, err)
			continue
		}

		if strings.HasSuffix(currentTx.Payers[0].Name, "Ashley") {
			ynabTx.Amount = ynabTx.Amount / 2
		} else if strings.HasSuffix(currentTx.Payers[0].Name, *accountName) {
			ynabTx.Amount = -1 * ynabTx.Amount / 2
		} else {
			log.Errorf("Unrecognized payer %q", currentTx.Payers[0])
			continue
		}

		ynabTx.AccountName = *accountName
		ynabTx.AccountID = *accountID
		ynabTx.Cleared = "cleared"

		if err = ynabClient.CreateTransactions(ctx, []ynab.Transaction{ynabTx}); err != nil {
			log.Errorf("Failed to create transaction %+v: %v", ynabTx, err)
		}
	}
}
