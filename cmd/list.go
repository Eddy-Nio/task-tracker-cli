/*
Copyright © 2025 Eddy Nio
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and filter tasks in your task list",
	Long: `The 'list' command displays tasks from your task list with powerful filtering options.

Each task entry shows:
  • ID: Unique identifier for the task
  • Title: Brief task name
  • Description: Detailed task information
  • Status: Current state (TODO, IN_PROGRESS, DONE)
  • Created At: Task creation timestamp
  • Updated At: Last modification timestamp

Available Flags:
  -s, --status string   Filter tasks by their current status:
                       • TODO - Show only pending tasks
                       • IN_PROGRESS - Show tasks being worked on
                       • DONE - Show completed tasks
                       If omitted, shows all tasks regardless of status

Examples:
  task list            # Lists all tasks
  task list -s TODO    # Lists only pending tasks
  task list -s DONE    # Lists only completed tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	var status string
	listCmd.Flags().StringVarP(&status, "status", "s", "", "Filter tasks by status (TODO, IN_PROGRESS, DONE)")
	listCmd.Flags().SortFlags = false

	listCmd.Run = func(cmd *cobra.Command, args []string) {
		storage, err := task.NewTaskStorage("tasks.json")
		if err != nil {
			log.Fatalf("Error initializing storage: %v", err)
		}

		var tasks []task.Task
		if status != "" {
			tasks = storage.ListTasksByStatus(task.Status(status))
		} else {
			tasks = storage.ListTasks()
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}

		storage.PrintTasks()
	}

}
