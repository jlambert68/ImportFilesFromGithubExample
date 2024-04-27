package testDataSelector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showNewOrEditGroupWindow(
	app fyne.App,
	parent fyne.Window,
	isNew bool,
	responseChannel *chan responseChannelStruct,
	incomingGroupName testDataPointGroupNameType,
	newOrEditedChosenTestDataPointsThisGroupMapPtr *map[testDataPointGroupNameType]*testDataPointNameMapType) {

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

	// Sample data for demonstration
	allPoints := []string{"Point_1", "Point_2", "Point_3", "Point_4", "Point_5", "Point_6", "Point_7", "Point_8", "Point_9", "Point_10"}
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
	for _, point := range allPoints {

		// Check if the point exists in the map with chosen points
		_, existsInMap = selectedPoints[testDataPointNameType(point)]
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

		allAvailablePointsList.Refresh()
		selectedPointsList.Refresh()
	}

	// Functionality to remove a point from 'selectedPoints'
	selectedPointsList.OnSelected = func(id widget.ListItemID) {
		allPointsAvailable = append(allPointsAvailable, allSelectedPoints[id])
		allSelectedPoints = append(allSelectedPoints[:id], allSelectedPoints[id+1:]...)

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
			newTestDataPointNameMap[testDataPointNameType(selectedPoint)] = testDataPointNameType(selectedPoint)
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
	listsSplitContainer := container.NewHSplit(allAvailablePointsList, selectedPointsList)
	buttonsContainer := container.NewHBox(saveButton, cancelButton)
	entryContainer := container.NewBorder(nil, nil, nil, nameStatusLabel, nameEntry)
	content := container.NewBorder(container.NewVBox(entryContainer, buttonsContainer), nil, nil, nil, listsSplitContainer)

	newOrEditTestDataPointGroupWindow.SetContent(content)
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
