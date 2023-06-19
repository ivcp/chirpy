package db

import "time"

func (db *DB) RevokeToken(token string) error {
	DBStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	DBStruct.RevokedTokens[token] = RevokedToken{
		Time:  time.Now().String(),
		Token: token,
	}

	err = db.writeDB(DBStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckIfTokenRevoked(token string) (bool, error) {
	DBStruct, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := DBStruct.RevokedTokens[token]
	if !ok {
		return false, nil
	}

	return true, nil
}
