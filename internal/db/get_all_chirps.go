package db

import (
	"encoding/json"
	"os"
)

func (db *DB) GetChirps() ([]Chirp, error) {
	dbFile, err := os.ReadFile(db.path)
	dbStructure := DBStructure{}
	if err != nil {
		return []Chirp{}, err
	}
	err = json.Unmarshal(dbFile, &dbStructure)
	if err != nil {
		return []Chirp{}, err
	}
	chirps := []Chirp{}

	for _, val := range dbStructure.Chirps {
		chirps = append(chirps, val)
	}

	return chirps, nil
}
