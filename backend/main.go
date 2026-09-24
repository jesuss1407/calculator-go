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

	addr := ":" + port
	log.Printf("calculator API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, api.NewRouter()))
}
