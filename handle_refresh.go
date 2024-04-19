package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rpstvs/webservergo/internals/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString := r.Header.Get("Authorization")

	refToken := CheckRefreshToken(tokenString)

	issuer, _ := auth.GetIssuerr(refToken, cfg.secret)

	fmt.Println(issuer)

	if issuer == "chirpy-refresh" {
		idString, _ := auth.ValidateToken(refToken, cfg.secret)
		id, _ := strconv.Atoi(idString)
		token, _ := auth.CreateToken(id, cfg.secret)

		respondwithJSON(w, http.StatusOK, response{
			Token: token,
		})
	} else {
		respondwithError(w, http.StatusBadRequest, "Not valid refresh token")
	}

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {

	//tokenString := r.Header.Get("Authorization")

	//token := CheckRefreshToken(tokenString)

	//issuer, _ := auth.GetIssuer(token, cfg.secret)

}

func CheckRefreshToken(token string) string {
	if token == "" {
		return ""
	}

	tokenString := token[len("Bearer "):]

	return tokenString
}
