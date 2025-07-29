package idgen

import "github.com/google/uuid"

// GenerateTaskID returns a new UUID string for task IDs.
func GenerateTaskID() string {
	return uuid.NewString()
}
