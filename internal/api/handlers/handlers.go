package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danny-molnar/crossword/internal/store"
	"github.com/danny-molnar/crossword/internal/tools"
)

type Handler struct {
	store *store.MemoryStore
	wl    *tools.Wordlist
}

func New(st *store.MemoryStore, wl *tools.Wordlist) *Handler {
	return &Handler{
		store: st,
		wl:    wl,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}
