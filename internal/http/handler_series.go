package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/store"
)

type SeriesHandler struct {
	store        store.Store
	metadata     metadata.Enricher
	imageBaseDir string
}

func NewSeriesHandler(s store.Store, metadata metadata.Enricher, imageBaseDir string) *SeriesHandler {
	return &SeriesHandler{store: s, metadata: metadata, imageBaseDir: imageBaseDir}
}

func (h *SeriesHandler) ListSeries(w http.ResponseWriter, r *http.Request) {
	seriesList, err := h.store.ListSeries()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// seriesList is []*domain.Series
	out := make([]*dto.Series, 0, len(seriesList))
	for _, s := range seriesList {
		out = append(out, dto.NewSeries(
			s,
			h.seriesHasPoster(s.ID),
			h.seriesHasBackdrop(s.ID),
		))
	}

	writeJSON(w, out)
}

func (h *SeriesHandler) GetSeries(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	sr, err := h.store.GetSeries(id)
	if err != nil || sr == nil {
		http.Error(w, "not found", 404)
		return
	}

	writeJSON(w, dto.NewSeries(
		sr,
		h.seriesHasPoster(sr.ID),
		h.seriesHasBackdrop(sr.ID),
	))
}

func (h *SeriesHandler) EnrichSeries(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.metadata.EnrichSeries(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	series, err := h.store.GetSeries(id)
	if err != nil || series == nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]any{
		"status":   "enriched",
		"provider": "tmdb",
		"series": dto.NewSeries(
			series,
			h.seriesHasPoster(series.ID),
			h.seriesHasBackdrop(series.ID),
		),
	})
}

func (h *SeriesHandler) seriesHasPoster(id int64) bool {
	_, err := os.Stat(filepath.Join(
		h.imageBaseDir,
		"series",
		strconv.FormatInt(id, 10),
		"poster.jpg",
	))
	return err == nil
}

func (h *SeriesHandler) seriesHasBackdrop(id int64) bool {
	_, err := os.Stat(filepath.Join(
		h.imageBaseDir,
		"series",
		strconv.FormatInt(id, 10),
		"backdrop.jpg",
	))
	return err == nil
}
