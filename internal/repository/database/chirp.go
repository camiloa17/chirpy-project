package database

import (
	"errors"
	"sort"

	"github.com/camiloa17/chirpy-project/internal/models"
)

// GetChirp return a chirp from the database
func (db *DB) GetChirp(id int) (models.Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}
	chirp, ok := database.Chirps[id]
	if !ok {
		return models.Chirp{}, errors.New("no chirp found")
	}
	return chirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]models.Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return []models.Chirp{}, err
	}
	chirps := []models.Chirp{}
	for _, chirp := range database.Chirps {
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(a, b int) bool {
		return chirps[a].ID < chirps[b].ID
	})
	return chirps, nil

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (models.Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}
	newId := len(database.Chirps) + 1

	chirp := models.Chirp{
		ID:       newId,
		Body:     body,
		AuthorID: authorID,
	}
	database.Chirps[chirp.ID] = chirp
	err = db.writeDB(database)
	if err != nil {
		return models.Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}
	chirp, ok := database.Chirps[id]
	if !ok {
		return errors.New("no chirp found")
	}
	delete(database.Chirps, chirp.ID)
	return nil
}
