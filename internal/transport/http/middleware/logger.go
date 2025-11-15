package middleware

import (
	"net/http"

	"github.com/shanth1/gotools/log"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func WithLogger(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// start := time.Now()

			reqLogger := logger.With(
				log.Str("method", r.Method),
				log.Str("path", r.URL.Path),
				log.Str("remote_addr", r.RemoteAddr),
				log.Str("user_agent", r.UserAgent()),
			)

			ctx := log.NewContext(r.Context(), reqLogger)
			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r.WithContext(ctx))

			reqLogger.Info().
				Int("status_code", rw.statusCode).
				// Dur("duration", time.Since(start)). // TODO!
				Msg("request completed")
		})
	}
}
