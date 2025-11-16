# src/libsディレクトリのルール

## 概要

このドキュメントは、`apps/web/src/libs`ディレクトリでライブラリ関連のユーティリティ関数を作成・管理する際のルールを定義します。一貫性のあるコード品質と保守性を確保することを目的としています。

## 基本原則

### 1. 命名規則

- **ファイル名**: `[ライブラリ名].ts` の形式（ケバブケース）
- **関数名**: ライブラリの機能を明確に表現する名前
- **型名**: `[ライブラリ名][機能名]Config`, `[ライブラリ名][機能名]Options` の形式

```typescript
// ✅ 正しい例
// ファイル名: fetcher.ts
export const createFetcher = (config: FetcherConfig) => {
  // 実装
};

// ❌ 間違った例
// ファイル名: fetchUtils.ts
export const fetchData = () => {
  // 実装
};
```

### 2. 型定義

すべてのライブラリ関数は以下の型定義を含む必要があります：

- **設定型**: `[ライブラリ名]Config`
- **オプション型**: `[ライブラリ名]Options`
- **レスポンス型**: `[ライブラリ名]Response`
- **型の明確性**: すべてのプロパティに適切な型注釈を付与

```typescript
// ✅ 正しい例
type FetcherConfig = {
  baseUrl: string;
  timeout?: number;
  retries?: number;
};

type FetcherOptions = {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  headers?: Record<string, string>;
  body?: unknown;
};
```

### 3. JSDocコメント

すべてのライブラリ関数には以下のJSDocコメントを含む必要があります：

- **概要**: 関数の目的と機能の説明
- **使用例**: 実際のコード例を含む
- **パラメータ**: 各パラメータの説明
- **リターン値**: リターン値の説明

```typescript
/**
 * HTTPリクエストを送信するためのフェッチャーを作成
 *
 * 設定されたベースURLとタイムアウトを使用してHTTPリクエストを送信します。
 * リトライ機能とエラーハンドリングを含みます。
 *
 * @example
 * ```typescript
 * const fetcher = createFetcher({
 *   baseUrl: 'https://api.example.com',
 *   timeout: 5000,
 *   retries: 3
 * });
 *
 * const data = await fetcher('/users', {
 *   method: 'GET',
 *   headers: { 'Authorization': 'Bearer token' }
 * });
 * ```
 *
 * @param config - フェッチャーの設定
 * @param config.baseUrl - ベースURL
 * @param config.timeout - タイムアウト時間（ミリ秒）
 * @param config.retries - リトライ回数
 * @returns 設定されたフェッチャー関数
 */
export const createFetcher = (config: FetcherConfig) => {
  // 実装
};
```

## ディレクトリ構造

```
src/libs/
├── fetcher.ts              # HTTPリクエスト関連
├── msw.ts                  # MSW（Mock Service Worker）関連
├── openapi.ts              # OpenAPI関連
├── authjs.ts               # Auth.js関連
├── antd.ts                 # Ant Design関連
├── date-fns.ts             # date-fns関連
├── zod.ts                  # Zod関連
└── utils.ts                # その他のユーティリティ
```

## 各ライブラリのルール

### 1. Fetcher関連

HTTPリクエストを送信するための関数を配置します。

```typescript
// fetcher.ts
type FetcherConfig = {
  baseUrl: string;
  timeout?: number;
  retries?: number;
};

type FetcherOptions = {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  headers?: Record<string, string>;
  body?: unknown;
};

export const createFetcher = (config: FetcherConfig) => {
  // 実装
};

export const createApiClient = (baseUrl: string) => {
  // 実装
};
```

### 2. MSW関連

Mock Service Workerの設定とハンドラーを配置します。

```typescript
// msw.ts
import { setupServer } from 'msw/node';
import { rest } from 'msw';

export const createMockServer = (handlers: RequestHandler[]) => {
  // 実装
};

export const createMockHandlers = () => {
  // 実装
};
```

### 3. OpenAPI関連

OpenAPIスキーマから生成された型や関数を配置します。

```typescript
// openapi.ts
import type { operations } from '@/types/openapi';

export const createApiClient = (baseUrl: string) => {
  // 実装
};

export const createApiError = (error: unknown) => {
  // 実装
};
```

### 4. Auth.js関連

認証関連のユーティリティ関数を配置します。

```typescript
// authjs.ts
import { getServerSession } from 'next-auth';

export const getSession = async () => {
  // 実装
};

export const requireAuth = async () => {
  // 実装
};
```

### 5. Ant Design関連

Ant Designコンポーネントの設定やユーティリティを配置します。

```typescript
// antd.ts
import type { ConfigProviderProps } from 'antd';

export const createAntdConfig = (): ConfigProviderProps => {
  // 実装
};

export const createAntdTheme = () => {
  // 実装
};
```

## 共通パターン

### 1. 設定オブジェクトの作成

ライブラリの設定を作成する関数は以下のパターンに従います：

```typescript
type LibraryConfig = {
  baseUrl: string;
  timeout?: number;
  retries?: number;
};

export const createLibraryConfig = (config: LibraryConfig) => {
  return {
    baseUrl: config.baseUrl,
    timeout: config.timeout ?? 5000,
    retries: config.retries ?? 3,
  };
};
```

### 2. エラーハンドリング

すべてのライブラリ関数には適切なエラーハンドリングを含む必要があります：

```typescript
export const safeApiCall = async <T>(
  apiCall: () => Promise<T>,
  fallback: T
): Promise<T> => {
  try {
    return await apiCall();
  } catch (error) {
    console.error('API call failed:', error);
    return fallback;
  }
};
```

### 3. 型安全性

すべての関数は型安全である必要があります：

```typescript
export const createTypedFetcher = <TResponse, TRequest = unknown>(
  baseUrl: string
) => {
  return async (endpoint: string, options?: FetcherOptions<TRequest>): Promise<TResponse> => {
    // 実装
  };
};
```

## 禁止事項

### 1. ビジネスロジックの混入

ライブラリ関数にビジネスロジックを含めてはいけません：

```typescript
// ❌ 間違った例
export const createFetcher = (config: FetcherConfig) => {
  return async (endpoint: string) => {
    // ビジネスロジックは含めない
    if (endpoint.includes('user')) {
      // ユーザー関連の特別な処理
    }
  };
};
```

### 2. グローバル状態の直接操作

グローバル状態を直接操作してはいけません：

```typescript
// ❌ 間違った例
export const createFetcher = (config: FetcherConfig) => {
  return async (endpoint: string) => {
    // グローバル状態の直接操作は禁止
    globalStore.setLoading(true);
  };
};
```

### 3. 副作用の不適切な管理

副作用を適切に管理しない場合：

```typescript
// ❌ 間違った例
export const createFetcher = (config: FetcherConfig) => {
  // 副作用の不適切な管理
  document.title = 'Loading...';

  return async (endpoint: string) => {
    // 実装
  };
};
```

## テスト

ライブラリ関数には適切なテストを記述します：

```typescript
// fetcher.test.ts
import { createFetcher } from './fetcher';

describe('createFetcher', () => {
  it('should create fetcher with correct config', () => {
    const config = {
      baseUrl: 'https://api.example.com',
      timeout: 5000,
      retries: 3
    };

    const fetcher = createFetcher(config);

    expect(fetcher).toBeDefined();
  });

  it('should handle errors correctly', async () => {
    const fetcher = createFetcher({
      baseUrl: 'https://api.example.com'
    });

    await expect(fetcher('/invalid')).rejects.toThrow();
  });
});
```

## メンテナンス

### 1. 定期的な見直し

- 月次でライブラリ関数の使用状況を確認
- 未使用の関数の特定と削除
- パフォーマンスの最適化

### 2. ドキュメントの更新

- 新しい関数追加時のルール適用確認
- 既存関数の変更時のドキュメント更新
- 使用例の追加・更新

### 3. コードレビュー

- 新しいライブラリ関数の作成時は必ずコードレビューを実施
- ルールに準拠しているかの確認
- テストカバレッジの確認

## 参考資料

- [TypeScript公式ドキュメント](https://www.typescriptlang.org/docs/)
- [Next.js公式ドキュメント](https://nextjs.org/docs)
- [Ant Design公式ドキュメント](https://ant.design/)
- [MSW公式ドキュメント](https://mswjs.io/)
- [Auth.js公式ドキュメント](https://authjs.dev/)
