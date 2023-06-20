package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ivcp/chirpy/internal/auth"
)

func (cfg *appConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	chirpIdStr := chi.URLParam(req, "chirpId")

	chirpId, err := strconv.Atoi(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}
	token := authHeader[7:]

	idStr, err := auth.ValidateJwt(token, cfg.jwtSecret, "access")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.database.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if userId != chirp.AuthorId {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}

	err = cfg.database.DeleteChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
