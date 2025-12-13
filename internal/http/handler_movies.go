package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
)

type MoviesHandler struct {
	store store.Store
}

func NewMoviesHandler(s store.Store) *MoviesHandler {
	return &MoviesHandler{store: s}
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
	for _, m := range movies {
		out = append(out, dto.NewMovie(&m))
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

	writeJSON(w, dto.NewMovie(movie))
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
