package repository

import (
	"github.com/camiloa17/chirpy-project/internal/models"
)

type DatabaseRepository interface {
	// CreateChirp creates a new chirp and saves it to disk
	CreateChirp(body string) (models.Chirp, error)
	// GetChirps returns all chirps in the database
	GetChirps() ([]models.Chirp, error)
}
