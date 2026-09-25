# Adlaire CI — Archive 詳細仕様

[`docs/details/archive.md`](archive.md) は `archive` owner component の詳細本文責務として、`archive` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`archive` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、表示、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `archive` |
| collaborator component | 機能ごとの接続境界と担当処理だけを定義し、ファイル全体の collaborator 一覧は定義しない。 |
| 持つ内容 | `archive` owner が主本文として定義する build log archive / cleanup の実体処理、snapshot 保存形式、保存済み tar.gz の検証・配信、snapshot 世代削除、snapshot delete 実体処理、rollback 用 artifact の展開・転送・temporary cleanup 実体処理。 |
| 持たない内容 | runner の通常 build 実行、snapshot 作成トリガー判定、rollback build の lock、ID 採番、log / history / status / pending 書込、API 共通 request / response、`.config_log`、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 手順、fixture 証跡責務。 |

archive owner は、保存済み build log と snapshot artifact の圧縮、展開、列挙、削除、転送の実体処理と、呼び出し元が状態を確定できる処理結果の返却だけを担当する。API 境界は [`docs/details/api.md`](api.md)、SDK 境界は [`docs/details/sdk.md`](sdk.md)、UI 境界は [`docs/details/ui.md`](ui.md)、runner 境界は [`docs/details/runner.md`](runner.md) の各詳細本文責務を参照する。

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [§27.7](archive.md#sec-27-7) | ビルドログのアーカイブ圧縮。 |
| [§27.15](archive.md#sec-27-15) | ビルドアーティファクト管理。 |

---

<a id="sec-27-7"></a>
**27.7 ビルドログのアーカイブ圧縮：**

owner component は `archive` とする。collaborator component は `runner`、`api`、`statefile` とする。

[§27.7](archive.md#sec-27-7) で archive / cleanup API の HTTP endpoint、query、response body、HTTP response body への変換を述べる場合は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を共通参照先とする。[§27.7](archive.md#sec-27-7) では archive log の探索、展開、圧縮、削除、処理結果だけを定義する。

archive owner は、runner または `POST /api/logs/archive` から呼び出された場合に、`.server_config.log_archive_after_days > 0` で対象日数より古い `.build_logs/{id}.json` を gzip 圧縮し、`.build_logs/archive/{id}.json.gz` へ保存する。圧縮成功後、元の `.build_logs/{id}.json` を削除する。`.build_logs/archive/` 内のファイルを再圧縮してはならない。

gzip は Go 標準ライブラリ `compress/gzip` を使用し、mtime は元ファイル mtime ではなく圧縮実行時刻でよい。圧縮前 JSON を読み込めないファイルは archive 対象外とし、WARN `LOG_ARCHIVE_SKIP_CORRUPT: id=<id>` を出す。実行中 build の `current_build_id` と一致する log は対象外とする。

archive owner は、通常 log が存在しない場合に archive log を gzip 展開し、通常 `.build_logs/{id}.json` と同じ schema の JSON object として呼び出し元へ渡す。

archive owner は、`POST /api/logs/cleanup` から呼び出された場合に、archive 済みファイルも `log_retention_days` の削除対象に含める。archive owner の処理結果は `archived_count`、`deleted_count`、`failed_count` を持つ。

**archive / cleanup 固定契約：**

| 項目 | 仕様 |
|------|------|
| archive 対象判定 | build log JSON の `finished_at` を基準にする。欠落時は file mtime を使わず対象外。 |
| archive id | file 名 `{id}.json` の id と JSON 内 `id` が一致する場合だけ対象。 |
| gzip path | `.build_logs/archive/{id}.json.gz`。既存 archive がある場合は上書きせず skip する。 |
| cleanup 順 | retention 対象の通常 log を basename の ASCII 昇順で各 1 回削除 → retention 対象の archive log を basename の ASCII 昇順で各 1 回削除 → archive directory を 1 回再読込 → entry が 0 件の場合だけ archive directory を `os.Remove` で 1 回削除。 |
| 削除失敗 | log 1 file の削除成功と `os.IsNotExist` は `deleted_count` を 1 増やす。それ以外の削除失敗は当該 file だけ `failed_count` を 1 増やし、`LOG_CLEANUP_DELETE_FAILED` と id だけを WARN で記録し、再試行せず次の file へ進む。archive directory の再読込または削除失敗は `LOG_ARCHIVE_DIRECTORY_CLEANUP_FAILED` と basename だけを WARN で記録し、`deleted_count` / `failed_count` を変更しない。 |
| 処理結果 | archive は `archived_count`、cleanup は `deleted_count` と `failed_count` を処理結果として返す。 |

**archive / cleanup 実装確認ゲート：**

| 観点 | 入力 | 合格条件 | 禁止条件 |
|------|------|----------|----------|
| 対象列挙 | `.build_logs/*.json`、`.build_state.current_build_id`、`.server_config.log_archive_after_days` | 対象候補を file 名辞書順で走査し、`finished_at` が閾値より古く、実行中 build でない log だけを対象にする。 | file mtime だけで archive 対象にすること、`.build_logs/archive/` 配下を再対象化すること。 |
| JSON 検証 | 通常 build log JSON | JSON object、`id`、`finished_at`、file 名 id 一致、UTC ISO 8601 秒精度を確認する。 | 破損 log の修復、未知 key の削除保存、対象外 log の削除。 |
| gzip 作成 | 対象 `.build_logs/{id}.json` | `.build_logs/archive/{id}.json.gz.tmp.{pid}` へ gzip 出力し、close 後に `.json.gz` へ rename する。 | 未完了 gzip を公開 path に置くこと、既存 `.json.gz` の上書き。 |
| 元 log 削除 | gzip 作成成功済み対象 | gzip を展開して JSON parse と id 一致を再確認した後、元 `.json` だけ削除する。 | gzip 検証前の元 log 削除、archive 失敗時の元 log 削除。 |
| archive 読取 | 通常 log 不在、archive log あり | gzip 展開後、通常 log と同じ schema の JSON object として返す。 | 展開済み JSON を通常 log として再保存すること。 |
| cleanup | 通常 log、archive log | retention 対象の通常 log、archive log を固定順で削除し、失敗を `failed_count` へ計上する。 | 一部失敗時の処理中断、失敗対象の自動 chmod / rename 修復。 |
| secret / log | WARN / ERROR 出力 | 固定 code、id、path basename、HTTP status 相当だけを出す。 | build log 本文、token、Authorization header、credential 付き URL、gzip 内容の出力。 |

**archive fixture 参照：**

archive / snapshot fixture の fixture 名、expected file、effects、fake filesystem、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。[`docs/details/archive.md`](archive.md) 詳細本文責務では、archive owner の保存、読取、download、delete、rollback 実体処理と状態差分だけを扱う。

検証観点:

| ケース | 期待結果 |
|--------|----------|
| 対象ログあり | `.json.gz` 作成、元 `.json` 削除、API 参照可。 |
| 実行中ログ | archive しない。 |
| 破損ログ | archive しない、WARN、処理継続。 |
| archive API | 件数を返し、disk usage に archive bytes を含める。 |
| cleanup | 通常 log と archive log の両方を保持期間で削除する。 |
| archive 既存 | 上書きせず skip。 |
| cleanup 一部失敗 | `failed_count` に計上し処理継続。 |

<a id="sec-27-15"></a>
**27.15 ビルドアーティファクト管理：**

本機能の目的は、`.snapshots/` に保存する build artifact の作成、世代削除、一覧読取、保存済み tar.gz の検証・配信、削除、rollback 展開・転送の実体処理を archive owner に固定することである。API endpoint、SDK method、UI 操作表示の境界は [§27.15 API / SDK / UI 共通参照先](#2715-api--sdk--ui-共通参照先) を参照する。

owner component は `archive` とする。collaborator component は `api`、`sdk`、`ui`、`runner`、`statefile` とする。snapshot 作成の呼出条件と入力引渡しは [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理)、snapshot の保存実体と世代削除は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を正本とする。

archive owner は snapshot の保存形式、保存、世代削除、一覧読取、保存済み tar.gz の検証・配信、delete 実体処理、rollback 展開・転送実体処理を担当する。api / sdk / ui の境界は [§27.15 API / SDK / UI 共通参照先](#2715-api--sdk--ui-共通参照先)、runner の通常 build 実行、snapshot 作成タイミング、build history / status finalizer は [`docs/details/runner.md`](runner.md) 詳細本文責務を参照する。

<a id="2715-api--sdk--ui-共通参照先"></a>
**[§27.15 API / SDK / UI 共通参照先](archive.md#2715-api--sdk--ui-共通参照先)：**

[§27.15](archive.md#sec-27-15) で HTTP status、JSON error、streaming response、API endpoint、request / response を述べる場合は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、SDK method と error 変換は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI 表示と直接操作禁止は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) を共通参照先とする。各表では archive owner が担当する実体処理と状態差分だけを記載する。

**API 呼び出し境界参照：**

| API | 処理 |
|-----|------|
| `GET /api/snapshots` | archive owner は snapshot 一覧読取結果だけを返す。 |
| `GET /api/snapshots/{id}/download` | archive owner は保存済み `site.tar.gz` の事前検証と、検証した同一 file descriptor からの byte stream 引渡しだけを担当する。HTTP header / status は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) を参照する。 |
| `DELETE /api/snapshots/{id}` | archive owner は事前検証と snapshot delete 実体処理を行い、`deleted`、`delete_partial`、`delete_failed` のいずれかを返す。`.config_log`、`.audit_log`、HTTP response は `api` owner の責務とする。 |
| `POST /api/history/{id}/rollback` | archive owner は rollback 用 artifact の事前検証、展開、転送、temporary cleanup を行い、deploy 結果と pending 候補値を `runner` owner へ返す。rollback build の lock、ID、log、history、status、pending は [`docs/details/runner.md`](runner.md) 詳細本文責務を参照する。 |

`id` は build id と一致するものだけ許可する。snapshot 専用 id は採番しない。`/`、`..`、空文字、URL decode 後に path separator を含む値は失敗扱いとする。

**snapshot 保存固定契約：**

| 項目 | 仕様 |
|------|------|
| snapshot id | build id と同一。`snap{YYYYMMDDHHmmss}` 形式の専用 id は作成しない。 |
| 保存 path | `{StateDir}/.snapshots/{build_id}`。`snapshots_keep > 0` の save 呼出時に `.snapshots` が不在なら archive owner が mode `0700` で作成する。既存 path が directory 以外または symlink、mode が `0700` より広い場合は自動補正せず `failed` とする。保存中は `{StateDir}/.snapshots/{build_id}.tmp.{pid}` を `O_EXCL` 相当で mode `0700` の新規 directory として作成し、existing path を再利用しない。 |
| 入力 | `build_id`、output site 絶対 path、`saved_at`、`output_sha256`、`snapshots_keep`。値の引渡し条件は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。 |
| 保存対象 | output site 配下の通常 file と、それを含むために必要な directory だけ。archive 内では output site 直下を root とし、追加の top-level wrapper directory を作らない。 |
| 禁止 source | symlink、hardlink、socket、device、fifo、`.git`、`.github_token`、`.admin_credentials`、`.api_tokens`、`.smtp_secret`、`.webhook_secret`、runner 状態 file、lock、pending queue を検出した場合は保存全体を `failed` とし、skip して不完全な snapshot を作成しない。 |
| source 整合性 | output site root は symlink でない既存 directory とする。列挙時と file open 後の `fstat` で device / inode / type / size / mtime を比較し、通常 file は link count `1` を必須とする。open は Linux `O_NOFOLLOW` 相当で行い、読取後の `fstat` で size / mtime が変化した場合も保存全体を `failed` とする。列挙後に追加・削除・置換された path を推測で継続しない。 |
| archive 形式 | `.snapshots/{build_id}/site.tar.gz` と `.snapshots/{build_id}/meta.json` の 2 file だけを mode `0600` で作成する。tar writer は `archive/tar.FormatUSTAR`、gzip writer は `compress/gzip.BestCompression` に固定する。USTAR で表現できない path または size は保存全体を `failed` とし、PAX / GNU 拡張へ fallback しない。 |
| entry path | output site root からの valid UTF-8、slash 区切り相対 path。root 自身は entry にしない。directory entry だけ末尾 `/` を 1 個付け、通常 file entry は末尾 `/` なしとする。先頭 `/`、`.`、`..`、空 segment、NUL、CR、LF、backslash、Windows drive prefix は禁止する。 |
| entry 順 | directory と通常 file を区別せず、entry path の byte 辞書順に書き込む。同一 path は 1 entry だけとする。 |
| tar header | format は USTAR。directory は `TypeDir`、size `0`、mode `0755`、通常 file は `TypeReg`、実 byte size、source に実行 bit が 1 つでもあれば mode `0755`、なければ `0644`。uid / gid / device major / device minor は `0`、user / group / link name は空、mtime は全 entry で `saved_at`、atime / ctime は zero value とする。 |
| gzip header | mtime は `saved_at`、name / comment は空、OS field は `255` とする。 |
| `meta.json` | `id`、`build_id`、`saved_at`、`size_bytes`、`file_count`、`output_sha256` を必須 key とする。UTF-8、BOM なし、JSON object 1 個、末尾 LF 1 個で保存し、保存 key 順は `id`、`build_id`、`saved_at`、`size_bytes`、`file_count`、`output_sha256` に固定する。`size_bytes` は tar 内の通常 file の非圧縮 byte 合計、`file_count` は通常 file entry 数とし、`site.tar.gz` と `meta.json` 自身を含めない。 |
| atomic publish | tmp directory 内で `site.tar.gz` と `meta.json` を作成し、tar / gzip writer close → underlying archive file sync → archive file close → metadata write 完了 → metadata file sync → metadata file close → tmp directory sync → [保存済み archive 検証・配信契約](#snapshot-archive-validation) の事前検証・metadata 再計算 → `{build_id}` へ rename → `.snapshots` directory sync の順に行う。rename 前の失敗は public path を作成せず `failed`、rename 後の `.snapshots` sync 失敗は検証済み public snapshot を維持し、結果 `saved`、warning `SNAPSHOT_DIRECTORY_SYNC_FAILED`として prune を実行しない。成功・失敗のどちらでも tmp path 削除を 1 回試行し、削除失敗は `SNAPSHOT_TMP_CLEANUP_FAILED` warning を返す。 |
| 既存 snapshot | 同じ build id が存在する場合は上書きしない。`meta.json` と `site.tar.gz` を [保存済み archive 検証・配信契約](#snapshot-archive-validation) で検証し、`id`、`build_id`、`output_sha256` が呼出入力と一致する場合だけ `exists_valid` と `SNAPSHOT_EXISTS` を返す。破損または不一致は `failed` と `SNAPSHOT_EXISTS_MISMATCH` を返し、既存 directory を変更しない。 |
| 世代削除 | 新 snapshot の publish 成功後だけ実行する。schema-valid な snapshot を `saved_at` 昇順、同時刻は id 昇順に並べ、`snapshots_keep` 超過分だけ古い順に削除する。破損 snapshot は自動削除しない。削除は [snapshot delete 副作用固定契約](#snapshot-delete-side-effect-contract) と同じ tombstone rename、directory sync、cleanup、directory sync の実装原語を使用する。`deleted` 以外は当該 id を失敗として記録し、残りの超過分を続行する。保存結果は `saved` のまま、失敗 id が 1 件以上なら warning `SNAPSHOT_PRUNE_FAILED` を 1 回返す。`delete_partial` の public path を復元せず、`delete_failed` の public path を自動削除しない。 |

<a id="snapshot-archive-validation"></a>
**保存済み archive 検証・配信契約：**

download のために新しい tar.gz を生成してはならない。archive owner は `.snapshots/{id}/site.tar.gz` を 1 回だけ open し、同じ file descriptor を事前検証後に先頭へ seek して API owner へ引き渡す。検証と配信の間に path から reopen してはならない。

| 検証項目 | 合格条件 |
|----------|----------|
| snapshot directory | caller が渡す検証対象 directory と expected build id を別値として扱う。directory、`meta.json`、`site.tar.gz` が symlink でなく、後二者が通常 file であることを必須とする。public snapshot 検証では directory 名が expected build id と一致し、save 中の tmp directory 検証では tmp 名の一致を要求しない。 |
| metadata | `meta.json` が schema-valid、`id=build_id={id}`、未知 key なし。 |
| gzip | header が保存固定契約に一致し、CRC と終端まで読み取れる。 |
| tar entry | output site 内容を root とする相対 path、相対 path 辞書順、重複なし、通常 file / directory だけ、header 固定値一致。symlink、hardlink、device、socket、fifo、GNU / PAX 拡張 entry は不許可。 |
| secret | 各 path segment が保存固定契約の禁止 source 名に一致しない。 |
| 再計算 | 通常 file entry の非圧縮 byte 合計と件数が `meta.json.size_bytes` と `file_count` に一致する。`meta.json.output_sha256` が非 `null` の場合は、tar 内の通常 file entry から [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の出力成果物 manifest SHA-256 を再計算し、一致を必須とする。 |

事前検証に失敗した場合は stream を返さない。stream 開始後に同一 file descriptor の read error が発生した場合は stream を中断し、`SNAPSHOT_STREAM_FAILED: id={id} entry=site.tar.gz` を返す。archive owner は HTTP header、HTTP status、JSON error body を作成しない。状態 file、snapshot directory、history、build log、config log は変更しない。

| archive 結果 | 状態差分 | 必須 log code |
|----------------|----------|---------------|
| `stream_ready` | なし。 | なし。 |
| `unsafe_entry` | なし。 | `SNAPSHOT_UNSAFE_ENTRY: id={id} entry={path}` |
| `secret_entry` | なし。 | `SNAPSHOT_SECRET_ENTRY: id={id} entry={path}` |
| `metadata_mismatch` | なし。 | `SNAPSHOT_META_MISMATCH: id={id}` |
| `stream_failed` | なし。 | `SNAPSHOT_STREAM_FAILED: id={id} entry=site.tar.gz` |

**Rollback 仕様：**

rollback で archive owner が受け取る値は `snapshot_id`、[`docs/details/runner.md` 詳細本文責務 build id 契約](runner.md#build-id-契約) により runner owner が採番した `new_build_id`、deploy target 一覧とする。archive owner は元 snapshot を変更せず、展開・転送結果を `success`、`failure_build`、`success_deploy_pending` のいずれかで返す。`success_deploy_pending` では [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の PendingTransfer object に必要な値を pending 候補として返すが、`.pending_transfers` へは書き込まない。

rollback は snapshot 内の成果物を deploy target へ再転送する操作であり、以下を行ってはならない。

| 禁止副作用 | 理由 |
|----------|------|
| `.last_sha` 更新 | rollback は監視対象 SHA の処理完了ではない。 |
| `.server_config`、`.branch_config`、`.notify_config` の復元 | 設定 rollback ではない。 |
| `.build_history` の過去行書き換え | rollback は新規履歴として追記する。 |
| `.build_logs/{元id}.json` の変更 | 元 build の証跡を保持する。 |
| 元 snapshot の削除または上書き | rollback 成否に関係なく元成果物を保持する。 |
| secret / token / credentials の復元 | snapshot に secret を含めないため復元対象外。 |

rollback build log / history / status / pending の値と書込順は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。archive owner は `.build_lock`、`.build_state`、`.build_status.json`、`.build_logs`、`.build_history`、`.pending_transfers` を書き込まず、rollback の queue も作成しない。

rollback は [保存済み archive 検証・配信契約](#snapshot-archive-validation) の事前検証を通過した後、`{StateDir}/.snapshots/.rollback.{new_build_id}.tmp` へ `site.tar.gz` を展開する。temporary directory は `.snapshots` 直下に existing path 非上書き、symlink 非追従、mode `0700` で新規作成し、existing path がある場合は展開せず失敗する。展開は事前検証で使用した同一 file descriptor を先頭へ seek して行い、entry ごとに再度 path と type を検証する。directory は `0755`、file は tar header の `0644` または `0755` で作成し、existing file、symlink、hardlink を追従または上書きしない。deploy はこの展開済み temporary directory だけを転送元とする。成功・失敗のどちらでも temporary directory の削除を 1 回試行し、削除失敗は rollback 結果を反転させず `SNAPSHOT_ROLLBACK_TMP_CLEANUP_FAILED` warning を返す。

**snapshot 一覧・削除固定契約：**

| 項目 | 仕様 |
|------|------|
| 一覧対象 | `.snapshots/{id}/meta.json` と `.snapshots/{id}/site.tar.gz` がとも存在し、[保存済み archive 検証・配信契約](#snapshot-archive-validation) の事前検証を通過する build id 形式の directory だけ。`.` で始まる save / rollback / pending / delete temporary path は一覧、破損 warning、prune 対象から除外する。 |
| size | `meta.json.size_bytes` の値。directory 使用量、圧縮後 `site.tar.gz` size、`meta.json` size を合算しない。 |
| delete 順 | api owner の id validation → runner owner の snapshot delete guard 取得・state 再確認 → archive owner の snapshot 事前検証・delete → 三値結果返却 → api owner の結果別 log 処理 → runner owner の guard 解放 → HTTP response。archive owner は guard を取得または解放しない。 |
| delete partial / log 失敗 | public snapshot 削除後の tombstone cleanup / directory sync 失敗は archive owner が `delete_partial` を返す。api owner の `.config_log` / `.audit_log` 書込順、書込可否、guard 解放、HTTP は [`docs/details/api.md` 詳細本文責務 snapshot delete API 契約](api.md#snapshot-delete-api-contract) を正とする。必須 log 失敗の場合も `500` とし、public snapshot と先行 log を巻き戻さない。 |
| rollback pending | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の PendingTransfer object を使用し、`trigger="rollback"`、`source_kind="snapshot"`、`out=null`、`rollback_from=snapshot_id`、`output_sha256=meta.json.output_sha256` とする。`output_sha256=null` の snapshot は pending 候補を返さず `failure_build` とする。 |

<a id="snapshot-delete-side-effect-contract"></a>
**snapshot delete 副作用固定契約：**

delete は destructive 操作であるため、archive owner が行う事前検証、削除、結果返却と失敗時副作用を固定する。

| 段階 | 成功条件 | 失敗時副作用 |
|------|----------|--------------|
| id validation | build id 形式、path separator なし、URL decode 後も安全。 | snapshot、config log、history、build log、pending、state 差分なし。 |
| delete guard | runner owner が `.build_lock` を取得後に `.build_state.running=false`、`current_build_id=null` を再確認し、guard を返す。 | guard 未取得で archive owner を呼び出さず、snapshot、config log、history、build log、pending 差分なし。 |
| 存在確認 | `.snapshots/{id}/meta.json` と `site.tar.gz` が存在し、[保存済み archive 検証・配信契約](#snapshot-archive-validation) を通過する。 | 差分なし。 |
| delete publish | 対象 snapshot directory を同一 `.snapshots` 内の `.delete.{id}.{pid}.tmp` へ existing path 非上書きで rename し、`.snapshots` を sync する。rename 成功後は public snapshot 削除済みとする。 | rename 前失敗は対象 snapshot を維持し、config / audit log 追記なし。rename 後の sync 失敗は public path を戻さず `delete_partial` を返す。 |
| tombstone cleanup | rename と最初の parent sync 成功後に tombstone を再帰削除し、`.snapshots` を再度 sync する。symlink を追従しない。 | cleanup または最終 sync 失敗は public snapshot 削除済みのまま `delete_partial`、warning `SNAPSHOT_DELETE_CLEANUP_FAILED` を返す。 |
| result handoff | 全削除段階成功は `deleted`、public path 削除後の失敗は `delete_partial`、public path 維持の失敗は `delete_failed` と snapshot id を api owner へ返す。 | HTTP status / body、`.config_log`、`.audit_log` を archive owner が生成または追記しない。 |

**artifact 実装確認固定契約：**

| 操作 | 確認条件 | 失敗時副作用 |
|------|----------|--------------|
| snapshot save | 固定順の `site.tar.gz` と schema-valid な `meta.json` を tmp directory に作成し、再検証後に `.snapshots/{build_id}` へ rename する。`meta.json.output_sha256` は runner が引き渡した値と一致する。 | tmp 作成中の失敗では public snapshot directory を作らない。既存 snapshot は変更しない。 |
| snapshot list | `meta.json` と `site.tar.gz` が [保存済み archive 検証・配信契約](#snapshot-archive-validation) を通過した snapshot だけを `saved_at` 降順、同時刻 id 降順で返す。 | 破損 snapshot は除外し、WARN `SNAPSHOT_META_CORRUPT`。修復しない。 |
| download | 保存済み `site.tar.gz` の全 entry を事前検証し、同じ file descriptor から保存済み byte 列をそのまま stream する。 | unsafe entry、secret file、metadata 不一致では stream を開始せず、固定 code を返す。stream 開始後の read error では stream を中断する。状態は変更しない。 |
| delete | 対象 snapshot directory を tombstone へ atomic rename 後に cleanup し、`deleted`、`delete_partial`、`delete_failed` のいずれかを api owner へ返す。 | rename 前失敗は public snapshot を維持する。rename 後失敗は public snapshot を巻き戻さず partial とする。archive owner は `.config_log`、`.audit_log`、history、build log、pending を変更しない。 |
| rollback success | 検証済み archive を展開し、deploy result、pending 候補、warning を runner owner へ返す。 | 展開または deploy 失敗を `failure_build` または `success_deploy_pending` で返し、元 snapshot と `.last_sha` を変更しない。 |

**snapshot pending retry 固定契約：**

runner owner が `source_kind="snapshot"` の PendingTransfer object を再試行する場合、archive owner は当該 entry だけを入力とし、元 snapshot から再送可能な一時成果物を毎回作成する。rollback 開始時の temporary directory または branch target の現在 `out` を再利用してはならない。

| 段階 | 契約 |
|------|------|
| schema 検証 | `trigger="rollback"`、`source_kind="snapshot"`、`out=null`、`rollback_from=snapshot_id`、`output_sha256` 非 `null` を必須とする。不一致は filesystem と SSH を変更せず schema error を runner owner へ返す。 |
| snapshot 検証 | `.snapshots/{snapshot_id}` を [保存済み archive 検証・配信契約](#snapshot-archive-validation) で検証し、`meta.json.output_sha256` と entry の `output_sha256` を一致させる。不在、破損、checksum 不一致で SSH を実行しない。 |
| 展開 | 事前検証に使用した同一 `site.tar.gz` file descriptor を先頭へ seek し、`{StateDir}/.snapshots/.pending.{build_id}.{target_id}.tmp` を existing path 非上書き、symlink 非追従、mode `0700` の新規 directory として作成して展開する。existing path がある場合は展開せず失敗する。entry の directory は `0755`、file は tar header の `0644` または `0755` で新規作成し、existing file、symlink、hardlink を追従または上書きしない。 |
| 転送 | entry に保存された `host`、`user`、`dest_dir`、`target_id` の当該 1 target だけへ、展開済み temporary directory を転送元として転送・remote checksum 検証する。現在の `.branch_config` で転送先を上書きしない。 |
| cleanup | 転送成功、転送失敗、展開失敗のいずれでも temporary directory 削除を 1 回試行する。cleanup 失敗は転送結果を反転させず `SNAPSHOT_PENDING_TMP_CLEANUP_FAILED` warning を runner owner へ返す。 |
| 結果 | 転送と remote checksum 検証成功は `success`、再送可能な SSH 失敗は `retryable_failure`、snapshot 不在は `snapshot_missing`、snapshot 破損は `snapshot_corrupt`、checksum 不一致は `snapshot_changed`を返す。archive owner は `.pending_transfers` を更新しない。 |

**snapshot `meta.json` schema 固定契約：**

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `id` | string | 必須 | build id と同一。 |
| `build_id` | string | 必須 | build id と同一。 |
| `saved_at` | string | 必須 | UTC ISO 8601 秒精度。 |
| `size_bytes` | integer | 必須 | snapshot 対象通常ファイル合計 bytes。 |
| `file_count` | integer | 必須 | snapshot 対象通常ファイル数。 |
| `output_sha256` | string/null | 必須 | build history の `output_sha256`。不明時 `null`。 |

未知 key は read 時に無視せず `SNAPSHOT_META_CORRUPT` としてその snapshot を一覧から除外する。`size_bytes` と `file_count` は [保存済み archive 検証・配信契約](#snapshot-archive-validation) に従って非圧縮の通常 file entry から再計算し、`meta.json` と不一致なら `metadata_mismatch` とする。

**rollback owner 接続順固定契約：**

1. api owner が path / auth / rate limit / maintenance / circuit を検証し、runner owner の rollback coordinator を呼び出す。
2. runner owner が lock を取得し、`new_build_id` を採番する。
3. archive owner が [保存済み archive 検証・配信契約](#snapshot-archive-validation) に従って `meta.json` と `site.tar.gz` を検証し、rollback temporary directory へ展開する。
4. runner owner が rollback 開始状態と running build log を保存し、api owner へ一回性 `prepared` handle と `new_build_id` を返す。
5. api owner が rollback の `build_trigger` audit を追記し、成功時だけ runner owner の worker 開始境界を呼び出す。
6. archive owner が展開済み artifact を deploy target へ転送し、temporary directory を cleanup し、deploy 結果、pending 候補、warning を runner owner へ返す。
7. runner owner が rollback build log、history、pending、status、build state を順番に確定し、lock を解放する。
8. api owner は worker 開始境界が `accepted` を返した場合のみ `202` を返し、archive owner が HTTP response を作成しない。

手順 4 の `prepared` 前に失敗した場合の部分書込、手順 5 の audit 失敗時の abort、accepted 後の finalizer は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。どの場合も、元 snapshot、元 build log、過去 history、`.last_sha` は変更しない。

**rollback 実装確認ゲート：**

| 観点 | 合格条件 |
|------|----------|
| lock / running | archive owner は runner owner が lock を取得済みであることを呼出前提とし、lock や running state を直接変更しない。 |
| id / trigger | runner owner から `new_build_id` と rollback 元 id を受け取り、archive owner は採番、build log / history / pending transfer 書込を行わない。 |
| extract / deploy | 検証済み `site.tar.gz` を state directory 内の新規 temporary directory に展開し、展開済み artifact だけを deploy target へ転送する。builder、GitHub read、Commit Status、SHA cache 更新、通常 snapshot 作成を行わない。 |
| success | `success`、deploy target 別結果、warning を runner owner へ返し、元 snapshot、元 build log、過去 history、`.last_sha` が unchanged。 |
| deploy failure | 再送可能な失敗は `success_deploy_pending` と pending 候補、再送不可の失敗は `failure_build` を返し、元 snapshot、`.last_sha`、元 build log を変更しない。 |
| pending | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の PendingTransfer object に必要な `source_kind="snapshot"`、`out=null`、snapshot id、snapshot manifest SHA-256、初期 `retry_count=0` を pending 候補で返し、archive owner は `.pending_transfers` を書き込まない。 |
| finalizer failure | archive owner は temporary cleanup warning だけを返す。runner finalizer、server log、lock 解放は [`docs/details/runner.md`](runner.md) 詳細本文責務を参照する。 |

**archive / snapshot fixture 証跡参照：**

archive / snapshot の fixture 名、合格条件、expected / effects、stream failure、secret absence、状態差分、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。[`docs/details/archive.md`](archive.md) 詳細本文責務では、snapshot 保存形式、`meta.json` schema、一覧 sort、保存済み tar.gz の entry 順序・検証・配信、delete 実体、rollback 用 artifact の展開・転送・temporary cleanup、元 snapshot 維持、`.last_sha` 非変更など archive owner の実体処理観点だけを扱う。HTTP header / status / body、`.config_log`、rollback 状態書込は各 owner を参照する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 一覧 | snapshot id、build id、保存日時、size が返る。 |
| download | 保存済み `site.tar.gz` と byte 一致する stream を返し、snapshot 外の file を参照しない。 |
| delete | 対象 id だけを削除し、api owner へ `deleted` を返す。config log と audit は api owner の fixture で確認する。 |
| rollback 成功 | 展開済み artifact の deploy 結果 `success`、target 別結果、warning を runner owner へ返し、元 snapshot を維持する。 |
| 不正 id | 状態差分なし。 |
| download symlink | symlink entry を含めず、secret 名検出時は stream 開始前に中止する。 |
| rollback pending | `success_deploy_pending` と [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema に一致する pending 候補を runner owner へ返し、archive owner は状態 file を書き込まない。 |
| rollback running | api / runner owner が archive owner を呼び出さず、archive 状態差分なし。 |
| delete log failure | archive owner は削除成功後の log 失敗に関与せず、削除済み snapshot を復元しない。 |
