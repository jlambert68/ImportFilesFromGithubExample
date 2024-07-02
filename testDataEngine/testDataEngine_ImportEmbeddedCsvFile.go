package testDataEngine

import (
	"bytes"
	"embed"
	"encoding/csv"
	"io"
	"log"
)

// ImportEmbeddedCsvTestDataFile
// Imports an embedded csv-file with relative path and name in 'fileNameAndRelativePath'
// and having a data divider of type 'divider'
// The first row must consist of column headers
func ImportEmbeddedCsvTestDataFile(
	embeddedFilePtr *embed.FS,
	fileNameAndRelativePath string,
	divider rune) (
	testDataHeaders []string,
	testDataRows [][]string) {

	var err error

	// Read the embedded file
	data, err := embeddedFilePtr.ReadFile(fileNameAndRelativePath)
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	// Parse the CSV file
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = divider

	// Read the header row first
	testDataHeaders, err = r.Read()
	if err != nil {
		log.Fatalf("Error reading headers: %v", err)
	}

	// Iterate through the records and extract rows
	for {
		rowRecord, errOrEOF := r.Read()

		// Check if we reach end of file
		if errOrEOF == io.EOF {
			break
		}

		// Check for error
		if errOrEOF != nil {
			log.Fatal(err)
		}

		// Loop all records in  row and extract them
		var testDataRow []string
		for _, recordItem := range rowRecord {
			testDataRow = append(testDataRow, recordItem)
		}

		// Add row to TestDataRows
		testDataRows = append(testDataRows, testDataRow)

	}

	return testDataHeaders, testDataRows
}
