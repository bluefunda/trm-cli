package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// codeInspectCmd represents the code-inspect command
var codeInspectCmd = &cobra.Command{
	Use:   "qa-check [package|object] [value]",
	Short: "Run QA checks for the specified package or object",
	Args:  cobra.ExactArgs(2), // Expect exactly two arguments
	Run: func(cmd *cobra.Command, args []string) {
		inspectType := args[0]
		value := args[1]

		var response string
		var err error

		if inspectType == "package" {
			response, err = codeInspect(value, nil) // Inspecting by package
		} else if inspectType == "object" {
			response, err = codeInspect("", []string{value}) // Inspecting by object
		} else {
			fmt.Println("Error: First argument must be 'package' or 'object'")
			return
		}

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	// Add codeInspectCmd as a subcommand of the root command
	rootCmd.AddCommand(codeInspectCmd)
}
