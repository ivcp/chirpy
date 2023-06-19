package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}
	token := authHeader[7:]

	if !isValidEmail(params.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if !isValidPassword((params.Password)) {
		respondWithError(w, http.StatusBadRequest, "Password must be 6 or more characters long")
		return
	}

	idStr, err := auth.ValidateJwt(token, cfg.jwtSecret, "access")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.database.UpdateUser(id, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Id:    user.Id,
		Email: user.Email,
	})
}
