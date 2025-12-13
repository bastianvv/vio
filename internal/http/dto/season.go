package dto

import "github.com/bastianvv/vio/internal/domain"

type Season struct {
	ID       int64  `json:"id"`
	Number   int    `json:"number"`
	Title    string `json:"title,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Overview string `json:"overview,omitempty"`
}

func NewSeason(s *domain.Season) *Season {
	if s == nil {
		return nil
	}

	return &Season{
		ID:       s.ID,
		Number:   s.Number,
		Title:    s.Title,
		Poster:   strPtr(s.PosterPath),
		Overview: s.Overview,
	}
}
