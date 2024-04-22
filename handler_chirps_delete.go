package main

import (
	"net/http"
	"strconv"

	"github.com/rpstvs/webservergo/internals/auth"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	dbChirps, _ := cfg.DB.GetChirps()

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	tokenString = tokenString[len("Bearer "):]
	id, _ := auth.ValidateToken(tokenString, cfg.secret)

	author_id, _ := strconv.Atoi(id)

	id2 := r.PathValue("chirpsid")
	chirpId, _ := strconv.Atoi(id2)

	if chirpId > len(dbChirps) {
		w.WriteHeader(http.StatusNotFound)
	}

	chirp, _ := cfg.DB.GetChirpById(chirpId)

	if chirp.Author_id == author_id {
		cfg.DB.DeleteChirp(chirpId)
		w.WriteHeader(http.StatusOK)
	} else {
		respondwithError(w, http.StatusForbidden, "Forbidden")
	}

}
