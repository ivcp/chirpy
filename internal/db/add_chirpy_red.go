package db

import (
	"errors"
)

func (db *DB) AddChirpyRed(id int) error {
	DBStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := DBStruct.Users[id]
	if !ok {
		return errors.New("User does not exist")
	}

	user.IsChirpyRed = ok

	DBStruct.Users[id] = user

	err = db.writeDB(DBStruct)
	if err != nil {
		return err
	}

	return nil
}
