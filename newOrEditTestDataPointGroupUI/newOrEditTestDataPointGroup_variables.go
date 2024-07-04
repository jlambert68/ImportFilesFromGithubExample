package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
// Message sent back when a Group is Created or is Edited
type ResponseChannelStruct struct {
	ShouldBeUpdated        bool
	TestDataPointGroupName testDataEngine.TestDataPointGroupNameType
}

*/

// Current selected index for List items
var SelectedIndexForGroups int = -1
var SelectedIndexForGroupTestDataPoints int = -1

// The Select-items for Groups ans TestDataPoints for a Group
var TestDataPointGroupsSelect *widget.Select
var TestDataPointsForAGroupSelect *widget.Select

// The List-widget holding all available TestDataPoints from Search
var allAvailablePointsList *widget.List

// The List-widget holding all selected TestDataPoints from Search
var selectedPointsList *widget.List

// *** Create the selection boxes for selecting TestDataValues values
var testDataSelectionsContainer *fyne.Container

// Create the container for Search- and Clear- buttons
var searchAndClearButtonsContainer *fyne.Container

// Layout configuration for the new/edit window
// Create the UpperAndLowerSplitContainer
var upperAndLowerSplitContainer *container.Split
var listsSplitContainer *container.Split

var upperSplitContainer *fyne.Container

var lowerRightSideContainer *fyne.Container

// The List-items for Groups ans TestDataPoints for a Group
var TestDataPointGroupsList *widget.List
var TestDataPointsForAGroupList *widget.List

var setStateForSaveButtonAndGroupNameTextEntryExternalCall func()

// Slices used to keep track of filtered, available and selected DataPoints
var allPointsAvailable []testDataEngine.DataPointTypeForGroupsStruct
var allSelectedPoints []testDataEngine.DataPointTypeForGroupsStruct
