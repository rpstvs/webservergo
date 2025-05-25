package main

import (
	"encoding/json"
	"net/http"

	"github.com/rpstvs/webservergo/internals/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode parameters", nil)
		return
	}

	userLogging, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", nil)
	}

	err = auth.CheckPasswordHash(params.Password, params.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", nil)
		return
	}

	userToken, _ := auth.CreateToken(userLogging.ID, cfg.tokenSecret)
	refreshToken, _ := auth.CreateRefreshToken(userLogging.ID, cfg.tokenSecret)

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:    userLogging.ID,
			Email: userLogging.Email,
		},
		Token:        userToken,
		RefreshToken: refreshToken,
	})

}
