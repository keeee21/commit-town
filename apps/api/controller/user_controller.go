package controller

import (
	"net/http"

	"github.com/keeee21/commit-town/api/dto"
	"github.com/keeee21/commit-town/api/usecase"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

// UpsertUser ユーザーを作成または更新
func (userController *UserController) UpsertUser(ctx echo.Context) error {
	var req dto.UpsertUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// 簡易バリデーション
	if req.GitHubUserID == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "github_user_id is required",
		})
	}
	if req.GitHubUsername == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "github_username is required",
		})
	}

	user, err := userController.userUsecase.UpsertUser(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to upsert user",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}
