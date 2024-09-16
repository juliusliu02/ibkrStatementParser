package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
	"regexp"
	"strings"
	"time"
)

/*
	Parser reads from CSV data and parses it into corresponding

Transaction data types.
*/
type parser struct {
	reader    *csv.Reader
	header    []string
	fields    map[string]string
	trades    []*Trade
	forexes   []*Forex
	cashes    []*Cash
	dividends []*Dividend
}

func (p *parser) read() []Transaction {
	data, err := p.reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var reportDate string

	for _, row := range data {

		if row[1] == "Header" {
			p.fields = make(map[string]string)
			p.header = row
			continue
		}
		for i, v := range row {
			// populate elements in map
			p.fields[p.header[i]] = v
		}
		switch tType := findMatch(row); tType {
		case meta:
			// Parse metadata
			if row[2] == "Period" {
				_, reportDate, _ = strings.Cut(row[3], " - ")
			}
		case trades:
			{
				trade, err := makeTrade(p.fields)
				if err == nil {
					p.trades = append(p.trades, &trade)
				}
			}
		case forex:
			{
				forex, err := makeForex(p.fields)
				if err == nil {
					forex.time, err = time.Parse("January 2, 2006", reportDate)
					if err != nil {
						fmt.Println("failed to parse time: " + reportDate)
					}
					p.forexes = append(p.forexes, &forex)
				}
			}
		case cash:
			{
				cash, err := makeCash(p.fields)
				if err == nil {
					p.cashes = append(p.cashes, &cash)
				}
			}
		case dividend:
			{
				dividend, err := makeDividend(p.fields)
				if err == nil {
					p.dividends = append(p.dividends, &dividend)
				}
			}
		case tax:
			{
				if dividend, err := findDividend(p.dividends, p.fields["Date"],
					strings.Split(p.fields["Description"], "(")[0]); err == nil {
					t, err := decimal.NewFromString(p.fields["Amount"])
					if err == nil {
						(*dividend).setTax(t)
					}
				}
			}
		case feeAdjust:
			regex, _ := regexp.Compile(`\((.*)\)`)
			// FindStringSubmatch returns the captured groups in an array.
			symbol := regex.FindStringSubmatch(p.fields["Description"])[1]
			if symbol == "" {
				fmt.Println("An unrecognized commission adjustment is reported. Please modify the value accordingly:" + strings.Join(row, " "))
			} else {
				t, err := findTrades(p.trades, p.fields["Date"], symbol)
				if err == nil {
					feeAdjust, err := decimal.NewFromString(p.fields["Amount"])
					if err == nil {
						t.fee = t.fee.Add(feeAdjust)
					}
				}
			}
		case count, none:
			continue
		}
	}

	transactions := make([]Transaction, 0)
	// Go's type system prevents type casting from []type to []interface.
	for _, trade := range p.trades {
		transactions = append(transactions, trade)
	}
	for _, forex := range p.forexes {
		transactions = append(transactions, forex)
	}
	for _, cash := range p.cashes {
		transactions = append(transactions, cash)
	}
	for _, dividend := range p.dividends {
		transactions = append(transactions, dividend)
	}
	return transactions
}

/* Constants and type alias. */

type recordType = int

const (
	none recordType = iota // sentinel value
	meta
	trades
	forex
	cash
	dividend
	tax
	feeAdjust
	count
)

func getTemplate(tType recordType) []string {
	switch tType {
	case meta:
		return []string{"Statement", "Data"}
	case trades:
		return []string{"Trades", "Data", "Order", "Stocks"}
	case forex:
		return []string{"Trades", "SubTotal", "", "Forex"}
	case cash:
		return []string{"Deposits & Withdrawals", "Data"}
	case dividend:
		return []string{"Dividends", "Data"}
	case tax:
		return []string{"Withholding Tax", "Data", "USD"}
	case feeAdjust:
		return []string{"Commission Adjustments", "Data", "USD"}
	// return empty string for invalid recordType.
	case count, none:
		return []string{}
	}
	// shouldn't execute.
	return []string{}
}

/* Functions to match records to corresponding transactionTypes. */

func findMatch(row []string) recordType {
	for i := range count {
		if matchStart(row, i) {
			return i
		}
	}
	return none
}

/* Match the start of the string to a template. Return false if recordType is invalid. */
func matchStart(row []string, tType recordType) bool {
	tem := getTemplate(tType)
	// len(tem) == 0 when tType is invalid.
	if len(tem) == 0 || len(row) < len(tem) {
		return false
	}
	for i, v := range tem {
		if row[i] != v {
			return false
		}
	}
	return true
}

/* Struct initializers that read from fields and populate data into structs.
Change keys if there are any changes to the CSV source file. */

func makeTrade(fields map[string]string) (Trade, error) {
	const timeFormat = "2006-01-02, 15:04:05"

	t, err0 := time.Parse(timeFormat, fields["Date/Time"])
	c, err1 := currency.ParseISO(fields["Currency"])
	q, err2 := decimal.NewFromString(fields["Quantity"])
	a, err3 := decimal.NewFromString(fields["Proceeds"])
	f, err4 := decimal.NewFromString(fields["Comm/Fee"])
	errs := []error{err0, err1, err2, err3, err4}
	for _, err := range errs {
		if err != nil {
			return Trade{}, errors.New("Invalid data:" + err.Error())
		}
	}
	return Trade{
		time:     t,
		curr:     c,
		symbol:   fields["Symbol"],
		quantity: q,
		amount:   a,
		fee:      f,
	}, nil
}

func makeForex(fields map[string]string) (Forex, error) {
	var tCurr string
	var slices = strings.Split(fields["Symbol"], ".")
	if slices[0] == fields["Currency"] {
		tCurr = slices[1]
	} else {
		tCurr = slices[0]
	}

	c, err0 := currency.ParseISO(fields["Currency"])
	tc, err1 := currency.ParseISO(tCurr)
	q, err2 := decimal.NewFromString(fields["Quantity"])
	a, err3 := decimal.NewFromString(fields["Proceeds"])
	f, err4 := decimal.NewFromString(fields["Comm in USD"])
	errs := []error{err0, err1, err2, err3, err4}
	for _, err := range errs {
		if err != nil {
			return Forex{}, errors.New("Invalid data:" + err.Error())
		}
	}

	return Forex{
		curr:       c,
		targetCurr: tc,
		quantity:   q,
		amount:     a,
		fee:        f,
	}, nil
}

func makeCash(fields map[string]string) (Cash, error) {
	const timeFormat = "2006-01-02"

	t, err0 := time.Parse(timeFormat, fields["Settle Date"])
	c, err1 := currency.ParseISO(fields["Currency"])
	a, err2 := decimal.NewFromString(fields["Amount"])
	errs := []error{err0, err1, err2}
	for _, err := range errs {
		if err != nil {
			return Cash{}, errors.New("Invalid data:" + err.Error())
		}
	}
	return Cash{date: t, curr: c, amount: a}, nil
}

func makeDividend(fields map[string]string) (Dividend, error) {
	const timeFormat = "2006-01-02"

	d, err0 := time.Parse(timeFormat, fields["Date"])
	c, err1 := currency.ParseISO(fields["Currency"])
	a, err2 := decimal.NewFromString(fields["Amount"])
	errs := []error{err0, err1, err2}
	for _, err := range errs {
		if err != nil {
			return Dividend{}, errors.New("Invalid data:" + err.Error())
		}
	}

	return Dividend{
		date:   d,
		curr:   c,
		symbol: strings.Split(fields["Description"], "(")[0],
		amount: a,
	}, nil
}

/* Returns the first transaction that matches with the conditions. */
func findDividend(ds []*Dividend, t string, symbol string) (*Dividend, error) {
	date, err := time.Parse(time.DateOnly, t)
	if err != nil {
		return &Dividend{}, errors.New("invalid date:" + err.Error())
	}
	for _, v := range ds {
		if v.date.Truncate(24*time.Hour) == date.Truncate(24*time.Hour) &&
			v.symbol == symbol {
			return v, nil
		}
	}
	return &Dividend{}, errors.New("dividend not found")
}

func findTrades(ds []*Trade, t string, symbol string) (*Trade, error) {
	date, err := time.Parse(time.DateOnly, t)
	if err != nil {
		return &Trade{}, errors.New("invalid date:" + err.Error())
	}
	for _, v := range ds {
		if v.time.Truncate(24*time.Hour) == date.Truncate(24*time.Hour) &&
			v.symbol == symbol {
			return v, nil
		}
	}
	return &Trade{}, errors.New("dividend not found")
}
