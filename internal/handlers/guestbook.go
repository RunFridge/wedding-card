package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
)

var OnGuestbookEntryCreated func(ctx context.Context, entryID int64)

func GetGuestbook(w http.ResponseWriter, r *http.Request) {
	cursor, _ := strconv.ParseInt(r.URL.Query().Get("cursor"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	entries, err := models.GetVisibleGuestbookEntriesPaginated(cursor, limit+1)
	if err != nil {
		http.Error(w, "Failed to fetch entries", http.StatusInternalServerError)
		return
	}

	hasMore := len(entries) > limit
	if hasMore {
		entries = entries[:limit]
	}

	for i := range entries {
		if entries[i].Secret {
			entries[i].Nickname = "비밀"
			entries[i].Message = "비밀 메시지입니다 💌"
		}
	}

	var nextCursor *int64
	if hasMore {
		nextCursor = &entries[len(entries)-1].ID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"items":       entries,
		"next_cursor": nextCursor,
	})
}

func CreateGuestbookEntry(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGuestbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Message = strings.TrimSpace(req.Message)
	req.Password = strings.TrimSpace(req.Password)

	if req.Nickname == "" {
		errorJSON(w, http.StatusBadRequest, "nickname_required", "Nickname is required")
		return
	}

	if req.Message == "" || len(req.Message) > 500 {
		errorJSON(w, http.StatusBadRequest, "message_length", "Message must be 1-500 characters")
		return
	}

	if len(req.Password) < 5 || len(req.Password) > 30 {
		errorJSON(w, http.StatusBadRequest, "password_length", "Password must be 5-30 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), config.Cfg.BcryptCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	ip := config.IPHash(r.RemoteAddr)
	entry, err := models.CreateGuestbookEntry(req.Nickname, req.Message, ip, string(hash), req.Secret)
	if err != nil {
		http.Error(w, "Failed to create entry", http.StatusInternalServerError)
		return
	}

	if entry.Secret {
		entry.Nickname = "비밀"
		entry.Message = "비밀 메시지입니다 💌"
		models.SetGuestbookEvaluated(entry.ID, true, false)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)

	if !req.Secret && OnGuestbookEntryCreated != nil {
		go OnGuestbookEntryCreated(context.Background(), entry.ID)
	}
}

func VerifyGuestbookPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	var req models.DeleteGuestbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	entry, err := models.GetGuestbookEntryByID(id)
	if err == sql.ErrNoRows {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch entry", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(entry.PasswordHash), []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateGuestbookEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	var req models.UpdateGuestbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" || len(req.Message) > 500 {
		errorJSON(w, http.StatusBadRequest, "message_length", "Message must be 1-500 characters")
		return
	}

	entry, err := models.GetGuestbookEntryByID(id)
	if err == sql.ErrNoRows {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch entry", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(entry.PasswordHash), []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}

	if err := models.UpdateGuestbookEntry(id, req.Message); err != nil {
		http.Error(w, "Failed to update entry", http.StatusInternalServerError)
		return
	}

	entry.Message = req.Message
	if entry.Secret {
		entry.Nickname = "비밀"
		entry.Message = "비밀 메시지입니다 💌"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)

	if !entry.Secret {
		models.ResetGuestbookEvaluated(id)
		if OnGuestbookEntryCreated != nil {
			go OnGuestbookEntryCreated(context.Background(), id)
		}
	}
}

func DeleteGuestbookEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid entry ID")
		return
	}

	var req models.DeleteGuestbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	entry, err := models.GetGuestbookEntryByID(id)
	if err == sql.ErrNoRows {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch entry", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(entry.PasswordHash), []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}

	if err := models.DeleteGuestbookEntry(id); err != nil {
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
