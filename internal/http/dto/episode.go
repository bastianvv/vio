package dto

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type Episode struct {
	ID       int64      `json:"id"`
	SeasonID int64      `json:"season_id"`
	Number   int        `json:"number"`
	TMDBID   *string    `json:"tmdb_id,omitempty"`
	Title    string     `json:"title,omitempty"`
	Overview string     `json:"overview,omitempty"`
	AirDate  *time.Time `json:"air_date,omitempty"`
	Runtime  int        `json:"runtime_min,omitempty"`
	HasStill bool       `json:"has_still"`
}

func NewEpisode(e *domain.Episode, hasStill bool) *Episode {
	if e == nil {
		return nil
	}

	return &Episode{
		ID:       e.ID,
		SeasonID: e.SeasonID,
		Number:   e.Number,
		TMDBID:   e.TMDBID,
		Title:    e.Title,
		Overview: e.Overview,
		AirDate:  e.AirDate,
		Runtime:  e.RuntimeMin,
		HasStill: hasStill, // ← SAME FIX AS SEASONS
	}
}
