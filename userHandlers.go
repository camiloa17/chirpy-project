package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}
	w.Header().Set("Content-Type", "application/json")
	email := userEmail{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode payload")
		return
	}

	_, err = mail.ParseAddress(email.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, err := cfg.DBRepo.CreateUser(strings.ToLower(email.Email))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, newUser)

}
