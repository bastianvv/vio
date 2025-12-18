package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/bastianvv/vio/internal/http"
	"github.com/bastianvv/vio/internal/store"
)

func main() {
	dbPath := envOr("VIO_DB", "vio.db")
	addr := envOr("VIO_ADDR", ":8080")

	// Open database
	s, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() { _ = s.Close() }()

	// Ensure schema exists
	if err := s.EnsureSchema(); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	// Initialize router
	r := apphttp.NewRouter(s)

	log.Printf("VIO listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
