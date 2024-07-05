package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	testDataDomainLabelText   string = "Available Domains for TestData"
	testDataTestAreaLabelText string = "Available TestAreas for domain "
)

func ShowNewOrEditGroupWindow(
	app fyne.App,
	parent fyne.Window,
	isNew bool,
	responseChannel *chan testDataEngine.ResponseChannelStruct,
	incomingGroupName testDataEngine.TestDataPointGroupNameType,
	newOrEditedChosenTestDataPointsThisGroupMapPtr *map[testDataEngine.TestDataPointGroupNameType]*testDataEngine.TestDataPointNameMapType,
	testDataForGroupObject *testDataEngine.TestDataForGroupObjectStruct) {

	var testDataModelMap map[testDataEngine.TestDataDomainUuidType]*testDataEngine.TestDataDomainModelStruct
	testDataModelMap = *testDataEngine.TestDataModel.TestDataModelMap

	//var testDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	//testDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap

	parent.Hide()

	// Set Default value
	testDataForGroupObject.ShouldUpdateMainWindow = testDataEngine.ResponseChannelStruct{
		ShouldBeUpdated:        false,
		TestDataPointGroupName: "",
	}

	// Set name for the Window
	newOrEditTestDataPointGroupWindow := app.NewWindow("Edit Group")
	if isNew {
		newOrEditTestDataPointGroupWindow.SetTitle("New Group")
	}

	newOrEditTestDataPointGroupWindow.Resize(fyne.NewSize(600, 1000))

	// When this window closed then show parent and send response to parent window
	newOrEditTestDataPointGroupWindow.SetOnClosed(func() {
		parent.Show()
		*responseChannel <- testDataForGroupObject.ShouldUpdateMainWindow
	})

	// Create and configure the list-component of all TestDataPoints
	generateAllAvailablePointsListUIComponent(&newOrEditTestDataPointGroupWindow, &testDataEngine.TestDataModel)

	// *** Create the selection boxes for selecting TestDataValues values
	generateTestDataSelectionsUIComponent(&testDataEngine.TestDataModel, testDataModelMap)

	// Create and configure the list-component of selected TestDataPoints
	generateSelectedPointsListUIComponent(
		&newOrEditTestDataPointGroupWindow,
		incomingGroupName,
		isNew,
		newOrEditedChosenTestDataPointsThisGroupMapPtr,
		testDataForGroupObject)

	var tempTestDataPointsLabel *widget.Label
	tempTestDataPointsLabel = widget.NewLabel("TestDataPoints based on filter")
	tempTestDataPointsLabel.TextStyle.Bold = true

	var lowerLeftSideContainer *fyne.Container
	lowerLeftSideContainer = container.NewBorder(tempTestDataPointsLabel, nil, nil, nil, allAvailablePointsList)

	listsSplitContainer = container.NewHSplit(lowerLeftSideContainer, lowerRightSideContainer)

	upperSplitContainer = container.NewBorder(nil, searchAndClearButtonsContainer, nil, nil, testDataSelectionsContainer)
	//lowerSplitContainer = container.NewVBox(entryContainer, buttonsContainer, searchAndClearButtonsContainer, listsSplitContainer)

	upperAndLowerSplitContainer = container.NewVSplit(upperSplitContainer, listsSplitContainer)

	//windowContent = container.NewBorder(container.NewVBox(entryContainer, buttonsContainer, testDataSelectionsContainer, searchAndClearButtonsContainer), nil, nil, nil, listsSplitContainer)

	newOrEditTestDataPointGroupWindow.SetContent(upperAndLowerSplitContainer)
	newOrEditTestDataPointGroupWindow.Show()

}

// testDataPointIntersectionOfTwoSlices returns a new slice containing only the elements that appear in both a and b.
func testDataPointIntersectionOfTwoSlices(firstSlice, secondSlice []testDataEngine.TestDataPointRowUuidType) []testDataEngine.TestDataPointRowUuidType {
	// Use firstSlice map to count occurrences of elements in the first slice
	elemCount := make(map[testDataEngine.TestDataPointRowUuidType]bool)

	// Fill the map with elements from the first slice
	for _, item := range firstSlice {
		elemCount[item] = true
	}

	// Create firstSlice slice to hold the intersectionSlice
	var intersectionSlice []testDataEngine.TestDataPointRowUuidType

	// Check each element in the second slice; if it's in the map, add to the intersectionSlice
	for _, item := range secondSlice {
		if _, found := elemCount[item]; found {
			intersectionSlice = append(intersectionSlice, item)
			// Optional: Remove item from map if you don't expect duplicates or don't need to count them
			delete(elemCount, item)
		}
	}

	return intersectionSlice
}
