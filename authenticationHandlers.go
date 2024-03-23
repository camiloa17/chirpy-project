package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/camiloa17/chirpy-project/internal/authentication"
)

// loginUserHandler handles user login authentication
func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type userPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type userResponse struct {
		Email        string `json:"email"`
		ID           int    `json:"id"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	token, err := authentication.GenerateJWTToken(user.ID, "access", cfg.JwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	rToken, err := authentication.GenerateJWTToken(user.ID, "refresh", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, 200, userResponse{Email: user.Email, ID: user.ID, Token: token, RefreshToken: rToken})

}

// refreshTokenHandler
func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type respRefreshToken struct {
		Token string `json:"token"`
	}
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetBearerToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	strUserId, err := authentication.ValidateToken(cleanToken, "refresh", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = cfg.DBRepo.FindRevokedRefreshToken(cleanToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	rToken, err := authentication.GenerateJWTToken(userId, "access", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, respRefreshToken{Token: rToken})

}

// revokeRefreshToken
func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	cleanToken, err := authentication.GetBearerToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = authentication.ValidateToken(cleanToken, "refresh", cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = cfg.DBRepo.RevokeRefreshToken(cleanToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
