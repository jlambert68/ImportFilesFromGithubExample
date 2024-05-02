package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
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
	testDataModelMap *map[TestDataDomainUuidType]*TestDataDomainModelStruct) {

	parent.Hide()

	var filteredTestDataPoints []string

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
		testDataSelectionLabel *widget.Label
		testDataCheckGroup     *widget.CheckGroup
	}
	var testDataValueSelections []*testDataValueSelectionStruct
	var testDataValuesSelectionContainer *fyne.Container

	// Create label for Domains
	domainsLabel = widget.NewLabel(testDataDomainLabelText)
	domainsLabel.TextStyle.Bold = true
	testAreasLabel = widget.NewLabel(testDataTestAreaLabelText)
	testAreasLabel.TextStyle.Bold = true

	// Extract TestData on Domain-level
	for _, tempTestDataDomainModel := range *testDataModelMap {
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

				// Loop 'testDataColumnsMetaDataToBeSorted' for Columns to present
				for _, testDataColumnsMetaData := range testDataColumnsMetaDataToBeSorted {

					// Check if column should be used for filtering TestData
					if testDataColumnsMetaData.ShouldColumnBeUsedForFindingTestData == true {

						var checkGroupOptions []string
						var tempTestDataColumnContainer *fyne.Container

						// Set Label
						var newColumnFilterLabel *widget.Label
						newColumnFilterLabel = widget.NewLabel(string(testDataColumnsMetaData.TestDataColumnUIName))
						newColumnFilterLabel.TextStyle.Bold = true

						var testDataValueSelection *testDataValueSelectionStruct
						testDataValueSelection = &testDataValueSelectionStruct{
							testDataSelectionLabel: newColumnFilterLabel,
							testDataCheckGroup:     nil,
						}

						// Extract the Map with the values
						var uniqueTestDataValuesForColumnMapPtr *map[TestDataValueType]TestDataValueType
						UniqueTestDataValuesForColumnMap := *testDataArea.UniqueTestDataValuesForColumnMap

						uniqueTestDataValuesForColumnMapPtr = UniqueTestDataValuesForColumnMap[testDataColumnsMetaData.TestDataColumnUuid]

						// Loop Values in Column and create Checkboxes
						for _, uniqueTestDataValue := range *uniqueTestDataValuesForColumnMapPtr {

							// Add value to slice for CheckBox-labels
							checkGroupOptions = append(checkGroupOptions, string(uniqueTestDataValue))

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

	//

	// Create the separate TestData-selection-containers
	testDomainContainer = container.NewVBox(domainsLabel, domainsSelect)
	testAreasContainer = container.NewVBox(testAreasLabel, testAreaSelect)

	// Create the main TestData-selection-container
	testDataSelectionsContainer = container.NewHBox(testDomainContainer, testAreasContainer, testDataValuesSelectionContainer)

	// Create Search TestData-button
	var searchTestDataButton *widget.Button
	searchTestDataButton = widget.NewButton("Search for TestDataPoints", func() {

		//var searchResult []string

	})

	// Create Clear checkboxes-button
	var clearTestDataFilterCheckBoxesButton *widget.Button
	clearTestDataFilterCheckBoxesButton = widget.NewButton("Clear checkboxes", func() {

		selected := []string{}

		// Loop all Columns and clear all checkboxes in the CheckGroups
		for _, testDataValueSelection := range testDataValueSelections {

			testDataValueSelection.testDataCheckGroup.SetSelected(selected)

		}

	})

	// Create the container for Search- and Clear- buttons
	var searchAndClearButtonsContainer *fyne.Container
	searchAndClearButtonsContainer = container.NewHBox(searchTestDataButton, clearTestDataFilterCheckBoxesButton)

	// Sample data for demonstration
	filteredTestDataPoints = []string{} // {"Point_1", "Point_2", "Point_3", "Point_4", "Point_5", "Point_6", "Point_7", "Point_8", "Point_9", "Point_10"}
	var allPointsAvailable []string
	var allSelectedPoints []string

	var existsInMap bool

	// If existing groupToEdit then extract points from it otherwise create an empty selected points slice
	var selectedPointsPtr *testDataPointNameMapType
	var selectedPoints testDataPointNameMapType

	if isNew == false {

		selectedPointsPtr = newOrEditedChosenTestDataPointsThisGroupMap[incomingGroupName]
		selectedPoints = *selectedPointsPtr

	} else {

	}

	// Create the list that holds all points that are available to chose from
	// Create the list that holds all points that are chosen
	for _, point := range filteredTestDataPoints {

		// Check if the point exists in the map with chosen points
		_, existsInMap = selectedPoints[testDataPointUuidType(point)]
		if existsInMap == false {
			// Add it to the list of available points
			allPointsAvailable = append(allPointsAvailable, point)

		} else {
			allSelectedPoints = append(allSelectedPoints, point)
		}

	}

	// Create and configure the list of all TestDataPoints
	allAvailablePointsList := widget.NewList(
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
		allSelectedPoints = append(allSelectedPoints, allPointsAvailable[id])
		allPointsAvailable = append(allPointsAvailable[:id], allPointsAvailable[id+1:]...)

		allAvailablePointsList.UnselectAll()

		allAvailablePointsList.Refresh()
		selectedPointsList.Refresh()

	}

	// Functionality to remove a point from 'selectedPoints'
	selectedPointsList.OnSelected = func(id widget.ListItemID) {
		allPointsAvailable = append(allPointsAvailable, allSelectedPoints[id])
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
