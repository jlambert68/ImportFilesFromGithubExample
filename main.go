package main

import (
	"ImportFilesFromGithub/importFilesFromGitHub"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

func main() {

	repoOwner := "jlambert68"    // Replace with the repository owner's username
	repoName := "FenixTesterGui" // Replace with the repository name
	repoPath := ""               // Replace with the path in the repository, if any

	originalApiUrl := "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/contents" + repoPath

	var myApp fyne.App
	myApp = app.New()

	var myMainWindow fyne.Window
	myMainWindow = myApp.NewWindow("This is Fenix Main Window")

	// Set initial size of the window
	myMainWindow.Resize(fyne.NewSize(400, 500))

	myButton := widget.NewButton("Import files from GitHub", func() {
		myMainWindow.Hide()
		_ = importFilesFromGitHub.InitiateImportFilesFromGitHubWindow(originalApiUrl, myMainWindow, myApp)
	})

	myVBoxContainer := container.NewVBox(myButton)
	myMainWindow.SetContent(myVBoxContainer)

	go func() {
		time.Sleep(2 * time.Second)

	}()

	myMainWindow.ShowAndRun()

}
