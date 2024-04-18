package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	Arr := strings.Split(tokenString, " ")

	if len(Arr) < 2 {
		return
	}

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	tokenString = tokenString[len("Bearer: "):]

}
