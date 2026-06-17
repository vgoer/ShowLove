// Package storage provides a unified storage interface supporting MinIO and RustFS.
package storage

import (
	"context"
	"io"
	"time"
)

// Provider identifies the storage backend.
type Provider string

const (
	ProviderMinIO  Provider = "minio"
	ProviderRustFS Provider = "rustfs"
)

// UploadOptions contains metadata for file uploads.
type UploadOptions struct {
	ContentType string
	Size        int64
	Public      bool
}

// Storage is the unified storage interface.
// It supports both MinIO (S3-compatible) and RustFS (SeaweedFS) backends.
type Storage interface {
	// Upload uploads a file and returns its access URL.
	Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (string, error)

	// Delete removes a file from storage.
	Delete(ctx context.Context, key string) error

	// GetURL returns a signed/presigned URL for accessing the file.
	// If ttl is 0, returns the permanent URL.
	GetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// NoOpStorage is a no-op storage implementation used for testing.
type NoOpStorage struct{}

// NewNoOpStorage creates a new no-op storage.
func NewNoOpStorage() *NoOpStorage {
	return &NoOpStorage{}
}

// Upload implements Storage. Always returns a placeholder URL.
func (s *NoOpStorage) Upload(_ context.Context, key string, _ io.Reader, _ UploadOptions) (string, error) {
	return "http://localhost:9000/showlove/" + key, nil
}

// Delete implements Storage. Always succeeds.
func (s *NoOpStorage) Delete(_ context.Context, _ string) error {
	return nil
}

// GetURL implements Storage. Returns a placeholder URL.
func (s *NoOpStorage) GetURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "http://localhost:9000/showlove/" + key, nil
}
