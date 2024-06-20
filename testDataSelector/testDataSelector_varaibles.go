package testDataSelector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TestDataModelStruct
// The structure holding all TestData
type TestDataModelStruct struct {
	TestDataDomainAndAreaNameToUuidMap *map[TestDataDomainOrAreaNameType]TestDataDomainOrAreaUuidType
	TestDataModelMap                   *map[TestDataDomainUuidType]*TestDataDomainModelStruct
}

// TestDataDomainModelStruct
// DataData for one Domain
type TestDataDomainModelStruct struct {
	TestDataDomainUuid TestDataDomainUuidType
	TestDataDomainName TestDataDomainNameType
	TestDataAreasMap   *map[TestDataAreaUuidType]*TestDataAreaStruct
}

// TestDataAreaStruct
// TestData for one Area within one Domain
type TestDataAreaStruct struct {
	TestDataDomainUuid                   TestDataDomainUuidType
	TestDataDomainName                   TestDataDomainNameType
	TestDataAreaUuid                     TestDataAreaUuidType
	TestDataAreaName                     TestDataAreaNameType
	TestDataValuesForRowMap              *map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	TestDataValuesForRowNameMap          *map[TestDataValueNameType]*[]TestDataPointRowUuidType
	TestDataValuesForColumnMap           *map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	TestDataValuesForColumnAndRowUuidMap *map[TestDataColumnAndRowUuidType]*TestDataPointValueStruct
	TestDataColumnsMetaDataMap           *map[TestDataColumnUuidType]*TestDataColumnMetaDataStruct
	UniqueTestDataValuesForColumnMap     *map[TestDataColumnUuidType]*map[TestDataValueType][]TestDataPointRowUuidType
}

// TestDataPointValueStruct
// Holding all information about one TestDataValue
type TestDataPointValueStruct struct {
	TestDataDomainUuid           TestDataDomainUuidType
	TestDataDomainName           TestDataDomainNameType
	TestDataAreaUuid             TestDataAreaUuidType
	TestDataAreaName             TestDataAreaNameType
	TestDataColumnUuid           TestDataColumnUuidType
	TestDataColumnDataName       TestDataColumnDataNameType
	TestDataColumnUIName         TestDataColumnUINameType
	TestDataPointRowUuid         TestDataPointRowUuidType
	TestDataColumnAndRowUuid     TestDataColumnAndRowUuidType
	TestDataValue                TestDataValueType
	TestDataValueNameDescription TestDataValueNameDescriptionType
	TestDataValueName            TestDataValueNameType
}

// TestDataColumnMetaDataStruct
// Holds the information about a column regarding if data should be visible when selecting for data, and some other stuff
type TestDataColumnMetaDataStruct struct {
	TestDataColumnUuid     TestDataColumnUuidType
	TestDataColumnDataName TestDataColumnDataNameType
	TestDataColumnUIName   TestDataColumnUINameType
	//TestDataPointRowsUuid                   []TestDataPointRowUuidType
	ShouldColumnBeUsedForFindingTestData    bool
	ShouldColumnBeUsedWithinTestDataSetName bool
}

// TestDataDomainUuidType
// The UUID of the Domain that owns the TestData
type TestDataDomainUuidType string

// TestDataDomainNameType
// The Name of the Domain that owns the TestData
type TestDataDomainNameType string

// TestDataAreaUuidType
// The UUID for a specific TestData-area within the Domain
type TestDataAreaUuidType string

// TestDataAreaNameType
// The Name for a specific TestData-area within the Domain
type TestDataAreaNameType string

// TestDataDomainOrAreaUuidType
// The UUID for a specific TestData-domain or TestData-area
type TestDataDomainOrAreaUuidType string

// TestDataDomainOrAreaNameType
// The Name for a specific TestData-domain or TestData-area
type TestDataDomainOrAreaNameType string

// TestDataPointRowUuidType
// Identifies one TestDataPoint/TestDataRow
type TestDataPointRowUuidType string

// TestDataColumnUuidType
// Identifies a Column that goes across TestDataPoints/TestDataRows
type TestDataColumnUuidType string

// TestDataColumnDataNameType
// The DataName used for Placeholder and other data manipulations connected to a specific column in the TestDataSet
type TestDataColumnDataNameType string

// TestDataColumnUINameType
// The Name used in UI connected to a specific column in the TestDataSet
type TestDataColumnUINameType string

// TestDataColumnAndRowUuidType
// The hash ColumnUuid 'concat' RowUuid: SHA256(TestDataColumnUuidType 'concat' TestDataPointRowUuidType)
type TestDataColumnAndRowUuidType string

// TestDataValueType
// The value for specific TestDataPoint-value
type TestDataValueType string

// TestDataValueNameDescriptionType
// The Description of how the name for a TestDataPoint is constructed. Used to show/describe the TestDataPoint in the UI
type TestDataValueNameDescriptionType string

// TestDataValueNameType
// The name for a TestDataPoint. Used to show/describe the TestDataPoint in the UI
type TestDataValueNameType string

// The slices for Groups ans TestDataPoints for a Group
var testDataPointGroups []testDataPointGroupNameType   // Define testDataPointGroups
var testDataPointsForAGroup []TestDataPointRowUuidType // Define testDataPointGroups

// The List-items for Groups ans TestDataPoints for a Group
var testDataPointGroupsList *widget.List
var testDataPointsForAGroupList *widget.List

// Current selected index for List items
var selectedIndexForGroups int = -1
var selectedIndexForGroupTestDataPoints int = -1

// The Select-items for Groups ans TestDataPoints for a Group
var testDataPointGroupsSelect *widget.Select
var testDataPointsForAGroupSelect *widget.Select

// The type for Group name
type testDataPointGroupNameType string

// The map holding all data for Groups and TestDataPoints for a group
var chosenTestDataPointsPerGroupMap map[testDataPointGroupNameType]*testDataPointNameMapType

// The type for the map that holds the connection between TestDataÅpintName and all the RowUUids connectoed to it
type testDataPointNameMapType map[TestDataValueNameType]*[]*dataPointTypeForGroupsStruct

// Types used for data structures for a specific rowValue
/*type testDataPointRowStruct struct {
	testDataDomainUuid           TestDataDomainUuidType
	testDataDomainName           TestDataDomainNameType
	testDataAreaUuid             TestDataAreaUuidType
	testDataAreaName             TestDataAreaNameType
	testDataPointUuid            TestDataPointRowUuidType
	testDataPointName            TestDataValueNameType
	testDataPointNameDescription TestDataValueNameDescriptionType
	testDatapointValue           TestDataValueType
} */

// Error/warning texts for control of Group Name
const (
	groupNameIsUnique                  string = "Group name is OK!"
	groupNameIsNotUnique               string = "Group name already exists!"
	groupNameIsEmpty                   string = "Group name can't be empty"
	groupNameCanNotStartOrEndWithSpace string = "Group name can not start or end with a space"
)

// Message sent back when a Group is Created or is Edited
type responseChannelStruct struct {
	shouldBeUpdated        bool
	testDataPointGroupName testDataPointGroupNameType
}

// Structure for the Group
type dataPointTypeForGroupsStruct struct {
	testDataDomainUuid            TestDataDomainUuidType
	testDataDomainName            TestDataDomainNameType
	testDataAreaUuid              TestDataAreaUuidType
	testDataAreaName              TestDataAreaNameType
	testDataPointName             TestDataValueNameType
	searchResultDataPointUuidMap  map[TestDataPointRowUuidType]TestDataPointRowUuidType
	availableTestDataPointUuidMap map[TestDataPointRowUuidType]TestDataPointRowUuidType
	selectedTestDataPointUuidMap  map[TestDataPointRowUuidType]TestDataPointRowUuidType
}

// Slices used to keep track of filtered, available and selected DataPoints
var allPointsAvailable []dataPointTypeForGroupsStruct
var allSelectedPoints []dataPointTypeForGroupsStruct

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
