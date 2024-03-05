package main

import (
	"fmt"
	"io"
	"net/http"
)

// Handler for sending metrics to admins
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	blah := `
	<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
	`
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(blah, cfg.fileserverHit))
}

// Resets the metrics
func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHit = 0
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hits reset to 0")
}
