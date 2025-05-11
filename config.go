package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/rpstvs/webservergo/internal/database"
)

func GetConfig() apiConfig {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Couldnt load env variables")
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("Postgres", dbURL)

	if err != nil {
		log.Fatal("Couldnt open DB connection")
	}

	return apiConfig{
		Platform:       os.Getenv("PLATFORM"),
		tokenSecret:    os.Getenv("TOKEN_SECRET"),
		dbQueries:      database.New(db),
		fileServerHits: atomic.Int32{},
	}
}
