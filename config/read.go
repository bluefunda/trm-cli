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
		if token, exists := config["BDA_TRM_TOKEN"]; exists {
			return token, nil
		}
		return "", errors.New("BDA_TRM_TOKEN not found in configuration")
	case "key":
		if key, exists := config["KEY"]; exists {
			return key, nil
		}
		return "", errors.New("KEY not found in configuration")
	case "gitUser":
		if user, exists := config["GIT_USER"]; exists {
			return user, nil
		}
		return "", errors.New("GIT_USER not found in configuration")
	}

	return "", errors.New("unexpected error occurred")
}

// RepoConfig holds the repository configuration
type RepoConfig struct {
	Path     string
	Filename string
	Data     string
	Package  string
}

// ReadRepoConfig reads repository configuration from the specified file and returns a RepoConfig struct
func ReadRepoConfig() (RepoConfig, error) {
	fileName := "repo"

	config, err := readPath(fileName)
	if err != nil {
		return RepoConfig{}, err
	}

	return RepoConfig{
		Path:     config["FILE_PATH"],
		Filename: config["FILE_NAME"],
		Data:     config["FILE_DATA"],
		Package:  config["OBJECT_NAME"],
	}, nil
}
