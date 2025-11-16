# Container/Presentationalパターン

## 概要

Container/Presentationalパターンは、Next.js App RouterにおいてServer ComponentとClient Componentの役割を明確に分離する設計パターンです。データ取得とビジネスロジックをContainer層で、UI表示とユーザーインタラクションをPresentational層で担当します。

## アーキテクチャ

### Container層（Server Component）
- **役割**: データ取得、ビジネスロジック、状態管理
- **特徴**: サーバーサイドで実行、データベースアクセス、API呼び出し
- **責任**:
  - データの取得とキャッシュコントロール
  - ビジネスロジックの実行
  - Presentational層へのデータ渡し

### Presentational層（Client Component）
- **役割**: UI表示、ユーザーインタラクション、状態管理
- **特徴**: クライアントサイドで実行、ブラウザAPI利用可能
- **責任**:
  - データの表示とレンダリング
  - ユーザーインタラクションの処理
  - ローカル状態の管理

## 実装パターン

### 基本的な構造
```tsx
// Container層（Server Component）
async function UserListContainer() {
  const users = await fetchUsers();
  return <UserListPresenter users={users} />;
}

// Presentational層（Client Component）
function UserListPresenter({ users }: { users: User[] }) {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  return (
    <div>
      {users.map(user => (
        <UserCard
          key={user.id}
          user={user}
          onSelect={setSelectedUser}
        />
      ))}
    </div>
  );
}
```

### Hooks分離パターン
```tsx
// useUserList.ts
export function useUserList(users: User[]) {
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [filter, setFilter] = useState('');

  const filteredUsers = useMemo(() =>
    users.filter(user => user.name.includes(filter)),
    [users, filter]
  );

  return {
    selectedUser,
    setSelectedUser,
    filter,
    setFilter,
    filteredUsers
  };
}

// UserListPresenter.tsx
function UserListPresenter({ users }: { users: User[] }) {
  const {
    selectedUser,
    setSelectedUser,
    filter,
    setFilter,
    filteredUsers
  } = useUserList(users);

  return (
    <div>
      <input
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
      />
      {filteredUsers.map(user => (
        <UserCard
          key={user.id}
          user={user}
          onSelect={setSelectedUser}
        />
      ))}
    </div>
  );
}
```

## ルール

1. **明確な役割分離**: Container層はデータ取得、Presentational層はUI表示
2. **Props経由のデータ渡し**: Container層からPresentational層へのデータはPropsで渡す
3. **Hooks分離**: Presentational層のロジックは`use-xxx.ts`に分離
4. **UI Component分離**: 純粋なUI表示は別のComponentに分離

## 例外

- **Client Componentからのデータフェッチ**: ユーザーインタラクションをトリガーとする場合
  - SWRやReact Queryを使用してNext.js API Routeを呼び出し
  - リアルタイム性が重要な場合

## メリット

- **パフォーマンス向上**: Server Componentでのデータ取得により初期表示が高速
- **保守性向上**: 関心の分離によりコードの理解と修正が容易
- **テスタビリティ向上**: 各層を独立してテスト可能
- **再利用性向上**: Presentational層の再利用が容易

## 注意点

- **過度な分離**: 単純なコンポーネントで無理に分離しない
- **Props drilling**: 深い階層でのProps渡しは避ける
- **Server/Client境界**: 適切な境界設定が重要

## 詳細な実装ガイド

### Containerコンポーネント
詳細な実装ルールとファイル名規則については、[@CONTAINER.md](../directories-and-files/CONTAINER.md)を参照してください。

### Presenterコンポーネント
詳細な実装ルールとファイル名規則については、[@PRESENTER.md](../directories-and-files/PRESENTER.md)を参照してください。

## 参考

- [Next.jsの考え方 - Container/Presentationalパターン](https://zenn.dev/akfm/books/nextjs-basic-principle/viewer/part_2_container_presentational_pattern)
