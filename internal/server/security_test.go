package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHandler(t *testing.T) {
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, w.Header().Get("X-Frame-Options"), "DENY")
		assert.Contains(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
		assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	middleware := server.SecurityMiddleware(validationHandler)
	middleware.ServeHTTP(httptest.NewRecorder(), req)
}
