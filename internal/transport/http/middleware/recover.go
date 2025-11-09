package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/shanth1/gotools/log"
)

func WithRecovery(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().
						Any("panic", err).
						Str("stack", string(debug.Stack())).
						Msg("recovered from panic")
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
