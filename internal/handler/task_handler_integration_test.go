
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskmanager/internal/model"
	"taskmanager/internal/repository"
	"taskmanager/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupIntegrationHandler() *http.ServeMux {
	repo := repository.NewInMemoryTaskRepository(zap.NewNop())
	svc := service.NewTaskService(repo, zap.NewNop())
	h := NewTaskHandler(svc, zap.NewNop())
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestIntegration_CreateAndGetTask(t *testing.T) {
	mux := setupIntegrationHandler()
	task := &model.Task{Title: "Integration Test", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 Created")
	var created model.Task
	err := json.NewDecoder(w.Body).Decode(&created)
	require.NoError(t, err, "decoding created task")
	assert.NotEmpty(t, created.ID, "expected non-empty ID")
	assert.Equal(t, task.Title, created.Title, "expected title to match")

	// Now GET the created task
	r2 := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	require.Equal(t, http.StatusOK, w2.Code, "expected 200 OK")
	var fetched model.Task
	err = json.NewDecoder(w2.Body).Decode(&fetched)
	require.NoError(t, err, "decoding fetched task")
	assert.Equal(t, created.ID, fetched.ID, "expected ID to match")
}

func TestIntegration_ListTasks(t *testing.T) {
	mux := setupIntegrationHandler()
	// Create a task
	task := &model.Task{Title: "List Test", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 Created")
	// List tasks
	r2 := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	require.Equal(t, http.StatusOK, w2.Code, "expected 200 OK")
	var tasks []model.Task
	err := json.NewDecoder(w2.Body).Decode(&tasks)
	require.NoError(t, err, "decoding tasks list")
	assert.NotEmpty(t, tasks, "expected at least one task")
}

func TestIntegration_UpdateAndDeleteTask(t *testing.T) {
	mux := setupIntegrationHandler()
	// Create a task
	task := &model.Task{Title: "Update Test", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 Created")
	var created model.Task
	err := json.NewDecoder(w.Body).Decode(&created)
	require.NoError(t, err, "decoding created task")
	// Update
	update := &model.Task{Title: "Updated Title", Completed: true}
	updateBody, _ := json.Marshal(update)
	r2 := httptest.NewRequest(http.MethodPut, "/tasks/"+created.ID, bytes.NewReader(updateBody))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	require.Equal(t, http.StatusOK, w2.Code, "expected 200 OK")
	var updated model.Task
	err = json.NewDecoder(w2.Body).Decode(&updated)
	require.NoError(t, err, "decoding updated task")
	assert.Equal(t, "Updated Title", updated.Title, "expected updated title")
	// Delete
	r3 := httptest.NewRequest(http.MethodDelete, "/tasks/"+created.ID, nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)
	require.Equal(t, http.StatusNoContent, w3.Code, "expected 204 No Content")
	// Confirm deleted
	r4 := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, r4)
	require.Equal(t, http.StatusNotFound, w4.Code, "expected 404 after delete")
}

func TestIntegration_CreateTask_ValidationError(t *testing.T) {
	mux := setupIntegrationHandler()
	task := &model.Task{Title: "", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code, "expected 400 Bad Request")
	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err, "decoding error response")
	assert.Contains(t, resp, "error", "expected error in response")
}

func TestIntegration_CreateTask_WithUserSuppliedID(t *testing.T) {
	mux := setupIntegrationHandler()
	userID := "integration-id-456"
	task := &model.Task{ID: userID, Title: "Integration with user ID", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 Created")
	var created model.Task
	err := json.NewDecoder(w.Body).Decode(&created)
	require.NoError(t, err, "decoding created task")
	assert.Equal(t, userID, created.ID, "expected ID to match user supplied ID")
	assert.Equal(t, "Integration with user ID", created.Title, "expected title to match")
}
