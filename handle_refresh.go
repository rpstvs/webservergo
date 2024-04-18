package main

import (
	"net/http"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	token := CheckRefreshToken(tokenString)

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	token := CheckRefreshToken(tokenString)

}

func CheckRefreshToken(token string) string {
	if token == "" {
		return ""
	}

	tokenString := token[len("Bearer: "):]

	return tokenString
}
