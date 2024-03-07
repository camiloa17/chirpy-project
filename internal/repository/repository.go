package repository

import (
	"github.com/camiloa17/chirpy-project/internal/models"
)

type DatabaseRepository interface {
	// CreateChirp creates a new chirp and saves it to disk
	CreateChirp(body string) (models.Chirp, error)
	// GetChirps returns all chirps in the database
	GetChirps() ([]models.Chirp, error)
	// GetChirp returns a Chirp by an id.
	GetChirp(id int) (models.Chirp, error)
	// CreateUser creates a new user
	CreateUser(userEmail, password string) (models.User, error)
	// GetUserByEmail finds a user by its email
	GetUserByEmail(userEmail string) (models.User, error)
	// UpdateUser updates the user information
	UpdateUser(user models.User) (models.User, error)
}
