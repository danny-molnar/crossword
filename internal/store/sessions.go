package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/danny-molnar/crossword/internal/domain"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]domain.SolveSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]domain.SolveSession),
	}
}

func (s *SessionStore) Create(sess domain.SolveSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
}

func (s *SessionStore) Get(id string) (domain.SolveSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sessions[id]
	if !ok {
		return domain.SolveSession{}, fmt.Errorf("session not found")
	}
	return v, nil
}

func (s *SessionStore) Update(id string, update func(domain.SolveSession) domain.SolveSession) (domain.SolveSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cur, ok := s.sessions[id]
	if !ok {
		return domain.SolveSession{}, fmt.Errorf("session not found")
	}

	cur = update(cur)
	cur.UpdatedAt = time.Now().UTC()
	s.sessions[id] = cur
	return cur, nil
}
