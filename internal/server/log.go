package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
)

// LoggerMiddleware is a middleware that logs the request and response
func LoggerMiddleware(next http.Handler) http.Handler {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {

		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Int("status", status).
				Int("size", size).
				Dur("elapsed_ms", duration).
				Msg("http access request")
		})(next).ServeHTTP(rw, req)
	}
	return http.HandlerFunc(loggingFn)
}
