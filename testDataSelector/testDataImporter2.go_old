package testDataSelector

import (
	"bytes"
	"embed"
	"encoding/csv"
	"io"
	"log"
)

//go:embed testData/FenixRawTestdata_646rows_211220.csv
var embeddedFile embed.FS

// Define your struct based on the CSV file structure
type TestDataRowType []string

func ImportTestDataFromFile2() ([]string, []TestDataRowType) {
	// Read the embedded file
	data, err := embeddedFile.ReadFile("testData/FenixRawTestdata_646rows_211220.csv")
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	// Parse the CSV file
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = ';' // CSV is semicolon-delimited

	var testDataRows []TestDataRowType
	// Read the header row first
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Error reading headers: %v", err)
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

		var testDataRow TestDataRowType
		// Loop all records
		for _, recordItem := range record {
			testDataRow = append(testDataRow, recordItem)
		}

		// Add row to TestDataRows
		testDataRows = append(testDataRows, testDataRow)

	}

	return headers, testDataRows
}
