package testDataSelector

import (
	"ImportFilesFromGithub/newOrEditTestDataPointGroupUI"
	"embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jlambert68/FenixScriptEngine/testDataEngine"
	"log"
)

//go:embed testData/FenixRawTestdata_646rows_211220.csv
var embeddedFile_SubCustody_MainTestDataArea embed.FS

//go:embed testData/FenixRawTestdata_3rows_240705.csv
var embeddedFile_SubCustody_ExtraTestDataArea embed.FS

//go:embed testData/FenixRawTestdata_10rows_240705.csv
var embeddedFile_CustodyCash_MainTestDataArea embed.FS

const (
	testDataDomainUuid testDataEngine.TestDataDomainUuidType = "7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae"
	testDataDomainName testDataEngine.TestDataDomainNameType = "Sub Custody"
	testDataAreaUuid   testDataEngine.TestDataAreaUuidType   = "010cc994-a913-4041-96fe-a96d7e0c97e8"
	testDataAreaName   testDataEngine.TestDataAreaNameType   = "Main TestData Area"
)

func ImportTestData_SubCustody_MainTestDataArea() {

	// Read the embedded file
	data, err := embeddedFile_SubCustody_MainTestDataArea.ReadFile("testData/FenixRawTestdata_646rows_211220.csv")
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	var testDataFromTestDataArea testDataEngine.TestDataFromSimpleTestDataAreaStruct

	testDataFromTestDataArea = testDataEngine.ImportEmbeddedSimpleCsvTestDataFile(
		data, ';')

	testDataEngine.AddTestDataToTestDataModel(testDataFromTestDataArea)

}

func ImportTestData_SubCustody_ExtraTestDat20240702074005aArea() {

	// Read the embedded file
	data, err := embeddedFile_SubCustody_ExtraTestDataArea.ReadFile("testData/FenixRawTestdata_3rows_240705.csv")
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	var testDataFromTestDataArea testDataEngine.TestDataFromSimpleTestDataAreaStruct

	testDataFromTestDataArea = testDataEngine.ImportEmbeddedSimpleCsvTestDataFile(
		data, ';')

	testDataEngine.AddTestDataToTestDataModel(testDataFromTestDataArea)

}

func ImportTestData_CustodyCash_MainTestDataArea() {

	// Read the embedded file
	data, err := embeddedFile_CustodyCash_MainTestDataArea.ReadFile("testData/FenixRawTestdata_10rows_240705.csv")
	if err != nil {
		log.Fatalf("Error reading the embedded file: %v", err)
	}

	var testDataFromTestDataArea testDataEngine.TestDataFromSimpleTestDataAreaStruct

	testDataFromTestDataArea = testDataEngine.ImportEmbeddedSimpleCsvTestDataFile(
		data, ';')

	testDataEngine.AddTestDataToTestDataModel(testDataFromTestDataArea)

}

func MainTestDataSelector(
	app fyne.App,
	parent fyne.Window,
	testDataForGroupObject *testDataEngine.TestDataForGroupObjectStruct) {

	parent.Hide()

	myWindow := app.NewWindow("TestData Management")
	myWindow.Resize(fyne.NewSize(600, 500))

	// When this window closed then show parent and send response to parent window
	myWindow.SetOnClosed(func() {
		parent.Show()
	})

	// Initiate 'chosenTestDataPointsPerGroupMap'
	if testDataForGroupObject.ChosenTestDataPointsPerGroupMap == nil {
		testDataForGroupObject.ChosenTestDataPointsPerGroupMap = make(map[testDataEngine.TestDataPointGroupNameType]*testDataEngine.TestDataPointNameMapType)
	}

	// Create List UI for 'testDataPointGroups'
	testDataPointGroupsList = widget.NewList(
		func() int { return len(testDataForGroupObject.TestDataPointGroups) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataForGroupObject.TestDataPointGroups[id]))
		},
	)

	testDataPointGroupsList.OnSelected = func(id widget.ListItemID) {
		newOrEditTestDataPointGroupUI.SelectedIndexForGroups = id

		// Update List for  'testDataPointsForAGroup'
		updateTestDataPointsForAGroupList(testDataForGroupObject.TestDataPointGroups[id], testDataForGroupObject)

		// Select correct Group in Select-dropdown
		newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect.SetSelected(string(testDataForGroupObject.TestDataPointGroups[id]))

		// UnSelect in DropDown- and List for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
		testDataPointsForAGroupList.UnselectAll()
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = -1

	}

	// Create function that converts a GroupSlice into a string slice
	testDataPointGroupsToStringSliceFunction := func() []string {
		var tempStringSlice []string

		for _, testDataPointGroup := range testDataForGroupObject.TestDataPointGroups {
			tempStringSlice = append(tempStringSlice, string(testDataPointGroup))
		}

		return tempStringSlice
	}

	// Create function that converts a TestDataPointsSlice into a string slice
	testDataPointsToStringSliceFunction := func() []string {
		var tempStringSlice []string

		for _, testDataPointForAGroup := range testDataForGroupObject.TestDataPointsForAGroup {
			tempStringSlice = append(tempStringSlice, string(testDataPointForAGroup))
		}

		return tempStringSlice
	}

	// Create the Group dropdown
	newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect = widget.NewSelect(testDataPointGroupsToStringSliceFunction(), func(selected string) {

		// Find List-item to select
		for index, group := range testDataForGroupObject.TestDataPointGroups {
			if string(group) == selected {
				newOrEditTestDataPointGroupUI.SelectedIndexForGroups = index

			}
		}

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction())
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.Refresh()

		// Select the correct item in the Groups-List
		testDataPointGroupsList.Select(newOrEditTestDataPointGroupUI.SelectedIndexForGroups)
		testDataPointGroupsList.Refresh()

		// UnSelect in DropDown- and List for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
		testDataPointsForAGroupList.UnselectAll()
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = -1

	})

	// Create the Groups TestDataPoints dropdown
	newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect = widget.NewSelect(testDataPointsToStringSliceFunction(), func(selected string) {

		// Find List-item to select
		for index, group := range testDataForGroupObject.TestDataPointsForAGroup {
			if string(group) == selected {
				newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = index

			}
		}

		// Select the correct item in the TestDataPoints-List
		testDataPointsForAGroupList.Select(newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints)
		testDataPointsForAGroupList.Refresh()
	})

	// Create List UI for 'testDataPointsForAGroup'
	testDataPointsForAGroupList = widget.NewList(
		func() int { return len(testDataForGroupObject.TestDataPointsForAGroup) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(string(testDataForGroupObject.TestDataPointsForAGroup[id]))
		},
	)

	testDataPointsForAGroupList.OnSelected = func(id widget.ListItemID) {
		newOrEditTestDataPointGroupUI.SelectedIndexForGroupTestDataPoints = id

		// Select correct Group in Select-dropdown
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetSelected(string(testDataForGroupObject.TestDataPointsForAGroup[id]))

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction())
		newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.Refresh()

	}

	var testDataPointGroupsContainer *fyne.Container
	testDataPointGroupsContainer = container.NewBorder(newOrEditTestDataPointGroupUI.TestDataPointGroupsSelect,
		nil, nil, nil, testDataPointGroupsList)

	var testDataPointsForAGroupContainer *fyne.Container
	testDataPointsForAGroupContainer = container.NewBorder(newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect,
		nil, nil, nil, testDataPointsForAGroupList)

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
			&testDataForGroupObject.ChosenTestDataPointsPerGroupMap,
			testDataForGroupObject)
	})

	// Crete the 'Edit'-button for editing an existing Group for TestDataPoints
	editButton := widget.NewButton("Edit", func() {
		if newOrEditTestDataPointGroupUI.SelectedIndexForGroups == -1 || len(testDataForGroupObject.TestDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}
		myWindow.Hide()
		newOrEditTestDataPointGroupUI.ShowNewOrEditGroupWindow(
			app,
			myWindow,
			false,
			&responseChannel,
			testDataForGroupObject.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups],
			&testDataForGroupObject.ChosenTestDataPointsPerGroupMap,
			testDataForGroupObject)
	})

	// Crete the 'Delete'-button for deleting an existing Group for TestDataPoints
	deleteButton := widget.NewButton("Delete", func() {
		if newOrEditTestDataPointGroupUI.SelectedIndexForGroups == -1 || len(testDataForGroupObject.TestDataPointGroups) == 0 {
			dialog.ShowInformation("Error", "No selection made", myWindow)
			return
		}

		dialog.ShowConfirm("Confirm to Delete", fmt.Sprintf("Are you sure that you what to delete TestDataPointGroup '%s'?",
			testDataForGroupObject.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups]), func(confirm bool) {
			if confirm {

				// Get the GroupName from the List to be deleted
				var groupNameToDelete testDataEngine.TestDataPointGroupNameType
				groupNameToDelete = testDataForGroupObject.TestDataPointGroups[newOrEditTestDataPointGroupUI.SelectedIndexForGroups]

				// Delete the group
				delete(testDataForGroupObject.ChosenTestDataPointsPerGroupMap, groupNameToDelete)

				// Rebuild the TestDataPointGroup-list
				testDataForGroupObject.TestDataPointGroups = nil
				for testDataPointsGroupName, _ := range testDataForGroupObject.ChosenTestDataPointsPerGroupMap {

					testDataForGroupObject.TestDataPointGroups = append(testDataForGroupObject.TestDataPointGroups, testDataPointsGroupName)
				}

				newOrEditTestDataPointGroupUI.SelectedIndexForGroups = -1

				testDataPointGroupsList.Refresh()
				testDataPointGroupsList.UnselectAll()

				// Clear the TestDataPointsList
				testDataForGroupObject.TestDataPointsForAGroup = nil
				testDataPointsForAGroupList.Refresh()

				// UnSelect in DropDown- and List for TestDataPoints
				newOrEditTestDataPointGroupUI.TestDataPointsForAGroupSelect.ClearSelected()
				testDataPointsForAGroupList.UnselectAll()
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
				testDataForGroupObject.TestDataPointGroups = nil
				groupNameIndex = 0
				groupNameIndexToSelect = 0

				for testDataPointsGroupName, _ := range testDataForGroupObject.ChosenTestDataPointsPerGroupMap {

					testDataForGroupObject.TestDataPointGroups = append(testDataForGroupObject.TestDataPointGroups, testDataPointsGroupName)

					if testDataPointsGroupName == shouldListBeUpdated.TestDataPointGroupName {

						groupNameIndexToSelect = groupNameIndex

					}

					groupNameIndex = groupNameIndex + 1

				}
				testDataPointGroupsList.Refresh()
				testDataPointGroupsList.UnselectAll()
				testDataPointGroupsList.Select(groupNameIndexToSelect)
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
