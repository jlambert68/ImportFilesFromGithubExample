package newOrEditTestDataPointGroupUI

import (
	"ImportFilesFromGithub/testDataEngine"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strings"
)

// *** Create the selection boxes for selecting TestDataValues values
func generateTestDataSelectionsUIComponent(
	testDataModel *testDataEngine.TestDataModelStruct,
	testDataModelMap map[testDataEngine.TestDataDomainUuidType]*testDataEngine.TestDataDomainModelStruct,
	testDataDomainUuid testDataEngine.TestDataDomainUuidType,
	testDataAreaUuid testDataEngine.TestDataAreaUuidType) {

	var existInMap bool

	var searchResult []testDataEngine.TestDataPointRowUuidType

	// Variable to handel DropDown for Domains
	var domainOptions []string
	var domains []*testDataEngine.TestDataDomainModelStruct
	var domainsLabel *widget.Label
	var domainsSelect *widget.Select
	var testDomainContainer *fyne.Container

	// Variable to handel DropDown for TestDataAreas for a chosen Domain
	var testAreaOptions []string
	var testAreas []*testDataEngine.TestDataAreaStruct
	var testAreasLabel *widget.Label
	var testAreaSelect *widget.Select
	var testAreasContainer *fyne.Container
	var testAreaMap *map[testDataEngine.TestDataAreaUuidType]*testDataEngine.TestDataAreaStruct

	type testDataValueSelectionStruct struct {
		testDataSelectionLabel       *widget.Label
		testDataCheckGroup           *widget.CheckGroup
		TestDataColumnUuid           testDataEngine.TestDataColumnUuidType
		TestDataColumnDataName       testDataEngine.TestDataColumnDataNameType
		TestDataPointValueRowUuidMap *map[testDataEngine.TestDataValueType]*[]testDataEngine.TestDataPointRowUuidType
	}
	var testDataValueSelections []*testDataValueSelectionStruct
	var testDataValuesSelectionContainer *fyne.Container

	// Create Search TestData-button
	var searchTestDataButton *widget.Button
	var clearTestDataFilterCheckBoxesButton *widget.Button

	// Create label for Domains
	domainsLabel = widget.NewLabel(testDataDomainLabelText)
	domainsLabel.TextStyle.Bold = true
	testAreasLabel = widget.NewLabel(testDataTestAreaLabelText)
	testAreasLabel.TextStyle.Bold = true

	// Extract TestData on Domain-level
	for _, tempTestDataDomainModel := range testDataModelMap {
		domainOptions = append(domainOptions, string(tempTestDataDomainModel.TestDataDomainName))
		domains = append(domains, tempTestDataDomainModel)
	}

	// Create Domain-Select-DropDown
	domainsSelect = widget.NewSelect(domainOptions, func(selected string) {

		// Clear UI object that need to be recreated

		// Extract correct TestArea
		for index, domain := range domains {
			if selected == string(domain.TestDataDomainName) {
				testAreaMap = domains[index].TestDataAreasMap
				break
			}
		}

		// Extract TestData on TestArea Level
		for _, tempTestDataArea := range *testAreaMap {
			testAreaOptions = append(testAreaOptions, string(tempTestDataArea.TestDataAreaName))
			testAreas = append(testAreas, tempTestDataArea)
		}

		// Create TestArea-Select-DropDown
		testAreaSelect = widget.NewSelect(testAreaOptions, func(selected string) {

			// Create available TestDataSelections for TestArea
			for _, testDataArea := range *testAreaMap {

				// Clear UI component that holds 'TestDataValuesSelections'
				testDataValuesSelectionContainer = container.NewHBox()

				// Create a slice with 'testDataColumnsMetaData' that can be sorted
				var testDataColumnsMetaDataToBeSorted []*testDataEngine.TestDataColumnMetaDataStruct
				for _, testDataColumnsMetaData := range *testDataArea.TestDataColumnsMetaDataMap {
					testDataColumnsMetaDataToBeSorted = append(testDataColumnsMetaDataToBeSorted, testDataColumnsMetaData)
				}

				// Sort the slice based on TestDataColumnUIName
				sort.Slice(testDataColumnsMetaDataToBeSorted, func(i, j int) bool {
					return testDataColumnsMetaDataToBeSorted[i].TestDataColumnUIName < testDataColumnsMetaDataToBeSorted[j].TestDataColumnUIName
				})

				// Loop 'testDataColumnsMetaDataToBeSorted' for Columns to present as separate CheckGroups
				for _, testDataColumnsMetaData := range testDataColumnsMetaDataToBeSorted {

					// Check if column should be used for filtering TestData as a CheckGroup
					if testDataColumnsMetaData.ShouldColumnBeUsedForFindingTestData == true {

						var checkGroupOptions []string
						var tempTestDataColumnContainer *fyne.Container

						// Set Label
						var newColumnFilterLabel *widget.Label
						newColumnFilterLabel = widget.NewLabel(string(testDataColumnsMetaData.TestDataColumnUIName))
						newColumnFilterLabel.TextStyle.Bold = true

						var tempTestDataPointValueRowUuidMap map[testDataEngine.TestDataValueType]*[]testDataEngine.TestDataPointRowUuidType
						tempTestDataPointValueRowUuidMap = make(map[testDataEngine.TestDataValueType]*[]testDataEngine.TestDataPointRowUuidType)

						var testDataValueSelection *testDataValueSelectionStruct
						testDataValueSelection = &testDataValueSelectionStruct{
							testDataSelectionLabel:       newColumnFilterLabel,
							testDataCheckGroup:           nil,
							TestDataColumnUuid:           testDataColumnsMetaData.TestDataColumnUuid,
							TestDataColumnDataName:       testDataColumnsMetaData.TestDataColumnDataName,
							TestDataPointValueRowUuidMap: &tempTestDataPointValueRowUuidMap,
						}

						// Extract the Map with the values
						var uniqueTestDataValuesForColumnMapPtr *map[testDataEngine.TestDataValueType][]testDataEngine.TestDataPointRowUuidType
						UniqueTestDataValuesForColumnMap := *testDataArea.UniqueTestDataValuesForColumnMap

						uniqueTestDataValuesForColumnMapPtr = UniqueTestDataValuesForColumnMap[testDataColumnsMetaData.TestDataColumnUuid]

						// Loop Values in Column and create Checkboxes, and store RowUuids for unique values
						for uniqueTestDataValue, testDataPointRowsUuid := range *uniqueTestDataValuesForColumnMapPtr {

							// Add value to slice for CheckBox-labels
							checkGroupOptions = append(checkGroupOptions, string(uniqueTestDataValue))

							// Add 'TestDataPointRowUuid' to correct slice for each unique value in the column
							var testDataPointRowUuidSlicePtr *[]testDataEngine.TestDataPointRowUuidType
							var testDataPointRowUuidSlice []testDataEngine.TestDataPointRowUuidType
							testDataPointRowUuidSlicePtr, existInMap = tempTestDataPointValueRowUuidMap[uniqueTestDataValue]

							if existInMap == false {
								var tempTestDataPointRowUuidSlice []testDataEngine.TestDataPointRowUuidType
								testDataPointRowUuidSlice = tempTestDataPointRowUuidSlice
							} else {
								testDataPointRowUuidSlice = *testDataPointRowUuidSlicePtr
							}

							testDataPointRowUuidSlice = append(testDataPointRowUuidSlice, testDataPointRowsUuid...)

							tempTestDataPointValueRowUuidMap[uniqueTestDataValue] = &testDataPointRowUuidSlice

						}

						// Sort values in CheckGroup
						sort.Strings(checkGroupOptions)

						// Create the CheckGroup
						var tempTestDataCheckGroup *widget.CheckGroup
						tempTestDataCheckGroup = widget.NewCheckGroup(checkGroupOptions, func(changed []string) {
							// Handle check change
						})

						// Add the CheckGroup
						testDataValueSelection.testDataCheckGroup = tempTestDataCheckGroup

						// Add 'testDataValueSelections' to slice
						testDataValueSelections = append(testDataValueSelections, testDataValueSelection)

						// Get the minimum size of the check group
						var testDataCheckGroupMinSize fyne.Size
						testDataCheckGroupMinSize = testDataValueSelection.testDataCheckGroup.MinSize()

						// Create the container having scrollbar the TestDataCheckGroup
						testDataCheckGroupContainer := container.NewScroll(testDataValueSelection.testDataCheckGroup)

						// Set
						testDataCheckGroupContainer.SetMinSize(fyne.NewSize(testDataCheckGroupContainer.Size().Height, testDataCheckGroupMinSize.Width))

						// Add to TestDataColumn-container
						tempTestDataColumnContainer = container.NewBorder(
							testDataValueSelection.testDataSelectionLabel,
							nil, nil, nil,
							testDataCheckGroupContainer)

						// Add 'tempTestDataColumnContainer' to 'testDataValuesSelectionContainer'
						testDataValuesSelectionContainer.Add(tempTestDataColumnContainer)

					}
				}
			}
		})

		// Set label for TestAreas
		testAreasLabel.SetText(fmt.Sprintf(testDataTestAreaLabelText+"'%s'", domainOptions[0]))

		// If there is only one item in TestArea-item then select that one
		if len(testAreaOptions) == 1 {
			testAreaSelect.SetSelected(testAreaOptions[0])
			testAreaSelect.Refresh()
		}

	})

	// If there is only one item in Domains-dropdown then select that one
	if len(domainOptions) == 1 {
		domainsSelect.SetSelected(domainOptions[0])
		domainsSelect.Refresh()

		// Set label for TestAreas
		testAreasLabel.SetText(fmt.Sprintf(testDataTestAreaLabelText+"'%s'", domainOptions[0]))
	}

	// Create the separate TestData-selection-containers
	testDomainContainer = container.NewVBox(domainsLabel, domainsSelect)
	testAreasContainer = container.NewVBox(testAreasLabel, testAreaSelect)

	// Create the main TestData-selection-container
	testDataSelectionsContainer = container.NewHBox(testDomainContainer, testAreasContainer, testDataValuesSelectionContainer)

	// Create Search TestData-button
	searchTestDataButton = widget.NewButton("Search for TestDataPoints", func() {

		var tempTestDataModelMap map[testDataEngine.TestDataDomainUuidType]*testDataEngine.TestDataDomainModelStruct
		var tempTestDataDomainModel testDataEngine.TestDataDomainModelStruct
		var tempTestDataAreaMap map[testDataEngine.TestDataAreaUuidType]*testDataEngine.TestDataAreaStruct
		var tempTestDataArea testDataEngine.TestDataAreaStruct
		var tempTestDataValuesForRowMap map[testDataEngine.TestDataPointRowUuidType]*[]*testDataEngine.TestDataPointValueStruct

		tempTestDataModelMap = *testDataModel.TestDataModelMap
		tempTestDataDomainModel = *tempTestDataModelMap[testDataDomainUuid]
		tempTestDataAreaMap = *tempTestDataDomainModel.TestDataAreasMap
		tempTestDataArea = *tempTestDataAreaMap[testDataAreaUuid]
		tempTestDataValuesForRowMap = *tempTestDataArea.TestDataValuesForRowMap

		var tempTestDataPointValueSlice []*testDataEngine.TestDataPointValueStruct

		//var tempTestDataPointValueSlice *[]*TestDataPointValueStruct

		var allTestDataPointRowsUuid []testDataEngine.TestDataPointRowUuidType

		// Loop all TestData and extract all rows
		for tempTestDataPointRowUuid, _ := range tempTestDataValuesForRowMap {
			allTestDataPointRowsUuid = append(allTestDataPointRowsUuid, tempTestDataPointRowUuid)
		}

		searchResult = allTestDataPointRowsUuid

		// Loop all Columns and extract selected checkboxes in the CheckGroups
		for _, testDataValueSelection := range testDataValueSelections {

			// Extract the Selected CheckBoxes
			var selectedCheckBoxes []string
			selectedCheckBoxes = testDataValueSelection.testDataCheckGroup.Selected

			// Extract 'TestDataPointRowUuid' for the Selected CheckBox-value-rows
			var testDataPointRowUuidMap map[testDataEngine.TestDataValueType]*[]testDataEngine.TestDataPointRowUuidType
			testDataPointRowUuidMap = *testDataValueSelection.TestDataPointValueRowUuidMap

			var testDataPointRowsUuid []testDataEngine.TestDataPointRowUuidType

			for _, selectedCheckBox := range selectedCheckBoxes {
				tempTestDataPointRowsUuid, _ := testDataPointRowUuidMap[testDataEngine.TestDataValueType(selectedCheckBox)]

				testDataPointRowsUuid = append(testDataPointRowsUuid, *tempTestDataPointRowsUuid...)

			}

			// Intersect with full TestDataSet to minimize the rows
			if len(testDataPointRowsUuid) != 0 {

				searchResult = testDataPointIntersectionOfTwoSlices(allTestDataPointRowsUuid, testDataPointRowsUuid)

			}
		}

		// Convert all DataPoints in SearchResult to be used in Available TestDataPoints-list based on already selected datapoints
		var tempTestDataValueName testDataEngine.TestDataValueNameType
		var tempTestDataPointRowUuid testDataEngine.TestDataPointRowUuidType
		var existInSelectedPoints bool
		var tempMapForSearchResultDataPoints map[testDataEngine.TestDataValueNameType]testDataEngine.DataPointTypeForGroupsStruct
		tempMapForSearchResultDataPoints = make(map[testDataEngine.TestDataValueNameType]testDataEngine.DataPointTypeForGroupsStruct)
		allPointsAvailable = nil

		for _, testDataPointRowUuid := range searchResult {

			tempTestDataPointValueSlice = *tempTestDataValuesForRowMap[testDataPointRowUuid]

			// Get the TestDataValueName
			tempTestDataValueName = tempTestDataPointValueSlice[0].TestDataValueName

			// Get the TestDataPointRowUuid
			tempTestDataPointRowUuid = tempTestDataPointValueSlice[0].TestDataPointRowUuid

			// Check if RowUuid already exists in SelectedDataPoints-list
			existInSelectedPoints = false
			if len(allSelectedPoints) != 0 {
				for _, selectedPoint := range allSelectedPoints {

					_, existInSelectedPoints = selectedPoint.AvailableTestDataPointUuidMap[tempTestDataPointRowUuid]

					// If the row already exist then exit for-loop
					if existInSelectedPoints == true {
						break
					}
				}
			}

			// Add the 'TestDataPointRowUuid' to inner map in 'searchResultDataPoint' if it doesn't already exist in  SelectedDataPoints-list
			if existInSelectedPoints == false {
				// Doesn't exist in Selected Points

				// Create the DataPoint from the SerachResult
				var searchResultDataPoint testDataEngine.DataPointTypeForGroupsStruct

				// Try to find the DataPoint in the Map based on 'tempTestDataValueName'
				searchResultDataPoint, existInMap = tempMapForSearchResultDataPoints[tempTestDataValueName]
				if existInMap == false {
					// It doesn't exist so create the 'searchResultDataPoint'
					searchResultDataPoint = testDataEngine.DataPointTypeForGroupsStruct{
						TestDataDomainUuid:            tempTestDataPointValueSlice[0].TestDataDomainUuid,
						TestDataDomainName:            tempTestDataPointValueSlice[0].TestDataDomainName,
						TestDataAreaUuid:              tempTestDataPointValueSlice[0].TestDataAreaUuid,
						TestDataAreaName:              tempTestDataPointValueSlice[0].TestDataAreaName,
						TestDataPointName:             tempTestDataValueName,
						SearchResultDataPointUuidMap:  nil,
						AvailableTestDataPointUuidMap: make(map[testDataEngine.TestDataPointRowUuidType]testDataEngine.TestDataPointRowUuidStruct),
						SelectedTestDataPointUuidMap:  make(map[testDataEngine.TestDataPointRowUuidType]testDataEngine.TestDataPointRowUuidStruct),
					}
				}

				// Create the 'testDataPointRowUuidObject' to be added to the Map
				var testDataPointRowUuidObject testDataEngine.TestDataPointRowUuidStruct
				testDataPointRowUuidObject = testDataEngine.TestDataPointRowUuidStruct{
					TestDataPointRowUuid:          tempTestDataPointValueSlice[0].TestDataPointRowUuid,
					TestDataPointRowValuesSummary: "",
				}

				searchResultDataPoint.AvailableTestDataPointUuidMap[tempTestDataPointRowUuid] = testDataPointRowUuidObject

				// Add the 'searchResultDataPoint' back to the temporary map for SearchResultDataPoints
				tempMapForSearchResultDataPoints[tempTestDataValueName] = searchResultDataPoint

			} else {
				// Exist in Selected Points

				// Create the DataPoint from the SerachResult
				var searchResultDataPoint testDataEngine.DataPointTypeForGroupsStruct

				// Try to find the DataPoint in the Map based on 'tempTestDataValueName'
				searchResultDataPoint, existInMap = tempMapForSearchResultDataPoints[tempTestDataValueName]
				if existInMap == false {
					// It doesn't exist so create the 'searchResultDataPoint'
					searchResultDataPoint = testDataEngine.DataPointTypeForGroupsStruct{
						TestDataDomainUuid:            tempTestDataPointValueSlice[0].TestDataDomainUuid,
						TestDataDomainName:            tempTestDataPointValueSlice[0].TestDataDomainName,
						TestDataAreaUuid:              tempTestDataPointValueSlice[0].TestDataAreaUuid,
						TestDataAreaName:              tempTestDataPointValueSlice[0].TestDataAreaName,
						TestDataPointName:             tempTestDataValueName,
						SearchResultDataPointUuidMap:  nil,
						AvailableTestDataPointUuidMap: make(map[testDataEngine.TestDataPointRowUuidType]testDataEngine.TestDataPointRowUuidStruct),
						SelectedTestDataPointUuidMap:  make(map[testDataEngine.TestDataPointRowUuidType]testDataEngine.TestDataPointRowUuidStruct),
					}
				}

				// Create the 'testDataPointRowUuidObject' to be added to the Map
				var testDataPointRowUuidObject testDataEngine.TestDataPointRowUuidStruct
				testDataPointRowUuidObject = testDataEngine.TestDataPointRowUuidStruct{
					TestDataPointRowUuid:          tempTestDataPointValueSlice[0].TestDataPointRowUuid,
					TestDataPointRowValuesSummary: "",
				}

				searchResultDataPoint.SelectedTestDataPointUuidMap[tempTestDataPointRowUuid] = testDataPointRowUuidObject

				// Add the 'searchResultDataPoint' back to the temporary map for SearchResultDataPoints
				tempMapForSearchResultDataPoints[tempTestDataValueName] = searchResultDataPoint
			}
		}

		// Create temporary slice to sort
		var allPointsAvailableToBeSorted []testDataEngine.DataPointTypeForGroupsStruct
		// Create the list that holds all points that are available to chose from
		for _, point := range tempMapForSearchResultDataPoints {

			allPointsAvailableToBeSorted = append(allPointsAvailableToBeSorted, point)
		}

		// Sort the slice with DataPoints
		allPointsAvailableToBeSorted = sortDataPointsList(allPointsAvailableToBeSorted)

		// copy back from sorted slice
		allPointsAvailable = allPointsAvailableToBeSorted

		// Refresh the List-widget
		allAvailablePointsList.Refresh()

	})

	// Create Clear checkboxes-button
	clearTestDataFilterCheckBoxesButton = widget.NewButton("Clear checkboxes", func() {

		var selected []string

		// Loop all Columns and clear all checkboxes in the CheckGroups
		for _, testDataValueSelection := range testDataValueSelections {

			testDataValueSelection.testDataCheckGroup.SetSelected(selected)

		}

	})

	// Create the container for the Search- and Clear-buttons
	searchAndClearButtonsContainer = container.NewHBox(searchTestDataButton, clearTestDataFilterCheckBoxesButton)

	/*
		// Convert into all 'TestDataValueName' in []TestDataPointRowUuidType to be used in Available TestDataPoints-list
		// Slices used to keep track of filtered DataPoints
		var filteredTestDataPoints []dataPointTypeForGroupsStruct
		var tempTestDataPointRowUuidSliceInMap []TestDataPointRowUuidType
		filteredTestDataPoints = nil
		var tempTestDataValueName string
		tempTestDataValuesForRowMap := *tempTestDataValuesForRowMapPtr
		for _, testDataPointRowUuid := range searchResult {

			tempTestDataPointValueSlicePtr, _ := tempTestDataValuesForRowMap[testDataPointRowUuid]
			tempTestDataPointValueSlice := *tempTestDataPointValueSlicePtr

			tempTestDataValueName = string(tempTestDataPointValueSlice[0].TestDataValueName)

			tempTestDataPointRowUuidSliceInMap, _ = tempTestDataValueNameToRowUuidMap[TestDataValueNameType(tempTestDataValueName)]
			tempTestDataPointRowUuidSliceInMap = append(tempTestDataPointRowUuidSliceInMap, testDataPointRowUuid)
			tempTestDataValueNameToRowUuidMap[TestDataValueNameType(tempTestDataValueName)] = tempTestDataPointRowUuidSliceInMap
		}

		for tempTestDataValueNameInMap, tempTestDataPointRowUuidSliceFromMap := range tempTestDataValueNameToRowUuidMap {

			// Create a filtered TestDataPoint
			var filteredTestDataPoint dataPointTypeForGroupsStruct
			filteredTestDataPoint = dataPointTypeForGroupsStruct{
				testDataDomainUuid:            "",
				testDataDomainName:            "",
				testDataAreaUuid:              "",
				testDataAreaName:              "",
				testDataPointName:             tempTestDataValueNameInMap,
				availableTestDataPointUuidMap: nil,
			}

			// Add the 'TestDataPointUuid's' to the filtered TestDataPoint
			for _, tempTestDataPointUuid := range tempTestDataPointRowUuidSliceFromMap {
				filteredTestDataPoint.availableTestDataPointUuidMap[tempTestDataPointUuid] = tempTestDataPointUuid
			}

		}

		// Create the list that holds all points that are available to chose from
		allPointsAvailable = nil
		var rowUuidExistInSelectedPoints bool
		var nameExistInSelectedPoints bool
		var nameExistInAvailablePoints bool
		var tempSelectedTestDataPointUuid TestDataPointRowUuidType
		var availablePointsIndex int

		for _, filteredPoint := range filteredTestDataPoints {

			// Add it to the list of available points, if it doesn't exist in the Selected-List
			if len(allSelectedPoints) == 0 {
				allSelectedPoints = append(allSelectedPoints, filteredPoint)
			} else {

				// Clear flags for of TestDataPointName and TestDataPointRowUuid exist in SelectedPoints
				nameExistInSelectedPoints = false
				rowUuidExistInSelectedPoints = false

				// Clear the flag if the TestDataPointName exist in allPointsAvailable-slice
				nameExistInAvailablePoints = false

				for _, selectedPoint := range allSelectedPoints {

					if selectedPoint.testDataPointName == filteredPoint.testDataPointName {

						nameExistInSelectedPoints = true

						// Check if row-UUID exist in SelectedPoint
						for _, selectedTestDataPointUuid := range selectedPoint.selectedTestDataPointUuidMap {
							_, existInMap = selectedPoint.selectedTestDataPointUuidMap[selectedTestDataPointUuid]

							// Exit for-loop if the TestDataPointUuid exist
							if existInMap == false {
								tempSelectedTestDataPointUuid = selectedTestDataPointUuid
								rowUuidExistInSelectedPoints = true
								break
							}
						}

						// If the TestDataPointUuid doesn't exist in SelectedPoints then add to the Available TestDataPoints
						if rowUuidExistInSelectedPoints == false {

							// Check if the TestDataPointName exist in the allPointsAvailable slice
							for tempAvailablePointsIndex, availablePoint := range allPointsAvailable {

								if availablePoint.testDataPointName == filteredPoint.testDataPointName {
									nameExistInAvailablePoints = true
									availablePointsIndex = tempAvailablePointsIndex
									break
								}
							}

							// If TestDataPointName exist in the allPointsAvailable-slice, then add it to the TestDataPoint in allPointsAvailable-slice
							if nameExistInAvailablePoints == true {
								existingFilteredPoint := allPointsAvailable[availablePointsIndex]
								existingFilteredPoint.availableTestDataPointUuidMap[tempSelectedTestDataPointUuid] = tempSelectedTestDataPointUuid
								allPointsAvailable[availablePointsIndex] = existingFilteredPoint

							} else {
								// The TestDataPointName didn't exist so add the full TestDataPoint
								allPointsAvailable = append(allPointsAvailable, filteredPoint)

							}

						}

						// Exit the for-loop if the TestDataPointName exist SelectedPoints
						if nameExistInSelectedPoints == true {
							break
						}
					}
				}
			}
		}

		// Custom sort: we sort by splitting each string into parts and comparing the parts
		sort.Slice(tempAllPointsAvailable, func(i, j int) bool {
			// Split both strings by '/'
			partsI := strings.Split(string(tempAllPointsAvailable[i].testDataPointName), "/")
			partsJ := strings.Split(string(tempAllPointsAvailable[j].testDataPointName), "/")

			// Compare each part; the first non-equal part determines the order
			for k := 0; k < len(partsI) && k < len(partsJ); k++ {
				if partsI[k] != partsJ[k] {
					return partsI[k] < partsJ[k]
				}
			}

			// If all compared parts are equal, but one slice is shorter, it comes first
			return len(partsI) < len(partsJ)
		})

		// Write back to original from local copy of 'allPointsAvailable'
		allPointsAvailable = &tempAllPointsAvailable

		// Refresh the List-widget
		allAvailablePointsList.Refresh()

		return testDataSelectionsContainer, searchAndClearButtonsContainer
	*/

}

// Sort a slice with DataPoints
func sortDataPointsList(dataPointListToBeSorted []testDataEngine.DataPointTypeForGroupsStruct) []testDataEngine.DataPointTypeForGroupsStruct {

	// Custom sort: we sort by splitting each string into parts and comparing the parts
	sort.Slice(dataPointListToBeSorted, func(i, j int) bool {
		// Split both strings by '/'
		partsI := strings.Split(string(dataPointListToBeSorted[i].TestDataPointName), "/")
		partsJ := strings.Split(string(dataPointListToBeSorted[j].TestDataPointName), "/")

		// Compare each part; the first non-equal part determines the order
		for k := 0; k < len(partsI) && k < len(partsJ); k++ {
			if partsI[k] != partsJ[k] {
				return partsI[k] < partsJ[k]
			}
		}

		// If all compared parts are equal, but one slice is shorter, it comes first
		return len(partsI) < len(partsJ)
	})

	return dataPointListToBeSorted

}
