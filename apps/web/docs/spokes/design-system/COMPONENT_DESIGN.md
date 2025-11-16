# コンポーネント設計システム

このドキュメントでは、Webアプリケーションのコンポーネント設計システムについて説明します。

## コンポーネント階層

### 1. コンポーネント分類

```
components/
├── elements/          # 基本UI要素
│   ├── orion-button/
│   ├── orion-input/
│   └── orion-modal/
├── composites/        # 複合コンポーネント
│   ├── orion-float-drawer/
│   └── orion-form-item/
├── layouts/          # レイアウトコンポーネント
│   ├── header/
│   ├── sidebar/
│   └── global-navigation/
└── ui/               # ユーティリティコンポーネント
    └── scaled-icon/
```

### 2. コンポーネント命名規則

```typescript
// コンポーネント名の規則
// 形式: orion-[component-name]
// 例: orion-button, orion-input, orion-modal

// ファイル構造
orion-button/
├── orion-button.tsx           # メインコンポーネント
├── orion-button.stories.tsx   # Storybookストーリー
├── orion-button.test.tsx      # テストファイル
├── index.ts                   # エクスポート
└── types.ts                   # 型定義（必要に応じて）
```

## 基本コンポーネント設計

### 1. ボタンコンポーネント

```typescript
// orion-button.tsx
"use client";

import React from 'react';
import { Button as AntButton, ButtonProps as AntButtonProps } from 'antd';
import { cn } from '@/libs/cn';

export type OrionButtonProps = Omit<AntButtonProps, 'type'> & {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'small' | 'medium' | 'large';
  loading?: boolean;
  disabled?: boolean;
  children: React.ReactNode;
};

/**
 * Orionボタンコンポーネント
 * 
 * @description 統一されたスタイルとバリアントを持つボタンコンポーネント
 * @example
 * ```tsx
 * <OrionButton variant="primary" size="medium">
 *   クリック
 * </OrionButton>
 * ```
 */
export const OrionButton: React.FC<OrionButtonProps> = ({
  variant = 'primary',
  size = 'medium',
  loading = false,
  disabled = false,
  className,
  children,
  ...props
}) => {
  const baseClasses = 'orion-button';
  const variantClasses = {
    primary: 'orion-button--primary',
    secondary: 'orion-button--secondary',
    outline: 'orion-button--outline',
    ghost: 'orion-button--ghost',
    danger: 'orion-button--danger',
  };
  const sizeClasses = {
    small: 'orion-button--small',
    medium: 'orion-button--medium',
    large: 'orion-button--large',
  };

  return (
    <AntButton
      type={variant === 'primary' ? 'primary' : 'default'}
      size={size}
      loading={loading}
      disabled={disabled}
      className={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        className
      )}
      {...props}
    >
      {children}
    </AntButton>
  );
};
```

### 2. 入力コンポーネント

```typescript
// orion-input.tsx
"use client";

import React from 'react';
import { Input as AntInput, InputProps as AntInputProps } from 'antd';
import { cn } from '@/libs/cn';

export type OrionInputProps = Omit<AntInputProps, 'size'> & {
  size?: 'small' | 'medium' | 'large';
  error?: boolean;
  helperText?: string;
  label?: string;
  required?: boolean;
};

/**
 * Orion入力コンポーネント
 * 
 * @description 統一されたスタイルとバリデーション機能を持つ入力コンポーネント
 * @example
 * ```tsx
 * <OrionInput 
 *   label="ユーザー名" 
 *   placeholder="名前を入力してください"
 *   required 
 * />
 * ```
 */
export const OrionInput: React.FC<OrionInputProps> = ({
  size = 'medium',
  error = false,
  helperText,
  label,
  required = false,
  className,
  ...props
}) => {
  const baseClasses = 'orion-input';
  const sizeClasses = {
    small: 'orion-input--small',
    medium: 'orion-input--medium',
    large: 'orion-input--large',
  };

  return (
    <div className="orion-input-wrapper">
      {label && (
        <label className="orion-input-label">
          {label}
          {required && <span className="orion-input-required">*</span>}
        </label>
      )}
      <AntInput
        size={size}
        status={error ? 'error' : undefined}
        className={cn(
          baseClasses,
          sizeClasses[size],
          error && 'orion-input--error',
          className
        )}
        {...props}
      />
      {helperText && (
        <div className={cn(
          'orion-input-helper',
          error && 'orion-input-helper--error'
        )}>
          {helperText}
        </div>
      )}
    </div>
  );
};
```

## 複合コンポーネント設計

### 1. フォームアイテムコンポーネント

```typescript
// orion-form-item.tsx
import React from 'react';
import { Form, FormItemProps as AntFormItemProps } from 'antd';
import { cn } from '@/libs/cn';

export type OrionFormItemProps = AntFormItemProps & {
  children: React.ReactNode;
  className?: string;
};

export function OrionFormItem({
  children,
  className,
  ...props
}: OrionFormItemProps) {
  return (
    <Form.Item
      className={cn('orion-form-item', className)}
      {...props}
    >
      {children}
    </Form.Item>
  );
}

// サブコンポーネント
export function OrionFormField({ children }: { children: React.ReactNode }) {
  return <div className="orion-form-field">{children}</div>;
}

export function OrionFormLabel({ children }: { children: React.ReactNode }) {
  return <div className="orion-form-label">{children}</div>;
}

export function OrionFormValidation({
  error
}: {
  error?: string
}) {
  if (!error) return null;

  return <div className="orion-form-validation">{error}</div>;
}
```

### 2. フロートドロワーコンポーネント

```typescript
// orion-float-drawer.tsx
"use client";

import React, { useState } from 'react';
import { Drawer, DrawerProps as AntDrawerProps } from 'antd';
import { cn } from '@/libs/cn';

export type OrionFloatDrawerProps = Omit<AntDrawerProps, 'open'> & {
  children: React.ReactNode;
  trigger?: React.ReactNode;
  className?: string;
};

/**
 * Orionフロートドロワーコンポーネント
 * 
 * @description トリガー要素をクリックして開くドロワーコンポーネント
 * @example
 * ```tsx
 * <OrionFloatDrawer 
 *   trigger={<Button>メニューを開く</Button>}
 *   title="設定"
 * >
 *   <p>ドロワーの内容</p>
 * </OrionFloatDrawer>
 * ```
 */
export const OrionFloatDrawer: React.FC<OrionFloatDrawerProps> = ({
  children,
  trigger,
  className,
  ...props
}) => {
  const [open, setOpen] = useState(false);

  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  return (
    <>
      {trigger && (
        <div onClick={handleOpen} className="orion-float-drawer-trigger">
          {trigger}
        </div>
      )}
      <Drawer
        open={open}
        onClose={handleClose}
        className={cn('orion-float-drawer', className)}
        {...props}
      >
        {children}
      </Drawer>
    </>
  );
};
```

## レイアウトコンポーネント設計

### 1. ヘッダーコンポーネント

```typescript
// header.tsx
"use client";

import React from 'react';
import { Layout, Avatar, Dropdown, Button } from 'antd';
import { UserOutlined, SettingOutlined } from '@ant-design/icons';
import { cn } from '@/libs/cn';

export type HeaderProps = {
  user?: {
    name: string;
    avatar?: string;
  };
  onLogout?: () => void;
  onSettings?: () => void;
  className?: string;
};

/**
 * Headerコンポーネント
 * 
 * @description アプリケーションのヘッダー部分を表示するコンポーネント
 * @example
 * ```tsx
 * <Header 
 *   user={{ name: "田中太郎", avatar: "/avatar.jpg" }}
 *   onLogout={() => console.log('ログアウト')}
 *   onSettings={() => console.log('設定')}
 * />
 * ```
 */
export const Header: React.FC<HeaderProps> = ({
  user,
  onLogout,
  onSettings,
  className
}) => {
  const menuItems = [
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '設定',
      onClick: onSettings,
    },
    {
      key: 'logout',
      label: 'ログアウト',
      onClick: onLogout,
    },
  ];

  return (
    <Layout.Header className={cn('orion-header', className)}>
      <div className="orion-header-content">
        <div className="orion-header-logo">
          <h1>System</h1>
        </div>
        <div className="orion-header-actions">
          {user && (
            <Dropdown
              menu={{ items: menuItems }}
              placement="bottomRight"
            >
              <Button type="text" className="orion-header-user">
                <Avatar
                  size="small"
                  src={user.avatar}
                  icon={<UserOutlined />}
                />
                <span>{user.name}</span>
              </Button>
            </Dropdown>
          )}
        </div>
      </div>
    </Layout.Header>
  );
};
```

### 2. サイドバーコンポーネント

```typescript
// sidebar.tsx
"use client";

import React from 'react';
import { Layout, Menu } from 'antd';
import { cn } from '@/libs/cn';

export type SidebarProps = {
  items: Array<{
    key: string;
    label: string;
    icon?: React.ReactNode;
    children?: Array<{
      key: string;
      label: string;
    }>;
  }>;
  selectedKeys?: string[];
  onSelect?: (key: string) => void;
  className?: string;
};

/**
 * Sidebarコンポーネント
 * 
 * @description アプリケーションのサイドバーナビゲーション
 * @example
 * ```tsx
 * <Sidebar 
 *   items={[
 *     { key: 'dashboard', label: 'ダッシュボード', icon: <DashboardOutlined /> },
 *     { key: 'users', label: 'ユーザー管理', icon: <UserOutlined /> }
 *   ]}
 *   selectedKeys={['dashboard']}
 *   onSelect={(key) => console.log(key)}
 * />
 * ```
 */
export const Sidebar: React.FC<SidebarProps> = ({
  items,
  selectedKeys = [],
  onSelect,
  className
}) => {
  return (
    <Layout.Sider
      width={240}
      className={cn('orion-sidebar', className)}
    >
      <Menu
        mode="inline"
        selectedKeys={selectedKeys}
        items={items}
        onSelect={({ key }) => onSelect?.(key)}
        className="orion-sidebar-menu"
      />
    </Layout.Sider>
  );
};
```

## スタイリング戦略

### 1. CSS Modules

```css
/* orion-button.module.css */
.orion-button {
  @apply inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors;
}

.orion-button--primary {
  @apply bg-blue-600 text-white hover:bg-blue-700;
}

.orion-button--secondary {
  @apply bg-gray-200 text-gray-900 hover:bg-gray-300;
}

.orion-button--outline {
  @apply border border-gray-300 bg-transparent hover:bg-gray-50;
}

.orion-button--ghost {
  @apply bg-transparent hover:bg-gray-100;
}

.orion-button--danger {
  @apply bg-red-600 text-white hover:bg-red-700;
}

.orion-button--small {
  @apply h-8 px-3 text-xs;
}

.orion-button--medium {
  @apply h-10 px-4 text-sm;
}

.orion-button--large {
  @apply h-12 px-6 text-base;
}
```

### 2. Tailwind CSSクラス

```typescript
// ユーティリティ関数
export function getButtonClasses(variant: string, size: string) {
  const baseClasses = 'inline-flex items-center justify-center rounded-md font-medium transition-colors';

  const variantClasses = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700',
    secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300',
    outline: 'border border-gray-300 bg-transparent hover:bg-gray-50',
    ghost: 'bg-transparent hover:bg-gray-100',
    danger: 'bg-red-600 text-white hover:bg-red-700',
  };

  const sizeClasses = {
    small: 'h-8 px-3 text-xs',
    medium: 'h-10 px-4 text-sm',
    large: 'h-12 px-6 text-base',
  };

  return `${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]}`;
}
```

## アクセシビリティ

### 1. ARIA属性の追加

```typescript
// アクセシビリティ対応のボタン
export const AccessibleButton: React.FC<OrionButtonProps & {
  ariaLabel?: string;
  ariaDescribedBy?: string;
}> = ({
  children,
  ariaLabel,
  ariaDescribedBy,
  ...props
}) => {
  return (
    <OrionButton
      aria-label={ariaLabel}
      aria-describedby={ariaDescribedBy}
      {...props}
    >
      {children}
    </OrionButton>
  );
};
```

### 2. キーボードナビゲーション

```typescript
// キーボードナビゲーション対応
export const KeyboardNavigableButton: React.FC<OrionButtonProps & {
  onKeyDown?: (event: React.KeyboardEvent) => void;
}> = ({
  onKeyDown,
  ...props
}) => {
  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      props.onClick?.(event as any);
    }
    onKeyDown?.(event);
  };

  return (
    <OrionButton
      onKeyDown={handleKeyDown}
      {...props}
    />
  );
};
```

## テスト戦略

### 1. コンポーネントテスト

```typescript
// orion-button.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { OrionButton } from './orion-button';

describe('OrionButton', () => {
  it('renders with correct text', () => {
    render(<OrionButton>Click me</OrionButton>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const handleClick = jest.fn();
    render(<OrionButton onClick={handleClick}>Click me</OrionButton>);

    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('applies correct variant classes', () => {
    render(<OrionButton variant="primary">Primary</OrionButton>);
    expect(screen.getByText('Primary')).toHaveClass('orion-button--primary');
  });
});
```

### 2. Storybookストーリー

```typescript
// orion-button.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { OrionButton } from './orion-button';

const meta: Meta<typeof OrionButton> = {
  title: 'Elements/OrionButton',
  component: OrionButton,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: { type: 'select' },
      options: ['primary', 'secondary', 'outline', 'ghost', 'danger'],
    },
    size: {
      control: { type: 'select' },
      options: ['small', 'medium', 'large'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    variant: 'primary',
    children: 'Primary Button',
  },
};

export const Secondary: Story = {
  args: {
    variant: 'secondary',
    children: 'Secondary Button',
  },
};

export const AllVariants: Story = {
  render: () => (
    <div className="flex gap-4">
      <OrionButton variant="primary">Primary</OrionButton>
      <OrionButton variant="secondary">Secondary</OrionButton>
      <OrionButton variant="outline">Outline</OrionButton>
      <OrionButton variant="ghost">Ghost</OrionButton>
      <OrionButton variant="danger">Danger</OrionButton>
    </div>
  ),
};
```

## 参考資料

- [Ant Design Components](https://ant.design/components/overview)
- [Material Design Components](https://m3.material.io/components)
- [Chakra UI Components](https://chakra-ui.com/docs/components)
