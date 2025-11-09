package router

import (
	"encoding/json"
	"net/http"
)

func (rt *Router) handleNotifierList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rt.notifierTypes); err != nil {
		rt.logger.Error().Err(err).Msg("failed to encode adapters list to json")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
