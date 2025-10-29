package server

import (
	"context"
	"html/template"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/hookrelay/internal/config"
	"github.com/shanth1/hookrelay/internal/notifier"
	"github.com/shanth1/hookrelay/internal/processor"
)

type server struct {
	cfg        *config.Config
	logger     log.Logger
	templates  *template.Template
	notifier   *notifier.Notifier
	httpServer *http.Server
	processors map[config.WebhookType]processor.WebhookProcessor
}

func New(cfg *config.Config, templates *template.Template, notifier *notifier.Notifier, logger log.Logger) *server {
	s := &server{
		cfg:       cfg,
		logger:    logger,
		templates: templates,
		notifier:  notifier,
	}

	s.processors = map[config.WebhookType]processor.WebhookProcessor{
		config.WebhookTypeGitHub:   processor.NewGithubProcessor(templates),
		config.WebhookTypeKanboard: processor.NewKanboardProcessor(templates),
		config.WebhookTypeCustom:   processor.NewCustomProcessor(),
	}

	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: s.routes(),
	}

	s.httpServer = httpServer

	return s
}

func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
