package dto

import "github.com/bastianvv/vio/internal/domain"

type Series struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Overview   string `json:"overview,omitempty"`
	PosterPath string `json:"poster,omitempty"`
	Backdrop   string `json:"backdrop,omitempty"`
}

func NewSeries(s *domain.Series) *Series {
	if s == nil {
		return nil
	}

	return &Series{
		ID:         s.ID,
		Title:      s.Title,
		Overview:   s.Overview,
		PosterPath: strPtr(s.PosterPath),
		Backdrop:   strPtr(s.BackdropPath),
	}
}
