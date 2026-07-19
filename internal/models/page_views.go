package models

import (
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type PageView struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func RecordPageView() {
	today := time.Now().Format("2006-01-02")
	database.DB.Exec(
		`INSERT INTO page_views (date, count) VALUES (?, 1)
		 ON CONFLICT(date) DO UPDATE SET count = count + 1`,
		today,
	)
}

func GetPageViews(days int) ([]PageView, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	rows, err := database.DB.Query(
		`SELECT date, count FROM page_views WHERE date >= ? ORDER BY date ASC`,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []PageView
	for rows.Next() {
		var v PageView
		if err := rows.Scan(&v.Date, &v.Count); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

func CleanupOldPageViews(months int) error {
	cutoff := time.Now().AddDate(0, -months, 0).Format("2006-01-02")
	_, err := database.DB.Exec(`DELETE FROM page_views WHERE date < ?`, cutoff)
	return err
}
