package database

import (
	"errors"
	"time"
)

func (db *DB) FindRevokedRefreshToken(refreshToken string) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := database.RevokedTokens[refreshToken]
	if ok {
		return errors.New("token revoked")
	}
	return nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	database, err := db.loadDB()
	if err != nil {
		return err
	}

	database.RevokedTokens[refreshToken] = time.Now().UTC()
	err = db.writeDB(database)
	if err != nil {
		return err
	}

	return nil
}
