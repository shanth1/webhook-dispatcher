package handler

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/telehook/internal/config"
	"github.com/shanth1/telehook/internal/service"
)

type handler struct {
	cfg       *config.Config
	logger    log.Logger
	templates *template.Template
	sender    *service.TelegramSender
}

func New(cfg *config.Config, templates *template.Template, sender *service.TelegramSender, logger log.Logger) *handler {
	return &handler{
		cfg:       cfg,
		logger:    logger,
		sender:    sender,
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
