package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// unitTestCmd represents the unit-test command
var unitTestCmd = &cobra.Command{
	Use:   "unit-test [package]",
	Short: "Run unit tests for the specified package",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		// Call the function to send the request with the provided package name
		response, err := unitTest(pkg)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response:", response)
		}
	},
}

func init() {
	// Add unitTestCmd as a subcommand of the root command
	rootCmd.AddCommand(unitTestCmd)
}
