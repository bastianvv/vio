package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type EpisodesHandler struct {
	store        store.Store
	imageBaseDir string
}

func NewEpisodesHandler(s store.Store, imageBaseDir string) *EpisodesHandler {
	return &EpisodesHandler{store: s, imageBaseDir: imageBaseDir}
}

func (h *EpisodesHandler) GetEpisode(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ep, err := h.store.GetEpisode(id)
	if err != nil || ep == nil {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, dto.NewEpisode(ep, h.episodeHasStill(ep.ID)))
}

func (h *EpisodesHandler) ListEpisodesBySeason(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	seasonID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid season id", http.StatusBadRequest)
		return
	}

	episodes, err := h.store.ListEpisodesBySeason(seasonID)
	if err != nil {
		http.Error(w, "failed to list episodes", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.Episode, 0, len(episodes))
	for i := range episodes {
		ep := &episodes[i]
		out = append(out, dto.NewEpisode(ep, h.episodeHasStill(ep.ID)))
	}

	writeJSON(w, out)
}

func (h *EpisodesHandler) ListEpisodeFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	episodeID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid episode id", http.StatusBadRequest)
		return
	}

	files, err := h.store.ListMediaFilesByEpisode(episodeID)
	if err != nil {
		http.Error(w, "failed to list episode files", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.MediaFile, 0, len(files))
	for i := range files {
		out = append(out, dto.NewMediaFile(&files[i]))
	}

	writeJSON(w, out)
}

func (h *EpisodesHandler) episodeHasStill(id int64) bool {
	dir := filepath.Join(h.imageBaseDir, "episodes", strconv.FormatInt(id, 10))
	_, err := os.Stat(filepath.Join(dir, "still.jpg"))
	return err == nil
}
