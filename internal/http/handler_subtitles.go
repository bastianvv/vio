package http

import (
	"net/http"
	"strconv"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type SubtitlesHandler struct {
	store store.Store
}

func NewSubtitlesHandler(s store.Store) *SubtitlesHandler {
	return &SubtitlesHandler{store: s}
}

func (h *SubtitlesHandler) ListSubtitleTracks(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	mediaFileID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid media file id", http.StatusBadRequest)
		return
	}

	subs, err := h.store.ListSubtitleTracks(mediaFileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]*dto.SubtitleTrack, 0, len(subs))
	for i := range subs {
		out = append(out, dto.NewSubtitleTrack(&subs[i]))
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, out)
}
