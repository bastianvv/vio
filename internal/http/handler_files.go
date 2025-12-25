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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mf, err := h.store.GetMediaFile(id)
	if err != nil || mf == nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(mf.Path)
	if err != nil {
		http.Error(w, "cannot open file", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "cannot stat file", http.StatusInternalServerError)
		return
	}

	// IMPORTANT: do NOT set Content-Type manually
	http.ServeContent(
		w,
		r,
		filepath.Base(mf.Path), // correct name
		stat.ModTime(),         // better than mf.UpdatedAt
		f,
	)
}
