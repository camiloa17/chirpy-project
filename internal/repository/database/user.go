package database

import (
	"errors"

	"github.com/camiloa17/chirpy-project/internal/models"
)

// CreateUser creates a user in the DB
func (db *DB) CreateUser(userEmail string, password string) (models.User, error) {
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
		ID:       newId,
		Email:    userEmail,
		Password: password,
	}
	database.Users[newId] = user
	err = db.writeDB(database)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// GetUserByEmail finds a user by its email
func (db *DB) GetUserByEmail(userEmail string) (models.User, error) {
	database, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	foundUser := models.User{}
	found := false
	for _, user := range database.Users {
		if user.Email == userEmail {
			foundUser = user
			found = true
			break
		}
	}

	if !found {
		return foundUser, errors.New("no user found")
	}

	return foundUser, nil
}

func (db *DB) UpdateUser(user models.User) (models.User, error) {
	database, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	userExist := false

	for _, savedUser := range database.Users {
		if savedUser.Email == user.Email && savedUser.ID != user.ID {
			userExist = true
			break
		}
	}

	if userExist {
		return models.User{}, errors.New("email already in use")
	}
	database.Users[user.ID] = user
	err = db.writeDB(database)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
