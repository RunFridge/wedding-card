package models

import (
	"database/sql"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestCreatePhotoUpload(t *testing.T) {
	cleanTables(t)

	id, err := CreatePhotoUpload("Alice", "127.0.0.1", "abc123.jpg", "", "photo.jpg", "data:image/jpeg;base64,thumb", "")
	if err != nil {
		t.Fatalf("CreatePhotoUpload failed: %v", err)
	}
	if id < 1 {
		t.Error("expected positive ID")
	}
}

func TestCreatePhotoUploadWithoutThumbnail(t *testing.T) {
	cleanTables(t)

	id, err := CreatePhotoUpload("Bob", "127.0.0.1", "def456.jpg", "", "photo2.jpg", "", "")
	if err != nil {
		t.Fatalf("CreatePhotoUpload failed: %v", err)
	}
	if id < 1 {
		t.Error("expected positive ID")
	}
}

func TestGetVisiblePhotoUploadsEmpty(t *testing.T) {
	cleanTables(t)

	uploads, err := GetVisiblePhotoUploads("")
	if err != nil {
		t.Fatalf("GetVisiblePhotoUploads failed: %v", err)
	}
	if len(uploads) != 0 {
		t.Errorf("expected 0 uploads, got %d", len(uploads))
	}
}

func TestGetVisiblePhotoUploadsHiddenByDefault(t *testing.T) {
	cleanTables(t)

	CreatePhotoUpload("Alice", "127.0.0.1", "abc.jpg", "", "photo.jpg", "", "")

	uploads, err := GetVisiblePhotoUploads("")
	if err != nil {
		t.Fatalf("GetVisiblePhotoUploads failed: %v", err)
	}
	if len(uploads) != 0 {
		t.Errorf("expected 0 visible uploads (hidden by default), got %d", len(uploads))
	}
}

func TestGetVisiblePhotoUploadsAfterUnhide(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "abc.jpg", "", "photo.jpg", "data:thumb", "")
	SetPhotoUploadHidden(id, false)

	uploads, err := GetVisiblePhotoUploads("")
	if err != nil {
		t.Fatalf("GetVisiblePhotoUploads failed: %v", err)
	}
	if len(uploads) != 1 {
		t.Fatalf("expected 1 visible upload, got %d", len(uploads))
	}
	if uploads[0].Name != "Alice" {
		t.Errorf("Name = %q, want %q", uploads[0].Name, "Alice")
	}
	if uploads[0].Thumbnail != "data:thumb" {
		t.Errorf("Thumbnail = %q, want %q", uploads[0].Thumbnail, "data:thumb")
	}
}

func TestGetVisiblePhotoUploadsNewestFirst(t *testing.T) {
	cleanTables(t)

	database.DB.Exec(`INSERT INTO photo_uploads (name, ip_address, hashname, original_filename, thumbnail, hidden, upload_date) VALUES (?, ?, ?, ?, ?, 0, ?)`,
		"First", "127.0.0.1", "a.jpg", "a.jpg", "", "2025-01-01 00:00:00")
	database.DB.Exec(`INSERT INTO photo_uploads (name, ip_address, hashname, original_filename, thumbnail, hidden, upload_date) VALUES (?, ?, ?, ?, ?, 0, ?)`,
		"Second", "127.0.0.1", "b.jpg", "b.jpg", "", "2025-01-02 00:00:00")

	uploads, err := GetVisiblePhotoUploads("")
	if err != nil {
		t.Fatalf("GetVisiblePhotoUploads failed: %v", err)
	}
	if len(uploads) != 2 {
		t.Fatalf("expected 2 uploads, got %d", len(uploads))
	}
	if uploads[0].Name != "Second" {
		t.Errorf("expected newest first, got %q", uploads[0].Name)
	}
}

func TestGetAllPhotoUploads(t *testing.T) {
	cleanTables(t)

	CreatePhotoUpload("Alice", "10.0.0.1", "abc.jpg", "", "photo.jpg", "", "")
	CreatePhotoUpload("Bob", "10.0.0.2", "def.jpg", "", "photo2.jpg", "", "")

	uploads, err := GetAllPhotoUploads()
	if err != nil {
		t.Fatalf("GetAllPhotoUploads failed: %v", err)
	}
	if len(uploads) != 2 {
		t.Errorf("expected 2 uploads, got %d", len(uploads))
	}
	if uploads[0].IP == "" {
		t.Error("expected IP to be populated in AdminPhotoUpload")
	}
	if uploads[0].Hashname == "" {
		t.Error("expected Hashname to be populated in AdminPhotoUpload")
	}
}

func TestGetPhotoUploadByID(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "10.0.0.1", "abc.jpg", "", "photo.jpg", "", "")
	photo, err := GetPhotoUploadByID(id)
	if err != nil {
		t.Fatalf("GetPhotoUploadByID failed: %v", err)
	}
	if photo.Name != "Alice" {
		t.Errorf("Name = %q, want %q", photo.Name, "Alice")
	}
	if photo.Hashname != "abc.jpg" {
		t.Errorf("Hashname = %q, want %q", photo.Hashname, "abc.jpg")
	}
}

func TestGetPhotoUploadByIDNotFound(t *testing.T) {
	cleanTables(t)

	_, err := GetPhotoUploadByID(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestSetPhotoUploadHidden(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "abc.jpg", "", "photo.jpg", "", "")

	// Default is hidden=1
	photo, _ := GetPhotoUploadByID(id)
	if !photo.Hidden {
		t.Error("expected photo to be hidden by default")
	}

	// Unhide
	if err := SetPhotoUploadHidden(id, false); err != nil {
		t.Fatalf("SetPhotoUploadHidden(false) failed: %v", err)
	}
	photo, _ = GetPhotoUploadByID(id)
	if photo.Hidden {
		t.Error("expected photo to be visible after unhide")
	}

	// Hide again
	if err := SetPhotoUploadHidden(id, true); err != nil {
		t.Fatalf("SetPhotoUploadHidden(true) failed: %v", err)
	}
	photo, _ = GetPhotoUploadByID(id)
	if !photo.Hidden {
		t.Error("expected photo to be hidden again")
	}
}

func TestSetPhotoUploadHiddenNotFound(t *testing.T) {
	cleanTables(t)

	err := SetPhotoUploadHidden(99999, false)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestDeletePhotoUpload(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "abc.jpg", "", "photo.jpg", "", "")
	if err := DeletePhotoUpload(id); err != nil {
		t.Fatalf("DeletePhotoUpload failed: %v", err)
	}

	_, err := GetPhotoUploadByID(id)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestDeletePhotoUploadNotFound(t *testing.T) {
	cleanTables(t)

	err := DeletePhotoUpload(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
