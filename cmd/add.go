/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	tasks "github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/spf13/cobra"
)

var title, description string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long: `The 'add' command allows you to add a new task to your task list in the system.

You can provide a detailed description of the task you want to add, which will be stored
in a local JSON file. The unique identifier (UUID) for the task will be automatically generated
and the task will be marked as 'todo' by default.`,

	Run: func(cmd *cobra.Command, args []string) {
		storage, err := tasks.NewTaskStorage("tasks.json")
		if err != nil {
			log.Fatalf("Error initializating storage file: %v", err)
		}

		task, err := storage.AddTask(title, description)
		if err != nil {
			log.Fatalf("Error when adding a new task: %v", err)
		}

		fmt.Printf("Task added: %v", task)
	},
}

func init() {
	addCmd.Flags().StringVarP(&title, "title", "t", "", "Task title")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Task description")
	addCmd.MarkFlagsRequiredTogether("title", "description")
	addCmd.Flags().SortFlags = false
	rootCmd.AddCommand(addCmd)
}
