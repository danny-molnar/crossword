package store

import (
	"fmt"
	"sync"

	"github.com/danny-molnar/crossword/internal/domain"
)

type PuzzleStore struct {
	mu      sync.RWMutex
	puzzles map[string]domain.Puzzle
}

func NewPuzzleStore() *PuzzleStore {
	return &PuzzleStore{
		puzzles: make(map[string]domain.Puzzle),
	}
}

func (s *PuzzleStore) PutPuzzle(p domain.Puzzle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.puzzles[p.ID] = p
}

func (s *PuzzleStore) GetPuzzle(id string) (domain.Puzzle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.puzzles[id]
	if !ok {
		return domain.Puzzle{}, fmt.Errorf("puzzle not found")
	}
	return p, nil
}
