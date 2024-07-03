package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Create and configure the list-component of selected TestDataPoints
func generateSelectedPointsListUIComponent(
	newOrEditTestDataPointGroupWindowPtr *fyne.Window,
	incomingGroupName testDataEngine.TestDataPointGroupNameType,
	isNew bool,
	newOrEditedChosenTestDataPointsThisGroupMapPtr *map[testDataEngine.TestDataPointGroupNameType]*testDataEngine.TestDataPointNameMapType) {

	var newOrEditTestDataPointGroupWindow fyne.Window
	newOrEditTestDataPointGroupWindow = *newOrEditTestDataPointGroupWindowPtr

	//var selectedPointTypes testDataPointNameMapType

	var newOrEditedChosenTestDataPointsThisGroupMap map[testDataEngine.TestDataPointGroupNameType]*testDataEngine.TestDataPointNameMapType
	newOrEditedChosenTestDataPointsThisGroupMap = *newOrEditedChosenTestDataPointsThisGroupMapPtr

	// WHen a Group already exist then move over the TestDataPoints to the SelectedPointsList
	if isNew == false {

		// If existing groupToEdit then extract points from it otherwise create an empty selected points slice
		var selectedPointsPtr *testDataEngine.TestDataPointNameMapType
		var dataPointsForGroup []*testDataEngine.DataPointTypeForGroupsStruct

		// Extract the TestDataPoints for the Group
		selectedPointsPtr = newOrEditedChosenTestDataPointsThisGroupMap[incomingGroupName]

		// Extract the indTestDataRows
		for _, testDataPointsForGroupPtr := range *selectedPointsPtr {
			dataPointsForGroup = *testDataPointsForGroupPtr

			// Loop the TestDataPointNames
			for _, tempDataPointForGroup := range dataPointsForGroup {

				testDataEngine.AllSelectedPoints = append(testDataEngine.AllSelectedPoints, *tempDataPointForGroup)

			}

		}
	}

	// Create and configure the list-component of selected TestDataPoints
	selectedPointsList = widget.NewList(
		func() int { return len(testDataEngine.AllSelectedPoints) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {

			obj.(*widget.Label).SetText(fmt.Sprintf(
				"%s [%d(%d)]",
				string(testDataEngine.AllSelectedPoints[id].TestDataPointName),
				len(testDataEngine.AllSelectedPoints[id].SelectedTestDataPointUuidMap),
				len(testDataEngine.AllSelectedPoints[id].AvailableTestDataPointUuidMap)+
					len(testDataEngine.AllSelectedPoints[id].SelectedTestDataPointUuidMap)))
		},
	)

	// the Entry for the name of the TestDataPointsGroup
	nameEntry := widget.NewEntry()
	nameStatusLabel := widget.NewLabel(testDataEngine.GroupNameIsUnique)

	// Buttons for Save and Cancel actions
	var saveButton *widget.Button
	saveButton = widget.NewButton("Save", func() {
		// Logic to add new newTestDataPointNameMap
		var newTestDataPointNameMap testDataEngine.TestDataPointNameMapType
		newTestDataPointNameMap = make(testDataEngine.TestDataPointNameMapType)

		// When GroupName is changed and the Group is in 'Edit'-mode the remove the old Group
		if isNew == false && nameEntry.Text != string(incomingGroupName) {
			delete(newOrEditedChosenTestDataPointsThisGroupMap, incomingGroupName)
		}

		for _, selectedPoint := range testDataEngine.AllSelectedPoints {

			var dataPointsForGroup []*testDataEngine.DataPointTypeForGroupsStruct
			dataPointsForGroup = append(dataPointsForGroup, &selectedPoint)
			newTestDataPointNameMap[selectedPoint.TestDataPointName] = &dataPointsForGroup

		}

		// Add the TestDataPoints to the GroupName used
		newOrEditedChosenTestDataPointsThisGroupMap[testDataEngine.TestDataPointGroupNameType(nameEntry.Text)] = &newTestDataPointNameMap
		newOrEditedChosenTestDataPointsThisGroupMapPtr = &newOrEditedChosenTestDataPointsThisGroupMap

		// Inform calling window that an update is done
		testDataEngine.ShouldUpdateMainWindow = testDataEngine.ResponseChannelStruct{
			ShouldBeUpdated:        true,
			TestDataPointGroupName: testDataEngine.TestDataPointGroupNameType(nameEntry.Text),
		}

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

	lowerRightSideContainer = container.NewBorder(
		container.NewVBox(tempTestGroupLabel, entryContainer, buttonsContainer),
		nil, nil, nil,
		selectedPointsList)

	// Create function call to 'setStateForSaveButtonAndGroupNameTextEntry' from outside
	testDataEngine.SetStateForSaveButtonAndGroupNameTextEntryExternalCall = func() {
		setStateForSaveButtonAndGroupNameTextEntry(nameEntry.Text, nameStatusLabel, saveButton, isNew, incomingGroupName)
	}

}

// Set the State for Save Button and the GroupName Entry
func setStateForSaveButtonAndGroupNameTextEntry(
	entryValue string,
	nameStatusLabel *widget.Label,
	saveButton *widget.Button,
	isNew bool,
	incomingGroupName testDataEngine.TestDataPointGroupNameType) {

	// Handle when the Selected List is empty
	if len(testDataEngine.AllSelectedPoints) == 0 {
		nameStatusLabel.SetText(testDataEngine.SelectedListIsEmpty)
		saveButton.Disable()

		return
	}

	// Handle when first or last character in Group Name is a 'space'
	if len(entryValue) > 0 && (entryValue[:1] == " " || entryValue[len(entryValue)-1:] == " ") {
		nameStatusLabel.SetText(testDataEngine.GroupNameCanNotStartOrEndWithSpace)
		saveButton.Disable()

		return
	}

	// Handle when this is a new Group and it is the first control
	if isNew == true && len(entryValue) == 0 {
		nameStatusLabel.SetText(testDataEngine.GroupNameIsEmpty)
		saveButton.Disable()

		return
	}

	// Handle when this there are no existing Groups in the map
	if len(testDataEngine.ChosenTestDataPointsPerGroupMap) == 0 {
		if isNew == true && len(entryValue) == 0 {
			nameStatusLabel.SetText(testDataEngine.GroupNameIsEmpty)
			saveButton.Disable()

			return
		} else {
			nameStatusLabel.SetText(testDataEngine.GroupNameIsUnique)
			saveButton.Enable()
		}
	}

	for existingTestDataPointGroupName, _ := range testDataEngine.ChosenTestDataPointsPerGroupMap {

		if len(entryValue) == 0 {
			nameStatusLabel.SetText(testDataEngine.GroupNameIsEmpty)
			saveButton.Disable()

		} else if entryValue == string(existingTestDataPointGroupName) &&
			entryValue != string(incomingGroupName) {

			nameStatusLabel.SetText(testDataEngine.GroupNameIsNotUnique)
			saveButton.Disable()

		} else {

			nameStatusLabel.SetText(testDataEngine.GroupNameIsUnique)
			saveButton.Enable()
		}
	}
}
