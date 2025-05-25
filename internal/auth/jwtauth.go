package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func CreateToken(userid uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	signinKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   string(userid.String()),
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

func ValidateJWT(tokenstring, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenstring, &claimsStruct, func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid userid")
	}

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

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", errors.New("no token")
	}

	parts := strings.Split(bearer, " ")

	if len(parts) != 2 {
		return "", errors.New("malformed auth token")
	}

	return parts[1], nil
}
