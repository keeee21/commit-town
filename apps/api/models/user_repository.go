package models

import (
	"time"
)

// UserRepository ユーザーがGUIで登録したGitHubリポジトリ情報
type UserRepository struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement"`
	UserID        uint64     `gorm:"index"`
	RepoOwner     string     `gorm:"size:100"`
	RepoName      string     `gorm:"size:100"`
	IsPublic      bool       `gorm:"default:true"`
	DeactivatedAt *time.Time
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`

	// Relations
	User                 User                   `gorm:"foreignKey:UserID;references:ID"`
	RepoDailyCommitLogs  []RepoDailyCommitLog   `gorm:"foreignKey:UserRepoID"`
}
