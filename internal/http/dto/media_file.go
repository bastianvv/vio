package dto

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type MediaFile struct {
	ID            int64  `json:"id"`
	LibraryID     int64  `json:"library_id"`
	MovieID       *int64 `json:"movie_id,omitempty"`
	EpisodeID     *int64 `json:"episode_id,omitempty"`
	Path          string `json:"path"`
	Container     string `json:"container"`
	VideoCodec    string `json:"video_codec"`
	AudioCodec    string `json:"audio_codec"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	AudioChannels int    `json:"audio_channels"`
	DurationSec   int    `json:"duration_sec"`

	IsMissing    bool       `json:"is_missing"`
	MissingSince *time.Time `json:"missing_since,omitempty"`
	LastSeenAt   *time.Time `json:"last_seen_at,omitempty"`
}

func NewMediaFile(m *domain.MediaFile) *MediaFile {
	if m == nil {
		return nil
	}

	return &MediaFile{
		ID:            m.ID,
		LibraryID:     m.LibraryID,
		MovieID:       m.MovieID,
		EpisodeID:     m.EpisodeID,
		Path:          m.Path,
		Container:     m.Container,
		VideoCodec:    m.VideoCodec,
		AudioCodec:    m.AudioCodec,
		Width:         m.VideoWidth,
		Height:        m.VideoHeight,
		AudioChannels: m.AudioChannels,
		DurationSec:   m.DurationSec,
		IsMissing:     m.IsMissing,
		MissingSince:  m.MissingSince,
		LastSeenAt:    m.LastSeenAt,
	}
}
