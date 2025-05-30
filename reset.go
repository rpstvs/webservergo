package main

import (
	"net/http"
)

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {

	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.dbQueries.DeleteUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt delete users", err)
		return
	}
}
