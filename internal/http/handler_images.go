package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	basePath string
}

func NewImageHandler(basePath string) *ImageHandler {
	return &ImageHandler{basePath: basePath}
}

func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity") // movies, series, seasons, episodes
	id := chi.URLParam(r, "id")
	kind := chi.URLParam(r, "kind") // poster, backdrop, still

	// Directory: data/images/{entity}/{id}
	dir := filepath.Join(h.basePath, entity, id)

	// Try known extensions (order matters)
	candidates := []string{
		kind + ".webp",
		kind + ".jpg",
		kind + ".jpeg",
		kind + ".png",
	}

	var filePath string
	for _, name := range candidates {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			filePath = p
			break
		}
	}

	if filePath == "" {
		http.NotFound(w, r)
		return
	}

	// Let net/http set Content-Type correctly
	http.ServeFile(w, r, filePath)
}
