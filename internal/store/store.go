package store

import (
	"time"

	"github.com/bastianvv/vio/internal/domain"
)

type Store interface {
	Close() error

	// Libraries
	CreateLibrary(lib *domain.Library) error
	ListLibraries() ([]domain.Library, error)
	GetLibrary(id int64) (*domain.Library, error)
	UpdateLibrary(lib *domain.Library) error

	// Movies
	CreateMovie(m *domain.Movie) error
	ListMoviesByLibrary(libraryID int64) ([]domain.Movie, error)
	GetMovie(id int64) (*domain.Movie, error)
	GetMovieByTitleAndYear(title string, year int, libraryID int64) (*domain.Movie, error)

	// Series
	CreateSeries(s *domain.Series) error
	GetSeries(id int64) (*domain.Series, error)
	GetSeriesByTitle(title string, libraryID int64) (*domain.Series, error)
	ListSeries() ([]*domain.Series, error)

	// Seasons
	CreateSeason(season *domain.Season) error
	GetSeason(id int64) (*domain.Season, error)
	GetSeasonBySeriesAndNumber(seriesID int64, number int) (*domain.Season, error)
	ListSeasonsBySeries(seriesID int64) ([]domain.Season, error)

	// Episodes
	CreateEpisode(ep *domain.Episode) error
	GetEpisode(id int64) (*domain.Episode, error)
	ListEpisodesBySeason(seasonID int64) ([]domain.Episode, error)
	GetEpisodeBySeasonAndNumber(seasonID int64, number int) (*domain.Episode, error)

	// Media files
	CreateMediaFile(mf *domain.MediaFile) error
	GetMediaFile(id int64) (*domain.MediaFile, error)
	ListMediaFilesByMovie(movieID int64) ([]domain.MediaFile, error)
	ListMediaFilesByEpisode(episodeID int64) ([]domain.MediaFile, error)
	CreateMediaFileEpisode(link *domain.MediaFileEpisode) error
	ListEpisodesByMediaFile(mediaFileID int64) ([]domain.Episode, error)
	GetMediaFileByPath(path string) (*domain.MediaFile, error)
	UpdateMediaFile(mf *domain.MediaFile) error
	MarkMissingMediaFiles(libraryID int64, scanStartedAt time.Time) (int64, error)
	MarkMediaFileSeen(id int64, seenAt time.Time) error

	// Subtitles
	CreateSubtitleTrack(st *domain.SubtitleTrack) error
	ListSubtitleTracks(mediaFileID int64) ([]domain.SubtitleTrack, error)

	// Cleanup
	CleanupMissingMediaFileLinks(libraryID int64) (int64, error)
	CleanupEmptyEpisodes(libraryID int64) (int64, error)
	CleanupEmptySeasons(libraryID int64) (int64, error)
	CleanupEmptySeries(libraryID int64) (int64, error)
	UnlinkMissingMediaFiles(libraryId int64) (int64, error)
}
