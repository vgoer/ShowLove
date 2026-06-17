package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig holds MinIO connection configuration.
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MinIOStorage implements Storage using MinIO (S3-compatible).
type MinIOStorage struct {
	client *minio.Client
	config MinIOConfig
}

// NewMinIOStorage creates a new MinIO-backed storage.
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio storage: create client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio storage: check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio storage: create bucket: %w", err)
		}
	}

	return &MinIOStorage{client: client, config: cfg}, nil
}

// Upload implements Storage. Uploads a file to MinIO.
func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (string, error) {
	_, err := s.client.PutObject(ctx, s.config.Bucket, key, reader, opts.Size, minio.PutObjectOptions{
		ContentType: opts.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio storage: upload %s: %w", key, err)
	}

	return s.GetURL(ctx, key, 0)
}

// Delete implements Storage. Removes a file from MinIO.
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.config.Bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio storage: delete %s: %w", key, err)
	}
	return nil
}

// GetURL implements Storage. Returns a presigned or permanent URL.
func (s *MinIOStorage) GetURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	scheme := "http"
	if s.config.UseSSL {
		scheme = "https"
	}

	if ttl > 0 {
		url, err := s.client.PresignedGetObject(ctx, s.config.Bucket, key, ttl, nil)
		if err != nil {
			return "", fmt.Errorf("minio storage: presigned url %s: %w", key, err)
		}
		return url.String(), nil
	}

	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.config.Endpoint, s.config.Bucket, key), nil
}

// Thumbnail generates a 200x200 thumbnail for the given image.
// key is the original image key; returns the thumbnail key.
func ThumbnailKey(key string) string {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 2 {
		return parts[0] + "_thumb." + parts[1]
	}
	return key + "_thumb"
}

// GenerateThumbnail creates a thumbnail from an image reader.
// Returns the thumbnail bytes and the content type.
func GenerateThumbnail(reader io.Reader, contentType string) ([]byte, string, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("thumbnail: decode image: %w", err)
	}

	// Resize to 200x200 (simple nearest-neighbor resize)
	thumb := resizeImage(img, 200, 200)

	var buf bytes.Buffer
	switch contentType {
	case "image/png":
		if err := png.Encode(&buf, thumb); err != nil {
			return nil, "", fmt.Errorf("thumbnail: encode png: %w", err)
		}
		return buf.Bytes(), "image/png", nil
	default:
		if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80}); err != nil {
			return nil, "", fmt.Errorf("thumbnail: encode jpeg: %w", err)
		}
		return buf.Bytes(), "image/jpeg", nil
	}
}

// resizeImage performs a simple resize using nearest-neighbor interpolation.
func resizeImage(src image.Image, width, height int) *image.RGBA {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcW / width
			srcY := y * srcH / height
			dst.Set(x, y, src.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))
		}
	}

	return dst
}
