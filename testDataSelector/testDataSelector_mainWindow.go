package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func MainTestDataSelector(
	app fyne.App,
	parent fyne.Window) {

	parent.Hide()

	myWindow := app.NewWindow("TestData Management")
	myWindow.Resize(fyne.NewSize(600, 500))

	// When this window closed then show parent and send response to parent window
	myWindow.SetOnClosed(func() {
		parent.Show()
	})

	// Initiate 'chosenTestDataPointsPerGroupMap'
	if chosenTestDataPointsPerGroupMap == nil {
		chosenTestDataPointsPerGroupMap = make(map[testDataPointGroupNameType]*testDataPointNameMapType)
	}

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
			app,
			myWindow,
			true,
			&responseChannel,
			"",
			&chosenTestDataPointsPerGroupMap)
	})

	// Crete the 'Edit'-button for editing an existing Group for TestDataPoints
	editButton := widget.NewButton("Edit", func() {
		if selectedIndexForGroups == -1 || len(testDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}
		myWindow.Hide()
		showNewOrEditGroupWindow(
			app,
			myWindow,
			false,
			&responseChannel,
			testDataPointGroups[selectedIndexForGroups],
			&chosenTestDataPointsPerGroupMap)
	})

	// Crete the 'Delete'-button for deleting an existing Group for TestDataPoints
	deleteButton := widget.NewButton("Delete", func() {
		if selectedIndexForGroups == -1 || len(testDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}

		dialog.ShowConfirm("Confirm to Delete", fmt.Sprintf("Are you sure that you what to delete TestDataPointGroup '%s'?", testDataPointGroups[selectedIndexForGroups]), func(confirm bool) {
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

				selectedIndexForGroups = -1

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

	myWindow.Show()

}
