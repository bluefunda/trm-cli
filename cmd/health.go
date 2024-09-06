package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"time"
)

// Create an HTTP client with a timeout
var client = http.Client{
	Timeout: 5 * time.Second,
}

func health() (string, error) {
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return "", fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Concatenate the base URL with the endpoint
	healthURL := baseURL + "/__health"
	resp, err := client.Get(healthURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		// Directly return the body as the response if healthy
		return string(body), nil
	}
	return "", fmt.Errorf("service is not healthy. Status code: %d", resp.StatusCode)
}
