package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bluefunda/trm-cli/config"
)

const pushURL = "https://abapdev.bluefunda.com:8080/rest/git/sap/v1/push"

// PostRequestBody represents the structure of the request body
type PostRequestBody struct {
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials"`
	Comment struct {
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Comment string `json:"comment"`
	} `json:"comment"`
	Path     string `json:"path"`
	URL      string `json:"url"`
	RepoData []struct {
		Filename   string `json:"filename"`
		Key        string `json:"key"`
		BranchName string `json:"branchName"`
		Package    string `json:"package"`
		Data       string `json:"data"`
	} `json:"repoData"`
}

func promptInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func pushGit(username, password, authorName, authorEmail, comment string) (int, error) {
	// Retrieve repo data from configuration or other methods
	repoConfig, err := config.ReadRepoConfig() // Assume this function exists to read repo config data
	if err != nil {
		return 0, fmt.Errorf("error reading repo config: %v", err)
	}

	key, err := config.ReadToken("key")
	if err != nil {
		return 0, fmt.Errorf("error reading key: %v", err)
	}

	// Split semicolon-separated values into slices
	paths := strings.Split(repoConfig.Path, ";")
	filenames := strings.Split(repoConfig.Filename, ";")
	datas := strings.Split(repoConfig.Data, ";")
	packages := strings.Split(repoConfig.Package, ";") // Split the Package field by comma

	// Prepare the request body
	requestBody := PostRequestBody{
		URL: "https://github.com/bluefunda/abap-apim.git", // This should be set according to your repo data
	}
	requestBody.Credentials.Username = username
	requestBody.Credentials.Password = password
	requestBody.Comment.Committer.Name = authorName
	requestBody.Comment.Committer.Email = authorEmail
	requestBody.Comment.Author.Name = authorName
	requestBody.Comment.Author.Email = authorEmail
	requestBody.Comment.Comment = comment

	// Loop through the slices and populate RepoData
	for i := range paths {
		// Check if the current index is valid for filenames, datas, paths, and packages slices
		if i < len(filenames) && i < len(datas) && i < len(paths) && i < len(packages) {
			repoData := struct {
				Filename   string `json:"filename"`
				Key        string `json:"key"`
				BranchName string `json:"branchName"`
				Package    string `json:"package"`
				Data       string `json:"data"`
			}{
				Filename:   filenames[i],
				Key:        key,
				BranchName: "mvp-1",     // Or fetch the branch name from the repo config
				Package:    packages[i], // Use the corresponding package value
				Data:       datas[i],
			}

			// Assign the current path to the Path field in repoData
			requestBody.Path = paths[i]

			// Add repoData to the requestBody's RepoData slice
			requestBody.RepoData = append(requestBody.RepoData, repoData)
		}
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Retrieve the CSRF token
	token, err := getCSRFToken() // Define this function to get the CSRF token
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	// Read the token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", pushURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", token)

	// Make the POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("error reading response: %v", err)
	}

	// Print the request body for debugging
	// Pretty-print the request body for debugging
	prettyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Printf("error pretty printing JSON: %v\n", err)
		return resp.StatusCode, nil
	}

	fmt.Println("RequestBody:", string(prettyJSON))

	fmt.Println("Response:", string(responseData))
	return resp.StatusCode, nil
}
