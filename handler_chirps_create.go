package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rpstvs/webservergo/internal/auth"
	"github.com/rpstvs/webservergo/internal/database"
)

type Chirp struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}

func (api *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No token", nil)
		return
	}

	user, err := auth.ValidateJWT(token, api.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var params parameters

	err = decoder.Decode(&params)

	if err != nil {
		log.Println("Couldnt decode the info")
		return
	}

	str, err := validateChirp(params.Body)

	if err != nil {
		log.Println("invalid chirp")
		return

	}

	chirp, err := api.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   str,
		UserID: user,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create chirp", err)
		return
	}
	respondWithJson(w, 201, response{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLen = 140

	if len(body) > maxChirpLen {
		return "", errors.New("chirp too long")
	}

	badWords := map[string]struct{}{
		"fornax":    {},
		"sharbert":  {},
		"kerfuffle": {},
	}

	cleaned := cleanBody(body, badWords)
	return cleaned, nil
}

func cleanBody(chirp string, badWords map[string]struct{}) string {
	replacement := "****"
	words := strings.Split(chirp, " ")
	for i, word := range words {
		_, ok := badWords[strings.ToLower(word)]

		if ok {
			words[i] = replacement
		}
	}

	new_string := strings.Join(words, " ")

	return new_string
}
