package config

import (
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_PATH", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.DatabasePath != "./wedding.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "./wedding.db")
	}
	if cfg.BcryptCost != 10 {
		t.Errorf("BcryptCost = %d, want %d", cfg.BcryptCost, 10)
	}
	if cfg.GameTimerMs != 30000 {
		t.Errorf("GameTimerMs = %d, want %d", cfg.GameTimerMs, 30000)
	}
}

func TestLoadCustomPortAndDB(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("DATABASE_PATH", "/tmp/test.db")

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "/tmp/test.db")
	}
}

func TestLoadSetsCfgGlobal(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_PATH", "")

	cfg := Load()

	if Cfg != cfg {
		t.Error("Cfg global was not set to returned config")
	}
}

func TestCharterBusNoticeOverride(t *testing.T) {
	base := loadWeddingDefaults()

	if base.CharterBusNotice != "" {
		t.Errorf("default CharterBusNotice = %q, want empty", base.CharterBusNotice)
	}

	notice := "좌석이 만석으로 예상됩니다.\n탑승 여부에 변동이 있으신 분은 미리 연락 부탁드립니다."
	applied := ApplyOverrides(base, map[string]string{"charter_bus_notice": notice})
	if applied.CharterBusNotice != notice {
		t.Errorf("CharterBusNotice = %q, want %q", applied.CharterBusNotice, notice)
	}

	diff := DiffToOverrideMap(base, applied)
	if diff["charter_bus_notice"] != notice {
		t.Errorf("diff[charter_bus_notice] = %q, want %q", diff["charter_bus_notice"], notice)
	}

	if _, ok := DiffToOverrideMap(base, base)["charter_bus_notice"]; ok {
		t.Error("diff of unchanged config should not contain charter_bus_notice")
	}
}

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := database.Init(":memory:"); err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}
	t.Cleanup(database.Close)
}

func TestInitAdminPasswordFromDB(t *testing.T) {
	setupTestDB(t)
	Load()

	// First call auto-generates
	if err := InitAdminPassword(); err != nil {
		t.Fatalf("first InitAdminPassword() error: %v", err)
	}
	firstHash := Cfg.AdminPasswordHash

	// Second call should load from DB
	Load()
	if err := InitAdminPassword(); err != nil {
		t.Fatalf("second InitAdminPassword() error: %v", err)
	}

	if string(Cfg.AdminPasswordHash) != string(firstHash) {
		t.Error("AdminPasswordHash should match the DB-stored hash on second load")
	}
	if Cfg.AdminPasswordNeedsChange {
		t.Error("AdminPasswordNeedsChange should be false when loading from DB")
	}
}

func TestInitAdminPasswordAutoGenerate(t *testing.T) {
	setupTestDB(t)
	Load()

	if err := InitAdminPassword(); err != nil {
		t.Fatalf("InitAdminPassword() error: %v", err)
	}

	if Cfg.AdminPasswordHash == nil {
		t.Error("AdminPasswordHash should not be nil after auto-generation")
	}
	if !Cfg.AdminPasswordNeedsChange {
		t.Error("AdminPasswordNeedsChange should be true after auto-generation")
	}
}

func TestInitSetupStateExistingDeployment(t *testing.T) {
	setupTestDB(t)
	Load()

	// Simulate existing deployment: password already set (not auto-generated)
	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), 10)
	models.SetSingleConfigOverride("admin_password_hash", string(hash))
	Cfg.AdminPasswordHash = hash
	Cfg.AdminPasswordNeedsChange = false

	if err := InitSetupState(); err != nil {
		t.Fatalf("InitSetupState() error: %v", err)
	}

	if Cfg.SetupRequired {
		t.Error("SetupRequired should be false for existing deployment")
	}

	val, err := models.GetSingleConfigOverride("sys:setup_completed")
	if err != nil || val != "true" {
		t.Error("setup_completed should be auto-set for existing deployment")
	}
}

func TestInitSetupStateFreshInstall(t *testing.T) {
	setupTestDB(t)
	Load()

	if err := InitAdminPassword(); err != nil {
		t.Fatalf("InitAdminPassword() error: %v", err)
	}

	// Fresh install: password was auto-generated
	if !Cfg.AdminPasswordNeedsChange {
		t.Fatal("expected AdminPasswordNeedsChange to be true for fresh install")
	}

	if err := InitSetupState(); err != nil {
		t.Fatalf("InitSetupState() error: %v", err)
	}

	if !Cfg.SetupRequired {
		t.Error("SetupRequired should be true for fresh install")
	}
}

func TestCompleteSetup(t *testing.T) {
	setupTestDB(t)
	Load()

	Cfg.SetupRequired = true
	Cfg.AdminPasswordNeedsChange = true

	if err := CompleteSetup(); err != nil {
		t.Fatalf("CompleteSetup() error: %v", err)
	}

	if Cfg.SetupRequired {
		t.Error("SetupRequired should be false after CompleteSetup")
	}
	if Cfg.AdminPasswordNeedsChange {
		t.Error("AdminPasswordNeedsChange should be false after CompleteSetup")
	}

	val, err := models.GetSingleConfigOverride("sys:setup_completed")
	if err != nil || val != "true" {
		t.Error("setup_completed should be set in DB after CompleteSetup")
	}
}
