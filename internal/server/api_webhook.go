package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/shanth1/gitrelay/internal/utils"
	"github.com/shanth1/gitrelay/internal/verifier"
	"github.com/shanth1/gotools/log"
)

func (s *server) handleGithubWebhook(verifier verifier.Verifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.With(
			log.Str("method", r.Method),
			log.Str("remote_addr", r.RemoteAddr),
		)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("read request body")
			http.Error(w, "Can't read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if ok := verifier.Verify(r, body); !ok {
			logger.Error().Msg("invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}

		form, err := url.ParseQuery(string(body))
		if err != nil {
			logger.Error().Err(err).Msg("parsing form payload")
			http.Error(w, "Error parsing form payload", http.StatusBadRequest)
			return
		}

		payloadJSON := form.Get("payload")
		if payloadJSON == "" {
			logger.Error().Msg("payload field is empty")
			http.Error(w, "Payload field is empty", http.StatusBadRequest)
			return
		}

		eventName := r.Header.Get("X-GitHub-Event")
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			logger.Error().Err(err).Msg("parsing json payload")
			http.Error(w, "Error parsing JSON payload", http.StatusBadRequest)
			return
		}
		payload["eventName"] = eventName

		var message bytes.Buffer
		tmpl := s.templates.Lookup(utils.GetTemplatePath("github", eventName))
		if tmpl == nil {
			tmpl = s.templates.Lookup(utils.GetTemplatePath("github", "default"))
		}
		if err := tmpl.Execute(&message, payload); err != nil {
			logger.Error().Err(err).Msg("executing template")
			http.Error(w, "Error formatting message", http.StatusInternalServerError)
			return
		}

		s.notifier.Broadcast(r.Context(), s.cfg.Recipients, message.String())

		logger.Info().Str("event", eventName).Int("recipients", len(s.cfg.Recipients)).Msg("Event processed and broadcasted.")
		w.WriteHeader(http.StatusOK)
	}
}
