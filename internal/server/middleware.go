package server

import "net/http"

type Middleware func(http.Handler) http.Handler

func WrapHandler(h http.Handler, m ...Middleware) http.Handler {

	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nextride-shortcut")
		next.ServeHTTP(w, r)
	})
}
