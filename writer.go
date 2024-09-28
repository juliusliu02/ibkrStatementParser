package main

import (
	"encoding/csv"
	"fmt"
	"os"
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
		tx.write()
	}
	csvWriter.Flush()
}
