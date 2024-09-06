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

	// Use bufio.Reader instead of bufio.Scanner
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// End of file
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Ignore malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	return config, nil
}

// ReadToken reads a token or key from the specified file and returns a map with the value
// ReadToken reads a token or key from the specified file and returns the value as a string
func ReadToken(input string) (string, error) {
	var fileName string
	var key string

	// Determine the file name and key based on the input
	switch input {
	case "token":
		fileName = "token"
		key = "BDA_TRM_TOKEN"
	case "gitUser":
		fileName = "gitUser"
		key = "GIT_USER"
	case "url":
		fileName = "url"
		key = "URL"
	default:
		return "", errors.New("invalid input specified")
	}

	config, err := readPath(fileName)
	if err != nil {
		return "", err
	}

	if value, exists := config[key]; exists {
		return value, nil
	}

	return "", errors.New(key + " not found in configuration")
}

// ReadKeyConfig reads key and Git URL from the "key" configuration file
func ReadKeyConfig() (string, string, error) {
	fileName := "key"

	config, err := readPath(fileName)
	if err != nil {
		return "", "", err
	}

	key, keyExists := config["KEY"]
	gitURL, urlExists := config["GIT_URL"]

	if keyExists && urlExists {
		return key, gitURL, nil
	}
	if !keyExists {
		return "", "", errors.New("KEY not found in configuration")
	}
	if !urlExists {
		return "", "", errors.New("GIT_URL not found in configuration")
	}

	return "", "", errors.New("unexpected error occurred")
}

// RepoConfig holds the repository configuration
type RepoConfig struct {
	Path       string
	Filename   string
	Data       string
	ObjectName string
	Package    string
}

// ReadRepoConfig reads repository configuration from the specified file and returns a RepoConfig struct
func ReadRepoConfig() (RepoConfig, error) {
	fileName := "repo"

	config, err := readPath(fileName)
	if err != nil {
		return RepoConfig{}, err
	}

	return RepoConfig{
		Path:       config["FILE_PATH"],
		Filename:   config["FILE_NAME"],
		Data:       config["FILE_DATA"],
		ObjectName: config["OBJECT_NAME"],
		Package:    config["PACKAGE"],
	}, nil
}
