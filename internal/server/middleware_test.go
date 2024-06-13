package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestServerMiddlewareHandler(t *testing.T) {
	validateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, w.Header().Get("Server"), "nextride-shortcut")
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	middleware := server.HeaderMiddleware(validateHandler)
	middleware.ServeHTTP(httptest.NewRecorder(), req)
}
