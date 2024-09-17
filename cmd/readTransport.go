package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"net/url"
)

// fetchTransportData fetches transport data based on requestID, username, status, and function
func fetchTransportData(requestID, username, transportStatus, transportFunction string) (string, error) {
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return "", fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Base transport URL
	transportURL := baseURL + "/rest/trm/v1/transports"

	// Construct the URL with query parameters
	params := url.Values{}
	params.Add("requestid", requestID)
	params.Add("username", username)
	params.Add("transportStatus", transportStatus)
	params.Add("transportFunction", transportFunction)

	// Final URL with query parameters
	requestURL := fmt.Sprintf("%s?%s", transportURL, params.Encode())

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
