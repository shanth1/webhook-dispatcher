package router

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/common"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/transport/http/middleware"
)

type Router struct {
	cfg           *config.Config
	logger        log.Logger
	service       ports.Service
	webhookTypes  []config.WebhookType
	notifierTypes []config.NotifierType
}

func New(
	cfg *config.Config,
	service ports.Service,
	logger log.Logger,
) *Router {
	webhookTypes := common.GetUniqueValues(cfg.Webhooks, func(c config.WebhookConfig) config.WebhookType { return c.Type })
	notifierTypes := common.GetUniqueValues(cfg.Notifiers, func(c config.NotifierConfig) config.NotifierType { return c.Type })

	return &Router{
		cfg:           cfg,
		service:       service,
		logger:        logger,
		webhookTypes:  webhookTypes,
		notifierTypes: notifierTypes,
	}
}

func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	rt.registerWebhookRoutes(mux)

	mux.Handle("GET /health", http.HandlerFunc(rt.handleHealthCheck))
	mux.Handle("GET /webhooks", http.HandlerFunc(rt.handleWebhookList))
	mux.Handle("GET /", http.HandlerFunc(rt.handleRoot))

	return middleware.Chain(
		mux,
		middleware.WithRecovery(rt.logger),
		middleware.WithLogger(rt.logger),
	)
}

func resolveRecipients(names []string, recipientMap map[string]config.Recipient) []config.Recipient {
	resolved := make([]config.Recipient, 0, len(names))
	for _, name := range names {
		if r, ok := recipientMap[name]; ok {
			resolved = append(resolved, r)
		}
	}
	return resolved
}
