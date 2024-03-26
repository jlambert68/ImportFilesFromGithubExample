package fileViewer

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	// Set the dropdown change handler
	dropdown.OnChanged = func(selected string) {
		for _, file := range *importedFilesPtr {
			if file.Name == selected {
				urlLabel.SetText(file.URL)

				myContainerObjects := leftContainer.Objects
				for index, object := range myContainerObjects {
					if object == richText {
						richText = parseAndFormatText(file.FileContetAsString)
						myContainerObjects[index] = richText
						leftContainer.Refresh()
						break
					}
				}

				break
			}
		}
	}

	// Placeholder for rightContainer - add your form view here
	rightContainer := widget.NewLabel("Right side content")

	leftContainer = container.NewBorder(container.NewVBox(dropdown, urlLabel), nil, nil, nil, richText)

	// Create split container
	split := container.NewHSplit(leftContainer, rightContainer)
	split.Offset = 0.5 // Adjust as needed

	fileViewerWindow.SetContent(split)

	// Set the callback function for window close event to show the Main window again
	fileViewerWindow.SetOnClosed(func() {
		fenixMainWindow.Show()
	})

	// Show the File Viewe Window
	fileViewerWindow.Show()
}

func parseAndFormatText(inputText string) (tempRichText *widget.RichText) {
	var segments []widget.RichTextSegment

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
			}

			// Add the styled text between {{ and }}
			currentText = inputText[startIndex : endIndex+2] // +2 to include the closing braces
			segments = append(segments, &widget.TextSegment{
				Text: currentText,
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
	return tempRichText
}
