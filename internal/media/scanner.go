package media

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bastianvv/vio/internal/domain"
	"github.com/bastianvv/vio/internal/store"
	"github.com/bastianvv/vio/internal/util"
)

var (
	reSxxExx      = regexp.MustCompile(`(?i)s(\d{1,2})e(\d{1,3})`)
	reSxxExxRange = regexp.MustCompile(`(?i)s(\d{1,2})e(\d{1,3})[-_](\d{1,3})`)
	reNxxN        = regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,3})`)
	reAnimeEp     = regexp.MustCompile(`(?i)(.*?)[\s._-]+(\d{1,4})(?:\D|$)`)
)

type ScanMode int

const (
	ScanModeIncremental ScanMode = iota
	ScanModeRescan
)

type ScanResult struct {
	LibraryID int64

	FilesScanned  int
	MoviesAdded   int
	SeriesAdded   int
	EpisodesAdded int

	Errors []error
}

type attachResult struct {
	SeriesCreated   bool
	SeasonCreated   bool
	EpisodesCreated int
	MovieCreated    bool
}

type Scanner interface {
	ScanLibrary(lib *domain.Library, mode ScanMode) (*ScanResult, error)
}

type FSScanner struct {
	store store.Store
}

func NewScanner(s store.Store) *FSScanner {
	return &FSScanner{store: s}
}

// Recognized video extensions (MVP).
var videoExt = map[string]bool{
	".mkv": true,
	".mp4": true,
	".avi": true,
}

// ScanLibrary walks the filesystem starting from lib.Path and
// processes all supported video files.
func (s *FSScanner) ScanLibrary(lib *domain.Library, mode ScanMode) (*ScanResult, error) {

	result := &ScanResult{
		LibraryID: lib.ID,
	}

	scanStartedAt := time.Now().UTC()

	walkErr := filepath.WalkDir(lib.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !videoExt[ext] {
			return nil
		}

		result.FilesScanned++

		err = s.store.WithTx(func(tx store.Store) error {
			return s.processVideoFileTx(
				tx,
				lib,
				mode,
				path,
				scanStartedAt,
				result,
			)
		})
		if err != nil {
			result.Errors = append(result.Errors, err)
		}

		return nil
	})

	if result.FilesScanned == 0 {
		return result, nil
	}

	// ONE cleanup pass, ONE transaction
	_ = s.store.WithTx(func(tx store.Store) error {
		if _, err := tx.MarkMissingMediaFiles(lib.ID, scanStartedAt); err != nil {
			return err
		}
		if _, err := tx.UnlinkMissingMediaFiles(lib.ID); err != nil {
			return err
		}
		if _, err := tx.CleanupEmptyEpisodes(lib.ID); err != nil {
			return err
		}
		if _, err := tx.CleanupEmptySeasons(lib.ID); err != nil {
			return err
		}
		if _, err := tx.CleanupEmptySeries(lib.ID); err != nil {
			return err
		}
		return nil
	})

	return result, walkErr
}

func (s *FSScanner) processVideoFileTx(
	tx store.Store,
	lib *domain.Library,
	mode ScanMode,
	path string,
	scanStartedAt time.Time,
	result *ScanResult,
) error {

	// Check existing media file
	existingMF, err := tx.GetMediaFileByPath(path)
	if err != nil {
		return err
	}

	if existingMF != nil && mode == ScanModeIncremental {
		return tx.MarkMediaFileSeen(existingMF.ID, scanStartedAt)
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	var hash string
	if mode == ScanModeRescan && existingMF != nil {
		hash, err = util.HashFile(path)
		if err != nil {
			return err
		}
		if existingMF.Hash == hash {
			return tx.MarkMediaFileSeen(existingMF.ID, scanStartedAt)
		}
	}

	if hash == "" {
		hash, err = util.HashFile(path)
		if err != nil {
			return err
		}
	}

	ffdata, err := RunFFProbe(path)
	if err != nil {
		return err
	}

	container := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")

	var (
		videoCodec    string
		audioCodec    string
		width         int
		height        int
		audioChannels int
		durationSec   int
	)

	if ffdata.Format.Duration != "" {
		if f, err := strconv.ParseFloat(ffdata.Format.Duration, 64); err == nil {
			durationSec = int(f + 0.5)
		}
	}

	for _, st := range ffdata.Streams {
		switch st.CodecType {
		case "video":
			videoCodec = st.CodecName
			width, height = st.Width, st.Height
		case "audio":
			audioCodec = st.CodecName
			audioChannels = st.Channels
		}
	}

	mf := &domain.MediaFile{
		LibraryID:     lib.ID,
		Path:          path,
		SizeBytes:     info.Size(),
		Hash:          hash,
		LastSeenAt:    &scanStartedAt,
		IsMissing:     false,
		Container:     container,
		VideoCodec:    videoCodec,
		AudioCodec:    audioCodec,
		VideoWidth:    width,
		VideoHeight:   height,
		AudioChannels: audioChannels,
		DurationSec:   durationSec,
	}

	if existingMF != nil {
		mf.ID = existingMF.ID
	}

	var episodes []*domain.Episode

	switch lib.Type {
	case domain.LibraryTypeMovies:
		ar, err := s.attachMovieTx(tx, lib, mf)
		if err != nil {
			return err
		}
		if ar.MovieCreated {
			result.MoviesAdded++
		}

	case domain.LibraryTypeSeries, domain.LibraryTypeAnime:
		eps, ar, err := s.attachSeriesEpisodeTx(tx, lib, mf)
		if err != nil {
			return err
		}
		if ar.SeriesCreated {
			result.SeriesAdded++
		}
		result.EpisodesAdded += ar.EpisodesCreated
		episodes = eps
		if len(eps) > 0 {
			mf.EpisodeID = &eps[0].ID
		}
	}

	now := time.Now().UTC()
	if mf.ID == 0 {
		mf.CreatedAt = now
		err = tx.CreateMediaFile(mf)
	} else {
		mf.UpdatedAt = now
		err = tx.UpdateMediaFile(mf)
	}
	if err != nil {
		return err
	}

	for _, ep := range episodes {
		link := &domain.MediaFileEpisode{
			MediaFileID: mf.ID,
			EpisodeID:   ep.ID,
		}
		if err := tx.CreateMediaFileEpisode(link); err != nil {
			return err
		}
	}

	if err := s.createAudioTracksTx(tx, mf, ffdata); err != nil {
		return err
	}

	return s.createSubtitleTracksTx(tx, mf, ffdata)
}

func (s *FSScanner) attachMovie(lib *domain.Library, mf *domain.MediaFile) (*attachResult, error) {

	title, year := guessMovieTitleAndYear(filepath.Base(mf.Path))

	existing, err := s.store.GetMovieByTitleAndYear(title, year, lib.ID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		mf.MovieID = &existing.ID
		return &attachResult{}, nil
	}

	m := &domain.Movie{
		LibraryID:     lib.ID,
		Title:         title,
		OriginalTitle: title,
		Year:          year,
		RuntimeMin:    mf.DurationSec / 60,
	}

	if err := s.store.CreateMovie(m); err != nil {
		return nil, err
	}

	mf.MovieID = &m.ID
	return &attachResult{MovieCreated: true}, nil
}

func (s *FSScanner) attachMovieTx(tx store.Store, lib *domain.Library, mf *domain.MediaFile) (*attachResult, error) {

	title, year := guessMovieTitleAndYear(filepath.Base(mf.Path))

	existing, err := tx.GetMovieByTitleAndYear(title, year, lib.ID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		mf.MovieID = &existing.ID
		return &attachResult{}, nil
	}

	m := &domain.Movie{
		LibraryID:     lib.ID,
		Title:         title,
		OriginalTitle: title,
		Year:          year,
		RuntimeMin:    mf.DurationSec / 60,
	}

	if err := tx.CreateMovie(m); err != nil {
		return nil, err
	}

	mf.MovieID = &m.ID
	return &attachResult{MovieCreated: true}, nil
}

// guessMovieTitleAndYear tries to parse "Title (2020).mkv" style names.
func guessMovieTitleAndYear(filename string) (string, int) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Example matches:
	// "Movie Title (2021)"
	// "Movie.Title.2021.1080p"
	re := regexp.MustCompile(`(?i)^(.*?)[\s\.\-_]*\(?((19|20)\d{2})\)?`)

	m := re.FindStringSubmatch(base)
	if len(m) >= 3 {
		rawTitle := strings.TrimSpace(m[1])
		yearStr := m[2]

		year, _ := strconv.Atoi(yearStr)
		title := normalizeTitle(rawTitle)
		return title, year
	}

	// Fallback: just normalize whole base.
	return normalizeTitle(base), 0
}

func normalizeTitle(s string) string {
	s = strings.ReplaceAll(s, ".", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.TrimSpace(s)
	return s
}

func (s *FSScanner) attachSeriesEpisode(
	lib *domain.Library,
	mf *domain.MediaFile,
) ([]*domain.Episode, *attachResult, error) {
	filename := filepath.Base(mf.Path)

	// 1) Range pattern: S01E01-02
	if m := reSxxExxRange.FindStringSubmatch(filename); len(m) == 4 {
		season, _ := strconv.Atoi(m[1])
		startEp, _ := strconv.Atoi(m[2])
		endEp, _ := strconv.Atoi(m[3])
		if endEp < startEp {
			endEp = startEp
		}
		return s.linkSeriesRange(lib, season, startEp, endEp, mf)
	}

	// 2) Single SxxExx
	if m := reSxxExx.FindStringSubmatch(filename); len(m) == 3 {
		season, _ := strconv.Atoi(m[1])
		episode, _ := strconv.Atoi(m[2])
		ep, ar, err := s.linkSeriesSingle(lib, season, episode, mf)
		if err != nil {
			return nil, nil, err
		}

		return []*domain.Episode{ep}, ar, nil

	}

	// 3) 1x02
	if m := reNxxN.FindStringSubmatch(filename); len(m) == 3 {
		season, _ := strconv.Atoi(m[1])
		episode, _ := strconv.Atoi(m[2])
		ep, ar, err := s.linkSeriesSingle(lib, season, episode, mf)
		if err != nil {
			return nil, nil, err
		}

		return []*domain.Episode{ep}, ar, nil

	}

	// 4) Anime-style
	if m := reAnimeEp.FindStringSubmatch(filename); len(m) == 3 {
		episode, _ := strconv.Atoi(m[2])
		ep, ar, err := s.linkSeriesSingle(lib, 1, episode, mf)
		if err != nil {
			return nil, nil, err
		}
		return []*domain.Episode{ep}, ar, nil

	}

	// 5) Fallback
	ep, ar, err := s.fallbackSeriesDetection(lib, mf)
	if err != nil {
		return nil, nil, err
	}

	return []*domain.Episode{ep}, ar, nil

}

func (s *FSScanner) attachSeriesEpisodeTx(
	tx store.Store,
	lib *domain.Library,
	mf *domain.MediaFile,
) ([]*domain.Episode, *attachResult, error) {

	filename := filepath.Base(mf.Path)

	// 1) Range pattern: S01E01-02
	if m := reSxxExxRange.FindStringSubmatch(filename); len(m) == 4 {
		season, _ := strconv.Atoi(m[1])
		startEp, _ := strconv.Atoi(m[2])
		endEp, _ := strconv.Atoi(m[3])
		if endEp < startEp {
			endEp = startEp
		}
		return s.linkSeriesRangeTx(tx, lib, season, startEp, endEp, mf)
	}

	// 2) Single SxxExx
	if m := reSxxExx.FindStringSubmatch(filename); len(m) == 3 {
		season, _ := strconv.Atoi(m[1])
		episode, _ := strconv.Atoi(m[2])

		ep, ar, err := s.linkSeriesSingleTx(tx, lib, season, episode, mf)
		if err != nil {
			return nil, nil, err
		}

		return []*domain.Episode{ep}, ar, nil
	}

	// 3) 1x02
	if m := reNxxN.FindStringSubmatch(filename); len(m) == 3 {
		season, _ := strconv.Atoi(m[1])
		episode, _ := strconv.Atoi(m[2])

		ep, ar, err := s.linkSeriesSingleTx(tx, lib, season, episode, mf)
		if err != nil {
			return nil, nil, err
		}

		return []*domain.Episode{ep}, ar, nil
	}

	// 4) Anime-style
	if m := reAnimeEp.FindStringSubmatch(filename); len(m) == 3 {
		episode, _ := strconv.Atoi(m[2])

		ep, ar, err := s.linkSeriesSingleTx(tx, lib, 1, episode, mf)
		if err != nil {
			return nil, nil, err
		}

		return []*domain.Episode{ep}, ar, nil
	}

	// 5) Fallback
	ep, ar, err := s.fallbackSeriesDetectionTx(tx, lib, mf)
	if err != nil {
		return nil, nil, err
	}

	return []*domain.Episode{ep}, ar, nil
}

// createAudioTracks stores embedded audio stream metadata.
func (s *FSScanner) createAudioTracks(mf *domain.MediaFile, ffdata *FFProbeOutput) error {
	for _, st := range ffdata.Streams {
		if st.CodecType != "audio" {
			continue
		}

		track := &domain.AudioTrack{
			MediaFileID: mf.ID,
			StreamIndex: st.Index,
			Language:    strings.TrimSpace(st.Tags.Language),
			Codec:       st.CodecName,
			Channels:    st.Channels,
			IsDefault:   st.Disposition.Default == 1,
		}

		if err := s.store.CreateAudioTrack(track); err != nil {
			return err
		}
	}

	return nil
}

func (s *FSScanner) createAudioTracksTx(tx store.Store, mf *domain.MediaFile, ffdata *FFProbeOutput) error {
	for _, st := range ffdata.Streams {
		if st.CodecType != "audio" {
			continue
		}

		track := &domain.AudioTrack{
			MediaFileID: mf.ID,
			StreamIndex: st.Index,
			Language:    strings.TrimSpace(st.Tags.Language),
			Codec:       st.CodecName,
			Channels:    st.Channels,
			IsDefault:   st.Disposition.Default == 1,
		}

		if err := tx.CreateAudioTrack(track); err != nil {
			return err
		}
	}

	return nil
}

// createSubtitleTracks stores both embedded and external subtitles.
func (s *FSScanner) createSubtitleTracks(mf *domain.MediaFile, ffdata *FFProbeOutput) error {
	// 1) Embedded subtitles from ffprobe.
	for _, st := range ffdata.Streams {
		if st.CodecType != "subtitle" {
			continue
		}

		language := strings.TrimSpace(st.Tags.Language)
		format := st.CodecName

		isDefault := st.Disposition.Default == 1
		isForced := st.Disposition.Forced == 1

		streamIndex := st.Index

		sub := &domain.SubtitleTrack{
			MediaFileID: mf.ID,
			Source:      domain.SubtitleSourceEmbedded,
			StreamIndex: &streamIndex,
			Language:    language,
			IsForced:    isForced,
			IsDefault:   isDefault,
			Format:      format,
		}

		if err := s.store.CreateSubtitleTrack(sub); err != nil {
			return err
		}
	}

	// 2) External sidecar subtitles (.srt, .ass, .vtt, .sub).
	externalSubs, err := findExternalSubtitles(mf.Path)
	if err != nil {
		return err
	}

	for _, ex := range externalSubs {
		lang := ex.Language
		format := strings.TrimPrefix(strings.ToLower(filepath.Ext(ex.Path)), ".")

		sub := &domain.SubtitleTrack{
			MediaFileID:  mf.ID,
			Source:       domain.SubtitleSourceExternal,
			ExternalPath: &ex.Path,
			Language:     lang,
			IsForced:     false,
			IsDefault:    false,
			Format:       format,
		}

		if err := s.store.CreateSubtitleTrack(sub); err != nil {
			return err
		}
	}

	return nil
}

func (s *FSScanner) createSubtitleTracksTx(tx store.Store, mf *domain.MediaFile, ffdata *FFProbeOutput) error {
	// 1) Embedded subtitles from ffprobe.
	for _, st := range ffdata.Streams {
		if st.CodecType != "subtitle" {
			continue
		}

		language := strings.TrimSpace(st.Tags.Language)
		format := st.CodecName

		isDefault := st.Disposition.Default == 1
		isForced := st.Disposition.Forced == 1

		streamIndex := st.Index

		sub := &domain.SubtitleTrack{
			MediaFileID: mf.ID,
			Source:      domain.SubtitleSourceEmbedded,
			StreamIndex: &streamIndex,
			Language:    language,
			IsForced:    isForced,
			IsDefault:   isDefault,
			Format:      format,
		}

		if err := tx.CreateSubtitleTrack(sub); err != nil {
			return err
		}
	}

	// 2) External sidecar subtitles (.srt, .ass, .vtt, .sub).
	externalSubs, err := findExternalSubtitles(mf.Path)
	if err != nil {
		return err
	}

	for _, ex := range externalSubs {
		lang := ex.Language
		format := strings.TrimPrefix(strings.ToLower(filepath.Ext(ex.Path)), ".")

		sub := &domain.SubtitleTrack{
			MediaFileID:  mf.ID,
			Source:       domain.SubtitleSourceExternal,
			ExternalPath: &ex.Path,
			Language:     lang,
			IsForced:     false,
			IsDefault:    false,
			Format:       format,
		}

		if err := tx.CreateSubtitleTrack(sub); err != nil {
			return err
		}
	}

	return nil
}

type externalSubtitle struct {
	Path     string
	Language string
}

var subtitleExt = map[string]bool{
	".srt": true,
	".ass": true,
	".vtt": true,
	".sub": true,
}

// findExternalSubtitles looks for files like:
//
//	Movie.mkv
//	Movie.srt
//	Movie.en.srt
//	Movie_es.srt
func findExternalSubtitles(videoPath string) ([]externalSubtitle, error) {
	dir := filepath.Dir(videoPath)
	base := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []externalSubtitle

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if !subtitleExt[ext] {
			continue
		}

		if !strings.HasPrefix(name, base) {
			continue
		}

		// Derive language from anything after base but before extension.
		// e.g. "Movie.en.srt" -> ".en"
		rest := strings.TrimSuffix(name, ext)
		rest = strings.TrimPrefix(rest, base)

		rest = strings.TrimLeft(rest, "._-")
		lang := strings.TrimSpace(rest) // may be "" if plain Movie.srt

		result = append(result, externalSubtitle{
			Path:     filepath.Join(dir, name),
			Language: lang,
		})
	}

	return result, nil
}

// linkSeriesSingle: one season/episode → returns the Episode.
func (s *FSScanner) linkSeriesSingle(
	lib *domain.Library,
	season, episode int,
	mf *domain.MediaFile,
) (*domain.Episode, *attachResult, error) {

	res := &attachResult{}
	title := extractSeriesTitle(mf.Path)

	// 1) Series
	sr, err := s.store.GetSeriesByTitle(title, lib.ID)
	if err != nil {
		return nil, nil, err
	}
	if sr == nil {
		sr = &domain.Series{
			LibraryID: lib.ID,
			Title:     title,
		}
		if err := s.store.CreateSeries(sr); err != nil {
			return nil, nil, err
		}
		res.SeriesCreated = true
	}

	// 2) Season
	se, err := s.store.GetSeasonBySeriesAndNumber(sr.ID, season)
	if err != nil {
		return nil, nil, err
	}
	if se == nil {
		se = &domain.Season{
			SeriesID: sr.ID,
			Number:   season,
		}
		if err := s.store.CreateSeason(se); err != nil {
			return nil, nil, err
		}
		res.SeasonCreated = true
	}

	// 3) Episode
	ep, err := s.store.GetEpisodeBySeasonAndNumber(se.ID, episode)
	if err != nil {
		return nil, nil, err
	}
	if ep == nil {
		ep = &domain.Episode{
			SeasonID: se.ID,
			Number:   episode,
		}
		if err := s.store.CreateEpisode(ep); err != nil {
			return nil, nil, err
		}
		res.EpisodesCreated = 1
	}

	return ep, res, nil
}

func (s *FSScanner) linkSeriesSingleTx(
	tx store.Store,
	lib *domain.Library,
	season, episode int,
	mf *domain.MediaFile,
) (*domain.Episode, *attachResult, error) {

	res := &attachResult{}
	title := extractSeriesTitle(mf.Path)

	// 1) Series
	sr, err := tx.GetSeriesByTitle(title, lib.ID)
	if err != nil {
		return nil, nil, err
	}
	if sr == nil {
		sr = &domain.Series{
			LibraryID: lib.ID,
			Title:     title,
		}
		if err := tx.CreateSeries(sr); err != nil {
			return nil, nil, err
		}
		res.SeriesCreated = true
	}

	// 2) Season
	se, err := tx.GetSeasonBySeriesAndNumber(sr.ID, season)
	if err != nil {
		return nil, nil, err
	}
	if se == nil {
		se = &domain.Season{
			SeriesID: sr.ID,
			Number:   season,
		}
		if err := tx.CreateSeason(se); err != nil {
			return nil, nil, err
		}
		res.SeasonCreated = true
	}

	// 3) Episode
	ep, err := tx.GetEpisodeBySeasonAndNumber(se.ID, episode)
	if err != nil {
		return nil, nil, err
	}
	if ep == nil {
		ep = &domain.Episode{
			SeasonID: se.ID,
			Number:   episode,
		}
		if err := tx.CreateEpisode(ep); err != nil {
			return nil, nil, err
		}
		res.EpisodesCreated = 1
	}

	return ep, res, nil
}

// linkSeriesRange: S01E01-02 → creates/returns episodes 1 and 2.
func (s *FSScanner) linkSeriesRange(
	lib *domain.Library,
	season, startEp, endEp int,
	mf *domain.MediaFile,
) ([]*domain.Episode, *attachResult, error) {

	res := &attachResult{}
	var eps []*domain.Episode

	for n := startEp; n <= endEp; n++ {
		ep, ar, err := s.linkSeriesSingle(lib, season, n, mf)
		if err != nil {
			return nil, nil, err
		}

		if ar.SeriesCreated {
			res.SeriesCreated = true
		}
		if ar.SeasonCreated {
			res.SeasonCreated = true
		}
		res.EpisodesCreated += ar.EpisodesCreated

		eps = append(eps, ep)
	}

	return eps, res, nil
}

func (s *FSScanner) linkSeriesRangeTx(
	tx store.Store,
	lib *domain.Library,
	season, startEp, endEp int,
	mf *domain.MediaFile,
) ([]*domain.Episode, *attachResult, error) {

	res := &attachResult{}
	var eps []*domain.Episode

	for n := startEp; n <= endEp; n++ {
		ep, ar, err := s.linkSeriesSingleTx(tx, lib, season, n, mf)
		if err != nil {
			return nil, nil, err
		}

		if ar.SeriesCreated {
			res.SeriesCreated = true
		}
		if ar.SeasonCreated {
			res.SeasonCreated = true
		}
		res.EpisodesCreated += ar.EpisodesCreated

		eps = append(eps, ep)
	}

	return eps, res, nil
}

func extractSeriesTitle(path string) string {
	dir := filepath.Dir(path)
	parts := strings.Split(dir, string(os.PathSeparator))

	if len(parts) == 0 {
		return ""
	}

	// Get last folder
	last := parts[len(parts)-1]

	// If last folder is a season folder, remove it
	if isSeasonFolder(last) {
		parts = parts[:len(parts)-1]
	}

	if len(parts) == 0 {
		return ""
	}

	title := parts[len(parts)-1]

	title = strings.ReplaceAll(title, ".", " ")
	title = strings.ReplaceAll(title, "_", " ")
	return strings.TrimSpace(title)
}

func isSeasonFolder(name string) bool {
	nameLower := strings.ToLower(name)

	// Common season folder patterns:
	if strings.HasPrefix(nameLower, "season ") {
		return true
	}
	if strings.HasPrefix(nameLower, "season") {
		return true
	}
	if strings.HasPrefix(nameLower, "s") && len(nameLower) <= 3 {
		// S01, S1, s02, etc.
		return true
	}

	// Numeric folders: "1", "01", "02"
	if m, _ := regexp.MatchString(`^\d{1,2}$`, nameLower); m {
		return true
	}

	return false
}

func (s *FSScanner) fallbackSeriesDetection(
	lib *domain.Library,
	mf *domain.MediaFile,
) (*domain.Episode, *attachResult, error) {
	epNum := guessEpisodeNumber(filepath.Base(mf.Path))

	ep, ar, err := s.linkSeriesSingle(lib, 1, epNum, mf)
	if err != nil {
		return nil, nil, err
	}

	return ep, ar, nil

}

func (s *FSScanner) fallbackSeriesDetectionTx(
	tx store.Store,
	lib *domain.Library,
	mf *domain.MediaFile,
) (*domain.Episode, *attachResult, error) {

	epNum := guessEpisodeNumber(filepath.Base(mf.Path))

	ep, ar, err := s.linkSeriesSingleTx(tx, lib, 1, epNum, mf)
	if err != nil {
		return nil, nil, err
	}

	return ep, ar, nil
}

func guessEpisodeNumber(filename string) int {
	re := regexp.MustCompile(`\b(\d{1,4})\b`)
	m := re.FindStringSubmatch(filename)
	if len(m) == 2 {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 1
}
