package dto

// UpsertUserRequest ユーザー作成/更新リクエスト
type UpsertUserRequest struct {
	GitHubUserID   uint64 `json:"github_user_id" validate:"required"`
	GitHubUsername string `json:"github_username" validate:"required"`
	Email          string `json:"email"`
}

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID             uint64 `json:"id"`
	GitHubUserID   uint64 `json:"github_user_id"`
	GitHubUsername string `json:"github_username"`
	Email          string `json:"email"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
