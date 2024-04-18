package database

type Token struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
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
