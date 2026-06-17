package model

import "time"

type DeviceToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(512);uniqueIndex;not null" json:"token"`
	Platform  string    `gorm:"type:varchar(10);not null;check:platform IN ('ios','android','web')" json:"platform"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (DeviceToken) TableName() string { return "device_tokens" }
