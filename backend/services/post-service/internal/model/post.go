// Package model defines the data models for the post service.
package model

import (
	"time"
)

// Post represents a community post.
type Post struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AuthorID         string    `gorm:"type:uuid;not null;index" json:"author_id"`
	AuthorNickname   string    `gorm:"type:varchar(100);not null" json:"author_nickname"`
	AuthorAvatar     string    `gorm:"type:text;default:''" json:"author_avatar"`
	Content          string    `gorm:"type:text;not null" json:"content"`
	MoodTag          string    `gorm:"type:varchar(50);not null;index" json:"mood_tag"`
	Images           string    `gorm:"type:jsonb;default:'[]'" json:"images"` // JSON array of URLs
	VoiceURL         string    `gorm:"type:text" json:"voice_url"`
	StickerHug       int32     `gorm:"default:0" json:"sticker_hug"`
	StickerCheer     int32     `gorm:"default:0" json:"sticker_cheer"`
	StickerUnderstand int32    `gorm:"default:0" json:"sticker_understand"`
	CommentCount     int32     `gorm:"default:0;index" json:"comment_count"`
	HasAIReply       bool      `gorm:"default:false" json:"has_ai_reply"`
	IsReported       bool      `gorm:"default:false" json:"is_reported"`
	IsHidden         bool      `gorm:"default:false" json:"is_hidden"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default table name.
func (Post) TableName() string {
	return "posts"
}
