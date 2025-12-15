package dto

import "github.com/bastianvv/vio/internal/domain"

type Season struct {
	ID         int64   `json:"id"`
	Number     int     `json:"number"`
	Title      string  `json:"title,omitempty"`
	Overview   string  `json:"overview,omitempty"`
	PosterPath *string `json:"poster_path,omitempty"`
}

func NewSeason(s *domain.Season) *Season {
	if s == nil {
		return nil
	}

	return &Season{
		ID:         s.ID,
		Number:     s.Number,
		Title:      s.Title,
		Overview:   s.Overview,
		PosterPath: s.PosterPath,
	}
}
