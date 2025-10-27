package server

import (
	"net/http"

	"github.com/shanth1/gitrelay/internal/middleware"
	"github.com/shanth1/gitrelay/internal/verifier"
)

func (s *server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", middleware.WithMethod(http.MethodPost,
		s.handleGithubWebhook(&verifier.GithubVerifier{Secret: s.cfg.WebhookSecret})),
	)
	mux.HandleFunc("GET /health", middleware.WithMethod(http.MethodGet, s.handleHealthCheck))

	mux.HandleFunc("/", s.handleRoot)

	return mux
}
