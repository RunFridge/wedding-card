package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/models"
)

func clearSimpleRedirect(t *testing.T) {
	t.Helper()
	if _, err := database.DB.Exec("DELETE FROM wedding_config_overrides WHERE key = 'simple_redirect_url'"); err != nil {
		t.Fatalf("failed to clear override: %v", err)
	}
}

func TestSimplePage(t *testing.T) {
	t.Run("renders page when redirect unset", func(t *testing.T) {
		clearSimpleRedirect(t)
		w := httptest.NewRecorder()
		SimplePage(w, httptest.NewRequest("GET", "/simple", nil))

		if w.Code != 200 {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
			t.Errorf("expected text/html content type, got %q", ct)
		}
		if !strings.Contains(w.Body.String(), "모바일 청첩장 보기") {
			t.Error("expected rendered page body")
		}
	})

	t.Run("redirects when url set", func(t *testing.T) {
		if err := models.SetSingleConfigOverride("simple_redirect_url", "https://example.com/card"); err != nil {
			t.Fatalf("failed to set override: %v", err)
		}
		defer clearSimpleRedirect(t)

		w := httptest.NewRecorder()
		SimplePage(w, httptest.NewRequest("GET", "/simple", nil))

		if w.Code != 302 {
			t.Fatalf("expected 302, got %d", w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "https://example.com/card" {
			t.Errorf("expected redirect location, got %q", loc)
		}
	})

	t.Run("renders page on invalid redirect url", func(t *testing.T) {
		if err := models.SetSingleConfigOverride("simple_redirect_url", "javascript:alert(1)"); err != nil {
			t.Fatalf("failed to set override: %v", err)
		}
		defer clearSimpleRedirect(t)

		w := httptest.NewRecorder()
		SimplePage(w, httptest.NewRequest("GET", "/simple", nil))

		if w.Code != 200 {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})
}

func TestFormatKoreanDatetime(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2030-01-01T11:00:00+09:00", "2030년 1월 1일 화요일 오전 11시"},
		{"2026-07-19T12:42:00+09:00", "2026년 7월 19일 일요일 오후 12시 42분"},
		{"2026-07-19T00:30:00+09:00", "2026년 7월 19일 일요일 오전 12시 30분"},
		{"not a date", "not a date"},
	}
	for _, c := range cases {
		if got := formatKoreanDatetime(c.in); got != c.want {
			t.Errorf("formatKoreanDatetime(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFamilyLine(t *testing.T) {
	if got := familyLine("김아버지", "김어머니", "장남", "김철수"); got != "김아버지 · 김어머니의 장남 김철수" {
		t.Errorf("unexpected family line: %q", got)
	}
	if got := familyLine("", "김어머니", "", "김철수"); got != "김어머니의 김철수" {
		t.Errorf("unexpected single-parent line: %q", got)
	}
	if got := familyLine("", "", "장남", "김철수"); got != "" {
		t.Errorf("expected empty line without parents, got %q", got)
	}
}
