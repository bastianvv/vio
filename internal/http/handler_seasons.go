package http

import (
	"net/http"
	"strconv"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type SeasonsHandler struct {
	store store.Store
}

func NewSeasonsHandler(s store.Store) *SeasonsHandler {
	return &SeasonsHandler{store: s}
}

func (h *SeasonsHandler) GetSeason(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	se, err := h.store.GetSeason(id)
	if err != nil || se == nil {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, dto.NewSeason(se))
}

func (h *SeasonsHandler) ListEpisodesBySeason(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	episodes, err := h.store.ListEpisodesBySeason(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out := make([]*dto.Episode, 0, len(episodes))
	for i := range episodes {
		out = append(out, dto.NewEpisode(&episodes[i]))
	}

	writeJSON(w, out)
}
func (h *SeasonsHandler) ListSeasonsBySeries(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	seriesID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid series id", http.StatusBadRequest)
		return
	}

	seasons, err := h.store.ListSeasonsBySeries(seriesID)
	if err != nil {
		http.Error(w, "failed to list seasons", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.Season, 0, len(seasons))
	for i := range seasons {
		out = append(out, dto.NewSeason(&seasons[i]))
	}

	writeJSON(w, out)
}
