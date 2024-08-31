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

func pushGit(username, password, authorName, authorEmail, comment string) error {
	// Retrieve repo data from configuration or other methods
	repoConfig, err := config.ReadRepoConfig() // Assume this function exists to read repo config data
	if err != nil {
		return fmt.Errorf("error reading repo config: %v", err)
	}

	key, err := config.ReadToken("key")
	if err != nil {
		return fmt.Errorf("error reading key: %v", err)
	}

	// Split semicolon-separated values into slices
	paths := strings.Split(repoConfig.Path, ";")
	filenames := strings.Split(repoConfig.Filename, ";")
	datas := strings.Split(repoConfig.Data, ";")

	// Prepare the request body
	requestBody := PostRequestBody{
		URL: "https://github.com/bluefunda/abap-git.git", // This should be set according to your repo data
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
		repoData := struct {
			Filename   string `json:"filename"`
			Key        string `json:"key"`
			BranchName string `json:"branchName"`
			Package    string `json:"package"`
			Data       string `json:"data"`
		}{
			Filename:   filenames[i],
			Key:        key,
			BranchName: "mvp-1",            // Or fetch the branch name from the repo config
			Package:    repoConfig.Package, // Assuming Package is not semicolon-separated
			Data:       datas[i],
		}
		requestBody.RepoData = append(requestBody.RepoData, repoData)
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Retrieve the CSRF token
	token, err := getCSRFToken() // Define this function to get the CSRF token
	if err != nil {
		return fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	// Read the token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodPost, pushURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
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
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	fmt.Println("Response:", string(responseData))
	return nil
}
