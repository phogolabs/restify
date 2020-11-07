package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/log"
	"go.opentelemetry.io/otel/api/trace"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(next http.Handler) http.Handler {
	fields := func(r *http.Request) log.Map {
		ctx := r.Context()

		scheme := func(r *http.Request) string {
			proto := "http"

			if r.TLS != nil {
				proto = "https"
			}

			return proto
		}

		fields := log.Map{
			"scheme":      scheme(r),
			"host":        r.Host,
			"url":         r.RequestURI,
			"proto":       r.Proto,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"request_id":  middleware.GetReqID(r.Context()),
		}

		span := trace.
			SpanFromContext(ctx).
			SpanContext()

		if span.HasTraceID() {
			fields["trace_id"] = span.TraceID
		}

		if span.HasSpanID() {
			fields["span_id"] = span.SpanID
		}

		return fields
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(fields(r))
		ctx := log.SetContext(r.Context(), logger)

		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(writer, r.WithContext(ctx))

		logger = logger.WithFields(log.Map{
			"status":   writer.Status(),
			"size":     writer.BytesWritten(),
			"duration": time.Since(start).String(),
		})

		switch {
		case writer.Status() >= 500:
			logger.Error("request handling fail")
		case writer.Status() >= 400:
			logger.Warn("request handling warning")
		default:
			logger.Info("request handling success")
		}
	}

	return http.HandlerFunc(fn)
}

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Logger {
	return log.GetContext(r.Context())
}
