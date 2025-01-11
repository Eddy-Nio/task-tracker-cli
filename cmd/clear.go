/*
Copyright © 2025 Eddy Nio
*/
package cmd

import (
	"fmt"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all tasks",
	Long:  `Clear all tasks from the storage file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storage, err := task.NewTaskStorage("tasks.json")
		if err != nil {
			return fmt.Errorf("error initializing storage: file might be corrupt: %v", err)
		}

		if err := storage.ClearTasks(); err != nil {
			return fmt.Errorf("error clearing tasks: %v", err)
		}

		fmt.Println("All tasks successfully deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
