package testDataSelector

import (
	uuidGenerator "github.com/google/uuid"
)

const (
	testDataDomainUuid TestDataDomainUuidType = "7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae"
	testDataDomainName TestDataDomainNameType = "Sub Custody"
	testDataAreaUuid   TestDataAreaUuidType   = "010cc994-a913-4041-96fe-a96d7e0c97e8"
	testDataAreaName   TestDataAreaNameType   = "Main TestData Area"
)

func buildTestDataMap(headers []string, testData []TestDataRowType) *map[TestDataDomainUuidType]*TestDataDomainModelStruct {

	// Define a namespace UUID; this could be any valid UUID that you choose to use as a namespace for your IDs.
	// Here, we use the DNS namespace provided by the UUID package for demonstration purposes.
	namespace := uuidGenerator.NameSpaceDNS

	var testDataHeaders TestDataRowType
	testDataHeaders = headers

	// Initialize your map
	var testDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	testDataModelMap = make(map[TestDataDomainUuidType]*TestDataDomainModelStruct)

	// Keeping track of all TestDataPoints for each Column
	var testDataPointsForColumns [][]*TestDataPointValueStruct

	// Initiate the maps used
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnMap map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnAndRowUuidMap map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	var tempTestDataColumnsMetaDataMap map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct
	var tempUniqueTestDataValuesForColumnMap map[TestDataColumnUuidType]*map[TestDataValueType]TestDataValueType
	tempTestDataValuesForRowMap = make(map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnAndRowUuidMap = make(map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct)
	tempTestDataColumnsMetaDataMap = make(map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct)
	tempUniqueTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*map[TestDataValueType]TestDataValueType)

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
		UniqueTestDataValuesForColumnMap:     &tempUniqueTestDataValuesForColumnMap,
	}

	// Create the TestDataAreaMap and TestData for Area
	var tempTestDataAreasMap map[TestDataAreaUuidType]*TestDataAreaStruct
	tempTestDataAreasMap = make(map[TestDataAreaUuidType]*TestDataAreaStruct)

	var tempTestDataDomainModel TestDataDomainModelStruct
	tempTestDataDomainModel = TestDataDomainModelStruct{
		TestDataDomainUuid: testDataDomainUuid,
		TestDataDomainName: testDataDomainName,
		TestDataAreasMap:   &tempTestDataAreasMap,
	}

	tempTestDataAreasMap[testDataAreaUuid] = testDataArea

	// Add the TestArea to under the TestDataDomain in the full TestDataModelMap
	testDataModelMap[testDataDomainUuid] = &tempTestDataDomainModel

	// Columns that require true for specific properties
	trueColumns := map[string]bool{
		"AccountCurrency":               true,
		"AccountEnvironment":            true,
		"ClientJuristictionCountryCode": true,
		"DebitOrCredit":                 true,
		"MarketCountry":                 true,
		"MarketName":                    true,
		"MarketSubType":                 true,
		"MarketCurrency":                true,
		"InterimCurrency":               true,
		"ContraCurrency":                true,
	}

	// Iterate through the CSV records to extract the TestDataPoints
	for _, tempTestDataRow := range testData {

		rowUuid := uuidGenerator.NewSHA1(namespace, []byte(tempTestDataRow[0]))

		var testDataPointsForRow []*TestDataPointValueStruct

		// The Description of how the name is constructed for one TestDataPoint (row)
		var testDataPointNameDescription string

		// The name for one TestDataPoint (row)
		var testDataPointName string

		// Loop over all TestDataPoints in the row
		for testDataColumnIndex, tempTestDataPoint := range tempTestDataRow {

			columnUuid := uuidGenerator.NewSHA1(namespace, []byte(testDataHeaders[testDataColumnIndex]))

			columnAndRowUuid := uuidGenerator.NewSHA1(namespace, []byte(columnUuid.String()+rowUuid.String()))

			// Create the TestDataPoint
			var testDataPoint *TestDataPointValueStruct
			testDataPoint = &TestDataPointValueStruct{
				TestDataDomainUuid:           testDataDomainUuid,
				TestDataDomainName:           testDataDomainName,
				TestDataAreaUuid:             testDataAreaUuid,
				TestDataAreaName:             testDataAreaName,
				TestDataColumnUuid:           TestDataColumnUuidType(columnUuid.String()),
				TestDataColumnDataName:       TestDataColumnDataNameType(testDataHeaders[testDataColumnIndex]),
				TestDataColumnUIName:         TestDataColumnUINameType(testDataHeaders[testDataColumnIndex]),
				TestDataPointRowUuid:         TestDataPointRowUuidType(rowUuid.String()),
				TestDataColumnAndRowUuid:     TestDataColumnAndRowUuidType(columnAndRowUuid.String()),
				TestDataValue:                TestDataValueType(tempTestDataPoint),
				TestDataValueNameDescription: "",
				TestDataValueName:            "",
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
				if len(testDataPointNameDescription) == 0 {
					testDataPointNameDescription = string(testDataPoint.TestDataColumnDataName)
					testDataPointName = string(testDataPoint.TestDataValue)
				} else {
					testDataPointNameDescription = testDataPointNameDescription + "/" + string(testDataPoint.TestDataColumnDataName)
					testDataPointName = testDataPointName + "/" + string(testDataPoint.TestDataValue)
				}
			}

		}

		// Loop the Values in the row and add 'TestDataPointName'
		for _, testDataPoint := range testDataPointsForRow {
			testDataPoint.TestDataValueNameDescription = TestDataValueNameDescriptionType(testDataPointNameDescription)
			testDataPoint.TestDataValueName = TestDataValueNameType(testDataPointName)
		}

		// Add 'testDataPointsForRow' to Map for TestDataPoints in one row
		tempTestDataValuesForRowMap[TestDataPointRowUuidType(rowUuid.String())] = &testDataPointsForRow

	}

	// Loop 'testDataPointsForColumns' and add to column Map and create ColumnMetaData
	var tempTestDataColumnUuid TestDataColumnUuidType
	for _, testDataPointsForColumn := range testDataPointsForColumns {

		tempTestDataColumnUuid = testDataPointsForColumn[0].TestDataColumnUuid
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn

		// initialized the 'inner' map for unique values in the column
		var innerMapWithUniqueColumnValues map[TestDataValueType]TestDataValueType
		innerMapWithUniqueColumnValues = make(map[TestDataValueType]TestDataValueType)

		// Loop all Column data and add to the 'inner' Map, witch will create a map with only one value per occurrence of the value
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn
		for _, testDataPoint := range testDataPointsForColumn {
			innerMapWithUniqueColumnValues[testDataPoint.TestDataValue] = testDataPoint.TestDataValue
		}

		// Add the 'inner' map to the Column-Map
		tempUniqueTestDataValuesForColumnMap[tempTestDataColumnUuid] = &innerMapWithUniqueColumnValues

		// Create column MetaData
		var tempTestDataColumnMetaDataStruct *TestDataColumnMetaDataStruct
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

	return &testDataModelMap
}
