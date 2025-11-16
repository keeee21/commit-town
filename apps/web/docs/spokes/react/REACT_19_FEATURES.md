# React 19 新機能とベストプラクティス

このドキュメントでは、React 19の新機能とWebアプリケーションでの活用方法について説明します。

## React 19の主要新機能

### 1. カスタム要素の完全サポート

```typescript
// カスタム要素の使用例
export function CustomElementExample() {
  const [value, setValue] = useState('');

  return (
    <div>
      {/* カスタム要素の直接使用が可能 */}
      <my-custom-input
        value={value}
        onchange={(e) => setValue(e.target.value)}
      />
    </div>
  );
}
```

### 2. 新しいフック: use()

```typescript
import { use } from 'react';

// Promise の使用
export function UserProfile({ userId }: { userId: string }) {
  const userPromise = fetchUser(userId);
  const user = use(userPromise); // Promise を直接使用

  return <div>{user.name}</div>;
}

// Context の使用
const ThemeContext = createContext('light');

export function ThemedComponent() {
  const theme = use(ThemeContext); // Context を直接使用
  return <div className={theme}>Themed content</div>;
}
```

### 3. アクションとフォームの改善

```typescript
// Server Actions との統合
export function ContactForm() {
  return (
    <form action={submitContact}>
      <input name="name" required />
      <input name="email" type="email" required />
      <button type="submit">送信</button>
    </form>
  );
}

async function submitContact(formData: FormData) {
  'use server';

  const name = formData.get('name') as string;
  const email = formData.get('email') as string;

  // サーバーサイドでの処理
  await saveContact({ name, email });
}
```

### 4. ref の改善

```typescript
// ref の自動転送
export function MyInput({ placeholder }: { placeholder: string }) {
  return <input placeholder={placeholder} />; // ref が自動で転送される
}

// 使用例
export function ParentComponent() {
  const inputRef = useRef<HTMLInputElement>(null);

  const focusInput = () => {
    inputRef.current?.focus();
  };

  return (
    <div>
      <MyInput ref={inputRef} placeholder="入力してください" />
      <button onClick={focusInput}>フォーカス</button>
    </div>
  );
}
```

## パフォーマンス改善

### 1. 自動バッチングの拡張

```typescript
// 非同期処理でも自動バッチング
export function AsyncBatchingExample() {
  const [count, setCount] = useState(0);
  const [flag, setFlag] = useState(false);

  const handleClick = async () => {
    // これらは自動的にバッチングされる
    setCount(c => c + 1);
    setFlag(f => !f);

    // 非同期処理後もバッチングされる
    await new Promise(resolve => setTimeout(resolve, 100));
    setCount(c => c + 1);
    setFlag(f => !f);
  };

  return <button onClick={handleClick}>更新</button>;
}
```

### 2. 並行レンダリングの改善

```typescript
import { Suspense } from 'react';

export function ConcurrentRendering() {
  return (
    <div>
      <Suspense fallback={<div>読み込み中...</div>}>
        <HeavyComponent />
      </Suspense>
      <Suspense fallback={<div>読み込み中...</div>}>
        <AnotherHeavyComponent />
      </Suspense>
    </div>
  );
}
```

## 型安全性の向上

### 1. より厳密な型チェック

```typescript
// イベントハンドラーの型安全性
export function TypedEventHandler() {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // e.target.value の型が string として推論される
    console.log(e.target.value);
  };

  return <input onChange={handleChange} />;
}
```

### 2. コンポーネントの型定義

```typescript
// より厳密なコンポーネント型定義
type ButtonProps = {
  children: React.ReactNode;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
  disabled?: boolean;
};

export function Button({
  children,
  onClick,
  variant = 'primary',
  disabled = false
}: ButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`btn btn-${variant}`}
    >
      {children}
    </button>
  );
}
```

## 開発者体験の改善

### 1. より良いエラーメッセージ

```typescript
// React 19ではより詳細なエラーメッセージが表示される
export function ErrorExample() {
  const [data, setData] = useState(null);

  // より明確なエラーメッセージが表示される
  return <div>{data.name}</div>; // data が null の場合のエラー
}
```

### 2. 開発者ツールの改善

```typescript
// React DevTools でのデバッグが改善
export function DebugExample() {
  const [state, setState] = useState({ count: 0, name: '' });

  // DevTools でより詳細な情報が表示される
  return (
    <div>
      <p>Count: {state.count}</p>
      <p>Name: {state.name}</p>
    </div>
  );
}
```

## 実践的な使用例

### 1. フォーム管理の改善

```typescript
import { useFormStatus } from 'react-dom';

export function SubmitButton() {
  const { pending } = useFormStatus();

  return (
    <button type="submit" disabled={pending}>
      {pending ? '送信中...' : '送信'}
    </button>
  );
}

export function ContactForm() {
  return (
    <form action={submitContact}>
      <input name="name" required />
      <input name="email" type="email" required />
      <SubmitButton />
    </form>
  );
}
```

### 2. データフェッチングの改善

```typescript
import { use } from 'react';

export function UserProfile({ userPromise }: { userPromise: Promise<User> }) {
  const user = use(userPromise);

  return (
    <div>
      <h1>{user.name}</h1>
      <p>{user.email}</p>
    </div>
  );
}

// 使用例
export function UserPage({ userId }: { userId: string }) {
  const userPromise = fetchUser(userId);

  return (
    <Suspense fallback={<div>読み込み中...</div>}>
      <UserProfile userPromise={userPromise} />
    </Suspense>
  );
}
```

## 移行ガイド

### 1. 既存コードの更新

```typescript
// 古い書き方
export function OldComponent() {
  const [data, setData] = useState(null);

  useEffect(() => {
    fetchData().then(setData);
  }, []);

  if (!data) return <div>読み込み中...</div>;

  return <div>{data.name}</div>;
}

// 新しい書き方（React 19）
export function NewComponent() {
  const dataPromise = fetchData();

  return (
    <Suspense fallback={<div>読み込み中...</div>}>
      <DataDisplay dataPromise={dataPromise} />
    </Suspense>
  );
}

function DataDisplay({ dataPromise }: { dataPromise: Promise<Data> }) {
  const data = use(dataPromise);
  return <div>{data.name}</div>;
}
```

### 2. 型定義の更新

```typescript
// より厳密な型定義
type ComponentProps = {
  children: React.ReactNode;
  className?: string;
  onClick?: (event: React.MouseEvent<HTMLDivElement>) => void;
}

export function Component({ children, className, onClick }: ComponentProps) {
  return (
    <div className={className} onClick={onClick}>
      {children}
    </div>
  );
}
```

## 参考資料

- [React 19 Release Notes](https://react.dev/blog/2024/12/05/react-19)
- [React 19 Upgrade Guide](https://react.dev/blog/2024/12/05/react-19-upgrade-guide)
- [React 19 New Features](https://react.dev/learn)
