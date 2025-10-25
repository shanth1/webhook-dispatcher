package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/shanth1/gotools/log"
)

func (h *handler) webhookHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(
		log.Str("method", r.Method),
		log.Str("remote_addr", r.RemoteAddr),
	)

	if r.Method != http.MethodPost {
		logger.Warn().Msg("invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error().Err(err).Msg("read request body")
		http.Error(w, "Can't read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		logger.Warn().Msg("signature is empty")
		http.Error(w, "Signature is empty", http.StatusBadRequest)
		return
	}
	if !verifySignature(body, h.cfg.WebhookSecret, signature) {
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
	templateName := eventName + ".tmpl"
	tmpl := h.templates.Lookup(templateName)
	if tmpl == nil {
		tmpl = h.templates.Lookup("default.tmpl")
	}
	if err := tmpl.Execute(&message, payload); err != nil {
		logger.Error().Err(err).Msg("executing template")
		http.Error(w, "Error formatting message", http.StatusInternalServerError)
		return
	}

	h.notifier.Broadcast(r.Context(), h.cfg.Recipients, message.String())

	logger.Info().Str("event", eventName).Int("recipients", len(h.cfg.Recipients)).Msg("Event processed and broadcasted.")
	w.WriteHeader(http.StatusOK)
}

func verifySignature(body []byte, secret string, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
