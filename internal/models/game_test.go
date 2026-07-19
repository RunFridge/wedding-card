package models

import (
	"testing"
)

func TestCreateGameScore(t *testing.T) {
	cleanTables(t)

	score, err := CreateGameScore("ABC", 5000, "127.0.0.1")
	if err != nil {
		t.Fatalf("CreateGameScore failed: %v", err)
	}
	if score.ID < 1 {
		t.Error("expected positive ID")
	}
	if score.Nickname != "ABC" {
		t.Errorf("Nickname = %q, want %q", score.Nickname, "ABC")
	}
	if score.TimeMs != 5000 {
		t.Errorf("TimeMs = %d, want %d", score.TimeMs, 5000)
	}
}

func TestGetTopScoresEmpty(t *testing.T) {
	cleanTables(t)

	scores, err := GetTopScores(10)
	if err != nil {
		t.Fatalf("GetTopScores failed: %v", err)
	}
	if len(scores) != 0 {
		t.Errorf("expected 0 scores, got %d", len(scores))
	}
}

func TestGetTopScoresOrdering(t *testing.T) {
	cleanTables(t)

	CreateGameScore("AAA", 10000, "127.0.0.1")
	CreateGameScore("BBB", 3000, "127.0.0.1")
	CreateGameScore("CCC", 7000, "127.0.0.1")

	scores, err := GetTopScores(10)
	if err != nil {
		t.Fatalf("GetTopScores failed: %v", err)
	}
	if len(scores) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(scores))
	}
	if scores[0].Nickname != "BBB" {
		t.Errorf("expected fastest first (BBB), got %q", scores[0].Nickname)
	}
	if scores[2].Nickname != "AAA" {
		t.Errorf("expected slowest last (AAA), got %q", scores[2].Nickname)
	}
}

func TestGetTopScoresLimit(t *testing.T) {
	cleanTables(t)

	CreateGameScore("AAA", 1000, "127.0.0.1")
	CreateGameScore("BBB", 2000, "127.0.0.1")
	CreateGameScore("CCC", 3000, "127.0.0.1")

	scores, err := GetTopScores(2)
	if err != nil {
		t.Fatalf("GetTopScores failed: %v", err)
	}
	if len(scores) != 2 {
		t.Errorf("expected 2 scores with limit=2, got %d", len(scores))
	}
}
