package testDataSelector

// Updates the list that show the TestDataPoints for a specific Group in main window
func updateTestDataPointsForAGroupList(testDataPointGroupName testDataPointGroupNameType) {

	// Clear the slice that holds all TestDataPoints
	testDataPointsForAGroup = nil

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap testDataPointNameMapType
	tempTestDataPointNameMap = *chosenTestDataPointsPerGroupMap[testDataPointGroupName]

	// Refill the slice with all TestDataPoints
	for testDataValueName, _ := range tempTestDataPointNameMap {

		testDataPointsForAGroup = append(testDataPointsForAGroup, testDataValueName)

	}

	// Refresh the List in the UI
	testDataPointsForAGroupList.Refresh()

}
