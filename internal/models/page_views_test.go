package models

import (
	"testing"
	"time"

	"github.com/RunFridge/wedding-card/internal/database"
)

func TestRecordPageView(t *testing.T) {
	cleanTables(t)

	RecordPageView()

	today := time.Now().Format("2006-01-02")
	var count int
	err := database.DB.QueryRow(`SELECT count FROM page_views WHERE date = ?`, today).Scan(&count)
	if err != nil {
		t.Fatalf("query page_views failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

func TestRecordPageViewMultiple(t *testing.T) {
	cleanTables(t)

	RecordPageView()
	RecordPageView()
	RecordPageView()

	today := time.Now().Format("2006-01-02")
	var count int
	err := database.DB.QueryRow(`SELECT count FROM page_views WHERE date = ?`, today).Scan(&count)
	if err != nil {
		t.Fatalf("query page_views failed: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestGetDailyPageViews(t *testing.T) {
	cleanTables(t)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	database.DB.Exec(`INSERT INTO page_views (date, count) VALUES (?, ?)`, yesterday, 5)
	database.DB.Exec(`INSERT INTO page_views (date, count) VALUES (?, ?)`, today, 10)

	views, err := GetPageViews(7)
	if err != nil {
		t.Fatalf("GetPageViews failed: %v", err)
	}
	if len(views) != 2 {
		t.Fatalf("expected 2 views, got %d", len(views))
	}
	if views[0].Date != yesterday {
		t.Errorf("first date = %q, want %q", views[0].Date, yesterday)
	}
	if views[0].Count != 5 {
		t.Errorf("first count = %d, want 5", views[0].Count)
	}
	if views[1].Date != today {
		t.Errorf("second date = %q, want %q", views[1].Date, today)
	}
	if views[1].Count != 10 {
		t.Errorf("second count = %d, want 10", views[1].Count)
	}
}

func TestGetDailyPageViewsEmpty(t *testing.T) {
	cleanTables(t)

	views, err := GetPageViews(30)
	if err != nil {
		t.Fatalf("GetPageViews failed: %v", err)
	}
	if views != nil {
		t.Errorf("expected nil slice for empty result, got %v", views)
	}
}
