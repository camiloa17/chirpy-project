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

var issuerMap = map[string]string{
	"access":  "chirpy-access",
	"refresh": "chirpy-refresh",
}

// GenerateJWTToken it generates a "access" or "refresh" token type
func GenerateJWTToken(userId int, tokeType, hash string) (string, error) {
	issuer, ok := issuerMap[tokeType]
	if !ok {
		return "", errors.New("wrong token type")
	}
	expiresIn := time.Duration(1 * time.Hour)
	if issuer == issuerMap["refresh"] {
		expiresIn = time.Duration(60 * 24 * time.Hour)
	}

	claim := jwt.RegisteredClaims{
		Issuer:    issuer,
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

// ValidateToken it validates a "access" or "refresh" token type
func ValidateToken(token, tokeType, hash string) (string, error) {
	issuerForToken, ok := issuerMap[tokeType]

	if !ok {
		return "", errors.New("wrong token type")
	}

	claim := jwt.RegisteredClaims{}
	validToken, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
		return []byte(hash), nil
	})

	if err != nil {
		return "", err
	}

	tokenIssuer, err := validToken.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if tokenIssuer != issuerForToken {
		return "", errors.New("issuer is not correct")
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
