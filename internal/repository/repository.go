package repository

import (
	"github.com/camiloa17/chirpy-project/internal/models"
)

type DatabaseRepository interface {
	// CreateChirp creates a new chirp and saves it to disk
	CreateChirp(body string, authorID int) (models.Chirp, error)
	// GetChirps returns all chirps in the database
	GetChirps() ([]models.Chirp, error)
	// GetChirp returns a Chirp by an id.
	GetChirp(id int) (models.Chirp, error)
	// DeleteChirp deletes a chirp by id
	DeleteChirp(id int) error
	// CreateUser creates a new user
	CreateUser(userEmail, password string) (models.User, error)
	// GetUserByEmail finds a user by its email
	GetUserByEmail(userEmail string) (models.User, error)
	// GetUserByID finds a user by its id
	GetUserByID(userID int) (models.User, error)
	// UpdateUser updates the user information
	UpdateUser(user models.User) (models.User, error)
	// FindRevokedRefreshToken
	FindRevokedRefreshToken(refreshToken string) error
	// RevokeRefreshToken
	RevokeRefreshToken(refreshToken string) error
}
