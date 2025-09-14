package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/stdx-space/syojctl/api"
	"github.com/stdx-space/syojctl/credentials"
)

// showProblemCmd represents the show-problem command
var showProblemCmd = &cobra.Command{
	Use:   "show-problem [problem-id]",
	Short: "Show details of a specific problem",
	Long:  `Show details of a specific problem from SYOJ by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a logger
		logger := log.New(os.Stderr)

		// Get the problem ID from arguments
		problemID := args[0]

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

		// Fetch problem details
		problem, err := client.GetProblem(problemID)
		if err != nil {
			logger.Error("Failed to fetch problem", "error", err)
			return
		}

		// Convert problem to markdown
		markdown := problemToMarkdown(problem)

		// Render with glamour
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
		)
		if err != nil {
			logger.Error("Failed to create renderer", "error", err)
			return
		}

		rendered, err := renderer.Render(markdown)
		if err != nil {
			logger.Error("Failed to render markdown", "error", err)
			return
		}

		fmt.Print(rendered)
	},
}

// problemToMarkdown converts a problem to markdown format
func problemToMarkdown(problem *api.Problem) string {
	var markdown strings.Builder

	// Title
	markdown.WriteString("# ")
	markdown.WriteString(problem.Title)
	markdown.WriteString("\n\n")

	// Problem metadata
	markdown.WriteString("- **Problem ID:** ")
	markdown.WriteString(problem.ID)
	markdown.WriteString("  \n")
	markdown.WriteString("- **Time limit:** ")
	markdown.WriteString(fmt.Sprintf("%d ms", problem.TimeLimit))
	markdown.WriteString("  \n")
	markdown.WriteString("- **Memory limit:** ")
	markdown.WriteString(fmt.Sprintf("%d MB", problem.MemoryLimit))
	markdown.WriteString("  \n")
	markdown.WriteString("- **Difficulty:** ")
	markdown.WriteString(problem.Difficulty)
	markdown.WriteString("\n\n")

	// Description
	if problem.Description != "" {
		markdown.WriteString("## Description\n\n")
		markdown.WriteString(problem.Description)
		markdown.WriteString("\n\n")
	}

	// Problem sections (Input, Output, Constraints, etc.)
	for _, section := range problem.ProblemSection {
		markdown.WriteString("## ")
		markdown.WriteString(section.Title)
		markdown.WriteString("\n\n")
		markdown.WriteString(section.Content)
		markdown.WriteString("\n\n")
	}

	// Sample test cases - show all non-empty test case groups
	testCaseGroupIndex := 0
	for i, group := range problem.TestCases {
		if len(group.TestCases) > 0 {
			markdown.WriteString(fmt.Sprintf("## Sample Testcase Group %d\n\n", i))
			
			for j, testCase := range group.TestCases {
				markdown.WriteString(fmt.Sprintf("### Testcase %d.%d\n\n", i, j))
				
				if testCase.Input != "" {
					markdown.WriteString("#### Input\n\n```\n")
					markdown.WriteString(testCase.Input)
					markdown.WriteString("\n```\n\n")
				}
				
				if testCase.Output != "" {
					markdown.WriteString("#### Output\n\n```\n")
					markdown.WriteString(testCase.Output)
					markdown.WriteString("\n```\n\n")
				}
			}
			testCaseGroupIndex++
		}
	}

	// Notes
	if problem.Notes != "" {
		markdown.WriteString("## Notes\n\n")
		markdown.WriteString(problem.Notes)
		markdown.WriteString("\n\n")
	}

	return markdown.String()
}

func init() {
	rootCmd.AddCommand(showProblemCmd)
}
