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

type SeasonsHandler struct {
	store        store.Store
	imageBaseDir string
}

func NewSeasonsHandler(s store.Store, imageBaseDir string) *SeasonsHandler {
	return &SeasonsHandler{store: s, imageBaseDir: imageBaseDir}
}

func (h *SeasonsHandler) GetSeason(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	se, err := h.store.GetSeason(id)
	if err != nil || se == nil {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, dto.NewSeason(se, h.seasonHasPoster(se.ID)))
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
		se := &seasons[i]
		out = append(out, dto.NewSeason(se, h.seasonHasPoster(se.ID)))
	}

	writeJSON(w, out)
}

func (h *SeasonsHandler) seasonHasPoster(seasonID int64) bool {
	dir := filepath.Join(
		h.imageBaseDir,
		"seasons",
		strconv.FormatInt(seasonID, 10),
	)
	candidates := []string{
		"poster.webp",
		"poster.jpg",
		"poster.jpeg",
		"poster.png",
	}

	for _, name := range candidates {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}

	return false
}
