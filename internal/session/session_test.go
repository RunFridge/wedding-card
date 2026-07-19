package session

import (
	"testing"
	"time"
)

func TestCreateAndValid(t *testing.T) {
	s := NewStore(time.Hour)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if len(token) != 64 {
		t.Fatalf("expected 64-char hex token, got %d chars", len(token))
	}
	if !s.Valid(token) {
		t.Fatal("token should be valid")
	}
}

func TestRejectUnknown(t *testing.T) {
	s := NewStore(time.Hour)
	if s.Valid("nonexistent") {
		t.Fatal("unknown token should be invalid")
	}
}

func TestExpiry(t *testing.T) {
	s := NewStore(time.Millisecond)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if s.Valid(token) {
		t.Fatal("expired token should be invalid")
	}
}

func TestRevoke(t *testing.T) {
	s := NewStore(time.Hour)
	token, _ := s.Create()
	s.Revoke(token)
	if s.Valid(token) {
		t.Fatal("revoked token should be invalid")
	}
}

func TestPurge(t *testing.T) {
	s := NewStore(time.Millisecond)
	s.Create()
	s.Create()
	time.Sleep(5 * time.Millisecond)

	s.Create()

	s.mu.RLock()
	count := len(s.tokens)
	s.mu.RUnlock()
	if count != 1 {
		t.Fatalf("expected 1 token after purge, got %d", count)
	}
}
