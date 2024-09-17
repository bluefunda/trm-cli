package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// createTransportCmd represents the create transport command
var createTransportCmd = &cobra.Command{
	Use:   "create-transport",
	Short: "Create a new SAP transport",
	Args:  cobra.NoArgs, // No arguments expected, we will prompt the user for input
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Prompt for request type
		fmt.Println("Enter 'K' for Workbench or 'W' for Customizing transport:")
		requestType, _ := reader.ReadString('\n')
		requestType = strings.TrimSpace(requestType)

		// Prompt for author
		fmt.Println("Enter author name:")
		author, _ := reader.ReadString('\n')
		author = strings.TrimSpace(author)

		// Prompt for text
		fmt.Println("Enter transport description text:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		// Call the function to create the transport
		response, err := createTransport(requestType, author, text)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

func init() {
	// Add sapUserCmd as a subcommand of rootCmd
	rootCmd.AddCommand(createTransportCmd)
}
