package models

import (
	"database/sql"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestSetPhotoEvaluated(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "a.jpg", "", "photo.jpg", "", "")

	if err := SetPhotoEvaluated(id, true, false); err != nil {
		t.Fatalf("SetPhotoEvaluated failed: %v", err)
	}

	photo, _ := GetPhotoUploadByID(id)
	if !photo.Evaluated {
		t.Error("expected Evaluated=true")
	}
	if photo.Hidden {
		t.Error("expected Hidden=false")
	}
}

func TestGetPhotoHashnameByID(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "myhash.jpg", "", "photo.jpg", "", "")

	hashname, err := GetPhotoHashnameByID(id)
	if err != nil {
		t.Fatalf("GetPhotoHashnameByID failed: %v", err)
	}
	if hashname != "myhash.jpg" {
		t.Errorf("hashname = %q, want %q", hashname, "myhash.jpg")
	}
}

func TestGetPhotoHashnameByIDNotFound(t *testing.T) {
	cleanTables(t)

	_, err := GetPhotoHashnameByID(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetPhotoPasswordHash(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "a.jpg", "", "photo.jpg", "", "mypwhash")

	hash, err := GetPhotoPasswordHash(id)
	if err != nil {
		t.Fatalf("GetPhotoPasswordHash failed: %v", err)
	}
	if hash != "mypwhash" {
		t.Errorf("hash = %q, want %q", hash, "mypwhash")
	}
}

func TestIncrementPhotoHearts(t *testing.T) {
	cleanTables(t)

	id1, _ := CreatePhotoUpload("Alice", "127.0.0.1", "a.jpg", "", "p1.jpg", "", "")
	id2, _ := CreatePhotoUpload("Bob", "127.0.0.1", "b.jpg", "", "p2.jpg", "", "")

	updates := map[int64]int64{id1: 3, id2: 5}
	if err := IncrementPhotoHearts(updates); err != nil {
		t.Fatalf("IncrementPhotoHearts failed: %v", err)
	}

	hearts, _ := GetPhotoHearts([]int64{id1, id2})
	if hearts[id1] != 3 {
		t.Errorf("id1 hearts = %d, want 3", hearts[id1])
	}
	if hearts[id2] != 5 {
		t.Errorf("id2 hearts = %d, want 5", hearts[id2])
	}

	// Increment again
	IncrementPhotoHearts(map[int64]int64{id1: 2})
	hearts, _ = GetPhotoHearts([]int64{id1})
	if hearts[id1] != 5 {
		t.Errorf("id1 hearts = %d, want 5 after second increment", hearts[id1])
	}
}

func TestResetPhotoHearts(t *testing.T) {
	cleanTables(t)

	id, _ := CreatePhotoUpload("Alice", "127.0.0.1", "a.jpg", "", "p.jpg", "", "")
	IncrementPhotoHearts(map[int64]int64{id: 10})

	if err := ResetPhotoHearts(id); err != nil {
		t.Fatalf("ResetPhotoHearts failed: %v", err)
	}

	hearts, _ := GetPhotoHearts([]int64{id})
	if hearts[id] != 0 {
		t.Errorf("hearts = %d, want 0 after reset", hearts[id])
	}
}

func TestResetPhotoHeartsNotFound(t *testing.T) {
	cleanTables(t)

	err := ResetPhotoHearts(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetTotalHearts(t *testing.T) {
	cleanTables(t)

	total, err := GetTotalHearts()
	if err != nil {
		t.Fatalf("GetTotalHearts failed: %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}

	id1, _ := CreatePhotoUpload("A", "127.0.0.1", "a.jpg", "", "p1.jpg", "", "")
	id2, _ := CreatePhotoUpload("B", "127.0.0.1", "b.jpg", "", "p2.jpg", "", "")
	SetPhotoUploadHidden(id1, false)
	SetPhotoUploadHidden(id2, false)
	IncrementPhotoHearts(map[int64]int64{id1: 3, id2: 7})

	total, _ = GetTotalHearts()
	if total != 10 {
		t.Errorf("total = %d, want 10", total)
	}
}

func TestGetTotalHeartsExcludesHidden(t *testing.T) {
	cleanTables(t)

	id1, _ := CreatePhotoUpload("A", "127.0.0.1", "a.jpg", "", "p1.jpg", "", "")
	SetPhotoUploadHidden(id1, false)
	IncrementPhotoHearts(map[int64]int64{id1: 5})

	id2, _ := CreatePhotoUpload("B", "127.0.0.1", "b.jpg", "", "p2.jpg", "", "")
	// id2 is hidden by default
	IncrementPhotoHearts(map[int64]int64{id2: 10})

	total, _ := GetTotalHearts()
	if total != 5 {
		t.Errorf("total = %d, want 5 (excluding hidden)", total)
	}
}

func TestGetPhotoHeartsEmpty(t *testing.T) {
	cleanTables(t)

	hearts, err := GetPhotoHearts([]int64{})
	if err != nil {
		t.Fatalf("GetPhotoHearts failed: %v", err)
	}
	if len(hearts) != 0 {
		t.Errorf("expected empty map, got %d entries", len(hearts))
	}
}

func TestGetVisiblePhotoUploadsPopularSort(t *testing.T) {
	cleanTables(t)

	id1, _ := CreatePhotoUpload("Low", "127.0.0.1", "a.jpg", "", "p1.jpg", "", "")
	id2, _ := CreatePhotoUpload("High", "127.0.0.1", "b.jpg", "", "p2.jpg", "", "")
	SetPhotoUploadHidden(id1, false)
	SetPhotoUploadHidden(id2, false)

	database.DB.Exec(`UPDATE photo_uploads SET hearts = 1 WHERE id = ?`, id1)
	database.DB.Exec(`UPDATE photo_uploads SET hearts = 10 WHERE id = ?`, id2)

	uploads, err := GetVisiblePhotoUploads("popular")
	if err != nil {
		t.Fatalf("GetVisiblePhotoUploads(popular) failed: %v", err)
	}
	if len(uploads) != 2 {
		t.Fatalf("expected 2, got %d", len(uploads))
	}
	if uploads[0].Name != "High" {
		t.Errorf("expected most popular first (High), got %q", uploads[0].Name)
	}
}
