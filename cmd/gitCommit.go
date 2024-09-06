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

type StatusData struct {
	Package string `json:"package"`
}

// ApiResponse represents the structure of the entire API response
type ApiResponse struct {
	ApiResult []struct {
		File   FileData   `json:"file"`
		Item   ItemData   `json:"item"`
		Status StatusData `json:"status"`
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
	// Read the base URL from the environment or config file
	baseURL, err := config.ReadToken("url")
	if err != nil || baseURL == "" {
		return fmt.Errorf("failed to retrieve base URL from config file")
	}

	// Concatenate the base URL with the endpoint
	stageURL := baseURL + "/rest/git/sap/v1/stage"

	key, gitURL, err := config.ReadKeyConfig()
	if err != nil {
		return fmt.Errorf("failed to retrieve repo key from config: %w", err)
	}

	fullURL := fmt.Sprintf("%s?url=%s&username=%s&password=%s&key=%s", stageURL, gitURL, username, password, key)

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
	var fileDataList, fileNameList, objectNameList, filePathList, packageList []string
	for _, result := range apiResponse.ApiResult {
		if objectName == "." || result.Item.ObjectName == objectName {
			fileDataList = append(fileDataList, result.File.Data)
			fileNameList = append(fileNameList, result.File.Filename)
			objectNameList = append(objectNameList, result.Item.ObjectName)
			filePathList = append(filePathList, result.File.Path)

			// Assuming that package information can be retrieved from result.Status
			packageList = append(packageList, result.Status.Package)
		}
	}

	// Store the slices in envVars
	envVars := map[string]string{
		"FILE_DATA":   strings.Join(fileDataList, ";"),
		"FILE_NAME":   strings.Join(fileNameList, ";"),
		"OBJECT_NAME": strings.Join(objectNameList, ";"),
		"FILE_PATH":   strings.Join(filePathList, ";"),
		"PACKAGE":     strings.Join(packageList, ";"),
	}

	// Update the configuration
	if err := config.UpdateEnvVars("repo", envVars); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}

	return nil
}
