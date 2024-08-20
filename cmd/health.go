package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const healthURL = "https://abapdev.bluefunda.com:8080/__health"

// Create an HTTP client with a timeout
var client = http.Client{
	Timeout: 5 * time.Second,
}

func health() (string, error) {
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
