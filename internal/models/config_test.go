package models

import (
	"database/sql"
	"testing"
)

func TestGetConfigOverridesEmpty(t *testing.T) {
	cleanTables(t)

	overrides, err := GetConfigOverrides()
	if err != nil {
		t.Fatalf("GetConfigOverrides failed: %v", err)
	}
	if len(overrides) != 0 {
		t.Errorf("expected 0 overrides, got %d", len(overrides))
	}
}

func TestSetAndGetSingleConfigOverride(t *testing.T) {
	cleanTables(t)

	if err := SetSingleConfigOverride("groom_eng_name", "John"); err != nil {
		t.Fatalf("SetSingleConfigOverride failed: %v", err)
	}

	val, err := GetSingleConfigOverride("groom_eng_name")
	if err != nil {
		t.Fatalf("GetSingleConfigOverride failed: %v", err)
	}
	if val != "John" {
		t.Errorf("value = %q, want %q", val, "John")
	}
}

func TestGetSingleConfigOverrideNotFound(t *testing.T) {
	cleanTables(t)

	_, err := GetSingleConfigOverride("nonexistent")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestSetSingleConfigOverrideUpsert(t *testing.T) {
	cleanTables(t)

	SetSingleConfigOverride("key1", "val1")
	SetSingleConfigOverride("key1", "val2")

	val, _ := GetSingleConfigOverride("key1")
	if val != "val2" {
		t.Errorf("expected upserted value %q, got %q", "val2", val)
	}
}

func TestSetConfigOverrides(t *testing.T) {
	cleanTables(t)

	overrides := map[string]string{
		"groom_eng_name": "John",
		"bride_eng_name": "Jane",
	}
	if err := SetConfigOverrides(overrides); err != nil {
		t.Fatalf("SetConfigOverrides failed: %v", err)
	}

	result, err := GetConfigOverrides()
	if err != nil {
		t.Fatalf("GetConfigOverrides failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 overrides, got %d", len(result))
	}
	if result["groom_eng_name"] != "John" {
		t.Errorf("groom_eng_name = %q, want %q", result["groom_eng_name"], "John")
	}
}

func TestSetConfigOverridesSkipsEmpty(t *testing.T) {
	cleanTables(t)

	overrides := map[string]string{
		"key1": "val1",
		"key2": "",
	}
	SetConfigOverrides(overrides)

	result, _ := GetConfigOverrides()
	if _, ok := result["key2"]; ok {
		t.Error("expected empty value to be skipped")
	}
}

func TestSetConfigOverridesPreservesSysKeys(t *testing.T) {
	cleanTables(t)

	SetSingleConfigOverride("sys:s3_bucket", "my-bucket")
	SetSingleConfigOverride("admin_password_hash", "hash123")

	SetConfigOverrides(map[string]string{"groom_eng_name": "John"})

	// sys: and admin_password_hash should survive
	val, err := GetSingleConfigOverride("sys:s3_bucket")
	if err != nil {
		t.Fatalf("sys: key was deleted: %v", err)
	}
	if val != "my-bucket" {
		t.Errorf("sys:s3_bucket = %q, want %q", val, "my-bucket")
	}

	val, err = GetSingleConfigOverride("admin_password_hash")
	if err != nil {
		t.Fatalf("admin_password_hash was deleted: %v", err)
	}
	if val != "hash123" {
		t.Errorf("admin_password_hash = %q, want %q", val, "hash123")
	}
}

func TestGetSystemConfigOverrides(t *testing.T) {
	cleanTables(t)

	SetSingleConfigOverride("sys:s3_bucket", "bucket1")
	SetSingleConfigOverride("sys:s3_region", "us-east-1")
	SetSingleConfigOverride("groom_eng_name", "John")

	result, err := GetSystemConfigOverrides()
	if err != nil {
		t.Fatalf("GetSystemConfigOverrides failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 system overrides, got %d", len(result))
	}
	if result["s3_bucket"] != "bucket1" {
		t.Errorf("s3_bucket = %q, want %q", result["s3_bucket"], "bucket1")
	}
}

func TestSetSystemConfigOverrides(t *testing.T) {
	cleanTables(t)

	SetSingleConfigOverride("sys:old_key", "old_val")

	overrides := map[string]string{
		"s3_bucket": "new-bucket",
		"s3_region": "eu-west-1",
	}
	if err := SetSystemConfigOverrides(overrides); err != nil {
		t.Fatalf("SetSystemConfigOverrides failed: %v", err)
	}

	result, _ := GetSystemConfigOverrides()
	if result["s3_bucket"] != "new-bucket" {
		t.Errorf("s3_bucket = %q, want %q", result["s3_bucket"], "new-bucket")
	}
	if _, ok := result["old_key"]; ok {
		t.Error("expected old sys: key to be deleted")
	}
}
