# TypeScript ベストプラクティス

このドキュメントでは、プロジェクトでTypeScriptを使用する際のベストプラクティスを定義します。

## 基本原則

### 1. anyの使用を禁止する

`any`型は型安全性を完全に破壊するため、使用を禁止します。

```typescript
// ❌ 悪い例
function processData(data: any) {
  return data.someProperty;
}

// ✅ 良い例
function processData(data: { someProperty: string }) {
  return data.someProperty;
}

// または、より具体的な型を定義
type DataType = {
  someProperty: string;
  otherProperty?: number;
};

function processData(data: DataType) {
  return data.someProperty;
}
```

### 2. 型キャストを可能な限り避ける

型キャスト（`as`）は型安全性を損なうため、可能な限り避けます。

```typescript
// ❌ 悪い例
const element = document.getElementById('myElement') as HTMLInputElement;

// ✅ 良い例
const element = document.getElementById('myElement');
if (element instanceof HTMLInputElement) {
  // ここでelementはHTMLInputElement型として安全に使用できる
  element.value = 'test';
}

// または、組み込みの型チェックを使用
const element = document.getElementById('myElement');
if (element instanceof HTMLInputElement) {
  element.value = 'test';
}
```

### 3. unknown型はブロックスコープ内での利用にとどめる

`unknown`型は型安全な`any`として使用できますが、ブロックスコープ内での利用に限定します。

```typescript
// ✅ 良い例
function processUnknownData(data: unknown) {
  // ブロックスコープ内で型チェックを行う
  if (typeof data === 'string') {
    return data.toUpperCase();
  }

  if (typeof data === 'number') {
    return data * 2;
  }

  if (data && typeof data === 'object' && 'message' in data) {
    return (data as { message: string }).message;
  }

  return 'Unknown data type';
}

// ❌ 悪い例 - unknown型を関数の外に漏らす
let globalUnknown: unknown;
function badExample(data: unknown) {
  globalUnknown = data; // 型安全性を損なう
}
```

### 4. ユーザー型ガードの利用を避ける

ユーザー定義の型ガードは複雑になりがちで、型安全性を損なう可能性があるため、可能な限り避けます。

```typescript
// ❌ 悪い例 - 複雑なユーザー型ガード
// 型ガード関数は複雑になりがちで、型安全性を損なう可能性があります

// ✅ 良い例 - 組み込みの型チェックを使用
function processUserData(data: unknown) {
  if (data && typeof data === 'object' && 'id' in data && 'name' in data) {
    // ここでdataは適切に型チェックされている
    const userData = data as { id: string; name: string };
    return userData.name;
  }
  throw new Error('Invalid user data');
}

// または、バリデーションライブラリを使用
import { z } from 'zod';

const UserSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string().email(),
});

function processUserData(data: unknown) {
  const user = UserSchema.parse(data); // 型安全なバリデーション
  return user.name; // userは自動的にUser型になる
}
```

### 5. ジェネリクスの利用を最小限に抑える

複雑なジェネリクスは可読性を損なうため、可能な限り避けます。一般的なジェネリクス（`Record<T>`、`Array<T>`など）は使用可能です。

```typescript
// ✅ 良い例 - 一般的なジェネリクス
type UserRecord = Record<string, User>;
type UserArray = Array<User>;
type UserMap = Map<string, User>;

// ✅ 良い例 - シンプルなジェネリクス
function createArray<T>(item: T): T[] {
  return [item];
}

// ❌ 悪い例 - 複雑なジェネリクス
type ComplexGeneric<T extends Record<string, any>, K extends keyof T, V extends T[K]> = {
  [P in K]: V extends T[P] ? T[P] : never;
} & {
  metadata: {
    keys: K[];
    values: V[];
  };
};

// ✅ 良い例 - 複雑なジェネリクスを避けて、具体的な型を定義
type UserWithMetadata = {
  id: string;
  name: string;
  email: string;
  metadata: {
    keys: string[];
    values: string[];
  };
};
```

## 型定義のベストプラクティス

### 1. typeの使用を優先する

`interface`はライブラリの型を拡張する場合を除き、可能な限り避けて`type`を使用します。

```typescript
// ✅ 良い例 - typeを使用
type User = {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
};

type CreateUserRequest = {
  name: string;
  email: string;
};

type UpdateUserRequest = {
  name?: string;
  email?: string;
};

// ✅ 良い例 - ライブラリの型拡張の場合のみinterfaceを使用
declare global {
  namespace Express {
    interface Request {
      user?: User;
    }
  }
}

// ❌ 悪い例 - 通常の型定義でinterfaceを使用
interface User {
  id: string;
  name: string;
}
```

### 2. 型エイリアスの適切な使用

```typescript
// ✅ 良い例 - プリミティブ型の組み合わせ
type UserId = string;
type UserEmail = string;
type UserStatus = 'active' | 'inactive' | 'pending';

// ✅ 良い例 - 複雑な型の組み合わせ
type UserWithStatus = User & {
  status: UserStatus;
};
```

### 3. ユニオン型の活用

```typescript
// ✅ 良い例
type Theme = 'light' | 'dark' | 'auto';
type LoadingState = 'idle' | 'loading' | 'success' | 'error';

// ✅ 良い例 - 判別可能なユニオン型
type ApiResponse =
  | { status: 'success'; data: User[] }
  | { status: 'error'; message: string }
  | { status: 'loading' };
```

## 関数とメソッドの型定義

### 1. 関数の型定義

```typescript
// ✅ 良い例
type EventHandler = (event: Event) => void;
type AsyncHandler<T> = (data: T) => Promise<void>;

// 関数の型定義
function processUser(user: User): Promise<User> {
  // 実装
}

// アロー関数の型定義
const processUser = (user: User): Promise<User> => {
  // 実装
};
```

### 2. オプショナルパラメータの適切な使用

```typescript
// ✅ 良い例
function createUser(name: string, email: string, options?: {
  isActive?: boolean;
  role?: string;
}) {
  // 実装
}

// ✅ 良い例 - デフォルト値の使用
function createUser(
  name: string,
  email: string,
  isActive: boolean = true
) {
  // 実装
}
```

## エラーハンドリング

### 1. 型安全なエラーハンドリング

```typescript
// ✅ 良い例
type Result<T, E = Error> =
  | { success: true; data: T }
  | { success: false; error: E };

function fetchUser(id: string): Promise<Result<User, string>> {
  return fetch(`/api/users/${id}`)
    .then(response => {
      if (!response.ok) {
        return { success: false, error: 'User not found' };
      }
      return response.json();
    })
    .then(data => ({ success: true, data }))
    .catch(error => ({ success: false, error: error.message }));
}
```

## モジュールとインポート

### 1. 型のみのインポート

```typescript
// ✅ 良い例
import type { User, CreateUserRequest } from './types';
import { createUser } from './api';

// ❌ 悪い例
import { User, createUser } from './types'; // 型と実装を混在
```

### 2. 名前空間の適切な使用

```typescript
// ✅ 良い例 - 名前空間の使用（typeを使用）
namespace Api {
  export type User = {
    id: string;
    name: string;
  };

  export function fetchUser(id: string): Promise<User> {
    // 実装
  }
}

// 使用時
const user: Api.User = await Api.fetchUser('123');
```

## 設定とツール

### 1. tsconfig.jsonの推奨設定

```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true
  }
}
```
