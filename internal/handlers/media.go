package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/RunFridge/wedding-card/internal/storage"
)

// PhotoURL returns the proxy URL for a user-uploaded photo.
func PhotoURL(hashname string) string {
	if hashname == "" {
		return ""
	}
	return "/api/photos/photo/" + hashname
}

// AssetURL returns the proxy URL for a curated asset photo.
func AssetURL(hashname string) string {
	if hashname == "" {
		return ""
	}
	return "/api/photos/asset/" + hashname
}

// GetMedia streams a photo from the configured Store.
// Route: GET /api/photos/{kind}/{hashname}
// kind: "photo" → photos/ prefix, "asset" → assets/ prefix.
func GetMedia(w http.ResponseWriter, r *http.Request) {
	if Store == nil {
		http.NotFound(w, r)
		return
	}

	hashname := chi.URLParam(r, "hashname")
	if hashname == "" || strings.Contains(hashname, "..") {
		http.NotFound(w, r)
		return
	}

	var prefix string
	switch chi.URLParam(r, "kind") {
	case "photo":
		prefix = "photos/"
	case "asset":
		prefix = "assets/"
	default:
		http.NotFound(w, r)
		return
	}

	etag := `"` + hashname + `"`
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	reader, info, err := Store.GetReader(r.Context(), prefix+hashname)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if info.ContentType != "" {
		w.Header().Set("Content-Type", info.ContentType)
	}
	if info.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("ETag", etag)
	io.Copy(w, reader)
}
