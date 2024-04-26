package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// The slices for Groups ans TestDataPoints for a Group
var testDataPointGroups = []testDataPointGroupNameType{} // Define testDataPointGroups
var testDataPointsForAGroup = []testDataPointNameType{}  // Define testDataPointGroups

// The List-items for Groups ans TestDataPoints for a Group
var testDataPointGroupsList *widget.List
var testDataPointsForAGroupList *widget.List

// The map holding all data for Groups and TestDataPoints for a group
var chosenTestDataPointsPerGroupMap map[testDataPointGroupNameType]*testDataPointNameMapType

// Types used for data structures
type testDataPointNameType string
type testDataPointGroupNameType string
type testDataPointNameMapType map[testDataPointNameType]testDataPointNameType

const (
	groupNameIsUnique    string = "Group name is OK!"
	groupNameIsNotUnique string = "Group name already exists!"
	groupNameIsEmpty     string = "Group name can't be empty"
)

type responseChannelStruct struct {
	shouldBeUpdated        bool
	testDataPointGroupName testDataPointGroupNameType
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TestData Management")
	myWindow.Resize(fyne.NewSize(600, 500))

	// Initiate 'chosenTestDataPointsPerGroupMap'
	chosenTestDataPointsPerGroupMap = make(map[testDataPointGroupNameType]*testDataPointNameMapType)

	// Create List UI for 'testDataPointGroups'
	var selectedIndexForGroups int = -1

	testDataPointGroupsList = widget.NewList(
		func() int { return len(testDataPointGroups) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataPointGroups[id]))
		},
	)

	testDataPointGroupsList.OnSelected = func(id widget.ListItemID) {
		selectedIndexForGroups = id

		// Update List for  'testDataPointsForAGroup'
		updateTestDataPointsForAGroupList(testDataPointGroups[id])
	}

	// Create List UI for 'testDataPointsForAGroup'
	testDataPointsForAGroupList = widget.NewList(
		func() int { return len(testDataPointsForAGroup) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataPointsForAGroup[id]))
		},
	)

	// Create Split Container used for 'testDataPointGroups' and 'testDataPointsForAGroup'
	var testDataGroupsAndPointsContainer *container.Split
	testDataGroupsAndPointsContainer = container.NewHSplit(testDataPointGroupsList, testDataPointsForAGroupList)

	var responseChannel chan responseChannelStruct
	responseChannel = make(chan responseChannelStruct)

	// The structure holding Group and TestDataPoints together
	//var newOrEditedChosenTestDataPointsPerGroupMap map[testDataPointGroupNameType][]testDataPointNameType
	//newOrEditedChosenTestDataPointsPerGroupMap = make(map[testDataPointGroupNameType][]testDataPointNameType)

	// Crete the 'New'-button for creating a new Group for TestDataPoints
	newButton := widget.NewButton("New", func() {
		myWindow.Hide()
		showNewOrEditGroupWindow(
			myApp,
			myWindow,
			true,
			&responseChannel,
			"",
			&chosenTestDataPointsPerGroupMap)
	})

	// Crete the 'Edit'-button for editing an existing Group for TestDataPoints
	editButton := widget.NewButton("Edit", func() {
		if selectedIndexForGroups == -1 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}
		myWindow.Hide()
		showNewOrEditGroupWindow(
			myApp,
			myWindow,
			false,
			&responseChannel,
			testDataPointGroups[selectedIndexForGroups],
			&chosenTestDataPointsPerGroupMap)
	})

	// Crete the 'Delete'-button for deleting an existing Group for TestDataPoints
	deleteButton := widget.NewButton("Delete", func() {
		if selectedIndexForGroups == -1 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}
		dialog.ShowConfirm("Confirm", "Are you sure?", func(confirm bool) {
			if confirm {

				// Get the GroupName from the List to be deleted
				var groupNameToDelete testDataPointGroupNameType
				groupNameToDelete = testDataPointGroups[selectedIndexForGroups]

				// Delete the group
				delete(chosenTestDataPointsPerGroupMap, groupNameToDelete)

				// Rebuild the TestDataPointGroup-list
				testDataPointGroups = nil
				for testDataPointsGroupName, _ := range chosenTestDataPointsPerGroupMap {

					testDataPointGroups = append(testDataPointGroups, testDataPointsGroupName)
				}

				testDataPointGroupsList.Refresh()
				testDataPointGroupsList.UnselectAll()

				// Clear the TestDataPointsList
				testDataPointsForAGroup = nil
				testDataPointsForAGroupList.Refresh()
			}
		}, myWindow)
	})

	// Create the container for handling TestDataGroups
	var buttonsContainer *fyne.Container
	buttonsContainer = container.NewHBox(newButton, editButton, deleteButton)

	// Create the container that holds all UI components used for Groups and Points
	myContainer := container.NewBorder(nil, buttonsContainer, nil, nil, testDataGroupsAndPointsContainer)

	myWindow.SetContent(myContainer)

	// Function that updates new or changes lists in the UI
	go func() {

		var shouldListBeUpdated responseChannelStruct
		var groupNameIndex int
		var groupNameIndexToSelect int

		for {

			shouldListBeUpdated = <-responseChannel

			// Update the List in main window if true as response
			if shouldListBeUpdated.shouldBeUpdated == true {

				// Clear slice and variables used
				testDataPointGroups = nil
				groupNameIndex = 0
				groupNameIndexToSelect = 0

				for testDataPointsGroupName, _ := range chosenTestDataPointsPerGroupMap {

					testDataPointGroups = append(testDataPointGroups, testDataPointsGroupName)

					if testDataPointsGroupName == shouldListBeUpdated.testDataPointGroupName {

						groupNameIndexToSelect = groupNameIndex

					}

					groupNameIndex = groupNameIndex + 1

				}
				testDataPointGroupsList.Refresh()
				testDataPointGroupsList.UnselectAll()
				testDataPointGroupsList.Select(groupNameIndexToSelect)
				selectedIndexForGroups = groupNameIndexToSelect
			}
		}
	}()

	myWindow.ShowAndRun()

}

// Updates the list show the TestDataPoints for a specific Group in main window
func updateTestDataPointsForAGroupList(testDataPointGroupName testDataPointGroupNameType) {

	// Clear the slice that holds all TestDataPoints
	testDataPointsForAGroup = nil

	// Extract the map with the TestDataPoints
	var tempTestDataPointNameMap testDataPointNameMapType
	tempTestDataPointNameMap = *chosenTestDataPointsPerGroupMap[testDataPointGroupName]

	// Refill the slice with all TestDataPoints
	for testDataPoint, _ := range tempTestDataPointNameMap {
		testDataPointsForAGroup = append(testDataPointsForAGroup, testDataPoint)

	}

	// Refresh the List in the UI
	testDataPointsForAGroupList.Refresh()
}

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
