package importFilesFromGitHub

import (
	"fmt"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// Generate the Button that moves upwards in the folder structure in GitHub
func generateMoveUpInFolderStructureButton() {

	moveUpInFolderStructureButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		// Handle the button click - go back in your navigation, for instance

		if currentApiUrl == rootApiUrl {
			return
		}

		newPath, err := moveUpInPath(currentApiUrl)
		if err == nil || len(newPath) > 0 {
			currentApiUrl = newPath

			currentPathShowedinGUI.Set(strings.Split(currentApiUrl, "?")[0])
			getFileListFromGitHub(currentApiUrl)
			filterFileListFromGitHub()
			filteredFileList.Refresh() // Refresh the list to update it with the new contents

		}
	})
}

// Move one step in the folder structure
func moveUpInPath(currentPath string) (string, error) {
	// Trim any trailing slashes
	trimmedPath := strings.TrimRight(currentPath, "/")

	// Split the path into components
	pathComponents := strings.Split(trimmedPath, "/")

	// Check if it's already at the root or has no parent
	if len(pathComponents) <= 1 {
		return "", fmt.Errorf("already at the root or invalid path")
	}

	// Remove the last component to move up one directory
	newPathComponents := pathComponents[:len(pathComponents)-1]

	// Join components back into a path
	newPath := strings.Join(newPathComponents, "/")
	if newPath == "" {
		newPath = "/" // Ensure root is represented correctly
	}

	return newPath, nil
}
