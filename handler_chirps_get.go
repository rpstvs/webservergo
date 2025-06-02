package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) RetrieveChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.dbQueries.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, 500, "no chirps", nil)
		return
	}

	authorId := uuid.Nil

	authorIDString := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	if authorIDString != "" {
		authorId, err = uuid.Parse(authorIDString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

	}

	var responses []response
	for _, chirp := range chirps {
		if authorId != uuid.Nil && chirp.UserID != authorId {
			continue
		}
		responseTmp := response{
			Id:         chirp.ID,
			Created_at: chirp.CreatedAt,
			Updated_at: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_id:    chirp.UserID,
		}
		responses = append(responses, responseTmp)
	}

	sort.Slice(responses, func(i, j int) bool {
		if sorting == "desc" {
			return responses[i].Created_at.After(responses[j].Updated_at)
		}
		return responses[i].Created_at.Before(responses[j].Updated_at)
	})

	respondWithJson(w, http.StatusOK, responses)

}

func (cfg *apiConfig) retrieveChirpsId(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	val, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, " invalid id", nil)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpId(r.Context(), val)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "no chirp", nil)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	})
}
