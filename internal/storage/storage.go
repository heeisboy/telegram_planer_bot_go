package storage

import (
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrUserNotFound = errors.New("user not found")
)

// Task represents a single to-do item for a user.
type Task struct {
	ID         int
	UserID     int64
	Text       string
	Done       bool
	ReminderAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Store is the storage interface for tasks.
type Store interface {
	CreateTask(userID int64, text string) (Task, error)
	ListTasks(userID int64) ([]Task, error)
	GetTask(userID int64, id int) (Task, bool, error)
	UpdateTask(userID int64, task Task) error
	DeleteTask(userID int64, id int) error
}
