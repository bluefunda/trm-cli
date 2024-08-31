package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Define the base directory globally
var baseDir = filepath.Join(getDir(), ".config", "bda")

func getDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE") // for Windows
}

// readPath reads configuration from the specified file in the base directory
func readPath(fileName string) (map[string]string, error) {
	config := make(map[string]string)
	filePath := filepath.Join(baseDir, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// ReadToken reads a token or key from the specified file and returns a map with the value
func ReadToken(input string) (string, error) {
	var fileName string
	// Determine the file name based on the input
	switch input {
	case "token":
		fileName = "token"
	case "key":

		fileName = "key"
	case "repo":
		fileName = "repo"
	case "gitUser":
		fileName = "gitUser"
	default:
		return "", errors.New("invalid input specified")
	}

	config, err := readPath(fileName)
	if err != nil {
		return "", err
	}

	// Return the relevant configuration
	switch input {
	case "token":
		// If input is empty, return only the BDA_TRM_TOKEN value
		if token, exists := config["BDA_TRM_TOKEN"]; exists {
			return token, nil
		}
		return "", errors.New("BDA_TRM_TOKEN not found in configuration")
	case "key":
		// If input is "key", return the KEY value
		if key, exists := config["KEY"]; exists {
			return key, nil
		}
	case "gitUser":
		// If input is "key", return the KEY value
		if key, exists := config["GIT_USER"]; exists {
			return key, nil
		}
		return "", errors.New("KEY not found in configuration")
	}

	return "", errors.New("unexpected error occurred")
}

// ReadRepoConfig reads repository configuration from the specified file and returns a map with key-value pairs
func ReadRepoConfig() (map[string]string, error) {
	fileName := "repo"

	config, err := readPath(fileName)
	if err != nil {
		return nil, err
	}

	// Return FILE_DATA, FILE_NAME, OBJECT_NAME, FILE_PATH, and KEY
	envVars := make(map[string]string)
	keys := []string{"FILE_DATA", "FILE_NAME", "OBJECT_NAME", "FILE_PATH", "KEY"}
	for _, key := range keys {
		if value, exists := config[key]; exists {
			envVars[key] = value
		} else {
			return nil, errors.New(key + " not found in configuration")
		}
	}
	return envVars, nil
}
