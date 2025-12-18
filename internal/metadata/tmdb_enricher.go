package metadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/bastianvv/vio/internal/metadata/tmdb"
	"github.com/bastianvv/vio/internal/store"
)

type TMDBEnricher struct {
	store store.Store
	tmdb  *tmdb.Client
}

func NewTMDBEnricher(store store.Store, tmdbClient *tmdb.Client) *TMDBEnricher {
	return &TMDBEnricher{
		store: store,
		tmdb:  tmdbClient,
	}
}

func (e *TMDBEnricher) EnrichMovie(ctx context.Context, movieID int64) error {
	movie, err := e.store.GetMovie(movieID)
	if err != nil {
		return err
	}
	if movie == nil {
		return errors.New("movie not found")
	}

	// 1) Search TMDB (basic results)
	results, err := e.tmdb.SearchMovie(ctx, movie.Title, movie.Year)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("no tmdb match for %q (%d)", movie.Title, movie.Year)
	}

	// MVP heuristic: first result.
	best := results[0]

	// 2) Fetch full details
	details, err := e.tmdb.GetMovie(ctx, best.ID)
	if err != nil {
		return err
	}
	if details == nil {
		return errors.New("tmdb details not found")
	}

	// 3) Apply metadata to domain model
	// Only overwrite fields we actually want to enrich.
	movie.TMDBID = &details.ID

	if details.OriginalTitle != "" {
		movie.OriginalTitle = details.OriginalTitle
	}
	if details.Overview != "" {
		movie.Overview = details.Overview
	}
	if details.RuntimeMin > 0 {
		movie.RuntimeMin = details.RuntimeMin
	}

	// TMDB returns poster/backdrop paths as strings like "/abc.jpg"
	// Your domain uses *string, so assign directly.
	movie.PosterPath = details.PosterPath
	movie.BackdropPath = details.BackdropPath

	// Optional: if year was missing locally, you can fill it.
	if movie.Year == 0 && details.ReleaseYear > 0 {
		movie.Year = details.ReleaseYear
	}

	// 4) Persist
	return e.store.UpdateMovie(movie)
}
