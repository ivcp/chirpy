package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

func (cfg *appConfig) handlerAddUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if !isValidEmail(params.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	// add user
	newUser, err := cfg.database.AddUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save user: %s", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, newUser)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
