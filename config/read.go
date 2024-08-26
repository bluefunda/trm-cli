package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Define configPath globally
var configFilePath = filepath.Join(getDir(), ".config", "bda", "token")

func getDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE") // for Windows
}

// ReadPath reads configuration from the token file
func readPath() (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(configFilePath)
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

// GetToken retrieves the token value from the configuration file
func ReadToken() (string, error) {
	config, err := readPath()
	if err != nil {
		return "", err
	}

	token, exists := config["BDA_TRM_TOKEN"]
	if !exists {
		return "", errors.New("token not found in configuration")
	}

	return token, nil
}
