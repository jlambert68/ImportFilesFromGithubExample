package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

// CustomWidget represents a custom component that can switch between a label and an icon with a background color
type CustomWidget struct {
	widget.BaseWidget
	//isIcon     bool
	label      *widget.Label
	icon       *widget.Icon
	background *canvas.Rectangle
	hovered    bool
	onHover    func(bool, int)
	onTapped   func(widget.TableCellID)
	cellID     widget.TableCellID
	isSelected bool
	tableRef   *CustomTableWidget
}

// NewCustomWidget creates a new CustomWidget
func NewCustomWidget(isSelected bool, text string, tableRef *CustomTableWidget) *CustomWidget {
	w := &CustomWidget{
		isSelected: isSelected,
		label:      widget.NewLabel(text),
		icon:       widget.NewIcon(theme.CheckButtonCheckedIcon()), // Replace with desired icon or picture
		background: canvas.NewRectangle(theme.BackgroundColor()),
		tableRef:   tableRef,
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

	var existInMap bool
	var rowIsSelected bool

	rowIsSelected, existInMap = w.tableRef.rowIsSelectedMap[w.cellID.Row]
	if existInMap == false {
		rowIsSelected = false
	}

	if w.cellID.Col == 0 {
		if rowIsSelected {
			w.icon.Show()
			w.label.Hide()
		} else {
			w.icon.Hide()
			w.label.Show()
		}
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
	w.label.SetText(text)
}

// SetCellID sets the position of the cell in the Table
func (w *CustomWidget) SetCellID(cellID widget.TableCellID) {
	w.cellID = cellID
}

func (w *CustomWidget) MouseIn(*desktop.MouseEvent) {
	w.hovered = true
	if w.onHover != nil {
		w.onHover(true, w.cellID.Row)
	}
	w.Refresh()
}

func (w *CustomWidget) MouseOut() {
	w.hovered = false
	if w.onHover != nil {
		w.onHover(false, w.cellID.Row)
	}
	w.Refresh()
}

func (w *CustomWidget) MouseMoved(*desktop.MouseEvent) {}

// Tapped handles tap events
func (w *CustomWidget) Tapped(*fyne.PointEvent) {
	if w.onTapped != nil {
		w.onTapped(w.cellID)
	}
}

// CustomTableWidget represents the custom table with row double-click handling and hover effects
type CustomTableWidget struct {
	*widget.Table
	cellObjects      map[widget.TableCellID]*CustomWidget
	rowIsSelectedMap map[int]bool
	//lastTap     time.Time
	//tapCount    int
	//hoveredRow  int
}

// Use a mutex to synchronize access to the map
var tableMutex sync.Mutex

func NewCustomTableWidget(data [][]string) *CustomTableWidget {
	table := &CustomTableWidget{
		Table:            &widget.Table{},
		cellObjects:      make(map[widget.TableCellID]*CustomWidget),
		rowIsSelectedMap: make(map[int]bool),
		//rowStatus:   make([]bool, len(data)),
		//hoveredRow:  -1,
	}
	table.Length = func() (int, int) {
		return len(data), len(data[0]) + 1 // Adding one more column for the status
	}
	table.CreateCell = func() fyne.CanvasObject {
		return NewCustomWidget(false, "", table)
	}
	table.UpdateCell = func(cellID widget.TableCellID, obj fyne.CanvasObject) {

		obj.(*CustomWidget).SetCellID(cellID)

		if cellID.Col == 0 {
			/*
				// Update the first column to show the double-click status
				if table.rowStatus[cellID.Row] {
					customWidget.SetText("Clicked")
				} else {
					customWidget.SetText("Not Clicked")
				}
				customWidget.isIcon = false

			*/
			obj.(*CustomWidget).SetText("")
			if table.rowIsSelectedMap[cellID.Row] == true {
				obj.(*CustomWidget).icon.Show()
				obj.(*CustomWidget).label.Hide()
			} else {
				obj.(*CustomWidget).icon.Hide()
				obj.(*CustomWidget).label.Show()
			}

		} else {
			// Update other columns with data
			obj.(*CustomWidget).SetText(data[cellID.Row][cellID.Col-1])
			obj.(*CustomWidget).icon.Hide()
			obj.(*CustomWidget).label.Show()

			if cellID.Row == 0 {

				obj.(*CustomWidget).label.TextStyle.Bold = true
			} else {
				obj.(*CustomWidget).label.TextStyle.Bold = false
			}
		}
		obj.(*CustomWidget).onHover = func(hovered bool, row int) {
			if hovered {
				table.hoverRow(row)
			} else {
				table.unhoverRow(row)
			}
		}
		obj.(*CustomWidget).onTapped = func(cellID widget.TableCellID) {
			table.handleCellTapped(cellID, table)
		}

		// Hinder concurrent map writes
		tableMutex.Lock()

		table.cellObjects[cellID] = obj.(*CustomWidget)
		obj.(*CustomWidget).Refresh()

		// Release map
		tableMutex.Unlock()

	}
	table.ExtendBaseWidget(table)
	/*table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			table.toggleRowIcon(id.Row)
		}

	}*/

	setColumnWidths(table.Table, data)

	return table
}

func (t *CustomTableWidget) handleCellTapped(cellID widget.TableCellID, table *CustomTableWidget) {
	// Handle cell click logic here

	println("Cell tapped:", cellID.Row, cellID.Col)

	if cellID.Row == 0 {
		return
	}

	isSelected := t.rowIsSelectedMap[cellID.Row]
	isSelected = !isSelected
	t.rowIsSelectedMap[cellID.Row] = isSelected

	for k, v := range t.rowIsSelectedMap {
		fmt.Println(k, v)
	}

	t.Refresh()

	//firstCellInRow := widget.TableCellID{Row: cellID.Row, Col: 0}

	//customWidget := t.cellObjects[firstCellInRow]
	//customWidget.isSelected = !customWidget.isSelected

	//if cellID.Row > 0 {
	//	table.toggleRowIcon(cellID.Row)
	//}
}

func (t *CustomTableWidget) hoverRow(row int) {
	if row == 0 { //t.hoveredRow == row || row == 0 {
		return
	}
	/*
		if t.hoveredRow != -1 {
			t.unhoverRow(t.hoveredRow)
		}
		t.hoveredRow = row

	*/
	_, cols := t.Length()
	for col := 0; col < cols; col++ {
		cellID := widget.TableCellID{Row: row, Col: col}

		customWidget := t.cellObjects[cellID]

		// Only change stuff if the column(row) is visible and has got an "object value"
		if customWidget != nil && customWidget.cellID.Row == row {
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
	//t.hoveredRow = -1
}
