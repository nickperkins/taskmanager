package service

import (
	"context"
	"taskmanager/internal/model"
)

// TaskService defines the business logic interface for tasks.
type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	GetTask(ctx context.Context, id string) (*model.Task, error)
	ListTasks(ctx context.Context) ([]*model.Task, error)
	UpdateTask(ctx context.Context, id string, update *model.Task) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}
