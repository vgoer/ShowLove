package model

// MoodEntry represents a daily mood record.
type MoodEntry struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_moods_user_date" json:"user_id"`
	MoodLevel int32  `gorm:"not null;check:mood_level >= 1 AND mood_level <= 10" json:"mood_level"`
	MoodLabel string `gorm:"type:varchar(50);not null" json:"mood_label"`
	Note      string `gorm:"type:text" json:"note"`
	CreatedAt string `gorm:"type:date;not null;uniqueIndex:idx_moods_user_date;default:CURRENT_DATE" json:"created_at"`
}

func (MoodEntry) TableName() string { return "mood_entries" }
