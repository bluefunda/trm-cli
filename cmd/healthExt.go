package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health of application",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the health function and handle its output
		response, err := health()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd) // Correct command name
}
