package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RunFridge/wedding-card/internal/models"
)

func TestGetGuestbookEmpty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/guestbook", nil)
	rr := httptest.NewRecorder()

	GetGuestbook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct{ Items []models.GuestbookEntry }
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	entries := resp.Items
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetGuestbookWithEntries(t *testing.T) {
	setupTestDB(t)

	createTestEntry(t, "Alice", "Hello", "pass1")
	createTestEntry(t, "Bob", "World", "pass2")

	req := httptest.NewRequest(http.MethodGet, "/api/guestbook", nil)
	rr := httptest.NewRecorder()

	GetGuestbook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct{ Items []models.GuestbookEntry }
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	entries := resp.Items
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestCreateGuestbookEntrySuccess(t *testing.T) {
	setupTestDB(t)

	body := `{"nickname":"Alice","message":"Hello world","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/guestbook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateGuestbookEntry(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	var entry models.GuestbookEntry
	if err := json.NewDecoder(rr.Body).Decode(&entry); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if entry.Nickname != "Alice" {
		t.Errorf("Nickname = %q, want %q", entry.Nickname, "Alice")
	}
}

func TestCreateGuestbookEntryValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty nickname", `{"nickname":"","message":"Hi","password":"secret123"}`},
		{"empty message", `{"nickname":"Alice","message":"","password":"secret123"}`},
		{"message too long", fmt.Sprintf(`{"nickname":"Alice","message":"%s","password":"secret123"}`, strings.Repeat("x", 501))},
		{"password too short", `{"nickname":"Alice","message":"Hi","password":"ab"}`},
		{"password too long", fmt.Sprintf(`{"nickname":"Alice","message":"Hi","password":"%s"}`, strings.Repeat("x", 31))},
		{"invalid json", `not json`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestDB(t)
			req := httptest.NewRequest(http.MethodPost, "/api/guestbook", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			CreateGuestbookEntry(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d for %s", rr.Code, http.StatusBadRequest, tt.name)
			}
		})
	}
}

func TestVerifyGuestbookPasswordCorrect(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Hi", "mypassword")

	body := `{"password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/guestbook/"+fmt.Sprint(entry.ID)+"/verify", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	VerifyGuestbookPassword(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestVerifyGuestbookPasswordWrong(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Hi", "mypassword")

	body := `{"password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/guestbook/"+fmt.Sprint(entry.ID)+"/verify", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	VerifyGuestbookPassword(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestVerifyGuestbookPasswordNotFound(t *testing.T) {
	setupTestDB(t)

	body := `{"password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/guestbook/99999/verify", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "99999"})
	rr := httptest.NewRecorder()

	VerifyGuestbookPassword(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestVerifyGuestbookPasswordInvalidID(t *testing.T) {
	body := `{"password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/guestbook/abc/verify", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	VerifyGuestbookPassword(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestUpdateGuestbookEntrySuccess(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Old message", "mypassword")

	body := `{"message":"New message","password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPut, "/api/guestbook/"+fmt.Sprint(entry.ID), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	UpdateGuestbookEntry(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var updated models.GuestbookEntry
	if err := json.NewDecoder(rr.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if updated.Message != "New message" {
		t.Errorf("Message = %q, want %q", updated.Message, "New message")
	}
}

func TestUpdateGuestbookEntryWrongPassword(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Old", "mypassword")

	body := `{"message":"New","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPut, "/api/guestbook/"+fmt.Sprint(entry.ID), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	UpdateGuestbookEntry(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestUpdateGuestbookEntryNotFound(t *testing.T) {
	setupTestDB(t)

	body := `{"message":"New","password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPut, "/api/guestbook/99999", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "99999"})
	rr := httptest.NewRecorder()

	UpdateGuestbookEntry(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestUpdateGuestbookEntryInvalidMessage(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Old", "mypassword")

	body := `{"message":"","password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPut, "/api/guestbook/"+fmt.Sprint(entry.ID), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	UpdateGuestbookEntry(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestDeleteGuestbookEntrySuccess(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Goodbye", "mypassword")

	body := `{"password":"mypassword"}`
	req := httptest.NewRequest(http.MethodDelete, "/api/guestbook/"+fmt.Sprint(entry.ID), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	DeleteGuestbookEntry(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestDeleteGuestbookEntryWrongPassword(t *testing.T) {
	setupTestDB(t)

	entry := createTestEntry(t, "Alice", "Goodbye", "mypassword")

	body := `{"password":"wrong"}`
	req := httptest.NewRequest(http.MethodDelete, "/api/guestbook/"+fmt.Sprint(entry.ID), strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": fmt.Sprint(entry.ID)})
	rr := httptest.NewRecorder()

	DeleteGuestbookEntry(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestDeleteGuestbookEntryNotFound(t *testing.T) {
	setupTestDB(t)

	body := `{"password":"mypassword"}`
	req := httptest.NewRequest(http.MethodDelete, "/api/guestbook/99999", strings.NewReader(body))
	req = newChiContext(req, map[string]string{"id": "99999"})
	rr := httptest.NewRecorder()

	DeleteGuestbookEntry(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}
