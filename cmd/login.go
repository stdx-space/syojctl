package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stdx-space/syojctl/api"
	"github.com/stdx-space/syojctl/credentials"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the SYOJ platform",
	Long:  `Login to the SYOJ platform using your email and password to obtain authentication credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a logger
		logger := log.New(os.Stderr)

		// Get username and password from flags or environment variables
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		// If not provided via flags, check environment variables
		if username == "" {
			username = os.Getenv("SYOJ_USERNAME")
		}
		if password == "" {
			password = os.Getenv("SYOJ_PASSWORD")
		}

		// Check if credentials are provided
		if username == "" || password == "" {
			logger.Error("Username and password are required. Please provide them via flags or environment variables.")
			return
		}

		// Create API client
		client, err := api.NewClient()
		if err != nil {
			logger.Error("Failed to create API client", "error", err)
			return
		}

		// Perform login
		apiCreds, err := client.Login(username, password)
		if err != nil {
			logger.Error("Failed during login", "error", err)
			return
		}

		// Convert API credentials to our credentials struct
		creds := &credentials.Credentials{
			Token:   apiCreds.Token,
			TokenId: apiCreds.TokenId,
		}

		// Save credentials to XDG config directory
		err = creds.Save()
		if err != nil {
			logger.Error("Failed to save credentials to file", "error", err)
			return
		}

		logger.Info("Credentials saved to XDG config directory")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Define flags for username and password
	loginCmd.Flags().StringP("username", "u", "", "Username for SYOJ login (defaults to SYOJ_USERNAME environment variable)")
	loginCmd.Flags().StringP("password", "p", "", "Password for SYOJ login (defaults to SYOJ_PASSWORD environment variable)")
}