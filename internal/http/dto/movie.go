package dto

import "github.com/bastianvv/vio/internal/domain"

type Movie struct {
	ID            int64   `json:"id"`
	LibraryID     int64   `json:"library_id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title,omitempty"`
	Year          int     `json:"year"`
	TMDBID        *string `json:"tmdb_id,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	RuntimeMin    int     `json:"runtime_min,omitempty"`
	PosterPath    *string `json:"poster_path,omitempty"`
	BackdropPath  *string `json:"backdrop_path,omitempty"`
}

func NewMovie(m *domain.Movie) *Movie {
	return &Movie{
		ID:            m.ID,
		LibraryID:     m.LibraryID,
		Title:         m.Title,
		OriginalTitle: m.OriginalTitle,
		Year:          m.Year,
		TMDBID:        m.TMDBID,
		Overview:      m.Overview,
		RuntimeMin:    m.RuntimeMin,
		PosterPath:    m.PosterPath,
		BackdropPath:  m.BackdropPath,
	}
}
