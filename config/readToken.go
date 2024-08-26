package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// ReadToken loads the token from the config file
func ReadToken() (string, error) {
	// Load config.env file
	err := godotenv.Load("config/config.env")
	if err != nil {
		return "", fmt.Errorf("error loading config.env file: %w", err)
	}

	// Get the token value from the environment variable
	token := os.Getenv("BDA_TRM_TOKEN")
	if token == "" {
		return "", fmt.Errorf("BDA_TRM_TOKEN not set in config.env file")
	}

	return token, nil
}
