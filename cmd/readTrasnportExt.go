package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// readTransportCmd represents the read command under transport
var readTransportCmd = &cobra.Command{
	Use:   "read-transport",
	Short: "Read transport data with user-provided input",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt for parameters if they are not provided via flags or are empty
		requestID, _ := cmd.Flags().GetString("requestid")
		username, _ := cmd.Flags().GetString("username")
		transportStatus, _ := cmd.Flags().GetString("transportStatus")
		transportFunction, _ := cmd.Flags().GetString("transportFunction")

		// Ask for inputs if values are not passed via flags
		requestID = getInputOrDefault("Enter request ID (or leave blank to skip): ", requestID)
		username = getInputOrDefault("Enter username (or leave blank to skip): ", username)
		transportStatus = getInputOrDefault("Enter transport status (or leave blank to skip): ", transportStatus)
		transportFunction = getInputOrDefault("Enter transport function (or leave blank to skip): ", transportFunction)

		// Call the function to fetch transport data
		response, err := fetchTransportData(requestID, username, transportStatus, transportFunction)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(response)
		}
	},
}

// getInputOrDefault prompts the user for input if not provided and returns either the input or the default value
func getInputOrDefault(prompt, defaultValue string) string {
	if defaultValue != "" {
		return defaultValue
	}

	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func init() {
	// Add flags for optional parameters
	readTransportCmd.Flags().String("requestid", "", "Request ID")
	readTransportCmd.Flags().String("username", "", "Username")
	readTransportCmd.Flags().String("transportStatus", "", "Transport status")
	readTransportCmd.Flags().String("transportFunction", "", "Transport function")

	// Add transportCmd as a subcommand of rootCmd
	rootCmd.AddCommand(readTransportCmd)
}
