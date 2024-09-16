package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var (
	account   = "IBKR"
	csvWriter *csv.Writer
	header    = []string{
		"Date",
		"Account",
		"Type",
		"Instrument Type",
		"Ticker Symbol",
		"Quantity",
		"Amount",
		"Currency",
		"Fees",
		"Fees Currency",
		"Taxes",
		"Taxes Currency",
		"Converted",
		"Converted Currency",
	}
)

func writeTransaction(transactions []Transaction) {
	f, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	csvWriter = csv.NewWriter(f)
	err = csvWriter.Write(header)
	if err != nil {
		return
	}
	for _, tx := range transactions {
		switch v := tx.(type) {
		case *Cash:
			writeCash(v)
		case *Trade:
			writeTrade(v)
		case *Forex:
			writeForex(v)
		case *Dividend:
			writeDividend(v)
		}
	}
	csvWriter.Flush()
}

func writeCash(c *Cash) {
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

func writeTrade(t *Trade) {
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

func writeForex(f *Forex) {
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

func writeDividend(d *Dividend) {
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
