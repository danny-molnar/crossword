package store

type MemoryStore struct {
	Puzzles  *PuzzleStore
	Sessions *SessionStore
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Puzzles:  NewPuzzleStore(),
		Sessions: NewSessionStore(),
	}
}
