package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		if err := os.Setenv("BDA_TRM_TOKEN", token); err != nil {
			return fmt.Errorf("failed to set environment variable: %v", err)
		}
		fmt.Println("Login response:", string(body))
		return nil
	}

	return fmt.Errorf("login failed, status code: %d, response body: %s", resp.StatusCode, string(body))
}

func parseTokenResponse(responseBody []byte) (string, error) {
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var tokenRes TokenResponse
	if err := json.Unmarshal(responseBody, &tokenRes); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	return tokenRes.AccessToken, nil
}
