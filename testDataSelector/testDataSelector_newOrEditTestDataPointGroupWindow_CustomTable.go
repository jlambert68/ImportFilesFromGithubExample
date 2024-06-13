package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"

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

// customWidget represents a custom component that can switch between a label and an icon with a background color
type customWidget struct {
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

// NewCustomWidget creates a new customWidget
func NewCustomWidget(isSelected bool, text string, tableRef *CustomTableWidget) *customWidget {
	w := &customWidget{
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

// CreateRenderer implements fyne.WidgetRenderer for customWidget
func (w *customWidget) CreateRenderer() fyne.WidgetRenderer {
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

// customWidgetRenderer is the renderer for customWidget
type customWidgetRenderer struct {
	objects []fyne.CanvasObject
	widget  *customWidget
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
func (w *customWidget) updateVisibility() {

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
func (w *customWidget) SetText(text string) {
	w.label.SetText(text)
}

// SetCellID sets the position of the cell in the Table
func (w *customWidget) SetCellID(cellID widget.TableCellID) {
	w.cellID = cellID
}

func (w *customWidget) MouseIn(*desktop.MouseEvent) {
	w.hovered = true
	if w.onHover != nil {
		w.onHover(true, w.cellID.Row)
	}
	w.Refresh()
}

func (w *customWidget) MouseOut() {
	w.hovered = false
	if w.onHover != nil {
		w.onHover(false, w.cellID.Row)
	}
	w.Refresh()
}

func (w *customWidget) MouseMoved(*desktop.MouseEvent) {}

// Tapped handles tap events
func (w *customWidget) Tapped(*fyne.PointEvent) {
	if w.onTapped != nil {
		w.onTapped(w.cellID)
	}
}

// CustomTableWidget represents the custom table with row double-click handling and hover effects
type CustomTableWidget struct {
	tableData [][]string
	*widget.Table
	cellObjects        map[widget.TableCellID]*customWidget
	rowIsSelectedMap   map[int]bool
	rowUuidIsSectedMap map[string]bool
	//lastTap     time.Time
	//tapCount    int
	//hoveredRow  int
}

// Use a mutex to synchronize access to the map
var tableMutex sync.Mutex

func NewCustomTableWidget(data [][]string) *CustomTableWidget {
	table := &CustomTableWidget{
		tableData:          data,
		Table:              &widget.Table{},
		cellObjects:        make(map[widget.TableCellID]*customWidget),
		rowIsSelectedMap:   make(map[int]bool),
		rowUuidIsSectedMap: make(map[string]bool),
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

		obj.(*customWidget).SetCellID(cellID)

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
			obj.(*customWidget).SetText("")
			if table.rowIsSelectedMap[cellID.Row] == true {
				obj.(*customWidget).icon.Show()
				obj.(*customWidget).label.Hide()
			} else {
				obj.(*customWidget).icon.Hide()
				obj.(*customWidget).label.Show()
			}

		} else {
			// Update other columns with data
			obj.(*customWidget).SetText(table.tableData[cellID.Row][cellID.Col-1])
			obj.(*customWidget).icon.Hide()
			obj.(*customWidget).label.Show()

			if cellID.Row == 0 {

				obj.(*customWidget).label.TextStyle.Bold = true
			} else {
				obj.(*customWidget).label.TextStyle.Bold = false
			}
		}
		obj.(*customWidget).onHover = func(hovered bool, row int) {
			if hovered {
				table.hoverRow(row)
			} else {
				table.unhoverRow(row)
			}
		}
		obj.(*customWidget).onTapped = func(cellID widget.TableCellID) {
			table.handleCellTapped(cellID, table)
		}

		// Hinder concurrent map writes
		tableMutex.Lock()

		table.cellObjects[cellID] = obj.(*customWidget)
		obj.(*customWidget).Refresh()

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
		table.Refresh()
		updateRowsSelectedMap(table)
		table.Refresh()

		return
	}

	// Not the header row
	var isSelected bool
	isSelected = t.rowIsSelectedMap[cellID.Row]
	isSelected = !isSelected
	t.rowIsSelectedMap[cellID.Row] = isSelected

	// Get RowUUID for selected row
	var rowUuid string
	rowUuid = table.tableData[cellID.Row][len(table.tableData[0])-1]

	// Update selected row for Row Uuid
	t.rowUuidIsSectedMap[rowUuid] = isSelected

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

	var isRowSelected bool
	var rowUuid string
	var uuidColumn int

	// Get column for Row Uuid
	uuidColumn = len(table.tableData[0]) - 1

	// Loop the rows
	for rowIndex, _ := range table.tableData {

		// Get Row Uuid
		rowUuid = table.tableData[rowIndex][uuidColumn]

		// Extract from row is selected
		isRowSelected = table.rowUuidIsSectedMap[rowUuid]

		// Recreate values in Map holding if row is selected or not
		table.rowIsSelectedMap[rowIndex] = isRowSelected

	}

}

// setColumnWidths adapt all columns in the popup window to fit its headers and data
func setColumnWidths(table *widget.Table, data [][]string) {
	maxWidths := make([]float32, len(data[0]))
	var width float32
	for col := range maxWidths {
		for row := range data {
			if row == 0 {
				width = fyne.MeasureText(data[row][col], theme.TextSize(), fyne.TextStyle{Bold: true}).Width
			} else {
				width = fyne.MeasureText(data[row][col], theme.TextSize(), fyne.TextStyle{}).Width
			}
			if width > maxWidths[col] {
				maxWidths[col] = width
			}
		}
		// Add some padding to the maximum width found
		maxWidths[col] += theme.Padding() * 4
	}

	for col, columnWidth := range maxWidths {
		table.SetColumnWidth(col+1, columnWidth)
	}

	table.SetColumnWidth(0, fyne.MeasureText("xxxx", theme.TextSize(), fyne.TextStyle{}).Width)
}

// showTable creates and shows a table for the selected node with data
func showTable(w fyne.Window, data [][]string) {

	// Create a new table with an extra column for checkboxes
	table := NewCustomTableWidget(data)

	// Set minimum size for the table to ensure it's larger
	table.Resize(fyne.NewSize(600, 500)) // Set the minimum size to 400x300 pixels

	// Use a scroll container to make the table scrollable in case it has more data
	scrollContainer := container.NewScroll(table)

	modal := dialog.NewCustomConfirm("Chose TestDataPoints", "Select TestDataPoints", "Cancel",
		scrollContainer,
		func(response bool) {
			if response {
				println("User confirmed action")

				// Extract rows that were selecte
				for row, isSelected := range table.rowIsSelectedMap {
					if isSelected == true {
						cellID := widget.TableCellID{row, 1}
						fmt.Println(table.cellObjects[cellID])
					}
				}

			} else {
				println("User canceled action")
			}

		}, w)
	modal.Resize(fyne.NewSize(800, 600))
	modal.Show()

}
