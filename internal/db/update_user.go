package db

import (
	"errors"

	"github.com/ivcp/chirpy/internal/auth"
)

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	DBStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := DBStruct.Users[id]
	if !ok {
		return User{}, errors.New("User does not exist")
	}

	pass, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	updatedUser := User{
		Id:           user.Id,
		Email:        email,
		PasswordHash: string(pass),
	}

	DBStruct.Users[id] = updatedUser

	err = db.writeDB(DBStruct)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}
