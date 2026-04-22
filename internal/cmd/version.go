package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the trm-cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("trm version %s\n", Version)
	},
}
