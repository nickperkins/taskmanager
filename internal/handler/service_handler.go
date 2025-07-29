package handler

import (
	"encoding/json"
	"net/http"
)

// ServiceHandler handles the root endpoint.
type ServiceHandler struct{}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

func (h *ServiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.handleRoot)
}

func (h *ServiceHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"service": "taskmanager"})
}
