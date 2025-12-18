package http

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/bastianvv/vio/internal/domain"
	"github.com/bastianvv/vio/internal/http/dto"
	"github.com/bastianvv/vio/internal/media"
	"github.com/bastianvv/vio/internal/scan"
	"github.com/bastianvv/vio/internal/store"
	"github.com/go-chi/chi/v5"
)

type LibrariesHandler struct {
	store   store.Store
	scanner media.Scanner
	scans   *scan.Registry
}

type CreateLibraryRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
}

type UpdateLibraryRequest struct {
	Name string `json:"name"`
}

func NewLibrariesHandler(s store.Store, sc media.Scanner, scans *scan.Registry) *LibrariesHandler {
	return &LibrariesHandler{
		store:   s,
		scanner: sc,
		scans:   scans,
	}
}

func (h *LibrariesHandler) ListLibraries(w http.ResponseWriter, r *http.Request) {
	libs, err := h.store.ListLibraries()
	if err != nil {
		http.Error(w, "failed to list libraries", http.StatusInternalServerError)
		return
	}

	out := make([]*dto.Library, 0, len(libs))
	for i := range libs {
		out = append(out, dto.NewLibrary(&libs[i]))
	}
	writeJSON(w, out)
}

func (h *LibrariesHandler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid library id", http.StatusBadRequest)
		return
	}

	lib, err := h.store.GetLibrary(id)
	if err != nil {
		http.Error(w, "library not found", http.StatusNotFound)
		return
	}

	writeJSON(w, dto.NewLibrary(lib))
}

func (h *LibrariesHandler) CreateLibrary(w http.ResponseWriter, r *http.Request) {
	var req CreateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Path == "" || req.Type == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// optional but recommended
	if _, err := os.Stat(req.Path); err != nil {
		http.Error(w, "library path does not exist", http.StatusBadRequest)
		return
	}

	lib := &domain.Library{
		Name: req.Name,
		Type: domain.LibraryType(req.Type),
		Path: req.Path,
	}

	if err := h.store.CreateLibrary(lib); err != nil {
		http.Error(w, "failed to create library", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, dto.NewLibrary(lib))
}

func (h *LibrariesHandler) UpdateLibrary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid library id", http.StatusBadRequest)
		return
	}

	lib, err := h.store.GetLibrary(id)
	if err != nil {
		http.Error(w, "library not found", http.StatusNotFound)
		return
	}

	var req UpdateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != "" {
		lib.Name = req.Name
	}

	if err := h.store.UpdateLibrary(lib); err != nil {
		http.Error(w, "failed to update library", http.StatusInternalServerError)
		return
	}

	writeJSON(w, dto.NewLibrary(lib))
}

func (h *LibrariesHandler) ScanLibrary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid library id", http.StatusBadRequest)
		return
	}

	lib, err := h.store.GetLibrary(id)
	if err != nil || lib == nil {
		http.Error(w, "library not found", http.StatusNotFound)
		return
	}

	job := h.scans.Start(id)

	go func(lib *domain.Library, jobID string) {
		_, err := h.scanner.ScanLibrary(lib, media.ScanModeIncremental)
		if err != nil {
			h.scans.Fail(jobID, err)
			return
		}
		h.scans.Finish(jobID)
	}(lib, job.ID)

	writeJSON(w, map[string]any{
		"job_id": job.ID,
		"status": job.Status,
	})
}

func (h *LibrariesHandler) RescanLibrary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid library id", http.StatusBadRequest)
		return
	}

	lib, err := h.store.GetLibrary(id)
	if err != nil || lib == nil {
		http.Error(w, "library not found", http.StatusNotFound)
		return
	}

	job := h.scans.Start(id)

	go func(lib *domain.Library, jobID string) {
		_, err := h.scanner.ScanLibrary(lib, media.ScanModeRescan)
		if err != nil {
			h.scans.Fail(jobID, err)
			return
		}
		h.scans.Finish(jobID)
	}(lib, job.ID)

	writeJSON(w, map[string]any{
		"job_id": job.ID,
		"status": job.Status,
	})
}

func (h *LibrariesHandler) GetScanJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")

	job, ok := h.scans.Get(jobID)
	if !ok {
		http.Error(w, "scan job not found", http.StatusNotFound)
		return
	}

	writeJSON(w, job)
}
