package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
)

// Base URL for fetching repository information
const (
	baseRepoURL = "https://abapdev.bluefunda.com:8080/rest/git/sap/v1/repo"
)

// fetchRepo fetches all repository data
func fetchRepo() (string, error) {
	// Use the base URL to fetch all repositories
	requestURL := baseRepoURL

	// Create a new HTTP request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Read the token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	return string(body), nil
}
