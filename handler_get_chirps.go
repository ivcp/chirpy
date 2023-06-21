package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/ivcp/chirpy/internal/db"
)

func (cfg *appConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")
	sortParam := req.URL.Query().Get("sort")
	sortDirection := "asc"

	if sortParam != "" && sortParam != "asc" && sortParam != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort input")
		return
	}

	if sortParam == "desc" {
		sortDirection = "desc"
	}

	if authorId != "" {
		id, err := strconv.Atoi(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Invalid author ID")
			return
		}
		chirps, err := cfg.database.GetChirpsByAuthor(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
			return
		}

		sortChirps(chirps, sortDirection)

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	chirps, err := cfg.database.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
		return
	}

	sortChirps(chirps, sortDirection)

	respondWithJSON(w, http.StatusOK, chirps)
}

func sortChirps(chirps []db.Chirp, dir string) {
	sort.Slice(chirps, func(i, j int) bool {
		if dir == "desc" {
			return chirps[i].Id > chirps[j].Id
		}
		return chirps[i].Id < chirps[j].Id
	})
}
