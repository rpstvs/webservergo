package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rpstvs/webservergo/internals/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	secret         string
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	db, err := database.NewDB("database.json")

	if err != nil {
		log.Fatal(err)
	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		secret:         jwtSecret,
	}
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.midlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.counterHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetCounterHits)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpsid}", apiCfg.retrieveChirpsId)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	//mux.HandleFunc("POST /api/revoke")

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	server.ListenAndServe()

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
