# Next.js ルーティング戦略

このドキュメントでは、Next.js App Routerを使用したルーティング設計について説明します。

## ルート構造

### 基本ルート

```
app/
├── page.tsx                 # / (ホーム)
├── layout.tsx              # ルートレイアウト
├── globals.css             # グローバルスタイル
├── not-found.tsx           # 404ページ
├── error.tsx               # エラーページ
├── loading.tsx             # ローディングページ
└── (authenticated)/        # 認証が必要なページ
    ├── layout.tsx          # 認証レイアウト
    └── dashboard/
        └── page.tsx        # /dashboard
```

### 動的ルート

```
app/
├── talent/
│   ├── page.tsx            # /talent (一覧)
│   ├── [talentId]/
│   │   ├── page.tsx        # /talent/[id] (詳細)
│   │   ├── edit/
│   │   │   └── page.tsx    # /talent/[id]/edit (編集)
│   │   └── layout.tsx      # タレント詳細レイアウト
│   └── layout.tsx          # タレントレイアウト
```

## レイアウトの設計

### ネストしたレイアウト

```typescript
// app/(authenticated)/layout.tsx
import { auth } from '@/libs/auth';

export default async function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await auth();

  if (!session) {
    redirect('/login');
  }

  return (
    <div className="flex h-screen">
      <Sidebar />
      <main className="flex-1">
        {children}
      </main>
    </div>
  );
}
```

### 条件付きレイアウト

```typescript
// app/talent/[talentId]/layout.tsx
import { getTalent } from '@/libs/api/talent';

export default async function TalentLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { talentId: string };
}) {
  const talent = await getTalent(params.talentId);

  if (!talent) {
    notFound();
  }

  return (
    <div>
      <TalentHeader talent={talent} />
      {children}
    </div>
  );
}
```

## ナビゲーション

### Link コンポーネントの使用

```typescript
import Link from 'next/link';

export function Navigation() {
  return (
    <nav>
      <Link href="/talent" className="nav-link">
        タレント一覧
      </Link>
      <Link href="/dashboard" className="nav-link">
        ダッシュボード
      </Link>
    </nav>
  );
}
```

### プログラム的なナビゲーション

```typescript
'use client';

import { useRouter } from 'next/navigation';

export function NavigationButton() {
  const router = useRouter();

  const handleNavigate = () => {
    router.push('/talent');
    // または
    router.replace('/talent');
    // または
    router.back();
  };

  return (
    <button onClick={handleNavigate}>
      ナビゲート
    </button>
  );
}
```

## パラメータとクエリ

### 動的パラメータ

```typescript
// app/talent/[talentId]/page.tsx
export default function TalentPage({
  params,
  searchParams,
}: {
  params: { talentId: string };
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  const talentId = params.talentId;
  const filter = searchParams.filter;

  return <TalentDetail talentId={talentId} filter={filter} />;
}
```

### クエリパラメータの検証

```typescript
import { z } from 'zod';

const searchParamsSchema = z.object({
  page: z.string().optional().default('1'),
  limit: z.string().optional().default('10'),
  filter: z.enum(['active', 'inactive']).optional(),
});

export default function TalentListPage({
  searchParams,
}: {
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  const { page, limit, filter } = searchParamsSchema.parse(searchParams);

  return <TalentList page={page} limit={limit} filter={filter} />;
}
```

## リダイレクト

```typescript
// app/old-page/page.tsx
import { redirect } from 'next/navigation';

export default function OldPage() {
  redirect('/new-page');
}
```

## ミドルウェア

### 認証ミドルウェア

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { auth } from '@/libs/auth';

export async function middleware(request: NextRequest) {
  const session = await auth();

  // 認証が必要なページ
  if (request.nextUrl.pathname.startsWith('/dashboard')) {
    if (!session) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  // 認証済みユーザーのリダイレクト
  if (request.nextUrl.pathname === '/login' && session) {
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/login'],
};
```

## エラーハンドリング

### 404エラー

```typescript
// app/not-found.tsx
import Link from 'next/link';

export default function NotFound() {
  return (
    <div>
      <h2>ページが見つかりません</h2>
      <Link href="/">ホームに戻る</Link>
    </div>
  );
}
```

### 動的ルートの404

```typescript
// app/talent/[talentId]/not-found.tsx
export default function TalentNotFound() {
  return (
    <div>
      <h2>タレントが見つかりません</h2>
      <p>指定されたタレントは存在しないか、削除されています。</p>
    </div>
  );
}
```

## パフォーマンス最適化

### プリフェッチ

```typescript
import Link from 'next/link';

export function TalentList({ talents }: { talents: Talent[] }) {
  return (
    <div>
      {talents.map((talent) => (
        <Link
          key={talent.id}
          href={`/talent/${talent.id}`}
          prefetch={true} // デフォルトで有効
        >
          {talent.name}
        </Link>
      ))}
    </div>
  );
}
```

### 動的インポート

```typescript
import dynamic from 'next/dynamic';

const HeavyChart = dynamic(() => import('./HeavyChart'), {
  loading: () => <div>チャート読み込み中...</div>,
  ssr: false,
});
```

## 参考資料

- [Next.js App Router Documentation](https://nextjs.org/docs/app)
- [Routing and Navigation](https://nextjs.org/docs/app/building-your-application/routing)
- [Middleware](https://nextjs.org/docs/app/building-your-application/routing/middleware)
