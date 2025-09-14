package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syojctl",
	Short: "A command-line tool for interacting with the SYOJ platform",
	Long: `syojctl is a command-line tool that allows you to interact with the 
SYOJ (Sing Yin Online Judge) platform, including logging in and managing
your account credentials.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Create a logger
	logger := log.New(os.Stderr)
	
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}