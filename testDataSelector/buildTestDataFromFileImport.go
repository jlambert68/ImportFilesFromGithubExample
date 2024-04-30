package testDataSelector

import (
	uuidGenerator "github.com/google/uuid"
	"log"
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

func buildTestDataMap(testData []TestDataRowType) *map[TestDataDomainUuidType]*TestDataDomainModeStruct {

	var testDataHeaders TestDataRowType
	testDataHeaders = testData[0]

	// Initialize your map
	var testDataModelMap map[TestDataDomainUuidType]*TestDataDomainModeStruct
	testDataModelMap = make(map[TestDataDomainUuidType]*TestDataDomainModeStruct)

	// Keeping track of all TestDataPoints for each Column
	var testDataPointsForColumns [][]*TestDataPointValueStruct

	// Initiate the maps used
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnMap map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnAndRowUuidMap map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	var tempTestDataColumnsMetaDataMap map[TestDataColumnUuidType]*[]*TestDataColumnMetaDataStruct
	tempTestDataValuesForRowMap = make(map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnAndRowUuidMap = make(map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct)
	tempTestDataColumnsMetaDataMap = make(map[TestDataColumnUuidType]*[]*TestDataColumnMetaDataStruct)

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

	// Iterate through the CSV records to extract the TestDataPoints
	for testDataRowIndex, tempTestDataRow := range testData {

		// Don't process header row
		if testDataRowIndex == 0 {
			continue
		}
		rowUuid, err := uuidGenerator.Parse(tempTestDataRow[0])
		if err != nil {
			log.Fatalln(err)
		}

		var testDataPointsForRow []*TestDataPointValueStruct

		// Loop over all TestDataPoints in the row
		for testDataColumnIndex, tempTestDataPoint := range tempTestDataRow {

			columnUuid, err := uuidGenerator.Parse(testDataHeaders[testDataColumnIndex])
			if err != nil {
				log.Fatalln(err)
			}

			columnAndRowUuid, err := uuidGenerator.Parse(columnUuid.String() + rowUuid.String())
			if err != nil {
				log.Fatalln(err)
			}

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
			}

			// Add TestDataPoint to 'testDataPointsForRow'
			testDataPointsForRow = append(testDataPointsForRow, testDataPoint)

			// Add the TestDataPoint to the Map for all TestDataPoints
			tempTestDataValuesForColumnAndRowUuidMap[TestDataColumnAndRowUuidType(columnAndRowUuid.String())] = testDataPoint

			// Add the TestDataPoint to correct slice for columns
			testDataPointsForColumns[testDataColumnIndex] = append(testDataPointsForColumns[testDataColumnIndex], testDataPoint)

		}

		// Add 'testDataPointsForRow' to Map for TestDataPoints in one row
		tempTestDataValuesForRowMap[TestDataPointRowUuidType(rowUuid.String())] = &testDataPointsForRow

	}

	// Loop 'testDataPointsForColumns' and add to column Map
	var tempTestDataColumnUuid TestDataColumnUuidType
	for _, testDataPointsForColumn := range testDataPointsForColumns {

		tempTestDataColumnUuid = testDataPointsForColumn[0].TestDataColumnUuid
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn
	}

	return &testDataModelMap
}
