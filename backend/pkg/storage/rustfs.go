package storage

import (
	"context"
	"fmt"
	"io"
	"time"
)

// RustFSConfig holds RustFS/SeaweedFS connection configuration.
type RustFSConfig struct {
	MasterURL string
	FilerURL  string
}

// RustFSStorage implements Storage using RustFS (SeaweedFS).
// This is a placeholder implementation that returns "not implemented" errors.
// Full implementation will be done when RustFS is available in the environment.
type RustFSStorage struct {
	config RustFSConfig
}

// NewRustFSStorage creates a new RustFS-backed storage (placeholder).
func NewRustFSStorage(cfg RustFSConfig) (*RustFSStorage, error) {
	return &RustFSStorage{config: cfg}, nil
}

// Upload implements Storage. Placeholder.
func (s *RustFSStorage) Upload(_ context.Context, _ string, _ io.Reader, _ UploadOptions) (string, error) {
	return "", fmt.Errorf("rustfs: upload not yet implemented")
}

// Delete implements Storage. Placeholder.
func (s *RustFSStorage) Delete(_ context.Context, _ string) error {
	return fmt.Errorf("rustfs: delete not yet implemented")
}

// GetURL implements Storage. Placeholder.
func (s *RustFSStorage) GetURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return fmt.Sprintf("http://%s/%s", s.config.FilerURL, key), nil
}

// NewStorage creates the appropriate storage implementation based on provider.
func NewStorage(provider Provider, minioCfg *MinIOConfig, rustfsCfg *RustFSConfig) (Storage, error) {
	switch provider {
	case ProviderMinIO:
		if minioCfg == nil {
			return nil, fmt.Errorf("storage: MinIO config is required")
		}
		return NewMinIOStorage(*minioCfg)
	case ProviderRustFS:
		if rustfsCfg == nil {
			rustfsCfg = &RustFSConfig{}
		}
		return NewRustFSStorage(*rustfsCfg)
	default:
		return nil, fmt.Errorf("storage: unknown provider %q", provider)
	}
}
