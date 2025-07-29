package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceHandler(t *testing.T) {
	h := NewServiceHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if got := rw.Body.String(); got != "{\"service\":\"taskmanager\"}\n" {
		t.Errorf("unexpected body: %s", got)
	}
}
