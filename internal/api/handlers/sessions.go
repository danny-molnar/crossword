package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/danny-molnar/crossword/internal/domain"
	"github.com/danny-molnar/crossword/internal/util"
)

type createSessionResponse struct {
	Session domain.SolveSession `json:"session"`
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	puzzleID := chi.URLParam(r, "id")

	// Ensure puzzle exists
	if _, err := h.store.Puzzles.GetPuzzle(puzzleID); err != nil {
		writeErr(w, http.StatusNotFound, "puzzle not found")
		return
	}

	now := time.Now().UTC()
	sess := domain.SolveSession{
		ID:        util.NewID(),
		PuzzleID:  puzzleID,
		CreatedAt: now,
		UpdatedAt: now,
		GridState: map[string]string{},
		Pencil:    map[string]bool{},
	}

	h.store.Sessions.Create(sess)
	writeJSON(w, http.StatusCreated, createSessionResponse{Session: sess})
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "sid")

	sess, err := h.store.Sessions.Get(sid)
	if err != nil {
		writeErr(w, http.StatusNotFound, "session not found")
		return
	}

	writeJSON(w, http.StatusOK, sess)
}

type updateSessionRequest struct {
	GridState map[string]string `json:"gridState"`
	Pencil    map[string]bool   `json:"pencil"`
}

func (h *Handler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "sid")

	var req updateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	updated, err := h.store.Sessions.Update(sid, func(cur domain.SolveSession) domain.SolveSession {
		if req.GridState != nil {
			cur.GridState = req.GridState
		}
		if req.Pencil != nil {
			cur.Pencil = req.Pencil
		}
		return cur
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, "session not found")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
