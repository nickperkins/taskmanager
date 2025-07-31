
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	h := NewHealthHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	require.Equal(t, http.StatusOK, rw.Code, "expected 200 OK")
	assert.Equal(t, "{\"ok\":true}\n", rw.Body.String(), "unexpected body")
}
