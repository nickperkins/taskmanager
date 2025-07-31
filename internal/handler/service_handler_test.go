
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceHandler(t *testing.T) {
	h := NewServiceHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	require.Equal(t, http.StatusOK, rw.Code, "expected 200 OK")
	assert.Equal(t, "{\"service\":\"taskmanager\"}\n", rw.Body.String(), "unexpected body")
}
