package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func CreateToken(id int, tokenSecret string, expiresIn time.Duration) (string, error) {

	signinKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", id),
	})

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

	return id, nil
}

func GetIssuerr(tokenstring, tokenSecret string) (string, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return "vazio1", err
	}

	if !token.Valid {
		return "vazio2", errors.New("invalid token")
	}

	issuer, _ := token.Claims.GetIssuer()

	return issuer, nil
}
