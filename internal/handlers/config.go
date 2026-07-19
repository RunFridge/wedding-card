package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
)

func GetWeddingConfig(w http.ResponseWriter, r *http.Request) {
	merged := config.Cfg.Wedding

	overrides, err := models.GetConfigOverrides()
	if err != nil {
		log.Printf("Failed to load config overrides: %v", err)
	} else if len(overrides) > 0 {
		merged = config.ApplyOverrides(merged, overrides)
	}

	raw, err := json.Marshal(merged)
	if err != nil {
		http.Error(w, "Failed to encode config", http.StatusInternalServerError)
		return
	}

	var result map[string]any
	json.Unmarshal(raw, &result)
	delete(result, "moderation_thresholds")

	var mainPhotoURL string
	mainPhoto, err := models.GetMainAssetPhoto()
	if err == nil {
		mainPhotoURL = AssetURL(mainPhoto.Hashname)
	} else if err != sql.ErrNoRows {
		log.Printf("Failed to fetch main asset photo: %v", err)
	}
	result["main_photo_url"] = mainPhotoURL

	gameCount, err := models.CountGameAssetPhotos()
	if err != nil {
		log.Printf("Failed to count game asset photos: %v", err)
	}
	result["game_photo_count"] = gameCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
