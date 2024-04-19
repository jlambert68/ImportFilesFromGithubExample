package importFilesFromGitHub

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// Channel used to inform when this window closes and there are values to use
// True = files where picked
// False = No file where picked
var sharedResponseChannel *chan bool

// A pointer to Fenix Main Window
var fenixMainWindow fyne.Window

// The window for the File Importer
var githubFileImporterWindow fyne.Window

// THe root ApiUrl
var rootApiUrl string

// The current ApiUrl tp fetch files and folders from
var currentApiUrl string

var fileRegExFilterMap map[string]string

var githubFiles, githubFilesFiltered, selectedFiles []GitHubFile

// Create a string data binding
var currentPathShowedinGUI binding.String

var selectedFilesTable *widget.Table
var filteredFileList *widget.List
var fileFilterPopupButton *widget.Button

// Create a label with data binding used for showing current GitHub path
var pathLabel *widget.Label

// The Button that moves upwards in the folder structure in GitHub
var moveUpInFolderStructureButton *widget.Button

// The button that imports the selected files from GitHub
var importSelectedFilesFromGithubButton *widget.Button

// The button that cancel everything and closes the window
var cancelButton *widget.Button
