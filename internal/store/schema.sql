PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS libraries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

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
    FOREIGN KEY(library_id) REFERENCES libraries(id)
);

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
    FOREIGN KEY(library_id) REFERENCES libraries(id)
);

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
    FOREIGN KEY(series_id) REFERENCES series(id)
);

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
    FOREIGN KEY(season_id) REFERENCES seasons(id)
);

CREATE TABLE IF NOT EXISTS media_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    library_id INTEGER NOT NULL,
    is_missing BOOLEAN NOT NULL DEFAULT FALSE,
    missing_since DATETIME NULL,
    last_seen_at DATETIME NULL,
    movie_id INTEGER,
    episode_id INTEGER,
    path TEXT NOT NULL,
    size_bytes INTEGER,
    hash TEXT UNIQUE ON CONFLICT IGNORE,
    container TEXT,
    video_codec TEXT,
    audio_codec TEXT,
    video_width INTEGER,
    video_height INTEGER,
    audio_channels INTEGER,
    duration_sec INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(library_id) REFERENCES libraries(id),
    FOREIGN KEY(movie_id) REFERENCES movies(id),
    FOREIGN KEY(episode_id) REFERENCES episodes(id)
);

CREATE TABLE IF NOT EXISTS media_file_episodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_file_id INTEGER NOT NULL,
    episode_id INTEGER NOT NULL,
    FOREIGN KEY(media_file_id) REFERENCES media_files(id),
    FOREIGN KEY(episode_id) REFERENCES episodes(id)
);

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
    FOREIGN KEY(media_file_id) REFERENCES media_files(id)
);
