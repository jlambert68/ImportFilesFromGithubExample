package testDataSelector

import (
	"sort"
)

// ListTestDataGroups
// List the current TestDataGroups that the User has
func ListTestDataGroups() (testDataPointGroupsAsStringSlice []string) {

	// Loop all 'testDataPointGroups'
	for _, tempTestDataPointGroup := range testDataPointGroups {
		testDataPointGroupsAsStringSlice = append(testDataPointGroupsAsStringSlice, string(tempTestDataPointGroup))
	}

	// Sort the Groups
	sort.Strings(testDataPointGroupsAsStringSlice)

	return testDataPointGroupsAsStringSlice

}

// ListTestDataGroupPointsForAGroup
// List the current TestDataGroupPoints for a specific TestDataGroup
func ListTestDataGroupPointsForAGroup(testDataGroup string) (testDataPointGroupsAsStringSlice []string) {

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap testDataPointNameMapType
	tempTestDataPointNameMap = *chosenTestDataPointsPerGroupMap[testDataPointGroupNameType(testDataGroup)]

	// Refill the slice with all TestDataPoints
	for testDataPoint, _ := range tempTestDataPointNameMap {
		testDataPointGroupsAsStringSlice = append(testDataPointGroupsAsStringSlice, string(testDataPoint))

	}

	// Sort the GroupPoints
	sort.Strings(testDataPointGroupsAsStringSlice)

	return testDataPointGroupsAsStringSlice

}

// ListTestDataRowsForAGroupPoint
// List the current TestDataRow for a specific TestDataGroupPoint
func ListTestDataRowsForAGroupPoint(testDataGroup string, testDataGroupPoint string) (testDataGroupPointRowsAsStringSlice []string) {

	//fixa denna

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap testDataPointNameMapType
	var dataPointRowsSlicePtr *[]*dataPointTypeForGroupsStruct
	var dataPointRowsSlice []*dataPointTypeForGroupsStruct

	// Extract DataPoints from for Group
	tempTestDataPointNameMap = *chosenTestDataPointsPerGroupMap[testDataPointGroupNameType(testDataGroup)]

	// Extract Rows for DataPoint
	dataPointRowsSlicePtr = tempTestDataPointNameMap[TestDataValueNameType(testDataGroupPoint)]
	dataPointRowsSlice = *dataPointRowsSlicePtr

	// Refill the slice with all TestDataPoints
	for _, testDataPointRowUuiObject := range dataPointRowsSlice[0].selectedTestDataPointUuidMap {
		testDataGroupPointRowsAsStringSlice = append(testDataGroupPointRowsAsStringSlice, string(testDataPointRowUuiObject.testDataPointRowUuid))

	}

	// Sort the GroupPoints
	sort.Strings(testDataGroupPointRowsAsStringSlice)

	return testDataGroupPointRowsAsStringSlice

}

// GetTestDataPointValues
// Generate a map with 'TestDataColumnDataName' as key and 'TestDataValue' as value
func GetTestDataPointValues(
	testDataPointName string,
	testDataPointRowUuid string) (
	testDataColumnDataNameMap map[string]string) { // map[TestDataColumnDataNameType]TestDataValueType

	// Initiate response-map
	testDataColumnDataNameMap = make(map[string]string)

	if testDataPointName == "" || testDataPointRowUuid == "" {
		return testDataColumnDataNameMap
	}

	// Create the data table for all matching 'testDataPointRowUuid'
	var tableData [][]string
	tableData = buildPopUpTableDataFromTestDataPointName(testDataPointName, testDataModelPtr)

	var headerSlice []string

	// Loop alla rows
	for rowIndex, rowData := range tableData {

		// Loop alla values for row
		for columnIndex, columnValue := range rowData {

			// Create a header slice
			if rowIndex == 0 {
				// Header row
				headerSlice = append(headerSlice, columnValue)

			} else {
				// Only process if this is the correct row
				if rowData[len(rowData)-1] == testDataPointRowUuid {

					testDataColumnDataNameMap[headerSlice[columnIndex]] = columnValue

				}
			}

		}

	}

	return testDataColumnDataNameMap

}
