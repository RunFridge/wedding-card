package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/models"
)

func TestGetHallOfFameEmpty(t *testing.T) {
	setupTestDB(t)
	cleanHallOfFame(t)

	req := httptest.NewRequest(http.MethodGet, "/api/hall-of-fame", nil)
	rr := httptest.NewRecorder()

	GetHallOfFame(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var entries []models.HallOfFameEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetHallOfFameWithEntries(t *testing.T) {
	setupTestDB(t)
	cleanHallOfFame(t)

	models.CreateHallOfFameEntry("Alice", "10.0.0.1")
	models.CreateHallOfFameEntry("Bob", "10.0.0.2")

	req := httptest.NewRequest(http.MethodGet, "/api/hall-of-fame", nil)
	rr := httptest.NewRecorder()

	GetHallOfFame(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var entries []models.HallOfFameEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestCreateHallOfFameEntrySuccess(t *testing.T) {
	setupTestDB(t)
	cleanHallOfFame(t)

	body := `{"nickname":"TestHero"}`
	req := httptest.NewRequest(http.MethodPost, "/api/hall-of-fame", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateHallOfFameEntry(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	var entry models.HallOfFameEntry
	if err := json.NewDecoder(rr.Body).Decode(&entry); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if entry.Nickname != "TestHero" {
		t.Errorf("Nickname = %q, want %q", entry.Nickname, "TestHero")
	}
}

func TestCreateHallOfFameEntryEmpty(t *testing.T) {
	setupTestDB(t)
	cleanHallOfFame(t)

	body := `{"nickname":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/hall-of-fame", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateHallOfFameEntry(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCreateHallOfFameEntryTooLong(t *testing.T) {
	setupTestDB(t)
	cleanHallOfFame(t)

	longName := strings.Repeat("a", 21)
	body := `{"nickname":"` + longName + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/hall-of-fame", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateHallOfFameEntry(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func cleanHallOfFame(t *testing.T) {
	t.Helper()
	if _, err := database.DB.Exec("DELETE FROM hall_of_fame"); err != nil {
		t.Fatalf("failed to clean hall_of_fame: %v", err)
	}
}
