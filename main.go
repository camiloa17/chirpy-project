package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/camiloa17/chirpy-project/internal/repository"
	"github.com/camiloa17/chirpy-project/internal/repository/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHit int
	DBRepo        repository.DatabaseRepository
	JwtSecret     string
}

func main() {
	const port = "8080"
	godotenv.Load()

	workDir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	db := database.NewDB(filepath.Join(workDir, "database.json"))

	apiConfig := apiConfig{
		fileserverHit: 0,
		DBRepo:        db,
		JwtSecret:     os.Getenv("JWT_SECRET"),
	}

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)
	r.Use(middlewareCors)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", healthHandler)
		r.Get("/reset", apiConfig.metricsResetHandler)
		r.Get("/chirps", apiConfig.getChirpsHandler)
		r.Get("/chirps/{chirpID}", apiConfig.getAChirpHandler)
		r.Delete("/chirps/{chirpID}", apiConfig.deleteChirpsHandler)
		r.Group(func(r chi.Router) {
			// r.Use(middleware.AllowContentType("application/json"))
			r.Post("/chirps", apiConfig.addChirpsHandler)
			r.Post("/users", apiConfig.createUserHandler)
			r.Put("/users", apiConfig.updateUserHandler)
			r.Post("/login", apiConfig.loginUserHandler)
			r.Post("/refresh", apiConfig.refreshTokenHandler)
			r.Post("/revoke", apiConfig.revokeRefreshToken)
			r.Post("/polka/webhooks", apiConfig.polkaPaymentEventHandler)
		})

	})
	r.Route("/admin", func(r chi.Router) {
		r.Get("/metrics", apiConfig.metricsHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(apiConfig.middlewareMetricsInc)
		fileServeDict := http.Dir(filepath.Join(workDir, "public"))
		r.Get("/app", func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(fileServeDict))
			fs.ServeHTTP(w, r)
		})
		r.Get("/app/*", func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(fileServeDict))
			fs.ServeHTTP(w, r)
		})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	if statusCode > 499 {
		fmt.Printf("responding with 5xx error %s\n", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respBody := errResponse{
		Error: msg,
	}
	respondWithJSON(w, statusCode, respBody)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("error marshalling JSON: %s\n", err)
	}
	if statusCode != 200 {
		w.WriteHeader(statusCode)
	}
	w.Write(dat)
}
