package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var id int

func createChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
		Id          int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondwithError(w, 400, "Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanMsg := badWordReplacement(params.Body, badWords)
	id++
	respondwithJSON(w, 201, returnVals{CleanedBody: cleanMsg, Id: id})

}

func respondwithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondwithJSON(w, code, errorResponse{Error: msg})
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func badWordReplacement(s string, badWords map[string]struct{}) string {

	replacement := "****"
	s2 := strings.Split(s, " ")
	for i, word := range s2 {

		_, ok := badWords[strings.ToLower(word)]
		if ok {
			s2[i] = replacement
		}

	}
	str3 := strings.Join(s2, " ")

	return str3

}
