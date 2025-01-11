package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var ErrNoUpdatesProvided = errors.New("no updates provided")

// Package task provides functionality for managing tasks in a task tracking system.
type TaskStorage struct {
	mu       sync.RWMutex // Protects concurrent access to tasks
	tasks    []Task       // Slice of tasks in memory
	filePath string       // Path to the JSON storage file
}

// NewTaskStorage creates a new TaskStorage instance with the specified file path.
// It loads existing tasks from the file if it exists, or creates a new file if it doesn't.
// Returns an error if the file operations fail.
func NewTaskStorage(filepath string) (*TaskStorage, error) {
	ts := &TaskStorage{
		tasks:    []Task{},
		filePath: filepath,
	}

	if err := ts.loadFromFile(); err != nil {
		return nil, err
	}

	return ts, nil
}

// AddTask creates a new task with the given title and description
func (ts *TaskStorage) AddTask(title, description string) (*Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, err := NewTask(title, description)
	if err != nil {
		return nil, err
	}

	ts.tasks = append(ts.tasks, *task)

	if err := ts.saveToFile(); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

func (ts *TaskStorage) ListTasks() []Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.tasks
}

func (ts *TaskStorage) GetTask(id string) (Task, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	idx, task, err := ts.findTaskById(id)
	if err != nil {
		return Task{}, err
	}

	fmt.Println("Task found:", task)

	return ts.tasks[idx], nil
}

func (ts *TaskStorage) findTaskById(id string) (int, Task, error) {
	for i, task := range ts.tasks {
		if task.ID == id {
			return i, task, nil
		}
	}

	return -1, Task{}, fmt.Errorf("task not found")
}

func (ts *TaskStorage) UpdateTask(ctx context.Context, taskID string, updates map[string]interface{}) (*Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		ts.mu.Lock()
		defer ts.mu.Unlock()

		// Find the task by ID
		idx, _, err := ts.findTaskById(taskID)
		if err != nil {
			return nil, err
		}

		// Apply updates
		for field, value := range updates {
			switch field {
			case "title":
				if title, ok := value.(string); ok {
					if title == "" {
						ts.tasks[idx].Title = "Untitled Task"
					} else {
						ts.tasks[idx].Title = title
					}
				}
			case "description":
				if desc, ok := value.(string); ok {
					ts.tasks[idx].Description = desc
				}
			case "status":
				if statusStr, ok := value.(string); ok {
					status, err := ValidateStatus(statusStr)
					if err != nil {
						return nil, err
					}
					ts.tasks[idx].Status = status
				} else if status, ok := value.(Status); ok {
					if _, err := ValidateStatus(string(status)); err != nil {
						return nil, err
					}
					ts.tasks[idx].Status = status
				}
			}
		}

		// Update timestamp and save
		ts.tasks[idx].UpdatedAt = time.Now()
		if err := ts.saveToFile(); err != nil {
			return nil, fmt.Errorf("failed to save updates: %w", err)
		}

		return &ts.tasks[idx], nil
	}
}

func (ts *TaskStorage) DeleteTask(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		ts.mu.Lock()
		defer ts.mu.Unlock()

		idx, _, err := ts.findTaskById(id)
		if err != nil {
			return err
		}

		ts.tasks = append(ts.tasks[:idx], ts.tasks[idx+1:]...)

		if err := ts.saveToFile(); err != nil {
			return err
		}

		return nil
	}
}

func (ts *TaskStorage) ListTasksByStatus(status Status) []Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var filteredTasks []Task
	for _, task := range ts.tasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}

func (ts *TaskStorage) ClearTasks() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	fmt.Print("Are you sure you want to delete all tasks? This action cannot be undone (y/n): ")
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return fmt.Errorf("operation cancelled by user")
	}

	ts.tasks = []Task{}
	if err := ts.saveToFile(); err != nil {
		return fmt.Errorf("error deleting tasks: %v", err)
	}

	fmt.Println("All tasks have been successfully deleted")
	return nil
}

func (ts *TaskStorage) loadFromFile() error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if _, err := os.Stat(ts.filePath); err != nil {
		if os.IsNotExist(err) {
			return ts.saveToFile()
		}
		return fmt.Errorf("checking file status: %w", err)
	}

	data, err := os.ReadFile(ts.filePath)
	if err != nil {
		return fmt.Errorf("error reading the file: %v", err)
	}

	if len(data) == 0 {
		fmt.Println("file is empty, creating a new one...")
		return ts.saveToFile()
	}

	if err := json.Unmarshal(data, &ts.tasks); err != nil {
		fmt.Println("error deserializing file, it may be corrupt.")

		if err := os.Remove(ts.filePath); err != nil {
			return fmt.Errorf("error removing corrupted file: %v", err)
		}

		fmt.Println("corrupted file removed, creating a new one...")

		return ts.saveToFile()
	}

	return nil
}

func (ts *TaskStorage) saveToFile() error {
	data, err := json.MarshalIndent(ts.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing tasks: %v", err)
	}

	err = os.WriteFile(ts.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing on the file: %v", err)
	}

	fmt.Println("Task saved successfully on the file:", ts.filePath)

	return nil
}

func (ts *TaskStorage) PrintTasks() {
	for _, task := range ts.tasks {
		fmt.Printf("------\nID: %s\nTitle: %s\nDescription: %s\nStatus: %s\nCreated: %s\nUpdated: %s\n------\n\n",
			task.ID, task.Title, task.Description, task.Status, task.CreatedAt.Format(time.RFC3339),
			task.UpdatedAt.Format(time.RFC3339))
	}
}

func (ts *TaskStorage) PrintTask(t Task) {
	fmt.Printf("------\nID: %s\nTitle: %s\nDescription: %s\nStatus: %s\nCreated: %s\nUpdated: %s\n------\n\n",
		t.ID, t.Title, t.Description, t.Status, t.CreatedAt.Format(time.RFC3339),
		t.UpdatedAt.Format(time.RFC3339))
}
