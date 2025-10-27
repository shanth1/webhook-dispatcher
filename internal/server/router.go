package server

import (
	"net/http"

	"github.com/shanth1/gitrelay/internal/verifier"
)

func (s *server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", s.handleGithubWebhook(&verifier.GithubVerifier{Secret: s.cfg.WebhookSecret}))
	mux.HandleFunc("GET /health", s.handleHealthCheck)

	mux.HandleFunc("/", s.handleRoot)

	return mux
}
