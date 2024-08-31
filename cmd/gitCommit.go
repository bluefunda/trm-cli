package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"io"
	"net/http"
	"os"
	"strings"
)

const stageURL = "https://abapdev.bluefunda.com:8080/rest/git/sap/v1/stage?url=https://github.com/bluefunda/abap-git.git"

// FileData represents the structure of the file data in the API response
type FileData struct {
	Data     string `json:"data"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

// ItemData represents the structure of the item data in the API response
type ItemData struct {
	ObjectName string `json:"objectName"`
}

// ApiResponse represents the structure of the entire API response
type ApiResponse struct {
	ApiResult []struct {
		File FileData `json:"file"`
		Item ItemData `json:"item"`
	} `json:"apiResult"`
}

func promptForCredentials(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}

// fetchData retrieves data from the specified URL with Bearer token authorization
func fetchData(fullURL string) (string, error) {
	// Read the token from the config
	bearerToken, err := config.ReadToken("token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token from env file: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Perform the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// StoreData processes the API response and stores key-value pairs in the configuration file
func StoreData(objectName, username, password string) error {
	key, err := config.ReadToken("key")
	if err != nil {
		return fmt.Errorf("failed to retrieve repo key from env file: %w", err)
	}

	fullURL := fmt.Sprintf("%s&username=%s&password=%s&key=%s", stageURL, username, password, key)

	// Debugging: Print the full URL
	fmt.Println("Full URL:", fullURL)

	responseData, err := fetchData(fullURL)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	if len(responseData) == 0 {
		return fmt.Errorf("received empty response from the API")
	}

	var apiResponse ApiResponse
	err = json.Unmarshal([]byte(responseData), &apiResponse)
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	// Collect data in slices
	var fileDataList, fileNameList, objectNameList, filePathList []string
	for _, result := range apiResponse.ApiResult {
		if objectName == "." || result.Item.ObjectName == objectName {
			fileDataList = append(fileDataList, result.File.Data)
			fileNameList = append(fileNameList, result.File.Filename)
			objectNameList = append(objectNameList, result.Item.ObjectName)
			filePathList = append(filePathList, result.File.Path)
		}
	}

	// Store the slices in envVars
	envVars := map[string]string{
		"FILE_DATA":   strings.Join(fileDataList, ";"), // Assuming `;` as a delimiter
		"FILE_NAME":   strings.Join(fileNameList, ";"),
		"OBJECT_NAME": strings.Join(objectNameList, ";"),
		"FILE_PATH":   strings.Join(filePathList, ";"),
	}

	// Update the configuration
	if err := config.UpdateEnvVars("repo", envVars); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}

	return nil
}
