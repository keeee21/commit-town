# React コンポーネントパターン

このドキュメントでは、Webアプリケーションで使用するReactコンポーネントの設計パターンについて説明します。

## コンポーネントの分類

### 1. Presentational Components（表示コンポーネント）

```typescript
// 純粋な表示コンポーネント
type UserCardProps = {
  user: User;
  onEdit?: (user: User) => void;
  onDelete?: (userId: string) => void;
}

export function UserCard({ user, onEdit, onDelete }: UserCardProps) {
  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <div className="actions">
        {onEdit && (
          <button onClick={() => onEdit(user)}>編集</button>
        )}
        {onDelete && (
          <button onClick={() => onDelete(user.id)}>削除</button>
        )}
      </div>
    </div>
  );
}
```

### 2. Container Components（コンテナコンポーネント）

```typescript
// ロジックを含むコンテナコンポーネント
export function UserListContainer() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await api.getUsers();
      setUsers(data);
    } catch (err) {
      setError('ユーザーの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (user: User) => {
    // 編集ロジック
  };

  const handleDelete = async (userId: string) => {
    // 削除ロジック
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  if (loading) return <div>読み込み中...</div>;
  if (error) return <div>エラー: {error}</div>;

  return (
    <div>
      {users.map(user => (
        <UserCard
          key={user.id}
          user={user}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ))}
    </div>
  );
}
```

## カスタムフックパターン

### 1. データフェッチングフック

```typescript
// useUsers.ts
export function useUsers() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await api.getUsers();
      setUsers(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  }, []);

  const addUser = useCallback(async (userData: CreateUserData) => {
    try {
      const newUser = await api.createUser(userData);
      setUsers(prev => [...prev, newUser]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ユーザーの作成に失敗しました');
    }
  }, []);

  const updateUser = useCallback(async (id: string, userData: UpdateUserData) => {
    try {
      const updatedUser = await api.updateUser(id, userData);
      setUsers(prev => prev.map(user =>
        user.id === id ? updatedUser : user
      ));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ユーザーの更新に失敗しました');
    }
  }, []);

  const deleteUser = useCallback(async (id: string) => {
    try {
      await api.deleteUser(id);
      setUsers(prev => prev.filter(user => user.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'ユーザーの削除に失敗しました');
    }
  }, []);

  return {
    users,
    loading,
    error,
    fetchUsers,
    addUser,
    updateUser,
    deleteUser,
  };
}
```

### 2. フォーム管理フック

```typescript
// useForm.ts
export function useForm<T extends Record<string, any>>(
  initialValues: T,
  validationSchema?: z.ZodSchema<T>
) {
  const [values, setValues] = useState<T>(initialValues);
  const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});

  const setValue = useCallback((field: keyof T, value: any) => {
    setValues(prev => ({ ...prev, [field]: value }));

    // バリデーション
    if (validationSchema) {
      try {
        validationSchema.parse({ ...values, [field]: value });
        setErrors(prev => ({ ...prev, [field]: undefined }));
      } catch (err) {
        if (err instanceof z.ZodError) {
          const fieldError = err.errors.find(e => e.path[0] === field);
          setErrors(prev => ({
            ...prev,
            [field]: fieldError?.message
          }));
        }
      }
    }
  }, [values, validationSchema]);

  const setFieldTouched = useCallback((field: keyof T) => {
    setTouched(prev => ({ ...prev, [field]: true }));
  }, []);

  const reset = useCallback(() => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
  }, [initialValues]);

  const isValid = Object.keys(errors).length === 0;

  return {
    values,
    errors,
    touched,
    setValue,
    setFieldTouched,
    reset,
    isValid,
  };
}
```

## レンダープロップパターン

```typescript
// Render Props パターン
type DataFetcherProps<T> = {
  url: string;
  children: (data: {
    data: T | null;
    loading: boolean;
    error: string | null;
    refetch: () => void;
  }) => React.ReactNode;
}

export function DataFetcher<T>({ url, children }: DataFetcherProps<T>) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(url);
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  }, [url]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return <>{children({ data, loading, error, refetch: fetchData })}</>;
}

// 使用例
export function UserList() {
  return (
    <DataFetcher<User[]> url="/api/users">
      {({ data, loading, error, refetch }) => {
        if (loading) return <div>読み込み中...</div>;
        if (error) return <div>エラー: {error}</div>;
        if (!data) return <div>データがありません</div>;

        return (
          <div>
            <button onClick={refetch}>再読み込み</button>
            {data.map(user => (
              <div key={user.id}>{user.name}</div>
            ))}
          </div>
        );
      }}
    </DataFetcher>
  );
}
```

## 高階コンポーネント（HOC）パターン

```typescript
// HOC パターン
type WithLoadingProps = {
  loading: boolean;
};

export function withLoading<P extends object>(
  Component: React.ComponentType<P>
) {
  return function WithLoadingComponent(props: P & WithLoadingProps) {
    const { loading, ...rest } = props;

    if (loading) {
      return <div>読み込み中...</div>;
    }

    return <Component {...(rest as P)} />;
  };
}

// 使用例
const UserCardWithLoading = withLoading(UserCard);

export function UserListWithLoading() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUsers().then(data => {
      setUsers(data);
      setLoading(false);
    });
  }, []);

  return (
    <div>
      {users.map(user => (
        <UserCardWithLoading
          key={user.id}
          user={user}
          loading={loading}
        />
      ))}
    </div>
  );
}
```

## コンテキストパターン

```typescript
// コンテキストの作成
type UserContextType = {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);

  const login = useCallback((userData: User) => {
    setUser(userData);
  }, []);

  const logout = useCallback(() => {
    setUser(null);
  }, []);

  const isAuthenticated = user !== null;

  const value = {
    user,
    login,
    logout,
    isAuthenticated,
  };

  return (
    <UserContext.Provider value={value}>
      {children}
    </UserContext.Provider>
  );
}

// カスタムフック
export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
}

// 使用例
export function UserProfile() {
  const { user, logout } = useUser();

  if (!user) {
    return <div>ログインしてください</div>;
  }

  return (
    <div>
      <h1>{user.name}</h1>
      <button onClick={logout}>ログアウト</button>
    </div>
  );
}
```

## コンポーネントの最適化

### 1. React.memo

```typescript
// メモ化されたコンポーネント
export const UserCard = React.memo(function UserCard({
  user,
  onEdit,
  onDelete
}: UserCardProps) {
  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <div className="actions">
        {onEdit && (
          <button onClick={() => onEdit(user)}>編集</button>
        )}
        {onDelete && (
          <button onClick={() => onDelete(user.id)}>削除</button>
        )}
      </div>
    </div>
  );
});
```

### 2. useMemo と useCallback

```typescript
export function UserList({ users, searchTerm }: UserListProps) {
  // フィルタリング結果をメモ化
  const filteredUsers = useMemo(() => {
    return users.filter(user =>
      user.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [users, searchTerm]);

  // イベントハンドラーをメモ化
  const handleEdit = useCallback((user: User) => {
    // 編集ロジック
  }, []);

  const handleDelete = useCallback((userId: string) => {
    // 削除ロジック
  }, []);

  return (
    <div>
      {filteredUsers.map(user => (
        <UserCard
          key={user.id}
          user={user}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ))}
    </div>
  );
}
```

## 参考資料

- [React Patterns](https://reactpatterns.com/)
- [React Component Patterns](https://reactjs.org/docs/thinking-in-react.html)
- [Custom Hooks](https://reactjs.org/docs/hooks-custom.html)
