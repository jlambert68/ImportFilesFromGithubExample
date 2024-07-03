package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"strings"
)

func generateAllAvailablePointsListUIComponent(
	newOrEditTestDataPointGroupWindow *fyne.Window,
	testDataModel *testDataEngine.TestDataModelStruct) {

	// Create and configure the list-component of all TestDataPoints
	allAvailablePointsList = widget.NewList(
		func() int { return len(testDataEngine.AllPointsAvailable) },
		func() fyne.CanvasObject {

			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {

			obj.(*widget.Label).SetText(fmt.Sprintf(
				"%s [%d(%d)]",
				string(testDataEngine.AllPointsAvailable[id].TestDataPointName),
				len(testDataEngine.AllPointsAvailable[id].AvailableTestDataPointUuidMap),
				len(testDataEngine.AllPointsAvailable[id].AvailableTestDataPointUuidMap)+
					len(testDataEngine.AllPointsAvailable[id].SelectedTestDataPointUuidMap)))
		},
	)

	allAvailablePointsList.OnSelected = func(id widget.ListItemID) {

		// Remove the number part of the visible name
		var clickedDataPointName string
		clickedDataPointName = string(testDataEngine.AllPointsAvailable[id].TestDataPointName) //filterToRemoveNumberOfSimilarTestDataPointsInName(string(allPointsAvailable[id].testDataPointName))

		var tableData [][]string
		tableData = buildPopUpTableDataFromTestDataPointName(clickedDataPointName, testDataModel)

		showTable(
			*newOrEditTestDataPointGroupWindow,
			tableData,
			testDataEngine.AllPointsAvailable[id].SelectedTestDataPointUuidMap)

		allAvailablePointsList.UnselectAll()

		allAvailablePointsList.Refresh()
		selectedPointsList.Refresh()

	}

}

// Removes the part of the name that specifies the number similar TestDataPoints, i.e.
// Sub Custody/Main TestData Area/SEK/AccTest/SE/CRDT/CH/Switzerland/BBH/EUR/EUR/SEK [2] ->
// Sub Custody/Main TestData Area/SEK/AccTest/SE/CRDT/CH/Switzerland/BBH/EUR/EUR/SEK
func filterToRemoveNumberOfSimilarTestDataPointsInName(dataPointNameToClean string) (cleanedName string) {
	re := regexp.MustCompile(`\[(\d+)\]`)

	// FindStringSubmatch returns an array of matches where the first element is the full match,
	// and the subsequent ones are the captured groups.
	matches := re.FindStringSubmatch(dataPointNameToClean)
	if len(matches) > 1 { // matches[0] is the full match, matches[1] would be the first captured group
		fmt.Println("Extracted number:", matches[1])
	} else {
		fmt.Println("No number found")
	}

	re = regexp.MustCompile(`^(.*?)\[\d+\]`)

	// FindStringSubmatch returns an array of matches where the first element is the full match,
	// and the subsequent ones are the captured groups.
	matches = re.FindStringSubmatch(dataPointNameToClean)

	if len(matches) > 1 { // matches[0] is the full match, matches[1] would be the first captured group
		cleanedName = strings.Trim(matches[1], " ")
		fmt.Println(fmt.Sprintf("Extracted text to the left: '%s'", cleanedName))
	} else {
		fmt.Println("No matching text found")
	}

	return cleanedName

}
