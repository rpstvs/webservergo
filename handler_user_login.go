package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rpstvs/webservergo/internals/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
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

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	userToken, _ := cfg.createToken(userLogging.ID, time.Duration(params.ExpiresInSeconds)*time.Second)

	fmt.Println(userToken)

	respondwithJSON(w, http.StatusOK, response{
		User: User{
			ID:    userLogging.ID,
			Email: userLogging.Email,
		},
		Token: userToken,
	})

}
