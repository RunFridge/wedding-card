package models

import (
	"strings"

	"github.com/RunFridge/wedding-card/internal/database"
)

func GetConfigOverrides() (map[string]string, error) {
	rows, err := database.DB.Query(`SELECT key, value FROM wedding_config_overrides`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	overrides := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		overrides[k] = v
	}
	return overrides, rows.Err()
}

func SetSingleConfigOverride(key, value string) error {
	_, err := database.DB.Exec(
		`INSERT INTO wedding_config_overrides (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value,
	)
	return err
}

func SetConfigOverrides(overrides map[string]string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM wedding_config_overrides WHERE key != 'admin_password_hash' AND key NOT LIKE 'sys:%'`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO wedding_config_overrides (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for k, v := range overrides {
		if v == "" {
			continue
		}
		if _, err := stmt.Exec(k, v); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetSingleConfigOverride(key string) (string, error) {
	var val string
	err := database.DB.QueryRow(`SELECT value FROM wedding_config_overrides WHERE key = ?`, key).Scan(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetSystemConfigOverrides() (map[string]string, error) {
	rows, err := database.DB.Query(`SELECT key, value FROM wedding_config_overrides WHERE key LIKE 'sys:%'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	overrides := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		overrides[strings.TrimPrefix(k, "sys:")] = v
	}
	return overrides, rows.Err()
}

func SetSystemConfigOverrides(overrides map[string]string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM wedding_config_overrides WHERE key LIKE 'sys:%'`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO wedding_config_overrides (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for k, v := range overrides {
		if _, err := stmt.Exec("sys:"+k, v); err != nil {
			return err
		}
	}

	return tx.Commit()
}
