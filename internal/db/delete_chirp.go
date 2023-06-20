package db

import (
	"encoding/json"
	"os"
)

func (db *DB) DeleteChirp(id int) error {
	dbFile, err := os.ReadFile(db.path)
	dbStructure := DBStructure{}
	if err != nil {
		return err
	}
	err = json.Unmarshal(dbFile, &dbStructure)
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
