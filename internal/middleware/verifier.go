package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/verifier"
)

func WithVerification(v verifier.Verifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.FromContext(r.Context())

			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error().Err(err).Msg("failed to read request body for verification")
				http.Error(w, "Can't read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			if ok := v.Verify(r, body); !ok {
				logger.Warn().Msg("invalid webhook signature")
				http.Error(w, "Invalid signature", http.StatusForbidden)
				return
			}

			logger.Debug().Msg("webhook signature verified successfully")
			next.ServeHTTP(w, r)
		})
	}
}
