package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOpStorage_Upload(t *testing.T) {
	s := NewNoOpStorage()
	url, err := s.Upload(context.Background(), "images/test.jpg", strings.NewReader("data"), UploadOptions{
		ContentType: "image/jpeg",
		Size:        4,
		Public:      true,
	})

	assert.NoError(t, err)
	assert.Contains(t, url, "images/test.jpg")
}

func TestNoOpStorage_Delete(t *testing.T) {
	s := NewNoOpStorage()
	err := s.Delete(context.Background(), "images/test.jpg")
	assert.NoError(t, err)
}

func TestNoOpStorage_GetURL(t *testing.T) {
	s := NewNoOpStorage()
	url, err := s.GetURL(context.Background(), "images/test.jpg", 0)
	assert.NoError(t, err)
	assert.Contains(t, url, "images/test.jpg")
}

func TestProviderConstants(t *testing.T) {
	assert.Equal(t, Provider("minio"), ProviderMinIO)
	assert.Equal(t, Provider("rustfs"), ProviderRustFS)
}
