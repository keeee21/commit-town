# src/utilsディレクトリのルール

## 概要

このドキュメントは、`apps/web/src/utils`ディレクトリで一般的なユーティリティ関数を作成・管理する際のルールを定義します。`src/libs`ディレクトリとは異なり、外部ライブラリに依存しない純粋なユーティリティ関数を配置します。

## LIBSディレクトリとの違い

| 項目 | `src/utils` | `src/libs` |
|------|-------------|------------|
| **目的** | 一般的なユーティリティ関数 | 外部ライブラリのラッパー・設定 |
| **依存関係** | 外部ライブラリに依存しない | 特定の外部ライブラリに依存 |
| **スコープ** | アプリケーション全体で使用 | 特定のライブラリの機能に特化 |
| **例** | 文字列操作、配列操作、日付処理 | MSW設定、Auth.js設定、Ant Design設定 |

## 基本原則

### 1. 命名規則

- **ファイル名**: `[機能名].ts` の形式（ケバブケース）
- **関数名**: 機能を明確に表現する名前
- **型名**: `[機能名]Options`, `[機能名]Result` の形式

```typescript
// ✅ 正しい例
// ファイル名: string-utils.ts
export const formatString = (str: string, options: FormatStringOptions) => {
  // 実装
};

// ❌ 間違った例
// ファイル名: stringHelpers.ts
export const format = (str: string) => {
  // 実装
};
```

### 2. 型定義

すべてのユーティリティ関数は以下の型定義を含む必要があります：

- **オプション型**: `[機能名]Options`
- **結果型**: `[機能名]Result`
- **型の明確性**: すべてのプロパティに適切な型注釈を付与

```typescript
// ✅ 正しい例
type FormatStringOptions = {
  maxLength?: number;
  ellipsis?: string;
  case?: 'upper' | 'lower' | 'title';
};

type FormatStringResult = {
  formatted: string;
  originalLength: number;
  truncated: boolean;
};
```

### 3. JSDocコメント

すべてのユーティリティ関数には以下のJSDocコメントを含む必要があります：

```typescript
/**
 * 文字列を指定されたオプションに従ってフォーマット
 *
 * 文字列の長さ制限、大文字小文字の変換、省略記号の追加などの
 * 基本的な文字列操作を提供します。
 *
 * @example
 * ```typescript
 * const result = formatString('Hello World', {
 *   maxLength: 10,
 *   ellipsis: '...',
 *   case: 'upper'
 * });
 * // result: { formatted: 'HELLO W...', originalLength: 11, truncated: true }
 * ```
 *
 * @param str - フォーマットする文字列
 * @param options - フォーマットオプション
 * @returns フォーマット結果
 */
export const formatString = (str: string, options: FormatStringOptions): FormatStringResult => {
  // 実装
};
```

## ディレクトリ構造

```
src/utils/
├── string-utils.ts          # 文字列操作関連
├── array-utils.ts           # 配列操作関連
├── date-utils.ts            # 日付操作関連
├── number-utils.ts          # 数値操作関連
├── object-utils.ts          # オブジェクト操作関連
├── validation-utils.ts      # バリデーション関連
├── format-utils.ts          # フォーマット関連
└── test-utils.ts            # テスト用ユーティリティ
```

## 各ユーティリティのルール

### 1. 文字列操作関連

```typescript
// string-utils.ts
type TruncateOptions = {
  length: number;
  ellipsis?: string;
  wordBoundary?: boolean;
};

export const truncateString = (str: string, options: TruncateOptions): string => {
  // 実装
};

export const capitalizeFirst = (str: string): string => {
  // 実装
};
```

### 2. 配列操作関連

```typescript
// array-utils.ts
type GroupByOptions<T> = {
  key: keyof T;
  groupName?: string;
};

export const groupBy = <T>(array: T[], options: GroupByOptions<T>): Record<string, T[]> => {
  // 実装
};

export const uniqueBy = <T>(array: T[], key: keyof T): T[] => {
  // 実装
};
```

### 3. 日付操作関連

```typescript
// date-utils.ts
type FormatDateOptions = {
  format?: 'YYYY-MM-DD' | 'MM/DD/YYYY' | 'DD/MM/YYYY';
  locale?: string;
};

export const formatDate = (date: Date, options: FormatDateOptions): string => {
  // 実装
};

export const isToday = (date: Date): boolean => {
  // 実装
};
```

## 共通パターン

### 1. 純粋関数の原則

すべてのユーティリティ関数は純粋関数である必要があります：

```typescript
// ✅ 正しい例
export const add = (a: number, b: number): number => {
  return a + b;
};

// ❌ 間違った例
export const add = (a: number, b: number): number => {
  console.log('Adding numbers'); // 副作用
  return a + b;
};
```

### 2. エラーハンドリング

適切なエラーハンドリングを含む必要があります：

```typescript
export const safeParseInt = (value: string, defaultValue: number = 0): number => {
  try {
    const parsed = parseInt(value, 10);
    return isNaN(parsed) ? defaultValue : parsed;
  } catch {
    return defaultValue;
  }
};
```

### 3. 型安全性

すべての関数は型安全である必要があります：

```typescript
export const createTypedMap = <K extends string, V>(
  entries: Array<[K, V]>
): Record<K, V> => {
  return Object.fromEntries(entries) as Record<K, V>;
};
```

## 禁止事項

### 1. 外部ライブラリへの依存

外部ライブラリに依存する関数は`src/libs`に配置してください：

```typescript
// ❌ 間違った例（src/utilsに配置）
import { format } from 'date-fns';

export const formatDate = (date: Date) => {
  return format(date, 'yyyy-MM-dd');
};

// ✅ 正しい例（src/libsに配置）
// src/libs/date-fns.ts
import { format } from 'date-fns';

export const createDateFormatter = (formatString: string) => {
  return (date: Date) => format(date, formatString);
};
```

### 2. ビジネスロジックの混入

ビジネスロジックを含めてはいけません：

```typescript
// ❌ 間違った例
export const formatString = (str: string) => {
  // ビジネスロジックは含めない
  if (str.includes('user')) {
    return str.toUpperCase();
  }
  return str;
};
```

### 3. グローバル状態の操作

グローバル状態を操作してはいけません：

```typescript
// ❌ 間違った例
export const formatString = (str: string) => {
  // グローバル状態の操作は禁止
  window.localStorage.setItem('lastFormatted', str);
  return str.toUpperCase();
};
```

## テスト

ユーティリティ関数には適切なテストを記述します：

```typescript
// string-utils.test.ts
import { truncateString, capitalizeFirst } from './string-utils';

describe('string-utils', () => {
  describe('truncateString', () => {
    it('should truncate string when longer than max length', () => {
      const result = truncateString('Hello World', { length: 5 });
      expect(result).toBe('Hello');
    });

    it('should add ellipsis when truncating', () => {
      const result = truncateString('Hello World', { length: 5, ellipsis: '...' });
      expect(result).toBe('He...');
    });
  });

  describe('capitalizeFirst', () => {
    it('should capitalize first letter', () => {
      const result = capitalizeFirst('hello');
      expect(result).toBe('Hello');
    });
  });
});
```

## メンテナンス

### 1. 定期的な見直し

- 月次でユーティリティ関数の使用状況を確認
- 未使用の関数の特定と削除
- パフォーマンスの最適化

### 2. ドキュメントの更新

- 新しい関数追加時のルール適用確認
- 既存関数の変更時のドキュメント更新
- 使用例の追加・更新

### 3. コードレビュー

- 新しいユーティリティ関数の作成時は必ずコードレビューを実施
- ルールに準拠しているかの確認
- テストカバレッジの確認

## 参考資料

- [TypeScript公式ドキュメント](https://www.typescriptlang.org/docs/)
- [Next.js公式ドキュメント](https://nextjs.org/docs)
- [関数型プログラミングの原則](https://en.wikipedia.org/wiki/Functional_programming)
- [プロジェクト内TypeScriptベストプラクティス](./typescript/TYPESCRIPT.md)
