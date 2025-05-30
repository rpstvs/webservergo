package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rpstvs/webservergo/internal/auth"
	"github.com/rpstvs/webservergo/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create user", nil)
		return
	}

	passHashed, _ := auth.PassHash(params.Password)

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: passHashed,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user in db", nil)
		return
	}

	respondWithJson(w, 201, returnVals{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})
}
