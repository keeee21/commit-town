# クライアントサイドでのHooksベストプラクティス

## 概要

Server Actionsを呼び出す際のクライアントサイドでのHooks（`useActionState`、`useTransition`など）の適切な使用方法について説明します。これらのHooksを正しく使用することで、パフォーマンスの向上とユーザー体験の改善を実現できます。

## useActionState vs useTransition の使い分け

### useActionState を使用する場合

- **フォーム送信**: フォームの状態管理と非同期アクションの処理を統合したい場合
- **状態の永続化**: アクションの結果を状態として保持し、UIに反映したい場合
- **エラーハンドリング**: アクションの成功/失敗状態を一元管理したい場合

### useTransition を使用する場合

- **UIの応答性**: 非同期処理中もUIの応答性を維持したい場合
- **状態更新の最適化**: 複数の状態更新を一つのトランジションとして扱いたい場合
- **ローディング状態の管理**: 処理中の状態を簡単に管理したい場合

## ベストプラクティス

### 1. useActionState の適切な使用

#### 基本的な使用方法

```typescript
import { useActionState } from 'react'
import { updateWorkspaceAction } from './server-actions'

function WorkspaceForm() {
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

  return (
    <form action={action}>
      <input name="name" required />
      <button type="submit" disabled={isPending}>
        {isPending ? 'Updating...' : 'Update'}
      </button>
      {lastResult?.message && (
        <p className={lastResult.status === 'SUCCESS' ? 'text-green-500' : 'text-red-500'}>
          {lastResult.message}
        </p>
      )}
    </form>
  )
}
```

#### withCallbacks を使用した高度な使用方法

```typescript
import { useActionState } from 'react'
import { withCallbacks } from './utils'
import { updateWorkspaceAction } from './server-actions'

function WorkspaceForm() {
  const [lastResult, action, isPending] = useActionState(
    withCallbacks(updateWorkspaceAction, {
      onStart: () => {
        // 処理開始時の処理（例：ローディング表示）
        console.log('Update started')
        return Date.now() // 参照値を返す
      },
      onEnd: (reference) => {
        // 処理終了時の処理
        console.log(`Update completed in ${Date.now() - reference}ms`)
      },
      onSuccess: (result) => {
        toast.success('Workspace updated successfully')
        setOpen(false)
        // 必要に応じてページリロードやリダイレクト
        router.refresh()
      },
      onError: (result) => {
        toast.error(`Failed to update workspace: ${result.message}`)
        // エラー時の追加処理
        console.error('Update failed:', result)
      }
    }),
    null
  )

  return (
    <form action={action}>
      {/* フォーム内容 */}
    </form>
  )
}
```

### 2. useTransition の適切な使用

#### 基本的な使用方法

```typescript
import { useState, useTransition } from 'react'
import { updateQuantityAction } from './server-actions'

function CartItem({ itemId, initialQuantity }) {
  const [quantity, setQuantity] = useState(initialQuantity)
  const [isPending, startTransition] = useTransition()

  const handleQuantityChange = (newQuantity) => {
    startTransition(async () => {
      try {
        const result = await updateQuantityAction(itemId, newQuantity)
        if (result.status === 'SUCCESS') {
          setQuantity(newQuantity)
          toast.success('Quantity updated')
        } else {
          toast.error('Failed to update quantity')
        }
      } catch (error) {
        toast.error('An error occurred')
        console.error('Update failed:', error)
      }
    })
  }

  return (
    <div>
      <input
        type="number"
        value={quantity}
        onChange={(e) => handleQuantityChange(Number(e.target.value))}
        disabled={isPending}
        min="1"
      />
      {isPending && <span>Updating...</span>}
    </div>
  )
}
```

#### 複数の状態更新を統合する場合

```typescript
import { useState, useTransition } from 'react'
import { updateUserProfileAction } from './server-actions'

function UserProfileForm({ user }) {
  const [name, setName] = useState(user.name)
  const [email, setEmail] = useState(user.email)
  const [isPending, startTransition] = useTransition()

  const handleSubmit = () => {
    startTransition(async () => {
      try {
        const result = await updateUserProfileAction({ name, email })
        if (result.status === 'SUCCESS') {
          // 複数の状態を一括更新
          setName(result.data.name)
          setEmail(result.data.email)
          toast.success('Profile updated successfully')
        } else {
          toast.error('Failed to update profile')
        }
      } catch (error) {
        toast.error('An error occurred')
      }
    })
  }

  return (
    <form onSubmit={(e) => { e.preventDefault(); handleSubmit() }}>
      <input
        value={name}
        onChange={(e) => setName(e.target.value)}
        disabled={isPending}
      />
      <input
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        disabled={isPending}
      />
      <button type="submit" disabled={isPending}>
        {isPending ? 'Updating...' : 'Update Profile'}
      </button>
    </form>
  )
}
```

### 3. パフォーマンス最適化のベストプラクティス

#### 競合状態の防止

```typescript
// ❌ 悪い例：複数のアクションが同時に実行される可能性
const handleMultipleActions = () => {
  updateAction1()
  updateAction2()
  updateAction3()
}

// ✅ 良い例：useActionState は自動的に競合を防ぐ
const [result1, action1, isPending1] = useActionState(updateAction1, null)
const [result2, action2, isPending2] = useActionState(updateAction2, null)
const [result3, action3, isPending3] = useActionState(updateAction3, null)

// または、useTransition で順次実行
const handleSequentialActions = () => {
  startTransition(async () => {
    await updateAction1()
    await updateAction2()
    await updateAction3()
  })
}
```

#### メモ化による最適化

```typescript
import { useMemo, useCallback } from 'react'

function OptimizedForm() {
  const [lastResult, action, isPending] = useActionState(updateAction, null)

  // コールバック関数をメモ化
  const callbacks = useMemo(() => ({
    onSuccess: (result) => {
      toast.success('Updated successfully')
      setOpen(false)
    },
    onError: (result) => {
      toast.error('Update failed')
    }
  }), [])

  // アクション関数をメモ化
  const memoizedAction = useCallback(
    withCallbacks(updateAction, callbacks),
    [callbacks]
  )

  return (
    <form action={memoizedAction}>
      {/* フォーム内容 */}
    </form>
  )
}
```

### 4. エラーハンドリングのベストプラクティス

#### 統一されたエラーハンドリング

```typescript
type ErrorHandler = {
  onSuccess?: (result: unknown) => void
  onError?: (error: unknown) => void
  onNetworkError?: (error: unknown) => void
}

function createErrorHandler(handlers: ErrorHandler) {
  return {
    onSuccess: (result: unknown) => {
      if (result?.status === 'SUCCESS') {
        handlers.onSuccess?.(result)
      } else {
        handlers.onError?.(result)
      }
    },
    onError: (error: unknown) => {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        handlers.onNetworkError?.(error)
      } else {
        handlers.onError?.(error)
      }
    }
  }
}

// 使用例
const [lastResult, action, isPending] = useActionState(
  withCallbacks(updateAction, createErrorHandler({
    onSuccess: (result) => {
      toast.success('Success!')
      router.refresh()
    },
    onError: (error) => {
      toast.error('Operation failed')
    },
    onNetworkError: (error) => {
      toast.error('Network error. Please check your connection.')
    }
  })),
  null
)
```

## 注意点と制限事項

### useActionState の制限

1. **フォーム専用**: `useActionState` は主にフォーム送信に最適化されている
2. **状態の永続化**: アクションの結果は状態として保持される
3. **競合防止**: 同時実行されるアクションは自動的に制御される

### useTransition の制限

1. **状態更新の最適化**: 複数の状態更新を一つのトランジションとして扱う
2. **UIの応答性**: 非同期処理中もUIの応答性を維持する
3. **手動制御**: 競合状態は手動で制御する必要がある

### パフォーマンス考慮事項

1. **不要な再レンダリング**: コールバック関数は適切にメモ化する
2. **メモリリーク**: コンポーネントのアンマウント時に適切にクリーンアップする
3. **ネットワーク最適化**: 重複するリクエストを防ぐ

## まとめ

- **useActionState**: フォーム送信と状態管理に最適
- **useTransition**: UIの応答性と状態更新の最適化に最適
- **適切なエラーハンドリング**: ユーザー体験の向上に重要
- **パフォーマンス最適化**: メモ化と競合状態の制御が重要
- **アクセシビリティ**: すべてのユーザーが利用できるUIの実装
