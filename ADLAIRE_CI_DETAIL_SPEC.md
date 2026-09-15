# Adlaire CI — 詳細仕様

本ファイルは `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様であり、実装の具体的詳細に関する正本である。

方針、ポリシー、実装状態、正本関係は `ADLAIRE_CI_SPEC.md` を正とする。

---

# Part 3 — 仕様
> 実装の具体的詳細を定める。「どのように動作・実装するか」に答える。

---

## 0a. 詳細仕様の記載基準

本ファイルの仕様項目は、実装者が追加の設計判断や推測を行わずに実装できる粒度で記載する。

仕様項目を追加または改訂する場合は、対象範囲に応じて以下を明記する。

| 項目 | 記載する内容 |
|------|-------------|
| 対象 | 対象コンポーネント、対象ファイル、対象機能、責務境界 |
| 入力 | 引数、HTTP リクエスト、設定値、読み込みファイル、環境前提 |
| 出力 | 戻り値、HTTP レスポンス、生成ファイル、ログ、通知、標準出力 |
| 状態 | 状態ファイル、メモリ上の状態、ロック、キャッシュ、更新タイミング |
| データ構造 | JSON キー、型、必須/任意、既定値、許容値、例 |
| 処理順序 | 正常系フロー、分岐条件、ループ条件、疑似コード |
| 異常系 | エラー条件、例外処理、ステータスコード、ログレベル、通知条件 |
| 運用条件 | タイムアウト、再試行、冪等性、排他制御、世代管理、削除条件 |
| セキュリティ | 認証、認可、秘密情報の保存禁止、権限、外部公開可否 |
| 検証 | 構文確認、実行確認、API 確認、生成物確認、整合性確認 |

未確定の内容は、実装可能な詳細仕様として記載してはならない。未確定の場合は、該当項目を「未仕様化」または「将来計画」として扱い、何が未確定かを明示する。

仕様化済み・未実装の内容は、実装ファイルが存在しなくても、本節の基準に従って実装可能な粒度まで具体化する。

---

## 0b. 成熟度棚卸し

本節は、`ADLAIRE_CI_SPEC.md` Part 1 §4b および Part 2 §0a に基づく詳細仕様の一次棚卸しである。

成熟度は、詳細仕様の記載粒度、実装ファイルの存在、現行リポジトリ上で確認できる責務境界を基準に判定する。実装済み判定は、実装ファイルの存在だけでは成立しない。詳細なコード突合、構文確認、実行確認、生成物確認は、各実装作業または別途の仕様突合作業で行う。

| 対象範囲 | 対象コンポーネント | 成熟度 | 判定理由 | 次に必要な作業 |
|----------|-------------------|--------|----------|----------------|
| §0〜§9 | `build_spec.py` | 実装済み | 現行リポジトリに `build_spec.py` が存在し、Markdown から HTML を生成する現行ビルドスクリプトとして扱う。 | 詳細仕様と実装の関数・定数・出力レポートを突合し、差分があれば仕様または実装を改訂する。 |
| §10〜§20 | `runner.py` | 改訂予定 | 現行リポジトリに `runner.py` は存在する。§10a で現行実装範囲と仕様化済み・未実装範囲を再分類済みだが、§11〜§20 本文には拡張仕様が混在している。 | §11〜§20 本文を、現行実装範囲と未実装拡張仕様に分割して再構成する。 |
| §21〜§22 | `api_server.py` | 仕様化済み・未実装 | 管理 API サーバーの責務、設定値、systemd、エンドポイント、レスポンス、エラー形式が定義されているが、現行リポジトリに `api_server.py` は存在しない。 | 実装前に API エンドポイントごとの入出力、状態ファイル、エラー条件の不足を確認する。 |
| §23 | `adlaire-ci-sdk.js` | 仕様化済み・未実装 | SDK のクラス、メソッド、戻り値、HTTP 対応関係が定義されているが、現行リポジトリに `adlaire-ci-sdk.js` は存在しない。 | API 仕様と SDK メソッド一覧を同期確認し、未定義の戻り値型があれば具体化する。 |
| §24 | `admin/index.html` | 仕様化済み・未実装 | 標準管理ツールの画面構成、表示条件、パネル責務が定義されているが、現行リポジトリに `admin/index.html` は存在しない。 | API・SDK と UI 操作の対応を確認し、各操作の成功/失敗表示を具体化する。 |
| §25 | `api_server.py` | 仕様化済み・未実装 | 認証情報ファイル、パスワードハッシュ、ログイン回数、パスワード変更フローが定義されているが、現行リポジトリに認証実装は存在しない。 | セッション管理、トークン生成、ファイル権限、異常系を API 仕様と突合する。 |
| §26 | `runner.py` / `api_server.py` | 改訂予定 | セットアップ・アップデート手順は、未実装の `api_server.py` を含む導入手順であり、現行実装だけでは完了手順として扱えない。 | 現行実装向け手順と、管理 API 導入後の手順を分離する。 |
| 概要内の MCP 記載 | `mcp_server.py` | 将来計画 | `mcp_server.py` は将来構成として言及されるが、詳細な入出力、ツール定義、起動手順、認証仕様は本ファイル内で実装可能な粒度まで定義されていない。 | 実装対象にする場合は、先に改訂予定へ昇格し、MCP 詳細仕様を新設する。 |

成熟度棚卸しの結果、現時点で優先して整合すべき対象は `runner.py` 詳細仕様である。`runner.py` は現行実装ファイルが存在する一方で、詳細仕様側に拡張済みの項目が多いため、§10a の分類に従って実装済み範囲と未実装範囲を区別して扱う。

---

## 0. システム概要

Adlaire CI の現行実装は `build_spec.py` と `runner.py` の 2 つのスクリプトで構成される。

本仕様では、管理 API サーバー `api_server.py`、標準管理ツール `admin/index.html`、JavaScript SDK `adlaire-ci-sdk.js` も仕様化済み・未実装コンポーネントとして定義する。将来的には `mcp_server.py` を加えた構成へ移行予定（→ §13 将来計画 MCP サーバー実装）。

**`build_spec.py`（ビルドスクリプト）**
GitHub リポジトリ上の Markdown 仕様書（`adlaire-db-spec.md`）を HTML に変換してローカルパスへ出力する。入力（`SRC`）と出力（`OUT`）は `build_spec.py` の定数で管理する。

**`runner.py`（CI ランナー）**
GitHub の Git Trees API / Git Blobs API を使用し、単一対象ファイルの blob SHA 変更を検出する。変更があった場合のみ Markdown 本文を `SRC` へ書き出し、`pipeline.sh` を介してビルドを起動し、成功時に SHA キャッシュを更新する。systemd タイマー（5 分間隔）で定期実行する oneshot 設計。

SSH 転送、ペンディングキュー、スナップショット、Webhook 通知、マルチブランチ、ビルドログ保存、サーキットブレーカーは §10a・§12〜§14b に定義する仕様化済み・未実装の拡張機能であり、現行 `runner.py` の実装済み範囲には含めない。

**`api_server.py`（管理 API サーバー、仕様化済み・未実装）**
常駐 HTTP サーバー（`http.server`）。管理ツールからの API リクエストを受け付け、認証・状態取得・手動ビルドトリガーを処理する。`adlaire-admin.service` として systemd に登録し、`runner.py` とは独立して常駐する。

**現行実装の実行フロー：**
```
systemd timer
  └─ runner.py（oneshot）
       ├─ 変更なし → スキップ
       └─ 変更あり → pipeline.sh → build_spec.py → HTML 生成
```

**仕様化済み・未実装コンポーネントおよび拡張機能を含む想定フロー：**
```
runner.py（拡張後）
  └─ SSH 転送 / スナップショット / Webhook 通知 / ビルドログ保存

adlaire-admin.service（常駐）
  └─ api_server.py → SDK → 管理ツール
```

外部ライブラリ・git・nginx 不要。Python 標準ライブラリと GitHub PAT（`contents: read`）のみで動作する。

---

## 1. 要件

| 項目 | 内容 |
|------|------|
| Python バージョン | 3.9 以上（型ヒント `dict[str, int]`、`list[tuple]` を使用） |
| 外部依存 | **なし** — `re`・`html`・`unicodedata` の標準ライブラリ 3 モジュールのみ使用。`pip install` 不要 |
| 入力 | UTF-8 エンコードの Markdown ファイル |
| 出力 | UTF-8 エンコードの単一 HTML ファイル |

---

## 2. ファイルパス設定

スクリプト冒頭の定数で入出力パスを管理する。

```python
SRC = "/opt/adlaire-builder/repo/adlaire-db-spec.md"  # 入力 Markdown
OUT = "/opt/adlaire-builder/dist/Adlaire-db-spec.html"  # 出力 HTML（CI サーバーローカル）
```

別の環境で実行する場合はこの 2 変数を書き換える。

---

## 3. 処理パイプライン

```
┌────────────────────────────────────────────────────────────────┐
│ 1. MDファイル読み込み（raw_lines）                              │
├────────────────────────────────────────────────────────────────┤
│ 2. 見出し抽出パス（フェンス内を除外）                           │
│    → headings: list[(level, text, slug, line_number)]          │
│    → slug_by_line: dict[line_number → slug]                    │
├────────────────────────────────────────────────────────────────┤
│ 2b. 脚注定義収集パス                                            │
│    → _fn_defs: dict[id → text]（脚注 ID とテキストの対応）      │
│    ※ _fn_order は convert() 実行中に inline() が参照順に更新    │
├────────────────────────────────────────────────────────────────┤
│ 3. TOC HTML 生成（build_toc()）                                 │
│    → toc_html: str                                             │
├────────────────────────────────────────────────────────────────┤
│ 4. MD → HTML 変換（convert()）                                  │
│    → body_html: str                                            │
├────────────────────────────────────────────────────────────────┤
│ 5. HTML テンプレート合成（PAGE f-string）                        │
│    CSS トークン・レイアウト・JS をすべてインライン埋め込み        │
├────────────────────────────────────────────────────────────────┤
│ 6. ファイル書き出し（OUT）                                       │
└────────────────────────────────────────────────────────────────┘
```

---

## 4. 関数リファレンス

### 4.1 `slugify(text: str) → str`

見出しテキストから HTML アンカー用のスラグ（`id` 属性値）を生成する。

**処理手順：**
1. Markdown 記法文字（`` ` * _ ~ [ ] ``）を除去
2. 文字を 1 文字ずつ走査し、空白・ハイフン・ドット・アンダースコアは `-` に変換、英数字・Unicode 文字（カテゴリ `L*`、`N*`）はそのまま保持、それ以外は除去
3. 連続する `-` を 1 つに正規化、前後の `-` をトリム
4. 空文字列になった場合は `section` にフォールバック
5. 同一スラグが複数回出現した場合、2 回目以降は `-1`、`-2` と連番サフィックスを付与

**グローバル状態：** `_seen: dict[str, int]` — 同一スラグの出現回数を追跡

**注意：** `_seen` はモジュールレベルのグローバル変数であり、スクリプト実行中に状態が蓄積される。スクリプトを複数回 `import` して使用する場合、`_seen` をリセットする必要がある。

---

### 4.2 `esc(s: str) → str`

`html.escape(s, quote=True)` のラッパー。HTML 特殊文字（`<`、`>`、`&`、`"`、`'`）をエスケープする。

---

### 4.3 `inline(text: str) → str`

インライン Markdown 記法を HTML に変換する。コードスパンを先にキャラクターレベルのスキャンで処理することで、後続の正規表現がコードスパン内の記法を誤って変換するのを防ぐ。

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

`url` が `#` で始まる内部アンカーリンクを処理する際、アンカー部分（`#` 以降）をモジュールレベルのグローバル変数 `_internal_link_refs: set[str]` に追記する。`convert()` 末尾で生成済みスラグの全集合（`_slug_count` のキー）と照合し、一致しないアンカーを次のように報告する：

```
[WARN] BROKEN_LINK: #anchor-text  (in: [label](#anchor-text))
```

不一致件数は `[REPORT]` の `broken_links` フィールドに反映される。照合は変換完了後に行うため、ドキュメント内の順序（前方参照・後方参照）を問わず検証できる。

**グローバル状態（脚注・内部リンク）：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `_fn_defs` | `dict[str, str]` | 脚注 ID → テキストの対応（`convert()` 呼び出し前に収集） |
| `_fn_order` | `list[str]` | 本文中での参照順脚注 ID リスト（`inline()` が参照時に追記） |
| `_internal_link_refs` | `set[str]` | 本文中に出現した内部アンカー参照（`#` 除いた文字列）の集合 |

**制約：** ネストしたインライン記法（`**_text_**` など）は限定的にサポート。

---

### 4.4 `build_toc(headings) → str`

見出しリストから TOC（目次）の HTML を生成する。h1〜h3 のみを対象とし（h4 は除外）、子見出しを持つ見出しはグループとしてアコーディオン形式に構築する。

**引数：** `headings: list[tuple[int, str, str, int]]` — `(level, text, slug, line_number)` のリスト

**出力：** TOC の `<li>` 要素群の HTML 文字列（`<ul>` ルートは HTML テンプレート側で定義）

**グループ（`.tg`）とリーフ（`.ti`）の判定：**
現在の見出しの次の見出しレベルが現在より深い場合、グループとして扱い、展開ボタン（`.tg-btn`）付きの `<ul>` をネストする。それ以外はリーフ（`.ti`）として `<li>` 1 個を出力する。

**スタック管理：** 内部スタック `stack: list[tuple[int, str]]` でグループの開閉を追跡し、`close_to(target_lv)` で不要になった `</ul></li>` を閉じる。スタックにはグループ（`'group'`）だけでなく、レベル管理のためリーフ（`'item'`）も積まれる。`close_to()` はリーフをサイレントに捨て、グループのみ `</ul></li>` を出力する。

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

---

### 4.5 `convert(lines, slug_by_line) → str`

Markdown の行リストを走査し、HTML コンテンツ文字列を生成するメインコンバーター。

**引数：**
- `lines: list[str]` — Markdown の全行
- `slug_by_line: dict[int, str]` — 行番号から見出しスラグへのマッピング

**内部バッファと状態変数：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `out` | `list[str]` | 出力 HTML 断片の蓄積 |
| `fence_active` | `bool` | コードフェンス内かどうか |
| `fence_lang` | `str` | コードフェンスの言語識別子 |
| `fence_buf` | `list[str]` | フェンス内の行バッファ |
| `fence_marker` | `str` | フェンス開始マーカー（`` ``` `` or `~~~`） |
| `para_buf` | `list[str]` | 段落テキストの行バッファ |
| `list_stack` | `list[tuple[str, int]]` | リストのネスト状態スタック |
| `table_buf` | `list[str]` | テーブル行のバッファ |

**内部ヘルパー関数（`convert` のクロージャ）：**

- `flush_para()` — `para_buf` を `<p class="mp">` として出力し、バッファをクリア
- `flush_list()` — `list_stack` を巻き戻し、すべての `</ul>` / `</ol>` を閉じる
- `flush_table()` — `table_buf` をパースし `<div class="tw"><table class="mt">` として出力
- `emit_code()` — `fence_buf` を `<div class="cb-wrap">` 構造として出力

**ブロック要素の検出優先順位（1 行ずつ処理）：**

1. フェンスコードブロック開始・終了（`` ``` `` / `~~~` で始まる行）
2. 見出し（`#` で始まる行）
3. 水平線（`---`、`***`、`___`）
4. テーブル（`|` で始まる行、連続する `|` 行をまとめて処理）
5. 引用（`>` で始まる行、連続する `>` でネスト可能）
6. 定義リスト（`: 定義` 形式の行かつ `para_buf` に用語がある場合）
7. リスト項目（`- * +` または `1.` 形式、`[ ]`/`[x]` プレフィックスでタスクリスト）
8. 空行（バッファのフラッシュトリガー）
9. 脚注定義行（`[^id]:` で始まる行、`_fn_defs` 収集済みのためスキップ）
10. 段落（上記以外の非空行、連続行を 1 つの `<p>` にまとめる）。先読みループは次のいずれかに該当する行で停止する：`#`（見出し）、`|`（テーブル）、`` ` ``×3以上（フェンス）、`~`×3以上（フェンス）、`>`（引用）、リストマーカー（`[-*+]` または `\d+[.)]`）、`: `（定義リストマーカー）、水平線（`---+`・`***+`・`___+`）

**テーブル変換の詳細：**
セパレーター行（`:---:`、`---` などで構成された行）のインデックスを自動検出し、セパレーター行より前の行をヘッダー（`<th>`）、それ以降を本文（`<td>`）として出力する。セパレーター行自体は出力しない。

**引用ネストの詳細：**
`>` で始まる連続行をまとめて収集し、内部関数 `_render_bq(blines)` が再帰的にネストを処理する。1 レベル分の `>` を剥いた後、内側行を先頭から走査し、`>` で始まる連続する行は `_render_bq()` を再帰呼び出し、それ以外の行は `inline()` でレンダリングして結合する。これにより、単一行・複数行・混在ネスト（同一ブロック内で `>` 行と `>>` 行が混在する場合）をすべて正しく処理する。例：`>> text` → `<blockquote class="mbq"><blockquote class="mbq">text</blockquote></blockquote>`。

**タスクリストの詳細：**
リスト項目のコンテンツが `r'^\[([ xX])\]\s+'` に一致する場合、`<li class="ml-task">` として出力する。チェック済み（`[x]` / `[X]`）は `checked` 属性付き、未チェック（`[ ]`）は属性なしの `<input type="checkbox" disabled>` を先頭に配置する。

**定義リストの詳細：**
`: 定義` 行（`startswith(': ')`、コロン＋スペース1文字）を検出したとき `para_buf` に内容があれば、`para_buf` の末尾要素を用語（`<dt>`）として取り出し、`<dl class="mdl"><dt>用語</dt><dd>定義</dd></dl>` を出力する。連続する `: ` 行は同一 `<dl>` 内の追加 `<dd>` としてまとめて処理し、その後に `</dl>` を閉じる。

**脚注定義行のスキップ：**
`r'^\[\^[^\]]+\]:'` に一致する行は `_fn_defs` への収集が完了しているためスキップし、本文への出力を行わない。

**脚注セクションの末尾出力：**
`convert()` 末尾で `_fn_order` が非空の場合、`<section class="fn-section">` 内に参照順番号付きの脚注リスト（`<ol class="fn-list">`）を出力する。各脚注には本文への戻りリンク（`<a class="fn-back">↩</a>`）を付与する。

**フェンスコードブロックの未閉鎖フォールバック：**
ファイル末尾まで読んだ時点で `fence_active` が `True` のままの場合（閉じる `` ``` `` がない場合）、`fence_buf` にコンテンツがあれば `emit_code()` を呼び出して強制出力する。`fence_buf` が空（フェンス開始直後に EOF）の場合は何も出力しない。

**見出し階層スキップ警告：**

`convert()` 内で `_prev_heading_level: int = 0` をローカル変数として保持する。見出し行（`#` で始まる行）を処理するたびに現在レベルと前回レベルを比較し、2 段以上の降順スキップ（例：h1→h3、h2→h4）を検出した場合に次の形式で `[WARN]` を出力する：

```
[WARN] HEADING_SKIP: h1→h3 "見出しテキスト"
```

- 昇順への復帰（例：h3→h1）はスキップに該当しない（章の区切りとして正常）
- 同レベルの連続（h2→h2）・1段降順（h2→h3）もスキップに該当しない
- 件数は `[REPORT]` の `heading_skips` フィールドに反映される

**読了時間集計：**

`convert()` 内で `_char_count: int = 0` をローカル変数として保持する。段落・リスト・引用テキストを `inline()` 処理する直前に、元の Markdown テキスト文字数（スペース・改行を含む）を加算する。コードブロック・フェンス内テキスト・見出しテキスト・テーブルは集計対象外とする。

読了時間の算出：
```python
reading_time_minutes = math.ceil(_char_count / 200)  # 200文字/分、切り上げ
```

算出した `reading_time_minutes` は `convert()` の戻り値と並んで呼び出し元（`runner.py` / `pipeline.sh` 経由）に渡し、`[REPORT]` 行と `.build_logs/{id}.json` に記録する。また HTML ヘッダーへの静的埋め込み（§5）にも使用する。

**見出し出力 HTML 構造：**
`#` で始まる行を `h1`〜`h4` に変換する際、末尾に `.hn-link` ボタンを付与する。

```html
<h2 id="slug" class="mh h2">見出しテキスト<button class="hn-link" data-href="#slug" aria-label="リンクをコピー">¶</button></h2>
```

- `data-href` 属性：`#` + `slugify()` で生成したスラグ
- `¶`（U+00B6 PILCROW SIGN）を使用
- CSS で通常時 `opacity: 0`、親見出し要素のホバー時に `opacity: 1` に変化する

**見出しスラグ重複解決：**

`convert()` 内では `_slug_count: dict[str, int]` をローカル変数として保持し、同一スラグが複数の見出しに割り当てられる場合に一意化する。

| 条件 | スラグ |
|------|--------|
| 初出 | `{slug}` そのまま |
| 2 回目 | `{slug}-2` |
| 3 回目 | `{slug}-3` |
| N 回目 | `{slug}-{N}` |

```python
_slug_count: dict[str, int] = {}

def _unique_slug(base: str) -> str:
    n = _slug_count.get(base, 0) + 1
    _slug_count[base] = n
    return base if n == 1 else f"{base}-{n}"
```

一意化後のスラグは `id` 属性・`data-href` 属性・TOC リンク `href`（§4.2）・`¶` ボタン・全文検索インデックス（§7.9）・前後章ナビゲーション（§7.15）のすべてで共通使用する。

---

### 4.6 `emit_code()` （`convert` 内クロージャ）

フェンスコードブロックを HTML に変換して `out` に追記する。

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
- コード内容：`\n`.join(fence_buf) を `esc()` でエスケープしてから出力

---

## 5. HTML 出力構造

```html
<!DOCTYPE html>
<html lang="ja">
<head>
  <!-- インライン CSS（ADS トークン + レイアウト + コンポーネント） -->
</head>
<body>
  <div id="progress-bar"></div>   <!-- 読み取り進捗バー（ページ上端固定、高さ 3px、幅 = スクロール率 % → §7.13） -->
  <header id="hdr">        <!-- 固定ヘッダー（高さ 52px、背景 --adlaire-surface-accent） -->
    <span id="reading-time">約 87 分</span>  <!-- 読了時間（ビルド時に静的埋め込み → §4.5・§6） -->
  <div id="lay">           <!-- フレックスコンテナ -->
    <nav id="sb">          <!-- サイドバー（幅 260px、固定） -->
      <div class="sb-search-wrap">
        <input id="sb-search">  <!-- TOC 検索 -->
      </div>
      <div id="sb-toc">
        <ul class="tr" id="toc-root">  <!-- TOC リスト -->
      </div>
    </nav>
    <main id="ct">         <!-- コンテンツエリア -->
      <div class="ci">     <!-- 最大幅 760px センタリングコンテナ -->
        <!-- h2 章の区切りごとに以下の <nav class="ch-nav"> が挿入される（→ §7.15） -->
        <!-- 各 h2 章の末尾（次の h2 または文書末）に静的生成 -->
        <nav class="ch-nav">
          <!-- 先頭章は .ch-prev なし、最終章は .ch-next なし -->
          <a class="ch-prev" href="#{prev-slug}">← {prev-title}</a>
          <a class="ch-next" href="#{next-slug}">{next-title} →</a>
        </nav>
        {body_html}
      </div>
    </main>
  </div>
  <button id="btt">        <!-- トップへ戻るボタン -->
  <script>…インライン JS…</script>
</body>
</html>
```

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
出力 HTML はライトモード固定（`prefers-color-scheme` 非対応）。

### 7.2 サイドバー開閉

`localStorage` キー `adb-sb` に `"1"`（開）または `"0"`（閉）を保存する。

- デスクトップ（`> 768px`）：`#sb.closed` クラスと JS インラインスタイル（`ct.style.marginLeft`）で幅を制御。CSS の `#sb.closed ~ #ct { margin-left: 0 }` ルールは初期レンダリング時のみ効く。以降の開閉操作はすべて JS インラインスタイルが CSS クラスより優先する
- モバイル（`≤ 768px`）：`#sb.open` / `transform: translateX` で画面外から引き出す
- モバイルでは TOC リンククリック時に自動的にサイドバーを閉じる

### 7.3 TOC グループ展開

`.tg-btn` クリックで対応する `<ul id="tg-{slug}">` の `hidden` 属性をトグルし、`aria-expanded` 属性を更新する。`data-target` 属性で対象 `<ul>` の ID を指定する。

**開閉状態の永続化：**
展開操作のたびに現在展開中のグループのスラグ配列を `localStorage` キー `adlaire-toc-state` に JSON 文字列で保存する。ページ読み込み時（`DOMContentLoaded`）に同キーを読み込み、保存済みスラグのグループを展開状態で描画する。`localStorage` アクセスはすべて `try/catch` で保護し、失敗時はデフォルト状態（初期展開なし）にフォールバックする。

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

コードブロックに言語別の色分けをインライン JS で適用する。外部ライブラリ不要。

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

ビルド時に検索インデックスを生成し、インライン JSON として HTML に埋め込む。TOC 検索フィルター（§7.4）と検索 UI を統合し、本文ヒット箇所へのジャンプを提供する。

**インデックス生成仕様（build_spec.py）：**
ビルド時に全見出しと各段落の先頭 200 文字を抽出し、以下の配列形式で `<script id="search-index">` タグに埋め込む。

```json
[
  { "id": "anchor-slug", "title": "見出しテキスト", "body": "段落先頭200文字..." },
  ...
]
```

**検索 UI の配置：**
§7.4 の TOC 検索フィルター入力欄を兼用する。入力値が 2 文字以上になった時点でインデックスに対して部分一致検索を実行する。

**ヒット箇所ハイライト：**
一致したエントリの見出しを TOC 内でハイライト（`.toc-hit` クラス付与）。クリックで対象アンカーへスクロールし、`<mark>` 要素でヒット文字列をページ内マーキングする（外部依存なし・標準 DOM 操作のみ）。

**クリア：**
入力欄を空にすると TOC ハイライトおよびページ内マーキングをすべて解除する。

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

**HTML 構造：** `flush_table()` が生成する全 `<th>` に `data-sort="{col_index}"` 属性（0始まりの列インデックス）と `aria-sort="none"` を付与する。

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

---

### 7.15 前後章ナビゲーションボタン

h2 見出し単位で「← 前の章」「次の章 →」ボタンを各章末尾に静的生成する（→ §4.5・§5・§6 CSS）。

**生成方法：** `convert()` の第2パスとして実装する。本文 HTML 生成後、h2 見出しの位置（`id` スラグ・タイトルテキスト）を収集し、各 h2 章の末尾（次の h2 の直前、または文書末）に `<nav class="ch-nav">` を挿入する。

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

**印刷時：** `@media print` で `.ch-nav { display: none }` とする（§6 CSS 参照）。

---

## 8. 実行方法

```bash
python3 build_spec.py
```

**標準出力：**
```
Converting MD...
Building TOC...
Assembling HTML...
Done → /opt/adlaire-builder/dist/Adlaire-db-spec.html  (1,713,731 bytes / 1,673 KB)
[REPORT] headings=342 tables=128 code_blocks=64 warnings=3 size_warn=false broken_links=1 heading_skips=0 reading_time=87
```

**変換レポート行（`[REPORT]` プレフィックス）：**
ビルド完了直後に 1 行で出力する。フィールドはスペース区切りの `key=value` 形式で固定順。

| フィールド | 内容 |
|---|---|
| `headings` | 出力 HTML 内の見出し要素（`h1`〜`h6`）の総数 |
| `tables` | 変換したテーブルの総数 |
| `code_blocks` | 変換したコードブロックの総数 |
| `warnings` | ビルド中に発生した警告件数 |
| `size_warn` | 出力 HTML が `OUTPUT_SIZE_WARN_MB` を超えた場合 `true`、それ以外 `false`（`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`） |
| `broken_links` | 参照先スラグが存在しない内部アンカーリンク（`[label](#anchor)`）の件数 |
| `heading_skips` | 見出しレベルが 2 段以上の降順スキップとなった件数 |
| `reading_time` | 推計読了時間（分、切り上げ）。200文字/分で算出 |

警告が発生した場合、`[REPORT]` 行の直前に `[WARN] {メッセージ}` 形式で 1 件ずつ出力する。

**runner.py による取り込み：**
runner.py は `pipeline.sh` の標準出力から `[REPORT]` 行と `[WARN]` 行を抽出し、パースした結果を `.build_logs/{id}.json` のビルドログエントリに追記する。

```json
{
  "build_id": "20260915-100000",
  "status": "success",
  "report": {
    "headings": 342,
    "tables_count": 128,
    "code_blocks_count": 64,
    "warnings": ["未対応記法: admonition (3箇所)"]
  }
}
```

> **フィールド名の対応：** stdout の `[REPORT]` 行は `tables=` / `code_blocks=` の短縮キーを使用するが、`.build_logs/{id}.json` への保存時および `GET /api/output-meta` レスポンスでは `tables_count` / `code_blocks_count` に変換する（→ §22）。

**再実行時の注意：** `_seen`（スラグ重複カウンタ）・`_fn_order`（脚注参照順）・`_fn_defs`（脚注定義）はいずれもモジュールレベル変数であり、スクリプトを起動するたびに初期化される。通常の `python3 build_spec.py` 実行では複数回実行しても出力は同一になる。ただし本スクリプトを `import` して `convert()` を複数回呼ぶ場合は、呼び出し前に `_fn_order.clear()` および `_seen.clear()` を明示的にリセットする必要がある。

---

## 9. 既知の制限

| 制限 | 詳細 |
|------|------|
| インデント付き閉じフェンス | CommonMark では `  ``` ` のようなインデント付き閉じフェンスが有効だが、本実装は `raw.strip() == fence_marker` で完全一致を要求するため非対応。インデント付き閉じフェンスはフェンス内の行として取り込まれる |
| 生 HTML のパススルー | Markdown 中の生 HTML（`<div>`、`<span>` 等）は `esc()` でエスケープされテキストとして出力される。HTML をそのまま通過させる機能はない |

---

### — CI ランナー —

## 10. CI ランナー 要件

| 項目 | 内容 |
|------|------|
| Python バージョン | 3.9 以上 |
| 外部依存 | **なし** — `urllib.request`・`base64`・`json`・`subprocess`・`os`・`sys`・`logging`（すべて標準ライブラリ） |
| 対象 OS | Linux（systemd 対応環境） |
| ネットワーク | サーバーから `api.github.com` への HTTPS 送信のみ |

---

## 10a. CI ランナー 実装突合・再分類

本節は、現行 `runner.py` と §10〜§20 の詳細仕様を突合した再分類である。

現行 `runner.py` は、GitHub API で単一対象ファイルの blob SHA を確認し、変更がある場合に blob 本文を取得して `pipeline.sh` を実行し、成功後に SHA を更新する最小 CI ランナーである。

### 現行実装済み範囲

| 項目 | 現行実装 | 根拠 |
|------|----------|------|
| 起動形式 | oneshot 実行。`main()` が 1 回の変更確認とビルド実行を行って終了する。 | `main()` |
| 設定値 | `TOKEN_FILE`、`OWNER`、`REPO`、`BRANCH`、`TARGET_FILE`、`SHA_FILE`、`SRC`、`BUILD_SCRIPT`、`LOG_LEVEL`。 | `runner.py` 冒頭定数 |
| GitHub PAT 読み込み | `TOKEN_FILE` を読み込み、不在または空の場合は ERROR ログ後に終了する。 | `read_token()` |
| Git Trees API | `BRANCH` の tree を取得し、`TARGET_FILE` の blob SHA を検索する。 | `get_blob_sha()` |
| SHA 比較 | `SHA_FILE` の前回 SHA と現在 SHA を比較し、一致時はビルドをスキップする。 | `read_last_sha()`、`main()` |
| Git Blobs API | blob 本文を取得し、Base64 をデコードする。 | `fetch_blob()` |
| ソース書き出し | 取得した Markdown を `SRC` へ書き出す。 | `write_src()` |
| ビルド起動 | `SRC` と同じディレクトリ配下の `.ci/pipeline.sh` を `bash` で実行する。 | `run_pipeline()` |
| 成功時 SHA 更新 | `pipeline.sh` が exit 0 の場合のみ `SHA_FILE` を現在 SHA で更新する。 | `save_sha()`、`main()` |
| ログ出力 | Python `logging` を stdout へ出力する。 | `setup_logging()` |

### 仕様化済み・未実装範囲

以下は §10〜§20 に詳細仕様が存在するが、現行 `runner.py` には未実装である。実装する場合は、該当仕様を再確認し、必要に応じて詳細仕様を補強してから実装する。

| 項目 | 関連節 | 未実装内容 |
|------|--------|------------|
| JSON 形式の SHA キャッシュ | §11〜§13 | 現行 `SHA_FILE` はプレーンテキスト SHA として読み書きされる。仕様上の `{"sha": "..."}` JSON 形式は未実装。 |
| `BRANCH_TARGETS` | §12〜§13 | 複数ブランチ、複数 target file、複数出力先の設定構造は未実装。現行は `BRANCH` / `TARGET_FILE` / `SHA_FILE` / `SRC` の単一設定。 |
| GitHub API リトライ | §12〜§13 | `API_RETRY_MAX`、`API_RETRY_BASE_SECONDS`、指数バックオフ、レート制限待機は未実装。 |
| ビルドロック | §11〜§13 | `.build_lock` による多重起動防止は未実装。 |
| ビルドクールダウン | §12〜§13 | `BUILD_COOLDOWN_SECONDS` による起動抑制は未実装。 |
| 強制再ビルド間隔 | §12〜§13 | `FORCE_BUILD_INTERVAL` による変更なし時の定期強制ビルドは未実装。 |
| コミット情報記録 | §13 | ビルドトリガー commit の SHA、message、author、date 取得は未実装。 |
| 事前チェック | §13 | ディスク空き容量、Python バージョン、`build_spec.py` 存在確認は未実装。 |
| Webhook 通知 | §13 | `.notify_config` 読み込み、成功/失敗/転送失敗/週次サマリー通知、`.notify_pending` 再送は未実装。 |
| ビルドログ保存 | §11〜§15 | `.build_logs/{id}.json` への stdout/stderr、変換レポート、所要時間保存は未実装。 |
| SSH 転送 | §14a | SHA256 差分検出、stdin パイプ転送、転送後整合性検証、ペンディングキューは未実装。 |
| スナップショット | §14b | `.snapshots/` への成果物保存、世代管理、ロールバック連携は未実装。 |
| サーキットブレーカー | §13 | `.build_circuit_state` による連続失敗停止は未実装。 |
| ログ世代管理 | §13 | `LOG_KEEP_N` による `.build_logs/` 削除は未実装。 |
| 出力サイズ警告 | §12〜§13 | `OUTPUT_SIZE_WARN_MB` による WARN ログと `size_warn` 記録は未実装。 |

### 将来計画として扱う範囲

§10〜§20 には、現行 runner の直接責務ではなく管理 API、標準管理ツール、将来の運用機能と結合して成立する項目が含まれる。これらは実装対象へ進める前に、API・SDK・UI との責務境界を再確認する。

| 項目 | 理由 |
|------|------|
| API 経由の動的ブランチ設定 | `runner.py` 単体では設定 API を持たないため、`api_server.py` 実装と合わせて扱う。 |
| API 経由のロールバック | `POST /api/history/{id}/rollback` は `api_server.py` のエンドポイント実装が前提となる。 |
| 管理画面からのスケジュール操作 | systemd timer の変更 API と標準管理ツール UI が前提となる。 |

---

## 11. CI ランナー ファイル構成

本節のファイル構成は、現行実装で使用するファイル、CI ランナー拡張で追加されるファイル、管理 API / SDK / UI 側のファイルを分離して示す。

現行リポジトリに存在する実装ファイルは `build_spec.py` と `runner.py` のみである。`api_server.py`、`admin/index.html`、`adlaire-ci-sdk.js` は仕様化済み・未実装であり、現行実装済みファイルとして扱ってはならない。

### 現行実装で使用するファイル

| パス | 成熟度 | 用途 |
|------|--------|------|
| `/opt/adlaire-builder/runner.py` | 実装済み | CI ランナー本体。単一ブランチの SHA 検出、blob 取得、`pipeline.sh` 起動、SHA 更新を行う。 |
| `/opt/adlaire-builder/build_spec.py` | 実装済み | Markdown から HTML を生成するビルドスクリプト。 |
| `/opt/adlaire-builder/.github_token` | 実装済み | GitHub PAT。現行 `runner.py` が読み込む。 |
| `/opt/adlaire-builder/.last_sha` | 実装済み | 前回取得した blob SHA。現行実装ではプレーンテキストで保存する。 |
| `/opt/adlaire-builder/repo/adlaire-db-spec.md` | 実装済み | GitHub Blobs API から取得した Markdown の書き出し先。 |
| `/opt/adlaire-builder/repo/.ci/pipeline.sh` | 実装済み | `runner.py` が `bash` で起動するビルド手順。 |

```
/opt/adlaire-builder/
├── runner.py
├── build_spec.py
├── .github_token
├── .last_sha
└── repo/
    ├── adlaire-db-spec.md
    └── .ci/
        └── pipeline.sh
```

### 仕様化済み・未実装の CI ランナー拡張ファイル

以下は §10a の「仕様化済み・未実装範囲」に対応するファイルである。現行 `runner.py` では作成・読み書きしない。

| パス | 成熟度 | 用途 |
|------|--------|------|
| `/opt/adlaire-builder/.build_history` | 仕様化済み・未実装 | ビルド履歴。 |
| `/opt/adlaire-builder/.notify_config` | 仕様化済み・未実装 | Webhook 通知設定。 |
| `/opt/adlaire-builder/.notify_log` | 仕様化済み・未実装 | Webhook 送信履歴。 |
| `/opt/adlaire-builder/.notify_pending` | 仕様化済み・未実装 | Webhook 通知失敗時の再送キュー。 |
| `/opt/adlaire-builder/.pending_transfers` | 仕様化済み・未実装 | SSH 転送失敗時の再送キュー。 |
| `/opt/adlaire-builder/.build_lock` | 仕様化済み・未実装 | 実行中ビルドの PID ロック。 |
| `/opt/adlaire-builder/.branch_config` | 仕様化済み・未実装 | ブランチターゲット設定。 |
| `/opt/adlaire-builder/.build_state` | 仕様化済み・未実装 | ビルド実行状態、週次サマリー送信日等。 |
| `/opt/adlaire-builder/.build_circuit_state` | 仕様化済み・未実装 | サーキットブレーカー状態。 |
| `/opt/adlaire-builder/.build_logs/` | 仕様化済み・未実装 | ビルドごとの個別ログ。 |
| `/opt/adlaire-builder/.snapshots/` | 仕様化済み・未実装 | ビルド成果物スナップショット。 |

### 管理 API / SDK / UI 側ファイル

以下は `api_server.py`、`adlaire-ci-sdk.js`、`admin/index.html` の仕様に属する。CI ランナー拡張と連携するものを含むが、現行 `runner.py` 単体の実装済み範囲には含めない。

```
/opt/adlaire-builder/
├── api_server.py        # 管理 API サーバー（常駐、仕様化済み・未実装）
├── .admin_credentials   # 認証情報ファイル（JSON、パーミッション 600）
├── .server_config       # サーバー設定（JSON）
├── .access_log          # ログイン履歴（JSON）
├── .repo_config         # リポジトリ監視設定（JSON）
├── .config_log          # 設定変更履歴（JSON）
├── .access_control      # IP アクセス制限設定（JSON）
├── .hooks               # ビルド前後フック設定（JSON）
├── .maintenance         # メンテナンスモード状態（JSON）
├── .alert_rules         # カスタムアラートルール（JSON）
├── .tag_rules           # 自動タグ付けルール（JSON）
├── .pipeline_config     # パイプライン設定（JSON）
├── .notes               # 運用ノート（テキスト）
├── .smtp_config         # SMTP 設定（JSON、パスワード除く）
├── .smtp_secret         # SMTP パスワード（プレーンテキスト、パーミッション 600）
├── .dashboard_layout    # ダッシュボードウィジェットレイアウト（JSON）
├── .webhook_events.json # Webhook 受信イベントログ（JSON Lines 形式、1行1イベント）
└── admin/
    ├── index.html           # 管理画面（単一ファイル完結、仕様化済み・未実装）
    └── adlaire-ci-sdk.js    # JavaScript SDK（仕様化済み・未実装）
```

### 出力先・配信先

```
/opt/adlaire-builder/dist/
└── Adlaire-db-spec.html # HTML 出力先

/var/www/html/           # 仕様化済み・未実装の SSH 転送先
```

### systemd

```
/etc/systemd/system/
├── adlaire-ci.service      # systemd ユニット（oneshot）
├── adlaire-ci.timer        # systemd タイマー（定期実行）
└── adlaire-admin.service   # 管理 API サーバー（常駐、仕様化済み・未実装）
```

### リポジトリ側

```
<repo>/
├── adlaire-db-spec.md   # ソース Markdown（GitHub 上のマスター）
└── .ci/
    └── pipeline.sh      # ビルド手順定義（実行権限付き）
```

---

## 12. 設定値（`runner.py` 冒頭）

現行 `runner.py` に実装済みの設定値は、`TOKEN_FILE`、`OWNER`、`REPO`、`BRANCH`、`TARGET_FILE`、`SHA_FILE`、`SRC`、`BUILD_SCRIPT`、`LOG_LEVEL` のみである。

### 現行実装済み設定

```python
TOKEN_FILE   = "/opt/adlaire-builder/.github_token"   # GitHub PAT
OWNER        = "<GitHubオーナー名>"                    # リポジトリオーナー
REPO         = "<リポジトリ名>"                        # リポジトリ名
BRANCH       = "main"                                  # 対象ブランチ
TARGET_FILE  = "adlaire-db-spec.md"                   # 監視対象ファイル
SHA_FILE     = "/opt/adlaire-builder/.last_sha"       # blob SHA キャッシュ
SRC          = "/opt/adlaire-builder/repo/adlaire-db-spec.md"  # 書き出し先
BUILD_SCRIPT = "/opt/adlaire-builder/build_spec.py"
LOG_LEVEL    = "INFO"
```

### 仕様化済み・未実装の拡張設定

以下の `BRANCH_TARGETS`、`PENDING_FILE`、`API_RETRY_MAX`、`BUILD_COOLDOWN_SECONDS`、`HISTORY_KEEP_N`、`FORCE_BUILD_INTERVAL`、`LOG_KEEP_N`、`API_CIRCUIT_BREAKER_THRESHOLD`、`OUTPUT_SIZE_WARN_MB`、`WEEKLY_SUMMARY_*` は仕様化済み・未実装の拡張設定である。

```python
PENDING_FILE           = "/opt/adlaire-builder/.pending_transfers"   # SSH 転送ペンディングキュー（JSON）
API_RETRY_MAX          = 5    # GitHub API 失敗時の最大再試行回数（指数バックオフ）
API_RETRY_BASE_SECONDS = 1    # バックオフ基底秒数（1→2→4→8→16 秒。0 = リトライ無効）
BUILD_COOLDOWN_SECONDS = 60   # 前回ビルド完了から次ビルドまでの最小間隔（秒。0 = 無効）→ §13
HISTORY_KEEP_N         = 10   # スナップショット保持世代数（0 = 無制限）→ §14b
FORCE_BUILD_INTERVAL   = 0    # 強制再ビルド間隔（時間。0 = 無効）→ §13
LOG_KEEP_N             = 50   # ビルドログ保持件数（0 = 無制限）→ §13
API_CIRCUIT_BREAKER_THRESHOLD = 3    # 全ブランチ連続失敗の許容周回数（0 = 無効）→ §13
OUTPUT_SIZE_WARN_MB           = 5    # 出力 HTML サイズ警告閾値（MB。0 = 無効）→ §13・§8
WEEKLY_SUMMARY_ENABLED        = True # 週次サマリー Webhook の有効/無効 → §13
WEEKLY_SUMMARY_DAY            = 0    # 送信曜日（0=月曜〜6=日曜） → §13
WEEKLY_SUMMARY_HOUR           = 9    # 送信時刻（0〜23、ローカル時刻） → §13

# ブランチターゲット設定（→ §14a）
# 複数エントリを定義した場合はリスト順に順次処理する（並列処理は対象外）
BRANCH_TARGETS = [
    {
        "branch":      "main",                                        # 監視対象ブランチ
        "target_file": "adlaire-db-spec.md",                         # 監視対象ファイル
        "sha_file":    "/opt/adlaire-builder/.last_sha",             # blob SHA キャッシュ（JSON 形式: {"sha": "..."}）
        "src":         "/opt/adlaire-builder/repo/adlaire-db-spec.md",  # Blobs API 書き出し先
        "out":         "/opt/adlaire-builder/dist/Adlaire-db-spec.html", # ビルド成果物パス
        "deploy_targets": [
            {
                "host":     "<配信サーバーIP>",                       # SSH 転送先ホスト
                "user":     "deploy",                                 # SSH 接続ユーザー
                "dest_dir": "/var/www/html",                          # 転送先ディレクトリ
            }
        ],
    }
]
```

**拡張設定への移行対応（仕様化済み・未実装）：**

現行 `runner.py` の単一ターゲット設定を `BRANCH_TARGETS` へ拡張する場合の対応は以下とする。現行実装では、`BRANCH`、`TARGET_FILE`、`SHA_FILE`、`SRC` が引き続き使用される。

| 現行項目 | 拡張後の移行先 |
|--------|--------|
| `BRANCH` | `BRANCH_TARGETS[n]["branch"]` |
| `TARGET_FILE` | `BRANCH_TARGETS[n]["target_file"]` |
| `SHA_FILE` | `BRANCH_TARGETS[n]["sha_file"]` |
| `SRC` | `BRANCH_TARGETS[n]["src"]` |
| 出力 HTML パス | `BRANCH_TARGETS[n]["out"]`。現行 `runner.py` では未管理であり、`pipeline.sh` / `build_spec.py` 側の責務。 |
| SSH 転送先 | `BRANCH_TARGETS[n]["deploy_targets"][m]`。現行 `runner.py` では未実装。 |

---

## 13. 処理フロー

本節の処理フローは、現行実装済みの最小フローと、仕様化済み・未実装の拡張フローに分けて扱う。

### 現行実装済みフロー

現行 `runner.py` の実装済みフローは以下である。

```
runner.py 起動
    │
    ├─ .github_token 読み込み
    ├─ Git Trees API で TARGET_FILE の blob SHA を取得
    ├─ SHA_FILE の前回 SHA と比較
    │   ├─ 一致 → INFO ログを出して正常終了
    │   └─ 不一致 → 続行
    ├─ Git Blobs API で blob 本文を取得
    ├─ Base64 デコード後、SRC へ書き出し
    ├─ bash {dirname(SRC)}/.ci/pipeline.sh を実行
    │   ├─ exit 0 → SHA_FILE を現在 SHA で更新
    │   └─ exit 0 以外 → SHA_FILE を更新せず終了コード 1 で終了
    └─ 正常終了
```

### 仕様化済み・未実装の拡張フロー

以下は、仕様化済み・未実装の拡張フローである。現行 `runner.py` はこのフローを実装していない。

```
runner.py 起動（systemd タイマーから呼び出し）
    │
    ├─ .github_token 読み込み（不在の場合は起動失敗）
    │
    ├─ [重複チェック] .build_lock が存在する場合
    │   ├─ ファイル内 PID が実行中 → INFO ログ（`BUILD_SKIP: already running (PID N)`）、正常終了
    │   └─ PID が存在しない（前回の異常終了） → .build_lock を削除して続行
    │
    ├─ .build_lock に自プロセスの PID を書き込み
    │   （以降、正常終了・例外終了いずれの場合も finally で .build_lock を削除）
    │
    ├─ [ペンディングキュー再試行] PENDING_FILE が存在する場合（→ §14a）
    │   └─ ペンディングエントリごとに SSH 転送を再試行
    │       ├─ 成功 → エントリを PENDING_FILE から削除
    │       └─ 失敗 → ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [Webhook 通知ペンディング再試行] .notify_pending が存在する場合
    │   └─ ペンディングエントリごとに Webhook 送信を再試行
    │       ├─ 成功（HTTP 2xx）→ INFO ログ、エントリを .notify_pending から削除
    │       └─ 失敗（HTTP エラー・接続失敗）→ ERROR ログ、エントリを保持（次回起動時に再試行）
    │
    ├─ [クールダウンチェック] BUILD_COOLDOWN_SECONDS > 0 の場合
    │   └─ 前回ビルド完了（.build_history の最終 finished_at）から BUILD_COOLDOWN_SECONDS 秒未満
    │       → INFO ログ（`COOLDOWN: skip, last_build N秒前`）、正常終了
    │
    ├─ BRANCH_TARGETS の各エントリを順次処理：
    │   │
    │   ├─ Step 1: Git Trees API
    │   │   GET /repos/{OWNER}/{REPO}/git/trees/{branch}?recursive=1
    │   │   → target_file の blob SHA を取得
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフ（API_RETRY_BASE_SECONDS × 2^n 秒）で再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリをスキップ
    │   │        └─ レスポンスヘッダー X-RateLimit-Remaining = 0 の場合
    │   │             → X-RateLimit-Reset（Unix 時刻）まで待機してから再試行
    │   │               INFO ログ（`RATE_LIMIT: waiting until {reset_time}`）
    │   │   [PAT 有効期限チェック] レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合
    │   │        ├─ 残日数 ≤ 7 日 → WARN ログ（`PAT_EXPIRY_WARN: expires_at={date} remaining_days={N}`）
    │   │        └─ 残日数 > 7 日 → 処理継続（チェックのみ）
    │   │
    │   ├─ SHA 比較（sha_file の前回 SHA と比較）
    │   │   ├─ 一致（変更なし）かつ FORCE_BUILD_INTERVAL = 0 → INFO ログ、このエントリをスキップ
    │   │   ├─ 一致（変更なし）かつ FORCE_BUILD_INTERVAL > 0 → 前回ビルドから指定時間以上経過していれば強制ビルド続行
    │   │   └─ 不一致（変更あり）→ 続行
    │   │
    │   ├─ Step 2: Git Blobs API
    │   │   GET /repos/{OWNER}/{REPO}/git/blobs/{sha}
    │   │   → Base64 デコード → src パスへ書き出し
    │   │   └─ API 失敗時：API_RETRY_MAX 回まで指数バックオフで再試行
    │   │        ├─ 全試行失敗時：ERROR ログ、このエントリをスキップ
    │   │        └─ レスポンスヘッダー GitHub-Authentication-Token-Expiration が存在する場合は Step 1 と同様に PAT 有効期限チェックを行う
    │   │
    │   ├─ [コミット情報取得] ビルドトリガーとなったコミット情報を取得し、ビルドログへ記録する
    │   │   GET /repos/{OWNER}/{REPO}/commits?path={BRANCH_TARGET.src}&sha={branch}&per_page=1
    │   │   → 先頭エントリから取得：
    │   │       commit_sha    = commit.sha
    │   │       commit_message = commit.commit.message（先頭1行のみ）
    │   │       commit_author  = commit.commit.author.name
    │   │       commit_at      = commit.commit.author.date
    │   │   └─ API 失敗時：各フィールドを null として記録し、処理続行（ビルドは妨げない）
    │   │
    │   ├─ [事前チェック] pipeline.sh 実行前に以下を確認し、不足時は ERROR ログ＋deploy_failure Webhook 通知、このエントリをスキップ
    │   │   ├─ ディスク空き容量 ≥ 出力ファイル推定サイズ × 3（`shutil.disk_usage`）
    │   │   ├─ Python バージョン ≥ 3.9（`sys.version_info`）
    │   │   └─ `build_spec.py` が存在すること（`os.path.exists`）
    │   │
    │   ├─ pipeline.sh 実行（bash {src の親ディレクトリ}/.ci/pipeline.sh）
    │   │   ├─ 成功（exit 0）：INFO ログ
    │   │   │   └─ [Webhook 通知送信] on: ["success"] 設定時
    │   │   │       → .notify_config の Webhook 宛先へ POST（ペイロード: event="success", branch, build_id 等）
    │   │   │       → 送信失敗（HTTP エラー・接続失敗・タイムアウト）の場合：
    │   │   │           ERROR ログ（`NOTIFY_FAIL: url={url} status={code}`）
    │   │   │           .notify_pending へ `{"event":"success","url":"...","payload":{...},"queued_at":"<ISO8601>"}` を追記
    │   │   └─ 失敗（exit ≠ 0）：ERROR ログ、sha_file 更新せず、このエントリをスキップ
    │   │       └─ [Webhook 通知送信] on: ["failure"] 設定時
    │   │           → .notify_config の Webhook 宛先へ POST（ペイロード: event="failure", branch, build_id 等）
    │   │           → 送信失敗の場合：ERROR ログ、.notify_pending へキューイング（success と同一形式）
    │   │
    │   └─ sha_file を新 SHA で更新（`{"sha": "<new_sha>"}` を JSON 書き込み）
    │        └─ SSH ファイル転送（deploy_targets リストの各エントリへ転送 → §14a）
    │             ├─ [転送後整合性検証] ssh user@host "sha256sum /dest/file" でリモート SHA を取得
    │             │   ├─ ローカル sha256 と一致 → 転送成功
    │             │   └─ 不一致またはコマンド失敗 → ERROR ログ、ペンディングキューへ再投入（§14a）
    │             └─ 整合性検証成功後 → スナップショット保存（→ §14b）
    │
    │        [出力サイズチェック] OUTPUT_SIZE_WARN_MB > 0 の場合
    │        出力 HTML ファイルのサイズを取得し、閾値と比較：
    │            size_mb = os.path.getsize(output_path) / (1024 * 1024)
    │            size_mb > OUTPUT_SIZE_WARN_MB の場合：
    │            → WARN ログ（`OUTPUT_SIZE_WARN: size={size_mb:.1f}MB threshold={OUTPUT_SIZE_WARN_MB}MB`）
    │            → ビルドログの size_warn フィールドを true に設定（§8）
    │
    └─ [サーキットブレーカー判定] API_CIRCUIT_BREAKER_THRESHOLD > 0 の場合
        全ブランチの今周回結果を集計し、全ブランチが失敗（API エラー・スキップを除くビルド失敗）の場合：
        → .build_circuit_state（JSON）の consecutive_failures を +1
        いずれか成功した場合：
        → consecutive_failures を 0 にリセット
        consecutive_failures ≥ API_CIRCUIT_BREAKER_THRESHOLD の場合：
        → ERROR ログ（`CIRCUIT_OPEN: consecutive_failures={N} threshold={API_CIRCUIT_BREAKER_THRESHOLD}`）
        → deploy_failure Webhook 通知（reason="circuit_open"）
        → .build_circuit_state に `{"open": true, "consecutive_failures": N, "opened_at": "<ISO8601>"}` を書き込み
        → 以降のポーリング周回はビルドをスキップ（SHA チェックも行わない）
        → POST /api/circuit-breaker/reset でリセット可能（open: false、consecutive_failures: 0 に戻す）

    ├─ [ログクリーンアップ] LOG_KEEP_N > 0 の場合
    │   .build_logs/ 内の {id}.json を mtime 昇順でソートし、
    │   件数が LOG_KEEP_N を超えた分を古いものから削除
    │
    └─ [週次サマリー判定] WEEKLY_SUMMARY_ENABLED = True かつ on: ["weekly_summary"] 設定の Webhook 宛先が存在する場合
        現在の曜日が WEEKLY_SUMMARY_DAY かつ現在時刻が WEEKLY_SUMMARY_HOUR:00±30分以内の場合：
        ├─ 二重送信防止チェック：.build_state の weekly_summary_sent_date と当日の日付（YYYY-MM-DD）を比較
        │   → 一致（本日送信済み）→ スキップ（INFO ログ: `WEEKLY_SUMMARY_SKIP: already sent today`）
        │   → 不一致（未送信）→ 以下を実行
        ├─ 過去 7 日間の .build_logs/{id}.json を集計：
        │       成功件数・失敗件数・成功率（%）・平均ビルド時間（秒）・最長ビルド（秒・ID）
        ├─ on: ["weekly_summary"] 設定の Webhook 宛先へ POST
        │       ペイロード: event="weekly_summary", period="7d", success_count, failure_count, success_rate, avg_duration_seconds, max_duration_seconds, max_duration_build_id
        ├─ 送信成功 → .build_state に `weekly_summary_sent_date: "YYYY-MM-DD"` を記録（INFO ログ）
        └─ 送信失敗（HTTP エラー・接続失敗）→ ERROR ログ、.notify_pending へキューイング（他の Webhook 失敗と同一形式）
```

---

## 14. `pipeline.sh`

リポジトリの `.ci/pipeline.sh` にビルド手順を記述する。

- 実行権限（`chmod +x`）が必要
- `set -e` を先頭に記述し、ステップ失敗時に即座に終了させることを推奨
- 終了コード `0` で成功、`0` 以外で失敗とみなす

**例：**
```bash
#!/bin/bash
set -e
python3 /opt/adlaire-builder/build_spec.py
```

`build_spec.py` のパスは `pipeline.sh` 内に直接記述する（`runner.py` は参照しない）。`build_spec.py` はサーバー固定（`/opt/adlaire-builder/`）のため、リポジトリには含めない。

---

## 14a. SSH ファイル転送

本節は、仕様化済み・未実装の CI ランナー拡張仕様である。

現行 `runner.py` は SSH 転送を実行しない。現行実装は `pipeline.sh` 起動と成功時 SHA 更新までを担当し、生成 HTML の静的コンテンツ配信サーバーへの転送は未実装である。

本機能を実装する場合、`runner.py` は `pipeline.sh` 成功後に、出力ファイルを SSH 経由で静的コンテンツ配信サーバーへ転送する。scp・rsync は使用しない。

### 設定値

`BRANCH_TARGETS` 各エントリの `deploy_targets` リスト内で管理する（→ §12）。

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `host` | 配信サーバーのホスト名 / IP | `"192.0.2.1"` |
| `user` | SSH 接続ユーザー | `"deploy"` |
| `dest_dir` | 配信サーバー上の転送先ディレクトリ | `"/var/www/html"` |

転送対象ファイルは `BRANCH_TARGETS[n]["out"]` から自動導出する。`deploy_targets` に複数エントリを定義した場合は全ての転送先へ順次転送する。

### 差分検出

転送前にリモートサーバーで対象ファイルの SHA256 ハッシュを取得し、ローカルファイルのハッシュと比較する。

```bash
# runner.py が subprocess 経由で実行
ssh <user>@<host> "sha256sum <dest_dir>/<filename>"
```

- ハッシュが一致 → スキップ（`SKIP` ログを記録）
- ハッシュが不一致、またはリモートにファイルが存在しない → 転送実行

### 転送

stdin パイプ経由で SSH 転送する。

```bash
# runner.py が subprocess（stdin=PIPE）経由で実行
ssh <user>@<host> "cat > <dest_dir>/<filename>" < <localfile>
```

### ペンディングキュー

転送失敗時（接続エラー・認証失敗等）は `PENDING_FILE`（JSON）へエントリを追記する。

```json
[
  {
    "branch_idx": 0,
    "deploy_idx": 0,
    "out": "/opt/adlaire-builder/dist/Adlaire-db-spec.html",
    "host": "192.0.2.1",
    "user": "deploy",
    "dest_dir": "/var/www/html",
    "failed_at": "2026-09-15T10:00:00",
    "retry_count": 1
  }
]
```

- `runner.py` 起動時（`BRANCH_TARGETS` 処理前）に `PENDING_FILE` を読み込み、エントリごとに再試行する（→ §13 処理フロー）
- 再試行成功時にエントリを削除する。失敗時は `retry_count` をインクリメントして保持する
- SSH 転送失敗 Webhook 通知（`deploy_failure` イベント）を送信する（on: `["deploy_failure"]` 設定時）

### 転送後整合性検証

SSH 転送完了後に、リモートファイルの SHA-256 チェックサムをローカルのものと照合する。

**検証コマンド：**
```
ssh {user}@{host} "sha256sum {dest_dir}/{filename}"
```
出力形式 `{hash}  {filename}` の最初のフィールドをローカル `hashlib.sha256` の hex digest と比較する。

| 項目 | 仕様 |
|---|---|
| タイムアウト | 30 秒（SSH 転送タイムアウトとは独立） |
| 検証失敗時 | ERROR ログ＋ペンディングキューへ再投入。スナップショット保存はしない |
| ログフィールド | `transfer_verified: false`（`.build_logs/{id}.json` に記録） |
| 正常時 | `transfer_verified: true`（`.build_logs/{id}.json` に記録） |

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スキップ（差分なし） | `INFO` | `SKIP Adlaire-db-spec.html: no change` |
| 転送成功 | `INFO` | `DEPLOY Adlaire-db-spec.html → 192.0.2.1` |
| 転送失敗→キューイング | `ERROR` | `DEPLOY FAILED Adlaire-db-spec.html: <reason> (queued)` |
| ペンディング再試行成功 | `INFO` | `PENDING RETRY OK Adlaire-db-spec.html → 192.0.2.1` |
| ペンディング再試行失敗 | `ERROR` | `PENDING RETRY FAILED Adlaire-db-spec.html: <reason>` |
| 整合性検証失敗→再投入 | `ERROR` | `VERIFY FAILED Adlaire-db-spec.html → 192.0.2.1: checksum mismatch (queued)` |

---

## 14b. スナップショット管理

本節は、仕様化済み・未実装の CI ランナー拡張仕様である。

現行 `runner.py` は `.snapshots/` ディレクトリを作成・更新しない。スナップショット保存、世代管理、ロールバックは未実装である。

本機能を実装する場合、`runner.py` は SSH 転送成功後に、ビルド成果物を `.snapshots/` ディレクトリへアーカイブする。`HISTORY_KEEP_N = 0` の場合はスナップショット機能を無効化する。

### ディレクトリ構造

```
/opt/adlaire-builder/
└── .snapshots/
    ├── b20260915100000/           # build_id = b{YYYYMMDDHHmmss}
    │   └── Adlaire-db-spec.html  # ビルド成果物のコピー
    ├── b20260914180000/
    │   └── Adlaire-db-spec.html
    └── ...
```

- ディレクトリ名は `b{YYYYMMDDHHmmss}` 形式のビルド ID（`.build_history` の `id` と一致する）
- `BRANCH_TARGETS` に複数エントリがある場合は、同一ビルド ID ディレクトリ内に各エントリの成果物をまとめて保存する

### 世代管理

- スナップショット保存後、`.snapshots/` 内のディレクトリ数が `HISTORY_KEEP_N` を超えた場合、最古のディレクトリから順に削除する
- 削除対象ディレクトリの特定は作成日時降順ソートで行う（ディレクトリ名の辞書順 = 時系列順）

### ロールバック

`POST /api/history/{id}/rollback`（→ §22）で指定ビルド ID のスナップショットから SSH 転送を再実行する。

- ロールバック API は `api_server.py` の実装を前提とする。現行リポジトリに `api_server.py` は存在しないため、現行実装済み機能として扱ってはならない
- `.snapshots/{id}/` が存在しない場合は `404` を返す
- 転送成功時は `.build_history` に rollback エントリを追記する

### ログ

| 状態 | ログレベル | メッセージ例 |
|------|-----------|------------|
| スナップショット保存 | `INFO` | `SNAPSHOT b20260915100000 saved` |
| 古世代削除 | `INFO` | `SNAPSHOT b20260910000000 pruned (keep_n=10)` |
| ロールバック成功 | `INFO` | `ROLLBACK b20260914180000 → 192.0.2.1 OK` |
| ロールバック失敗 | `ERROR` | `ROLLBACK b20260914180000: <reason>` |

---

## 15. ログ

本節は、現行実装済みの stdout ログと、仕様化済み・未実装の構造化ログ拡張を分けて定義する。

### 現行実装済みログ

Python 標準の `logging` モジュールを使用する。出力先は stdout（systemd が journald に転送）。

| レベル | 出力条件 |
|--------|---------|
| `INFO` | 起動、変更なしスキップ、ビルド開始・完了、SHA 更新 |
| `WARNING` | — |
| `ERROR` | トークン読み込み失敗、API 失敗、ビルド失敗 |
| `DEBUG` | API レスポンス詳細等（`LOG_LEVEL = "DEBUG"` 時のみ） |

現行 `runner.py` は `.build_logs/{id}.json` を作成しない。`pipeline.sh` の stdout / stderr 保存、`[REPORT]` / `[WARN]` の取り込み、ビルド所要時間、転送検証結果、コミット情報、ログ世代管理は現行実装済み機能として扱ってはならない。

### 仕様化済み・未実装のログ拡張

以下は CI ランナー拡張として仕様化済みだが、現行 `runner.py` には未実装である。

| 項目 | 内容 |
|------|------|
| ビルドログファイル | ビルドごとに `.build_logs/{id}.json` を作成する。 |
| stdout / stderr 保存 | `pipeline.sh` の標準出力・標準エラーをビルドログへ保存する。 |
| 変換レポート取り込み | `build_spec.py` が出力する `[REPORT]` 行をパースし、`tables_count`、`code_blocks_count` 等へ変換して保存する。 |
| 警告取り込み | `[WARN]` 行を配列として保存し、`warnings` 件数と整合させる。 |
| ビルド所要時間 | `started_at`、`finished_at`、`duration_seconds` を保存する。 |
| コミット情報 | ビルド対象 commit の SHA、message、author、date を保存する。 |
| 転送検証結果 | SSH 転送後整合性検証の結果として `transfer_verified` を保存する。 |
| 出力サイズ警告 | `OUTPUT_SIZE_WARN_MB` 超過時に `size_warn: true` を保存する。 |
| ログ世代管理 | `LOG_KEEP_N` を超過した `.build_logs/{id}.json` を古いものから削除する。 |

これらの拡張ログを実装する場合は、§10a の未実装範囲、§12 の拡張設定、§13 の拡張フロー、§22 の API レスポンス仕様と整合させる。

---

## 16. systemd タイマー

### `adlaire-ci.service`（oneshot）

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=deploy
ExecStart=/usr/bin/python3 /opt/adlaire-builder/runner.py
StandardOutput=journal
StandardError=journal
```

### `adlaire-ci.timer`（5分ごと定期実行）

```ini
[Unit]
Description=Adlaire CI Runner Timer
After=network-online.target
Wants=network-online.target

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
```

```bash
sudo systemctl enable --now adlaire-ci.timer  # タイマー登録・起動
sudo systemctl list-timers adlaire-ci          # 次回実行時刻確認
sudo journalctl -u adlaire-ci -f               # ログ確認
```

---

## 17. GitHub 側設定

| 項目 | 内容 |
|------|------|
| PAT スコープ | `contents: read`（読み取り専用）のみ |
| PAT の種類 | Fine-grained PAT（特定リポジトリのみ許可）を推奨 |
| Webhook 設定（ポーリング方式） | **不要**（デフォルト。`BRANCH_TARGETS` によるポーリングのみ使用する場合） |
| Webhook 設定（受信方式） | GitHub リポジトリ設定 → Webhooks → Add webhook で `POST /api/webhook` の URL・Secret を設定する（→ §22）。`push` イベントのみ選択を推奨。**外部公開エンドポイントが必要**（リバースプロキシ経由） |

---

## 18. 初回セットアップ手順

```bash
# 1. deploy ユーザー作成
sudo useradd -m -s /bin/bash deploy

# 2. 作業ディレクトリ作成
sudo mkdir -p /opt/adlaire-builder/repo
sudo chown -R deploy:deploy /opt/adlaire-builder

# 3. GitHub PAT を保存（Fine-grained PAT、contents: read のみ）
echo "<PAT>" | sudo -u deploy tee /opt/adlaire-builder/.github_token
sudo chmod 600 /opt/adlaire-builder/.github_token

# 3b. CI サーバー → 配信サーバー SSH 鍵設定
#     deploy ユーザーの SSH 鍵を生成（既存鍵がある場合はスキップ）
sudo -u deploy ssh-keygen -t ed25519 -f /home/deploy/.ssh/id_ed25519 -N ""
#     公開鍵を配信サーバーへ登録（配信サーバー側で実行）
#     cat /home/deploy/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
#     初回接続時の known_hosts 登録
sudo -u deploy ssh-keyscan -H <配信サーバーIP> >> /home/deploy/.ssh/known_hosts

# 4. SHA キャッシュファイルを初期化（JSON 形式）
echo '{"sha": ""}' | sudo -u deploy tee /opt/adlaire-builder/.last_sha
sudo chmod 600 /opt/adlaire-builder/.last_sha

# 5. build_spec.py を配置
sudo cp build_spec.py /opt/adlaire-builder/build_spec.py
sudo chown deploy:deploy /opt/adlaire-builder/build_spec.py

# 6. runner.py を配置
sudo cp runner.py /opt/adlaire-builder/runner.py
sudo chown deploy:deploy /opt/adlaire-builder/runner.py

# 7. systemd ユニットを登録・タイマー起動
sudo cp adlaire-ci.service /etc/systemd/system/
sudo cp adlaire-ci.timer   /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now adlaire-ci.timer

# 以降は仕様化済み・未実装の管理 API サーバー導入手順
# 8. api_server.py を配置
sudo cp api_server.py /opt/adlaire-builder/api_server.py
sudo chown deploy:deploy /opt/adlaire-builder/api_server.py

# 9. 認証情報ファイルを初期化（初期パスワード: admin）
sudo -u deploy python3 /opt/adlaire-builder/api_server.py --init-credentials
sudo chmod 600 /opt/adlaire-builder/.admin_credentials

# 10. systemd ユニットを登録・起動
sudo cp adlaire-admin.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now adlaire-admin
```

---

## 19. 管理 API サーバー 既知の制限

| 制限 | 詳細 |
|------|------|
| セッションはインメモリ管理 | 再起動で全セッションが消去される |
| HTTPS 非対応 | TLS ターミネーションは nginx 等リバースプロキシで行う。`api_server.py` 単体では HTTP のみ |
| シングルユーザー専用 | 現行実装はユーザー名固定（admin）。マルチユーザー対応は Part 1 §13 参照 |
| 並列リクエストの制限 | Python `http.server` ベースのため、高負荷並列リクエストには非対応 |

---

## 20. CI ランナー 既知の制限

| 制限 | 詳細 |
|------|------|
| ポーリング遅延 | 変更検出はタイマー間隔（デフォルト 5 分）に依存する。即時反応は不可（Webhook 受信で補完可能 → §22） |
| `pipeline.sh` のみ対応 | YAML 形式のパイプライン定義には非対応（シェルスクリプト固定） |
| ネットワーク断時の挙動 | GitHub API 失敗時は `API_RETRY_MAX` 回まで指数バックオフで再試行する。全試行失敗時のみ ERROR ログを記録してエントリをスキップする |
| `BRANCH_TARGETS` 直列処理 | 複数エントリはリスト順に順次処理する。並列処理には非対応 |
| ペンディングキューは再試行のみ | ペンディング再試行が連続失敗した場合の上限・放棄ポリシーは未定義 |
| Webhook 受信の外部公開 | `POST /api/webhook` は `api_server.py`（`127.0.0.1` バインド）のため、GitHub からの受信にはリバースプロキシが必要 |

---

### — 管理ツール —

## 21. 管理ツール システム構成

```
systemd timer
  └─ runner.py（変更検出・ビルド起動・SSH 転送）
       └─ SSH → 静的コンテンツ配信サーバー

api_server.py（常駐 HTTP サーバー、仕様化済み・未実装）

admin/index.html（標準管理ツール、仕様化済み・未実装）
  └─ adlaire-ci-sdk.js（SDK）─── HTTP ───► api_server.py
```

仕様化済み・未実装コンポーネント `api_server.py` は、Python 標準ライブラリ（`http.server`）で実装し、管理ツールからの API リクエストを受け付ける。`runner.py` とは独立して常駐する。

**`api_server.py` 設定値（スクリプト冒頭）：**

```python
HOST                 = "127.0.0.1"                                    # バインドアドレス（外部公開禁止）
PORT                 = 8765                                            # リッスンポート
CREDENTIALS_FILE     = "/opt/adlaire-builder/.admin_credentials"      # 認証情報ファイル
OUTPUT_URL           = "https://example.com/Adlaire-db-spec.html"     # 出力ファイルの公開 URL
HISTORY_FILE         = "/opt/adlaire-builder/.build_history"          # ビルド履歴ファイル
NOTIFY_CONFIG_FILE   = "/opt/adlaire-builder/.notify_config"          # Webhook 通知設定
SERVER_CONFIG_FILE   = "/opt/adlaire-builder/.server_config"          # サーバー設定
ACCESS_LOG_FILE      = "/opt/adlaire-builder/.access_log"             # ログイン履歴
NOTIFY_LOG_FILE      = "/opt/adlaire-builder/.notify_log"             # Webhook 送信履歴
WEBHOOK_SECRET_FILE  = "/opt/adlaire-builder/.webhook_secret"         # GitHub Webhook HMAC-SHA256 Secret（→ §22）
SNAPSHOT_DIR         = "/opt/adlaire-builder/.snapshots"              # スナップショット保存ディレクトリ（→ §14b）
LOG_LEVEL            = "INFO"
OWNER                = "<GitHubオーナー名>"                            # 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
REPO                 = "<リポジトリ名>"                                # 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
```

**systemd ユニット（常駐型、タイマー不要）：**

```ini
[Unit]
Description=Adlaire Admin API Server
After=network.target

[Service]
Type=simple
User=deploy
ExecStart=/usr/bin/python3 /opt/adlaire-builder/api_server.py
Restart=on-failure
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now adlaire-admin  # 登録・起動
sudo journalctl -u adlaire-admin -f        # ログ確認
```

---

## 22. バックエンド API 仕様

**ベース URL：** `http://localhost:{PORT}/api`
**認証：** `Authorization: Bearer {SESSION_TOKEN}`（`/api/login` で取得したセッショントークン）
**レスポンス形式：** JSON

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `POST` | `/api/login` | 不要 | ログイン（セッショントークン返却） |
| `POST` | `/api/logout` | 要 | ログアウト（セッション破棄） |
| `POST` | `/api/change-password` | 要 | パスワード変更 |
| `GET` | `/api/status` | 要 | 最終ビルド時刻・SHA・成否・実行中フラグを返す |
| `POST` | `/api/build` | 要 | 手動ビルドトリガー（`runner.py` を即時起動） |
| `GET` | `/api/logs?n=100&q=<keyword>` | 要 | 最新ビルドログを n 行返す（`q` 省略時は全行） |
| `GET` | `/api/logs/export` | 要 | ビルドログ全件を JSON 形式でエクスポートする |
| `POST` | `/api/logs/cleanup` | 要 | 保持期間（`log_retention_days`）を超えた `.build_logs/` エントリを削除する |
| `GET` | `/api/history?page=<n>&per_page=<n>` | 要 | 過去ビルド履歴一覧をページ指定で返す（省略時: `page=1`, `per_page=20`） |
| `POST` | `/api/reset-sha` | 要 | SHA キャッシュをクリアし次回強制ビルドを起こす |
| `GET` | `/api/sysinfo` | 要 | 出力ファイルサイズ・更新日時・稼働時間を返す |
| `GET` | `/api/health`        | 不要 | 死活監視用ヘルスチェック |
| `GET` | `/api/schedule` | 要 | systemd timer の次回実行予定時刻を返す |
| `GET` | `/api/notify-config` | 要 | Webhook 通知設定を返す |
| `POST` | `/api/notify-config` | 要 | Webhook 通知設定を更新する |
| `GET` | `/api/config` | 要 | サーバー設定を返す |
| `POST` | `/api/config` | 要 | サーバー設定を更新する |
| `GET` | `/api/pat-status` | 要 | GitHub PAT の有効性確認 |
| `GET` | `/api/access-log` | 要 | ログイン履歴一覧を返す |
| `GET` | `/api/stats?days=7` | 要 | ビルド統計（成功率・回数・平均間隔）を返す |
| `GET` | `/api/repo-info` | 要 | リポジトリ設定（OWNER/REPO/BRANCH/TARGET_FILE）を返す |
| `GET` | `/api/backup` | 要 | 全設定（通知設定・サーバー設定）を JSON 形式でエクスポートする |
| `POST` | `/api/restore` | 要 | JSON 形式の設定をインポートし全設定を上書き復元する |
| `POST` | `/api/notify-test` | 要 | Webhook にテスト通知を送信し疎通を確認する |
| `POST` | `/api/build/force` | 要 | SHA リセットとビルドをアトミックに実行する（強制ビルド） |
| `POST` | `/api/pat-verify` | 要 | GitHub API を呼び出し PAT の有効性をリアルタイム検証する |
| `GET` | `/api/history/{id}/log` | 要 | 指定ビルド ID のログを取得する |
| `POST` | `/api/build/cancel` | 要 | 実行中のビルドを強制停止する（`running: true` のときのみ有効） |
| `POST` | `/api/log-level` | 要 | `api_server.py` の LOG_LEVEL をランタイムで変更する |
| `POST` | `/api/pat-update` | 要 | `.github_token` ファイルを更新し PAT を差し替える |
| `GET` | `/api/dashboard` | 要 | ステータス・システム情報・統計・スケジュールを一括返却する |
| `GET` | `/api/notify-log` | 要 | Webhook 送信履歴（日時・イベント・HTTP ステータス・成否）を返す |
| `GET` | `/api/sessions` | 要 | 有効セッション一覧（作成日時・有効期限）を返す |
| `POST` | `/api/sessions/revoke-all` | 要 | 現セッション以外の全セッションを強制無効化する |
| `POST` | `/api/schedule/interval` | 要 | systemd タイマーのポーリング間隔を動的変更する |
| `GET` | `/api/logs/search?q=<keyword>&from=<date>&to=<date>&level=<warn\|error>` | 要 | 日付範囲・重大度を指定して過去ビルドログ（`.build_logs/`）を横断検索する |
| `GET` | `/api/output-meta` | 要 | 出力ファイルのサイズ・見出し数・生成日時・前回比サイズ差分を返す |
| `GET` | `/api/stats/timeline?days=30` | 要 | 日別ビルド成功/失敗件数の時系列配列を返す |
| `GET` | `/api/stats/build-duration?n=20` | 要 | 過去 N 件のビルド所要時間統計（平均・最小・最大・直近リスト）を返す |
| `GET` | `/api/webhook-events?limit=50&offset=0` | 要 | 受信 Webhook イベント一覧を新しい順にページネーション付きで返す（→ `.webhook_events.json`） |
| `POST` | `/api/circuit-breaker/reset` | 要 | サーキットブレーカーをリセットする（`open: false`・`consecutive_failures: 0` に戻しポーリングを再開） |
| `GET` | `/api/branch-config` | 要 | 現在有効なブランチターゲット設定を返す（`.branch_config` 存在時はその内容、なければ `BRANCH_TARGETS` のデフォルト値） |
| `POST` | `/api/branch-config` | 要 | ブランチターゲット設定を `.branch_config` へ書き込む（`runner.py` 再起動不要で次回ポーリングから反映） |
| `POST` | `/api/notify/weekly-summary` | 要 | 週次サマリー Webhook を即時手動送信する（過去 7 日間の統計を集計して送信） |
| `GET` | `/api/diagnostics` | 要 | PAT・GitHub API・出力ファイル・systemd・Webhook の一括自己診断結果を返す |
| `GET` | `/api/rate-limit` | 要 | GitHub API のレート制限残量・上限・リセット時刻を返す |
| `GET` | `/api/disk-usage` | 要 | ビルドログ合計・出力ファイルのディスク使用量を返す |
| `GET` | `/api/config-log` | 要 | 設定変更履歴（変更日時・種別・変更前後の値）を返す |
| `GET` | `/api/history/{id}/comment` | 要 | 指定ビルドのコメントを取得する |
| `POST` | `/api/history/{id}/comment` | 要 | 指定ビルドにコメントを付与・更新する |
| `POST` | `/api/repo-config` | 要 | リポジトリ監視設定（OWNER / REPO / BRANCH / TARGET_FILE）を更新する |
| `GET` | `/api/history/export` | 要 | ビルド履歴一覧を JSON 形式でエクスポートする |
| `POST` | `/api/history/{id}/flag` | 要 | 指定ビルドに重要フラグを設定・解除する |
| `POST` | `/api/history/{id}/rollback` | 要 | 指定ビルド ID のスナップショットから SSH 転送を再実行する（→ §14b） |
| `POST` | `/api/schedule/force-interval` | 要 | `FORCE_BUILD_INTERVAL`（強制再ビルド間隔）を動的変更する |
| `POST` | `/api/schedule/cooldown` | 要 | `BUILD_COOLDOWN_SECONDS`（ビルドクールダウン秒数）を動的変更する |
| `POST` | `/api/webhook` | 不要（Secret 検証） | GitHub push Webhook を受信し、署名検証後にビルドをトリガーする（→ §22 Webhook 受信仕様） |

**`POST /api/login` リクエスト / レスポンス：**
```json
// リクエスト
{ "password": "admin" }

// レスポンス
{ "token": "<session_token>", "must_change": "prompt" }
```

`must_change` の有効値：`"none"`（変更不要）| `"prompt"`（促す：初回ログイン時）| `"forced"`（強制：5 回目以降。変更完了まで管理画面の操作を制限）

**`POST /api/logout` リクエスト / レスポンス：**
```json
// リクエスト: なし（Bearer トークンのみ）
// レスポンス: 200
{ "message": "Logged out" }
```

**`POST /api/change-password` リクエスト / レスポンス：**
```json
// リクエスト
{ "current_password": "...", "new_password": "..." }
// レスポンス: 200
{ "message": "Password changed" }
```

**`GET /api/status` レスポンス例：**
```json
{
  "last_sha": "abc123",
  "last_build_at": "2026-09-14T10:00:00",
  "last_build_status": "success",
  "output_url": "https://example.com/Adlaire-db-spec.html",
  "running": false
}
```

`last_build_status` の有効値：`"success"` | `"failure"` | `"none"`（初回未実行時）
`running` の有効値：`true`（ビルド実行中）| `false`（待機中）
`running` の判定：`api_server.py` が `systemctl is-active adlaire-ci.service` を実行し、`active` の場合 `true` を返す。

**`GET /api/logs` レスポンス例：**
```json
{ "lines": ["2026-09-14T10:00:00 [INFO] Build start", "..."] }
```

**`GET /api/logs/export` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00",
  "lines": ["2026-09-14T10:00:00 [INFO] Build start", "..."]
}
```

**`GET /api/history` レスポンス例：**
```json
{
  "total": 42, "page": 1, "per_page": 20, "pages": 3,
  "history": [
    { "id": "b001", "build_at": "2026-09-15T10:00:00", "sha": "abc123", "status": "success", "output_size_bytes": 2048576, "trigger": "auto",    "duration_seconds": 42, "flagged": false, "tags": ["release"] },
    { "id": "b002", "build_at": "2026-09-14T18:30:00", "sha": "def456", "status": "failure", "output_size_bytes": null,    "trigger": "manual", "duration_seconds": 7,  "flagged": true,  "tags": [] }
  ]
}
```

`page` は 1 始まり。`per_page` の最大値は 100。範囲外ページを指定した場合は `history: []` を返す。

`id` はビルド実行時に生成するユニーク識別子（形式：`b{YYYYMMDDHHmmss}`）。`.build_logs/{id}.json` に対応するログファイルが保存される。

**`POST /api/reset-sha` レスポンス例：**
```json
{ "message": "SHA reset" }
```

**`GET /api/sysinfo` レスポンス例：**
```json
{
  "output_size_bytes": 2048576,
  "output_mtime": "2026-09-15T10:00:00",
  "uptime_seconds": 86400
}
```

**`GET /api/schedule` レスポンス例：**
```json
{ "next_run_at": "2026-09-15T10:05:00", "interval": "5min", "paused": false, "allowed_hours": { "from": 9, "to": 18 } }
```

`paused` が `true` のとき、ポーリングは停止中で `next_run_at` は `null` を返す。

`allowed_hours`：自動ビルドを許可する時間帯（UTC）。`null` = 無制限。`from` 以上 `to` 未満の時刻のみビルドを実行する。許可時間帯外のポーリングでは変更を検出しても実行を保留し、次の許可時間帯に入った時点で実行する。

**`POST /api/schedule/allowed-hours` リクエスト / レスポンス：**
```json
// 設定
{ "from": 9, "to": 18 }
// 解除（無制限に戻す）
{ "from": null, "to": null }
// レスポンス: 200
{ "message": "Allowed hours updated", "allowed_hours": { "from": 9, "to": 18 } }
```

**`POST /api/schedule/pause` / `POST /api/schedule/resume` レスポンス例：**
```json
{ "message": "Schedule paused" }
{ "message": "Schedule resumed" }
```

既に一時停止中に `pause`、または稼働中に `resume` を呼び出した場合は `409 Conflict` を返す。

> **責務分担：** Webhook 通知の**送信責務は `runner.py`** にある。`runner.py` はビルド完了時に `.notify_config` を読み込んで Webhook を送信する。`api_server.py`（通知 API）は設定の読み書きのみを担い、自身では通知を送信しない。

**`GET /api/notify-config` レスポンス例：**
```json
{
  "webhooks": [
    { "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }
  ],
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": false, "to": [], "on": [] }
}
```

`secret`：Webhook 署名シークレット。未設定時は `null`、設定済み時は `"***"`（マスク）を返す（→ 16E 参照）。

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"weekly_summary"`（定期サマリー送信時）。複数指定可。

`summary`：定期サマリー通知の設定。`enabled: true` のとき指定スケジュールで統計サマリーを Webhook 送信する。`interval` の有効値：`"daily"` | `"weekly"`。`hour` は 0〜23（UTC）。`day_of_week` は `"weekly"` 時のみ有効（0 = 日曜〜6 = 土曜）。

**`POST /api/notify-summary` レスポンス例：**
```json
{ "message": "Summary sent", "webhook_url": "https://hooks.example.com/..." }
```

即時サマリー送信。Webhook 未設定または無効時は `422` を返す。

`payload_template`：Webhook 送信 JSON ペイロードのテンプレート文字列。`null` = デフォルトペイロードを使用。テンプレート内で使用可能な変数は以下の通り。

| 変数 | 内容 |
|------|------|
| `{{id}}` | ビルド ID |
| `{{status}}` | ビルド結果（`success` / `failure`） |
| `{{sha}}` | 対象コミット SHA |
| `{{duration_seconds}}` | ビルド所要時間（秒） |
| `{{build_at}}` | ビルド実行日時（ISO 8601） |

`secret`（16E）：Webhook 送信時の HMAC-SHA256 署名用シークレット文字列。設定時はリクエストヘッダーに `X-Adlaire-Signature: sha256=<hmac>` を付与する。`null` = 署名なし。`GET /api/notify-config` で返却する際、設定済みの場合は `"***"` でマスクし、未設定の場合は `null` を返す。`POST /api/notify-config` で更新可能。

**`GET /api/config` レスポンス例：**
```json
{ "log_max_lines": 500, "history_max_count": 100, "build_timeout_seconds": 300, "log_retention_days": 30, "pat_expires_at": null, "snapshots_keep": 5, "queue_max_size": 3 }

`pat_expires_at`：PAT の有効期限日（`YYYY-MM-DD` 形式）。`null` = 未設定。`GET /api/diagnostics` の `pat` 項目で 7 日以内なら `"warn"`、期限当日以前なら `"error"` に変更。
```

**`GET /api/health` レスポンス例：**
```json
{
  "status": "ok",
  "last_build_at": "2026-09-15T10:00:00",
  "last_build_status": "success",
  "last_deploy_at": "2026-09-15T10:01:00",
  "last_deploy_status": "success",
  "pending_transfers": 0,
  "uptime_seconds": 86400
}
```

- `status`：常に `"ok"`（サーバーが応答している限り）
- `last_build_at`：最終ビルド完了日時（未実行時 `null`）
- `last_build_status`：`"success"` | `"failure"` | `"none"`
- `last_deploy_at`：最終 SSH 転送完了日時（未実行時 `null`）
- `last_deploy_status`：`"success"` | `"failure"` | `"skipped"` | `"none"`
- `pending_transfers`：ペンディングキューのエントリ数
- `uptime_seconds`：`api_server.py` 起動からの経過秒数

**`GET /api/pat-status` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:00:00" }
```

**`GET /api/access-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00", "result": "success" },
    { "at": "2026-09-15T09:00:00", "result": "failure" }
]}
```

**`GET /api/stats` レスポンス例：**
```json
{
  "days": 7,
  "total_builds": 42,
  "success_count": 40,
  "failure_count": 2,
  "success_rate": 0.952,
  "avg_interval_minutes": 240,
  "avg_duration_seconds": 38,
  "max_duration_seconds": 91
}
```

**`GET /api/repo-info` レスポンス例：**
```json
{
  "owner": "fqwink",
  "repo": "Adlaire-Design-System",
  "branch": "main",
  "target_file": "adlaire-db-spec.md"
}
```

**`GET /api/backup` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00",
  "notify_config": {
    "webhooks": [{ "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }],
    "on": ["failure"],
    "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
    "email": { "enabled": false, "to": [], "on": [] }
  },
  "server_config": { "log_max_lines": 500, "history_max_count": 100 }
}
```

**`POST /api/notify-test` レスポンス例：**
```json
{ "message": "Test notification sent", "webhook_url": "https://hooks.example.com/..." }
```

**`POST /api/build/force` レスポンス例：**
```json
{ "message": "SHA reset and build triggered" }
```

**`POST /api/pat-verify` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:05:00", "scopes": ["contents:read"] }
```

**`GET /api/history/{id}/log` レスポンス例：**
```json
{
  "id": "b001",
  "build_at": "2026-09-15T10:00:00",
  "sha": "abc123",
  "status": "success",
  "output_size_bytes": 2048576,
  "trigger": "auto",
  "comment": null,
  "flagged": false,
  "tags": ["release"],
  "lines": ["2026-09-15T10:00:00 [INFO] Build start", "..."]
}
```

**`POST /api/restore` リクエスト / レスポンス：**
```json
// リクエスト（GET /api/backup と同一形式）
{
  "notify_config": {
    "webhooks": [{ "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }],
    "on": ["failure"],
    "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
    "email": { "enabled": false, "to": [], "on": [] }
  },
  "server_config": { "log_max_lines": 500, "history_max_count": 100 }
}
// レスポンス: 200
{ "message": "Restored" }
```

**`GET /api/notify-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00", "event": "failure", "http_status": 200, "result": "success", "attempt": 1, "error": null },
    { "at": "2026-09-14T18:30:00", "event": "start",   "http_status": 500, "result": "failure", "attempt": 3, "error": "HTTP 500" }
]}
```

`event` の有効値：`"start"` | `"success"` | `"failure"` | `"weekly_summary"`（`GET /api/notify-config` の `on` と同一）。`result` の有効値：`"success"` | `"failure"`（Webhook 送信の成否）。

**`GET /api/sessions` レスポンス例：**
```json
{ "sessions": [
    { "created_at": "2026-09-15T09:00:00", "expires_at": "2026-09-15T17:00:00", "current": true },
    { "created_at": "2026-09-15T08:00:00", "expires_at": "2026-09-15T16:00:00", "current": false }
]}
```

**`POST /api/sessions/revoke-all` レスポンス例：**
```json
{ "message": "All other sessions revoked", "revoked_count": 1 }
```

**`POST /api/schedule/interval` リクエスト / レスポンス：**
```json
// リクエスト
{ "interval_seconds": 300 }
// レスポンス: 200
{ "message": "Interval updated", "interval_seconds": 300 }
```

**`POST /api/build/cancel` レスポンス例：**
```json
{ "message": "Build cancelled" }
```

`running: false` のときに呼び出した場合は `409 Conflict` → `{"error": "No build is running"}` を返す。

**`GET /api/build/stream` — SSE ストリーミング：**

`Content-Type: text/event-stream` で接続を維持し、ビルドログを逐次配信する。認証トークンをクエリパラメータ（`?token=<session_token>`）で受け付ける。

```
data: {"type": "log",  "line": "2026-09-15T10:00:01 [INFO] Build start"}

data: {"type": "log",  "line": "2026-09-15T10:00:42 [INFO] Build success"}

data: {"type": "end",  "status": "success", "duration_seconds": 42}
```

`type` の有効値：`"log"`（ログ行）| `"end"`（ビルド完了）。ビルドが未実行時に接続した場合は即時 `{"type": "end", "status": null}` を送信して切断する。

**`POST /api/log-level` リクエスト / レスポンス：**
```json
// リクエスト
{ "level": "DEBUG" }
// レスポンス: 200
{ "message": "Log level changed", "level": "DEBUG" }
```

`level` の有効値：`"DEBUG"` | `"INFO"` | `"WARNING"` | `"ERROR"`

**`POST /api/pat-update` リクエスト / レスポンス：**
```json
// リクエスト
{ "token": "github_pat_..." }
// レスポンス: 200
{ "message": "PAT updated" }
```

**`GET /api/dashboard` レスポンス例：**
```json
{
  "status": {
    "last_sha": "abc123",
    "last_build_at": "2026-09-15T10:00:00",
    "last_build_status": "success",
    "output_url": "https://example.com/Adlaire-db-spec.html",
    "running": false
  },
  "sysinfo": {
    "output_size_bytes": 2048576,
    "output_mtime": "2026-09-15T10:00:00",
    "uptime_seconds": 86400
  },
  "stats": {
    "days": 7,
    "total_builds": 42,
    "success_count": 40,
    "failure_count": 2,
    "success_rate": 0.952,
    "avg_interval_minutes": 240
  },
  "schedule": { "next_run_at": "2026-09-15T10:05:00", "interval": "5min", "paused": false },
  "alerts": [
    { "level": "warn", "message": "PAT expires in 5 days" }
  ]
}
```

`alerts`：診断異常（`warn` / `error`）および PAT 期限切れ間近の場合に項目を返す。異常なしのときは空配列 `[]`。`level` の有効値：`"warn"` | `"error"`。パネルのナビゲーション項目に `alerts` の最高深刻度（`error` > `warn`）のバッジを表示する。

**`GET /api/logs/search` レスポンス例：**
```json
{
  "query": "ERROR",
  "from": "2026-09-10",
  "to": "2026-09-15",
  "results": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00", "lines": ["2026-09-15T10:00:01 [ERROR] Build failed"] },
    { "id": "b20260912183000", "build_at": "2026-09-12T18:30:00", "lines": ["2026-09-12T18:30:05 [ERROR] Timeout"] }
  ]
}
```

`from` / `to` は `YYYY-MM-DD` 形式。省略時は全期間。`q` 省略時は全行返却。
`level=warn` で `[WARN]` 行のみ、`level=error` で `[ERROR]` 行のみを絞り込む。省略時は全レベルを返却する。

**`GET /api/output-meta` レスポンス例：**
```json
{
  "size_bytes": 2048576,
  "mtime": "2026-09-15T10:00:00",
  "heading_count": 342,
  "size_diff_bytes": 1024,
  "tables_count": 128,
  "code_blocks_count": 64,
  "build_warnings": ["未対応記法: admonition (3箇所)"],
  "size_warn": false
}
```

`size_diff_bytes`：前回ビルド時との差分（正＝増加、負＝減少、`null`＝比較不能）。
前回サイズは `.build_history` の直近エントリに記録された `output_size_bytes` フィールドから取得する。

`tables_count` / `code_blocks_count`：直近ビルドの変換レポート（§8）より取得。ビルド前は `null`。
`build_warnings`：直近ビルドで発生した警告メッセージの配列（§8 参照）。ビルド前は空配列 `[]`。
値は runner.py が `.build_logs/{id}.json` から最新エントリを読み取って返す。

**`GET /api/stats/timeline` レスポンス例：**
```json
{
  "days": 30,
  "timeline": [
    { "date": "2026-09-15", "success": 3, "failure": 0 },
    { "date": "2026-09-14", "success": 2, "failure": 1 },
    { "date": "2026-09-13", "success": 4, "failure": 0 }
  ]
}
```

日付降順。`days` 日分のうちビルドが 0 件の日はエントリなし。

**`GET /api/stats/build-duration` レスポンス例：**

クエリパラメータ `n`（デフォルト 20）で対象件数を指定する。`.build_logs/{id}.json` の `duration_seconds` フィールドを集計する。

```json
{
  "n": 20,
  "count": 18,
  "avg_seconds": 38.5,
  "min_seconds": 22,
  "max_seconds": 67,
  "recent": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00", "duration_seconds": 42, "status": "success" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00", "duration_seconds": 7,  "status": "failure" }
  ]
}
```

`count` は `duration_seconds` が記録されているビルドの件数（`n` 以下）。`recent` は新しい順。

**`GET /api/webhook-events` レスポンス例：**

クエリパラメータ `limit`（デフォルト 50、上限 200）と `offset` でページネーションする。`.webhook_events.json` を逆順（新しい順）で返す。

```json
{
  "total": 128,
  "offset": 0,
  "limit": 50,
  "events": [
    {
      "timestamp": "2026-09-15T10:00:00",
      "delivery_id": "abc-123-def",
      "event": "push",
      "ref": "refs/heads/main",
      "sha": "abc123def456",
      "build_triggered": true
    }
  ]
}
```

**`POST /api/circuit-breaker/reset` レスポンス例：**

```json
{ "message": "Circuit breaker reset", "open": false, "consecutive_failures": 0 }
```

サーキットブレーカーが既に閉じている（`open: false`）場合も同じレスポンスを返す（冪等）。

**`GET /api/branch-config` レスポンス例：**

```json
{
  "source": "file",
  "branches": [
    {
      "branch": "main",
      "src": "/opt/adlaire-builder/Adlaire-db-spec.md",
      "out": "/opt/adlaire-builder/Adlaire-db-spec.html",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
```

`source` は設定の出所を示す。`.branch_config` ファイルが存在する場合は `"file"`、存在しない場合（`BRANCH_TARGETS` デフォルト値を使用）は `"default"` を返す。

**`POST /api/branch-config` リクエスト / レスポンス：**

```json
// リクエスト（GET /api/branch-config の branches と同一形式）
{
  "branches": [
    {
      "branch": "main",
      "src": "/opt/adlaire-builder/Adlaire-db-spec.md",
      "out": "/opt/adlaire-builder/Adlaire-db-spec.html",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
// レスポンス: 200
{ "message": "Branch config updated", "branches_count": 1 }
```

`branches` が空配列 `[]` の場合は `.branch_config` ファイルを削除し、`BRANCH_TARGETS` のデフォルト値に戻す（`source: "default"` に戻る）。変更は次回ポーリング周回から反映される。

**`POST /api/notify/weekly-summary` レスポンス例：**

```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

`on: ["weekly_summary"]` 設定の Webhook 宛先がない場合は `422 Unprocessable Entity` を返す。

**`.build_logs/{id}.json` 追加フィールド（ビルド所要時間・コミット情報・サイズ警告）：**

```json
{
  "id": "b20260915100000",
  "started_at": "2026-09-15T09:59:18",
  "finished_at": "2026-09-15T10:00:00",
  "duration_seconds": 42,
  "status": "success",
  "commit_sha": "abc123def456",
  "commit_message": "fix: typo in §4.3 description",
  "commit_author": "Kazuhiro Kurata",
  "commit_at": "2026-09-15T09:58:00",
  "size_warn": false,
  "broken_links": 0,
  "heading_skips": 0,
  "reading_time": 87
}
```

`started_at` は `pipeline.sh` 実行開始時刻、`finished_at` は完了（または失敗）時刻。`duration_seconds` は整数（小数点以下切り捨て）。
`commit_sha` / `commit_message` / `commit_author` / `commit_at` はコミット情報取得 API（§13）の結果を記録する。API 失敗時は `null`。
`size_warn` は出力 HTML が `OUTPUT_SIZE_WARN_MB` 超過時 `true`、それ以外 `false`。`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`。

**`GET /api/diagnostics` レスポンス例：**
```json
{
  "checked_at": "2026-09-15T10:00:00",
  "items": [
    { "name": "pat",         "status": "ok",   "message": "PAT is valid" },
    { "name": "github_api",  "status": "ok",   "message": "GitHub API reachable" },
    { "name": "output_file", "status": "ok",   "message": "Output file exists (2.0 MB)" },
    { "name": "systemd",     "status": "ok",   "message": "adlaire-ci.service is active" },
    { "name": "webhook",     "status": "warn", "message": "Webhook URL not configured" }
  ]
}
```

`status` の有効値：`"ok"` | `"warn"` | `"error"`。

**`GET /api/rate-limit` レスポンス例：**
```json
{ "limit": 5000, "remaining": 4823, "reset_at": "2026-09-15T11:00:00", "used": 177 }
```

**`GET /api/disk-usage` レスポンス例：**
```json
{
  "build_logs_bytes": 10485760,
  "build_logs_count": 42,
  "output_file_bytes": 2048576,
  "total_bytes": 12534296
}
```

**`GET /api/config-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00", "type": "server_config", "diff": { "log_max_lines": [500, 1000] }, "diff_text": "- log_max_lines: 500\n+ log_max_lines: 1000" },
    { "at": "2026-09-14T18:00:00", "type": "notify_config", "diff": { "enabled": [false, true] },    "diff_text": "- enabled: false\n+ enabled: true" },
    { "at": "2026-09-13T12:00:00", "type": "repo_config",   "diff": { "branch": ["main", "develop"] }, "diff_text": "- branch: main\n+ branch: develop" }
]}
```

`type` の有効値：`"server_config"` | `"notify_config"` | `"repo_config"`
`diff` の形式：`{ フィールド名: [変更前, 変更後] }`。設定変更時に `.config_log` へ追記する。
`diff_text`：`diff` を `"- key: old_value\n+ key: new_value"` 形式の文字列に変換したフィールド。複数フィールド変更時は行を連結する。設定変更時に `diff` と同時に記録する。

**`GET /api/history/{id}/comment` レスポンス例：**
```json
{ "id": "b20260914183000", "comment": "ネットワーク障害による失敗。再ビルド済み。", "updated_at": "2026-09-14T19:00:00" }
```

コメント未設定時は `"comment": null`。コメントは `.build_logs/{id}.json` の `comment` フィールドに保存する。

**`POST /api/history/{id}/comment` リクエスト / レスポンス：**
```json
// リクエスト
{ "comment": "ネットワーク障害による失敗。再ビルド済み。" }
// レスポンス: 200
{ "message": "Comment saved" }
```

**`POST /api/history/{id}/flag` リクエスト / レスポンス：**
```json
// リクエスト
{ "flagged": true }
// レスポンス: 200
{ "message": "Flag updated" }
```

フラグは `.build_logs/{id}.json` の `flagged` フィールドに保存する。

**`POST /api/history/{id}/tags` リクエスト / レスポンス：**
```json
// リクエスト
{ "tags": ["release", "hotfix"] }
// レスポンス: 200
{ "message": "Tags updated" }
```

タグは `.build_logs/{id}.json` の `tags` フィールド（`string[]`）に保存する。空配列 `[]` を指定するとタグをすべて削除する。

**`GET /api/tokens` レスポンス例：**
```json
{ "tokens": [
    { "id": "tok001", "label": "監視用", "scope": "read", "created_at": "2026-09-15T10:00:00", "last_used_at": "2026-09-15T11:00:00" }
]}
```

**`POST /api/tokens` リクエスト / レスポンス：**
```json
// リクエスト
{ "label": "監視用", "scope": "read" }
// レスポンス: 201
{ "id": "tok001", "token": "act_...", "label": "監視用", "scope": "read", "created_at": "2026-09-15T10:00:00" }
```

`token` はレスポンス時のみ返却し、以後は取得不可。`scope` の有効値：`"read"`（読み取り専用）。読み取り専用トークンは `GET` 系エンドポイントのみ許可し、`POST` / `DELETE` 系は `403 Forbidden` を返す。トークンは `Authorization: Bearer <token>` ヘッダーで送信する。

**`DELETE /api/tokens/{id}` レスポンス例：**
```json
{ "message": "Token revoked" }
```

---

### スナップショット（14A）

ビルド成功時に出力ファイルを `.snapshots/` へ自動保存する。保持世代数は `GET /api/config` の `snapshots_keep`（デフォルト `5`、`0` = 機能無効）で制御し、超過した古い世代は自動削除する。

**`GET /api/snapshots` レスポンス例：**
```json
{ "snapshots": [
    { "id": "snap001", "build_id": "b20260915100000", "saved_at": "2026-09-15T10:00:00", "size_bytes": 2048576 },
    { "id": "snap002", "build_id": "b20260914183000", "saved_at": "2026-09-14T18:30:00", "size_bytes": 2031616 }
]}
```

**`GET /api/snapshots/{id}/download`**
バイナリレスポンス。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="Adlaire-db-spec.html"` を付与する。

**`DELETE /api/snapshots/{id}` レスポンス例：**
```json
{ "message": "Snapshot deleted" }
```

**`POST /api/history/{id}/rollback` リクエスト / レスポンス：**
```json
// リクエスト: なし（パスパラメーターのみ）
// レスポンス: 200
{ "message": "Rollback started", "build_id": "b20260914183000" }
```

- `.snapshots/{id}/` が存在しない場合は `404 Not Found` を返す
- ロールバックは非同期で SSH 転送を実行する（`running: true` 中は `409 Conflict` を返す）
- 転送成功時は `.build_history` に `trigger: "rollback"` のエントリを追記する

---

### Webhook 受信仕様（22-W）

`POST /api/webhook` は GitHub からの push イベントを受信し、署名検証後にビルドをトリガーする。認証ヘッダー（`Authorization: Bearer`）は不要だが、`X-Hub-Signature-256` ヘッダーによる HMAC-SHA256 署名検証が必須である。

**署名検証：**
```python
# api_server.py の実装例
import hmac, hashlib
expected = "sha256=" + hmac.new(secret.encode(), body, hashlib.sha256).hexdigest()
if not hmac.compare_digest(expected, request_header["X-Hub-Signature-256"]):
    return 403
```

- Secret は `WEBHOOK_SECRET_FILE`（`/opt/adlaire-builder/.webhook_secret`）から読み込む
- Secret 未設定時（ファイル不在）は Webhook 受信を `501 Not Implemented` で拒否する

**`POST /api/webhook` リクエストヘッダー：**
```
X-GitHub-Event: push
X-Hub-Signature-256: sha256=<hmac_hex>
Content-Type: application/json
```

**`POST /api/webhook` レスポンス：**
```json
// 200: ビルドトリガー成功
{ "message": "Build triggered", "ref": "refs/heads/main" }

// 400: ペイロード不正（ref フィールド欠損等）
{ "error": "Invalid payload" }

// 403: 署名検証失敗
{ "error": "Invalid signature" }

// 409: ビルド実行中
{ "error": "Build already running" }

// 501: Secret 未設定
{ "error": "Webhook not configured" }
```

**ビルドトリガー条件：**
- `ref` フィールドが `BRANCH_TARGETS` のいずれかの `branch` と一致する push イベントのみビルドをトリガーする
- 一致する branch が存在しない場合は `200` + `{ "message": "No matching branch" }` を返す（ビルドはしない）

**イベントログ：**
署名検証成功後、受信イベントを `.webhook_events.json` に JSON Lines 形式（1行1イベント）で追記する。ビルドの実行可否によらず全受信イベントを記録する。

記録フォーマット（1行）：
```json
{"timestamp": "2026-09-15T10:00:00", "delivery_id": "abc-123-def", "event": "push", "ref": "refs/heads/main", "sha": "abc123def456", "build_triggered": true}
```

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `timestamp` | string (ISO 8601) | イベント受信日時 |
| `delivery_id` | string | `X-GitHub-Delivery` ヘッダー値 |
| `event` | string | `X-GitHub-Event` ヘッダー値（`"push"` 等） |
| `ref` | string | ペイロードの `ref` フィールド |
| `sha` | string | ペイロードの `after` フィールド（push 後の HEAD SHA） |
| `build_triggered` | boolean | ビルドをトリガーしたか（`true` / `false`） |

**Webhook Secret 設定 API：**

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `GET` | `/api/webhook-config` | 要 | Secret 設定状態（設定済み/未設定）を返す |
| `POST` | `/api/webhook-config` | 要 | `WEBHOOK_SECRET_FILE` を更新して Secret を設定する |

**`GET /api/webhook-config` レスポンス例：**
```json
{ "configured": true }
```

**`POST /api/webhook-config` リクエスト / レスポンス：**
```json
// リクエスト
{ "secret": "<新しいSecret文字列>" }
// レスポンス: 200
{ "message": "Webhook secret updated" }
```

---

### スケジュール強制再ビルド仕様（22-S）

**`POST /api/schedule/force-interval` リクエスト / レスポンス：**
```json
// リクエスト
{ "hours": 24 }
// 無効化
{ "hours": 0 }
// レスポンス: 200
{ "message": "Force build interval updated", "hours": 24 }
```

- `runner.py` 側の `FORCE_BUILD_INTERVAL` を動的変更する（`.server_config` に保存し、起動時に読み込む）
- `hours` は 0 以上の整数。0 で機能無効化

**`POST /api/schedule/cooldown` リクエスト / レスポンス：**
```json
// リクエスト
{ "seconds": 120 }
// 無効化
{ "seconds": 0 }
// レスポンス: 200
{ "message": "Build cooldown updated", "seconds": 120 }
```

- `runner.py` 側の `BUILD_COOLDOWN_SECONDS` を動的変更する（`.server_config` に保存し、起動時に読み込む）
- `seconds` は 0 以上の整数。0 で機能無効化

---

### メンテナンスモード（14B）

メンテナンスモード有効中は手動ビルド（`POST /api/build`・`POST /api/build/force`）・スケジュールビルドの両方を拒否し、`503 Service Unavailable` + `{ "error": "maintenance" }` を返す。`GET /api/health` は制限対象外とする。

**`GET /api/maintenance` レスポンス例：**
```json
{ "enabled": false, "reason": null, "since": null }
```
メンテナンス中は `{ "enabled": true, "reason": "定期メンテナンス", "since": "2026-09-15T10:00:00" }`。

**`POST /api/maintenance/enable` リクエスト / レスポンス：**
```json
// リクエスト
{ "reason": "定期メンテナンス" }
// レスポンス: 200
{ "message": "Maintenance mode enabled", "since": "2026-09-15T10:00:00" }
```

**`POST /api/maintenance/disable` レスポンス：**
```json
{ "message": "Maintenance mode disabled" }
```

---

### IP アクセス制限（14C）

`allow` に IPv4 アドレスまたは CIDR 表記のリストを設定する。空リストは制限なし（全接続許可）を意味する。制限に一致しない接続元からのリクエストは `403 Forbidden` を返す。`GET /api/health` は制限対象外とする。設定は `.access_control` に保存する。

**`GET /api/access-control` レスポンス例：**
```json
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
```
制限なしの場合: `{ "allow": [] }`

**`POST /api/access-control` リクエスト / レスポンス：**
```json
// リクエスト
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
// レスポンス: 200
{ "message": "Access control updated", "allow": ["192.168.1.0/24", "10.0.0.1"] }
```

---

### ビルドフック（14E）

ビルド実行の直前（`pre`）・直後（`post`）に任意のシェルコマンドを実行する。フック設定は `.hooks` に保存する。

- `pre` フックが失敗（`exit_code != 0`）し `abort_on_failure: true` の場合、ビルドを中断しステータスを `hook_error` とする。
- `post` フックは `abort_on_failure` 設定に関わらずビルド結果（`success` / `failure`）を変更しない。
- フックの実行ログは `.build_logs/{build_id}_hook_{id}.json` に保存する。

**`GET /api/hooks` レスポンス例：**
```json
{ "hooks": [
    { "id": "h001", "phase": "pre",  "command": "echo build start", "enabled": true, "abort_on_failure": true },
    { "id": "h002", "phase": "post", "command": "echo build end",   "enabled": true, "abort_on_failure": false }
]}
```

**`POST /api/hooks` リクエスト / レスポンス：**
```json
// リクエスト
{ "phase": "pre", "command": "echo build start", "abort_on_failure": true }
// レスポンス: 201
{ "id": "h001", "phase": "pre", "command": "echo build start", "enabled": true, "abort_on_failure": true }
```

**`DELETE /api/hooks/{id}` レスポンス：**
```json
{ "message": "Hook deleted" }
```

**`GET /api/hooks/{id}/log` レスポンス例：**
```json
{ "id": "h001", "runs": [
    { "build_id": "b20260915100000", "ran_at": "2026-09-15T10:00:00", "exit_code": 0, "output": "build start\n" },
    { "build_id": "b20260914183000", "ran_at": "2026-09-14T18:30:00", "exit_code": 1, "output": "Error: command not found\n" }
]}
```
`runs` は直近 20 件を返す（新しい順）。

---

### カスタムアラートルール（15A）

`GET /api/dashboard` 取得時にルールを評価し、条件を満たすものを `alerts` 配列へ自動追加する。`.alert_rules` に保存する。

対応メトリクス：`success_rate_7d`（7日間成功率 %）/ `avg_duration_seconds`（7日間平均ビルド時間 秒）/ `last_build_age_hours`（最終ビルドからの経過時間 時間）/ `disk_usage_bytes`（ディスク使用量 バイト）
対応演算子：`lt`（未満）/ `gt`（超過）/ `lte`（以下）/ `gte`（以上）

**`GET /api/alert-rules` レスポンス例：**
```json
{ "rules": [
    { "id": "r001", "metric": "success_rate_7d", "operator": "lt", "threshold": 90, "level": "warn", "message": "7日間成功率が90%を下回っています" },
    { "id": "r002", "metric": "last_build_age_hours", "operator": "gt", "threshold": 48, "level": "error", "message": "48時間以上ビルドが実行されていません" }
]}
```

**`POST /api/alert-rules` リクエスト / レスポンス：**
```json
// リクエスト
{ "metric": "avg_duration_seconds", "operator": "gt", "threshold": 120, "level": "warn", "message": "平均ビルド時間が2分を超えています" }
// レスポンス: 201
{ "id": "r003", "metric": "avg_duration_seconds", "operator": "gt", "threshold": 120, "level": "warn", "message": "平均ビルド時間が2分を超えています" }
```

**`DELETE /api/alert-rules/{id}` レスポンス：**
```json
{ "message": "Alert rule deleted" }
```

---

### 自動タグ付けルール（15B）

ビルド完了時に条件式を評価し、マッチしたルールのタグを `.build_history` のエントリへ自動追記する（手動タグと共存する）。`.tag_rules` に保存する。

条件式で使用可能な変数：`status`（`"success"` / `"failure"`）/ `duration_seconds`（整数）/ `trigger`（`"auto"` / `"manual"` / `"force"`）
演算子：`==`・`!=`・`>`・`<`・`>=`・`<=`

**`GET /api/tag-rules` レスポンス例：**
```json
{ "rules": [
    { "id": "t001", "condition": "status == 'failure'", "tags": ["要確認"] },
    { "id": "t002", "condition": "duration_seconds > 120", "tags": ["低速"] },
    { "id": "t003", "condition": "trigger == 'force'", "tags": ["強制実行"] }
]}
```

**`POST /api/tag-rules` リクエスト / レスポンス：**
```json
// リクエスト
{ "condition": "status == 'success'", "tags": ["green"] }
// レスポンス: 201
{ "id": "t004", "condition": "status == 'success'", "tags": ["green"] }
```

**`DELETE /api/tag-rules/{id}` レスポンス：**
```json
{ "message": "Tag rule deleted" }
```

---

### 出力ファイルチェックサム（15C）

ビルド成功時に出力 HTML ファイルの SHA-256 ハッシュを算出し `.build_history` の該当エントリに `output_sha256` として記録する。

**`GET /api/output-meta` レスポンス変更（`sha256` フィールド追加）：**
```json
{ "size_bytes": 2048576, "mtime": "2026-09-15T10:00:00", "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }
```

**`GET /api/history/{id}/log` レスポンス変更（`output_sha256` フィールド追加）：**
```json
{ "id": "b001", "build_at": "2026-09-15T10:00:00", "sha": "abc123", "status": "success",
  "output_size_bytes": 2048576, "output_sha256": "e3b0c44298fc1c149afbf4c8996fb924...",
  "trigger": "auto", "comment": null, "flagged": false, "tags": ["release"],
  "lines": ["2026-09-15T10:00:00 [INFO] Build start", "..."] }
```

**`POST /api/verify-output` レスポンス例：**
```json
// 一致時
{ "match": true, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "e3b0c44298fc1c149afbf4c8996fb924..." }
// 不一致時
{ "match": false, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "f4a2d5591c8f3a742f902e3b6f7c1c3d..." }
```
`expected` は `.build_history` の最終成功エントリに記録された `output_sha256`。現在の出力ファイルが存在しない場合は `404` を返す。

---

### ビルドパイプライン設定（15D）

`runner.py` がビルド実行時に `.pipeline_config` を読み込み、`build_spec.py` の呼び出しに `extra_args`・`env` を適用する。`.pipeline_config` に保存する。

**`GET /api/pipeline-config` レスポンス例：**
```json
{ "extra_args": ["--verbose"], "env": { "DEBUG": "1" } }
```
初期値（未設定時）：`{ "extra_args": [], "env": {} }`

**`POST /api/pipeline-config` リクエスト / レスポンス：**
```json
// リクエスト
{ "extra_args": ["--verbose"], "env": { "DEBUG": "1" } }
// レスポンス: 200
{ "message": "Pipeline config updated" }
```

---

### 運用ノート（15E）

システム全体の運用メモを Markdown テキストで保存・取得する。`.notes` に保存する。認証必須。

**`GET /api/notes` レスポンス例：**
```json
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00\nPAT 更新期限: 2026-12-01", "updated_at": "2026-09-15T10:00:00" }
```
初回（未作成）時：`{ "content": "", "updated_at": null }`

**`POST /api/notes` リクエスト / レスポンス：**
```json
// リクエスト
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00" }
// レスポンス: 200
{ "message": "Notes updated", "updated_at": "2026-09-15T10:00:00" }
```

---

### メール通知（SMTP）（16B）

Webhook に加えてメールでビルド結果を通知できる機能。SMTP 接続設定は `.smtp_config` に、パスワードは `.smtp_secret`（パーミッション 600）に分離して保存する。`GET /api/notify-config` のレスポンスに `email` セクションを追加する。

**`GET /api/smtp-config` レスポンス例：**
```json
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true, "password_set": true }
```
`password_set`：`.smtp_secret` が存在するかを真偽値で返す。パスワード本体は返却しない。
未設定時：`{ "host": null, "port": 587, "user": null, "tls": true, "from": null, "to": [], "on": [], "enabled": false, "password_set": false }`

**`POST /api/smtp-config` リクエスト / レスポンス：**
```json
// リクエスト（変更するフィールドのみ指定可。password フィールドは省略可能）
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "password": "s3cr3t", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true }
// レスポンス: 200
{ "message": "SMTP config updated" }
```
`on` の有効値：`"start"` / `"success"` / `"failure"`

**`POST /api/smtp-test` レスポンス例：**
```json
// 成功時
{ "result": "success", "message": "Test email sent to ops@example.com" }
// 失敗時
{ "result": "failure", "message": "Connection refused: smtp.example.com:587" }
```
SMTP 未設定または `enabled: false` の場合は `422` を返す。

**`GET /api/notify-config` への追加（`email` セクション）：**
```json
{
  "webhooks": [ { "url": "...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": "***" } ],
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": true, "to": ["ops@example.com"], "on": ["failure"] }
}
```

---

### ビルドキューイング（16C）

ビルド実行中に `POST /api/build`・`POST /api/build/force` を受信した場合、キューに追加して順次実行する。キューの最大長は `GET /api/config` の `queue_max_size`（デフォルト `3`、`0` = キューなし）で制御する。キューが満杯の場合は `429 Too Many Requests` + `{ "error": "queue_full" }` を返す。

**`GET /api/queue` レスポンス例：**
```json
{ "queued": [
    { "id": "q001", "trigger": "manual", "queued_at": "2026-09-15T10:01:00" },
    { "id": "q002", "trigger": "auto",   "queued_at": "2026-09-15T10:02:00" }
  ], "max_size": 3 }
```
キューが空の場合：`{ "queued": [], "max_size": 3 }`

**`DELETE /api/queue` レスポンス：**
```json
{ "message": "Queue cleared", "cleared_count": 2 }
```
実行中のビルドは停止しない（`POST /api/build/cancel` を別途使用する）。

`POST /api/config` に `queue_max_size`（整数、`0` = キューなし）追加。

---

### ダッシュボードウィジェットカスタマイズ（16D）

ダッシュボードに表示するウィジェットの種類・順序を設定できる機能。`.dashboard_layout` に保存する。

有効なウィジェット識別子：`status`（ビルド状態）/ `stats`（統計サマリー）/ `schedule`（次回実行）/ `alerts`（アラート）/ `disk`（ディスク使用量）/ `rate_limit`（GitHub API レート制限）/ `snapshots`（スナップショット件数）/ `maintenance`（メンテナンス状態）/ `queue`（キュー状態）

**`GET /api/dashboard-layout` レスポンス例：**
```json
{ "widgets": ["status", "alerts", "stats", "schedule", "disk"] }
```
未設定時はデフォルト順（全ウィジェット）を返す。

**`POST /api/dashboard-layout` リクエスト / レスポンス：**
```json
// リクエスト
{ "widgets": ["status", "alerts", "schedule", "stats"] }
// レスポンス: 200
{ "message": "Dashboard layout updated" }
```
`widgets` に未知の識別子が含まれる場合は `400` を返す。

---

**`POST /api/logs/cleanup` レスポンス例：**
```json
{ "message": "Cleanup completed", "deleted_count": 12 }
```

`log_retention_days` が `0` の場合は削除せず `deleted_count: 0` を返す。

**`GET /api/history/export` レスポンス：**

`Content-Type: application/json` で返却される。

```json
{ "export_at": "2026-09-15T10:00:00", "history": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00", "sha": "abc123", "status": "success", "trigger": "auto", "duration_seconds": 42, "flagged": false, "tags": ["release"], "comment": "" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00", "sha": "def456", "status": "failure", "trigger": "manual", "duration_seconds": 7, "flagged": true, "tags": [], "comment": "ネットワーク障害による失敗" }
]}
```

**`POST /api/repo-config` リクエスト / レスポンス：**
```json
// リクエスト（変更するフィールドのみ指定可）
{ "owner": "fqwink", "repo": "Adlaire-Design-System", "branch": "main", "target_file": "adlaire-db-spec.md" }
// レスポンス: 200
{ "message": "Repo config updated" }
```

設定は `.repo_config` に保存し、`GET /api/repo-info` もこのファイルを参照する。

`trigger` の有効値：`"auto"`（定期ポーリング）| `"manual"`（`POST /api/build`）| `"force"`（`POST /api/build/force`）

**エラーレスポンス形式：**

| ステータス | 条件 | レスポンス |
|-----------|------|-----------|
| `401` | 認証失敗・セッション期限切れ | `{"error": "Unauthorized"}` |
| `404` | リソース不存在 | `{"error": "Not Found"}` |
| `409` | 状態競合（`POST /api/build/cancel` でビルド未実行時等） | `{"error": "No build is running"}` |
| `422` | Webhook 送信失敗（`POST /api/notify-test`） | `{"error": "Webhook delivery failed"}` |
| `429` | ビルドキューが満杯（`POST /api/build` / `POST /api/build/force`） | `{"error": "queue_full"}` |
| `500` | サーバー内部エラー | `{"error": "Internal Server Error"}` |
| `503` | メンテナンスモード中のリクエスト（`POST /api/build` 等） | `{"error": "maintenance"}` |

---

## 23. JavaScript SDK 仕様

本節は、仕様化済み・未実装の `adlaire-ci-sdk.js` に関する仕様である。

**ファイル：** `adlaire-ci-sdk.js`（単一ファイル、外部依存なし）
**モジュール形式：** ES Module（`import` / `export`）

```js
class AdlaireCI {
  constructor({ baseUrl })
  // this._token でセッショントークンを管理。login() 後の全リクエストに自動付与

  login(password)                               // POST /api/login → {token, must_change}; this._token にセット
  logout()                                      // POST /api/logout; this._token をクリア
  changePassword(currentPassword, newPassword)  // POST /api/change-password

  getStatus()               // GET /api/status               → Promise<StatusObject>
  triggerBuild()            // POST /api/build               → Promise<void>
  getLogs(n = 100, q = '')  // GET /api/logs?n={n}&q={q}     → Promise<{lines: string[]}>
  getHistory({ page = 1, perPage = 20 } = {}) // GET /api/history?page={page}&per_page={perPage} → Promise<HistoryPageObject>
  resetSha()                // POST /api/reset-sha           → Promise<void>
  getSysinfo()              // GET /api/sysinfo              → Promise<SysinfoObject>
  getSchedule()             // GET /api/schedule             → Promise<ScheduleObject>
  getNotifyConfig()         // GET /api/notify-config        → Promise<NotifyConfig>
  setNotifyConfig(config)   // POST /api/notify-config       → Promise<void>
  getConfig()               // GET /api/config               → Promise<ConfigObject>
  setConfig(config)         // POST /api/config              → Promise<void>
  health()                  // GET /api/health               → Promise<{status: string}>
  getPatStatus()            // GET /api/pat-status           → Promise<PatStatusObject>
  getAccessLog()            // GET /api/access-log           → Promise<{log: AccessRecord[]}>
  getStats(days = 7)        // GET /api/stats?days={days}    → Promise<StatsObject>
  exportLogs()              // GET /api/logs/export          → Promise<{exported_at: string, lines: string[]}>
  cleanupLogs()             // POST /api/logs/cleanup        → Promise<{message: string, deleted_count: number}>
  getRepoInfo()             // GET /api/repo-info            → Promise<RepoInfoObject>
  backup()                  // GET /api/backup               → Promise<BackupObject>
  restore(config)           // POST /api/restore             → Promise<void>
  notifyTest()              // POST /api/notify-test         → Promise<{message: string, webhook_url: string}>
  notifySummary()           // POST /api/notify-summary      → Promise<{message: string, webhook_url: string}>
  buildForce()              // POST /api/build/force         → Promise<void>
  patVerify()               // POST /api/pat-verify          → Promise<PatVerifyObject>
  getHistoryLog(id)         // GET /api/history/{id}/log     → Promise<HistoryLogObject>
  cancelBuild()             // POST /api/build/cancel        → Promise<void>
  resetCircuitBreaker()     // POST /api/circuit-breaker/reset → Promise<{message: string, open: boolean, consecutive_failures: number}>
  streamBuild(onLine, onEnd) // GET /api/build/stream (SSE)  → EventSource（onLine(line), onEnd({status, duration_seconds}) コールバック）
  setLogLevel(level)        // POST /api/log-level           → Promise<{message: string, level: string}>
  updatePat(token)          // POST /api/pat-update          → Promise<void>
  getDashboard()            // GET /api/dashboard            → Promise<DashboardObject>
  getNotifyLog()            // GET /api/notify-log           → Promise<{log: NotifyRecord[]}>
  getSessions()             // GET /api/sessions             → Promise<{sessions: SessionRecord[]}>
  revokeAllSessions()       // POST /api/sessions/revoke-all → Promise<{message: string, revoked_count: number}>
  setScheduleInterval(seconds) // POST /api/schedule/interval → Promise<{message: string, interval_seconds: number}>
  pauseSchedule()              // POST /api/schedule/pause   → Promise<void>
  resumeSchedule()             // POST /api/schedule/resume  → Promise<void>
  setAllowedHours(from, to)    // POST /api/schedule/allowed-hours {from, to} → Promise<void>
  clearAllowedHours()          // POST /api/schedule/allowed-hours {from:null, to:null} → Promise<void>
  setForceInterval(hours)      // POST /api/schedule/force-interval → Promise<{message: string, force_interval_hours: number}>
  searchLogs(q = '', from = '', to = '') // GET /api/logs/search?q={q}&from={from}&to={to} → Promise<SearchResult>
  getOutputMeta()              // GET /api/output-meta       → Promise<OutputMetaObject>
  getStatsTimeline(days = 30)  // GET /api/stats/timeline?days={days} → Promise<TimelineObject>
  getStatsBuildDuration(n = 10) // GET /api/stats/build-duration?n={n} → Promise<BuildDurationStats>
  getDiagnostics()             // GET /api/diagnostics       → Promise<DiagnosticsObject>
  getRateLimit()               // GET /api/rate-limit        → Promise<RateLimitObject>
  getDiskUsage()               // GET /api/disk-usage        → Promise<DiskUsageObject>
  getConfigLog()               // GET /api/config-log        → Promise<{log: ConfigLogRecord[]}>
  getBranchConfig()            // GET /api/branch-config     → Promise<{source: string, branches: BranchTargetRecord[]}>
  setBranchConfig(branches)    // POST /api/branch-config    → Promise<{message: string, branches_count: number}>
  notifyWeeklySummary()        // POST /api/notify/weekly-summary → Promise<{message: string, period: string, success_count: number, failure_count: number, success_rate: number}>
  getWebhookEvents(limit = 50, offset = 0) // GET /api/webhook-events?limit={limit}&offset={offset} → Promise<{events: WebhookEventRecord[], total: number}>
  getHistoryComment(id)        // GET /api/history/{id}/comment  → Promise<CommentObject>
  setHistoryComment(id, comment) // POST /api/history/{id}/comment → Promise<void>
  setRepoConfig(config)        // POST /api/repo-config      → Promise<void>
  exportHistory()              // GET /api/history/export    → Promise<ExportObject>（JSON）
  setHistoryFlag(id, flagged)  // POST /api/history/{id}/flag → Promise<void>
  setHistoryTags(id, tags)     // POST /api/history/{id}/tags → Promise<void>
  rollbackHistory(id)          // POST /api/history/{id}/rollback → Promise<void>
  getTokens()                  // GET /api/tokens            → Promise<{tokens: TokenRecord[]}>
  createToken(label)           // POST /api/tokens           → Promise<TokenCreateResult>
  revokeToken(id)              // DELETE /api/tokens/{id}    → Promise<void>
  // 14A スナップショット
  getSnapshots()               // GET /api/snapshots         → Promise<{snapshots: SnapshotRecord[]}>
  downloadSnapshot(id)         // GET /api/snapshots/{id}/download → Promise<Blob>
  deleteSnapshot(id)           // DELETE /api/snapshots/{id} → Promise<void>
  // 14B メンテナンスモード
  getMaintenance()             // GET /api/maintenance       → Promise<MaintenanceObject>
  enableMaintenance(reason)    // POST /api/maintenance/enable → Promise<{message: string, since: string}>
  disableMaintenance()         // POST /api/maintenance/disable → Promise<void>
  // 14C IP アクセス制限
  getAccessControl()           // GET /api/access-control    → Promise<{allow: string[]}>
  setAccessControl(allowList)  // POST /api/access-control   → Promise<{message: string, allow: string[]}>
  // 14E フック
  getHooks()                   // GET /api/hooks             → Promise<{hooks: HookRecord[]}>
  addHook(phase, command, abortOnFailure = true) // POST /api/hooks → Promise<HookRecord>
  deleteHook(id)               // DELETE /api/hooks/{id}     → Promise<void>
  getHookLog(id)               // GET /api/hooks/{id}/log    → Promise<{id: string, runs: HookRunRecord[]}>
  // 15A アラートルール
  getAlertRules()              // GET /api/alert-rules       → Promise<{rules: AlertRule[]}>
  addAlertRule(metric, operator, threshold, level, message) // POST /api/alert-rules → Promise<AlertRule>
  deleteAlertRule(id)          // DELETE /api/alert-rules/{id} → Promise<void>
  // 15B 自動タグ付けルール
  getTagRules()                // GET /api/tag-rules         → Promise<{rules: TagRule[]}>
  addTagRule(condition, tags)  // POST /api/tag-rules        → Promise<TagRule>
  deleteTagRule(id)            // DELETE /api/tag-rules/{id} → Promise<void>
  // 15C チェックサム
  verifyOutput()               // POST /api/verify-output    → Promise<{match: boolean, expected: string, actual: string}>
  // 15D パイプライン設定
  getPipelineConfig()          // GET /api/pipeline-config   → Promise<PipelineConfig>
  setPipelineConfig(config)    // POST /api/pipeline-config  → Promise<void>
  // 15E 運用ノート
  getNotes()                   // GET /api/notes             → Promise<{content: string, updated_at: string|null}>
  setNotes(content)            // POST /api/notes            → Promise<{message: string, updated_at: string}>
  // 16B メール通知
  getSmtpConfig()              // GET /api/smtp-config       → Promise<SmtpConfig>
  setSmtpConfig(config)        // POST /api/smtp-config      → Promise<void>
  smtpTest()                   // POST /api/smtp-test        → Promise<{result: string, message: string}>
  // 16C ビルドキュー
  getQueue()                   // GET /api/queue             → Promise<{queued: QueueEntry[], max_size: number}>
  clearQueue()                 // DELETE /api/queue          → Promise<{message: string, cleared_count: number}>
  // 16D ダッシュボードレイアウト
  getDashboardLayout()         // GET /api/dashboard-layout  → Promise<{widgets: string[]}>
  setDashboardLayout(widgets)  // POST /api/dashboard-layout → Promise<void>
}

export { AdlaireCI };
```

全メソッドは `Promise` を返す（`streamBuild` は `EventSource` を返す）。HTTP エラー（4xx / 5xx）は `Error` としてスロー。`401` 受信時はセッション期限切れとして `this._token` をクリアする。合計 95 メソッド。

---

## 24. 標準管理ツール 仕様

本節は、仕様化済み・未実装の `admin/index.html` に関する仕様である。

**ファイル構成：**
```
/opt/adlaire-builder/admin/
├── index.html          # 管理画面（単一ファイル完結、仕様化済み・未実装）
└── adlaire-ci-sdk.js   # SDK（標準管理ツールに同梱、仕様化済み・未実装）
```

**画面構成：**

| パネル | 表示内容 | 表示条件 |
|-------|---------|---------|
| ログイン | パスワード入力フォーム | 未ログイン時のみ |
| パスワード変更 | 現在・新パスワード入力フォーム | `must_change: true` 時（強制時は他パネル非表示） |
| ステータス | 最終ビルド時刻・SHA・成否・出力ファイルリンク・ダッシュボードウィジェット編集モード（表示するウィジェットをチェックボックスで選択・並び替え・保存） | ログイン済み |
| 手動実行 | ビルドトリガーボタン・強制ビルドボタン・キャンセルボタン（実行中のみ有効）・SHA リセットボタン・実行結果表示・リアルタイムログ表示エリア（SSE ストリーミング）・キュー状態表示（待機中件数・クリアボタン） | ログイン済み |
| ログビューア | 最新ビルドログ（n 行・キーワードフィルター・ログレベルフィルターボタン（INFO / WARNING / ERROR / DEBUG）・JSON エクスポートボタン・横断検索フォーム（期間指定）・検索結果一覧） | ログイン済み |
| ビルド履歴 | 過去ビルド一覧（日時・SHA・成否・トリガー種別・所要時間・重要フラグ列・タグ列・ログ表示リンク・コメント入力欄）・フラグ付きのみ表示フィルター・タグフィルター・ページネーション UI・JSON エクスポートボタン | ログイン済み |
| システム情報 | 出力ファイルサイズ・更新日時・稼働時間・ディスク使用量（ログ合計・出力ファイル）・PAT 即時検証ボタン・PAT 更新フォーム・PAT 有効期限表示（設定フォーム・期限切れ間近で警告表示）・GitHub API レート制限表示 | ログイン済み |
| 通知設定     | Webhook 一覧（追加/削除/ラベル/有効無効切り替え/リトライ回数・間隔設定/シークレット入力欄）・通知条件設定（ビルド開始時・成功時・失敗時）・各 Webhook ペイロードテンプレート編集フォーム（変数一覧表示）・テスト送信ボタン・定期サマリー設定（間隔・時刻・曜日・即時送信ボタン）・送信履歴（試行回数・エラー内容列含む）・メール通知セクション（SMTP 設定フォーム・宛先リスト・通知条件・テスト送信ボタン） | ログイン済み |
| 設定         | ログ保持行数・履歴保持件数の設定変更・ビルドタイムアウト設定・ログレベル変更（INFO / DEBUG）・ログ保持期間（日数、0 = 無制限）・スナップショット保持世代数設定・ビルドキュー最大長設定・手動クリーンアップボタン・設定変更履歴（変更日時・項目・変更前後の値）・IP アクセス制限セクション（許可 IP / CIDR 一覧・追加フォーム・削除ボタン）・フック設定セクション（pre / post フック一覧・追加フォーム・実行ログリンク・有効無効切り替え）・アラートルール設定セクション（メトリクス・演算子・しきい値・レベル・メッセージの入力フォーム・ルール一覧・削除ボタン）・自動タグ付けルールセクション（条件式・タグ入力フォーム・ルール一覧・削除ボタン）・パイプライン設定セクション（追加引数入力欄・環境変数テーブル） | ログイン済み |
| アクセスログ | ログイン履歴（日時・成否）           | ログイン済み |
| 統計         | ビルド回数・成功率・平均間隔・平均・最大ビルド時間・日別時系列データ（グラフ表示対応） | ログイン済み |
| リポジトリ情報 | 監視対象リポジトリ・ブランチ・ファイルの確認・設定変更フォーム（OWNER / REPO / BRANCH / TARGET_FILE）・ポーリング間隔変更フォーム・ポーリング一時停止／再開ボタン・許可時間帯設定（from〜to、解除ボタン）・メンテナンスモード有効化フォーム（理由入力）・解除ボタン・現在の状態表示 | ログイン済み |
| セッション管理 | 有効セッション一覧・全セッション強制無効化ボタン | ログイン済み |
| システム診断   | PAT・GitHub API・出力ファイル・systemd・Webhook の診断項目一覧（ok / warn / error）・診断実行ボタン・アラートバッジ（`GET /api/dashboard` の `alerts` に基づき warn / error を表示）・出力整合性チェック項目（`POST /api/verify-output` 結果表示）・メンテナンスモード中はバナーを全パネル上部に表示 | ログイン済み |
| ビルド比較     | ビルド履歴から 2 件を選択・ログ並列表示・差分ハイライト（クライアントサイド処理） | ログイン済み |
| API トークン管理 | 発行済みトークン一覧（ラベル・スコープ・作成日時・最終使用日時）・新規発行フォーム（ラベル入力・スコープ選択）・発行時のみトークン文字列を表示・失効ボタン | ログイン済み |
| 運用ノート     | Markdown レンダリング表示・編集モード切替・保存ボタン・最終更新日時表示 | ログイン済み |
| スナップショット | ビルド成果物の世代一覧（最大`snapshots_keep`件）・個別ダウンロード・削除 | ログイン済み |
| メンテナンス   | メンテナンスモードの有効化（理由テキスト付き）・無効化・状態・開始時刻表示 | ログイン済み |
| アクセス制御   | 許可 IP / CIDR 一覧・CIDR 追加フォーム・削除ボタン（ブロック時は 403 を返す） | ログイン済み |
| フック         | Pre/Post ビルドフック一覧・コマンド追加・削除・実行ログ（直近 N 件）確認 | ログイン済み |

**パスワード変更フロー：**
- `must_change: "prompt"`: パスワード変更パネルを表示。他パネルも操作可能
- `must_change: "forced"`: パスワード変更パネルのみ表示。変更完了後に通常画面へ遷移

**カスタマイズポイント：**
- SDK の `baseUrl` は `<script>` タグ内の設定変数で外出し
- CSS カスタムプロパティで外観変更可能（ADS トークン準拠）
- 各パネルは独立した `<section>` 単位で差し替え可能な構造とする

---

## 25. 認証 実装仕様

**認証情報ファイル形式（JSON）：**
```json
{
  "password_hash": "<pbkdf2_hmac_sha256_hex>",
  "salt": "<hex>",
  "login_count": 0
}
```

**ハッシュアルゴリズム：** Python 標準ライブラリ `hashlib.pbkdf2_hmac`
```python
hashlib.pbkdf2_hmac('sha256', password.encode('utf-8'), bytes.fromhex(salt), 260000)
```

**セッショントークン生成：**
```python
import secrets
token = secrets.token_hex(32)  # 256bit ランダムトークン
```

**セッション管理：** `api_server.py` 内のインメモリ辞書で管理。有効期限 8 時間。再起動で全セッション破棄。同一ユーザーの複数同時セッションを許容する。

**セッション期限切れ時：** `401 Unauthorized` を返す。クライアント（SDK）は `this._token` をクリアし、再ログインを促す。

**ログインフロー：**
```
POST /api/login
  └─ パスワードハッシュ検証
       ├─ 失敗 → 401
       └─ 成功 → login_count + 1 → ファイル更新
                  ├─ login_count == 1 → must_change: "prompt"（促す）
                  ├─ login_count >= 5 → must_change: "forced"（強制）
                  └─ それ以外      → must_change: "none"
                  → セッショントークン生成・返却
```

**パスワード変更時：** `login_count` を 0 にリセット。新しい salt を生成しハッシュを更新。変更完了後に現セッション以外のセッションを破棄。

**`--init-credentials` オプション：** `api_server.py` を `--init-credentials` 引数で起動した場合、初期パスワード `admin` で `.admin_credentials` を生成して終了する（HTTP サーバーは起動しない）。

---

## 26. セットアップ・アップデート手順

> **安定版ポリシー：** タグ付き安定版リリース（例：`v1.0.0`）のみをサポートする。開発ブランチ（`main` 等）の直接追従は非対応。`git pull` は使用しない。

### §26.1 要件

| 項目 | 要件 |
|------|------|
| Python | 3.9 以上（標準ライブラリのみ、追加インストール不要） |
| init システム | systemd（Linux） |
| バージョン管理 | git |
| ネットワーク | デプロイ先への SSH 接続（SSH 転送機能を使用する場合） |

### §26.2 設定変数

スクリプト内で以下の変数をカスタマイズする。

| 変数 | デフォルト値 | 説明 |
|------|------------|------|
| `REPO_URL` | —（必須） | GitHub 等のリポジトリ URL |
| `INSTALL_DIR` | `/opt/adlaire-builder` | インストール先ディレクトリ |
| `SERVICE_USER` | `root` | systemd サービスの実行ユーザー |
| `VERSION` | —（必須） | セットアップ・アップデート対象の安定版タグ（例：`v1.0.0`） |

### §26.3 初回セットアップ手順

```bash
# ── 変数設定 ──────────────────────────────────────────
REPO_URL="https://github.com/<owner>/<repo>.git"
INSTALL_DIR="/opt/adlaire-builder"
VERSION="v1.0.0"

# ── 1. リポジトリ取得 ─────────────────────────────────
git clone "$REPO_URL" "$INSTALL_DIR"
git -C "$INSTALL_DIR" checkout "$VERSION"

# ── 2. 必要ディレクトリ作成 ───────────────────────────
mkdir -p "$INSTALL_DIR/.build_logs"
mkdir -p "$INSTALL_DIR/.snapshots"

# ── 3. 初期認証情報生成（初期パスワード: admin）────────
python3 "$INSTALL_DIR/api_server.py" --init-credentials

# ── 4. パーミッション設定 ─────────────────────────────
chmod 600 "$INSTALL_DIR/.admin_credentials"

# ── 5. systemd サービスファイル配置 ───────────────────
# §26.4 のファイル内容を /etc/systemd/system/ に配置した上で:
systemctl daemon-reload

# ── 6. サービス有効化・起動 ───────────────────────────
systemctl enable adlaire-ci.timer adlaire-api
systemctl start  adlaire-ci.timer adlaire-api

# ── 7. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer adlaire-api
```

### §26.4 systemd サービスファイル

**`/etc/systemd/system/adlaire-ci.service`**（`runner.py`）：

```ini
[Unit]
Description=Adlaire CI Runner

[Service]
Type=oneshot
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/bin/python3 /opt/adlaire-builder/runner.py
```

**`/etc/systemd/system/adlaire-ci.timer`**（`runner.py` 定期起動タイマー）：

```ini
[Unit]
Description=Adlaire CI Runner Timer

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
Unit=adlaire-ci.service

[Install]
WantedBy=timers.target
```

**`/etc/systemd/system/adlaire-api.service`**（`api_server.py`）：

```ini
[Unit]
Description=Adlaire CI API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/adlaire-builder
ExecStart=/usr/bin/python3 /opt/adlaire-builder/api_server.py
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

`User` / `WorkingDirectory` / `ExecStart` のパスは §26.2 の設定変数に合わせて変更する。

### §26.5 アップデート手順

`git pull` は使用しない。安定版タグを指定してチェックアウトし、サービスを再起動する。

```bash
# ── 変数設定 ──────────────────────────────────────────
INSTALL_DIR="/opt/adlaire-builder"
NEW_VERSION="v1.2.0"

# ── 1. 最新タグ一覧を確認 ─────────────────────────────
git -C "$INSTALL_DIR" fetch --tags
git -C "$INSTALL_DIR" tag --list --sort=-v:refname

# ── 2. 対象バージョンへ切り替え ───────────────────────
git -C "$INSTALL_DIR" checkout "$NEW_VERSION"

# ── 3. サービス再起動 ─────────────────────────────────
systemctl restart adlaire-ci.timer adlaire-api

# ── 4. 起動確認 ───────────────────────────────────────
systemctl status adlaire-ci.timer adlaire-api
```

### §26.6 サービス操作リファレンス

| 操作 | コマンド |
|------|---------|
| 状態確認 | `systemctl status adlaire-ci.timer adlaire-api` |
| 起動 | `systemctl start adlaire-ci.timer adlaire-api` |
| 停止 | `systemctl stop adlaire-ci.timer adlaire-api` |
| 再起動 | `systemctl restart adlaire-ci.timer adlaire-api` |
| ログ確認（runner） | `journalctl -u adlaire-ci.service -f` |
| ログ確認（API） | `journalctl -u adlaire-api -f` |
