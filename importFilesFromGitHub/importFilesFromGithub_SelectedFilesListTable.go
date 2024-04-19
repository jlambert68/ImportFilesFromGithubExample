package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

type clickableLabel struct {
	widget.Label
	onDoubleTap func()
	lastTapTime time.Time
	isClickable bool
}

func newClickableLabel(text string, onDoubleTap func(), tempIsClickable bool) *clickableLabel {
	l := &clickableLabel{
		widget.Label{Text: text},
		onDoubleTap,
		time.Now(),
		tempIsClickable}
	l.ExtendBaseWidget(l)
	return l
}

func (l *clickableLabel) Tapped(e *fyne.PointEvent) {
	if l.isClickable == false {
		return
	}

	if time.Since(l.lastTapTime) < 500*time.Millisecond {
		if l.onDoubleTap != nil {
			l.onDoubleTap()
		}
	}
	l.lastTapTime = time.Now()
}

func (l *clickableLabel) TappedSecondary(*fyne.PointEvent) {
	// Implement if you need right-click (secondary tap) actions.
}

func (l *clickableLabel) MouseIn(*desktop.MouseEvent)    {}
func (l *clickableLabel) MouseMoved(*desktop.MouseEvent) {}
func (l *clickableLabel) MouseOut()                      {}

// Create the UI-list that holds the selected files
func generateSelectedFilesListTable(parentWindow fyne.Window) {
	// Correctly initialize the selectedFilesTable as a new table
	selectedFilesTable = widget.NewTable(
		func() (int, int) { return 0, 2 }, // Start with zero rows, 2 columns
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Create cells as labels
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// This should be filled when updating the table
		},
	)

	/*
			selectedFilesTable = widget.NewList(
				func() int { return len(selectedFiles) },
				func() fyne.CanvasObject { return widget.NewLabel("") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(selectedFiles[i].Name)
				},
			)



		selectedFilesTable.OnSelected = func(id widget.ListItemID) {

		}
	*/

}

func UpdateSelectedFilesTable() {

	selectedFilesTable.Length = func() (int, int) {
		return len(selectedFiles), 2
	}
	selectedFilesTable.CreateCell = func() fyne.CanvasObject {
		return newClickableLabel("", func() {}, false)

	}
	selectedFilesTable.UpdateCell = func(id widget.TableCellID, cell fyne.CanvasObject) {
		switch id.Col {
		case 0:
			// For the "Name" column, use the clickable label
			clickable := cell.(*clickableLabel)
			clickable.SetText(selectedFiles[id.Row].Name)
			clickable.isClickable = true

			clickable.onDoubleTap = func() {

				// Remove the file from selectedFiles and refresh the list
				for fileIndex, file := range selectedFiles {
					if file.URL == selectedFiles[id.Row].URL {
						selectedFiles = append(selectedFiles[:fileIndex], selectedFiles[fileIndex+1:]...)
						selectedFilesTable.Unselect(id)
						selectedFilesTable.Refresh()
						UpdateSelectedFilesTable()
						break
					}
				}

			}

		case 1:
			// For the "URL" column, use a regular label

			nonClickable := cell.(*clickableLabel)
			nonClickable.SetText(selectedFiles[id.Row].URL)
		}
	}

	maxNameWidth := float32(150) // Start with a minimum width
	maxUrlWidth := float32(250)  // Start with a minimum width
	for _, file := range selectedFiles {
		textNameWidth := fyne.MeasureText(file.Name, theme.TextSize(), fyne.TextStyle{}).Width
		textUrlWidth := fyne.MeasureText(file.URL, theme.TextSize(), fyne.TextStyle{}).Width
		if textNameWidth > maxNameWidth {
			maxNameWidth = textNameWidth
		}
		if textUrlWidth > maxUrlWidth {
			maxUrlWidth = textUrlWidth
		}
	}

	selectedFilesTable.SetColumnWidth(0, maxNameWidth+theme.Padding()*4) // Add padding
	selectedFilesTable.SetColumnWidth(1, maxUrlWidth+theme.Padding()*4)  // Path column width can be static or calculated similarly

	selectedFilesTable.Refresh()

}

type customLabel struct {
	widget.Label
	onDoubleTap func()
	lastTap     time.Time
}

func newCustomLabel(text string, onDoubleTap func()) *customLabel {
	l := &customLabel{Label: widget.Label{Text: text}, onDoubleTap: onDoubleTap, lastTap: time.Now()}
	l.ExtendBaseWidget(l)
	return l
}

func (l *customLabel) Tapped(e *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(l.lastTap) < 500*time.Millisecond { // 500 ms as double-click interval
		if l.onDoubleTap != nil {
			l.onDoubleTap()
		}
	}
	l.lastTap = now
}

func (l *customLabel) TappedSecondary(*fyne.PointEvent) {
	// Implement if you need right-click (secondary tap) actions.
}

func (l *customLabel) MouseIn(*desktop.MouseEvent)    {}
func (l *customLabel) MouseMoved(*desktop.MouseEvent) {}
func (l *customLabel) MouseOut()                      {}

/*
type coloredLabelItem struct {
	text  string
	color color.Color
}

func newColoredLabelItem(text string, color color.Color) *coloredLabelItem {
	return &coloredLabelItem{text: text, color: color}
}

func (item *coloredLabelItem) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(item.text)
	label.color = item.color
	label.Refresh()

	return widget.NewSimpleRenderer(label)
}

*/
