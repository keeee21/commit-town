package repository

import (
	"github.com/keeee21/commit-town/api/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByGitHubUserID GitHub User IDでユーザーを検索
func (userRepo *UserRepository) FindByGitHubUserID(githubUserID uint64) (*models.User, error) {
	var user models.User
	err := userRepo.db.Where("github_user_id = ?", githubUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 新規ユーザーを作成
func (userRepo *UserRepository) Create(user *models.User) error {
	return userRepo.db.Create(user).Error
}

// Update ユーザー情報を更新
func (userRepo *UserRepository) Update(user *models.User) error {
	return userRepo.db.Save(user).Error
}

// Upsert ユーザーを作成または更新（GitHub User IDで判定）
func (userRepo *UserRepository) Upsert(user *models.User) error {
	existing, err := userRepo.FindByGitHubUserID(user.GitHubUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 新規作成
			return userRepo.Create(user)
		}
		return err
	}

	// 既存レコードを更新
	user.ID = existing.ID
	user.CreatedAt = existing.CreatedAt
	return userRepo.Update(user)
}
