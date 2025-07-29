package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskmanager/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockTaskService is a testify mock for TaskService.
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*model.Task), args.Error(1)
}
func (m *MockTaskService) GetTask(ctx context.Context, id string) (*model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Task), args.Error(1)
}
func (m *MockTaskService) ListTasks(ctx context.Context) ([]*model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Task), args.Error(1)
}
func (m *MockTaskService) UpdateTask(ctx context.Context, id string, task *model.Task) (*model.Task, error) {
	args := m.Called(ctx, id, task)
	return args.Get(0).(*model.Task), args.Error(1)
}
func (m *MockTaskService) DeleteTask(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskHandler_CreateTask_Success(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	task := &model.Task{Title: "Test", Completed: false}
	ts.On("CreateTask", mock.Anything, task).Return(task, nil)
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp model.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, task.Title, resp.Title)
	ts.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_ValidationError(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	task := &model.Task{Title: "", Completed: false}
	ts.On("CreateTask", mock.Anything, task).Return((*model.Task)(nil), errors.New("title is required"))
	body, _ := json.Marshal(task)
	r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"].(map[string]interface{})["message"], "title is required")
	ts.AssertExpectations(t)
}

func TestTaskHandler_ListTasks(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	tasks := []*model.Task{{ID: "task-1", Title: "A", Completed: false}}
	ts.On("ListTasks", mock.Anything).Return(tasks, nil)
	r := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []model.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Len(t, resp, 1)
	ts.AssertExpectations(t)
}

func TestTaskHandler_GetTask_Success(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	task := &model.Task{ID: id, Title: "Test", Completed: false}
	ts.On("GetTask", mock.Anything, id).Return(task, nil)
	r := httptest.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, id, resp.ID)
	ts.AssertExpectations(t)
}

func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	ts.On("GetTask", mock.Anything, id).Return((*model.Task)(nil), errors.New("not found"))
	r := httptest.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"].(map[string]interface{})["message"], "not found")
	ts.AssertExpectations(t)
}

func TestTaskHandler_UpdateTask_Success(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	update := &model.Task{Title: "New", Completed: true}
	updated := &model.Task{ID: id, Title: "New", Completed: true}
	ts.On("UpdateTask", mock.Anything, id, update).Return(updated, nil)
	body, _ := json.Marshal(update)
	r := httptest.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.Task
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "New", resp.Title)
	ts.AssertExpectations(t)
}

func TestTaskHandler_UpdateTask_ValidationError(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	update := &model.Task{Title: ""}
	ts.On("UpdateTask", mock.Anything, id, update).Return((*model.Task)(nil), errors.New("title is required"))
	body, _ := json.Marshal(update)
	r := httptest.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"].(map[string]interface{})["message"], "title is required")
	ts.AssertExpectations(t)
}

func TestTaskHandler_DeleteTask_Success(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	ts.On("DeleteTask", mock.Anything, id).Return(nil)
	r := httptest.NewRequest(http.MethodDelete, "/tasks/"+id, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
	ts.AssertExpectations(t)
}

func TestTaskHandler_DeleteTask_NotFound(t *testing.T) {
	ts := new(MockTaskService)
	logger := zap.NewNop()
	h := NewTaskHandler(ts, logger)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	id := "task-1"
	ts.On("DeleteTask", mock.Anything, id).Return(errors.New("not found"))
	r := httptest.NewRequest(http.MethodDelete, "/tasks/"+id, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"].(map[string]interface{})["message"], "not found")
	ts.AssertExpectations(t)
}
