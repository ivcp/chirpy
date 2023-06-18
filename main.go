package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/ivcp/chirpy/internal/db"
)

type appConfig struct {
	fileserverHits int
	database       *db.DB
}

func main() {
	debug()
	const filepathRoot = "."
	const port = "8080"

	database, err := db.NewDb("database.json")
	if err != nil {
		log.Fatal(err)
	}
	appCfg := &appConfig{
		database: database,
	}

	r := chi.NewRouter()

	fsHandler := appCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)

	apiRouter.Get("/chirps", appCfg.handlerGetChirps)
	apiRouter.Get("/chirps/{chirpId}", appCfg.handlerGetOneChirp)
	apiRouter.Post("/chirps", appCfg.handlerAddChirp)

	apiRouter.Post("/users", appCfg.handlerAddUser)
	apiRouter.Post("/login", appCfg.handlerLogin)

	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", appCfg.handlerHits)
	r.Mount("/admin", adminRouter)

	cors := middlewareCors(r)

	server := &http.Server{
		Handler: cors,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func debug() {
	debug := flag.Bool("debug", false, "remove db")
	flag.Parse()

	if *debug {
		os.Remove("database.json")
	}
}
