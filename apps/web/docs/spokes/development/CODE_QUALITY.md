# コード品質管理

このドキュメントでは、Webアプリケーションのコード品質を維持するためのベストプラクティスについて説明します。

## Biome設定

### 1. Biome設定ファイル

```jsonc
// biome.jsonc
{
  "$schema": "https://biomejs.dev/schemas/1.9.4/schema.json",
  "organizeImports": {
    "enabled": true
  },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true,
      "correctness": {
        "useExhaustiveDependencies": "warn"
      },
      "style": {
        "noNonNullAssertion": "error",
        "useConst": "error"
      },
      "suspicious": {
        "noExplicitAny": "warn",
        "noArrayIndexKey": "warn"
      },
      "complexity": {
        "noExcessiveCognitiveComplexity": "warn"
      }
    }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2,
    "lineWidth": 100,
    "lineEnding": "lf"
  },
  "javascript": {
    "formatter": {
      "quoteStyle": "single",
      "semicolons": "asNeeded",
      "trailingCommas": "es5"
    }
  },
  "files": {
    "include": ["src/**/*"],
    "ignore": [
      "node_modules/**",
      "dist/**",
      "build/**",
      "coverage/**",
      "*.config.*"
    ]
  }
}
```

### 2. コード品質チェック

```bash
# コードの品質チェック
pnpm biome check --write ./src

# フォーマットのみ
pnpm biome format --write ./src

# リンターのみ
pnpm biome lint --write ./src

# 特定のファイルのみ
pnpm biome check --write ./src/components/Button.tsx
```

## TypeScript設定

### 1. 厳密な型チェック

```json
// tsconfig.json
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

### 2. 型定義のベストプラクティス

```typescript
// 良い例: 明確な型定義
type User = {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
};

// 悪い例: anyの使用
function processUser(user: any) {
  return user.name; // 型安全性がない
}

// 良い例: 適切な型定義
function processUser(user: User): string {
  return user.name; // 型安全
}
```


## コードレビューガイドライン

### 1. レビューチェックリスト

```markdown
## コードレビューチェックリスト

### 機能性
- [ ] 要件を満たしているか
- [ ] エラーハンドリングが適切か
- [ ] エッジケースが考慮されているか

### コード品質
- [ ] 型定義が適切か
- [ ] 関数が単一責任を持っているか
- [ ] 変数名が分かりやすいか
- [ ] コメントが適切か

### パフォーマンス
- [ ] 不要な再レンダリングがないか
- [ ] メモ化が適切に使用されているか
- [ ] 非同期処理が適切か

### セキュリティ
- [ ] 入力値の検証が適切か
- [ ] 機密情報が露出していないか
- [ ] XSS対策がされているか

### テスト
- [ ] テストが適切に書かれているか
- [ ] テストカバレッジが十分か
- [ ] テストが意味のある内容か
```

### 2. レビューコメント例

```typescript
// 良いレビューコメント
// ❌ 悪い例
// この関数は長すぎる

// ✅ 良い例
// この関数は複数の責任を持っているようです。
// データの取得とフォーマットを分離することを検討してください。

// 具体的な改善提案
// この部分を別の関数に抽出することで、可読性とテスタビリティが向上します。
```

## パフォーマンス監視

### 1. バンドルサイズ監視

```typescript
// next.config.ts
const nextConfig: NextConfig = {
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
      };
    }
    return config;
  },
  experimental: {
    optimizePackageImports: ['antd', 'lucide-react'],
  },
};
```

### 2. パフォーマンスメトリクス

```typescript
// utils/performance.ts
export function measurePerformance(name: string, fn: () => void) {
  const start = performance.now();
  fn();
  const end = performance.now();
  console.log(`${name} took ${end - start} milliseconds`);
}

// 使用例
measurePerformance('Data processing', () => {
  processLargeDataset(data);
});
```

## メモリリーク対策

### 1. イベントリスナーのクリーンアップ

```typescript
// 良い例: useEffectでのクリーンアップ
useEffect(() => {
  const handleResize = () => {
    // リサイズ処理
  };

  window.addEventListener('resize', handleResize);

  return () => {
    window.removeEventListener('resize', handleResize);
  };
}, []);
```

### 2. タイマーのクリーンアップ

```typescript
// 良い例: タイマーのクリーンアップ
useEffect(() => {
  const interval = setInterval(() => {
    // 定期的な処理
  }, 1000);

  return () => {
    clearInterval(interval);
  };
}, []);
```

## エラーハンドリング

### 1. エラーバウンダリ

```typescript
// components/ErrorBoundary.tsx
import React from 'react';

type ErrorBoundaryState = {
  hasError: boolean;
  error?: Error;
};

export class ErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  ErrorBoundaryState
> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>エラーが発生しました</h2>
          <p>ページを再読み込みしてください。</p>
        </div>
      );
    }

    return this.props.children;
  }
}
```

### 2. 非同期エラーハンドリング

```typescript
// 良い例: 非同期エラーハンドリング
async function fetchUserData(userId: string) {
  try {
    const response = await fetch(`/api/users/${userId}`);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Failed to fetch user data:', error);
    throw error;
  }
}
```

## セキュリティ対策

### 1. 入力値検証

```typescript
// 良い例: 入力値検証
import { z } from 'zod';

const userSchema = z.object({
  name: z.string().min(1, '名前は必須です'),
  email: z.string().email('有効なメールアドレスを入力してください'),
  age: z.number().min(0, '年齢は0以上である必要があります'),
});

function validateUserInput(input: unknown) {
  try {
    return userSchema.parse(input);
  } catch (error) {
    if (error instanceof z.ZodError) {
      throw new Error(error.errors.map(e => e.message).join(', '));
    }
    throw error;
  }
}
```

### 2. XSS対策

```typescript
// 良い例: XSS対策
import DOMPurify from 'dompurify';

function sanitizeHTML(html: string): string {
  return DOMPurify.sanitize(html);
}

// 危険な例: innerHTMLの直接使用
// element.innerHTML = userInput; // ❌

// 安全な例: サニタイズ後の使用
element.innerHTML = sanitizeHTML(userInput); // ✅
```

## テストカバレッジ

### 1. カバレッジ設定

```typescript
// vitest.config.ts
export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
        '**/*.d.ts',
        '**/*.stories.tsx',
        '**/*.config.*',
      ],
      thresholds: {
        global: {
          branches: 80,
          functions: 80,
          lines: 80,
          statements: 80,
        },
      },
    },
  },
});
```

### 2. カバレッジレポート

```bash
# カバレッジレポートの生成
pnpm test:cov

# HTMLレポートの確認
open coverage/index.html
```

## 参考資料

- [Biome Documentation](https://biomejs.dev/)
- [TypeScript Best Practices](https://typescript-eslint.io/rules/)
- [React Best Practices](https://react.dev/learn)
