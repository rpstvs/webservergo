package main

import (
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	var apiCfg apiConfig
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.midlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("GET /admin/metrics", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.resetCounterHits)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	server.ListenAndServe()

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

/*
func (apic *apiConfig) counterHandler(next http.Handler) http.Handler {
	tmp := fmt.Sprintf("Hits: %d", apic.fileserverHits)
	w.Write([]byte(tmp))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
*/

func (apic *apiConfig) resetCounterHits(w http.ResponseWriter, r *http.Request) {
	apic.fileserverHits = 0
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) midlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
