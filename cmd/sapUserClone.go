package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"net/http"
)

// cloneUser creates a new user with the provided username and returns a response message
func cloneUser(usernameFrom, username string) (string, error) {
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return "", fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Concatenate the base URL with the endpoint
	baseCloneURL := baseURL + "/rest/apim/v1/system/users"

	// Define your User struct
	type userclone struct {
		UserNameFrom string `json:"userNameFrom"`
		UserName     string `json:"userName"`
		Password     string `json:"password"`
	}
	// Placeholder for password, replace with actual values or logic
	tempPassword := "Welcome123"

	userData := userclone{
		UserNameFrom: usernameFrom,
		UserName:     username,
		Password:     tempPassword,
	}

	requestBody, err := json.Marshal(userData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %v", err)
	}

	token, err := getCSRFToken()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseCloneURL, bytes.NewBuffer(requestBody))
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
	req.Header.Set("X-CSRF-Token", token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Optional: Read and log response body for more details
		var respBody bytes.Buffer
		if _, err := respBody.ReadFrom(resp.Body); err == nil {
			return "", fmt.Errorf("failed to create user: %s, response: %s", resp.Status, respBody.String())
		}
		return "", fmt.Errorf("failed to create user: %s", resp.Status)
	}

	return "User created successfully", nil
}
