package server_test

import (
	"net/http"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestHealthzHandler(t *testing.T) {
	assert.HTTPBodyContains(t, server.HealthzHandler, "GET", "/healthz", nil, "OK")
	assert.HTTPStatusCode(t, server.HealthzHandler, "GET", "/healthz", nil, http.StatusOK)
}
