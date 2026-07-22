package demo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/handlers"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/bcrypt"
)

const (
	Password        = "demo_1234!"
	seededMarkerKey = "sys:demo_seeded"
)

var wipedTables = []string{
	"guestbook_entries",
	"game_scores",
	"photo_uploads",
	"asset_photos",
	"hall_of_fame",
	"page_views",
	"game_beats",
	"wedding_config_overrides",
}

func Init(cronSpec string) error {
	schedule, err := cron.ParseStandard(cronSpec)
	if err != nil {
		return fmt.Errorf("invalid DEMO_RESET_CRON %q: %w", cronSpec, err)
	}

	if err := bootstrap(); err != nil {
		return err
	}

	scheduler := cron.New()
	if _, err := scheduler.AddFunc(cronSpec, runScheduledReset); err != nil {
		return err
	}
	scheduler.Start()

	log.Println("========================================")
	log.Println("  DEMO MODE ENABLED")
	log.Printf("  Admin password: %s", Password)
	log.Printf("  Next reset: %s", schedule.Next(time.Now()).Format(time.RFC3339))
	log.Println("========================================")
	return nil
}

func runScheduledReset() {
	log.Println("Demo reset starting")
	if err := Reset(); err != nil {
		log.Printf("Demo reset failed: %v", err)
		return
	}
	log.Println("Demo reset completed")
}

func bootstrap() error {
	if err := forceAdminState(); err != nil {
		return err
	}
	if marker, _ := models.GetSingleConfigOverride(seededMarkerKey); marker == "true" {
		return nil
	}
	if err := seed(); err != nil {
		return err
	}
	return models.SetSingleConfigOverride(seededMarkerKey, "true")
}

// ponytail: requests during the few seconds of a reset see empty data; add an
// RWMutex gate in middleware if that ever matters for a demo site.
func Reset() error {
	if err := wipe(); err != nil {
		return err
	}
	config.Load()
	if err := config.InitIPHashKey(); err != nil {
		return err
	}
	if err := forceAdminState(); err != nil {
		return err
	}
	if handlers.ReinitStorage != nil {
		if err := handlers.ReinitStorage("", "", "", "", ""); err != nil {
			return err
		}
	}
	if handlers.UpdateMapCSP != nil {
		handlers.UpdateMapCSP("")
	}
	if err := seed(); err != nil {
		return err
	}
	return models.SetSingleConfigOverride(seededMarkerKey, "true")
}

func forceAdminState() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), config.Cfg.BcryptCost)
	if err != nil {
		return err
	}
	if err := models.SetSingleConfigOverride("admin_password_hash", string(hash)); err != nil {
		return err
	}
	config.Cfg.AdminPasswordHash = hash
	config.Cfg.AdminPasswordNeedsChange = false
	if err := config.CompleteSetup(); err != nil {
		return err
	}
	config.RemovePasswordFile()
	return nil
}

func wipe() error {
	deleteStoredObjects()

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, table := range wipedTables {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func deleteStoredObjects() {
	if handlers.Store == nil {
		return
	}
	ctx := context.Background()
	for _, key := range storedObjectKeys() {
		if err := handlers.Store.Delete(ctx, key); err != nil {
			log.Printf("Failed to delete stored object %s: %v", key, err)
		}
	}
}

func storedObjectKeys() []string {
	var keys []string
	collect := func(query, prefix string) {
		rows, err := database.DB.Query(query)
		if err != nil {
			log.Printf("Failed to list stored objects: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var first, second string
			if err := rows.Scan(&first, &second); err != nil {
				continue
			}
			if first != "" {
				keys = append(keys, prefix+first)
			}
			if second != "" {
				keys = append(keys, prefix+second)
			}
		}
	}
	collect(`SELECT hashname, original_hashname FROM photo_uploads`, "photos/")
	collect(`SELECT hashname, thumb_hashname FROM asset_photos`, "assets/")
	return keys
}
