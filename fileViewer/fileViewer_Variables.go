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

const testCaseExecutionUuid string = "07f8c5db-5a2a-4f1a-87ca-0c2e11f747a2"
