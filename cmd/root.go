/*
Copyright © 2025 Eddy Nio
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task-tracker-cli",
	Short: "A CLI tool to manage and track your tasks",
	Long: `Task Tracker CLI is a command-line tool designed to help you manage and track your tasks efficiently.

With this tool you can:
- Add new tasks with title and description
- List all your tasks and filter them by status
- Update task details and status
- Delete tasks when completed
- Clear all tasks when needed

All tasks are stored locally in a JSON file for easy access and persistence.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
