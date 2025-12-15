package dto

import "github.com/bastianvv/vio/internal/domain"

type SubtitleTrack struct {
	ID        int64  `json:"id"`
	Source    string `json:"source"` // "EMBEDDED" | "EXTERNAL"
	Language  string `json:"language,omitempty"`
	Format    string `json:"format,omitempty"`
	IsForced  bool   `json:"is_forced"`
	IsDefault bool   `json:"is_default"`

	// Only for embedded subtitles
	StreamIndex *int `json:"stream_index,omitempty"`

	// Only for external subtitles
	External bool `json:"external"`
}

func NewSubtitleTrack(st *domain.SubtitleTrack) *SubtitleTrack {
	if st == nil {
		return nil
	}

	return &SubtitleTrack{
		ID:          st.ID,
		Source:      string(st.Source),
		Language:    st.Language,
		Format:      st.Format,
		IsForced:    st.IsForced,
		IsDefault:   st.IsDefault,
		StreamIndex: st.StreamIndex,
		External:    st.Source == domain.SubtitleSourceExternal,
	}
}
