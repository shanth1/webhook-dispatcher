package router

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/core/ports"
	"github.com/shanth1/hookrelay/internal/service/webhook"
	"github.com/shanth1/hookrelay/internal/transport/http/middleware"
)

type Router struct {
	logger     log.Logger
	cfg        *config.Config
	service    ports.WebhookService
	processors map[config.WebhookType]webhook.Processor
}

func New(
	cfg *config.Config,
	service ports.WebhookService,
	processors map[config.WebhookType]webhook.Processor,
	logger log.Logger,
) *Router {
	return &Router{
		cfg:        cfg,
		service:    service,
		processors: processors,
		logger:     logger,
	}
}

func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	rt.registerWebhookRoutes(mux)

	mux.Handle("GET /health", http.HandlerFunc(rt.handleHealthCheck))
	mux.Handle("GET /adapters", http.HandlerFunc(rt.handleAdaptersList))
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
