package main

import (
	"net/http"
	"time"

	"github.com/rpstvs/webservergo/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt get refresh token from header", nil)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), tokenString)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token", nil)

		return
	}

	accessToken, err := auth.CreateToken(user.ID, cfg.tokenSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldt create token", nil)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no refresh token", nil)
		return
	}

	_, err = cfg.dbQueries.RevokeRefreshToken(r.Context(), tokenString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant revoke token", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
