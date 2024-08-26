package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var configPath = filepath.Join(getHomeDir(), ".config", "bda", "token")

func getHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE") // for Windows
}

// EnsureConfigDir ensures that the directory for the config file exists
func ensureConfigDir() error {
	dir := filepath.Dir(configPath)

	// Check if the directory already exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Prompt the user for permission to create the directory
		fmt.Printf("Directory %s does not exist. Do you want to create it? (y/n): ", dir)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("user did not grant permission to create directory")
		}

		// Create the directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}
		fmt.Println("Directory created successfully")
	}
	return nil
}

// ReadConfig reads configuration from the token file
func readConfigPath() (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
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

// WriteConfig writes configuration to the token file
func writeConfigPath(config map[string]string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range config {
		_, err := writer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// UpdateToken updates the BDA_TRM_TOKEN in the config file, creating the file if it doesn't exist
func UpdateToken(newToken string) error {
	// Ensure the config directory exists
	if err := ensureConfigDir(); err != nil {
		return fmt.Errorf("error ensuring config directory: %v", err)
	}

	// Check if the config file exists, and create it if it doesn't
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file does not exist. Creating a new one.")
		// Create an empty config map
		config := make(map[string]string)
		// Set the new token
		config["BDA_TRM_TOKEN"] = newToken
		// Write the new config file
		return writeConfigPath(config)
	}

	// Read existing config
	config, err := readConfigPath()
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Update config values
	config["BDA_TRM_TOKEN"] = newToken

	// Write updated config
	err = writeConfigPath(config)
	if err != nil {
		return fmt.Errorf("error updating config file: %v", err)
	}

	return nil
}
