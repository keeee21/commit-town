# Client Componentでの "use client" ディレクティブの必須化ルール

## ルール

以下のような場合は"use client" ディレクティブの利用を強制する。

- Hooks を使用する場合
- イベントハンドラー（onClick, onChange等）を使用する場合
- ブラウザAPI（localStorage, window等）を使用する場合
- 状態管理（useState, useReducer等）を使用する場合
- ライフサイクル（useEffect, useLayoutEffect等）を使用する場合
- フォームの制御（controlled components）を行う場合

## 背景

### RSC Payload転送量の削減

Server Componentでクライアント側の機能を使用すると、以下の問題が発生する：

1. **RSC Payloadの増大**: Server Componentでクライアント側の機能を使用すると、Next.jsが自動的にClient Componentに変換する際に、RSC Payloadが大幅に増加する
2. **レンダリング性能の低下**: 大きなRSC Payloadは、初期レンダリング時の転送時間を増加させ、ユーザー体験を悪化させる
3. **メモリ使用量の増加**: クライアント側で不要なデータが保持され、メモリ効率が悪化する

### 具体的な影響

- **Hooks使用時**: Server ComponentでHooksを使用すると、Next.jsが自動的にClient Componentに変換する際に、コンポーネント全体がクライアント側で実行される
- **イベントハンドラー使用時**: イベントハンドラーを含むコンポーネントは、必然的にクライアント側で実行される必要がある
- **ブラウザAPI使用時**: localStorageやwindowオブジェクトなどは、サーバー側では利用できないため、適切なClient Component化が必要

## 実装ガイドライン

### 必須パターン

```tsx
// ✅ 正しい実装
"use client";

import { useState, useEffect } from 'react';

export function InteractiveComponent() {
  const [count, setCount] = useState(0);

  useEffect(() => {
    // クライアント側での処理
  }, []);

  return (
    <button onClick={() => setCount(count + 1)}>
      Count: {count}
    </button>
  );
}
```

### 避けるべきパターン

```tsx
// ❌ 避けるべき実装（"use client"なし）
import { useState } from 'react';

export function BadComponent() {
  const [count, setCount] = useState(0); // エラーが発生する可能性

  return <div>{count}</div>;
}
```

### 最適化のポイント

1. **最小限のClient Component化**: 必要な部分のみをClient Componentに分離
2. **Server Componentの活用**: 静的な部分は可能な限りServer Componentで実装
3. **適切な境界の設定**: インタラクティブな部分と静的な部分を明確に分離

## 検証方法

### 開発時の確認

1. ブラウザの開発者ツールでNetworkタブを確認
2. RSC Payloadのサイズを監視
3. 不要なClient Component化がないかチェック

### パフォーマンス測定

```tsx
// パフォーマンス測定の例
"use client";

import { useEffect } from 'react';

export function PerformanceMonitor() {
  useEffect(() => {
    // RSC Payloadサイズの測定
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (entry.name.includes('rsc')) {
          console.log('RSC Payload size:', entry.transferSize);
        }
      }
    });

    observer.observe({ entryTypes: ['resource'] });

    return () => observer.disconnect();
  }, []);

  return null;
}
```

## 参考

- [Next.jsの考え方 - Client Componentsのユースケース](https://zenn.dev/akfm/books/nextjs-basic-principle/viewer/part_2_client_components_usecase#rsc-payload%E8%BB%A2%E9%80%81%E9%87%8F%E3%81%AE%E5%89%8A%E6%B8%9B)
- [Next.js公式ドキュメント - Client Components](https://nextjs.org/docs/app/building-your-application/rendering/client-components)
