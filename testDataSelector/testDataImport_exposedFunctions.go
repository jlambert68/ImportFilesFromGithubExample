package testDataSelector

import "sort"

// ListTestDataGroups
// List the current TestDataGroups that the User has
func ListTestDataGroups() (testDataGroups []string) {

	// Convert TestDataGroups into '[]string'
	var testDataPointGroupsAsStringSlice []string

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
func ListTestDataGroupPointsForAGroup(testDataGroup string) (testDataGroupPoints []string) {

	var testDataPointGroupsAsStringSlice []string

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

// GetTestDataPointValues
// Generate a map with 'TestDataColumnDataName' as key and 'TestDataValue' as value
func GetTestDataPointValues(testDataGroupPoint string) (testDataColumnDataNameMap map[string]string) { // map[TestDataColumnDataNameType]TestDataValueType

	return testDataColumnDataNameMap

	/*
		// Initiate response-map
		testDataColumnDataNameMap = make(map[string]string)

		if testDataGroupPoint == "" {
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

	*/
}
