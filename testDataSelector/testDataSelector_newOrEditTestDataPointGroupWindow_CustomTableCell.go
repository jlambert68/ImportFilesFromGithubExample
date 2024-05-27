package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

// CustomWidget represents a custom component that can switch between a label and an icon
type CustomWidget struct {
	widget.BaseWidget
	isIcon  bool
	label   *widget.Label
	icon    *widget.Icon
	lastTap time.Time
}

// NewCustomWidget creates a new CustomWidget
func NewCustomWidget(isIcon bool, text string) *CustomWidget {
	w := &CustomWidget{
		isIcon: isIcon,
		label:  widget.NewLabel(text),
		icon:   widget.NewIcon(theme.ContentAddIcon()), // Replace with desired icon or picture
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer implements fyne.WidgetRenderer for CustomWidget
func (w *CustomWidget) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{w.label, w.icon}
	w.updateVisibility()
	return &customWidgetRenderer{objects: objects, widget: w}
}

// customWidgetRenderer is the renderer for CustomWidget
type customWidgetRenderer struct {
	objects []fyne.CanvasObject
	widget  *CustomWidget
}

func (r *customWidgetRenderer) Layout(size fyne.Size) {
	for _, obj := range r.objects {
		obj.Resize(size)
	}
}

func (r *customWidgetRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

func (r *customWidgetRenderer) Refresh() {
	r.widget.updateVisibility()
	canvas.Refresh(r.widget)
}

func (r *customWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *customWidgetRenderer) Destroy() {}

// updateVisibility updates the visibility of the label and icon based on the isIcon flag
func (w *CustomWidget) updateVisibility() {
	if w.isIcon {
		w.icon.Show()
		w.label.Hide()
	} else {
		w.icon.Hide()
		w.label.Show()
	}
}

// SetText sets the text of the label
func (w *CustomWidget) SetText(text string) {
	if !w.isIcon {
		w.label.SetText(text)
	}
}

// CustomTableWidget represents the custom table with row double-click handling
type CustomTableWidget struct {
	*widget.Table
	cellObjects map[widget.TableCellID]fyne.CanvasObject
	rowStatus   []bool
	lastTap     time.Time
	tapCount    int
}

func NewCustomTableWidget(data [][]string) *CustomTableWidget {
	table := &CustomTableWidget{
		Table:       &widget.Table{},
		cellObjects: make(map[widget.TableCellID]fyne.CanvasObject),
		rowStatus:   make([]bool, len(data)),
	}
	table.Length = func() (int, int) {
		return len(data), len(data[0]) + 1 // Adding one more column for the status
	}
	table.CreateCell = func() fyne.CanvasObject {
		return NewCustomWidget(false, "")
	}
	table.UpdateCell = func(cellID widget.TableCellID, obj fyne.CanvasObject) {
		customWidget := obj.(*CustomWidget)
		if cellID.Col == 0 {
			// Update the first column to show the double-click status
			if table.rowStatus[cellID.Row] {
				customWidget.SetText("Clicked")
			} else {
				customWidget.SetText("Not Clicked")
			}
			customWidget.isIcon = false
			customWidget.Refresh()
		} else {
			// Update other columns with data
			customWidget.SetText(data[cellID.Row][cellID.Col-1])
			table.cellObjects[cellID] = obj
		}
	}
	table.ExtendBaseWidget(table)
	table.OnSelected = func(id widget.TableCellID) {

		fmt.Println(id.Row, id.Col)

		now := time.Now()
		if now.Sub(table.lastTap) < 500*time.Millisecond { // 500ms for double-tap
			table.tapCount++
		} else {
			table.tapCount = 1
		}
		table.lastTap = now

		if table.tapCount == 2 {
			table.tapCount = 0
			table.toggleRowIcon(id.Row)
		}
	}

	setColumnWidths(table.Table, data)

	return table
}

func (t *CustomTableWidget) toggleRowIcon(row int) {
	t.rowStatus[row] = !t.rowStatus[row]
	cellID := widget.TableCellID{Row: row, Col: 0}
	customWidget := t.cellObjects[cellID].(*CustomWidget)
	if t.rowStatus[row] {
		customWidget.SetText("Clicked")
	} else {
		customWidget.SetText("Not Clicked")
	}
	customWidget.Refresh()

	_, cols := t.Length()
	for col := 1; col < cols; col++ {
		cellID := widget.TableCellID{Row: row, Col: col}
		customWidget := t.cellObjects[cellID].(*CustomWidget)
		customWidget.isIcon = !customWidget.isIcon
		customWidget.Refresh()
	}
}
