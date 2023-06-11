package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	appCfg := &apiConfig{}
	r := chi.NewRouter()

	fsHandler := appCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()

	apiRouter.Get("/healthz", handlerReadiness)

	apiRouter.Get("/metrics", appCfg.handlerHits)

	r.Mount("/api", apiRouter)

	cors := middlewareCors(r)

	server := &http.Server{
		Handler: cors,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
