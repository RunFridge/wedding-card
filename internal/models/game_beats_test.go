package models

import (
	"testing"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestRecordGameBeat(t *testing.T) {
	cleanTables(t)

	RecordGameBeat()

	today := time.Now().Format("2006-01-02")
	var count int
	err := database.DB.QueryRow(`SELECT count FROM game_beats WHERE date = ?`, today).Scan(&count)
	if err != nil {
		t.Fatalf("query game_beats failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

func TestRecordGameBeatMultiple(t *testing.T) {
	cleanTables(t)

	RecordGameBeat()
	RecordGameBeat()
	RecordGameBeat()

	today := time.Now().Format("2006-01-02")
	var count int
	err := database.DB.QueryRow(`SELECT count FROM game_beats WHERE date = ?`, today).Scan(&count)
	if err != nil {
		t.Fatalf("query game_beats failed: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestGetGameBeats(t *testing.T) {
	cleanTables(t)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	database.DB.Exec(`INSERT INTO game_beats (date, count) VALUES (?, ?)`, yesterday, 2)
	database.DB.Exec(`INSERT INTO game_beats (date, count) VALUES (?, ?)`, today, 7)

	beats, err := GetGameBeats(7)
	if err != nil {
		t.Fatalf("GetGameBeats failed: %v", err)
	}
	if len(beats) != 2 {
		t.Fatalf("expected 2 beats, got %d", len(beats))
	}
	if beats[0].Date != yesterday || beats[0].Count != 2 {
		t.Errorf("first = %+v, want %s/2", beats[0], yesterday)
	}
	if beats[1].Date != today || beats[1].Count != 7 {
		t.Errorf("second = %+v, want %s/7", beats[1], today)
	}
}

func TestGetGameBeatsEmpty(t *testing.T) {
	cleanTables(t)

	beats, err := GetGameBeats(30)
	if err != nil {
		t.Fatalf("GetGameBeats failed: %v", err)
	}
	if beats != nil {
		t.Errorf("expected nil slice for empty result, got %v", beats)
	}
}

func TestCleanupOldGameBeats(t *testing.T) {
	cleanTables(t)

	old := time.Now().AddDate(0, -7, 0).Format("2006-01-02")
	recent := time.Now().AddDate(0, 0, -10).Format("2006-01-02")

	database.DB.Exec(`INSERT INTO game_beats (date, count) VALUES (?, ?)`, old, 99)
	database.DB.Exec(`INSERT INTO game_beats (date, count) VALUES (?, ?)`, recent, 5)

	if err := CleanupOldGameBeats(6); err != nil {
		t.Fatalf("CleanupOldGameBeats failed: %v", err)
	}

	var n int
	database.DB.QueryRow(`SELECT COUNT(*) FROM game_beats`).Scan(&n)
	if n != 1 {
		t.Errorf("expected 1 row after cleanup, got %d", n)
	}
}
