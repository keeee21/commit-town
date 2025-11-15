# 開発ガイド

## 新機能の追加手順

このドキュメントでは、クリーンアーキテクチャに従って新機能を追加する方法を説明します。

### 例: ユーザー管理機能の追加

#### 1. ドメインモデルの作成

`model/user.go`:
```go
package model

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID
    Name      string
    Email     string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### 2. リポジトリの作成（インターフェース + DBエンティティ + GORM実装）

`repository/user_repository.go`:
```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/keeee21/commit-town/api/model"
    "gorm.io/gorm"
)

// インターフェース定義
type IUserRepository interface {
    FindAll(ctx context.Context) ([]model.User, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id string) error
}

// DBエンティティ
type UserEntity struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (UserEntity) TableName() string {
    return "users"
}

func (e *UserEntity) BeforeCreate(tx *gorm.DB) error {
    if e.ID == uuid.Nil {
        e.ID = uuid.New()
    }
    return nil
}

// Entity → Model 変換
func (e *UserEntity) ToModel() *model.User {
    return &model.User{
        ID:        e.ID,
        Name:      e.Name,
        Email:     e.Email,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}

// Model → Entity 変換
func (e *UserEntity) FromModel(m *model.User) {
    e.ID = m.ID
    e.Name = m.Name
    e.Email = m.Email
    e.CreatedAt = m.CreatedAt
    e.UpdatedAt = m.UpdatedAt
}

// GORM実装
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
    var entities []UserEntity
    if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
        return nil, err
    }

    users := make([]model.User, len(entities))
    for i, entity := range entities {
        users[i] = *entity.ToModel()
    }
    return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    uid, err := uuid.Parse(id)
    if err != nil {
        return nil, err
    }

    var entity UserEntity
    if err := r.db.WithContext(ctx).Where("id = ?", uid).First(&entity).Error; err != nil {
        return nil, err
    }
    return entity.ToModel(), nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    var entity UserEntity
    entity.FromModel(user)

    if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
        return err
    }
    *user = *entity.ToModel()
    return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
    var entity UserEntity
    entity.FromModel(user)

    if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
        return err
    }
    *user = *entity.ToModel()
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
    uid, err := uuid.Parse(id)
    if err != nil {
        return err
    }
    return r.db.WithContext(ctx).Where("id = ?", uid).Delete(&UserEntity{}).Error
}
```

**重要ポイント**:
- Repository層で**インターフェース、DBエンティティ、実装**の全てを定義
- Gateway層はDB操作には使用しない（外部サービス専用）

#### 3. ユースケースの作成

`usecase/user_usecase.go`:
```go
package usecase

import (
    "context"
    "github.com/keeee21/commit-town/api/model"
    "github.com/keeee21/commit-town/api/repository"
)

type IUserUsecase interface {
    GetAll(ctx context.Context) ([]model.User, error)
    GetByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, name, email string) (*model.User, error)
    Update(ctx context.Context, id, name, email string) (*model.User, error)
    Delete(ctx context.Context, id string) error
}

type userUsecase struct {
    userRepo repository.IUserRepository
}

func NewUserUsecase(userRepo repository.IUserRepository) IUserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetAll(ctx context.Context) ([]model.User, error) {
    return u.userRepo.FindAll(ctx)
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*model.User, error) {
    return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, name, email string) (*model.User, error) {
    user := &model.User{
        Name:  name,
        Email: email,
    }

    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (u *userUsecase) Update(ctx context.Context, id, name, email string) (*model.User, error) {
    user, err := u.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    user.Name = name
    user.Email = email

    if err := u.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
    return u.userRepo.Delete(ctx, id)
}
```

#### 4. OpenAPIスキーマの定義

`packages/openapi/schema.yaml`:
```yaml
components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
          minLength: 2
          maxLength: 100
          description: ユーザー名
        email:
          type: string
          format: email
          description: メールアドレス

    UpdateUserRequest:
      type: object
      properties:
        name:
          type: string
          minLength: 2
          maxLength: 100
        email:
          type: string
          format: email
```

**OpenAPIでのバリデーション**:
- 型チェック（string, number, etc.）
- 必須項目（required）
- 文字列長（minLength, maxLength）
- フォーマット（email, uuid, date-time, etc.）

#### 5. バリデーターの作成（ビジネスルール）

`validator/user_validator.go`:
```go
package validator

import (
    "context"
    "fmt"
    "strings"

    "github.com/keeee21/commit-town/api/repository"
)

type UserValidator struct {
    userRepo repository.IUserRepository
}

func NewUserValidator(userRepo repository.IUserRepository) *UserValidator {
    return &UserValidator{
        userRepo: userRepo,
    }
}

type CreateUserInput struct {
    Name  string
    Email string
}

type UpdateUserInput struct {
    Name  string
    Email string
}

// ValidateCreateUser - OpenAPIでは表現できないビジネスルールをチェック
func (v *UserValidator) ValidateCreateUser(ctx context.Context, input CreateUserInput) error {
    // DB整合性チェック: メールアドレスの重複確認
    exists, err := v.userRepo.ExistsByEmail(ctx, input.Email)
    if err != nil {
        return fmt.Errorf("failed to check email existence: %w", err)
    }
    if exists {
        return fmt.Errorf("email already registered")
    }

    // ビジネスルール: 名前に@を含めない
    if strings.Contains(input.Name, "@") {
        return fmt.Errorf("name cannot contain @ symbol")
    }

    return nil
}

func (v *UserValidator) ValidateUpdateUser(ctx context.Context, userID string, input UpdateUserInput) error {
    // メールアドレスを変更する場合は重複チェック
    if input.Email != "" {
        exists, err := v.userRepo.ExistsByEmailExcludingUser(ctx, input.Email, userID)
        if err != nil {
            return fmt.Errorf("failed to check email existence: %w", err)
        }
        if exists {
            return fmt.Errorf("email already registered")
        }
    }

    // ビジネスルール
    if input.Name != "" && strings.Contains(input.Name, "@") {
        return fmt.Errorf("name cannot contain @ symbol")
    }

    return nil
}
```

**Validator層の役割**:
- データベース整合性チェック
- 複数フィールドにまたがる検証
- OpenAPIでは表現できない複雑なビジネスルール

#### 6. コントローラーの作成

`controller/user_controller.go`:
```go
package controller

import (
    "net/http"
    "github.com/keeee21/commit-town/api/usecase"
    "github.com/keeee21/commit-town/api/validator"
    "github.com/labstack/echo/v4"
)

type UserController struct {
    userUsecase   usecase.IUserUsecase
    userValidator *validator.UserValidator
}

func NewUserController(userUsecase usecase.IUserUsecase, userValidator *validator.UserValidator) *UserController {
    return &UserController{
        userUsecase:   userUsecase,
        userValidator: userValidator,
    }
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (c *UserController) GetAll(ctx echo.Context) error {
    users, err := c.userUsecase.GetAll(ctx.Request().Context())
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetByID(ctx echo.Context) error {
    id := ctx.Param("id")

    user, err := c.userUsecase.GetByID(ctx.Request().Context(), id)
    if err != nil {
        return ctx.JSON(http.StatusNotFound, map[string]string{
            "error": "User not found",
        })
    }
    return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Create(ctx echo.Context) error {
    var req CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request",
        })
    }

    // ① OpenAPIで基本的なバリデーションは自動実施済み（型、必須項目、文字列長など）

    // ② Validator層でビジネスルールをチェック
    if err := c.userValidator.ValidateCreateUser(ctx.Request().Context(), validator.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    user, err := c.userUsecase.Create(ctx.Request().Context(), req.Name, req.Email)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) Update(ctx echo.Context) error {
    id := ctx.Param("id")

    var req UpdateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request",
        })
    }

    // バリデーション
    if err := c.userValidator.ValidateUpdateUser(validator.UpdateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    user, err := c.userUsecase.Update(ctx.Request().Context(), id, req.Name, req.Email)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Delete(ctx echo.Context) error {
    id := ctx.Param("id")

    if err := c.userUsecase.Delete(ctx.Request().Context(), id); err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return ctx.NoContent(http.StatusNoContent)
}
```

**バリデーションの流れ**:
```
リクエスト
  ↓
① OpenAPIバリデーション（自動）
  - 型チェック
  - 必須項目
  - 文字列長
  - フォーマット（email, uuidなど）
  ↓
② Validator層（手動実装）
  - DB整合性チェック
  - ビジネスルール
  - 複数フィールドの関連チェック
  ↓
Usecase
```

#### 7. ルーティングの追加

`router/router.go`:
```go
func SetupRoutes(
    e *echo.Echo,
    healthController *controller.HealthController,
    userController *controller.UserController, // 追加
) {
    e.GET("/health", healthController.Check)

    // User routes
    users := e.Group("/users")
    users.GET("", userController.GetAll)
    users.GET("/:id", userController.GetByID)
    users.POST("", userController.Create)
    users.PUT("/:id", userController.Update)
    users.DELETE("/:id", userController.Delete)
}
```

#### 8. main.goで依存性注入

`main.go`:
```go
func main() {
    // ... (データベース接続など)

    // Repository
    userRepo := repository.NewUserRepository(database)

    // Usecase
    userUsecase := usecase.NewUserUsecase(userRepo)

    // Validator（Repositoryを注入）
    userValidator := validator.NewUserValidator(userRepo)

    // Controller
    userController := controller.NewUserController(userUsecase, userValidator)

    // ... (Echo初期化など)

    // Routes
    router.SetupRoutes(e, healthController, userController)

    // ... (サーバー起動)
}
```

#### 9. マイグレーションの追加

`db/database.go`:
```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &repository.UserEntity{}, // 追加
    )
}
```

---

## 外部サービス連携機能の追加（Gateway層の使用例）

### 例: Google認証機能の追加

#### 1. Gatewayでインターフェースと実装を作成

`gateway/google_auth.go`:
```go
package gateway

import (
    "context"
    "fmt"
    "os"

    "github.com/keeee21/commit-town/api/model"
    "google.golang.org/api/idtoken"
)

// インターフェース定義
type IGoogleAuthGateway interface {
    VerifyIDToken(idToken string) (*model.GoogleAuthPayload, error)
}

// 実装
type googleAuthGateway struct{}

func NewGoogleAuthGateway() IGoogleAuthGateway {
    return &googleAuthGateway{}
}

func (g *googleAuthGateway) VerifyIDToken(idToken string) (*model.GoogleAuthPayload, error) {
    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    payload, err := idtoken.Validate(context.Background(), idToken, clientID)
    if err != nil {
        return nil, fmt.Errorf("token validation failed: %w", err)
    }

    return &model.GoogleAuthPayload{
        Email: fmt.Sprintf("%v", payload.Claims["email"]),
        Name:  fmt.Sprintf("%v", payload.Claims["name"]),
        Sub:   payload.Subject,
    }, nil
}
```

#### 2. Usecaseで使用

```go
type authUsecase struct {
    googleAuth gateway.IGoogleAuthGateway
    userRepo   repository.IUserRepository
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context, idToken string) (*model.User, error) {
    // Gatewayで外部サービス（Google）と連携
    authPayload, err := u.googleAuth.VerifyIDToken(idToken)
    if err != nil {
        return nil, err
    }

    // Repositoryでデータベース操作
    user, err := u.userRepo.FindByEmail(ctx, authPayload.Email)
    if err != nil {
        // 新規ユーザー作成
        user = &model.User{
            Email: authPayload.Email,
            Name:  authPayload.Name,
        }
        if err := u.userRepo.Create(ctx, user); err != nil {
            return nil, err
        }
    }

    return user, nil
}
```

---

## コーディング規約

### 命名規則

- **ファイル名**: `snake_case.go`
- **パッケージ名**: 小文字の単一単語
- **構造体**: `PascalCase`
- **インターフェース**: `I` プレフィックス + `PascalCase` (例: `IUserRepository`)
- **関数/メソッド**: `PascalCase`（公開）、`camelCase`（非公開）
- **変数**: `camelCase`

### エラーハンドリング

```go
// Good
user, err := r.userRepo.FindByID(ctx, id)
if err != nil {
    return nil, fmt.Errorf("failed to find user: %w", err)
}

// Bad - エラーを無視しない
user, _ := r.userRepo.FindByID(ctx, id)
```

### コンテキストの使用

すべてのDB操作、外部API呼び出しには `context.Context` を渡す:

```go
func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    var entity UserEntity
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
        return nil, err
    }
    return entity.ToModel(), nil
}
```

## テストの書き方

### Repositoryのテスト

```go
package repository_test

import (
    "testing"
    "context"
    // ...
)

func TestUserRepository_Create(t *testing.T) {
    // テスト用DBセットアップ
    // ...

    repo := repository.NewUserRepository(db)

    user := &model.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, user.ID)
}
```

### Usecaseのテスト

```go
package usecase_test

import (
    "testing"
    "context"
    // モックを使用
)

func TestUserUsecase_Create(t *testing.T) {
    mockRepo := &MockUserRepository{}
    uc := usecase.NewUserUsecase(mockRepo)

    // テスト実装
}
```

### Gatewayのテスト

```go
package gateway_test

import (
    "testing"
    // モックを使用
)

func TestGoogleAuthGateway_VerifyIDToken(t *testing.T) {
    // 外部サービスをモック
    // ...
}
```

## デバッグ

### ログの追加

```go
import "log"

log.Printf("User created: %+v", user)
```

### GORMのクエリログ

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

## パフォーマンス最適化

### N+1問題の回避

```go
// Bad
users, _ := db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}

// Good
db.Preload("Posts").Find(&users)
```

### インデックスの追加

```go
type UserEntity struct {
    Email string `gorm:"unique;not null;index"`
}
```

## OpenAPI との連携

OpenAPI スキーマを更新したら、型を再生成:

```bash
cd apps/api
make generate
```

## よくある質問

**Q: 新しいミドルウェアを追加するには？**

A: `middleware/` ディレクトリに作成し、`main.go` で登録します。

**Q: バリデーションはどこで行う？**

A: Controller層でリクエストのバリデーション、Usecase層でビジネスルールのバリデーションを行います。

**Q: トランザクションはどう扱う？**

A: Usecase層でトランザクションを開始し、複数のRepository操作をラップします。

**Q: Repository層とGateway層の使い分けは？**

A:
- **Repository層**: データベース操作（GORM使用）
- **Gateway層**: 外部サービス連携（Google Auth、AWS S3、外部APIなど）
