package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

func (cfg *appConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.database.GetChirps()
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		w.WriteHeader(500)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
