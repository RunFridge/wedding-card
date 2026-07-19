package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalClient struct {
	basePath string
}

func NewLocalClient(basePath string) (*LocalClient, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	return &LocalClient{basePath: basePath}, nil
}

func (c *LocalClient) Healthy(_ context.Context) error {
	_, err := os.Stat(c.basePath)
	return err
}

func (c *LocalClient) Upload(_ context.Context, key string, data io.Reader, _ string) error {
	fullPath := filepath.Join(c.basePath, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (c *LocalClient) GetReader(_ context.Context, key string) (io.ReadCloser, MediaInfo, error) {
	fullPath := filepath.Join(c.basePath, key)
	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, MediaInfo{}, ErrNotFound
		}
		return nil, MediaInfo{}, err
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, MediaInfo{}, err
	}
	return f, MediaInfo{
		ContentType: guessContentType(key),
		Size:        stat.Size(),
	}, nil
}

func (c *LocalClient) Delete(_ context.Context, key string) error {
	err := os.Remove(filepath.Join(c.basePath, key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func guessContentType(key string) string {
	switch strings.ToLower(filepath.Ext(key)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	}
	return "application/octet-stream"
}
