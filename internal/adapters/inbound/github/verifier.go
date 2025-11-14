package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/shanth1/hookrelay/internal/core/ports"
)

func (h *Handler) verify(req ports.WebhookRequest) bool {
	signature := req.Headers["x-hub-signature-256"]
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(req.Payload)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
