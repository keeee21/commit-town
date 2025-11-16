package models

import (
	"time"
)

// UserStreak 連続コミット期間(streak)を管理する履歴テーブル
type UserStreak struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"index"`
	StartDate time.Time
	EndDate   *time.Time
	Length    int
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
