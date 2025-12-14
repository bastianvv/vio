package http

import (
	"net/http"
	"strconv"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type EpisodesHandler struct {
	store store.Store
}

func NewEpisodesHandler(s store.Store) *EpisodesHandler {
	return &EpisodesHandler{store: s}
}

func (h *EpisodesHandler) GetEpisode(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	ep, err := h.store.GetEpisode(id)
	if err != nil || ep == nil {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, dto.NewEpisode(ep))
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
		out = append(out, dto.NewEpisode(&episodes[i]))
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
