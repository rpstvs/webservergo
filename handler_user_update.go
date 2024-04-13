package main

import "net/http"

func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
