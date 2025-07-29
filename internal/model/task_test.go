package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidation_Valid(t *testing.T) {
	task := &Task{
		ID:          "task-123",
		Title:       "Test Task",
		Description: "A valid description.",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	assert.NoError(t, task.Validate(), "valid task should not return error")
}

func TestTaskValidation_TitleRequired(t *testing.T) {
	task := &Task{
		ID:        "task-123",
		Title:     "",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "title is required")
}

func TestTaskValidation_TitleLength(t *testing.T) {
	longTitle := make([]byte, 201)
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	task := &Task{
		ID:        "task-123",
		Title:     string(longTitle),
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "title must be between 1 and 200 characters")
}

func TestTaskValidation_DescriptionLength(t *testing.T) {
	longDesc := make([]byte, 1001)
	for i := range longDesc {
		longDesc[i] = 'a'
	}
	task := &Task{
		ID:          "task-123",
		Title:       "Valid Title",
		Description: string(longDesc),
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "description must be at most 1000 characters")
}

// ID validation tests
func TestTaskValidation_IDRequired(t *testing.T) {
	task := &Task{
		ID:        "",
		Title:     "Valid Title",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "id is required")
}

func TestTaskValidation_IDTooLong(t *testing.T) {
	task := &Task{
		ID:        "abcdefghijklmnopqrstuvwxyz12345678901234567890",
		Title:     "Valid Title",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "id must be at most 36 characters")
}

func TestTaskValidation_IDInvalidChars(t *testing.T) {
	task := &Task{
		ID:        "invalid$id!",
		Title:     "Valid Title",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.ErrorContains(t, task.Validate(), "id must be alphanumeric or dash")
}
