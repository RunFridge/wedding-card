package models

import (
	"database/sql"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

type AssetPhoto struct {
	ID               int64     `json:"id"`
	Label            string    `json:"label"`
	Hashname         string    `json:"hashname"`
	ThumbHashname    string    `json:"thumb_hashname"`
	OriginalFilename string    `json:"original_filename"`
	UseForGame       bool      `json:"use_for_game"`
	IsMainPhoto      bool      `json:"is_main_photo"`
	SortOrder        int       `json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	URL              string    `json:"url,omitempty"`
	ThumbnailURL     string    `json:"thumbnail_url,omitempty"`
}

func CreateAssetPhoto(label, hashname, thumbHashname, originalFilename string) (int64, error) {
	result, err := database.DB.Exec(`
		INSERT INTO asset_photos (label, hashname, thumb_hashname, original_filename)
		VALUES (?, ?, ?, ?)
	`, label, hashname, thumbHashname, originalFilename)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetAllAssetPhotos() ([]AssetPhoto, error) {
	rows, err := database.DB.Query(`
		SELECT id, label, hashname, thumb_hashname, original_filename, use_for_game, is_main_photo, sort_order, created_at
		FROM asset_photos
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []AssetPhoto
	for rows.Next() {
		var p AssetPhoto
		if err := rows.Scan(&p.ID, &p.Label, &p.Hashname, &p.ThumbHashname, &p.OriginalFilename, &p.UseForGame, &p.IsMainPhoto, &p.SortOrder, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	if photos == nil {
		photos = []AssetPhoto{}
	}
	return photos, rows.Err()
}

func GetGameAssetPhotos() ([]AssetPhoto, error) {
	rows, err := database.DB.Query(`
		SELECT id, label, hashname, thumb_hashname, original_filename, use_for_game, is_main_photo, sort_order, created_at
		FROM asset_photos
		WHERE use_for_game = 1
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []AssetPhoto
	for rows.Next() {
		var p AssetPhoto
		if err := rows.Scan(&p.ID, &p.Label, &p.Hashname, &p.ThumbHashname, &p.OriginalFilename, &p.UseForGame, &p.IsMainPhoto, &p.SortOrder, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	if photos == nil {
		photos = []AssetPhoto{}
	}
	return photos, rows.Err()
}

func GetMainAssetPhoto() (*AssetPhoto, error) {
	var p AssetPhoto
	err := database.DB.QueryRow(`
		SELECT id, label, hashname, thumb_hashname, original_filename, use_for_game, is_main_photo, sort_order, created_at
		FROM asset_photos
		WHERE is_main_photo = 1
		LIMIT 1
	`).Scan(&p.ID, &p.Label, &p.Hashname, &p.ThumbHashname, &p.OriginalFilename, &p.UseForGame, &p.IsMainPhoto, &p.SortOrder, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func SetAssetPhotoGameFlag(id int64, useForGame bool) error {
	result, err := database.DB.Exec(`UPDATE asset_photos SET use_for_game = ? WHERE id = ?`, useForGame, id)
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

func SetAssetPhotoAsMain(id int64) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE asset_photos SET is_main_photo = 0 WHERE is_main_photo = 1`); err != nil {
		return err
	}

	result, err := tx.Exec(`UPDATE asset_photos SET is_main_photo = 1 WHERE id = ?`, id)
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

	return tx.Commit()
}

func ClearMainAssetPhoto(id int64) error {
	result, err := database.DB.Exec(`UPDATE asset_photos SET is_main_photo = 0 WHERE id = ?`, id)
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

func GetAssetPhotoByID(id int64) (*AssetPhoto, error) {
	var p AssetPhoto
	err := database.DB.QueryRow(`
		SELECT id, label, hashname, thumb_hashname, original_filename, use_for_game, is_main_photo, sort_order, created_at
		FROM asset_photos
		WHERE id = ?
	`, id).Scan(&p.ID, &p.Label, &p.Hashname, &p.ThumbHashname, &p.OriginalFilename, &p.UseForGame, &p.IsMainPhoto, &p.SortOrder, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func DeleteAssetPhoto(id int64) (*AssetPhoto, error) {
	photo, err := GetAssetPhotoByID(id)
	if err != nil {
		return nil, err
	}

	if _, err := database.DB.Exec(`DELETE FROM asset_photos WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return photo, nil
}

func UpdateAssetPhoto(id int64, label string, sortOrder int) error {
	result, err := database.DB.Exec(`UPDATE asset_photos SET label = ?, sort_order = ? WHERE id = ?`, label, sortOrder, id)
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

func CountGameAssetPhotos() (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM asset_photos WHERE use_for_game = 1`).Scan(&count)
	return count, err
}
