package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func createClaims(w http.ResponseWriter, r *http.Request) {

	token := jwt.NewWithClaims()

}
