package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bastianvv/vio/internal/domain"
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
		streamURL := "/api/subtitles/" + strconv.FormatInt(subs[i].ID, 10) + "/stream"
		out = append(out, dto.NewSubtitleTrack(&subs[i], streamURL))
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, out)
}

func (h *SubtitlesHandler) StreamSubtitleTrack(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	subtitleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid subtitle id", http.StatusBadRequest)
		return
	}

	st, err := h.store.GetSubtitleTrack(subtitleID)
	if err != nil {
		http.Error(w, "failed to load subtitle track", http.StatusInternalServerError)
		return
	}
	if st == nil {
		http.Error(w, "subtitle track not found", http.StatusNotFound)
		return
	}

	if st.Source != domain.SubtitleSourceExternal || st.ExternalPath == nil {
		http.Error(w, "embedded subtitle streaming not implemented", http.StatusNotImplemented)
		return
	}

	f, err := os.Open(*st.ExternalPath)
	if err != nil {
		http.Error(w, "cannot open subtitle file", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "cannot stat subtitle file", http.StatusInternalServerError)
		return
	}

	http.ServeContent(
		w,
		r,
		filepath.Base(*st.ExternalPath),
		stat.ModTime(),
		f,
	)
}
