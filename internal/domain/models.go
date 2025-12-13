package domain

import "time"

type LibraryType string

const (
	LibraryTypeMovies LibraryType = "movies"
	LibraryTypeSeries LibraryType = "series"
	LibraryTypeAnime  LibraryType = "anime"
	LibraryTypeOther  LibraryType = "other"
)

type Library struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Type      LibraryType `json:"type"`
	Path      string      `json:"path"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Movie struct {
	ID            int64     `json:"id"`
	LibraryID     int64     `json:"library_id"`
	Title         string    `json:"title"`
	OriginalTitle string    `json:"original_title"`
	Year          int       `json:"year"`
	TMDBID        *string   `json:"tmdb_id,omitempty"`
	Overview      string    `json:"overview"`
	RuntimeMin    int       `json:"runtime_min"`
	PosterPath    *string   `json:"poster_path,omitempty"`
	BackdropPath  *string   `json:"backdrop_path,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Series struct {
	ID            int64     `json:"id"`
	LibraryID     int64     `json:"library_id"`
	Title         string    `json:"title"`
	OriginalTitle string    `json:"original_title"`
	TMDBID        *string   `json:"tmdb_id,omitempty"`
	Overview      string    `json:"overview"`
	Status        string    `json:"status"` // ongoing, ended, etc.
	PosterPath    *string   `json:"poster_path,omitempty"`
	BackdropPath  *string   `json:"backdrop_path,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Season struct {
	ID         int64      `json:"id"`
	SeriesID   int64      `json:"series_id"`
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	Overview   string     `json:"overview"`
	PosterPath *string    `json:"poster_path,omitempty"`
	AirDate    *time.Time `json:"air_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Episode struct {
	ID         int64      `json:"id"`
	SeasonID   int64      `json:"season_id"`
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	Overview   string     `json:"overview"`
	AirDate    *time.Time `json:"air_date,omitempty"`
	RuntimeMin int        `json:"runtime_min"`
	StillPath  *string    `json:"still_path,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type MediaFile struct {
	ID            int64      `json:"id"`
	LibraryID     int64      `json:"library_id"`
	MovieID       *int64     `json:"movie_id,omitempty"`
	EpisodeID     *int64     `json:"episode_id,omitempty"`
	Path          string     `json:"path"`
	SizeBytes     int64      `json:"size_bytes"`
	Hash          string     `json:"hash"`
	LastSeenAt    *time.Time `json:"last_seen_at"`
	IsMissing     bool       `json:"is_missing"`
	MissingSince  *time.Time `json:"missing_since"`
	Container     string     `json:"container"`
	VideoCodec    string     `json:"video_codec"`
	AudioCodec    string     `json:"audio_codec"`
	VideoWidth    int        `json:"video_width"`
	VideoHeight   int        `json:"video_height"`
	AudioChannels int        `json:"audio_channels"`
	DurationSec   int        `json:"duration_sec"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type MediaFileEpisode struct {
	ID          int64 `json:"id"`
	MediaFileID int64 `json:"media_file_id"`
	EpisodeID   int64 `json:"episode_id"`
}

type SubtitleSource string

const (
	SubtitleSourceEmbedded SubtitleSource = "EMBEDDED"
	SubtitleSourceExternal SubtitleSource = "EXTERNAL"
)

type SubtitleTrack struct {
	ID           int64          `json:"id"`
	MediaFileID  int64          `json:"media_file_id"`
	Source       SubtitleSource `json:"source"`
	ExternalPath *string        `json:"external_path,omitempty"`
	StreamIndex  *int           `json:"stream_index,omitempty"`
	Language     string         `json:"language"`
	IsForced     bool           `json:"is_forced"`
	IsDefault    bool           `json:"is_default"`
	Format       string         `json:"format"`
}
