package testDataSelector

import (
	//"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"sort"
	"sync"
)

// The color of a row in the TestData-popup-table when it is selected
var selectedRowColor = color.NRGBA{R: 0xa5, G: 0xf2, B: 0xa2, A: 0xff}

// The color of a row when it is selected and the user is hovering the row
var selectedAndHoveredRowColor = color.NRGBA{R: 0x60, G: 0xb8, B: 0xf7, A: 0xff}

// Icons used for Ascending and Descending sort indicator
var ascendingSortIndicatorIcon = widget.NewIcon(theme.MoveUpIcon())
var descendingSortIndicatorIcon = widget.NewIcon(theme.MoveDownIcon())
var notFocusFortSortingIcon = widget.NewIcon(theme.MediaPauseIcon())

// CustomWidget represents a custom component that can switch between a label and an icon with a background color
type CustomWidget struct {
	widget.BaseWidget
	//isIcon     bool
	label          *widget.Label
	icon           *widget.Icon
	headerSortIcon *widget.Icon
	background     *canvas.Rectangle
	hovered        bool
	onHover        func(bool, int)
	onTapped       func(widget.TableCellID)
	cellID         widget.TableCellID
	isSelected     bool
	tableRef       *CustomTableWidget
	sortOrder      sortDirection
}

// NewCustomWidget creates a new CustomWidget
func NewCustomWidget(isSelected bool, text string, tableRef *CustomTableWidget) *CustomWidget {
	w := &CustomWidget{
		isSelected:     isSelected,
		label:          widget.NewLabel(text),
		icon:           widget.NewIcon(theme.CheckButtonCheckedIcon()), // Replace with desired icon or picture
		headerSortIcon: notFocusFortSortingIcon,
		background:     canvas.NewRectangle(theme.BackgroundColor()),
		tableRef:       tableRef,
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer implements fyne.WidgetRenderer for CustomWidget
func (w *CustomWidget) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{w.background, w.label, w.icon, w.headerSortIcon}
	w.updateVisibility()

	if w.cellID.Col == 0 {
		// First column
		return &customWidgetRenderer{objects: objects, widget: w, layout: container.NewHBox(w.icon)}
	} else {
		// Data column
		return &customWidgetRenderer{objects: objects, widget: w, layout: container.NewHBox(w.headerSortIcon)}
	}
}

// customWidgetRenderer is the renderer for CustomWidget
type customWidgetRenderer struct {
	objects []fyne.CanvasObject
	widget  *CustomWidget
	layout  *fyne.Container
}

func (r *customWidgetRenderer) Layout(size fyne.Size) {
	r.widget.background.Resize(size)
	r.widget.label.Resize(size)
	r.widget.icon.Resize(size)
	r.layout.Resize(size)
}

func (r *customWidgetRenderer) MinSize() fyne.Size {

	// label's + sort icon width and labels height
	var newMinSize fyne.Size

	tempLabel := widget.NewLabel("Measure Height")
	tempLabel.Refresh()

	if r.widget.cellID.Col == 0 {
		// First column
		newMinSize = fyne.NewSize(0, tempLabel.MinSize().Height) // Doesn't matter
	} else {
		// Data column
		newMinSize = fyne.NewSize(r.objects[1].Size().Width+r.objects[3].Size().Width, tempLabel.MinSize().Height)
	}

	return newMinSize
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
		// Hovering the row
		if rowIsSelected == true {
			// Hovering and row is selected
			w.background.FillColor = selectedAndHoveredRowColor

		} else {
			// Hovering and row is not selected
			w.background.FillColor = theme.PrimaryColor()
		}

	} else {
		// Not hovering the row
		if rowIsSelected == true {
			// Not Hovering and row is selected
			w.background.FillColor = selectedRowColor

		} else {
			// Not Hovering and row is not selected
			w.background.FillColor = theme.BackgroundColor()
		}
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
	tableData [][]string
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
		tableData:        data,
		Table:            &widget.Table{},
		cellObjects:      make(map[widget.TableCellID]*CustomWidget),
		rowIsSelectedMap: make(map[int]bool),
		//rowStatus:   make([]bool, len(data)),
		//hoveredRow:  -1,
	}
	table.Length = func() (int, int) {
		return len(table.tableData), len(table.tableData[0]) + 1 // Adding one more column for the status
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
			obj.(*CustomWidget).SetText(table.tableData[cellID.Row][cellID.Col-1])
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

	// Tapped on Header
	if cellID.Row == 0 {

		// Extract the TableCell
		tableCell := t.cellObjects[cellID]

		// Pick sort order
		var sortOrder sortDirection
		if tableCell.sortOrder == dataSortOrderNotSelected || tableCell.sortOrder == dataSortDescending {

			sortOrder = dataSortAscending
			tableCell.sortOrder = dataSortAscending
			tableCell.headerSortIcon = ascendingSortIndicatorIcon

		} else {
			sortOrder = dataSortDescending
			tableCell.sortOrder = dataSortDescending
			tableCell.headerSortIcon = descendingSortIndicatorIcon
		}

		// Clear out all other sort order icons
		for columnIndex := 1; columnIndex < len(table.tableData[0]); columnIndex++ {
			if columnIndex != cellID.Col {
				newCellId := widget.TableCellID{
					Row: 0,
					Col: columnIndex,
				}
				// Only set data if teh cell has been initialized
				_, existInMap := table.cellObjects[newCellId]
				if existInMap == true {
					table.cellObjects[newCellId].sortOrder = dataSortOrderNotSelected
					table.cellObjects[newCellId].headerSortIcon = notFocusFortSortingIcon
				}
			}
		}

		sortTable(table.tableData, cellID.Col-1, sortOrder)
		updateRowsSelectedMap(table)
		table.Refresh()

		return
	}

	// Not the header row
	var isSelected bool
	isSelected = t.rowIsSelectedMap[cellID.Row]
	isSelected = !isSelected
	t.rowIsSelectedMap[cellID.Row] = isSelected

	// Update all cells on the row with selected or not
	for tempCellId, tableCell := range t.cellObjects {
		// Only update cells on the celected row
		if tempCellId.Row == cellID.Row {
			tableCell.isSelected = isSelected
		}
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

// sortDirection defines the sorting order
type sortDirection int

const (
	dataSortOrderNotSelected sortDirection = iota
	dataSortAscending
	dataSortDescending
)

// sortTable sorts a 2D string slice based on the specified column and direction,
// keeping the first row intact.
func sortTable(data [][]string, column int, direction sortDirection) {
	// Only sort the rows starting from the second row
	sort.Slice(data[1:], func(i, j int) bool {
		if direction == dataSortOrderNotSelected {
			direction = dataSortAscending
		}

		if direction == dataSortAscending {
			return data[i+1][column] < data[j+1][column]
		}
		return data[i+1][column] > data[j+1][column]
	})
}

// updateRowsSelectedMap updates the map holding which row that is selected
func updateRowsSelectedMap(table *CustomTableWidget) {

	var cellId widget.TableCellID
	var isRowSelected bool

	// Loop the rows
	for rowIndex, _ := range table.tableData {

		// Create a CellId to use
		cellId = widget.TableCellID{
			Row: rowIndex,
			Col: 0,
		}

		// Extract from Cell if row is selected
		isRowSelected = table.cellObjects[cellId].isSelected

		// Recreate Map holding if row is selected or not
		table.rowIsSelectedMap[rowIndex] = isRowSelected

	}

}
