# Adlaire CI 生成 HTML デザイン仕様

**対象出力：** `adlaire-ci-build` が生成する静的 Web サイト HTML
**ビルドコンポーネント：** `components/builder.go` から生成する `adlaire-ci-build`
**デザインシステム：** [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)
**更新履歴：** 日付本文を正本化しない。デザイン変更の時系列は Git 履歴と Pull Request を正とする。

---

本ファイルは、生成 HTML のデザイン関係を定義するデザイン責務の正本である。生成 HTML のデザイン方針、視覚仕様、レイアウト、色、タイポグラフィ、TOC、コードブロック、トップへ戻るボタンの判断は本ファイルを参照する。

本ファイルは、機能仕様、運用仕様、API 仕様、状態 schema、実装状態、ロードマップ状態を定義しない。

Adlaire CI 全体の方針、ポリシー、正本参照先は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務、実装状態、実装可否、Phase、将来計画は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、builder の入出力、HTML / CSS / JavaScript 生成、theme component、検証条件は [`docs/details/builder.md`](details/builder.md) 詳細本文責務を参照する。

## 1. デザイン方針

docs.rs / MDN に倣った技術ドキュメントレイアウト。14,000 行超の仕様書を快適に閲覧するため、**構造の明快さ**と**情報密度への耐性**を最優先とする。

- ヘッダーのみアクセントカラーを使う。コンテンツ・サイドバーは中性色ベース
- CSS カスタムプロパティは [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)（`Tokens/`）定義の `--adlaire-*` トークンのみ使用
- **ライトモード固定**（`prefers-color-scheme` 非対応、ダークモードなし）
- 外部フォント不使用。システムフォントスタックで日本語環境の可読性を確保

---

## 2. カラートークン（ADS トークン）

| トークン | 値 | 用途 |
|---|---|---|
| `--adlaire-surface-accent` | `#0066cc` | ヘッダー背景 |
| `--adlaire-surface-accent-mid` | `#0055aa` | ホバー・二次アクセント |
| `--adlaire-surface-accent-strong` | `#004499` | コードテキスト色 |
| `--adlaire-surface-page` | `#f5f5f5` | ページ背景 |
| `--adlaire-surface-card` | `#ffffff` | カード・サイドバー背景 |
| `--adlaire-surface-soft` | `#f0f7ff` | ホバー背景・引用背景 |
| `--adlaire-surface-soft-strong` | `#e8f2ff` | コードブロック背景・インラインコード背景 |
| `--adlaire-surface-border` | `#e0e0e0` | ボーダー全般 |
| `--adlaire-surface-text` | `#333333` | 本文テキスト |
| `--adlaire-surface-text-muted` | `#555555` | 補助テキスト（TOC lv2 等） |
| `--adlaire-surface-text-subtle` | `#666666` | 三次テキスト（プレースホルダー等） |
| `--adlaire-color-primary` | `#0066cc` | リンク・アクセントボーダー |
| `--adlaire-color-secondary` | `#0055aa` | アクティブ TOC リンク色 |

---

## 3. タイポグラフィ

| トークン | 値 |
|---|---|
| `--adlaire-font-family-base` | `"Helvetica Neue", Helvetica, Arial, sans-serif` |
| `--adlaire-font-family-mono` | `"JetBrains Mono", "Courier New", Courier, monospace` |

外部フォント（Google Fonts 等）は使用しない。

### 見出し階層

| 要素 | サイズ | ウェイト | 装飾 |
|---|---|---|---|
| h1 | `--adlaire-font-size-2xl`（2rem） | 700 | 下線（`--adlaire-color-primary` 2px） |
| h2 | `--adlaire-font-size-xl`（1.5rem） | 600 | 下線（`--adlaire-surface-border` 1px）、上マージン `--adlaire-space-12` |
| h3 | `--adlaire-font-size-lg`（1.125rem） | 600 | 装飾なし |
| h4 | `--adlaire-font-size-sm`（0.875rem） | 500 | 左ボーダー（`--adlaire-color-primary` 3px）＋背景（`--adlaire-surface-soft`）、モノスペースフォント |

本文行長は `max-width: 68ch` でキャップ。

---

## 4. レイアウト

```
┌────────────────────────────────────────────────────────────┐
│  HEADER（固定・高さ 52px・--adlaire-surface-accent 背景）    │
├──────────┬─────────────────────────────────────────────────┤
│          │                                                 │
│ SIDEBAR  │  CONTENT                                        │
│  260px   │  max-width: 760px                               │
│  固定    │  padding: 3rem 3rem 96px                        │
│          │                                                 │
└──────────┴─────────────────────────────────────────────────┘
```

| 変数 | 値 | 用途 |
|---|---|---|
| `--hh` | `52px` | ヘッダー高さ |
| `--sw` | `var(--adlaire-layout-sidebar-compact)` = `260px` | サイドバー幅 |
| `--adlaire-layout-container-narrow` | `760px` | コンテンツ最大幅（`.ci`） |

### レスポンシブ

| ブレークポイント | 動作 |
|---|---|
| `≤ 768px` | サイドバーを画面外に収納（`86vw` 幅）、トグルで引き出し |
| `≤ 400px` | コンテンツ padding を縮小（`1.5rem 1rem`） |

---

## 5. サイドバー（TOC）

### 開閉制御

デスクトップとモバイルで制御方法が異なる。

| モード | 開 | 閉 |
|---|---|---|
| デスクトップ（> 768px） | `.closed` なし、`ct.style.marginLeft = 'var(--sw)'`（JS インライン） | `#sb.closed`、`ct.style.marginLeft = '0'`（JS インライン） |
| モバイル（≤ 768px） | `#sb.open`（`translateX(0)`） | `.open` なし（`translateX(-86vw)` で画面外） |

> **注意：** CSS の `#sb.closed ~ #ct { margin-left: 0 }` ルールは初期レンダリング時のみ機能する。以降の開閉操作はすべて JS のインラインスタイルが CSS クラスより優先して制御する。

`localStorage` キー `adb-sb`（`"1"` = 開、`"0"` = 閉）で開閉状態を永続化。

### TOC 構造

```
┌─ 検索ボックス ─────────────────────────────┐
│  [🔍 セクションを検索…]                     │
└─────────────────────────────────────────────┘
┌─ TOC リスト ────────────────────────────────┐
│  ▼ 1. 概要                    ← .tg（グループ）
│      1.1 プロジェクト概要      ← .ti（リーフ）
│  ▶ 2. 機能スコープ
└─────────────────────────────────────────────┘
```

### TOC リンクスタイル

| クラス | `padding-left` | フォントサイズ | 色 |
|---|---|---|---|
| `.lv1` | `14px` | `--adlaire-font-size-sm`（0.875rem） | `--adlaire-surface-text` |
| `.lv2` | `20px` | `0.8125rem` | `--adlaire-surface-text-muted` |
| `.lv3` | `34px` | `--adlaire-font-size-xs`（0.75rem） | `--adlaire-surface-text-subtle` |

アクティブ状態：`color: --adlaire-color-secondary`、`background: --adlaire-surface-soft`、`border-left: 2px solid --adlaire-color-primary`

---

## 6. コンポーネント

### コードブロック

```
┌─ .cb-wrap ──────────────────────────────────────────┐
│                         [rust] [コピー]  ← .cb-meta │
│  fn main() { … }                                    │
└─────────────────────────────────────────────────────┘
```

- 背景 `--adlaire-surface-soft-strong`、ボーダー `--adlaire-surface-border` 1px、角丸 `--adlaire-radius-lg`（8px）
- `.cb-meta`：`position: absolute; top: 8px; right: 10px`
- 言語ラベル（`.cl`）：モノフォント、`--adlaire-font-size-xs`、`--adlaire-surface-text-subtle`、大文字
- コピーボタン（`.cb-copy`）：通常 `opacity: 0`、ホバーで表示。クリック後「✓ 完了」→ 1.8 秒後に「コピー」へ復元
- `pre.cb`：フォントサイズ `--adlaire-font-size-sm`（0.875rem）、行高 `1.65`

### テーブル

- `.tw`（ラッパー）：`overflow-x: auto`、ボーダー・角丸 `--adlaire-radius-lg`・シャドウ
- `th`：モノフォント、`--adlaire-font-size-xs`、背景 `--adlaire-surface-soft`
- 偶数行：背景 `--adlaire-surface-soft`。ホバー行：背景 `--adlaire-surface-soft-strong`

### インラインコード

- モノフォント、サイズ `0.83em`、背景 `--adlaire-surface-soft-strong`、テキスト `--adlaire-surface-accent-strong`、ボーダー `--adlaire-surface-border` 1px、角丸 `--adlaire-radius-sm`（4px）

### 引用（blockquote）

- 左ボーダー `3px solid --adlaire-color-primary`、背景 `--adlaire-surface-soft`、角丸 `0 --adlaire-radius-lg --adlaire-radius-lg 0`、斜体テキスト

### 定義リスト

- `dt`：`font-weight: semibold`、`--adlaire-surface-text`
- `dd`：`margin-left: --adlaire-space-6`、`--adlaire-surface-text-muted`

---

## 7. トップへ戻るボタン

- 画面右下（`bottom: 28px / right: 24px`）に固定。`z-index: --adlaire-z-sticky`
- 直径 `38px` 円形、背景 `--adlaire-surface-accent`、色 `#fff`
- `scrollY > 400` で表示（`opacity: 0 → 1`、`pointer-events: none → auto`）
- ホバー時 `translateY(-2px)` で浮き上がり

---

## 8. ビルド方法

```bash
adlaire-ci-build --src <source.md> --out site/
```

入出力パス、既定値、終了コード、レポート出力は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務および [`docs/details/builder.md`](details/builder.md) builder 詳細本文責務を参照する。

CSS トークンの変更は Go 版 `components/builder.go` の HTML テンプレート内 `:root { }` ブロックに反映して再ビルドする。
