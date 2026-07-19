package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/moderation"
	"github.com/RunFridge/wedding-card/internal/storage"
)

var ReinitStorage func(endpoint, region, bucket, accessKey, secretKey string) error

func AdminGetSystemSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.MaskedSystemSettings())
}

func AdminUpdateSystemSettings(w http.ResponseWriter, r *http.Request) {
	var updates map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	// Validate bcrypt_cost
	if v, ok := updates["bcrypt_cost"]; ok {
		if f, ok := v.(float64); ok {
			c := int(f)
			if c < 4 || c > 31 {
				errorJSON(w, http.StatusBadRequest, "bcrypt_cost_range", "bcrypt_cost must be between 4 and 31")
				return
			}
		}
	}

	// Validate game_timer_ms
	if v, ok := updates["game_timer_ms"]; ok {
		if f, ok := v.(float64); ok {
			t := int(f)
			if t < 1000 || t > 120000 {
				errorJSON(w, http.StatusBadRequest, "game_timer_range", "game_timer_ms must be between 1000 and 120000")
				return
			}
		}
	}

	// Snapshot S3 values before save to detect changes
	prevBucket := config.Cfg.S3Bucket
	prevRegion := config.Cfg.S3Region
	prevEndpoint := config.Cfg.S3Endpoint
	prevAccessKey := config.Cfg.S3AccessKey
	prevSecretKey := config.Cfg.S3SecretKey
	prevModeration := config.Cfg.UseModeration
	prevOpenAIKey := config.Cfg.OpenAIAPIKey

	if err := config.SaveSystemSettings(updates); err != nil {
		http.Error(w, "Failed to save settings", http.StatusInternalServerError)
		return
	}

	s3Changed := config.Cfg.S3Bucket != prevBucket ||
		config.Cfg.S3Region != prevRegion ||
		config.Cfg.S3Endpoint != prevEndpoint ||
		config.Cfg.S3AccessKey != prevAccessKey ||
		config.Cfg.S3SecretKey != prevSecretKey

	restartRequired := config.Cfg.UseModeration != prevModeration ||
		config.Cfg.OpenAIAPIKey != prevOpenAIKey

	s3Reinitialized := false
	if s3Changed && ReinitStorage != nil {
		if err := ReinitStorage(
			config.Cfg.S3Endpoint, config.Cfg.S3Region, config.Cfg.S3Bucket,
			config.Cfg.S3AccessKey, config.Cfg.S3SecretKey,
		); err == nil {
			s3Reinitialized = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":            "ok",
		"restart_required":  restartRequired,
		"s3_reinitialized":  s3Reinitialized,
	})
}

func AdminTestModeration(w http.ResponseWriter, r *http.Request) {
	var req struct {
		APIKey string `json:"openai_api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	apiKey := req.APIKey
	if apiKey == "" {
		apiKey = config.Cfg.OpenAIAPIKey
	}
	if apiKey == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "No API key provided",
		})
		return
	}

	client := moderation.NewClient(apiKey)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	_, err := client.ModerateText(ctx, "hello")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func AdminTestS3Connection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Endpoint  string `json:"s3_endpoint"`
		Region    string `json:"s3_region"`
		Bucket    string `json:"s3_bucket"`
		AccessKey string `json:"s3_access_key"`
		SecretKey string `json:"s3_secret_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	if req.Bucket == "" {
		errorJSON(w, http.StatusBadRequest, "s3_bucket_required", "s3_bucket is required")
		return
	}

	// If access key or secret key are empty, use current config values
	accessKey := req.AccessKey
	if accessKey == "" {
		accessKey = config.Cfg.S3AccessKey
	}
	secretKey := req.SecretKey
	if secretKey == "" {
		secretKey = config.Cfg.S3SecretKey
	}

	client, err := storage.NewS3Client(
		req.Endpoint, req.Region, req.Bucket,
		accessKey, secretKey,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := client.Healthy(ctx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
