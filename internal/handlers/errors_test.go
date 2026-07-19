package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	errorJSON(rec, http.StatusBadRequest, "nickname_required", "Nickname is required")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["code"] != "nickname_required" {
		t.Errorf("code = %q, want %q", body["code"], "nickname_required")
	}
	if body["error"] != "Nickname is required" {
		t.Errorf("error = %q, want %q", body["error"], "Nickname is required")
	}
}
