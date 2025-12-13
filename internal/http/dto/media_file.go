package dto

import "github.com/bastianvv/vio/internal/domain"

type MediaFile struct {
	ID            int64  `json:"id"`
	Path          string `json:"path"`
	Container     string `json:"container"`
	VideoCodec    string `json:"video_codec"`
	AudioCodec    string `json:"audio_codec"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	AudioChannels int    `json:"audio_channels"`
	Duration      int    `json:"duration"`
}

func NewMediaFile(m *domain.MediaFile) *MediaFile {
	if m == nil {
		return nil
	}

	return &MediaFile{
		ID:            m.ID,
		Path:          m.Path,
		Container:     m.Container,
		VideoCodec:    m.VideoCodec,
		AudioCodec:    m.AudioCodec,
		Width:         m.VideoWidth,
		Height:        m.VideoHeight,
		AudioChannels: m.AudioChannels,
		Duration:      m.DurationSec,
	}
}
