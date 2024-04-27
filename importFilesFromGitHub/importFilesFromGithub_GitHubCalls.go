package importFilesFromGitHub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// List files and folders for a certain GitHub url
func getFileListFromGitHub(apiUrl string) {

	var tempGithubFiles []GitHubFile
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %s", err.Error())
	}

	// Add the API token in the request header
	apiToken := gitHubApiKey // Replace with your actual API token
	req.Header.Add("Authorization", "token "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occurred while calling GitHub API: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error in response from Github: Status code '%d' with message '%s'", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading API response: %s", err.Error())
	}

	if err := json.Unmarshal(body, &tempGithubFiles); err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err.Error())
	}

	githubFiles = tempGithubFiles
}

// Load the files content from GitHub
func loadFileContent(file GitHubFile) ([]byte, error) {
	// Assuming file.URL is the URL to the raw content of the file
	resp, err := http.Get(file.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
