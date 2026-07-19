package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/session"
)

var ModerationQueueLen func(ctx context.Context) (int64, error)
var ShutdownFunc func()
var UpdateMapCSP func(embedProvider string)

func AdminVerify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Password == "" {
		errorJSON(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	if bcrypt.CompareHashAndPassword(config.Cfg.AdminPasswordHash, []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	token, err := session.Global.Create()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    token,
		Path:     "/api/admin",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isHTTPS(r),
		MaxAge:   86400,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":                "ok",
		"password_needs_change": config.Cfg.AdminPasswordNeedsChange,
		"setup_required":        config.Cfg.SetupRequired,
	})
}

func AdminSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":                "ok",
		"role":                  "admin",
		"password_needs_change": config.Cfg.AdminPasswordNeedsChange,
		"setup_required":        config.Cfg.SetupRequired,
	})
}

func AdminCompleteSetup(w http.ResponseWriter, r *http.Request) {
	if config.Cfg.AdminPasswordNeedsChange {
		errorJSON(w, http.StatusBadRequest, "setup_password_unchanged", "Password must be changed before completing setup")
		return
	}

	if err := config.CompleteSetup(); err != nil {
		http.Error(w, "Failed to complete setup", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if req.NewPassword == "" {
		errorJSON(w, http.StatusBadRequest, "new_password_required", "New password is required")
		return
	}

	if !config.Cfg.SetupRequired {
		if req.CurrentPassword == "" {
			errorJSON(w, http.StatusBadRequest, "current_password_required", "Current password is required")
			return
		}
		if bcrypt.CompareHashAndPassword(config.Cfg.AdminPasswordHash, []byte(req.CurrentPassword)) != nil {
			errorJSON(w, http.StatusUnauthorized, "current_password_incorrect", "Current password is incorrect")
			return
		}
	}

	if len(req.NewPassword) < 8 {
		errorJSON(w, http.StatusBadRequest, "password_min_length", "New password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), config.Cfg.BcryptCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := models.SetSingleConfigOverride("admin_password_hash", string(hash)); err != nil {
		http.Error(w, "Failed to save password", http.StatusInternalServerError)
		return
	}

	config.Cfg.AdminPasswordHash = hash
	config.Cfg.AdminPasswordNeedsChange = false
	config.RemovePasswordFile()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminGetGuestbook(w http.ResponseWriter, r *http.Request) {
	entries, err := models.GetAllGuestbookEntries()
	if err != nil {
		http.Error(w, "Failed to fetch entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func AdminToggleGuestbookVisibility(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	var req struct {
		Hidden bool `json:"hidden"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if err := models.SetGuestbookEntryHidden(id, req.Hidden); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Entry not found")
			return
		}
		http.Error(w, "Failed to update entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminGetGameRankings(w http.ResponseWriter, r *http.Request) {
	scores, err := models.GetAllGameScores(100)
	if err != nil {
		http.Error(w, "Failed to fetch rankings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func AdminDeleteGameScore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid score ID")
		return
	}

	if err := models.DeleteGameScore(id); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Score not found")
			return
		}
		http.Error(w, "Failed to delete score", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AdminPurgeGameRankings(w http.ResponseWriter, r *http.Request) {
	if err := models.DeleteAllGameScores(); err != nil {
		http.Error(w, "Failed to purge rankings", http.StatusInternalServerError)
		return
	}
	_ = models.SetSingleConfigOverride("game_reset_at", strconv.FormatInt(time.Now().UnixMilli(), 10))
	HeartsHub.BroadcastGameReset()
	w.WriteHeader(http.StatusNoContent)
}

func AdminDeleteGuestbookEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	if err := models.DeleteGuestbookEntry(id); err != nil {
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AdminGetConfig(w http.ResponseWriter, r *http.Request) {
	merged := config.Cfg.Wedding

	overrides, err := models.GetConfigOverrides()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	if len(overrides) > 0 {
		merged = config.ApplyOverrides(merged, overrides)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(merged)
}

func AdminUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var updated config.WeddingConfig
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	diff := config.DiffToOverrideMap(config.Cfg.Wedding, updated)

	if err := models.SetConfigOverrides(diff); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	if UpdateMapCSP != nil {
		merged := config.ApplyOverrides(config.Cfg.Wedding, diff)
		UpdateMapCSP(merged.MapProviders.EmbedProvider)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("admin_token"); err == nil {
		session.Global.Revoke(c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    "",
		Path:     "/api/admin",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isHTTPS(r),
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}

func isHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func AdminGetModerationStatus(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Enabled       bool                    `json:"enabled"`
		QueueDepth    int64                   `json:"queue_depth"`
		Guestbook     *models.ModerationStats `json:"guestbook"`
		Photos        *models.ModerationStats `json:"photos"`
		WsConnections int                     `json:"ws_connections"`
		TotalHearts   int64                   `json:"total_hearts"`
	}{}

	if ModerationQueueLen != nil {
		resp.Enabled = true
		depth, err := ModerationQueueLen(r.Context())
		if err == nil {
			resp.QueueDepth = depth
		}
	}

	gbStats, err := models.GetGuestbookModerationStats()
	if err != nil {
		http.Error(w, "Failed to fetch moderation stats", http.StatusInternalServerError)
		return
	}
	resp.Guestbook = gbStats

	phStats, err := models.GetPhotoModerationStats()
	if err != nil {
		http.Error(w, "Failed to fetch moderation stats", http.StatusInternalServerError)
		return
	}
	resp.Photos = phStats

	if HeartsHub != nil {
		resp.WsConnections = HeartsHub.ConnectionCount()
	}
	totalHearts, err := models.GetTotalHearts()
	if err == nil {
		resp.TotalHearts = totalHearts
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func AdminGetPageViews(w http.ResponseWriter, r *http.Request) {
	views, err := models.GetPageViews(180)
	if err != nil {
		http.Error(w, "Failed to fetch page views", http.StatusInternalServerError)
		return
	}
	if views == nil {
		views = []models.PageView{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(views)
}

func AdminGetGameBeats(w http.ResponseWriter, r *http.Request) {
	beats, err := models.GetGameBeats(180)
	if err != nil {
		http.Error(w, "Failed to fetch game beats", http.StatusInternalServerError)
		return
	}
	if beats == nil {
		beats = []models.GameBeat{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beats)
}

func AdminRestartServer(w http.ResponseWriter, r *http.Request) {
	if ShutdownFunc == nil {
		http.Error(w, "Restart not supported", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		ShutdownFunc()
	}()
}
