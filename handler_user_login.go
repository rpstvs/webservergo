package main

import (
	"encoding/json"
	"net/http"

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

	userToken, _ := cfg.createToken(User{
		userLogging.ID,
		params.Email,
		params.Password,
		"",
	})

	respondwithJSON(w, http.StatusOK, response{User: User{
		Email: params.Email,
		ID:    userLogging.ID,
		Token: userToken,
	}})
}

func (cfg *apiConfig) protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	err := cfg.ValidateToken(tokenString)

	if err != nil {
		respondwithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

}
