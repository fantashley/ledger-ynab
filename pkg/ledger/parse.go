package ledger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fantashley/ledger-ynab/pkg/usd"
)

func loadTransactions(path string) ([]Transaction, error) {
	filePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for file %s: %w", path, err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ledger file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("Failed to close ledger file: %v", err)
		}
	}()

	var (
		scanner      = bufio.NewScanner(file)
		transactions = []Transaction{}
		currentTx    Transaction
		currentLine  string
		lineNum      int
	)

	for scanner.Scan() {
		lineNum++
		currentLine = scanner.Text()

		trimmed := strings.TrimSpace(currentLine)
		if strings.HasPrefix(trimmed, ";") {
			if lineNum > 1 {
				currentTx.Comment += strings.TrimPrefix(trimmed, ";")
			}
			continue
		}

		chunks := strings.Fields(currentLine)
		if len(chunks) < 2 {
			log.Printf("Could not classify line %d: %q", lineNum, currentLine)
			continue
		}

		if date, err := time.Parse("2006/01/02", chunks[0]); err == nil {
			if err = currentTx.Validate(); err == nil {
				transactions = append(transactions, currentTx)
			} else if !currentTx.Empty() {
				log.Printf("Invalid transaction before line %d: %+v: %v", lineNum, currentTx, err)
			}

			currentTx = Transaction{
				Date:  date,
				Payee: strings.Join(chunks[1:], " "),
			}

			continue
		}

		if strings.ContainsRune(chunks[0], ':') && strings.ContainsRune(chunks[1], '$') {
			amt, err := usd.ParseUSD(chunks[1])
			if err != nil {
				log.Printf("Failed to detect dollar amount of line %d: %q: %v", lineNum, currentLine, err)
				continue
			}

			account := Entry{
				Name:   chunks[0],
				Amount: amt,
			}

			if strings.HasPrefix(chunks[0], "Expense") {
				currentTx.Categories = append(currentTx.Categories, account)
				continue
			}
			if strings.HasPrefix(chunks[0], "Liabilities") {
				currentTx.Payers = append(currentTx.Payers, account)
				continue
			}
		}

		log.Printf("Could not classify line %d: %q", lineNum, currentLine)
	}

	if err = currentTx.Validate(); err == nil {
		transactions = append(transactions, currentTx)
	}

	return transactions, nil
}
