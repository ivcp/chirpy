package db

import (
	"errors"
	"fmt"
)

func (db *DB) GetUserByEmail(email string) (User, error) {
	DBStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, err := userInDb(DBStruct.Users, email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func userInDb(db map[int]User, email string) (User, error) {
	for _, val := range db {
		if val.Email == email {
			return val, nil
		}
	}
	return User{}, errors.New(fmt.Sprintf("User with email %s does not exist", email))
}
