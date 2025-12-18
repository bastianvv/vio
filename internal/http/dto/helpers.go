package dto

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func parseID(r *http.Request) int64 {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}
