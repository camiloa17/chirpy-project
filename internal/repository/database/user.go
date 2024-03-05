package database

import (
	"errors"

	"github.com/camiloa17/chirpy-project/internal/models"
)

// CreateUser creates a user in the DB
func (db *DB) CreateUser(userEmail string) (models.User, error) {
	database, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	userExist := false

	for _, user := range database.Users {
		if user.Email == userEmail {
			userExist = true
			break
		}
	}
	if userExist {
		return models.User{}, errors.New("user already created")
	}
	newId := len(database.Users) + 1

	user := models.User{
		ID:    newId,
		Email: userEmail,
	}
	database.Users[newId] = user
	err = db.writeDB(database)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
