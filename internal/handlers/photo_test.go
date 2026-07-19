package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/models"
)

func createTestPhoto(t *testing.T, name string, hidden bool) int64 {
	t.Helper()
	id, err := models.CreatePhotoUpload(name, "127.0.0.1", "hash_"+name+".jpg", "", "original.jpg", "data:thumb", "")
	if err != nil {
		t.Fatalf("CreatePhotoUpload failed: %v", err)
	}
	if !hidden {
		if err := models.SetPhotoUploadHidden(id, false); err != nil {
			t.Fatalf("SetPhotoUploadHidden failed: %v", err)
		}
	}
	return id
}

// --- GetVisiblePhotos ---

func TestGetVisiblePhotosEmpty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/photos", nil)
	rr := httptest.NewRecorder()

	GetVisiblePhotos(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct{ Items []models.PhotoUpload }
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	uploads := resp.Items
	if len(uploads) != 0 {
		t.Errorf("expected 0 uploads, got %d", len(uploads))
	}
}

func TestGetVisiblePhotosOnlyVisible(t *testing.T) {
	setupTestDB(t)

	createTestPhoto(t, "visible", false)
	createTestPhoto(t, "hidden", true)

	req := httptest.NewRequest(http.MethodGet, "/api/photos", nil)
	rr := httptest.NewRecorder()

	GetVisiblePhotos(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct{ Items []models.PhotoUpload }
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	uploads := resp.Items
	if len(uploads) != 1 {
		t.Fatalf("expected 1 visible upload, got %d", len(uploads))
	}
	if uploads[0].Name != "visible" {
		t.Errorf("Name = %q, want %q", uploads[0].Name, "visible")
	}
	if uploads[0].Thumbnail != "data:thumb" {
		t.Errorf("Thumbnail = %q, want %q", uploads[0].Thumbnail, "data:thumb")
	}
}

// --- UploadPhoto ---

func TestUploadPhotoNoStore(t *testing.T) {
	setupTestDB(t)

	oldStore := Store
	Store = nil
	defer func() { Store = oldStore }()

	req := httptest.NewRequest(http.MethodPost, "/api/photos/upload", nil)
	rr := httptest.NewRecorder()

	UploadPhoto(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}

// --- AdminGetPhotos ---

func TestAdminGetPhotosEmpty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/photos", nil)
	rr := httptest.NewRecorder()

	AdminGetPhotos(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var uploads []models.AdminPhotoUpload
	if err := json.NewDecoder(rr.Body).Decode(&uploads); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(uploads) != 0 {
		t.Errorf("expected 0 uploads, got %d", len(uploads))
	}
}

func TestAdminGetPhotosAll(t *testing.T) {
	setupTestDB(t)

	createTestPhoto(t, "visible", false)
	createTestPhoto(t, "hidden", true)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/photos", nil)
	rr := httptest.NewRecorder()

	AdminGetPhotos(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var uploads []models.AdminPhotoUpload
	if err := json.NewDecoder(rr.Body).Decode(&uploads); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(uploads) != 2 {
		t.Errorf("expected 2 uploads (admin sees all), got %d", len(uploads))
	}
}

// --- AdminTogglePhotoVisibility ---

func TestAdminTogglePhotoVisibility(t *testing.T) {
	setupTestDB(t)

	id := createTestPhoto(t, "Alice", true)

	body := `{"hidden":false}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admin/photos/%d/visibility", id), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(id)})
	rr := httptest.NewRecorder()

	AdminTogglePhotoVisibility(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	photo, _ := models.GetPhotoUploadByID(id)
	if photo.Hidden {
		t.Error("expected photo to be visible after toggle")
	}
}

func TestAdminTogglePhotoVisibilityNotFound(t *testing.T) {
	setupTestDB(t)

	body := `{"hidden":false}`
	req := httptest.NewRequest(http.MethodPatch, "/api/admin/photos/99999/visibility", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "99999"})
	rr := httptest.NewRecorder()

	AdminTogglePhotoVisibility(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestAdminTogglePhotoVisibilityInvalidID(t *testing.T) {
	body := `{"hidden":false}`
	req := httptest.NewRequest(http.MethodPatch, "/api/admin/photos/abc/visibility", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	AdminTogglePhotoVisibility(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestAdminTogglePhotoVisibilityInvalidBody(t *testing.T) {
	setupTestDB(t)

	id := createTestPhoto(t, "Alice", true)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admin/photos/%d/visibility", id), strings.NewReader("not json"))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(id)})
	rr := httptest.NewRecorder()

	AdminTogglePhotoVisibility(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// --- AdminDeletePhoto ---

func TestAdminDeletePhoto(t *testing.T) {
	setupTestDB(t)

	id := createTestPhoto(t, "Alice", true)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/admin/photos/%d", id), nil)
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(id)})
	rr := httptest.NewRecorder()

	AdminDeletePhoto(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusNoContent, rr.Body.String())
	}

	// Verify deleted
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM photo_uploads WHERE id = ?", id).Scan(&count)
	if count != 0 {
		t.Error("expected photo to be deleted from database")
	}
}

func TestAdminDeletePhotoNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/photos/99999", nil)
	req = newChiContext(req, map[string]string{"id": "99999"})
	rr := httptest.NewRecorder()

	AdminDeletePhoto(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestAdminDeletePhotoInvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/photos/abc", nil)
	req = newChiContext(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	AdminDeletePhoto(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// --- generateUUID ---

func TestGenerateUUID(t *testing.T) {
	uuid := generateUUID()

	// UUID v4 format: 8-4-4-4-12 hex chars
	if len(uuid) != 36 {
		t.Errorf("UUID length = %d, want 36", len(uuid))
	}

	// Check version nibble (should be 4)
	if uuid[14] != '4' {
		t.Errorf("UUID version nibble = %c, want '4'", uuid[14])
	}

	// Check variant (should be 8, 9, a, or b)
	variant := uuid[19]
	if variant != '8' && variant != '9' && variant != 'a' && variant != 'b' {
		t.Errorf("UUID variant nibble = %c, want one of 8/9/a/b", variant)
	}

	// Check uniqueness
	uuid2 := generateUUID()
	if uuid == uuid2 {
		t.Error("two consecutive UUIDs should not be equal")
	}
}
