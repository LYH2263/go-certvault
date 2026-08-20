package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/certvault"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	switch {
	case errors.Is(err, certvault.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, certvault.ErrClosed):
		code = http.StatusServiceUnavailable
	case errors.Is(err, certvault.ErrNoSigner):
		code = http.StatusFailedDependency
	}
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

func readJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
