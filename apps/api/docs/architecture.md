# API アーキテクチャ

## クリーンアーキテクチャの概要

このAPIはクリーンアーキテクチャに基づいて設計されています。各レイヤーは明確な責務を持ち、依存関係は外側から内側に向かって一方向になっています。

## ディレクトリ構成

```
apps/api/
├── config/          # アプリケーション設定
├── controller/      # コントローラー層（プレゼンテーション層）
├── db/              # データベース接続管理
├── gateway/         # 外部アクセス層（インフラ層）
├── middleware/      # HTTPミドルウェア
├── migrate/         # データベースマイグレーション
├── model/           # ドメインモデル（エンティティ）
├── repository/      # リポジトリ（インターフェース + 実装）
├── router/          # ルーティング設定
├── usecase/         # ユースケース層（ビジネスロジック）
├── util/            # ユーティリティ関数
├── validator/       # バリデーション層
└── docs/            # ドキュメント
```

## レイヤーの責務

### 1. Model（ドメイン層）
**パス**: `model/`

- ビジネスロジックの中心となるドメインモデル
- 他のレイヤーに依存しない純粋なビジネスオブジェクト
- データベースやフレームワークの詳細から独立

**例**:
```go
type User struct {
    ID        uuid.UUID
    Name      string
    Email     string
    CreatedAt time.Time
}
```

### 2. Repository（データアクセス層）
**パス**: `repository/`

- データアクセスのインターフェース定義
- データベース操作の実装（GORM使用）
- DBエンティティの定義とドメインモデルとの変換
- ユースケース層がこのインターフェースに依存

**Repository層に含まれるもの**:
1. インターフェース定義
2. GORM実装
3. DBエンティティ定義

**例**:
```go
// インターフェース定義
type IUserRepository interface {
    FindByID(ctx context.Context, id string) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
}

// DBエンティティ
type UserEntity struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    CreatedAt time.Time
}

// ToModel: Entity → Domain Model 変換
func (e *UserEntity) ToModel() *model.User { ... }

// GORM実装
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    var entity UserEntity
    if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return entity.ToModel(), nil
}
```

### 3. Usecase（ビジネスロジック層）
**パス**: `usecase/`

- アプリケーションのビジネスロジックを実装
- Repositoryインターフェースを使用してデータにアクセス
- ドメインモデルを操作して業務要件を実現

**例**:
```go
type UserUsecase interface {
    GetUser(ctx context.Context, id string) (*model.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}
```

### 4. Controller（プレゼンテーション層）
**パス**: `controller/`

- HTTPリクエストとレスポンスの処理
- リクエストのバリデーション
- Usecaseを呼び出してレスポンスを返す
- HTTPに関する詳細のみを扱う

**例**:
```go
type UserController struct {
    userUsecase usecase.UserUsecase
}

func (c *UserController) GetUser(ctx echo.Context) error {
    // HTTPリクエストの処理
    // Usecaseの呼び出し
    // HTTPレスポンスの返却
}
```

### 5. Gateway（外部サービス連携層）
**パス**: `gateway/`

- 外部サービス（Google認証、AWS S3、外部APIなど）との連携
- 外部サービスのインターフェース定義と実装
- 外部サービス固有のロジックをカプセル化

**重要**: データベース操作はRepository層で行い、Gateway層は外部サービス専用

**例**:
```go
// Google認証リポジトリのインターフェース
type IGoogleAuthRepository interface {
    VerifyIDToken(idToken string) (*model.GoogleAuthPayload, error)
}

// Google認証の実装
type googleAuthRepo struct{}

func NewGoogleAuthRepository() IGoogleAuthRepository {
    return &googleAuthRepo{}
}

func (r *googleAuthRepo) VerifyIDToken(idToken string) (*model.GoogleAuthPayload, error) {
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

**Gateway層の使用例**:
- Google OAuth認証
- AWS S3へのファイルアップロード
- 外部決済API連携
- メール送信サービス連携

### 6. Router（ルーティング層）
**パス**: `router/`

- HTTPエンドポイントとコントローラーのマッピング
- ルーティング設定の一元管理

**例**:
```go
func SetupRoutes(e *echo.Echo, userController *controller.UserController) {
    e.GET("/users/:id", userController.GetUser)
}
```

### 7. DB（データベース管理）
**パス**: `db/`

- データベース接続の初期化
- マイグレーションの実行

### 8. Middleware（ミドルウェア）
**パス**: `middleware/`

- 認証、ロギング、CORS等の横断的関心事
- リクエスト/レスポンスの前処理・後処理

### 9. Config（設定管理）
**パス**: `config/`

- アプリケーション設定の読み込み
- 環境変数の管理

### 10. Validator（バリデーション層）
**パス**: `validator/`

- **ビジネスロジックに関わるバリデーション**を担当
- Controllerの肥大化を防ぐ
- OpenAPIでは表現できない複雑なバリデーションを実装

**OpenAPIとの役割分担**:
- **OpenAPI**: 基本的なデータ構造の検証（型、必須項目、文字列長、フォーマット）
- **Validator層**: ビジネスルールに基づく検証（DB整合性、複数フィールドの関連チェック）

**例**:
```go
type UserValidator struct {
    userRepo repository.IUserRepository
}

type CreateUserInput struct {
    Name  string
    Email string
}

func (v *UserValidator) ValidateCreateUser(ctx context.Context, input CreateUserInput) error {
    // OpenAPIでは表現できないDB整合性チェック
    exists, err := v.userRepo.ExistsByEmail(ctx, input.Email)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("email already registered")
    }

    // 複雑なビジネスルール
    if strings.Contains(input.Name, "@") {
        return fmt.Errorf("name cannot contain @ symbol")
    }

    return nil
}
```

**使用場所**:
- Controllerから呼び出される
- OpenAPIバリデーション後に実行される
- 単体テストが容易

### 11. Util（ユーティリティ）
**パス**: `util/`

- 共通で使用されるヘルパー関数
- 文字列操作、日付変換などの汎用的な処理

## 依存関係の方向

```
┌─────────────────────────────────────────────────────────┐
│                     main.go (起動)                        │
└───────────────────┬─────────────────────────────────────┘
                    │
        ┌───────────┼───────────────┐
        │           │               │
        ▼           ▼               ▼
   ┌────────┐  ┌────────┐  ┌──────────┐
   │ Router │  │  DB    │  │  Config  │
   └────┬───┘  └───┬────┘  └──────────┘
        │          │
        ▼          │
   ┌──────────┐   │
   │Controller│   │
   └────┬─────┘   │
        │         │
        ▼         │
   ┌─────────┐   │
   │ Usecase │   │
   └────┬────┘   │
        │        │
   ┌────┴────────┴───────┐
   │                     │
   ▼                     ▼
┌──────┐          ┌────────────┐
│Model │          │ Repository │
└──────┘          │(interface  │
   ▲              │   +        │
   │              │implementation)
   │              └─────┬──────┘
   │                    │
   └────────────────────┘

   外部サービス連携
   ┌─────────┐
   │ Gateway │ (Google Auth, S3, etc.)
   └─────────┘
```

**重要なポイント**:
- 内側の層は外側の層に依存しない
- **Repository は同じファイル内でインターフェース定義と実装の両方を持つ**
- Usecase は Repository のインターフェースに依存（具体的な実装には依存しない）
- **Gateway はDB操作ではなく、外部サービス連携専用**
- この構造により、テストが容易で、変更に強い設計になる

## データフロー

### リクエスト → レスポンス

**データベース操作の場合**:
```
HTTP Request
    ↓
Router
    ↓
Controller (リクエストの受け取り)
    ↓
Usecase (ビジネスロジック)
    ↓
Repository Interface
    ↓
Repository Implementation (GORM操作)
    ↓
Database
    ↓
Repository (Entity → Model変換)
    ↓
Usecase (結果の処理)
    ↓
Controller (レスポンス生成)
    ↓
HTTP Response
```

**外部サービス連携の場合**:
```
HTTP Request
    ↓
Router
    ↓
Controller (リクエストの受け取り)
    ↓
Usecase (ビジネスロジック)
    ↓
Gateway Interface
    ↓
Gateway Implementation (外部API呼び出し)
    ↓
External Service (Google Auth, S3, etc.)
    ↓
Gateway (レスポンス変換)
    ↓
Usecase (結果の処理)
    ↓
Controller (レスポンス生成)
    ↓
HTTP Response
```

## テスタビリティ

クリーンアーキテクチャにより、各レイヤーを独立してテスト可能:

1. **Usecase のテスト**: Repository や Gateway インターフェースをモックして、ビジネスロジックのみをテスト
2. **Controller のテスト**: Usecase をモックして、HTTP処理のみをテスト
3. **Repository のテスト**: 実際のDBまたはテスト用DBを使用してテスト
4. **Gateway のテスト**: 外部サービスをモックして、連携ロジックをテスト

## 新機能の追加方法

### データベース操作が必要な機能の場合
1. **Model を定義** (`model/`)
2. **Repository でインターフェースと実装を作成** (`repository/`)
   - インターフェース定義
   - DBエンティティ定義
   - GORM実装
3. **Usecase を作成** (`usecase/`)
4. **Controller を作成** (`controller/`)
5. **Router に登録** (`router/`)
6. **main.go で依存性注入**

### 外部サービス連携が必要な機能の場合
1. **Model を定義** (`model/`) - 必要に応じて
2. **Gateway でインターフェースと実装を作成** (`gateway/`)
   - インターフェース定義
   - 外部サービス連携実装
3. **Usecase を作成** (`usecase/`)
4. **Controller を作成** (`controller/`)
5. **Router に登録** (`router/`)
6. **main.go で依存性注入**

## 環境変数

- `DATABASE_URL`: PostgreSQL接続文字列
- `PORT`: APIサーバーのポート番号（デフォルト: 8080）

## 使用技術

- **Web Framework**: Echo v4
- **ORM**: GORM
- **Database**: PostgreSQL
- **環境変数**: godotenv
- **UUID**: google/uuid
