package models

import (
	"database/sql"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type HallOfFameEntry struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminHallOfFameEntry struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

func GetHallOfFameEntries() ([]HallOfFameEntry, error) {
	rows, err := database.DB.Query(`
		SELECT id, nickname, created_at
		FROM hall_of_fame
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []HallOfFameEntry
	for rows.Next() {
		var e HallOfFameEntry
		if err := rows.Scan(&e.ID, &e.Nickname, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []HallOfFameEntry{}
	}

	return entries, rows.Err()
}

func GetAllHallOfFameEntries() ([]AdminHallOfFameEntry, error) {
	rows, err := database.DB.Query(`
		SELECT id, nickname, ip, created_at
		FROM hall_of_fame
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AdminHallOfFameEntry
	for rows.Next() {
		var e AdminHallOfFameEntry
		if err := rows.Scan(&e.ID, &e.Nickname, &e.IP, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []AdminHallOfFameEntry{}
	}

	return entries, rows.Err()
}

func CreateHallOfFameEntry(nickname, ip string) (*HallOfFameEntry, error) {
	result, err := database.DB.Exec(`
		INSERT INTO hall_of_fame (nickname, ip)
		VALUES (?, ?)
	`, nickname, ip)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &HallOfFameEntry{
		ID:        id,
		Nickname:  nickname,
		CreatedAt: time.Now(),
	}, nil
}

func DeleteHallOfFameEntry(id int64) error {
	result, err := database.DB.Exec(`DELETE FROM hall_of_fame WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
