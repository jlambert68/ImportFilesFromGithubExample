package importFilesFromGitHub

// Struct for parsing JSON response
type GitHubFile struct {
	Name               string `json:"name"`
	Type               string `json:"type"` // "file" or "dir"
	URL                string `json:"url"`  // URL to fetch contents if it's a directory
	Content            []byte `json:"content"`
	fileContetAsString string
}
type GitHubFileDetail struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	// Include other fields as needed
}
