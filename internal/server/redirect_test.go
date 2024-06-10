package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestRedirectHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	handler := http.HandlerFunc(server.RedirectHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Contains(t, rr.Header().Get("Location"), "/next")
	assert.HTTPStatusCode(t, server.RedirectHandler, "GET", "/", nil, http.StatusSeeOther)
}
