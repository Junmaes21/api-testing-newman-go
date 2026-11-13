package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/junmaes21/api-testing-newman-go/internal/api"
)

func main() {
	addr := ":" + envOr("PORT", "8080")
	server := &http.Server{
		Addr:              addr,
		Handler:           api.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("API listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
