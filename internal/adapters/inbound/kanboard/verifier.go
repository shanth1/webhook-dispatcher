package kanboard

import "github.com/shanth1/hookrelay/internal/core/ports"

func (h *Handler) verify(req ports.WebhookRequest) bool {
	tokenFromRequest := req.Headers["token"]
	return tokenFromRequest != "" && tokenFromRequest == h.secret
}
