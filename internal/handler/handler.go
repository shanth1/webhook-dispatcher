package handler

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/shanth1/gitrelay/internal/config"
	"github.com/shanth1/gitrelay/internal/service"
	"github.com/shanth1/gotools/log"
)

type handler struct {
	cfg       *config.Config
	logger    log.Logger
	templates *template.Template
	notifier  *service.Notifier
}

func New(cfg *config.Config, templates *template.Template, notifier *service.Notifier, logger log.Logger) *handler {
	return &handler{
		cfg:       cfg,
		logger:    logger,
		notifier:  notifier,
		templates: templates,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/webhook") {
		h.webhookHandler(w, r)
		return
	}
}
