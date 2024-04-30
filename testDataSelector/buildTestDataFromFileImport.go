package testDataSelector

import (
	uuidGenerator "github.com/google/uuid"
)

// ... (All your type definitions here)

// Placeholder function to generate domain UUID from domain name
func getDomainUUID(domainName string) TestDataDomainUuidType {
	// In a real application, you would have a proper mechanism to generate or retrieve UUIDs
	return TestDataDomainUuidType("domain-uuid-" + domainName)
}

const (
	testDataDomainUuid TestDataDomainUuidType = "7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae"
	testDataDomainName TestDataDomainNameType = "Sub Custody"
	testDataAreaUuid   TestDataAreaUuidType   = "010cc994-a913-4041-96fe-a96d7e0c97e8"
	testDataAreaName   TestDataAreaNameType   = "Main TestData Area"
)

func buildTestDataMap(headers []string, testData []TestDataRowType) *map[TestDataDomainUuidType]*TestDataDomainModeStruct {

	// Define a namespace UUID; this could be any valid UUID that you choose to use as a namespace for your IDs.
	// Here, we use the DNS namespace provided by the UUID package for demonstration purposes.
	namespace := uuidGenerator.NameSpaceDNS

	var testDataHeaders TestDataRowType
	testDataHeaders = headers

	// Initialize your map
	var testDataModelMap map[TestDataDomainUuidType]*TestDataDomainModeStruct
	testDataModelMap = make(map[TestDataDomainUuidType]*TestDataDomainModeStruct)

	// Keeping track of all TestDataPoints for each Column
	var testDataPointsForColumns [][]*TestDataPointValueStruct

	// Initiate the maps used
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnMap map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnAndRowUuidMap map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	var tempTestDataColumnsMetaDataMap map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct
	tempTestDataValuesForRowMap = make(map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnAndRowUuidMap = make(map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct)
	tempTestDataColumnsMetaDataMap = make(map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct)

	// TestData for one Area within one Domain
	var testDataArea *TestDataAreaStruct
	testDataArea = &TestDataAreaStruct{
		TestDataDomainUuid:                   testDataDomainUuid,
		TestDataDomainName:                   testDataDomainName,
		TestDataAreaUuid:                     testDataAreaUuid,
		TestDataAreaName:                     testDataAreaName,
		TestDataValuesForRowMap:              &tempTestDataValuesForRowMap,
		TestDataValuesForColumnMap:           &tempTestDataValuesForColumnMap,
		TestDataValuesForColumnAndRowUuidMap: &tempTestDataValuesForColumnAndRowUuidMap,
		TestDataColumnsMetaDataMap:           &tempTestDataColumnsMetaDataMap,
	}

	// Create the TestDataAreaMap and TestData for Area
	var tempTestDataAreasMap map[TestDataAreaUuidType]*TestDataAreaStruct
	tempTestDataAreasMap = make(map[TestDataAreaUuidType]*TestDataAreaStruct)

	var tempTestDataDomainModel TestDataDomainModeStruct
	tempTestDataDomainModel = TestDataDomainModeStruct{
		TestDataDomainUuid: testDataDomainUuid,
		TestDataDomainName: testDataDomainName,
		TestDataAreasMap:   &tempTestDataAreasMap,
	}

	tempTestDataAreasMap[testDataAreaUuid] = testDataArea

	// Add the TestArea to under the TestDataDomain in the full TestDataModelMap
	testDataModelMap[testDataDomainUuid] = &tempTestDataDomainModel


	// Columns that require true for specific properties
	trueColumns := map[string]bool{
		"AccountCurrency":              true,
		"AccountEnvironment":           true,
		"ClientJuristictionCountryCode": true,
		"DebitOrCredit":                true,
		"MarketCountry":                true,
		"MarketName":                   true,
		"MarketSubType":                true,
		"MarketCurrency":               true,
		"InterimCurrency":              true,
		"ContraCurrency":               true,
	}

	// SLice that holds the UI-name for alla TestDataRows
	var testDataPointNames []string

	// Iterate through the CSV records to extract the TestDataPoints
	for _, tempTestDataRow := range testData {

		rowUuid := uuidGenerator.NewSHA1(namespace, []byte(tempTestDataRow[0]))

		var testDataPointsForRow []*TestDataPointValueStruct

		// The name for one TestDataPoint (row)
		var testDataPointName string

		// Loop over all TestDataPoints in the row
		for testDataColumnIndex, tempTestDataPoint := range tempTestDataRow {

			columnUuid := uuidGenerator.NewSHA1(namespace, []byte(testDataHeaders[testDataColumnIndex]))

			columnAndRowUuid := uuidGenerator.NewSHA1(namespace, []byte(columnUuid.String()+rowUuid.String()))

			// Create the TestDataPoint
			var testDataPoint *TestDataPointValueStruct
			testDataPoint = &TestDataPointValueStruct{
				TestDataDomainUuid:       testDataDomainUuid,
				TestDataDomainName:       testDataDomainName,
				TestDataAreaUuid:         testDataAreaUuid,
				TestDataAreaName:         testDataAreaName,
				TestDataColumnUuid:       TestDataColumnUuidType(columnUuid.String()),
				TestDataColumnDataName:   TestDataColumnDataNameType(testDataHeaders[testDataColumnIndex]),
				TestDataColumnUIName:     TestDataColumnUINameType(testDataHeaders[testDataColumnIndex]),
				TestDataPointRowUuid:     TestDataPointRowUuidType(rowUuid.String()),
				TestDataColumnAndRowUuid: TestDataColumnAndRowUuidType(columnAndRowUuid.String()),
				TestDataValue:            TestDataValueType(tempTestDataPoint),
				TestDataValueName: "",
			}

			// Add TestDataPoint to 'testDataPointsForRow'
			testDataPointsForRow = append(testDataPointsForRow, testDataPoint)

			// Add the TestDataPoint to the Map for all TestDataPoints
			tempTestDataValuesForColumnAndRowUuidMap[TestDataColumnAndRowUuidType(columnAndRowUuid.String())] = testDataPoint

			// Add the TestDataPoint to correct slice for columns
			if testDataPointsForColumns == nil || len(testDataPointsForColumns) == testDataColumnIndex {
				testDataPointsForColumns = append(testDataPointsForColumns, []*TestDataPointValueStruct{})
			}
			testDataPointsForColumns[testDataColumnIndex] = append(testDataPointsForColumns[testDataColumnIndex], testDataPoint)

			// If this column is in the TestDataPointName the add it
			if trueColumns[string(testDataPoint.TestDataColumnDataName)] == true {
				if len(testDataPointName) == 0 {
					testDataPointName = string(testDataPoint.TestDataColumnDataName)
				} else {
					testDataPointName = testDataPointName + "/" + string(testDataPoint.TestDataColumnDataName)
				}
			}

		}

		// Add 'TestDataPointName' to slice of names
		testDataPointNames = append(testDataPointNames, testDataPointName)

		// Add 'testDataPointsForRow' to Map for TestDataPoints in one row
		tempTestDataValuesForRowMap[TestDataPointRowUuidType(rowUuid.String())] = &testDataPointsForRow

	}

	// Loop 'testDataPointsForColumns' and add to column Map and create ColumnMetaData
	var tempTestDataColumnUuid TestDataColumnUuidType
	for _, testDataPointsForColumn := range testDataPointsForColumns {

		tempTestDataColumnUuid = testDataPointsForColumn[0].TestDataColumnUuid
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn

		// Create column MetaData
		var tempTestDataColumnMetaDataStruct  *TestDataColumnMetaDataStruct
		tempTestDataColumnMetaDataStruct = &TestDataColumnMetaDataStruct{
			TestDataColumnUuid:                      testDataPointsForColumn[0].TestDataColumnUuid,
			TestDataColumnDataName:                  testDataPointsForColumn[0].TestDataColumnDataName,
			TestDataColumnUIName:                    testDataPointsForColumn[0].TestDataColumnUIName,
			ShouldColumnBeUsedForFindingTestData:    trueColumns[string(testDataPointsForColumn[0].TestDataColumnDataName)],
			ShouldColumnBeUsedWithinTestDataSetName: trueColumns[string(testDataPointsForColumn[0].TestDataColumnDataName)],
		}

		// Add Column MetaData to Map
		tempTestDataColumnsMetaDataMap[tempTestDataColumnUuid] = tempTestDataColumnMetaDataStruct
	}

	// Loop all TestDataPoints and add 'TestDataPointName'
	for


	return &testDataModelMap
}
