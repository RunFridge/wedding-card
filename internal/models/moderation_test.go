package models

import (
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestGetGuestbookModerationStats(t *testing.T) {
	cleanTables(t)

	stats, err := GetGuestbookModerationStats()
	if err != nil {
		t.Fatalf("GetGuestbookModerationStats failed: %v", err)
	}
	if stats.Total != 0 {
		t.Errorf("Total = %d, want 0", stats.Total)
	}

	// Add entries with various states
	CreateGuestbookEntry("A", "msg1", "127.0.0.1", "hash", false) // pending (evaluated=0)
	CreateGuestbookEntry("B", "msg2", "127.0.0.1", "hash", false)
	database.DB.Exec(`UPDATE guestbook_entries SET evaluated = 1, hidden = 0 WHERE nickname = 'B'`) // approved
	CreateGuestbookEntry("C", "msg3", "127.0.0.1", "hash", false)
	database.DB.Exec(`UPDATE guestbook_entries SET evaluated = 1, hidden = 1 WHERE nickname = 'C'`) // flagged

	stats, err = GetGuestbookModerationStats()
	if err != nil {
		t.Fatalf("GetGuestbookModerationStats failed: %v", err)
	}
	if stats.Total != 3 {
		t.Errorf("Total = %d, want 3", stats.Total)
	}
	if stats.Pending != 1 {
		t.Errorf("Pending = %d, want 1", stats.Pending)
	}
	if stats.Approved != 1 {
		t.Errorf("Approved = %d, want 1", stats.Approved)
	}
	if stats.Flagged != 1 {
		t.Errorf("Flagged = %d, want 1", stats.Flagged)
	}
}

func TestGetPhotoModerationStats(t *testing.T) {
	cleanTables(t)

	stats, err := GetPhotoModerationStats()
	if err != nil {
		t.Fatalf("GetPhotoModerationStats failed: %v", err)
	}
	if stats.Total != 0 {
		t.Errorf("Total = %d, want 0", stats.Total)
	}

	CreatePhotoUpload("Alice", "127.0.0.1", "a.jpg", "", "photo.jpg", "", "")
	CreatePhotoUpload("Bob", "127.0.0.1", "b.jpg", "", "photo2.jpg", "", "")

	stats, _ = GetPhotoModerationStats()
	if stats.Total != 2 {
		t.Errorf("Total = %d, want 2", stats.Total)
	}
	if stats.Pending != 2 {
		t.Errorf("Pending = %d, want 2", stats.Pending)
	}
}
