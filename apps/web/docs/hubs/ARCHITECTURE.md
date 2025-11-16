# アーキテクチャ

このドキュメントでは、Webアプリケーションのアーキテクチャ設計について説明します。

## ディレクトリ構造

- [ディレクトリ構造](./DIRECTORY_STRUCTURE.md) - プロジェクト全体のディレクトリ設計

## コンポーネント設計

- [Container/Presentational パターン](../spokes/nextjs/CONTAINER_PRESENTATION_ARCHITECTURE.md) - コンポーネントの責務分離
- [コンポーネント設計](../spokes/directories-and-files/COMPONENTS.md) - 再利用可能なUIコンポーネント
- [Container設計](../spokes/directories-and-files/CONTAINER.md) - ロジック層の設計
- [Presenter設計](../spokes/directories-and-files/PRESENTER.md) - UI層の設計

## 状態管理

- [カスタムフック](../spokes/directories-and-files/HOOKS.md) - ロジックの再利用
- [Server Actions](../spokes/react/SERVER_ACTIONS.md) - サーバーサイド処理
- [useEffect ルール](../spokes/react/USEEFFECT_RULES.md) - 副作用の管理

## データフロー

- [Features設計](../spokes/directories-and-files/FEATURES.md) - ドメイン固有の機能
- [Actions設計](../spokes/directories-and-files/ACTIONS.md) - サーバーアクション
- [型定義](../spokes/typescript/TYPESCRIPT.md) - 型安全性の確保

## ユーティリティ

- [ライブラリ設計](../spokes/directories-and-files/LIBS.md) - 共通ユーティリティ
- [Utils設計](../spokes/directories-and-files/UTILS.md) - ヘルパー関数
