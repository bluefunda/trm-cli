package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"net/http"
)

// Constants for configuration
const (
	csrfTokenURL = "https://abapdev.bluefunda.com:8080/rest/apim/v1/csrf-token" // Replace with actual URL
)

// getCSRFToken retrieves the CSRF token from the API and returns it as a string
func getCSRFToken() (string, error) {
	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, csrfTokenURL, nil)
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

	// Set the header to fetch the CSRF token
	req.Header.Set("X-Csrf-Token", "fetch")

	// Assuming httpClient is already defined
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve CSRF token: %s", resp.Status)
	}

	// Retrieve the CSRF token from the response headers
	csrfToken := resp.Header.Get("X-Csrf-Token")
	if csrfToken == "" {
		return "", fmt.Errorf("CSRF token is empty")
	}

	// Return the CSRF token
	return csrfToken, nil
}
