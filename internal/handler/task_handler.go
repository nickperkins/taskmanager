package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"taskmanager/internal/model"
	"taskmanager/internal/service"
	"time"

	"go.uber.org/zap"
)

// TaskHandler handles HTTP requests for /tasks endpoints.
// TaskHandler handles HTTP requests for /tasks endpoints.
type TaskHandler struct {
	service service.TaskService
	logger  *zap.Logger
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(service service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{service: service, logger: logger}
}

// RegisterRoutes registers the /tasks routes to the given mux.
func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
	mux.HandleFunc("/tasks/", h.handleTaskByID)
}

// handleTasks handles POST (create) and GET (list) on /tasks.
func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.listTasks(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleTaskByID handles GET, PUT, DELETE on /tasks/{id}.
func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req model.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	created, err := h.service.CreateTask(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "task not found")
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req model.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	updated, err := h.service.UpdateTask(r.Context(), id, &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		h.writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeError writes a JSON error response.
func (h *TaskHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":      status,
			"message":   message,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
	h.logger.Warn("http error", zap.Int("status", status), zap.String("message", message))
}
