package session

import (
	"testing"
	"time"
)

func TestGameCreateAndConsume(t *testing.T) {
	s := NewGameStore(time.Hour)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if len(token) != 64 {
		t.Fatalf("expected 64-char hex token, got %d chars", len(token))
	}
	ok, reason := s.Consume(token, 0)
	if !ok {
		t.Fatalf("Consume should succeed, got reason: %s", reason)
	}
}

func TestGameConsumeUnknown(t *testing.T) {
	s := NewGameStore(time.Hour)
	ok, reason := s.Consume("nonexistent", 0)
	if ok {
		t.Fatal("unknown token should fail")
	}
	if reason != "invalid token" {
		t.Fatalf("expected reason %q, got %q", "invalid token", reason)
	}
}

func TestGameConsumeExpired(t *testing.T) {
	s := NewGameStore(time.Millisecond)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	ok, reason := s.Consume(token, 0)
	if ok {
		t.Fatal("expired token should fail")
	}
	if reason != "token expired" {
		t.Fatalf("expected reason %q, got %q", "token expired", reason)
	}
}

func TestGameConsumeTooFast(t *testing.T) {
	s := NewGameStore(time.Hour)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	ok, reason := s.Consume(token, 10000)
	if ok {
		t.Fatal("too-fast consume should fail")
	}
	if reason != "completed too fast" {
		t.Fatalf("expected reason %q, got %q", "completed too fast", reason)
	}
}

func TestGameConsumeSingleUse(t *testing.T) {
	s := NewGameStore(time.Hour)
	token, err := s.Create()
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	ok, _ := s.Consume(token, 0)
	if !ok {
		t.Fatal("first Consume should succeed")
	}
	ok, reason := s.Consume(token, 0)
	if ok {
		t.Fatal("second Consume should fail")
	}
	if reason != "invalid token" {
		t.Fatalf("expected reason %q, got %q", "invalid token", reason)
	}
}
