# useEffectの正しい使い方ルール

## 基本原則

### 1. ReactはUIライブラリである
- ReactにはUIの管理をさせるべきであり、その他のことはReactの役目ではない
- アプリケーションのロジックの管理やそれに付随するステートの管理はReactの外部で処理するべき
- そのためにはRecoilなどのライブラリが有用

### 2. Reactはコンポーネントベースのライブラリである
- コンポーネント内に書かれるコードは、useEffectも含めてすべてそのコンポーネントのロジックでなければならない
- そのコンポーネントに閉じないロジックを書くのは良くない

### 3. クリーンアップ関数の必須化
- **クリーンアップ関数の無いuseEffectは不適格**
- コンポーネントがアンマウントされたときにそのコンポーネントの影響が元に戻されないため、コンポーネントベースの原則から外れる

## ルール

### ✅ 許容されるuseEffectの使い方

#### 1. イベントハンドラを登録する系
**許容度:** 😃 文句なし。望ましいuseEffectの使い方

```tsx
const Component = () => {
  const [h, setH] = useState(0);

  useEffect(() => {
    const handler = () => {
      setH(h => (h + 1) % 360);
    };
    window.addEventListener("pointermove", handler);
    return () => {
      window.removeEventListener("pointermove", handler);
    };
  }, []);

  return (
    <div
      style={{
        backgroundColor: `hsl(${h}, 100%, 50%)`,
        width: "100px",
        height: "100px"
      }}
    />
  );
};
```

**理由:**
- Reactの他の機能では賄えないようなDOM操作が必要
- 高頻度で発生するイベントを処理する場合のエスケープハッチ
- 適切なクリーンアップ関数が実装されている

#### 2. データを取得する系（限定的）
**許容度:** 🙃 場合により許容できるケースもある。しかし、より良い代替手段がそのうち登場するので将来性は乏しい

```tsx
const Component = () => {
  const [user, setUser] = useState<null | { login: string }>(null);

  useEffect(() => {
    const abortController = new AbortController();
    fetch("https://api.github.com/users/uhyo", {
      signal: abortController.signal
    })
      .then(res => res.json())
      .then(user => setUser(user))
      .catch(reportError);

    return () => {
      abortController.abort();
    };
  }, []);

  return <div>{user?.login}</div>;
};
```

**理由:**
- 適切なクリーンアップ関数（AbortController）が実装されている
- ただし、より良い代替手段（React Query、SWR等）の使用を推奨

### ❌ 避けるべきuseEffectの使い方

#### 1. 値の変化に反応するためのuseEffect
**許容度:** 😡 不適格

```tsx
// ❌ 悪い例
const Component = ({ value }) => {
  const [processedValue, setProcessedValue] = useState(value);

  useEffect(() => {
    setProcessedValue(processValue(value));
  }, [value]);

  return <div>{processedValue}</div>;
};

// ✅ 良い例
const Component = ({ value }) => {
  const processedValue = useMemo(() => processValue(value), [value]);
  return <div>{processedValue}</div>;
};
```

#### 2. 依存配列のlintエラーを無視するuseEffect
**許容度:** 😡 不適格

```tsx
// ❌ 悪い例
useEffect(() => {
  // 何かの処理
}, []); // eslint-disable-line react-hooks/exhaustive-deps

// ✅ 良い例
useEffect(() => {
  // 何かの処理
}, [dependency1, dependency2]);
```

#### 3. クリーンアップ関数のないuseEffect
**許容度:** 😡 不適格

```tsx
// ❌ 悪い例
useEffect(() => {
  // 何かの処理（クリーンアップなし）
}, []);

// ✅ 良い例
useEffect(() => {
  // 何かの処理
  return () => {
    // クリーンアップ処理
  };
}, []);
```

## 実装ガイドライン

### 1. useEffectの依存配列について
- 依存配列は最適化のためのものであり、最適化を超えた意味を持った依存配列は不適格
- 値の変化に反応するためにuseEffectを使うのは良くない
- 依存配列のlintエラーを無視するのはご法度

### 2. コンポーネントロジックの原則
- useEffectはコンポーネントロジックの一部である
- useEffectは"副作用"のためのものではない
- そのコンポーネントに閉じないロジックを書くのは良くない

### 3. UI管理の原則
- useEffectはUIの管理という目的のために使う
- Reactの他の機能では賄えないようなDOM操作をしたい場合に使用する

## チェックリスト

### useEffectを書く前に確認すること
- [ ] この処理はUIの管理のために必要か？
- [ ] この処理はコンポーネントロジックの一部か？
- [ ] クリーンアップ関数は必要か？必要なら実装されているか？
- [ ] 依存配列は適切に設定されているか？
- [ ] 値の変化に反応するためにuseEffectを使っていないか？

### コードレビュー時の確認事項
- [ ] クリーンアップ関数が実装されているか？
- [ ] 依存配列のlintエラーがないか？
- [ ] コンポーネントロジックの範囲内で処理されているか？
- [ ] UI管理の目的で使用されているか？

## 参考

- [過激派が教える！ useEffectの正しい使い方](https://zenn.dev/uhyo/articles/useeffect-taught-by-extremist)
- [React公式ドキュメント - useEffect](https://react.dev/reference/react/useEffect)
