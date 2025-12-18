package dto

import "github.com/bastianvv/vio/internal/domain"

type Series struct {
	ID            int64   `json:"id"`
	LibraryID     int64   `json:"library_id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title,omitempty"`
	TMDBID        *string `json:"tmdb_id,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	Status        string  `json:"status,omitempty"`
	PosterPath    *string `json:"poster_path,omitempty"`
	BackdropPath  *string `json:"backdrop_path,omitempty"`
}

func NewSeries(s *domain.Series) *Series {
	if s == nil {
		return nil
	}

	return &Series{
		ID:            s.ID,
		LibraryID:     s.LibraryID,
		Title:         s.Title,
		OriginalTitle: s.OriginalTitle,
		TMDBID:        s.TMDBID,
		Overview:      s.Overview,
		Status:        s.Status,
		PosterPath:    s.PosterPath,
		BackdropPath:  s.BackdropPath,
	}
}
