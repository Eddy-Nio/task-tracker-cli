package cmd

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Error ejecutando rootCmd: %v", err)
	}
}

func TestAddCmd(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		title       string
		description string
	}{
		{
			name:        "Valid task",
			args:        []string{"add", "-t", "Test Task", "-d", "Test Description"},
			wantErr:     false,
			title:       "Test Task",
			description: "Test Description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)

			if tt.wantErr {
				if err := rootCmd.Execute(); err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err := rootCmd.Execute(); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
