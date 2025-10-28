package server

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/processor"
)

func (s *server) handleWebhook(p processor.WebhookProcessor, recipients []config.Recipient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.FromContext(r.Context())

		message, err := p.Process(r)
		if err != nil {
			logger.Error().Err(err).Msg("failed to process webhook")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.notifier.Broadcast(r.Context(), recipients, message)

		logger.Info().Int("recipients", len(recipients)).Msg("Webhook processed and broadcasted.")
		w.WriteHeader(http.StatusOK)
	}
}
