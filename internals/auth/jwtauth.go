package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id int, tokenSecret string) (string, error) {

	signinKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   fmt.Sprintf("%d", id),
	})
	fmt.Println("vou criar um token")

	return token.SignedString(signinKey)

}

func CreateRefreshToken(id int, tokenSecret string) (string, error) {

	signinKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(60 * 24 * time.Hour)),
		Subject:   fmt.Sprintf("%d", id),
	})
	fmt.Println("vou criar um token")

	return token.SignedString(signinKey)

}

func ValidateToken(tokenstring, tokenSecret string) (string, error) {

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {

		return []byte(tokenSecret), nil
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
