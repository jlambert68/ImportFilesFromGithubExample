package testDataSelector

import "fyne.io/fyne/v2/widget"

// The slices for Groups ans TestDataPoints for a Group
var testDataPointGroups []testDataPointGroupNameType // Define testDataPointGroups
var testDataPointsForAGroup []testDataPointNameType  // Define testDataPointGroups

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
	groupNameIsUnique                  string = "Group name is OK!"
	groupNameIsNotUnique               string = "Group name already exists!"
	groupNameIsEmpty                   string = "Group name can't be empty"
	groupNameCanNotStartOrEndWithSpace string = "Group name can not start or end with a space"
)

type responseChannelStruct struct {
	shouldBeUpdated        bool
	testDataPointGroupName testDataPointGroupNameType
}
