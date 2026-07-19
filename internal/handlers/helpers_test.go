package handlers

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/session"
)

func TestMain(m *testing.M) {
	if err := database.Init(":memory:"); err != nil {
		panic("failed to init test db: " + err.Error())
	}
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("testadmin"), 4)
	config.Cfg = &config.Config{
		BcryptCost:        4,
		GameTimerMs:       30000,
		AdminPasswordHash: adminHash,
	}
	session.Init()
	session.InitGame()
	code := m.Run()
	database.Close()
	os.Exit(code)
}

func setupTestDB(t *testing.T) {
	t.Helper()
	for _, table := range []string{"guestbook_entries", "game_scores", "photo_uploads"} {
		if _, err := database.DB.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func newChiContext(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func createTestEntry(t *testing.T, nickname, message, password string) *models.GuestbookEntry {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), config.Cfg.BcryptCost)
	if err != nil {
		t.Fatalf("bcrypt hash failed: %v", err)
	}
	entry, err := models.CreateGuestbookEntry(nickname, message, "127.0.0.1", string(hash), false)
	if err != nil {
		t.Fatalf("CreateGuestbookEntry failed: %v", err)
	}
	return entry
}
