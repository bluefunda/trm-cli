package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"github.com/spf13/cobra"
)

// gitPush represents the push command under git
var gitPush = &cobra.Command{
	Use:   "push",
	Short: "Push changes to Git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt for credentials and commit details
		username, err := config.ReadToken("gitUser")
		if err != nil {
			fmt.Println("Error reading Git username:", err)
			return
		}

		password := promptInput("Enter GitHub password: ")
		authorName := promptInput("Enter author name: ")
		authorEmail := promptInput("Enter author email: ")
		comment := promptInput("Enter comment: ")

		// Call the function to push the changes
		statusCode, err := pushGit(username, password, authorName, authorEmail, comment)
		if err != nil {
			fmt.Printf("Error pushing to Git: %v\n", err)
			return
		}

		// Handle the status code
		fmt.Printf("Push successful. HTTP Status Code: %d\n", statusCode)
	},
}

func init() {
	// Add gitPush as a subcommand of gitCmd
	gitCmd.AddCommand(gitPush)

	// Add gitCmd as a subcommand of rootCmd
	rootCmd.AddCommand(gitCmd)
}
