package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rpstvs/webservergo/internal/database"
)

func GetConfig() apiConfig {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Couldnt load env variables")
	}

	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)

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
