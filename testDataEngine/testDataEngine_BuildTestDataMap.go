package testDataEngine

import uuidGenerator "github.com/google/uuid"

func AddTestDataToTestDataModel(
	testDataDomainUuid TestDataDomainUuidType,
	testDataDomainName TestDataDomainNameType,
	testDataAreaUuid TestDataAreaUuidType,
	testDataAreaName TestDataAreaNameType,
	testDataHeaders []string,
	testDataRows [][]string) {

	var existInMap bool

	// Define a namespace UUID; this could be any valid UUID that you choose to use as a namespace for your IDs.
	// Here, we use the DNS namespace provided by the UUID package for demonstration purposes.
	namespace := uuidGenerator.NameSpaceDNS

	// Initialize your map
	var localTestDataDomainModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	localTestDataDomainModelMap = make(map[TestDataDomainUuidType]*TestDataDomainModelStruct)

	// Keeping track of all TestDataPoints for each Column
	var testDataPointsForColumns [][]*TestDataPointValueStruct

	// Initiate the maps used
	var tempTestDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForRowNameMap map[TestDataValueNameType]*[]TestDataPointRowUuidType
	var tempTestDataValuesForColumnMap map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForColumnAndRowUuidMap map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	var tempTestDataColumnsMetaDataMap map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct
	var tempUniqueTestDataValuesForColumnMap map[TestDataColumnUuidType]*map[TestDataValueType][]TestDataPointRowUuidType
	tempTestDataDomainAndAreaNameToUuidMap = make(map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType)
	tempTestDataValuesForRowMap = make(map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForRowNameMap = make(map[TestDataValueNameType]*[]TestDataPointRowUuidType)
	tempTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*[]*TestDataPointValueStruct)
	tempTestDataValuesForColumnAndRowUuidMap = make(map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct)
	tempTestDataColumnsMetaDataMap = make(map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct)
	tempUniqueTestDataValuesForColumnMap = make(map[TestDataColumnUuidType]*map[TestDataValueType][]TestDataPointRowUuidType)

	// Add the TestDataDomain and TestDataArea to map that map for Name to UUID conversion
	tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(testDataDomainName)] = TestDataDomainOrAreaUuidType(testDataDomainUuid)
	tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(testDataAreaName)] = TestDataDomainOrAreaUuidType(testDataAreaUuid)

	// TestData for one Area within one Domain
	var testDataArea *TestDataAreaStruct
	testDataArea = &TestDataAreaStruct{
		TestDataDomainUuid:                   testDataDomainUuid,
		TestDataDomainName:                   testDataDomainName,
		TestDataAreaUuid:                     testDataAreaUuid,
		TestDataAreaName:                     testDataAreaName,
		TestDataValuesForRowMap:              &tempTestDataValuesForRowMap,
		TestDataValuesForRowNameMap:          &tempTestDataValuesForRowNameMap,
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
	localTestDataDomainModelMap[testDataDomainUuid] = &tempTestDataDomainModel

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
	for _, tempTestDataRow := range testDataRows {

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
				TestDataValueNameDescription: "TestDataDomainName/TestDataAreaName/",
				TestDataValueName:            TestDataValueNameType(testDataDomainName) + "/" + TestDataValueNameType(testDataAreaName) + "/",
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

		var tempTestDataValuesForRowUuidSlice []*map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct

		// Loop the Values in the row and add 'TestDataPointName'
		for _, testDataPoint := range testDataPointsForRow {
			testDataPoint.TestDataValueNameDescription = testDataPoint.TestDataValueNameDescription + TestDataValueNameDescriptionType(testDataPointNameDescription)
			testDataPoint.TestDataValueName = testDataPoint.TestDataValueName + TestDataValueNameType(testDataPointName)
		}

		// Add 'testDataPointsForRow' to Map for TestDataPoints in one row
		tempTestDataValuesForRowMap[TestDataPointRowUuidType(rowUuid.String())] = &testDataPointsForRow

		// Add the map the slices of maps
		tempTestDataValuesForRowUuidSlice = append(tempTestDataValuesForRowUuidSlice, &tempTestDataValuesForRowMap)

	}

	// Loop 'testDataPointsForColumns' and add to column Map and create ColumnMetaData
	var tempTestDataColumnUuid TestDataColumnUuidType
	for _, testDataPointsForColumn := range testDataPointsForColumns {

		tempTestDataColumnUuid = testDataPointsForColumn[0].TestDataColumnUuid
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn

		// initialized the 'inner' map for unique values in the column
		var innerMapWithUniqueColumnValues map[TestDataValueType][]TestDataPointRowUuidType
		innerMapWithUniqueColumnValues = make(map[TestDataValueType][]TestDataPointRowUuidType)

		// Loop all Column data and add to the 'inner' Map, witch will create a map with only one value per occurrence of the value
		tempTestDataValuesForColumnMap[tempTestDataColumnUuid] = &testDataPointsForColumn
		for _, testDataPoint := range testDataPointsForColumn {

			// Extract existing slice
			var testDataPointRowUuidSlice []TestDataPointRowUuidType
			testDataPointRowUuidSlice, _ = innerMapWithUniqueColumnValues[testDataPoint.TestDataValue]

			// Append to slice
			testDataPointRowUuidSlice = append(testDataPointRowUuidSlice, testDataPoint.TestDataPointRowUuid)

			// Store back the slice
			innerMapWithUniqueColumnValues[testDataPoint.TestDataValue] = testDataPointRowUuidSlice
		}

		// Add the 'inner' map to the Column-Map
		tempUniqueTestDataValuesForColumnMap[tempTestDataColumnUuid] = &innerMapWithUniqueColumnValues

		// Create column MetaData
		var tempTestDataColumnMetaDataStruct *TestDataColumnMetaDataStruct
		tempTestDataColumnMetaDataStruct = &TestDataColumnMetaDataStruct{
			TestDataColumnUuid:     testDataPointsForColumn[0].TestDataColumnUuid,
			TestDataColumnDataName: testDataPointsForColumn[0].TestDataColumnDataName,
			TestDataColumnUIName:   testDataPointsForColumn[0].TestDataColumnUIName,
			//TestDataPointRowsUuid: 						,
			ShouldColumnBeUsedForFindingTestData:    trueColumns[string(testDataPointsForColumn[0].TestDataColumnDataName)],
			ShouldColumnBeUsedWithinTestDataSetName: trueColumns[string(testDataPointsForColumn[0].TestDataColumnDataName)],
		}

		// Add Column MetaData to Map
		tempTestDataColumnsMetaDataMap[tempTestDataColumnUuid] = tempTestDataColumnMetaDataStruct
	}

	// Loop 'TestDataValuesForRowMap' to be able to create 'TestDataValuesForRowNameMap'
	var testDataPointRowName TestDataValueNameType
	var tempTestDataValuesForRowSlice []*TestDataPointValueStruct
	//var localTestDataValuesForRowMapSlice []TestDataPointRowUuidType
	for _, tempTestDataValuesForRowSlicePtr := range tempTestDataValuesForRowMap {
		tempTestDataValuesForRowSlice = *tempTestDataValuesForRowSlicePtr
		//localTestDataPointRowUuid
		testDataPointRowName = tempTestDataValuesForRowSlice[0].TestDataValueName
		localTestDataValuesForRowMapSlicePtr, existInMap := tempTestDataValuesForRowNameMap[testDataPointRowName]
		if existInMap == false {
			var tempLocalTestDataValuesForRowMapSlice []TestDataPointRowUuidType
			tempLocalTestDataValuesForRowMapSlice = append(tempLocalTestDataValuesForRowMapSlice, tempTestDataValuesForRowSlice[0].TestDataPointRowUuid)

			tempTestDataValuesForRowNameMap[testDataPointRowName] = &tempLocalTestDataValuesForRowMapSlice
			localTestDataValuesForRowMapSlicePtr = &tempLocalTestDataValuesForRowMapSlice
		} else {
			//localTestDataValuesForRowMapSlice = *localTestDataValuesForRowMapSlicePtr

			*localTestDataValuesForRowMapSlicePtr = append(*localTestDataValuesForRowMapSlicePtr, tempTestDataValuesForRowSlice[0].TestDataPointRowUuid)
		}

	}

	// Check how to add values to full 'testDataModel'
	if TestDataModel.TestDataModelMap == nil {

		// This is the first time the TestDataModel is populated
		TestDataModel = TestDataModelStruct{
			TestDataDomainAndAreaNameToUuidMap: &tempTestDataDomainAndAreaNameToUuidMap,
			TestDataModelMap:                   &localTestDataDomainModelMap,
		}

	} else {

		// Add values from 'tempTestDataDomainAndAreaNameToUuidMap' into full 'testDataModel'
		var currentTestDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
		currentTestDataDomainAndAreaNameToUuidMap = *TestDataModel.TestDataDomainAndAreaNameToUuidMap

		// Only need to loop and add to map wíthout checking if it exists
		for domainOrAreaName, domainOrAreaUuid := range tempTestDataDomainAndAreaNameToUuidMap {
			currentTestDataDomainAndAreaNameToUuidMap[domainOrAreaName] = domainOrAreaUuid
		}
		TestDataModel.TestDataDomainAndAreaNameToUuidMap = &currentTestDataDomainAndAreaNameToUuidMap

		// Add values from 'localTestDataDomainModelMap'  into full 'testDataModel'
		var currentTestDataDomainModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
		currentTestDataDomainModelMap = *TestDataModel.TestDataModelMap

		for localTestDataDomainUuid, localTestDataDomainModel := range localTestDataDomainModelMap {

			// Check if 'localTestDataDomainUuid' exists within Full TestDataModel
			var oneFullTestDataDomainModelPtr *TestDataDomainModelStruct
			oneFullTestDataDomainModelPtr, existInMap = currentTestDataDomainModelMap[localTestDataDomainUuid]
			if existInMap == true {
				var oneFullTestDataDomainModel TestDataDomainModelStruct
				oneFullTestDataDomainModel = *oneFullTestDataDomainModelPtr

				// Loop 'localTestDataDomainModel' and add each TestDataAreasMap to 'oneFullTestDataDomainModel'
				for oneLocalTestDataAreaUuid, oneLocalTestDataArea := range *localTestDataDomainModel.TestDataAreasMap {

					// Check if oneLocalTestDataAreaUuid exists within Full TestDataModel for current localTestDataDomainUuid
					var oneFullTestDataAreasMap map[TestDataAreaUuidType]*TestDataAreaStruct
					oneFullTestDataAreasMap = *oneFullTestDataDomainModel.TestDataAreasMap

					// Add 'oneLocalTestDataArea' to oneFullTestDataAreasMap
					oneFullTestDataAreasMap[oneLocalTestDataAreaUuid] = oneLocalTestDataArea

				}

			} else {

				// Nothing exist for this 'TestDataDomainUuid' so just add it
				currentTestDataDomainModelMap[localTestDataDomainUuid] = localTestDataDomainModel

			}

		}

		TestDataModel.TestDataModelMap = &currentTestDataDomainModelMap

	}

}
