package router

import (
	"net/http"

	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/transport/http/handler"
	"github.com/shanth1/hookrelay/internal/transport/http/middleware"
)

func (rt *Router) registerWebhookRoutes(mux *http.ServeMux) {
	recipientMap := make(map[string]config.Recipient)
	for _, r := range rt.cfg.Recipients {
		recipientMap[r.Name] = r
	}

	for _, hookCfg := range rt.cfg.Webhooks {
		hookCfg := hookCfg

		resolvedRecipients := resolveRecipients(hookCfg.Recipients, recipientMap)
		handler := handler.NewWebhookHandler(rt.service, hookCfg, resolvedRecipients)
		chain := middleware.Chain(
			handler,
			middleware.WithMethod(http.MethodPost),
		)

		rt.logger.Info().Str("path", hookCfg.Path).Str("type", string(hookCfg.Type)).Msg("registering webhook handler")
		mux.Handle("POST "+hookCfg.Path, chain)
	}
}
