package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid gmail", "test@gmail.com", true},
		{"valid domain", "user@example.co.uk", true},
		{"valid plus", "user+tag@example.com", true},
		{"valid dot", "user.name@example.com", true},
		{"empty", "", false},
		{"no at", "notanemail", false},
		{"no domain", "test@", false},
		{"no user", "@example.com", false},
		{"double at", "a@b@c.com", false},
		{"spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"valid 8 chars", "Abc12345", true},
		{"valid long", "securePassword123!", true},
		{"too short", "Abc1", false},
		{"empty", "", false},
		{"only letters", "abcdefgh", false},
		{"only numbers", "12345678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateNickname(t *testing.T) {
	tests := []struct {
		name     string
		nickname string
		valid    bool
	}{
		{"valid", "小温暖", true},
		{"valid english", "WarmHeart", true},
		{"valid mixed", "Hello123", true},
		{"empty", "", false},
		{"too long", strings.Repeat("a", 101), false},
		{"just spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNickname(tt.nickname)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		maxLen  int
		valid   bool
	}{
		{"valid", "这是一段正常的内容", 5000, true},
		{"empty", "", 5000, false},
		{"too long", strings.Repeat("a", 5001), 5000, false},
		{"exactly max", strings.Repeat("a", 5000), 5000, true},
		{"whitespace only", "   \n  \t  ", 5000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContent(tt.content, tt.maxLen)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
