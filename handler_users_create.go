package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rpstvs/webservergo/internals/auth"
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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldnt create user")
		return
	}

	passHashed, _ := auth.PassHash(params.Password)

	user, err := cfg.DB.CreateUser(params.Email, passHashed)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondwithJSON(w, http.StatusCreated, User{
		Email:         user.Email,
		ID:            user.ID,
		Is_Chirpy_Red: user.Is_Chirpy_Red,
	})

}
