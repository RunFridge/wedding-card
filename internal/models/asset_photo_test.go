package models

import (
	"database/sql"
	"testing"
)

func TestCreateAssetPhoto(t *testing.T) {
	cleanTables(t)

	id, err := CreateAssetPhoto("Wedding", "abc.jpg", "abc_thumb.jpg", "photo.jpg")
	if err != nil {
		t.Fatalf("CreateAssetPhoto failed: %v", err)
	}
	if id < 1 {
		t.Error("expected positive ID")
	}
}

func TestGetAllAssetPhotosEmpty(t *testing.T) {
	cleanTables(t)

	photos, err := GetAllAssetPhotos()
	if err != nil {
		t.Fatalf("GetAllAssetPhotos failed: %v", err)
	}
	if len(photos) != 0 {
		t.Errorf("expected 0 photos, got %d", len(photos))
	}
}

func TestGetAllAssetPhotos(t *testing.T) {
	cleanTables(t)

	CreateAssetPhoto("Photo1", "a.jpg", "a_t.jpg", "orig1.jpg")
	CreateAssetPhoto("Photo2", "b.jpg", "b_t.jpg", "orig2.jpg")

	photos, err := GetAllAssetPhotos()
	if err != nil {
		t.Fatalf("GetAllAssetPhotos failed: %v", err)
	}
	if len(photos) != 2 {
		t.Errorf("expected 2 photos, got %d", len(photos))
	}
}

func TestGetAssetPhotoByID(t *testing.T) {
	cleanTables(t)

	id, _ := CreateAssetPhoto("Wedding", "abc.jpg", "abc_t.jpg", "orig.jpg")
	photo, err := GetAssetPhotoByID(id)
	if err != nil {
		t.Fatalf("GetAssetPhotoByID failed: %v", err)
	}
	if photo.Label != "Wedding" {
		t.Errorf("Label = %q, want %q", photo.Label, "Wedding")
	}
	if photo.Hashname != "abc.jpg" {
		t.Errorf("Hashname = %q, want %q", photo.Hashname, "abc.jpg")
	}
}

func TestGetAssetPhotoByIDNotFound(t *testing.T) {
	cleanTables(t)

	_, err := GetAssetPhotoByID(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestSetAssetPhotoGameFlag(t *testing.T) {
	cleanTables(t)

	id, _ := CreateAssetPhoto("Photo", "a.jpg", "a_t.jpg", "o.jpg")

	if err := SetAssetPhotoGameFlag(id, true); err != nil {
		t.Fatalf("SetAssetPhotoGameFlag failed: %v", err)
	}

	photo, _ := GetAssetPhotoByID(id)
	if !photo.UseForGame {
		t.Error("expected UseForGame=true")
	}

	SetAssetPhotoGameFlag(id, false)
	photo, _ = GetAssetPhotoByID(id)
	if photo.UseForGame {
		t.Error("expected UseForGame=false")
	}
}

func TestSetAssetPhotoGameFlagNotFound(t *testing.T) {
	cleanTables(t)

	err := SetAssetPhotoGameFlag(99999, true)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetGameAssetPhotos(t *testing.T) {
	cleanTables(t)

	id1, _ := CreateAssetPhoto("Game1", "a.jpg", "a_t.jpg", "o1.jpg")
	CreateAssetPhoto("NotGame", "b.jpg", "b_t.jpg", "o2.jpg")
	id3, _ := CreateAssetPhoto("Game2", "c.jpg", "c_t.jpg", "o3.jpg")

	SetAssetPhotoGameFlag(id1, true)
	SetAssetPhotoGameFlag(id3, true)

	photos, err := GetGameAssetPhotos()
	if err != nil {
		t.Fatalf("GetGameAssetPhotos failed: %v", err)
	}
	if len(photos) != 2 {
		t.Errorf("expected 2 game photos, got %d", len(photos))
	}
}

func TestCountGameAssetPhotos(t *testing.T) {
	cleanTables(t)

	id1, _ := CreateAssetPhoto("P1", "a.jpg", "a_t.jpg", "o1.jpg")
	CreateAssetPhoto("P2", "b.jpg", "b_t.jpg", "o2.jpg")

	count, _ := CountGameAssetPhotos()
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	SetAssetPhotoGameFlag(id1, true)
	count, _ = CountGameAssetPhotos()
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestSetAssetPhotoAsMain(t *testing.T) {
	cleanTables(t)

	id1, _ := CreateAssetPhoto("P1", "a.jpg", "a_t.jpg", "o1.jpg")
	id2, _ := CreateAssetPhoto("P2", "b.jpg", "b_t.jpg", "o2.jpg")

	if err := SetAssetPhotoAsMain(id1); err != nil {
		t.Fatalf("SetAssetPhotoAsMain failed: %v", err)
	}

	p1, _ := GetAssetPhotoByID(id1)
	if !p1.IsMainPhoto {
		t.Error("expected id1 to be main photo")
	}

	// Setting another as main should clear the first
	SetAssetPhotoAsMain(id2)
	p1, _ = GetAssetPhotoByID(id1)
	p2, _ := GetAssetPhotoByID(id2)
	if p1.IsMainPhoto {
		t.Error("expected id1 to no longer be main photo")
	}
	if !p2.IsMainPhoto {
		t.Error("expected id2 to be main photo")
	}
}

func TestSetAssetPhotoAsMainNotFound(t *testing.T) {
	cleanTables(t)

	err := SetAssetPhotoAsMain(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetMainAssetPhoto(t *testing.T) {
	cleanTables(t)

	_, err := GetMainAssetPhoto()
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows when no main, got %v", err)
	}

	id, _ := CreateAssetPhoto("Main", "a.jpg", "a_t.jpg", "o.jpg")
	SetAssetPhotoAsMain(id)

	photo, err := GetMainAssetPhoto()
	if err != nil {
		t.Fatalf("GetMainAssetPhoto failed: %v", err)
	}
	if photo.Label != "Main" {
		t.Errorf("Label = %q, want %q", photo.Label, "Main")
	}
}

func TestClearMainAssetPhoto(t *testing.T) {
	cleanTables(t)

	id, _ := CreateAssetPhoto("Main", "a.jpg", "a_t.jpg", "o.jpg")
	SetAssetPhotoAsMain(id)

	if err := ClearMainAssetPhoto(id); err != nil {
		t.Fatalf("ClearMainAssetPhoto failed: %v", err)
	}

	photo, _ := GetAssetPhotoByID(id)
	if photo.IsMainPhoto {
		t.Error("expected IsMainPhoto=false after clear")
	}
}

func TestClearMainAssetPhotoNotFound(t *testing.T) {
	cleanTables(t)

	err := ClearMainAssetPhoto(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestDeleteAssetPhoto(t *testing.T) {
	cleanTables(t)

	id, _ := CreateAssetPhoto("Delete Me", "a.jpg", "a_t.jpg", "o.jpg")
	photo, err := DeleteAssetPhoto(id)
	if err != nil {
		t.Fatalf("DeleteAssetPhoto failed: %v", err)
	}
	if photo.Label != "Delete Me" {
		t.Errorf("Label = %q, want %q", photo.Label, "Delete Me")
	}

	_, err = GetAssetPhotoByID(id)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestDeleteAssetPhotoNotFound(t *testing.T) {
	cleanTables(t)

	_, err := DeleteAssetPhoto(99999)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestUpdateAssetPhoto(t *testing.T) {
	cleanTables(t)

	id, _ := CreateAssetPhoto("Old", "a.jpg", "a_t.jpg", "o.jpg")
	if err := UpdateAssetPhoto(id, "New", 5); err != nil {
		t.Fatalf("UpdateAssetPhoto failed: %v", err)
	}

	photo, _ := GetAssetPhotoByID(id)
	if photo.Label != "New" {
		t.Errorf("Label = %q, want %q", photo.Label, "New")
	}
	if photo.SortOrder != 5 {
		t.Errorf("SortOrder = %d, want %d", photo.SortOrder, 5)
	}
}

func TestUpdateAssetPhotoNotFound(t *testing.T) {
	cleanTables(t)

	err := UpdateAssetPhoto(99999, "New", 0)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
