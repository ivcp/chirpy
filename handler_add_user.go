package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerAddUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing username or password")
		return
	}

	if !isValidEmail(params.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if !isValidPassword((params.Password)) {
		respondWithError(w, http.StatusBadRequest, "Password must be 6 or more characters long")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	newUser, err := cfg.database.AddUser(params.Email, string(hashedPass))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save user: %s", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		Id:          newUser.Id,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	})
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) bool {
	return len(password) >= 6
}
