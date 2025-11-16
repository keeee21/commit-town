# API 型定義パターン

このドキュメントでは、WebアプリケーションのAPI型定義のベストプラクティスについて説明します。

## OpenAPI との連携

### 1. 自動生成された型の活用

```typescript
// api から生成された型を使用
import type {
  User,
  Talent,
  CreateUserRequest,
  UpdateUserRequest,
  ApiError
} from 'openapi';

// 型の拡張
type ExtendedUser = User & {
  fullName: string;
  isActive: boolean;
};

// 型の変換
type UserFormData = Omit<User, 'id' | 'createdAt' | 'updatedAt'>;
type UserTableData = Pick<User, 'id' | 'name' | 'email' | 'status'>;
```

### 2. API レスポンスの型定義

```typescript
// 基本的なAPIレスポンス型
type ApiResponse<T> = {
  data: T;
  message: string;
  status: 'success' | 'error';
  timestamp: string;
};

// ページネーション付きレスポンス
type PaginatedResponse<T> = {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
  message: string;
  status: 'success' | 'error';
  timestamp: string;
};

// エラーレスポンス
type ApiErrorResponse = {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
    field?: string;
  };
  status: 'error';
  timestamp: string;
};
```

## リクエスト・レスポンス型

### 1. CRUD 操作の型定義

```typescript
// 作成リクエスト
type CreateRequest<T> = Omit<T, 'id' | 'createdAt' | 'updatedAt'>;

// 更新リクエスト
type UpdateRequest<T> = Partial<CreateRequest<T>> & { id: string };

// 削除リクエスト
type DeleteRequest = {
  id: string;
};

// 一覧取得リクエスト
type ListRequest = {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
  search?: string;
  filter?: Record<string, unknown>;
};

// 使用例
type CreateUserRequest = CreateRequest<User>;
type UpdateUserRequest = UpdateRequest<User>;
type UserListRequest = ListRequest;
```

### 2. 検索・フィルタリング型

```typescript
// 検索条件の型
type SearchCriteria = {
  query?: string;
  fields?: string[];
  exact?: boolean;
};

// フィルター条件の型
type FilterCriteria = {
  [key: string]: {
    operator: 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'nin' | 'like';
    value: unknown;
  };
};

// ソート条件の型
type SortCriteria = {
  field: string;
  direction: 'asc' | 'desc';
};

// 複合クエリ型
type QueryParams = SearchCriteria & FilterCriteria & {
  sort?: SortCriteria[];
  page?: number;
  limit?: number;
};
```

## エラーハンドリング型

### 1. エラー型の定義

```typescript
// 基本的なエラー型
type BaseError = {
  code: string;
  message: string;
  timestamp: string;
};

// バリデーションエラー
type ValidationError = BaseError & {
  code: 'VALIDATION_ERROR';
  details: {
    field: string;
    message: string;
    value: unknown;
  }[];
};

// 認証エラー
type AuthenticationError = BaseError & {
  code: 'AUTHENTICATION_ERROR';
  details: {
    reason: 'INVALID_TOKEN' | 'TOKEN_EXPIRED' | 'MISSING_TOKEN';
  };
};

// 認可エラー
type AuthorizationError = BaseError & {
  code: 'AUTHORIZATION_ERROR';
  details: {
    required: string[];
    provided: string[];
  };
};

// ビジネスロジックエラー
type BusinessLogicError = BaseError & {
  code: 'BUSINESS_LOGIC_ERROR';
  details: {
    constraint: string;
    value: unknown;
  };
};

// エラーの判別共用体
type ApiError =
  | ValidationError
  | AuthenticationError
  | AuthorizationError
  | BusinessLogicError;
```

### 2. エラーハンドリング関数

```typescript
// エラーハンドリング関数（型ガードを使わない方法）
function handleApiError(error: ApiError): string {
  if (error.code === 'VALIDATION_ERROR') {
    return error.details.map(d => `${d.field}: ${d.message}`).join(', ');
  }

  if (error.code === 'AUTHENTICATION_ERROR') {
    return '認証に失敗しました。再度ログインしてください。';
  }

  return error.message;
}
```

## フォーム型定義

### 1. フォームデータ型

```typescript
// フォームフィールドの状態
type FormField<T> = {
  value: T;
  error?: string;
  touched: boolean;
  dirty: boolean;
};

// フォーム全体の状態
type FormState<T> = {
  values: T;
  errors: Partial<Record<keyof T, string>>;
  touched: Partial<Record<keyof T, boolean>>;
  dirty: Partial<Record<keyof T, boolean>>;
  isValid: boolean;
  isSubmitting: boolean;
  isDirty: boolean;
};

// フォームアクション
type FormActions<T> = {
  setValue: <K extends keyof T>(field: K, value: T[K]) => void;
  setError: (field: keyof T, error: string) => void;
  clearError: (field: keyof T) => void;
  setTouched: (field: keyof T, touched: boolean) => void;
  reset: () => void;
  submit: () => Promise<void>;
};
```

### 2. バリデーション型

```typescript
// バリデーションルール
type ValidationRule<T> = {
  required?: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  custom?: (value: T) => string | undefined;
};

// フィールドバリデーション
type FieldValidation<T> = {
  [K in keyof T]?: ValidationRule<T[K]>;
};

// フォームバリデーション
type FormValidation<T> = {
  fields: FieldValidation<T>;
  submit?: (values: T) => Promise<string | undefined>;
};

// 使用例
const userFormValidation: FormValidation<UserFormData> = {
  fields: {
    name: {
      required: true,
      minLength: 2,
      maxLength: 50,
    },
    email: {
      required: true,
      pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    },
    age: {
      required: true,
      custom: (value) => value < 0 ? '年齢は0以上である必要があります' : undefined,
    },
  },
};
```

## 状態管理型

### 1. 非同期状態の型

```typescript
// 非同期状態の基本型
type AsyncState<T> = {
  data: T | null;
  loading: boolean;
  error: string | null;
};

// 非同期状態の操作
type AsyncActions<T> = {
  fetch: () => Promise<void>;
  refetch: () => Promise<void>;
  reset: () => void;
  setData: (data: T) => void;
  setError: (error: string) => void;
  setLoading: (loading: boolean) => void;
};

// 完全な非同期状態
type AsyncStateWithActions<T> = AsyncState<T> & AsyncActions<T>;
```

### 2. キャッシュ型

```typescript
// キャッシュエントリ
type CacheEntry<T> = {
  data: T;
  timestamp: number;
  ttl: number; // Time to live in milliseconds
};

// キャッシュの状態
type CacheState<T> = {
  entries: Map<string, CacheEntry<T>>;
  maxSize: number;
  defaultTtl: number;
};

// キャッシュ操作
type CacheActions<T> = {
  get: (key: string) => T | null;
  set: (key: string, data: T, ttl?: number) => void;
  delete: (key: string) => void;
  clear: () => void;
  has: (key: string) => boolean;
  isExpired: (key: string) => boolean;
};
```

## 型のテスト

### 1. 型のテストユーティリティ

```typescript
// 型の等価性チェック
type Equals<X, Y> = (<T>() => T extends X ? 1 : 2) extends <T>() => T extends Y ? 1 : 2 ? true : false;

// 型のテスト
type Test1 = Equals<CreateRequest<User>, Omit<User, 'id' | 'createdAt' | 'updatedAt'>>; // true
type Test2 = Equals<UpdateRequest<User>, Partial<CreateRequest<User>> & { id: string }>; // true

// 型の存在チェック
type HasProperty<T, K extends PropertyKey> = K extends keyof T ? true : false;

type Test3 = HasProperty<ApiResponse<User>, 'data'>; // true
type Test4 = HasProperty<ApiResponse<User>, 'invalid'>; // false
```

## 参考資料

- [OpenAPI TypeScript Generator](https://openapi-ts.pages.dev/)
- [TypeScript API Design Patterns](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html)
- [REST API TypeScript Patterns](https://blog.logrocket.com/typescript-rest-api/)
