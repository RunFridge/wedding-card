package models

import (
	"database/sql"
	"testing"
)

func TestGetAllGuestbookEntries(t *testing.T) {
	cleanTables(t)

	CreateGuestbookEntry("Alice", "Hello", "10.0.0.1", "hash1", false)
	CreateGuestbookEntry("Bob", "World", "10.0.0.2", "hash2", true)

	entries, err := GetAllGuestbookEntries()
	if err != nil {
		t.Fatalf("GetAllGuestbookEntries failed: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	// Should include IP (admin view)
	if entries[0].IP == "" {
		t.Error("expected IP to be populated")
	}
	// Verify both entries present (ordering may vary with same timestamp in :memory:)
	names := map[string]bool{entries[0].Nickname: true, entries[1].Nickname: true}
	if !names["Alice"] || !names["Bob"] {
		t.Errorf("expected Alice and Bob, got %q and %q", entries[0].Nickname, entries[1].Nickname)
	}
}

func TestGetAllGuestbookEntriesSecret(t *testing.T) {
	cleanTables(t)

	CreateGuestbookEntry("Secret", "Hidden msg", "10.0.0.1", "hash", true)

	entries, err := GetAllGuestbookEntries()
	if err != nil {
		t.Fatalf("GetAllGuestbookEntries failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1, got %d", len(entries))
	}
	if !entries[0].Secret {
		t.Error("expected Secret=true")
	}
	// Admin should see real content
	if entries[0].Message != "Hidden msg" {
		t.Errorf("Message = %q, want %q", entries[0].Message, "Hidden msg")
	}
}

func TestSetGuestbookEntryHidden(t *testing.T) {
	cleanTables(t)

	entry, _ := CreateGuestbookEntry("Alice", "Hello", "127.0.0.1", "hash", false)

	if err := SetGuestbookEntryHidden(entry.ID, true); err != nil {
		t.Fatalf("SetGuestbookEntryHidden failed: %v", err)
	}

	visible, _ := GetVisibleGuestbookEntries()
	if len(visible) != 0 {
		t.Errorf("expected 0 visible after hiding, got %d", len(visible))
	}

	SetGuestbookEntryHidden(entry.ID, false)
	visible, _ = GetVisibleGuestbookEntries()
	if len(visible) != 1 {
		t.Errorf("expected 1 visible after unhiding, got %d", len(visible))
	}
}

func TestSetGuestbookEntryHiddenNotFound(t *testing.T) {
	cleanTables(t)

	err := SetGuestbookEntryHidden(99999, true)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestSetGuestbookEvaluated(t *testing.T) {
	cleanTables(t)

	entry, _ := CreateGuestbookEntry("Alice", "Hello", "127.0.0.1", "hash", false)

	if err := SetGuestbookEvaluated(entry.ID, true, false); err != nil {
		t.Fatalf("SetGuestbookEvaluated failed: %v", err)
	}

	all, _ := GetAllGuestbookEntries()
	if !all[0].Evaluated {
		t.Error("expected Evaluated=true")
	}
	if all[0].Hidden {
		t.Error("expected Hidden=false")
	}
}

func TestGetGuestbookContentByID(t *testing.T) {
	cleanTables(t)

	entry, _ := CreateGuestbookEntry("Alice", "Secret content", "127.0.0.1", "hash", false)

	nickname, message, err := GetGuestbookContentByID(entry.ID)
	if err != nil {
		t.Fatalf("GetGuestbookContentByID failed: %v", err)
	}
	if nickname != "Alice" {
		t.Errorf("nickname = %q, want %q", nickname, "Alice")
	}
	if message != "Secret content" {
		t.Errorf("message = %q, want %q", message, "Secret content")
	}
}

func TestGetGuestbookContentByIDNotFound(t *testing.T) {
	cleanTables(t)

	_, _, err := GetGuestbookContentByID(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestResetGuestbookEvaluated(t *testing.T) {
	cleanTables(t)

	entry, _ := CreateGuestbookEntry("Alice", "Hello", "127.0.0.1", "hash", false)
	SetGuestbookEvaluated(entry.ID, true, false)

	if err := ResetGuestbookEvaluated(entry.ID); err != nil {
		t.Fatalf("ResetGuestbookEvaluated failed: %v", err)
	}

	all, _ := GetAllGuestbookEntries()
	if all[0].Evaluated {
		t.Error("expected Evaluated=false after reset")
	}
	if !all[0].Hidden {
		t.Error("expected Hidden=true after reset")
	}
}

func TestCreateGuestbookEntrySecret(t *testing.T) {
	cleanTables(t)

	entry, err := CreateGuestbookEntry("Alice", "Secret!", "127.0.0.1", "hash", true)
	if err != nil {
		t.Fatalf("CreateGuestbookEntry failed: %v", err)
	}
	if !entry.Secret {
		t.Error("expected Secret=true")
	}

	// Should appear in visible entries
	visible, _ := GetVisibleGuestbookEntries()
	if len(visible) != 1 {
		t.Fatalf("expected 1, got %d", len(visible))
	}
	if !visible[0].Secret {
		t.Error("expected Secret=true in visible entries")
	}
}

func TestGetGuestbookEntryByIDSecret(t *testing.T) {
	cleanTables(t)

	created, _ := CreateGuestbookEntry("Alice", "Secret!", "127.0.0.1", "hash", true)
	entry, err := GetGuestbookEntryByID(created.ID)
	if err != nil {
		t.Fatalf("GetGuestbookEntryByID failed: %v", err)
	}
	if !entry.Secret {
		t.Error("expected Secret=true")
	}
}
