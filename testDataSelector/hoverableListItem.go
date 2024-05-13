package testDataSelector

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// HoverableListItem represents a list item that reacts to hover and right-click events.
type HoverableListItem struct {
	widget.BaseWidget
	label      *canvas.Text
	background *canvas.Rectangle
	hovered    bool
	onTap      func()
	onRightTap func() // Function to call on right-click
}

func NewHoverableListItem(text string, onTap func(), onRightTap func()) *HoverableListItem {
	label := canvas.NewText(text, theme.ForegroundColor())
	label.TextStyle.Bold = true

	item := &HoverableListItem{
		label:      label,
		background: canvas.NewRectangle(theme.BackgroundColor()),
		onTap:      onTap,
		onRightTap: onRightTap,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (h *HoverableListItem) CreateRenderer() fyne.WidgetRenderer {
	return &hoverableListItemRenderer{
		item:    h,
		label:   h.label,
		rect:    h.background,
		objects: []fyne.CanvasObject{h.background, h.label},
	}
}

type hoverableListItemRenderer struct {
	item    *HoverableListItem
	label   *canvas.Text
	rect    *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *hoverableListItemRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
	r.label.Resize(size)
	r.label.Move(fyne.NewPos(10, 0)) // Small padding from the left
}

func (r *hoverableListItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.label.MinSize().Height+16) // Plus some padding vertically
}

func (r *hoverableListItemRenderer) Refresh() {
	r.rect.FillColor = theme.BackgroundColor()
	if r.item.hovered {
		r.rect.FillColor = theme.HoverColor()
	}
	r.label.Color = theme.ForegroundColor()
	r.rect.Refresh()
	r.label.Refresh()
}

func (r *hoverableListItemRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *hoverableListItemRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *hoverableListItemRenderer) Destroy() {}

func (h *HoverableListItem) Tapped(*fyne.PointEvent) {
	if h.onTap != nil {
		h.onTap()
	}
}

func (h *HoverableListItem) TappedSecondary(*fyne.PointEvent) {
	if h.onRightTap != nil {
		h.onRightTap()
	}
}

func (h *HoverableListItem) MouseIn(e *desktop.MouseEvent) {
	h.hovered = true
	h.Refresh()
}

func (h *HoverableListItem) MouseOut() {
	h.hovered = false
	h.Refresh()
}

func (h *HoverableListItem) MouseMoved(e *desktop.MouseEvent) {}

func generateListContainer(allAvailablePointsList *[]string, listContainer *fyne.Container) {

	// Remove children from container
	listContainer.RemoveAll()

	// Add all the children to the container
	for _, availablePoint := range *allAvailablePointsList {

		newHoverableListItem := NewHoverableListItem("Item 1", func() { fmt.Println(fmt.Sprintf("Left-clicked : %s", string(availablePoint))) }, func() { fmt.Println(fmt.Sprintf("Right-clicked : %s", string(availablePoint))) })
		listContainer.Add(newHoverableListItem)

	}

	// Refresh the ListContainer
	listContainer.Refresh()

}
