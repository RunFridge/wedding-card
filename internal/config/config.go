package config

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RunFridge/wedding-card/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Port                     string
	DatabasePath             string
	BcryptCost               int
	GameTimerMs              int
	AdminPasswordHash        []byte
	AdminPasswordNeedsChange bool
	IPHashKey                []byte

	S3Endpoint  string
	S3Region    string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string

	UseModeration bool
	OpenAIAPIKey  string

	RateLimitEnabled bool

	SetupRequired bool

	Wedding WeddingConfig
}

type BusInfoEntry struct {
	Stop   string `json:"stop"`
	Routes string `json:"routes"`
}

type CharterBusEntry struct {
	Location  string `json:"location"`
	Company   string `json:"company"`
	BusNumber string `json:"bus_number"`
	Departure string `json:"departure"`
}

type MapLinks struct {
	Google string `json:"google"`
	Kakao  string `json:"kakao"`
	Naver  string `json:"naver"`
	Tmap   string `json:"tmap"`
}

type MapProviders struct {
	EmbedProvider string   `json:"embed_provider"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	APIKey        string   `json:"api_key"`
	Links         MapLinks `json:"links"`
}

type WeddingConfig struct {
	GroomEngName           string             `json:"groom_eng_name"`
	GroomKorName           string             `json:"groom_kor_name"`
	BrideEngName           string             `json:"bride_eng_name"`
	BrideKorName           string             `json:"bride_kor_name"`
	GroomFatherKorName     string             `json:"groom_father_kor_name"`
	GroomMotherKorName     string             `json:"groom_mother_kor_name"`
	BrideFatherKorName     string             `json:"bride_father_kor_name"`
	BrideMotherKorName     string             `json:"bride_mother_kor_name"`
	GroomBankAccount       string             `json:"groom_bank_account"`
	BrideBankAccount       string             `json:"bride_bank_account"`
	GroomFatherBankAccount string             `json:"groom_father_bank_account"`
	GroomMotherBankAccount string             `json:"groom_mother_bank_account"`
	BrideFatherBankAccount string             `json:"bride_father_bank_account"`
	BrideMotherBankAccount string             `json:"bride_mother_bank_account"`
	WeddingDatetime        string             `json:"wedding_datetime"`
	VenueName              string             `json:"venue_name"`
	VenueAddress           string             `json:"venue_address"`
	VenueFloor             string             `json:"venue_floor"`
	VenueHall              string             `json:"venue_hall"`
	VenuePhone             string             `json:"venue_phone"`
	MapProviders           MapProviders       `json:"map_providers"`
	SubwayInfo             []string           `json:"subway_info"`
	BusInfo                []BusInfoEntry     `json:"bus_info"`
	CarInfo                string             `json:"car_info"`
	CharterBus             []CharterBusEntry  `json:"charter_bus"`
	CharterBusNotice       string             `json:"charter_bus_notice"`
	GroomBirthOrder        string             `json:"groom_birth_order"`
	BrideBirthOrder        string             `json:"bride_birth_order"`
	CardGameTimer          int                `json:"card_game_timer"`
	GameNpcMessage         string             `json:"game_npc_message"`
	AvatarColors           string             `json:"avatar_colors"`
	ShortGreeting          string             `json:"short_greeting"`
	MainGreetText          string             `json:"main_greet_text"`
	PhotoUploadEnabled     bool               `json:"photo_upload_enabled"`
	PhotoUploadHoursBefore float64            `json:"photo_upload_hours_before"`
	HeartsFlushIntervalMs  int                `json:"hearts_flush_interval_ms"`
	HeartsFlushBatchSize   int                `json:"hearts_flush_batch_size"`
	ModerationThresholds   map[string]float64 `json:"moderation_thresholds"`
}

var Cfg *Config

func Load() *Config {
	cfg := &Config{
		Port:             "8080",
		DatabasePath:     "./wedding.db",
		BcryptCost:       10,
		GameTimerMs:      30000,
		RateLimitEnabled: true,
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.DatabasePath = dbPath
	}

	cfg.Wedding = loadWeddingDefaults()

	Cfg = cfg
	return cfg
}

func defaultModerationThresholds() map[string]float64 {
	return map[string]float64{
		"sexual/minors":          0.10,
		"violence/graphic":       0.30,
		"self-harm/instructions": 0.30,
		"self-harm/intent":       0.30,
		"hate/threatening":       0.30,
		"illicit/violent":        0.30,
		"harassment/threatening": 0.40,
		"self-harm":              0.40,
		"sexual":                 0.50,
		"violence":               0.60,
		"hate":                   0.80,
		"harassment":             0.80,
		"illicit":                0.90,
	}
}

func loadWeddingDefaults() WeddingConfig {
	return WeddingConfig{
		GroomEngName:           "Groom",
		GroomKorName:           "김철수",
		BrideEngName:           "Bride",
		BrideKorName:           "이영희",
		GroomFatherKorName:     "김아버지",
		GroomMotherKorName:     "박어머니",
		BrideFatherKorName:     "이아버지",
		BrideMotherKorName:     "최어머니",
		GroomBankAccount:       "카카오뱅크 0000-00-0000000",
		BrideBankAccount:       "은행 000-000000-00000",
		GroomFatherBankAccount: "은행 000-000000-00000",
		GroomMotherBankAccount: "은행 000-000000-00000",
		BrideFatherBankAccount: "은행 000-000000-00000",
		BrideMotherBankAccount: "은행 000-000000-00000",
		WeddingDatetime:        "2030-01-01T11:00:00+09:00",
		VenueName:              "OO웨딩홀",
		VenueAddress:           "서울특별시 중구 세종대로 110",
		VenueFloor:             "3층",
		VenueHall:              "그랜드홀",
		VenuePhone:             "02-000-0000",
		MapProviders: MapProviders{
			Links: MapLinks{
				Google: "",
				Kakao:  "",
				Naver:  "",
				Tmap:   "",
			},
		},
		SubwayInfo: []string{"1호선", "시청역 1번 출구 도보 5분 거리"},
		BusInfo: []BusInfoEntry{
			{Stop: "시청역 1번 출구 정류장", Routes: "100, 200, 300"},
		},
		CharterBus: []CharterBusEntry{
			{Location: "OO경기장 주차장", Company: "OO관광", BusNumber: "00가0000", Departure: "오전 7시 30분 출발"},
		},
		CarInfo:                "웨딩홀 전용 주차장 이용 가능",
		GroomBirthOrder:        "장남",
		BrideBirthOrder:        "장녀",
		CardGameTimer:          30000,
		GameNpcMessage:         "안녕! 게임을 클리어하면 신랑 신부의 특별한 사진들을 볼 수 있어요... 도전해 보세요!",
		AvatarColors:           "8B6914,A0722A,C4943A,5C3A0E,F5E6C8",
		ShortGreeting:          "저희 결혼합니다",
		MainGreetText:          "소중한 분들을 초대합니다.\n함께 축복해 주시면 더없는 기쁨으로 간직하겠습니다.",
		PhotoUploadHoursBefore: 1,
		HeartsFlushIntervalMs:  2000,
		HeartsFlushBatchSize:   50,
		ModerationThresholds:   defaultModerationThresholds(),
	}
}

func ApplyOverrides(base WeddingConfig, overrides map[string]string) WeddingConfig {
	c := base

	str := func(key string, dst *string) {
		if v, ok := overrides[key]; ok {
			*dst = v
		}
	}

	str("groom_eng_name", &c.GroomEngName)
	str("groom_kor_name", &c.GroomKorName)
	str("bride_eng_name", &c.BrideEngName)
	str("bride_kor_name", &c.BrideKorName)
	str("groom_father_kor_name", &c.GroomFatherKorName)
	str("groom_mother_kor_name", &c.GroomMotherKorName)
	str("bride_father_kor_name", &c.BrideFatherKorName)
	str("bride_mother_kor_name", &c.BrideMotherKorName)
	str("groom_bank_account", &c.GroomBankAccount)
	str("bride_bank_account", &c.BrideBankAccount)
	str("groom_father_bank_account", &c.GroomFatherBankAccount)
	str("groom_mother_bank_account", &c.GroomMotherBankAccount)
	str("bride_father_bank_account", &c.BrideFatherBankAccount)
	str("bride_mother_bank_account", &c.BrideMotherBankAccount)
	str("wedding_datetime", &c.WeddingDatetime)
	str("venue_name", &c.VenueName)
	str("venue_address", &c.VenueAddress)
	str("venue_floor", &c.VenueFloor)
	str("venue_hall", &c.VenueHall)
	str("venue_phone", &c.VenuePhone)
	if v, ok := overrides["map_providers"]; ok {
		var parsed MapProviders
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			c.MapProviders = parsed
		}
	}

	if _, hasNew := overrides["map_providers"]; !hasNew {
		if v, ok := overrides["kakao_map_link"]; ok {
			c.MapProviders.Links.Kakao = v
		}
		if v, ok := overrides["naver_map_link"]; ok {
			c.MapProviders.Links.Naver = v
		}
		if v, ok := overrides["tmap_link"]; ok {
			c.MapProviders.Links.Tmap = v
		}
	}
	str("car_info", &c.CarInfo)
	str("groom_birth_order", &c.GroomBirthOrder)
	str("bride_birth_order", &c.BrideBirthOrder)
	str("avatar_colors", &c.AvatarColors)
	str("short_greeting", &c.ShortGreeting)
	str("main_greet_text", &c.MainGreetText)

	if v, ok := overrides["photo_upload_enabled"]; ok {
		c.PhotoUploadEnabled = strings.EqualFold(v, "true")
	}
	if v, ok := overrides["photo_upload_hours_before"]; ok {
		if h, err := strconv.ParseFloat(v, 64); err == nil && h >= 0 {
			c.PhotoUploadHoursBefore = h
		}
	}

	if v, ok := overrides["card_game_timer"]; ok {
		if t, err := strconv.Atoi(v); err == nil && t >= 1000 && t <= 120000 {
			c.CardGameTimer = t
		}
	}
	if v, ok := overrides["game_npc_message"]; ok && v != "" {
		c.GameNpcMessage = v
	}

	if v, ok := overrides["hearts_flush_interval_ms"]; ok {
		if t, err := strconv.Atoi(v); err == nil && t >= 500 && t <= 30000 {
			c.HeartsFlushIntervalMs = t
		}
	}

	if v, ok := overrides["hearts_flush_batch_size"]; ok {
		if t, err := strconv.Atoi(v); err == nil && t >= 1 && t <= 1000 {
			c.HeartsFlushBatchSize = t
		}
	}

	if v, ok := overrides["subway_info"]; ok {
		var parsed []string
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			c.SubwayInfo = parsed
		}
	}

	if v, ok := overrides["bus_info"]; ok {
		var parsed []BusInfoEntry
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			c.BusInfo = parsed
		}
	}
	if v, ok := overrides["charter_bus"]; ok {
		var parsed []CharterBusEntry
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			c.CharterBus = parsed
		}
	}
	str("charter_bus_notice", &c.CharterBusNotice)

	if v, ok := overrides["moderation_thresholds"]; ok {
		var parsed map[string]float64
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			if c.ModerationThresholds == nil {
				c.ModerationThresholds = defaultModerationThresholds()
			}
			maps.Copy(c.ModerationThresholds, parsed)
		}
	}

	return c
}

func DiffToOverrideMap(base, current WeddingConfig) map[string]string {
	diff := make(map[string]string)

	strDiff := func(key, baseVal, curVal string) {
		if baseVal != curVal {
			diff[key] = curVal
		}
	}

	strDiff("groom_eng_name", base.GroomEngName, current.GroomEngName)
	strDiff("groom_kor_name", base.GroomKorName, current.GroomKorName)
	strDiff("bride_eng_name", base.BrideEngName, current.BrideEngName)
	strDiff("bride_kor_name", base.BrideKorName, current.BrideKorName)
	strDiff("groom_father_kor_name", base.GroomFatherKorName, current.GroomFatherKorName)
	strDiff("groom_mother_kor_name", base.GroomMotherKorName, current.GroomMotherKorName)
	strDiff("bride_father_kor_name", base.BrideFatherKorName, current.BrideFatherKorName)
	strDiff("bride_mother_kor_name", base.BrideMotherKorName, current.BrideMotherKorName)
	strDiff("groom_bank_account", base.GroomBankAccount, current.GroomBankAccount)
	strDiff("bride_bank_account", base.BrideBankAccount, current.BrideBankAccount)
	strDiff("groom_father_bank_account", base.GroomFatherBankAccount, current.GroomFatherBankAccount)
	strDiff("groom_mother_bank_account", base.GroomMotherBankAccount, current.GroomMotherBankAccount)
	strDiff("bride_father_bank_account", base.BrideFatherBankAccount, current.BrideFatherBankAccount)
	strDiff("bride_mother_bank_account", base.BrideMotherBankAccount, current.BrideMotherBankAccount)
	strDiff("wedding_datetime", base.WeddingDatetime, current.WeddingDatetime)
	strDiff("venue_name", base.VenueName, current.VenueName)
	strDiff("venue_address", base.VenueAddress, current.VenueAddress)
	strDiff("venue_floor", base.VenueFloor, current.VenueFloor)
	strDiff("venue_hall", base.VenueHall, current.VenueHall)
	strDiff("venue_phone", base.VenuePhone, current.VenuePhone)
	strDiff("car_info", base.CarInfo, current.CarInfo)
	strDiff("charter_bus_notice", base.CharterBusNotice, current.CharterBusNotice)
	strDiff("groom_birth_order", base.GroomBirthOrder, current.GroomBirthOrder)
	strDiff("bride_birth_order", base.BrideBirthOrder, current.BrideBirthOrder)
	strDiff("avatar_colors", base.AvatarColors, current.AvatarColors)
	strDiff("short_greeting", base.ShortGreeting, current.ShortGreeting)
	strDiff("main_greet_text", base.MainGreetText, current.MainGreetText)

	if base.PhotoUploadHoursBefore != current.PhotoUploadHoursBefore {
		diff["photo_upload_hours_before"] = strconv.FormatFloat(current.PhotoUploadHoursBefore, 'f', -1, 64)
	}
	if base.PhotoUploadEnabled != current.PhotoUploadEnabled {
		if current.PhotoUploadEnabled {
			diff["photo_upload_enabled"] = "true"
		} else {
			diff["photo_upload_enabled"] = "false"
		}
	}

	if base.CardGameTimer != current.CardGameTimer {
		diff["card_game_timer"] = strconv.Itoa(current.CardGameTimer)
	}
	if base.GameNpcMessage != current.GameNpcMessage {
		diff["game_npc_message"] = current.GameNpcMessage
	}

	if base.HeartsFlushIntervalMs != current.HeartsFlushIntervalMs {
		diff["hearts_flush_interval_ms"] = strconv.Itoa(current.HeartsFlushIntervalMs)
	}

	if base.HeartsFlushBatchSize != current.HeartsFlushBatchSize {
		diff["hearts_flush_batch_size"] = strconv.Itoa(current.HeartsFlushBatchSize)
	}

	jsonDiff := func(key string, baseVal, curVal any) {
		baseJSON, _ := json.Marshal(baseVal)
		curJSON, _ := json.Marshal(curVal)
		if string(baseJSON) != string(curJSON) {
			diff[key] = string(curJSON)
		}
	}

	jsonDiff("subway_info", base.SubwayInfo, current.SubwayInfo)
	jsonDiff("bus_info", base.BusInfo, current.BusInfo)
	jsonDiff("charter_bus", base.CharterBus, current.CharterBus)
	jsonDiff("map_providers", base.MapProviders, current.MapProviders)
	jsonDiff("moderation_thresholds", base.ModerationThresholds, current.ModerationThresholds)

	return diff
}

func InitIPHashKey() error {
	key, err := models.GetSingleConfigOverride("sys:ip_hash_key")
	if err == nil {
		decoded, err := base64.RawURLEncoding.DecodeString(key)
		if err == nil && len(decoded) == 32 {
			Cfg.IPHashKey = decoded
			return nil
		}
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return fmt.Errorf("failed to generate IP hash key: %w", err)
	}
	if err := models.SetSingleConfigOverride("sys:ip_hash_key", base64.RawURLEncoding.EncodeToString(b)); err != nil {
		return fmt.Errorf("failed to persist IP hash key: %w", err)
	}
	Cfg.IPHashKey = b
	return nil
}

func IPHash(remoteAddr string) string {
	ip := strings.TrimSpace(remoteAddr)
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	mac := hmac.New(sha256.New, Cfg.IPHashKey)
	mac.Write([]byte(ip))
	return "hmac-sha256:" + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func InitAdminPassword() error {
	overrides, err := models.GetConfigOverrides()
	if err != nil {
		return fmt.Errorf("failed to read config overrides: %w", err)
	}
	if stored, ok := overrides["admin_password_hash"]; ok {
		Cfg.AdminPasswordHash = []byte(stored)
		log.Println("Admin password loaded from database")
		return nil
	}

	password, err := generateSecurePassword(16)
	if err != nil {
		return fmt.Errorf("failed to generate admin password: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), Cfg.BcryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash generated password: %w", err)
	}
	Cfg.AdminPasswordHash = hash
	if err := models.SetSingleConfigOverride("admin_password_hash", string(hash)); err != nil {
		return fmt.Errorf("failed to persist admin password hash: %w", err)
	}

	Cfg.AdminPasswordNeedsChange = true

	log.Println("========================================")
	log.Println("  AUTO-GENERATED ADMIN PASSWORD")
	log.Printf("  %s", password)
	log.Println("========================================")
	log.Println("Change this password via the admin panel.")

	writePasswordFile(password)

	return nil
}

// writePasswordFile writes the auto-generated password to a file next to the
// database so it can be retrieved when running in detached mode (docker compose up -d).
func writePasswordFile(password string) {
	path := passwordFilePath()
	if err := os.WriteFile(path, []byte(password+"\n"), 0600); err != nil {
		log.Printf("Warning: failed to write password file %s: %v", path, err)
		return
	}
	log.Printf("Admin password written to %s", path)
}

// RemovePasswordFile deletes the auto-generated password file if it exists.
func RemovePasswordFile() {
	path := passwordFilePath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to remove password file %s: %v", path, err)
	}
}

func passwordFilePath() string {
	return filepath.Join(filepath.Dir(Cfg.DatabasePath), "admin_password.txt")
}

func generateSecurePassword(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// InitSetupState determines whether the setup wizard should be shown.
// For existing deployments (password already set, no setup_completed key),
// it auto-marks setup as completed. For fresh installs, it sets SetupRequired.
func InitSetupState() error {
	if !Cfg.AdminPasswordNeedsChange {
		val, err := models.GetSingleConfigOverride("sys:setup_completed")
		if err != nil || val != "true" {
			if err := models.SetSingleConfigOverride("sys:setup_completed", "true"); err != nil {
				return fmt.Errorf("failed to auto-complete setup: %w", err)
			}
		}
		return nil
	}

	val, _ := models.GetSingleConfigOverride("sys:setup_completed")
	if val != "true" {
		Cfg.SetupRequired = true
	}
	return nil
}

// CompleteSetup marks the setup wizard as completed.
func CompleteSetup() error {
	if err := models.SetSingleConfigOverride("sys:setup_completed", "true"); err != nil {
		return fmt.Errorf("failed to complete setup: %w", err)
	}
	Cfg.SetupRequired = false
	Cfg.AdminPasswordNeedsChange = false
	return nil
}

var sensitiveFields = map[string]bool{
	"s3_access_key":  true,
	"s3_secret_key":  true,
	"openai_api_key": true,
}

// LoadAndApplySystemSettings reads sys: overrides from DB and applies them to Cfg.
func LoadAndApplySystemSettings() {
	overrides, err := models.GetSystemConfigOverrides()
	if err != nil {
		log.Printf("Failed to load system settings from DB: %v", err)
		return
	}

	for key, val := range overrides {
		applySystemSetting(key, val)
	}
}

func applySystemSetting(key, val string) {
	switch key {
	case "bcrypt_cost":
		if c, err := strconv.Atoi(val); err == nil && c >= 4 && c <= 31 {
			Cfg.BcryptCost = c
		}
	case "game_timer_ms":
		if t, err := strconv.Atoi(val); err == nil && t >= 1000 && t <= 120000 {
			Cfg.GameTimerMs = t
		}
	case "rate_limit_enabled":
		Cfg.RateLimitEnabled = val == "true"
	case "s3_bucket":
		Cfg.S3Bucket = val
	case "s3_region":
		Cfg.S3Region = val
	case "s3_endpoint":
		Cfg.S3Endpoint = val
	case "s3_access_key":
		Cfg.S3AccessKey = val
	case "s3_secret_key":
		Cfg.S3SecretKey = val
	case "use_moderation":
		Cfg.UseModeration = strings.EqualFold(val, "true")
	case "openai_api_key":
		Cfg.OpenAIAPIKey = val
	}
}

// MaskedSystemSettings returns the current system settings with sensitive
// fields masked.
func MaskedSystemSettings() map[string]any {
	return map[string]any{
		"bcrypt_cost":           Cfg.BcryptCost,
		"game_timer_ms":         Cfg.GameTimerMs,
		"rate_limit_enabled":    Cfg.RateLimitEnabled,
		"s3_bucket":             Cfg.S3Bucket,
		"s3_region":             Cfg.S3Region,
		"s3_endpoint":           Cfg.S3Endpoint,
		"s3_access_key":         "",
		"s3_access_key_set":     Cfg.S3AccessKey != "",
		"s3_secret_key":         "",
		"s3_secret_key_set":     Cfg.S3SecretKey != "",
		"use_moderation":        Cfg.UseModeration,
		"openai_api_key":        "",
		"openai_api_key_set":    Cfg.OpenAIAPIKey != "",
	}
}

// SaveSystemSettings persists system settings to the DB and updates Cfg.
// Sensitive fields with empty values are kept as-is (not overwritten).
func SaveSystemSettings(updates map[string]any) error {
	existing, err := models.GetSystemConfigOverrides()
	if err != nil {
		return fmt.Errorf("failed to load existing system settings: %w", err)
	}

	result := make(map[string]string, len(existing))
	maps.Copy(result, existing)

	for key, val := range updates {
		if strings.HasSuffix(key, "_set") {
			continue
		}

		strVal := fmt.Sprintf("%v", val)

		if sensitiveFields[key] && strVal == "" {
			continue
		}

		if key == "use_moderation" || key == "rate_limit_enabled" {
			switch v := val.(type) {
			case bool:
				strVal = strconv.FormatBool(v)
			}
		}

		if key == "bcrypt_cost" || key == "game_timer_ms" {
			switch v := val.(type) {
			case float64:
				strVal = strconv.Itoa(int(v))
			}
		}

		result[key] = strVal
	}

	if err := models.SetSystemConfigOverrides(result); err != nil {
		return fmt.Errorf("failed to save system settings: %w", err)
	}

	for key, val := range result {
		applySystemSetting(key, val)
	}

	return nil
}
