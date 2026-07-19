package models

import (
	"database/sql"
	"testing"
)

func TestGetAllGameScores(t *testing.T) {
	cleanTables(t)

	CreateGameScore("AAA", 5000, "10.0.0.1")
	CreateGameScore("BBB", 3000, "10.0.0.2")
	CreateGameScore("CCC", 8000, "10.0.0.3")

	scores, err := GetAllGameScores(100)
	if err != nil {
		t.Fatalf("GetAllGameScores failed: %v", err)
	}
	if len(scores) != 3 {
		t.Errorf("expected 3 scores, got %d", len(scores))
	}
	// Should be ordered by time_ms ASC
	if scores[0].Nickname != "BBB" {
		t.Errorf("expected fastest first (BBB), got %q", scores[0].Nickname)
	}
	// Admin scores should include IP
	if scores[0].IP == "" {
		t.Error("expected IP to be populated")
	}
}

func TestGetAllGameScoresLimit(t *testing.T) {
	cleanTables(t)

	CreateGameScore("AAA", 1000, "127.0.0.1")
	CreateGameScore("BBB", 2000, "127.0.0.1")
	CreateGameScore("CCC", 3000, "127.0.0.1")

	scores, _ := GetAllGameScores(2)
	if len(scores) != 2 {
		t.Errorf("expected 2 scores with limit=2, got %d", len(scores))
	}
}

func TestGetAllGameScoresEmpty(t *testing.T) {
	cleanTables(t)

	scores, err := GetAllGameScores(10)
	if err != nil {
		t.Fatalf("GetAllGameScores failed: %v", err)
	}
	if len(scores) != 0 {
		t.Errorf("expected 0, got %d", len(scores))
	}
}

func TestDeleteGameScore(t *testing.T) {
	cleanTables(t)

	score, _ := CreateGameScore("AAA", 5000, "127.0.0.1")
	if err := DeleteGameScore(score.ID); err != nil {
		t.Fatalf("DeleteGameScore failed: %v", err)
	}

	scores, _ := GetAllGameScores(10)
	if len(scores) != 0 {
		t.Errorf("expected 0 after delete, got %d", len(scores))
	}
}

func TestDeleteGameScoreNotFound(t *testing.T) {
	cleanTables(t)

	err := DeleteGameScore(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestHasPlayedFromIP(t *testing.T) {
	cleanTables(t)

	played, err := HasPlayedFromIP("10.0.0.1")
	if err != nil {
		t.Fatalf("HasPlayedFromIP failed: %v", err)
	}
	if played {
		t.Error("expected false for new IP")
	}

	CreateGameScore("AAA", 5000, "10.0.0.1")

	played, _ = HasPlayedFromIP("10.0.0.1")
	if !played {
		t.Error("expected true after creating score")
	}

	played, _ = HasPlayedFromIP("10.0.0.2")
	if played {
		t.Error("expected false for different IP")
	}
}
