package server

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var botUserAgentPattern = regexp.MustCompile(`(?i)bot|crawler|spider|scraper|uptime-kuma|pingdom|uptimerobot|statuscake|monitor|curl|wget|httpclient|go-http-client|python-requests|headlesschrome|phantomjs|slurp|facebookexternalhit|twitterbot|discordbot|slackbot|telegrambot|whatsapp|embedly|quora link preview|outbrain|vkshare|w3c_validator|nuzzel|bitlybot|redditbot|applebot|yeti|naver|daum|bingbot|googlebot|yandex|baidu|duckduck`)

var (
	visitedPerDay sync.Map
	lastPruneDay  atomic.Value
)

// shouldCountPageView reports whether the request is a genuine unique page
// view. It filters out bots/monitoring agents and repeat visits from the same
// IP on the same day.
func shouldCountPageView(r *http.Request) bool {
	ua := strings.TrimSpace(r.UserAgent())
	if ua == "" || botUserAgentPattern.MatchString(ua) {
		return false
	}
	ip := clientIP(r)
	if ip == "" {
		return true
	}
	today := time.Now().Format("2006-01-02")
	pruneVisitedIPs(today)
	if _, loaded := visitedPerDay.LoadOrStore(ip+"|"+today, struct{}{}); loaded {
		return false
	}
	return true
}

func clientIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// pruneVisitedIPs drops dedup entries from prior days. Runs at most once per
// day per process — the first call after a date change sweeps stale keys.
func pruneVisitedIPs(today string) {
	if prev, _ := lastPruneDay.Load().(string); prev == today {
		return
	}
	lastPruneDay.Store(today)
	suffix := "|" + today
	visitedPerDay.Range(func(k, _ any) bool {
		if key, ok := k.(string); ok && !strings.HasSuffix(key, suffix) {
			visitedPerDay.Delete(k)
		}
		return true
	})
}
