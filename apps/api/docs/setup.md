# セットアップガイド

## 前提条件

- Go 1.25.2 以上
- PostgreSQL 16
- Docker & Docker Compose（推奨）
- Make

## 環境構築

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd commit-tree
```

### 2. 環境変数の設定

```bash
# apps/api/.env を作成
cd apps/api
cp .env.example .env
```

`.env` ファイルを編集:
```env
DATABASE_URL=postgresql://commit_tree:commit_tree@localhost:5432/commit_tree?sslmode=disable
PORT=8080
```

### 3. データベースの起動

プロジェクトルートから:

```bash
# Docker Composeでデータベースを起動
./compose.sh up -d

# ログの確認
./compose.sh logs -f postgres

# 停止する場合
./compose.sh down
```

### 4. 依存関係のインストール

```bash
cd apps/api
go mod download
```

### 5. アプリケーションの起動

```bash
# 開発モード
make run

# または直接実行
go run main.go
```

アプリケーションは `http://localhost:8080` で起動します。

## ビルド

### 開発ビルド

```bash
make build
```

バイナリは `bin/main` として生成されます。

### 本番ビルド

```bash
go build -ldflags="-w -s" -o bin/main main.go
```

## データベースマイグレーション

GORMのオートマイグレーション機能を使用しています。

- アプリケーション起動時に自動的にマイグレーションが実行されます
- 新しいモデルを追加した場合は `db/database.go` の `AutoMigrate` に登録してください

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &gateway.YourNewEntity{},
    )
}
```

## API エンドポイント

### Health Check

```bash
curl http://localhost:8080/health
```

レスポンス:
```json
{
  "status": "ok"
}
```

## 開発コマンド

### すべてのコマンド

```bash
# アプリケーション実行
make run

# ビルド
make build

# OpenAPI型生成（Go側）
make generate
```

## トラブルシューティング

### データベース接続エラー

**エラー**: `Failed to connect to database`

**解決方法**:
1. PostgreSQLが起動しているか確認
   ```bash
   ./compose.sh ps
   ```
2. `.env` の `DATABASE_URL` が正しいか確認
3. データベースコンテナのログを確認
   ```bash
   ./compose.sh logs postgres
   ```

### ポートが使用中

**エラー**: `bind: address already in use`

**解決方法**:
1. 別のポートを使用（`.env` の `PORT` を変更）
2. または既存のプロセスを停止
   ```bash
   lsof -ti:8080 | xargs kill -9
   ```

### Go モジュールのエラー

**エラー**: `cannot find module`

**解決方法**:
```bash
go mod tidy
go mod download
```

## ディレクトリ構造

```
apps/api/
├── config/          # 設定ファイル
├── controller/      # HTTPコントローラー
├── db/              # データベース接続
├── docs/            # ドキュメント
├── gateway/         # リポジトリ実装
├── middleware/      # ミドルウェア
├── migrate/         # マイグレーション
├── model/           # ドメインモデル
├── repository/      # リポジトリインターフェース
├── router/          # ルーティング
├── usecase/         # ビジネスロジック
├── util/            # ユーティリティ
├── .env             # 環境変数（gitignore）
├── .env.example     # 環境変数テンプレート
├── go.mod           # Go モジュール
├── go.sum           # Go 依存関係
├── main.go          # エントリーポイント
└── Makefile         # Makeコマンド
```

## 次のステップ

1. [アーキテクチャドキュメント](./architecture.md)を読む
2. [開発ガイド](./development.md)を確認する
3. 新しい機能を追加する
