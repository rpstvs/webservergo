package main

import (
	"net/http"
	"sort"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondwithJSON(w, http.StatusOK, chirps)

}
