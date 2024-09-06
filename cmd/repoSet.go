package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bluefunda/trm-cli/config"
)

// Repository represents the structure of the repository data
type Repository struct {
	LocalSettings struct {
		DisplayName string `json:"displayName"`
	} `json:"localSettings"`
	Key string `json:"key"`
	URL string `json:"url"`
}

// RepoData represents the structure of the entire JSON response
type RepoData struct {
	Repo []Repository `json:"repo"`
}

// findRepo searches for a repository by its display name and returns its key and URL
func findRepo(repos []Repository, repoName string) (string, string, error) {
	for _, repo := range repos {
		if repo.LocalSettings.DisplayName == repoName {
			return repo.Key, repo.URL, nil
		}
	}
	return "", "", fmt.Errorf("repository with display name '%s' not found", repoName)
}

// RepoKey fetches all repositories and finds the key and URL for the given display name,
// then updates the token with that key and URL.
func RepoKey(repoName string) error {
	// Fetch the repository data
	repoData, err := fetchRepo()
	if err != nil {
		return fmt.Errorf("failed to fetch repository data: %w", err)
	}

	// Parse the repository data
	var data RepoData
	err = json.Unmarshal([]byte(repoData), &data)
	if err != nil {
		return fmt.Errorf("failed to parse repository data: %w", err)
	}

	// Find the repository key and URL by display name
	key, url, err := findRepo(data.Repo, repoName)
	if err != nil {
		return err
	}

	// Create a map with the key-value pair(s)
	envVars := map[string]string{
		"KEY":     key,
		"GIT_URL": url,
	}

	// Update the token with the found key and URL
	if err := config.UpdateEnvVars("key", envVars); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	// On success, return nil (no message)
	return nil
}
