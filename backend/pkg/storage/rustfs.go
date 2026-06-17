package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// RustFSConfig holds S3-compatible storage connection configuration.
// Despite the name "RustFS", this uses the standard S3 protocol
// and is compatible with MinIO, Ceph, SeaweedFS S3 gateway, etc.
type RustFSConfig struct {
	Endpoint  string // e.g. "192.168.77.100:9000"
	AccessKey string
	SecretKey string
	Region    string // e.g. "us-east-1"
	Bucket    string // e.g. "packages"
	UseSSL    bool
}

// RustFSStorage implements Storage using S3-compatible protocol.
type RustFSStorage struct {
	client *minio.Client
	config RustFSConfig
}

// NewRustFSStorage creates a new S3-compatible storage backed by minio-go.
func NewRustFSStorage(cfg RustFSConfig) (*RustFSStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("rustfs storage: create client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("rustfs storage: check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region}); err != nil {
			return nil, fmt.Errorf("rustfs storage: create bucket: %w", err)
		}
	}

	return &RustFSStorage{client: client, config: cfg}, nil
}

// Upload implements Storage.
func (s *RustFSStorage) Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (string, error) {
	_, err := s.client.PutObject(ctx, s.config.Bucket, key, reader, opts.Size, minio.PutObjectOptions{
		ContentType: opts.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("rustfs storage: upload %s: %w", key, err)
	}
	return s.GetURL(ctx, key, 0)
}

// Delete implements Storage.
func (s *RustFSStorage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.config.Bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("rustfs storage: delete %s: %w", key, err)
	}
	return nil
}

// GetURL implements Storage.
func (s *RustFSStorage) GetURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	scheme := "http"
	if s.config.UseSSL {
		scheme = "https"
	}

	if ttl > 0 {
		url, err := s.client.PresignedGetObject(ctx, s.config.Bucket, key, ttl, nil)
		if err != nil {
			return "", fmt.Errorf("rustfs storage: presigned url %s: %w", key, err)
		}
		return url.String(), nil
	}

	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.config.Endpoint, s.config.Bucket, key), nil
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
			return nil, fmt.Errorf("storage: RustFS config is required")
		}
		return NewRustFSStorage(*rustfsCfg)
	default:
		return nil, fmt.Errorf("storage: unknown provider %q", provider)
	}
}
