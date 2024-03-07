package authentication

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the users password to be saved on the db.
func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 2)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

// ValidatePassword validates the password of the user login in.
func ValidatePassword(password, hash string) error {

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWTToken
func GenerateJWTToken(userId int, expiresInSeconds int, hash string) (string, error) {
	expiresInHours := time.Duration(24 * time.Hour)
	expiresIn := time.Duration(24 * time.Hour)

	if expiresInSeconds != 0 && expiresInSeconds < int(expiresInHours.Seconds()) {
		expiresIn = time.Duration(expiresInSeconds) * time.Second
	}

	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprint(userId),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := jwtToken.SignedString([]byte(hash))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken
func ValidateToken(token, hash string) (string, error) {
	claim := jwt.RegisteredClaims{}
	validToken, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
		return []byte(hash), nil
	})
	if err != nil {
		return "", err
	}

	strUserId, err := validToken.Claims.GetSubject()

	if err != nil {
		return "", err
	}

	return strUserId, nil
}

// GetBearerToken gets the auth token from the Authorization header
func GetBearerToken(headerToken string) (string, error) {
	tokenPrefix := "Bearer"
	cleanBearerToken, hasPrefix := strings.CutPrefix(headerToken, tokenPrefix)
	if !hasPrefix {
		return "", errors.New("no Bearer token provides")
	}
	return strings.TrimSpace(cleanBearerToken), nil
}
