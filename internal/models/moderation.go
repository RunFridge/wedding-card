package models

import "github.com/RunFridge/wedding-card/internal/database"

type ModerationStats struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Flagged  int `json:"flagged"`
}

func GetGuestbookModerationStats() (*ModerationStats, error) {
	var s ModerationStats
	err := database.DB.QueryRow(`
		SELECT
			IFNULL(COUNT(*), 0),
			IFNULL(SUM(CASE WHEN evaluated = 0 THEN 1 ELSE 0 END), 0),
			IFNULL(SUM(CASE WHEN evaluated = 1 AND hidden = 0 THEN 1 ELSE 0 END), 0),
			IFNULL(SUM(CASE WHEN evaluated = 1 AND hidden = 1 THEN 1 ELSE 0 END), 0)
		FROM guestbook_entries
	`).Scan(&s.Total, &s.Pending, &s.Approved, &s.Flagged)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetPhotoModerationStats() (*ModerationStats, error) {
	var s ModerationStats
	err := database.DB.QueryRow(`
		SELECT
			IFNULL(COUNT(*), 0),
			IFNULL(SUM(CASE WHEN evaluated = 0 THEN 1 ELSE 0 END), 0),
			IFNULL(SUM(CASE WHEN evaluated = 1 AND hidden = 0 THEN 1 ELSE 0 END), 0),
			IFNULL(SUM(CASE WHEN evaluated = 1 AND hidden = 1 THEN 1 ELSE 0 END), 0)
		FROM photo_uploads
	`).Scan(&s.Total, &s.Pending, &s.Approved, &s.Flagged)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
