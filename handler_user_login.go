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
		respondwithError(w, http.StatusInternalServerError, "couldnt decode parameters")
		return
	}

	userLogging, err := cfg.DB.GetuserByEmail(params.Email)

	if err != nil {
		respondwithError(w, http.StatusNotFound, "user not found")
	}

	err = auth.CheckPasswordHash(params.Password, userLogging.Password)

	if err != nil {
		respondwithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	userToken, _ := auth.CreateToken(userLogging.ID, cfg.secret)
	refreshToken, _ := auth.CreateRefreshToken(userLogging.ID, cfg.secret)

	cfg.DB.CreateTokenDB(refreshToken)

	respondwithJSON(w, http.StatusOK, response{
		User: User{
			ID:    userLogging.ID,
			Email: userLogging.Email,
		},
		Token:        userToken,
		RefreshToken: refreshToken,
	})

}
