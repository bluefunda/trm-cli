package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// sapUserCmd represents the sap-user command
var sapUserCmd = &cobra.Command{
	Use:   "sap-user",
	Short: "Manage SAP users",
}

// readCmd represents the read command under sap-user
var readCmd = &cobra.Command{
	Use:   "read [input]",
	Short: "Read SAP users based on input",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		// Call the function to fetch all users
		response, err := fetchUsers(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	// Add readCmd as a subcommand of sapUserCmd
	sapUserCmd.AddCommand(readCmd)

	// Add sapUserCmd as a subcommand of rootCmd
	rootCmd.AddCommand(sapUserCmd)
}
