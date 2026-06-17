package storage

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThumbnailKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"images/photo.jpg", "images/photo_thumb.jpg"},
		{"avatar.png", "avatar_thumb.png"},
		{"noext", "noext_thumb"},
		{"path/to/file.jpeg", "path/to/file_thumb.jpeg"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ThumbnailKey(tt.input))
		})
	}
}

func TestGenerateThumbnail(t *testing.T) {
	// Create a 400x300 test image
	src := image.NewRGBA(image.Rect(0, 0, 400, 300))
	var buf bytes.Buffer
	err := png.Encode(&buf, src)
	require.NoError(t, err)

	thumbBytes, contentType, err := GenerateThumbnail(&buf, "image/png")
	require.NoError(t, err)
	assert.Equal(t, "image/png", contentType)
	assert.NotEmpty(t, thumbBytes)

	// Verify the thumbnail is 200x200
	thumb, _, err := image.Decode(bytes.NewReader(thumbBytes))
	require.NoError(t, err)
	bounds := thumb.Bounds()
	assert.Equal(t, 200, bounds.Dx())
	assert.Equal(t, 200, bounds.Dy())
}

func TestGenerateThumbnail_JPG(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 800, 600))
	var buf bytes.Buffer
	err := png.Encode(&buf, src)
	require.NoError(t, err)

	thumbBytes, contentType, err := GenerateThumbnail(&buf, "image/jpeg")
	require.NoError(t, err)
	assert.Equal(t, "image/jpeg", contentType)
	assert.NotEmpty(t, thumbBytes)
}

func TestNewStorage_UnknownProvider(t *testing.T) {
	_, err := NewStorage("unknown", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown provider")
}

func TestNewStorage_RustFS(t *testing.T) {
	cfg := &RustFSConfig{MasterURL: "localhost:9333", FilerURL: "localhost:8888"}
	s, err := NewStorage(ProviderRustFS, nil, cfg)
	require.NoError(t, err)
	assert.NotNil(t, s)

	// RustFS operations are placeholders
	_, err = s.Upload(nil, "test", nil, UploadOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}
