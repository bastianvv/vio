package http

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/store"
)

type MoviesHandler struct {
	store        store.Store
	metadata     metadata.Enricher
	imageBaseDir string
}

func NewMoviesHandler(
	store store.Store,
	metadata metadata.Enricher,
	imageBaseDir string,
) *MoviesHandler {
	return &MoviesHandler{
		store:        store,
		metadata:     metadata,
		imageBaseDir: imageBaseDir,
	}
}

// GET /api/movies?library_id=X
func (h *MoviesHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	lidStr := r.URL.Query().Get("library_id")
	if lidStr == "" {
		http.Error(w, "library_id required", http.StatusBadRequest)
		return
	}
	lid, err := strconv.ParseInt(lidStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid library_id", http.StatusBadRequest)
		return
	}

	movies, err := h.store.ListMoviesByLibrary(lid)
	if err != nil {
		http.Error(w, "failed to list movies", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.Movie, 0, len(movies))
	for i := range movies {
		m := &movies[i]
		out = append(out, dto.NewMovie(
			m,
			h.movieHasPoster(m.ID),
			h.movieHasBackdrop(m.ID),
		))
	}

	writeJSON(w, out)
}

// GET /api/movies/{id}
func (h *MoviesHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	movie, err := h.store.GetMovie(id)
	if err != nil || movie == nil {
		http.Error(w, "movie not found", http.StatusNotFound)
		return
	}

	writeJSON(w, dto.NewMovie(
		movie,
		h.movieHasPoster(movie.ID),
		h.movieHasBackdrop(movie.ID),
	))

}

// GET /api/movies/{id}/files
func (h *MoviesHandler) ListMovieFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	files, err := h.store.ListMediaFilesByMovie(id)
	if err != nil {
		http.Error(w, "failed to list movie files", 500)
		return
	}

	out := make([]*dto.MediaFile, 0, len(files))
	for _, f := range files {
		out = append(out, dto.NewMediaFile(&f))
	}

	writeJSON(w, out)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func (h *MoviesHandler) ListMediaFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	movieID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	files, err := h.store.ListMediaFilesByMovie(movieID)
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.MediaFile, 0, len(files))
	for i := range files {
		out = append(out, dto.NewMediaFile(&files[i]))
	}

	writeJSON(w, out)
}

func (h *MoviesHandler) EnrichMovie(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.metadata.EnrichMovie(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	movie, err := h.store.GetMovie(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]any{
		"status":   "enriched",
		"provider": "tmdb",
		"movie":    dto.NewMovie(movie, h.movieHasPoster(movie.ID), h.movieHasBackdrop(movie.ID)),
	})

}

func (h *MoviesHandler) movieHasPoster(id int64) bool {
	_, err := os.Stat(filepath.Join(
		h.imageBaseDir,
		"movies",
		strconv.FormatInt(id, 10),
		"poster.jpg",
	))
	return err == nil
}

func (h *MoviesHandler) movieHasBackdrop(id int64) bool {
	_, err := os.Stat(filepath.Join(
		h.imageBaseDir,
		"movies",
		strconv.FormatInt(id, 10),
		"backdrop.jpg",
	))
	return err == nil
}
