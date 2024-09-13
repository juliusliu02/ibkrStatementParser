package main

import (
	"encoding/csv"
	"os"
)

func main() {
	// Open the CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	p := parser{reader: reader}
	transactions := p.read()

	writeTransaction(transactions)
}
