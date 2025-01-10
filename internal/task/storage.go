package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Eddy-Nio/task-tracker-cli/kit/utils"
)

type TaskStorage struct {
	mu       sync.RWMutex
	tasks    []Task
	filePath string
}

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

func (ts *TaskStorage) AddTask(title, description string) (Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	id, err := utils.GenerateUUID()
	if err != nil {
		return Task{}, err
	}

	now := time.Now().UTC()

	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ts.tasks = append(ts.tasks, task)

	fmt.Println("Task added to storage:", task)

	if err := ts.saveToFile(); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (ts *TaskStorage) loadFromFile() error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if _, err := os.Stat(ts.filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist, creating a new one...")
		return ts.saveToFile()
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

	fmt.Println("Loaded tasks from file:", ts.tasks)

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
