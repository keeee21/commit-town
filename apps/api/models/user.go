package models

import (
	"time"

	"gorm.io/gorm"
)

// User GitHubアカウントを基にしたアプリユーザー情報
type User struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement"`
	GitHubUserID   uint64         `gorm:"uniqueIndex;column:github_user_id"` // GitHub API の profile.id (変更不可、一意)
	GitHubUsername string         `gorm:"size:100;column:github_username"`   // GitHub API の profile.login (変更可能)
	Email          string         `gorm:"size:255"`                           // メールアドレス
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relations
	Repositories    []UserRepository     `gorm:"foreignKey:UserID"`
	DailyCommitLogs []UserDailyCommitLog `gorm:"foreignKey:UserID"`
	Streaks         []UserStreak         `gorm:"foreignKey:UserID"`
}
