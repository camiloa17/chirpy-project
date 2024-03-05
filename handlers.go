package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Health handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, http.StatusText(http.StatusOK))
}

// Handler for sending metrics to admins
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	blah := `
	<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
	`
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(blah, cfg.fileserverHit))
}

// Resets the metrics
func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHit = 0
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hits reset to 0")
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
