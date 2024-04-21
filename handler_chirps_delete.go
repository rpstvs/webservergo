package main

import (
	"net/http"
	"strconv"

	"github.com/rpstvs/webservergo/internals/auth"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	dbChirps, _ := cfg.DB.GetChirps()

	tokenString := w.Header().Get("Authorization")

	tokenString = tokenString[len("Bearer "):]

	id, _ := auth.ValidateToken(tokenString, cfg.secret)

	author_id, _ := strconv.Atoi(id)

	id2 := r.PathValue("chirpsid")
	chirpId, _ := strconv.Atoi(id2)

	if chirpId > len(dbChirps) {
		w.WriteHeader(http.StatusNotFound)
	}

	for _, dbChirp := range dbChirps {

		if dbChirp.Author_id == author_id {
			delete(dbChirp, chirpId)
		}
	}

}
