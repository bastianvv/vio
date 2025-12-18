package dto

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type Season struct {
	ID         int64      `json:"id"`
	SeriesID   int64      `json:"series_id"`
	Number     int        `json:"number"`
	TMDBID     *string    `json:"tmdb_id,omitempty"`
	Title      string     `json:"title,omitempty"`
	Overview   string     `json:"overview,omitempty"`
	PosterPath *string    `json:"poster_path,omitempty"`
	AirDate    *time.Time `json:"air_date,omitempty"`
}

func NewSeason(s *domain.Season) *Season {
	if s == nil {
		return nil
	}

	return &Season{
		ID:         s.ID,
		SeriesID:   s.SeriesID,
		Number:     s.Number,
		TMDBID:     s.TMDBID,
		Title:      s.Title,
		Overview:   s.Overview,
		PosterPath: s.PosterPath,
		AirDate:    s.AirDate,
	}
}
