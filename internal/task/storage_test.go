package task

import (
	"context"
	"os"
	"testing"

	"github.com/Eddy-Nio/task-tracker-cli/config"
)

func setupTestStorage(t *testing.T) (*TaskStorage, string) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "tasks_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	ts, err := NewTaskStorage(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to create TaskStorage: %v", err)
	}

	return ts, tmpFile.Name()
}

func TestNewTaskStorage(t *testing.T) {
	scenarios := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Valid file path",
			filePath:    "tasks_test.json",
			expectError: false,
		},
		{
			name:        "Empty file path",
			filePath:    "",
			expectError: true,
		},
		{
			name:        "Invalid directory path",
			filePath:    "/invalid/directory/tasks.json",
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ts, err := NewTaskStorage(scenario.filePath)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for file path %s, got nil", scenario.filePath)
			}

			if !scenario.expectError {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if ts == nil {
					t.Error("Expected non-nil TaskStorage")
				}
			}
		})
	}
}

func TestTaskStorage_AddTask(t *testing.T) {
	ts, tmpFile := setupTestStorage(t)
	defer os.Remove(tmpFile)

	scenarios := []struct {
		name          string
		title         string
		description   string
		expectError   bool
		expectedTitle string
	}{
		{
			name:          "Valid task",
			title:         "Test Task",
			description:   "Test Description",
			expectError:   false,
			expectedTitle: "Test Task",
		},
		{
			name:          "Empty title",
			title:         "",
			description:   "Test Description",
			expectError:   false,
			expectedTitle: "Untitled Task",
		},
		{
			name:          "Very long title",
			title:         string(make([]byte, config.DefaultConfig.Task.MaxTitleLength+1)),
			description:   "Test Description",
			expectError:   true,
			expectedTitle: "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			task, err := ts.AddTask(scenario.title, scenario.description)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for scenario %s, got nil", scenario.name)
				return
			}

			if !scenario.expectError {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if task.Title != scenario.expectedTitle {
					t.Errorf("Expected title %q, got %q", scenario.expectedTitle, task.Title)
				}

				if task.ID == "" {
					t.Error("Task ID should not be empty")
				}

				savedTask, err := ts.GetTask(task.ID)
				if err != nil {
					t.Errorf("Failed to retrieve saved task: %v", err)
				}

				if savedTask.Title != task.Title {
					t.Errorf("Saved task title %q doesn't match original %q", savedTask.Title, task.Title)
				}
			}
		})
	}
}

func TestTaskStorage_GetTask(t *testing.T) {
	ts, tmpFile := setupTestStorage(t)
	defer os.Remove(tmpFile)

	// Add a task for testing
	addedTask, err := ts.AddTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	scenarios := []struct {
		name        string
		taskID      string
		expectError bool
	}{
		{
			name:        "Existing task",
			taskID:      addedTask.ID,
			expectError: false,
		},
		{
			name:        "Non-existent task",
			taskID:      "non-existent-id",
			expectError: true,
		},
		{
			name:        "Empty task ID",
			taskID:      "",
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			task, err := ts.GetTask(scenario.taskID)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for task ID %s, got nil", scenario.taskID)
			}

			if !scenario.expectError {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if task.ID != scenario.taskID {
					t.Errorf("Expected task ID %s, got %s", scenario.taskID, task.ID)
				}
			}
		})
	}
}

func TestTaskStorage_UpdateTask(t *testing.T) {
	ts, tmpFile := setupTestStorage(t)
	defer os.Remove(tmpFile)

	// Add a task for testing
	originalTask, err := ts.AddTask("Original Task", "Original Description")
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	scenarios := []struct {
		name        string
		taskID      string
		updates     map[string]interface{}
		expectError bool
	}{
		{
			name:   "Valid update - title only",
			taskID: originalTask.ID,
			updates: map[string]interface{}{
				"title": "Updated Title",
			},
			expectError: false,
		},
		{
			name:   "Valid update - status",
			taskID: originalTask.ID,
			updates: map[string]interface{}{
				"status": StatusInProgress,
			},
			expectError: false,
		},
		{
			name:   "Invalid task ID",
			taskID: "non-existent-id",
			updates: map[string]interface{}{
				"title": "Updated Title",
			},
			expectError: true,
		},
		{
			name:   "Invalid status update",
			taskID: originalTask.ID,
			updates: map[string]interface{}{
				"status": "INVALID_STATUS",
			},
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			updatedTask, err := ts.UpdateTask(context.Background(), scenario.taskID, scenario.updates)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for scenario %s, got nil", scenario.name)
			}

			if !scenario.expectError {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				for field, value := range scenario.updates {
					switch field {
					case "title":
						if updatedTask.Title != value {
							t.Errorf("Expected title %v, got %v", value, updatedTask.Title)
						}
					case "status":
						if updatedTask.Status != value {
							t.Errorf("Expected status %v, got %v", value, updatedTask.Status)
						}
					}
				}

				if !updatedTask.UpdatedAt.After(originalTask.UpdatedAt) {
					t.Error("UpdatedAt should be more recent than original task")
				}
			}
		})
	}
}

func TestTaskStorage_DeleteTask(t *testing.T) {
	ts, tmpFile := setupTestStorage(t)
	defer os.Remove(tmpFile)

	// Add a task for testing
	task, err := ts.AddTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	scenarios := []struct {
		name        string
		taskID      string
		expectError bool
	}{
		{
			name:        "Existing task",
			taskID:      task.ID,
			expectError: false,
		},
		{
			name:        "Non-existent task",
			taskID:      "non-existent-id",
			expectError: true,
		},
		{
			name:        "Empty task ID",
			taskID:      "",
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := ts.DeleteTask(context.Background(), scenario.taskID)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for task ID %s, got nil", scenario.taskID)
			}

			if !scenario.expectError {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				_, err := ts.GetTask(scenario.taskID)
				if err == nil {
					t.Error("Task should not exist after deletion")
				}
			}
		})
	}
}

func TestListTasksByStatus(t *testing.T) {
	ts, tmpFile := setupTestStorage(t)
	defer os.Remove(tmpFile)

	task1, err := ts.AddTask("Tarea 1", "Descripción 1")
	if err != nil {
		t.Fatalf("Error adding task1: %v", err)
	}

	task2, err := ts.AddTask("Tarea 2", "Descripción 2")
	if err != nil {
		t.Fatalf("Error adding task2: %v", err)
	}

	_, err = ts.UpdateTask(context.Background(), task2.ID, map[string]interface{}{
		"status": StatusInProgress,
	})
	if err != nil {
		t.Fatalf("Error updating task2 status: %v", err)
	}

	todoTasks := ts.ListTasksByStatus(StatusTodo)
	if len(todoTasks) != 1 || todoTasks[0].ID != task1.ID {
		t.Errorf("Expected 1 TODO task with ID %s, got %d tasks", task1.ID, len(todoTasks))
	}

	inProgressTasks := ts.ListTasksByStatus(StatusInProgress)
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 IN_PROGRESS task, got %d", len(inProgressTasks))
	}
}

func TestClearTasks(t *testing.T) {
	storage, _ := NewTaskStorage("test_tasks.json")
	defer os.Remove("test_tasks.json")

	storage.AddTask("Test Task", "Test Description")

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.Write([]byte("y\n"))
		w.Close()
	}()

	err := storage.ClearTasks()
	os.Stdin = oldStdin

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tasks := storage.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after clearing, got %d", len(tasks))
	}
}

func TestPrintTasks(t *testing.T) {
	storage, _ := NewTaskStorage("test_tasks.json")
	defer os.Remove("test_tasks.json")

	task, _ := NewTask("Test Task", "Test Description")
	storage.AddTask(task.Title, task.Description)

	storage.PrintTasks()

	storage.PrintTask(*task)
}
