package middleware

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
)

// Tracer is a middleware that uses openconsensus to trace a http requests
func Tracer(next http.Handler) http.Handler {
	tracer := ochttp.Handler{
		Handler: next,
	}

	return http.HandlerFunc(tracer.ServeHTTP)
}
