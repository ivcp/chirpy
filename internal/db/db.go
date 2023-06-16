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
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

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

func (db *DB) CreateChirp(body string) (Chirp, error) {
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

	err = db.writeDB(DBStruct)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

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

func (db *DB) AddUser(email string) (User, error) {
	DBStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, val := range DBStruct.Users {
		if val.Email == email {
			return User{}, errors.New(fmt.Sprintf("%s already taken", email))
		}
	}

	id := len(DBStruct.Users) + 1
	newUser := User{
		Id:    id,
		Email: email,
	}
	DBStruct.Users[id] = newUser

	err = db.writeDB(DBStruct)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
