package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/camiloa17/chirpy-project/internal/authentication"
)

type polkaPayload struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaPaymentEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetAuthToken(token, "ApiKey")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	apiKey := os.Getenv("POLKA_API_KEY")
	if apiKey != cleanToken {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	polkaEvent := polkaPayload{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&polkaEvent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode message")
		return
	}

	if polkaEvent.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "")
		return
	}
	user, err := cfg.DBRepo.GetUserByID(polkaEvent.Data.UserId)
	if err != nil && err.Error() == "no user found" {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.IsChirpyRed = true
	_, err = cfg.DBRepo.UpdateUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
