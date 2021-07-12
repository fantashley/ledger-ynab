package usd

import (
	"fmt"
	"strconv"
	"strings"
)

type USD struct {
	Negative bool
	Dollars  int
	Cents    int
}

func (d USD) Add(amount USD) USD {
	return ToUSD(d.CentsTotal() + amount.CentsTotal())
}

func (d USD) CentsTotal() int {
	total := 100*d.Dollars + d.Cents
	if d.Negative {
		total *= -1
	}

	return total
}

func ToUSD(cents int) USD {
	var (
		usd      USD
		absCents int
	)

	absCents = cents
	if cents < 0 {
		absCents = -1 * cents
		usd.Negative = true
	}

	usd.Dollars = absCents / 100
	usd.Cents = absCents - 100*usd.Dollars

	return usd
}

func ParseUSD(amount string) (USD, error) {
	var (
		usd USD
		err error
	)

	withoutSign := strings.TrimPrefix(amount, "$")
	usd.Negative = withoutSign[0] == '-'
	if usd.Negative {
		withoutSign = withoutSign[1:]
	}

	split := strings.Split(withoutSign, ".")
	if len(split) != 2 {
		return usd, fmt.Errorf("incorrect string format for parsing dollar amount: %q", amount)
	}

	usd.Dollars, err = strconv.Atoi(split[0])
	if err == nil {
		usd.Cents, err = strconv.Atoi(split[1])
	}
	if err != nil {
		return usd, fmt.Errorf("failed to convert %s to integers: %w", withoutSign, err)
	}

	return usd, nil
}
