package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fantashley/ledger-ynab/pkg/ledger"
	"github.com/fantashley/ledger-ynab/pkg/ynab"
	"github.com/peterbourgon/ff/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		fs          = flag.NewFlagSet("ynab", flag.ExitOnError)
		ledgerFile  = fs.String("f", "ledger.dat", "Path to ledger file")
		accessToken = fs.String("ynab-access-token", "", "YNAB personal access token")
		budget      = fs.String("budget-id", "default", "YNAB budget ID")
		accountID   = fs.String("account-id", "", "Account to create transactions in")
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

	fmt.Printf("Last Transaction: %+v\n", ledger.Transactions[len(ledger.Transactions)-1])

	_, err = ynab.NewClient(http.DefaultClient, ynab.Config{
		PersonalAccessToken: *accessToken,
		BudgetID:            *budget,
	})
	if err != nil {
		log.Fatalf("Error creating YNAB client: %v", err)
	}

}
