package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	Email    string `json:"email"`
	Password string `json:"-,omitempty"`
	ID       int    `json:"id"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldnt create user")
		return
	}

	passHashed, _ := passHash(params.Password)

	user, err := cfg.DB.CreateUser(params.Email, passHashed)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondwithJSON(w, http.StatusCreated, users{
		Email: user.Email,
		ID:    user.ID,
	})

}

func passHash(password string) (string, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("couldnt hash the password")
	}

	return string(passHash), nil
}
