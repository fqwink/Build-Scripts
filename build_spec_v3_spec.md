# build_spec_v3.py / runner.py — ビルド・CI ランナー仕様ドキュメント

**対象スクリプト：** `build_spec_v3.py`（ビルド）/ `runner.py`（CI ランナー）  
**出力ファイル：** `Adlaire-db-spec.html`  
**スクリプトバージョン：** v3（Adlaire Design System ブルートークン正式採用）  
**ドキュメントバージョン：** V.N（累積インクリメント、リセット禁止 → Part 2 §2 参照）  
**最終更新：** 2026-09-14  

---

## 本ドキュメントの構成

| Part | 名称 | 責務の問い | 記載する内容 |
|------|------|-----------|------------|
| Part 1 | 方針 | **なぜ・何を** | 目的、設計思想、方向性の原則。変更頻度が低く、判断の拠り所となる指針 |
| Part 2 | ポリシー | **しなければならない／してはならない** | 遵守義務のある規則・制約・禁止事項。セキュリティ要件・運用ルール・バージョン管理規則 |
| Part 3 | 仕様 | **どのように** | 実装の具体的詳細。要件・構成・API・アルゴリズム・設定値・手順 |

新しい記載内容は「この内容はどの責務の問いに答えるか」を基準に Part を決定する。

---

# Part 1 — 方針
> 目的・設計思想・方向性の原則を定める。「なぜこう作るか」に答える。

## 1. 目的

`build_spec_v3.py` は、Adlaire DB 仕様書の Markdown ソースを単一の自己完結型 HTML ドキュメントへ変換する Python スクリプトである。外部ライブラリに依存せず、標準ライブラリ（`re`、`html`、`unicodedata`）のみで動作する。

- 14,000 行超の大規模 Markdown 仕様書を、快適に閲覧できる HTML ドキュメントサイトへ変換する
- CSS・JS をすべてインラインに埋め込み、単一 HTML ファイルとして配布可能にする
- Adlaire Design System（ADS）のトークンを採用し、一貫したデザイン言語を維持する

## 2. 開発方針

本プロジェクトは **仕様駆動開発（Spec-Driven Development）** を採用する。

- **仕様書が唯一の真実（Single Source of Truth）**：実装の追加・変更はすべて本ドキュメントへの反映を先行させる
- **仕様→実装の順序**：設計上の決定は本ドキュメントに記録してから `build_spec_v3.py` に反映する
- **仕様との乖離は不整合**：スクリプトの動作が本仕様書と食い違う場合、どちらかに誤りがある
- **仕様書のバージョン管理**：スクリプトの変更履歴と仕様書の改訂履歴は同期して管理する

## 3. デザイン方針

docs.rs / MDN に倣った技術ドキュメントレイアウト。14,000 行超の仕様書を快適に閲覧するため、**構造の明快さ**と**情報密度への耐性**を最優先とする。

- ヘッダーのみアクセントカラーを使う。コンテンツ・サイドバーは中性色ベース
- CSS カスタムプロパティは [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)（`Tokens/`）定義の `--adlaire-*` トークンのみ使用
- **ライトモード固定**（`prefers-color-scheme` 非対応、ダークモードなし）
- 外部フォント不使用。システムフォントスタックで日本語環境の可読性を確保

### — CI ランナー —

## 4. CI ランナーの目的

GitHub API（Git Blobs API）を定期的にポーリングし、対象ファイルの変更を検出してビルドパイプラインを自動実行する自己ホスト型 CI ランナー（`runner.py`）。GitHub Actions・Webhook・外部 CI サービスへの依存をゼロにする。

- GitHub Git Trees API で対象ファイルの blob SHA を取得し、前回 SHA と比較して変更を検出する
- 変更検出時のみ Git Blobs API でファイル本文を取得し、ビルドを実行する
- 外部公開エンドポイント・リバースプロキシ不要
- 標準ライブラリのみで実装し、`pip install` 不要

## 5. CI ランナーの開発方針

- **仕様書が唯一の真実**：実装の変更は本ドキュメントへの反映を先行させる
- **ゼロ外部依存**：Python 標準ライブラリのみ使用
- **単一ファイル実装**：`runner.py` 1 ファイルで完結
- **シンプル性優先**：HTTP サーバー不要。1 回実行して終了する oneshot 設計
- **差分検出**：SHA キャッシュにより変更がない場合はビルドをスキップ

## 6. GitHub Actions との対応関係

| GitHub Actions | 内製ランナー |
|---|---|
| `.github/workflows/*.yml` | `.ci/pipeline.sh` |
| `runs-on: ubuntu-latest` | 自前サーバー（固定） |
| `steps:` の各ステップ | shell の各コマンド |
| `secrets.*` | PAT ファイル（`.github_token`）・サーバー上の環境変数 |
| Actions Marketplace | なし（自前実装のみ） |
| push トリガー | systemd タイマーによる定期ポーリング + SHA 差分検出 |

---

# Part 2 — ポリシー
> 遵守義務のある規則と制約を定める。「何をしなければならないか／してはならないか」に答える。

## 1. デザイントークン準拠

> ⚠️ **準拠義務：** `--adlaire-*` トークンの値は [Adlaire Design System](https://github.com/fqwink/Adlaire-Design-System)（`Tokens/` ディレクトリ）で定義された値に準拠すること。

- スクリプト側での独自トークンの追加・変更は行わない
- ライトモード固定のため、ダークモード用トークンブロックは不要
- トークン値の変更は ADS 側のアップデートに追従する形でのみ実施する

## 2. バージョン管理

**スクリプトバージョン（`v3`）とドキュメントバージョン（`V.N`）は独立して管理する。** スクリプトバージョンはビルドロジックの世代を示し、ドキュメントバージョンはコンテンツの累積改訂番号を示す。

| 項目 | 内容 |
|------|------|
| 形式 | `V.N`（`V` は固定、`N` は正の整数） |
| 更新方針 | コンテンツの変更・追記のたびに `N` を 1 以上インクリメントする |
| リセット禁止 | スクリプトバージョン（`v3`→`v4` 等）が上がっても `N` は継続する。`V.1` に戻してはならない |
| 例 | `V.205` → `V.206` → `V.207`（スクリプトが v4 に移行しても `V.1` に戻さない） |

## 3. カスタマイズ可能範囲

以下の項目はスクリプト内で変更可能な設定ポイントである。

| 設定項目 | 場所 | 変更方法 |
|---------|------|---------|
| 入出力パス | スクリプト冒頭の `SRC` / `OUT` 定数 | 値を書き換えて再実行 |
| デザイントークン値 | `PAGE` f-string 内 `:root { }` ブロック | ADS 準拠の範囲内で変更可 |
| ドキュメントタイトル | `PAGE` f-string 内 `<title>` タグ | 任意の文字列に変更可 |
| ヘッダー表示名 | `PAGE` f-string 内 `<span class="hdr-title">` | 任意の文字列に変更可 |
| バージョンバッジ | `PAGE` f-string 内 `<span class="hdr-ver">` | `V.N` 形式で累積インクリメント |
| TOC 対象見出しレベル | `build_toc()` 内のフィルタ行（`lv <= 3`） | 上限レベルを変更可 |

### — CI ランナー —

## 4. CI ランナー セキュリティポリシー

- GitHub PAT（Personal Access Token）はファイル（`.github_token`）に保存し、パーミッションを `600` に設定する。スクリプト内にハードコードしない
- PAT のスコープは `contents: read`（読み取り専用）のみ付与する。書き込みスコープは不要
- ランナーは外部公開エンドポイントを持たない。サーバーから GitHub API への送信のみで動作する

## 5. CI ランナー 実行ポリシー

- blob SHA が前回実行時と同一の場合はビルドをスキップする（差分なしと判断）
- `runner.py` は 1 回実行して終了する oneshot 設計とし、多重起動は systemd タイマーの設定（`OnUnitActiveSec`）で防ぐ
- `pipeline.sh` の終了コードが `0` 以外の場合はビルド失敗としてログに記録する
- SHA ファイルはビルド成功後にのみ更新する。ビルド失敗時は前回 SHA を保持し、次回起動時に再試行する

## 6. CI ランナー ブランチポリシー

- デフォルトでは `main` ブランチの push のみを対象とする
- 対象ブランチは設定値（`BRANCH`）で変更可能

---

# Part 3 — 仕様
> 実装の具体的詳細を定める。「どのように動作・実装するか」に答える。

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
OUT = "/var/www/html/Adlaire-db-spec.html"             # 出力 HTML
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
| `[label](url)` | `<a href="url">label</a>` |
| `[^id]` | `<sup><a href="#fn-id" id="fnref-id" class="fn-ref">[N]</a></sup>`（N は参照順の番号） |

**処理順の注意：** 画像記法（`![alt](url)`）はリンク記法（`[label](url)`）より先にマッチングする。脚注参照（`[^id]`）はリンク置換後に適用する。

**グローバル状態（脚注）：**

| 変数 | 型 | 用途 |
|------|-----|------|
| `_fn_defs` | `dict[str, str]` | 脚注 ID → テキストの対応（`convert()` 呼び出し前に収集） |
| `_fn_order` | `list[str]` | 本文中での参照順脚注 ID リスト（`inline()` が参照時に追記） |

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
  <header id="hdr">        <!-- 固定ヘッダー（高さ 52px、背景 --adlaire-surface-accent） -->
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

---

## 8. 実行方法

```bash
python3 build_spec_v3.py
```

**標準出力：**
```
Converting MD...
Building TOC...
Assembling HTML...
Done → /var/www/html/Adlaire-db-spec.html  (1,713,731 bytes / 1,673 KB)
```

**再実行時の注意：** `_seen`（スラグ重複カウンタ）・`_fn_order`（脚注参照順）・`_fn_defs`（脚注定義）はいずれもモジュールレベル変数であり、スクリプトを起動するたびに初期化される。通常の `python3 build_spec_v3.py` 実行では複数回実行しても出力は同一になる。ただし本スクリプトを `import` して `convert()` を複数回呼ぶ場合は、呼び出し前に `_fn_order.clear()` および `_seen.clear()` を明示的にリセットする必要がある。

---

## 9. 既知の制限

| 制限 | 詳細 |
|------|------|
| インデント付き閉じフェンス | CommonMark では `  ``` ` のようなインデント付き閉じフェンスが有効だが、本実装は `raw.strip() == fence_marker` で完全一致を要求するため非対応。インデント付き閉じフェンスはフェンス内の行として取り込まれる |
| 生 HTML のパススルー | Markdown 中の生 HTML（`<div>`、`<span>` 等）は `esc()` でエスケープされテキストとして出力される。HTML をそのまま通過させる機能はない |

---

## 10. CI ランナー 要件

| 項目 | 内容 |
|------|------|
| Python バージョン | 3.9 以上 |
| 外部依存 | **なし** — `urllib.request`・`base64`・`json`・`subprocess`・`os`・`logging`（すべて標準ライブラリ） |
| 対象 OS | Linux（systemd 対応環境） |
| ネットワーク | サーバーから `api.github.com` への HTTPS 送信のみ |

---

## 11. CI ランナー ファイル構成

### サーバー側

```
/opt/adlaire-builder/
├── runner.py            # CI ランナー本体（単一ファイル、oneshot）
├── build_spec_v3.py     # ビルドスクリプト（サーバー固定）
├── .github_token        # GitHub PAT（パーミッション 600）
└── .last_sha            # 前回取得時の blob SHA キャッシュ

/opt/adlaire-builder/repo/
└── adlaire-db-spec.md   # API 取得後に書き出されるソース Markdown

/var/www/html/           # HTML 出力先（nginx / Apache が配信）

/etc/systemd/system/
├── adlaire-ci.service   # systemd ユニット（oneshot）
└── adlaire-ci.timer     # systemd タイマー（定期実行）
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

```python
TOKEN_FILE   = "/opt/adlaire-builder/.github_token"   # GitHub PAT
OWNER        = "<GitHubオーナー名>"                    # リポジトリオーナー
REPO         = "<リポジトリ名>"                        # リポジトリ名
BRANCH       = "main"                                  # 対象ブランチ
TARGET_FILE  = "adlaire-db-spec.md"                   # 監視対象ファイル
SHA_FILE     = "/opt/adlaire-builder/.last_sha"       # blob SHA キャッシュ
SRC          = "/opt/adlaire-builder/repo/adlaire-db-spec.md"  # 書き出し先
BUILD_SCRIPT = "/opt/adlaire-builder/build_spec_v3.py"
LOG_LEVEL    = "INFO"
```

---

## 13. 処理フロー

```
runner.py 起動（systemd タイマーから呼び出し）
    │
    ├─ .github_token 読み込み（不在の場合は起動失敗）
    │
    ├─ Step 1: Git Trees API
    │   GET /repos/{OWNER}/{REPO}/git/trees/{BRANCH}?recursive=1
    │   → TARGET_FILE の blob SHA を取得
    │   └─ API 失敗時：ERROR ログ、終了
    │
    ├─ SHA_FILE の前回 SHA と比較
    │   └─ 一致（変更なし）→ INFO ログ、正常終了
    │
    ├─ Step 2: Git Blobs API
    │   GET /repos/{OWNER}/{REPO}/git/blobs/{sha}
    │   → Base64 デコード → SRC パスへ書き出し
    │   └─ API 失敗時：ERROR ログ、終了
    │
    ├─ pipeline.sh 実行（bash {SRC の親ディレクトリ}/.ci/pipeline.sh）
    │   ├─ 成功（exit 0）：INFO ログ
    │   └─ 失敗（exit ≠ 0）：ERROR ログ、SHA_FILE 更新せず終了
    │
    └─ SHA_FILE を新 SHA で更新
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
python3 /opt/adlaire-builder/build_spec_v3.py
```

`build_spec_v3.py` はサーバー固定（`/opt/adlaire-builder/`）のため、リポジトリには含めない。

---

## 15. ログ

Python 標準の `logging` モジュールを使用する。出力先は stdout（systemd が journald に転送）。

| レベル | 出力条件 |
|--------|---------|
| `INFO` | 起動、変更なしスキップ、ビルド開始・完了、SHA 更新 |
| `WARNING` | — |
| `ERROR` | トークン読み込み失敗、API 失敗、ビルド失敗 |
| `DEBUG` | API レスポンス詳細等（`LOG_LEVEL = "DEBUG"` 時のみ） |

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
| Webhook 設定 | **不要** |

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

# 4. SHA キャッシュファイルを初期化
sudo -u deploy touch /opt/adlaire-builder/.last_sha

# 5. build_spec_v3.py を配置
sudo cp build_spec_v3.py /opt/adlaire-builder/build_spec_v3.py
sudo chown deploy:deploy /opt/adlaire-builder/build_spec_v3.py

# 6. runner.py を配置
sudo cp runner.py /opt/adlaire-builder/runner.py
sudo chown deploy:deploy /opt/adlaire-builder/runner.py

# 7. systemd ユニットを登録・タイマー起動
sudo cp adlaire-ci.service /etc/systemd/system/
sudo cp adlaire-ci.timer   /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now adlaire-ci.timer
```

---

## 19. 拡張ポイント

| 区分 | 機能 | 概要 |
|------|------|------|
| 拡張 | 複数ファイル監視 | `TARGET_FILE` をリストにし、複数 MD ファイルの変更を検出 |
| 拡張 | ビルドログのファイル保存 | `subprocess` の stdout/stderr をファイルに記録 |
| 拡張 | GitHub Commit Status API | ビルド結果（成功/失敗）を対象コミットに紐付けて GitHub に通知 |
| 拡張 | ビルドタイムアウト | `subprocess` に `timeout` を設定し、長時間ビルドを強制終了 |
| 拡張 | ポーリング間隔の動的変更 | systemd タイマーの `OnUnitActiveSec` を変更して間隔を調整 |
| 将来対応 | Webhook 方式 | 即時反応が必要な場合の代替方式。現行の Git Blobs API ポーリング方式を置き換える。サーバーに Git・nginx（リバースプロキシ）が必要。仕様化は対応時に本ドキュメントへ追記する |

---

## 20. CI ランナー 既知の制限

| 制限 | 詳細 |
|------|------|
| 単一ファイル監視 | 初期実装は `TARGET_FILE` 1 ファイル固定。複数対応は §19 拡張ポイント参照 |
| ポーリング遅延 | 変更検出はタイマー間隔（デフォルト 5 分）に依存する。即時反応は不可 |
| `pipeline.sh` のみ対応 | YAML 形式のパイプライン定義には非対応（シェルスクリプト固定） |
| ネットワーク断時の挙動 | API 失敗時はエラーログのみ記録し、次回タイマー起動まで再試行しない |
