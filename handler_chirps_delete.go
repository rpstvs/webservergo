package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rpstvs/webservergo/internal/auth"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token on request", nil)
		return
	}

	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization", nil)
		return
	}

	id, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not valid", nil)
		return
	}

	id2 := r.PathValue("chirpID")
	chirpId, _ := uuid.Parse(id2)

	chirp, err := cfg.dbQueries.GetChirpId(r.Context(), chirpId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", nil)
		return
	}

	if chirp.UserID == id {
		cfg.dbQueries.DeleteChirp(r.Context(), chirpId)
		w.WriteHeader(http.StatusNoContent)
	} else {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
	}

}
