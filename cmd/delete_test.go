package cmd

import (
	"os"
	"testing"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCommand(t *testing.T) {
	cleanupTestFile(t)

	tests := []struct {
		name       string
		taskID     string
		setupFile  string
		setupTasks []task.Task
		wantErr    bool
	}{
		{
			name:    "Delete existing task",
			taskID:  "test-id-1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFile == "corrupted_tasks.json" {
				err := os.WriteFile(tt.setupFile, []byte("{invalid json}"), 0644)
				assert.NoError(t, err)
			}

			deleteCmd.SetArgs([]string{"--id", tt.taskID})
			err := deleteCmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
