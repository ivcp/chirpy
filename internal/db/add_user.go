package db

import (
	"errors"
	"fmt"
)

func (db *DB) AddUser(email string, password string) (User, error) {
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
		Id:           id,
		Email:        email,
		PasswordHash: password,
		IsChirpyRed:  false,
	}
	DBStruct.Users[id] = newUser

	err = db.writeDB(DBStruct)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
