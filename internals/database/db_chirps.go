package database

type Chirp struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_id int    `json:"author_id"`
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Body:      body,
		ID:        id,
		Author_id: author_id,
	}
	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	for _, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			return chirp, nil
		}
	}
	return Chirp{}, err
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}

	return nil

}
