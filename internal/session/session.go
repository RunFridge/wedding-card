package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type entry struct {
	expiresAt time.Time
}

type Store struct {
	mu      sync.RWMutex
	tokens  map[string]entry
	ttl     time.Duration
}

var Global *Store

func Init() {
	Global = NewStore(24 * time.Hour)
}

func NewStore(ttl time.Duration) *Store {
	return &Store{
		tokens: make(map[string]entry),
		ttl:    ttl,
	}
}

func (s *Store) Create() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked()
	s.tokens[token] = entry{expiresAt: time.Now().Add(s.ttl)}
	return token, nil
}

func (s *Store) Valid(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.tokens[token]
	return ok && time.Now().Before(e.expiresAt)
}

func (s *Store) Revoke(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}

func (s *Store) purgeExpiredLocked() {
	now := time.Now()
	for k, e := range s.tokens {
		if now.After(e.expiresAt) {
			delete(s.tokens, k)
		}
	}
}
