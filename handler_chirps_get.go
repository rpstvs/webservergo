package main

import (
	"net/http"
	"sort"
	"strconv"
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

func (cfg *apiConfig) retrieveChirpsId(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	idString := r.PathValue("chirpsid")
	id, err := strconv.Atoi(idString)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "Couldn't retrieve with that id")
		return
	}

	chirp := Chirp{}

	if len(dbChirps) < id {
		w.WriteHeader(http.StatusNotFound)
	}

	for _, dbChirp := range dbChirps {

		if dbChirp.ID == id {
			chirp = Chirp{
				ID:        dbChirp.ID,
				Body:      dbChirp.Body,
				Author_id: dbChirp.Author_id,
			}

		}

	}

	respondwithJSON(w, http.StatusOK, chirp)
}
