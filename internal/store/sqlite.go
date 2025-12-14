package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bastianvv/vio/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteExec interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type SQLiteStore struct {
	db   *sql.DB
	exec sqliteExec
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, err
	}

	_, _ = db.Exec(`PRAGMA busy_timeout = 5000;`)
	_, _ = db.Exec(`PRAGMA journal_mode = WAL;`)

	return &SQLiteStore{
		db:   db,
		exec: db,
	}, nil
}

func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLiteStore) WithTx(fn func(tx Store) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, _ = tx.Exec(`PRAGMA foreign_keys = ON;`)

	txStore := &SQLiteStore{
		db:   nil, 
		exec: tx,
	}

	if err := fn(txStore); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ============================================================================
// Libraries
// ============================================================================

func (s *SQLiteStore) CreateLibrary(lib *domain.Library) error {
	now := time.Now()
	lib.CreatedAt = now
	lib.UpdatedAt = now

	res, err := s.exec.Exec(`
        INSERT INTO libraries (name, type, path, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `, lib.Name, lib.Type, lib.Path, lib.CreatedAt, lib.UpdatedAt)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	lib.ID = id
	return nil
}

func (s *SQLiteStore) ListLibraries() ([]domain.Library, error) {
	rows, err := s.exec.Query(`
        SELECT id, name, type, path, created_at, updated_at
        FROM libraries
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var libs []domain.Library
	for rows.Next() {
		var l domain.Library
		if err := rows.Scan(
			&l.ID, &l.Name, &l.Type, &l.Path,
			&l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		libs = append(libs, l)
	}
	return libs, rows.Err()
}

func (s *SQLiteStore) GetLibrary(id int64) (*domain.Library, error) {
	var l domain.Library
	err := s.exec.QueryRow(`
        SELECT id, name, type, path, created_at, updated_at
        FROM libraries
        WHERE id = ?
    `, id).Scan(
		&l.ID, &l.Name, &l.Type, &l.Path,
		&l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *SQLiteStore) UpdateLibrary(lib *domain.Library) error {
	if lib.ID == 0 {
		return errors.New("cannot update library without ID")
	}

	now := time.Now().UTC()

	res, err := s.exec.Exec(`
		UPDATE libraries
		SET name = ?, type = ?, path = ?, updated_at = ?
		WHERE id = ?
	`,
		lib.Name,
		string(lib.Type),
		lib.Path,
		now,
		lib.ID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	lib.UpdatedAt = now
	return nil
}

// ============================================================================
// Movies
// ============================================================================

func (s *SQLiteStore) CreateMovie(m *domain.Movie) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	res, err := s.exec.Exec(`
        INSERT INTO movies (library_id, title, original_title, year, tmdb_id,
                            overview, runtime_min, poster_path, backdrop_path,
                            created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, m.LibraryID, m.Title, m.OriginalTitle, m.Year, m.TMDBID,
		m.Overview, m.RuntimeMin, m.PosterPath, m.BackdropPath,
		m.CreatedAt, m.UpdatedAt)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	m.ID = id
	return nil
}

func (s *SQLiteStore) GetMovie(id int64) (*domain.Movie, error) {
	var m domain.Movie
	err := s.exec.QueryRow(`
        SELECT id, library_id, title, original_title, year, tmdb_id,
               overview, runtime_min, poster_path, backdrop_path,
               created_at, updated_at
        FROM movies
        WHERE id = ?
    `, id).Scan(
		&m.ID, &m.LibraryID, &m.Title, &m.OriginalTitle, &m.Year, &m.TMDBID,
		&m.Overview, &m.RuntimeMin, &m.PosterPath, &m.BackdropPath,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *SQLiteStore) GetMovieByTitleAndYear(title string, year int, libraryID int64) (*domain.Movie, error) {
	row := s.exec.QueryRow(`
        SELECT id, library_id, title, original_title, year, tmdb_id, overview,
               runtime_min, poster_path, backdrop_path, created_at, updated_at
        FROM movies
        WHERE title = ? AND year = ? AND library_id = ?
    `, title, year, libraryID)

	var m domain.Movie
	err := row.Scan(
		&m.ID,
		&m.LibraryID,
		&m.Title,
		&m.OriginalTitle,
		&m.Year,
		&m.TMDBID,
		&m.Overview,
		&m.RuntimeMin,
		&m.PosterPath,
		&m.BackdropPath,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // no movie found — not an error
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *SQLiteStore) ListMoviesByLibrary(libraryID int64) ([]domain.Movie, error) {
	rows, err := s.exec.Query(`
        SELECT id, library_id, title, original_title, year, tmdb_id,
               overview, runtime_min, poster_path, backdrop_path,
               created_at, updated_at
        FROM movies
        WHERE library_id = ?
        ORDER BY title
    `, libraryID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []domain.Movie
	for rows.Next() {
		var m domain.Movie
		err := rows.Scan(
			&m.ID, &m.LibraryID, &m.Title, &m.OriginalTitle, &m.Year, &m.TMDBID,
			&m.Overview, &m.RuntimeMin, &m.PosterPath, &m.BackdropPath,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// ============================================================================
// Series / Seasons / Episodes
// ============================================================================

func (s *SQLiteStore) CreateSeries(sr *domain.Series) error {
	now := time.Now()
	sr.CreatedAt = now
	sr.UpdatedAt = now

	res, err := s.exec.Exec(`
        INSERT INTO series (library_id, title, original_title, tmdb_id, overview,
                            status, poster_path, backdrop_path, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, sr.LibraryID, sr.Title, sr.OriginalTitle, sr.TMDBID, sr.Overview,
		sr.Status, sr.PosterPath, sr.BackdropPath, sr.CreatedAt, sr.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	sr.ID = id
	return nil
}

func (s *SQLiteStore) GetSeriesByTitle(title string, libraryID int64) (*domain.Series, error) {
	row := s.exec.QueryRow(`
        SELECT id, library_id, title, original_title, tmdb_id, overview, status,
               poster_path, backdrop_path, created_at, updated_at
        FROM series
        WHERE title = ? AND library_id = ?
    `, title, libraryID)

	var sr domain.Series
	err := row.Scan(
		&sr.ID,
		&sr.LibraryID,
		&sr.Title,
		&sr.OriginalTitle,
		&sr.TMDBID,
		&sr.Overview,
		&sr.Status,
		&sr.PosterPath,
		&sr.BackdropPath,
		&sr.CreatedAt,
		&sr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func (s *SQLiteStore) GetSeries(id int64) (*domain.Series, error) {
	var sr domain.Series
	err := s.exec.QueryRow(`
        SELECT id, library_id, title, original_title, tmdb_id, overview,
               status, poster_path, backdrop_path, created_at, updated_at
        FROM series
        WHERE id = ?
    `, id).Scan(
		&sr.ID, &sr.LibraryID, &sr.Title, &sr.OriginalTitle, &sr.TMDBID,
		&sr.Overview, &sr.Status, &sr.PosterPath, &sr.BackdropPath,
		&sr.CreatedAt, &sr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (s *SQLiteStore) ListSeries() ([]*domain.Series, error) {
	rows, err := s.exec.Query(`
        SELECT id, library_id, title, overview, poster_path, backdrop_path,
               created_at, updated_at
        FROM series
        ORDER BY title
    `)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []*domain.Series
	for rows.Next() {
		sr := &domain.Series{}
		err := rows.Scan(
			&sr.ID,
			&sr.LibraryID,
			&sr.Title,
			&sr.Overview,
			&sr.PosterPath,
			&sr.BackdropPath,
			&sr.CreatedAt,
			&sr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, sr)
	}

	return out, rows.Err()
}

func (s *SQLiteStore) CreateSeason(se *domain.Season) error {
	now := time.Now()
	se.CreatedAt = now
	se.UpdatedAt = now

	res, err := s.exec.Exec(`
        INSERT INTO seasons (series_id, season_number, title, overview,
                             poster_path, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, se.SeriesID, se.Number, se.Title, se.Overview,
		se.PosterPath, se.CreatedAt, se.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	se.ID = id
	return nil
}

func (s *SQLiteStore) GetSeasonBySeriesAndNumber(seriesID int64, number int) (*domain.Season, error) {
	row := s.exec.QueryRow(`
        SELECT id, series_id, season_number, title, overview, poster_path,
               created_at, updated_at
        FROM seasons
        WHERE series_id = ? AND season_number = ?
    `, seriesID, number)

	var se domain.Season
	err := row.Scan(
		&se.ID,
		&se.SeriesID,
		&se.Number,
		&se.Title,
		&se.Overview,
		&se.PosterPath,
		&se.CreatedAt,
		&se.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &se, nil
}

func (s *SQLiteStore) GetSeason(id int64) (*domain.Season, error) {
	row := s.exec.QueryRow(`
        SELECT id, series_id, season_number, title, overview, poster_path,
               air_date, created_at, updated_at
        FROM seasons
        WHERE id = ?
    `, id)

	var se domain.Season
	err := row.Scan(
		&se.ID,
		&se.SeriesID,
		&se.Number,
		&se.Title,
		&se.Overview,
		&se.PosterPath,
		&se.AirDate,
		&se.CreatedAt,
		&se.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &se, nil
}

func (s *SQLiteStore) ListSeasonsBySeries(seriesID int64) ([]domain.Season, error) {
	rows, err := s.exec.Query(`
        SELECT id, series_id, season_number, title, overview, poster_path,
               air_date, created_at, updated_at
        FROM seasons
        WHERE series_id = ?
        ORDER BY season_number
    `, seriesID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []domain.Season
	for rows.Next() {
		var se domain.Season
		err := rows.Scan(
			&se.ID,
			&se.SeriesID,
			&se.Number,
			&se.Title,
			&se.Overview,
			&se.PosterPath,
			&se.AirDate,
			&se.CreatedAt,
			&se.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, se)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CreateEpisode(ep *domain.Episode) error {
	now := time.Now()
	ep.CreatedAt = now
	ep.UpdatedAt = now

	res, err := s.exec.Exec(`
        INSERT INTO episodes (season_id, episode_number, title, overview,
                              air_date, runtime_min, still_path, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, ep.SeasonID, ep.Number, ep.Title, ep.Overview,
		ep.AirDate, ep.RuntimeMin, ep.StillPath,
		ep.CreatedAt, ep.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	ep.ID = id
	return nil
}

func (s *SQLiteStore) GetEpisode(id int64) (*domain.Episode, error) {
	var ep domain.Episode
	err := s.exec.QueryRow(`
        SELECT id, season_id, episode_number, title, overview,
               air_date, runtime_min, still_path, created_at, updated_at
        FROM episodes
        WHERE id = ?
    `, id).Scan(
		&ep.ID, &ep.SeasonID, &ep.Number, &ep.Title, &ep.Overview,
		&ep.AirDate, &ep.RuntimeMin, &ep.StillPath,
		&ep.CreatedAt, &ep.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

func (s *SQLiteStore) GetEpisodeBySeasonAndNumber(seasonID int64, number int) (*domain.Episode, error) {
	row := s.exec.QueryRow(`
        SELECT id, season_id, episode_number, title, overview, air_date, runtime_min,
               still_path, created_at, updated_at
        FROM episodes
        WHERE season_id = ? AND episode_number = ?
    `, seasonID, number)

	var ep domain.Episode
	err := row.Scan(
		&ep.ID,
		&ep.SeasonID,
		&ep.Number,
		&ep.Title,
		&ep.Overview,
		&ep.AirDate,
		&ep.RuntimeMin,
		&ep.StillPath,
		&ep.CreatedAt,
		&ep.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ep, nil
}

func (s *SQLiteStore) ListEpisodesBySeason(seasonID int64) ([]domain.Episode, error) {
	rows, err := s.exec.Query(`
        SELECT id, season_id, episode_number, title, overview, air_date,
               runtime_min, still_path, created_at, updated_at
        FROM episodes
        WHERE season_id = ?
        ORDER BY episode_number
    `, seasonID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var episodes []domain.Episode
	for rows.Next() {
		var ep domain.Episode
		err := rows.Scan(
			&ep.ID, &ep.SeasonID, &ep.Number, &ep.Title, &ep.Overview,
			&ep.AirDate, &ep.RuntimeMin, &ep.StillPath,
			&ep.CreatedAt, &ep.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, ep)
	}
	return episodes, rows.Err()
}

// ============================================================================
// Media Files
// ============================================================================
func (s *SQLiteStore) GetMediaFile(id int64) (*domain.MediaFile, error) {
	row := s.exec.QueryRow(`
        SELECT id, library_id, movie_id, episode_id, path, size_bytes,
               hash, is_missing, last_seen_at, missing_since, container, video_codec, audio_codec,
               video_width, video_height, audio_channels, duration_sec,
               created_at, updated_at
        FROM media_files
        WHERE id = ?
    `, id)

	var mf domain.MediaFile

	err := row.Scan(
		&mf.ID,
		&mf.LibraryID,
		&mf.MovieID,
		&mf.EpisodeID,
		&mf.Path,
		&mf.SizeBytes,
		&mf.Hash,
		&mf.IsMissing,
		&mf.LastSeenAt,
		&mf.MissingSince,
		&mf.Container,
		&mf.VideoCodec,
		&mf.AudioCodec,
		&mf.VideoWidth,
		&mf.VideoHeight,
		&mf.AudioChannels,
		&mf.DurationSec,
		&mf.CreatedAt,
		&mf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &mf, nil
}

func (s *SQLiteStore) CreateMediaFile(mf *domain.MediaFile) error {
	now := time.Now().UTC()
	mf.CreatedAt = now
	mf.UpdatedAt = now

	res, err := s.exec.Exec(`
		INSERT INTO media_files (
		    library_id,
		    path,
		    size_bytes,
		    hash,
		    is_missing,
		    last_seen_at,
		    missing_since,
		    container,
		    video_codec,
		    audio_codec,
		    video_width,
		    video_height,
		    audio_channels,
		    duration_sec,
		    created_at,
		    updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		mf.LibraryID,
		mf.Path,
		mf.SizeBytes,
		mf.Hash,
		mf.IsMissing,
		mf.LastSeenAt,
		mf.MissingSince,
		mf.Container,
		mf.VideoCodec,
		mf.AudioCodec,
		mf.VideoWidth,
		mf.VideoHeight,
		mf.AudioChannels,
		mf.DurationSec,
		mf.CreatedAt,
		mf.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	mf.ID = id
	return nil
}

func (s *SQLiteStore) CreateMediaFileEpisode(link *domain.MediaFileEpisode) error {
	res, err := s.exec.Exec(`
        INSERT INTO media_file_episodes (media_file_id, episode_id)
        VALUES (?, ?)
    `, link.MediaFileID, link.EpisodeID)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	link.ID = id
	return nil
}

func (s *SQLiteStore) ListEpisodesByMediaFile(mediaFileID int64) ([]domain.Episode, error) {
	rows, err := s.exec.Query(`
        SELECT e.id, e.season_id, e.episode_number, e.title, e.overview,
               e.air_date, e.runtime_min, e.still_path, e.created_at, e.updated_at
        FROM media_file_episodes mfe
        JOIN episodes e ON e.id = mfe.episode_id
        WHERE mfe.media_file_id = ?
        ORDER BY e.episode_number
    `, mediaFileID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var eps []domain.Episode
	for rows.Next() {
		var ep domain.Episode
		err := rows.Scan(
			&ep.ID,
			&ep.SeasonID,
			&ep.Number,
			&ep.Title,
			&ep.Overview,
			&ep.AirDate,
			&ep.RuntimeMin,
			&ep.StillPath,
			&ep.CreatedAt,
			&ep.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}
	return eps, rows.Err()
}

func (s *SQLiteStore) ListMediaFilesByEpisode(episodeID int64) ([]domain.MediaFile, error) {
	rows, err := s.exec.Query(`
        SELECT id, library_id, movie_id, episode_id, path, size_bytes, hash,
               container, video_codec, audio_codec, video_width, video_height,
               audio_channels, duration_sec,
               created_at, updated_at
        FROM media_files
        WHERE episode_id = ?
        ORDER BY id
    `, episodeID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []domain.MediaFile
	for rows.Next() {
		var mf domain.MediaFile
		err := rows.Scan(
			&mf.ID, &mf.LibraryID, &mf.MovieID, &mf.EpisodeID, &mf.Path, &mf.SizeBytes,
			&mf.Hash, &mf.Container, &mf.VideoCodec, &mf.AudioCodec,
			&mf.VideoWidth, &mf.VideoHeight, &mf.AudioChannels, &mf.DurationSec,
			&mf.CreatedAt, &mf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, mf)
	}
	return result, rows.Err()
}

func (s *SQLiteStore) ListMediaFilesByMovie(movieID int64) ([]domain.MediaFile, error) {
	rows, err := s.exec.Query(`
        SELECT id, library_id, movie_id, episode_id, path, size_bytes,
               hash, container, video_codec, audio_codec,
               video_width, video_height, audio_channels, duration_sec,
               created_at, updated_at
        FROM media_files
        WHERE movie_id = ?
        ORDER BY id
    `, movieID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []domain.MediaFile

	for rows.Next() {
		var mf domain.MediaFile
		err := rows.Scan(
			&mf.ID,
			&mf.LibraryID,
			&mf.MovieID,
			&mf.EpisodeID,
			&mf.Path,
			&mf.SizeBytes,
			&mf.Hash,
			&mf.Container,
			&mf.VideoCodec,
			&mf.AudioCodec,
			&mf.VideoWidth,
			&mf.VideoHeight,
			&mf.AudioChannels,
			&mf.DurationSec,
			&mf.CreatedAt,
			&mf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, mf)
	}

	return list, rows.Err()
}

func (s *SQLiteStore) GetMediaFileByPath(path string) (*domain.MediaFile, error) {
	const q = `
		SELECT
			id,
			library_id,
			path,
			size_bytes,
			hash,
			container,
			video_codec,
			audio_codec,
			video_width,
			video_height,
			audio_channels,
			duration_sec,
			movie_id,
			episode_id
		FROM media_files
		WHERE path = ?
		LIMIT 1
	`

	var mf domain.MediaFile

	err := s.exec.QueryRow(q, path).Scan(
		&mf.ID,
		&mf.LibraryID,
		&mf.Path,
		&mf.SizeBytes,
		&mf.Hash,
		&mf.Container,
		&mf.VideoCodec,
		&mf.AudioCodec,
		&mf.VideoWidth,
		&mf.VideoHeight,
		&mf.AudioChannels,
		&mf.DurationSec,
		&mf.MovieID,
		&mf.EpisodeID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &mf, nil
}

func (s *SQLiteStore) UpdateMediaFile(mf *domain.MediaFile) error {
	const q = `
	UPDATE media_files SET
	    size_bytes = ?,
	    hash = ?,
	    container = ?,
	    video_codec = ?,
	    audio_codec = ?,
	    video_width = ?,
	    video_height = ?,
	    audio_channels = ?,
	    duration_sec = ?,
	    is_missing = FALSE,
	    last_seen_at = ?,
	    missing_since = NULL,
	    updated_at = ?
	WHERE id = ?
	`

	_, err := s.exec.Exec(
		q,
		mf.SizeBytes,
		mf.Hash,
		mf.Container,
		mf.VideoCodec,
		mf.AudioCodec,
		mf.VideoWidth,
		mf.VideoHeight,
		mf.AudioChannels,
		mf.DurationSec,
		mf.LastSeenAt,
		time.Now().UTC(),
		mf.ID,
	)
	return err
}

func (s *SQLiteStore) MarkMissingMediaFiles(
	libraryID int64,
	scanStartedAt time.Time,
) (int64, error) {

	const q = `
		UPDATE media_files
		SET
			is_missing = TRUE,
			missing_since = COALESCE(missing_since, ?)
		WHERE
			library_id = ?
			AND (
				last_seen_at IS NULL
				OR last_seen_at < ?
			)
			AND is_missing = FALSE
	`

	res, err := s.exec.Exec(q, scanStartedAt, libraryID, scanStartedAt)
	if err != nil {
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (s *SQLiteStore) MarkMediaFileSeen(
	id int64,
	seenAt time.Time,
) error {

	const q = `
		UPDATE media_files
		SET
			last_seen_at = ?,
			is_missing = FALSE,
			missing_since = NULL,
			updated_at = ?
		WHERE id = ?
	`

	_, err := s.exec.Exec(q, seenAt, seenAt, id)
	return err
}

// ============================================================================
// Subtitles
// ============================================================================

func (s *SQLiteStore) CreateSubtitleTrack(st *domain.SubtitleTrack) error {
	res, err := s.exec.Exec(`
        INSERT INTO subtitle_tracks (
            media_file_id, source, external_path, stream_index, language,
            is_forced, is_default, format
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, st.MediaFileID, st.Source, st.ExternalPath, st.StreamIndex, st.Language,
		st.IsForced, st.IsDefault, st.Format)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	st.ID = id
	return nil
}

func (s *SQLiteStore) ListSubtitleTracks(mediaFileID int64) ([]domain.SubtitleTrack, error) {
	rows, err := s.exec.Query(`
        SELECT id, media_file_id, source, external_path, stream_index, language,
               is_forced, is_default, format
        FROM subtitle_tracks
        WHERE media_file_id = ?
        ORDER BY id
    `, mediaFileID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []domain.SubtitleTrack
	for rows.Next() {
		var st domain.SubtitleTrack
		err := rows.Scan(
			&st.ID, &st.MediaFileID, &st.Source, &st.ExternalPath, &st.StreamIndex,
			&st.Language, &st.IsForced, &st.IsDefault, &st.Format,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, st)
	}
	return list, rows.Err()
}

// ============================================================================
// Clean-up
// ============================================================================

func (s *SQLiteStore) CleanupEmptySeries(libraryID int64) (int64, error) {
	const q = `
	DELETE FROM series
	WHERE id IN (
		SELECT sr.id
		FROM series sr
		LEFT JOIN seasons s ON s.series_id = sr.id
		WHERE sr.library_id = ?
		GROUP BY sr.id
		HAVING COUNT(s.id) = 0
	)
	`
	res, err := s.exec.Exec(q, libraryID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) CleanupEmptySeasons(libraryID int64) (int64, error) {
	const q = `
	DELETE FROM seasons
	WHERE id IN (
		SELECT s.id
		FROM seasons s
		JOIN series sr ON sr.id = s.series_id
		LEFT JOIN episodes e ON e.season_id = s.id
		WHERE sr.library_id = ?
		GROUP BY s.id
		HAVING COUNT(e.id) = 0
	)
	`
	res, err := s.exec.Exec(q, libraryID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) CleanupEmptyEpisodes(libraryID int64) (int64, error) {
	const q = `
	DELETE FROM episodes
	WHERE id IN (
		SELECT e.id
		FROM episodes e
		JOIN seasons s ON s.id = e.season_id
		JOIN series sr ON sr.id = s.series_id
		LEFT JOIN media_files mf
			ON mf.episode_id = e.id
			AND mf.is_missing = 0
		WHERE sr.library_id = ?
		GROUP BY e.id
		HAVING COUNT(mf.id) = 0
	)
	`
	res, err := s.exec.Exec(q, libraryID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) CleanupMissingMediaFileLinks(libraryID int64) (int64, error) {
	const q = `
	DELETE FROM media_file_episodes
	WHERE media_file_id IN (
		SELECT id
		FROM media_files
		WHERE library_id = ?
		  AND is_missing = 1
	)
	`
	res, err := s.exec.Exec(q, libraryID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) UnlinkMissingMediaFiles(libraryID int64) (int64, error) {
	// 1) Delete range links for missing files
	const delLinks = `
	DELETE FROM media_file_episodes
	WHERE media_file_id IN (
		SELECT id FROM media_files
		WHERE library_id = ? AND is_missing = 1
	)
	`
	if _, err := s.exec.Exec(delLinks, libraryID); err != nil {
		return 0, err
	}

	// 2) Null direct links
	const upd = `
	UPDATE media_files
	SET episode_id = NULL,
	    movie_id = NULL,
	    updated_at = ?
	WHERE library_id = ?
	  AND is_missing = 1
	  AND (episode_id IS NOT NULL OR movie_id IS NOT NULL)
	`
	res, err := s.exec.Exec(upd, time.Now().UTC(), libraryID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
