package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	username string
	password string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "trm-login",
	Short: "CLI tool for logging into the application",
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the application",
	Run: func(cmd *cobra.Command, args []string) {
		username = promptForInput("Enter username: ")
		password = promptForInput("Enter password: ")

		if err := login(username, password); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Login successful!")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func login(username, password string) error {
	// URL for the login endpoint
	loginURL := "https://abapdev.bluefunda.com/trm"

	// Create form data
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	// Make the HTTP request
	resp, err := http.PostForm(loginURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and handle response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		token := extractTokenFromResponse(body)
		if token != "" {
			if err := os.Setenv("BDA_TRM_TOKEN", token); err != nil {
				return fmt.Errorf("failed to set environment variable: %v", err)
			}
		}
		fmt.Println("Login response:", string(body))
		return nil
	}

	return fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
}

func extractTokenFromResponse(responseBody []byte) string {
	// This function should parse the response body and extract the token.
	type Response struct {
		Token string `json:"token"`
	}

	var res Response
	if err := json.Unmarshal(responseBody, &res); err != nil {
		fmt.Println("Failed to parse response:", err)
		return ""
	}

	return res.Token
}
