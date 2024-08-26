package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const configFilePath = "config/config.env"

// ReadConfig reads configuration from the .env file
func readConfig() (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(configFilePath)
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

// WriteConfig writes configuration to the .env file
func writeConfig(config map[string]string) error {
	file, err := os.Create(configFilePath)
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

// UpdateToken updates the BDA_TRM_TOKEN in the config file
func UpdateToken(newToken string) error {
	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %v", err)
	}

	// Read existing config
	config, err := readConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Update config values
	config["BDA_TRM_TOKEN"] = newToken

	// Write updated config
	err = writeConfig(config)
	if err != nil {
		return fmt.Errorf("error updating config file: %v", err)
	}

	return nil
}
