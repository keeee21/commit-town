# Web ドキュメント

このディレクトリには、Webアプリケーションの開発・運用に関するドキュメントが含まれています。

## ドキュメント構造

このドキュメントは **hubs & spokes** アーキテクチャに基づいて整理されています：

- **hubs**: 主要なトピックの概要と関連ドキュメントへのリンク
- **spokes**: 具体的な実装詳細や技術仕様

## 主要ドキュメント (Hubs)

### 🏗️ [アーキテクチャ](./hubs/ARCHITECTURE.md)
システム設計、コンポーネント設計、データフローについて

### 🛠️ [技術スタック](./hubs/TECHNOLOGY_STACK.md)
使用している技術、ライブラリ、フレームワークについて

### 📁 [ディレクトリ構造](./hubs/DIRECTORY_STRUCTURE.md)
プロジェクトのディレクトリ設計と各ディレクトリの役割

### 💻 [開発ワークフロー](./hubs/DEVELOPMENT.md)
開発プロセス、コーディング規約、開発環境について

### 🧪 [テスト戦略](./hubs/TESTING.md)
テスト実装、テスト環境、テスト実行について

### 🚀 [デプロイメント](./hubs/DEPLOYMENT.md)
デプロイ戦略、環境設定、パフォーマンス最適化について

## 詳細ドキュメント (Spokes)

### ディレクトリ・ファイル設計
- [Actions](./spokes/directories-and-files/ACTIONS.md) - Server Actions設計
- [Components](./spokes/directories-and-files/COMPONENTS.md) - UIコンポーネント設計
- [Container](./spokes/directories-and-files/CONTAINER.md) - ロジック層設計
- [Features](./spokes/directories-and-files/FEATURES.md) - ドメイン機能設計
- [Hooks](./spokes/directories-and-files/HOOKS.md) - カスタムフック設計
- [Libs](./spokes/directories-and-files/LIBS.md) - ライブラリ設計
- [Presenter](./spokes/directories-and-files/PRESENTER.md) - UI層設計
- [Storybook](./spokes/directories-and-files/STORYBOOK.md) - コンポーネント開発環境
- [Utils](./spokes/directories-and-files/UTILS.md) - ユーティリティ設計

### Next.js
- [Container/Presentational アーキテクチャ](./spokes/nextjs/CONTAINER_PRESENTATION_ARCHITECTURE.md) - コンポーネント設計パターン

### React
- [useClient ディレクティブ](./spokes/react/ENFORCE_USE_CLIENT_DIRECTIVE.md) - クライアントコンポーネント
- [JSDocコメント](./spokes/react/JS_DOC_COMMENT_FOR_UI.md) - UIドキュメント化
- [Post Hooks](./spokes/react/POST_HOOKS.md) - カスタムフック設計
- [Server Actions](./spokes/react/SERVER_ACTIONS.md) - サーバーアクション
- [useEffect ルール](./spokes/react/USEEFFECT_RULES.md) - 副作用管理

### Storybook
- [自動ドキュメント化](./spokes/storybook/ENFORCE_AUTO_DOCS_FOR_STORYBOOK.md) - コンポーネントドキュメント

### TypeScript
- [TypeScript設定](./spokes/typescript/TYPESCRIPT.md) - 型定義と設定

## 使用方法

1. **概要を理解したい場合**: 対応するhubファイルを参照
2. **具体的な実装を知りたい場合**: hubファイルから関連するspokeファイルを参照
3. **特定の技術について詳しく知りたい場合**: 直接spokeファイルを参照

## AI コンテキスト最適化

このドキュメント構造により、AIに渡すコンテキストを効率的に管理できます：

- **hubファイル**: 全体像と関連ドキュメントの概要
- **spokeファイル**: 具体的な実装詳細と技術仕様

必要に応じて、関連するhubファイルとspokeファイルを組み合わせて使用してください。
