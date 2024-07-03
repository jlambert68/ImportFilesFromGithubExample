package testDataSelector

import (
	"ImportFilesFromGithub/newOrEditTestDataPointGroupUI"
	"ImportFilesFromGithub/testDataEngine"
	"embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//go:embed testData/FenixRawTestdata_646rows_211220.csv
var embeddedFile embed.FS

const (
	testDataDomainUuid testDataEngine.TestDataDomainUuidType = "7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae"
	testDataDomainName testDataEngine.TestDataDomainNameType = "Sub Custody"
	testDataAreaUuid   testDataEngine.TestDataAreaUuidType   = "010cc994-a913-4041-96fe-a96d7e0c97e8"
	testDataAreaName   testDataEngine.TestDataAreaNameType   = "Main TestData Area"
)

func ImportTestData() {

	var headers []string
	var testDataRows [][]string

	headers, testDataRows = testDataEngine.ImportEmbeddedCsvTestDataFile(&embeddedFile, "testData/FenixRawTestdata_646rows_211220.csv", ';')

	testDataEngine.AddTestDataToTestDataModel(
		testDataDomainUuid, testDataDomainName,
		testDataAreaUuid, testDataAreaName,
		headers, testDataRows)

}

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
	if testDataEngine.ChosenTestDataPointsPerGroupMap == nil {
		testDataEngine.ChosenTestDataPointsPerGroupMap = make(map[testDataEngine.TestDataPointGroupNameType]*testDataEngine.TestDataPointNameMapType)
	}

	// Create List UI for 'testDataPointGroups'
	newOrEditTestDataPointGroupUI.TestDataPointGroupsList = widget.NewList(
		func() int { return len(testDataEngine.TestDataPointGroups) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataEngine.TestDataPointGroups[id]))
		},
	)

	newOrEditTestDataPointGroupUI.TestDataPointGroupsList.OnSelected = func(id widget.ListItemID) {
		newOrEditTestDataPointGroupUI.SelectedIndexForGroups = id

		// Update List for  'testDataPointsForAGroup'
		updateTestDataPointsForAGroupList(testDataEngine.TestDataPointGroups[id])

		// Select correct Group in Select-dropdown
		newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect.SetSelected(string(testDataEngine.TestDataPointGroups[id]))

		// UnSelect in DropDown- and List for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.UnselectAll()
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = -1

	}

	// Create function that converts a GroupSlice into a string slice
	testDataPointGroupsToStringSliceFunction := func() []string {
		var tempStringSlice []string

		for _, testDataPointGroup := range testDataEngine.TestDataPointGroups {
			tempStringSlice = append(tempStringSlice, string(testDataPointGroup))
		}

		return tempStringSlice
	}

	// Create function that converts a TestDataPointsSlice into a string slice
	testDataPointsToStringSliceFunction := func() []string {
		var tempStringSlice []string

		for _, testDataPointForAGroup := range testDataEngine.TestDataPointsForAGroup {
			tempStringSlice = append(tempStringSlice, string(testDataPointForAGroup))
		}

		return tempStringSlice
	}

	// Create the Group dropdown
	newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect = widget.NewSelect(testDataPointGroupsToStringSliceFunction(), func(selected string) {

		// Find List-item to select
		for index, group := range testDataEngine.TestDataPointGroups {
			if string(group) == selected {
				newOrEditTestDataPointGroupUI.SelectedIndexForGroups = index

			}
		}

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction())
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.Refresh()

		// Select the correct item in the Groups-List
		newOrEditTestDataPointGroupUI.TestDataPointGroupsList.Select(newOrEditTestDataPointGroupUI.SelectedIndexForGroups)
		newOrEditTestDataPointGroupUI.TestDataPointGroupsList.Refresh()

		// UnSelect in DropDown- and List for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.UnselectAll()
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = -1

	})

	// Create the Groups TestDataPoints dropdown
	newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect = widget.NewSelect(testDataPointsToStringSliceFunction(), func(selected string) {

		// Find List-item to select
		for index, group := range testDataEngine.TestDataPointsForAGroup {
			if string(group) == selected {
				newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = index

			}
		}

		// Select the correct item in the TestDataPoints-List
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.Select(newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints)
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.Refresh()
	})

	// Create List UI for 'testDataPointsForAGroup'
	newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList = widget.NewList(
		func() int { return len(testDataEngine.TestDataPointsForAGroup) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataEngine.TestDataPointsForAGroup[id]))
		},
	)

	newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.OnSelected = func(id widget.ListItemID) {
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = id

		// Select correct Group in Select-dropdown
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetSelected(string(testDataEngine.TestDataPointsForAGroup[id]))

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction())
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.Refresh()

	}

	var testDataPointGroupsContainer *fyne.Container
	testDataPointGroupsContainer = container.NewBorder(newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect,
		nil, nil, nil, newOrEditTestDataPointGroupUI.TestDataPointGroupsList)

	var testDataPointsForAGroupContainer *fyne.Container
	testDataPointsForAGroupContainer = container.NewBorder(newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect,
		nil, nil, nil, newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList)

	// Create Split Container used for 'testDataPointGroups' and 'testDataPointsForAGroup'
	var testDataGroupsAndPointsContainer *container.Split
	testDataGroupsAndPointsContainer = container.NewHSplit(testDataPointGroupsContainer, testDataPointsForAGroupContainer)

	var responseChannel chan testDataEngine.ResponseChannelStruct
	responseChannel = make(chan testDataEngine.ResponseChannelStruct)

	// The structure holding Group and TestDataPoints together
	//var newOrEditedChosenTestDataPointsPerGroupMap map[testDataEngine.TestDataPointGroupNameType][]TestDataPointRowUuidType
	//newOrEditedChosenTestDataPointsPerGroupMap = make(map[testDataEngine.TestDataPointGroupNameType][]TestDataPointRowUuidType)

	// Crete the 'New'-button for creating a new Group for TestDataPoints
	newButton := widget.NewButton("New", func() {
		myWindow.Hide()
		newOrEditTestDataPointGroupUI.ShowNewOrEditGroupWindow(
			app,
			myWindow,
			true,
			&responseChannel,
			"",
			&testDataEngine.ChosenTestDataPointsPerGroupMap,
			testDataDomainUuid,
			testDataAreaUuid)
	})

	// Crete the 'Edit'-button for editing an existing Group for TestDataPoints
	editButton := widget.NewButton("Edit", func() {
		if newOrEditTestDataPointGroupUI.SelectedIndexForGroups == -1 || len(testDataEngine.TestDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}
		myWindow.Hide()
		newOrEditTestDataPointGroupUI.ShowNewOrEditGroupWindow(
			app,
			myWindow,
			false,
			&responseChannel,
			testDataEngine.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups],
			&testDataEngine.ChosenTestDataPointsPerGroupMap,
			testDataDomainUuid,
			testDataAreaUuid)
	})

	// Crete the 'Delete'-button for deleting an existing Group for TestDataPoints
	deleteButton := widget.NewButton("Delete", func() {
		if newOrEditTestDataPointGroupUI.SelectedIndexForGroups == -1 || len(testDataEngine.TestDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}

		dialog.ShowConfirm("Confirm to Delete", fmt.Sprintf("Are you sure that you what to delete TestDataPointGroup '%s'?", testDataEngine.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups]), func(confirm bool) {
			if confirm {

				// Get the GroupName from the List to be deleted
				var groupNameToDelete testDataEngine.TestDataPointGroupNameType
				groupNameToDelete = testDataEngine.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups]

				// Delete the group
				delete(testDataEngine.ChosenTestDataPointsPerGroupMap, groupNameToDelete)

				// Rebuild the TestDataPointGroup-list
				testDataEngine.TestDataPointGroups = nil
				for testDataPointsGroupName, _ := range testDataEngine.ChosenTestDataPointsPerGroupMap {

					testDataEngine.TestDataPointGroups = append(testDataEngine.TestDataPointGroups, testDataPointsGroupName)
				}

				newOrEditTestDataPointGroupUI.SelectedIndexForGroups = -1

				newOrEditTestDataPointGroupUI.TestDataPointGroupsList.Refresh()
				newOrEditTestDataPointGroupUI.TestDataPointGroupsList.UnselectAll()

				// Clear the TestDataPointsList
				testDataEngine.TestDataPointsForAGroup = nil
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.Refresh()

				// UnSelect in DropDown- and List for TestDataPoints
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupList.UnselectAll()
				newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = -1
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

		var shouldListBeUpdated testDataEngine.ResponseChannelStruct
		var groupNameIndex int
		var groupNameIndexToSelect int

		for {

			shouldListBeUpdated = <-responseChannel

			// Update the List in main window if true as response
			if shouldListBeUpdated.ShouldBeUpdated == true {

				// Clear slice and variables used
				testDataEngine.TestDataPointGroups = nil
				groupNameIndex = 0
				groupNameIndexToSelect = 0

				for testDataPointsGroupName, _ := range testDataEngine.ChosenTestDataPointsPerGroupMap {

					testDataEngine.TestDataPointGroups = append(testDataEngine.TestDataPointGroups, testDataPointsGroupName)

					if testDataPointsGroupName == shouldListBeUpdated.TestDataPointGroupName {

						groupNameIndexToSelect = groupNameIndex

					}

					groupNameIndex = groupNameIndex + 1

				}
				newOrEditTestDataPointGroupUI.TestDataPointGroupsList.Refresh()
				newOrEditTestDataPointGroupUI.TestDataPointGroupsList.UnselectAll()
				newOrEditTestDataPointGroupUI.TestDataPointGroupsList.Select(groupNameIndexToSelect)
				newOrEditTestDataPointGroupUI.SelectedIndexForGroups = groupNameIndexToSelect

				// Select the correct group in the dropdown for groups
				newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect.SetOptions(testDataPointGroupsToStringSliceFunction())
				newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect.SetSelected(string(shouldListBeUpdated.TestDataPointGroupName))
				newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect.Refresh()

				// Select the correct TestDataPoint in the dropdown for TestDataPoints
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction())
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.Refresh()

			}
		}
	}()

	myWindow.Show()

}
