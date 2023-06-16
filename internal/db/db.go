package db

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

func NewDb(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	_, err := os.ReadFile(path)

	if errors.Is(err, fs.ErrNotExist) {
		err = db.ensureDB()
		if err != nil {
			return db, err
		}
	}
	if err != nil {
		return db, err
	}

	return db, nil
}

func (db *DB) ensureDB() error {
	emptyDb := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}
	err := db.writeDB(emptyDb)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	dbFile, err := os.ReadFile(db.path)
	dbStructure := DBStructure{}
	if err != nil {
		return dbStructure, err
	}
	err = json.Unmarshal(dbFile, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	dbJson, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	db.mux.Lock()
	err = os.WriteFile(db.path, dbJson, 0o666)
	if err != nil {
		return err
	}
	db.mux.Unlock()
	return nil
}
