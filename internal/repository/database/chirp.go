package database

import (
	"sort"

	"github.com/camiloa17/chirpy-project/internal/models"
)

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]models.Chirp, error) {
	dbChirps, err := db.loadDB()
	if err != nil {
		return []models.Chirp{}, err
	}
	chirps := []models.Chirp{}
	for _, chirp := range dbChirps.Chirps {
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(a, b int) bool {
		return chirps[a].ID < chirps[b].ID
	})
	return chirps, nil

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (models.Chirp, error) {
	chirps, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}
	newId := len(chirps.Chirps) + 1

	chirp := models.Chirp{
		ID:   newId,
		Body: body,
	}
	chirps.Chirps[chirp.ID] = chirp
	err = db.writeDB(chirps)
	if err != nil {
		return models.Chirp{}, err
	}

	return chirp, nil
}
