package tasks

import (
	"time"

	"github.com/Eddy-Nio/task-tracker-cli/kit/utils"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "inProgress"
	StatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTask(title, description string) (*Task, error) {
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Task{
		ID:          uuid,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
