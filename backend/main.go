// Command server runs the calculator HTTP API.
package main

import (
	"log"
	"net/http"
	"os"

	"calculator-go/backend/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	staticDir := os.Getenv("STATIC_DIR")

	addr := ":" + port
	log.Printf("calculator API listening on %s", addr)
	if staticDir != "" {
		log.Printf("serving frontend from %s", staticDir)
	}
	log.Fatal(http.ListenAndServe(addr, newHandler(staticDir)))
}

// newHandler returns the API routes and, when staticDir is set, also serves the
// built frontend from that directory. The Docker image uses this to run both apps
// from one process on one port; in local development the Vite dev server does it.
func newHandler(staticDir string) http.Handler {
	apiRoutes := api.NewRouter()
	if staticDir == "" {
		return apiRoutes
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", apiRoutes)
	mux.Handle("/health", apiRoutes)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}
