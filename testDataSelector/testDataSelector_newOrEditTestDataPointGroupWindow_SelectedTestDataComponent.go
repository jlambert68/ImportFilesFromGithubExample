package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Create and configure the list-component of selected TestDataPoints
func generateSelectedPointsListUIComponent(
	selectedPointsList *widget.List,
	allPointsAvailable *[]dataPointTypeForGroupsStruct,
	allSelectedPoints *[]dataPointTypeForGroupsStruct,
	newOrEditTestDataPointGroupWindowPtr *fyne.Window,
	lowerRightSideContainer *fyne.Container,
	incomingGroupName testDataPointGroupNameType,
	isNew bool,
	newOrEditedChosenTestDataPointsThisGroupMapPtr *map[testDataPointGroupNameType]*testDataPointNameMapType) {

	var newOrEditTestDataPointGroupWindow fyne.Window
	newOrEditTestDataPointGroupWindow = *newOrEditTestDataPointGroupWindowPtr

	//var selectedPointTypes testDataPointNameMapType

	var newOrEditedChosenTestDataPointsThisGroupMap map[testDataPointGroupNameType]*testDataPointNameMapType
	newOrEditedChosenTestDataPointsThisGroupMap = *newOrEditedChosenTestDataPointsThisGroupMapPtr

	// WHen a Group already exist then move over the TestDataPoints to the SelectedPointsList
	if isNew == false {

		// If existing groupToEdit then extract points from it otherwise create an empty selected points slice
		var selectedPointsPtr *testDataPointNameMapType
		var dataPointsForGroup []*dataPointTypeForGroupsStruct

		// Extract the TestDataPoints for the Group
		selectedPointsPtr = newOrEditedChosenTestDataPointsThisGroupMap[incomingGroupName]

		// Extract the indTestDataRows
		for _, testDataPointsForGroupPtr := range *selectedPointsPtr {
			dataPointsForGroup = *testDataPointsForGroupPtr

			// Loop the TestDataPointNames
			for _, tempDataPointForGroup := range dataPointsForGroup {

				*allSelectedPoints = append(*allSelectedPoints, *tempDataPointForGroup)

			}

		}
	}

	// Create and configure the list-component of selected TestDataPoints
	selectedPointsList = widget.NewList(
		func() int { return len(*allSelectedPoints) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {

			var localCopyAllSelectedPoints []dataPointTypeForGroupsStruct
			localCopyAllSelectedPoints = *allSelectedPoints

			obj.(*widget.Label).SetText(fmt.Sprintf("%s [%d]", string(localCopyAllSelectedPoints[id].testDataPointName), len(localCopyAllSelectedPoints[id].testDataPointUuidMap)))
		},
	)

	/*
		// Functionality to remove a pointTypeSlicePtr from 'selectedPointTypes'
		selectedPointsList.OnSelected = func(id widget.ListItemID) {
			//allPointsAvailable = append(allPointsAvailable, allSelectedPoints[id])
			allSelectedPoints = append(allSelectedPoints[:id], allSelectedPoints[id+1:]...)

			selectedPointsList.UnselectAll()

			allAvailablePointsList.Refresh()
			selectedPointsList.Refresh()
		}

	*/

	// the Entry for the name of the TestDataPointsGroup
	nameEntry := widget.NewEntry()
	nameStatusLabel := widget.NewLabel(groupNameIsUnique)

	// Buttons for Save and Cancel actions
	var saveButton *widget.Button
	saveButton = widget.NewButton("Save", func() {
		// Logic to add new newTestDataPointNameMap
		var newTestDataPointNameMap testDataPointNameMapType
		newTestDataPointNameMap = make(testDataPointNameMapType)
		var testDataPointRowSlice []*testDataPointRowStruct
		//var tempSelectedPointUuid TestDataPointRowUuidType

		// Loop all points and add them the 'newTestDataPointNameMap'
		for _, selectedPoint := range *allSelectedPoints {

			for _, selectedPointUuid := range selectedPoint.testDataPointUuidMap {

				// Create a TestDataPoint to be stored in a TestDataGroup
				var testDataPoint testDataPointRowStruct
				testDataPoint = testDataPointRowStruct{
					testDataDomainUuid:           selectedPoint.testDataDomainUuid,
					testDataDomainName:           selectedPoint.testDataDomainName,
					testDataAreaUuid:             selectedPoint.testDataAreaUuid,
					testDataAreaName:             selectedPoint.testDataAreaName,
					testDataPointUuid:            selectedPointUuid,
					testDataPointName:            selectedPoint.testDataPointName,
					testDataPointNameDescription: "< Not set, yet >",
					testDatapointValue:           "< Not set, yet >",
				}

				testDataPointRowSlice = append(testDataPointRowSlice, &testDataPoint)

			}

			newTestDataPointNameMap[testDataPointRowSlice[0].testDataPointName] = &testDataPointRowSlice

		}

		// When GroupName is changed and the Group is in 'Edit'-mode the remove the old Group
		if isNew == false && nameEntry.Text != string(incomingGroupName) {
			delete(newOrEditedChosenTestDataPointsThisGroupMap, testDataPointGroupNameType(incomingGroupName))
		}

		// Add the TestDataPoints to the GroupName used
		newOrEditedChosenTestDataPointsThisGroupMap[testDataPointGroupNameType(nameEntry.Text)] = &newTestDataPointNameMap
		newOrEditedChosenTestDataPointsThisGroupMapPtr = &newOrEditedChosenTestDataPointsThisGroupMap

		//shouldUpdateMainWindow = responseChannelStruct{
		//	shouldBeUpdated:        true,
		//	testDataPointGroupName: testDataPointGroupNameType(nameEntry.Text),
		//}

		newOrEditTestDataPointGroupWindow.Close()
	})

	var cancelButton *widget.Button
	cancelButton = widget.NewButton("Cancel", func() {
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

	var buttonsContainer *fyne.Container
	buttonsContainer = container.NewHBox(saveButton, cancelButton)

	var entryContainer *fyne.Container
	entryContainer = container.NewBorder(nil, nil, nil, nameStatusLabel, nameEntry)

	var tempTestGroupLabel *widget.Label
	tempTestGroupLabel = widget.NewLabel("TestDataGroup and its TestDataPoints")
	tempTestGroupLabel.TextStyle.Bold = true

	lowerRightSideContainer = container.NewBorder(container.NewVBox(tempTestGroupLabel, entryContainer, buttonsContainer), nil, nil, nil, selectedPointsList)

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
