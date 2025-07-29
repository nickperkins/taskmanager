// Package repository provides an in-memory implementation of TaskRepository.
package repository

import (
	"context"
	"errors"
	"sync"
	"taskmanager/internal/model"

	"go.uber.org/zap"
)

// InMemoryTaskRepository is a thread-safe in-memory implementation of TaskRepository.
type InMemoryTaskRepository struct {
	mu     sync.RWMutex
	tasks  map[string]*model.Task
	logger *zap.Logger
}

// NewInMemoryTaskRepository creates a new InMemoryTaskRepository.
func NewInMemoryTaskRepository(logger *zap.Logger) *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  make(map[string]*model.Task),
		logger: logger,
	}
}

// CreateTask adds a new task to the repository.
func (r *InMemoryTaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		r.logger.Warn("task already exists", zap.String("id", task.ID))
		return errors.New("task already exists")
	}
	r.tasks[task.ID] = task
	r.logger.Info("task created", zap.String("id", task.ID))
	return nil
}

// GetTask retrieves a task by its ID.
func (r *InMemoryTaskRepository) GetTask(ctx context.Context, id string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		r.logger.Warn("task not found", zap.String("id", id))
		return nil, errors.New("task not found")
	}
	r.logger.Debug("task retrieved", zap.String("id", id))
	return task, nil
}

// ListTasks returns all tasks in the repository.
func (r *InMemoryTaskRepository) ListTasks(ctx context.Context) ([]*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	r.logger.Debug("listed tasks", zap.Int("count", len(tasks)))
	return tasks, nil
}

// UpdateTask updates an existing task in the repository.
func (r *InMemoryTaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; !exists {
		r.logger.Warn("task not found for update", zap.String("id", task.ID))
		return errors.New("task not found")
	}
	r.tasks[task.ID] = task
	r.logger.Info("task updated", zap.String("id", task.ID))
	return nil
}

// DeleteTask removes a task by its ID from the repository.
func (r *InMemoryTaskRepository) DeleteTask(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[id]; !exists {
		r.logger.Warn("task not found for delete", zap.String("id", id))
		return errors.New("task not found")
	}
	delete(r.tasks, id)
	r.logger.Info("task deleted", zap.String("id", id))
	return nil
}
