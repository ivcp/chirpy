package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *appConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")

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

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	chirps, err := cfg.database.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
