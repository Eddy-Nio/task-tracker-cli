package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/Eddy-Nio/task-tracker-cli/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	cleanupTestFile(t)

	tests := []struct {
		name       string
		status     string
		setupFile  string
		setupTasks []task.Task
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "List empty tasks",
			status:     "",
			setupTasks: []task.Task{},
			wantOutput: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFile == "corrupted_tasks.json" {
				err := os.WriteFile(tt.setupFile, []byte("{invalid json}"), 0644)
				assert.NoError(t, err)
			}
			buf := new(bytes.Buffer)

			listCmd.SetOut(buf)
			args := []string{}
			if tt.status != "" {
				args = append(args, "--status", tt.status)
			}
			listCmd.SetArgs(args)

			err := listCmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantOutput != "" {
					assert.Contains(t, buf.String(), tt.wantOutput)
				}
			}
		})
	}
}
