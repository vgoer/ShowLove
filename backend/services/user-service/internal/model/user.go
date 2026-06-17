// Package model defines the data models for the user service.
package model

import (
	"time"
)

// User represents a registered user.
type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"` // bcrypt hash, never serialized
	Nickname     string    `gorm:"type:varchar(100);not null" json:"nickname"`
	AvatarURL    string    `gorm:"type:text;default:''" json:"avatar_url"`
	Bio          string    `gorm:"type:text;default:''" json:"bio"`
	KindnessScore int     `gorm:"default:0" json:"kindness_score"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}

// RefreshToken represents a stored refresh token.
type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(512);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides the default table name.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
