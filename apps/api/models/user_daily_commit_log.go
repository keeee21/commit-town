package models

import (
	"time"
)

// UserDailyCommitLog 全リポジトリを合算した、ユーザー単位の日次活動集計
type UserDailyCommitLog struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	UserID       uint64    `gorm:"index"`
	Date         time.Time `gorm:"index"`
	TotalCommits int
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
