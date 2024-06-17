package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"regexp"
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

	var testDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	testDataModelMap = *testDataModel.TestDataModelMap

	//var testDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	//testDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap

	parent.Hide()

	var shouldUpdateMainWindow responseChannelStruct

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

	// Slices used to keep track of filtered, available and selected DataPoints
	var allPointsAvailable []dataPointTypeForListsStruct
	var allSelectedPoints []dataPointTypeForListsStruct

	// The List-widget holding all available TestDataPoints from Search
	var allAvailablePointsList *widget.List

	// The List-widget holding all selected TestDataPoints from Search
	var selectedPointsList *widget.List

	// *** Create the selection boxes for selecting TestDataValues values
	var testDataSelectionsContainer *fyne.Container

	// Create the container for Search- and Clear- buttons
	var searchAndClearButtonsContainer *fyne.Container

	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct

	// Layout configuration for the new/edit window
	// Create the UpperAndLowerSplitContainer
	var upperAndLowerSplitContainer *container.Split
	var listsSplitContainer *container.Split

	var upperSplitContainer *fyne.Container

	var lowerRightSideContainer *fyne.Container

	// *** Create the selection boxes for selecting TestDataValues values
	generateTestDataSelectionsUIComponent(
		testDataSelectionsContainer,
		testDataModelMap,
		testDataModel,
		&tempTestDataValuesForRowMap,
		&allPointsAvailable,
		&allSelectedPoints,
		allAvailablePointsList,
		searchAndClearButtonsContainer)

	// Create and configure the list-component of all TestDataPoints
	generateAllAvailablePointsListUIComponent(
		allAvailablePointsList,
		&allPointsAvailable,
		&allSelectedPoints,
		&newOrEditTestDataPointGroupWindow,
		selectedPointsList,
		testDataModel)

	// Create and configure the list-component of selected TestDataPoints
	generateSelectedPointsListUIComponent(
		selectedPointsList,
		&allPointsAvailable,
		&allSelectedPoints,
		&newOrEditTestDataPointGroupWindow,
		lowerRightSideContainer,
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

// Build the Table Data, based on TestDataPointName, to be used when the popup table is shown to the user to pick from
func buildPopUpTableDataFromTestDataPointName(
	tempTestDataPointRowName string,
	testDataModel *TestDataModelStruct) (
	tableData [][]string) {

	re := regexp.MustCompile(`^([^/]+)/([^/]+)`)

	matches := re.FindStringSubmatch(tempTestDataPointRowName)
	if len(matches) > 2 {
		fmt.Println("First part:", matches[1])
		fmt.Println("Second part:", matches[2])
	} else {
		fmt.Println("No matching parts found")
	}

	var tempTestDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	var tempTestDataDomainModel TestDataDomainModelStruct
	var tempTestDataAreaMap map[TestDataAreaUuidType]*TestDataAreaStruct
	var tempTestDataArea TestDataAreaStruct
	var tempTestDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	var tempTestDataValuesForRowNameMap map[TestDataValueNameType]*[]TestDataPointRowUuidType
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForRowUuidMapBaseOnNameSlice []TestDataPointRowUuidType

	var tempTestDataDomainOrAreaUuid TestDataDomainOrAreaUuidType
	var tempTestDataDomainUuid TestDataDomainUuidType
	var tempTestDataAreaUuid TestDataAreaUuidType

	tempTestDataModelMap = *testDataModel.TestDataModelMap

	// Extract the UUID for Domain and Area
	tempTestDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(matches[1])]
	tempTestDataDomainUuid = TestDataDomainUuidType(tempTestDataDomainOrAreaUuid)
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(matches[2])]
	tempTestDataAreaUuid = TestDataAreaUuidType(tempTestDataDomainOrAreaUuid)

	// Extract Domain and Area maps
	tempTestDataDomainModel = *tempTestDataModelMap[tempTestDataDomainUuid]
	tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
	tempTestDataArea = *tempTestDataAreaMap[tempTestDataAreaUuid]
	tempTestDataValuesForRowNameMap = *tempTestDataArea.TestDataValuesForRowNameMap
	tempTestDataValuesForRowMap = *tempTestDataArea.TestDataValuesForRowMap

	var tempTestDataPointRowNameToSearchFor string
	tempTestDataPointRowNameToSearchFor = tempTestDataPointRowName //[len(matches[0]+"/"):]

	tempTestDataValuesForRowUuidMapBaseOnNameSlice = *tempTestDataValuesForRowNameMap[TestDataValueNameType(tempTestDataPointRowNameToSearchFor)]

	fmt.Println(tempTestDataValuesForRowUuidMapBaseOnNameSlice)

	// Loop the slice to extract the RowUUids

	var headerSlice []string
	for rowIndex, tempTestDataPointRowUuid := range tempTestDataValuesForRowUuidMapBaseOnNameSlice {

		var rowSlice []string
		tempTestDataValuesForRowSlice := tempTestDataValuesForRowMap[tempTestDataPointRowUuid]

		// Loop the slice with RowValue
		for _, tempTestDataPointValue := range *tempTestDataValuesForRowSlice {

			// Create a header slice
			if rowIndex == 0 {
				headerSlice = append(headerSlice, string(tempTestDataPointValue.TestDataColumnUIName))
			}

			rowSlice = append(rowSlice, string(tempTestDataPointValue.TestDataValue))
		}

		// Append data for 'TestDataPointRowName'
		rowSlice = append(rowSlice, tempTestDataPointRowName)

		// Append data for unique 'TestDataPointRowUuid'
		rowSlice = append(rowSlice, string(tempTestDataPointRowUuid))

		// Add a header when first row
		if rowIndex == 0 {

			// Add headers for 'TestDataPointRowName' and 'TestDataPointRowUuid'
			headerSlice = append(headerSlice, "TestDataPointRowName")
			headerSlice = append(headerSlice, "TestDataPointRowUuid")

			tableData = append(tableData, headerSlice)
		}

		// Add the data
		tableData = append(tableData, rowSlice)

	}

	return tableData
}

/*
// Build the Table Data, based on TestDataPointUuid, to be used when the popup table is shown to the user to pick from
func buildPopUpTableDataFromTestDataPointUuid(
	tempTestDataPointRowUuids []string,
	testDataModel *TestDataModelStruct) (
	tableData [][]string) {

	re := regexp.MustCompile(`^([^/]+)/([^/]+)`)

	matches := re.FindStringSubmatch(tempTestDataPointRowName)
	if len(matches) > 2 {
		fmt.Println("First part:", matches[1])
		fmt.Println("Second part:", matches[2])
	} else {
		fmt.Println("No matching parts found")
	}

	var tempTestDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
	var tempTestDataDomainModel TestDataDomainModelStruct
	var tempTestDataAreaMap map[TestDataAreaUuidType]*TestDataAreaStruct
	var tempTestDataArea TestDataAreaStruct
	var tempTestDataDomainAndAreaNameToUuidMap map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	var tempTestDataValuesForRowNameMap map[TestDataValueNameType]*[]TestDataPointRowUuidType
	var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForRowUuidMapBaseOnNameSlice []TestDataPointRowUuidType

	var tempTestDataDomainOrAreaUuid TestDataDomainOrAreaUuidType
	var tempTestDataDomainUuid TestDataDomainUuidType
	var tempTestDataAreaUuid TestDataAreaUuidType

	tempTestDataModelMap = *testDataModel.TestDataModelMap

	// Extract the UUID for Domain and Area
	tempTestDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(matches[1])]
	tempTestDataDomainUuid = TestDataDomainUuidType(tempTestDataDomainOrAreaUuid)
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[TestDataDomainOrAreaNameType(matches[2])]
	tempTestDataAreaUuid = TestDataAreaUuidType(tempTestDataDomainOrAreaUuid)

	// Extract Domain and Area maps
	tempTestDataDomainModel = *tempTestDataModelMap[tempTestDataDomainUuid]
	tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
	tempTestDataArea = *tempTestDataAreaMap[tempTestDataAreaUuid]
	tempTestDataValuesForRowNameMap = *tempTestDataArea.TestDataValuesForRowNameMap
	tempTestDataValuesForRowMap = *tempTestDataArea.TestDataValuesForRowMap

	var tempTestDataPointRowNameToSearchFor string
	tempTestDataPointRowNameToSearchFor = tempTestDataPointRowName //[len(matches[0]+"/"):]

	tempTestDataValuesForRowUuidMapBaseOnNameSlice = *tempTestDataValuesForRowNameMap[TestDataValueNameType(tempTestDataPointRowNameToSearchFor)]

	fmt.Println(tempTestDataValuesForRowUuidMapBaseOnNameSlice)

	// Loop the slice to extract the RowUUids

	var headerSlice []string
	for rowIndex, tempTestDataPointRowUuid := range tempTestDataValuesForRowUuidMapBaseOnNameSlice {

		var rowSlice []string
		tempTestDataValuesForRowSlice := tempTestDataValuesForRowMap[tempTestDataPointRowUuid]

		// Loop the slice with RowValue
		for _, tempTestDataPointValue := range *tempTestDataValuesForRowSlice {

			// Create a header slice
			if rowIndex == 0 {
				headerSlice = append(headerSlice, string(tempTestDataPointValue.TestDataColumnUIName))
			}

			rowSlice = append(rowSlice, string(tempTestDataPointValue.TestDataValue))
		}

		// Append data for 'TestDataPointRowName'
		rowSlice = append(rowSlice, tempTestDataPointRowName)

		// Append data for unique 'TestDataPointRowUuid'
		rowSlice = append(rowSlice, string(tempTestDataPointRowUuid))

		// Add a header when first row
		if rowIndex == 0 {

			// Add headers for 'TestDataPointRowName' and 'TestDataPointRowUuid'
			headerSlice = append(headerSlice, "TestDataPointRowName")
			headerSlice = append(headerSlice, "TestDataPointRowUuid")

			tableData = append(tableData, headerSlice)
		}

		// Add the data
		tableData = append(tableData, rowSlice)

	}

	return tableData
}

*/
