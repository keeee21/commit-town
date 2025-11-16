# Components ディレクトリ構造とルール

## ディレクトリ構成

```
src/components/
├── ui/           # 原子的なUIコンポーネント
├── elements/     # 分子的なコンポーネント
├── composites/   # 複合的なコンポーネント
├── layouts/      # レイアウトコンポーネント
└── demo/         # デモ・サンプルコンポーネント
```

## 各ディレクトリの役割とルール

### ui/

**役割**: これ以上分割できない原子的なUIコンポーネント

**特徴**:
- スタイルの記述が含まれている
- 再利用可能な最小単位のコンポーネント
- 外部UIライブラリのコンポーネントを直接使用

**ルール**:
- **Client Rendering必須**: スタイリングが含まれるため、`"use client"`ディレクティブを必須とする
- 詳細は [@ENFORCE_USE_CLIENT_DIRECTIVE.md](../react/ENFORCE_USE_CLIENT_DIRECTIVE.md) を参照
- ビジネスロジックを含まない
- プロップスは最小限に留める

**例**:
```tsx
"use client";

import { Button } from "@/components/ui/button";

export function CustomButton({ children, variant, ...props }) {
  return (
    <Button
      className="custom-styles"
      variant={variant}
      {...props}
    >
      {children}
    </Button>
  );
}
```

### elements/

**役割**: ui/のコンポーネントや外部UIライブラリのコンポーネントをラップした分子的なコンポーネント

**特徴**:
- ui/コンポーネントを組み合わせて作成
- 基本的なスタイリングが含まれる
- 単一の責任を持つ

**ルール**:
- **Client Rendering推奨**: スタイリングが発生するため、基本的には`"use client"`ディレクティブを使用
- ビジネスロジックは最小限
- 再利用可能な単位として設計
- 型定義は [@TYPESCRIPT.md](../../../TYPESCRIPT.md) に従う

**例**:
```tsx
"use client";

import { CustomButton } from "@/components/ui/custom-button";
import { CustomInput } from "@/components/ui/custom-input";

export function SearchForm({ onSearch, placeholder }) {
  return (
    <form className="search-form-styles">
      <CustomInput placeholder={placeholder} />
      <CustomButton type="submit" onClick={onSearch}>
        検索
      </CustomButton>
    </form>
  );
}
```

### composites/

**役割**: elements/を複数呼び出した上位レイヤーのコンポーネント

**特徴**:
- elements/コンポーネントを組み合わせて作成
- より複雑な機能を持つ
- ページ固有のロジックを含む場合がある

**ルール**:
- Server ComponentまたはClient Componentの使い分けを適切に行う
- ビジネスロジックを含む場合は適切に分離
- 状態管理が必要な場合は`"use client"`ディレクティブを使用
- 型定義は [@TYPESCRIPT.md](../../../TYPESCRIPT.md) に従う

**例**:
```tsx
import { SearchForm } from "@/components/elements/search-form";
import { FilterPanel } from "@/components/elements/filter-panel";
import { DataTable } from "@/components/elements/data-table";

export function UserManagementPage() {
  return (
    <div className="user-management-layout">
      <SearchForm onSearch={handleSearch} />
      <FilterPanel onFilterChange={handleFilter} />
      <DataTable data={users} />
    </div>
  );
}
```

## 設計原則

### 1. 単一責任の原則
各コンポーネントは一つの明確な責任を持つ

### 2. 再利用性の最大化
- ui/ → プロジェクト全体で再利用
- elements/ → 機能単位で再利用
- composites/ → ページ単位で再利用

### 3. 適切な抽象化レベル
- ui/ → 最低レベルの抽象化
- elements/ → 中レベルの抽象化
- composites/ → 高レベルの抽象化

### 4. 型安全性の確保
- すべてのコンポーネントで適切な型定義を行う
- 型定義の詳細は [@TYPESCRIPT.md](../../../TYPESCRIPT.md) を参照

## 命名規則

### ファイル名
- PascalCaseを使用
- コンポーネント名と一致させる

### コンポーネント名
- ディレクトリ名と一致させる
- 明確で説明的な名前を使用

### 例
```
ui/
├── custom-button.tsx     # CustomButton
├── custom-input.tsx      # CustomInput
└── loading-spinner.tsx   # LoadingSpinner

elements/
├── search-form.tsx       # SearchForm
├── filter-panel.tsx      # FilterPanel
└── data-table.tsx        # DataTable

composites/
├── user-management-page.tsx  # UserManagementPage
├── product-catalog-page.tsx  # ProductCatalogPage
└── dashboard-page.tsx        # DashboardPage
```

## 参考

- [@ENFORCE_USE_CLIENT_DIRECTIVE.md](../react/ENFORCE_USE_CLIENT_DIRECTIVE.md) - Client Componentのルール
- [@TYPESCRIPT.md](../../../TYPESCRIPT.md) - 型定義のガイドライン
