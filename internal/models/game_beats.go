package models

import (
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type GameBeat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func RecordGameBeat() {
	today := time.Now().Format("2006-01-02")
	database.DB.Exec(
		`INSERT INTO game_beats (date, count) VALUES (?, 1)
		 ON CONFLICT(date) DO UPDATE SET count = count + 1`,
		today,
	)
}

func GetGameBeats(days int) ([]GameBeat, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	rows, err := database.DB.Query(
		`SELECT date, count FROM game_beats WHERE date >= ? ORDER BY date ASC`,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beats []GameBeat
	for rows.Next() {
		var b GameBeat
		if err := rows.Scan(&b.Date, &b.Count); err != nil {
			return nil, err
		}
		beats = append(beats, b)
	}
	return beats, nil
}

func CleanupOldGameBeats(months int) error {
	cutoff := time.Now().AddDate(0, -months, 0).Format("2006-01-02")
	_, err := database.DB.Exec(`DELETE FROM game_beats WHERE date < ?`, cutoff)
	return err
}
