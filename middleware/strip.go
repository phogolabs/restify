package middleware

import (
	"net/http"

	"github.com/go-chi/chi"
)

// StripSlashes is a middleware that will match request paths with a trailing
// slash, strip it from the path and continue routing through the mux, if a route
// matches, then it will serve the handler.
func StripSlashes(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			path string
			rctx = chi.RouteContext(r.Context())
		)

		if rctx.RoutePath != "" {
			path = rctx.RoutePath
		} else {
			path = r.URL.Path
		}

		if len(path) > 1 && path[len(path)-1] == '/' {
			rctx.RoutePath = path[:len(path)-1]
			r.URL.Path = path
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
