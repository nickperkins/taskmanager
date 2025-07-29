// Package model contains the core data structures for the Task Management API.
package model

import (
	"errors"
	"strings"
	"time"
)

// Task represents a task in the Task Management API.
//
// Fields:
//   - ID: unique identifier (string, max 24 chars, alphanumeric/dash)
//   - Title: required, 1-200 characters
//   - Description: optional, max 1000 characters
//   - Completed: boolean, required (default false)
//   - CreatedAt: timestamp when task was created
//   - UpdatedAt: timestamp when task was last updated
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks the Task fields for correctness according to business rules.
func (t *Task) Validate() error {
	// Validate ID: non-empty, alphanumeric/dash, max 36 chars (UUID allowed)
	id := strings.TrimSpace(t.ID)
	if id == "" {
		return errors.New("id is required")
	}
	if len(id) > 36 {
		return errors.New("id must be at most 36 characters")
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
			return errors.New("id must be alphanumeric or dash")
		}
	}

	// Validate Title: required, 1-200 chars
	title := strings.TrimSpace(t.Title)
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) < 1 || len(title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}

	// Validate Description: optional, max 1000 chars
	if len(t.Description) > 1000 {
		return errors.New("description must be at most 1000 characters")
	}

	// Completed: required (bool, default false)
	// No validation needed for bool, but check for presence if needed in JSON unmarshalling elsewhere

	return nil
}
