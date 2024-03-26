package main

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

func main() {

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

	myButton := widget.NewButton("Import files from GitHub", func() {
		myMainWindow.Hide()
		_ = importFilesFromGitHub.InitiateImportFilesFromGitHubWindow(originalApiUrl, myMainWindow, myApp)
	})

	inputText := "This is {{yellow}} text and this is {{also yellow}} and normal again."

	// Create canvas.Text objects for different segments
	regularText := canvas.NewText("This is regular text. ", color.Black)
	yellowText := canvas.NewText("This is yellow text. ", color.NRGBA{R: 255, G: 255, B: 0, A: 255})
	moreText := canvas.NewText("And this is more regular text.", color.Black)

	regularText.TextSize = 18
	yellowText.TextSize = 14
	moreText.TextSize = 10
	/*
		rt := widget.NewRichText(
			&widget.TextSegment{
				Text: "A Title\r\n",
				Style: widget.RichTextStyle{
					Inline:    true,
					TextStyle: fyne.TextStyle{Bold: true},
				},
			},
			&widget.TextSegment{
				Text: "Here is some text,\r\nwith a line break",
			},
		)

		btn := widget.NewButton("Refresh", func() {
			rt.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text: "A Title\n",
					Style: widget.RichTextStyle{
						Inline:    true,
						TextStyle: fyne.TextStyle{Bold: true},
						ColorName: fyne.ThemeColorName("ColorNameShadow"),
					},
				},
				&widget.TextSegment{
					Text: "Here is some text,\n\n\n\nwith a line break",
				},
			}
			rt.Refresh()
		})
	*/

	myContainer := container.NewBorder(myButton, nil, nil, nil, parseAndColorText(inputText))
	myMainWindow.SetContent(myContainer)

	myMainWindow.ShowAndRun()

}

func parseAndColorText(inputText string) (tempRichText *widget.RichText) {
	var segments []widget.RichTextSegment

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
