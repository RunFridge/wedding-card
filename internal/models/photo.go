package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type PhotoUpload struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	UploadDate       time.Time `json:"upload_date"`
	URL              string    `json:"url"`
	Thumbnail        string    `json:"thumbnail"`
	Hearts           int64     `json:"hearts"`
	IP               string    `json:"-"`
	Hashname         string    `json:"-"`
	OriginalFilename string    `json:"-"`
	Hidden           bool      `json:"-"`
}

type AdminPhotoUpload struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	UploadDate       time.Time `json:"upload_date"`
	IP               string    `json:"ip_address"`
	Hashname         string    `json:"hashname"`
	OriginalHashname string    `json:"original_hashname"`
	OriginalFilename string    `json:"original_filename"`
	Hidden           bool      `json:"hidden"`
	Evaluated        bool      `json:"evaluated"`
	Hearts           int64     `json:"hearts"`
	URL              string    `json:"url,omitempty"`
	OriginalURL      string    `json:"original_url,omitempty"`
}

func CreatePhotoUpload(name, ip, hashname, originalHashname, originalFilename, thumbnail, passwordHash string) (int64, error) {
	result, err := database.DB.Exec(`
		INSERT INTO photo_uploads (name, ip_address, hashname, original_hashname, original_filename, thumbnail, password_hash)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, name, ip, hashname, originalHashname, originalFilename, thumbnail, passwordHash)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetVisiblePhotoUploads(sort string) ([]PhotoUpload, error) {
	return GetVisiblePhotoUploadsPaginated(sort, 0, 0)
}

func GetVisiblePhotoUploadsPaginated(sort string, offset, limit int) ([]PhotoUpload, error) {
	orderClause := "ORDER BY upload_date DESC"
	if sort == "popular" {
		orderClause = "ORDER BY hearts DESC, upload_date DESC"
	}

	query := `SELECT id, name, upload_date, hashname, thumbnail, hearts
		FROM photo_uploads
		WHERE hidden = 0
		` + orderClause
	args := []any{}
	if limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []PhotoUpload
	for rows.Next() {
		var p PhotoUpload
		if err := rows.Scan(&p.ID, &p.Name, &p.UploadDate, &p.Hashname, &p.Thumbnail, &p.Hearts); err != nil {
			return nil, err
		}
		uploads = append(uploads, p)
	}
	if uploads == nil {
		uploads = []PhotoUpload{}
	}
	return uploads, rows.Err()
}

func GetAllPhotoUploads() ([]AdminPhotoUpload, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, upload_date, ip_address, hashname, original_hashname, original_filename, hidden, evaluated, hearts
		FROM photo_uploads
		ORDER BY upload_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []AdminPhotoUpload
	for rows.Next() {
		var p AdminPhotoUpload
		if err := rows.Scan(&p.ID, &p.Name, &p.UploadDate, &p.IP, &p.Hashname, &p.OriginalHashname, &p.OriginalFilename, &p.Hidden, &p.Evaluated, &p.Hearts); err != nil {
			return nil, err
		}
		uploads = append(uploads, p)
	}
	if uploads == nil {
		uploads = []AdminPhotoUpload{}
	}
	return uploads, rows.Err()
}

func GetPhotoUploadByID(id int64) (*AdminPhotoUpload, error) {
	var p AdminPhotoUpload
	err := database.DB.QueryRow(`
		SELECT id, name, upload_date, ip_address, hashname, original_hashname, original_filename, hidden, evaluated
		FROM photo_uploads
		WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.UploadDate, &p.IP, &p.Hashname, &p.OriginalHashname, &p.OriginalFilename, &p.Hidden, &p.Evaluated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func SetPhotoUploadHidden(id int64, hidden bool) error {
	result, err := database.DB.Exec(`
		UPDATE photo_uploads SET hidden = ? WHERE id = ?
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

func DeletePhotoUpload(id int64) error {
	result, err := database.DB.Exec(`
		DELETE FROM photo_uploads WHERE id = ?
	`, id)
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

func SetPhotoEvaluated(id int64, evaluated bool, hidden bool) error {
	_, err := database.DB.Exec(`
		UPDATE photo_uploads SET evaluated = ?, hidden = ? WHERE id = ?
	`, evaluated, hidden, id)
	return err
}

func GetPhotoHashnameByID(id int64) (string, error) {
	var hashname string
	err := database.DB.QueryRow(`SELECT hashname FROM photo_uploads WHERE id = ?`, id).Scan(&hashname)
	return hashname, err
}

func GetUnevaluatedPhotoIDs() ([]int64, error) {
	rows, err := database.DB.Query(`SELECT id FROM photo_uploads WHERE evaluated = 0 AND hidden = 1 AND hashname != ''`)
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

func GetPhotoPasswordHash(id int64) (string, error) {
	var hash string
	err := database.DB.QueryRow(`SELECT password_hash FROM photo_uploads WHERE id = ?`, id).Scan(&hash)
	return hash, err
}

func IncrementPhotoHearts(updates map[int64]int64) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE photo_uploads SET hearts = hearts + ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, count := range updates {
		if _, err := stmt.Exec(count, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func ResetPhotoHearts(id int64) error {
	result, err := database.DB.Exec(`UPDATE photo_uploads SET hearts = 0 WHERE id = ?`, id)
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

func GetTotalHearts() (int64, error) {
	var total int64
	err := database.DB.QueryRow(`SELECT COALESCE(SUM(hearts), 0) FROM photo_uploads WHERE hidden = 0`).Scan(&total)
	return total, err
}

func GetPhotoHearts(ids []int64) (map[int64]int64, error) {
	if len(ids) == 0 {
		return map[int64]int64{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT id, hearts FROM photo_uploads WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int64)
	for rows.Next() {
		var id, hearts int64
		if err := rows.Scan(&id, &hearts); err != nil {
			return nil, err
		}
		result[id] = hearts
	}
	return result, rows.Err()
}
