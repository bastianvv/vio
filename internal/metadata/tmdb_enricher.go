package metadata

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/bastianvv/vio/internal/metadata/tmdb"
	"github.com/bastianvv/vio/internal/store"
)

type TMDBEnricher struct {
	store     store.Store
	tmdb      *tmdb.Client
	imageBase string
}

func NewTMDBEnricher(store store.Store, tmdbClient *tmdb.Client, imageBase string) *TMDBEnricher {
	return &TMDBEnricher{
		store:     store,
		tmdb:      tmdbClient,
		imageBase: imageBase,
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

	// Cache poster
	if details.PosterPath != nil {
		local, err := cacheTMDBImage(
			ctx,
			e.imageBase,
			"movies",
			movie.ID,
			"poster",
			*details.PosterPath,
		)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(e.imageBase, local)
		movie.PosterPath = &rel
	}

	// Cache backdrop
	if details.BackdropPath != nil {
		local, err := cacheTMDBImage(
			ctx,
			e.imageBase,
			"movies",
			movie.ID,
			"backdrop",
			*details.BackdropPath,
		)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(e.imageBase, local)
		movie.BackdropPath = &rel
	}

	// 4) Persist
	return e.store.UpdateMovie(movie)
}

func (e *TMDBEnricher) EnrichSeries(ctx context.Context, seriesID int64) error {

	series, err := e.store.GetSeries(seriesID)
	if err != nil {
		return err
	}
	if series == nil {
		return errors.New("series not found")
	}

	// 1) Resolve TMDB ID
	if series.TMDBID == nil {
		results, err := e.tmdb.SearchTV(ctx, series.Title)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			return fmt.Errorf("no tmdb tv match for %q", series.Title)
		}

		best := results[0]
		series.TMDBID = &best.ID
	}

	// 2) Fetch TV details
	details, err := e.tmdb.GetTV(ctx, *series.TMDBID)
	if err != nil {
		return err
	}

	series.OriginalTitle = details.OriginalName
	series.Overview = details.Overview
	series.Status = details.Status

	// 3) Seasons
	for _, s := range details.Seasons {

		season, err := e.store.GetSeasonBySeriesAndNumber(series.ID, s.Number)
		if err != nil {
			return err
		}
		if season == nil {
			continue // created by scanner only
		}

		if season.TMDBID == nil {
			season.TMDBID = &s.ID
		}

		season.Title = s.Title
		season.Overview = s.Overview

		if s.PosterPath != nil {
			_, err := cacheTMDBImage(
				ctx,
				e.imageBase,
				"seasons",
				season.ID,
				"poster",
				*s.PosterPath,
			)
			if err != nil {
				return err
			}
		}

		if err := e.store.UpdateSeason(season); err != nil {
			return err
		}

		// 4) Episodes
		eps, err := e.tmdb.GetSeason(ctx, *series.TMDBID, s.Number)
		if err != nil {
			return err
		}

		for _, te := range eps {
			ep, err := e.store.GetEpisodeBySeasonAndNumber(season.ID, te.Number)
			if err != nil || ep == nil {
				continue
			}

			if ep.TMDBID == nil {
				ep.TMDBID = &te.ID
			}

			ep.Title = te.Title
			ep.Overview = te.Overview
			ep.RuntimeMin = te.RuntimeMin

			if te.StillPath != nil {
				_, err := cacheTMDBImage(
					ctx,
					e.imageBase,
					"episodes",
					ep.ID,
					"still",
					*te.StillPath,
				)
				if err != nil {
					return err
				}
			}

			if err := e.store.UpdateEpisode(ep); err != nil {
				return err
			}
		}
	}

	// Cache poster
	if details.PosterPath != nil {
		local, err := cacheTMDBImage(
			ctx,
			e.imageBase,
			"series",
			series.ID,
			"poster",
			*details.PosterPath,
		)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(e.imageBase, local)
		series.PosterPath = &rel
	}

	// Cache backdrop
	if details.BackdropPath != nil {
		local, err := cacheTMDBImage(
			ctx,
			e.imageBase,
			"series",
			series.ID,
			"backdrop",
			*details.BackdropPath,
		)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(e.imageBase, local)
		series.BackdropPath = &rel
	}

	if err := e.store.UpdateSeries(series); err != nil {
		return err
	}

	return nil
}
