package importFilesFromGitHub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
)

// Generate the button that imports the selected files from Github
func generateImportSelectedFilesFromGithubButton(parentWindow fyne.Window) {

	importSelectedFilesFromGithubButton = widget.NewButton("Import Files", func() {
		for fileIndex, file := range selectedFiles {
			content, err := loadFileContent(file)
			if err != nil {
				dialog.ShowError(err, parentWindow)
				continue
			}
			// Do something with the content, e.g., display it, process it, etc.
			selectedFiles[fileIndex].Content = content

			extractedContent, err := extractContentFromJson(string(content))
			if err != nil {
				log.Fatalf("Error parsing JSON: %s", err)
			}

			contentAsString, err := decodeBase64Content(string(extractedContent))
			if err != nil {
				log.Fatalf("Failed to decode content: %s", err)
			}
			// 'content' now contains the decoded file content as a string
			fmt.Println(contentAsString)

			// Save the file content
			file.fileContetAsString = contentAsString
		}

		fenixMainWindow.Show()
		parentWindow.Close()
	})
}

// Extra the file content from the json
func extractContentFromJson(jsonData string) (string, error) {
	var fileDetail GitHubFileDetail
	err := json.Unmarshal([]byte(jsonData), &fileDetail)
	if err != nil {
		return "", err
	}

	return fileDetail.Content, nil
}

// Decode the file content from base64 to string
func decodeBase64Content(encodedContent string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
