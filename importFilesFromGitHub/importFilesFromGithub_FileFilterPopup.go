package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type checklistItem struct {
	Label   string
	Checked bool
}

func generateFileFilterPopup(parentWindow fyne.Window) {

	var checkBoxitems = []checklistItem{
		{"*.go", false},            // Not prechecked
		{"*.txt", false},           // Not prechecked
		{"*.json", false},          // Not prechecked
		{"*.xml", false},           // Not prechecked
		{"*.fenixFunction", false}, // Not prechecked
		{"*.*", true},              // Prechecked
		// Add more items as needed
	}

	// Map to track selected options
	selectedOptions := make(map[string]bool)

	// Create checkboxes for each option
	checkboxList := container.NewVBox()
	for _, checkBoxitem := range checkBoxitems {
		checkbox := widget.NewCheck(checkBoxitem.Label, func(checked bool) {
			selectedOptions[checkBoxitem.Label] = checked
			if checked == true {
				fileRegExFilterMap[checkBoxitem.Label] = checkBoxitem.Label
			} else {
				delete(fileRegExFilterMap, checkBoxitem.Label)
			}
			filterFileListFromGitHub()
			filteredFileList.Refresh()
		})
		if checkBoxitem.Checked == true {
			fileRegExFilterMap[checkBoxitem.Label] = checkBoxitem.Label
			checkbox.SetChecked(true)
		}
		checkboxList.Add(checkbox)
	}

	// Button to show the multi-select dropdown
	fileFilterPopupButton = widget.NewButton("Select Options", func() {
		popUp := widget.NewPopUp(checkboxList, parentWindow.Canvas())
		popUp.Show()
		popUp.Resize(fyne.NewSize(200, 200)) // Adjust size as needed
	})

}
