package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rpstvs/webservergo/internal/auth"
)

func (cfg *apiConfig) PolkaHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	headerKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "malformed key", nil)
		return
	}

	if headerKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key on header", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode params", nil)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Printf("User to be upgraded %s \n", params.Data.UserID.String())
	_, err = cfg.dbQueries.UpgradeUser(r.Context(), params.Data.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldnt update user", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
