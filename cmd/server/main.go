package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bastianvv/vio/internal/domain"
	apphttp "github.com/bastianvv/vio/internal/http"
	"github.com/bastianvv/vio/internal/media"
	"github.com/bastianvv/vio/internal/store"
)

func main() {
	dbPath := envOr("VIO_DB", "vio.db")
	addr := envOr("VIO_ADDR", ":8080")

	s, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() { _ = s.Close() }()

	// --------------------------------------------------------
	// 1. CREATE A TEST LIBRARY ENTRY (ONLY IF NOT EXISTS)
	// --------------------------------------------------------

	path := "/home/bastianv/NAS/Media/vio-test/tv"

	lib := &domain.Library{
		Name: "TV Shows",
		Type: domain.LibraryTypeSeries,
		Path: path,
	}
	if err := s.CreateLibrary(lib); err != nil {
		log.Fatalf("failed to create library: %v", err)
	}

	scanner := media.NewScanner(s)
	log.Println("Scanning library...")

	result, err := scanner.ScanLibrary(lib, media.ScanModeIncremental)
	if err != nil {
		log.Fatalf("scan failed: %v", err)
	}

	log.Printf(
		"Scan done: files=%d movies=%d series=%d episodes=%d errors=%d",
		result.FilesScanned,
		result.MoviesAdded,
		result.SeriesAdded,
		result.EpisodesAdded,
		len(result.Errors),
	)

	// --------------------------------------------------------
	// 3. START HTTP SERVER (optional for now)
	// --------------------------------------------------------
	r := apphttp.NewRouter(s)

	log.Printf("listening on %s", addr)
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
