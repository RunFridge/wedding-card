package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
)

func GetHallOfFame(w http.ResponseWriter, r *http.Request) {
	entries, err := models.GetHallOfFameEntries()
	if err != nil {
		http.Error(w, "Failed to fetch hall of fame", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func CreateHallOfFameEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		errorJSON(w, http.StatusBadRequest, "nickname_required", "Nickname is required")
		return
	}
	if len([]rune(req.Nickname)) > 20 {
		errorJSON(w, http.StatusBadRequest, "nickname_length", "Nickname must be 20 characters or less")
		return
	}

	ip := config.IPHash(r.RemoteAddr)
	entry, err := models.CreateHallOfFameEntry(req.Nickname, ip)
	if err != nil {
		http.Error(w, "Failed to create hall of fame entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func AdminGetHallOfFame(w http.ResponseWriter, r *http.Request) {
	entries, err := models.GetAllHallOfFameEntries()
	if err != nil {
		http.Error(w, "Failed to fetch hall of fame", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func AdminDeleteHallOfFameEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	if err := models.DeleteHallOfFameEntry(id); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Entry not found")
			return
		}
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
