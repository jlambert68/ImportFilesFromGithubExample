package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"time"
)

// Create the UI-list that holds the chpsed
func generateSelectedFilesList(parentWindow fyne.Window) {

	selectedFilesList = widget.NewList(
		func() int { return len(selectedFiles) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(selectedFiles[i].Name)
		},
	)

	selectedFilesList.OnSelected = func(id widget.ListItemID) {
		// Remove the file from selectedFiles and refresh the list
		selectedFiles = append(selectedFiles[:id], selectedFiles[id+1:]...)
		selectedFilesList.Unselect(id)
		selectedFilesList.Refresh()
	}

}

type customLabel struct {
	widget.Label
	onDoubleTap func()
	lastTap     time.Time
}

func newCustomLabel(text string, onDoubleTap func()) *customLabel {
	l := &customLabel{Label: widget.Label{Text: text}, onDoubleTap: onDoubleTap}
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
