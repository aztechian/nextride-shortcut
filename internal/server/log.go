package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// LoggerMiddleware is a middleware that logs the request and response
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			logger := hlog.FromRequest(r)
			var event *zerolog.Event
			if r.URL.Path == "/healthz" { // log health checks as debug
				event = logger.Debug()
			} else {
				event = logger.Info()
			}
			event.
				Int("status", status).
				Int("size", size).
				Dur("elapsed_ms", duration).
				Msg("http access request")
		})(next).ServeHTTP(rw, req)
	})
}
