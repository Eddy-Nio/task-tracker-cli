package cmd

import (
	"os"
	"testing"
)

func TestAddCommand(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	addCmd.SetArgs([]string{
		"--title", "Test Task",
		"--description", "Test Description",
	})

	if err := addCmd.Execute(); err != nil {
		t.Errorf("Error executing add command: %v", err)
	}
}
