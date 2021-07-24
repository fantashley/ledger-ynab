package ynab

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	URL                 string
	PersonalAccessToken string
	BudgetID            string
}

func (c Config) FillDefaults() Config {
	if c.URL == "" {
		c.URL = "https://api.youneedabudget.com/"
	}

	return c
}

func (c Config) Validate() error {
	var err error
	if c.URL == "" {
		err = multierror.Append(err, errors.New("api url missing"))
	}
	if c.PersonalAccessToken == "" {
		err = multierror.Append(err, errors.New("personal access token missing"))
	}
	if c.BudgetID == "" {
		err = multierror.Append(err, errors.New("budget id missing"))
	}

	return err
}
