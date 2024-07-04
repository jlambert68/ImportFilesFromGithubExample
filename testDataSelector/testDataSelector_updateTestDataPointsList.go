package testDataSelector

import (
	"ImportFilesFromGithub/newOrEditTestDataPointGroupUI"
	"ImportFilesFromGithub/testDataEngine"
)

// Updates the list that show the TestDataPoints for a specific Group in main window
func updateTestDataPointsForAGroupList(
	testDataPointGroupName testDataEngine.TestDataPointGroupNameType,
	testDataModelObject *testDataEngine.TestDataModelObjectStruct) {

	// Clear the slice that holds all TestDataPoints
	testDataModelObject.TestDataPointsForAGroup = nil

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap testDataEngine.TestDataPointNameMapType
	tempTestDataPointNameMap = *testDataModelObject.ChosenTestDataPointsPerGroupMap[testDataPointGroupName]

	// Refill the slice with all TestDataPoints
	for testDataValueName, _ := range tempTestDataPointNameMap {

		testDataModelObject.TestDataPointsForAGroup = append(testDataModelObject.TestDataPointsForAGroup, testDataValueName)

	}

	// Refresh the List in the UI
	newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.Refresh()

}
