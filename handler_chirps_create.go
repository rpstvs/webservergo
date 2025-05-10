package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/rpstvs/webservergo/internal/auth"
)

type Chirp struct {
	Body      string `json:"body"`
	ID        int    `json:"id"`
	Author_id int    `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		respondwithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	tokenString = tokenString[len("Bearer "):]

	id, err := auth.ValidateToken(tokenString, cfg.secret)

	if err != nil {
		respondwithError(w, http.StatusUnauthorized, "este burro nao entra")
		return
	}

	realId, _ := strconv.Atoi(id)

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondwithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, realId)
	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondwithJSON(w, http.StatusCreated, Chirp{
		Body:      chirp.Body,
		ID:        chirp.ID,
		Author_id: chirp.Author_id,
	})

}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedbody(body, badWords)
	return cleaned, nil
}

func getCleanedbody(body string, badWords map[string]struct{}) string {
	replacement := "****"
	s2 := strings.Split(body, " ")
	for i, word := range s2 {
		_, ok := badWords[strings.ToLower(word)]
		if ok {
			s2[i] = replacement
		}
	}
	str3 := strings.Join(s2, " ")
	return str3

}
