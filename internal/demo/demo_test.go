package demo

import (
	"bytes"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/handlers"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var storageDir string

func TestMain(m *testing.M) {
	os.Setenv("DEMO_MODE", "1")
	config.Load()
	config.Cfg.BcryptCost = bcrypt.MinCost

	if err := database.Init(":memory:"); err != nil {
		panic("failed to init test db: " + err.Error())
	}

	var err error
	storageDir, err = os.MkdirTemp("", "demo-storage-")
	if err != nil {
		panic(err)
	}
	local, err := storage.NewLocalClient(storageDir)
	if err != nil {
		panic(err)
	}
	handlers.Store = local

	code := m.Run()
	database.Close()
	os.RemoveAll(storageDir)
	os.Exit(code)
}

func tableCount(t *testing.T, table string) int {
	t.Helper()
	var n int
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

func storageFiles(t *testing.T) []string {
	t.Helper()
	var files []string
	filepath.WalkDir(storageDir, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func TestResetSeedsEverything(t *testing.T) {
	if err := Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	config.Cfg.BcryptCost = bcrypt.MinCost

	for _, table := range []string{"guestbook_entries", "game_scores", "photo_uploads", "asset_photos", "hall_of_fame", "page_views", "game_beats"} {
		if tableCount(t, table) == 0 {
			t.Errorf("table %s is empty after Reset", table)
		}
	}

	if marker, _ := models.GetSingleConfigOverride(seededMarkerKey); marker != "true" {
		t.Errorf("seeded marker = %q, want true", marker)
	}
	if err := bcrypt.CompareHashAndPassword(config.Cfg.AdminPasswordHash, []byte(Password)); err != nil {
		t.Errorf("admin password is not the demo password: %v", err)
	}
	if config.Cfg.SetupRequired {
		t.Error("setup still required after Reset")
	}

	overrides, err := models.GetConfigOverrides()
	if err != nil {
		t.Fatalf("GetConfigOverrides: %v", err)
	}
	if overrides["groom_kor_name"] == "" {
		t.Error("wedding config not seeded")
	}

	var gameAssets int
	database.DB.QueryRow(`SELECT COUNT(*) FROM asset_photos WHERE use_for_game = 1`).Scan(&gameAssets)
	if gameAssets < gameAssetPhotoCount {
		t.Errorf("game asset photos = %d, want >= %d", gameAssets, gameAssetPhotoCount)
	}
	var mainPhotos int
	database.DB.QueryRow(`SELECT COUNT(*) FROM asset_photos WHERE is_main_photo = 1`).Scan(&mainPhotos)
	if mainPhotos != 1 {
		t.Errorf("main photos = %d, want 1", mainPhotos)
	}

	var hiddenPhotos int
	database.DB.QueryRow(`SELECT COUNT(*) FROM photo_uploads WHERE hidden = 1`).Scan(&hiddenPhotos)
	if hiddenPhotos != 0 {
		t.Errorf("hidden guest photos = %d, want 0", hiddenPhotos)
	}

	if len(storageFiles(t)) == 0 {
		t.Error("no files in storage after Reset")
	}
}

func TestResetRestoresMutatedState(t *testing.T) {
	if err := Reset(); err != nil {
		t.Fatalf("first Reset: %v", err)
	}
	config.Cfg.BcryptCost = bcrypt.MinCost
	guestbookCount := tableCount(t, "guestbook_entries")
	oldFiles := storageFiles(t)

	if _, err := models.CreateGuestbookEntry("불청객", "이 글은 리셋되어야 한다", "ip", "hash", false); err != nil {
		t.Fatalf("CreateGuestbookEntry: %v", err)
	}
	if err := models.SetSingleConfigOverride("groom_kor_name", "변경된이름"); err != nil {
		t.Fatalf("SetSingleConfigOverride: %v", err)
	}

	if err := Reset(); err != nil {
		t.Fatalf("second Reset: %v", err)
	}
	config.Cfg.BcryptCost = bcrypt.MinCost

	if got := tableCount(t, "guestbook_entries"); got != guestbookCount {
		t.Errorf("guestbook count after reset = %d, want %d", got, guestbookCount)
	}
	overrides, _ := models.GetConfigOverrides()
	if overrides["groom_kor_name"] == "변경된이름" {
		t.Error("mutated override survived Reset")
	}
	for _, f := range oldFiles {
		if _, err := os.Stat(f); err == nil {
			t.Errorf("old storage file %s survived Reset", f)
		}
	}
	if len(storageFiles(t)) == 0 {
		t.Error("no new storage files after Reset")
	}
}

func TestPlaceholdersAreDistinct(t *testing.T) {
	a, err := placeholderJPEG(0)
	if err != nil {
		t.Fatalf("placeholderJPEG(0): %v", err)
	}
	b, err := placeholderJPEG(1)
	if err != nil {
		t.Fatalf("placeholderJPEG(1): %v", err)
	}
	if _, err := jpeg.Decode(bytes.NewReader(a)); err != nil {
		t.Errorf("placeholder 0 is not a valid JPEG: %v", err)
	}
	if _, err := jpeg.Decode(bytes.NewReader(b)); err != nil {
		t.Errorf("placeholder 1 is not a valid JPEG: %v", err)
	}
	if bytes.Equal(a, b) {
		t.Error("placeholders 0 and 1 are identical")
	}
}
