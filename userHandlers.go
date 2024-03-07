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
	Email string `json:"email"`
	ID    int    `json:"id"`
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

	respondWithJSON(w, http.StatusCreated, userResponse{Email: newUser.Email, ID: newUser.ID})

}

// updateUserHandler handles user updates
func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type userPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetBearerToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	strUserId, err := authentication.ValidateToken(cleanToken, cfg.JwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	userId, err := strconv.Atoi(strUserId)
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
		ID:       userId,
		Password: hashedPass,
		Email:    userPl.Email,
	}

	user, err := cfg.DBRepo.UpdateUser(userUpdate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, userResponse{Email: user.Email, ID: user.ID})
}

// loginUserHandler handles user login authentication
func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type userPayload struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}
	type userResponse struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
		Token string `json:"token"`
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

	user, err := cfg.DBRepo.GetUserByEmail(strings.ToLower(userPl.Email))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = authentication.ValidatePassword(userPl.Password, user.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := authentication.GenerateJWTToken(user.ID, userPl.ExpiresInSeconds, cfg.JwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, 200, userResponse{Email: user.Email, ID: user.ID, Token: token})

}
