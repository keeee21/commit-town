# TypeScript 高度な型定義

このドキュメントでは、Webアプリケーションで使用する高度なTypeScriptの型定義パターンについて説明します。

## ユーティリティ型の活用

### 1. 基本的なユーティリティ型

```typescript
// 既存の型から新しい型を作成
type User = {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
};

// 部分的な型
type PartialUser = Partial<User>;
// { id?: string; name?: string; email?: string; ... }

// 必須フィールドのみ
type RequiredUser = Required<Pick<User, 'id' | 'name'>>;
// { id: string; name: string }

// 特定のフィールドを除外
type UserWithoutTimestamps = Omit<User, 'createdAt' | 'updatedAt'>;
// { id: string; name: string; email: string }

// 特定のフィールドのみ
type UserBasicInfo = Pick<User, 'id' | 'name' | 'email'>;
// { id: string; name: string; email: string }
```

### 2. 条件型（Conditional Types）

```typescript
// 条件に基づく型選択
type ApiResponse<T> = T extends string
  ? { message: T }
  : { data: T };

type StringResponse = ApiResponse<string>; // { message: string }
type DataResponse = ApiResponse<User>; // { data: User }

// 配列かどうかの判定
type IsArray<T> = T extends any[] ? true : false;

type Test1 = IsArray<string[]>; // true
type Test2 = IsArray<string>; // false
```

### 3. テンプレートリテラル型

```typescript
// 文字列テンプレート型
type EventName<T extends string> = `on${Capitalize<T>}`;

type ClickEvent = EventName<'click'>; // 'onClick'
type ChangeEvent = EventName<'change'>; // 'onChange'

// API エンドポイントの型
type ApiEndpoint<T extends string> = `/api/${T}`;

type UserEndpoint = ApiEndpoint<'users'>; // '/api/users'
type TalentEndpoint = ApiEndpoint<'talents'>; // '/api/talents'
```

## ジェネリクスの高度な使用

### 1. 制約付きジェネリクス

```typescript
// オブジェクトのキーを制約
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

const user = { id: '1', name: 'John', age: 30 };
const name = getProperty(user, 'name'); // string
const id = getProperty(user, 'id'); // string

// 配列の要素型を制約
function getFirstElement<T extends readonly unknown[]>(arr: T): T[0] | undefined {
  return arr[0];
}

const numbers = [1, 2, 3] as const;
const first = getFirstElement(numbers); // 1
```

### 2. 再帰的な型定義

```typescript
// ネストしたオブジェクトの型
type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

type NestedUser = {
  id: string;
  profile: {
    name: string;
    avatar: {
      url: string;
      alt: string;
    };
  };
};

type PartialNestedUser = DeepPartial<NestedUser>;
// {
//   id?: string;
//   profile?: {
//     name?: string;
//     avatar?: {
//       url?: string;
//       alt?: string;
//     };
//   };
// }
```

## 型ガードと型述語

### 1. 型ガード関数

```typescript
// 型ガード関数
function isString(value: unknown): value is string {
  return typeof value === 'string';
}

function isUser(value: unknown): value is User {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'name' in value &&
    'email' in value
  );
}

// 使用例
function processValue(value: unknown) {
  if (isString(value)) {
    // value は string として扱われる
    console.log(value.toUpperCase());
  }

  if (isUser(value)) {
    // value は User として扱われる
    console.log(value.name);
  }
}
```

### 2. 判別共用体（Discriminated Union）

```typescript
// 判別共用体の定義
type ApiResponse<T> =
  | { status: 'success'; data: T }
  | { status: 'error'; error: string }
  | { status: 'loading' };

// 型ガード関数
function isSuccessResponse<T>(response: ApiResponse<T>): response is { status: 'success'; data: T } {
  return response.status === 'success';
}

// 使用例
function handleResponse<T>(response: ApiResponse<T>) {
  if (isSuccessResponse(response)) {
    // response.data が利用可能
    console.log(response.data);
  } else if (response.status === 'error') {
    // response.error が利用可能
    console.error(response.error);
  } else {
    // response.status === 'loading'
    console.log('読み込み中...');
  }
}
```

## モジュール拡張

### 1. グローバル型の拡張

```typescript
// グローバル型の拡張
declare global {
  interface Window {
    gtag: (command: string, targetId: string, config?: object) => void;
  }
}

// 使用例
function trackEvent(eventName: string) {
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', eventName);
  }
}
```

### 2. モジュール宣言の拡張

```typescript
// CSS モジュールの型定義
declare module '*.module.css' {
  const classes: { [key: string]: string };
  export default classes;
}

// 画像ファイルの型定義
declare module '*.png' {
  const src: string;
  export default src;
}

declare module '*.jpg' {
  const src: string;
  export default src;
}

declare module '*.svg' {
  const src: string;
  export default src;
}
```

## 高度な型操作

### 1. 型の変換

```typescript
// 文字列リテラル型の変換
type UppercaseKeys<T> = {
  [K in keyof T as Uppercase<string & K>]: T[K];
};

type UserWithUppercaseKeys = UppercaseKeys<User>;
// { ID: string; NAME: string; EMAIL: string; ... }

// オプショナルなキーの作成
type OptionalKeys<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;

type UserWithOptionalEmail = OptionalKeys<User, 'email'>;
// { id: string; name: string; email?: string; ... }
```

### 2. 型の検証

```typescript
// 型の等価性チェック
type Equals<X, Y> = (<T>() => T extends X ? 1 : 2) extends <T>() => T extends Y ? 1 : 2 ? true : false;

type Test1 = Equals<string, string>; // true
type Test2 = Equals<string, number>; // false

// 型の存在チェック
type HasProperty<T, K extends PropertyKey> = K extends keyof T ? true : false;

type HasId = HasProperty<User, 'id'>; // true
type HasPhone = HasProperty<User, 'phone'>; // false
```

## 実践的な型定義例

### 1. API レスポンスの型定義

```typescript
// 基本的なAPIレスポンス型
type ApiResponse<T> = {
  data: T;
  message: string;
  status: number;
};

// ページネーション付きレスポンス
type PaginatedResponse<T> = ApiResponse<T[]> & {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
};

// エラーレスポンス
type ApiError = {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
  status: number;
};

// 使用例
type UserListResponse = PaginatedResponse<User>;
type UserDetailResponse = ApiResponse<User>;
```

### 2. フォームの型定義

```typescript
// フォームフィールドの型
type FormField<T> = {
  value: T;
  error?: string;
  touched: boolean;
};

// フォーム全体の型
type FormState<T> = {
  [K in keyof T]: FormField<T[K]>;
} & {
  isValid: boolean;
  isSubmitting: boolean;
};

// 使用例
type UserFormData = {
  name: string;
  email: string;
  age: number;
};

type UserFormState = FormState<UserFormData>;
// {
//   name: FormField<string>;
//   email: FormField<string>;
//   age: FormField<number>;
//   isValid: boolean;
//   isSubmitting: boolean;
// }
```

### 3. イベントハンドラーの型定義

```typescript
// イベントハンドラーの型
type EventHandler<T = HTMLElement> = (event: React.SyntheticEvent<T>) => void;

// 特定のイベントハンドラー
type ClickHandler = EventHandler<HTMLButtonElement>;
type ChangeHandler = (event: React.ChangeEvent<HTMLInputElement>) => void;

// フォームイベントハンドラー
type FormSubmitHandler<T> = (data: T) => void | Promise<void>;

// 使用例
type FormProps<T> = {
  onSubmit: FormSubmitHandler<T>;
  onChange?: ChangeHandler;
  onReset?: () => void;
};
```

## 型のテスト

### 1. 型のテストユーティリティ

```typescript
// 型のテスト用ヘルパー
type Expect<T extends true> = T;
type Equal<X, Y> = (<T>() => T extends X ? 1 : 2) extends <T>() => T extends Y ? 1 : 2 ? true : false;

// 型のテスト例
type Test1 = Expect<Equal<string, string>>; // ✅
type Test2 = Expect<Equal<string, number>>; // ❌ Type error

// 型の存在チェック
type Test3 = Expect<HasProperty<User, 'id'>>; // ✅
type Test4 = Expect<HasProperty<User, 'phone'>>; // ❌ Type error
```

## 参考資料

- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Advanced Types](https://www.typescriptlang.org/docs/handbook/2/types.html)
- [TypeScript Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html)
