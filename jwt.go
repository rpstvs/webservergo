package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createToken(user User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(200 * time.Second)),
		Subject:   string(user.ID),
	})

	tokenstring, err := token.SignedString([]byte(cfg.secret))

	if err != nil {
		return "", errors.New("token invalid")
	}

	return tokenstring, nil

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
	fmt.Println("acesso valido")

	return id, nil
}
