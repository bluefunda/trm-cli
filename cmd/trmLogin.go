package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	clientID     = "admin-trm"
	clientSecret = "O8xQ9P0tjMPIEJqtzEWNsqlvCwgpxt9I"
)

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}

func login(username, password string) error {
	// URL for the token endpoint
	tokenURL := "https://abapdev.bluefunda.com:7443/realms/trm/protocol/openid-connect/token"

	// Create form data
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")

	// Make the HTTP request to get the token
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	// Read and handle response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		token, err := parseTokenResponse(body)
		if err != nil {
			return fmt.Errorf("failed to parse token response: %v", err)
		}

		// Create a map with the key-value pair(s)
		envVars := map[string]string{
			"BDA_TRM_TOKEN": token,
		}

		// Call UpdateEnvVars with the map
		if err := config.UpdateEnvVars("login", envVars); err != nil {
			return fmt.Errorf("failed to update token in config file: %w", err)
		}

		fmt.Println("Login response:", string(body))
		fmt.Println("Token updated successfully!")
		return nil
	}

	return fmt.Errorf("login failed, status code: %d, response body: %s", resp.StatusCode, string(body))
}

func parseTokenResponse(responseBody []byte) (string, error) {
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var tokenRes tokenResponse
	if err := json.Unmarshal(responseBody, &tokenRes); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	return tokenRes.AccessToken, nil
}
