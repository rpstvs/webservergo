package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type userLogin struct {
	Password string `json:"-,omitempty"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	/*
	 receive a request, procurar um user por email
	 comparar passwords,
	 dar autorizaçao ou nao
	*/
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
	passwordDb, id, err := cfg.lookupEmail(params.Email)

	if err != nil {
		respondwithError(w, http.StatusNotFound, "user not found")
	}

	ok := bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(params.Password))

	if ok == nil {
		respondwithJSON(w, http.StatusOK, userLogin{
			Email: params.Email,
			ID:    id,
		})
	} else {
		respondwithError(w, http.StatusUnauthorized, "not authorized")
	}

	//
}

func (cfg *apiConfig) lookupEmail(email string) (string, int, error) {
	dbUsers, err := cfg.DB.GetUsers()

	if err != nil {
		return "", 0, errors.New("cenas")
	}

	for _, dbUser := range dbUsers {
		if dbUser.Email == email {
			return dbUser.Password, dbUser.ID, nil
		}
	}
	return "", 0, errors.New("email not found")
}
