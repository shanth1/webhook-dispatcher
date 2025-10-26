package server

import (
	"net/http"

	"github.com/shanth1/gotools/log"
)

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With(
		log.Str("method", r.Method),
		log.Str("remote_addr", r.RemoteAddr),
	)

	if r.URL.Path != "/" {
		logger.Info().Str("path", r.URL.Path).Msg("unknown root")
		http.NotFound(w, r)
		return
	}
	logger.Info().Msg("root path requested")
	w.WriteHeader(http.StatusOK)
}
