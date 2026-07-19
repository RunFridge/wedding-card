package server

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func resetPageViewFilterState() {
	visitedPerDay = sync.Map{}
	lastPruneDay.Store("")
}

func TestShouldCountPageView_FiltersBotUserAgents(t *testing.T) {
	bots := []string{
		"Uptime-Kuma/1.23.11",
		"UptimeRobot/2.0",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"curl/8.5.0",
		"Wget/1.21",
		"Go-http-client/1.1",
		"python-requests/2.31.0",
		"facebookexternalhit/1.1",
		"Twitterbot/1.0",
		"Slackbot-LinkExpanding 1.0",
		"Discordbot/2.0",
		"HeadlessChrome/120.0.0.0",
	}
	for _, ua := range bots {
		resetPageViewFilterState()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("User-Agent", ua)
		req.RemoteAddr = "10.0.0.1:1234"
		if shouldCountPageView(req) {
			t.Errorf("expected bot UA %q to be filtered", ua)
		}
	}
}

func TestShouldCountPageView_FiltersEmptyUserAgent(t *testing.T) {
	resetPageViewFilterState()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Del("User-Agent")
	req.RemoteAddr = "10.0.0.1:1234"
	if shouldCountPageView(req) {
		t.Error("expected empty UA to be filtered")
	}
}

func TestShouldCountPageView_CountsRealBrowser(t *testing.T) {
	resetPageViewFilterState()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1")
	req.RemoteAddr = "10.0.0.1:1234"
	if !shouldCountPageView(req) {
		t.Error("expected real browser UA to count")
	}
}

func TestShouldCountPageView_DedupesSameIPSameDay(t *testing.T) {
	resetPageViewFilterState()
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36"

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", ua)
	req.RemoteAddr = "203.0.113.5:4001"
	if !shouldCountPageView(req) {
		t.Fatal("first visit should count")
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("User-Agent", ua)
	req2.RemoteAddr = "203.0.113.5:4002"
	if shouldCountPageView(req2) {
		t.Error("second visit from same IP should be deduped")
	}
}

func TestShouldCountPageView_DifferentIPsBothCount(t *testing.T) {
	resetPageViewFilterState()
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.0 Safari/605.1.15"

	r1 := httptest.NewRequest("GET", "/", nil)
	r1.Header.Set("User-Agent", ua)
	r1.RemoteAddr = "198.51.100.1:5000"

	r2 := httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("User-Agent", ua)
	r2.RemoteAddr = "198.51.100.2:5000"

	if !shouldCountPageView(r1) {
		t.Error("first IP should count")
	}
	if !shouldCountPageView(r2) {
		t.Error("second IP should count independently")
	}
}

func TestShouldCountPageView_HandlesRealIPWithoutPort(t *testing.T) {
	resetPageViewFilterState()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.RemoteAddr = "203.0.113.9" // chi middleware.RealIP strips the port

	if !shouldCountPageView(req) {
		t.Fatal("first visit should count")
	}
	if shouldCountPageView(req) {
		t.Error("second visit with same IP (no port) should be deduped")
	}
}

func TestPruneVisitedIPs_ClearsPriorDayEntries(t *testing.T) {
	resetPageViewFilterState()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	visitedPerDay.Store("10.0.0.1|"+yesterday, struct{}{})
	visitedPerDay.Store("10.0.0.2|"+today, struct{}{})

	pruneVisitedIPs(today)

	if _, ok := visitedPerDay.Load("10.0.0.1|" + yesterday); ok {
		t.Error("yesterday's entry should have been pruned")
	}
	if _, ok := visitedPerDay.Load("10.0.0.2|" + today); !ok {
		t.Error("today's entry should survive prune")
	}
}
