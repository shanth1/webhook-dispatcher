package httptransport

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
)

type webhookHandler struct {
	service     ports.WebhookService
	adapter     InboundAdapter
	webhookType config.WebhookType
	recipients  []config.Recipient
}

func NewWebhookHandler(
	service ports.WebhookService,
	adapter InboundAdapter,
	webhookType config.WebhookType,
	recipients []config.Recipient,
) http.Handler {
	return &webhookHandler{
		service:     service,
		adapter:     adapter,
		webhookType: webhookType,
		recipients:  recipients,
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

	verificationData := VerificationData{
		Body:        body,
		Headers:     extractHeaders(r.Header),
		QueryParams: extractQueryParams(r.URL.Query()),
	}

	if ok := h.adapter.Verify(verificationData); !ok {
		logger.Warn().Msg("invalid webhook signature or token")
		http.Error(w, "Invalid signature or token", http.StatusForbidden)
		return
	}

	coreRequest := ports.WebhookRequest{
		WebhookType: h.webhookType,
		Payload:     body,
		Headers:     verificationData.Headers,
	}

	if err := h.service.ProcessWebhook(r.Context(), coreRequest, h.recipients); err != nil {
		logger.Error().Err(err).Msg("failed to process webhook")
		http.Error(w, err.Error(), http.StatusBadRequest)
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
