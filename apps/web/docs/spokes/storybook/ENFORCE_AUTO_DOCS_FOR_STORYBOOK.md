# Storybook自動ドキュメント強制ルール

## 概要
Storybookでコンポーネントを作成・更新する際に、`tags: ["autodocs"]`を使用して自動的にドキュメントを生成することを強制する。

## 強制事項

### 1. 必須プロパティ
すべてのStorybookストーリーファイル（`.stories.tsx`）には`tags: ["autodocs"]`を使用した自動ドキュメント生成を必須とする設定が必須です：

```typescript
export default {
  title: 'Library/ComponentName', // 適切な階層構造
  component: ComponentName,
  tags: ['autodocs'], // 自動ドキュメント生成を有効化
  argTypes: {
    // すべてのpropsの説明
  }
} satisfies Meta<typeof ComponentName>;
```

## 参考資料

- [Storybook Autodocs](https://storybook.js.org/docs/7.0/writing-docs/autodocs)
- [Storybook Controls](https://storybook.js.org/docs/essentials/controls/)
- [Storybook Args](https://storybook.js.org/docs/essentials/args/)
