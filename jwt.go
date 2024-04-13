package main

import (
	"errors"
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

	tokenstring, err := token.SignedString(cfg.secret)

	if err != nil {
		return "", nil
	}

	return tokenstring, nil

}

func (cfg *apiConfig) ValidateToken(tokenstring string) error {

	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {

		return cfg.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
