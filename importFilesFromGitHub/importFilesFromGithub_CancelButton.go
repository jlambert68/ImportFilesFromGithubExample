package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Generate the button that cancel everything and closes the window
func generateCancelButton(parentWindow fyne.Window) {

	cancelButton = widget.NewButton("Cancel", func() {
		fenixMainWindow.Show()
		parentWindow.Close()
	})
}
