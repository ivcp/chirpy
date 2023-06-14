package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerChirpValidator(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorT struct {
		Error string `json:"error"`
	}
	type resp struct {
		Id   int    `json:"id"`
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	w.Header().Set("Content-Type", "application/json")
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

	respondWithJSON(w, http.StatusCreated, resp{
		Body: cleanedBody,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
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
