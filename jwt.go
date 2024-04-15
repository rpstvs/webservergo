package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createToken(id int, ExpiresIn time.Duration) (string, error) {

	signinKey := []byte(cfg.secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ExpiresIn)),
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
	fmt.Println("acesso valido")

	return id, nil
}
