package main

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
	"time"
)

type Transaction interface {
	getTransactionType() string
}

type Trade struct {
	time     time.Time
	curr     currency.Unit
	symbol   string
	quantity decimal.Decimal
	amount   decimal.Decimal
	fee      decimal.Decimal
}

type Forex struct {
	time       time.Time
	curr       currency.Unit
	targetCurr currency.Unit
	quantity   decimal.Decimal
	amount     decimal.Decimal
	fee        decimal.Decimal
}

type Cash struct {
	date time.Time
	curr currency.Unit
	// Use negative amounts to represent withdrawals.
	amount decimal.Decimal
}

type Dividend struct {
	date   time.Time
	curr   currency.Unit
	symbol string
	amount decimal.Decimal
	tax    decimal.Decimal
}

/* Trade */

func (t *Trade) getTransactionType() string {
	if t.quantity.IsPositive() {
		return "BUY"
	} else {
		return "SELL"
	}
}

func (f *Forex) getTransactionType() string {
	return "CONVERSION"
}

func (c *Cash) getTransactionType() string {
	if c.amount.IsPositive() {
		return "DEPOSIT"
	} else {
		return "WITHDRAWAL"
	}
}

func (d *Dividend) getTransactionType() string {
	return "DIVIDEND"
}

func (d *Dividend) setTax(dec decimal.Decimal) {
	d.tax = dec
}
