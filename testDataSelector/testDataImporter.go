package testDataSelector

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

////go:embed testData/FenixRawTestdata_646rows_211220.csv
//var embeddedFile embed.FS

// Define your struct based on the CSV file structure
type TestData struct {
	TestDataId                    string
	AccountId                     string
	DirectClientCustodyAccountId  string
	AccountCurrency               string
	AccountEnvironment            string
	ClientJuristictionCountryCode string
	DebitOrCredit                 string
	PSET                          string
	MarketCountry                 string
	MarketName                    string
	MarketSubType                 string
	AccountType                   string
	PooledCashAccount             string
	Nostro                        string
	MarketCurrency                string
	ISIN                          string
	SecurityType                  string
	InterimCurrency               string
	ContraCurrency                string
	PrincipalOrIncome             string
	SecProgram                    string
	ProvisionalIncome             string
	Random                        string
}

func ImportTestDataFromFile() []TestData {
	// Read the embedded file
	data, err := embeddedFile.ReadFile("FenixRawTestdata_646rows_211220.csv")
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	// Parse the CSV file
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = ';' // CSV is semicolon-delimited

	var testData []TestData
	// Skip header row if necessary
	if _, err := r.Read(); err != nil {
		log.Fatal(err)
	}

	// Iterate through the records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Assuming the order of your CSV columns matches the struct
		td := TestData{
			TestDataId:                    record[0],
			AccountId:                     record[1],
			DirectClientCustodyAccountId:  record[2],
			AccountCurrency:               record[3],
			AccountEnvironment:            record[4],
			ClientJuristictionCountryCode: record[5],
			DebitOrCredit:                 record[6],
			PSET:                          record[7],
			MarketCountry:                 record[8],
			MarketName:                    record[9],
			MarketSubType:                 record[10],
			AccountType:                   record[11],
			PooledCashAccount:             record[12],
			Nostro:                        record[13],
			MarketCurrency:                record[14],
			ISIN:                          record[15],
			SecurityType:                  record[16],
			InterimCurrency:               record[17],
			ContraCurrency:                record[18],
			PrincipalOrIncome:             record[19],
			SecProgram:                    record[20],
			ProvisionalIncome:             record[21],
			Random:                        record[22],
		}

		testData = append(testData, td)
	}

	// Now `testData` slice has all the data from the CSV
	for _, td := range testData {
		fmt.Printf("Test Data ID: %s, Account ID: %s, ... (other fields)\n", td.TestDataId, td.AccountId)
		// Print other fields as needed
	}

	return testData
}
