package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createToken(id int) (string, error) {

	signinKey := []byte(cfg.secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   fmt.Sprintf("%d", id),
	})
	fmt.Println("vou criar um token")

	return token.SignedString(signinKey)

}

func (cfg *apiConfig) createRefreshToken(id int) (string, error) {

	signinKey := []byte(cfg.secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(60 * 24 * time.Hour)),
		Subject:   fmt.Sprintf("%d", id),
	})
	fmt.Println("vou criar um token")

	return token.SignedString(signinKey)

}

func (cfg *apiConfig) ValidateToken(tokenstring string) (string, error) {

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {

		return []byte(cfg.secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	id, _ := token.Claims.GetSubject()
	issuer, _ := token.Claims.GetIssuer()

	if issuer == "chirpy-refresh" {
		return "", errors.New("no access")
	}

	fmt.Println("acesso valido")

	return id, nil
}
