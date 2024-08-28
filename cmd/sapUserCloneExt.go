package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command under sap-user
var cloneCmd = &cobra.Command{
	Use:   "clone [usernameFrom] [username]",
	Short: "Create a new SAP user",
	Args:  cobra.ExactArgs(2), // Expect exactly 2 argument
	Run: func(cmd *cobra.Command, args []string) {
		usernameFrom := args[0]
		username := args[1]
		// Call the function to create the user
		response, err := cloneUser(usernameFrom, username)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	// Add createCmd as a subcommand of sapUserCmd
	sapUserCmd.AddCommand(cloneCmd)

	// Add sapUserCmd as a subcommand of rootCmd
	rootCmd.AddCommand(sapUserCmd)
}
