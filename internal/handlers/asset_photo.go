package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/RunFridge/wedding-card/internal/imaging"
	"github.com/RunFridge/wedding-card/internal/models"
)

func AdminGetAssetPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := models.GetAllAssetPhotos()
	if err != nil {
		http.Error(w, "Failed to fetch asset photos", http.StatusInternalServerError)
		return
	}

	for i := range photos {
		photos[i].URL = AssetURL(photos[i].Hashname)
		photos[i].ThumbnailURL = AssetURL(photos[i].ThumbHashname)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

func AdminUploadAssetPhoto(w http.ResponseWriter, r *http.Request) {
	if Store == nil {
		http.Error(w, "Photo storage not configured", http.StatusServiceUnavailable)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorJSON(w, http.StatusBadRequest, "file_too_large", "File too large (max 10MB)")
		return
	}

	label := strings.TrimSpace(r.FormValue("label"))

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

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	uuid := generateUUID()
	hashname := uuid + ".jpg"
	thumbHashname := uuid + "_thumb.jpg"

	if err := Store.Upload(r.Context(), "assets/"+hashname, bytes.NewReader(imageData), "image/jpeg"); err != nil {
		log.Printf("Store upload failed for asset %s: %v", hashname, err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// Use client-provided thumbnail if available, otherwise fall back to server-side generation
	var thumbData []byte
	if thumbFile, _, thumbErr := r.FormFile("thumbnail"); thumbErr == nil {
		defer thumbFile.Close()
		thumbData, _ = io.ReadAll(thumbFile)
	}
	if len(thumbData) == 0 {
		processed, procErr := imaging.ProcessAssetUpload(bytes.NewReader(imageData))
		if procErr != nil {
			Store.Delete(r.Context(), "assets/"+hashname)
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}
		thumbData = processed.Thumbnail
	}

	if err := Store.Upload(r.Context(), "assets/"+thumbHashname, bytes.NewReader(thumbData), "image/jpeg"); err != nil {
		log.Printf("Store upload failed for asset thumb %s: %v", thumbHashname, err)
		Store.Delete(r.Context(), "assets/"+hashname)
		http.Error(w, "Failed to upload thumbnail", http.StatusInternalServerError)
		return
	}

	id, err := models.CreateAssetPhoto(label, hashname, thumbHashname, header.Filename)
	if err != nil {
		log.Printf("DB insert failed for asset photo %s: %v", hashname, err)
		Store.Delete(r.Context(), "assets/"+hashname)
		Store.Delete(r.Context(), "assets/"+thumbHashname)
		http.Error(w, "Failed to save asset photo record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": id})
}

func AdminToggleAssetPhotoGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid asset photo ID")
		return
	}

	var req struct {
		UseForGame bool `json:"use_for_game"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if err := models.SetAssetPhotoGameFlag(id, req.UseForGame); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Asset photo not found")
			return
		}
		http.Error(w, "Failed to update asset photo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminSetMainPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid asset photo ID")
		return
	}

	var req struct {
		IsMainPhoto bool `json:"is_main_photo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if req.IsMainPhoto {
		if err := models.SetAssetPhotoAsMain(id); err != nil {
			if err == sql.ErrNoRows {
				errorJSON(w, http.StatusNotFound, "not_found", "Asset photo not found")
				return
			}
			http.Error(w, "Failed to set main photo", http.StatusInternalServerError)
			return
		}
	} else {
		if err := models.ClearMainAssetPhoto(id); err != nil {
			if err == sql.ErrNoRows {
				errorJSON(w, http.StatusNotFound, "not_found", "Asset photo not found")
				return
			}
			http.Error(w, "Failed to clear main photo", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminUpdateAssetPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid asset photo ID")
		return
	}

	var req struct {
		Label     string `json:"label"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if err := models.UpdateAssetPhoto(id, req.Label, req.SortOrder); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Asset photo not found")
			return
		}
		http.Error(w, "Failed to update asset photo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AdminDeleteAssetPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid asset photo ID")
		return
	}

	photo, err := models.DeleteAssetPhoto(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Asset photo not found")
			return
		}
		http.Error(w, "Failed to delete asset photo", http.StatusInternalServerError)
		return
	}

	if Store != nil {
		Store.Delete(r.Context(), "assets/"+photo.Hashname)
		Store.Delete(r.Context(), "assets/"+photo.ThumbHashname)
	}

	w.WriteHeader(http.StatusNoContent)
}
