package main

import (
	"ImportFilesFromGithub/fileViewer"
	"ImportFilesFromGithub/importFilesFromGitHub"
	"ImportFilesFromGithub/testDataEngine"
	"ImportFilesFromGithub/testDataSelector"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jlambert68/FenixScriptEngine/luaEngine"
	"log"
	"strings"
)

func main() {

	var err error
	// Initiate Lua Script Engine
	err = luaEngine.InitiateLuaScriptEngine([][]byte{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer luaEngine.CloseDownLuaScriptEngine()

	repoOwner := "jlambert68"    // Replace with the repository owner's username
	repoName := "FenixTesterGui" // Replace with the repository name
	repoPath := ""               // Replace with the path in the repository, if any

	originalApiUrl := "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/contents" + repoPath

	var myApp fyne.App
	myApp = app.New()

	var myMainWindow fyne.Window
	myMainWindow = myApp.NewWindow("This is Fenix Main Window")

	// Set initial size of the window
	myMainWindow.Resize(fyne.NewSize(800, 600))

	// Initiate testDataForGroupObject used for keeping Groups TestData separate in different TestCases
	var testDataForGroupObject *testDataEngine.TestDataForGroupObjectStruct
	testDataForGroupObject = &testDataEngine.TestDataForGroupObjectStruct{}

	//Add TestData to global TestDataModel
	testDataSelector.ImportTestData()

	var fileTable *widget.Table

	var responseChannel chan bool
	responseChannel = make(chan bool)
	var selectedFiles *[]importFilesFromGitHub.GitHubFile
	selectedFiles = &[]importFilesFromGitHub.GitHubFile{}
	githubFilesImporterButton := widget.NewButton("Import files from GitHub", func() {
		myMainWindow.Hide()
		var tempSelectedFiles []importFilesFromGitHub.GitHubFile
		tempSelectedFiles = *selectedFiles
		selectedFiles = importFilesFromGitHub.InitiateImportFilesFromGitHubWindow(originalApiUrl, myMainWindow, myApp, &responseChannel, tempSelectedFiles)
	})

	filesViewerButton := widget.NewButton("View imported files", func() {
		myMainWindow.Hide()
		fileViewer.InitiateFileViewer(myMainWindow, myApp, selectedFiles, testDataForGroupObject)

	})

	// Create the button for handling TestDataGroups and TestDataPoints
	openTestDataWindow := widget.NewButton("Open TestData Window", func() {
		myMainWindow.Hide()
		testDataSelector.MainTestDataSelector(
			myApp,
			myMainWindow,
			testDataForGroupObject)
	})

	//inputText := "This is {{bold}} text and this is {{also bold}} and this normal again."
	//var tempRichText *widget.RichText
	//tempRichText = parseAndFormatText(inputText)

	// Correctly initialize the fileList as a new table
	fileTable = widget.NewTable(
		func() (int, int) { return 0, 2 }, // Start with zero rows, 2 columns
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Create cells as labels
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// This should be filled when updating the table
		},
	)

	buttonContainer := container.NewVBox(githubFilesImporterButton, filesViewerButton, openTestDataWindow)

	var files []importFilesFromGitHub.GitHubFile
	files = []importFilesFromGitHub.GitHubFile{}
	updateTable(fileTable, files)

	myContainer := container.NewBorder(buttonContainer, nil, nil, nil, fileTable)
	myMainWindow.SetContent(myContainer)

	go func() {

		//var responseValue bool

		for {

			_ = <-responseChannel

			files = *selectedFiles

			if len(files) > 0 {
				//var fileContent string
				//var file importFilesFromGitHub.GitHubFile

				updateTable(fileTable, files)
				importFilesFromGitHub.UpdateSelectedFilesTable()
				/*
					file = files[0]
					fileContent = file.FileContentAsString

					myContainerObjects := myContainer.Objects
					for index, object := range myContainerObjects {
						if object == tempRichText {
							tempRichText = parseAndFormatText(fileContent)
							myContainerObjects[index] = tempRichText
							myContainer.Refresh()
							break
						}
					}

				*/

			}
		}
	}()

	// Define a slice of interfaces
	//var mySlice []interface{}

	// Add an integer to the slice
	//mySlice = append(mySlice, 1)
	//tengoScriptExecuter.ExecuteScripte("SubCustody_RandomFloatValue", mySlice)

	myMainWindow.ShowAndRun()

}

func updateTable(fileList *widget.Table, files []importFilesFromGitHub.GitHubFile) {

	maxNameWidth := float32(150) // Start with a minimum width
	maxUrlWidth := float32(250)  // Start with a minimum width
	for _, file := range files {
		textNameWidth := fyne.MeasureText(file.Name, theme.TextSize(), fyne.TextStyle{}).Width
		textUrlWidth := fyne.MeasureText(file.URL, theme.TextSize(), fyne.TextStyle{}).Width
		if textNameWidth > maxNameWidth {
			maxNameWidth = textNameWidth
		}
		if textUrlWidth > maxUrlWidth {
			maxUrlWidth = textUrlWidth
		}
	}

	fileList.SetColumnWidth(0, maxNameWidth+theme.Padding()*4) // Add padding
	fileList.SetColumnWidth(1, maxUrlWidth+theme.Padding()*4)  // Path column width can be static or calculated similarly

	//fileList.SetColumnWidth(0, 150) // Set width of file name column
	//fileList.SetColumnWidth(1, 400) // Set width of path column
	fileList.Length = func() (int, int) {
		return len(files), 2
	}
	fileList.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	fileList.UpdateCell = func(id widget.TableCellID, cell fyne.CanvasObject) {
		switch id.Col {
		case 0:
			cell.(*widget.Label).SetText(files[id.Row].Name)
		case 1:
			cell.(*widget.Label).SetText(files[id.Row].URL)
		}
	}
	fileList.Refresh()
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

	/*
		// Splitting the string at each "{{" and "}}"
		parts := strings.FieldsFunc(inputText, func(r rune) bool {
			return r == '{' || r == '}'
		})

		for i, part := range parts {
			var txt *widget.TextSegment

			// Handle text inside "{{...}}"
			if i%2 == 1 {
				txt = &widget.TextSegment{
					Text: "{{" + part + "}}",
					Style: widget.RichTextStyle{
						Inline:    true,
						TextStyle: fyne.TextStyle{Bold: true, Italic: true},
					},
				}
			} else {
				// Handle regular text
				txt = &widget.TextSegment{
					Text: part,
					Style: widget.RichTextStyle{
						Inline: true,
					},
				}
			}

			segments = append(segments, txt)
		}

	*/
	/*
		parts := strings.Split(inputText, "#")
		for i, part := range parts {

			//txt := canvas.NewText(part, color.Black)
			var txt *widget.TextSegment
			txt = &widget.TextSegment{
				Text: part,
				Style: widget.RichTextStyle{
					Inline: true,
				},
			}

			if i%2 == 1 { // Color odd parts (between '#')
				txt = &widget.TextSegment{
					Text: "#" + part + "#",
					Style: widget.RichTextStyle{
						Inline:    true,
						TextStyle: fyne.TextStyle{Bold: true, Italic: true},
					},
				}
			}
			segments = append(segments, txt)
		}

	*/

	tempRichText = &widget.RichText{
		BaseWidget: widget.BaseWidget{},
		Segments:   segments,
		Wrapping:   0,
		Scroll:     0,
		Truncation: 0,
	}
	return tempRichText
}
