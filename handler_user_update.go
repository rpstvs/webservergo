package main

import (
	"encoding/json"
	"net/http"

	"github.com/rpstvs/webservergo/internal/auth"
	"github.com/rpstvs/webservergo/internal/database"
)

func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode parameters", nil)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldt get token", nil)
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

	passHashed, _ := auth.PassHash(params.Password)

	updated, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: passHashed,
		ID:             id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "nao deu", nil)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		Email: params.Email,
		ID:    updated.ID,
	})
}
