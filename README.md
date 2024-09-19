# `ibkrStatementParser` Readme

`ibkrStatementParser` is a tool designed to parse Interactive Broker's statements in CSV format and extracts relevant trade records. With this tool, users can easily convert their statements into CSV files that can be directly imported into [Portfolio Profit](https://portfolioprofit.app/).

## Usage

Export your data from Interactive Broker in CSV format, place it under the project directory named as `data.csv`, and execute the program. The output will be stored in the same directory.

## Project structure

The functionality of `ibkrStatementParser` is divided into three separate files: `parser`, `writer`, and `transaction`.

### `parser`

The parser is responsible for reading from CSV file, filters out the data rows, and parses them into transaction structs.

#### Variables

- `reader`: reads from CSV.
- `header`: store headers.
- `map`: maps header to data.
- `trades`, `forexes`, `cashes`, `dividends`: arrays that stores pointers to `transaction` structs.

The parser recognizes data rows by matching the starting fields of rows with pre-defined templates. If a match is found, the row is considered valid.

### `writer`

The writer takes an array of transactions and writes them into a new CSV file that adheres to [Portfolio Profit’s CSV import documentation](https://portfolioprofit.app/docs/import/csv). There is a dedicated write method for each transaction struct.

### `transaction`

All custom structs for transactions, including Trade, Forex, Cash, and Dividend, implement the Transaction interface.

The structs use `time.Time` to represent date/time, [currency.Unit](https://pkg.go.dev/golang.org/x/text/currency) for currency symbols, and [decimal.Decimal](https://pkg.go.dev/github.com/shopspring/decimal) to represent numbers in decimal.

**The program assumes all currency-related fields (such as `amount`, `fee`) in one transaction have the same currency symbol except in `Forex`.**

For a detailed description of the structs, see Structs.

## Data structure

### Type alias

`recordType` is an alias of `int` used to indicate types of record. The following constants are defined as `recordType`s:

- `none`: no match.
- `meta`: global header.
- `trades`: stock transactions.
- `forex`: foreign exchange transactions.
- `cash`: deposit / withdrawal.
- `dividend`: dividend transactions.
- `tax`: tax transactions (attributed to a transaction).
- `feeAdjust`: commission adjustments.
- `count`: a functional value.

### Structs

#### Trade

`Trade` represents stock-trading transactions:

- time: time.Time
- curr: currency.Unit
- symbol: string
- quantity, amount, fee: decimal.Decimal

The sign of `quantity` decides if it is a `BUY` or a `SELL`.

#### Forex

`Forex` represents transactions for currency exchange:

- time: time.Time
- curr, targetCurr: currency.Unit
- quantity, amount, fee: decimal.Decimal

#### Cash

`Cash` represents deposits and withdrawals of cash of the account:

- date: time.Time
- curr: currency.Unit
- amount: decimal.Decimal

The sign of `amount` decides if it is a `DEPOSIT` or a `WITHDRAWAL`.

#### Dividend

`Dividend` represents dividend transactions:

- time: time.Time
- curr: currency.Unit
- symbol: string
- amount, tax: decimal.Decimal

## Algorithm

The program reads data from `data.csv` in the project directory, processes valid data rows, and generates an `output.csv` file.

Transaction records are read from different sections of the Interactive Broker statement and then used to create structs. When adjustment records are encountered, such as fee adjustments or taxes, they are incorporated into the existing transactions.

Finally, write methods are used to export transactions to the output file with the schema defined in `writer.go`.
