package router

import (
	"net/http"

	"github.com/shanth1/hookrelay/internal/adapters/inbound/custom"
	"github.com/shanth1/hookrelay/internal/adapters/inbound/github"
	"github.com/shanth1/hookrelay/internal/adapters/inbound/kanboard"
	"github.com/shanth1/hookrelay/internal/config"
	httptransport "github.com/shanth1/hookrelay/internal/transport/http"
	"github.com/shanth1/hookrelay/internal/transport/http/middleware"
)

func (rt *Router) registerWebhookRoutes(mux *http.ServeMux) {
	recipientMap := make(map[string]config.Recipient)
	for _, r := range rt.cfg.Recipients {
		recipientMap[r.Name] = r
	}

	for _, hookCfg := range rt.cfg.Webhooks {
		hookCfg := hookCfg

		var adapter httptransport.InboundAdapter
		switch hookCfg.Type {
		case config.WebhookTypeGitHub:
			adapter = github.NewAdapter(hookCfg)
		case config.WebhookTypeKanboard:
			adapter = kanboard.NewAdapter(hookCfg)
		case config.WebhookTypeCustom:
			adapter = custom.NewAdapter(hookCfg)
		default:
			rt.logger.Error().Str("type", string(hookCfg.Type)).Msg("unknown webhook type")
			continue
		}

		resolvedRecipients := resolveRecipients(hookCfg.Recipients, recipientMap)
		handler := httptransport.NewWebhookHandler(rt.service, adapter, hookCfg.Type, resolvedRecipients)
		chain := middleware.Chain(
			handler,
			middleware.WithMethod(http.MethodPost),
		)

		rt.logger.Info().Str("path", hookCfg.Path).Str("type", string(hookCfg.Type)).Msg("registering webhook handler")
		mux.Handle("POST "+hookCfg.Path, chain)
	}
}
