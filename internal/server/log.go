package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter // compose http.ResponseWriter
	size                int
	status              int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.size += size                         // capture size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(rw)

		defer func() {
			panicVal := recover()
			if panicVal != nil {
				lrw.status = http.StatusInternalServerError // ensure that the status code is updated
				panic(panicVal)                             // continue panicking
			}

			log.
				Info().
				Str("method", req.Method).
				Str("url", req.URL.RequestURI()).
				Str("user_agent", req.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Int("status_code", lrw.status).
				Int("size", lrw.size).
				Str("host", req.Host).
				Msg("access request")
		}()

		next.ServeHTTP(lrw, req) // inject our logging implementation of http.ResponseWriter
	}
	return http.HandlerFunc(loggingFn)
}
