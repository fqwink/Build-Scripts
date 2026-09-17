# Adlaire CI — Builder 詳細仕様

本ファイルは `docs/DETAIL_INDEX.md` から分割した `builder` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、正本関係、実装状態、ロードマップ状態、実装可否の上位判断を記載してはならない。方針、ポリシー、正本関係は `docs/SPEC.md`、実装状態、ロードマップ状態、実装可否は `docs/ROADMAP.md` を正とする。

本ファイルを読む前に、`docs/SPEC.md` で方針とポリシーを確認し、`docs/ROADMAP.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `builder` owner component の主本文であり、collaborator component の仕様は呼び出し境界、状態、fixture、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `builder` |
| collaborator component | `runner`、`api`、`statefile` |
| 持つ内容 | `builder` owner が主本文として定義する Markdown 変換、静的 Web サイト出力、HTML / CSS / JavaScript、theme component、builder fixture、builder owner 追加機能。 |
| 持たない内容 | GitHub read、runner 状態更新、API endpoint、SDK method 実装、UI DOM 詳細、状態 schema、admin 静的配信、setup / release 手順、fixture / PR 証跡正本。 |

---

## 1. 要件

| 項目 | 内容 |
|------|------|
| Go バージョン | Go `1.22` 以上。 |
| 外部依存 | なし。Go 標準ライブラリのみを使用する。外部依存が必要になった場合は実装せず、先に `docs/SPEC.md` Part 1 §4.1 と Part 2 §4 に従って仕様改訂する。 |
| 入力 | UTF-8 エンコードの Markdown ファイル、または Markdown ファイルを含むディレクトリ |
| 出力 | 静的 Web サイトディレクトリ（HTML / CSS / JavaScript / search index） |

---

## 2. ファイルパス設定

`builder` は、以下の既定値を持つ設定構造体で入出力パスを管理する。

```go
type BuildConfig struct {
    Src     string
    Out     string
    Title   string
    Theme   string
    BaseDir string
    Strict  bool
}

var DefaultBuildConfig = BuildConfig{
    Src: "/opt/adlaire-builder/repo/docs",
    Out: "/opt/adlaire-builder/dist/site",
    Title: "Adlaire Documentation",
    Theme: "adlaire-default",
    BaseDir: "",
    Strict: false,
}
```

別の環境で実行する場合は、この既定値を CLI 引数で上書きする。`builder` は設定ファイルを読み込まない。

**CLI 引数仕様：**

| 引数 | 必須 | 既定値 | 説明 |
|------|------|--------|------|
| `--src <path>` | 任意 | `DefaultBuildConfig.Src` | 入力 Markdown ファイルまたは Markdown ディレクトリの絶対パスまたは相対パス。相対パスはカレントディレクトリ基準で解決する。 |
| `--out <path>` | 任意 | `DefaultBuildConfig.Out` | 出力サイトディレクトリの絶対パスまたは相対パス。存在しない場合は作成する。 |
| `--title <text>` | 任意 | `DefaultBuildConfig.Title` | サイト名、`index.html` の `<title>`、header 表示名に使用する。空文字は禁止。 |
| `--theme <name>` | 任意 | `DefaultBuildConfig.Theme` | 初期仕様では `adlaire-default` のみ許可する。 |
| `--base-dir <path>` | 任意 | `DefaultBuildConfig.BaseDir` | 相対リンク・画像解決の基準ディレクトリ。空の場合は `--src` がファイルなら親ディレクトリ、ディレクトリなら `--src` 自身を使用する。 |
| `--strict` | 任意 | `false` | 警告をビルド失敗として扱う。警告が 1 件以上ある場合は終了コード `2` とする。 |
| `--build-id <id>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-build-id` として埋め込む。空文字の場合も空 content の meta を出力する。 |
| `--commit-sha <sha>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-commit-sha` として埋め込む。空文字の場合も空 content の meta を出力する。 |
| `--build-at <iso8601>` | 任意 | 空文字 | 出力 HTML の `<head>` に `adlaire-build-at` として埋め込む。値がある場合は UTC ISO 8601 のみ許可する。 |
| `--version` | 任意 | なし | バイナリ名、仕様名、Go build 情報を 1 行で標準出力へ表示して終了する。 |
| `--help` | 任意 | なし | 引数一覧を標準出力へ表示して終了する。 |

**CLI 引数の異常系：**

| 条件 | 終了コード | 出力 |
|------|------------|------|
| 未知の引数 | `2` | stderr に `unknown option: <name>` |
| `--src` / `--out` / `--title` / `--theme` / `--base-dir` / `--build-id` / `--commit-sha` / `--build-at` の値欠落 | `2` | stderr に `missing value: <name>` |
| `--src` が存在しない | `2` | stderr に `source not found: <path>` |
| `--src` が Markdown ファイルでもディレクトリでもない | `2` | stderr に `source is not markdown file or directory: <path>` |
| `--src` 内の Markdown が UTF-8 として読めない | `2` | stderr に `source is not valid UTF-8: <path>` |
| `--title` が空文字 | `2` | stderr に `title must not be empty` |
| `--theme` が `adlaire-default` 以外 | `2` | stderr に `unknown theme: <name>` |
| `--build-id` が空文字以外で `b{YYYYMMDDHHmmss}` または `b{YYYYMMDDHHmmss}-NNN` 形式でない | `2` | stderr に `invalid build id: <value>` |
| `--commit-sha` が空文字以外で 7〜40 文字の lowercase hex でない | `2` | stderr に `invalid commit sha: <value>` |
| `--build-at` が空文字以外で UTC ISO 8601 でない | `2` | stderr に `invalid build at: <value>` |
| `--out` ディレクトリ作成失敗 | `1` | stderr に `cannot create output directory: <path>` |
| `--out` が既存ファイル | `1` | stderr に `output path is not directory: <path>` |
| `--out` 書き込み失敗 | `1` | stderr に `cannot write output: <path>` |

`--help` と `--version` は他の引数より優先し、成功時は終了コード `0` とする。

**CLI パース固定仕様：**

- 引数は `flag` package 互換の `--name value` と `--name=value` の両方を許可する。
- 短縮オプション（例：`-s`、`-o`）は禁止する。指定された場合は未知の引数として扱う。
- 同一引数が複数回指定された場合は最後の値を採用する。ただし `--strict` は 1 回以上指定されれば `true` とする。
- `--src`、`--out`、`--base-dir` の相対パスは `os.Getwd()` の戻り値を基準に `filepath.Abs()` で絶対パスへ変換する。
- `--base-dir` が空ではない場合、存在するディレクトリでなければならない。存在しない場合は終了コード `2`、stderr に `base directory not found: <path>` を出力する。
- `--base-dir` がファイルの場合は終了コード `2`、stderr に `base path is not directory: <path>` を出力する。
- stderr のエラー行は末尾に改行 1 つを付ける。複数エラーをまとめて出力せず、最初に検出したエラー 1 件で終了する。

**固定出力：**

| 条件 | stdout |
|------|--------|
| `--help` | `Usage: adlaire-ci-build [--src path] [--out path] [--title text] [--theme name] [--base-dir path] [--strict] [--build-id id] [--commit-sha sha] [--build-at iso8601] [--version] [--help]` |
| `--version` | `adlaire-ci-build ADLAIRE_CI_SPEC go=<runtime.Version()>` |

`--help` と `--version` の stdout は 1 行固定とし、末尾に改行 1 つを付ける。`--help` または `--version` を指定した場合、`--src` の存在確認、`--theme` 検証、出力ディレクトリ作成は行わない。

**CLI 値正規化・path 安全契約：**

| 対象 | 正規化 | 禁止 / 失敗条件 |
|------|--------|-----------------|
| `--src` | `filepath.Abs` → `filepath.Clean` | NUL、空文字、存在しない path。 |
| `--out` | `filepath.Abs` → `filepath.Clean` | NUL、空文字、親ディレクトリ不存在、既存通常ファイル。 |
| `--base-dir` | 空なら §2a の規則で決定。指定時は `filepath.Abs` → `filepath.Clean` | NUL、空文字、存在しない path、通常ファイル。 |
| `--title` | 前後空白を除去せず入力値をそのまま使用 | 空文字だけ禁止。空白だけの文字列は空 title として扱い `2`。 |
| `--theme` | 前後空白を除去せず完全一致 | `adlaire-default` 以外。 |

`--src` と `--out` が同一 path、または `--out` が `--src` 配下にある場合は終了コード `2` とし、stderr に `output path must be outside source: <path>` を出力する。`--src` が `--out` 配下にある場合も同じ扱いとする。実装者判断で入力ディレクトリ内へ生成物を混在させてはならない。

---

## 2a. 入力収集・出力パス決定

`builder` は、`--src` がファイルかディレクトリかで入力収集方法を切り替える。

| `--src` 種別 | 入力収集 | 出力 |
|--------------|----------|------|
| Markdown ファイル | 指定ファイル 1 件のみを入力とする。拡張子は `.md` または `.markdown` のみ許可する。 | `index.html` 1 件と `assets/` を出力する。 |
| ディレクトリ | 配下の `.md` / `.markdown` ファイルを再帰収集する。隠しディレクトリ、`.git`、`.ci`、`node_modules`、`vendor`、`dist` は収集対象外とする。 | `index.html` をサイト目次、各 Markdown を `pages/{slug}.html` として出力する。 |

入力ファイルの並び順は、`--base-dir` からの相対パスを `/` 区切りに正規化した文字列の昇順とする。OS やファイルシステムの列挙順に依存してはならない。

Markdown ディレクトリ入力で Markdown ファイルが 0 件の場合は終了コード `2` とし、stderr に `no markdown files found: <path>` を出力する。

`PageData.Title` は、各 Markdown ファイルの最初の h1 見出しを使用する。h1 が存在しない場合は、拡張子を除いたファイル名をタイトル化して使用する。ファイル名タイトル化では、`-` と `_` を空白に置換し、前後空白を除去する。空になった場合は `Untitled` とする。

ページ slug は以下の順で決定する。

1. `--base-dir` からの相対パスから拡張子を除く。
2. パス区切り `/`、空白、`_` を `-` に置換する。
3. §4.1 `slugify` と同じ文字種ルールで正規化する。
4. 空になった場合は `page` とする。
5. 同一 slug が重複した場合は、2 件目以降に `-2`、`-3` のように連番を付与する。

出力パスは以下とする。

| 条件 | `PageData.OutputPath` |
|------|-----------------------|
| 単一 Markdown 入力 | `index.html` |
| ディレクトリ入力のサイト目次 | `index.html` |
| ディレクトリ入力の各 Markdown | `pages/{pageSlug}.html` |

`assets/style.css`、`assets/app.js`、`assets/search-index.json` は常に出力する。`assets/` 配下へ Markdown 由来ファイルを出力してはならない。

**リンク・画像・相対パス解決：**

| 入力 | 処理 |
|------|------|
| `http://` / `https://` URL | 外部 URL としてそのまま出力する。 |
| `mailto:` / `tel:` URL | 内部リンク検証対象外とし、`href` はそのまま出力する。 |
| `#anchor` | 同一ページ内 anchor として、ページ内の一意化後 slug と照合する。 |
| `./doc.md` / `../dir/doc.md` / `dir/doc.markdown#x` | `--base-dir` 基準で Markdown 入力ファイルへ解決し、該当ページの HTML パスへ変換する。anchor がある場合は `#x` を維持し、出力先ページの slug と照合する。 |
| `.md` / `.markdown` 以外の相対リンク | ファイルをコピーせず、元の相対 URL を保持する。存在確認は警告対象外。 |
| 画像 `![alt](path)` | ファイルコピーを行わず、`src` は元 URL を `esc()` して出力する。`assets/` へ画像を複製してはならない。 |

Markdown 間リンクの解決に失敗した場合、HTML は元 URL のまま出力し、`[WARN] BROKEN_PAGE_LINK: <url> (in: <source>)` を出力する。`--strict` が `true` の場合、警告出力後に終了コード `2` とする。ページ間リンク解決で使用するパス比較は、絶対パス化、`filepath.Clean()`、パス区切り `/` 正規化を行った文字列で比較する。

**入力収集固定契約：**

| 項目 | 仕様 |
|------|------|
| 対象拡張子 | `.md`、`.markdown` だけ。大文字拡張子 `.MD`、`.Markdown` は対象外。 |
| 隠しファイル | ファイル名が `.` で始まる Markdown は収集対象外。 |
| symlink | ファイル・ディレクトリとも追跡しない。`os.Lstat` で symlink と判定した path は無視する。 |
| 最大ファイルサイズ | 1 ファイル 10 MiB。超過時は終了コード `2`、stderr `source file too large: <path>`。 |
| 改行 | `\r\n` と `\r` は読み込み時に `\n` へ正規化する。 |
| BOM | 先頭 UTF-8 BOM は除去する。本文途中の BOM は通常文字として扱う。 |

ディレクトリ再帰では、除外ディレクトリに入る前に prune する。除外対象配下で Markdown が見つかっても警告を出さない。

**入力収集・出力生成の機能単位契約：**

`builder` は、以下の機能単位を順番に実行する。各機能単位は前段の出力だけを入力とし、失敗時は後続機能を実行しない。

| 機能単位 | 入力 | 出力 | 失敗条件 | 失敗時の状態 |
|----------|------|------|----------|--------------|
| CLI resolver | `os.Args`、`DefaultBuildConfig` | `BuildConfig` | 未知引数、値欠落、空 title、未知 theme | ファイル読込・出力作成を行わず終了コード `2`。 |
| source collector | `BuildConfig.Src`、`BuildConfig.BaseDir` | `[]PageInput` | source 不在、Markdown 0 件、UTF-8 不正、許可外拡張子 | 出力ディレクトリを変更せず終了コード `2`。 |
| page planner | `[]PageInput` | `[]PageData`、slug map、output path map | slug / output path 衝突を連番解決後も一意化できない | 出力ディレクトリを変更せず終了コード `1`。 |
| markdown renderer | `PageInput`、heading map、footnote map | page HTML fragment、warning list、report counters | 未閉鎖 fence、broken link などの警告があり `--strict=true` | 出力ディレクトリを変更せず終了コード `2`。 |
| site assembler | `[]PageData`、theme component | `index.html`、`pages/*.html`、`assets/*` のメモリ上生成物 | theme component 欠落、template 合成失敗 | 出力ディレクトリを変更せず終了コード `1`。 |
| atomic output writer | メモリ上生成物、`BuildConfig.Out` | 出力サイトディレクトリ | tmp 作成失敗、書き込み失敗、rename 失敗、sync 失敗 | 既存正常出力を維持し終了コード `1`。 |
| report emitter | report counters、warning list | stdout `[WARN]`、`[REPORT]` | stdout 書き込み失敗 | 終了コード `1`。出力済みファイルは巻き戻さない。 |

`PageInput` は `{source_path, relative_path, raw_text}` を持つメモリ上構造とする。`PageData` は `{title, slug, source_path, output_path, html, headings, warnings, reading_time}` を持つメモリ上構造とする。これらの構造は状態ファイルとして保存しない。

**atomic output writer 詳細：**

1. `BuildConfig.Out` と同じ親ディレクトリに `{basename}.tmp.{pid}` を作成する。
2. tmp 内に `index.html`、必要な `pages/`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` をすべて書き込む。
3. 全ファイルを close し、通常ファイルは `0644`、ディレクトリは `0755` に補正する。
4. tmp 配下の各ファイルを `Sync` し、tmp ディレクトリを `Sync` する。
5. 既存 `BuildConfig.Out` が存在する場合は `{basename}.prev.{pid}` へ rename する。
6. tmp を `BuildConfig.Out` へ rename する。
7. 親ディレクトリを `Sync` する。
8. 手順 6 まで成功した後に prev を削除する。prev 削除失敗は `[WARN] OUTPUT_PREV_CLEANUP_FAILED: path={path}` とし、終了コードは `0` のままとする。

手順 1〜4 で失敗した場合は tmp を削除し、既存 `BuildConfig.Out` を変更しない。手順 5 成功後から手順 6 失敗までの間に失敗した場合は、prev を `BuildConfig.Out` へ戻す。復元に失敗した場合は stderr に `cannot restore previous output: <path>` を出力し終了コード `1` とする。

---

## 3. 処理パイプライン

```
┌────────────────────────────────────────────────────────────────┐
│ 1. MDファイル読み込み（raw_lines）                              │
├────────────────────────────────────────────────────────────────┤
│ 2. 見出し抽出パス（フェンス内を除外）                           │
│    → headings: []Heading（slug は重複解決済み）                 │
│    → slugByLine: map[int]string                                │
├────────────────────────────────────────────────────────────────┤
│ 2b. 脚注定義収集パス                                            │
│    → ctx.FootnoteDefs: map[string]string                       │
│    ※ ctx.FootnoteOrder は convert() 実行中に inline() が更新    │
├────────────────────────────────────────────────────────────────┤
│ 3. TOC HTML 生成（buildTOC(headings)）                          │
│    → tocHTML: string                                           │
├────────────────────────────────────────────────────────────────┤
│ 4. MD → HTML 変換（convert(lines, headings, ctx)）              │
│    → ConvertResult.HTML                                        │
├────────────────────────────────────────────────────────────────┤
│ 5. サイトテンプレート合成                                      │
│    theme component、CSS、JS、search index を生成                 │
├────────────────────────────────────────────────────────────────┤
│ 6. サイトディレクトリ書き出し（OUT）                            │
└────────────────────────────────────────────────────────────────┘
```

---

## 4. 関数リファレンス

### 4.1 `slugify(text string) string`

見出しテキストから HTML アンカー用のスラグ（`id` 属性値）を生成する。

**責務：** `slugify` は入力テキストから重複解決前の base slug だけを返す。重複解決は §4.5 の見出し収集処理で `slugCount map[string]int` を使って行う。`slugify` 内で状態を保持してはならない。

**処理手順：**
1. Markdown 記法文字（`` ` * _ ~ [ ] ``）を除去
2. 文字を 1 文字ずつ走査し、空白・ハイフン・ドット・アンダースコアは `-` に変換、英数字・Unicode 文字（カテゴリ `L*`、`N*`）はそのまま保持、それ以外は除去
3. 連続する `-` を 1 つに正規化、前後の `-` をトリム
4. 空文字列になった場合は `section` にフォールバック

**入力：**

| 引数 | 型 | 条件 |
|------|----|------|
| `text` | `string` | UTF-8 文字列。空文字を許可する。 |

**出力：**

| 戻り値 | 型 | 条件 |
|--------|----|------|
| `baseSlug` | `string` | 空文字は禁止。入力が空または除去後に空になる場合は `section`。 |

**エラー：** 戻り値エラーは持たない。入力が UTF-8 不正になる可能性は §2 の入力ファイル読み込み段階で排除する。

**禁止事項：**
- package 変数、グローバル map、呼び出し間で残る状態を使用しない。
- 重複時の `-2` 付与を `slugify` 内で行わない。
- URL encode は行わない。`id` と `href` には本関数の戻り値または一意化後 slug をそのまま使用する。

---

### 4.2 `esc(s string) string`

Go 標準ライブラリ `html.EscapeString(s)` 相当の処理を行う。HTML 特殊文字（`<`、`>`、`&`、`"`、`'`）をエスケープする。

**実装契約：**
- Go 実装では `html.EscapeString` を使用する。
- 追加の独自置換、Markdown 記法変換、改行変換を行わない。
- 空文字は空文字を返す。

---

### 4.3 `inline(text string, ctx *RenderContext) string`

インライン Markdown 記法を HTML に変換する。コードスパンを先にキャラクターレベルのスキャンで処理することで、後続の正規表現がコードスパン内の記法を誤って変換するのを防ぐ。

**関連型：**

```go
type RenderContext struct {
    FootnoteDefs      map[string]string
    FootnoteOrder     []string
    FootnoteSeen      map[string]bool
    InternalLinkRefs  map[string]bool
    BrokenLinks       []string
    HeadingSkipCount  int
    CharCount         int
}
```

`ctx` は `nil` 禁止。呼び出し元は `convert()` 実行前に `RenderContext` を作成し、脚注定義を `FootnoteDefs` に収集しておく。

**処理手順：**

1. **コードスパン抽出（文字列スキャン）**
   - バッククォート（`` ` ``）の位置を順次検索
   - `` `` `` で始まる場合はダブルバッククォートコードスパンとして処理
   - 単一バッククォートの場合はシングルバッククォートコードスパンとして処理
   - コードスパン内テキストは `esc()` した上で `<code class="ic">` でラップ
   - コードスパン外テキストは `esc()` してそのままセグメントに追加

2. **インライン記法の正規表現置換**（コードスパン処理後のテキストに適用）

| パターン | 出力 |
|---------|------|
| `***text***` | `<strong><em>text</em></strong>` |
| `**text**` | `<strong>text</strong>` |
| `*text*`（単独 `*`） | `<em>text</em>` |
| `__text__` | `<strong>text</strong>` |
| `_text_`（単独 `_`） | `<em>text</em>` |
| `~~text~~` | `<del>text</del>` |
| `![alt](url)` | `<img src="url" alt="alt" style="max-width:100%">` |
| `[label](url)`（`url` が `http://` または `https://` で始まる場合） | `<a href="url" target="_blank" rel="noopener noreferrer">label</a>` |
| `[label](url)`（上記以外 — 内部リンク・アンカー） | `<a href="url">label</a>` |
| `[^id]` | `<sup><a href="#fn-id" id="fnref-id" class="fn-ref">[N]</a></sup>`（N は参照順の番号） |

**処理順の注意：** 画像記法（`![alt](url)`）はリンク記法（`[label](url)`）より先にマッチングする。リンク記法はマッチング後に `url` が `http://` または `https://` で始まるかを判定し、外部リンクと内部リンクを区別する。脚注参照（`[^id]`）はリンク置換後に適用する。

**内部リンク整合性チェック：**

`url` が `#` で始まる内部アンカーリンクを処理する際、アンカー部分（`#` 以降）を `ctx.InternalLinkRefs[anchor] = true` として記録する。`convert()` 末尾で生成済み slug の集合と照合し、一致しないアンカーを次のように報告する：

```
[WARN] BROKEN_LINK: #anchor-text  (in: [label](#anchor-text))
```

不一致件数は `[REPORT]` の `broken_links` フィールドに反映される。照合は変換完了後に行うため、ドキュメント内の順序（前方参照・後方参照）を問わず検証できる。

**脚注・内部リンク状態：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `ctx.FootnoteDefs` | `map[string]string` | 脚注 ID → テキストの対応（`convert()` 呼び出し前に収集） |
| `ctx.FootnoteOrder` | `[]string` | 本文中での参照順脚注 ID リスト（`inline()` が参照時に追記） |
| `ctx.FootnoteSeen` | `map[string]bool` | 同一脚注 ID の重複登録防止 |
| `ctx.InternalLinkRefs` | `map[string]bool` | 本文中に出現した内部アンカー参照（`#` 除いた文字列）の集合 |

**脚注番号付与：**
- 初回参照時のみ `FootnoteOrder` へ ID を append する。
- 2 回目以降の同一 ID 参照では既存番号を再利用する。
- `FootnoteDefs` に存在しない ID でも HTML 参照は出力し、脚注本文は空文字として扱わず、末尾脚注出力時に `[WARN] MISSING_FOOTNOTE: id` を出す。

**インライン記法のネスト仕様：**

| 入力 | 出力 |
|------|------|
| `***text***` | `<strong><em>text</em></strong>` |
| `**_text_**` | `<strong><em>text</em></strong>` |
| `__*text*__` | `<strong><em>text</em></strong>` |
| `*__text__*` | `<em><strong>text</strong></em>` |
| `_**text**_` | `<em><strong>text</strong></em>` |

上表以外のネストした強調・削除・リンクの組み合わせは追加変換しない。未対応ネストは、先にマッチした外側または内側の単一記法だけを変換し、残った Markdown 記号は `esc()` 済みテキストとして出力する。実装者判断で CommonMark 全互換のネスト処理を追加してはならない。

---

### 4.4 `buildTOC(headings []Heading) string`

見出しリストから TOC（目次）の HTML を生成する。h1〜h3 のみを対象とし（h4 は除外）、子見出しを持つ見出しはグループとしてアコーディオン形式に構築する。

**関連型：**

```go
type Heading struct {
    Level      int
    Text       string
    Slug       string
    LineNumber int
}
```

**引数：** `headings []Heading`

**出力：** TOC の `<li>` 要素群の HTML 文字列（`<ul>` ルートは HTML テンプレート側で定義）

**グループ（`.tg`）とリーフ（`.ti`）の判定：**
現在の見出しの次の見出しレベルが現在より深い場合、グループとして扱い、展開ボタン（`.tg-btn`）付きの `<ul>` をネストする。それ以外はリーフ（`.ti`）として `<li>` 1 個を出力する。

**スタック管理：** 内部スタックは `[]tocStackItem` とする。

```go
type tocStackItem struct {
    Level int
    Kind  string // "group" または "item"
}
```

`closeTo(targetLevel int)` で不要になった `</ul></li>` を閉じる。スタックにはグループ（`"group"`）だけでなく、レベル管理のためリーフ（`"item"`）も積む。`closeTo()` はリーフをサイレントに捨て、グループのみ `</ul></li>` を出力する。

**生成 HTML 構造（グループの場合）：**
```html
<li class="tg">
  <div class="tg-row">
    <a href="#slug" class="tl lv1" data-slug="slug">見出しテキスト</a>
    <button class="tg-btn" aria-expanded="false" data-target="tg-slug" aria-label="展開">
      <svg>…</svg>
    </button>
  </div>
  <ul id="tg-slug" class="tc" hidden>
    …子アイテム…
  </ul>
</li>
```

**生成 HTML 構造（リーフの場合）：**
```html
<li class="ti">
  <a href="#slug" class="tl lv1" data-slug="slug">見出しテキスト</a>
</li>
```

**エラー：** 戻り値エラーは持たない。`Level` が 1〜4 以外の `Heading` は無視し、`[WARN] INVALID_HEADING_LEVEL: line={LineNumber}` を stdout へ出力する。警告件数は `[REPORT]` の `warnings` に含める。

---

### 4.5 `convert(lines []string, headings []Heading, ctx *RenderContext) ConvertResult`

Markdown の行リストを走査し、HTML コンテンツ文字列を生成するメインコンバーター。

**関連型：**

```go
type ConvertResult struct {
    HTML               string
    ReadingTimeMinutes int
    BrokenLinks        int
    HeadingSkips       int
    Warnings           []string
}
```

**引数：**
- `lines []string` — Markdown の全行。末尾改行は含めない。
- `headings []Heading` — 見出し収集済みリスト。各 `Heading.Slug` は重複解決済みとする。
- `ctx *RenderContext` — `inline()` と共有する変換状態。`nil` 禁止。

**戻り値：**
- `HTML` — 本文 HTML 文字列。
- `ReadingTimeMinutes` — §4.5 の読了時間算出結果。
- `BrokenLinks` — 内部リンク不一致件数。
- `HeadingSkips` — 見出し階層スキップ件数。
- `Warnings` — `[WARN]` として stdout 出力した警告本文の一覧。

**エラー：** 戻り値エラーは持たない。入力ファイル不存在、UTF-8 不正、書き込み失敗などの異常は §8 の実行方法と終了コードで扱う。Markdown 構文上の不足は警告またはフォールバック出力で処理する。

**内部バッファと状態変数：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `out` | `strings.Builder` | 本文 HTML 断片の蓄積 |
| `fence_active` | `bool` | コードフェンス内かどうか |
| `fence_lang` | `string` | コードフェンスの言語識別子 |
| `fence_buf` | `[]string` | フェンス内の行バッファ |
| `fence_marker` | `string` | フェンス開始マーカー（`` ``` `` or `~~~`） |
| `para_buf` | `[]string` | 段落テキストの行バッファ |
| `list_stack` | `[]listStackItem` | リストのネスト状態スタック |
| `table_buf` | `[]string` | テーブル行のバッファ |
| `slugByLine` | `map[int]string` | `headings` から生成する行番号 → slug の対応 |

**内部ヘルパー関数：**

- `flushPara()` — `para_buf` を `<p class="mp">` として出力し、バッファをクリア
- `flushList()` — `list_stack` を巻き戻し、すべての `</ul>` / `</ol>` を閉じる
- `flushTable()` — `table_buf` をパースし `<div class="tw"><table class="mt">` として出力
- `emitCode(&out, fence_lang, fence_buf)` — `fence_buf` を `<div class="cb-wrap">` 構造として出力

**ブロック要素の検出優先順位（1 行ずつ処理）：**

1. フェンスコードブロック開始・終了（`` ``` `` / `~~~` で始まる行）
2. 見出し（`#` で始まる行）
3. 水平線（`---`、`***`、`___`）
4. テーブル（`|` で始まる行、連続する `|` 行をまとめて処理）
5. 引用（`>` で始まる行、連続する `>` でネスト可能）
6. 定義リスト（`: 定義` 形式の行かつ `para_buf` に用語がある場合）
7. リスト項目（`- * +` または `1.` 形式、`[ ]`/`[x]` プレフィックスでタスクリスト）
8. 空行（バッファのフラッシュトリガー）
9. 脚注定義行（`[^id]:` で始まる行、`ctx.FootnoteDefs` 収集済みのためスキップ）
10. 段落（上記以外の非空行、連続行を 1 つの `<p>` にまとめる）。先読みループは次のいずれかに該当する行で停止する：`#`（見出し）、`|`（テーブル）、`` ` ``×3以上（フェンス）、`~`×3以上（フェンス）、`>`（引用）、リストマーカー（`[-*+]` または `\d+[.)]`）、`: `（定義リストマーカー）、水平線（`---+`・`***+`・`___+`）

**Markdown passthrough 禁止契約：**

| 入力 | 出力 |
|------|------|
| 生 HTML 行 `<div>text</div>` | `<p class="mp">&lt;div&gt;text&lt;/div&gt;</p>` |
| HTML comment `<!-- x -->` | `<p class="mp">&lt;!-- x --&gt;</p>` |
| script/style tag | tag 全体を text として `esc()` し、実行可能 HTML にしない。 |
| unknown Markdown 記法 | text として `esc()` し、独自 HTML を生成しない。 |

`builder` は Markdown 入力由来の HTML を信頼済みとして扱ってはならない。`PageData.BodyHTML` に入る HTML は、本仕様で生成すると定義したタグと属性だけで構成する。

**テーブル変換の詳細：**
セパレーター行（`:---:`、`---` などで構成された行）のインデックスを自動検出し、セパレーター行より前の行をヘッダー（`<th>`）、それ以降を本文（`<td>`）として出力する。セパレーター行自体は出力しない。

テーブル列数はヘッダー行のセル数を正とする。本文行のセル数が不足する場合は空文字セルを補い、超過する場合は超過分を最後のセルへ ` | ` で連結する。ヘッダー行が存在しない、またはセパレーター行だけの場合はテーブルとして扱わず、段落として出力する。

**引用ネストの詳細：**
`>` で始まる連続行をまとめて収集し、`renderBlockquote(lines []string, ctx *RenderContext) string` が再帰的にネストを処理する。1 レベル分の `>` を剥いた後、内側行を先頭から走査し、`>` で始まる連続する行は `renderBlockquote()` を再帰呼び出し、それ以外の行は `inline(text, ctx)` でレンダリングして結合する。これにより、単一行・複数行・混在ネスト（同一ブロック内で `>` 行と `>>` 行が混在する場合）をすべて正しく処理する。例：`>> text` → `<blockquote class="mbq"><blockquote class="mbq">text</blockquote></blockquote>`。

**タスクリストの詳細：**
リスト項目のコンテンツが正規表現 `^\[([ xX])\]\s+` に一致する場合、`<li class="ml-task">` として出力する。チェック済み（`[x]` / `[X]`）は `checked` 属性付き、未チェック（`[ ]`）は属性なしの `<input type="checkbox" disabled>` を先頭に配置する。

リストネストは先頭空白 2 文字を 1 レベルとして扱い、tab は 4 空白へ展開してから判定する。最大ネストは 6 レベルとし、7 レベル以上は 6 レベルとして出力し `[WARN] LIST_NESTING_CLAMPED: line={line}` を出す。

**定義リストの詳細：**
`: 定義` 行（`strings.HasPrefix(line, ": ")`、コロン＋スペース1文字）を検出したとき `para_buf` に内容があれば、`para_buf` の末尾要素を用語（`<dt>`）として取り出し、`<dl class="mdl"><dt>用語</dt><dd>定義</dd></dl>` を出力する。連続する `: ` 行は同一 `<dl>` 内の追加 `<dd>` としてまとめて処理し、その後に `</dl>` を閉じる。

**脚注定義行のスキップ：**
`^\[\^[^\]]+\]:` に一致する行は `ctx.FootnoteDefs` への収集が完了しているためスキップし、本文への出力を行わない。

**主要ブロックの固定 HTML 断片：**

| Markdown | HTML |
|----------|------|
| `# Title` | `<h1 id="title" class="mh h1">Title<button class="hn-link" data-href="#title" aria-label="リンクをコピー">¶</button></h1>` |
| `text` | `<p class="mp">text</p>` |
| `---` | `<hr class="mr">` |
| `> quote` | `<blockquote class="mbq">quote</blockquote>` |
| `- item` | `<ul class="ml"><li>item</li></ul>` |
| `1. item` | `<ol class="ml"><li>item</li></ol>` |
| `- [x] done` | `<ul class="ml"><li class="ml-task"><input type="checkbox" disabled checked>done</li></ul>` |
| `term` + 次行 `: desc` | `<dl class="mdl"><dt>term</dt><dd>desc</dd></dl>` |

上表の属性順、class 名、button 文言、`checked` 属性の位置は固定する。テスト比較では、タグ間および text node 内の連続空白を 1 つの ASCII space へ正規化してから比較する。タグ名、属性名、属性値、親子構造は完全一致させる。

**脚注セクションの末尾出力：**
`convert()` 末尾で `ctx.FootnoteOrder` が非空の場合、`<section class="fn-section">` 内に参照順番号付きの脚注リスト（`<ol class="fn-list">`）を出力する。各脚注には本文への戻りリンク（`<a class="fn-back">↩</a>`）を付与する。

**フェンスコードブロックの未閉鎖フォールバック：**
ファイル末尾まで読んだ時点で `fence_active` が `true` のままの場合（閉じる `` ``` `` がない場合）、`fence_buf` にコンテンツがあれば `emitCode(&out, fence_lang, fence_buf)` を呼び出して強制出力する。`fence_buf` が空（フェンス開始直後に EOF）の場合は何も出力しない。未閉鎖フェンスは `[WARN] UNCLOSED_FENCE: line={startLine}` を出力し、`Warnings` に追加する。

**見出し階層スキップ警告：**

`convert()` 内で `prevHeadingLevel int = 0` をローカル変数として保持する。見出し行（`#` で始まる行）を処理するたびに現在レベルと前回レベルを比較し、2 段以上の降順スキップ（例：h1→h3、h2→h4）を検出した場合に次の形式で `[WARN]` を出力する：

```
[WARN] HEADING_SKIP: h1→h3 "見出しテキスト"
```

- 昇順への復帰（例：h3→h1）はスキップに該当しない（章の区切りとして正常）
- 同レベルの連続（h2→h2）・1段降順（h2→h3）もスキップに該当しない
- 件数は `[REPORT]` の `heading_skips` フィールドに反映される

**読了時間集計：**

`convert()` 内で `ctx.CharCount = 0` に初期化する。段落・リスト・引用テキストを `inline(text, ctx)` 処理する直前に、元の Markdown テキスト文字数（スペース・改行を含む）を加算する。コードブロック・フェンス内テキスト・見出しテキスト・テーブルは集計対象外とする。

読了時間の算出：
```go
readingTimeMinutes := int(math.Ceil(float64(ctx.CharCount) / 200.0)) // 200文字/分、切り上げ
```

算出した `readingTimeMinutes` は `ConvertResult.ReadingTimeMinutes` として返し、呼び出し元が `[REPORT]` 行と `.build_logs/{id}.json` に記録する。また HTML ヘッダーへの静的埋め込み（§5）にも使用する。

**見出し出力 HTML 構造：**
`#` で始まる行を `h1`〜`h4` に変換する際、末尾に `.hn-link` ボタンを付与する。

```html
<h2 id="slug" class="mh h2">見出しテキスト<button class="hn-link" data-href="#slug" aria-label="リンクをコピー">¶</button></h2>
```

- `data-href` 属性：`#` + 一意化後 slug
- `¶`（U+00B6 PILCROW SIGN）を使用
- CSS で通常時 `opacity: 0`、親見出し要素のホバー時に `opacity: 1` に変化する

**見出しスラグ重複解決：**

見出し収集処理では `slugCount map[string]int` をローカル変数として保持し、同一 base slug が複数の見出しに割り当てられる場合に一意化する。

| 条件 | スラグ |
|------|--------|
| 初出 | `{slug}` そのまま |
| 2 回目 | `{slug}-2` |
| 3 回目 | `{slug}-3` |
| N 回目 | `{slug}-{N}` |

```go
slugCount := map[string]int{}

uniqueSlug := func(base string) string {
    n := slugCount[base] + 1
    slugCount[base] = n
    if n == 1 {
        return base
    }
    return fmt.Sprintf("%s-%d", base, n)
}
```

一意化後のスラグは `id` 属性・`data-href` 属性・TOC リンク `href`（§4.4）・`¶` ボタン・全文検索インデックス（§7.9）・前後章ナビゲーション（§7.15）のすべてで共通使用する。

---

### 4.6 `emitCode(out *strings.Builder, lang string, buf []string)`

フェンスコードブロックを HTML に変換して `out` に追記する。`convert()` 内クロージャではなく、Go ファイル内の非公開補助関数として実装する。

**入力：**

| 引数 | 型 | 条件 |
|------|----|------|
| `out` | `*strings.Builder` | `nil` 禁止。呼び出し元の本文 HTML 出力先。 |
| `lang` | `string` | フェンス開始行の言語識別子。未指定時は空文字。 |
| `buf` | `[]string` | フェンス内本文。空配列を許可する。 |

**出力 HTML 構造：**
```html
<div class="cb-wrap" data-lang="rust">
  <div class="cb-meta">
    <span class="cl">rust</span>
    <button class="cb-copy" aria-label="コピー">コピー</button>
  </div>
  <pre class="cb"><code>…エスケープ済みコード…</code></pre>
</div>
```

- `data-lang` 属性：言語識別子（なければ空文字列）
- `.cl`（言語ラベル）：言語識別子がない場合は出力しない
- コード内容：`strings.Join(buf, "\n")` を `esc()` でエスケープしてから出力
- `buf` が空の場合でも空の `<pre class="cb"><code></code></pre>` を出力する

**エラー：** 戻り値エラーは持たない。`out == nil` はプログラム不備として `panic("nil output builder")` を許可する。

---

## 5. 静的 Web サイト出力構造

サイト全体の合成は `assembleSite(site SiteData, theme Theme) error` が担当する。`convert()` は各 Markdown ページの本文 HTML を生成し、`assembleSite()` はページ HTML、共通 CSS、共通 JavaScript、検索インデックス、テーマコンポーネントを出力サイトディレクトリへ書き出す。

**関連型：**

```go
type SiteData struct {
    Title          string
    Pages          []PageData
    SearchIndex    []SearchIndexEntry
    GeneratedAtUTC string
}

type PageData struct {
    Title              string
    SourcePath         string
    OutputPath         string
    Layout             string
    TocHTML            string
    BodyHTML           string
    ReadingTimeMinutes int
}

type Theme struct {
    Name       string
    Components ThemeComponents
    CSS        string
    JS         string
}

type ThemeComponents struct {
    Header     string
    Sidebar    string
    Breadcrumb string
    TOC        string
    Search     string
    Footer     string
    CodeBlock  string
    Table      string
    Pagination string
}

type SearchIndexEntry struct {
    URL   string
    ID    string
    Title string
    Body  string
}
```

**入力契約：**

| フィールド | 型 | 条件 |
|------------|----|------|
| `SiteData.Title` | `string` | 空は禁止。HTML 出力時は `esc()` する。 |
| `SiteData.Pages` | `[]PageData` | 1 件以上。単一 Markdown 入力の場合も 1 ページのサイトとして扱う。 |
| `SiteData.SearchIndex` | `[]SearchIndexEntry` | サイト内検索用。未生成時は空配列。 |
| `SearchIndexEntry.URL` | `string` | `PageData.OutputPath` または `PageData.OutputPath + "#" + headingSlug`。空は禁止。 |
| `SearchIndexEntry.ID` | `string` | 見出し slug。ページ単位エントリの場合は空文字を許可する。 |
| `SearchIndexEntry.Title` | `string` | 検索結果に表示するタイトル。空は禁止。 |
| `SearchIndexEntry.Body` | `string` | 段落先頭 200 文字以内の説明文。該当本文がない場合は空文字を許可する。 |
| `Theme.Name` | `string` | 初期仕様では `adlaire-default` 固定。 |
| `PageData.Title` | `string` | 空の場合は `SiteData.Title` を使用する。HTML 出力時は `esc()` する。 |
| `PageData.SourcePath` | `string` | 入力 Markdown の絶対パスまたは `--base-dir` からの相対パス。 |
| `PageData.OutputPath` | `string` | `--out` からの相対 HTML パス。トップページは `index.html`。 |
| `PageData.Layout` | `string` | 初期仕様では `document` または `index`。未知値は禁止。 |
| `TocHTML` | `string` | `buildTOC(headings)` の戻り値。`<ul id="toc-root">` の内側へ挿入する。 |
| `BodyHTML` | `string` | `ConvertResult.HTML`。`<div class="ci">` の内側へ挿入する。 |
| `ReadingTimeMinutes` | `int` | 1 以上。0 以下の場合は `1` として表示する。 |
| `GeneratedAtUTC` | `string` | UTC ISO 8601。空の場合は生成日時 meta を出力しない。 |

`assembleSite()` は上記フィールドを結合してファイルを書き出すだけとし、Markdown 変換、slug 生成、TOC 生成、検索インデックス抽出、警告集計を行ってはならない。

**出力ディレクトリ構造：**

```text
{out}/
  index.html
  pages/
    {slug}.html
  assets/
    style.css
    app.js
    search-index.json
```

単一 Markdown 入力の場合、本文ページを `index.html` として出力し、`pages/` は作成しなくてよい。Markdown ディレクトリ入力の場合、`index.html` はサイト目次ページとし、各 Markdown ファイルを `pages/{slug}.html` として出力する。

**必須生成物内容契約：**

| ファイル | 内容 | 空許可 |
|----------|------|--------|
| `index.html` | 完全な HTML document。`<!DOCTYPE html>` から始まる。 | 不可 |
| `pages/{slug}.html` | 完全な HTML document。directory 入力時のみ。 | 不可 |
| `assets/style.css` | §6 の class / selector を含む CSS。 | 不可 |
| `assets/app.js` | §7 の DOMContentLoaded handler と機能実装。 | 不可 |
| `assets/search-index.json` | `SearchIndexEntry[]` の JSON。entry 0 件でも `[]`。 | 可 |

生成物は UTF-8、LF 改行とする。HTML、CSS、JS、JSON の末尾には LF を 1 つ付ける。BOM は出力しない。

**相対 root 算出：**

`relativeRoot` は、各 HTML ファイルから `assets/` へ到達するための相対 prefix とする。

| `PageData.OutputPath` | `relativeRoot` |
|-----------------------|----------------|
| `index.html` | `""` |
| `pages/example.html` | `"../"` |
| `pages/dir-example.html` | `"../"` |

初期仕様ではページ HTML を `pages/` 直下に平坦化するため、`relativeRoot` は上表の 2 種類のみとする。`pages/dir/page.html` のような階層出力を追加する場合は、先に本表、リンク解決、検索 index URL、breadcrumb 仕様を改訂する。

**出力更新手順：**

`assembleSite()` は `--out` を直接途中更新してはならない。以下の順で一時ディレクトリへ完全生成してから置換する。

1. `--out` と同じ親ディレクトリに `{basename}.tmp.{pid}` を作成する。
2. 一時ディレクトリ配下へ `index.html`、ページ HTML、`assets/` をすべて書き出す。
3. すべてのファイルについて `file.Sync()` と `file.Close()` を完了する。
4. 一時ディレクトリ内の必須ファイル存在を検証する。
5. 既存 `--out` が存在する場合は `{basename}.prev.{pid}` へ `os.Rename` する。
6. 一時ディレクトリを `--out` へ `os.Rename` する。
7. 置換成功後、旧 `{basename}.prev.{pid}` を削除する。

手順 1〜4 で失敗した場合は一時ディレクトリを削除し、既存 `--out` を変更してはならない。手順 6 で失敗した場合は `{basename}.prev.{pid}` が存在するか確認し、存在する場合は `{basename}.prev.{pid}` を `--out` へ戻す復旧 rename を必ず 1 回試行する。復旧 rename が失敗した場合は stderr に `cannot restore previous output directory: <path>` を出力した後、続けて `cannot replace output directory: <path>` を出力し、終了コード `1` とする。`{basename}.prev.{pid}` が存在しない場合は `cannot replace output directory: <path>` のみを出力し、終了コード `1` とする。

`--out` が既存ファイルでディレクトリではない場合は終了コード `1` とし、stderr に `output path is not directory: <path>` を出力する。

**初期テーマコンポーネント：**

初期実装の theme は `adlaire-default` のみとする。theme component は Go コード内の内製テンプレートとして保持し、外部テンプレートファイルを読み込んではならない。

| コンポーネント | 責務 |
|----------------|------|
| `header` | サイト名、現在ページ名、読了時間を表示する。 |
| `sidebar` | サイト内ページ一覧と現在ページの TOC を表示する。 |
| `breadcrumb` | `index.html` から現在ページまでの階層を表示する。 |
| `toc` | 現在ページの見出し TOC を表示する。 |
| `search` | `assets/search-index.json` を読み込み、サイト内検索 UI を提供する。 |
| `footer` | 生成時刻とサイト名を表示する。 |
| `codeblock` | 言語ラベル、コピー操作、折りたたみ状態を提供する。 |
| `table` | 横スクロール wrapper と列 sort を提供する。 |
| `pagination` | 前後ページへのリンクを表示する。 |

カスタムテーマ、外部テンプレート、テーマパッケージ、theme component 差し替え、複数 theme 同梱の入出力、状態、検証条件は本ファイルでは定義しない。

**ページ HTML の必須 DOM 構造：**

```html
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="adlaire-build-id" content="{BuildMeta.BuildID}">
  <meta name="adlaire-commit-sha" content="{BuildMeta.CommitSHA}">
  <meta name="adlaire-build-at" content="{BuildMeta.BuildAt}">
  <title>{PageData.Title}</title>
  <link rel="stylesheet" href="{relativeRoot}assets/style.css">
</head>
<body>
  <div id="progress-bar"></div>   <!-- 読み取り進捗バー（ページ上端固定、高さ 3px、幅 = スクロール率 % → §7.13） -->
  <header id="hdr">        <!-- 固定ヘッダー（高さ 52px、背景 --adlaire-surface-accent） -->
    <span id="doc-title">{PageData.Title}</span>
    <span id="reading-time">約 {PageData.ReadingTimeMinutes} 分</span>  <!-- 読了時間（ビルド時に静的埋め込み → §4.5・§6） -->
  </header>
  <div id="lay">           <!-- フレックスコンテナ -->
    <nav id="sb">          <!-- サイドバー（幅 260px、固定） -->
      <div class="sb-search-wrap">
        <input id="sb-search" type="search" autocomplete="off" aria-label="セクション検索">  <!-- TOC 検索 -->
      </div>
      <div id="sb-toc">
        <ul class="tr" id="toc-root">{PageData.TocHTML}</ul>  <!-- TOC リスト -->
        <p id="sb-none" class="sb-none" hidden>一致するセクションはありません</p>
      </div>
    </nav>
    <main id="ct">         <!-- コンテンツエリア -->
      <div class="ci">     <!-- 最大幅 760px センタリングコンテナ -->
        {PageData.BodyHTML}
      </div>
    </main>
  </div>
  <button id="btt" type="button" aria-label="トップへ戻る">↑</button>        <!-- トップへ戻るボタン -->
  <script src="{relativeRoot}assets/app.js"></script>
</body>
</html>
```

**ディレクトリ入力時の `index.html` 仕様：**

`index.html` はサイト目次ページとして生成し、Markdown 本文を持たない。`PageData.Layout` は `index` とする。

| 項目 | 仕様 |
|------|------|
| `<title>` | `SiteData.Title` |
| `#doc-title` | `SiteData.Title` |
| `.ci` 内 | `<h1 id="site-index" class="mh h1">{SiteData.Title}<button class="hn-link" data-href="#site-index" aria-label="リンクをコピー">¶</button></h1>` の後に `<ul class="site-page-list">` を出力する。 |
| ページ一覧 | `SiteData.Pages` のうち `Layout == "document"` のページを入力ソート順で `<li><a href="{OutputPath}">{Title}</a></li>` として出力する。 |
| TOC | サイト目次ページの TOC は `site-index` 1 件だけを含める。 |
| 検索 index | サイト目次ページのエントリを 1 件追加し、`url` は `index.html#site-index`、`id` は `site-index`、`title` は `SiteData.Title`、`body` は空文字とする。 |

**必須 DOM ID / class 契約：**

| セレクター | 個数 | 使用者 | 変更可否 |
|------------|------|--------|----------|
| `#progress-bar` | 各 HTML 1 個 | `assets/app.js` §7.13 | 変更禁止 |
| `#hdr` | 各 HTML 1 個 | CSS / layout | 変更禁止 |
| `#doc-title` | 各 HTML 1 個 | header | 変更禁止 |
| `#reading-time` | 各 HTML 1 個 | header | 変更禁止 |
| `#lay` | 各 HTML 1 個 | CSS / sidebar layout | 変更禁止 |
| `#sb` | 各 HTML 1 個 | sidebar JS | 変更禁止 |
| `#sb-search` | 各 HTML 1 個 | TOC/search JS | 変更禁止 |
| `#toc-root` | 各 HTML 1 個 | TOC JS | 変更禁止 |
| `#ct` | 各 HTML 1 個 | content layout | 変更禁止 |
| `.ci` | 各 HTML 1 個 | content width | 変更禁止 |
| `#btt` | 各 HTML 1 個 | top button JS | 変更禁止 |

実装は上表のセレクターを追加、削除、リネームしてはならない。UI 改善で新しいセレクターが必要な場合は、先に本表へ追加し、§6 と §7 の CSS / JavaScript 契約を同時に更新する。

**サイト合成の禁止事項：**
- `PageData.BodyHTML`、`PageData.TocHTML` はすでに HTML として生成済みのため、`assembleSite()` 内で再エスケープしない。
- `assets/search-index.json` は `encoding/json` の出力だけを受け付け、文字列連結で JSON を自作しない。
- `<header id="hdr">`、`<nav id="sb">`、`<main id="ct">`、`<div class="ci">` の id / class を変更しない。
- `assets/style.css`、`assets/app.js` 以外の CSS / JavaScript を生成してはならない。
- 外部 CSS、外部 JavaScript、外部フォント参照を追加しない。

**HTML head / asset 参照固定契約：**

| 項目 | 仕様 |
|------|------|
| `<title>` | `PageData.Title` を `esc()` した値。site index は `SiteData.Title`。 |
| meta build id | `--build-id` 未指定時も `content=""` で出力する。 |
| meta commit sha | `--commit-sha` 未指定時も `content=""` で出力する。 |
| meta build at | `--build-at` 未指定時も `content=""` で出力する。 |
| generated at | `GeneratedAtUTC` が空でない場合だけ `<meta name="adlaire-generated-at" content="{GeneratedAtUTC}">` を出力する。 |
| CSS link | HTML ごとに `{relativeRoot}assets/style.css` 1 件だけ。 |
| JS script | `</body>` 直前に `{relativeRoot}assets/app.js` 1 件だけ。 |

HTML には inline `<style>`、inline `<script>`、外部 CDN、外部 font、画像 preload を出力しない。

---

## 6. CSS クラス一覧

### コンテンツ要素

| クラス | 要素 | 説明 |
|--------|------|------|
| `.mh` | `h1`〜`h4` | 見出し共通スタイル |
| `.h1`〜`.h4` | `h1`〜`h4` | 見出しレベル別スタイル |
| `.mp` | `p` | 本文段落（`max-width: 68ch`） |
| `.mr` | `hr` | 水平線 |
| `.mbq` | `blockquote` | 引用ブロック |
| `.ic` | `code` | インラインコード |
| `.ml` | `ul`/`ol` | リスト |

### コードブロック

| クラス | 説明 |
|--------|------|
| `.cb-wrap` | コードブロック外枠（`position: relative`） |
| `.cb-meta` | 言語ラベル＋コピーボタンのオーバーレイ（`position: absolute; top: 8px; right: 10px`） |
| `.cl` | 言語ラベル（`var(--adlaire-font-family-mono)`、大文字、`var(--adlaire-font-size-xs)`） |
| `.cb-copy` | コピーボタン（ホバーまで非表示） |
| `.cb-copy.copied` | コピー完了状態（1.8 秒間） |
| `.cb` | `<pre>` 要素（背景 `var(--adlaire-surface-soft-strong)`） |

### タスクリスト

| クラス | 説明 |
|--------|------|
| `.ml-task` | タスクリスト `<li>`（`list-style: none`、チェックボックス付き） |
| `.ml-task input[type="checkbox"]` | チェックボックス（`disabled`、`accent-color: var(--adlaire-color-primary)`） |

### 定義リスト

| クラス | 説明 |
|--------|------|
| `.mdl` | `<dl>` 要素 |
| `.mdl dt` | 定義用語（`font-weight: semibold`） |
| `.mdl dd` | 定義本文（`margin-left: 1.5rem`、`color: var(--adlaire-surface-text-muted)`） |

### テーブル

| クラス | 説明 |
|--------|------|
| `.tw` | テーブルラッパー（`overflow-x: auto`） |
| `.mt` | `<table>` 要素 |

### 脚注

| クラス | 説明 |
|--------|------|
| `.fn-ref` | 脚注参照リンク（`<sup>` 内、`font-size: var(--adlaire-font-size-xs)`） |
| `.fn-section` | 脚注セクション全体（本文末尾、`border-top` で区切り） |
| `.fn-list` | 脚注 `<ol>`（`font-size: var(--adlaire-font-size-sm)`） |
| `.fn-item` | 脚注 `<li>` |
| `.fn-back` | 本文への戻りリンク（`↩`） |

### TOC（サイドバー）

| クラス | 説明 |
|--------|------|
| `.tr` | TOC ルート `<ul>` |
| `.tg` | グループ（子を持つ見出し） |
| `.tg-row` | グループの行（リンク＋展開ボタン） |
| `.tg-btn` | グループ展開トグルボタン |
| `.tc` | グループの子 `<ul>`（`hidden` で折りたたみ） |
| `.ti` | リーフ（子を持たない見出し） |
| `.tl` | TOC リンク（`<a>` 要素） |
| `.tl.lv1`〜`.tl.lv3` | TOC リンクのレベル別スタイル |
| `.tl.active` | アクティブな TOC リンク |
| `.sb-none` | 検索結果なしメッセージ |

### シンタックスハイライト CSS クラス（§7.8）

| クラス | 対象トークン |
|--------|------------|
| `.hl-kw` | キーワード |
| `.hl-str` | 文字列リテラル |
| `.hl-num` | 数値リテラル |
| `.hl-cmt` | コメント |
| `.hl-key` | JSON オブジェクトキー |
| `.hl-op` | diff 追加行（`+`）・削除行（`-`） |

### 見出しアンカーリンクコピー CSS クラス（§7.11）

| クラス | 要素 | 説明 |
|--------|------|------|
| `.hn-link` | `<button>` | 見出し末尾に付与する `¶` ボタン（通常時 `opacity: 0`、ホバー時 `opacity: 1`） |

### 進捗バー・テーブルソート CSS クラス（§7.13・§7.14）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `#progress-bar` | `<div>` | ページ上端固定（`position: fixed; top: 0; left: 0`）、高さ 3px、幅は JS で設定、`background: var(--adlaire-color-primary)`、`z-index: 1000`、`transition: width 0.1s linear` |
| `.mt th[data-sort]` | `<th>` | ソート可能列ヘッダー（`cursor: pointer; user-select: none`） |
| `.mt th[aria-sort="ascending"]::after` | `::after` 疑似要素 | `content: " ▲"` |
| `.mt th[aria-sort="descending"]::after` | `::after` 疑似要素 | `content: " ▼"` |

### 読了時間表示 CSS クラス（§4.5）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `#reading-time` | `<span>` | 固定ヘッダー右端に表示（`margin-left: auto`）。`color: var(--adlaire-text-secondary)`、`font-size: var(--adlaire-font-size-sm)`、`white-space: nowrap` |

### 前後章ナビゲーション CSS クラス（§7.15）

| クラス / セレクター | 要素 | 説明 |
|--------------------|------|------|
| `.ch-nav` | `<nav>` | 章末尾ナビゲーション（`display: flex; justify-content: space-between; padding: 1rem 0; margin-top: 2rem; border-top: 1px solid var(--adlaire-border-default)`） |
| `.ch-prev` | `<a>` | 前の章リンク（`color: var(--adlaire-color-primary)`、テキスト装飾なし） |
| `.ch-next` | `<a>` | 次の章リンク（`color: var(--adlaire-color-primary)`、テキスト装飾なし） |

### 印刷スタイル（§6 @media print）

`@media print` ブロックでの主な規則：

| 対象 | 印刷時の処理 |
|------|------------|
| `#hdr`（固定ヘッダー）、`#sb`（サイドバー）、`.cb-copy`、`.hn-link`、`#btt`、`.expand-code`、`#progress-bar`、`.ch-nav` | `display: none` |
| `#lay`、`#main` | ブロック表示・幅 100%・余白 0 |
| `pre.cb[data-collapsible]` | `max-height: none`（折りたたみ解除） |
| `a[href^="http"]::after`、`a[href^="https"]::after` | `content: " (" attr(href) ")"` で URL を末尾に表示 |
| `h2`、`h3` | `page-break-before: avoid` |

---

## 7. JavaScript 機能

### 7.1 テーマ切り替え（廃止）

ADS 採用により、ダークモードおよびテーマトグルボタンは廃止。
出力サイトはライトモード固定（`prefers-color-scheme` 非対応）。

### 7.2 サイドバー開閉

`localStorage` キー `adb-sb` に `"1"`（開）または `"0"`（閉）を保存する。

- デスクトップ（`> 768px`）：`#sb.closed` クラスと JS インラインスタイル（`ct.style.marginLeft`）で幅を制御。CSS の `#sb.closed ~ #ct { margin-left: 0 }` ルールは初期レンダリング時のみ効く。以降の開閉操作はすべて JS インラインスタイルが CSS クラスより優先する
- モバイル（`≤ 768px`）：`#sb.open` / `transform: translateX` で画面外から引き出す
- モバイルでは TOC リンククリック時に自動的にサイドバーを閉じる

`localStorage` 読み書きは `try/catch` で保護する。読み込み失敗、保存失敗、保存値が `"1"` / `"0"` 以外の場合は、デスクトップでは開、モバイルでは閉を初期状態とする。`assets/app.js` は localStorage 以外の永続 storage、cookie、IndexedDB を使用してはならない。

### 7.3 TOC グループ展開

`.tg-btn` クリックで対応する `<ul id="tg-{slug}">` の `hidden` 属性をトグルし、`aria-expanded` 属性を更新する。`data-target` 属性で対象 `<ul>` の ID を指定する。

**開閉状態の永続化：**
展開操作のたびに現在展開中のグループのスラグ配列を `localStorage` キー `adlaire-toc-state` に JSON 文字列で保存する。ページ読み込み時（`DOMContentLoaded`）に同キーを読み込み、保存済みスラグのグループを展開状態で描画する。`localStorage` アクセスはすべて `try/catch` で保護し、失敗時はデフォルト状態（初期展開なし）にフォールバックする。

保存値が JSON 配列でない場合、存在しない slug を含む場合、100 件を超える場合は保存値を無視し、書き戻しは行わない。`aria-expanded` と `hidden` は常に逆状態に保ち、`aria-expanded="true"` のとき対象 `<ul>` から `hidden` を外す。

### 7.4 TOC 検索フィルター

`#sb-search` の `input` イベントで TOC 項目をリアルタイムフィルタリングする。

- 各 `.tl` のテキストコンテンツと検索クエリ（小文字化）を照合
- 非一致の `<li>` は `hidden = true` で非表示
- グループは子に一致項目があれば表示を維持し、`<ul>` を強制展開
- 入力がクリアされた場合は全項目を復元し、展開状態を `aria-expanded` に従って復元
- 一致なしの場合は `#sb-none` メッセージを表示

**状態変数 `searchActive`：** 検索中はアクティブ TOC リンクの自動スクロールを抑制する。

### 7.5 アクティブ見出し追跡

`IntersectionObserver` で `.mh[id]` 要素のビューポート内への進入を監視する。

- `rootMargin: "-8% 0px -78% 0px"` — 画面上部 8% 〜 22% の帯域内の見出しを「アクティブ」とみなす
- 進入した見出しの `id` に対応する TOC リンクに `.active` クラスを付与
- 親グループが折りたたまれている場合は自動展開し `aria-expanded="true"` を設定
- 検索中でない場合は対応 TOC リンクを `scrollIntoView` で可視範囲にスクロール

**進捗バー連動：** `scroll` イベントリスナー（`passive: true`）を同一リスナーで共有し、スクロールのたびに読み取り進捗バーの幅を更新する（→ §7.13）。

### 7.6 コピーボタン

各 `.cb-copy` ボタンのクリックで、親 `.cb-wrap` 内の `<code>` の `innerText` をクリップボードにコピーする。

```
navigator.clipboard.writeText() → 成功: done()
                                 → 失敗: fallbackCopy() → done()
fallbackCopy(): <textarea> を一時生成して execCommand('copy') を実行
done(): ボタンテキストを "✓ 完了" に変更、.copied クラス付与、1.8 秒後に "コピー" に戻す
```

### 7.7 トップへ戻るボタン

`scroll` イベント（`passive: true`）で `scrollY > 400` の場合に `#btt.visible` クラスを付与する。クリックで `window.scrollTo({ top: 0, behavior: 'smooth' })` を実行する。

### 7.8 シンタックスハイライト

コードブロックに言語別の色分けを `assets/app.js` で適用する。外部ライブラリ不要。

**対応言語：** `python` / `bash` / `json` / `sql` / `ini` / `diff`

**実装方式：**
各言語ごとにトークン正規表現パターンを定義し、`<code>` 要素のテキストに対して順次マッチを走らせる。マッチしたトークンを `<span class="hl-{type}">` でラップしてから `innerHTML` に書き戻す。

| CSS クラス | 対象トークン |
|---|---|
| `hl-kw` | キーワード（`def`, `class`, `if`, `SELECT` 等） |
| `hl-str` | 文字列リテラル（シングル／ダブルクォート） |
| `hl-num` | 数値リテラル |
| `hl-cmt` | コメント（`#`・`//`・`--` 行コメント） |
| `hl-key` | オブジェクトキー（JSON の `"key":` パターン） |
| `hl-op` | diff の追加行（`+`）・削除行（`-`） |

**適用タイミング：** `DOMContentLoaded` 後に全 `.cb-wrap` 要素を走査し、`data-lang` 属性の値を参照して言語を判定する（`data-lang` は `.cb-wrap` div に付与されており、`<code>` 要素には付与されない）。ハイライト処理は配下の `<code>` 要素のテキストに対して行う。`data-lang` が未設定または対応外の場合はハイライトをスキップする。

### 7.9 本文内全文検索

ビルド時に `assets/search-index.json` を生成し、`assets/app.js` の検索 UI から読み込む。TOC 検索フィルター（§7.4）と検索 UI を統合し、本文ヒット箇所へのジャンプを提供する。

**インデックス生成仕様（builder）：**
ビルド時に全ページの見出しと各段落の先頭 200 文字を抽出し、以下の配列形式で `assets/search-index.json` に書き出す。

```json
[
  { "url": "pages/example.html#anchor-slug", "id": "anchor-slug", "title": "見出しテキスト", "body": "段落先頭200文字..." },
  ...
]
```

**検索 index 固定契約：**

- JSON の top-level は配列とする。object wrapper は使用しない。
- entry の key 順は `url`、`id`、`title`、`body` とする。
- `encoding/json` の標準エスケープを使用し、独自 pretty print は行わない。
- entry の並び順は、ページ入力順、同一ページ内では見出し出現順とする。ページ単位エントリは各ページの先頭に 1 件だけ追加する。
- `body` は HTML タグ除去後のプレーンテキストを Unicode rune 単位で最大 200 文字とする。200 文字を超える場合は 200 文字で切り、三点リーダーを追加しない。
- 空白、タブ、改行の連続は半角スペース 1 つへ正規化し、前後空白を削除する。
- `title` と `body` に HTML entity を残してはならない。検索 index は表示前に JS 側で `textContent` として挿入し、`innerHTML` へ直接挿入しない。

**検索 UI の配置：**
§7.4 の TOC 検索フィルター入力欄を兼用する。入力値が 2 文字以上になった時点でインデックスに対して部分一致検索を実行する。

**ヒット箇所ハイライト：**
一致したエントリの見出しを検索結果として表示する。現在ページ内のヒットは TOC 内でハイライト（`.toc-hit` クラス付与）し、クリックで対象アンカーへスクロールする。別ページのヒットは `url` へ遷移する。現在ページ内では `<mark>` 要素でヒット文字列をページ内マーキングする（外部依存なし・標準 DOM 操作のみ）。

**クリア：**
入力欄を空にすると TOC ハイライトおよびページ内マーキングをすべて解除する。

**検索 UI 実行時契約：**

| 項目 | 仕様 |
|------|------|
| fetch path | 現在 HTML の `relativeRoot + "assets/search-index.json"`。 |
| fetch 失敗 | TOC 検索だけを継続し、console に `search index unavailable` を 1 回だけ出す。 |
| 最小文字数 | 2 文字未満は index 検索を実行しない。TOC filter は 1 文字から実行する。 |
| 最大結果 | 20 件。entry 順を保持する。 |
| 表示先 | `#sb-toc` 内に `<div id="search-results">` を 1 個だけ作成し、結果更新時に中身を置換する。 |
| 挿入方法 | result title/body/url は `textContent` または `setAttribute` で設定し、`innerHTML` へ検索 index 由来値を入れない。 |
| mark 解除 | 検索ごとに前回の `<mark data-search-hit="true">` を text node へ戻してから新規 mark を挿入する。 |

### 7.10 コードブロック折りたたみ

30 行超のコードブロックを初期折りたたみ状態でレンダリングし、ユーザー操作で全行表示に切り替える。

**閾値：** 30 行（空行を含む総行数）。30 行以下のブロックは折りたたみ UI を生成しない。

**初期状態：**
折りたたみ対象の `<pre>` に `data-collapsible="true"` と `data-total-lines="{N}"` を付与する。CSS で `max-height` を制限（先頭 10 行相当）し、下端をグラデーションフェードでマスクする。

**展開リンク：**
ブロック末尾に `<button class="expand-code">全 {N} 行を表示</button>` を配置する。クリックで `max-height` を解除し、ボタンを非表示にする（折りたたみへの再折りたたみ機能は提供しない）。

**コピーボタンとの共存：**
§7.6 のコピーボタンは折りたたみ状態でも常時表示する。コピー操作は全行テキストを対象とする（表示行のみではない）。

### 7.11 見出しアンカーリンクコピー

見出しにホバーすると表示される `¶` ボタン（`.hn-link`）をクリックすると、その見出しのアンカー URL をクリップボードにコピーする。

**コピー対象 URL：**
`window.location.origin + window.location.pathname + button.dataset.href`

**実装：**
§7.6 と同じ Clipboard API（`navigator.clipboard.writeText()`）を使用し、同一のフォールバック（`execCommand('copy')`）を流用する。コピー完了フィードバックは `aria-label` を `"コピーしました"` に一時変更し、1.8 秒後に `"リンクをコピー"` へ復元する（§7.6 のコピーボタンと同じタイミング）。

**イベント登録：**
`DOMContentLoaded` 後に `document.querySelectorAll('.hn-link')` を走査して `click` リスナーを登録する。

### 7.12 キーボードショートカット

`keydown` イベントで以下のショートカットを処理する。フォーカスが `<input>`・`<textarea>`・`<select>` にある場合は `/` と `t` を無効にする（`Escape` のみ有効）。

| キー | 動作 |
|------|------|
| `/` | TOC 検索欄（`#sb-search`）にフォーカスを移動し、ページスクロールを抑止（`event.preventDefault()`） |
| `Escape` | TOC 検索欄の内容をクリアし、検索結果をリセットする（§7.4 の検索リセット処理と同等） |
| `t` | ページ先頭へスクロール（`window.scrollTo({ top: 0, behavior: 'smooth' })`） |

**イベント登録：** `document.addEventListener('keydown', handler)` を `DOMContentLoaded` 後に登録する。

### 7.13 読み取り進捗バー

ページ上端に高さ 3px の進捗バー（`<div id="progress-bar">`）を固定表示する（→ §5 HTML 出力構造）。

**幅の計算：**
```js
const scrolled = document.documentElement.scrollTop;
const total    = document.documentElement.scrollHeight
                 - document.documentElement.clientHeight;
const pct      = total > 0 ? (scrolled / total * 100) : 0;
document.getElementById('progress-bar').style.width = pct + '%';
```

**スクロールイベント共有：** §7.5 のアクティブ見出し追跡が登録する `scroll` イベントリスナー（`passive: true`）内で処理する。リスナーを別途登録しない。

**CSS：**

```css
#progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  width: 0%;
  background: var(--adlaire-color-primary);
  z-index: 1000;
  transition: width 0.1s linear;
}
```

### 7.14 テーブル列ソート

`<th>` クリックで列を昇順／降順ソートする。

**HTML 構造：** `flushTable()` が生成する全 `<th>` に `data-sort="{col_index}"` 属性（0始まりの列インデックス）と `aria-sort="none"` を付与する。

```html
<th data-sort="0" aria-sort="none">列名</th>
```

**ソートアルゴリズム：**

1. クリックされた `<th>` の `aria-sort` を確認し、`"ascending"` → `"descending"`、それ以外 → `"ascending"` に決定
2. 同テーブル内の他の全 `<th>` の `aria-sort` を `"none"` にリセット
3. 対象列のセルテキスト（`textContent.trim()`）を取得し、`Number()` で数値変換可能ならば数値比較、それ以外は文字列比較（`localeCompare`）でソート
4. `<tbody>` の子 `<tr>` を並び替えて再挿入
5. クリックされた `<th>` の `aria-sort` を決定した値に更新（CSS `::after` でインジケーター表示）

**イベント登録：** `DOMContentLoaded` 後に `document.querySelectorAll('.mt th[data-sort]')` を走査して `click` リスナーを登録する。

**スコープ：** 同一テーブル内のソートのみ。複数列ソートは対象外。

同値比較の場合は元の行順を保持する安定ソートとする。数値比較では空文字は文字列として扱い、`Number("")` による `0` 扱いを禁止する。

---

### 7.15 前後章ナビゲーションボタン

h2 見出し単位で「← 前の章」「次の章 →」ボタンを各章末尾に静的生成する（→ §4.5・§5・§6 CSS）。

**生成方法：** `injectChapterNavigation(bodyHTML string, headings []Heading) string` として実装する。`convert()` が `ConvertResult.HTML` を返した後、呼び出し元が h2 見出しの `Slug` と `Text` を `headings` から抽出し、各 h2 章の末尾（次の h2 の直前、または文書末）に `<nav class="ch-nav">` を挿入する。

`injectChapterNavigation()` は `Heading.Level == 2` の見出しだけを対象とする。対象 h2 が 0 件または 1 件の場合、章ナビゲーションを挿入しない。

**HTML 構造：**

```html
<!-- 先頭章（.ch-prev なし） -->
<nav class="ch-nav">
  <span></span>
  <a class="ch-next" href="#next-slug">次の章タイトル →</a>
</nav>

<!-- 中間章 -->
<nav class="ch-nav">
  <a class="ch-prev" href="#prev-slug">← 前の章タイトル</a>
  <a class="ch-next" href="#next-slug">次の章タイトル →</a>
</nav>

<!-- 最終章（.ch-next なし） -->
<nav class="ch-nav">
  <a class="ch-prev" href="#prev-slug">← 前の章タイトル</a>
  <span></span>
</nav>
```

**スラグの参照：** `href` に使用するスラグは §4.5 のスラグ重複解決後の一意スラグを使用する。

**スコープ：** h2 レベルの見出しのみ。h3 以下の小節には生成しない。

**挿入失敗時の扱い：** 対象 h2 の HTML 位置を特定できない場合、本文 HTML を変更せず `[WARN] CHAPTER_NAV_SKIPPED: slug={slug}` を出力し、`[REPORT] warnings` に含める。

**印刷時：** `@media print` で `.ch-nav { display: none }` とする（§6 CSS 参照）。

**JavaScript 初期化順序固定契約：**

`assets/app.js` は `DOMContentLoaded` 後に以下の順で初期化する。

1. 必須 DOM 参照を取得し、存在しない要素があっても例外で停止せず該当機能だけ無効化する。
2. sidebar 状態を復元する。
3. TOC group 状態を復元する。
4. TOC 検索と search index fetch を初期化する。
5. syntax highlight を適用する。
6. copy button、heading anchor、table sort、code expand、top button、keyboard shortcut を登録する。
7. IntersectionObserver と scroll handler を登録し、進捗バーを 1 回更新する。

初期化中の 1 機能の失敗で他機能を停止してはならない。catch した例外は `console.warn("adlaire static init failed", name)` の形式で機能名だけを出し、Markdown 本文、search query、secret 相当値を出力しない。

---

## 8. 実行方法

```bash
/usr/local/bin/adlaire-ci-build
```

`adlaire-ci-build` の実行は、§2 の CLI 引数仕様に従う。引数なしの場合は `DefaultBuildConfig` の `Src`、`Out`、`Title`、`Theme`、`BaseDir`、`Strict` を使用する。

**実行順序契約：**

1. CLI 引数を検証する。`--help` / `--version` はここで処理し、Markdown 読み込みを行わない。
2. `--theme` を検証し、`adlaire-default` の `Theme` を選択する。
3. `--src` がファイルの場合は 1 ページ、ディレクトリの場合は配下の `.md` ファイルを辞書順に収集する。
4. 各 Markdown を UTF-8 として読み込み、行配列 `lines []string` を作成する。
5. 各ページの見出しを収集し、`[]Heading` と `slugByLine` を作成する。
6. 各ページの脚注定義を収集し、`RenderContext` を初期化する。
7. `buildTOC(headings)` で `PageData.TocHTML` を作成する。
8. `convert(lines, headings, ctx)` で `ConvertResult` を作成する。
9. `injectChapterNavigation(result.HTML, headings)` を適用し、`PageData.BodyHTML` を確定する。
10. サイト全体の検索インデックスを `encoding/json` で生成し、`SiteData.SearchIndex` を確定する。
11. `assembleSite(siteData, theme)` で `index.html`、ページ HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を生成する。
12. 出力サイトディレクトリのファイル数と合計サイズを取得し、stdout に完了行と `[REPORT]` 行を出力する。

途中で失敗した場合は、失敗段階以降を実行しない。一時ファイルが存在する場合は削除してから終了する。

**終了コード：**

| 終了コード | 条件 | 後続処理 |
|------------|------|----------|
| `0` | 静的 Web サイト生成に成功し、`[REPORT]` 行を出力した。 | `runner` は成功として扱う。 |
| `1` | 出力ディレクトリ作成、HTML / CSS / JavaScript / search index 書き込み、テンプレート合成など処理中の一般エラー。 | `runner` はビルド失敗として扱い、SHA を更新しない。 |
| `2` | CLI 引数不正、入力ファイル不存在、入力 UTF-8 不正、または `--strict` 指定時の警告発生。 | `runner` は設定または入力エラーとして扱い、SHA を更新しない。ただし `--strict` 警告時は `[REPORT]` を取り込む。 |

終了コード `0` の場合、stdout には必ず `Collecting Markdown...`、`Converting MD...`、`Building site...`、`Writing assets...`、`Done → ...`、`[REPORT] ...` をこの順序で出力する。警告がある場合は `[REPORT]` の直前に `[WARN] ...` を 1 件 1 行で出力する。

終了コード `1` または `2` の場合、stderr に原因を 1 行以上出力し、`[REPORT]` 行は出力しない。ただし `--strict` 警告による終了コード `2` だけは、stdout に通常進捗、`[WARN]`、`[REPORT]` を出力し、stderr は空とする。途中まで作成した出力サイトは公開用パスへ残してはならず、一時ディレクトリを削除して終了する。

**CLI 出力固定契約：**

| ケース | stdout | stderr | 終了コード | 副作用 |
|--------|--------|--------|------------|--------|
| `--help` | `Usage: adlaire-ci-build [--src path] [--out path] [--title text] [--theme name] [--base-dir path] [--strict] [--build-id id] [--commit-sha sha] [--build-at iso8601] [--version] [--help]` + LF | 空 | `0` | 入力読込、出力作成なし。 |
| `--version` | `adlaire-ci-build ADLAIRE_CI_SPEC go={version}` + LF | 空 | `0` | 入力読込、出力作成なし。 |
| 引数不正 | 空 | 固定エラー 1 行 + LF | `2` | 入力読込、出力作成なし。 |
| 入力不存在 | 空 | `source not found: {path}` + LF | `2` | 出力作成なし。 |
| UTF-8 不正 | 空 | `source is not valid UTF-8: {path}` + LF | `2` | 出力作成なし。 |
| strict 警告あり | 通常進捗、`[WARN]`、`[REPORT]` | 空 | `2` | 一時出力は削除し、公開用 `--out` は置換しない。 |
| 書込失敗 | 失敗前までの進捗行 | `write failed: {path}` + LF | `1` | 一時出力を削除し、公開用 `--out` は置換しない。 |

進捗 stdout は LF 改行だけを使用する。`Done → ...` 行の path は `--out` の解決後絶対パスとする。`--help` と `--version` が同時指定された場合は `--help` を優先する。`--strict` で警告が発生した場合も `[REPORT]` は出力し、runner は終了コード `2` と `[REPORT]` の両方を保存する。

**一時出力・置換契約：**

`adlaire-ci-build` は公開用 `--out` へ直接書き込まず、同一親ディレクトリに `{out}.tmp.{pid}` を作成して全ファイルを書き込む。全ファイルの write、close、sync、検索 index 生成、asset 生成が成功した場合だけ、既存 `--out` を `{out}.previous.{pid}` へ rename し、tmp を `--out` へ rename する。rename 後に親ディレクトリを sync する。置換成功後、旧 directory を削除する。置換前に失敗した場合は tmp だけ削除し、既存 `--out` を保持する。置換後の旧 directory 削除に失敗した場合は WARN を出すが終了コードは `0` のままとする。

**標準出力：**
```
Collecting Markdown...
Converting MD...
Building site...
Writing assets...
Done → /opt/adlaire-builder/dist/site  (pages=12 files=15 bytes=1713731)
[REPORT] pages=12 headings=342 tables=128 code_blocks=64 warnings=3 size_warn=false broken_links=1 heading_skips=0 reading_time=87 theme=adlaire-default
```

**変換レポート行（`[REPORT]` プレフィックス）：**
ビルド完了直後に 1 行で出力する。フィールドはスペース区切りの `key=value` 形式で固定順。

| フィールド | 内容 |
|---|---|
| `pages` | 出力した HTML ページ数 |
| `headings` | 出力サイト内の見出し要素（`h1`〜`h6`）の総数 |
| `tables` | 変換したテーブルの総数 |
| `code_blocks` | 変換したコードブロックの総数 |
| `warnings` | ビルド中に発生した警告件数 |
| `size_warn` | 出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` を超えた場合 `true`、それ以外 `false`（`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`） |
| `broken_links` | 参照先スラグが存在しない内部アンカーリンク（`[label](#anchor)`）の件数 |
| `heading_skips` | 見出しレベルが 2 段以上の降順スキップとなった件数 |
| `reading_time` | 推計読了時間（分、切り上げ）。200文字/分で算出 |
| `theme` | 使用した theme 名。初期仕様では `adlaire-default` |
| `build_id` | `--build-id` または環境変数 `ADLAIRE_BUILD_ID` の値。未指定時は空文字 |
| `commit_sha` | `--commit-sha` または環境変数 `ADLAIRE_COMMIT_SHA` の値。未指定時は空文字 |
| `build_at` | `--build-at` または環境変数 `ADLAIRE_BUILD_AT` の値。未指定時は空文字 |

固定順は以下とし、未使用フィールドの省略は禁止する。

```text
pages headings tables code_blocks warnings size_warn broken_links heading_skips reading_time theme build_id commit_sha build_at
```

`warnings` は出力した `[WARN]` 行数と一致しなければならない。`reading_time` は `ConvertResult.ReadingTimeMinutes`、`broken_links` は `ConvertResult.BrokenLinks`、`heading_skips` は `ConvertResult.HeadingSkips` を使用する。

警告が発生した場合、`[REPORT]` 行の直前に `[WARN] {メッセージ}` 形式で 1 件ずつ出力する。

**`runner` による取り込み：**
Go 版 CI ランナーでは、`runner` が `pipeline.sh` の標準出力から `[REPORT]` 行と `[WARN]` 行を抽出し、パースした結果を `.build_logs/{id}.json` のビルドログエントリに追記する。

```json
{
  "build_id": "b20260915100000",
  "status": "success",
  "report": {
    "headings": 342,
    "tables_count": 128,
    "code_blocks_count": 64,
    "warnings_count": 3,
    "size_warn": false,
    "broken_links": 1,
    "heading_skips": 0,
    "reading_time": 87,
    "theme": "adlaire-default",
    "build_id": "b20260915100000",
    "commit_sha": "abc1234",
    "build_at": "2026-09-15T10:00:00Z"
  }
}
```

> **フィールド名の対応：** stdout の `[REPORT]` 行は `tables=` / `code_blocks=` の短縮キーを使用するが、`.build_logs/{id}.json` への保存時および `GET /api/output-meta` レスポンスでは `tables_count` / `code_blocks_count` に変換する（→ §22）。

**再実行時の注意：** スラグ重複カウンタ、脚注参照順、脚注定義は `adlaire-ci-build` の 1 実行内で初期化する。通常の `/usr/local/bin/adlaire-ci-build` 実行では複数回実行しても出力は同一になる。Go 版では変換状態をパッケージグローバル変数として共有せず、変換処理ごとに専用の状態構造体を生成する。

---

## 8a. `builder` 受け入れ fixture

`builder` の初期実装は、本節の fixture をすべて満たすまで完了として扱わない。fixture ファイルは実装 PR で `testdata/builder/` 配下へ追加する。仕様 PR では fixture の期待値を本節で固定する。

### Fixture A: 単一 Markdown 入力

**入力ファイル：** `testdata/builder/single/source.md`

~~~markdown
# Title

Intro paragraph with [self](#title).

## Install

```bash
echo hello
```

- [x] done
- [ ] todo

| Name | Value |
| ---- | ----- |
| A | 1 |

[^n]: note body

See footnote[^n].
~~~

**実行：**

```bash
adlaire-ci-build --src testdata/builder/single/source.md --out /tmp/adlaire-ci-fixture-single --title "Fixture Site"
```

**期待結果：**

- 終了コード `0`。
- `/tmp/adlaire-ci-fixture-single/index.html`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` が存在する。
- `pages/` は存在しない。
- `index.html` に `<h1 id="title" class="mh h1">Title<button class="hn-link" data-href="#title" aria-label="リンクをコピー">¶</button></h1>` を含む。
- `index.html` に `<a href="#title">self</a>` を含み、`BROKEN_LINK` 警告を出さない。
- `index.html` に `<div class="cb-wrap" data-lang="bash">` と `<span class="cl">bash</span>` を含む。
- `index.html` に `<li class="ml-task"><input type="checkbox" disabled checked>done</li>` と `<li class="ml-task"><input type="checkbox" disabled>todo</li>` を含む。
- `[REPORT]` は `pages=1`、`theme=adlaire-default` を含む。

### Fixture B: ディレクトリ Markdown 入力

**入力ファイル：**

```text
testdata/builder/site/docs/intro.md
testdata/builder/site/docs/guide/setup.md
testdata/builder/site/docs/guide/setup_copy.md
```

`intro.md`:

```markdown
# Intro

Go to [setup](guide/setup.md#setup).
```

`guide/setup.md`:

```markdown
# Setup

Body.
```

`guide/setup_copy.md`:

```markdown
# Setup

Second.
```

**実行：**

```bash
adlaire-ci-build --src testdata/builder/site/docs --out /tmp/adlaire-ci-fixture-site --title "Docs"
```

**期待結果：**

- 終了コード `0`。
- `/tmp/adlaire-ci-fixture-site/index.html` がサイト目次ページである。
- document ページは `pages/guide-setup.html`、`pages/guide-setup-copy.html`、`pages/intro.html` として出力される。
- 入力順は `guide/setup.md`、`guide/setup_copy.md`、`intro.md` の辞書順とし、サイト目次の表示順も同一とする。
- `intro.html` 内の `guide/setup.md#setup` は、同じ `pages/` ディレクトリ内のページ間リンクとして `guide-setup.html#setup` に変換する。`guide/setup.html#setup`、`pages/guide-setup.html#setup`、元の `guide/setup.md#setup` のまま出力してはならない。
- 2 つの `# Setup` 見出しはページごとの slug 空間でそれぞれ `setup` とする。ページをまたいだ見出し slug に `-2` を付けてはならない。
- `assets/search-index.json` は `index.html#site-index`、`pages/guide-setup.html#setup`、`pages/guide-setup-copy.html#setup`、`pages/intro.html#intro` の entry を含む。

### Fixture C: 異常系

| 実行 | 終了コード | stderr |
|------|------------|--------|
| `adlaire-ci-build --theme unknown` | `2` | `unknown theme: unknown` |
| `adlaire-ci-build --src /path/not-found.md` | `2` | `source not found: /path/not-found.md` |
| `adlaire-ci-build --src testdata/builder/empty-dir` | `2` | `no markdown files found: testdata/builder/empty-dir` |
| `adlaire-ci-build --title ""` | `2` | `title must not be empty` |

異常系 fixture では `[REPORT]` を stdout へ出力してはならない。`--out` に既存の正常出力がある場合でも、異常系実行で既存出力を変更してはならない。

### Fixture D: 冪等性

同一入力、同一 CLI 引数で 2 回連続実行した場合、`GeneratedAtUTC` を含む meta 行を除き、全出力ファイルの内容が一致しなければならない。比較対象から除外できるのは、HTML 内の `name="adlaire-generated-at"` meta と footer の生成時刻表示だけとする。検索 index、ページ HTML の本文、CSS、JS、REPORT の数値は一致必須とする。

### Fixture E: path 安全性と既存出力保護

**事前状態：**

`/tmp/adlaire-ci-fixture-safe/index.html` に `previous output` を含む正常出力を作成しておく。

**実行と期待結果：**

| 実行 | 終了コード | 期待結果 |
|------|------------|----------|
| `adlaire-ci-build --src testdata/builder/site/docs --out testdata/builder/site/docs/out` | `2` | stderr `output path must be outside source: <path>`、出力作成なし。 |
| `adlaire-ci-build --src /tmp/adlaire-ci-fixture-safe --out /tmp/adlaire-ci-fixture-safe` | `2` | stderr `output path must be outside source: <path>`、既存 `index.html` 維持。 |
| 10 MiB 超の Markdown file を `--src` に指定 | `2` | stderr `source file too large: <path>`、`[REPORT]` なし。 |

### Fixture F: HTML escape と Markdown 境界

**入力：**

```markdown
# Unsafe

<script>alert(1)</script>

| A | B |
| - | - |
| 1 |
| 2 | 3 | 4 |

- item
        - too deep
              - deeper
```

**期待結果：**

- `<script>` は実行可能 tag にならず、`&lt;script&gt;alert(1)&lt;/script&gt;` として出力される。
- テーブル不足セルは、欠落 cell の個数分だけ空 `<td></td>` を出力する。
- テーブル超過セルは最後のセルに `3 | 4` として連結される。
- 7 レベル以上の list nesting は 6 レベルへ丸められ、`[WARN] LIST_NESTING_CLAMPED` が出る。

### Fixture G: search index / JavaScript contract

**入力：** Fixture B と同じ directory 入力。

**期待結果：**

- `assets/search-index.json` は top-level array で、entry key 順が `url`、`id`、`title`、`body`。
- `body` は HTML tag を含まず、200 文字を超えない。
- `assets/app.js` に `localStorage` access の `try` / `catch`、`search-results`、`data-search-hit`、`search index unavailable` が含まれる。
- `assets/app.js` に `document.cookie`、`indexedDB`、外部 URL fetch が含まれない。

### Fixture H: strict warning and atomic output

**事前状態：**

`/tmp/adlaire-ci-fixture-strict/index.html` に `previous output` を含む正常出力を作成しておく。

**入力：**

~~~markdown
# Title

[missing](#does-not-exist)

```bash
echo unclosed
~~~

上記 fixture では実ファイル上の fence を閉じずに EOF とする。

**実行：**

```bash
adlaire-ci-build --src testdata/builder/strict/source.md --out /tmp/adlaire-ci-fixture-strict --strict
```

**期待結果：**

- 終了コード `2`。
- stdout に `[WARN] BROKEN_LINK` と `[WARN] UNCLOSED_FENCE` と `[REPORT]` を出力する。
- stderr は空。
- `/tmp/adlaire-ci-fixture-strict/index.html` は `previous output` のままで置換されない。

**§8〜§20 builder / runner 中核機能別実装完全性固定契約：**

§8〜§20 の中核機能は、各節の本文と fixture に加えて下表を満たした場合だけ実装完了とする。下表は既存機能の実装時チェックリストであり、将来機能、MCP、外部公開構成、上位方針は扱わない。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §8 | builder CLI 実行 | CLI 引数、Markdown file / directory、theme、build meta。 | 静的 Web サイト、stdout 進捗、`[REPORT]`。 | 公開用 `--out` は tmp 完成後だけ置換する。 | 引数不正、UTF-8 不正、strict 警告、書込失敗時は既存出力を保持する。 | help/version、単一入力、directory 入力、strict、atomic output。 |
| §8a | builder fixture | `testdata/builder/` 入力一式。 | expected HTML / CSS / JS / search index / stdout / stderr。 | fixture 実行時だけ一時出力を作成する。 | 異常系 fixture で `[REPORT]` を出さず既存出力を変えない。 | Fixture A〜H 全件。 |
| §10〜§12 | runner 起動 / 設定 | `--state-dir`、secret、branch target、server config、systemd oneshot。 | runner 終了コード、slog、正規化設定。 | 必須検証成功後だけ lock / state を更新する。 | secret 不足、設定不正、insecure mode では build を開始しない。 | secret 不足、token mode 不正、branch config default、dry-run directory 作成なし。 |
| §13 | runner 処理フロー | SHA cache、GitHub API、build queue、cooldown、circuit、trigger。 | build log、history、status、queue 更新。 | finalizer で lock/state/status を固定順に更新する。 | GitHub 全失敗、lock 不正、state write 失敗時の副作用を固定する。 | R1〜R17、R21〜R27。 |
| §14 | pipeline 起動 | `.ci/pipeline.sh`、builder binary、timeout、stdout/stderr。 | pipeline result、`[REPORT]` parse、warnings。 | pipeline 成功後だけ deploy / snapshot へ進む。 | timeout / exit 非 0 で SHA cache、deploy、snapshot を更新しない。 | pipeline success、non-zero、timeout、duplicate report。 |
| §14a | SSH deploy | deploy target、local output、remote checksum。 | deploy result、pending transfer。 | 転送成功 target だけ success、失敗 target は pending へ保存する。 | checksum mismatch / SSH 失敗で snapshot を作成しない。 | SSH success、checksum mismatch、pending duplicate。 |
| §14b | snapshot | build output、history keep、snapshot keep。 | `.snapshots/{build_id}`、snapshot manifest。 | build / deploy 成功後に atomic save し、世代 prune する。 | snapshot 保存失敗は WARN とし、build success を反転しない。 | snapshot save/prune、snapshot failure remains success。 |
| §15 | logs/history | stdout/stderr、report、warnings、duration、target status。 | `.build_logs/{id}.json`、`.build_history`。 | build log 成功後だけ history を追記する。 | log write failure では history / SHA / deploy / snapshot を行わない。 | build log write failure、history append failure、report parse。 |
| §16〜§18 | systemd / GitHub / setup | unit file、PAT、binary path、timer。 | service/timer 設定、導入済み状態。 | setup 手順で明示された file / unit だけ作成する。 | PAT 不正、checksum 不一致、unit 失敗で後続手順を開始しない。 | setup success、checksum mismatch、service failure。 |
| §19〜§20 | 既知制限反映 | API / runner 制限事項。 | 実装対象外の明示。 | 制限を回避する隠れ機能を追加しない。 | 未定義 endpoint、外部認証、HTTPS listener、worker pool を実装しない。 | 実装 PR 本文で、対象外の節、未定義 endpoint、外部認証、HTTPS listener、worker pool が差分に含まれないことを列挙する。 |

**§8〜§20 中核機能 横断受け入れ固定契約：**

| 項目 | 合格条件 |
|------|----------|
| atomicity | builder output、runner state、build log、history、status、snapshot は、各節で定義した順序以外で確定しない。 |
| no hidden dependency | Go 標準ライブラリと既存 shell / systemd 契約以外の外部依存を追加しない。 |
| no silent success | write failure、history failure、state finalizer failure、checksum mismatch、pipeline timeout を成功扱いにしない。 |
| no secret leak | PAT、SSH secret、token、env secret を stdout、stderr、journal、build log、history、snapshot に平文保存しない。 |
| reproducibility | 同一入力、同一 CLI、同一 fake 外部応答では、時刻・build id を除き同じ状態差分になる。 |
| fixture completeness | §8a と §15a の対象 fixture を未実行または FAIL のまま該当コンポーネント完了扱いにしない。 |
| downstream handoff | runner が生成する `.build_status.json`、`.build_history`、`.build_logs/{id}.json` は §22 の API adapter が追加判断なしに読める schema とする。 |

---

## 9. 既知の制限

| 制限 | 詳細 |
|------|------|
| インデント付き閉じフェンス | CommonMark では `  ``` ` のようなインデント付き閉じフェンスが有効だが、本実装は `raw.strip() == fence_marker` で完全一致を要求するため非対応。インデント付き閉じフェンスはフェンス内の行として取り込まれる |
| 生 HTML のパススルー | Markdown 中の生 HTML（`<div>`、`<span>` 等）は `esc()` でエスケープされテキストとして出力される。HTML をそのまま通過させる機能はない |

---

---

## 27. Builder owner 追加仕様化機能 詳細仕様

### 27.4 出力サイトへのビルドメタ埋め込み

owner component は `builder` とする。collaborator component は `runner`、`api`、`statefile` とする。

`adlaire-ci-build` は `--build-id`、`--commit-sha`、`--build-at` を受け取り、全 HTML ページの `<head>` に次の meta を必ず出力する。

```html
<meta name="adlaire-build-id" content="b20260916010000">
<meta name="adlaire-commit-sha" content="abcdef1">
<meta name="adlaire-build-at" content="2026-09-16T01:00:00Z">
```

値が空文字の場合も meta tag は出力し、`content=""` とする。HTML escape は attribute escape とし、`"`、`&`、`<`、`>` を escape する。runner が pipeline を起動する場合は、同じ値を CLI 引数または環境変数 `ADLAIRE_BUILD_ID`、`ADLAIRE_COMMIT_SHA`、`ADLAIRE_BUILD_AT` のいずれかで渡す。CLI 引数を使える場合は CLI 引数を優先する。

`[REPORT]` には `build_id`、`commit_sha`、`build_at` を追加する。`.build_logs/{id}.json.build_meta`、`GET /api/output-meta` は HTML meta と同じ値を返す。

**build meta 入力固定契約：**

| 項目 | 仕様 |
|------|------|
| `build_id` | 空文字または `^[A-Za-z0-9_-]{1,64}$`。不正値は builder 終了コード `2`。 |
| `commit_sha` | 空文字、7〜40 文字 lowercase hex。その他は終了コード `2`。 |
| `build_at` | 空文字または UTC ISO 8601 秒精度。timezone offset、ミリ秒は終了コード `2`。 |
| HTML / REPORT / API | 3 箇所の値は byte 単位で一致させる。空文字は `""` として保持する。 |
| escape | HTML meta attribute は `esc()` ではなく attribute escape を使う。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 値あり | 全 HTML に 3 meta が同値で出力される。 |
| 値なし | 3 meta が `content=""` で出力される。 |
| 複数ページ | index と全ページで同じ build meta。 |
| output-meta | HTML meta、REPORT、build log、API response が一致する。 |
| 不正 sha | 終了コード `2`、出力差分なし。 |
| 空値 | HTML / REPORT / API がすべて空文字で一致する。 |

### 27.25 ビルドキャッシュ

本機能の目的は、複数ページ静的サイト生成時に未変更入力の変換結果を再利用し、build 時間を短縮することである。

owner component は `builder` とする。collaborator component は `runner`、`statefile` とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| CLI | `adlaire-ci-build --cache-dir <path>` |
| 設定 key | `.server_config.build_cache_enabled` |
| 既定値 | `false` |
| 状態 | `.build_cache.json` と `.build_cache/pages/` |
| cache key | input relative path、input sha256、builder version、theme、build config hash。 |

**正常系：**

1. cache 有効時、builder は入力 file ごとに cache key を計算する。
2. cache hit かつ依存 manifest が一致する場合、HTML fragment と page metadata を再利用する。
3. cache miss の場合、通常変換し、成功後に cache entry を atomic write する。
4. `[REPORT]` に `cache_hits`、`cache_misses`、`cache_disabled_reason` を出力する。

**cache entry 固定契約：**

| 項目 | 仕様 |
|------|------|
| schema version | `.build_cache.json.schema_version=1`。不一致時は全 entry miss。 |
| cache key | `sha256(input_relative_path + "\n" + input_sha256 + "\n" + builder_version + "\n" + theme + "\n" + build_config_hash)` の hex。 |
| page file | `.build_cache/pages/{cache_key}.json`。 |
| page schema | `{ "cache_key", "input_path", "input_sha256", "deps", "html_fragment", "metadata", "created_at" }`。 |
| hit 検証 | cache entry の `cache_key`、`input_path`、`input_sha256`、依存 SHA、builder version、theme、build config hash がすべて一致する場合だけ hit。 |
| 破損 entry | WARN を出し、該当 page file の削除を試みる。削除失敗でも build は継続する。 |
| report | `cache_hits`、`cache_misses` は integer、`cache_disabled_reason` は `null` または固定文字列。 |

**無効化条件：**

`--strict`、theme 変更、builder version 変更、依存 file 変更、cache schema version 不一致、cache entry 破損時は該当 entry を miss とする。

**異常系：**

| 条件 | 処理 |
|------|------|
| cache 読み取り失敗 | WARN、該当 entry miss。 |
| cache 書き込み失敗 | build は成功扱い、WARN と report に記録。 |
| cache 破損 | 該当 entry 削除を試み、miss。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 2 回目同一入力 | cache hit、出力内容一致。 |
| 1 file 変更 | 変更 file のみ miss。 |
| theme 変更 | 全対象 miss。 |
| cache 破損 | build 継続、WARN。 |
| entry 不一致 | miss として通常変換。 |
| cache write failure | build success、report に warning。 |

### 27.28 依存ファイルトラッキング

本機能の目的は、Markdown から参照される画像、相対リンク、include 対象を追跡し、関連する入力だけを再ビルド対象にすることである。

owner component は `builder` とする。collaborator component は `runner`、`statefile` とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.dependency_manifest.json` |
| 対象 | Markdown link、image、HTML `<img src>`、`{{ include "path" }}` |
| path | 相対 path のみ。絶対 URL、fragment-only link は対象外。 |
| schema | `{ "pages": { "<page>": { "deps": [{"path":"...","sha256":"..."}] } } }` |

**正常系：**

1. builder が page ごとに依存 path を抽出する。
2. base dir 基準で正規化し、`..` で base 外へ出る path は broken dependency とする。
3. 依存 file の SHA-256 を記録する。
4. runner は入力 SHA と依存 SHA を比較し、変更された dependency を参照する page を build 対象へ追加する。

**dependency manifest 固定契約：**

| 項目 | 仕様 |
|------|------|
| page key | 入力 Markdown の base dir 相対 path。`/` 区切り、辞書順。 |
| dep path | base dir 相対 path。URL、fragment-only、mailto、tel、data URI は対象外。 |
| include | `{{ include "path" }}` の double quote 形式だけを対象にする。single quote、式展開、glob は対象外。 |
| 重複 dep | page 内で同一 dep path が複数回出ても 1 件だけ保存する。 |
| 保存条件 | build 成功後だけ `.dependency_manifest.json` を置換する。failure build では既存 manifest を維持する。 |
| broken deps | `broken_dependencies[]` に page、path、reason を保存する。strict では終了コード `2`。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| dependency 不在 | builder は WARN、strict なら終了コード `2`。 |
| manifest 破損 | runner は full build。成功時に再作成。 |
| base 外参照 | WARN、strict なら終了コード `2`。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 画像変更 | 参照 page が再ビルド対象。 |
| 未参照画像変更 | build 対象にしない。 |
| manifest 破損 | full build、manifest 再作成。 |
| base 外参照 | broken dependency として記録。 |
| failure build | 既存 manifest を上書きしない。 |
