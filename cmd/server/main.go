package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	apphttp "github.com/bastianvv/vio/internal/http"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/metadata/tmdb"
	"github.com/bastianvv/vio/internal/store"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	dbPath := envOr("VIO_DB", "vio.db")
	addr := envOr("VIO_ADDR", ":8080")
	tmdbKey := os.Getenv("TMDB_API_KEY")

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
	imageBasePath := os.Getenv("IMAGE_CACHE_DIR")
	absImagePath, err := filepath.Abs(imageBasePath)
	if err != nil {
		panic(err)
	}

	tmdbClient := tmdb.New(tmdbKey)
	enricher := metadata.NewTMDBEnricher(s, tmdbClient, absImagePath)

	if err := os.MkdirAll(imageBasePath, 0755); err != nil {
		log.Fatalf("failed to create image cache dir: %v", err)
	}

	// Router
	r := apphttp.NewRouter(s, enricher, absImagePath)

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
