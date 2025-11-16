# featuresディレクトリのルール

## 概要

featuresディレクトリは、**Feature-Based Architecture（機能ベースアーキテクチャ）**と**Domain-Driven Design（ドメイン駆動設計）**の原則に基づいて設計されています。このディレクトリは、アプリケーションの機能をドメインごとに整理し、関連するコードを一箇所に集約することで、保守性とスケーラビリティを向上させます。

## 設計原則

### 1. Co-Location（共配置）の原則
関連するファイルを物理的に近くに配置することで、以下の利点を得られます：
- コードの関連性が視覚的に理解しやすい
- ファイル間の依存関係が明確
- リファクタリング時の影響範囲が把握しやすい

### 2. Single Responsibility Principle（単一責任の原則）
各featureディレクトリは、一つのドメインまたは機能にのみ責任を持ちます。

### 3. Encapsulation（カプセル化）
各featureは内部実装を隠蔽し、明確なインターフェースを通じてのみ外部とやり取りします。

## ルール

featuresディレクトリには以下のようなファイルを配置する。

- ドメインに関連するコンポーネント
- ドメインに関連するロジック
- ドメインに関連するテスト
- ドメインに関連するHooks
- ドメインに関連する型定義
- ドメインに関連する定数
- ドメインに関連する状態管理

また、featuresディレクトリに配置するファイルは Co-Location の原則に従いドメインの近くに配置する、かつ各ドメインはページルーティングに対応するように配置する。

## ディレクトリ構造

`apps/web/app/[domain-name]/page.tsx` 関連のファイルは以下のように配置する。

```
apps/web/features/
├── [domain-name]
│   ├── components/ # ドメインに関連するコンポーネント
│   │   ├── index.ts # Barrel export（推奨）
│   │   ├── [ComponentName].tsx
│   │   └── [ComponentName].test.tsx
│   ├── types/ # ドメインに関連する型定義
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [types].ts
│   ├── constants/ # ドメインに関連する定数
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [constants].ts
│   ├── hooks/ # ドメインに関連するHooks
│   │   ├── index.ts # Barrel export（推奨）
│   │   ├── use[hook-name].ts
│   │   └── use[hook-name].test.ts
│   ├── utils/ # ドメインに関連するユーティリティ（詳細は[UTILS.md](./UTILS.md)を参照）
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [utils].ts
│   ├── libs/ # ドメインに関連するライブラリ（詳細は[LIBS.md](./LIBS.md)を参照）
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [libs].ts
│   ├── actions/ # ドメインに関連するサーバーアクション
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [actions].ts
│   ├── schemas/ # ドメインに関連するスキーマ
│   │   ├── index.ts # Barrel export（推奨）
│   │   └── [schemas].ts
│   ├── test/ # ドメインに関連するテスト
│   │   ├── __mocks__/ # モックファイル
│   │   ├── fixtures/ # テストデータ
│   │   └── [test-files].test.ts
│   ├── index.ts # メインのBarrel export
│   └── ...any-directories/
```

## ベストプラクティス

### 1. Barrel Exports（バレルエクスポート）の活用
各サブディレクトリに`index.ts`ファイルを配置し、外部への公開インターフェースを明確にします。

```typescript
// features/user/components/index.ts
export { UserProfile } from './UserProfile';
export { UserList } from './UserList';
export { UserForm } from './UserForm';
```

### 2. 命名規則
- **ディレクトリ名**: kebab-case（例：`user-profile`）
- **コンポーネントファイル**: PascalCase（例：`UserProfile.tsx`）
- **Hookファイル**: camelCase with use prefix（例：`useUserProfile.ts`）
- **ユーティリティファイル**: camelCase（例：`userUtils.ts`）

### 3. 依存関係の管理
- 同じfeature内のファイルは自由に参照可能
- 異なるfeature間の依存は最小限に抑制
- 共通機能は`libs/`（[LIBS.md](./LIBS.md)を参照）や`utils/`（[UTILS.md](./UTILS.md)を参照）ディレクトリに配置

### 4. テスト戦略
- **Unit Tests**: 各ファイルと同じディレクトリに配置
- **Integration Tests**: `test/`ディレクトリに配置
- **E2E Tests**: プロジェクトルートの`e2e/`ディレクトリに配置

### 5. 型安全性の確保
- 各feature内で使用する型は`types/`ディレクトリに定義
- 外部APIとの型定義は`schemas/`ディレクトリに配置
- 型の再エクスポートは`index.ts`で管理

### 6. 状態管理
- グローバル状態は`libs/`ディレクトリの状態管理ライブラリで管理（詳細は[LIBS.md](./LIBS.md)を参照）
- ローカル状態は各コンポーネント内で管理
- 複雑な状態ロジックはカスタムHookに分離

## 禁止事項

### 1. 循環依存の禁止
feature間での循環依存は禁止。必要に応じて共通ライブラリを作成。

### 2. 直接的なimportの禁止
他のfeatureの内部実装に直接アクセスすることは禁止。公開されたインターフェースのみ使用。

### 3. グローバル変数の禁止
feature内でのグローバル変数の使用は禁止。適切な状態管理手法を使用。

## 移行ガイド

既存のコードをfeaturesディレクトリに移行する際の手順：

1. **ドメインの特定**: 機能の責任範囲を明確に定義
2. **ディレクトリの作成**: 適切なディレクトリ構造を作成
3. **ファイルの移動**: 関連ファイルを適切なディレクトリに移動
4. **Barrel exportの設定**: `index.ts`ファイルを作成
5. **import文の更新**: 新しいパスに合わせてimport文を更新
6. **テストの実行**: 全てのテストが正常に動作することを確認

## 参考資料

- [Feature-Based Architecture in React](https://reactjs.org/docs/thinking-in-react.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Co-location Principle](https://kentcdodds.com/blog/colocation)
