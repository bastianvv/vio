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
	HasPoster     bool    `json:"has_poster"`
	HasBackdrop   bool    `json:"has_backdrop"`
}

func NewSeries(s *domain.Series, hasPoster, hasBackdrop bool) *Series {
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
		HasPoster:     hasPoster,
		HasBackdrop:   hasBackdrop,
	}
}
