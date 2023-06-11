package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	appCfg := &apiConfig{}
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()

	mux.Handle("/app/", appCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("/healthz", handlerReadiness)

	mux.HandleFunc("/metrics", appCfg.handlerHits)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}
