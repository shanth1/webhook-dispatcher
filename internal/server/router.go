package server

import "net/http"

func (s *server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhook", s.handleWebhook)
	mux.HandleFunc("GET /health", s.handleHealthCheck)

	mux.HandleFunc("/", s.handleRoot)

	return mux
}
