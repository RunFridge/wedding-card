package models

import (
	"database/sql"
	"testing"
)

func TestCreateHallOfFameEntry(t *testing.T) {
	cleanTables(t)

	entry, err := CreateHallOfFameEntry("TestUser", "192.168.1.1")
	if err != nil {
		t.Fatalf("CreateHallOfFameEntry failed: %v", err)
	}
	if entry.ID < 1 {
		t.Error("expected positive ID")
	}
	if entry.Nickname != "TestUser" {
		t.Errorf("Nickname = %q, want %q", entry.Nickname, "TestUser")
	}
	if entry.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestGetHallOfFameEntries(t *testing.T) {
	cleanTables(t)

	CreateHallOfFameEntry("Alice", "10.0.0.1")
	CreateHallOfFameEntry("Bob", "10.0.0.2")
	CreateHallOfFameEntry("Charlie", "10.0.0.3")

	entries, err := GetHallOfFameEntries()
	if err != nil {
		t.Fatalf("GetHallOfFameEntries failed: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Nickname != "Alice" {
		t.Errorf("first entry = %q, want %q", entries[0].Nickname, "Alice")
	}
	if entries[2].Nickname != "Charlie" {
		t.Errorf("last entry = %q, want %q", entries[2].Nickname, "Charlie")
	}
}

func TestGetHallOfFameEntriesEmpty(t *testing.T) {
	cleanTables(t)

	entries, err := GetHallOfFameEntries()
	if err != nil {
		t.Fatalf("GetHallOfFameEntries failed: %v", err)
	}
	if entries == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetAllHallOfFameEntries(t *testing.T) {
	cleanTables(t)

	CreateHallOfFameEntry("Alice", "10.0.0.1")
	CreateHallOfFameEntry("Bob", "10.0.0.2")

	entries, err := GetAllHallOfFameEntries()
	if err != nil {
		t.Fatalf("GetAllHallOfFameEntries failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].IP != "10.0.0.1" {
		t.Errorf("first entry IP = %q, want %q", entries[0].IP, "10.0.0.1")
	}
	if entries[1].IP != "10.0.0.2" {
		t.Errorf("second entry IP = %q, want %q", entries[1].IP, "10.0.0.2")
	}
}

func TestDeleteHallOfFameEntry(t *testing.T) {
	cleanTables(t)

	entry, err := CreateHallOfFameEntry("ToDelete", "10.0.0.1")
	if err != nil {
		t.Fatalf("CreateHallOfFameEntry failed: %v", err)
	}

	if err := DeleteHallOfFameEntry(entry.ID); err != nil {
		t.Fatalf("DeleteHallOfFameEntry failed: %v", err)
	}

	entries, err := GetHallOfFameEntries()
	if err != nil {
		t.Fatalf("GetHallOfFameEntries failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after delete, got %d", len(entries))
	}
}

func TestDeleteHallOfFameEntryNotFound(t *testing.T) {
	cleanTables(t)

	err := DeleteHallOfFameEntry(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
