# Adlaire CI Web デザイン仕様

**対象 surface：** `adlaire-ci-build` が生成する静的 Web サイトと [`admin/index.html`](../admin/index.html) の標準管理 UI
**実装 artifact：** 生成静的 Web サイトは [`components/builder.go`](../components/builder.go)、標準管理 UI は [`admin/index.html`](../admin/index.html)
**デザインシステム：** [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)
**更新履歴：** 日付本文を正本化しない。デザイン変更の時系列は Git 履歴と Pull Request を正とする。

---

[`docs/DESIGN.md`](DESIGN.md) デザイン責務は、生成静的 Web サイトと標準管理 UI のデザイン関係を定義する正本である。生成静的 Web サイトのデザイン方針、視覚仕様、レイアウト、色、タイポグラフィ、TOC、コードブロック、トップへ戻るボタンは [§1〜§8](#1-デザイン方針)、標準管理 UI の視覚仕様は [標準管理 UI 視覚契約](#admin-ui-visual-contract) を確認する。

| surface | デザイン正本範囲 | 挙動の参照先 |
|---------|------------------|--------------|
| 生成静的 Web サイト | [§1〜§8](#1-デザイン方針) の token、selector、レイアウト、responsive、print。 | [`docs/details/builder.md`](details/builder.md) |
| 標準管理 UI | [標準管理 UI 視覚契約](#admin-ui-visual-contract) の token、selector、レイアウト、responsive。 | [`docs/details/ui.md`](details/ui.md) |

デザイン外の正本参照先は [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務を参照する。

## 1. デザイン方針

[`docs/DESIGN.md` デザイン責務 §1](DESIGN.md#1-デザイン方針)〜[§8](DESIGN.md#8-生成側参照) は生成静的 Web サイトへ適用する。docs.rs / MDN に倣った技術ドキュメントレイアウトとし、14,000 行超の仕様書を快適に閲覧するため、**構造の明快さ**と**情報密度への耐性**を最優先とする。

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

<a id="見出し階層"></a>

**見出し階層：**

| 要素 | サイズ | ウェイト | 装飾 |
|---|---|---|---|
| h1 | `--adlaire-font-size-2xl`（2rem） | 700 | 下線（`--adlaire-color-primary` 2px） |
| h2 | `--adlaire-font-size-xl`（1.5rem） | 600 | 下線（`--adlaire-surface-border` 1px）、上マージン `--adlaire-space-12` |
| h3 | `--adlaire-font-size-lg`（1.125rem） | 600 | 装飾なし |
| h4 | `--adlaire-font-size-sm`（0.875rem） | 500 | 左ボーダー（`--adlaire-color-primary` 3px）＋背景（`--adlaire-surface-soft`）、モノスペースフォント |
| h5 | `0.8125rem` | 600 | 装飾なし、本文色 |
| h6 | `--adlaire-font-size-xs`（0.75rem） | 600 | 装飾なし、`--adlaire-surface-text-muted` |

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

| 寸法 | 値 | 用途 |
|---|---|---|
| ヘッダー高さ | `52px` | 固定ヘッダーの高さ |
| サイドバー幅 | `var(--adlaire-layout-sidebar-compact)` = `260px` | デスクトップサイドバーの幅 |
| コンテンツ最大幅 | `var(--adlaire-layout-container-narrow)` = `760px` | `.ci` の最大幅 |

<a id="レスポンシブ"></a>

**レスポンシブ：**

| ブレークポイント | 動作 |
|---|---|
| `≤ 768px` | サイドバーを画面外に収納（`86vw` 幅）、トグルで引き出し |
| `≤ 400px` | コンテンツ padding を縮小（`1.5rem 1rem`） |

---

## 5. サイドバー（TOC）

<a id="開閉制御"></a>

**開閉制御：**

デスクトップとモバイルで制御方法が異なる。

| モード | 開 | 閉 |
|---|---|---|
| デスクトップ（> 768px） | サイドバーを表示し、コンテンツ開始位置を `var(--adlaire-layout-sidebar-compact)` 分確保する。 | サイドバーを非表示にし、コンテンツ開始位置を `0` にする。 |
| モバイル（≤ 768px） | サイドバーを `translateX(0)` で画面内に表示する。 | サイドバーを `translateX(-86vw)` で画面外に収納する。 |

開閉イベント、状態 class、永続化、復元処理は [`docs/details/builder.md` 詳細本文責務 §7.2](details/builder.md#sec-7-2) を参照する。

<a id="toc-構造"></a>

**TOC 構造：**

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

<a id="toc-リンクスタイル"></a>

**TOC リンクスタイル：**

| クラス | `padding-left` | フォントサイズ | 色 |
|---|---|---|---|
| `.lv1` | `14px` | `--adlaire-font-size-sm`（0.875rem） | `--adlaire-surface-text` |
| `.lv2` | `20px` | `0.8125rem` | `--adlaire-surface-text-muted` |
| `.lv3` | `34px` | `--adlaire-font-size-xs`（0.75rem） | `--adlaire-surface-text-subtle` |
| `.lv4` | `42px` | `--adlaire-font-size-xs`（0.75rem） | `--adlaire-surface-text-subtle` |
| `.lv5` | `50px` | `--adlaire-font-size-xs`（0.75rem） | `--adlaire-surface-text-subtle` |
| `.lv6` | `58px` | `--adlaire-font-size-xs`（0.75rem） | `--adlaire-surface-text-subtle` |

アクティブ状態：`color: --adlaire-color-secondary`、`background: --adlaire-surface-soft`、`border-left: 2px solid --adlaire-color-primary`

---

## 6. コンポーネント

<a id="コードブロック"></a>

**コードブロック：**

```
┌─ .cb-wrap ──────────────────────────────────────────┐
│                         [rust] [コピー]  ← .cb-meta │
│  fn main() { … }                                    │
└─────────────────────────────────────────────────────┘
```

- 背景 `--adlaire-surface-soft-strong`、ボーダー `--adlaire-surface-border` 1px、角丸 `--adlaire-radius-lg`（8px）
- `.cb-meta`：`position: absolute; top: 8px; right: 10px`
- 言語ラベル（`.cl`）：モノフォント、`--adlaire-font-size-xs`、`--adlaire-surface-text-subtle`、大文字
- コピーボタン（`.cb-copy`）：通常 `opacity: 0`、コードブロック hover または keyboard focus 時に表示する。操作後の文言と復元処理は [`docs/details/builder.md` 詳細本文責務 §7.6](details/builder.md#sec-7-6) を参照する
- `pre.cb`：フォントサイズ `--adlaire-font-size-sm`（0.875rem）、行高 `1.65`

<a id="テーブル"></a>

**テーブル：**

- `.tw`（ラッパー）：`overflow-x: auto`、ボーダー・角丸 `--adlaire-radius-lg`・シャドウ
- `th`：モノフォント、`--adlaire-font-size-xs`、背景 `--adlaire-surface-soft`
- 偶数行：背景 `--adlaire-surface-soft`。ホバー行：背景 `--adlaire-surface-soft-strong`

<a id="インラインコード"></a>

**インラインコード：**

- モノフォント、サイズ `0.83em`、背景 `--adlaire-surface-soft-strong`、テキスト `--adlaire-surface-accent-strong`、ボーダー `--adlaire-surface-border` 1px、角丸 `--adlaire-radius-sm`（4px）

<a id="引用blockquote"></a>

**引用（blockquote）：**

- 左ボーダー `3px solid --adlaire-color-primary`、背景 `--adlaire-surface-soft`、角丸 `0 --adlaire-radius-lg --adlaire-radius-lg 0`、斜体テキスト

<a id="定義リスト"></a>

**定義リスト：**

- `dt`：`font-weight: semibold`、`--adlaire-surface-text`
- `dd`：`margin-left: --adlaire-space-6`、`--adlaire-surface-text-muted`

<a id="builder-拡張コンポーネント視覚契約"></a>

**Builder 拡張コンポーネント視覚契約：**

[`docs/details/builder.md` 詳細本文責務 §28](details/builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) で追加される生成 HTML の視覚仕様、レイアウト、レスポンシブ、印刷時の見え方は [`docs/DESIGN.md`](DESIGN.md) デザイン責務を正本とする。出力対象 selector、設定値、REPORT、状態、検証条件は [`docs/details/builder.md`](details/builder.md) 詳細本文責務を参照する。

| 対象 | 視覚契約 |
|------|----------|
| light visual baseline | `:root` は light 固定の custom property を定義する。dark / auto selector、dark background、theme toggle の視覚表現を持たない。 |
| typography / block | admonition、badge、definition list、task list、footnote、math は本文幅内に収め、本文の行長、余白、読みやすさを壊さない。 |
| typography stability | font size を viewport width に比例させず、負の `letter-spacing` を使わない。hover / focus により border、padding、font weight、要素寸法を変えない。 |
| code extension | code title、line number、diff highlight は code block と一体で読める配置にし、copy 対象 text と装飾 text を視覚的に区別する。 |
| navigation runtime UI | section collapse、TOC active、hash focus、skip link は focus indicator と active indicator を常に可視にし、focus / active 化で layout 寸法を変えない。 |
| media UI | image lightbox、Mermaid placeholder / SVG wrapper、print QR は本文の流れを妨げない。外部画像取得や外部 script 読込を前提にした視覚状態を持たない。 |
| responsive | 幅 `320px` の viewport で、本文、見出し、TOC、skip link、admonition、badge、definition list、task list、footnote、math、code title、line numbers、diff highlight、lightbox、print QR の text が重ならず、切れず、親要素外へ不可視にはみ出さない。table と code block だけは既存 scroll wrapper 内の horizontal overflow を許可する。 |
| print | `@media print` では interactive controls、collapse toggle、lightbox trigger UI、TOC active indicator、skip link の画面専用装飾を非表示にする。本文、見出し、画像、code、table、footnote、definition list、task list は非表示にしない。印刷時は全 section を展開表示し、light 固定の背景と文字色を維持する。 |

| 対象 | layout 固定内容 |
|------|----------------|
| `.adlaire-admonition` | 本文幅内の block とし、他 card 内へ入れ子の card 表現を追加しない。title と body は縦積み、長い語は折り返す。 |
| `.adlaire-badge` | inline 要素として扱い、行高を不自然に拡大しない。前後 text を押し潰さない。 |
| `.code-title` | 対応する code block 直前にだけ表示し、code block と分離して floating 表示しない。 |
| `.code-lines` / `.line-no` | line number column と code text column の対応を維持し、折り返し時も行番号と本文が逆転しない。 |
| `.tok-inserted` / `.tok-deleted` / `.tok-context` | 背景色と text color の両方で状態を表し、色だけに依存しない記号または text を維持する。 |
| `.math-inline` / `.math-block` | inline math は行内、block math は本文幅内 block とし、未対応記法を画像化しない。 |
| `.adlaire-lightbox-dialog` | dialog 表示時は viewport 内に収め、画像は `max-inline-size: 100%`、`max-block-size: 100%` 相当で切らない。 |
| `.mermaid-diagram` / `.mermaid-source` | SVG wrapper は本文幅内に収め、fallback text を保持する。 |
| `.print-qr` / `.print-qr-svg` | print 専用 block とし、screen 表示では本文内の常時表示要素にしない。 |
| `.skip-link` / `.is-active` | focus outline は常に可視にし、focus / active 化で layout 寸法を変えない。 |

| 既存 selector | 視覚契約 |
|---------------|----------|
| `.mp` | 本文行長を `max-width: 68ch` で制限する。 |
| `.md-image` | `max-inline-size: 100%` とし、本文コンテナから横方向にはみ出さない。 |
| `.cb-wrap` / `.cb-meta` / `.cl` / `.cb-copy` / `.cb` | `.cb-wrap` を配置基準とし、`.cb-meta` は `top: 8px; right: 10px` に配置する。`.cl` はモノスペース、大文字、`--adlaire-font-size-xs` とする。`.cb-copy` は通常非表示、コードブロックの hover または keyboard focus 時に表示する。`.cb` の背景は `--adlaire-surface-soft-strong` とする。 |
| `.ml-task` / `.ml-task input[type="checkbox"]` | list marker を表示せず、checkbox の accent color は `--adlaire-color-primary` とする。 |
| `.mdl dt` / `.mdl dd` | `dt` は semibold、`dd` は `margin-left: 1.5rem`、文字色は `--adlaire-surface-text-muted` とする。 |
| `.tw` | 横方向 overflow は wrapper 内の scroll で扱う。 |
| `.fn-ref` / `.fn-section` / `.fn-list` | `.fn-ref` は `--adlaire-font-size-xs`、`.fn-section` は上 border で本文と分離し、`.fn-list` は `--adlaire-font-size-sm` とする。 |
| `.hn-link` | 通常時 `opacity: 0`、見出し hover または keyboard focus 時 `opacity: 1` とする。 |
| `#progress-bar` | viewport 上端の `top: 0; left: 0` に固定し、高さ `3px`、背景 `--adlaire-color-primary`、`z-index: 1000`、`transition: width 0.1s linear` とする。 |
| `.mt th[data-sort]` | pointer cursor を表示し、text selection を抑止する。 |
| `#reading-time` | header 右端へ配置し、`margin-left: auto`、文字色 `--adlaire-surface-text-muted`、`--adlaire-font-size-sm`、折り返しなしとする。 |
| `.ch-nav` / `.ch-prev` / `.ch-next` | `.ch-nav` は章末尾で前後 link を両端配置し、`padding: 1rem 0`、`margin-top: 2rem`、上 border を持つ。link は `--adlaire-color-primary` とし、既定の text decoration を表示しない。 |

**印刷表現：**

| 対象 | 印刷時の処理 |
|------|------------|
| `#hdr`、`#sb`、`.cb-copy`、`.hn-link`、`#btt`、`.expand-code`、`#progress-bar`、`.ch-nav` | `display: none` とする。 |
| `#lay`、`#ct` | block 表示、幅 `100%`、余白 `0` とする。 |
| `pre.cb[data-collapsible]` | `max-height: none` とし、内容をすべて表示する。 |
| `a[href^="http"]::after`、`a[href^="https"]::after` | `content: " (" attr(href) ")"` とし、URL を link text の後ろに表示する。 |
| `h2`、`h3` | `page-break-before: avoid` とする。 |

---

## 7. トップへ戻るボタン

- 画面右下（`bottom: 28px / right: 24px`）に固定。`z-index: --adlaire-z-sticky`
- 直径 `38px` 円形、背景 `--adlaire-surface-accent`、色 `#fff`
- 非表示状態は `opacity: 0`、`pointer-events: none`、表示状態は `opacity: 1`、`pointer-events: auto` とする。表示状態へ切り替える条件は [`docs/details/builder.md` 詳細本文責務 §7.7](details/builder.md#sec-7-7) を参照する
- ホバー時 `translateY(-2px)` で浮き上がり

---

## 8. 生成側参照

ビルド実行方法、入出力パス、既定値、終了コード、レポート出力は [`docs/details/builder.md`](details/builder.md) 詳細本文責務を参照する。

`builder` が生成する `assets/style.css` は、[`docs/DESIGN.md`](DESIGN.md) デザイン責務の token、selector、視覚値を実装しなければならない。生成処理と出力検証は [`docs/details/builder.md`](details/builder.md) 詳細本文責務を参照する。

---

<a id="admin-ui-visual-contract"></a>

## 9. 標準管理 UI 視覚契約

標準管理 UI は反復操作と状態比較を行う運用画面とし、装飾的な hero、入れ子の card、背景画像、gradient、外部配信 font、外部配信 CSS、外部配信 JavaScript を使用しない。repository 内の [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) は必須のローカル ES Module として読み込む。ライトモード固定とし、色切替 UI と `prefers-color-scheme` 分岐を持たない。DOM、操作、表示状態、SDK 呼び出しは [`docs/details/ui.md`](details/ui.md)、配布と静的配信は [`docs/details/admin.md`](details/admin.md) を参照する。

**標準管理 UI カラートークン：**

`:root` は `color-scheme: light` と次の custom property だけを持つ。標準管理 UI の JavaScript は値を書き換えない。

| token | 値 | 用途 |
|-------|----|------|
| `--bg` | `#f6f7f9` | page background |
| `--surface` | `#ffffff` | panel、message background |
| `--text` | `#1f2933` | primary text |
| `--muted` | `#687385` | label、empty state、secondary text |
| `--line` | `#d8dde6` | border |
| `--accent` | `#0f766e` | primary action |
| `--accent-dark` | `#115e59` | secondary action text |
| `--danger` | `#b42318` | destructive action、error |
| `--warn` | `#b45309` | warning banner |
| `--ok` | `#0f7a3b` | success message |
| `--focus` | `#2563eb` | keyboard focus indicator |

**標準管理 UI タイポグラフィ・レイアウト：**

| 対象 | 固定値 |
|------|--------|
| universal selector | `box-sizing: border-box`。 |
| `body` | margin `0`、background `var(--bg)`、text `var(--text)`、`font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`、`font-size: 14px`、`line-height: 1.5`。 |
| button / input / select / textarea | `font: inherit`。 |
| `.shell` | CSS grid、`grid-template-columns: 260px minmax(0, 1fr)`、`min-height: 100vh`。 |
| `.sidebar` | width は shell の `260px` track、padding `16px`、右 border `1px solid var(--line)`、background `#202a35`、text `#f9fafb`。 |
| `.brand` | margin `0 0 16px`、font-size `18px`。 |
| `.nav` | CSS grid、gap `6px`。 |
| `.content` | padding `20px`。 |
| `.panel` | CSS grid、gap `14px`、max-width `1180px`、margin `0 auto 24px`、border `1px solid var(--line)`、radius `8px`、background `var(--surface)`、padding `18px`。 |
| `.panel h2` | margin `0`、font-size `20px`。 |
| `.grid` | `repeat(auto-fit, minmax(210px, 1fr))`、gap `12px`。 |
| `.actions` | flex、wrap、gap `8px`、cross-axis center。 |
| `.nav button` | width `100%`、transparent border / background、text `#e5e7eb`、左寄せ。`aria-current="page"` の場合は background `#314154`、text `#ffffff`。 |
| button | min-height `34px`、border `1px solid var(--accent)`、radius `6px`、background `var(--accent)`、text `#ffffff`、padding `6px 12px`、cursor `pointer`。secondary は background `#ffffff` / text `var(--accent-dark)`、danger は border / background `var(--danger)`。 |
| input / select / textarea | width `100%`、border `1px solid var(--line)`、radius `6px`、background `#ffffff`、text `var(--text)`、padding `7px 9px`。textarea は min-height `88px`、vertical resize。checkbox は width `auto`。 |
| label | CSS grid、gap `4px`、text `var(--muted)`、font-size `12px`、font-weight `600`。 |
| `pre` | overflow `auto`、border `1px solid var(--line)`、radius `6px`、background `#f9fafb`、padding `10px`、`white-space: pre-wrap`、`word-break: break-word`。 |
| `.banner` / `.message` | margin `0 0 12px`、border `1px solid var(--line)`、radius `6px`、background `var(--surface)`、padding `10px 12px`。`.message.error` は border `#f3b6b0` / text `var(--danger)`、`.message.success` は border `#a7e0b8` / text `var(--ok)`、`.banner` は border `#f3d08c` / text `var(--warn)`。 |
| `.token-once` | border `#9cc8ff`、background `#f0f7ff`、text `#174ea6`。one-time secret 専用領域以外へ適用しない。 |
| `.output` | min-height `42px`。 |
| `.inline` | flex、gap `8px`、cross-axis center。内部 checkbox は flex `0 0 auto`。 |
| `.empty` | text `var(--muted)`。 |

**標準管理 UI responsive・状態表現：**

| 条件 | 固定動作 |
|------|----------|
| viewport `> 820px` | `.shell` は sidebar `260px` と content の 2 column。 |
| viewport `≤ 820px` | `.shell` は `1fr` の 1 column、`.sidebar` は `position: static`。 |
| `[hidden]` | `display: none !important`。表示切替で DOM 順序を変更しない。 |
| disabled | cursor `not-allowed`、opacity `0.55`。寸法を変更しない。 |
| keyboard focus | `button`、`input`、`select`、`textarea`、`a[href]` の `:focus-visible` は `outline: 2px solid var(--focus)` と `outline-offset: 2px` を使用する。border、padding、要素寸法を変更しない。 |
| error / success / warning | 色だけに依存せず、[`docs/details/ui.md`](details/ui.md) の固定 text と領域で状態を示す。 |
| overflow | `pre` は内部 scroll を許可し、通常 text、label、button text、field text は親要素外へ不可視にはみ出さない。 |

[`admin/index.html`](../admin/index.html) の inline CSS は、[標準管理 UI 視覚契約](#admin-ui-visual-contract) の token、selector、視覚値を実装しなければならない。視覚値を変更する場合は、先に [`docs/DESIGN.md`](DESIGN.md) デザイン責務を改訂する。
