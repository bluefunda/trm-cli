package cmd

import (
	"fmt"
	"github.com/bluefunda/trm-cli/config"
	"github.com/spf13/cobra"
)

var setEnvVar = &cobra.Command{
	Use:   "set",
	Short: "set environment variable ",
}

// urlAdd represents the add command under git
var setUrl = &cobra.Command{
	Use:   "url [url]",
	Short: "Set URL",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		inputUrl := args[0] // Get the URL from arguments

		// Check if the URL is stored in the config file
		storedUrl, err := config.ReadToken("url")
		if err != nil || storedUrl == "" {
			fmt.Println("No URL stored, adding new one.")
		} else {
			fmt.Printf("Stored URL: %s\n", storedUrl)

			// Prompt the user to confirm whether to overwrite the stored URL
			confirm := promptForConfirmation("Do you want to update the stored URL? (y/n): ")
			if confirm != "y" {
				fmt.Println("URL update canceled.")
				return
			}
		}

		// Store the URL in the config file
		envVars := map[string]string{
			"URL": inputUrl,
		}
		if err := config.UpdateEnvVars("url", envVars); err != nil {
			fmt.Printf("Failed to update URL in config file: %v\n", err)
			return
		}
		fmt.Println("URL successfully updated.")
	},
}

// Function to prompt for user confirmation
func promptForConfirmation(prompt string) string {
	var response string
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return response
}

func init() {

	setEnvVar.AddCommand(setUrl)
	// Add urlCmd as a subcommand of rootCmd
	rootCmd.AddCommand(setEnvVar)
}
