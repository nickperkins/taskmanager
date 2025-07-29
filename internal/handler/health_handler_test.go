package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	h := NewHealthHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if got := rw.Body.String(); got != "{\"ok\":true}\n" {
		t.Errorf("unexpected body: %s", got)
	}
}
