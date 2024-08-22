package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Constants for configuration
const (
	csrfTokenURL = "https://example.com/api/csrf-token" // Replace with actual URL
	csrfTokenEnv = "CSRF_TOKEN"
)

// getCSRFToken retrieves the CSRF token from the API and stores it in an environment variable
func getCSRFToken() error {
	req, err := http.NewRequest(http.MethodGet, csrfTokenURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-CSRF-Token", "fetch")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve CSRF token: %s", resp.Status)
	}

	var tokenResponse struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to decode CSRF token response: %v", err)
	}

	if tokenResponse.Token == "" {
		return fmt.Errorf("CSRF token is empty")
	}

	if err := os.Setenv(csrfTokenEnv, tokenResponse.Token); err != nil {
		return fmt.Errorf("failed to set CSRF token environment variable: %v", err)
	}

	// Optionally verify that the environment variable is set correctly
	if token := os.Getenv(csrfTokenEnv); token != tokenResponse.Token {
		return fmt.Errorf("unexpected value for CSRF token environment variable")
	}

	fmt.Println("CSRF token retrieved and set successfully")
	return nil
}
