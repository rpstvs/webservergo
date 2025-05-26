package main

import (
	"net/http"
	"sync/atomic"

	"github.com/rpstvs/webservergo/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	Platform       string
	tokenSecret    string
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := GetConfig()

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.midlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.RetrieveChirps)
	mux.HandleFunc("GET /api/chirps/{chirpsid}", apiCfg.retrieveChirpsId)
	mux.HandleFunc("DELETE /api/chirps/{chirpsid}", apiCfg.DeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.PolkaHandler)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUser)

	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.revoke)

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.counterHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	server.ListenAndServe()

}
