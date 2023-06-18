package main

import (
	"encoding/json"
	"net/http"

	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
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
	user, err := cfg.database.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:    user.Id,
		Email: user.Email,
	})
}
