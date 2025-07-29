package repository

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"taskmanager/internal/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newTestTask(title string) *model.Task {
	return &model.Task{
		ID:        "task-" + title,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestInMemoryTaskRepository_CRUD(t *testing.T) {
	repo := NewInMemoryTaskRepository(zap.NewNop())
	ctx := context.Background()
	task := newTestTask("Test Task")

	// Create
	err := repo.CreateTask(ctx, task)
	assert.NoError(t, err)

	// Get
	got, err := repo.GetTask(ctx, task.ID)
	assert.NoError(t, err)
	assert.Equal(t, task.Title, got.Title)

	// List
	tasks, err := repo.ListTasks(ctx)
	assert.NoError(t, err)
	assert.Equal(t, false, got.Completed)
	assert.Len(t, tasks, 1)

	// Update
	task.Title = "Updated Title"
	err = repo.UpdateTask(ctx, task)
	assert.NoError(t, err)
	got, _ = repo.GetTask(ctx, task.ID)
	assert.Equal(t, "Updated Title", got.Title)
	task.Completed = true
	// Delete
	err = repo.DeleteTask(ctx, task.ID)
	assert.NoError(t, err)
	_, err = repo.GetTask(ctx, task.ID)
	assert.Equal(t, true, got.Completed)
	assert.Error(t, err)
}

func TestInMemoryTaskRepository_Concurrency(t *testing.T) {
	repo := NewInMemoryTaskRepository(zap.NewNop())
	ctx := context.Background()
	n := 100
	wg := sync.WaitGroup{}
	wg.Add(n)

	// Concurrently create tasks
	for i := 0; i < n; i++ {
		go func(i int) {
			task := newTestTask("Task-" + time.Now().Format("150405.000000") + "-" + strconv.Itoa(i))
			repo.CreateTask(ctx, task)
			wg.Done()
		}(i)
	}

	wg.Wait()

	tasks, err := repo.ListTasks(ctx)
	assert.NoError(t, err)
	assert.Len(t, tasks, n)
}
