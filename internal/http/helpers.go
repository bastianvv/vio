package http

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type apiError struct {
	Error string `json:"error"`
}

func parseID(r *http.Request, key string) (int64, error) {
	idStr := r.PathValue(key) // Go 1.22+ mux
	return strconv.ParseInt(idStr, 10, 64)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiError{
		Error: err.Error(),
	})
}
