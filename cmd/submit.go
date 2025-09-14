package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stdx-space/syojctl/api"
	"github.com/stdx-space/syojctl/credentials"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit [problem-id]",
	Short: "Submit solution code for a problem",
	Long: `Submit solution code for a problem to SYOJ.
	
This command takes a problem ID and reads code from a file (specified with -i flag) 
or from standard input if no file is specified, then submits it to SYOJ for evaluation.

Examples:
  syojctl submit I001 -i solution.cpp
  cat solution.cpp | syojctl submit I001
  syojctl submit I001 < solution.cpp`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a logger
		logger := log.New(os.Stderr)

		// Get problem ID from arguments
		problemID := args[0]

		// Get language from flag
		language, _ := cmd.Flags().GetString("language")

		// Get input file from flag
		inputFile, _ := cmd.Flags().GetString("input")

		// Read source code
		var code string
		var err error

		if inputFile != "" {
			// Read from specified file
			codeBytes, err := os.ReadFile(inputFile)
			if err != nil {
				logger.Error("Failed to read input file", "error", err)
				return
			}
			code = string(codeBytes)

			// If language is not specified, try to infer it from file extension
			if language == "" {
				switch {
				case strings.HasSuffix(inputFile, ".cpp") || strings.HasSuffix(inputFile, ".cc") || strings.HasSuffix(inputFile, ".cxx"):
					language = "cpp20" // Use C++ 20
				case strings.HasSuffix(inputFile, ".c"):
					language = "c"
				case strings.HasSuffix(inputFile, ".py"):
					language = "python"
				case strings.HasSuffix(inputFile, ".js"):
					language = "javascript"
				case strings.HasSuffix(inputFile, ".go"):
					language = "go"
				case strings.HasSuffix(inputFile, ".java"):
					language = "java"
				case strings.HasSuffix(inputFile, ".rs"):
					language = "rust"
				case strings.HasSuffix(inputFile, ".cs"):
					language = "csharp"
				default:
					language = "cpp20" // Default to C++ 20
					logger.Warn("Could not infer language from file extension, defaulting to C++ 20")
				}
			}
		} else {
			// Read from standard input
			logger.Info("Reading code from standard input. Press Ctrl+D when finished.")
			reader := bufio.NewReader(os.Stdin)
			var builder strings.Builder
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				builder.WriteString(line)
			}
			code = builder.String()

			// If language is not specified and we're reading from stdin, default to C++ 20
			if language == "" {
				language = "cpp20"
				logger.Info("No language specified, defaulting to C++ 20")
			}
		}

		// Validate that we have code
		if code == "" {
			logger.Error("No code provided for submission")
			return
		}

		// Load credentials
		creds, err := credentials.Load()
		if err != nil {
			logger.Error("Failed to load credentials. Please login first.", "error", err)
			return
		}

		// Create API client with credentials
		client, err := api.NewClientWithCredentials(creds.Token, creds.TokenId)
		if err != nil {
			logger.Error("Failed to create API client", "error", err)
			return
		}

		// Submit code
		response, err := client.SubmitCode(code, language, problemID)
		if err != nil {
			logger.Error("Failed to submit code", "error", err)
			return
		}

		logger.Info("Code submitted successfully!", "message", response.Message)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Define flags for the submit command
	submitCmd.Flags().StringP("language", "l", "", "Programming language for submission")
	submitCmd.Flags().StringP("input", "i", "", "Input file containing source code (reads from stdin if not specified)")
}