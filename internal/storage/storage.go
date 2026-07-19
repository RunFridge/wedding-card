package storage

import (
	"context"
	"errors"
	"io"
)

var ErrNotFound = errors.New("storage: object not found")

type MediaInfo struct {
	ContentType string
	Size        int64
}

type Storage interface {
	Healthy(ctx context.Context) error
	Upload(ctx context.Context, key string, data io.Reader, contentType string) error
	GetReader(ctx context.Context, key string) (io.ReadCloser, MediaInfo, error)
	Delete(ctx context.Context, key string) error
}
