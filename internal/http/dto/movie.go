package dto

import "github.com/bastianvv/vio/internal/domain"

type Movie struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Year         int     `json:"year"`
	Overview     string  `json:"overview,omitempty"`
	PosterPath   *string `json:"poster_path,omitempty"`
	BackdropPath *string `json:"backdrop_path,omitempty"`
}

func NewMovie(m *domain.Movie) *Movie {
	return &Movie{
		ID:           m.ID,
		Title:        m.Title,
		Year:         m.Year,
		Overview:     m.Overview,
		PosterPath:   m.PosterPath,
		BackdropPath: m.BackdropPath,
	}
}
