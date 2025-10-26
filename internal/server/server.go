package server

import (
	"context"
	"html/template"
	"net/http"

	"github.com/shanth1/gitrelay/internal/config"
	"github.com/shanth1/gitrelay/internal/service"
	"github.com/shanth1/gotools/log"
)

type server struct {
	cfg        *config.Config
	logger     log.Logger
	templates  *template.Template
	notifier   *service.Notifier
	httpServer *http.Server
}

func New(cfg *config.Config, templates *template.Template, notifier *service.Notifier, logger log.Logger) *server {
	s := &server{
		cfg:       cfg,
		logger:    logger,
		templates: templates,
		notifier:  notifier,
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
