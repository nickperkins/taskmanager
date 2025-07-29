package service

import (
	"context"
	"errors"
	"testing"
	"time"

	//   "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"taskmanager/internal/model"
	// ...existing code...
)

// MockTaskRepository is a testify mock for TaskRepository.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}
func (m *MockTaskRepository) GetTask(ctx context.Context, id string) (*model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Task), args.Error(1)
}
func (m *MockTaskRepository) ListTasks(ctx context.Context) ([]*model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Task), args.Error(1)
}
func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}
func (m *MockTaskRepository) DeleteTask(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskService_CreateTask_Success(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	task := &model.Task{
		Title:       "Test Task",
		Description: "A test task",
		Completed:   false,
	}
	repo.On("CreateTask", ctx, mock.AnythingOfType("*model.Task")).Return(nil)
	created, err := ts.CreateTask(ctx, task)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, task.Title, created.Title)
	assert.Equal(t, task.Description, created.Description)
	assert.False(t, created.Completed)
	assert.WithinDuration(t, time.Now(), created.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), created.UpdatedAt, time.Second)
	repo.AssertExpectations(t)
}

func TestTaskService_CreateTask_ValidationError(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	task := &model.Task{Title: "", Completed: false}
	created, err := ts.CreateTask(ctx, task)
	assert.Error(t, err)
	assert.Nil(t, created)
}

func TestTaskService_GetTask_Success(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	task := &model.Task{ID: id, Title: "Test", Completed: false}
	repo.On("GetTask", ctx, id).Return(task, nil)
	got, err := ts.GetTask(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	repo.On("GetTask", ctx, id).Return((*model.Task)(nil), errors.New("task not found"))
	got, err := ts.GetTask(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestTaskService_ListTasks(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	tasks := []*model.Task{{ID: "task-1", Title: "A", Completed: false}}
	repo.On("ListTasks", ctx).Return(tasks, nil)
	got, err := ts.ListTasks(ctx)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_Success(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	existing := &model.Task{ID: id, Title: "Old", Completed: false}
	update := &model.Task{Title: "New", Completed: true}
	repo.On("GetTask", ctx, id).Return(existing, nil)
	repo.On("UpdateTask", ctx, mock.AnythingOfType("*model.Task")).Return(nil)
	got, err := ts.UpdateTask(ctx, id, update)
	assert.NoError(t, err)
	assert.Equal(t, "New", got.Title)
	assert.True(t, got.Completed)
	repo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_ValidationError(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	existing := &model.Task{ID: id, Title: "Old", Completed: false}
	update := &model.Task{Title: ""}
	repo.On("GetTask", ctx, id).Return(existing, nil)
	got, err := ts.UpdateTask(ctx, id, update)
	assert.Error(t, err)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_NotFound(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	update := &model.Task{Title: "New"}
	repo.On("GetTask", ctx, id).Return((*model.Task)(nil), errors.New("task not found"))
	got, err := ts.UpdateTask(ctx, id, update)
	assert.Error(t, err)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_Success(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	repo.On("DeleteTask", ctx, id).Return(nil)
	err := ts.DeleteTask(ctx, id)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_NotFound(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	id := "task-1"
	repo.On("DeleteTask", ctx, id).Return(errors.New("task not found"))
	err := ts.DeleteTask(ctx, id)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestTaskService_CreateTask_WithUserSuppliedID(t *testing.T) {
	repo := new(MockTaskRepository)
	logger := zap.NewNop()
	ts := NewTaskService(repo, logger)
	ctx := context.Background()
	userID := "custom-id-123"
	task := &model.Task{
		ID:          userID,
		Title:       "Task with user ID",
		Description: "Should use provided ID",
		Completed:   false,
	}
	repo.On("CreateTask", ctx, mock.AnythingOfType("*model.Task")).Return(nil)
	created, err := ts.CreateTask(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, userID, created.ID)
	assert.Equal(t, task.Title, created.Title)
	repo.AssertExpectations(t)
}
