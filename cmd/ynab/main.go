package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fantashley/ledger-ynab/pkg/ledger"
)

func main() {
	var (
		ledgerFile = flag.String("f", "ledger.dat", "Path to ledger file")
	)

	flag.Parse()

	ledger, err := ledger.New(*ledgerFile)
	if err != nil {
		log.Panicf("Failed to create ledger: %v", err)
	}

	fmt.Printf("Last Transaction: %+v\n", ledger.Transactions[len(ledger.Transactions)-1])
}
