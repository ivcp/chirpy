package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}
	token := authHeader[7:]

	idStr, err := auth.ValidateJwt(token, cfg.jwtSecret, "refresh")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	isRevoked, err := cfg.database.CheckIfTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token is revoked")
		return
	}

	newToken, err := auth.CreateJwt(id, cfg.jwtSecret, "access")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
