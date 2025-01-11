/*
Copyright © 2025 Eddy Nio
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing task",
	Long: `The 'update' command allows you to modify an existing task in your task list.

You can update various attributes of a task including its title, description, and status.
You'll need to provide the task ID and specify which fields you want to update.
The changes will be saved to the local JSON file and the task's 'updated_at' timestamp
will be automatically updated.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	var taskID, title, description, status string

	updateCmd.Flags().StringVarP(&taskID, "id", "i", "", "Task ID to update")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "New task title")
	updateCmd.Flags().StringVarP(&description, "desc", "d", "", "New task description")
	updateCmd.Flags().StringVarP(&status, "status", "s", "", "New task status (todo/t, in_progress/ip/p, done/d)")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().SortFlags = false

	updateCmd.RunE = func(cmd *cobra.Command, args []string) error {
		storage, err := task.NewTaskStorage("tasks.json")
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %v", err)
		}

		updates := make(map[string]interface{})
		if title != "" {
			updates["title"] = title
		}
		if description != "" {
			updates["description"] = description
		}
		if status != "" {
			if _, err := task.ValidateStatus(status); err != nil {
				return fmt.Errorf("invalid status: %v", err)
			}
			updates["status"] = status
		}

		if len(updates) == 0 {
			return fmt.Errorf("at least one field must be provided for update")
		}

		updatedTask, err := storage.UpdateTask(context.Background(), taskID, updates)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}

		fmt.Println("Task updated successfully:")
		storage.PrintTask(*updatedTask)
		return nil
	}
}
