package cmd

import (
	"encoding/json"
	"fmt"
)

// Repository represents the structure of the repository data
type Repository struct {
	LocalSettings struct {
		DisplayName string `json:"displayName"`
	} `json:"localSettings"`
	Key string `json:"key"`
}

// findRepoKeyByDisplayName searches for a repository by its display name and returns its key
func findRepo(repos []Repository, repoName string) (string, error) {
	for _, repo := range repos {
		if repo.LocalSettings.DisplayName == repoName {
			return repo.Key, nil
		}
	}
	return "", fmt.Errorf("repository with display name '%s' not found", repoName)
}

// RepoKey fetches all repositories and finds the key for the given display name
func RepoKey(repoName string) (string, error) {
	// Fetch the repository data
	repoData, err := fetchRepo()
	if err != nil {
		return "", fmt.Errorf("failed to fetch repository data: %w", err)
	}

	// Parse the repository data
	var repos []Repository
	err = json.Unmarshal([]byte(repoData), &repos)
	if err != nil {
		return "", fmt.Errorf("failed to parse repository data: %w", err)
	}

	// Find the repository key by display name
	return findRepo(repos, repoName)
}
