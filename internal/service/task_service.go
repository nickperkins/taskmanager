package service

import (
	"context"
	"errors"
	"sort"
	"taskmanager/internal/idgen"
	"taskmanager/internal/model"
	"taskmanager/internal/repository"
	"time"

	"go.uber.org/zap"
)

// TaskService defines the business logic interface for tasks.

// taskServiceImpl provides business logic for managing tasks.
type taskServiceImpl struct {
	repo   repository.TaskRepository
	logger *zap.Logger
}

// NewTaskService creates a new TaskService with the given repository and logger.
func NewTaskService(repo repository.TaskRepository, logger *zap.Logger) TaskService {
	return &taskServiceImpl{repo: repo, logger: logger}
}

// CreateTask validates and creates a new task.
func (s *taskServiceImpl) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	// If ID is empty, generate a new one (now uses UUID)
	if task.ID == "" {
		task.ID = idgen.GenerateTaskID()
	}
	if err := task.Validate(); err != nil {
		s.logger.Warn("validation failed", zap.Error(err))
		return nil, err
	}
	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	if err := s.repo.CreateTask(ctx, task); err != nil {
		s.logger.Error("failed to create task", zap.Error(err))
		return nil, err
	}
	return task, nil
}

// GetTask retrieves a task by ID.
func (s *taskServiceImpl) GetTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		s.logger.Warn("task not found", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return task, nil
}

// ListTasks returns all tasks.
func (s *taskServiceImpl) ListTasks(ctx context.Context) ([]*model.Task, error) {
	tasks, err := s.repo.ListTasks(ctx)
	if err != nil {
		s.logger.Error("failed to list tasks", zap.Error(err))
		return nil, err
	}
	// Sort tasks by CreatedAt ascending
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})
	return tasks, nil
}

// UpdateTask updates an existing task by ID. Allows partial updates.
func (s *taskServiceImpl) UpdateTask(ctx context.Context, id string, update *model.Task) (*model.Task, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		s.logger.Warn("task not found for update", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	// If update.Title is explicitly set to empty string, that's a validation error
	if update.Title == "" {
		err := errors.New("title cannot be empty")
		s.logger.Warn("validation failed on update", zap.Error(err))
		return nil, err
	}

	// Only update fields that are set (partial update)
	if update.Title != "" {
		task.Title = update.Title
	}
	if update.Description != "" {
		task.Description = update.Description
	}
	// Only update Completed if explicitly set (cannot distinguish false from unset in Go, so always update)
	task.Completed = update.Completed
	task.UpdatedAt = time.Now().UTC()

	// Validate before calling repo.UpdateTask
	if err := task.Validate(); err != nil {
		s.logger.Warn("validation failed on update", zap.Error(err))
		return nil, err
	}

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		s.logger.Error("failed to update task", zap.Error(err))
		return nil, err
	}
	return task, nil
}

// DeleteTask deletes a task by ID.
func (s *taskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	if err := s.repo.DeleteTask(ctx, id); err != nil {
		s.logger.Warn("failed to delete task", zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

// generateTaskID is a stub for generating a unique string ID (to be improved in later tasks)

// ErrTaskNotFound is returned when a task is not found.
var ErrTaskNotFound = errors.New("task not found")
