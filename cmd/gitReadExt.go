package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// gitRepoCmd represents the git-repo command
var gitRepoCmd = &cobra.Command{
	Use:   "git",
	Short: "Manage Git repositories",
}

// gitReadCmd represents the read command under git-repo
var gitReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read all Git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the function to fetch repositories
		response, err := fetchRepo()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Repositories:", response)
		}
	},
}

func init() {
	// Add gitReadCmd as a subcommand of gitRepoCmd
	gitRepoCmd.AddCommand(gitReadCmd)

	// Add gitRepoCmd as a subcommand of rootCmd
	rootCmd.AddCommand(gitRepoCmd)
}
