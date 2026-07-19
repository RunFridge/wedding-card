package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type gameEntry struct {
	createdAt time.Time
	expiresAt time.Time
}

type GameStore struct {
	mu     sync.Mutex
	tokens map[string]gameEntry
	ttl    time.Duration
}

var GameSessions *GameStore

func InitGame() {
	GameSessions = &GameStore{
		tokens: make(map[string]gameEntry),
		ttl:    2 * time.Minute,
	}
}

func NewGameStore(ttl time.Duration) *GameStore {
	return &GameStore{
		tokens: make(map[string]gameEntry),
		ttl:    ttl,
	}
}

func (s *GameStore) Create() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpired()
	s.tokens[token] = gameEntry{createdAt: now, expiresAt: now.Add(s.ttl)}
	return token, nil
}

func (s *GameStore) Consume(token string, minElapsedMs int) (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.tokens[token]
	if !ok {
		return false, "invalid token"
	}
	delete(s.tokens, token)

	if time.Now().After(entry.expiresAt) {
		return false, "token expired"
	}
	if time.Since(entry.createdAt).Milliseconds() < int64(minElapsedMs) {
		return false, "completed too fast"
	}
	return true, ""
}

func (s *GameStore) purgeExpired() {
	now := time.Now()
	for token, entry := range s.tokens {
		if now.After(entry.expiresAt) {
			delete(s.tokens, token)
		}
	}
}
