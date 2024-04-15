package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rpstvs/webservergo/internals/auth"
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
		respondwithError(w, http.StatusInternalServerError, "couldnt decode parameters")
		return
	}

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	id, err := cfg.ValidateToken(tokenString)

	if err != nil {
		respondwithError(w, http.StatusUnauthorized, "este burro nao entra")
		return
	}
	passHashed, _ := auth.PassHash(params.Password)
	realId, _ := strconv.Atoi(id)

	updated, err := cfg.DB.UpdateUser(realId, params.Email, passHashed)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "nao deu")
		return
	}

	respondwithJSON(w, http.StatusOK, User{
		Email: params.Email,
		ID:    updated.ID,
	})
}
