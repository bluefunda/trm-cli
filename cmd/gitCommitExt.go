package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Manage Git repositories",
}

// gitAdd represents the add command under git
var gitAdd = &cobra.Command{
	Use:   "add [object-name]",
	Short: "Set Git repositories",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		objectName := args[0] // Get the repository name from arguments

		// Check if the username is stored in the config file
		username, err := config.ReadToken("gitUser")
		if err != nil || username == "" {
			// Prompt the user for the username if not found
			username = promptForCredentials("Enter your GitHub username: ")

			// Store the username in the config file
			envVars := map[string]string{
				"GIT_USER": username,
			}
			if err := config.UpdateEnvVars("gitUser", envVars); err != nil {
				fmt.Printf("Failed to update username in config file: %v\n", err)
			}
		}

		// Prompt for password
		password := promptForCredentials("Enter GitHub password or token: ")

		// Call the function to fetch the repository key with credentials
		err = StoreData(objectName, username, password)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	// Add gitAdd as a subcommand of gitCmd
	gitCmd.AddCommand(gitAdd)

	// Add gitCmd as a subcommand of rootCmd
	rootCmd.AddCommand(gitCmd)
}
