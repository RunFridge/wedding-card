package models

import (
	"database/sql"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type GameScore struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	TimeMs    int       `json:"time_ms"`
	IP        string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateGameScoreRequest struct {
	Nickname  string `json:"nickname"`
	TimeMs    int    `json:"time_ms"`
	GameToken string `json:"game_token"`
}

type AdminGameScore struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	TimeMs    int       `json:"time_ms"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllGameScores(limit int) ([]AdminGameScore, error) {
	rows, err := database.DB.Query(`
		SELECT id, nickname, time_ms, ip, created_at
		FROM game_scores
		ORDER BY time_ms ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []AdminGameScore
	for rows.Next() {
		var s AdminGameScore
		if err := rows.Scan(&s.ID, &s.Nickname, &s.TimeMs, &s.IP, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	if scores == nil {
		scores = []AdminGameScore{}
	}

	return scores, rows.Err()
}

func DeleteGameScore(id int64) error {
	result, err := database.DB.Exec(`DELETE FROM game_scores WHERE id = ?`, id)
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

func DeleteAllGameScores() error {
	_, err := database.DB.Exec(`DELETE FROM game_scores`)
	return err
}

func GetTopScores(limit int) ([]GameScore, error) {
	rows, err := database.DB.Query(`
		SELECT id, nickname, time_ms, created_at
		FROM game_scores
		ORDER BY time_ms ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []GameScore
	for rows.Next() {
		var s GameScore
		if err := rows.Scan(&s.ID, &s.Nickname, &s.TimeMs, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	if scores == nil {
		scores = []GameScore{}
	}

	return scores, rows.Err()
}

func HasPlayedFromIP(ip string) (bool, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM game_scores WHERE ip = ?`, ip).Scan(&count)
	return count > 0, err
}

func CreateGameScore(nickname string, timeMs int, ip string) (*GameScore, error) {
	result, err := database.DB.Exec(`
		INSERT INTO game_scores (nickname, time_ms, ip)
		VALUES (?, ?, ?)
	`, nickname, timeMs, ip)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &GameScore{
		ID:        id,
		Nickname:  nickname,
		TimeMs:    timeMs,
		CreatedAt: time.Now(),
	}, nil
}
