package importFilesFromGitHub

// Struct for parsing JSON response
type GitHubFile struct {
	Name                string `json:"name"`
	Type                string `json:"type"` // "file" or "dir"
	URL                 string `json:"url"`  // URL to fetch contents if it's a directory
	DownloadURL         string `json:"download_url"`
	Content             []byte `json:"content"`
	SHA                 string `json:"sha"`
	Size                int    `json:"size"`
	FileContentAsString string
}
type GitHubFileDetail struct {
	Name                string `json:"name"`
	Path                string `json:"path"`
	URL                 string `json:"url"`
	DownloadURL         string `json:"download_url"`
	Type                string `json:"type"`
	Content             string `json:"content"`
	Encoding            string `json:"encoding"`
	SHA                 string `json:"sha"`
	Size                int    `json:"size"`
	FileContentAsString string
	// Include other fields as needed
}
