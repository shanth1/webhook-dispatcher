package kanboard

import (
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
	tokenFromRequest := data.Headers["token"]
	return tokenFromRequest != "" && tokenFromRequest == a.secret
}
