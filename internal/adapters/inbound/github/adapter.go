package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/shanth1/hookrelay/internal/config"
	httptransport "github.com/shanth1/hookrelay/internal/transport/http"
)

type Adapter struct {
	secret string
}

var _ httptransport.InboundAdapter = (*Adapter)(nil)

func NewAdapter(hookCfg config.WebhookConfig) httptransport.InboundAdapter {
	return &Adapter{secret: hookCfg.Secret}
}

func (a *Adapter) Verify(data httptransport.VerificationData) bool {
	signature := data.Headers["X-Hub-Signature-256"]
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(a.secret))
	mac.Write(data.Body)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
