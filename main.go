package main

import (
	"ImportFilesFromGithub/fileViewer"
	"ImportFilesFromGithub/importFilesFromGitHub"
	"ImportFilesFromGithub/luaScriptEngine"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

func main() {

	var err error
	// Initiate Lua Script Engine
	err = luaScriptEngine.InitiateLuaScriptEngine([][]byte{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer luaScriptEngine.CloseDownLuaScriptEngine()

	repoOwner := "jlambert68"    // Replace with the repository owner's username
	repoName := "FenixTesterGui" // Replace with the repository name
	repoPath := ""               // Replace with the path in the repository, if any

	originalApiUrl := "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/contents" + repoPath

	var myApp fyne.App
	myApp = app.New()

	var myMainWindow fyne.Window
	myMainWindow = myApp.NewWindow("This is Fenix Main Window")

	// Set initial size of the window
	myMainWindow.Resize(fyne.NewSize(400, 500))

	var responseChannel chan bool
	responseChannel = make(chan bool)
	var selectedFiles *[]importFilesFromGitHub.GitHubFile
	githubFilesImporterButton := widget.NewButton("Import files from GitHub", func() {
		myMainWindow.Hide()
		selectedFiles = importFilesFromGitHub.InitiateImportFilesFromGitHubWindow(originalApiUrl, myMainWindow, myApp, &responseChannel)
		fmt.Println(selectedFiles)
	})

	filesViewerButton := widget.NewButton("View imported files", func() {
		myMainWindow.Hide()
		fileViewer.InitiateFileViewer(myMainWindow, myApp, selectedFiles)

	})

	inputText := "This is {{yellow}} text and this is {{also yellow}} and normal again."
	var tempRichText *widget.RichText
	tempRichText = parseAndFormatText(inputText)

	buttonContainer := container.NewVBox(githubFilesImporterButton, filesViewerButton)

	myContainer := container.NewBorder(buttonContainer, nil, nil, nil, tempRichText)
	myMainWindow.SetContent(myContainer)

	go func() {

		for {
			var responseValue bool
			responseValue = <-responseChannel
			fmt.Println(responseValue)

			var files []importFilesFromGitHub.GitHubFile
			files = *selectedFiles

			if len(files) > 0 {
				var fileContent string
				var file importFilesFromGitHub.GitHubFile

				file = files[0]
				fileContent = file.FileContetAsString

				myContainerObjects := myContainer.Objects
				for index, object := range myContainerObjects {
					if object == tempRichText {
						tempRichText = parseAndFormatText(fileContent)
						myContainerObjects[index] = tempRichText
						myContainer.Refresh()
						break
					}
				}

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
