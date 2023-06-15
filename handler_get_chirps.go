package main

import (
	"net/http"
	"sort"
)

func (cfg *appConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
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
