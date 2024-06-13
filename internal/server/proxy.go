package server

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func ProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ProxyHeaders(next).ServeHTTP(w, r)
	})
}
