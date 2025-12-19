package dto

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type Season struct {
	ID        int64      `json:"id"`
	SeriesID  int64      `json:"series_id"`
	Number    int        `json:"number"`
	TMDBID    *string    `json:"tmdb_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Overview  string     `json:"overview,omitempty"`
	AirDate   *time.Time `json:"air_date,omitempty"`
	HasPoster bool       `json:"has_poster"`
}

func NewSeason(s *domain.Season, hasPoster bool) *Season {
	if s == nil {
		return nil
	}

	return &Season{
		ID:        s.ID,
		SeriesID:  s.SeriesID,
		Number:    s.Number,
		TMDBID:    s.TMDBID,
		Title:     s.Title,
		Overview:  s.Overview,
		AirDate:   s.AirDate,
		HasPoster: hasPoster,
	}
}
