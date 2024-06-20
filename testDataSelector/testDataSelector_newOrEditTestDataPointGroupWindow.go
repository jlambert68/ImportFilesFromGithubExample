package testDataSelector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	testDataDomainLabelText   string = "Available Domains for TestData"
	testDataTestAreaLabelText string = "Available TestAreas for domain "
)

func showNewOrEditGroupWindow(
	app fyne.App,
	parent fyne.Window,
	isNew bool,
	responseChannel *chan responseChannelStruct,
	incomingGroupName testDataPointGroupNameType,
	newOrEditedChosenTestDataPointsThisGroupMapPtr *map[testDataPointGroupNameType]*testDataPointNameMapType,
	testDataModel *TestDataModelStruct) {

	// Store reference to TestDataModel
	testDataModelPtr = testDataModel

	var testDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	testDataModelMap = *testDataModel.TestDataModelMap

	//var testDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	//testDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap

	parent.Hide()

	// Set Default value
	shouldUpdateMainWindow = responseChannelStruct{
		shouldBeUpdated:        false,
		testDataPointGroupName: "",
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
		*responseChannel <- shouldUpdateMainWindow
	})

	// Create and configure the list-component of all TestDataPoints
	generateAllAvailablePointsListUIComponent(&newOrEditTestDataPointGroupWindow, testDataModel)

	// *** Create the selection boxes for selecting TestDataValues values
	generateTestDataSelectionsUIComponent(testDataModel, testDataModelMap)

	// Create and configure the list-component of selected TestDataPoints
	generateSelectedPointsListUIComponent(
		&newOrEditTestDataPointGroupWindow,
		incomingGroupName,
		isNew,
		newOrEditedChosenTestDataPointsThisGroupMapPtr)

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
func testDataPointIntersectionOfTwoSlices(firstSlice, secondSlice []TestDataPointRowUuidType) []TestDataPointRowUuidType {
	// Use firstSlice map to count occurrences of elements in the first slice
	elemCount := make(map[TestDataPointRowUuidType]bool)

	// Fill the map with elements from the first slice
	for _, item := range firstSlice {
		elemCount[item] = true
	}

	// Create firstSlice slice to hold the intersectionSlice
	var intersectionSlice []TestDataPointRowUuidType

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
