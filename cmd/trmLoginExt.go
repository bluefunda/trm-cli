package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	username string
	password string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the application",
	Run: func(cmd *cobra.Command, args []string) {
		username = promptForInput("Enter username: ")
		password = promptForInput("Enter password: ")

		if err := login(username, password); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Login successful!")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
