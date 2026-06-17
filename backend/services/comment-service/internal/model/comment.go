package model

import "time"

// Comment represents a comment on a post.
type Comment struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PostID         string    `gorm:"type:uuid;not null;index:idx_comments_post" json:"post_id"`
	AuthorID       string    `gorm:"type:uuid;not null;index" json:"author_id"`
	AuthorNickname string    `gorm:"type:varchar(100);not null" json:"author_nickname"`
	AuthorAvatar   string    `gorm:"type:text;default:''" json:"author_avatar"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	IsAIGenerated  bool      `gorm:"default:false" json:"is_ai_generated"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}
