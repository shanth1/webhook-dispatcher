package middleware

import (
	"net/http"

	"github.com/shanth1/gotools/log"
)

func WithLogger(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logger.With(
				log.Str("method", r.Method),
				log.Str("remote_addr", r.RemoteAddr),
			)

			ctx := log.NewContext(r.Context(), requestLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
