package dto

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type Episode struct {
	ID      int64      `json:"id"`
	Number  int        `json:"number"`
	Title   string     `json:"title,omitempty"`
	AirDate *time.Time `json:"air_date,omitempty"`
	Runtime int        `json:"runtime_min,omitempty"`
	Still   *string    `json:"still,omitempty"`
}

func NewEpisode(e *domain.Episode) *Episode {
	if e == nil {
		return nil
	}

	return &Episode{
		ID:      e.ID,
		Number:  e.Number,
		Title:   e.Title,
		AirDate: e.AirDate,
		Runtime: e.RuntimeMin,
		Still:   e.StillPath,
	}
}
