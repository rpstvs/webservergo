package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	se := r.URL.Query().Get("sort")

	if se != "desc" {
		se = "asc"
	}

	if s != "" {
		idInt, _ := strconv.Atoi(s)

		chirps, err := cfg.DB.GetChirpAuthor(idInt)

		chirpz := []Chirp{}

		for _, dbChirp := range chirps {
			chirpz = append(chirpz, Chirp{
				ID:   dbChirp.ID,
				Body: dbChirp.Body,
			})
		}

		if err != nil {
			return
		}

		response := sortChirps(chirpz, se)

		respondwithJSON(w, http.StatusOK, response)
	} else {

		dbChirps, err := cfg.DB.GetChirps()
		if err != nil {
			respondwithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
			return
		}

		response := []Chirp{}

		for _, dbChirp := range dbChirps {
			response = append(response, Chirp{
				ID:   dbChirp.ID,
				Body: dbChirp.Body,
			})
		}

		response = sortChirps(response, se)

		respondwithJSON(w, http.StatusOK, response)
	}

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

func sortChirps(chirps []Chirp, method string) []Chirp {

	sort.Slice(chirps, func(i, j int) bool {

		if method == "asc" {
			return chirps[i].ID < chirps[j].ID
		} else {
			return chirps[i].ID > chirps[j].ID
		}
	})

	return chirps

}
