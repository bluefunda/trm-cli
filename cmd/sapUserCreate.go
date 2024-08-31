package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
)

// HTTPClient is a reusable HTTP client instance
var httpClient = &http.Client{}

// createUser creates a new user with the provided username and returns a response message
func createUser(username string) (string, error) {
	// Base URLs
	const (
		baseURL      = "https://abapdev.bluefunda.com:8080/rest/apim/v1/system/users"
		tempPassword = "Welcome123"
	)

	// User represents the structure of the user data to be sent in the request body
	type user struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	userData := user{
		UserName: username,
		Password: tempPassword,
	}

	requestBody, err := json.Marshal(userData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %v", err)
	}

	token, err := getCSRFToken()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Read the token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create user: %s, response: %s", resp.Status, string(bodyBytes))
	}

	return "User created successfully", nil
}
