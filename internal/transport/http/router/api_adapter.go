package router

import (
	"encoding/json"
	"net/http"
	"sort"
)

func (rt *Router) handleAdaptersList(w http.ResponseWriter, r *http.Request) {
	adapterKeys := make([]string, 0, len(rt.processors))
	for k := range rt.processors {
		adapterKeys = append(adapterKeys, string(k))
	}
	sort.Strings(adapterKeys)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(adapterKeys); err != nil {
		rt.logger.Error().Err(err).Msg("failed to encode adapters list to json")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
