package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"sort"
	"strings"
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

	var newOrEditedChosenTestDataPointsThisGroupMap map[testDataPointGroupNameType]*testDataPointNameMapType
	newOrEditedChosenTestDataPointsThisGroupMap = *newOrEditedChosenTestDataPointsThisGroupMapPtr

	var saveButton *widget.Button

	newOrEditTestDataPointGroupWindow := app.NewWindow("Edit Group")
	if isNew {
		newOrEditTestDataPointGroupWindow.SetTitle("New Group")
	}

	newOrEditTestDataPointGroupWindow.Resize(fyne.NewSize(600, 500))

	// When this window closed then show parent and send response to parent window
	newOrEditTestDataPointGroupWindow.SetOnClosed(func() {
		parent.Show()
		*responseChannel <- shouldUpdateMainWindow
	})

	var existInMap bool

	// Slices used to keep track of filtered, available and selected DataPoints
	var filteredTestDataPoints []string
	var allPointsAvailable []string
	var allSelectedPoints []string

	// The List-widget holding all available TestDataPoints from Search
	var allAvailablePointsList *widget.List

	// If existing groupToEdit then extract points from it otherwise create an empty selected points slice
	var selectedPointsPtr *testDataPointNameMapType
	var selectedPoints testDataPointNameMapType

	if isNew == false {

		selectedPointsPtr = newOrEditedChosenTestDataPointsThisGroupMap[incomingGroupName]
		selectedPoints = *selectedPointsPtr

	}

	// *** Create the selection boxes for selecting TestDataValues values
	var testDataSelectionsContainer *fyne.Container

	var domainOptions []string
	var domains []*TestDataDomainModelStruct
	var domainsLabel *widget.Label
	var domainsSelect *widget.Select
	var testDomainContainer *fyne.Container

	var testAreaOptions []string
	var testAreas []*TestDataAreaStruct
	var testAreasLabel *widget.Label
	var testAreaSelect *widget.Select
	var testAreasContainer *fyne.Container
	var testAreaMap *map[TestDataAreaUuidType]*TestDataAreaStruct

	type testDataValueSelectionStruct struct {
		testDataSelectionLabel       *widget.Label
		testDataCheckGroup           *widget.CheckGroup
		TestDataColumnUuid           TestDataColumnUuidType
		TestDataColumnDataName       TestDataColumnDataNameType
		TestDataPointValueRowUuidMap *map[TestDataValueType]*[]TestDataPointRowUuidType
	}
	var testDataValueSelections []*testDataValueSelectionStruct
	var testDataValuesSelectionContainer *fyne.Container

	var tempTestDataPointRowUuidSliceInMap []TestDataPointRowUuidType

	// Create label for Domains
	domainsLabel = widget.NewLabel(testDataDomainLabelText)
	domainsLabel.TextStyle.Bold = true
	testAreasLabel = widget.NewLabel(testDataTestAreaLabelText)
	testAreasLabel.TextStyle.Bold = true

	// Extract TestData on Domain-level
	for _, tempTestDataDomainModel := range testDataModelMap {
		domainOptions = append(domainOptions, string(tempTestDataDomainModel.TestDataDomainName))
		domains = append(domains, tempTestDataDomainModel)
	}

	// Create Domain-Select-DropDown
	domainsSelect = widget.NewSelect(domainOptions, func(selected string) {

		// Clear UI object that need to be recreated

		// Extract correct TestArea
		for index, domain := range domains {
			if selected == string(domain.TestDataDomainName) {
				testAreaMap = domains[index].TestDataAreasMap
				break
			}
		}

		// Extract TestData on TestArea Level
		for _, tempTestDataArea := range *testAreaMap {
			testAreaOptions = append(testAreaOptions, string(tempTestDataArea.TestDataAreaName))
			testAreas = append(testAreas, tempTestDataArea)
		}

		// Create TestArea-Select-DropDown
		testAreaSelect = widget.NewSelect(testAreaOptions, func(selected string) {

			// Create available TestDataSelections for TestArea
			for _, testDataArea := range *testAreaMap {

				// Clear UI component that holds 'TestDataValuesSelections'
				testDataValuesSelectionContainer = container.NewHBox()

				// Create a slice with 'testDataColumnsMetaData' that can be sorted
				var testDataColumnsMetaDataToBeSorted []*TestDataColumnMetaDataStruct
				for _, testDataColumnsMetaData := range *testDataArea.TestDataColumnsMetaDataMap {
					testDataColumnsMetaDataToBeSorted = append(testDataColumnsMetaDataToBeSorted, testDataColumnsMetaData)
				}

				// Sort the slice based on TestDataColumnUIName
				sort.Slice(testDataColumnsMetaDataToBeSorted, func(i, j int) bool {
					return testDataColumnsMetaDataToBeSorted[i].TestDataColumnUIName < testDataColumnsMetaDataToBeSorted[j].TestDataColumnUIName
				})

				// Loop 'testDataColumnsMetaDataToBeSorted' for Columns to present as separate CheckGroups
				for _, testDataColumnsMetaData := range testDataColumnsMetaDataToBeSorted {

					// Check if column should be used for filtering TestData as a CheckGroup
					if testDataColumnsMetaData.ShouldColumnBeUsedForFindingTestData == true {

						var checkGroupOptions []string
						var tempTestDataColumnContainer *fyne.Container

						// Set Label
						var newColumnFilterLabel *widget.Label
						newColumnFilterLabel = widget.NewLabel(string(testDataColumnsMetaData.TestDataColumnUIName))
						newColumnFilterLabel.TextStyle.Bold = true

						var tempTestDataPointValueRowUuidMap map[TestDataValueType]*[]TestDataPointRowUuidType
						tempTestDataPointValueRowUuidMap = make(map[TestDataValueType]*[]TestDataPointRowUuidType)

						var testDataValueSelection *testDataValueSelectionStruct
						testDataValueSelection = &testDataValueSelectionStruct{
							testDataSelectionLabel:       newColumnFilterLabel,
							testDataCheckGroup:           nil,
							TestDataColumnUuid:           testDataColumnsMetaData.TestDataColumnUuid,
							TestDataColumnDataName:       testDataColumnsMetaData.TestDataColumnDataName,
							TestDataPointValueRowUuidMap: &tempTestDataPointValueRowUuidMap,
						}

						// Extract the Map with the values
						var uniqueTestDataValuesForColumnMapPtr *map[TestDataValueType][]TestDataPointRowUuidType
						UniqueTestDataValuesForColumnMap := *testDataArea.UniqueTestDataValuesForColumnMap

						uniqueTestDataValuesForColumnMapPtr = UniqueTestDataValuesForColumnMap[testDataColumnsMetaData.TestDataColumnUuid]

						// Loop Values in Column and create Checkboxes, and store RowUuids for unique values
						for uniqueTestDataValue, testDataPointRowsUuid := range *uniqueTestDataValuesForColumnMapPtr {

							// Add value to slice for CheckBox-labels
							checkGroupOptions = append(checkGroupOptions, string(uniqueTestDataValue))

							// Add 'TestDataPointRowUuid' to correct slice for each unique value in the column
							var testDataPointRowUuidSlicePtr *[]TestDataPointRowUuidType
							var testDataPointRowUuidSlice []TestDataPointRowUuidType
							testDataPointRowUuidSlicePtr, existInMap = tempTestDataPointValueRowUuidMap[uniqueTestDataValue]

							if existInMap == false {
								var tempTestDataPointRowUuidSlice []TestDataPointRowUuidType
								testDataPointRowUuidSlice = tempTestDataPointRowUuidSlice
							} else {
								testDataPointRowUuidSlice = *testDataPointRowUuidSlicePtr
							}

							testDataPointRowUuidSlice = append(testDataPointRowUuidSlice, testDataPointRowsUuid...)

							tempTestDataPointValueRowUuidMap[uniqueTestDataValue] = &testDataPointRowUuidSlice

						}

						// Sort values in CheckGroup
						sort.Strings(checkGroupOptions)

						// Create the CheckGroup
						var tempTestDataCheckGroup *widget.CheckGroup
						tempTestDataCheckGroup = widget.NewCheckGroup(checkGroupOptions, func(changed []string) {
							// Handle check change
						})

						// Add the CheckGroup
						testDataValueSelection.testDataCheckGroup = tempTestDataCheckGroup

						// Add 'testDataValueSelections' to slice
						testDataValueSelections = append(testDataValueSelections, testDataValueSelection)

						// Get the minimum size of the check group
						var testDataCheckGroupMinSize fyne.Size
						testDataCheckGroupMinSize = testDataValueSelection.testDataCheckGroup.MinSize()

						// Create the container having scrollbar the TestDataCheckGroup
						testDataCheckGroupContainer := container.NewScroll(testDataValueSelection.testDataCheckGroup)

						// Set
						testDataCheckGroupContainer.SetMinSize(fyne.NewSize(testDataCheckGroupContainer.Size().Height, testDataCheckGroupMinSize.Width))

						// Add to TestDataColumn-container
						tempTestDataColumnContainer = container.NewBorder(
							testDataValueSelection.testDataSelectionLabel,
							nil, nil, nil,
							testDataCheckGroupContainer)

						// Add 'tempTestDataColumnContainer' to 'testDataValuesSelectionContainer'
						testDataValuesSelectionContainer.Add(tempTestDataColumnContainer)

					}
				}
			}
		})

		// Set label for TestAreas
		testAreasLabel.SetText(fmt.Sprintf(testDataTestAreaLabelText+"'%s'", domainOptions[0]))

		// If there is only one item in TestArea-item then select that one
		if len(testAreaOptions) == 1 {
			testAreaSelect.SetSelected(testAreaOptions[0])
			testAreaSelect.Refresh()
		}

	})

	// If there is only one item in Domains-dropdown then select that one
	if len(domainOptions) == 1 {
		domainsSelect.SetSelected(domainOptions[0])
		domainsSelect.Refresh()

		// Set label for TestAreas
		testAreasLabel.SetText(fmt.Sprintf(testDataTestAreaLabelText+"'%s'", domainOptions[0]))
	}

	// Create the separate TestData-selection-containers
	testDomainContainer = container.NewVBox(domainsLabel, domainsSelect)
	testAreasContainer = container.NewVBox(testAreasLabel, testAreaSelect)

	// Create the main TestData-selection-container
	testDataSelectionsContainer = container.NewHBox(testDomainContainer, testAreasContainer, testDataValuesSelectionContainer)

	// Create Search TestData-button
	var searchTestDataButton *widget.Button
	searchTestDataButton = widget.NewButton("Search for TestDataPoints", func() {

		//var allTestDataPointRowsUuid []TestDataPointRowUuidType
		var searchResult []TestDataPointRowUuidType

		var tempTestDataModelMap map[TestDataDomainUuidType]*TestDataDomainModelStruct
		var tempTestDataDomainModel TestDataDomainModelStruct
		var tempTestDataAreaMap map[TestDataAreaUuidType]*TestDataAreaStruct
		var tempTestDataArea TestDataAreaStruct
		var tempTestDataValuesForRowMap map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
		var tempTestDataPointValueSlice []*TestDataPointValueStruct

		tempTestDataModelMap = *testDataModel.TestDataModelMap
		tempTestDataDomainModel = *tempTestDataModelMap[testDataDomainUuid]
		tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
		tempTestDataArea = *tempTestDataAreaMap[testDataAreaUuid]
		tempTestDataValuesForRowMap = *tempTestDataArea.TestDataValuesForRowMap

		var allTestDataPointRowsUuid []TestDataPointRowUuidType
		var tempTestDataValueNameToRowUuidMap map[TestDataValueNameType][]TestDataPointRowUuidType

		tempTestDataValueNameToRowUuidMap = make(map[TestDataValueNameType][]TestDataPointRowUuidType)

		// Loop all TestData and extract all rows
		for tempTestDataPointRowUuid, _ := range tempTestDataValuesForRowMap {
			allTestDataPointRowsUuid = append(allTestDataPointRowsUuid, tempTestDataPointRowUuid)
		}

		searchResult = allTestDataPointRowsUuid

		// Loop all Columns and extract selected checkboxes in the CheckGroups
		for _, testDataValueSelection := range testDataValueSelections {

			// Extract the Selected CheckBoxes
			var selectedCheckBoxes []string
			selectedCheckBoxes = testDataValueSelection.testDataCheckGroup.Selected

			// Extract 'TestDataPointRowUuid' for the Selected CheckBox-value-rows
			var testDataPointRowUuidMap map[TestDataValueType]*[]TestDataPointRowUuidType
			testDataPointRowUuidMap = *testDataValueSelection.TestDataPointValueRowUuidMap

			var testDataPointRowsUuid []TestDataPointRowUuidType

			for _, selectedCheckBox := range selectedCheckBoxes {
				tempTestDataPointRowsUuid, _ := testDataPointRowUuidMap[TestDataValueType(selectedCheckBox)]

				testDataPointRowsUuid = append(testDataPointRowsUuid, *tempTestDataPointRowsUuid...)

			}

			// Intersect with full TestDataSet to minimize the rows
			if len(testDataPointRowsUuid) != 0 {

				searchResult = testDataPointIntersectionOfTwoSlices(searchResult, testDataPointRowsUuid)

			}
		}

		// Convert into all 'TestDataValueName' in []string to be used in Available TestDataPoints-list
		filteredTestDataPoints = nil
		var tempTestDataValueName string
		for _, testDataPointRowUuid := range searchResult {

			tempTestDataPointValueSlice = *tempTestDataValuesForRowMap[testDataPointRowUuid]

			tempTestDataValueName = string(tempTestDataPointValueSlice[0].TestDataValueName)

			tempTestDataPointRowUuidSliceInMap, _ = tempTestDataValueNameToRowUuidMap[TestDataValueNameType(tempTestDataValueName)]
			tempTestDataPointRowUuidSliceInMap = append(tempTestDataPointRowUuidSliceInMap, testDataPointRowUuid)
			tempTestDataValueNameToRowUuidMap[TestDataValueNameType(tempTestDataValueName)] = tempTestDataPointRowUuidSliceInMap
		}

		for tempTestDataValueNameInMap, tempTestDataPointRowUuidSliceInMap := range tempTestDataValueNameToRowUuidMap {

			filteredTestDataPoints = append(filteredTestDataPoints, fmt.Sprintf("%s [%d]", string(tempTestDataValueNameInMap), len(tempTestDataPointRowUuidSliceInMap)))
		}

		// Create the list that holds all points that are available to chose from
		// Create the list that holds all points that are chosen
		allPointsAvailable = nil
		var existInSelectedPoints bool
		for _, point := range filteredTestDataPoints {

			// Add it to the list of available points, if it doesn't exist in the Selected-List
			if len(allSelectedPoints) == 0 {
				allPointsAvailable = append(allPointsAvailable, point)
			} else {
				for _, selectedPoint := range allSelectedPoints {

					if selectedPoint == point {
						existInSelectedPoints = true
						break
					}
					if existInSelectedPoints == false {
						allPointsAvailable = append(allPointsAvailable, point)
					} else {
						existInSelectedPoints = false
					}
				}
			}
		}

		// Custom sort: we sort by splitting each string into parts and comparing the parts
		sort.Slice(allPointsAvailable, func(i, j int) bool {
			// Split both strings by '/'
			partsI := strings.Split(allPointsAvailable[i], "/")
			partsJ := strings.Split(allPointsAvailable[j], "/")

			// Compare each part; the first non-equal part determines the order
			for k := 0; k < len(partsI) && k < len(partsJ); k++ {
				if partsI[k] != partsJ[k] {
					return partsI[k] < partsJ[k]
				}
			}

			// If all compared parts are equal, but one slice is shorter, it comes first
			return len(partsI) < len(partsJ)
		})

		// Refresh the List-widget
		allAvailablePointsList.Refresh()

	})

	// Create Clear checkboxes-button
	var clearTestDataFilterCheckBoxesButton *widget.Button
	clearTestDataFilterCheckBoxesButton = widget.NewButton("Clear checkboxes", func() {

		var selected []string

		// Loop all Columns and clear all checkboxes in the CheckGroups
		for _, testDataValueSelection := range testDataValueSelections {

			testDataValueSelection.testDataCheckGroup.SetSelected(selected)

		}

	})

	// Create the container for Search- and Clear- buttons
	var searchAndClearButtonsContainer *fyne.Container
	searchAndClearButtonsContainer = container.NewHBox(searchTestDataButton, clearTestDataFilterCheckBoxesButton)

	// Create the list that holds all points that are chosen
	for _, point := range selectedPoints {

		allSelectedPoints = append(allSelectedPoints, string(point.testDataPointName))

	}

	// Create and configure the list of all TestDataPoints
	allAvailablePointsList = widget.NewList(
		func() int { return len(allPointsAvailable) },
		func() fyne.CanvasObject {

			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(allPointsAvailable[id])
		},
	)

	// Create and configure the list of selected TestDataPoints
	selectedPointsList := widget.NewList(
		func() int { return len(allSelectedPoints) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(allSelectedPoints[id])
		},
	)

	// Functionality to add a point from 'allPointsAvailable' to 'allSelectedPoints'

	allAvailablePointsList.OnSelected = func(id widget.ListItemID) {

		re := regexp.MustCompile(`\[(\d+)\]`)

		// FindStringSubmatch returns an array of matches where the first element is the full match,
		// and the subsequent ones are the captured groups.
		matches := re.FindStringSubmatch(allPointsAvailable[id])
		if len(matches) > 1 { // matches[0] is the full match, matches[1] would be the first captured group
			fmt.Println("Extracted number:", matches[1])
		} else {
			fmt.Println("No number found")
		}

		re = regexp.MustCompile(`^(.*?)\[\d+\]`)

		// FindStringSubmatch returns an array of matches where the first element is the full match,
		// and the subsequent ones are the captured groups.
		matches = re.FindStringSubmatch(allPointsAvailable[id])
		var clickedDataPointName string

		if len(matches) > 1 { // matches[0] is the full match, matches[1] would be the first captured group
			clickedDataPointName = strings.Trim(matches[1], " ")
			fmt.Println(fmt.Sprintf("Extracted text to the left: '%s'", clickedDataPointName))
		} else {
			fmt.Println("No matching text found")
		}

		var tableData [][]string
		tableData = buildTableData(clickedDataPointName, testDataModel)

		showTable(newOrEditTestDataPointGroupWindow, tableData)

		allSelectedPoints = append(allSelectedPoints, allPointsAvailable[id])
		allPointsAvailable = append(allPointsAvailable[:id], allPointsAvailable[id+1:]...)

		allAvailablePointsList.UnselectAll()

		allAvailablePointsList.Refresh()
		selectedPointsList.Refresh()

	}

	// Functionality to remove a point from 'selectedPoints'
	selectedPointsList.OnSelected = func(id widget.ListItemID) {
		//allPointsAvailable = append(allPointsAvailable, allSelectedPoints[id])
		allSelectedPoints = append(allSelectedPoints[:id], allSelectedPoints[id+1:]...)

		selectedPointsList.UnselectAll()

		allAvailablePointsList.Refresh()
		selectedPointsList.Refresh()
	}

	// the Entry for the name of the TestDataPointsGroup
	nameEntry := widget.NewEntry()
	nameStatusLabel := widget.NewLabel(groupNameIsUnique)

	// Buttons for Save and Cancel actions
	saveButton = widget.NewButton("Save", func() {
		// Logic to add new newTestDataPointNameMa
		var newTestDataPointNameMap testDataPointNameMapType
		newTestDataPointNameMap = make(testDataPointNameMapType)

		// Loop all points and add them the 'newTestDataPointNameMap'
		for _, selectedPoint := range allSelectedPoints {

			var testDataPoint testDataPointStruct
			testDataPoint = testDataPointStruct{
				testDataPointUuid:            testDataPointUuidType(selectedPoint),
				testDataPointName:            testDataPointNameType(selectedPoint),
				testDataPointNameDescription: testDataPointNameDescriptionType(selectedPoint),
				testDatapointValue:           testDatapointValueType(selectedPoint),
			}
			newTestDataPointNameMap[testDataPointUuidType(selectedPoint)] = testDataPoint
		}

		// When GroupName is changed and the Group is in 'Edit'-mode the remove the old Group
		if isNew == false && nameEntry.Text != string(incomingGroupName) {
			delete(newOrEditedChosenTestDataPointsThisGroupMap, testDataPointGroupNameType(incomingGroupName))
		}

		// Add the TestDataPoints to the GroupName used
		newOrEditedChosenTestDataPointsThisGroupMap[testDataPointGroupNameType(nameEntry.Text)] = &newTestDataPointNameMap
		newOrEditedChosenTestDataPointsThisGroupMapPtr = &newOrEditedChosenTestDataPointsThisGroupMap

		shouldUpdateMainWindow = responseChannelStruct{
			shouldBeUpdated:        true,
			testDataPointGroupName: testDataPointGroupNameType(nameEntry.Text),
		}

		newOrEditTestDataPointGroupWindow.Close()
	})
	cancelButton := widget.NewButton("Cancel", func() {
		newOrEditTestDataPointGroupWindow.Close()
	})

	// Function for checking that name is unique
	nameEntry.OnChanged = func(entryValue string) {

		// Trigger State change control for Save Button and GroupName Entry
		setStateForSaveButtonAndGroupNameTextEntry(entryValue, nameStatusLabel, saveButton, isNew, incomingGroupName)
	}

	// Set placeholder text for GroupName Entry
	nameEntry.SetPlaceHolder("<enter group name here>")

	// Extract Group Name
	if incomingGroupName != "" {

		// Set Group Name in entry
		nameEntry.SetText(string(incomingGroupName))
	} else {
		// Set Group Name in entry
		//nameEntry.SetText("<new group>")
	}

	// Trigger State change control for Save Button and GroupName Entry
	setStateForSaveButtonAndGroupNameTextEntry(nameEntry.Text, nameStatusLabel, saveButton, isNew, incomingGroupName)

	// Layout configuration for the new/edit window
	// Create the UpperAndLowerSplitContainer
	var upperAndLowerSplitContainer *container.Split
	var listsSplitContainer *container.Split
	var buttonsContainer *fyne.Container
	var entryContainer *fyne.Container
	var upperSplitContainer *fyne.Container
	//var lowerSplitContainer *fyne.Container

	var tempTestDataPointsLabel *widget.Label
	tempTestDataPointsLabel = widget.NewLabel("TestDataPoints based on filter")
	tempTestDataPointsLabel.TextStyle.Bold = true

	var lowerLeftSideContainer *fyne.Container
	lowerLeftSideContainer = container.NewBorder(tempTestDataPointsLabel, nil, nil, nil, allAvailablePointsList)

	buttonsContainer = container.NewHBox(saveButton, cancelButton)
	entryContainer = container.NewBorder(nil, nil, nil, nameStatusLabel, nameEntry)

	var tempTestGroupLabel *widget.Label
	tempTestGroupLabel = widget.NewLabel("TestDataGroup and its TestDataPoints")
	tempTestGroupLabel.TextStyle.Bold = true

	var lowerRightSideContainer *fyne.Container
	lowerRightSideContainer = container.NewBorder(container.NewVBox(tempTestGroupLabel, entryContainer, buttonsContainer), nil, nil, nil, selectedPointsList)
	listsSplitContainer = container.NewHSplit(lowerLeftSideContainer, lowerRightSideContainer)

	upperSplitContainer = container.NewBorder(nil, searchAndClearButtonsContainer, nil, nil, testDataSelectionsContainer)
	//lowerSplitContainer = container.NewVBox(entryContainer, buttonsContainer, searchAndClearButtonsContainer, listsSplitContainer)

	upperAndLowerSplitContainer = container.NewVSplit(upperSplitContainer, listsSplitContainer)

	//windowContent = container.NewBorder(container.NewVBox(entryContainer, buttonsContainer, testDataSelectionsContainer, searchAndClearButtonsContainer), nil, nil, nil, listsSplitContainer)

	newOrEditTestDataPointGroupWindow.SetContent(upperAndLowerSplitContainer)
	newOrEditTestDataPointGroupWindow.Show()

}

// Set the State for Save Button and the GroupName Entry
func setStateForSaveButtonAndGroupNameTextEntry(
	entryValue string,
	nameStatusLabel *widget.Label,
	saveButton *widget.Button,
	isNew bool,
	incomingGroupName testDataPointGroupNameType) {

	// Handle when first or last character in Group Name is a 'space'
	if len(entryValue) > 0 && (entryValue[:1] == " " || entryValue[len(entryValue)-1:] == " ") {
		nameStatusLabel.SetText(groupNameCanNotStartOrEndWithSpace)
		saveButton.Disable()

		return
	}

	// Handle when this is a new Group and it is the first control
	if isNew == true && len(entryValue) == 0 {
		nameStatusLabel.SetText(groupNameIsEmpty)
		saveButton.Disable()

		return
	}

	// Handle when this there are no existing Groups in the map
	if len(chosenTestDataPointsPerGroupMap) == 0 {
		if isNew == true && len(entryValue) == 0 {
			nameStatusLabel.SetText(groupNameIsEmpty)
			saveButton.Disable()

			return
		} else {
			nameStatusLabel.SetText(groupNameIsUnique)
			saveButton.Enable()
		}
	}

	for existingTestDataPointGroupName, _ := range chosenTestDataPointsPerGroupMap {

		if len(entryValue) == 0 {
			nameStatusLabel.SetText(groupNameIsEmpty)
			saveButton.Disable()

		} else if entryValue == string(existingTestDataPointGroupName) &&
			entryValue != string(incomingGroupName) {

			nameStatusLabel.SetText(groupNameIsNotUnique)
			saveButton.Disable()

		} else {

			nameStatusLabel.SetText(groupNameIsUnique)
			saveButton.Enable()
		}
	}
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

func buildTableData(
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
	var tempTestDataValuesForRowNameMap map[TestDataValueNameType]*[]*map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	var tempTestDataValuesForRowUuidMapBaseOnNameSlice []*map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct

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

	var tempTestDataPointRowNameToSearchFor string
	tempTestDataPointRowNameToSearchFor = tempTestDataPointRowName[len(matches[0]+"/"):]

	tempTestDataValuesForRowUuidMapBaseOnNameSlice = *tempTestDataValuesForRowNameMap[TestDataValueNameType(tempTestDataPointRowNameToSearchFor)]

	// Loop the slice to extract the map for one row of data
	for _, tempTestDataValueMapForRow := range tempTestDataValuesForRowUuidMapBaseOnNameSlice {

		var rowSlice []string

		// Loop the map to get all values for the row
		for _, tempTestDataPointSlice := range *tempTestDataValueMapForRow {
			for _, tempTestDataPoint := range *tempTestDataPointSlice {
				rowSlice = append(rowSlice, string(tempTestDataPoint.TestDataValue))
			}

		}

		tableData = append(tableData, rowSlice)

	}

	return tableData
}

// showTable creates and shows a table for the selected node with data
func showTable(w fyne.Window, data [][]string) {
	table := widget.NewTable(
		func() (int, int) { return len(data), len(data[0]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Create a label for each cell
		},
		func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(data[cellID.Row][cellID.Col]) // Set text based on data
		},
	)

	// Calculate and set column widths based on content
	setColumnWidths(table, data)

	// Set minimum size for the table to ensure it's larger
	table.Resize(fyne.NewSize(400, 300)) // Set the minimum size to 400x300 pixels

	// Use a scroll container to make the table scrollable in case it has more data
	scrollContainer := container.NewScroll(table)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300)) // Ensure the scroll container is also adequately sized

	// Show table in a pop-up and ensure the pop-up is appropriately sized
	popup := widget.NewModalPopUp(scrollContainer, w.Canvas())
	popup.Resize(fyne.NewSize(450, 350)) // Resize the popup to be slightly larger than the table and container
	popup.Show()
}
func setColumnWidths(table *widget.Table, data [][]string) {
	maxWidths := make([]float32, len(data[0]))
	for col := range maxWidths {
		for row := range data {
			width := fyne.MeasureText(data[row][col], theme.TextSize(), fyne.TextStyle{}).Width
			if width > maxWidths[col] {
				maxWidths[col] = width
			}
		}
		// Add some padding to the maximum width found
		maxWidths[col] += theme.Padding() * 4
	}

	for col, width := range maxWidths {
		table.SetColumnWidth(col, width)
	}
}
