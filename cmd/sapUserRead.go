package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"net/url"
)

// Base URL for fetching users
const (
	baseUserURL = "https://abapdev.bluefunda.com:8080/rest/apim/v1/system/users"
)

// fetchUsers fetches users data or a specific user if userID is provided
func fetchUsers(userID string) (string, error) {
	// Construct the URL based on whether userID is provided
	var requestURL string
	if userID == "all" {
		// Fetch all users
		requestURL = baseUserURL
	} else {
		// Fetch specific user
		// Escape the userID to be URL-safe
		escapedUserID := url.QueryEscape(userID)
		requestURL = fmt.Sprintf("%s?userName=%s", baseUserURL, escapedUserID)
	}

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
