# Adlaire CI — Obsidian 詳細仕様

[`docs/details/obsidian.md`](obsidian.md) は、`obsidian` owner component の詳細本文責務として、Obsidian local vault 連携と local vault 同期の入出力、状態、処理順序、異常系、検証条件だけを定義する。方針、状態語彙、Phase 現在状態、ディレクトリ構成、fixture schema は再定義しない。

| 確認対象 | 正本参照 |
|----------|----------|
| 方針、ポリシー、完了判定 | [`docs/SPEC.md`](../SPEC.md) |
| Phase 14 / Phase 15 の現在状態 | [`docs/ROADMAP.md`](../ROADMAP.md) |
| 詳細入口 | [`docs/DETAIL_INDEX.md` Phase 14 Obsidian Vault 連携参照](../DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry)、[`docs/DETAIL_INDEX.md` Phase 15 Obsidian local vault 同期参照](../DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) |
| fixture 証跡 | [`docs/details/fixture.md` Phase 14 Obsidian Vault 連携証跡](fixture.md#phase-14-obsidian-vault-integration-evidence)、[`docs/details/fixture.md` Phase 15 Obsidian local vault 同期証跡](fixture.md#phase-15-obsidian-local-sync-evidence) |
| 実在 / 未作成 path | [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) |

`obsidian` owner component は `builder`、`statefile`、`security`、`release`、`setup` と連携できる。ただし、vault 読取、wikilink / embed / tag / asset 正規化、sync plan、sync apply、sync rollback の仕様判断は `obsidian` owner component が所有する。`release` と `setup` は Phase 15 の `adlaire-ci-obsidian` 配布連携 collaborator であり、sync plan / apply / rollback の判断を所有してはならない。

<a id="obsidian-phase14-vault-integration-contract"></a>
**Phase 14 Obsidian Vault 連携契約：**

Phase 14 は `adlaire-ci-build --input-mode obsidian-vault` で起動し、Obsidian local vault 内の Markdown note と asset を deterministic に正規化して builder へ渡す。Phase 14 は入力 vault を読み取り専用として扱い、入力 vault 配下に file、directory、log、cache、state、report、temporary file を作成、変更、削除してはならない。

| 入力 | 固定契約 |
|------|----------|
| vault root | `--obsidian-vault <path>` で指定する既存 directory。symlink、存在しない path、file、device、socket、FIFO、permission denied は終了コード `2`。 |
| entry note | `--obsidian-entry <relative-path>` で指定する vault root 相対 Markdown file。先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF を禁止する。拡張子は `.md` だけを許可する。 |
| output root | builder の既存 output option を使用する。Phase 14 は builder に渡す前の中間 root を自動作成し、公開 output は builder が所有する。 |
| filter file | `--obsidian-filter-file <path>` を指定した場合だけ JSON object として読む。未指定時は entry note から到達する note と asset だけを対象にする。YAML、TOML、CSV、XML の filter を禁止する。 |
| strict | `--strict` がある場合、warning 分類の未解決 link、未使用 asset、duplicate tag も失敗にする。 |

`--obsidian-filter-file` は UTF-8 JSON object とし、top-level key は `include_notes`、`exclude_notes`、`include_assets`、`exclude_assets` だけを許可する。各 value は array of string とし、各 string は vault root 相対 path pattern で、先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF、`**` を禁止する。未知 key、重複 key、型不一致、空 string は終了コード `2`、stderr `OBSIDIAN_INVALID_FILTER` とする。

filter pattern は slash 区切りの segment 列とする。通常 segment は exact byte match、segment 全体が `*` の場合だけ同一階層の任意 1 segment に一致する。`*.md`、`note*`、`*note`、character class、brace、escape、`?`、recursive wildcard、正規表現を禁止し、検出した場合は終了コード `2`、stderr `OBSIDIAN_INVALID_FILTER_PATTERN` とする。filter array の順序は意味を持たず、同一 string は 1 件へ重複排除する。適用順は `entry note graph`、`include_notes` / `include_assets` による追加、`exclude_notes` / `exclude_assets` による除外の順に固定し、exclude は include に必ず優先する。filter file によって note が asset pattern へ一致した場合、または asset が note pattern へ一致した場合は推測補正せず、終了コード `2`、stderr `OBSIDIAN_FILTER_TYPE_MISMATCH` とする。

Phase 14 / Phase 15 の vault / project 相対 path は、UTF-8、slash 区切り、先頭 `/` なし、末尾 `/` なし、空 segment なし、`.` / `..` segment なし、backslash なし、NUL / CR / LF なし、drive prefix なし、`~` prefix なし、1 segment 255 bytes 以下、全体 4096 bytes 以下に固定する。root 判定は `EvalSymlinks` による解決後 path ではなく、open 対象の各 segment を `Lstat` し symlink を拒否する no-follow 境界で行う。directory、regular file 以外の device、socket、FIFO は拒否し、regular file の hardlink count が 2 以上の場合は `OBSIDIAN_HARDLINK_FORBIDDEN` とする。

Phase 14 / Phase 15 が出力する JSON は canonical JSON とする。canonical JSON は UTF-8、object key は本仕様で列挙した順序、array は各項目で定義した sort key 順、indent なし、末尾 LF 1 個、未知 key なし、`null` なし、空 array は `[]` とする。実装が Go `map` を使って key 順を暗黙化することを禁止し、struct または明示順 writer で出力する。

Phase 14 は `.obsidian/` 配下を読まない。`.obsidian/` 配下の設定、plugin data、workspace、canvas、sync metadata を、link 解決、tag 抽出、theme、alias、sort、出力に使ってはならない。

| 出力 | 固定契約 |
|------|----------|
| normalized note | 中間 root 配下へ builder 入力用 Markdown として出力する。入力 note の byte を直接公開 output へコピーしない。 |
| asset | 中間 root 配下へ digest 検証付きで copy する。asset path は vault root 相対 path を slash 区切りで保持し、公開 URL 用 path へ builder が変換する。 |
| map | `obsidian_map.json` を中間 root 配下へ出力し、source path、normalized path、resolved link、asset、tag、diagnostic を持つ。 |
| stdout | builder の標準 stdout 契約を維持し、Phase 14 追加 report key は `obsidian_assets`、`obsidian_diagnostics`、`obsidian_notes`、`obsidian_unresolved_links` だけを許可する。 |
| stderr | 失敗時は 1 行目に machine error code、2 行目以降に path と line / column を出力する。secret、absolute vault path、token、home directory は出力しない。 |

Phase 14 machine code は以下に固定する。終了コード `2` は入力、仕様外構文、path safety、filter、CLI contract の不合格、終了コード `1` は中間 root 書込み、report copy、builder handoff の実行失敗とする。終了コード列が `0 / 2` の code は non-strict では diagnostic `warning` として stdout report と `obsidian_map.json` に残し終了コード `0`、strict では同じ code を stderr へ出力し終了コード `2` とする。

| error code | 終了コード | 条件 |
|------------|------------|------|
| `OBSIDIAN_INPUT_MODE_REQUIRED` | `2` | `--obsidian-*` option が指定され、`--input-mode obsidian-vault` がない。 |
| `OBSIDIAN_INVALID_VAULT` | `2` | vault root が既存 directory ではない、または symlink / device / socket / FIFO / permission denied。 |
| `OBSIDIAN_INVALID_ENTRY` | `2` | entry note が path safety または `.md` 拡張子条件に合格しない。 |
| `OBSIDIAN_ENTRY_EXCLUDED` | `2` | filter exclude により entry note が対象外になる。 |
| `OBSIDIAN_INVALID_FILTER` | `2` | filter file の JSON object、key、型、重複 key、空 string が固定契約に合格しない。 |
| `OBSIDIAN_INVALID_FILTER_PATTERN` | `2` | filter pattern が固定 pattern grammar に合格しない。 |
| `OBSIDIAN_FILTER_TYPE_MISMATCH` | `2` | note pattern と asset pattern の対象種別が一致しない。 |
| `OBSIDIAN_PATH_ESCAPE` | `2` | 解決対象 path が vault root または中間 root の no-follow 境界外へ出る。 |
| `OBSIDIAN_HARDLINK_FORBIDDEN` | `2` | regular file の hardlink count が 2 以上。 |
| `OBSIDIAN_YAML_FRONTMATTER_UNSUPPORTED` | `2` | YAML frontmatter を検出した。 |
| `OBSIDIAN_YAML_BLOCK_UNSUPPORTED` | `2` | YAML / YML fenced code block を検出した。 |
| `OBSIDIAN_WIKILINK_MALFORMED` | `2` | wikilink / embed の bracket、target、alias、heading が grammar に合格しない。 |
| `OBSIDIAN_NOTE_EMBED_UNSUPPORTED` | `2` | note embed を検出した。 |
| `OBSIDIAN_AMBIGUOUS_NOTE_LINK` | `2` | basename link が複数 note に一致する。 |
| `OBSIDIAN_UNRESOLVED_LINK` | `0 / 2` | note link または asset embed の target が解決不能。 |
| `OBSIDIAN_DUPLICATE_TAG` | `0 / 2` | 1 note 内で同一 tag が重複した。 |
| `OBSIDIAN_UNUSED_ASSET` | `0 / 2` | filter include または reachable graph 上の asset が出力で参照されない。 |
| `OBSIDIAN_ASSET_UNSUPPORTED` | `2` | asset extension、media type、digest、path safety が不合格。 |
| `OBSIDIAN_REPORT_PATH_INVALID` | `2` | `--obsidian-report-file` が output root 外または path safety 不合格。 |
| `OBSIDIAN_INTERMEDIATE_WRITE_FAILED` | `1` | 中間 root への normalized note、asset、`obsidian_map.json` 書込みに失敗した。 |
| `OBSIDIAN_REPORT_COPY_FAILED` | `1` | `--obsidian-report-file` への copy に失敗した。 |
| `OBSIDIAN_BUILDER_FAILED` | `1` | builder handoff 後に builder owner component が失敗した。 |

<a id="obsidian-phase14-output-schema-contract"></a>
**Phase 14 output schema 契約：**

中間 root は `index.md`、`notes/`、`assets/`、`obsidian_map.json` だけを top-level に持つ。entry note の normalized path は `index.md` に固定し、entry 以外の note は `notes/<source_path>`、asset は `assets/<source_path>` に固定する。`notes/` と `assets/` 配下の path は source path を維持し、拡張子、case、space を変更しない。source path の path safety に不合格な入力を percent encode、rename、推測補正して出力してはならない。

`obsidian_map.json` は root object とし、key 順を `schema_version`、`vault_digest`、`entry_source_path`、`entry_normalized_path`、`notes`、`assets`、`diagnostics` に固定する。`schema_version` は `obsidian-map-v1`、`entry_normalized_path` は `index.md` とする。`vault_digest` は対象 regular file を source path の UTF-8 byte 昇順で並べ、各 file について `source_path` + LF + lowercase SHA-256 + LF + decimal size + LF を連結した byte 列の SHA-256 lowercase hex とする。

`notes` は source path の UTF-8 byte 昇順の array とし、各 note object の key 順は `source_path`、`normalized_path`、`title`、`source_sha256`、`source_trailing_lf`、`outgoing_links`、`asset_embeds`、`tags`、`diagnostics` とする。`title` は最初の ATX heading text、存在しない場合は拡張子を除いた basename、`source_trailing_lf` は入力 byte が LF で終わる場合だけ `true` とする。`tags` は重複排除後の UTF-8 byte 昇順 array とする。

`outgoing_links` は `line`、`column`、`raw` の順に昇順 sort し、各 object の key 順を `raw`、`target`、`target_path`、`heading`、`alias`、`status`、`error_code`、`line`、`column` とする。`status` は `resolved`、`unresolved`、`ambiguous`、`unsupported` のいずれか、`error_code` は成功時 empty string、失敗時は本詳細本文の `OBSIDIAN_*` code とする。`target_path`、`heading`、`alias` は該当なしの場合 empty string とし、key を省略してはならない。

`asset_embeds` と root `assets` は `source_path` の UTF-8 byte 昇順で sort し、各 object の key 順を `raw`、`source_path`、`normalized_path`、`sha256`、`size`、`media_type`、`status`、`error_code`、`line`、`column` とする。`media_type` は `.png=image/png`、`.jpg=image/jpeg`、`.jpeg=image/jpeg`、`.gif=image/gif`、`.webp=image/webp`、`.svg=image/svg+xml`、`.pdf=application/pdf` に固定する。禁止拡張子、vault 外参照、digest 不一致は `status=unsupported` または `status=unresolved` とし、strict でない場合も diagnostic record に残す。SVG は opaque asset として copy だけを行い、inline 展開、XML parse、script 除去、外部参照 fetch、data URI 変換をしてはならない。

`line` は LF 正規化後の 1-based line number、`column` は当該 line 先頭からの 1-based UTF-8 byte offset とする。multibyte 文字を rune 数、display width、grapheme cluster 数で数えてはならない。

`diagnostics` は `source_path`、`line`、`column`、`code` の順に昇順 sort し、各 object の key 順を `code`、`severity`、`source_path`、`line`、`column`、`target` とする。`severity` は `error` または `warning` だけを許可する。free-form message、absolute path、home directory、host user 名、secret を diagnostic に含めてはならない。

<a id="obsidian-phase14-cli-contract"></a>
**Phase 14 CLI 契約：**

| option | 必須 | 値 | 契約 |
|--------|------|----|------|
| `--input-mode` | 必須 | `obsidian-vault` | 他の値は Phase 14 対象外として builder 既存契約へ委譲する。 |
| `--obsidian-vault` | 必須 | path | 既存 directory。symlink なら失敗。 |
| `--obsidian-entry` | 必須 | relative path | vault root 相対 `.md`。 |
| `--obsidian-filter-file` | 任意 | path | JSON filter file。 |
| `--obsidian-report-file` | 任意 | path | `obsidian_map.json` の copy 出力先。output root 外なら失敗。 |
| `--strict` | 任意 | なし | warning を failure へ昇格する。 |

Phase 14 は `--obsidian-*` option が 1 つでも指定され、`--input-mode obsidian-vault` がない場合、終了コード `2`、stderr `OBSIDIAN_INPUT_MODE_REQUIRED` とする。

<a id="obsidian-phase14-normalization-contract"></a>
**Phase 14 Obsidian 正規化契約：**

Markdown file は UTF-8 とし、invalid UTF-8、NUL、HT / LF / CR を除く C0 制御文字は終了コード `2`。CRLF は LF へ正規化して中間 root へ出力する。入力 byte の末尾改行有無は `obsidian_map.json` に `source_trailing_lf` として記録する。

YAML は使用禁止である。file 先頭が `---` + LF または `---` + CRLF で始まる場合は YAML frontmatter と判定し、終了コード `2`、stderr `OBSIDIAN_YAML_FRONTMATTER_UNSUPPORTED` とする。info string が `yaml` または `yml` の fenced code block は終了コード `2`、stderr `OBSIDIAN_YAML_BLOCK_UNSUPPORTED` とする。

| 構文 | 固定契約 |
|------|----------|
| note wikilink | `[[target]]`、`[[target#heading]]`、`[[target|alias]]`、`[[target#heading|alias]]` を許可する。target は `.md` 省略可。 |
| asset embed | `![[asset.ext]]`、`![[asset.ext|alt]]` を許可する。asset extension は `.png`、`.jpg`、`.jpeg`、`.gif`、`.webp`、`.svg`、`.pdf` だけ。 |
| note embed | `![[note.md]]` と `![[note]]` は Phase 14 では未対応とし、終了コード `2`、stderr `OBSIDIAN_NOTE_EMBED_UNSUPPORTED`。 |
| tag | `#` の直前が start、space、tab、`(`、`[`、`{` のいずれかで、本文が `[A-Za-z0-9][A-Za-z0-9_/-]*` の場合だけ tag とする。code fence、inline code、URL fragment 内の `#` は tag ではない。 |
| heading link | `#heading` は解決先 note の builder heading slug と完全一致する場合だけ成功する。 |
| duplicate basename | target に `/` がなく、同じ basename を持つ note が複数ある場合は推測解決せず `OBSIDIAN_AMBIGUOUS_NOTE_LINK`。 |
| case mismatch | file system が case-insensitive でも path は byte の case-sensitive exact match とする。 |
| unsupported syntax | `.canvas`、Dataview code block、Templater marker、external fetch directive は失敗にする。 |

wikilink / embed parser は、code fence と inline code span を先に保護し、保護範囲内の `[[`、`![[`、`#` を通常 text として扱う。保護範囲外で `[[` または `![[` の直前 byte が backslash の場合は token 開始ではなく、backslash と bracket をそのまま通常 text として残す。保護範囲外で unescaped `[[` または `![[` を見つけた場合、同一 note 内の次の `]]` までを 1 token とし、その間に別の `[[` または `![[` が出現する場合は `OBSIDIAN_WIKILINK_MALFORMED` とする。target は trim 後に空であってはならず、target、heading、alias の各 field に NUL、CR、LF、backslash、`[`、`]` を含めてはならない。alias は output Markdown の link text へ使うだけで、target 解決、sort、digest、normalized path に影響してはならない。

Phase 14 の link 解決は、entry note から到達可能な note graph だけを対象にする。filter file の include は到達対象を増やせるが、exclude が entry note を除外する場合は終了コード `2`、stderr `OBSIDIAN_ENTRY_EXCLUDED`。

normalized Markdown への変換は次に固定する。変換後の Markdown は UTF-8、LF、末尾 LF は入力 `source_trailing_lf` と一致させる。変換後に builder へ渡す link destination は normalized note file から見た slash 区切りの相対 path とし、`filepath.Rel` 相当の結果を `/` 区切りへ固定し、空文字、backslash、absolute path、`..` で中間 root 外へ出る path を禁止する。

| 入力 token | resolved 時の normalized Markdown | unresolved / warning 時の normalized Markdown | failure 条件 |
|------------|------------------------------------|----------------------------------------------|--------------|
| note wikilink | `[display](relative-target.md)` または `[display](relative-target.md#slug)` に置換する。`display` は alias が非空なら alias、alias が空なら heading が非空なら heading、heading が空なら target の basename without extension とする。 | non-strict では raw token を byte 単位で保持し、`obsidian_map.json` diagnostics に記録する。strict では stderr `OBSIDIAN_UNRESOLVED_LINK`。 | target path safety 不合格、duplicate basename、heading slug 不一致、malformed token。 |
| asset image embed | `![display](relative-asset-path)` に置換する。`.png`、`.jpg`、`.jpeg`、`.gif`、`.webp`、`.svg` を image embed とする。`display` は alias が非空なら alias、alias が空なら asset basename とする。 | non-strict では raw token を byte 単位で保持し、`obsidian_map.json` diagnostics に記録する。strict では stderr `OBSIDIAN_UNRESOLVED_LINK`。 | asset path safety 不合格、許可外 extension、digest 読取失敗、vault 外参照。 |
| asset PDF embed | `[display](relative-asset-path)` に置換する。`display` は alias が非空なら alias、alias が空なら asset basename とする。 | non-strict では raw token を byte 単位で保持し、`obsidian_map.json` diagnostics に記録する。strict では stderr `OBSIDIAN_UNRESOLVED_LINK`。 | asset path safety 不合格、digest 読取失敗、vault 外参照。 |
| tag | Markdown 本文 byte は変更しない。抽出結果だけを `obsidian_map.json` の `tags` へ記録する。 | duplicate tag は non-strict で warning、strict で stderr `OBSIDIAN_DUPLICATE_TAG`。 | tag 本文が固定 grammar に合格しない場合は通常 text。 |

`relative-target.md` は target note の normalized path を source note の normalized path の parent から見た相対 path とする。entry note への link は `index.md`、entry 以外への link は `notes/<source_path>` を基準に計算する。heading link の `slug` は [`docs/details/builder.md` 詳細本文責務 §28 ID / slug / search index / JS state 決定性固定契約](builder.md#sec-28-common-determinism) と同じ slug base と duplicate slug 規則で target note 内の見出しから確定した値を使用する。Obsidian heading text を URL encode、percent decode、case fold、Unicode 正規化してはならない。

Markdown link text と image alt text に使う `display` は `[`、`]`、NUL、CR、LF を含んではならない。`display` がこの条件に合格しない場合は `OBSIDIAN_WIKILINK_MALFORMED` とする。変換処理は code fence、inline code span、YAML rejection、wikilink / embed 解析、tag 抽出、heading slug 作成、link / embed 置換の順に実行する。置換は source note 内の token 開始 byte offset 降順で行い、先に置換した token が後続 token の line / column、diagnostic、map sort に影響してはならない。

<a id="obsidian-phase14-builder-handoff-contract"></a>
**Phase 14 builder handoff 契約：**

Phase 14 は正規化済み中間 root を builder owner component へ渡す。builder 呼び出しは process 再起動ではなく owner component 境界呼び出しで行い、builder の Markdown parser、site output、theme、search index、report 契約を変更しない。

| handoff item | 契約 |
|--------------|------|
| source root | normalized note と asset だけを含む。vault metadata、`.obsidian/`、filter file、temporary file は含めない。 |
| entry file | `index.md` に固定する。entry note の source path は `obsidian_map.json` に記録する。 |
| site mode | builder の静的 Web サイト出力を使う。single HTML 専用出力を Phase 14 の標準にしない。 |
| failure | builder が失敗した場合、Phase 14 は入力 vault を変更せず、builder の終了コードと stderr を保持する。 |

<a id="obsidian-phase15-local-sync-contract"></a>
**Phase 15 Obsidian local vault 同期契約：**

Phase 15 は `adlaire-ci-obsidian sync plan`、`adlaire-ci-obsidian sync apply`、`adlaire-ci-obsidian sync rollback` で起動する。`sync plan` は project root、vault root、`sync_state.json` を変更せず、`--plan-file` への plan JSON atomic write だけを行う。`sync apply` は `sync plan` が生成した plan file と plan hash が一致する場合だけ書込みを行う。`sync rollback` は apply が生成した rollback record に基づいて、最後に成功または部分失敗した apply の影響を戻す。

| 同期対象 | 契約 |
|----------|------|
| project root | Adlaire CI project 側の Markdown / asset tree。path safety は Phase 14 と同じ。 |
| vault root | Obsidian local vault。symlink root、permission denied、device、socket、FIFO は失敗。 |
| state dir | statefile owner component の lock / atomic write 契約で更新する。 |
| direction | `import-only`、`export-only`、`bidirectional` だけを許可する。 |
| delete policy | `reject`、`tombstone` だけを許可する。hard delete は Phase 15 で禁止する。 |

`--conflict-dir` と `--tombstone-dir` は project root 相対 path とし、Phase 14 / Phase 15 共通 path safety に合格しなければならない。未指定時の既定値は `adlaire-ci-conflicts/` と `adlaire-ci-tombstones/` とする。指定値が vault root、project root 外、absolute path、symlink、hardlink、device、socket、FIFO、既存 regular file のいずれかに該当する場合は終了コード `2`、stderr `OBSIDIAN_SYNC_PATH_INVALID` とする。

`--plan-file` と `--rollback-file` は state dir 相対 path とし、Phase 14 / Phase 15 共通 path safety に合格しなければならない。absolute path、空 path、`.`、`..`、backslash、NUL、CR、LF、state dir 外、symlink、hardlink、device、socket、FIFO、既存 directory は終了コード `2`、stderr `OBSIDIAN_SYNC_PATH_INVALID` とする。`sync plan` の `--plan-file` は必須であり、既存 regular file がある場合は atomic rewrite する。`sync apply` の `--rollback-file` が未指定の場合は `rollback/<plan_hash>.json` を既定値とし、`rollback/` は state dir 配下に mode `0700` で作成する。`sync rollback` の `--rollback-file` は必須とし、既定値推測を行わない。

<a id="obsidian-phase15-cli-contract"></a>
**Phase 15 CLI 契約：**

| command | 必須 option | 任意 option | 書込み |
|---------|-------------|-------------|--------|
| `sync plan` | `--vault`、`--project-root`、`--state-dir`、`--plan-file`、`--direction` | `--delete-policy`、`--conflict-dir`、`--tombstone-dir` | `--plan-file` のみ |
| `sync apply` | `--vault`、`--project-root`、`--state-dir`、`--plan-file`、`--plan-hash` | `--rollback-file`、`--open-uri` | あり |
| `sync rollback` | `--vault`、`--project-root`、`--state-dir`、`--rollback-file` | `--open-uri` | あり |

`--delete-policy` の既定値は `reject`、`--conflict-dir` の既定値は `adlaire-ci-conflicts/`、`--tombstone-dir` の既定値は `adlaire-ci-tombstones/`、`--open-uri` の既定値は `false` とする。`sync apply` は conflict dir と tombstone dir を plan file からだけ読み、CLI option で再指定してはならない。`--open-uri` は `true` または `false` だけを許可する。`true` の場合でも、Obsidian URI の起動成功を sync 成功条件にしてはならない。URI 起動失敗は warning として report に記録し、apply / rollback の filesystem 結果を覆さない。

Phase 15 CLI は、未知 command、未知 option、重複 option、必須 option 欠落、`--direction` の許可値外、`--delete-policy` の許可値外、`--open-uri` の `true` / `false` 以外を終了コード `2`、stderr `OBSIDIAN_SYNC_INVALID_OPTION` とする。path safety に到達した option 値の不合格は `OBSIDIAN_SYNC_PATH_INVALID` とし、enum / boolean / command / option の不合格と混同してはならない。

Phase 15 の成功時 stdout は canonical JSON object 1 行と LF だけとする。`sync plan` は key 順を `command`、`plan_file`、`plan_hash`、`operations`、`conflicts`、`tombstones`、`applied` に固定し、`applied=false` とする。`sync apply` は key 順を `command`、`plan_file`、`plan_hash`、`operations_applied`、`conflicts`、`tombstones`、`rollback_file`、`state_digest` に固定する。`sync rollback` は key 順を `command`、`rollback_file`、`operations_rolled_back`、`conflicts`、`state_digest` に固定する。stdout の `plan_file` と `rollback_file` は state dir 相対 path とし、absolute path、home directory、temporary directory、host 固有 path を出力してはならない。失敗時 stdout は 0 byte、stderr は `obsidian: <error-code>` + LF の 1 行だけとし、path、digest、Go error、stack trace、URI、absolute path を出力してはならない。

<a id="obsidian-phase15-schema-contract"></a>
**Phase 15 schema 契約：**

`sync_state.json` は root object とし、key 順を `schema_version`、`project_root_digest`、`vault_root_digest`、`entries`、`tombstones`、`conflicts`、`last_apply_id` に固定する。`schema_version` は `obsidian-sync-state-v1` とする。`entries` は `path` の UTF-8 byte 昇順 array とし、各 object の key 順を `path`、`project_digest`、`vault_digest`、`last_sync_digest`、`last_sync_unix` とする。存在しない側の digest は empty string とし、mtime だけで一致扱いにしてはならない。`tombstones` は `tombstone_id` 昇順、`conflicts` は `conflict_id` 昇順とする。

`project_root_digest` と `vault_root_digest` は同期対象 regular file を normalized path の UTF-8 byte 昇順で並べ、各 file について `normalized_path` + LF + lowercase SHA-256 + LF + decimal size + LF を連結した byte 列の SHA-256 lowercase hex とする。`state_digest` は `sync_state.json` の canonical JSON byte 列の SHA-256 lowercase hex とする。state file が未作成の場合、実装は empty state object を canonical JSON として生成し、その byte 列から `state_digest` を算出する。empty state object の `entries`、`tombstones`、`conflicts` は `[]`、`last_apply_id` は empty string とし、`project_root_digest` と `vault_root_digest` は現在 tree digest を入れる。

`plan.json` は root object とし、key 順を `schema_version`、`plan_id`、`created_at_unix`、`direction`、`delete_policy`、`conflict_dir`、`tombstone_dir`、`project_root_digest`、`vault_root_digest`、`state_digest`、`operations`、`conflicts`、`tombstones` に固定する。`schema_version` は `obsidian-sync-plan-v1` とする。`conflict_dir` と `tombstone_dir` は project root 相対 path とし、plan 作成時に正規化した値だけを格納する。`plan_id` は `direction` + LF + `delete_policy` + LF + `conflict_dir` + LF + `tombstone_dir` + LF + `project_root_digest` + LF + `vault_root_digest` + LF + `state_digest` + LF の SHA-256 lowercase hex とする。`plan_hash` は `plan.json` canonical JSON byte 列の SHA-256 lowercase hex とし、`plan.json` 内には格納しない。

`operations` は `path` の UTF-8 byte 昇順、同一 path 内は `conflict`、`tombstone`、`create`、`update`、`noop` の順に sort する。各 object の key 順は `op`、`direction`、`path`、`source`、`destination`、`before_digest`、`after_digest`、`conflict_id`、`tombstone_id`、`rollback_required` とする。`op` は `conflict`、`tombstone`、`create`、`update`、`noop`、`direction` は `project-to-vault`、`vault-to-project`、`none` のいずれかとする。該当しない string field は empty string、`rollback_required` は boolean とし、key を省略してはならない。

operation array は、project tree、vault tree、sync state entries、tombstones、conflicts の path 和集合に対して 1 path 以上 1 operation 以下を生成する。変更がない path は `noop` として残し、plan 作成時の対象集合を audit 可能にする。`import-only` は vault 側を source、project 側を destination とし、project だけが state から変更されている path は `conflict` / `opposite-side-edit` とする。`export-only` は project 側を source、vault 側を destination とし、vault だけが state から変更されている path は `conflict` / `opposite-side-edit` とする。`bidirectional` は片側だけが state から変更された path を変更側から反対側への `create` または `update` とし、両側が state から変更された path は `conflict` / `both-side-edit` とする。削除は `delete_policy=reject` では `conflict` / `delete-vs-edit`、`delete_policy=tombstone` では `tombstone` とし、hard delete operation を生成してはならない。

`conflicts` は `conflict_id` 昇順 array とし、各 object の key 順を `conflict_id`、`path`、`reason`、`project_digest`、`vault_digest`、`state_digest`、`resolution` に固定する。`reason` は `both-side-edit`、`delete-vs-edit`、`rename-collision`、`read-only-target`、`digest-changed`、`opposite-side-edit` のいずれか、`resolution` は Phase 15 では常に `manual` とする。自動 merge、last-writer-wins、mtime 優先を禁止する。

`tombstones` は `tombstone_id` 昇順 array とし、各 object の key 順を `tombstone_id`、`path`、`digest`、`direction`、`created_at_unix`、`tombstone_path` に固定する。`tombstone_path` は tombstone dir 相対 path とし、vault root または project root の絶対 path を含めてはならない。

`rollback record` は root object とし、key 順を `schema_version`、`apply_id`、`plan_hash`、`started_at_unix`、`completed_at_unix`、`operations`、`state_before_digest`、`state_after_digest` に固定する。`schema_version` は `obsidian-sync-rollback-v1` とする。`operations` は apply 実行順で記録し、rollback 実行時は逆順で処理する。各 operation record の key 順は `op`、`path`、`destination`、`previous_digest`、`new_digest`、`backup_path`、`status`、`error_code` とし、`status` は `applied`、`skipped`、`failed` のいずれかとする。destination write 開始前の rollback record 初期版では、`completed_at_unix=0`、`state_after_digest=""`、全 operation `status=skipped`、`error_code=""` とする。

`apply_id` は `plan_hash` + LF + `state_before_digest` + LF + decimal `started_at_unix` + LF の SHA-256 lowercase hex とする。`conflict_id` は `path` + LF + `reason` + LF + `project_digest` + LF + `vault_digest` + LF + `state_digest` + LF の SHA-256 lowercase hex とする。`tombstone_id` は `path` + LF + `digest` + LF + `direction` + LF + decimal `created_at_unix` + LF の SHA-256 lowercase hex とする。ID 算出に absolute path、mtime、process id、random、map iteration order を含めてはならない。

rollback backup は state dir 配下の `rollback-backups/<apply_id>/<ordinal>-<path_sha256>.bak` に保存する。`ordinal` は apply 実行順の 6 桁 decimal zero padding、`path_sha256` は operation `path` の UTF-8 byte 列の SHA-256 lowercase hex とする。`backup_path` は state dir 相対 path を記録し、absolute path を記録してはならない。destination が apply 前に存在しない create operation では apply 時の record を `previous_digest=""`、`backup_path=""` とし、rollback 時は destination の current digest が `new_digest` と一致する場合だけ、destination を state dir 配下の `rollback-created/<apply_id>/<ordinal>-<path_sha256>.bak` へ atomic rename してから state を戻す。rollback-created path は rollback record の当該 operation `backup_path` へ atomic rewrite で追記する。Phase 15 rollback は destination content を unlink で破棄してはならない。

Phase 15 の `created_at_unix`、`started_at_unix`、`completed_at_unix`、`last_sync_unix` は UTC Unix seconds の decimal integer とする。fixture では fake clock で固定し、実装は現在時刻を直接参照する箇所を sync plan / apply / rollback の時刻 provider 1 箇所へ集約する。各 timestamp は処理開始時に 1 回だけ取得し、同一 report 内で同じ意味の timestamp を複数回取得して揺らしてはならない。

<a id="obsidian-phase15-apply-rollback-contract"></a>
**Phase 15 apply / rollback 契約：**

`sync plan` は project tree、vault tree、sync state を読み、operation array を `--plan-file` へ JSON で出力する。plan hash は canonical JSON byte 列の SHA-256 lowercase hex とする。operation order は path の UTF-8 byte 昇順、同一 path 内は `conflict`、`tombstone`、`create`、`update`、`noop` の順とする。

`sync apply` は state dir の process lock を取得し、plan file を再読込し、plan hash を照合し、plan file 内の `conflict_dir` と `tombstone_dir` を含む全 path field を再検証し、project tree / vault tree / sync state の digest が plan 作成時と一致することを確認してから staging へ書く。書込み順序は、rollback record 初期版を staging へ書く、file fsync、parent directory fsync、rollback record を `--rollback-file` へ atomic rename、各 destination の sibling staging へ新 content または tombstone content を書く、file fsync、parent directory fsync、destination atomic rename、operation record status 更新、rollback record atomic rewrite、`sync_state.json` staging 書込み、file fsync、parent directory fsync、state atomic rename、state parent directory fsync、rollback record の `completed_at_unix` と `state_after_digest` 更新の順に固定する。途中失敗では、成功済み rename と state update を rollback record に記録し、終了コード `1` とする。rollback record を作成できない場合は destination への write を開始してはならない。

`sync rollback` は rollback record の operation を逆順に処理する。rollback record の対象 path が現在 digest と一致しない場合は上書きせず conflict とし、終了コード `1`、stderr `OBSIDIAN_ROLLBACK_CONFLICT` とする。

rollback 実行時は operation `status` だけを信用してはならない。対象 destination の現在 digest が rollback record の `new_digest` と一致し、`previous_digest` または `backup_path` が存在する場合は、`status=skipped` のままでも crash recovery 対象として rollback を実行する。現在 digest が `previous_digest` と一致する operation は rollback 済みとして `skipped` 扱いにし、現在 digest が `previous_digest` とも `new_digest` とも一致しない operation は conflict とする。

<a id="obsidian-phase15-conflict-tombstone-contract"></a>
**Phase 15 conflict / tombstone 契約：**

| 状況 | 契約 |
|------|------|
| both-side edit | project と vault の両方が前回 digest から変更された場合、書込みせず conflict record を作る。 |
| delete vs edit | 片側削除、片側変更は conflict。delete policy が `tombstone` でも自動上書きしない。 |
| rename collision | 同じ destination path に複数 source が解決される場合 conflict。 |
| clock skew | mtime は参考情報だけとし、digest 不一致だけで変更判定する。 |
| tombstone | delete policy `tombstone` の削除は tombstone dir へ移動し、元 path、digest、timestamp、direction を state に記録する。 |
| hard delete | Phase 15 では禁止。hard delete option は存在しない。 |

<a id="obsidian-phase15-external-boundary-contract"></a>
**Phase 15 外部境界契約：**

Phase 15 は Obsidian Sync service、Obsidian cloud、remote vault API、plugin runtime、community plugin、external watcher、external diff library を使用しない。network access は行わない。Obsidian URI は apply / rollback 成功後に対象 vault または file を開く任意補助だけに使用できる。URI は `obsidian://open` だけを許可し、query は percent-encoding し、vault 名と file path 以外を含めない。

Phase 15 は credentials を扱わない。token、password、Obsidian account、Sync encryption password、remote endpoint、cloud credential を CLI、state、report、log、fixture に保存してはならない。

<a id="obsidian-phase15-distribution-contract"></a>
**Phase 15 配布連携契約：**

Phase 15 の実装完了には `adlaire-ci-obsidian` 実行バイナリの配布契約が必要である。Phase 15 完了状態では、[`docs/details/release.md` 詳細本文責務 Phase 15 Obsidian Release 配布拡張契約](release.md#phase-15-obsidian-release-extension-contract)、[`docs/details/setup.md` 詳細本文責務 §26.2a](setup.md#sec-26-2a)、[`docs/details/setup.md` 詳細本文責務 Phase 15 Obsidian CLI 導入手順](setup.md#phase-15-obsidian-setup-contract) が整合し、Release asset 生成、checksum、setup 取得対象、`install-obsidian` mode、version 出力、fixture 証跡に `adlaire-ci-obsidian-linux-amd64` が含まれる。この配布連携が未完了の場合、sync plan / apply / rollback の実装が合格していても Phase 15 を `実装済み` に遷移してはならない。

| error code | 終了コード | 条件 |
|------------|------------|------|
| `OBSIDIAN_SYNC_INVALID_OPTION` | `2` | Phase 15 CLI の command、option、必須 option、enum、boolean が固定契約に合格しない。 |
| `OBSIDIAN_SYNC_PLAN_HASH_MISMATCH` | `2` | apply の `--plan-hash` が plan file と一致しない。 |
| `OBSIDIAN_SYNC_STATE_CHANGED` | `1` | plan 作成後に project、vault、state の digest が変化した。 |
| `OBSIDIAN_SYNC_CONFLICT` | `1` | conflict が 1 件以上ある。 |
| `OBSIDIAN_SYNC_PATH_INVALID` | `2` | `--conflict-dir`、`--tombstone-dir`、destination、backup、または rollback path が path safety に合格しない。 |
| `OBSIDIAN_SYNC_ATOMIC_WRITE_FAILED` | `1` | staging、fsync、rename、state update、rollback record のいずれかに失敗した。 |
| `OBSIDIAN_ROLLBACK_CONFLICT` | `1` | rollback 対象 path の現在 digest が rollback record の `previous_digest` / `new_digest` と一致しない。 |
| `OBSIDIAN_SYNC_SERVICE_DEPENDENCY_FORBIDDEN` | `2` | Obsidian Sync service、cloud、remote API、plugin runtime への依存を検出した。 |
