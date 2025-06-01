package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rpstvs/webservergo/internal/auth"
	"github.com/rpstvs/webservergo/internal/database"
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
		return
	}

	err = auth.CheckPasswordHash(params.Password, userLogging.HashedPassword)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", nil)
		return
	}

	userToken, err := auth.CreateToken(userLogging.ID, cfg.tokenSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create access jwt", nil)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create refresh token", nil)
		return
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userLogging.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt save refresh token", nil)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:          userLogging.ID,
			Email:       userLogging.Email,
			IsChirpyRed: userLogging.IsChirpyRed,
		},
		Token:        userToken,
		RefreshToken: refreshToken,
	})

}
