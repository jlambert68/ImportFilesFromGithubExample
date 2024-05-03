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
