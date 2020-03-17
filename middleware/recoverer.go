package middleware

import (
	"net/http"

	"github.com/phogolabs/flaw"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if cause := recover(); cause != nil && cause != http.ErrAbortHandler {
				err := flaw.Errorf("%+v", cause)

				logger := GetLogger(r)
				logger.WithError(err).Error("critical panic error occurred")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
