package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type webhookHandler struct {
	service    ports.Service
	webhookCfg config.WebhookConfig
	recipients []config.Recipient
}

func NewWebhookHandler(
	service ports.Service,
	webhookCfg config.WebhookConfig,
	recipients []config.Recipient,
) http.Handler {
	return &webhookHandler{
		service:    service,
		webhookCfg: webhookCfg,
		recipients: recipients,
	}
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.ProcessWebhook(r.Context(), h.webhookCfg.Name, inboundReq, h.recipients); err != nil {
		logger.Warn().Err(err).Str("name", string(h.webhookCfg.Name)).Msg("failed to process webhook")
		http.Error(w, "Webhook processing failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
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
