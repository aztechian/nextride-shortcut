package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {
	globalLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.ObjectsAreEqual(globalLogger, hlog.FromRequest(r))
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	middleware := server.LoggerMiddleware(validationHandler)
	middleware.ServeHTTP(httptest.NewRecorder(), req)
}
