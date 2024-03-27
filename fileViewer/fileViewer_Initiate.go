package fileViewer

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"ImportFilesFromGithub/tengoScriptExecuter"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"strconv"
	"strings"
)

func InitiateFileViewe(
	mainWindow fyne.Window,
	myApp fyne.App,
	importedFilesPtr *[]importFilesFromGitHub.GitHubFile) {

	// Disable the main window
	mainWindow.Hide()

	// Store reference to Fenix Main Window
	fenixMainWindow = mainWindow

	// Create the window for GitHub files
	fileViewerWindow = myApp.NewWindow("Imported Files Viewer")
	// Set initial size of the window
	fileViewerWindow.Resize(fyne.NewSize(400, 500))

	var leftContainer *fyne.Container
	var rightContainer *fyne.Container

	// Extract filenames for the dropdown
	var fileNames []string
	for _, file := range *importedFilesPtr {
		fileNames = append(fileNames, file.Name)
	}

	// Create UI components
	dropdown := widget.NewSelect(fileNames, nil)
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

	// Set the dropdown change handler
	dropdown.OnChanged = func(selected string) {
		for _, file := range *importedFilesPtr {
			if file.Name == selected {
				urlLabel.SetText(file.URL)

				myContainerObjects := leftContainer.Objects
				for index, object := range myContainerObjects {
					if object == richText {
						richText, _ = parseAndFormatText(file.FileContetAsString)
						myContainerObjects[index] = richText
						leftContainer.Refresh()
					}
				}

				myContainerObjects = rightContainer.Objects
				for index, object := range myContainerObjects {
					if object == richTextWithValues {
						_, richTextWithValues = parseAndFormatText(file.FileContetAsString)
						myContainerObjects[index] = richTextWithValues
						rightContainer.Refresh()
					}
				}

				break
			}
		}
	}

	topContainer := container.NewVBox(dropdown, urlLabel)

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

func match(text string) (functionName string, functionValueSlice []interface{}) {
	//text := "{{SubCustody.Today(1)}}"
	pattern := `\{\{([a-zA-Z0-9_.]+)\((([-?\d+],?\s*)+)\)\}\}`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)

	if len(matches) >= 3 {
		functionName = matches[1]
		functionArgs := matches[2]

		// Split the arguments into a slice
		args := strings.Split(functionArgs, ",")
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg) // Trim whitespace from each argument

			argAsInt, _ := strconv.Atoi(args[i])
			functionValueSlice = append(functionValueSlice, argAsInt)
		}

		fmt.Println("Function Name:", functionName)
		fmt.Println("Function Arguments:", args)
	} else {
		fmt.Println("No match found")
	}

	// Add an integer to the slice
	//functionValueSlice = append(functionValueSlice, functionValue)

	return functionName, functionValueSlice
}

func parseAndFormatText(inputText string) (
	tempRichText *widget.RichText,
	tempRichTextWithValues *widget.RichText) {

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
			functionName, functionValueSlice := match(currentText)
			newTextFromScriptEngine := tengoScriptExecuter.ExecuteScripte(functionName, functionValueSlice)

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
