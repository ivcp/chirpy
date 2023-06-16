package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *appConfig) handlerAddChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
		return
	}

	const maxChirpLength = 140
	if len([]rune(params.Body)) > maxChirpLength {

		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanBody(params.Body)

	// add to db

	newChirp, err := cfg.database.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed so save chirp to DB")
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}

func cleanBody(body string) string {
	bannedWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

	bodySlice := strings.Split(body, " ")

	for i, word := range bodySlice {
		if _, ok := bannedWords[strings.ToLower(word)]; ok {
			bodySlice[i] = "****"
		}
	}

	return strings.Join(bodySlice, " ")
}
