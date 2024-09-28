package main

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
	"time"
)

type Transaction interface {
	getTransactionType() string
	write()
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

func (c *Cash) write() {
	err := csvWriter.Write([]string{
		c.date.Format(time.DateOnly), // Date
		account,                      // Account
		c.getTransactionType(),       // Transaction Type
		"",                           // Instrument Type
		"",                           // Symbol
		"",                           // Quantity
		c.amount.Abs().String(),      // Amount
		c.curr.String(),              // Currency
		"",                           // Fees
		"",                           // Fees Currency
		"",                           // Taxes
		"",                           // Taxes Currency
		"",                           // Converted
		"",                           // Converted Currency
	})
	if err != nil {
		panic(err)
	}
}

func (t *Trade) write() {
	err := csvWriter.Write([]string{
		t.time.Format(time.DateTime), // Date
		account,                      // Account
		t.getTransactionType(),       // Transaction Type
		"SECURITY",                   // Instrument Type
		t.symbol,                     // Symbol
		t.quantity.Abs().String(),    // Quantity
		t.amount.Abs().String(),      // Amount
		t.curr.String(),              // Currency
		t.fee.Abs().String(),         // Fees
		t.curr.String(),              // Fees Currency
		"",                           // Taxes
		"",                           // Taxes Currency
		"",                           // Converted
		"",                           // Converted Currency
	})
	if err != nil {
		panic(err)
	}
}

func (f *Forex) write() {
	err := csvWriter.Write([]string{
		f.time.Format(time.DateTime), // Date
		account,                      // Account
		f.getTransactionType(),       // Transaction Type
		"",                           // Instrument Type
		"",                           // Symbol
		"",                           // Quantity
		f.amount.Abs().String(),      // Amount
		f.curr.String(),              // Currency
		f.fee.Abs().String(),         // Fees
		"USD",                        // Fees Currency
		"",                           // Taxes
		"",                           // Taxes Currency
		f.quantity.String(),          // Converted
		f.targetCurr.String(),        // Converted Currency
	})
	if err != nil {
		panic(err)
	}
}

func (d *Dividend) write() {
	err := csvWriter.Write([]string{
		d.date.Format(time.DateOnly), // Date
		account,                      // Account
		d.getTransactionType(),       // Transaction Type
		"SECURITY",                   // Instrument Type
		d.symbol,                     // Symbol
		"",                           // Quantity
		d.amount.String(),            // Amount
		d.curr.String(),              // Currency
		"",                           // Fees
		"",                           // Fees Currency
		d.tax.Abs().String(),         // Taxes
		d.curr.String(),              // Taxes Currency
		"",                           // Converted
		"",                           // Converted Currency
	})
	if err != nil {
		panic(err)
	}
}
