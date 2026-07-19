package database

import (
	"testing"
)

func TestInitMemory(t *testing.T) {
	if err := Init(":memory:"); err != nil {
		t.Fatalf("Init(:memory:) failed: %v", err)
	}
	defer Close()

	if DB == nil {
		t.Fatal("DB is nil after Init")
	}
}

func TestMigrationsCreateTables(t *testing.T) {
	if err := Init(":memory:"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	tables := []string{"guestbook_entries", "game_scores", "photo_uploads"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestMigrationsIdempotent(t *testing.T) {
	if err := Init(":memory:"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	if err := RunMigrations(); err != nil {
		t.Fatalf("RunMigrations second call failed: %v", err)
	}
}
