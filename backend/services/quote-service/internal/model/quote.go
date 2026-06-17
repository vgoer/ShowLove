package model

import "time"

type DailyQuote struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TextZH        string    `gorm:"type:text;not null" json:"text_zh"`
	TextEN        string    `gorm:"type:text;not null" json:"text_en"`
	Author        string    `gorm:"type:varchar(100)" json:"author"`
	BackgroundURL string    `gorm:"type:text" json:"background_url"`
	ScheduledDate string    `gorm:"type:date;uniqueIndex;not null" json:"scheduled_date"`
	Pushed        bool      `gorm:"default:false" json:"pushed"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (DailyQuote) TableName() string { return "daily_quotes" }
