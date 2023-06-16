package db

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
