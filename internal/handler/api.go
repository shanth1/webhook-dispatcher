package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/telehook/internal/config"
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
		logger.Error().Msg("signature is empty")
		http.Error(w, "Signature is empty", http.StatusBadRequest)
	}
	if !verifySignature(body, h.cfg.WebhookSecret, signature) {
		logger.Error().Msg("invalid webhook signature")
		http.Error(w, "Invalid signature", http.StatusForbidden)
		return
	}

	eventName := r.Header.Get("X-GitHub-Event")
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
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

	var wg sync.WaitGroup
	for _, bot := range h.cfg.Telegram.Bots {
		for _, client := range bot.Clients {
			wg.Add(1)
			go func(client config.Client, token, text string) {
				defer wg.Done()
				if err := h.sender.Send(token, client.ChatID, text); err != nil {
					logger.Error().Str("client", client.Name).Str("bot", token[:5]).Err(err).Msg("send message")
				}
			}(client, bot.Token, message.String())
		}
	}
	wg.Wait()

	logger.Info().Msgf("Event '%s' processed and sent to configured chats.", eventName)
	w.WriteHeader(http.StatusOK)
}

func verifySignature(body []byte, secret string, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
