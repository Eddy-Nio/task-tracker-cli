package task

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	type testScenario struct {
		name  string
		input struct {
			title       string
			description string
		}
		expectedError bool
		expectedTitle string
	}

	scenarios := []testScenario{
		{
			name: "Valid task",
			input: struct {
				title       string
				description string
			}{
				title:       "Test task",
				description: "Test description",
			},
			expectedError: false,
			expectedTitle: "Test task",
		},
		{
			name: "Empty title",
			input: struct {
				title       string
				description string
			}{
				title:       "",
				description: "Test description",
			},
			expectedError: false,
			expectedTitle: "Untitled Task",
		},
		{
			name: "Very long title",
			input: struct {
				title       string
				description string
			}{
				title:       "This is a very long title that might exceed the maximum length allowed for a task title",
				description: "Test description",
			},
			expectedError: true,
			expectedTitle: "",
		},
		{
			name: "Empty description",
			input: struct {
				title       string
				description string
			}{
				title:       "Test task",
				description: "",
			},
			expectedError: false,
			expectedTitle: "Test task",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// When
			task, err := NewTask(scenario.input.title, scenario.input.description)

			// Then
			if (err != nil) != scenario.expectedError {
				t.Errorf("NewTask() error = %v, expectedError %v", err, scenario.expectedError)
				return
			}

			if !scenario.expectedError {
				if task == nil {
					t.Fatal("NewTask() returned nil task without error")
				}

				assertions := []struct {
					got     interface{}
					want    interface{}
					message string
				}{
					{task.Title, scenario.expectedTitle, "incorrect title"},
					{task.Description, scenario.input.description, "incorrect description"},
					{task.Status, StatusTodo, "incorrect status"},
				}

				for _, assertion := range assertions {
					if assertion.got != assertion.want {
						t.Errorf("%s: got %v, want %v", assertion.message, assertion.got, assertion.want)
					}
				}

				if task.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
			}
		})
	}
}

func TestValidateStatus(t *testing.T) {
	scenarios := []struct {
		name           string
		input          string
		expectedStatus Status
		expectError    bool
	}{
		{"Valid TODO status", "todo", StatusTodo, false},
		{"Valid TODO alias", "t", StatusTodo, false},
		{"Valid IN_PROGRESS status", "in_progress", StatusInProgress, false},
		{"Valid IN_PROGRESS alias", "ip", StatusInProgress, false},
		{"Valid DONE status", "done", StatusDone, false},
		{"Valid DONE alias", "d", StatusDone, false},
		{"Invalid status", "invalid", "", true},
		{"Empty status", "", "", true},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			status, err := ValidateStatus(scenario.input)

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for input %s, got nil", scenario.input)
			}

			if !scenario.expectError && err != nil {
				t.Errorf("Unexpected error for input %s: %v", scenario.input, err)
			}

			if status != scenario.expectedStatus {
				t.Errorf("Expected status %v, got %v", scenario.expectedStatus, status)
			}
		})
	}
}

func TestTask_Validate(t *testing.T) {
	now := time.Now()

	scenarios := []struct {
		name        string
		task        Task
		expectError bool
	}{
		{
			name: "Valid task",
			task: Task{
				ID:          "123",
				Title:       "Test Task",
				Description: "Test Description",
				Status:      StatusTodo,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expectError: false,
		},
		{
			name: "Empty title",
			task: Task{
				ID:          "123",
				Title:       "",
				Description: "Test Description",
				Status:      StatusTodo,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expectError: true,
		},
		{
			name: "Invalid status",
			task: Task{
				ID:          "123",
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "INVALID",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expectError: true,
		},
		{
			name: "Zero CreatedAt",
			task: Task{
				ID:          "123",
				Title:       "Test Task",
				Description: "Test Description",
				Status:      StatusTodo,
				CreatedAt:   time.Time{},
				UpdatedAt:   now,
			},
			expectError: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := scenario.task.Validate()

			if scenario.expectError && err == nil {
				t.Errorf("Expected error for scenario %s, got nil", scenario.name)
			}

			if !scenario.expectError && err != nil {
				t.Errorf("Unexpected error for scenario %s: %v", scenario.name, err)
			}
		})
	}
}

func TestGenerateTaskID(t *testing.T) {
	// Test multiple IDs to ensure uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := generateTaskID()

		if err != nil {
			t.Errorf("Unexpected error generating ID: %v", err)
		}

		if len(id) != 8 {
			t.Errorf("Expected ID length of 8, got %d", len(id))
		}

		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}

		ids[id] = true
	}
}
