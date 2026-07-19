package models

import (
	"database/sql"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestCreateGuestbookEntry(t *testing.T) {
	cleanTables(t)

	entry, err := CreateGuestbookEntry("Alice", "Hello!", "127.0.0.1", "hash123", false)
	if err != nil {
		t.Fatalf("CreateGuestbookEntry failed: %v", err)
	}
	if entry.ID < 1 {
		t.Error("expected positive ID")
	}
	if entry.Nickname != "Alice" {
		t.Errorf("Nickname = %q, want %q", entry.Nickname, "Alice")
	}
	if entry.Message != "Hello!" {
		t.Errorf("Message = %q, want %q", entry.Message, "Hello!")
	}
}

func TestGetGuestbookEntryByID(t *testing.T) {
	cleanTables(t)

	created, _ := CreateGuestbookEntry("Bob", "Hi", "127.0.0.1", "hash", false)
	entry, err := GetGuestbookEntryByID(created.ID)
	if err != nil {
		t.Fatalf("GetGuestbookEntryByID failed: %v", err)
	}
	if entry.Nickname != "Bob" {
		t.Errorf("Nickname = %q, want %q", entry.Nickname, "Bob")
	}
	if entry.PasswordHash != "hash" {
		t.Errorf("PasswordHash = %q, want %q", entry.PasswordHash, "hash")
	}
}

func TestGetGuestbookEntryByIDNotFound(t *testing.T) {
	cleanTables(t)

	_, err := GetGuestbookEntryByID(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestUpdateGuestbookEntry(t *testing.T) {
	cleanTables(t)

	created, _ := CreateGuestbookEntry("Carol", "Old", "127.0.0.1", "hash", false)
	if err := UpdateGuestbookEntry(created.ID, "New"); err != nil {
		t.Fatalf("UpdateGuestbookEntry failed: %v", err)
	}

	updated, _ := GetGuestbookEntryByID(created.ID)
	if updated.Message != "New" {
		t.Errorf("Message = %q, want %q", updated.Message, "New")
	}
}

func TestDeleteGuestbookEntry(t *testing.T) {
	cleanTables(t)

	created, _ := CreateGuestbookEntry("Dave", "Bye", "127.0.0.1", "hash", false)
	if err := DeleteGuestbookEntry(created.ID); err != nil {
		t.Fatalf("DeleteGuestbookEntry failed: %v", err)
	}

	_, err := GetGuestbookEntryByID(created.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestGetVisibleGuestbookEntriesEmpty(t *testing.T) {
	cleanTables(t)

	entries, err := GetVisibleGuestbookEntries()
	if err != nil {
		t.Fatalf("GetVisibleGuestbookEntries failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetVisibleGuestbookEntriesFiltersHidden(t *testing.T) {
	cleanTables(t)

	CreateGuestbookEntry("Visible", "see me", "127.0.0.1", "hash", false)
	CreateGuestbookEntry("Hidden", "hidden", "127.0.0.1", "hash", false)

	// Manually hide the second entry
	entries, _ := GetVisibleGuestbookEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries before hiding, got %d", len(entries))
	}

	// Get the second entry's ID by querying
	var hiddenID int64
	for _, e := range entries {
		if e.Nickname == "Hidden" {
			hiddenID = e.ID
		}
	}

	// Use raw SQL to set hidden flag
	if _, err := database.DB.Exec("UPDATE guestbook_entries SET hidden = 1 WHERE id = ?", hiddenID); err != nil {
		t.Fatalf("failed to hide entry: %v", err)
	}

	visible, err := GetVisibleGuestbookEntries()
	if err != nil {
		t.Fatalf("GetVisibleGuestbookEntries failed: %v", err)
	}
	if len(visible) != 1 {
		t.Errorf("expected 1 visible entry, got %d", len(visible))
	}
	if visible[0].Nickname != "Visible" {
		t.Errorf("expected visible entry nickname %q, got %q", "Visible", visible[0].Nickname)
	}
}

func TestGetVisibleGuestbookEntriesNewestFirst(t *testing.T) {
	cleanTables(t)

	// Insert with explicit timestamps to guarantee ordering
	database.DB.Exec(`INSERT INTO guestbook_entries (nickname, message, ip, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`,
		"First", "first", "127.0.0.1", "hash", "2025-01-01 00:00:00")
	database.DB.Exec(`INSERT INTO guestbook_entries (nickname, message, ip, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`,
		"Second", "second", "127.0.0.1", "hash", "2025-01-02 00:00:00")

	entries, err := GetVisibleGuestbookEntries()
	if err != nil {
		t.Fatalf("GetVisibleGuestbookEntries failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Nickname != "Second" {
		t.Errorf("expected newest first, got %q", entries[0].Nickname)
	}
}
