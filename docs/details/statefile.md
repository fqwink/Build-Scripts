# Adlaire CI — Statefile 詳細仕様

本ファイルは `statefile` owner component の詳細仕様正本である。

本ファイルの詳細仕様ファイル管理条件は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 に従う。本ファイルは `statefile` owner component の主本文であり、collaborator component の仕様は読み書き境界、業務処理、表示、security、fixture、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `statefile` |
| collaborator component | `runner`、`api`、`archive`、`security` |
| 持つ内容 | `statefile` owner が主本文として定義する状態ファイル共通仕様、lock、atomic write、JSON Lines、破損時処理、状態読取 adapter、主要 schema。 |
| 持たない内容 | API endpoint の request / response、runner の業務処理、SDK method 実装、UI 表示判断、setup / release 手順、fixture / PR 証跡正本、個別 component の業務判断。 |

---

## 対象範囲

| 範囲 | 内容 |
|------|------|
| §22.0a | 状態ファイル共通仕様、更新手順、schema 厳格化、状態読取 adapter、状態読取 priority。 |
| §22.0c | 主要状態ファイル schema。 |

### 22.0a 状態ファイル共通仕様

`api` および拡張後 `runner` が読み書きする状態ファイルは、下表の初期値、形式、更新責務に従う。表にない状態ファイルを追加してはならない。追加が必要な場合は、先に本節へパス、形式、初期値、更新責務、破損時の扱いを追記する。

| パス | 形式 | 初期値 | 更新責務 | 破損時の扱い |
|------|------|--------|----------|--------------|
| `.admin_credentials` | JSON object | `--init-credentials` で生成 | `api` | 起動時に ERROR ログを出し、HTTP サーバーを起動しない。 |
| `.totp_secret` | JSON object | `{"enabled":false,"secret_base32":null,"confirmed_at":null,"last_accepted_step":null}` | `api` | 読み込み不能時は TOTP 有効 login を `500` で拒否する。破損時は退避するが自動再生成で認証を弱めてはならない。 |
| `.audit_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。追記不能時は対象操作を失敗扱いにする。 |
| `.api_rate_state` | JSON object | `{"windows":{}}` | `api` | `.api_rate_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 window で再生成する。 |
| `.server_config` | JSON object | `{}` | `api` | `.server_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 object で再生成する。 |
| `.notify_config` | JSON object | `{"webhooks":[],"channels":[],"on":[],"summary":{"enabled":false,"interval":"weekly","hour":9,"day_of_week":1},"email":{"enabled":false,"to":[],"on":[]}}` | `runner` / `api` | `.notify_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.notify_log` | JSON Lines | 空ファイル | `runner` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.notify_pending` | JSON array | `[]` | `runner` | `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.pending_transfers` | JSON array | `[]` | `runner` | `.pending_transfers.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.build_history` | JSON Lines | 空ファイル | `runner` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_logs/{id}.json` | JSON object | ビルドごとに新規作成 | `runner` | 対象 ID の API は `500` を返し、既存ファイルは上書きしない。 |
| `.build_lock` | text | 不在 | `runner` | 内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。PID が存在しない場合は stale lock として削除し、存在する場合は `409` 相当の実行中として扱う。形式不正または PID 判定不能の場合は上書きせず `409` を返す。 |
| `.last_sha` / `BranchTarget.SHAFile` | JSON object | `{"sha":""}` | `runner` | JSON 破損、object 以外、`sha` key 不在、`sha` 型不一致は当該 target の decode failure とし、成功時まで更新しない。 |
| `.branch_config` | JSON object | 不在 | `runner` / `api` | `.branch_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、再生成せず `BRANCH_TARGETS` デフォルトへフォールバックする。 |
| `.build_state` | JSON object | `{"running":false,"current_build_id":null,"queued":[],"last_started_at":null,"last_finished_at":null,"weekly_summary_last_sent_at":null,"weekly_summary_sent_date":null}` | `runner` / `api` | `.build_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_status.json` | JSON object | `{"schema_version":1,"updated_at":null,"status":"none","running":false,"current_build_id":null,"last_build_id":null,"last_trigger":null,"last_target_status":null,"last_branch":null,"last_target_file":null,"last_blob_sha":null,"last_commit_sha":null,"last_started_at":null,"last_finished_at":null,"last_duration_seconds":null,"last_error":null,"last_deploy_status":null,"pending_transfers_count":0,"notify_pending_count":0,"circuit_open":false,"circuit_consecutive_failures":0,"output_sha256":null,"size_warn":false}` | `runner` | `.build_status.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_circuit_state` | JSON object | `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` | `runner` / `api` | `.build_circuit_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.local_watch_state.json` | JSON object | `{"files":{}}` | `runner` | `.local_watch_state.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、full build 後に再生成する。 |
| `.build_cache.json` | JSON object | `{"schema_version":1,"entries":{}}` | `builder` | 破損時は WARN を出し、cache miss として扱い、成功後に再生成する。 |
| `.build_cache/pages/` | directory | 空ディレクトリ | `builder` | entry 不一致または読み取り不能 file は miss とし、他 entry は継続使用する。 |
| `.dependency_manifest.json` | JSON object | `{"pages":{}}` | `builder` / `runner` | 破損時は full build とし、成功後に再生成する。 |
| `.approval_queue` | JSON Lines | 空ファイル | `runner` / `api` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_trends.json` | JSON object | `{"schema_version":1,"samples":[],"summary":{"count":0,"avg_seconds":null,"median_seconds":null,"p95_seconds":null}}` | `runner` / `api` | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.build_history` から再集計する。 |
| `.build_chain_config` | JSON object | `{"chains":[]}` | `runner` / `api` | `.build_chain_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、chain 無効として通常 build のみ継続する。 |
| `.repo_config` | JSON object | `{}` | `api` | `.repo_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、スクリプト定数へフォールバックする。 |
| `.config_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.api_access_log` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。秘密情報は記録しない。 |
| `.webhook_secret` | text | 不在 | `api` | 読み込み不能時は Webhook 受信を `501` で拒否する。 |
| `.webhook_events.json` | JSON Lines | 空ファイル | `api` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_control` | JSON object | `{"allow":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.hooks` | JSON object | `{"hooks":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.maintenance` | JSON object | `{"enabled":false,"reason":null,"since":null}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.api_tokens` | JSON object | `{"tokens":[]}` | `api` | `500` を返し、自動再生成しない。token 管理情報の消失による意図しない再許可を防ぐため、破損ファイルは上書きしない。 |
| `.alert_rules` | JSON object | `{"rules":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.tag_rules` | JSON object | `{"rules":[]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.pipeline_config` | JSON object | `{"extra_args":[],"env":{}}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.notes` | UTF-8 text | 空文字列 | `api` | 読み込み不能時は `500` を返し、自動上書きしない。 |
| `.smtp_config` | JSON object | SMTP 未設定値 | `api` | 初期値で再生成し、ERROR ログを記録する。 |
| `.smtp_secret` | text | 不在 | `api` | 読み込み不能時は SMTP 送信を `422` で拒否する。 |
| `.dashboard_layout` | JSON object | `{"widgets":["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]}` | `api` | 初期値で再生成し、ERROR ログを記録する。 |

`.build_logs/archive/` は gzip 圧縮済み build log の保存先ディレクトリである。初期値は空ディレクトリとし、`runner` または `POST /api/logs/archive` が archive 書き込み前に存在確認し、不在の場合だけ作成する。圧縮済みファイル名は `{id}.json.gz` 固定とし、通常 `.build_logs/{id}.json` と同じ build id を表す。

JSON Lines ファイルは、1 行につき 1 JSON object とする。追記時は末尾に改行を必ず付ける。秘密情報を含む可能性のある `.admin_credentials`、`.totp_secret`、`.github_token`、`.webhook_secret`、`.smtp_secret` は mode `600` を必須とする。

**状態ファイル更新手順：**

1. 対象ファイルの `{name}.lock` を `O_CREATE|O_EXCL` で作成する。
2. ロック取得に失敗した場合は 100ms 間隔で最大 10 秒待つ。
3. 現在値を読み込み、schema と入力値を検証する。
4. 更新後 JSON を `{name}.tmp.{pid}` に UTF-8 / LF で書き出す。
5. ファイルを close し、通常状態ファイルは `0644`、秘密情報ファイルと lock file は `0600` に chmod する。
6. `os.Rename(tmp, target)` で置換する。
7. target file を open して `Sync` し、続けて親ディレクトリを open して `Sync` する。
8. ロックファイルを削除する。

手順 3〜7 の途中で失敗した場合は target を変更せず、tmp を削除し、ロックを削除して `500 Internal Server Error` を返す。`os.Rename` 後の `Sync` に失敗した場合は target を維持し、ERROR ログと `.config_log` へ失敗を記録して `500` を返す。複数ファイル更新 API は §22.0d の Write 列順にこの手順を実行し、途中失敗時は未処理ファイルを書き込まない。既に書き込んだファイルの自動ロールバックは行わず、`.config_log` に失敗内容を記録する。

**状態ファイル schema 厳格化契約：**

状態ファイルの読込、正規化、保存は以下に固定する。§22.0c または個別機能節で例外を明記していない限り、実装者判断で旧形式、未知 key、null、欠落配列を成功扱いにしてはならない。

| 対象 | 読込時 | 保存時 | 失敗時 |
|------|--------|--------|--------|
| 未知 key | JSON object に schema 未定義 key がある場合は破損扱いとする。例外は §22.0c で明記した旧形式正規化だけ。 | 未知 key を保存しない。既存未知 key を黙って削除して保存しない。 | read endpoint は `500 {"error":"State file is corrupted"}`。write endpoint は target を変更しない。 |
| 必須 key 不足 | 個別節に「欠落時に適用する既定値」と「保存するか読み取り時だけか」が明記されていない場合は破損扱いとする。 | 必須 key はすべて明示保存する。 | 初期値再生成が §22.0a 表で指定されたファイルだけ再生成する。 |
| `null` | 型欄が `string/null`、`object/null`、`integer/null` 等で明示した key だけ許可する。 | nullable でない key に `null` を保存しない。 | validation error または破損扱い。 |
| 配列 | `[]` を既定値とする key は read adapter の戻り値で空配列を返す。 | 保存 API は配列 key を省略せず、空の場合も `[]` を明示する。 | 型不一致は `422` または `500`。 |
| 数値 | 整数 key は JSON number の整数だけ許可する。小数、指数表記由来の非整数、文字列数値は拒否する。 | 整数は JSON number として保存する。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| 時刻 | UTC ISO 8601 秒精度 `Z` だけ許可する。 | 保存前に UTC 秒精度へ丸める。ミリ秒、local timezone、offset 付き文字列を保存しない。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| mode | 秘密情報ファイルは `0600`、通常 JSON / JSON Lines は `0644`、directory は `0755` を標準とする。 | chmod 失敗時は成功扱いにしない。 | chmod 失敗は `500`。target を更新した後の chmod 失敗は ERROR ログに残す。 |
| 改行 | text / JSON / JSON Lines は LF で保存する。JSON object / array ファイルは末尾 LF 1 個を付ける。 | CRLF、BOM、末尾余分空白を新規保存しない。 | 入力 text が CRLF を含む場合の扱いは個別機能節に従う。 |

旧 schema からの正規化は、本ファイルに「旧 key」「変換後 key」「削除する key」「保存するか読み取り時だけか」を明記した場合だけ実装する。明記がない旧形式は破損扱いとし、黙って推測変換してはならない。

**API 状態読取アダプタ固定契約：**

`api` は、Phase 3 endpoint の状態読取を下表の adapter 名と戻り値で実装する。各 adapter は Go 内部関数名として固定し、同じ状態ファイルを endpoint ごとに別ロジックで直接 parse してはならない。

| Adapter | 読取対象 | 正常戻り値 | 不在時 | 破損時 / 読込不能時 |
|---------|----------|------------|--------|---------------------|
| `readBuildStatus()` | `.build_status.json` | `BuildStatus` | `(nil, false, nil)` を返し、fallback 判定へ渡す。 | `(nil, true, ErrStateCorrupted)` または `ErrStateReadFailed`。 |
| `readBuildState()` | `.build_state` | `BuildState` | §22.0a の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readBuildHistory()` | `.build_history` | `[]BuildHistoryEntry` | 空配列を返す。 | 行単位破損は除外し、ファイル読込不能だけ `ErrStateReadFailed`。 |
| `readBuildLog(id)` | `.build_logs/{id}.json`、`.build_logs/archive/{id}.json.gz` | `BuildLog` | 通常ログ不在時は archive を読む。両方不在は `ErrNotFound`。 | 対象 ID の通常ログまたは archive が破損している場合は `ErrStateCorrupted`。 |
| `readLatestBuildLogs(n,q)` | `.build_logs/`、`.build_logs/archive/` | `[]LogLine` | 空配列を返す。 | 個別ログ破損は除外し、`LOG_SKIP_CORRUPT` を server log へ記録する。ディレクトリ読込不能は `ErrStateReadFailed`。 |
| `readPendingTransfers()` | `.pending_transfers` | `[]PendingTransfer` | 空配列を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readCircuitState()` | `.build_circuit_state` | `BuildCircuitState` | §22.0a の初期値を返す。 | `ErrStateCorrupted` または `ErrStateReadFailed`。 |
| `readBuildLock()` | `.build_lock` | `BuildLockState` | `running=false` を返す。 | 形式不正、PID 判定不能、OS 判定失敗は `running=true, stale=false, valid=false` として返し、build command は `409`。 |

`GET` endpoint は上表の adapter を read-only で呼び出し、状態ファイルの作成、削除、退避、chmod、正規化、再生成、破損行の除去書き戻しを行ってはならない。`GET` endpoint が `{name}.lock` を検出しても、`.build_lock` 以外の lock file は待機条件やエラー条件にせず、rename 済み target をそのまま読む。write endpoint は §22.0a の状態ファイル更新手順に従う。

**API 状態読取 priority：**

| Endpoint | 読取順 | 正常時 response 算出 | 不在時 | 破損時 / 読込不能時 |
|----------|--------|----------------------|--------|---------------------|
| `GET /api/status` | `readBuildStatus()` → 不在時だけ `readBuildHistory()`、`readBuildState()`、`readBuildLock()`、`readPendingTransfers()`、`readCircuitState()` | `.build_status.json` がある場合は同ファイルを正とし、`running` だけ `.build_lock` が valid running の場合に `true` へ上書きする。 | `.build_status.json` 不在時は fallback で `status`、`last_*`、`running`、`pending_transfers_count`、`circuit_*` を算出する。履歴なしは `status:"none"`。 | `.build_status.json` 破損は `500 {"error":"State file is corrupted"}`。fallback 中の必須読取破損も `500`。 |
| `GET /api/health` | `readBuildStatus()`、`readBuildLock()`、出力サイト確認 | API process、出力サイト、直近 build 状態を `ok` / `warn` / `error` で返す。 | `.build_status.json` 不在は `status:"degraded"`。 | `.build_status.json` 破損は例外的に `200` とし、`status:"degraded"`、該当 item を `error` にする。 |
| `GET /api/history` | `readBuildHistory()` | 有効行だけを新しい順に sort し、query filter 後に paging する。 | `total:0`、`pages:0`、`history:[]`。 | 行単位破損は除外し、`BUILD_HISTORY_SKIP_CORRUPT` を server log へ記録する。ファイル読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/history/{id}/log` | `readBuildLog(id)` | 対象 ID の log object を返す。 | 通常ログと archive の両方が不在なら `404 {"error":"Not found"}`。 | 対象 ID の log 破損は `500 {"error":"State file is corrupted"}`。 |
| `GET /api/logs` | `readLatestBuildLogs(n,q)` | 最新 log line を時系列順へ正規化し、`q` 指定時は部分一致で絞り込む。 | `{"lines":[]}`。 | 個別 log 破損は除外する。ディレクトリ読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/queue` | `readBuildState()`、queue 上限算出時に `.server_config` | `queued`、`running`、`queue_max_size` を返す。 | `.build_state` 不在は `queued:[]`、`running:false`。`.server_config` 不在は既定 queue 上限。 | `.build_state` 破損は `500 {"error":"State file is corrupted"}`。 |
| `POST /api/build` / `POST /api/build/force` | `readCircuitState()` → `readBuildLock()` → `readBuildState()` | circuit closed かつ lock 非実行なら `.build_state` を更新し、build 開始または queue 追加を返す。 | `.build_state` / `.build_circuit_state` 不在は初期値。 | circuit / state 破損は `500`。`.build_lock` が valid running、形式不正、PID 判定不能の場合は `409 {"error":"Conflict"}`。 |
| `POST /api/cancel` | `readBuildLock()`、`readBuildState()` | 実行中 build を cancel request 状態へ更新する。 | lock 不在かつ running false は `409 {"error":"Conflict"}`。 | `.build_state` 破損は `500`。`.build_lock` 形式不正または PID 判定不能は `409`。 |
| `POST /api/circuit-breaker/reset` | `readCircuitState()` | `open:false`、`consecutive_failures:0`、`opened_at:null`、`last_error:null` を atomic write する。 | 不在は初期値から reset 後値を書き込む。 | 読取破損は `500 {"error":"State file is corrupted"}` とし、上書きしない。 |

`ErrStateCorrupted` は JSON parse 失敗、schema_version 不一致、必須 key 不足、型不一致、列挙値不一致、UTC 時刻形式不一致のいずれかで返す。`ErrStateReadFailed` は permission denied、通常ファイルではない path、gzip 読込失敗、I/O error で返す。API response body はそれぞれ `{"error":"State file is corrupted"}`、`{"error":"State file read failed"}` 固定とし、path、Go error、ファイル内容を含めない。

JSON Lines adapter は空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致の行を壊れた行として除外する。除外後に sort、filter、paging、`total`、`pages` を算出する。壊れた行の存在は response body に含めず、server log に固定コード、path、1 始まりの line number だけを記録する。

`.build_lock` の PID が存在しない場合、`readBuildLock()` は `running=false, stale=true, valid=true` を返す。read-only endpoint は stale lock を削除しない。build command は開始前に `.build_lock` を再読込し、同じ stale 判定なら `.build_lock` だけを削除してから新規 lock を作成する。削除失敗時は `409 {"error":"Conflict"}` とし、`.build_state` を変更しない。

### 22.0c 主要状態ファイル schema

本節の schema は、API 実装、SDK 型、標準管理ツール表示、バックアップ/リストアの基準である。ここに定義したキー以外を保存してはならない。追加キーを追加する場合は、型、既定値、読み書き API、既存データの扱いを本節へ追記してから実装する。

**`.server_config` schema：**

| キー | 型 | 既定値 | 許容値 | 読み書き API | 説明 |
|------|----|--------|--------|--------------|------|
| `log_max_lines` | integer | `500` | 1〜10000 | `GET/POST /api/config` | `GET /api/logs` が返す最大行数。 |
| `history_max_count` | integer | `100` | 1〜10000 | `GET/POST /api/config` | `.build_history` の通常表示上限。削除処理の上限ではない。 |
| `build_timeout_seconds` | integer | `300` | 1〜86400 | `GET/POST /api/config` | 手動/自動ビルドのタイムアウト秒数。 |
| `log_retention_days` | integer | `30` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/cleanup` | `0` は自動削除なし。 |
| `log_level` | string | `"INFO"` | `"INFO"` / `"DEBUG"` / `"WARNING"` / `"ERROR"` | `GET/POST /api/config`, `POST /api/log-level` | `api` のランタイムログレベル。 |
| `pat_expires_at` | string/null | `null` | `YYYY-MM-DD` または `null` | `GET/POST /api/config` | PAT 期限表示・診断用。 |
| `snapshots_keep` | integer | `5` | 0〜100 | `GET/POST /api/config` | `0` はスナップショット保存無効。 |
| `queue_max_size` | integer | `3` | 0〜100 | `GET/POST /api/config`, `GET /api/queue` | `0` はキュー無効。 |
| `build_retry_max` | integer | `0` | 0〜10 | `GET/POST /api/config` | ビルド失敗時の自動リトライ最大回数。`0` は無効。 |
| `build_retry_base_seconds` | integer | `5` | 1〜3600 | `GET/POST /api/config` | 自動リトライ backoff 基底秒数。待機秒数は `base * attempt` とする。 |
| `commit_status_enabled` | boolean | `false` | `true` / `false` | `GET/POST /api/config` | GitHub Commit Status API 送信の有効/無効。 |
| `commit_status_context` | string | `"Adlaire CI"` | 1〜100 文字 | `GET/POST /api/config` | GitHub commit status の `context`。 |
| `commit_status_target_url` | string/null | `null` | `http://` または `https://` の URL、または `null` | `GET/POST /api/config` | Commit Status の `target_url`。`null` の場合は送信 payload から省略する。 |
| `log_archive_after_days` | integer | `0` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/archive` | `0` は archive 無効。指定日数より古い通常 build log を gzip 圧縮する。 |
| `build_trend_keep_count` | integer | `1000` | 10〜10000 | `GET/POST /api/config` | `.build_trends.json` に保持する trend sample 件数。 |
| `duration_anomaly` | object | `{"enabled":false,"min_samples":20,"avg_multiplier":2.0,"p95_multiplier":1.5}` | §27.38 | `GET/POST /api/config` | build 所要時間異常検知の設定。 |
| `force_build_interval_hours` | integer | `0` | 0〜8760 | `POST /api/schedule/force-interval`, `GET /api/schedule` | `0` は強制再ビルド無効。 |
| `build_cooldown_seconds` | integer | `0` | 0〜86400 | `POST /api/schedule/cooldown`, `GET /api/schedule` | `0` はクールダウン無効。 |
| `schedule_interval_seconds` | integer | `300` | 30〜86400 | `POST /api/schedule/interval`, `GET /api/schedule` | systemd timer 更新値。 |
| `schedule_paused` | boolean | `false` | `true` / `false` | `POST /api/schedule/pause`, `POST /api/schedule/resume`, `GET /api/schedule` | 自動ポーリング停止状態。 |
| `allowed_hours` | object/null | `null` | `{"from":0〜23,"to":0〜23}` または `null` | `POST /api/schedule/allowed-hours`, `GET /api/schedule` | UTC の自動ビルド許可時間帯。 |
| `session_timeout_seconds` | integer | `28800` | 300〜2592000 | `GET/POST /api/config` | 新規 session の有効期限秒数。既存 session の `expires_at` は変更しない。 |
| `api_rate_limit` | object | `{"enabled":true,"groups":{"login":{"window_seconds":60,"max_requests":10},"read":{"window_seconds":60,"max_requests":600},"trigger":{"window_seconds":60,"max_requests":60},"operate":{"window_seconds":60,"max_requests":120},"config":{"window_seconds":60,"max_requests":60},"admin":{"window_seconds":60,"max_requests":60}}}` | §27.47 | `GET /api/api-rate-limit`, `POST /api/api-rate-limit`, `GET/POST /api/config` | API rate limit の endpoint group 別固定窓設定。 |

`.server_config` の `POST /api/config` では `force_build_interval_hours`、`build_cooldown_seconds`、`schedule_interval_seconds`、`schedule_paused`、`allowed_hours` を直接更新してはならない。これらは専用スケジュール API からのみ更新する。

**`.notify_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `webhooks` | object[] | `[]` | 下記 Webhook object | 互換通知先一覧。`channels` が空の場合、runner は `webhooks` を webhook channel として扱う。 |
| `channels` | object[] | `[]` | 下記 Channel object | 統一通知 channel 一覧。`channels` が存在する場合、runner は `channels` を優先し、`webhooks` / `email` は互換表示用として扱う。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"` | 通知イベント。重複は除去する。 |
| `summary` | object | 下記 Summary object | 下記 | 定期サマリー設定。 |
| `email` | object | 下記 Email object | 下記 | メール通知設定。SMTP 詳細は `.smtp_config` / `.smtp_secret` を正とする。 |

Channel object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `id` | string | 自動採番 | `n` + 数字、または 1〜64 文字の英数字 `_` `-` | channel 識別子。 |
| `type` | string | 必須 | `"webhook"` / `"email"` / `"command"` | 送信方式。 |
| `label` | string | `""` | 0〜64 文字 | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の channel へは送信しない。 |
| `on` | string[] | `[]` | top-level `on` と同じ、または `"*"` | 空配列の場合は top-level `on` に従う。 |
| `config` | object | `{}` | type 別 schema | webhook `url`、email `to`、command `command_args`。 |
| `retry_count` | integer | `2` | 0〜10 | retry 対象失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | 1〜3600 | 再試行間隔。 |

Webhook object:

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `url` | string | 必須 | URL 検証に従う | 送信先 URL。 |
| `label` | string | `""` | 0〜64 文字 | 管理画面表示名。 |
| `enabled` | boolean | `true` | boolean | `false` の宛先へは送信しない。 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"deploy_failure"`, `"weekly_summary"`, `"approval_required"`, `"duration_anomaly"`, `"config_corrupt"`, `"*"` | この宛先が受け取るイベント。空配列の場合は top-level `on` に従う。 |
| `payload_template` | string/null | `null` | 0〜10000 文字または `null` | `null` は標準 payload。 |
| `retry_count` | integer | `2` | 0〜10 | 送信失敗時の追加試行回数。 |
| `retry_interval_seconds` | integer | `30` | 1〜3600 | 再試行間隔。 |
| `secret` | string/null | `null` | 1〜256 文字または `null` | 保存時は平文保存可。ただし GET/backup では `"***"` へマスクする。 |

Summary object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `interval` | string | `"weekly"` | `"weekly"` 固定 |
| `hour` | integer | `9` | 0〜23 |
| `day_of_week` | integer | `1` | 0〜6 |

Email object:

| キー | 型 | 既定値 | 許容値 |
|------|----|--------|--------|
| `enabled` | boolean | `false` | boolean |
| `to` | string[] | `[]` | メールアドレス配列、最大 50 件 |
| `on` | string[] | `[]` | `"start"`, `"success"`, `"failure"`, `"duration_anomaly"` |

**`.branch_config` schema：**

```json
{
  "branch_targets": [
    {
      "branch": "main",
      "target_file": "docs",
      "sha_file": "/opt/adlaire-builder/.last_sha",
      "src": "/opt/adlaire-builder/repo/docs",
      "out": "/opt/adlaire-builder/dist/site",
      "deploy_targets": [
        { "host": "192.0.2.1", "user": "deploy", "dest_dir": "/var/www/html/" }
      ]
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `branch_targets` | object[] | 必須 | 0〜50 件 | 永続ファイルの正本 key。空配列は `.branch_config` 削除と同義。 |
| `branch` | string | 必須 | 1〜128 文字、`refs/heads/` は含めない | GitHub branch 名。 |
| `target_file` | string | 必須 | 相対パス、`..` 禁止 | GitHub リポジトリ内の監視対象ファイル。 |
| `sha_file` | string | 必須 | 絶対パス | 対象 branch/file の SHA キャッシュ。 |
| `src` | string | 必須 | 絶対パス | blob 本文の書き出し先。 |
| `out` | string | 必須 | 絶対パス | ビルド成果物パス。 |
| `deploy_targets` | object[] | 必須 | 0〜20 件 | SSH 転送先。空配列は転送なし。 |
| `deploy_targets[].host` | string | 必須 | 1〜255 文字 | SSH host。 |
| `deploy_targets[].user` | string | 必須 | 1〜64 文字 | SSH user。 |
| `deploy_targets[].dest_dir` | string | 必須 | 絶対パス | 転送先ディレクトリ。 |

`.branch_config` が不在の場合、`GET /api/branch-config` は `source: "default"` と `BRANCH_TARGETS` の定数値を返す。`.branch_config` が存在する場合、`source: "file"` とファイル内容を返す。API request / response の表示名として `branches` を使う場合でも、永続ファイルへ保存する key は必ず `branch_targets` とする。API は `branches: []` または `branch_targets: []` を `.branch_config` の空配列保存として扱ってはならない。`POST /api/branch-config` で空配列を受け取った場合は `.branch_config` を削除し、default 復帰として扱う。

**`.last_sha` / `BranchTarget.SHAFile` schema：**

```json
{"sha":""}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `sha` | string | 必須 | 空文字または Git blob SHA | 空文字は初回実行扱い。 |

SHA cache は target ごとの処理済み Git blob SHA を保存する JSON object とする。legacy text 形式は自動変換しない。

| 状態 | 読込時の扱い | 保存時の扱い |
|------|--------------|--------------|
| 不在 | 初回実行として `previous_blob_sha=""` を扱う。 | build 成功時に `{"sha":"<new_sha>"}` を atomic write する。 |
| `{"sha":""}` | 初回実行として扱う。 | build 成功時に新 SHA で上書きする。 |
| `{"sha":"<value>"}` | `<value>` を前回 SHA として比較する。 | build 成功時に新 SHA で上書きする。 |
| JSON 破損 | 当該 target の decode failure として扱う。 | 更新しない。 |
| object 以外 | JSON 破損と同じ扱い。 | 更新しない。 |
| `sha` key 不在または string 以外 | JSON 破損と同じ扱い。 | 更新しない。 |
| 未知 key あり | 前回 SHA として `sha` だけを読む。 | build 成功時に `sha` だけの object で上書きする。 |

SHA cache の更新タイミング、skip / failure 時の更新可否、複数 target 時の個別更新は [`docs/details/runner.md`](runner.md) §13 の SHA cache 読み書き契約を正とする。

**`.repo_config` schema：**

| キー | 型 | 既定値 | 許容値 | 説明 |
|------|----|--------|--------|------|
| `owner` | string | スクリプト定数 `OWNER` | 1〜100 文字 | GitHub owner。 |
| `repo` | string | スクリプト定数 `REPO` | 1〜100 文字 | GitHub repository。 |
| `branch` | string | スクリプト定数 `BRANCH` | 1〜128 文字 | 単一ターゲット用 branch。 |
| `target_file` | string | スクリプト定数 `TARGET_FILE` | 相対パス、`..` 禁止 | 単一ターゲット用監視ファイル。 |
| `updated_at` | string | 更新時刻 | ISO 8601 | 最終更新日時。 |

`POST /api/repo-config` は指定されたキーのみ更新する。未指定キーは既存値を保持する。全キーが未指定の場合は `422` を返す。

**`.totp_secret` schema：**

```json
{
  "enabled": true,
  "secret_base32": "JBSWY3DPEHPK3PXP",
  "confirmed_at": "2026-09-15T10:00:00Z",
  "last_accepted_step": 59652320
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | TOTP 有効状態。 |
| `secret_base32` | string/null | 必須 | RFC 4648 base32、padding なし、16〜64 文字、または `null` | TOTP secret。API response、log、backup へ平文出力しない。 |
| `confirmed_at` | string/null | 必須 | ISO 8601 または `null` | TOTP 有効化完了日時。 |
| `last_accepted_step` | integer/null | 必須 | Unix time 30 秒 step または `null` | 同一 code 再利用防止。 |

**`.admin_credentials` schema：**

```json
{
  "password_hash": "<sha256_iter_v1_hex>",
  "salt": "<hex>",
  "algorithm": "sha256_iter_v1",
  "iterations": 260000,
  "must_change": true,
  "login_count": 0,
  "last_login_at": null,
  "updated_at": "2026-09-15T10:00:00Z"
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `password_hash` | string | 必須 | 64 文字 lowercase hex | password 本体は保存しない。 |
| `salt` | string | 必須 | 64 文字 lowercase hex | 32 bytes salt。 |
| `algorithm` | string | 必須 | `"sha256_iter_v1"` 固定 | 他 algorithm は初期実装で拒否する。 |
| `iterations` | integer | 必須 | `260000` 固定 | 値が異なる場合は認証を `500` で拒否する。 |
| `must_change` | boolean | 必須 | boolean | 初期生成時 `true`、パスワード変更後 `false`。 |
| `login_count` | integer | 必須 | 0 以上 | session token 発行成功時だけ +1。TOTP ticket 発行時は増やさない。 |
| `last_login_at` | string/null | 必須 | UTC ISO 8601 または `null` | session token 発行成功時だけ更新する。 |
| `updated_at` | string | 必須 | UTC ISO 8601 | password hash 更新時刻。 |

`.admin_credentials` に未知 key がある場合は credentials 破損として扱い、自動削除しない。必須 key 不足、型不一致、hex 不正、`algorithm` 不一致、`iterations` 不一致もすべて credentials 破損とする。API 起動時検証で credentials 破損を検出した場合は、§22.0a に従って ERROR ログを出し、HTTP サーバーを起動しない。HTTP サーバー稼働中の読込時検証で credentials 破損を検出した場合、`POST /api/login` と `POST /api/change-password` は `500 {"error":"Internal server error"}` を返す。API response、`.access_log`、`.audit_log`、journal に破損内容、hash、salt を出してはならない。

**`.audit_log` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `timestamp` | string | 必須 | ISO 8601 | 発生日時。 |
| `request_id` | string | 必須 | 16 byte hex | API request 単位の識別子。 |
| `actor_type` | string | 必須 | `"admin"` / `"api_token"` / `"system"` / `"anonymous"` | 操作者種別。 |
| `actor_id` | string/null | 必須 | `"admin"`、token id、`"system"`、または `null` | 操作者。secret 本体は保存しない。 |
| `action` | string | 必須 | §27.44 | 操作種別。 |
| `target_type` | string | 必須 | §27.44 | 対象種別。 |
| `target_id` | string/null | 必須 | 対象 id または `null` | 対象識別子。 |
| `result` | string | 必須 | `"success"` / `"failure"` / `"denied"` | 結果。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。 |
| `message` | string/null | 必須 | 0〜500 文字または `null` | 固定文言。secret、token、password は保存しない。 |

**`.api_rate_state` schema：**

```json
{
  "windows": {
    "ip:127.0.0.1:login": {
      "window_start": "2026-09-15T10:00:00Z",
      "count": 1
    }
  }
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `windows` | object | 必須 | key は `{dimension}:{id}:{endpoint_group}` | 固定窓状態。 |
| `window_start` | string | 必須 | ISO 8601 | 現在窓の開始時刻。 |
| `count` | integer | 必須 | 0 以上 | 現在窓内リクエスト数。 |

**`.api_tokens` schema：**

```json
{
  "tokens": [
    {
      "id": "tok000001",
      "label": "監視用",
      "scopes": ["read"],
      "token_hash": "<sha256_hex>",
      "created_at": "2026-09-15T10:00:00Z",
      "last_used_at": null,
      "expires_at": null,
      "revoked_at": null
    }
  ]
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `tokens` | object[] | 必須 | 0〜100 件 | 発行済み API token 一覧。 |
| `id` | string | 必須 | `tok` + 3 桁以上の数字 | token 識別子。 |
| `label` | string | 必須 | 1〜64 文字 | 表示名。 |
| `scopes` | string[] | 必須 | `read`, `trigger`, `operate`, `config`, `admin` の 1〜5 件 | token に許可する scope。 |
| `token_hash` | string | 必須 | SHA-256 hex | token 本体は保存しない。 |
| `created_at` | string | 必須 | ISO 8601 | 作成日時。 |
| `last_used_at` | string/null | 必須 | ISO 8601 または `null` | 最終使用日時。 |
| `expires_at` | string/null | 必須 | ISO 8601 または `null` | 有効期限。`null` は無期限。 |
| `revoked_at` | string/null | 必須 | ISO 8601 または `null` | 失効日時。`null` は有効。 |

`POST /api/tokens` は token 本体を `act_` + 32 byte 相当のランダム文字列として生成し、レスポンス時に 1 回だけ返す。保存する値は `token_hash` のみとする。`DELETE /api/tokens/{id}` は物理削除せず、`revoked_at` を現在時刻へ更新する。旧 `scope` 文字列が存在する場合は読み込み時に `scopes:[scope]` へ正規化して保存し直す。

**`.api_tokens` 実装固定値：**

| 項目 | 仕様 |
|------|------|
| token id 採番 | `tok` + 6 桁連番とする。既存最大番号が `tok000123` の場合、次は `tok000124` とする。連番抽出不能な id は衝突確認対象には含めるが、最大番号算出には使わない。 |
| token 本体 | `crypto/rand` 32 bytes を `encoding/base64.RawURLEncoding` で文字列化し、先頭に `act_` を付ける。保存前 hash は prefix を含む token 全体に対して `sha256` を計算する。 |
| `label` | 1〜64 文字。前後空白は保存前に除去する。除去後が空文字なら `422`。 |
| `scopes` | 1〜5 件。重複は除去し、保存値は `read`, `trigger`, `operate`, `config`, `admin` の順に正規化する。 |
| `expires_at` | `null` または現在時刻より後の UTC ISO 8601。過去または現在時刻は `422`。 |
| 一覧順 | `GET /api/tokens` は `created_at` 降順、同時刻は `id` 昇順で返す。 |
| 失効済み表示 | `GET /api/tokens` は失効済み token も返す。token 本体と `token_hash` は返さない。 |
| 認証時更新 | 有効 token 認証成功時だけ `last_used_at` を現在時刻へ更新する。期限切れ、失効済み、hash 不一致では更新しない。 |
| 破損行相当 | `.api_tokens.tokens` 内の個別 record が schema 不正の場合、認証と一覧は `500` を返し、自動補正しない。旧 `scope` から `scopes` への正規化だけは例外として許可する。 |

**`.maintenance` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `enabled` | boolean | 必須 | boolean | メンテナンス有効状態。 |
| `reason` | string/null | 必須 | 0〜500 文字または `null` | 理由。 |
| `since` | string/null | 必須 | ISO 8601 または `null` | 有効化日時。 |

`enabled: false` の場合、`reason` と `since` は `null` とする。`POST /api/maintenance/enable` は `enabled: true`、`reason`、`since` を同時に保存する。

**`.access_control` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `allow` | string[] | 必須 | IPv4 address または CIDR、0〜100 件 | 空配列は制限なし。保存時は入力順を保持し、重複は除去する。 |

`allow` の各要素は前後空白を除去してから検証する。空文字、IPv6、hostname、URL、CIDR prefix が 0〜32 以外、parse 不能な値は `422` とし、既存 `.access_control` を変更しない。

**`.hooks` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `hooks` | object[] | 必須 | 0〜50 件 | 登録済み hook。 |
| `id` | string | 必須 | §22.0e.2 の hook id | hook 識別子。 |
| `phase` | string | 必須 | `"pre"` / `"post"` | 実行 phase。 |
| `command_args` | string[] | 必須 | 1〜20 件 | shell を介さず `exec.CommandContext` に渡す引数配列。 |
| `enabled` | boolean | 必須 | boolean | `false` の hook は実行しない。 |
| `abort_on_failure` | boolean | 必須 | boolean | `pre` hook 失敗時だけ参照する。 |
| `timeout_seconds` | integer | 必須 | 1〜3600 | hook 単体の timeout。未指定作成時は `300`。 |
| `created_at` | string | 必須 | UTC ISO 8601 | hook 作成日時。 |

`command_args[0]` は 1〜256 文字、`command_args[1:]` の各要素は 1〜500 文字とし、NUL、改行、CR を禁止する。`command_args[0]` は絶対 path または PATH 解決可能なコマンド名に限定する。`phase` と `command_args` が既存 enabled hook と完全一致する場合、`POST /api/hooks` は `409 {"error":"Conflict"}` を返す。

**`.alert_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | dashboard alert rule。 |
| `id` | string | 必須 | §22.0e.2 の alert rule id | rule 識別子。 |
| `metric` | string | 必須 | `success_rate_7d`, `avg_duration_seconds`, `last_build_age_hours`, `disk_usage_bytes` | 評価対象。 |
| `operator` | string | 必須 | `lt`, `gt`, `lte`, `gte` | 比較演算子。 |
| `threshold` | number | 必須 | 0 以上 | 比較値。 |
| `level` | string | 必須 | `info`, `warn`, `error` | alert severity。 |
| `message` | string | 必須 | 1〜200 文字 | UI 表示文。secret を含めない。 |

**`.tag_rules` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `rules` | object[] | 必須 | 0〜100 件 | 自動タグ付け rule。 |
| `id` | string | 必須 | §22.0e.2 の tag rule id | rule 識別子。 |
| `condition` | string | 必須 | §15B の条件式 grammar | 評価条件。 |
| `tags` | string[] | 必須 | 1〜20 件、各 1〜50 文字 | 付与するタグ。重複は除去する。 |

`condition` は `変数 空白 演算子 空白 値` の 1 条件だけを許可する。`&&`、`||`、括弧、関数呼び出し、正規表現、算術式は `422` とする。

**`.pipeline_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `extra_args` | string[] | 必須 | 0〜50 件 | `builder` に渡す追加 CLI 引数。 |
| `env` | object | 必須 | key/value は下記 | builder process に追加する環境変数。 |

`extra_args` は空文字、NUL、改行、CR を禁止し、`--src`、`--out`、`--build-id`、`--commit-sha`、`--build-at`、`--version`、`--help` を指定してはならない。`env` key は `^[A-Z_][A-Z0-9_]{0,63}$`、value は 0〜1000 文字とし、`PATH`、`HOME`、`SHELL`、`USER`、`GITHUB_TOKEN`、`ADLAIRE_TOKEN` は上書き禁止とする。

**`.dashboard_layout` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `widgets` | string[] | 必須 | §16D の widget id、1〜9 件 | 表示 widget 順序。重複禁止。 |

**`.smtp_config` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `host` | string/null | 必須 | hostname または IP、または `null` | SMTP server。 |
| `port` | integer | 必須 | 1〜65535 | SMTP port。 |
| `user` | string/null | 必須 | 0〜255 文字または `null` | SMTP user。 |
| `tls` | boolean | 必須 | boolean | STARTTLS または TLS 使用。 |
| `from` | string/null | 必須 | email address または `null` | From address。 |
| `to` | string[] | 必須 | email address、0〜50 件 | 送信先。 |
| `on` | string[] | 必須 | `start`, `success`, `failure`, `duration_anomaly` | 送信イベント。 |
| `enabled` | boolean | 必須 | boolean | メール通知有効状態。 |

`.smtp_secret` は UTF-8 text とし、末尾 LF なしで password 本体だけを保存する。mode は `0600` 固定。`POST /api/smtp-config` で `password` が未指定の場合、既存 `.smtp_secret` を変更しない。`password:null` は secret 削除を意味し、`.smtp_secret` が存在する場合だけ削除する。

**`.build_state` schema：**

```json
{
  "running": false,
  "current_build_id": null,
  "queued": [],
  "last_started_at": null,
  "last_finished_at": null,
  "weekly_summary_last_sent_at": null,
  "weekly_summary_sent_date": null
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `running` | boolean | 必須 | boolean | ビルド実行中状態。 |
| `current_build_id` | string/null | 必須 | build id または `null` | 実行中 build id。 |
| `queued` | object[] | 必須 | 0〜`queue_max_size` 件 | 待機中 build queue。 |
| `last_started_at` | string/null | 必須 | ISO 8601 または `null` | 最終開始日時。 |
| `last_finished_at` | string/null | 必須 | ISO 8601 または `null` | 最終完了日時。 |
| `weekly_summary_last_sent_at` | string/null | 必須 | ISO 8601 または `null` | 週次サマリー最終送信日時。 |
| `weekly_summary_sent_date` | string/null | 必須 | `YYYY-MM-DD` または `null` | 週次サマリー二重送信防止日付。 |

Queue entry schema:

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | §22.0e.2 の queue id | queue id。 |
| `trigger` | string | 必須 | `"manual"`, `"webhook"`, `"approval"` | 起動種別。 |
| `queued_at` | string | 必須 | UTC ISO 8601 | queue 追加日時。 |
| `requested_by` | string | 必須 | `"admin"`, `"webhook"`, `"approval"`, API token id | queue 追加元。 |
| `priority` | string | 必須 | §27.35 の値 | 優先度。未指定作成時は `"normal"`。 |
| `created_seq` | integer | 必須 | 1 以上 | 既存最大 + 1。 |
| `payload` | object | 必須 | trigger ごとの固定 payload | 不要時は `{}`。 |

Queue entry `payload` は trigger ごとに以下を許可する。未知 key は API 追加時 `422`、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

| trigger | payload |
|---------|---------|
| `manual` | `{ "force": boolean }`。 |
| `webhook` | `{ "delivery_id": string, "branch": string, "sha": string }`。 |
| `approval` | `{ "approval_id": string, "branch": string, "sha": string, "target": string }`。 |

`running: false` の場合、`current_build_id` は `null` とする。`DELETE /api/queue` は `queued` を空配列へ置換し、`running` と `current_build_id` は変更しない。

**`.build_circuit_state` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `open` | boolean | 必須 | boolean | `true` の間は自動ポーリングを停止する。 |
| `consecutive_failures` | integer | 必須 | 0 以上 | 連続失敗回数。 |
| `opened_at` | string/null | 必須 | ISO 8601 または `null` | open に遷移した日時。 |
| `last_failure_at` | string/null | 必須 | ISO 8601 または `null` | 最終失敗日時。 |
| `last_error` | string/null | 必須 | 文字列または `null` | 最終失敗理由。 |

初期値は `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` とする。`POST /api/circuit-breaker/reset` は初期値へ戻す。

**`.config_log` JSON Lines schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | 変更日時。 |
| `type` | string | 必須 | 状態ファイル種別 | 例: `"server_config"`, `"notify_config"`, `"repo_config"`。 |
| `action` | string | 必須 | `"create"`, `"update"`, `"delete"` | 変更種別。 |
| `diff` | object | 必須 | `{key:[before,after]}` | 変更前後。秘密情報は `"***"`。 |
| `diff_text` | string | 必須 | 1 文字以上 | 人間向け差分。秘密情報は `"***"`。 |

`diff` は `{key:[before,after]}` とする。`diff_text` は 1 行以上の文字列とし、値は JSON 表現で記録する。キー名に `password`、`token`、`secret`、`pat`、`smtp_password` を含む値は before / after とも `"***"` に置換する。配列や object の内部 key も同じ規則で再帰的にマスクする。

**`.access_log` JSON Lines schema：**

各行はログイン、ログアウト、API token 作成/失効、read token 認証の監査イベントを表す JSON object とする。秘密情報、セッショントークン、API token 本体を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | イベント日時。 |
| `action` | string | 必須 | `"login"`, `"logout"`, `"token_create"`, `"token_revoke"`, `"token_auth"` | 監査イベント種別。 |
| `result` | string | 必須 | `"success"`, `"failure"` | 成否。 |
| `session_id` | string/null | 必須 | 文字列または `null` | セッション識別用の短縮 ID。token 本体ではない。 |
| `token_id` | string/null | 必須 | API token id または `null` | API token 関連イベントの対象。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。取得不能時は `null`。 |
| `reason` | string/null | 必須 | 文字列または `null` | 失敗理由。秘密情報を含めない。 |

**`.api_access_log` JSON Lines schema：**

各行は認証後 API、認証失敗 API、Webhook API の HTTP 呼び出し 1 件を表す JSON object とする。password、token、Webhook secret、SMTP password、request body の secret 値を保存してはならない。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `at` | string | 必須 | ISO 8601 | response 送信直前の日時。 |
| `request_id` | string | 必須 | `req{YYYYMMDDHHmmss}-NNN` | API 呼び出し識別子。 |
| `method` | string | 必須 | HTTP method | `GET` / `POST` / `PUT` / `PATCH` / `DELETE`。 |
| `path` | string | 必須 | `/api/...` | query を含まない path。 |
| `query` | object | 必須 | JSON object | 許可済み query key と値。秘密値は禁止。 |
| `status` | integer | 必須 | HTTP status code | response status。 |
| `duration_ms` | integer | 必須 | 0 以上 | handler 開始から response 確定までのミリ秒。 |
| `auth_type` | string | 必須 | `"session"`, `"api_token"`, `"webhook"`, `"none"` | 認証種別。 |
| `actor` | string/null | 必須 | `"admin"`、token id、`"webhook"`、または `null` | 操作者。token 本体は保存しない。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。 |
| `user_agent` | string/null | 必須 | 文字列または `null` | 取得不能時は `null`。 |
| `error` | string/null | 必須 | エラーコードまたは `null` | 成功時は `null`。 |

**`.webhook_events.json` JSON Lines schema：**

`.webhook_events.json` は 1 行 1 event を追記する。mode は `0600` とする。request header 全体、署名値、Webhook secret、payload 全体を保存してはならない。

| key | 型 | 必須 | 許容値 / 説明 |
|-----|----|------|---------------|
| `timestamp` | string | 必須 | UTC ISO 8601。 |
| `delivery_id` | string/null | 必須 | `X-GitHub-Delivery`。1〜200 文字または `null`。 |
| `event` | string | 必須 | GitHub event 名。1〜100 文字。 |
| `ref` | string/null | 必須 | push ref または `null`。 |
| `branch` | string/null | 必須 | `refs/heads/` を除いた branch または `null`。 |
| `sha` | string/null | 必須 | push `after`。40 文字 lowercase hex または `null`。 |
| `repository` | string/null | 必須 | `owner/repo` 形式または `null`。 |
| `build_triggered` | boolean | 必須 | queue 追加済みなら `true`。 |
| `queued_id` | string/null | 必須 | queue id または `null`。 |
| `result` | string | 必須 | `"queued"`, `"duplicate"`, `"ignored_event"`, `"ignored_branch"`, `"queue_full"`, `"error"`。 |

**`.approval_queue` JSON Lines schema：**

`.approval_queue` は append-only とし、1 行 1 record を追記する。同一 id の最新 record を有効状態として扱い、古い record は監査履歴として残す。

| キー | 型 | 必須 | 許容値 / 説明 |
|------|----|------|---------------|
| `id` | string | 必須 | `appr{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| `status` | string | 必須 | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| `branch` | string | 必須 | branch target 名。 |
| `sha` | string | 必須 | 40 文字 lowercase hex。 |
| `target` | string | 必須 | branch target id または target file。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `expires_at` | string | 必須 | UTC ISO 8601。 |
| `decided_at` | string/null | 必須 | approve / reject / expire 時刻または `null`。 |
| `decided_by` | string/null | 必須 | 管理 session は `"admin"`、API token は token id、未決定時は `null`。 |
| `queue_id` | string/null | 必須 | approve で追加した queue id または `null`。 |
| `reason` | string/null | 必須 | reject 理由、expire 理由、または `null`。 |

**`.build_history` JSON Lines schema：**

各行は以下の JSON object とする。

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | `b{YYYYMMDDHHmmss}` | ビルド ID。 |
| `build_at` | string | 必須 | ISO 8601 | ビルド完了日時。 |
| `sha` | string/null | 必須 | Git SHA または `null` | 対象 blob / commit SHA。 |
| `status` | string | 必須 | `"success"`, `"failure"`, `"cancelled"`, `"hook_error"` | ビルド結果。 |
| `trigger` | string | 必須 | `"polling"`, `"force_interval"`, `"manual"`, `"webhook"`, `"retry_pending_transfer"`, `"startup_config_integrity"`, `"rollback"`, `"local_watch"`, `"approval"` | 起動種別。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 任意 | SHA-256 hex または `null` | 成果物チェックサム。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | 所要時間。 |
| `retry_count` | integer | 任意 | 0 以上 | 最終成功または最終失敗までに実行した追加 retry 回数。未記録時は `0` と扱う。 |
| `commit_status_state` | string/null | 任意 | `"pending"`, `"success"`, `"failure"`, `"error"`, `null` | 最終 GitHub Commit Status 送信状態。未送信時は `null`。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | タグ検証に従う | 手動/自動タグ。 |
| `comment` | string/null | 必須 | コメント検証に従う | コメント。 |
| `rollback_from` | string/null | 任意 | build id または `null` | rollback の元 build id。 |

**`.build_logs/{id}.json` schema：**

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | build id | ファイル名 `{id}.json` と一致する。 |
| `started_at` | string | 必須 | ISO 8601 | ビルド開始日時。 |
| `finished_at` | string/null | 必須 | ISO 8601 または `null` | 完了前は `null`。 |
| `duration_seconds` | integer/null | 必須 | 0 以上または `null` | 完了前は `null`。 |
| `status` | string | 必須 | `"running"`, `"success"`, `"failure"`, `"cancelled"`, `"hook_error"` | 現在/最終状態。 |
| `trigger` | string | 必須 | `.build_history.trigger` と同じ | 起動種別。 |
| `branch` | string | 必須 | branch 名 | 対象 branch。 |
| `target_file` | string | 必須 | 相対パス | 対象ファイル。 |
| `sha` | string/null | 必須 | SHA または `null` | 対象 SHA。 |
| `commit_sha` | string/null | 必須 | SHA または `null` | commit SHA。 |
| `commit_message` | string/null | 必須 | 文字列または `null` | commit message。 |
| `commit_author` | string/null | 必須 | 文字列または `null` | commit author。 |
| `commit_at` | string/null | 必須 | ISO 8601 または `null` | commit 日時。 |
| `stdout` | string[] | 必須 | 0 件以上 | `pipeline.sh` stdout 行。 |
| `stderr` | string[] | 必須 | 0 件以上 | `pipeline.sh` stderr 行。 |
| `warnings` | string[] | 必須 | 0 件以上 | `[WARN]` 行または runner warning。 |
| `report` | object/null | 必須 | 下記 Report object または `null` | `[REPORT]` の解析結果。 |
| `output_size_bytes` | integer/null | 必須 | 0 以上または `null` | 成果物サイズ。 |
| `output_sha256` | string/null | 任意 | SHA-256 hex または `null` | 成果物チェックサム。 |
| `size_warn` | boolean | 必須 | boolean | サイズ警告。 |
| `transfer_verified` | boolean/null | 必須 | boolean または `null` | SSH 転送未実行時は `null`。 |
| `dry_run` | boolean | 必須 | boolean | 通常 build log では常に `false`。dry-run は build log を作成しないため、本 field が `true` の build log を新規作成してはならない。 |
| `attempts` | object[] | 必須 | 1 件以上 | build / deploy の試行履歴。dry-run は build log を作成しないため、本配列へ検証 attempt を保存しない。 |
| `retry_count` | integer | 必須 | 0 以上 | 追加 retry 回数。初回のみで終わった場合は `0`。 |
| `commit_status` | object/null | 必須 | CommitStatus object または `null` | GitHub Commit Status API 送信結果。無効時は `null`。 |
| `build_meta` | object | 必須 | BuildMeta object | 出力サイトへ埋め込んだ build metadata。 |
| `error` | string/null | 必須 | 文字列または `null` | 失敗理由。 |
| `comment` | string/null | 必須 | コメント検証に従う | コメント。 |
| `flagged` | boolean | 必須 | boolean | 重要フラグ。 |
| `tags` | string[] | 必須 | タグ検証に従う | タグ。 |

Report object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `headings` | integer | 必須 | 見出し数。 |
| `tables_count` | integer | 必須 | テーブル数。 |
| `code_blocks_count` | integer | 必須 | コードブロック数。 |
| `warnings_count` | integer | 必須 | 警告件数。 |
| `size_warn` | boolean | 必須 | 出力サイトサイズ警告。 |
| `broken_links` | integer | 必須 | 内部リンク不整合数。 |
| `heading_skips` | integer | 必須 | 見出しレベルスキップ数。 |
| `reading_time` | integer | 必須 | 推計読了時間。 |
| `theme` | string | 必須 | 使用 theme 名。 |
| `build_id` | string | 必須 | `[REPORT] build_id`。未指定時は空文字。 |
| `commit_sha` | string | 必須 | `[REPORT] commit_sha`。未指定時は空文字。 |
| `build_at` | string | 必須 | `[REPORT] build_at`。未指定時は空文字。 |

Attempt object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `attempt` | integer | 必須 | 初回は `1`。retry ごとに +1。 |
| `started_at` | string | 必須 | UTC ISO 8601。 |
| `finished_at` | string/null | 必須 | 完了時刻。dry-run は build log を作成しないため、本 field を保存しない。 |
| `stage` | string | 必須 | `"github"`, `"pipeline"`, `"deploy"` のいずれか。dry-run は build log を作成しないため `"dry_run"` stage を新規保存してはならない。 |
| `status` | string | 必須 | `"success"` または `"failure"`。 |
| `retryable` | boolean | 必須 | この失敗が retry 対象か。成功時は `false`。 |
| `error` | string/null | 必須 | 失敗理由。成功時は `null`。 |

CommitStatus object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `enabled` | boolean | 必須 | `.server_config.commit_status_enabled` の評価結果。 |
| `state` | string/null | 必須 | `"pending"`, `"success"`, `"failure"`, `"error"`、未送信時 `null`。 |
| `context` | string | 必須 | 送信 context。 |
| `target_url` | string/null | 必須 | 送信 target_url。省略時 `null`。 |
| `sent_at` | string/null | 必須 | 最終送信時刻。未送信時 `null`。 |
| `http_status` | integer/null | 必須 | GitHub API HTTP status。未送信時 `null`。 |
| `error` | string/null | 必須 | 送信失敗理由。成功時 `null`。 |

BuildMeta object:

| キー | 型 | 必須 | 説明 |
|------|----|------|------|
| `build_id` | string | 必須 | HTML meta `adlaire-build-id` と同じ値。 |
| `commit_sha` | string | 必須 | HTML meta `adlaire-commit-sha` と同じ値。 |
| `build_at` | string | 必須 | HTML meta `adlaire-build-at` と同じ値。 |

**`.build_status.json` schema：**

`.build_status.json` は runner の現在状態と直近結果を 1 ファイルで読むための要約である。API / UI / MCP は現在状態を表示する場合、`.build_status.json` を第一参照元とし、存在しない場合のみ `.build_state`、`.build_history`、`.build_lock` から後方互換の値を算出する。

```json
{
  "schema_version": 1,
  "updated_at": "2026-09-16T01:00:12Z",
  "status": "success",
  "running": false,
  "current_build_id": null,
  "last_build_id": "b20260916010000",
  "last_trigger": "polling",
  "last_target_status": "success",
  "last_branch": "main",
  "last_target_file": "docs",
  "last_blob_sha": "012345",
  "last_commit_sha": "abcdef",
  "last_started_at": "2026-09-16T01:00:00Z",
  "last_finished_at": "2026-09-16T01:00:12Z",
  "last_duration_seconds": 12,
  "last_error": null,
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "circuit_consecutive_failures": 0,
  "output_sha256": null,
  "size_warn": false
}
```

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `schema_version` | integer | 必須 | `1` 固定 | schema version。 |
| `updated_at` | string/null | 必須 | UTC ISO 8601 または `null` | 初期値のみ `null`。 |
| `status` | string | 必須 | `"none"`, `"running"`, `"success"`, `"failure"`, `"skipped_no_change"`, `"skipped_cooldown"`, `"success_deploy_pending"`, `"circuit_open"`, `"config_recovered"`, `"config_error"`, `"lock_skipped"` | `"none"` は初期値のみ。 |
| `running` | boolean | 必須 | boolean | runner が build 処理中なら `true`。 |
| `current_build_id` | string/null | 必須 | build id または `null` | 実行中 build id。実行中でなければ `null`。 |
| `last_build_id` | string/null | 必須 | build id または `null` | 最後に build log / history を作成した build id。未実行なら `null`。 |
| `last_trigger` | string/null | 必須 | `.build_history.trigger` と同じ値または `null` | 起動種別。 |
| `last_target_status` | string/null | 必須 | runner target status または `null` | runner 内部分類。 |
| `last_branch` | string/null | 必須 | branch 名または `null` | 最終対象 branch。 |
| `last_target_file` | string/null | 必須 | 相対 path または `null` | 最終対象 file。 |
| `last_blob_sha` | string/null | 必須 | Git blob SHA または `null` | 取得不能時は `null`。 |
| `last_commit_sha` | string/null | 必須 | Git commit SHA または `null` | 取得不能時は `null`。 |
| `last_started_at` | string/null | 必須 | UTC ISO 8601 または `null` | 最終開始時刻。 |
| `last_finished_at` | string/null | 必須 | UTC ISO 8601 または `null` | 最終終了時刻。 |
| `last_duration_seconds` | integer/null | 必須 | 0 以上または `null` | 最終所要時間。 |
| `last_error` | string/null | 必須 | 文字列または `null` | 成功・通常 skip は `null`。 |
| `last_deploy_status` | string/null | 必須 | `"success"`, `"pending"`, `"failed"`, `"skipped"`, `null` | 最終 deploy 状態。 |
| `pending_transfers_count` | integer | 必須 | 0 以上 | pending transfer 件数。 |
| `notify_pending_count` | integer | 必須 | 0 以上 | pending notification 件数。 |
| `circuit_open` | boolean | 必須 | boolean | `.build_circuit_state.open` と一致する。 |
| `circuit_consecutive_failures` | integer | 必須 | 0 以上 | `.build_circuit_state.consecutive_failures` と一致する。 |
| `output_sha256` | string/null | 必須 | SHA-256 hex または `null` | 直近成功成果物の manifest SHA-256。 |
| `size_warn` | boolean | 必須 | boolean | 直近 report の size warning。 |

`.build_status.json` の更新タイミング、各 `status` の選択条件、書き込み失敗時の runner 終了コードは [`docs/details/runner.md`](runner.md) §13 の build status 更新契約を正とする。

### §22.0s 状態ファイル実装完了固定契約

`statefile` owner component は、§22.0a〜§22.0c の schema、adapter、更新手順、破損時処理を実装単位として扱う。実装者は状態ファイルごとに別々の暗黙処理を追加せず、下表の共通契約を満たす。

| 観点 | 入力 | 必須処理 | 成功時出力 | 失敗時出力 / 副作用 |
|------|------|----------|------------|---------------------|
| path 解決 | state dir、§22.0a の状態ファイル path | state dir 外へ出る path、absolute user input、`..`、symlink 経由の secret 参照を拒否する。 | 正規化済み target path。 | `500` または呼び出し元の固定 error。target を作成・変更しない。 |
| schema 読取 | JSON object / JSON array / JSON Lines / text | §22.0c の key、型、nullable、列挙値、UTC 時刻形式を検証する。 | typed value。 | `ErrStateCorrupted` または行単位 skip。未知 key を削除して成功扱いにしない。 |
| 初期値 | file 不在 | §22.0a の「不在時」または初期値を返す。 | 初期 typed value。 | read-only 呼び出しでは file を作成しない。write 呼び出しだけ更新手順で作成する。 |
| atomic write | 更新後 JSON / text | `{name}.lock`、tmp、chmod、rename、file sync、parent sync、lock 削除を順に行う。 | target が完全な新内容へ置換される。 | target は旧内容を維持する。tmp と lock は best effort で削除し、呼び出し元へ `500` を返す。 |
| lock timeout | 既存 `{name}.lock` | 100ms 間隔、最大 10 秒待つ。 | lock 取得後に更新継続。 | `409 {"error":"Conflict"}` または呼び出し元の conflict error。target を変更しない。 |
| chmod | target file / directory | 秘密情報 `0600`、通常 file `0644`、directory `0755` を適用する。 | mode が固定値に一致する。 | chmod 失敗は成功扱いにしない。rename 前なら target 変更なし、rename 後なら ERROR ログへ記録する。 |
| corrupt backup | §22.0a で退避指定された破損 file | `{name}.corrupt.{YYYYMMDDHHMMSS}.bak` へ同一 directory 内で rename する。 | backup file と再生成初期値。 | backup 失敗時は再生成せず `500`。secret 内容を log / response に含めない。 |
| JSON Lines | JSON Lines file | 空行、JSON object 以外、必須 key 不足、型不一致行を除外する。 | 有効行だけの配列。 | 壊れた行は server log に固定 code、path、line number だけ記録する。response に壊れた行数を含めない。 |
| 複数ファイル更新 | 複数 state 書込 endpoint | API detail の Write 列順に 1 file ずつ atomic write する。 | 全対象が順に更新される。 | 未処理 file は変更しない。更新済み file は自動 rollback しない。`.config_log` に失敗を記録する。 |

**read / write 境界固定：**

| 呼び出し種別 | 許可する処理 | 禁止する処理 |
|--------------|--------------|--------------|
| `GET` endpoint | 既存 target の読取、typed value 変換、JSON Lines の有効行抽出、fallback 値算出。 | file 作成、chmod、corrupt backup、旧形式保存、lock 待機、lock 削除、tmp 作成。 |
| write endpoint | 入力 validation 後の atomic write、§22.0a で定義された初期値作成、定義済み corrupt backup。 | validation 前の状態変更、未定義 file 作成、未知 key 保存、secret 平文 log。 |
| runner write | build lifecycle に必要な state 更新、`.build_lock` stale 判定後の lock 削除。 | API 専用 credentials / token / auth state の直接変更。 |
| setup write | 初期配置に必要な `.github_token`、`.last_sha`、admin directory の配置。 | API runtime state、history、build log、session、token の生成。 |

**runner 連動状態更新固定ゲート：**

runner / archive / commitstatus / security / api が同じ実装 PR で状態更新を組み合わせる場合でも、statefile owner の契約は下表で固定する。呼び出し元 component は業務判断を持ち、statefile は path、schema、lock、atomic write、JSON Lines、破損時処理だけを担当する。

| ゲート | statefile 側の固定処理 | 呼び出し元が渡す値 | 失敗時境界 |
|--------|------------------------|-------------------|------------|
| schema precheck | 保存前に §22.0c の key、型、nullable、enum、UTC 時刻、配列要素 schema を検証する。 | 保存済みとして確定した typed value。 | schema 不一致は target 変更なしで `ErrStateCorrupted` または validation error を返す。 |
| unknown key rejection | 既存 file と新規 value の両方で未知 key を拒否する。例外は本節に明記済みの旧形式正規化だけ。 | 表示用 key、SDK 用 key、fixture 用 key を含まない object。 | 未知 key を削除して保存しない。既存未知 key も暗黙修復しない。 |
| write order evidence | 複数 file 更新では呼び出し元が決めた順に 1 file ずつ atomic write し、fixture の `write_order` と一致させる。 | 順序付き write plan。 | 失敗地点以降は実行しない。成功済み file は statefile が rollback しない。 |
| JSON Lines append | append 対象は 1 行 1 JSON object とし、末尾 newline を固定する。 | 1 record の typed value。 | append 失敗は対象操作へ返し、既存行の rewrite、sort、修復をしない。 |
| no mutation read | read-only adapter は fallback 値を返すだけで、file 作成、chmod、backup、lock 削除、旧形式保存を行わない。 | 読取対象 path と fallback 条件。 | 読取失敗は typed error を返し、filesystem 差分なし。 |
| corrupt handling | §22.0a に再生成指定がある file だけ backup → 初期値再生成を許可する。 | 破損判定結果と対象 path。 | backup 失敗時は再生成しない。再生成指定がない file は変更しない。 |
| secret path | secret を含む file は `0600`、通常 state は `0644`、directory は `0755` に固定する。 | secret file 判定。 | chmod 失敗を成功扱いにせず、secret 内容を log / response / fixture expected に出さない。 |

**状態ファイル fixture 合格ゲート：**

| fixture | 初期状態 | 操作 | 合格条件 |
|---------|----------|------|----------|
| state read missing | target 不在 | 対応する read adapter 呼び出し | §22.0a の不在時戻り値を返し、filesystem 差分なし。 |
| state corrupt object | JSON parse 不能または未知 key あり | read endpoint 呼び出し | `500 {"error":"State file is corrupted"}`。target 差分なし。 |
| state corrupt regenerates | §22.0a で再生成指定済み file が破損 | write endpoint または再生成を伴う操作 | corrupt backup が 1 件作成され、初期値だけが保存される。 |
| state lock timeout | `{name}.lock` が 10 秒以上残る | write endpoint 呼び出し | `409`、target/tmp 差分なし。 |
| state chmod failure | chmod を fake failure | write endpoint 呼び出し | 成功扱いにせず、target 更新有無が atomic write 表の失敗時動作と一致する。 |
| state fsync failure | file sync または parent sync を fake failure | write endpoint 呼び出し | `500`、ERROR log、secret 非表示。 |
| json lines partial corrupt | 有効行と破損行が混在 | list endpoint 呼び出し | 有効行だけ返し、server log に line number、response body に破損詳細なし。 |
| get no mutation | 破損なし state 一式 | 全 read-only endpoint 呼び出し | state dir の file list、mtime、mode、content が変化しない。 |
| multi write partial failure | 2 file 目の write を fake failure | 複数ファイル更新 endpoint 呼び出し | 1 file 目は保持、2 file 目以降は未変更、`.config_log` に失敗記録。 |

`statefile` は上表の fixture expected が用意され、成功系、validation failure、corrupt read、lock timeout、chmod failure、fsync failure、GET no mutation の差分が確認できるまで実装完了扱いにしてはならない。
