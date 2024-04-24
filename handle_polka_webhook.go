package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) PolkaHandler(w http.ResponseWriter, r *http.Request) {

	type PaymentEvent struct {
		UserID int `json:"user_id"`
	}
	type Body struct {
		Data  PaymentEvent `json:"data"`
		Event string       `json:"event"`
	}

	decoder := json.NewDecoder(r.Body)

	params := Body{}

	err := decoder.Decode(&params)

	if err != nil {
		respondwithError(w, http.StatusInternalServerError, "couldnt decode")
	}
	if params.Event == "user.upgraded" {
		err = cfg.DB.UpgradeUser(params.Data.UserID)

		if err != nil {
			respondwithError(w, http.StatusNotFound, "user not found")
		}

		w.WriteHeader(http.StatusOK)
	}
}
