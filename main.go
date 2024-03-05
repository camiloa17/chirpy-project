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
)

type apiConfig struct {
	fileserverHit int
	DBRepo        repository.DatabaseRepository
}

func main() {
	const port = "8080"
	workDir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	db := database.NewDB(filepath.Join(workDir, "database.json"))

	apiConfig := apiConfig{
		fileserverHit: 0,
		DBRepo:        db,
	}
	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)
	r.Use(middlewareCors)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", healthHandler)
		r.Get("/reset", apiConfig.metricsResetHandler)
		r.Get("/chirps", apiConfig.getChirpsHandler)
		r.Group(func(r chi.Router) {
			// r.Use(middleware.AllowContentType("application/json"))
			r.Post("/chirps", apiConfig.addChirpsHandler)
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
	w.WriteHeader(statusCode)
	w.Write(dat)
}

func hideNegativeWords(text string, negativeWords map[string]struct{}) string {
	bodyWords := strings.Fields(text)
	for idx, word := range bodyWords {
		lowerCase := strings.ToLower(word)
		_, ok := negativeWords[lowerCase]
		if ok {
			bodyWords[idx] = "****"
		}
	}
	return strings.Join(bodyWords, " ")
}
