package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/camiloa17/chirpy-project/internal/models"
	"github.com/camiloa17/chirpy-project/internal/repository"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(dbPath string) repository.DatabaseRepository {

	db := DB{
		path: dbPath,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		fmt.Println(err)
	}
	return &db
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	data, err := os.ReadFile(db.path)
	chirps := models.Chirps{Chirps: make(map[int]models.Chirp)}
	if errors.Is(err, os.ErrNotExist) {
		err := db.writeDB(chirps)
		if err != nil {
			return err
		}
	}
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}
	if len(data) == 0 {
		db.writeDB(chirps)
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (models.Chirps, error) {
	chirps := models.Chirps{Chirps: make(map[int]models.Chirp)}
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil || len(data) == 0 {
		err := db.ensureDB()
		if err != nil {
			return chirps, err
		}
		data, err = os.ReadFile(db.path)
		if err != nil {
			return chirps, err
		}
	}

	err = json.Unmarshal(data, &chirps)
	if err != nil {
		fmt.Println(err)
		return chirps, err
	}
	return chirps, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure models.Chirps) error {
	val, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, val, 0600)
	if err != nil {
		return err
	}
	return nil
}
