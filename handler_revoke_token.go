package main

import (
	"net/http"

	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.database.RevokeToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
