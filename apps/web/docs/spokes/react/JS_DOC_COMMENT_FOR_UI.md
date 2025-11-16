# UI ComponentでのJSDocに関するルール

## ルール

UI Componentを作成する時、Markdown形式でのコメントを強制する。

## 背景

- Storybookでのコメントの表示に対応するため
- コンポーネントの使用方法やプロパティの説明を明確にするため
- 開発者間でのコンポーネント理解を促進するため

## 必須項目

### 1. コンポーネントの説明
```typescript
/**
 * ボタンコンポーネント
 *
 * 様々なサイズとバリアントを持つ再利用可能なボタンコンポーネントです。
 * クリックイベントやローディング状態をサポートします。
 */
export const Button = ({ ... }) => {
  // ...
}
```

### 2. Propsの説明
```typescript
type ButtonProps = {
  /**
   * ボタンのテキスト
   * @default "ボタン"
   */
  children: React.ReactNode;

  /**
   * ボタンのサイズ
   * - `sm`: 小さいサイズ（32px）
   * - `md`: 中サイズ（40px）
   * - `lg`: 大きいサイズ（48px）
   * @default "md"
   */
  size?: 'sm' | 'md' | 'lg';

  /**
   * ボタンのバリアント
   * - `primary`: メインアクション用
   * - `secondary`: サブアクション用
   * - `danger`: 危険なアクション用
   * @default "primary"
   */
  variant?: 'primary' | 'secondary' | 'danger';

  /**
   * ローディング状態かどうか
   * @default false
   */
  loading?: boolean;

  /**
   * クリック時のコールバック関数
   */
  onClick?: () => void;
};
```

### 3. 使用例
```typescript
/**
 * 使用例:
 *
 * ```tsx
 * // 基本的な使用
 * <Button onClick={() => console.log('clicked')}>
 *   クリック
 * </Button>
 *
 * // ローディング状態
 * <Button loading={true}>
 *   処理中...
 * </Button>
 *
 * // 危険なアクション
 * <Button variant="danger" onClick={handleDelete}>
 *   削除
 * </Button>
 * ```
 */
```

## 推奨項目

### 1. 注意事項や制限
```typescript
/**
 * 注意事項:
 * - `loading`が`true`の時は`onClick`は実行されません
 * - `disabled`と`loading`は同時に指定しないでください
 */
```

### 2. アクセシビリティ情報
```typescript
/**
 * アクセシビリティ:
 * - キーボードナビゲーション対応
 * - スクリーンリーダー対応
 * - フォーカス表示あり
 */
```

### 3. パフォーマンス情報
```typescript
/**
 * パフォーマンス:
 * - React.memoでメモ化済み
 * - 不要な再レンダリングを防止
 */
```

## 型定義について

UI Componentの型定義については、[TypeScript ベストプラクティス](../../typescript/TYPESCRIPT.md)に従ってください。

## 禁止事項

- 日本語と英語が混在したコメント
- 実装詳細の説明（コードを見れば分かる内容）
- 古い情報や不正確な情報
- 過度に長いコメント（1つの項目で3行を超える場合は要検討）

## チェックリスト

- [ ] コンポーネントの目的が明確に説明されている
- [ ] すべてのpropsに説明がある
- [ ] デフォルト値が明記されている
- [ ] 使用例が含まれている
- [ ] 注意事項や制限が記載されている（必要に応じて）
- [ ] アクセシビリティ情報が含まれている（必要に応じて）
- [ ] 型定義は[TYPESCRIPT.md](../../typescript/TYPESCRIPT.md)に準拠している
