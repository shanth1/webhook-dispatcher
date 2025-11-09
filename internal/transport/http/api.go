package httptransport

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type API struct {
	service ports.Service
	logger  log.Logger
	config  *config.Config
}

func NewAPI(service ports.Service, logger log.Logger, cfg *config.Config) *API {
	return &API{
		service: service,
		logger:  logger,
		config:  cfg,
	}
}

func (a *API) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *API) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	a.handleListWebhooks(w, r)
}

func (a *API) handleListWebhooks(w http.ResponseWriter, r *http.Request) {
	webhookTypes := make(map[string]struct{})
	for _, hook := range a.config.Webhooks {
		webhookTypes[string(hook.Type)] = struct{}{}
	}

	result := make([]string, 0, len(webhookTypes))
	for t := range webhookTypes {
		result = append(result, t)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		a.logger.Error().Err(err).Msg("failed to encode webhooks list to json")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *API) webhookHandlerFactory(webhookCfg config.WebhookConfig) http.HandlerFunc {
	recipientMap := make(map[string]config.Recipient)
	for _, r := range a.config.Recipients {
		recipientMap[r.Name] = r
	}

	resolvedRecipients := make([]config.Recipient, 0, len(webhookCfg.Recipients))
	for _, name := range webhookCfg.Recipients {
		if r, ok := recipientMap[name]; ok {
			resolvedRecipients = append(resolvedRecipients, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.FromContext(r.Context())

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("failed to read request body")
			http.Error(w, "Cannot read request body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		inboundReq := ports.WebhookRequest{
			Payload: body,
			Headers: extractHeaders(r.Header),
			Params:  extractQueryParams(r.URL.Query()),
		}

		if err := a.service.ProcessWebhook(r.Context(), webhookCfg.Name, inboundReq, resolvedRecipients); err != nil {
			logger.Warn().Err(err).Str("webhook_name", string(webhookCfg.Name)).Msg("failed to process webhook")
			http.Error(w, "Webhook processing failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func extractHeaders(h http.Header) map[string]string {
	headers := make(map[string]string)
	for key, values := range h {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return headers
}

func extractQueryParams(q map[string][]string) map[string]string {
	params := make(map[string]string)
	for key, values := range q {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}
