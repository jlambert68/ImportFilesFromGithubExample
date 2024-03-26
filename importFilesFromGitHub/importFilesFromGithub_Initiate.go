package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
	"time"
)

func InitiateImportFilesFromGitHubWindow(
	originalApiUrl string,
	mainWindow fyne.Window,
	myApp fyne.App,
	responseChannel *chan bool) *[]GitHubFile {

	// Cleare list from Previous Import
	githubFilesFiltered = nil

	// Disable the main window
	mainWindow.Hide()

	// Store Channel reference in local varaible
	sharedResponseChannel = responseChannel

	// Store reference to Fenix Main Window
	fenixMainWindow = mainWindow

	// Read Environment variables
	runInit()

	// Create the window for GitHub files
	githubFileImporterWindow = myApp.NewWindow("GitHub file importer")

	rootApiUrl = originalApiUrl
	currentApiUrl = rootApiUrl

	currentPathShowedinGUI = binding.NewString()
	currentPathShowedinGUI.Set(strings.Split(currentApiUrl, "?")[0]) // Setting initial value

	fileRegExFilterMap = make(map[string]string)

	// Retrieve and Filter files from GitHub
	getFileListFromGitHub(currentApiUrl)
	filterFileListFromGitHub()

	// Create the UI-list that holds the selected files
	generateSelectedFilesList(githubFileImporterWindow)

	// Create the UI-list that holds the filtered files and folders from GitHub
	generateFilteredList(githubFileImporterWindow)

	// Generate the File filter PopUp
	generateFileFilterPopup(githubFileImporterWindow)

	// Create a label with data binding used for showing current GitHub path
	pathLabel = widget.NewLabelWithData(currentPathShowedinGUI)

	// Generate the Button that moves upwards in the folder structure in GitHub
	generateMoveUpInFolderStructureButton()

	// Generate the button that imports the selected files from GitHub
	generateImportSelectedFilesFromGithubButton(githubFileImporterWindow)

	// Generate the button that cancel everything and closes the window
	generateCancelButton(githubFileImporterWindow)

	// Set initial size of the window
	githubFileImporterWindow.Resize(fyne.NewSize(400, 500))

	// Generate the wow that holds the up/back button and the path itself
	var pathRowContainer *fyne.Container
	pathRowContainer = container.NewHBox(moveUpInFolderStructureButton, pathLabel)

	// Create the top element which has the Filter button and the path.label and the back/upp button
	myTopLayout := container.NewVBox(fileFilterPopupButton, pathRowContainer)

	// Generate the container which has the filtered folders and files to the left and the selected files to the right
	splitContainer := container.NewHSplit(filteredFileList, selectedFilesList)
	splitContainer.Offset = 0.5 // Adjust if you want different initial proportions

	// Generate the row that has the import button and the cancel button
	var importCancelRow *fyne.Container
	importCancelRow = container.NewHBox(layout.NewSpacer(), importSelectedFilesFromGithubButton, cancelButton)

	// Crate the full content that should be plaved in the window
	content := container.NewBorder(myTopLayout, importCancelRow, nil, nil, splitContainer)

	// Set content
	githubFileImporterWindow.SetContent(content)

	// Set the callback function for window close event to show the Main window again
	githubFileImporterWindow.SetOnClosed(func() {
		*sharedResponseChannel <- false
		fenixMainWindow.Show()
	})

	// Show the githubFileImporterWindow
	githubFileImporterWindow.Show()

	return &selectedFiles

}

type TappableLabel struct {
	widget.Label
	onDoubleTap func()
	lastTap     time.Time
}

func NewTappableLabel(text *binding.String, onDoubleTap func()) *TappableLabel {
	labelWithData := widget.NewLabelWithData(*text)

	l := &TappableLabel{
		Label:       *labelWithData,
		onDoubleTap: onDoubleTap,
	}

	//l.Label = *widget.NewLabelWithData(text)

	l.ExtendBaseWidget(l)
	return l
}

func (l *TappableLabel) Tapped(e *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(l.lastTap) < 500*time.Millisecond { // 500 ms as double-click interval
		if l.onDoubleTap != nil {
			l.onDoubleTap()
		}
	}
	l.lastTap = now
}
