package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/imaging"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/storage"
)

var Store storage.Storage

var OnPhotoUploaded func(ctx context.Context, photoID int64)

func PhotoStorageStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if Store == nil {
		json.NewEncoder(w).Encode(map[string]bool{"available": false})
		return
	}
	err := Store.Healthy(r.Context())
	json.NewEncoder(w).Encode(map[string]bool{"available": err == nil})
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func GetVisiblePhotos(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	uploads, err := models.GetVisiblePhotoUploadsPaginated(sort, offset, limit+1)
	if err != nil {
		http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
		return
	}

	hasMore := len(uploads) > limit
	if hasMore {
		uploads = uploads[:limit]
	}

	for i := range uploads {
		uploads[i].URL = PhotoURL(uploads[i].Hashname)
	}

	var nextOffset *int
	if hasMore {
		v := offset + limit
		nextOffset = &v
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"items":       uploads,
		"next_offset": nextOffset,
	})
}

func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if Store == nil {
		http.Error(w, "Photo upload is not configured", http.StatusServiceUnavailable)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorJSON(w, http.StatusBadRequest, "file_too_large", "File too large (max 10MB)")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" || len(name) > 50 {
		errorJSON(w, http.StatusBadRequest, "name_required", "Name is required (max 50 characters)")
		return
	}

	pw := strings.TrimSpace(r.FormValue("password"))
	if len(pw) < 5 || len(pw) > 30 {
		errorJSON(w, http.StatusBadRequest, "password_length", "Password must be 5-30 characters")
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pw), config.Cfg.BcryptCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "image_required", "Image file is required")
		return
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		errorJSON(w, http.StatusBadRequest, "image_files_only", "Only image files are allowed")
		return
	}

	processed, err := imaging.ProcessUpload(file)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "unsupported_image_format", "Unsupported image format")
		return
	}

	hashname := generateUUID() + ".jpg"
	key := "photos/" + hashname

	if err := Store.Upload(r.Context(), key, bytes.NewReader(processed.Optimized), "image/jpeg"); err != nil {
		log.Printf("Store upload failed for key %s: %v", key, err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	var originalHashname string
	originalFile, originalHeader, origErr := r.FormFile("original")
	if origErr == nil {
		defer originalFile.Close()
		origCT := originalHeader.Header.Get("Content-Type")
		if strings.HasPrefix(origCT, "image/") {
			origProcessed, procErr := imaging.ProcessUpload(originalFile)
			if procErr == nil {
				originalHashname = strings.TrimSuffix(hashname, ".jpg") + "_original.jpg"
				origKey := "photos/" + originalHashname
				if uploadErr := Store.Upload(r.Context(), origKey, bytes.NewReader(origProcessed.Optimized), "image/jpeg"); uploadErr != nil {
					log.Printf("Store upload failed for original key %s: %v", origKey, uploadErr)
					originalHashname = ""
				}
			}
		}
	}

	ip := config.IPHash(r.RemoteAddr)
	id, err := models.CreatePhotoUpload(name, ip, hashname, originalHashname, header.Filename, processed.Thumbnail, string(passwordHash))
	if err != nil {
		log.Printf("DB insert failed for photo %s: %v", hashname, err)
		Store.Delete(r.Context(), key)
		http.Error(w, "Failed to save photo record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": id})

	if OnPhotoUploaded != nil {
		go OnPhotoUploaded(context.Background(), id)
	}
}

func VerifyPhotoPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid photo ID")
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	hash, err := models.GetPhotoPasswordHash(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
			return
		}
		http.Error(w, "Failed to fetch photo", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UserDeletePhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid photo ID")
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	hash, err := models.GetPhotoPasswordHash(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
			return
		}
		http.Error(w, "Failed to fetch photo", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		errorJSON(w, http.StatusUnauthorized, "wrong_password", "Wrong password")
		return
	}

	photo, err := models.GetPhotoUploadByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch photo", http.StatusInternalServerError)
		return
	}

	if Store != nil {
		Store.Delete(r.Context(), "photos/"+photo.Hashname)
		if photo.OriginalHashname != "" {
			Store.Delete(r.Context(), "photos/"+photo.OriginalHashname)
		}
	}

	if err := models.DeletePhotoUpload(id); err != nil {
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AdminGetPhotos(w http.ResponseWriter, r *http.Request) {
	uploads, err := models.GetAllPhotoUploads()
	if err != nil {
		http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
		return
	}

	for i := range uploads {
		uploads[i].URL = PhotoURL(uploads[i].Hashname)
		uploads[i].OriginalURL = PhotoURL(uploads[i].OriginalHashname)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploads)
}

func AdminTogglePhotoVisibility(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid photo ID")
		return
	}

	var req struct {
		Hidden bool `json:"hidden"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if err := models.SetPhotoUploadHidden(id, req.Hidden); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Photo not found")
			return
		}
		http.Error(w, "Failed to update photo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminDeletePhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid photo ID")
		return
	}

	photo, err := models.GetPhotoUploadByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Photo not found")
			return
		}
		http.Error(w, "Failed to fetch photo", http.StatusInternalServerError)
		return
	}

	if Store != nil {
		Store.Delete(r.Context(), "photos/"+photo.Hashname)
		if photo.OriginalHashname != "" {
			Store.Delete(r.Context(), "photos/"+photo.OriginalHashname)
		}
	}

	if err := models.DeletePhotoUpload(id); err != nil {
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
