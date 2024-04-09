package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PassHash(password string) (string, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("couldnt hash the password")
	}

	return string(passHash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
