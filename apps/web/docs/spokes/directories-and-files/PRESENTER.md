# Presenterコンポーネント

## 概要

Presenterコンポーネントは、Container/PresentationalパターンにおけるUI表示とユーザーインタラクションを担当するClient Componentです。Container層から受け取ったデータを表示し、ユーザーの操作に応じて適切な処理を行います。

## 役割

- **UI表示**: 受け取ったデータの表示とレンダリング
- **ユーザーインタラクション**: クリック、入力、選択などの操作の処理
- **ローカル状態管理**: コンポーネント内での状態管理
- **イベントハンドリング**: ユーザーアクションに対する適切な処理

## ファイル名規則

### 基本パターン
```
{機能名}Presenter.tsx
```

### 具体例
```
UserListPresenter.tsx
ProductDetailPresenter.tsx
OrderHistoryPresenter.tsx
```

## 実装ルール

1. **Client Componentとして実装**: `'use client'`ディレクティブを使用
2. **Props経由のデータ受取**: Container層からデータを受け取る
3. **カスタムフックの活用**: ロジックは`use-{機能名}.ts`に分離
4. **components/の活用**: `@/components/`のコンポーネントを利用
5. **型安全性**: TypeScriptの型定義を活用
6. **単一責任**: UI表示とユーザーインタラクションのみに集中

## 実装例

```tsx
'use client';

import { useState } from 'react';
import { OrionButton } from '@/components/elements/orion-button';
import { OrionList } from '@/components/elements/orion-list';
import { OrionModal } from '@/components/elements/orion-modal';
import { useUserList } from '@/hooks/use-user-list';
import type { User } from '@/types/user';

type UserListPresenterProps = {
  users: User[];
  onUserSelect?: (user: User) => void;
};

export function UserListPresenter({ users, onUserSelect }: UserListPresenterProps) {
  const {
    selectedUser,
    setSelectedUser,
    filter,
    setFilter,
    filteredUsers,
    isModalOpen,
    setIsModalOpen
  } = useUserList(users);

  const handleUserClick = (user: User) => {
    setSelectedUser(user);
    setIsModalOpen(true);
    onUserSelect?.(user);
  };

  return (
    <div className="user-list-container">
      <div className="search-section">
        <input
          type="text"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          placeholder="ユーザーを検索..."
          className="search-input"
        />
      </div>

      <OrionList
        data={filteredUsers}
        renderItem={(user) => (
          <div
            key={user.id}
            className="user-item"
            onClick={() => handleUserClick(user)}
          >
            <span>{user.name}</span>
            <span>{user.email}</span>
          </div>
        )}
      />

      <OrionModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="ユーザー詳細"
      >
        {selectedUser && (
          <div className="user-detail">
            <p>名前: {selectedUser.name}</p>
            <p>メール: {selectedUser.email}</p>
            <OrionButton onClick={() => setIsModalOpen(false)}>
              閉じる
            </OrionButton>
          </div>
        )}
      </OrionModal>
    </div>
  );
}
```

## カスタムフック分離

### 基本パターン
```tsx
// use-{機能名}.ts
export function useUserList(users: User[]) {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [filter, setFilter] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const filteredUsers = useMemo(() =>
    users.filter(user =>
      user.name.toLowerCase().includes(filter.toLowerCase()) ||
      user.email.toLowerCase().includes(filter.toLowerCase())
    ),
    [users, filter]
  );

  return {
    selectedUser,
    setSelectedUser,
    filter,
    setFilter,
    filteredUsers,
    isModalOpen,
    setIsModalOpen
  };
}
```

## components/ディレクトリの活用

### 利用可能なコンポーネント

#### ui/ - 原子的なUIコンポーネント
```tsx
import { ScaledIcon } from '@/components/ui/scaled-icon';
```

#### elements/ - 分子的なコンポーネント
```tsx
import { OrionButton } from '@/components/elements/orion-button';
import { OrionInput } from '@/components/elements/orion-input';
import { OrionList } from '@/components/elements/orion-list';
import { OrionModal } from '@/components/elements/orion-modal';
import { OrionSelect } from '@/components/elements/orion-select';
```

#### composites/ - 複合的なコンポーネント
```tsx
import { OrionFloatDrawer } from '@/components/composites/orion-float-drawer';
```

#### layouts/ - レイアウトコンポーネント
```tsx
import { GlobalNavigation } from '@/components/layouts/global-navigation';
import { Header } from '@/components/layouts/header';
import { Sidebar } from '@/components/layouts/sidebar';
```

## 実装パターン

### 1. 基本的な表示コンポーネント
```tsx
'use client';

import { OrionList } from '@/components/elements/orion-list';

type DataListPresenterProps = {
  data: DataType[];
  onItemClick?: (item: DataType) => void;
};

export function DataListPresenter({ data, onItemClick }: DataListPresenterProps) {
  return (
    <OrionList
      data={data}
      renderItem={(item) => (
        <div onClick={() => onItemClick?.(item)}>
          {item.name}
        </div>
      )}
    />
  );
}
```

### 2. フォームコンポーネント
```tsx
'use client';

import { useState } from 'react';
import { OrionInput } from '@/components/elements/orion-input';
import { OrionButton } from '@/components/elements/orion-button';
import { useFormValidation } from '@/hooks/use-form-validation';

type FormPresenterProps = {
  onSubmit: (data: FormData) => void;
  initialData?: Partial<FormData>;
};

export function FormPresenter({ onSubmit, initialData }: FormPresenterProps) {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    email: '',
    ...initialData
  });

  const { errors, validate, isValid } = useFormValidation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isValid(formData)) {
      onSubmit(formData);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <OrionInput
        value={formData.name}
        onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
        error={errors.name}
        placeholder="名前"
      />
      <OrionInput
        value={formData.email}
        onChange={(e) => setFormData(prev => ({ ...prev, email: e.target.value }))}
        error={errors.email}
        placeholder="メールアドレス"
      />
      <OrionButton type="submit" disabled={!isValid(formData)}>
        送信
      </OrionButton>
    </form>
  );
}
```

## 注意点

- **過度な分離を避ける**: 単純なコンポーネントで無理に分離しない
- **Props drillingの回避**: 深い階層でのProps渡しは避ける
- **カスタムフックの適切な分離**: ロジックは`use-{機能名}.ts`に分離
- **components/の適切な利用**: 既存のコンポーネントを最大限活用
- **型安全性の確保**: 適切な型定義を行う

## 関連ファイル

- Container層: `{機能名}Container.tsx`
- カスタムフック: `use-{機能名}.ts`
- 型定義: `@/types/{機能名}.ts`
- 利用コンポーネント: `@/components/`内の各コンポーネント
