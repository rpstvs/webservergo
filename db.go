package main

import (
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:id`
	Body string `json: body`
}

func NewDB(path string) (*DB, error) {
	if path == "" {
		return nil, errors.New("path can't be null")
	}

	if _, err := os.Stat("database.json"); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	return db, nil
}
