package models

import (
	"os"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestMain(m *testing.M) {
	if err := database.Init(":memory:"); err != nil {
		panic("failed to init test db: " + err.Error())
	}
	code := m.Run()
	database.Close()
	os.Exit(code)
}

func cleanTables(t *testing.T) {
	t.Helper()
	for _, table := range []string{"guestbook_entries", "game_scores", "photo_uploads", "asset_photos", "wedding_config_overrides", "hall_of_fame", "page_views", "game_beats"} {
		if _, err := database.DB.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}
