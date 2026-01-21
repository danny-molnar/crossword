package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/danny-molnar/crossword/internal/api/handlers"
	"github.com/danny-molnar/crossword/internal/store"
	"github.com/danny-molnar/crossword/internal/tools"
)

func NewRouter(st *store.MemoryStore, wl *tools.Wordlist) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.New(st, wl)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		r.Get("/puzzles/{id}", h.GetPuzzle)

		r.Post("/puzzles/{id}/sessions", h.CreateSession)
		r.Get("/sessions/{sid}", h.GetSession)
		r.Put("/sessions/{sid}", h.UpdateSession)

		r.Get("/tools/anagram", h.Anagram)
		r.Get("/tools/pattern", h.Pattern)
	})

	return r
}
