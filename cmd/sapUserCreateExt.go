package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command under sap-user
var createCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Create a new SAP user",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		// Call the function to create the user
		response, err := createUser(username)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	// Add createCmd as a subcommand of sapUserCmd
	sapUserCmd.AddCommand(createCmd)

	// Add sapUserCmd as a subcommand of rootCmd
	rootCmd.AddCommand(sapUserCmd)
}
