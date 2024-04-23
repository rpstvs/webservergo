package main

import (
	"encoding/json"
	"net/http"

	"github.com/rpstvs/webservergo/internals/auth"
)

type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"-"`
	Token         string `json:"token,omitempty"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
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
