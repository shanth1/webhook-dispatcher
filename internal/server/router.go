package server

import (
	"net/http"

	"github.com/shanth1/gitrelay/internal/middleware"
	"github.com/shanth1/gitrelay/internal/verifier"
)

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", middleware.Chain(
		http.HandlerFunc(s.handleHealthCheck),
		middleware.WithLogger(s.logger),
		middleware.WithMethod(http.MethodGet),
	))

	mux.Handle("POST /webhook", middleware.Chain(
		http.HandlerFunc(s.handleGithubWebhook),
		middleware.WithLogger(s.logger),
		middleware.WithMethod(http.MethodPost),
		middleware.WithVerification(&verifier.GithubVerifier{Secret: s.cfg.WebhookSecret}),
	))

	mux.Handle("/", middleware.Chain(
		http.HandlerFunc(s.handleRoot),
		middleware.WithLogger(s.logger),
	))

	return mux
}
