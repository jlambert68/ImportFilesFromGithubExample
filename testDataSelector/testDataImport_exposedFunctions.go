package testDataSelector

import "sort"

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
	for _, testDataPointRowUuidPtr := range dataPointRowsSlice[0].selectedTestDataPointUuidMap {
		testDataGroupPointRowsAsStringSlice = append(testDataGroupPointRowsAsStringSlice, string(testDataPointRowUuidPtr))

	}

	// Sort the GroupPoints
	sort.Strings(testDataGroupPointRowsAsStringSlice)

	return testDataGroupPointRowsAsStringSlice

}

// GetTestDataPointValues
// Generate a map with 'TestDataColumnDataName' as key and 'TestDataValue' as value
func GetTestDataPointValues(
	testDataGroupPoint string,
	testDataGroupPointRowUuid string) (
	testDataColumnDataNameMap map[string]string) { // map[TestDataColumnDataNameType]TestDataValueType

	// Initiate response-map
	testDataColumnDataNameMap = make(map[string]string)

	if testDataGroupPoint == "" || testDataGroupPointRowUuid == "" {
		return testDataColumnDataNameMap
	}

	var tempTestDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	var tempTestDataDomainModel TestDataDomainModelStruct
	var tempTestDataAreaMap map[TestDataAreaUuidType]*TestDataAreaStruct
	var tempTestDataArea TestDataAreaStruct
	var tempTestDataValuesForRowNameMap map[TestDataValueNameType]*[]*map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataPointValueSlice []*TestDataPointValueStruct

	tempTestDataModelMap = *testDataModelRef.TestDataModelMap
	tempTestDataDomainModel = *tempTestDataModelMap[testDataDomainUuid]
	tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
	tempTestDataArea = *tempTestDataAreaMap[testDataAreaUuid]
	tempTestDataValuesForRowNameMap = *tempTestDataArea.TestDataValuesForRowNameMap

	// Initiate response-map
	testDataColumnDataNameMap = make(map[string]string)

	// Extract correct TestDataPointRow
	tempTestDataPointValueSlice = *tempTestDataValuesForRowNameMap[TestDataValueNameType(testDataGroupPoint)]

	// Loop Values in slice and create response-map
	for _, tempTestDataPoint := range tempTestDataPointValueSlice {

		testDataColumnDataNameMap[string(tempTestDataPoint.TestDataColumnDataName)] = string(tempTestDataPoint.TestDataValue)

	}

	return testDataColumnDataNameMap

}
