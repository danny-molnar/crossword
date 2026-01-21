package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/danny-molnar/crossword/internal/domain"
)

func (h *Handler) GetPuzzle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.store.Puzzles.GetPuzzle(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "puzzle not found")
		return
	}

	// Always return the public view (no answers / solutions)
	pub := domain.ToPublic(p)
	writeJSON(w, http.StatusOK, pub)
}
