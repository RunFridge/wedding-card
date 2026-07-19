package database

import "log"

func RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS guestbook_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT NOT NULL,
			message TEXT NOT NULL,
			password_hash TEXT DEFAULT '',
			ip TEXT,
			hidden INTEGER DEFAULT 0,
			evaluated INTEGER DEFAULT 0,
			secret INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS game_scores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT NOT NULL,
			time_ms INTEGER NOT NULL,
			ip TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_guestbook_hidden ON guestbook_entries(hidden)`,
		`CREATE INDEX IF NOT EXISTS idx_game_scores_time ON game_scores(time_ms ASC)`,
		`CREATE TABLE IF NOT EXISTS photo_uploads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			upload_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			ip_address TEXT,
			hashname TEXT NOT NULL,
			original_hashname TEXT DEFAULT '',
			original_filename TEXT NOT NULL,
			thumbnail TEXT DEFAULT '',
			password_hash TEXT DEFAULT '',
			hidden INTEGER DEFAULT 1,
			evaluated INTEGER DEFAULT 0,
			hearts INTEGER DEFAULT 0
		)`,
		`CREATE INDEX IF NOT EXISTS idx_photo_uploads_hidden ON photo_uploads(hidden)`,
		`CREATE TABLE IF NOT EXISTS wedding_config_overrides (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS asset_photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL DEFAULT '',
			hashname TEXT NOT NULL,
			thumb_hashname TEXT NOT NULL,
			original_filename TEXT NOT NULL,
			use_for_game INTEGER DEFAULT 0,
			is_main_photo INTEGER DEFAULT 0,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_asset_photos_game ON asset_photos(use_for_game)`,
		`CREATE TABLE IF NOT EXISTS hall_of_fame (id INTEGER PRIMARY KEY AUTOINCREMENT, nickname TEXT NOT NULL, ip TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE TABLE IF NOT EXISTS page_views (date TEXT PRIMARY KEY, count INTEGER DEFAULT 0)`,
		`CREATE TABLE IF NOT EXISTS game_beats (date TEXT PRIMARY KEY, count INTEGER DEFAULT 0)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return err
		}
	}

	log.Println("Database migrations completed")
	return nil
}
