package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/camiloa17/chirpy-project/internal/authentication"
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
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetAuthToken(token, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	strUserId, err := authentication.ValidateToken(cleanToken, "access", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type body struct {
		Body string `json:"body"`
	}
	w.Header().Set("Content-Type", "application/json")
	message := body{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&message)

	defer r.Body.Close()

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

	chirp, err := cfg.DBRepo.CreateChirp(cleanMessage, userId)

	if err != nil {
		fmt.Printf("error saving chirp, %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not save chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) deleteChirpsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetAuthToken(token, "Bearer")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	strUserId, err := authentication.ValidateToken(cleanToken, "access", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.DBRepo.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if chirp.AuthorID != userId {
		respondWithError(w, http.StatusForbidden, "")
		return
	}

	err = cfg.DBRepo.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}

func hideNegativeWords(text string, negativeWords map[string]struct{}) string {
	bodyWords := strings.Fields(text)
	for idx, word := range bodyWords {
		lowerCase := strings.ToLower(word)
		_, ok := negativeWords[lowerCase]
		if ok {
			bodyWords[idx] = "****"
		}
	}
	return strings.Join(bodyWords, " ")
}
