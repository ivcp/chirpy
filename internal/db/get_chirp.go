package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbFile, err := os.ReadFile(db.path)
	dbStructure := DBStructure{}
	if err != nil {
		return Chirp{}, err
	}
	err = json.Unmarshal(dbFile, &dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New(fmt.Sprintf("Chirp with id %v does not exist", id))
	}

	return chirp, nil
}
