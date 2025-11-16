package models

import (
	"time"

	"gorm.io/datatypes"
)

// RepoDailyCommitLog GitHub APIから取得した、リポジトリ単位×日次のコミット集計
type RepoDailyCommitLog struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	UserRepoID  uint64         `gorm:"index"`
	CommitDate  time.Time      `gorm:"index"`
	CommitCount int
	RawData     datatypes.JSON
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	// Relations
	UserRepository UserRepository `gorm:"foreignKey:UserRepoID;references:ID"`
}
