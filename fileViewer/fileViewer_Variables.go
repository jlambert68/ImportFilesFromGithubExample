package fileViewer

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"fyne.io/fyne/v2"
)

// A pointer to Fenix Main Window
var fenixMainWindow fyne.Window

// The window for the File Viewer
var fileViewerWindow fyne.Window

var importedFiles []importFilesFromGitHub.GitHubFile
