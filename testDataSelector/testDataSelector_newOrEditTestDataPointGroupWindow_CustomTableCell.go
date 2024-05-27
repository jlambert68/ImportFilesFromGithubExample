package testDataSelector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

// CustomWidget represents a custom component that can switch between a label and an icon with a background color
type CustomWidget struct {
	widget.BaseWidget
	isIcon     bool
	label      *widget.Label
	icon       *widget.Icon
	background *canvas.Rectangle
	hovered    bool
	onHover    func(bool)
	onTapped   func()
}

// NewCustomWidget creates a new CustomWidget
func NewCustomWidget(isIcon bool, text string) *CustomWidget {
	w := &CustomWidget{
		isIcon:     isIcon,
		label:      widget.NewLabel(text),
		icon:       widget.NewIcon(theme.ContentAddIcon()), // Replace with desired icon or picture
		background: canvas.NewRectangle(theme.BackgroundColor()),
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer implements fyne.WidgetRenderer for CustomWidget
func (w *CustomWidget) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{w.background, w.label, w.icon}
	w.updateVisibility()
	return &customWidgetRenderer{objects: objects, widget: w}
}

// customWidgetRenderer is the renderer for CustomWidget
type customWidgetRenderer struct {
	objects []fyne.CanvasObject
	widget  *CustomWidget
}

func (r *customWidgetRenderer) Layout(size fyne.Size) {
	r.widget.background.Resize(size)
	r.widget.label.Resize(size)
	r.widget.icon.Resize(size)
}

func (r *customWidgetRenderer) MinSize() fyne.Size {
	return r.objects[1].MinSize() // label's min size
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
	if w.hovered {
		w.background.FillColor = theme.PrimaryColor()
	} else {
		w.background.FillColor = theme.BackgroundColor()
	}
	w.background.Refresh()
}

// SetText sets the text of the label
func (w *CustomWidget) SetText(text string) {
	if !w.isIcon {
		w.label.SetText(text)
	}
}

func (w *CustomWidget) MouseIn(*desktop.MouseEvent) {
	w.hovered = true
	if w.onHover != nil {
		w.onHover(true)
	}
	w.Refresh()
}

func (w *CustomWidget) MouseOut() {
	w.hovered = false
	if w.onHover != nil {
		w.onHover(false)
	}
	w.Refresh()
}

func (w *CustomWidget) MouseMoved(*desktop.MouseEvent) {}

// Tapped handles tap events
func (w *CustomWidget) Tapped(*fyne.PointEvent) {
	if w.onTapped != nil {
		w.onTapped()
	}
}

// CustomTableWidget represents the custom table with row double-click handling and hover effects
type CustomTableWidget struct {
	*widget.Table
	cellObjects map[widget.TableCellID]*CustomWidget
	rowStatus   []bool
	lastTap     time.Time
	tapCount    int
	hoveredRow  int
}

func NewCustomTableWidget(data [][]string) *CustomTableWidget {
	table := &CustomTableWidget{
		Table:       &widget.Table{},
		cellObjects: make(map[widget.TableCellID]*CustomWidget),
		rowStatus:   make([]bool, len(data)),
		hoveredRow:  -1,
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
		} else {
			// Update other columns with data
			customWidget.SetText(data[cellID.Row][cellID.Col-1])
		}
		customWidget.onHover = func(hovered bool) {
			if hovered {
				table.hoverRow(cellID.Row)
			} else {
				table.unhoverRow(cellID.Row)
			}
		}
		customWidget.onTapped = func() {
			table.handleCellTapped(cellID)
		}
		table.cellObjects[cellID] = customWidget
		customWidget.Refresh()
	}
	table.ExtendBaseWidget(table)
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			table.toggleRowIcon(id.Row)
		}

	}

	setColumnWidths(table.Table, data)

	return table
}

func (t *CustomTableWidget) handleCellTapped(cellID widget.TableCellID) {
	// Handle cell click logic here
	println("Cell tapped:", cellID.Row, cellID.Col)
}

func (t *CustomTableWidget) toggleRowIcon(row int) {
	t.rowStatus[row] = !t.rowStatus[row]
	cellID := widget.TableCellID{Row: row, Col: 0}
	customWidget := t.cellObjects[cellID]
	if t.rowStatus[row] {
		customWidget.SetText("Clicked")
	} else {
		customWidget.SetText("Not Clicked")
	}
	customWidget.Refresh()

	_, cols := t.Length()
	for col := 1; col < cols; col++ {
		cellID := widget.TableCellID{Row: row, Col: col}
		customWidget := t.cellObjects[cellID]
		customWidget.isIcon = !customWidget.isIcon
		customWidget.Refresh()
	}
}

func (t *CustomTableWidget) hoverRow(row int) {
	if t.hoveredRow == row || row == 0 {
		return
	}
	if t.hoveredRow != -1 {
		t.unhoverRow(t.hoveredRow)
	}
	t.hoveredRow = row
	_, cols := t.Length()
	for col := 0; col < cols; col++ {
		cellID := widget.TableCellID{Row: row, Col: col}

		customWidget := t.cellObjects[cellID]

		// Only change stuff if the column(row) is visible and has got an "object value"
		if customWidget != nil {
			customWidget.hovered = true
			customWidget.Refresh()
		}
	}
}

func (t *CustomTableWidget) unhoverRow(row int) {
	_, cols := t.Length()
	for col := 0; col < cols; col++ {
		cellID := widget.TableCellID{Row: row, Col: col}
		customWidget := t.cellObjects[cellID]

		// Only change stuff if the column(row) is visible and has got an "object value"
		if customWidget != nil {
			customWidget.hovered = false
			customWidget.Refresh()
		}
	}
	t.hoveredRow = -1
}
