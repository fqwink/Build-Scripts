# Adlaire CI — Obsidian 詳細仕様

[`docs/details/obsidian.md`](obsidian.md) は、`obsidian` owner component の詳細本文責務として、Obsidian local vault 連携と local vault 同期の入出力、状態、処理順序、異常系、検証条件だけを定義する。方針、状態語彙、Phase 現在状態、ディレクトリ構成、fixture schema は再定義しない。

| 確認対象 | 正本参照 |
|----------|----------|
| 方針、ポリシー、完了判定 | [`docs/SPEC.md`](../SPEC.md) |
| Phase 14 / Phase 15 の現在状態 | [`docs/ROADMAP.md`](../ROADMAP.md) |
| 詳細入口 | [`docs/DETAIL_INDEX.md` Phase 14 Obsidian Vault 連携参照](../DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry)、[`docs/DETAIL_INDEX.md` Phase 15 Obsidian local vault 同期参照](../DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry) |
| fixture 証跡 | [`docs/details/fixture.md` Phase 14 Obsidian Vault 連携証跡](fixture.md#phase-14-obsidian-vault-integration-evidence)、[`docs/details/fixture.md` Phase 15 Obsidian local vault 同期証跡](fixture.md#phase-15-obsidian-local-sync-evidence) |
| 実在 / 未作成 path | [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) |

`obsidian` owner component は `builder`、`statefile`、`security` と連携できる。ただし、vault 読取、wikilink / embed / tag / asset 正規化、sync plan、sync apply、sync rollback の仕様判断は `obsidian` owner component が所有する。

<a id="obsidian-phase14-vault-integration-contract"></a>
**Phase 14 Obsidian Vault 連携契約：**

Phase 14 は `adlaire-ci-build --input-mode obsidian-vault` で起動し、Obsidian local vault 内の Markdown note と asset を deterministic に正規化して builder へ渡す。Phase 14 は入力 vault を読み取り専用として扱い、入力 vault 配下に file、directory、log、cache、state、report、temporary file を作成、変更、削除してはならない。

| 入力 | 固定契約 |
|------|----------|
| vault root | `--obsidian-vault <path>` で指定する既存 directory。symlink、存在しない path、file、device、socket、FIFO、permission denied は終了コード `2`。 |
| entry note | `--obsidian-entry <relative-path>` で指定する vault root 相対 Markdown file。先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF を禁止する。拡張子は `.md` だけを許可する。 |
| output root | builder の既存 output option を使用する。Phase 14 は builder に渡す前の中間 root を自動作成し、公開 output は builder が所有する。 |
| filter file | `--obsidian-filter-file <path>` を指定した場合だけ JSON object として読む。未指定時は entry note から到達する note と asset だけを対象にする。YAML、TOML、CSV、XML の filter を禁止する。 |
| strict | `--strict` がある場合、warning 分類の未解決 link、unsupported syntax、未使用 asset、duplicate tag も失敗にする。 |

`--obsidian-filter-file` は UTF-8 JSON object とし、top-level key は `include_notes`、`exclude_notes`、`include_assets`、`exclude_assets` だけを許可する。各 value は array of string とし、各 string は vault root 相対 path pattern で、先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF、`**` を禁止する。未知 key、重複 key、型不一致、空 string は終了コード `2`。

Phase 14 は `.obsidian/` 配下を読まない。`.obsidian/` 配下の設定、plugin data、workspace、canvas、sync metadata を、link 解決、tag 抽出、theme、alias、sort、出力に使ってはならない。

| 出力 | 固定契約 |
|------|----------|
| normalized note | 中間 root 配下へ builder 入力用 Markdown として出力する。入力 note の byte を直接公開 output へコピーしない。 |
| asset | 中間 root 配下へ digest 検証付きで copy する。asset path は vault root 相対 path を slash 区切りで保持し、公開 URL 用 path へ builder が変換する。 |
| map | `obsidian_map.json` を中間 root 配下へ出力し、source path、normalized path、resolved link、asset、tag、diagnostic を持つ。 |
| stdout | builder の標準 stdout 契約を維持し、Phase 14 追加 report key は `obsidian_assets`、`obsidian_diagnostics`、`obsidian_notes`、`obsidian_unresolved_links` だけを許可する。 |
| stderr | 失敗時は 1 行目に machine error code、2 行目以降に path と line / column を出力する。secret、absolute vault path、token、home directory は出力しない。 |

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

Markdown file は UTF-8 とし、invalid UTF-8、NUL、C0 制御文字は終了コード `2`。CRLF は LF へ正規化して中間 root へ出力する。入力 byte の末尾改行有無は `obsidian_map.json` に `source_trailing_lf` として記録する。

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

Phase 14 の link 解決は、entry note から到達可能な note graph だけを対象にする。filter file の include は到達対象を増やせるが、exclude が entry note を除外する場合は終了コード `2`、stderr `OBSIDIAN_ENTRY_EXCLUDED`。

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

Phase 15 は `adlaire-ci-obsidian sync plan`、`adlaire-ci-obsidian sync apply`、`adlaire-ci-obsidian sync rollback` で起動する。`sync plan` は書込みを行わない。`sync apply` は `sync plan` が生成した plan file と plan hash が一致する場合だけ書込みを行う。`sync rollback` は apply が生成した rollback record に基づいて、最後に成功または部分失敗した apply の影響を戻す。

| 同期対象 | 契約 |
|----------|------|
| project root | Adlaire CI project 側の Markdown / asset tree。path safety は Phase 14 と同じ。 |
| vault root | Obsidian local vault。symlink root、permission denied、device、socket、FIFO は失敗。 |
| state dir | statefile owner component の lock / atomic write 契約で更新する。 |
| direction | `import-only`、`export-only`、`bidirectional` だけを許可する。 |
| delete policy | `reject`、`tombstone` だけを許可する。hard delete は Phase 15 で禁止する。 |

<a id="obsidian-phase15-cli-contract"></a>
**Phase 15 CLI 契約：**

| command | 必須 option | 任意 option | 書込み |
|---------|-------------|-------------|--------|
| `sync plan` | `--vault`、`--project-root`、`--state-dir`、`--plan-file`、`--direction` | `--delete-policy`、`--conflict-dir`、`--tombstone-dir` | なし |
| `sync apply` | `--vault`、`--project-root`、`--state-dir`、`--plan-file`、`--plan-hash` | `--conflict-dir`、`--tombstone-dir`、`--rollback-file`、`--open-uri` | あり |
| `sync rollback` | `--vault`、`--project-root`、`--state-dir`、`--rollback-file` | `--open-uri` | あり |

`--open-uri` は `true` または `false` だけを許可する。`true` の場合でも、Obsidian URI の起動成功を sync 成功条件にしてはならない。URI 起動失敗は warning として report に記録し、apply / rollback の filesystem 結果を覆さない。

<a id="obsidian-phase15-apply-rollback-contract"></a>
**Phase 15 apply / rollback 契約：**

`sync plan` は project tree、vault tree、sync state を読み、operation array を `--plan-file` へ JSON で出力する。plan hash は canonical JSON byte 列の SHA-256 lowercase hex とする。operation order は path の UTF-8 byte 昇順、同一 path 内は `conflict`、`tombstone`、`create`、`update`、`noop` の順とする。

`sync apply` は state dir の process lock を取得し、plan file を再読込し、plan hash を照合し、project tree / vault tree / sync state の digest が plan 作成時と一致することを確認してから staging へ書く。staging 書込み、file fsync、parent directory fsync、atomic rename、state update、rollback record 書込みの順序を固定する。途中失敗では、成功済み rename と state update を rollback record に記録し、終了コード `1` とする。

`sync rollback` は rollback record の operation を逆順に処理する。rollback record の対象 path が現在 digest と一致しない場合は上書きせず conflict とし、終了コード `1`、stderr `OBSIDIAN_ROLLBACK_CONFLICT` とする。

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

| error code | 終了コード | 条件 |
|------------|------------|------|
| `OBSIDIAN_SYNC_PLAN_HASH_MISMATCH` | `2` | apply の `--plan-hash` が plan file と一致しない。 |
| `OBSIDIAN_SYNC_STATE_CHANGED` | `1` | plan 作成後に project、vault、state の digest が変化した。 |
| `OBSIDIAN_SYNC_CONFLICT` | `1` | conflict が 1 件以上ある。 |
| `OBSIDIAN_SYNC_ATOMIC_WRITE_FAILED` | `1` | staging、fsync、rename、state update、rollback record のいずれかに失敗した。 |
| `OBSIDIAN_SYNC_SERVICE_DEPENDENCY_FORBIDDEN` | `2` | Obsidian Sync service、cloud、remote API、plugin runtime への依存を検出した。 |
