package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE") // for Windows
}

var (
	loginConfigPath = filepath.Join(getHomeDir(), ".config", "bda", "token")
	repoConfigPath  = filepath.Join(getHomeDir(), ".config", "bda", "repo")
	keyConfigPath   = filepath.Join(getHomeDir(), ".config", "bda", "key")
	gitAuthPath     = filepath.Join(getHomeDir(), ".config", "bda", "gitUser")
)

func ensureConfigDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist. Do you want to create it? (y/n): ", dir)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("user did not grant permission to create directory")
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}
		fmt.Println("Directory created successfully")
	}
	return nil
}

func readConfigPath(path string) (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(path)
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

func writeConfigPath(path string, config map[string]string) error {
	file, err := os.Create(path)
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

// UpdateEnvVars updates multiple environment variables in the specified config file
func UpdateEnvVars(configType string, newVars map[string]string) error {
	var configPath string

	// Determine config path based on the configType input
	switch configType {
	case "login":
		configPath = loginConfigPath
	case "repo":
		configPath = repoConfigPath
	case "key":
		configPath = keyConfigPath
	case "gitUser":
		configPath = gitAuthPath
	default:
		return fmt.Errorf("invalid config type: %s", configType)
	}

	// Ensure the config directory exists
	if err := ensureConfigDir(configPath); err != nil {
		return fmt.Errorf("error ensuring config directory: %v", err)
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file does not exist for %s.\n", configType)
		fmt.Printf("Do you want to create the config file at %s? (y/n): ", configPath)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("user did not grant permission to create the config file")
		}

		// Write the new config file with provided key-value pairs
		if err := writeConfigPath(configPath, newVars); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}

		fmt.Println("Config file created successfully")
		return nil
	}

	// Read existing config
	config, err := readConfigPath(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Update config values with the new key-value pairs
	for key, value := range newVars {
		config[key] = value
	}

	// Write updated config
	if err := writeConfigPath(configPath, config); err != nil {
		return fmt.Errorf("error updating config file: %v", err)
	}

	return nil
}
