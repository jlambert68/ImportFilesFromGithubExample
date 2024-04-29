package testDataSelector

import "fyne.io/fyne/v2/widget"

// TestDataModelMapType
// Model holding Testdata for one or more Domains
type TestDataModelMapType map[TestDataDomainUuidType]*TestDataDomainModelMapStruct

// TestDataDomainModelMapStruct
// DataData for one Domain
type TestDataDomainModelMapStruct struct {
	TestDataDomainUuid TestDataDomainUuidType
	TestDataDomainName TestDataDomainNameType
	TestDataAreasMap   map[TestDataAreaUuidType]*TestDataAreaStruct
}

// TestDataAreaStruct
// TestData for one Area within one Domain
type TestDataAreaStruct struct {
	TestDataDomainUuid       TestDataDomainUuidType
	TestDataDomainName       TestDataDomainNameType
	TestDataAreaUuid         TestDataAreaUuidType
	TestDataAreaName         TestDataAreaNameType
	TestDataForDomainAreaMap map[TestDataRowMapUuidType]*[]*TestDataRowStruct // All TestData for a specific area in a specific Domain

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

// TestDataRowMapUuidType
// The Uuid for one TestDataRow. This uuid is unique within a Domain
type TestDataRowMapUuidType string

// TestDataRowStruct
// A specific Row of TestData
type TestDataRowStruct struct {
	TestDataDomainUuid          TestDataDomainUuidType
	TestDataDomainName          TestDataDomainNameType
	TestDataAreaUuid            TestDataAreaUuidType
	TestDataAreaName            TestDataAreaNameType
	TestDataValuesFromRowMap    *map[TestDataPointRowUuidType]*[]*TestDataPointValueStruct
	TestDataValuesFromColumnMap map[TestDataColumnUuidType]*[]*TestDataPointValueStruct
	TestDataColumnsMetaDataMap  map[TestDataColumnUuidType]*[]*TestDataColumnMetaDataStruct
}

// TestDataPointValueStruct
// Holding all information about one TestDataValue
type TestDataPointValueStruct struct {
	TestDataColumnUuid     TestDataColumnUuidType
	TestDataColumnDataName TestDataColumnDataNameType
	TestDataColumnUIName   TestDataColumnUINameType
	TestDataPointRowUuid   TestDataPointRowUuidType
	TestDataValue          TestDataValueType
}

// TestDataColumnMetaDataStruct
// Holds the information about a column regarding if data should be visible when selecting for data, and some other stuff
type TestDataColumnMetaDataStruct struct {
	TestDataColumnUuid                      TestDataColumnUuidType
	TestDataColumnDataName                  TestDataColumnDataNameType
	TestDataColumnUIName                    TestDataColumnUINameType
	ShouldColumnBeUsedForFindingTestData    bool
	ShouldColumnBeUsedWithinTestDataSetName bool
}

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

// TestDataValueType
// The value for specific TestDataPoint-value
type TestDataValueType string

// The slices for Groups ans TestDataPoints for a Group
var testDataPointGroups []testDataPointGroupNameType // Define testDataPointGroups
var testDataPointsForAGroup []testDataPointNameType  // Define testDataPointGroups

// The List-items for Groups ans TestDataPoints for a Group
var testDataPointGroupsList *widget.List
var testDataPointsForAGroupList *widget.List

// Current selected index for List items
var selectedIndexForGroups int = -1
var selectedIndexForGroupTestDataPoints int = -1

// The Select-items for Groups ans TestDataPoints for a Group
var testDataPointGroupsSelect *widget.Select
var testDataPointsForAGroupSelect *widget.Select

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
