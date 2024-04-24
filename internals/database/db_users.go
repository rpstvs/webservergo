package database

import (
	"errors"
)

type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, password string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		Email:         email,
		ID:            id,
		Password:      password,
		Is_Chirpy_Red: false,
	}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(dbStructure.Users))

	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetuserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err

	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("resource does not exist")
}

func (db *DB) GetUserId(id int) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err

	}

	user, ok := dbStructure.Users[id]

	if !ok {
		return User{}, errors.New("resource does not exist")
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()

	user := dbStructure.Users[id]

	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.Password = password
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUser(id int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	user := dbStructure.Users[id]

	user.Is_Chirpy_Red = true

	dbStructure.Users[id] = user

	db.writeDB(dbStructure)

	return nil

}
