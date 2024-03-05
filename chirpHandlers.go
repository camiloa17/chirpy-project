package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) getAChirpHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "no id sent")
	}

	chirp, err := cfg.DBRepo.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DBRepo.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fail to read the chirps")
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) addChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}
	w.Header().Set("Content-Type", "application/json")
	message := body{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode message")
		return

	}
	isValid := len(message.Body) <= 140

	if !isValid {
		respondWithError(w, http.StatusBadRequest, "Chirp is to long")
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanMessage := hideNegativeWords(message.Body, badWords)

	chirp, err := cfg.DBRepo.CreateChirp(cleanMessage)

	if err != nil {
		fmt.Printf("error saving chirp, %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not save chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}
