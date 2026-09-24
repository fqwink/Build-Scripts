# Adlaire CI — Archive 詳細仕様

本ファイルは `archive` owner component の詳細本文責務として、`archive` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`archive` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、表示、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `archive` |
| collaborator component | `runner`、`api`、`sdk`、`ui`、`statefile` |
| 持つ内容 | `archive` owner が主本文として定義する build log archive / cleanup の実体処理、snapshot 保存形式、download tar.gz 生成安全性、snapshot delete 実体処理、rollback 転送実体処理。 |
| 持たない内容 | runner の通常 build 実行、snapshot 作成トリガー判定、API 共通 request / response、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 手順、fixture 証跡責務。 |

archive owner は、保存済み build log と snapshot artifact の圧縮、展開、列挙、削除、転送の実体処理だけを担当する。API / SDK / UI / runner 境界は collaborator component の詳細本文責務を参照する。

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
| cleanup 順 | 通常 log 削除 → archive log 削除 → 空 archive directory 削除試行。 |
| 削除失敗 | 処理継続し、処理結果に `failed_count` を含める。 |
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

本機能の目的は、`.snapshots/` に保存された build artifact について、archive owner が一覧読取、download tar.gz 生成、削除、rollback 転送の実体処理を固定することである。API endpoint、SDK method、UI 操作表示の境界は [§27.15 API / SDK / UI 共通参照先](#2715-api--sdk--ui-共通参照先) を参照する。

owner component は `archive` とする。collaborator component は `api`、`sdk`、`ui`、`runner`、`statefile` とする。snapshot 作成は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。

archive owner は snapshot の保存形式、一覧読取、download tar.gz 生成、delete 実体処理、rollback 転送実体処理を担当する。api / sdk / ui の境界は [§27.15 API / SDK / UI 共通参照先](#2715-api--sdk--ui-共通参照先)、runner の通常 build 実行、snapshot 作成タイミング、build history / status finalizer は [`docs/details/runner.md`](runner.md) 詳細本文責務を参照する。

<a id="2715-api--sdk--ui-共通参照先"></a>
**[§27.15 API / SDK / UI 共通参照先](archive.md#2715-api--sdk--ui-共通参照先)：**

[§27.15](archive.md#sec-27-15) で HTTP status、JSON error、streaming response、API endpoint、request / response を述べる場合は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、SDK method と error 変換は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI 表示と直接操作禁止は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) を共通参照先とする。各表では archive owner が担当する実体処理と状態差分だけを記載する。

**API 呼び出し境界参照：**

| API | 処理 |
|-----|------|
| `GET /api/snapshots` | archive owner は snapshot 一覧読取結果だけを返す。 |
| `GET /api/snapshots/{id}/download` | archive owner は download tar.gz 生成だけを担当する。 |
| `DELETE /api/snapshots/{id}` | archive owner は snapshot delete 実体処理と `.config_log` 追記境界だけを担当する。 |
| `POST /api/history/{id}/rollback` | archive owner は rollback 転送実体処理を担当する。rollback build log/history の作成境界は [`docs/details/runner.md`](runner.md) 詳細本文責務を参照する。 |

`id` は build id と一致するものだけ許可する。snapshot 専用 id は採番しない。`/`、`..`、空文字、URL decode 後に path separator を含む値は失敗扱いとする。

**snapshot 保存固定契約：**

| 項目 | 仕様 |
|------|------|
| snapshot id | build id と同一。`snap{YYYYMMDDHHmmss}` 形式の専用 id は作成しない。 |
| 保存 path | `{StateDir}/.snapshots/{build_id}`。保存中は `{StateDir}/.snapshots/{build_id}.tmp.{pid}` を使い、完了後に rename する。 |
| 保存対象 | 出力サイトディレクトリ配下の通常ファイルだけ。symlink、hardlink、socket、device、fifo、`.git`、secret、runner 状態ファイル、lock、pending queue は保存対象外。 |
| archive 形式 | `.snapshots/{build_id}/site.tar.gz` と `.snapshots/{build_id}/meta.json` を作成する。 |
| `meta.json` | `id`、`build_id`、`saved_at`、`size_bytes`、`file_count`、`output_sha256` を必須 key とする。 |
| 既存 snapshot | 同じ build id の snapshot が既に存在する場合は上書きせず、WARN `SNAPSHOT_EXISTS: id={id}` を出して snapshot 保存を skip する。 |
| 世代削除 | 新 snapshot 作成成功後に `saved_at` 昇順で `snapshots_keep` 超過分だけ削除する。削除失敗は build 成功を失敗へ反転せず、WARN `SNAPSHOT_PRUNE_FAILED` を出す。 |

**download tar.gz 生成契約：**

| 項目 | 仕様 |
|------|------|
| root | `.snapshots/{id}/` を root とし、root 外を参照しない。 |
| entry path | root からの相対 path。`/` 始まり、`..`、空 segment、NUL、Windows drive prefix は禁止。 |
| entry 種別 | 通常ファイルと directory だけを含める。symlink、hardlink、device、socket、fifo は含めない。 |
| header | `Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"`。 |
| 順序 | directory、file とも相対 path 辞書順。 |
| mtime | snapshot 内 file の mtime を使用する。snapshot 内 file から mtime を取得できない場合は build log の `finished_at` を使用する。 |
| secret 除外 | `.github_token`、`.admin_credentials`、`.api_tokens`、`.smtp_secret`、`.webhook_secret`、runner 状態ファイル名は検出時点で download を中止する。 |

**download stream 固定契約：**

download は stream 開始前に snapshot directory 全体を走査し、entry path、entry 種別、secret 禁止名、`meta.json` schema、`size_bytes`、`file_count` を検証する。stream 開始前検証に失敗した場合は binary header を送信しない。stream 開始後に read error が発生した場合は stream を中断し、server log に `SNAPSHOT_STREAM_FAILED: id={id} entry={path}` を出す。stream 開始後は JSON error body を追加送信してはならない。状態ファイル、snapshot directory、history、build log、config log は変更しない。

| ケース | HTTP / stream | 状態差分 | 必須 log |
|--------|---------------|----------|----------|
| 事前検証成功 | `200`、binary tar.gz。 | なし。 | なし。 |
| unsafe entry | `500 {"error":"Snapshot download failed"}`。binary header なし。 | なし。 | `SNAPSHOT_UNSAFE_ENTRY: id={id} entry={path}` |
| secret file 検出 | `500 {"error":"Snapshot download failed"}`。binary header なし。 | なし。 | `SNAPSHOT_SECRET_ENTRY: id={id} entry={path}` |
| meta 不一致 | `500 {"error":"Snapshot metadata mismatch"}`。binary header なし。 | なし。 | `SNAPSHOT_META_MISMATCH: id={id}` |
| stream 中 read error | stream 中断。JSON body 追加なし。 | なし。 | `SNAPSHOT_STREAM_FAILED: id={id} entry={path}` |

**Rollback 仕様：**

rollback は新しい build id を採番し、`.build_history.trigger="rollback"`、`rollback_from=<元id>` を保存する。元 snapshot は変更しない。rollback 中に別 build が running の場合は状態差分なしで中止する。転送失敗時は rollback build log を `failure` とし、元 snapshot は削除しない。

rollback は snapshot 内の成果物を deploy target へ再転送する操作であり、以下を行ってはならない。

| 禁止副作用 | 理由 |
|----------|------|
| `.last_sha` 更新 | rollback は監視対象 SHA の処理完了ではない。 |
| `.server_config`、`.branch_config`、`.notify_config` の復元 | 設定 rollback ではない。 |
| `.build_history` の過去行書き換え | rollback は新規履歴として追記する。 |
| `.build_logs/{元id}.json` の変更 | 元 build の証跡を保持する。 |
| 元 snapshot の削除または上書き | rollback 成否に関係なく元成果物を保持する。 |
| secret / token / credentials の復元 | snapshot に secret を含めないため復元対象外。 |

rollback build log は `target_status="success"` または `failure_build` とし、`trigger="rollback"`、`rollback_from=<元id>`、`snapshot_id=<元id>` を含める。rollback 転送で pending が発生した場合は `success_deploy_pending` とし、`.pending_transfers` に rollback 用 entry を追加する。

rollback 開始時は `.build_lock` を取得し、取得できない場合は状態差分なしで中止する。`.build_lock` 取得後に `.build_state.running=true`、`current_build_id=<new_id>` を保存し、転送完了後に finalizer で `running=false` とする。rollback は queue に積まない。

**snapshot 一覧・削除固定契約：**

| 項目 | 仕様 |
|------|------|
| 一覧対象 | `.snapshots/{id}/meta.json` が存在する directory だけ。 |
| size | directory 配下の通常ファイル size 合計。symlink は size 集計前に異常扱い。 |
| delete 順 | id validation → running check → snapshot directory 確認 → delete → `.config_log` 追記 → response。 |
| delete log 失敗 | snapshot 削除済みのまま。削除は巻き戻さない。 |
| rollback pending | pending entry には `rollback_from`、`snapshot_id`、deploy target を保存する。 |

**snapshot delete 副作用固定契約：**

delete は destructive endpoint であるため、成功条件と失敗時副作用を [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) の固定表に固定する。

| 段階 | 成功条件 | 失敗時副作用 |
|------|----------|--------------|
| id validation | build id 形式、path separator なし、URL decode 後も安全。 | snapshot、config log、history、build log、pending、state 差分なし。 |
| running check | `.build_state.running=false`。 | snapshot、config log、history、build log、pending 差分なし。 |
| 存在確認 | `.snapshots/{id}/meta.json` が schema valid。 | 差分なし。 |
| delete | 対象 snapshot directory だけ削除成功。 | 対象 snapshot が残る。config log 追記なし。 |
| config log | `.config_log` に `target="snapshot_delete"`、`target_id={id}` を追記。 | snapshot は削除済みのまま。他 snapshot、history、build log、pending は変更しない。 |
| response | `{ "message":"Snapshot deleted" }`。 | 成功 response を返さない。 |

**artifact 実装確認固定契約：**

| 操作 | 確認条件 | 失敗時副作用 |
|------|----------|--------------|
| snapshot save | `site.tar.gz` と `meta.json` を tmp directory に作成し、検証後に `.snapshots/{build_id}` へ rename する。`meta.json.output_sha256` が build history の値と一致する。 | tmp 作成中の失敗では公開 snapshot directory を作らない。既存 snapshot は変更しない。 |
| snapshot list | `meta.json` が schema valid な snapshot だけを `saved_at` 降順、同時刻 id 降順で返す。 | 破損 snapshot は除外し、WARN `SNAPSHOT_META_CORRUPT`。修復しない。 |
| download | tar.gz entry が root 外を参照せず、相対 path 辞書順で stream される。 | unsafe entry、secret file、symlink 検出時は stream 開始前なら binary header を返さず、開始後なら stream を中断し server log に固定 code を出す。状態は変更しない。 |
| delete | 対象 snapshot directory だけを削除し、`.config_log` に target `snapshot_delete` を追記する。 | `.config_log` 失敗では snapshot は削除済みのまま。他 snapshot は変更しない。 |
| rollback success | 新規 build id、rollback build log、history、`.build_status.json` finalizer、deploy result が整合する。 | deploy 失敗時は rollback build を failure または success_deploy_pending として新規記録し、元 snapshot と `.last_sha` は変更しない。 |

**snapshot `meta.json` schema 固定契約：**

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `id` | string | 必須 | build id と同一。 |
| `build_id` | string | 必須 | build id と同一。 |
| `saved_at` | string | 必須 | UTC ISO 8601 秒精度。 |
| `size_bytes` | integer | 必須 | snapshot 対象通常ファイル合計 bytes。 |
| `file_count` | integer | 必須 | snapshot 対象通常ファイル数。 |
| `output_sha256` | string/null | 必須 | build history の `output_sha256`。不明時 `null`。 |

未知 key は read 時に無視せず `SNAPSHOT_META_CORRUPT` としてその snapshot を一覧から除外する。`size_bytes` と `file_count` は download 時にも再計算し、`meta.json` と不一致なら failure とする。

**rollback 状態更新順固定契約：**

1. snapshot id を検証する。
2. snapshot と `meta.json` を検証する。
3. `.build_lock` を取得する。
4. 新規 build id を採番する。
5. `.build_state.running=true` と `.build_status.json.status="running"` を保存する。
6. snapshot artifact を deploy target へ転送する。
7. rollback build log と `.build_history` を追記する。
8. pending transfer がある場合は `.pending_transfers` に rollback entry を保存する。
9. `.build_status.json` finalizer と `.build_state.running=false` を保存する。
10. `.build_lock` を解放する。

手順 3 より前の失敗は状態差分なしとする。手順 3 以後の失敗は rollback build log に失敗地点、`rollback_from`、`snapshot_id` を残し、`.build_lock` 解放と finalizer を必ず試行する。finalizer 失敗時も、元 snapshot、元 build log、過去 history、`.last_sha` は変更しない。

**rollback 実装確認ゲート：**

| 観点 | 合格条件 |
|------|----------|
| lock / running | `.build_lock` 取得前の validation failure は no-write。lock 取得後は `.build_state.running=true`、finalizer で `false`、lock 解放を必ず試行する。 |
| id / trigger | 新規 build id を採番し、rollback 元 id を build log / history / pending transfer に `rollback_from` と `snapshot_id` で保存する。 |
| deploy | snapshot artifact だけを deploy target へ転送し、builder、GitHub read、Commit Status、SHA cache 更新、通常 snapshot 作成を行わない。 |
| success | rollback build log、history、`.build_status.json` finalizer が成功し、元 snapshot、元 build log、過去 history、`.last_sha` が unchanged。 |
| deploy failure | rollback build log / history は failure または success_deploy_pending として新規保存し、元 snapshot、`.last_sha`、元 build log は unchanged。 |
| pending | `.pending_transfers` entry に `trigger="rollback"`、`rollback_from`、`snapshot_id`、deploy target、retry_count を保存する。 |
| finalizer failure | server log 固定 code、元 snapshot / 元 log / `.last_sha` unchanged。lock 解放は best effort。 |

**archive / snapshot fixture 証跡参照：**

archive / snapshot の fixture 名、合格条件、expected / effects、stream failure、secret absence、状態差分、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。[`docs/details/archive.md`](archive.md) 詳細本文責務では、`meta.json` schema、一覧 sort、download header、tar entry 順序、delete 順、rollback 状態更新順、元 snapshot 維持、`.last_sha` 非変更など archive owner の実体処理観点だけを扱う。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 一覧 | snapshot id、build id、保存日時、size が返る。 |
| download | tar.gz を返し、snapshot 外のファイルを含まない。 |
| delete | 対象 id だけ削除、config log 追記。 |
| rollback 成功 | 新規 build id、trigger rollback、rollback_from 保存。 |
| 不正 id | 状態差分なし。 |
| download symlink | symlink entry を含めず、secret 名検出時は stream 開始前に中止する。 |
| rollback pending | 新規 rollback log/history、pending entry、元 snapshot 維持。 |
| rollback running | 状態差分なし。 |
| delete log failure | snapshot は削除済み。 |
