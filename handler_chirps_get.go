package main

import (
	"net/http"
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

	var responses []response
	for _, chirp := range chirps {
		responseTmp := response{
			Id:         chirp.ID,
			Created_at: chirp.CreatedAt,
			Updated_at: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_id:    chirp.UserID,
		}
		responses = append(responses, responseTmp)
	}

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
