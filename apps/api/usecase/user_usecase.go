package usecase

import (
	"github.com/keeee21/commit-town/api/dto"
	"github.com/keeee21/commit-town/api/models"
	"github.com/keeee21/commit-town/api/repository"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// UpsertUser ユーザーを作成または更新
func (userUsecase *UserUsecase) UpsertUser(req *dto.UpsertUserRequest) (*dto.UserResponse, error) {
	user := &models.User{
		GitHubUserID:   req.GitHubUserID,
		GitHubUsername: req.GitHubUsername,
		Email:          req.Email,
	}

	if err := userUsecase.userRepo.Upsert(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:             user.ID,
		GitHubUserID:   user.GitHubUserID,
		GitHubUsername: user.GitHubUsername,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
