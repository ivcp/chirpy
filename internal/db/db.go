package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDb(path string) (*DB, error) {
	db := &DB{
		path: path,
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	// load db
	DBStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(DBStruct.Chirps) + 1
	newChirp := Chirp{
		Id:   id,
		Body: body,
	}
	DBStruct.Chirps[id] = newChirp
	// write db
	err = db.writeDB(DBStruct)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// func (db *DB) GetChirps() ([]Chirp, error)

func (db *DB) ensureDB() error {
	_, err := os.Create(db.path)
	if err != nil {
		return err
	}
	fmt.Println("file created")
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

	err = os.WriteFile("database.json", dbJson, 0o666)
	if err != nil {
		return err
	}

	return nil
}
