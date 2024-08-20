package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Base URL for fetching users
const (
	baseUserURL = "https://abapdev.bluefunda.com:8080/DEV/rest/apim/v1/system/users"
)

// fetchUsers fetches users data or a specific user if userID is provided
func fetchUsers(userID string) (string, error) {
	// Construct the URL based on whether userID is provided
	var requestURL string
	if userID == "" {
		// Fetch all users
		requestURL = baseUserURL
	} else {
		// Fetch specific user
		// Escape the userID to be URL-safe
		escapedUserID := url.QueryEscape(userID)
		requestURL = fmt.Sprintf("%s?userName=%s", baseUserURL, escapedUserID)
	}

	// Make an HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	return string(body), nil
}
