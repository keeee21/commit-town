# Actions ディレクトリ

## 概要

`actions/` ディレクトリには、Server Actions とクライアントサイドのアクション関数を配置します。これらのファイルは、フォーム送信やユーザーインタラクションに応じたサーバーサイド処理を定義します。

## 配置するファイル

### Server Actions (`[xxx]-actions.ts`)

サーバーサイドでのみ実行されるアクション関数を定義します。

```typescript
// actions/update-user-actions.ts
import 'server-only'

export async function updateUserAction(formData: FormData) {
  // サーバーサイドでのみ実行される処理
  const userId = formData.get('userId') as string
  const name = formData.get('name') as string

  // データベース更新処理など
  return { status: 'success', message: 'User updated' }
}
```

## ベストプラクティス

詳細な実装方法については [Server Actions のベストプラクティス](../react/SERVER_ACTIONS.md) を参照してください。

### ファイル命名規則

- `server-actions.ts`: Server Actions の定義
- `client-actions.ts`: クライアントアクションの定義
- `action-types.ts`: アクション関連の型定義
- `action-utils.ts`: アクション用のユーティリティ関数

### ディレクトリ構造例

```
actions/
├── server-actions.ts    # Server Actions
├── client-actions.ts    # クライアントアクション
├── action-types.ts      # 型定義
└── action-utils.ts      # ユーティリティ関数
```

## 注意点

- Server Actions には必ず `'server-only'` をインポートする
- クライアントアクションには `'use client'` ディレクティブを付与する
- 型定義は `action-types.ts` に集約し、他のファイルからインポートして使用する
- アクション関数は単一責任の原則に従い、1つのファイルに複数の関連するアクションを配置する
