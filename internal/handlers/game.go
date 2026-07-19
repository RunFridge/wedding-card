package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/session"
)

var nicknameRegex = regexp.MustCompile(`^[A-Z]{3}$`)

var MinGameTimeMs = 3000

var emojiPool = []struct {
	ID    string `json:"id"`
	Emoji string `json:"emoji"`
}{
	{"e1", "💑"}, {"e2", "💍"}, {"e3", "💐"}, {"e4", "🎂"},
	{"e5", "💒"}, {"e6", "❤️"}, {"e7", "💌"}, {"e8", "💓"},
	{"e9", "💕"}, {"e10", "💖"}, {"e11", "👫"}, {"e12", "💎"},
	{"e13", "🌹"}, {"e14", "🌟"}, {"e15", "🎀"}, {"e16", "🍾"},
	{"e17", "🎉"}, {"e18", "🎶"}, {"e19", "🌸"}, {"e20", "🌞"},
}

func GetGamePhotos(w http.ResponseWriter, r *http.Request) {
	gamePhotos, err := models.GetGameAssetPhotos()
	if err != nil {
		log.Printf("Failed to fetch game asset photos: %v", err)
		gamePhotos = nil
	}

	token, err := session.GameSessions.Create()
	if err != nil {
		log.Printf("Failed to create game token: %v", err)
		http.Error(w, "Failed to initialize game session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(gamePhotos) >= 6 {
		type photoEntry struct {
			ID           int64  `json:"id"`
			ThumbnailURL string `json:"thumbnail_url"`
			DetailURL    string `json:"detail_url"`
		}

		allPhotos := make([]photoEntry, len(gamePhotos))
		for i, p := range gamePhotos {
			allPhotos[i] = photoEntry{
				ID:           p.ID,
				ThumbnailURL: AssetURL(p.ThumbHashname),
				DetailURL:    AssetURL(p.Hashname),
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"type": "photo", "photos": allPhotos, "game_token": token})
	} else {
		json.NewEncoder(w).Encode(map[string]any{"type": "emoji", "photos": emojiPool, "game_token": token})
	}
}

func GetGameRankings(w http.ResponseWriter, r *http.Request) {
	scores, err := models.GetTopScores(10)
	if err != nil {
		http.Error(w, "Failed to fetch rankings", http.StatusInternalServerError)
		return
	}

	hasPlayed, _ := models.HasPlayedFromIP(config.IPHash(r.RemoteAddr))

	resetAt, _ := models.GetSingleConfigOverride("game_reset_at")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"rankings": scores, "has_played": hasPlayed, "game_reset_at": resetAt})
}

func CreateGameScore(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGameScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if req.GameToken == "" {
		errorJSON(w, http.StatusForbidden, "game_token_required", "Game token required")
		return
	}
	if ok, reason := session.GameSessions.Consume(req.GameToken, MinGameTimeMs); !ok {
		errorJSON(w, http.StatusForbidden, "invalid_game_session", "Invalid game session: "+reason)
		return
	}

	if !nicknameRegex.MatchString(req.Nickname) {
		errorJSON(w, http.StatusBadRequest, "nickname_format", "Nickname must be exactly 3 uppercase letters")
		return
	}

	if req.TimeMs < 1 || req.TimeMs > config.Cfg.GameTimerMs {
		errorJSON(w, http.StatusBadRequest, "time_range", fmt.Sprintf("Time must be between 1-%dms", config.Cfg.GameTimerMs))
		return
	}

	ip := config.IPHash(r.RemoteAddr)
	score, err := models.CreateGameScore(req.Nickname, req.TimeMs, ip)
	if err != nil {
		http.Error(w, "Failed to save score", http.StatusInternalServerError)
		return
	}

	if HeartsHub != nil {
		HeartsHub.BroadcastRankingUpdate(req.Nickname, req.TimeMs)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(score)
}

func RecordGameBeat(w http.ResponseWriter, r *http.Request) {
	go models.RecordGameBeat()
	w.WriteHeader(http.StatusNoContent)
}
