# Containerコンポーネント

## 概要

Containerコンポーネントは、Container/Presentationalパターンにおけるデータ取得とビジネスロジックを担当するServer Componentです。UI表示は行わず、Presentational層にデータを渡す役割を担います。

## 役割

- **データ取得**: サーバーサイドでのデータベースアクセス、API呼び出し
- **ビジネスロジック**: データの変換、バリデーション、計算処理
- **状態管理**: サーバーサイドでの状態管理
- **Presentational層へのデータ渡し**: Props経由でのデータ提供

## ファイル名規則

### 基本パターン
```
{機能名}Container.tsx
```

### 具体例
```
UserListContainer.tsx
ProductDetailContainer.tsx
OrderHistoryContainer.tsx
```

## 実装ルール

1. **Server Componentとして実装**: `'use client'`ディレクティブを使用しない
2. **async関数として定義**: データ取得のため非同期処理が必要
3. **データ取得のみ**: UI表示ロジックは含めない
4. **Presentational層へのProps渡し**: 取得したデータをPropsで渡す
5. **エラーハンドリング**: データ取得エラーの適切な処理
6. **型安全性**: TypeScriptの型定義を活用

## 実装例

```tsx
// UserListContainer.tsx
import { UserListPresenter } from './UserListPresenter';
import { fetchUsers } from '@/lib/api/users';
import type { User } from '@/types/user';

export async function UserListContainer() {
  try {
    const users = await fetchUsers();
    return <UserListPresenter users={users} />;
  } catch (error) {
    return <div>ユーザー情報の取得に失敗しました</div>;
  }
}
```

## 注意点

- **過度な分離を避ける**: 単純なコンポーネントで無理に分離しない
- **Props drillingの回避**: 深い階層でのProps渡しは避ける
- **適切な境界設定**: Server/Client境界を適切に設定する
- **パフォーマンス考慮**: データ取得の最適化を心がける

## 関連ファイル

- Presentational層: `{機能名}Presenter.tsx`
- カスタムフック: `use{機能名}.ts`
- 型定義: `@/types/{機能名}.ts`
