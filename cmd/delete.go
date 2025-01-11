/*
Copyright © 2025 Eddy Nio
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task from the task list",
	Long: `The 'delete' command allows you to remove a task from your task list in the system.

You can specify the task ID you want to delete and it will be permanently removed from
the local JSON file. Make sure to double check the ID before deleting as this action
cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	var taskID string
	deleteCmd.Flags().StringVarP(&taskID, "id", "i", "", "Task ID to delete")
	deleteCmd.MarkFlagRequired("id")
	deleteCmd.Flags().SortFlags = false

	deleteCmd.Run = func(cmd *cobra.Command, args []string) {
		storage, err := task.NewTaskStorage("tasks.json")
		if err != nil {
			log.Fatalf("Error initializing storage: %v", err)
		}

		if err := storage.DeleteTask(context.Background(), taskID); err != nil {
			log.Fatalf("Error deleting task: %v", err)
		}

		fmt.Printf("Task with ID %s deleted successfully\n", taskID)
	}

}
