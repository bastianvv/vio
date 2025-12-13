package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type FilesHandler struct {
	store store.Store
}

func NewFilesHandler(s store.Store) *FilesHandler {
	return &FilesHandler{store: s}
}

func (h *FilesHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	f, err := h.store.GetMediaFile(id)
	if err != nil || f == nil {
		http.Error(w, "file not found", 404)
		return
	}

	writeJSON(w, dto.NewMediaFile(f))
}

func (h *FilesHandler) StreamFile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	mf, err := h.store.GetMediaFile(id)
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}

	f, err := os.Open(mf.Path)
	if err != nil {
		http.Error(w, "cannot open file", 500)
		return
	}
	defer func() { _ = f.Close() }()

	w.Header().Set("Content-Type", "video/"+mf.Container)
	http.ServeContent(w, r, mf.Path, mf.UpdatedAt, f)
}
