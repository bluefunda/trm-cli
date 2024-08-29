package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// gitReadCmd represents the read command under git-repo
var repoSet = &cobra.Command{
	Use:   "set [repo-name]",
	Short: "Set Git repositories",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		repoName := args[0] // Get the repository name from arguments
		// Call the function to fetch the repository key
		key, err := RepoKey(repoName)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Repository Key:", key)
		}
	},
}

func init() {
	// Add gitReadCmd as a subcommand of gitRepoCmd
	gitRepoCmd.AddCommand(repoSet)

	// Add gitRepoCmd as a subcommand of rootCmd
	rootCmd.AddCommand(gitRepoCmd)
}
