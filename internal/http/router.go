package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/media"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/scan"
	"github.com/bastianvv/vio/internal/store"
)

func NewRouter(s store.Store, enricher metadata.Enricher) http.Handler {
	r := chi.NewRouter()

	// Initialize split handlers
	scanner := media.NewScanner(s)
	scans := scan.NewRegistry()
	seriesHandler := NewSeriesHandler(s, enricher)
	seasonsHandler := NewSeasonsHandler(s)
	episodesHandler := NewEpisodesHandler(s)
	moviesHandler := NewMoviesHandler(s, enricher)
	librariesHandler := NewLibrariesHandler(s, scanner, scans)
	filesHandler := NewFilesHandler(s)
	subtitlesHandler := NewSubtitlesHandler(s)

	// ---- Libraries ----
	r.Get("/api/libraries", librariesHandler.ListLibraries)
	r.Post("/api/libraries", librariesHandler.CreateLibrary)
	r.Get("/api/libraries/{id}", librariesHandler.GetLibrary)
	r.Put("/api/libraries/{id}", librariesHandler.UpdateLibrary)
	r.Post("/api/libraries/{id}/scan", librariesHandler.ScanLibrary)
	r.Post("/api/libraries/{id}/rescan", librariesHandler.RescanLibrary)

	// ---- Movies ----
	r.Get("/api/movies", moviesHandler.ListMovies)
	r.Get("/api/movies/{id}", moviesHandler.GetMovie)
	r.Get("/api/movies/{id}/files", moviesHandler.ListMediaFiles)
	r.Post("/api/movies/{id}/enrich", moviesHandler.EnrichMovie)

	// ---- Series ----
	r.Get("/api/series", seriesHandler.ListSeries)
	r.Get("/api/series/{id}", seriesHandler.GetSeries)
	r.Get("/api/series/{id}/seasons", seasonsHandler.ListSeasonsBySeries)
	r.Post("/api/series/{id}/enrich", seriesHandler.EnrichSeries)

	// ---- Seasons ----
	r.Get("/api/seasons/{id}", seasonsHandler.GetSeason)
	r.Get("/api/seasons/{id}/episodes", episodesHandler.ListEpisodesBySeason)

	// ---- Episodes ----
	r.Get("/api/episodes/{id}", episodesHandler.GetEpisode)
	r.Get("/api/episodes/{id}/files", episodesHandler.ListEpisodeFiles)

	// ---- Files ----
	r.Get("/api/files/{id}", filesHandler.GetFile)
	r.Get("/api/files/{id}/stream", filesHandler.StreamFile)

	// --- Subtitles ---
	r.Get("/api/files/{id}/subtitles", subtitlesHandler.ListSubtitleTracks)

	// --- Scanner ---
	r.Get("/api/scans/{job_id}", librariesHandler.GetScanJob)

	return r
}
