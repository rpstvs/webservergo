package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
