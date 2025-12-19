package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/metadata"
	"github.com/bastianvv/vio/internal/store"
)

type SeriesHandler struct {
	store    store.Store
	metadata metadata.Enricher
}

func NewSeriesHandler(s store.Store, metadata metadata.Enricher) *SeriesHandler {
	return &SeriesHandler{store: s, metadata: metadata}
}

func (h *SeriesHandler) ListSeries(w http.ResponseWriter, r *http.Request) {
	seriesList, err := h.store.ListSeries()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// seriesList is []*domain.Series
	out := make([]*dto.Series, 0, len(seriesList))
	for _, sr := range seriesList {
		out = append(out, dto.NewSeries(sr)) // sr is already *domain.Series
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

	writeJSON(w, dto.NewSeries(sr))
}

func (h *SeriesHandler) ListSeasonsBySeries(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	seasons, err := h.store.ListSeasonsBySeries(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// seasons is []domain.Season (value slice!)
	out := make([]*dto.Season, 0, len(seasons))
	for i := range seasons {
		out = append(out, dto.NewSeason(&seasons[i]))
	}

	writeJSON(w, out)
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
		"series":   dto.NewSeries(series),
	})
}
