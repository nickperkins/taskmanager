// Package repository provides interfaces for task storage.
package repository

import (
	"context"
	"taskmanager/internal/model"
)

// TaskReader defines read operations for tasks.
type TaskReader interface {
	// GetTask retrieves a task by its ID.
	GetTask(ctx context.Context, id string) (*model.Task, error)
	// ListTasks returns all tasks.
	ListTasks(ctx context.Context) ([]*model.Task, error)
}

// TaskWriter defines write operations for tasks.
type TaskWriter interface {
	// CreateTask adds a new task.
	CreateTask(ctx context.Context, task *model.Task) error
	// UpdateTask updates an existing task.
	UpdateTask(ctx context.Context, task *model.Task) error
	// DeleteTask removes a task by its ID.
	DeleteTask(ctx context.Context, id string) error
}

// TaskRepository combines read and write operations for tasks.
type TaskRepository interface {
	TaskReader
	TaskWriter
}
