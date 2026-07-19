package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/session"
)

func createGameToken(t *testing.T) string {
	t.Helper()
	token, err := session.GameSessions.Create()
	if err != nil {
		t.Fatalf("failed to create game token: %v", err)
	}
	return token
}

func TestGetGamePhotosEmojiFallback(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/game/photos", nil)
	rr := httptest.NewRecorder()

	GetGamePhotos(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var body struct {
		Type   string `json:"type"`
		Photos []struct {
			ID    string `json:"id"`
			Emoji string `json:"emoji"`
		} `json:"photos"`
		GameToken string `json:"game_token"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if body.Type != "emoji" {
		t.Errorf("type = %q, want %q", body.Type, "emoji")
	}
	if len(body.Photos) != 20 {
		t.Errorf("expected 20 emojis, got %d", len(body.Photos))
	}
	if body.GameToken == "" {
		t.Error("expected game_token in response")
	}

	seen := make(map[string]bool)
	for _, p := range body.Photos {
		if p.ID == "" {
			t.Error("emoji entry has empty ID")
		}
		if p.Emoji == "" {
			t.Error("emoji entry has empty emoji")
		}
		if seen[p.ID] {
			t.Errorf("duplicate emoji ID: %q", p.ID)
		}
		seen[p.ID] = true
	}
}

func TestGetGameRankingsEmpty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/game/rankings", nil)
	rr := httptest.NewRecorder()

	GetGameRankings(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct {
		Rankings []models.GameScore `json:"rankings"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Rankings) != 0 {
		t.Errorf("expected 0 scores, got %d", len(resp.Rankings))
	}
}

func TestGetGameRankingsOrdered(t *testing.T) {
	setupTestDB(t)

	models.CreateGameScore("AAA", 10000, "127.0.0.1")
	models.CreateGameScore("BBB", 3000, "127.0.0.1")
	models.CreateGameScore("CCC", 7000, "127.0.0.1")

	req := httptest.NewRequest(http.MethodGet, "/api/game/rankings", nil)
	rr := httptest.NewRecorder()

	GetGameRankings(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct {
		Rankings []models.GameScore `json:"rankings"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Rankings) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(resp.Rankings))
	}
	if resp.Rankings[0].Nickname != "BBB" {
		t.Errorf("expected fastest first (BBB), got %q", resp.Rankings[0].Nickname)
	}
}

func TestCreateGameScoreSuccess(t *testing.T) {
	setupTestDB(t)
	origMin := MinGameTimeMs
	MinGameTimeMs = 0
	defer func() { MinGameTimeMs = origMin }()

	token := createGameToken(t)
	body := fmt.Sprintf(`{"nickname":"ABC","time_ms":15000,"game_token":"%s"}`, token)
	req := httptest.NewRequest(http.MethodPost, "/api/game/rankings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateGameScore(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	var score models.GameScore
	if err := json.NewDecoder(rr.Body).Decode(&score); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if score.Nickname != "ABC" {
		t.Errorf("Nickname = %q, want %q", score.Nickname, "ABC")
	}
	if score.TimeMs != 15000 {
		t.Errorf("TimeMs = %d, want %d", score.TimeMs, 15000)
	}
}

func TestCreateGameScoreNicknameValidation(t *testing.T) {
	origMin := MinGameTimeMs
	MinGameTimeMs = 0
	defer func() { MinGameTimeMs = origMin }()

	tests := []struct {
		name     string
		nickname string
	}{
		{"lowercase", "abc"},
		{"too short", "AB"},
		{"too long", "ABCD"},
		{"digits", "A1B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestDB(t)
			token := createGameToken(t)
			body := fmt.Sprintf(`{"nickname":"%s","time_ms":15000,"game_token":"%s"}`, tt.nickname, token)
			req := httptest.NewRequest(http.MethodPost, "/api/game/rankings", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			CreateGameScore(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d for nickname %q", rr.Code, http.StatusBadRequest, tt.nickname)
			}
		})
	}
}

func TestCreateGameScoreTimeValidation(t *testing.T) {
	origMin := MinGameTimeMs
	MinGameTimeMs = 0
	defer func() { MinGameTimeMs = origMin }()

	tests := []struct {
		name   string
		timeMs int
		want   int
	}{
		{"zero", 0, http.StatusBadRequest},
		{"over limit", 30001, http.StatusBadRequest},
		{"min valid", 1, http.StatusCreated},
		{"max valid", 30000, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestDB(t)
			token := createGameToken(t)
			body := fmt.Sprintf(`{"nickname":"ABC","time_ms":%d,"game_token":"%s"}`, tt.timeMs, token)
			req := httptest.NewRequest(http.MethodPost, "/api/game/rankings", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			CreateGameScore(rr, req)

			if rr.Code != tt.want {
				t.Errorf("status = %d, want %d for time_ms=%d; body = %s", rr.Code, tt.want, tt.timeMs, rr.Body.String())
			}
		})
	}
}

func TestCreateGameScoreMissingToken(t *testing.T) {
	setupTestDB(t)

	body := `{"nickname":"ABC","time_ms":15000}`
	req := httptest.NewRequest(http.MethodPost, "/api/game/rankings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateGameScore(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusForbidden, rr.Body.String())
	}
}

func TestCreateGameScoreInvalidToken(t *testing.T) {
	setupTestDB(t)

	body := `{"nickname":"ABC","time_ms":15000,"game_token":"totally-bogus-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/game/rankings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateGameScore(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d; body = %s", rr.Code, http.StatusForbidden, rr.Body.String())
	}
}
