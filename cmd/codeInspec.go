package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
)

// sendRequest sends a POST request for code inspection
func codeInspect(pkg string, objNames []string) (string, error) {
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return "", fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Concatenate the base URL with the endpoint
	url := baseURL + "/rest/git/sap/v1/code-inspector"

	// RequestBody represents the structure of the data to be sent in the request body
	type request struct {
		ObjectName []string `json:"objectName"`
		Package    string   `json:"package"`
		Variant    string   `json:"variant"`
	}

	// Create the request body with provided objectName and package
	requestData := request{
		ObjectName: objNames,
		Package:    pkg,
		Variant:    "DEFAULT", // Assuming variant is always "DEFAULT"
	}

	// Marshal the request data to JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %v", err)
	}

	// Create a new POST request with the JSON body
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Get CSRF token
	token, err := getCSRFToken()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	// Read the bearer token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrf-token", token)

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and handle the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %s, response: %s", resp.Status, string(respBody))
	}

	return string(respBody), nil
}
