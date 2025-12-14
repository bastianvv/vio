PRAGMA foreign_keys = ON;

-- Libraries
CREATE TABLE IF NOT EXISTS libraries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(path) -- optional but useful: one library per root path
);

-- Movies
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    original_title TEXT,
    year INTEGER,
    tmdb_id TEXT,
    overview TEXT,
    runtime_min INTEGER,
    poster_path TEXT,
    backdrop_path TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(library_id) REFERENCES libraries(id) ON DELETE CASCADE,
    UNIQUE(library_id, title, year)
);

-- Series
CREATE TABLE IF NOT EXISTS series (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    original_title TEXT,
    tmdb_id TEXT,
    overview TEXT,
    status TEXT,
    poster_path TEXT,
    backdrop_path TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(library_id) REFERENCES libraries(id) ON DELETE CASCADE,
    UNIQUE(library_id, title)
);

-- Seasons
CREATE TABLE IF NOT EXISTS seasons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id INTEGER NOT NULL,
    season_number INTEGER NOT NULL,
    title TEXT,
    overview TEXT,
    poster_path TEXT,
    air_date DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE,
    UNIQUE(series_id, season_number)
);

-- Episodes
CREATE TABLE IF NOT EXISTS episodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    season_id INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,
    title TEXT,
    overview TEXT,
    air_date DATETIME,
    runtime_min INTEGER,
    still_path TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(season_id) REFERENCES seasons(id) ON DELETE CASCADE,
    UNIQUE(season_id, episode_number)
);

-- Media Files
CREATE TABLE IF NOT EXISTS media_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id INTEGER NOT NULL,

    is_missing BOOLEAN NOT NULL DEFAULT FALSE,
    missing_since DATETIME NULL,
    last_seen_at DATETIME NULL,

    movie_id INTEGER NULL,
    episode_id INTEGER NULL,

    path TEXT NOT NULL,
    size_bytes INTEGER,
    hash TEXT, -- NOT UNIQUE (same content across different paths is valid)

    container TEXT,
    video_codec TEXT,
    audio_codec TEXT,
    video_width INTEGER,
    video_height INTEGER,
    audio_channels INTEGER,
    duration_sec INTEGER,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY(library_id) REFERENCES libraries(id) ON DELETE CASCADE,
    FOREIGN KEY(movie_id) REFERENCES movies(id) ON DELETE SET NULL,
    FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE SET NULL,

    UNIQUE(path),

    CHECK (
      (is_missing = 0 AND missing_since IS NULL)
      OR
      (is_missing = 1 AND missing_since IS NOT NULL)
    )
);

-- Multi-episode links (range files)
CREATE TABLE IF NOT EXISTS media_file_episodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_file_id INTEGER NOT NULL,
    episode_id INTEGER NOT NULL,
    FOREIGN KEY(media_file_id) REFERENCES media_files(id) ON DELETE CASCADE,
    FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE CASCADE,
    UNIQUE(media_file_id, episode_id)
);

-- Subtitles
CREATE TABLE IF NOT EXISTS subtitle_tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_file_id INTEGER NOT NULL,
    source TEXT NOT NULL,
    external_path TEXT,
    stream_index INTEGER,
    language TEXT,
    is_forced BOOLEAN,
    is_default BOOLEAN,
    format TEXT,
    FOREIGN KEY(media_file_id) REFERENCES media_files(id) ON DELETE CASCADE
);
