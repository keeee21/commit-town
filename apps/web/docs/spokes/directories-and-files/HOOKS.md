# カスタムフック用のルール

## 概要

このドキュメントは、`apps/web/src/hooks`ディレクトリでカスタムフックを作成・管理する際のルールを定義します。一貫性のあるコード品質と保守性を確保することを目的としています。

## 基本原則

### 1. 命名規則

- **ファイル名**: `use-[機能名].ts` の形式（ケバブケース）
- **フック名**: `use[機能名]` の形式（パスカルケース）
- **型名**: `Use[機能名]Params`, `Use[機能名]Return` の形式

```typescript
// ✅ 正しい例
// ファイル名: use-personal-data.ts
export const usePersonalData = ({ talentId, initialData }: UsePersonalDataParams): UsePersonalDataReturn => {
  // 実装
};

// ❌ 間違った例
// ファイル名: personalDataHook.ts
export const personalDataHook = () => {
  // 実装
};
```

### 2. 型定義

すべてのカスタムフックは以下の型定義を含む必要があります：

- **パラメータ型**: `Use[機能名]Params`
- **リターン型**: `Use[機能名]Return`
- **型の明確性**: すべてのプロパティに適切な型注釈を付与

```typescript
// ✅ 正しい例
type UsePersonalDataParams = {
  talentId: operations["TalentController_get"]["parameters"]["path"]["talentId"];
  initialData: FindPersonalResponseSchema | null;
};

type UsePersonalDataReturn = {
  personal: FindPersonalResponseSchema | null;
  isPending: boolean;
  refetch: () => void;
};
```

### 3. JSDocコメント

すべてのカスタムフックには以下のJSDocコメントを含む必要があります：

- **概要**: フックの目的と機能の説明
- **使用例**: 実際のコード例を含む
- **パラメータ**: 各パラメータの説明
- **リターン値**: 各リターン値の説明

```typescript
/**
 * 個人データの取得・管理を行うカスタムフック
 *
 * 指定されたtalentIdに基づいて個人データを取得し、
 * ローディング状態とリフェッチ機能を提供します。
 *
 * @example
 * ```tsx
 * const { personal, isPending, refetch } = usePersonalData({
 *   talentId: "123",
 *   initialData: null
 * });
 *
 * if (isPending) return <Loading />;
 * return <div>{personal?.name}</div>;
 * ```
 *
 * @param params - フックのパラメータ
 * @param params.talentId - タレントID
 * @param params.initialData - 初期データ
 * @returns 個人データとローディング状態
 */
export const usePersonalData = ({ talentId, initialData }: UsePersonalDataParams): UsePersonalDataReturn => {
  // 実装
};
```

### 4. エラーハンドリング

すべての非同期処理には適切なエラーハンドリングを含む必要があります：

```typescript
// ✅ 正しい例
const refetch = useCallback(() => {
  startTransition(async () => {
    try {
      const response = await fetch(`/api/talent/${talentId}/personal`);
      const result = await response.json();
      setPersonal(result.success ? result.data : null);
    } catch (error) {
      console.error("Failed to fetch personal data:", error);
      setPersonal(null);
    }
  });
}, [talentId]);
```

### 5. 状態管理

適切なReactフックを使用して状態を管理します：

- **useState**: ローカル状態の管理
- **useTransition**: 非同期処理のローディング状態
- **useCallback**: 関数のメモ化

```typescript
// ✅ 正しい例
const [isPending, startTransition] = useTransition();
const [personal, setPersonal] = useState<FindPersonalResponseSchema | null>(initialData);

const refetch = useCallback(() => {
  startTransition(async () => {
    // 非同期処理
  });
}, [talentId]);
```

### 6. 依存関係の管理

useCallbackやuseEffectの依存関係配列を適切に管理します：

```typescript
// ✅ 正しい例
const refetch = useCallback(() => {
  // 実装
}, [talentId]); // 必要な依存関係のみを含める

// ❌ 間違った例
const refetch = useCallback(() => {
  // 実装
}, []); // 依存関係が不足
```

## ディレクトリ構造

```
src/hooks/
├── use-orion-notification.ts    # 通知機能
├── use-personal-data.ts         # 個人データ管理
├── use-sales-data.ts           # 売上データ管理
└── use-service-management-data.ts # サービス管理データ
```

## 共通パターン

### データ取得フック

データ取得を行うカスタムフックは以下のパターンに従います：

```typescript
type UseDataParams = {
  id: string;
  initialData: DataType | null;
};

type UseDataReturn = {
  data: DataType | null;
  isPending: boolean;
  refetch: () => void;
};

export const useData = ({ id, initialData }: UseDataParams): UseDataReturn => {
  const [isPending, startTransition] = useTransition();
  const [data, setData] = useState<DataType | null>(initialData);

  const refetch = useCallback(() => {
    startTransition(async () => {
      try {
        const response = await fetch(`/api/data/${id}`);
        const result = await response.json();
        setData(result.success ? result.data : null);
      } catch (error) {
        console.error("Failed to fetch data:", error);
        setData(null);
      }
    });
  }, [id]);

  return {
    data,
    isPending,
    refetch,
  };
};
```

### 通知フック

通知機能を提供するカスタムフックは以下のパターンに従います：

```typescript
type NotificationParams = {
  message: React.ReactNode;
  description?: React.ReactNode;
  placement?: "topLeft" | "topRight" | "bottomLeft" | "bottomRight";
  type?: "success" | "info" | "warning" | "error";
};

export const useNotification = () => {
  const [api, contextHolder] = notification.useNotification();

  const openNotification = (params: NotificationParams) => {
    const { message, description, placement = "bottomLeft", type = "info" } = params;
    api[type]({
      message,
      description,
      placement,
    });
  };

  return {
    openNotification,
    contextHolder,
  };
};
```

## 禁止事項

### 1. 直接的なDOM操作

カスタムフック内でDOMを直接操作してはいけません：

```typescript
// ❌ 間違った例
export const useBadHook = () => {
  useEffect(() => {
    document.getElementById('myElement')?.click(); // DOM操作は禁止
  }, []);
};
```

### 2. グローバル状態の直接操作

グローバル状態を直接操作してはいけません：

```typescript
// ❌ 間違った例
export const useBadHook = () => {
  const setGlobalState = useGlobalStore(state => state.setValue);
  // グローバル状態の直接操作は禁止
};
```

### 3. 副作用の不適切な管理

useEffectの依存関係を適切に管理しない場合：

```typescript
// ❌ 間違った例
export const useBadHook = () => {
  useEffect(() => {
    // 依存関係が不足している
  }, []); // 必要な依存関係が含まれていない
};
```

## テスト

カスタムフックには適切なテストを記述します：

```typescript
// use-personal-data.test.ts
import { renderHook, act } from '@testing-library/react';
import { usePersonalData } from './use-personal-data';

describe('usePersonalData', () => {
  it('should return initial data', () => {
    const { result } = renderHook(() =>
      usePersonalData({
        talentId: '123',
        initialData: { name: 'Test User' }
      })
    );

    expect(result.current.personal).toEqual({ name: 'Test User' });
    expect(result.current.isPending).toBe(false);
  });

  it('should refetch data', async () => {
    const { result } = renderHook(() =>
      usePersonalData({
        talentId: '123',
        initialData: null
      })
    );

    await act(async () => {
      result.current.refetch();
    });

    // アサーション
  });
});
```

## メンテナンス

### 1. 定期的な見直し

- 月次でカスタムフックの使用状況を確認
- 未使用のフックの特定と削除
- パフォーマンスの最適化

### 2. ドキュメントの更新

- 新しいフック追加時のルール適用確認
- 既存フックの変更時のドキュメント更新
- 使用例の追加・更新

### 3. コードレビュー

- 新しいカスタムフックの作成時は必ずコードレビューを実施
- ルールに準拠しているかの確認
- テストカバレッジの確認

## 参考資料

- [React Hooks公式ドキュメント](https://react.dev/reference/react)
- [TypeScript公式ドキュメント](https://www.typescriptlang.org/docs/)
- [Ant Design公式ドキュメント](https://ant.design/components/notification)
- [Next.js公式ドキュメント](https://nextjs.org/docs)

---

## useEffectのルール

useEffectのルールについては以下のドキュメントを参照してください。

- [useEffectのルール](../react/USEEFFECT_RULES.md)
