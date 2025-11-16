# Server Actions と useActionState のベストプラクティス

## 概要

React 19の`useActionState`と`Server Actions`を使用する際のベストプラクティスについて説明します。特に、`useEffect`に頼らずにUIを適切に更新する方法に焦点を当てます。

## useEffectを使わない理由

`useEffect`を使った実装には以下の問題があります：

1. **パフォーマンスとUXの問題**: レンダリング遅延や追加のレンダリングサイクルが発生する
2. **コードの可読性の問題**: トーストの処理がフォームアクションのロジックと分離している

詳細については [useEffectの正しい使い方ルール](./USEEFFECT_RULES.md) を参照してください。

## 推奨される実装方法

### 1. useActionStateのコールバックでのハンドリング

```typescript
const [lastResult, action, isPending] = useActionState(
  updateWorkspaceAction,
  null,
  {
    onSuccess: (result) => {
      toast.success('Workspace updated')
      setOpen(false)
    },
    onError: (result) => {
      toast.error('Failed to update workspace')
    }
  }
)
```

### 2. withCallbacksを使用する方法（推奨）

より柔軟で再利用可能なアプローチとして、`withCallbacks`ヘルパー関数を使用します。

```typescript
type Callbacks<T, R = void> = {
  onStart?: () => R
  onEnd?: (reference: R) => void
  onSuccess?: (result: T) => void
  onError?: (result: T) => void
}

// シンプルなジェネリクスを使用（TYPESCRIPT.mdのルールに従う）
function withCallbacks<T, R = void>(
  fn: (...args: unknown[]) => Promise<T>,
  callbacks: Callbacks<T, R>
) {
  return async (...args: unknown[]) => {
    const promise = fn(...args)
    const reference = callbacks.onStart?.()
    const result = await promise

    if (reference) {
      callbacks.onEnd?.(reference)
    }

    if (result.status === 'success') {
      callbacks.onSuccess?.(result)
    }

    if (result.status === 'error') {
      callbacks.onError?.(result)
    }

    return promise
  }
}
```

#### 使用例

```typescript
const [lastResult, action, isPending] = useActionState(
  withCallbacks(updateWorkspaceAction, {
    onError() {
      toast.error('Failed to update workspace')
    },
    onSuccess() {
      toast.success('Workspace updated')
      setOpen(false)
    },
  }),
  null,
)
```

## 型定義

プロジェクトで使用する型定義の例：

```typescript
export type ActionState =
  | {
      message: string
      status: "SUCCESS" | "ERROR"
    }
  | null
  | undefined
```

型定義の詳細については [TypeScript ベストプラクティス](../typescript/TYPESCRIPT.md) を参照してください。

## 利点

1. **パフォーマンス向上**: 不要なレンダリングサイクルを避ける
2. **UX向上**: トースト表示などの遅延を最小化
3. **コードの可読性**: フォームアクションとUI更新のロジックが一箇所に集約
4. **再利用性**: `withCallbacks`ヘルパーは様々なServer Actionsで再利用可能

## Server Actionsのセキュリティ

### server-onlyモジュールの使用

Server Actionsはサーバーサイドでのみ実行されるべきコードです。機密情報や内部ビジネスロジックがクライアントに漏洩することを防ぐため、`server-only`モジュールを使用します。

```typescript
// server-actions.ts
import 'server-only'
import { secretKey } from './secrets'

export async function updateUserData(userId: string, data: unknown) {
  // 機密情報を使用したサーバーサイド処理
  const result = await processWithSecretKey(userId, data, secretKey)
  return { status: 'success', data: result }
}
```

### 環境変数の適切な管理

- サーバー専用の環境変数: `process.env.SECRET_KEY`
- クライアント公開用: `process.env.NEXT_PUBLIC_API_URL`

```typescript
// server-actions.ts
import 'server-only'

export async function fetchUserData(userId: string) {
  // サーバー専用の環境変数を使用
  const apiKey = process.env.SECRET_API_KEY
  const response = await fetch(`/api/users/${userId}`, {
    headers: { Authorization: `Bearer ${apiKey}` }
  })
  return response.json()
}
```

### 実験的なTaint API（Next.js 14+）

機密データがクライアントに渡されることを防ぐため、実験的なTaint APIを使用できます。

```typescript
// next.config.js
module.exports = {
  experimental: {
    taint: true,
  }
}

// server-actions.ts
import { experimental_taintObjectReference, experimental_taintUniqueValue } from 'react'

export async function getUserData(userId: string) {
  const data = await fetchUserFromDB(userId)

  // オブジェクト全体をクライアントに渡すことを禁止
  experimental_taintObjectReference('Do not pass user data to client', data)

  // 特定の値（トークンなど）をクライアントに渡すことを禁止
  experimental_taintUniqueValue('Do not pass tokens to client', data, data.token)

  return data
}
```

## 注意点

- `withCallbacks`の`onEnd`は`onStart`の返り値を受け取るため、`onStart`で何らかの参照を返す必要がある
- Server Actionsの結果の型（`status`プロパティ）に応じて適切にコールバックを設定する
- エラーハンドリングは`onError`コールバックで行い、成功時の処理は`onSuccess`で行う
- **セキュリティ**: Server Actionsには必ず`'server-only'`をインポートし、機密情報の漏洩を防ぐ
- **データアクセスレイヤー**: データベースアクセスなどの機密処理は専用のレイヤーで実装し、Server Actionsから呼び出す

## 参考資料

- [React 19 useActionState](https://ja.react.dev/reference/react/useActionState)
- [You Might Not Need an Effect](https://ja.react.dev/learn/you-might-not-need-an-effect)
- [useEffect で苦しまない！useActionState × Server Actions ベストプラクティス](https://zenn.dev/sc30gsw/articles/6b43b44e04e89e)
- [Next.js Server Only機能](https://zenn.dev/417/scraps/0844b20e6b0952)
