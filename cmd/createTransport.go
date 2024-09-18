package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"strings"
)

// TransportData represents the structure of the transport data to be sent in the request body
type TransportData struct {
	RequestType string `json:"requestType"`
	Author      string `json:"author"`
	Text        string `json:"text"`
}

// createTransport creates a new transport request in SAP with the given requestType, author, and text
func createTransport(requestType, author, text string) (string, error) {
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return "", fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Concatenate the base URL with the endpoint
	baseTransportURL := baseURL + "/rest/trm/v1/transports"

	// Create the transport data
	transportData := TransportData{
		RequestType: requestType,
		Author:      author,
		Text:        text,
	}

	// Marshal the transport data to JSON
	requestBody, err := json.Marshal(transportData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transport data: %v", err)
	}

	// Print the request body
	fmt.Printf("Request body: %s\n", string(requestBody))

	// Retrieve the CSRF token
	token, cookies, err := getCSRFToken()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CSRF token: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, baseTransportURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Read the Bearer token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token from config file: %w", err)
	}

	// Set request headers
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("x-csrf-token", token)
	// Set the cookies in the request header
	var cookieStrings []string
	for name, value := range cookies {
		cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", name, value))
	}
	req.Header.Set("Cookie", strings.Join(cookieStrings, "; "))

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create transport: %s, response: %s", resp.Status, string(bodyBytes))
	}

	// Return success message
	return "Transport created successfully", nil
}
