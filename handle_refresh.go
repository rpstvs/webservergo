package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rpstvs/webservergo/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString := r.Header.Get("Authorization")

	refToken := CheckRefreshToken(tokenString)

	issuer, _ := auth.GetIssuerr(refToken, cfg.tokenSecret)

	tokenDb, _ := cfg.DB.GetToken(refToken)

	fmt.Println(issuer)

	if issuer == "chirpy-refresh" && !tokenDb.Revoke {
		idString, _ := auth.ValidateJWT(refToken, cfg.tokenSecret)
		id, _ := strconv.Atoi(idString)
		token, _ := auth.CreateToken(id, cfg.tokenSecret)

		respondWithJson(w, http.StatusOK, response{
			Token: token,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Not valid refresh token", nil)
	}

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	token := CheckRefreshToken(tokenString)

	cfg.DB.RevokeToken(token)

}

func CheckRefreshToken(token string) string {
	if token == "" {
		return ""
	}

	tokenString := token[len("Bearer "):]

	return tokenString
}
