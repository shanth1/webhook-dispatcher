package server

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/shanth1/gotools/log"
)

func (s *server) handleAdaptersList(w http.ResponseWriter, r *http.Request) {
	logger := log.FromContext(r.Context())

	adapterKeys := make([]string, 0, len(s.processors))
	for k := range s.processors {
		adapterKeys = append(adapterKeys, string(k))
	}

	sort.Strings(adapterKeys)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(adapterKeys); err != nil {
		logger.Error().Err(err).Msg("failed to encode adapters list to json")
	}
}
