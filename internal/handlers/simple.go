package handlers

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/models"
)

//go:embed simple.html
var simpleHTML string

var simpleTmpl = template.Must(template.New("simple").Parse(simpleHTML))

const simpleFontCDN = "https://cdn.jsdelivr.net"

var simpleCSP = strings.Join([]string{
	"default-src 'none'",
	"style-src 'unsafe-inline' " + simpleFontCDN,
	"font-src " + simpleFontCDN,
	"base-uri 'self'",
	"form-action 'none'",
	"frame-ancestors 'none'",
}, "; ")

type simpleAccount struct {
	Label string
	Value string
}

type simplePageData struct {
	config.WeddingConfig
	DateText        string
	GroomFamilyLine string
	BrideFamilyLine string
	Accounts        []simpleAccount
}

func SimplePage(w http.ResponseWriter, r *http.Request) {
	merged := config.Cfg.Wedding
	overrides, err := models.GetConfigOverrides()
	if err != nil {
		log.Printf("Failed to load config overrides: %v", err)
	} else if len(overrides) > 0 {
		merged = config.ApplyOverrides(merged, overrides)
	}

	if raw := merged.SimpleRedirectURL; raw != "" {
		if u, err := url.Parse(raw); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			http.Redirect(w, r, raw, http.StatusFound)
			return
		}
	}

	data := simplePageData{
		WeddingConfig:   merged,
		DateText:        formatKoreanDatetime(merged.WeddingDatetime),
		GroomFamilyLine: familyLine(merged.GroomFatherKorName, merged.GroomMotherKorName, merged.GroomBirthOrder, merged.GroomKorName),
		BrideFamilyLine: familyLine(merged.BrideFatherKorName, merged.BrideMotherKorName, merged.BrideBirthOrder, merged.BrideKorName),
		Accounts:        collectAccounts(merged),
	}

	var buf bytes.Buffer
	if err := simpleTmpl.Execute(&buf, data); err != nil {
		log.Printf("Failed to render simple page: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Security-Policy", simpleCSP)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func formatKoreanDatetime(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	weekdays := [...]string{"일", "월", "화", "수", "목", "금", "토"}
	ampm := "오전"
	hour := t.Hour()
	if hour >= 12 {
		ampm = "오후"
	}
	if hour == 0 {
		hour = 12
	} else if hour > 12 {
		hour -= 12
	}
	s := fmt.Sprintf("%d년 %d월 %d일 %s요일 %s %d시", t.Year(), int(t.Month()), t.Day(), weekdays[t.Weekday()], ampm, hour)
	if t.Minute() > 0 {
		s += fmt.Sprintf(" %d분", t.Minute())
	}
	return s
}

func familyLine(father, mother, birthOrder, name string) string {
	var parents []string
	for _, p := range []string{father, mother} {
		if p != "" {
			parents = append(parents, p)
		}
	}
	if len(parents) == 0 || name == "" {
		return ""
	}
	relation := name
	if birthOrder != "" {
		relation = birthOrder + " " + name
	}
	return strings.Join(parents, " · ") + "의 " + relation
}

func collectAccounts(c config.WeddingConfig) []simpleAccount {
	candidates := []simpleAccount{
		{Label: "신랑 " + c.GroomKorName, Value: c.GroomBankAccount},
		{Label: "신랑 아버지 " + c.GroomFatherKorName, Value: c.GroomFatherBankAccount},
		{Label: "신랑 어머니 " + c.GroomMotherKorName, Value: c.GroomMotherBankAccount},
		{Label: "신부 " + c.BrideKorName, Value: c.BrideBankAccount},
		{Label: "신부 아버지 " + c.BrideFatherKorName, Value: c.BrideFatherBankAccount},
		{Label: "신부 어머니 " + c.BrideMotherKorName, Value: c.BrideMotherBankAccount},
	}
	var accounts []simpleAccount
	for _, a := range candidates {
		if a.Value != "" {
			accounts = append(accounts, simpleAccount{Label: strings.TrimSpace(a.Label), Value: a.Value})
		}
	}
	return accounts
}
