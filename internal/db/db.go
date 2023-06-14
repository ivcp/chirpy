package db

import (
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
	Body string `json:"cleaned_body"`
}

func NewDb(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	db := &DB{}

	if errors.Is(err, fs.ErrNotExist) {
		err = db.ensureDB(path)
		if err != nil {
			return db, err
		}
	} else {
		return db, nil
	}

	return db, nil
}

// func (db *DB) CreateChirp(body string) (Chirp, error)

// func (db *DB) GetChirps() ([]Chirp, error)

func (db *DB) ensureDB(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	fmt.Println("file created")
	return nil
}

// func (db *DB) loadDB() (DBStructure, error)

// func (db *DB) writeDB(dbStructure DBStructure) error
