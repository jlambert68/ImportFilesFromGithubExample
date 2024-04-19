package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"log"
	"regexp"
	"strings"
	"time"
)

func generateFilteredList(parentWindow fyne.Window) {

	filteredFileList = widget.NewList(
		func() int {
			return len(githubFilesFiltered)
		},
		func() fyne.CanvasObject {
			// Create a customFilteredLabel for each item.
			label := newCustomFilteredLabel("Template", func() {
				// Define double-click action here.
			})
			return label
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// Update the label text and double-click action for each item.
			label := obj.(*customFilteredLabel)
			label.Text = githubFilesFiltered[id].Name

			if githubFilesFiltered[id].Type == "file" {
				label.TextStyle = fyne.TextStyle{Italic: true}
			}

			label.onDoubleTap = func() {

				selectedFile := githubFilesFiltered[id]
				if selectedFile.Type == "dir" {
					// The item is a directory; fetch its contents
					getFileListFromGitHub(selectedFile.URL)
					filterFileListFromGitHub()
					filteredFileList.Refresh() // Refresh the list to update it with the new contents
					currentPathShowedinGUI.Set(strings.Split(selectedFile.URL, "?")[0])

					currentApiUrl = selectedFile.URL
				} else if selectedFile.Type == "file" {
					// Add file to selectedFiles and refresh the list only when if it doesn't exist
					var shouldAddFile bool
					shouldAddFile = true
					for _, existingSelectedFile := range selectedFiles {
						if existingSelectedFile.URL == selectedFile.URL {
							shouldAddFile = false
							break
						}
					}

					if shouldAddFile == true {
						selectedFiles = append(selectedFiles, selectedFile)
						UpdateSelectedFilesTable()
						selectedFilesTable.Refresh()

					}

				} else {
					// Show a dialog when other.
					dialog.ShowInformation("Info", "Double-clicked on: "+githubFiles[id].Name+" with Type "+githubFiles[id].Type, parentWindow)
				}
			}
			label.Refresh()
		},
	)
}

func filterFileListFromGitHub() {

	var fullRegExFilter string
	var tempGithubFilesFiltered []GitHubFile

	var tempRegex string

	for fileFilter, _ := range fileRegExFilterMap {
		if fileFilter == "*.*" {
			tempRegex = ".*"
		} else {

			tempRegex = strings.ReplaceAll(fileFilter, "*", "\\")
		}
		tempRegex = tempRegex + "$"

		if len(fullRegExFilter) == 0 {
			fullRegExFilter = fullRegExFilter + tempRegex
		} else {
			fullRegExFilter = fullRegExFilter + "|" + tempRegex
		}
	}

	if len(tempRegex) == 0 {
		tempRegex = `.*`
	}

	combinedRegex, err := regexp.Compile(fullRegExFilter)
	if err != nil {
		log.Fatalln("Error compiling regex:", err)
		return
	}

	for _, githubFile := range githubFiles {
		if combinedRegex.MatchString(githubFile.Name) == true || githubFile.Type == "dir" {
			tempGithubFilesFiltered = append(tempGithubFilesFiltered, githubFile)
		}
	}

	githubFilesFiltered = tempGithubFilesFiltered
}

type customFilteredLabel struct {
	widget.Label
	onDoubleTap func()
	lastTap     time.Time
}

func newCustomFilteredLabel(text string, onDoubleTap func()) *customFilteredLabel {
	l := &customFilteredLabel{Label: widget.Label{Text: text}, onDoubleTap: onDoubleTap}
	l.ExtendBaseWidget(l)
	return l
}

func (l *customFilteredLabel) Tapped(e *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(l.lastTap) < 500*time.Millisecond { // 500 ms as double-click interval
		if l.onDoubleTap != nil {
			l.onDoubleTap()
		}
	}
	l.lastTap = now
}

func (l *customFilteredLabel) TappedSecondary(*fyne.PointEvent) {
	// Implement if you need right-click (secondary tap) actions.
}

func (l *customFilteredLabel) MouseIn(*desktop.MouseEvent)    {}
func (l *customFilteredLabel) MouseMoved(*desktop.MouseEvent) {}
func (l *customFilteredLabel) MouseOut()                      {}
