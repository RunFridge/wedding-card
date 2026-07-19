package models

import (
	"database/sql"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type GuestbookEntry struct {
	ID           int64     `json:"id"`
	Nickname     string    `json:"nickname"`
	Message      string    `json:"message"`
	Secret       bool      `json:"secret"`
	IP           string    `json:"-"`
	Hidden       bool      `json:"-"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateGuestbookRequest struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	Password string `json:"password"`
	Secret   bool   `json:"secret"`
}

type UpdateGuestbookRequest struct {
	Message  string `json:"message"`
	Password string `json:"password"`
}

type DeleteGuestbookRequest struct {
	Password string `json:"password"`
}

type AdminGuestbookEntry struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	Message   string    `json:"message"`
	Secret    bool      `json:"secret"`
	IP        string    `json:"ip"`
	Hidden    bool      `json:"hidden"`
	Evaluated bool      `json:"evaluated"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllGuestbookEntries() ([]AdminGuestbookEntry, error) {
	rows, err := database.DB.Query(`
		SELECT id, nickname, message, secret, ip, hidden, evaluated, created_at
		FROM guestbook_entries
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AdminGuestbookEntry
	for rows.Next() {
		var e AdminGuestbookEntry
		if err := rows.Scan(&e.ID, &e.Nickname, &e.Message, &e.Secret, &e.IP, &e.Hidden, &e.Evaluated, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []AdminGuestbookEntry{}
	}

	return entries, rows.Err()
}

func SetGuestbookEntryHidden(id int64, hidden bool) error {
	result, err := database.DB.Exec(`
		UPDATE guestbook_entries SET hidden = ? WHERE id = ?
	`, hidden, id)
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

func GetVisibleGuestbookEntries() ([]GuestbookEntry, error) {
	return GetVisibleGuestbookEntriesPaginated(0, 0)
}

func GetVisibleGuestbookEntriesPaginated(cursor int64, limit int) ([]GuestbookEntry, error) {
	query := `SELECT id, nickname, message, secret, created_at FROM guestbook_entries WHERE hidden = 0`
	args := []any{}
	if cursor > 0 {
		query += ` AND id < ?`
		args = append(args, cursor)
	}
	query += ` ORDER BY id DESC`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []GuestbookEntry
	for rows.Next() {
		var e GuestbookEntry
		if err := rows.Scan(&e.ID, &e.Nickname, &e.Message, &e.Secret, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []GuestbookEntry{}
	}

	return entries, rows.Err()
}

func GetGuestbookEntryByID(id int64) (*GuestbookEntry, error) {
	var e GuestbookEntry
	err := database.DB.QueryRow(`
		SELECT id, nickname, message, secret, password_hash, created_at
		FROM guestbook_entries
		WHERE id = ?
	`, id).Scan(&e.ID, &e.Nickname, &e.Message, &e.Secret, &e.PasswordHash, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func CreateGuestbookEntry(nickname, message, ip, passwordHash string, secret bool) (*GuestbookEntry, error) {
	result, err := database.DB.Exec(`
		INSERT INTO guestbook_entries (nickname, message, ip, password_hash, secret)
		VALUES (?, ?, ?, ?, ?)
	`, nickname, message, ip, passwordHash, secret)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &GuestbookEntry{
		ID:        id,
		Nickname:  nickname,
		Message:   message,
		Secret:    secret,
		CreatedAt: time.Now(),
	}, nil
}

func UpdateGuestbookEntry(id int64, message string) error {
	_, err := database.DB.Exec(`
		UPDATE guestbook_entries SET message = ? WHERE id = ?
	`, message, id)
	return err
}

func DeleteGuestbookEntry(id int64) error {
	_, err := database.DB.Exec(`
		DELETE FROM guestbook_entries WHERE id = ?
	`, id)
	return err
}

func SetGuestbookEvaluated(id int64, evaluated bool, hidden bool) error {
	_, err := database.DB.Exec(`
		UPDATE guestbook_entries SET evaluated = ?, hidden = ? WHERE id = ?
	`, evaluated, hidden, id)
	return err
}

func GetGuestbookContentByID(id int64) (nickname, message string, err error) {
	err = database.DB.QueryRow(`
		SELECT nickname, message FROM guestbook_entries WHERE id = ?
	`, id).Scan(&nickname, &message)
	return
}

func GetUnevaluatedGuestbookIDs() ([]int64, error) {
	rows, err := database.DB.Query(`SELECT id FROM guestbook_entries WHERE evaluated = 0 AND secret = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func ResetGuestbookEvaluated(id int64) error {
	_, err := database.DB.Exec(`
		UPDATE guestbook_entries SET evaluated = 0, hidden = 1 WHERE id = ?
	`, id)
	return err
}
