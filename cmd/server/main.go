package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/bastianvv/vio/internal/http"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/metadata/tmdb"
	"github.com/bastianvv/vio/internal/store"
)

func main() {
	dbPath := envOr("VIO_DB", "vio.db")
	addr := envOr("VIO_ADDR", ":8080")
	tmdbKey := "487de335e804022ee3538ab055bd2b53" //os.Getenv("TMDB_API_KEY")

	if tmdbKey == "" {
		log.Fatal("TMDB_API_KEY is required")
	}

	// Open database
	s, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() { _ = s.Close() }()

	if err := s.EnsureSchema(); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	// Metadata
	tmdbClient := tmdb.New(tmdbKey)
	enricher := metadata.NewTMDBEnricher(s, tmdbClient)

	// Router
	r := apphttp.NewRouter(s, enricher)

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
