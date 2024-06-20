package fileViewer

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"ImportFilesFromGithub/luaScriptEngine"
	"ImportFilesFromGithub/testDataSelector"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"strconv"
	"strings"
)

func InitiateFileViewer(
	mainWindow fyne.Window,
	myApp fyne.App,
	importedFilesPtr *[]importFilesFromGitHub.GitHubFile) {

	// Disable the main window
	mainWindow.Hide()

	// The Select-items for Groups ans TestDataPoints for a Group
	var testDataPointGroupsSelect *widget.Select
	var testDataPointGroupsSelectSelected string
	var testDataPointsForAGroupSelect *widget.Select
	var testDataPointForAGroupSelectSelected string
	var testDataRowsForTestDataPointsSelect *widget.Select
	var testDataRowsForTestDataPointsSelectSelected string

	// The slices for Groups, TestDataPoints for a Group and the specific TestDataRows for a TestDataPoint
	var testDataPointGroups []string
	var testDataPointsForAGroup []string
	var testDataRowsForATestDataPoint []string

	// Store reference to Fenix Main Window
	fenixMainWindow = mainWindow

	// Create the window for GitHub files
	fileViewerWindow = myApp.NewWindow("Imported Files Viewer")
	// Set initial size of the window
	fileViewerWindow.Resize(fyne.NewSize(800, 700))

	var leftContainer *fyne.Container
	var rightContainer *fyne.Container

	// Extract filenames for the fileSelectorDropdown
	var fileNames []string
	for _, file := range *importedFilesPtr {
		fileNames = append(fileNames, file.Name)
	}

	// Create UI component for 'fileSelectorDropdown'
	fileSelectorDropdown := widget.NewSelect(fileNames, nil)
	urlLabel := widget.NewLabel("")
	var richText *widget.RichText
	richText = &widget.RichText{
		BaseWidget: widget.BaseWidget{},
		Segments:   nil,
		Wrapping:   0,
		Scroll:     0,
		Truncation: 0,
	}
	var richTextWithValues *widget.RichText
	richTextWithValues = &widget.RichText{
		BaseWidget: widget.BaseWidget{},
		Segments:   nil,
		Wrapping:   0,
		Scroll:     0,
		Truncation: 0,
	}

	var testDataPointValues map[string]string // map[TestDataColumnDataNameType]TestDataValueType

	// Set the fileSelectorDropdown change handler
	var selectedFile string
	fileSelectorDropdown.OnChanged = func(selected string) {

		selectedFile = selected

		testDataPointValues = testDataSelector.GetTestDataPointValues(
			testDataPointForAGroupSelectSelected,
			testDataRowsForTestDataPointsSelectSelected)

		for _, file := range *importedFilesPtr {
			if file.Name == selected {
				urlLabel.SetText(file.URL)

				myContainerObjects := leftContainer.Objects
				for index, object := range myContainerObjects {
					if object == richText {
						richText, _ = parseAndFormatText(file.FileContentAsString, &testDataPointValues)
						myContainerObjects[index] = richText
						leftContainer.Refresh()
					}
				}

				myContainerObjects = rightContainer.Objects
				for index, object := range myContainerObjects {
					if object == richTextWithValues {
						_, richTextWithValues = parseAndFormatText(file.FileContentAsString, &testDataPointValues)
						myContainerObjects[index] = richTextWithValues
						rightContainer.Refresh()
					}
				}

				break
			}
		}
	}

	// Create function that converts a GroupSlice into a string slice
	getTestGroupsFromTestDataEngineFunction := func() []string {

		testDataPointGroups = testDataSelector.ListTestDataGroups()

		return testDataPointGroups
	}

	// Create function that converts a TestDataPointsSlice into a string slice
	testDataPointsToStringSliceFunction := func(testDataGroup string) []string {

		if testDataGroup == "" {
			return []string{}
		}

		testDataPointsForAGroup = testDataSelector.ListTestDataGroupPointsForAGroup(testDataGroup)

		return testDataPointsForAGroup
	}

	// Create function that converts a slice with the specific TestDataPoints into a string slice
	testDataRowSliceToStringSliceFunction := func(testDataGroup string, testDataGroupPoint string) []string {

		if testDataGroup == "" || testDataGroupPoint == "" {
			return []string{}
		}

		testDataRowsForATestDataPoint = testDataSelector.ListTestDataRowsForAGroupPoint(testDataGroup, testDataGroupPoint)

		return testDataRowsForATestDataPoint
	}

	// Create the Group dropdown - <Name of the group>
	testDataPointGroupsSelect = widget.NewSelect(getTestGroupsFromTestDataEngineFunction(), func(selected string) {

		testDataPointGroupsSelectSelected = selected

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		testDataPointsForAGroupSelect.SetOptions(testDataPointsToStringSliceFunction(selected))
		testDataPointsForAGroupSelect.Refresh()

		// UnSelect in DropDown- and List for TestDataPoints
		testDataPointsForAGroupSelect.ClearSelected()

	})

	// Create the Groups TestDataPoints dropdown - <Sub Custody/Main TestData Area/SEK/AccTest/SE/CRDT/CH/Switzerland/BBH/EUR/EUR/SEK>
	testDataPointsForAGroupSelect = widget.NewSelect(testDataPointsToStringSliceFunction(testDataPointGroupsSelectSelected), func(selected string) {

		testDataPointForAGroupSelectSelected = selected

		// Select the correct TestDataPoint in the dropdown for TestDataPoints
		testDataRowsForTestDataPointsSelect.SetOptions(testDataRowSliceToStringSliceFunction(testDataPointGroupsSelect.Selected, selected))
		testDataRowsForTestDataPointsSelect.Refresh()

		// UnSelect in DropDown- and List for Specific TestDataPoints
		testDataRowsForTestDataPointsSelect.ClearSelected()

	})

	// Create the Groups Specific TestDataPoint dropdown - <All the specific values>
	testDataRowsForTestDataPointsSelect = widget.NewSelect(testDataRowSliceToStringSliceFunction(testDataPointGroupsSelectSelected, testDataPointForAGroupSelectSelected), func(selected string) {

		testDataRowsForTestDataPointsSelectSelected = selected

		fileSelectorDropdown.SetSelected(selectedFile)

	})

	// Create UI component for 'TestDataGroupPointSelector'

	topContainer := container.NewVBox(fileSelectorDropdown, urlLabel, testDataPointGroupsSelect, testDataPointsForAGroupSelect, testDataRowsForTestDataPointsSelect)

	// Placeholder for rightContainer - add your form view here
	rightContainer = container.NewBorder(nil, nil, nil, nil, richTextWithValues)

	leftContainer = container.NewBorder(nil, nil, nil, nil, richText)

	// Create split container
	split := container.NewHSplit(leftContainer, rightContainer)
	split.Offset = 0.5 // Adjust as needed

	fullContentContainer := container.NewBorder(topContainer, nil, nil, nil, split)

	fileViewerWindow.SetContent(fullContentContainer)

	// Set the callback function for window close event to show the Main window again
	fileViewerWindow.SetOnClosed(func() {
		fenixMainWindow.Show()
	})

	// Show the File Viewe Window
	fileViewerWindow.Show()
}

func match(text string) (mainScriptInputSlice []interface{}, err error) {

	var arrayIndexSlice []interface{}
	var functionArgumentSlice []interface{}

	//regExPattern := `\{\{([a-zA-Z0-9_.]+)(?:\[(\d*(?:,\s*\d*)*)\])?\((.*?)\)\}(?:\((true|false)(?:,\s*(\d+))?\))?\}`
	regExPattern := `\{\{([a-zA-Z0-9_.]+)(?:\[([-+]?\d*(?:,\s*[-+]?\d*)*)\])?\((.*?)\)\}(?:\((true|false)(?:,\s*(\d+))?\))?\}`
	/*
		Explanation of Each Part
		\{\{: Matches literal {{.
		([a-zA-Z0-9_.]+): Matches and captures the function name consisting of one or more alphanumeric characters, underscores, or dots.
		(?:\[([-+]?\d*(?:,\s*[-+]?\d*)*)\])?: Non-capturing group for indices, which can now include negative or positive numbers. It's optional.
		[-+]?\d*: Matches an optional sign (+ or -), followed by any digits.
		(?:,\s*[-+]?\d*)*: Matches zero or more repetitions of a comma, optional spaces, an optional sign, and digits, allowing for lists of indices.
		\((.*?)\): Captures the arguments within parentheses. .*? is used for lazy matching to stop at the first ).
		\}: Matches literal }.
		(?:\((true|false)(?:,\s*(\d+))?\))?: Optional non-capturing group for additional parameters like boolean values and numbers, typically used for configurations or flags.
	*/
	re := regexp.MustCompile(regExPattern)

	matches := re.FindStringSubmatch(text)

	if len(matches) >= 2 {
		placeholder := matches[0]
		functionName := matches[1]
		arrayIndexes := matches[2]
		functionArgs := matches[3]
		useEntropyFromTestCaseExecutionUuid := matches[4]
		addExtraEntropyValue := matches[5]

		// Add 'placeholder' to 'mainScriptInputSlice'
		mainScriptInputSlice = append(mainScriptInputSlice, placeholder)

		functionName = strings.ReplaceAll(functionName, ".", "_")
		mainScriptInputSlice = append(mainScriptInputSlice, functionName)

		// Split the array indexes into a slices
		indexes := strings.Split(arrayIndexes, ",")

		// Create a ArrayIndex-array as s '[]interface{}'
		var indexAsInt int
		for i, index := range indexes {
			indexes[i] = strings.TrimSpace(index)

			// Only convert when there is some value
			if len(indexes[i]) > 0 {
				indexAsInt, err = strconv.Atoi(indexes[i])
				if err != nil {
					err = errors.New(fmt.Sprintf("Couldn't convert array index '%s' in '%s' to an integer. Placeholder = '%s'", indexes[i], indexes, placeholder))

					return nil, err

				}

				arrayIndexSlice = append(arrayIndexSlice, indexAsInt)
			}
		}

		// Add the FunctionArguments-array to the main input array
		mainScriptInputSlice = append(mainScriptInputSlice, arrayIndexSlice)

		// Split the function arguments into a slice
		args := strings.Split(functionArgs, ",")

		// Create a FunctionArguments-array as a '[]interface{}'
		//var argAsInt int
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)

			/*
				// Only convert when there is some value
				if len(args[i]) > 0 {
					argAsInt, err = strconv.Atoi(args[i])
					if err != nil {
						err = errors.New(fmt.Sprintf("Couldn't convert parameter '%s' in '%s' to an integer. Placeholder = '%s'", args[i], args, placeholder))

						return nil, err

					}

					functionArgumentSlice = append(functionArgumentSlice, argAsInt)
				}

			*/
			functionArgumentSlice = append(functionArgumentSlice, args[i]) //argAsInt)
		}

		// Add the FunctionArguments-array to the main input array
		mainScriptInputSlice = append(mainScriptInputSlice, functionArgumentSlice)

		// When there is no boolean for 'useEntropyFromTestCaseExecutionUuid' then always use 'true'
		if len(useEntropyFromTestCaseExecutionUuid) == 0 {

			mainScriptInputSlice = append(mainScriptInputSlice, true)

		} else {

			tempBoolean, _ := strconv.ParseBool(useEntropyFromTestCaseExecutionUuid)
			mainScriptInputSlice = append(mainScriptInputSlice, tempBoolean)
		}

		// When there is no value for 'addExtraEntropyValue' then always use '0'
		if len(addExtraEntropyValue) == 0 {

			mainScriptInputSlice = append(mainScriptInputSlice, uint64(0))

		} else {

			var tempExtraEntropy uint64
			tempExtraEntropy, err = strconv.ParseUint(addExtraEntropyValue, 10, 32)
			if err != nil {
				return nil, err
			}

			mainScriptInputSlice = append(mainScriptInputSlice, tempExtraEntropy)
		}

	} else {
		fmt.Println("No match found for:", text)
		err = errors.New(fmt.Sprintf("No match found for '%s'", text))
	}

	return mainScriptInputSlice, err
}

func parseAndFormatText(inputText string, testDataPointValuesPtr *map[string]string) (
	tempRichText *widget.RichText,
	tempRichTextWithValues *widget.RichText) {

	var testDataPointValues map[string]string
	testDataPointValues = *testDataPointValuesPtr

	var segments []widget.RichTextSegment
	var segmentsWithValues []widget.RichTextSegment

	var currentText string

	for len(inputText) > 0 {
		startIndex := strings.Index(inputText, "{{")
		endIndex := strings.Index(inputText, "}}")

		if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
			// Add the text before {{
			if startIndex > 0 {
				currentText = inputText[:startIndex]
				segments = append(segments,
					&widget.TextSegment{
						Text: currentText,
						Style: widget.RichTextStyle{
							Inline: true,
						}})

				segmentsWithValues = append(segmentsWithValues,
					&widget.TextSegment{
						Text: currentText,
						Style: widget.RichTextStyle{
							Inline: true,
						}})
			}

			// Add the styled text between {{ and }}
			currentText = inputText[startIndex : endIndex+2] // +2 to include the closing braces

			var newTextFromScriptEngine string
			if strings.Contains(currentText, ".TestData.") == true {

				// Filter out Start and End '{{' and '}}'

				var testDataToReplace string
				testDataToReplace = currentText[2 : len(currentText)-2]

				var existInMap bool

				// Substring to find
				substr := ".TestData."

				// Find the position of ".TestData."
				pos := strings.Index(testDataToReplace, substr)

				// Extract the text to the right of ".TestData." if it exists
				var testDataColumnDataName string
				if pos != -1 {
					// Adjust position to skip ".TestData."
					start := pos + len(substr)
					if start < len(testDataToReplace) {
						testDataColumnDataName = testDataToReplace[start:]
					}
				} else {
					testDataColumnDataName = ""
				}

				if testDataColumnDataName != "" {
					newTextFromScriptEngine, existInMap = testDataPointValues[testDataColumnDataName]
					if existInMap == false {
						newTextFromScriptEngine = fmt.Sprintf("TestDataColumnDataName '%s' does not exist in the TestDataMap", testDataColumnDataName)
					}

				} else {
					newTextFromScriptEngine = currentText + " - is not a correct TestData-reference"
				}

			} else {
				functionValueSlice, err := match(currentText)
				if err == nil {
					//newTextFromScriptEngine = tengoScriptExecuter.ExecuteScripte(functionValueSlice)
					newTextFromScriptEngine = luaScriptEngine.ExecuteLuaScriptBasedOnPlaceholder(functionValueSlice, testCaseExecutionUuid)

				} else {
					newTextFromScriptEngine = err.Error()
				}
			}
			segments = append(segments, &widget.TextSegment{
				Text: currentText,
				Style: widget.RichTextStyle{
					Inline:    true,
					TextStyle: fyne.TextStyle{Bold: true},
				},
			})

			segmentsWithValues = append(segmentsWithValues, &widget.TextSegment{
				Text: newTextFromScriptEngine,
				Style: widget.RichTextStyle{
					Inline:    true,
					TextStyle: fyne.TextStyle{Bold: true},
				},
			})

			// Move past this segment
			inputText = inputText[endIndex+2:]
		} else {
			// Add the remaining text, if any
			segments = append(segments, &widget.TextSegment{Text: inputText})
			segmentsWithValues = append(segmentsWithValues, &widget.TextSegment{Text: inputText})
			break
		}
	}

	tempRichText = &widget.RichText{
		BaseWidget: widget.BaseWidget{},
		Segments:   segments,
		Wrapping:   0,
		Scroll:     0,
		Truncation: 0,
	}

	tempRichTextWithValues = &widget.RichText{
		BaseWidget: widget.BaseWidget{},
		Segments:   segmentsWithValues,
		Wrapping:   0,
		Scroll:     0,
		Truncation: 0,
	}
	return tempRichText, tempRichTextWithValues
}
