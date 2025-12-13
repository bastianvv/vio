package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/media"
	"github.com/bastianvv/vio/internal/store"
)

func NewRouter(s store.Store) http.Handler {
	r := chi.NewRouter()

	// Initialize split handlers
	scanner := media.NewScanner(s)
	seriesHandler := NewSeriesHandler(s)
	seasonsHandler := NewSeasonsHandler(s)
	episodesHandler := NewEpisodesHandler(s)
	moviesHandler := NewMoviesHandler(s)
	librariesHandler := NewLibrariesHandler(s, scanner)
	filesHandler := NewFilesHandler(s)

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

	// ---- Series ----
	r.Get("/api/series", seriesHandler.ListSeries)
	r.Get("/api/series/{id}", seriesHandler.GetSeries)
	r.Get("/api/series/{id}/seasons", seasonsHandler.ListSeasonsBySeries)

	// ---- Seasons ----
	r.Get("/api/seasons/{id}", seasonsHandler.GetSeason)
	r.Get("/api/seasons/{id}/episodes", episodesHandler.ListEpisodesBySeason)

	// ---- Episodes ----
	r.Get("/api/episodes/{id}", episodesHandler.GetEpisode)

	// ---- Files ----
	r.Get("/api/files/{id}", filesHandler.GetFile)
	r.Get("/api/files/{id}/stream", filesHandler.StreamFile)

	return r
}
