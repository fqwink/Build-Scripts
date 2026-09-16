# Adlaire CI — Archive 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` から分割した `archive` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

`ADLAIRE_CI_DETAIL_SPEC.md` は、詳細仕様の入口、索引、共通固定値、責務 component 対応表を持つ。本ファイルを読む前に、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j を確認する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `archive` |
| collaborator component | `runner`、`api`、`sdk`、`ui`、`statefile` |
| 持つ内容 | build log archive / cleanup の実体処理、snapshot 保存形式、download tar.gz 生成安全性、snapshot delete 実体処理、rollback 転送実体処理。 |
| 持たない内容 | runner の通常 build 実行、snapshot 作成トリガー判定、API 共通 request / response、SDK method 実装、UI DOM 詳細、状態ファイル schema 定義。 |

archive owner は、保存済み build log と snapshot artifact を安全に圧縮、展開、列挙、削除、転送する実体処理だけを担当する。API は HTTP endpoint の request / response と archive owner 呼び出し境界、SDK は API method 呼び出し、UI は操作表示だけを担当する。runner の build 実行、build id 採番、通常 snapshot 作成タイミング、history / status finalizer は runner owner の詳細仕様を正とし、本ファイルへ重複定義しない。

---

## 対象範囲

| 範囲 | 内容 |
|------|------|
| §27.7 | ビルドログのアーカイブ圧縮。 |
| §27.15 | ビルドアーティファクト管理。 |

---

### 27.7 ビルドログのアーカイブ圧縮

owner component は `archive` とする。collaborator component は `runner`、`api`、`statefile` とする。

archive owner は、runner または `POST /api/logs/archive` から呼び出された場合に、`.server_config.log_archive_after_days > 0` で対象日数より古い `.build_logs/{id}.json` を gzip 圧縮し、`.build_logs/archive/{id}.json.gz` へ保存する。圧縮成功後、元の `.build_logs/{id}.json` を削除する。`.build_logs/archive/` 内のファイルを再圧縮してはならない。

gzip は Go 標準ライブラリ `compress/gzip` を使用し、mtime は元ファイル mtime ではなく圧縮実行時刻でよい。圧縮前 JSON を読み込めないファイルは archive 対象外とし、WARN `LOG_ARCHIVE_SKIP_CORRUPT: id=<id>` を出す。実行中 build の `current_build_id` と一致する log は対象外とする。

archive owner は、通常 log が存在しない場合に archive log を gzip 展開し、通常 `.build_logs/{id}.json` と同じ schema の JSON object として API へ返す。`GET /api/logs/search`、`GET /api/history/{id}/log`、`GET /api/output-meta`、`GET /api/disk-usage` の endpoint、query、response body は `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e を正とし、本節は archive log の探索、展開、除外、件数返却だけを定義する。

archive owner は、`POST /api/logs/cleanup` から呼び出された場合に、archive 済みファイルも `log_retention_days` の削除対象に含める。archive owner は `POST /api/logs/archive` へ `archived_count`、`POST /api/logs/cleanup` へ `deleted_count` と `failed_count` を返す。HTTP response body の形式は `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e を正とする。

**archive / cleanup 固定契約：**

| 項目 | 仕様 |
|------|------|
| archive 対象判定 | build log JSON の `finished_at` を基準にする。欠落時は file mtime を使わず対象外。 |
| archive id | file 名 `{id}.json` の id と JSON 内 `id` が一致する場合だけ対象。 |
| gzip path | `.build_logs/archive/{id}.json.gz`。既存 archive がある場合は上書きせず skip する。 |
| cleanup 順 | 通常 log 削除 → archive log 削除 → 空 archive directory 削除試行。 |
| 削除失敗 | 処理継続し、処理結果に `failed_count` を含める。 |
| 処理結果 | archive は `archived_count`、cleanup は `deleted_count` と `failed_count` を API へ返す。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 対象ログあり | `.json.gz` 作成、元 `.json` 削除、API 参照可。 |
| 実行中ログ | archive しない。 |
| 破損ログ | archive しない、WARN、処理継続。 |
| archive API | 件数を返し、disk usage に archive bytes を含める。 |
| cleanup | 通常 log と archive log の両方を保持期間で削除する。 |
| archive 既存 | 上書きせず skip。 |
| cleanup 一部失敗 | `failed_count` に計上し処理継続。 |

### 27.15 ビルドアーティファクト管理

本機能の目的は、`.snapshots/` に保存された build artifact を API、SDK、UI から一覧、download、削除、rollback できるようにすることである。

owner component は `archive` とする。collaborator component は `api`、`sdk`、`ui`、`runner`、`statefile` とする。snapshot 作成は `runner` の §14b を正とする。

archive owner は snapshot の保存形式、一覧読取、download tar.gz 生成、delete 実体処理、rollback 転送実体処理を担当する。API は下表 endpoint の request / response と archive owner 呼び出し境界だけを担当する。SDK は API method 呼び出し、UI は操作表示と disabled 判定だけを担当する。runner の通常 build 実行、通常 snapshot 作成タイミング、build history / status finalizer の共通処理は runner owner を正とし、本節へ重複定義しない。

**API 呼び出し境界：**

| API | 処理 |
|-----|------|
| `GET /api/snapshots` | archive owner の snapshot 一覧読取を呼び出し、response schema は API 仕様に従う。 |
| `GET /api/snapshots/{id}/download` | archive owner の download tar.gz 生成を呼び出し、streaming response を返す。 |
| `DELETE /api/snapshots/{id}` | archive owner の snapshot delete を呼び出し、成功後に `.config_log` 追記を行う。 |
| `POST /api/history/{id}/rollback` | archive owner の rollback 転送を呼び出し、rollback build log/history の作成境界を runner 共通処理と整合させる。 |

`id` は build id と一致するものだけ許可する。snapshot 専用 id は採番しない。`/`、`..`、空文字、URL decode 後に path separator を含む値は `422` とする。

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
| mtime | snapshot 内 file の mtime を使用してよい。存在しない場合は build log の `finished_at`。 |
| secret 除外 | `.github_token`、`.admin_credentials`、`.api_tokens`、`.smtp_secret`、`.webhook_secret`、runner 状態ファイル名は検出時点で `500` とし、download を中止する。 |

**Rollback 仕様：**

rollback は新しい build id を採番し、`.build_history.trigger="rollback"`、`rollback_from=<元id>` を保存する。元 snapshot は変更しない。rollback 中に別 build が running の場合は `409` とする。転送失敗時は rollback build log を `failure` とし、元 snapshot は削除しない。

rollback は snapshot 内の成果物を deploy target へ再転送する操作であり、以下を行ってはならない。

| 禁止対象 | 理由 |
|----------|------|
| `.last_sha` 更新 | rollback は監視対象 SHA の処理完了ではない。 |
| `.server_config`、`.branch_config`、`.notify_config` の復元 | 設定 rollback ではない。 |
| `.build_history` の過去行書き換え | rollback は新規履歴として追記する。 |
| `.build_logs/{元id}.json` の変更 | 元 build の証跡を保持する。 |
| 元 snapshot の削除または上書き | rollback 成否に関係なく元成果物を保持する。 |
| secret / token / credentials の復元 | snapshot に secret を含めないため復元対象外。 |

rollback build log は `target_status="success"` または `failure_build` とし、`trigger="rollback"`、`rollback_from=<元id>`、`snapshot_id=<元id>` を含める。rollback 転送で pending が発生した場合は `success_deploy_pending` とし、`.pending_transfers` に rollback 用 entry を追加する。

rollback 開始時は `.build_lock` を取得し、取得できない場合は `409 {"error":"Build is running"}` を返す。`.build_lock` 取得後に `.build_state.running=true`、`current_build_id=<new_id>` を保存し、転送完了後に finalizer で `running=false` とする。rollback は queue に積まない。

**snapshot 一覧・削除固定契約：**

| 項目 | 仕様 |
|------|------|
| 一覧対象 | `.snapshots/{id}/meta.json` が存在する directory だけ。 |
| size | directory 配下の通常ファイル size 合計。symlink は size 集計前に異常扱い。 |
| delete 順 | id validation → running check → snapshot directory 確認 → delete → `.config_log` 追記 → response。 |
| delete log 失敗 | snapshot 削除済みのまま `500`。削除は巻き戻さない。 |
| rollback pending | pending entry には `rollback_from`、`snapshot_id`、deploy target を保存する。 |

**SDK / UI 操作境界：**

SDK は `getSnapshots()`、`downloadSnapshot(id)`、`deleteSnapshot(id)`、`rollbackHistory(id)` を提供する。SDK は snapshot の存在、download 安全性、rollback 可否を状態ファイルから推測せず、API response / error をそのまま扱う。UI は snapshot 一覧に id、saved_at、size_bytes、download、delete、rollback 操作を表示する。delete と rollback は実行中 build がある場合 disabled とする。UI は snapshot directory、tar.gz、rollback state を直接操作してはならない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 一覧 | snapshot id、build id、保存日時、size が返る。 |
| download | tar.gz を返し、snapshot 外のファイルを含まない。 |
| delete | 対象 id だけ削除、config log 追記。 |
| rollback 成功 | 新規 build id、trigger rollback、rollback_from 保存。 |
| 不正 id | `422`、状態差分なし。 |
| download symlink | symlink entry を含めず、secret 名検出時は `500`。 |
| rollback pending | 新規 rollback log/history、pending entry、元 snapshot 維持。 |
| rollback running | `409`、状態差分なし。 |
| delete log failure | snapshot は削除済み、response は `500`。 |
