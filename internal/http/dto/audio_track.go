package dto

import "github.com/bastianvv/vio/internal/domain"

type AudioTrack struct {
	ID          int64  `json:"id"`
	MediaFileID int64  `json:"media_file_id"`
	StreamIndex int    `json:"stream_index"`
	Language    string `json:"language,omitempty"`
	Codec       string `json:"codec,omitempty"`
	Channels    int    `json:"channels,omitempty"`
	IsDefault   bool   `json:"is_default"`
}

func NewAudioTrack(at *domain.AudioTrack) *AudioTrack {
	if at == nil {
		return nil
	}

	return &AudioTrack{
		ID:          at.ID,
		MediaFileID: at.MediaFileID,
		StreamIndex: at.StreamIndex,
		Language:    at.Language,
		Codec:       at.Codec,
		Channels:    at.Channels,
		IsDefault:   at.IsDefault,
	}
}
