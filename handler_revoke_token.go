package main

import (
	"net/http"
	"strings"
)

func (cfg *appConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}
	token := authHeader[7:]

	err := cfg.database.RevokeToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
