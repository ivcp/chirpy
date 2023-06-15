package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ivcp/chirpy/internal/db"
)

type appConfig struct {
	fileserverHits int
	database       *db.DB
}

func main() {
	database, err := db.NewDb("database.json")
	if err != nil {
		log.Fatal(err)
	}
	const filepathRoot = "."
	const port = "8080"
	appCfg := &appConfig{
		database: database,
	}

	r := chi.NewRouter()

	fsHandler := appCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/chirps", appCfg.handlerAddChirp)

	adminRouter.Get("/metrics", appCfg.handlerHits)

	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	cors := middlewareCors(r)

	server := &http.Server{
		Handler: cors,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
