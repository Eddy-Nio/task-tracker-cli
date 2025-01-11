package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Eddy-Nio/task-tracker-cli/config"
	"github.com/google/uuid"
)

// Status represents the status of a task
type Status string

const (
	// Statuses
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"

	// Error messages
	ErrTitleEmpty   = "title cannot be empty"
	ErrTitleTooLong = "title exceeds maximum length of %d characters"
	ErrDescTooLong  = "description exceeds maximum length of %d characters"
)

var (
	// Status aliases
	StatusAliases = map[string]Status{
		"todo":        StatusTodo,
		"t":           StatusTodo,
		"in_progress": StatusInProgress,
		"ip":          StatusInProgress,
		"p":           StatusInProgress,
		"done":        StatusDone,
		"d":           StatusDone,
	}

	// Common errors
	ErrInvalidTaskID = errors.New("invalid task ID")
	ErrStorageAccess = errors.New("storage access error")
)

// Task represents a task
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTask creates a new task
func NewTask(title, description string) (*Task, error) {
	if title == "" {
		title = "Untitled Task"
	}
	if len(title) > config.DefaultConfig.Task.MaxTitleLength {
		return nil, fmt.Errorf(ErrTitleTooLong, config.DefaultConfig.Task.MaxTitleLength)
	}

	if len(description) > config.DefaultConfig.Task.MaxDescriptionLength {
		return nil, fmt.Errorf(ErrDescTooLong, config.DefaultConfig.Task.MaxDescriptionLength)
	}

	uuid, err := generateTaskID()
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

// generateTaskID generates a task ID
func generateTaskID() (string, error) {
	uuidObject, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuidObject.String()[:8], nil
}

// ValidateStatus validates the status of a task
func ValidateStatus(s string) (Status, error) {
	normalizedStatus := strings.ToLower(s)
	if status, exists := StatusAliases[normalizedStatus]; exists {
		return status, nil
	}
	return "", fmt.Errorf("invalid status: %s. Use one of: todo/t, in_progress/ip/p, done/d", s)
}

// Validate validates the task
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(t.Title) > config.DefaultConfig.Task.MaxTitleLength {
		return fmt.Errorf("title length exceeds maximum of %d characters", config.DefaultConfig.Task.MaxTitleLength)
	}
	if len(t.Description) > config.DefaultConfig.Task.MaxDescriptionLength {
		return fmt.Errorf("description length exceeds maximum of %d characters", config.DefaultConfig.Task.MaxDescriptionLength)
	}
	if t.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}
	if _, err := ValidateStatus(string(t.Status)); err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}
	if t.CreatedAt.IsZero() {
		return fmt.Errorf("created_at cannot be zero")
	}
	if t.UpdatedAt.IsZero() {
		return fmt.Errorf("updated_at cannot be zero")
	}
	return nil
}
