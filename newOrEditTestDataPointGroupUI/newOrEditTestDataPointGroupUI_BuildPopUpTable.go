package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fmt"
	"regexp"
)

// Build the Table Data, based on TestDataPointName, to be used when the popup table is shown to the user to pick from
func buildPopUpTableDataFromTestDataPointName(
	tempTestDataPointRowName string,
	testDataModel *testDataEngine.TestDataModelStruct) (
	tableData [][]string) {

	re := regexp.MustCompile(`^([^/]+)/([^/]+)`)

	matches := re.FindStringSubmatch(tempTestDataPointRowName)
	if len(matches) > 2 {
		fmt.Println("First part:", matches[1])
		fmt.Println("Second part:", matches[2])
	} else {
		fmt.Println("No matching parts found")
	}

	var tempTestDataModelMap map[testDataEngine.TestDataDomainUuidType]*testDataEngine.TestDataDomainModelStruct
	var tempTestDataDomainModel testDataEngine.TestDataDomainModelStruct
	var tempTestDataAreaMap map[testDataEngine.TestDataAreaUuidType]*testDataEngine.TestDataAreaStruct
	var tempTestDataArea testDataEngine.TestDataAreaStruct
	var tempTestDataDomainAndAreaNameToUuidMap map[testDataEngine.TestDataDomainOrAreaNameType]testDataEngine.TestDataDomainOrAreaUuidType
	var tempTestDataValuesForRowNameMap map[testDataEngine.TestDataValueNameType]*[]testDataEngine.TestDataPointRowUuidType
	var tempTestDataValuesForRowMap map[testDataEngine.TestDataPointRowUuidType]*[]*testDataEngine.TestDataPointValueStruct
	var tempTestDataValuesForRowUuidMapBaseOnNameSlice []testDataEngine.TestDataPointRowUuidType

	var tempTestDataDomainOrAreaUuid testDataEngine.TestDataDomainOrAreaUuidType
	var tempTestDataDomainUuid testDataEngine.TestDataDomainUuidType
	var tempTestDataAreaUuid testDataEngine.TestDataAreaUuidType

	tempTestDataModelMap = *testDataModel.TestDataModelMap

	// Extract the UUID for Domain and Area
	tempTestDataDomainAndAreaNameToUuidMap = *testDataModel.TestDataDomainAndAreaNameToUuidMap
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[testDataEngine.TestDataDomainOrAreaNameType(matches[1])]
	tempTestDataDomainUuid = testDataEngine.TestDataDomainUuidType(tempTestDataDomainOrAreaUuid)
	tempTestDataDomainOrAreaUuid, _ = tempTestDataDomainAndAreaNameToUuidMap[testDataEngine.TestDataDomainOrAreaNameType(matches[2])]
	tempTestDataAreaUuid = testDataEngine.TestDataAreaUuidType(tempTestDataDomainOrAreaUuid)

	// Extract Domain and Area maps
	tempTestDataDomainModel = *tempTestDataModelMap[tempTestDataDomainUuid]
	tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
	tempTestDataArea = *tempTestDataAreaMap[tempTestDataAreaUuid]
	tempTestDataValuesForRowNameMap = *tempTestDataArea.TestDataValuesForRowNameMap
	tempTestDataValuesForRowMap = *tempTestDataArea.TestDataValuesForRowMap

	var tempTestDataPointRowNameToSearchFor string
	tempTestDataPointRowNameToSearchFor = tempTestDataPointRowName //[len(matches[0]+"/"):]

	tempTestDataValuesForRowUuidMapBaseOnNameSlice = *tempTestDataValuesForRowNameMap[testDataEngine.TestDataValueNameType(tempTestDataPointRowNameToSearchFor)]

	// Loop the slice to extract the RowUuids

	var headerSlice []string
	for rowIndex, tempTestDataPointRowUuid := range tempTestDataValuesForRowUuidMapBaseOnNameSlice {

		var rowSlice []string
		tempTestDataValuesForRowSlice := tempTestDataValuesForRowMap[tempTestDataPointRowUuid]

		// Loop the slice with RowValue
		for _, tempTestDataPointValue := range *tempTestDataValuesForRowSlice {

			// Create a header slice
			if rowIndex == 0 {
				headerSlice = append(headerSlice, string(tempTestDataPointValue.TestDataColumnUIName))
			}

			rowSlice = append(rowSlice, string(tempTestDataPointValue.TestDataValue))
		}

		// Append data for 'TestDataPointRowName'
		rowSlice = append(rowSlice, tempTestDataPointRowName)

		// Append data for unique 'TestDataPointRowUuid'
		rowSlice = append(rowSlice, string(tempTestDataPointRowUuid))

		// Add a header when first row
		if rowIndex == 0 {

			// Add headers for 'TestDataPointRowName' and 'TestDataPointRowUuid'
			headerSlice = append(headerSlice, "TestDataPointRowName")
			headerSlice = append(headerSlice, "TestDataPointRowUuid")

			tableData = append(tableData, headerSlice)
		}

		// Add the data
		tableData = append(tableData, rowSlice)

	}

	return tableData
}
