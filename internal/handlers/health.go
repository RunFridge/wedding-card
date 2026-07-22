package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RunFridge/wedding-card/internal/config"
)

var Version = "dev"

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":         "ok",
		"version":        Version,
		"setup_required": config.Cfg.SetupRequired,
		"demo":           config.Cfg.DemoMode,
	})
}
