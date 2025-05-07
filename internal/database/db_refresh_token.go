package database

import "fmt"

type Token struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	Revoke bool   `json:"revoke"`
}

func (db *DB) CreateTokenDB(token string) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	id := len(dbStructure.Tokens) + 1

	tokenDb := Token{
		ID:    id,
		Token: token,
	}
	dbStructure.Tokens[id] = tokenDb

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RevokeToken(tokenIn string) error {
	dbStructucture, err := db.loadDB()
	if err != nil {
		return err
	}

	for i, token := range dbStructucture.Tokens {

		if token.Token == tokenIn {
			fmt.Println("estou aqui dentro a rolar")
			token.Revoke = true
			dbStructucture.Tokens[i] = token
		}
	}
	err = db.writeDB(dbStructucture)

	return err
}

func (db *DB) GetToken(tokenIn string) (Token, error) {
	dbStructucture, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}

	for _, token := range dbStructucture.Tokens {

		if token.Token == tokenIn {
			return token, nil
		}
	}

	return Token{}, err
}
