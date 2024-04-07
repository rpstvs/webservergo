package main

import (
	"encoding/json"
	"net/http"
)

type users struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldnt create user")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondwithJSON(w, http.StatusCreated, users{
		Email: user.Email,
		ID:    user.ID,
	})

}
