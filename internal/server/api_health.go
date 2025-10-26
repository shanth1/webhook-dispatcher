package server

import (
	"net/http"

	"github.com/shanth1/gotools/log"
)

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With(
		log.Str("method", r.Method),
		log.Str("remote_addr", r.RemoteAddr),
	)

	logger.Info().Msg("health check requested")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
