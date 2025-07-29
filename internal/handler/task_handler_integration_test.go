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
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)
	if created.ID == "" {
		t.Error("expected non-empty ID")
	}
	if created.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, created.Title)
	}

	// Now GET the created task
	r2 := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var fetched model.Task
	json.NewDecoder(w2.Body).Decode(&fetched)
	if fetched.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, fetched.ID)
	}
}

func TestIntegration_ListTasks(t *testing.T) {
	mux := setupIntegrationHandler()
	// Create a task
	task := &model.Task{Title: "List Test", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	// List tasks
	r2 := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var tasks []model.Task
	json.NewDecoder(w2.Body).Decode(&tasks)
	if len(tasks) == 0 {
		t.Error("expected at least one task")
	}
}

func TestIntegration_UpdateAndDeleteTask(t *testing.T) {
	mux := setupIntegrationHandler()
	// Create a task
	task := &model.Task{Title: "Update Test", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)
	// Update
	update := &model.Task{Title: "Updated Title", Completed: true}
	updateBody, _ := json.Marshal(update)
	r2 := httptest.NewRequest(http.MethodPut, "/tasks/"+created.ID, bytes.NewReader(updateBody))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var updated model.Task
	json.NewDecoder(w2.Body).Decode(&updated)
	if updated.Title != "Updated Title" {
		t.Errorf("expected updated title, got %q", updated.Title)
	}
	// Delete
	r3 := httptest.NewRequest(http.MethodDelete, "/tasks/"+created.ID, nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)
	if w3.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w3.Code)
	}
	// Confirm deleted
	r4 := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, r4)
	if w4.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w4.Code)
	}
}

func TestIntegration_CreateTask_ValidationError(t *testing.T) {
	mux := setupIntegrationHandler()
	task := &model.Task{Title: "", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["error"]; !ok {
		t.Error("expected error in response")
	}
}

func TestIntegration_CreateTask_WithUserSuppliedID(t *testing.T) {
	mux := setupIntegrationHandler()
	userID := "integration-id-456"
	task := &model.Task{ID: userID, Title: "Integration with user ID", Completed: false}
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created model.Task
	json.NewDecoder(w.Body).Decode(&created)
	if created.ID != userID {
		t.Errorf("expected ID %q, got %q", userID, created.ID)
	}
	if created.Title != "Integration with user ID" {
		t.Errorf("expected title %q, got %q", "Integration with user ID", created.Title)
	}
}
