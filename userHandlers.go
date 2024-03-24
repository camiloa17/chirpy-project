package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/camiloa17/chirpy-project/internal/authentication"
	"github.com/camiloa17/chirpy-project/internal/models"
)

type userResponse struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// createUserHandler handles the user creation request
func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	w.Header().Set("Content-Type", "application/json")
	userPl := userPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userPl)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode payload")
		return
	}

	_, err = mail.ParseAddress(userPl.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(userPl.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "no password sent")
		return
	}

	hashedPass, err := authentication.HashPassword(userPl.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUser, err := cfg.DBRepo.CreateUser(strings.ToLower(userPl.Email), string(hashedPass))

	if err != nil && err.Error() == "user already created" {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userResponse{Email: newUser.Email, ID: newUser.ID, IsChirpyRed: newUser.IsChirpyRed})

}

// updateUserHandler handles user updates
func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type userPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetAuthToken(token, "Bearer")

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	strUserId, err := authentication.ValidateToken(cleanToken, "access", cfg.JwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DBRepo.GetUserByID(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userPl := userPayload{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&userPl)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode payload")
		return
	}
	_, err = mail.ParseAddress(userPl.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(userPl.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "no password sent")
		return
	}

	hashedPass, err := authentication.HashPassword(userPl.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userUpdate := models.User{
		ID:          userId,
		Password:    hashedPass,
		Email:       userPl.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	uUser, err := cfg.DBRepo.UpdateUser(userUpdate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, userResponse{Email: uUser.Email, ID: uUser.ID, IsChirpyRed: uUser.IsChirpyRed})
}
