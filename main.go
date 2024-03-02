package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const publicFilePath = "./public"
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(publicFilePath))))

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Fatal(srv.ListenAndServe())
}
