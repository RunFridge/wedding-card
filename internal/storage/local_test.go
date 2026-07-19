package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLocalClient(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "storage")

	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}

func TestLocalUploadAndGetReader(t *testing.T) {
	dir := t.TempDir()
	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}

	ctx := context.Background()
	content := "hello world"
	key := "test/photo.jpg"

	if err := client.Upload(ctx, key, strings.NewReader(content), "image/jpeg"); err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, key))
	if err != nil {
		t.Fatalf("file not found after upload: %v", err)
	}
	if string(data) != content {
		t.Errorf("file content = %q, want %q", string(data), content)
	}

	reader, info, err := client.GetReader(ctx, key)
	if err != nil {
		t.Fatalf("GetReader failed: %v", err)
	}
	defer reader.Close()

	if info.ContentType != "image/jpeg" {
		t.Errorf("ContentType = %q, want %q", info.ContentType, "image/jpeg")
	}
	if info.Size != int64(len(content)) {
		t.Errorf("Size = %d, want %d", info.Size, len(content))
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(body) != content {
		t.Errorf("stream body = %q, want %q", string(body), content)
	}
}

func TestLocalGetReaderNotFound(t *testing.T) {
	dir := t.TempDir()
	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}

	_, _, err = client.GetReader(context.Background(), "missing.jpg")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestLocalDelete(t *testing.T) {
	dir := t.TempDir()
	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}

	ctx := context.Background()
	key := "deleteme.txt"

	if err := client.Upload(ctx, key, strings.NewReader("data"), "text/plain"); err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	if err := client.Delete(ctx, key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, key)); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestLocalDeleteNotFound(t *testing.T) {
	dir := t.TempDir()
	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}

	err = client.Delete(context.Background(), "nonexistent.txt")
	if err != nil {
		t.Errorf("expected nil error for non-existent file, got %v", err)
	}
}

func TestLocalHealthy(t *testing.T) {
	dir := t.TempDir()
	client, err := NewLocalClient(dir)
	if err != nil {
		t.Fatalf("NewLocalClient failed: %v", err)
	}

	if err := client.Healthy(context.Background()); err != nil {
		t.Errorf("Healthy failed: %v", err)
	}
}
