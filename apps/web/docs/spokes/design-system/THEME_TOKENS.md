# テーマトークン設計

このドキュメントでは、CWebアプリケーションのデザインシステムにおけるテーマトークンの設計について説明します。

## トークン階層

### 1. Primitive Tokens（基本トークン）

```typescript
// primitive-tokens.ts
export const primitiveTokens = {
  // カラーパレット
  gray: {
    15: "#0D0D0D",
    22: "#1B1B1B",
    30: "#2E2E2E",
    // ... 他の色調
  },
  blue: {
    22: "#00134c",
    30: "#00276d",
    38: "#0b3d8b",
    // ... 他の色調
  },
  // その他の色パレット
} as const;
```

### 2. Semantic Tokens（意味的トークン）

```typescript
// semantic-tokens.ts
export const semanticTokens = {
  color: {
    // 背景色
    background: {
      primary: primitiveTokens.gray[98],
      secondary: primitiveTokens.gray[96],
      tertiary: primitiveTokens.gray[93],
    },
    // テキスト色
    text: {
      primary: primitiveTokens.gray[15],
      secondary: primitiveTokens.gray[46],
      tertiary: primitiveTokens.gray[62],
      inverse: primitiveTokens.white[100],
    },
    // アクセント色
    accent: {
      primary: primitiveTokens.blue[54],
      secondary: primitiveTokens.blue[38],
      hover: primitiveTokens.blue[62],
      active: primitiveTokens.blue[46],
    },
    // 状態色
    status: {
      success: primitiveTokens.green[54],
      warning: primitiveTokens.yellow[54],
      error: primitiveTokens.red[54],
      info: primitiveTokens.blue[54],
    },
  },
  // その他のトークン
} as const;
```

## カラートークン

### 1. カラーパレット設計

```typescript
// カラーパレットの構造
type ColorPalette = {
  [shade: number]: string;
};

type ColorTokens = {
  gray: ColorPalette;
  blue: ColorPalette;
  red: ColorPalette;
  green: ColorPalette;
  yellow: ColorPalette;
  black: ColorPalette;
  white: ColorPalette;
}

// 色調の命名規則
// 15, 22, 30, 38, 46, 54, 62, 70, 78, 85, 90, 93, 96, 98
// 数値が大きいほど明るい色
```

### 2. 意味的カラー設計

```typescript
// 意味的カラーの定義
type SemanticColors = {
  // 背景色
  background: {
    primary: string;    // メイン背景
    secondary: string;  // セカンダリ背景
    tertiary: string;   // サード背景
    inverse: string;    // 反転背景
  };

  // テキスト色
  text: {
    primary: string;    // メインテキスト
    secondary: string;  // セカンダリテキスト
    tertiary: string;   // サードテキスト
    inverse: string;    // 反転テキスト
    disabled: string;   // 無効テキスト
  };

  // アクセント色
  accent: {
    primary: string;    // メインアクセント
    secondary: string;  // セカンダリアクセント
    hover: string;      // ホバー状態
    active: string;     // アクティブ状態
    focus: string;      // フォーカス状態
  };

  // 状態色
  status: {
    success: string;    // 成功状態
    warning: string;    // 警告状態
    error: string;      // エラー状態
    info: string;       // 情報状態
  };

  // ボーダー色
  border: {
    primary: string;    // メインボーダー
    secondary: string;  // セカンダリボーダー
    focus: string;      // フォーカスボーダー
    error: string;      // エラーボーダー
  };
}
```

## タイポグラフィトークン

### 1. フォントサイズ

```typescript
// フォントサイズトークン
export const fontSizeTokens = {
  xs: '12px',      // 12px
  sm: '14px',      // 14px
  base: '16px',    // 16px
  lg: '18px',      // 18px
  xl: '20px',      // 20px
  '2xl': '24px',   // 24px
  '3xl': '30px',   // 30px
  '4xl': '36px',   // 36px
  '5xl': '48px',   // 48px
} as const;
```

### 2. フォントウェイト

```typescript
// フォントウェイトトークン
export const fontWeightTokens = {
  thin: 100,
  extralight: 200,
  light: 300,
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
  extrabold: 800,
  black: 900,
} as const;
```

### 3. 行間

```typescript
// 行間トークン
export const lineHeightTokens = {
  none: 1,
  tight: 1.25,
  snug: 1.375,
  normal: 1.5,
  relaxed: 1.625,
  loose: 2,
} as const;
```

## スペーシングトークン

### 1. スペーシングスケール

```typescript
// スペーシングトークン
export const spacingTokens = {
  0: '0px',
  1: '4px',    // 0.25rem
  2: '8px',    // 0.5rem
  3: '12px',   // 0.75rem
  4: '16px',   // 1rem
  5: '20px',   // 1.25rem
  6: '24px',   // 1.5rem
  8: '32px',   // 2rem
  10: '40px',  // 2.5rem
  12: '48px',  // 3rem
  16: '64px',  // 4rem
  20: '80px',  // 5rem
  24: '96px',  // 6rem
  32: '128px', // 8rem
} as const;
```

### 2. ブレークポイント

```typescript
// ブレークポイントトークン
export const breakpointTokens = {
  xs: '0px',
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
} as const;
```

## シャドウトークン

### 1. ボックスシャドウ

```typescript
// シャドウトークン
export const shadowTokens = {
  none: 'none',
  sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
  base: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
  md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  xl: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
  '2xl': '0 25px 50px -12px rgb(0 0 0 / 0.25)',
  inner: 'inset 0 2px 4px 0 rgb(0 0 0 / 0.05)',
} as const;
```

## ボーダーラディウストークン

### 1. 角丸

```typescript
// ボーダーラディウストークン
export const borderRadiusTokens = {
  none: '0px',
  sm: '2px',
  base: '4px',
  md: '6px',
  lg: '8px',
  xl: '12px',
  '2xl': '16px',
  '3xl': '24px',
  full: '9999px',
} as const;
```

## トークンの使用

### 1. CSS変数としての使用

```typescript
// CSS変数の生成
export function generateCSSVariables(tokens: Record<string, any>): string {
  return Object.entries(tokens)
    .map(([key, value]) => `--${key}: ${value};`)
    .join('\n');
}

// 使用例
const cssVariables = generateCSSVariables(semanticTokens);
```

### 2. Tailwind CSSとの統合


### 2. TypeScript型の生成

```typescript
// トークンの型定義
export type PrimitiveTokens = typeof primitiveTokens;
export type SemanticTokens = typeof semanticTokens;
export type FontSizeTokens = typeof fontSizeTokens;
export type FontWeightTokens = typeof fontWeightTokens;
export type LineHeightTokens = typeof lineHeightTokens;
export type SpacingTokens = typeof spacingTokens;
export type ShadowTokens = typeof shadowTokens;
export type BorderRadiusTokens = typeof borderRadiusTokens;
```

## トークンの管理

### 1. トークンの検証

```typescript
// トークンの検証関数
export function validateTokens(tokens: Record<string, any>): boolean {
  // トークンの一貫性チェック
  // 色の値の妥当性チェック
  // 数値の範囲チェック
  return true;
}
```

### 2. トークンのドキュメント生成

```typescript
// トークンのドキュメント生成
export function generateTokenDocumentation(tokens: Record<string, any>): string {
  return Object.entries(tokens)
    .map(([key, value]) => `- **${key}**: \`${value}\``)
    .join('\n');
}
```

## 参考資料

- [Design Tokens W3C Community Group](https://design-tokens.github.io/community-group/)
- [Ant Design Token System](https://ant.design/docs/react/customize-theme)
- [Material Design Tokens](https://m3.material.io/foundations/design-tokens)
