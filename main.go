package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiConfig struct {
	fileserverHit int
}

func main() {
	const port = "8080"
	apiConfig := apiConfig{
		fileserverHit: 0,
	}
	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)
	r.Use(middlewareCors)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", healthHandler)
		r.Get("/reset", apiConfig.metricsResetHandler)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Get("/metrics", apiConfig.metricsHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(apiConfig.middlewareMetricsInc)
		workDir, err := os.Getwd()
		if err != nil {
			log.Panicln(err)
		}
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
