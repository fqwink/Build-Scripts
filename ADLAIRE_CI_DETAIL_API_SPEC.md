# Adlaire CI — API 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` から分割した `api` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

`ADLAIRE_CI_DETAIL_SPEC.md` は、詳細仕様の入口、索引、共通固定値、責務 component 対応表を持つ。本ファイルを読む前に、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j を確認する。

`ADLAIRE_CI_DETAIL_SPEC.md` §27.38a は runner / builder / API / SDK / UI にまたがる横断補足契約であり、本ファイルへ移動しない。API 連動機能を実装する場合は、本ファイルの個別節と合わせて `ADLAIRE_CI_DETAIL_SPEC.md` §27.38a を確認する。

---

### — 管理ツール —

## 21. 管理ツール システム構成

```
systemd timer
  └─ components/runner.go（変更検出・ビルド起動）
       └─ SSH 転送

components/api.go（常駐 HTTP サーバー）

admin/index.html（標準管理ツール）
  └─ adlaire-ci-sdk.js（SDK）─── HTTP ───► components/api.go
```

owner component `api` は、Go 標準ライブラリ `net/http` で実装し、管理ツールからの API リクエストを受け付ける。`runner` とは独立して常駐する。

**`components/api.go` 設定値（スクリプト冒頭）：**

```go
Host              = "127.0.0.1"                               // バインドアドレス（外部公開禁止）
Port              = 8765                                      // リッスンポート
CredentialsFile   = "/opt/adlaire-builder/.admin_credentials" // 認証情報ファイル
OutputURL         = "https://example.com/" // 出力サイトの公開 URL
HistoryFile       = "/opt/adlaire-builder/.build_history"     // ビルド履歴ファイル
NotifyConfigFile  = "/opt/adlaire-builder/.notify_config"     // Webhook 通知設定
ServerConfigFile  = "/opt/adlaire-builder/.server_config"     // サーバー設定
AccessLogFile     = "/opt/adlaire-builder/.access_log"        // ログイン履歴
NotifyLogFile     = "/opt/adlaire-builder/.notify_log"        // Webhook 送信履歴
WebhookSecretFile = "/opt/adlaire-builder/.webhook_secret"    // GitHub Webhook HMAC-SHA256 Secret（→ §22）
SnapshotDir       = "/opt/adlaire-builder/.snapshots"         // スナップショット保存ディレクトリ（→ §14b）
LogLevel          = "INFO"
Owner             = "<GitHubオーナー名>"                       // 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
Repo              = "<リポジトリ名>"                           // 初期値。POST /api/repo-config で動的変更可能（.repo_config に保存）
```

**systemd ユニット（常駐型、タイマー不要）：**

```ini
[Unit]
Description=Adlaire CI API Server
After=network.target

[Service]
Type=simple
User=deploy
ExecStart=/usr/local/bin/adlaire-ci-api
Restart=on-failure
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now adlaire-ci-api  # 登録・起動
sudo journalctl -u adlaire-ci-api -f        # ログ確認
```

---

## 22. バックエンド API 仕様

**ベース URL：** `http://localhost:{PORT}/api`
**認証：** `Authorization: Bearer {SESSION_TOKEN}`（`/api/login` で取得したセッショントークン）
**レスポンス形式：** JSON

### 22.0 API 共通契約

本節の API は `components/api.go` の対象仕様である。実装時は、エンドポイント固有仕様より先に以下の共通契約を満たす。

| 項目 | 仕様 |
|------|------|
| Go バージョン | Go `1.22` 以上。HTTP 実装は Go 標準ライブラリ `net/http` を使用する。 |
| bind | 既定値は `127.0.0.1:8765`。`--addr` で上書き可能。`--addr 0.0.0.0:<port>` を指定しても、`components/api.go` は TLS listener、origin 制限、IP allowlist、reverse proxy 設定生成を追加実行しない。 |
| 文字コード | リクエストボディ、レスポンスボディ、状態ファイルはいずれも UTF-8 とする。 |
| JSON レスポンス | JSON レスポンスには `Content-Type: application/json; charset=utf-8` を付与する。 |
| リクエスト body 上限 | JSON body は 1 MiB を上限とする。超過時は `413 Payload Too Large` と `{"error": "Payload too large"}` を返す。 |
| request body 禁止 | §22.0e で `Request` が `none` の endpoint に body がある場合は `400 Bad Request` と `{"error": "Request body is not allowed"}` を返す。 |
| 成功レスポンス | 各エンドポイント例に記載した JSON オブジェクトを返す。空レスポンスは使用しない。 |
| エラーレスポンス | エラー時は `{"error":"<message>"}` を返す。入力検証失敗時のみ `details` を配列 `[{ "field": "<field>", "message": "<reason>" }]` とし、複数エラーがある場合はリクエスト JSON の出現順、query、path parameter の順で並べる。入力検証以外の補足は `details` を使わず、`error` を実装者向けではない固定文言にする。 |
| 未知のパス | 定義されていない `/api/...` は `404 Not Found` と `{"error": "Not found"}` を返す。 |
| 未対応メソッド | パスは存在するがメソッドが異なる場合は `405 Method Not Allowed` と `{"error": "Method not allowed"}` を返す。 |
| JSON 不正 | JSON ボディのパースに失敗した場合は `400 Bad Request` と `{"error": "Invalid JSON"}` を返す。 |
| 入力検証失敗 | 型、必須キー、範囲、有効値が仕様と異なる場合は `422 Unprocessable Entity` と `{"error":"Validation failed","details":[...]}` を返す。`field` は JSON body key、query key、または path parameter 名とし、body 全体の形式不正は `field` を `"$"` とする。 |
| 認証なし | 認証必須エンドポイントで Bearer トークンがない、または無効な場合は `401 Unauthorized` と `{"error": "Unauthorized"}` を返す。 |
| 権限不足 | 認証済み API token の scope が不足する場合は `403 Forbidden` と `{"error":"Forbidden"}` を返す。管理 session は全 API 操作を許可する。API token は `read`、`trigger`、`operate`、`config`、`admin` の scope だけを許可し、token 作成時に指定された scope 外の endpoint は拒否する。 |
| 競合 | 現在状態と要求操作が両立しない場合は `409 Conflict` を返す。対象は、実行中ビルドへの二重開始、ビルド未実行時の cancel、停止済みスケジュールへの pause、稼働中スケジュールへの resume、lock 取得 10 秒超過、stale 判定不能な `.build_lock` である。 |
| 未設定機能 | endpoint の必須 secret、必須外部設定、必須状態ファイルが未設定で処理を開始できない場合は `501 Not Implemented` と `{"error":"Not configured"}` を返す。エンドポイント固有仕様で `422`、`503`、`500` を明記している場合のみ個別指定を優先する。 |
| 時刻形式 | API レスポンスと状態ファイルの機械処理用時刻は UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ` とする。明示オフセット、timezone なし文字列、ミリ秒付き文字列は保存しない。外部 API から取得した時刻も保存前に UTC `Z` へ正規化する。 |
| GET の副作用 | `GET` エンドポイントは状態ファイルを書き換えない。診断 API が外部確認を行う場合も、結果保存は行わない。 |
| 状態ファイル更新 | JSON 状態ファイルの更新は同一ディレクトリの一時ファイルへ書き出してから `os.Rename` で置換する。秘密情報を含むファイルは作成後に mode `600` を設定する。rename 後は対象ファイルと親ディレクトリを `Sync` し、永続化失敗時は `500` を返す。 |
| 秘密情報 | PAT、Webhook Secret、セッショントークン、API トークンはログ、バックアップ、GET レスポンスへ平文出力しない。設定済み表示は `"***"` または boolean で返す。 |
| 並列更新 | 同一状態ファイルを更新する API は、ファイル単位のロックを取得してから読み込み、検証、書き込みを行う。ロック取得待ちは最大 10 秒とし、超過時は `409 Conflict` を返す。 |
| 監査ログ | 設定変更 API は、変更前後の値を `.config_log` に追記する。ただし秘密情報の値は変更前後とも `"***"` にマスクする。 |
| CORS | 既定では CORS ヘッダーを付与しない。標準管理ツールは同一 origin から配信する。`OPTIONS` preflight は定義しない。CORS を有効化する拡張は本ファイルで未定義とし、実装してはならない。 |
| セキュリティヘッダー | すべての API レスポンスに `Cache-Control: no-store`、`X-Content-Type-Options: nosniff` を付与する。SSE は `Cache-Control: no-store` と `X-Accel-Buffering: no` を付与する。 |
| 判定順 | path 解決 → method 検証 → body 可否/サイズ検証 → JSON parse → 認証 → 権限 → 入力検証 → 状態競合 → 処理実行の順に判定する。 |

エンドポイント例に記載されたフィールド名、型、有効値、HTTP ステータスは規範とする。API、SDK、標準管理ツールのいずれかを変更する場合は、§22、§23、§24 の対応関係を同時に確認する。

**API 共通エラー固定文言：**

| 条件 | HTTP status | body |
|------|-------------|------|
| 未知 path | `404` | `{"error":"Not found"}` |
| method 不一致 | `405` | `{"error":"Method not allowed"}` |
| body 禁止 endpoint に body あり | `400` | `{"error":"Request body is not allowed"}` |
| body 上限超過 | `413` | `{"error":"Payload too large"}` |
| JSON parse 失敗 | `400` | `{"error":"Invalid JSON"}` |
| 認証なし / 無効 token / 期限切れ session | `401` | `{"error":"Unauthorized"}` |
| scope 不足 / 管理操作不可 | `403` | `{"error":"Forbidden"}` |
| rate limit 超過 | `429` | `{"error":"Too many requests"}` |
| 入力検証失敗 | `422` | `{"error":"Validation failed","details":[...]}` |
| 状態競合 | `409` | endpoint 固有文言。未定義の場合は `{"error":"Conflict"}` |
| 必須設定なし | `501` | `{"error":"Not configured"}` |
| 内部処理失敗 | `500` | `{"error":"Internal server error"}` |

上表の body は空白差分を除いて固定とする。`500` の response body に Go error、path、secret、状態ファイル内容、外部 API response body を含めてはならない。内部原因は server log にだけ固定コード付きで出力する。

**API 実行順・副作用境界固定契約：**

全 endpoint は、下表の順序と副作用境界に従う。endpoint 固有節で異なる順序を明記していない限り、実装者判断で検証順、状態読取、状態書込、外部呼び出し、ログ追記の順序を入れ替えてはならない。

| 段階 | 処理 | 失敗時 status | 失敗時副作用 |
|------|------|---------------|--------------|
| 1 | path 解決。未知 path を判定する。 | `404` | 状態ファイル、外部 API、監査ログ、設定ログを変更しない。`.api_access_log` だけ §27.6 に従って記録する。 |
| 2 | method 検証。path が存在し method が不一致か判定する。 | `405` | endpoint 固有処理を開始しない。`.api_access_log` 以外を変更しない。 |
| 3 | body 禁止、body size、JSON parse を検証する。`Request=none` endpoint では JSON parse を行わない。 | `400` / `413` | 認証、状態読取、状態書込、外部呼び出しを行わない。 |
| 4 | access control、maintenance、認証不要 endpoint 判定、Bearer/session/API token 認証を実行する。 | `401` / `403` / `503` | 認証失敗時は endpoint 状態を読まない。認証処理で定義された session/ticket 更新とログ追記だけを許可する。 |
| 5 | scope、rate limit、権限を判定する。 | `403` / `429` | endpoint 固有の状態書込、外部呼び出し、コマンド実行を行わない。rate limit 状態更新と監査ログは §27.47 に限定する。 |
| 6 | query、path parameter、body schema、enum、範囲を検証する。 | `422` | 状態ファイルを変更しない。外部 API、systemd、runner、hook、通知を呼び出さない。 |
| 7 | endpoint 固有の read adapter を呼び、状態競合を判定する。 | `409` / `500` | write lock を取得していても target を変更しない。tmp file があれば削除する。 |
| 8 | endpoint 固有処理を実行し、必要な状態ファイルを §22.0d の Write 列順に更新する。 | endpoint 固有 | 途中失敗時の巻き戻しは、endpoint 固有節または §26 rollback 節に明記された範囲だけ行う。 |
| 9 | `.config_log`、`.audit_log`、`.access_log`、`.api_access_log` を仕様順に追記する。 | `500` | response を成功扱いにしない。既に確定済みの endpoint 状態は自動推測で再変更しない。 |
| 10 | response body と header を確定する。 | - | response 生成時に追加の状態読取、状態書込、外部呼び出しを行わない。 |

`GET` endpoint は段階 8 で状態ファイルを書き換えない。`POST`、`DELETE` endpoint でも、段階 6 までに失敗した場合は endpoint 固有の状態書込を一切行わない。外部 API 送信、systemd 操作、hook 実行、通知送信、snapshot 操作は、endpoint 固有節の処理順に現れる場合だけ実行する。実装者判断で「先に外部確認してから validation error を返す」処理にしてはならない。

### 22.0a 状態ファイル共通仕様

`components/api.go` および拡張後 `components/runner.go` が読み書きする状態ファイルは、下表の初期値、形式、更新責務に従う。表にない状態ファイルを追加してはならない。追加が必要な場合は、先に本節へパス、形式、初期値、更新責務、破損時の扱いを追記する。

| パス | 形式 | 初期値 | 更新責務 | 破損時の扱い |
|------|------|--------|----------|--------------|
| `.admin_credentials` | JSON object | `--init-credentials` で生成 | `components/api.go` | 起動時に ERROR ログを出し、HTTP サーバーを起動しない。 |
| `.totp_secret` | JSON object | `{"enabled":false,"secret_base32":null,"confirmed_at":null,"last_accepted_step":null}` | `components/api.go` | 読み込み不能時は TOTP 有効 login を `500` で拒否する。破損時は退避するが自動再生成で認証を弱めてはならない。 |
| `.audit_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。追記不能時は対象操作を失敗扱いにする。 |
| `.api_rate_state` | JSON object | `{"windows":{}}` | `components/api.go` | `.api_rate_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 window で再生成する。 |
| `.server_config` | JSON object | `{}` | `components/api.go` | `.server_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、空 object で再生成する。 |
| `.notify_config` | JSON object | `{"webhooks":[],"channels":[],"on":[],"summary":{"enabled":false,"interval":"weekly","hour":9,"day_of_week":1},"email":{"enabled":false,"to":[],"on":[]}}` | `components/runner.go` / `components/api.go` | `.notify_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.notify_log` | JSON Lines | 空ファイル | `components/runner.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.notify_pending` | JSON array | `[]` | `components/runner.go` | `.notify_pending.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.pending_transfers` | JSON array | `[]` | `components/runner.go` | `.pending_transfers.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`[]` で再生成する。 |
| `.build_history` | JSON Lines | 空ファイル | `components/runner.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_logs/{id}.json` | JSON object | ビルドごとに新規作成 | `components/runner.go` | 対象 ID の API は `500` を返し、既存ファイルは上書きしない。 |
| `.build_lock` | text | 不在 | `components/runner.go` | 内容は `pid={pid}\nstarted_at={UTC_ISO8601}\n` とする。PID が存在しない場合は stale lock として削除し、存在する場合は `409` 相当の実行中として扱う。形式不正または PID 判定不能の場合は上書きせず `409` を返す。 |
| `.branch_config` | JSON object | 不在 | `components/runner.go` / `components/api.go` | `.branch_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、再生成せず `BRANCH_TARGETS` デフォルトへフォールバックする。 |
| `.build_state` | JSON object | `{"running":false,"current_build_id":null,"queued":[],"last_started_at":null,"last_finished_at":null,"weekly_summary_last_sent_at":null,"weekly_summary_sent_date":null}` | `components/runner.go` / `components/api.go` | `.build_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_status.json` | JSON object | `{"schema_version":1,"updated_at":null,"status":"none","running":false,"current_build_id":null,"last_build_id":null,"last_trigger":null,"last_target_status":null,"last_branch":null,"last_target_file":null,"last_blob_sha":null,"last_commit_sha":null,"last_started_at":null,"last_finished_at":null,"last_duration_seconds":null,"last_error":null,"last_deploy_status":null,"pending_transfers_count":0,"notify_pending_count":0,"circuit_open":false,"circuit_consecutive_failures":0,"output_sha256":null,"size_warn":false}` | `components/runner.go` | `.build_status.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.build_circuit_state` | JSON object | `{"open":false,"consecutive_failures":0,"opened_at":null,"last_failure_at":null,"last_error":null}` | `components/runner.go` / `components/api.go` | `.build_circuit_state.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、初期値で再生成する。 |
| `.local_watch_state.json` | JSON object | `{"files":{}}` | `components/runner.go` | `.local_watch_state.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、full build 後に再生成する。 |
| `.build_cache.json` | JSON object | `{"schema_version":1,"entries":{}}` | `components/builder.go` | 破損時は WARN を出し、cache miss として扱い、成功後に再生成する。 |
| `.build_cache/pages/` | directory | 空ディレクトリ | `components/builder.go` | entry 不一致または読み取り不能 file は miss とし、他 entry は継続使用する。 |
| `.dependency_manifest.json` | JSON object | `{"pages":{}}` | `components/builder.go` / `components/runner.go` | 破損時は full build とし、成功後に再生成する。 |
| `.approval_queue` | JSON Lines | 空ファイル | `components/runner.go` / `components/api.go` | 読み込み可能な行のみ使用し、壊れた行は ERROR ログへ記録して無視する。 |
| `.build_trends.json` | JSON object | `{"schema_version":1,"samples":[],"summary":{"count":0,"avg_seconds":null,"median_seconds":null,"p95_seconds":null}}` | `components/runner.go` / `components/api.go` | `.build_trends.json.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、`.build_history` から再集計する。 |
| `.build_chain_config` | JSON object | `{"chains":[]}` | `components/runner.go` / `components/api.go` | `.build_chain_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、chain 無効として通常 build のみ継続する。 |
| `.repo_config` | JSON object | `{}` | `components/api.go` | `.repo_config.corrupt.{YYYYMMDDHHMMSS}.bak` へ退避し、スクリプト定数へフォールバックする。 |
| `.config_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.api_access_log` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。秘密情報は記録しない。 |
| `.webhook_secret` | text | 不在 | `components/api.go` | 読み込み不能時は Webhook 受信を `501` で拒否する。 |
| `.webhook_events.json` | JSON Lines | 空ファイル | `components/api.go` | 読み込み可能な行のみ返し、壊れた行は無視する。 |
| `.access_control` | JSON object | `{"allow":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.hooks` | JSON object | `{"hooks":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.maintenance` | JSON object | `{"enabled":false,"reason":null,"since":null}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.api_tokens` | JSON object | `{"tokens":[]}` | `components/api.go` | `500` を返し、自動再生成しない。token 管理情報の消失による意図しない再許可を防ぐため、破損ファイルは上書きしない。 |
| `.alert_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.tag_rules` | JSON object | `{"rules":[]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.pipeline_config` | JSON object | `{"extra_args":[],"env":{}}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.notes` | UTF-8 text | 空文字列 | `components/api.go` | 読み込み不能時は `500` を返し、自動上書きしない。 |
| `.smtp_config` | JSON object | SMTP 未設定値 | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |
| `.smtp_secret` | text | 不在 | `components/api.go` | 読み込み不能時は SMTP 送信を `422` で拒否する。 |
| `.dashboard_layout` | JSON object | `{"widgets":["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]}` | `components/api.go` | 初期値で再生成し、ERROR ログを記録する。 |

`.build_logs/archive/` は gzip 圧縮済み build log の保存先ディレクトリである。初期値は空ディレクトリとし、`components/runner.go` または `POST /api/logs/archive` が必要時に作成する。圧縮済みファイル名は `{id}.json.gz` 固定とし、通常 `.build_logs/{id}.json` と同じ build id を表す。

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
| 必須 key 不足 | 既定値補完が個別節で明記されていない場合は破損扱いとする。 | 必須 key はすべて明示保存する。 | 初期値再生成が §22.0a 表で指定されたファイルだけ再生成する。 |
| `null` | 型欄が `string/null`、`object/null`、`integer/null` 等で明示した key だけ許可する。 | nullable でない key に `null` を保存しない。 | validation error または破損扱い。 |
| 配列 | `[]` を既定値とする key は read adapter の戻り値で空配列を返してよい。 | 保存 API は配列 key を省略せず、空の場合も `[]` を明示する。 | 型不一致は `422` または `500`。 |
| 数値 | 整数 key は JSON number の整数だけ許可する。小数、指数表記由来の非整数、文字列数値は拒否する。 | 整数は JSON number として保存する。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| 時刻 | UTC ISO 8601 秒精度 `Z` だけ許可する。 | 保存前に UTC 秒精度へ丸める。ミリ秒、local timezone、offset 付き文字列を保存しない。 | API 入力は `422`、状態ファイル読込は破損扱い。 |
| mode | 秘密情報ファイルは `0600`、通常 JSON / JSON Lines は `0644`、directory は `0755` を標準とする。 | chmod 失敗時は成功扱いにしない。 | chmod 失敗は `500`。target を更新した後の chmod 失敗は ERROR ログに残す。 |
| 改行 | text / JSON / JSON Lines は LF で保存する。JSON object / array ファイルは末尾 LF 1 個を付ける。 | CRLF、BOM、末尾余分空白を新規保存しない。 | 入力 text が CRLF を含む場合の扱いは個別機能節に従う。 |

旧 schema からの正規化は、本ファイルに「旧 key」「変換後 key」「削除する key」「保存するか読み取り時だけか」を明記した場合だけ実装する。明記がない旧形式は破損扱いとし、黙って推測変換してはならない。

**API 状態読取アダプタ固定契約：**

`components/api.go` は、P0 / P1 endpoint の状態読取を下表の adapter 名と戻り値で実装する。各 adapter は Go 内部関数名として固定し、同じ状態ファイルを endpoint ごとに別ロジックで直接 parse してはならない。

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
| `GET /api/queue` | `readBuildState()`、必要時 `.server_config` | `queued`、`running`、`queue_max_size` を返す。 | `.build_state` 不在は `queued:[]`、`running:false`。`.server_config` 不在は既定 queue 上限。 | `.build_state` 破損は `500 {"error":"State file is corrupted"}`。 |
| `POST /api/build` / `POST /api/build/force` | `readCircuitState()` → `readBuildLock()` → `readBuildState()` | circuit closed かつ lock 非実行なら `.build_state` を更新し、build 開始または queue 追加を返す。 | `.build_state` / `.build_circuit_state` 不在は初期値。 | circuit / state 破損は `500`。`.build_lock` が valid running、形式不正、PID 判定不能の場合は `409 {"error":"Conflict"}`。 |
| `POST /api/cancel` | `readBuildLock()`、`readBuildState()` | 実行中 build を cancel request 状態へ更新する。 | lock 不在かつ running false は `409 {"error":"Conflict"}`。 | `.build_state` 破損は `500`。`.build_lock` 形式不正または PID 判定不能は `409`。 |
| `POST /api/circuit-breaker/reset` | `readCircuitState()` | `open:false`、`consecutive_failures:0`、`opened_at:null`、`last_error:null` を atomic write する。 | 不在は初期値から reset 後値を書き込む。 | 読取破損は `500 {"error":"State file is corrupted"}` とし、上書きしない。 |

`ErrStateCorrupted` は JSON parse 失敗、schema_version 不一致、必須 key 不足、型不一致、列挙値不一致、UTC 時刻形式不一致のいずれかで返す。`ErrStateReadFailed` は permission denied、通常ファイルではない path、gzip 読込失敗、I/O error で返す。API response body はそれぞれ `{"error":"State file is corrupted"}`、`{"error":"State file read failed"}` 固定とし、path、Go error、ファイル内容を含めない。

JSON Lines adapter は空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致の行を壊れた行として除外する。除外後に sort、filter、paging、`total`、`pages` を算出する。壊れた行の存在は response body に含めず、server log に固定コード、path、1 始まりの line number だけを記録する。

`.build_lock` の PID が存在しない場合、`readBuildLock()` は `running=false, stale=true, valid=true` を返す。read-only endpoint は stale lock を削除しない。build command は開始前に `.build_lock` を再読込し、同じ stale 判定なら `.build_lock` だけを削除してから新規 lock を作成する。削除失敗時は `409 {"error":"Conflict"}` とし、`.build_state` を変更しない。

### 22.0b 入力検証共通仕様

API 実装は以下の検証を共通で行う。違反時は、エンドポイント固有の指定がない限り `422 Unprocessable Entity` と `{"error":"Validation failed","details":[{"field":"<field>","message":"<reason>"}]}` を返す。

| 対象 | 検証条件 |
|------|----------|
| `id` パスパラメータ | `^[A-Za-z0-9_-]{1,64}$` に一致すること。`/`、`.`、空文字は禁止。 |
| `page` | 1 以上の整数。 |
| `per_page` | 1 以上 100 以下の整数。 |
| `limit` | 1 以上 200 以下の整数。 |
| `offset` | 0 以上の整数。 |
| `days` | 1 以上 366 以下の整数。 |
| `n` | 1 以上 1000 以下の整数。 |
| 日付 | `YYYY-MM-DD` 形式で、存在する暦日であること。 |
| 時刻 | 0 以上 23 以下の整数。 |
| URL | `http://` または `https://` で始まること。Webhook URL は本番用途では `https://` のみ許可し、`http://` はローカル検証用途のみ許可する。 |
| ファイルパス | 絶対パスのみ許可する。`..` を含むパス、NUL 文字、空文字は禁止。 |
| タグ | 1 件 1〜32 文字、最大 20 件。重複は除去して保存する。 |
| コメント | 最大 2000 文字。空文字 `""` はコメント削除として扱う。 |
| メールアドレス | `local@domain` 形式で、空白を含まないこと。 |
| CIDR | IPv4 アドレスまたは IPv4 CIDR として解釈できること。 |
| コマンド引数配列 | `string[]` とし、1 要素以上 32 要素以下。各要素は 1〜256 文字。実行は `/bin/sh -c` を使わず、Go 標準ライブラリ `os/exec` の `exec.CommandContext(args[0], args[1:]...)` とする。 |

**endpoint 別 query 検証上書き：**

共通検証値と endpoint 個別節の値が異なる場合は、下表を優先する。下表にない query は §22.0b の共通検証を使用する。

| Endpoint | query | 既定値 | 許容値 | 補足 |
|----------|-------|--------|--------|------|
| `GET /api/history` | `page` | `1` | 1 以上 | 整数文字列だけ許可する。 |
| `GET /api/history` | `per_page` | `20` | 1〜100 | `0`、負数、小数、指数表記は禁止。 |
| `GET /api/logs` | `n` | `100` | 1〜1000 | `q` は空文字を許可する。 |
| `GET /api/logs/search` | `q` | `""` | 0〜500 文字 | 空文字は全件検索ではなく level/from/to のみ検索として扱う。 |
| `GET /api/logs/search` | `from`, `to` | `""` | 空文字または `YYYY-MM-DD` | `from > to` は `422`。 |
| `GET /api/logs/search` | `level` | `""` | `info`, `warn`, `warning`, `error`, `debug`, 空文字 | 大文字小文字は区別しない。 |
| `GET /api/api-access-log` | `limit` | `100` | 1〜1000 | §27.6 を優先する。 |
| `GET /api/api-access-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/webhook-events` | `limit` | `50` | 1〜1000 | §27.13 を優先する。 |
| `GET /api/webhook-events` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/audit-log` | `limit` | `100` | 1〜200 | §27.44 を優先する。 |
| `GET /api/audit-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/stats` | `days` | `7` | 1〜366 | 整数文字列だけ許可する。 |
| `GET /api/stats/timeline` | `days` | `30` | 1〜366 | 整数文字列だけ許可する。 |
| `GET /api/stats/build-duration` | `n` | `10` | 1〜1000 | 整数文字列だけ許可する。 |
| `GET /api/stats/build-trends` | `n` | `100` | 1〜1000 | 整数文字列だけ許可する。 |

**入力検証 details 固定：**

| ケース | `details[].field` | `details[].message` |
|--------|-------------------|---------------------|
| 必須 body key 不足 | key 名 | `required` |
| 未知 body key | key 名 | `unknown field` |
| 型不一致 | key 名 | `invalid type` |
| 範囲外 | key 名または query 名 | `out of range` |
| enum 不一致 | key 名または query 名 | `invalid value` |
| path parameter 不正 | parameter 名 | `invalid path parameter` |
| body 全体が object でない | `$` | `object required` |
| query が整数でない | query 名 | `integer required` |
| 日付の暦日不正 | query 名または key 名 | `invalid date` |

複数エラーがある場合、body key は JSON object の出現順、query は URL query の出現順、path parameter は route 定義順で並べる。body、query、path にまたがる場合は body → query → path の順とする。SDK と UI は `details[].field` と `details[].message` をそのまま扱うため、実装者判断で文言を言い換えてはならない。

### 22.0c 主要状態ファイル schema

本節の schema は、API 実装、SDK 型、標準管理ツール表示、バックアップ/リストアの基準である。ここに定義したキー以外を保存してはならない。追加キーを追加する場合は、型、既定値、読み書き API、既存データの扱いを本節へ追記してから実装する。

**`.server_config` schema：**

| キー | 型 | 既定値 | 許容値 | 読み書き API | 説明 |
|------|----|--------|--------|--------------|------|
| `log_max_lines` | integer | `500` | 1〜10000 | `GET/POST /api/config` | `GET /api/logs` が返す最大行数。 |
| `history_max_count` | integer | `100` | 1〜10000 | `GET/POST /api/config` | `.build_history` の通常表示上限。削除処理の上限ではない。 |
| `build_timeout_seconds` | integer | `300` | 1〜86400 | `GET/POST /api/config` | 手動/自動ビルドのタイムアウト秒数。 |
| `log_retention_days` | integer | `30` | 0〜3650 | `GET/POST /api/config`, `POST /api/logs/cleanup` | `0` は自動削除なし。 |
| `log_level` | string | `"INFO"` | `"INFO"` / `"DEBUG"` / `"WARNING"` / `"ERROR"` | `GET/POST /api/config`, `POST /api/log-level` | `components/api.go` のランタイムログレベル。 |
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
| `config` | object | `{}` | type 別 schema | webhook url、email to、command_args 等。 |
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
| `extra_args` | string[] | 必須 | 0〜50 件 | `components/builder.go` に渡す追加 CLI 引数。 |
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

Queue entry:

| キー | 型 | 必須 | 許容値 | 説明 |
|------|----|------|--------|------|
| `id` | string | 必須 | `q` + 3 桁以上の数字 | queue id。 |
| `trigger` | string | 必須 | `"manual"`, `"webhook"` | 起動種別。強制実行は `"manual"` と `payload.force=true` で表す。 |
| `queued_at` | string | 必須 | ISO 8601 | queue 追加日時。 |
| `requested_by` | string | 必須 | `"api"`, `"webhook"` | queue 追加元。 |
| `payload` | object | 必須 | JSON object | force/webhook 等の追加情報。不要時は `{}`。 |

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
| `method` | string | 必須 | HTTP method | `GET` / `POST` / `DELETE` 等。 |
| `path` | string | 必須 | `/api/...` | query を含まない path。 |
| `query` | object | 必須 | JSON object | 許可済み query key と値。秘密値は禁止。 |
| `status` | integer | 必須 | HTTP status code | response status。 |
| `duration_ms` | integer | 必須 | 0 以上 | handler 開始から response 確定までのミリ秒。 |
| `auth_type` | string | 必須 | `"session"`, `"api_token"`, `"webhook"`, `"none"` | 認証種別。 |
| `actor` | string/null | 必須 | `"admin"`、token id、`"webhook"`、または `null` | 操作者。token 本体は保存しない。 |
| `remote_addr` | string/null | 必須 | IP 文字列または `null` | 接続元。 |
| `user_agent` | string/null | 必須 | 文字列または `null` | 取得不能時は `null`。 |
| `error` | string/null | 必須 | エラーコードまたは `null` | 成功時は `null`。 |

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

### 22.0d API と状態ファイル対応表

API 実装では、下表の read/write 以外の状態ファイルを操作してはならない。複数ファイルを write する API は、表の順序で検証、バックアップ、書き込みを行い、途中失敗時は後続ファイルを書き込まない。

| API | Read | Write | 補足 |
|-----|------|-------|------|
| `POST /api/login` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.access_log`, `.audit_log` | TOTP 有効時は session を発行せず ticket を返す。 |
| `POST /api/login/totp` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.totp_secret`, `.access_log`, `.audit_log` | password verified ticket と TOTP code を検証し、成功時に session を発行する。 |
| `POST /api/logout` | メモリ上 session | メモリ上 session | ファイルは更新しない。 |
| `POST /api/change-password` | `.admin_credentials` | `.admin_credentials` | 現 session 以外をメモリから削除する。 |
| `GET /api/audit-log` | `.audit_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/access-log` | `.access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/api-access-log` | `.api_access_log` | なし | 壊れた行は無視し、新しい順で返す。 |
| `GET /api/sessions` | メモリ上 session | なし | token 本体は返さない。 |
| `POST /api/sessions/revoke-all` | メモリ上 session | メモリ上 session, `.access_log` | 現 session 以外を削除する。 |
| `GET /api/auth/totp-status` | `.totp_secret` | なし | 単一 admin の TOTP 状態を返す。 |
| `POST /api/auth/totp-setup` | `.totp_secret` | なし | 仮 secret と otpauth URI をメモリで生成し、未確認のまま永続化しない。 |
| `POST /api/auth/totp-confirm` | `.totp_secret` | `.totp_secret`, `.audit_log` | 仮 secret を code 検証後に保存し、TOTP を有効化する。 |
| `DELETE /api/auth/totp` | `.totp_secret` | `.totp_secret`, `.audit_log` | code 検証後に TOTP を無効化する。 |
| `GET /api/status` | `.build_status.json`, `.build_state`, `.build_lock`, `.build_history` | なし | `.build_status.json` を第一参照元とする。不在時のみ `.build_state` と `.build_history` から後方互換の値を算出する。 |
| `POST /api/build` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state` または queue | 実行中かつ queue 有効なら queue へ追加する。 |
| `POST /api/build/force` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | `.build_state`, SHA cache または queue | SHA reset と build trigger は同一ロック内で行う。 |
| `POST /api/build/cancel` | `.build_lock` | `.build_state`, `.build_logs/{id}.json` | 実行中でない場合は `409`。 |
| `GET /api/build/stream` | `.build_logs/{id}.json`, `.build_state` | なし | SSE 配信のみ。ログファイルは更新しない。 |
| `GET /api/logs` | `.build_logs/` | なし | 最新ログを読む。 |
| `GET /api/logs/search` | `.build_logs/` | なし | 横断検索のみ。 |
| `GET /api/logs/export` | `.build_logs/` | なし | JSON export。 |
| `POST /api/logs/cleanup` | `.server_config`, `.build_logs/` | `.build_logs/` | 削除対象のみ削除する。 |
| `POST /api/logs/archive` | `.server_config`, `.build_logs/` | `.build_logs/archive/`, `.build_logs/` | archive 対象を gzip 圧縮し、成功後に通常 log を削除する。 |
| `GET /api/history` | `.build_history` | なし | `page`、`per_page`、`trigger`、`tag`、`flagged` で絞り込み、ページングして返す。 |
| `GET /api/history/export` | `.build_history` | なし | 全件 export。 |
| `GET /api/history/{id}/log` | `.build_logs/{id}.json` | なし | ファイル破損時は `500`。 |
| `GET /api/history/{id}/comment` | `.build_logs/{id}.json` | なし | なし。 |
| `POST /api/history/{id}/comment` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | コメントだけ更新する。 |
| `POST /api/history/{id}/flag` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | flag だけ更新する。 |
| `POST /api/history/{id}/tags` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.build_history`, `.config_log` | `.build_history` の同一 id にも反映する。 |
| `POST /api/history/{id}/rollback` | `.snapshots/{id}/`, `.server_config` | `.build_history`, `.build_logs/{new_id}.json` | rollback エントリを新規追加する。 |
| `GET /api/notify-config` | `.notify_config`, `.smtp_config` | なし | secrets はマスクする。 |
| `POST /api/notify-config` | `.notify_config` | `.notify_config`, `.config_log` | secret は GET で返さない。 |
| `GET /api/notify-log` | `.notify_log` | なし | 壊れた行は無視する。 |
| `POST /api/notify-test` | `.notify_config` | `.notify_log` | 送信結果を追記する。 |
| `POST /api/notify/weekly-summary` | `.notify_config`, `.build_logs/` | `.notify_log` | 宛先なしは `422`。 |
| `GET /api/config` | `.server_config` | なし | 既定値を merge して返す。 |
| `POST /api/config/validate` | `.server_config`, request body | なし | 保存せず検証結果だけ返す。`.config_log` も更新しない。 |
| `POST /api/config` | `.server_config` | `.server_config`, `.config_log` | 許可キーのみ更新する。 |
| `POST /api/log-level` | `.server_config` | `.server_config`, `.config_log` | `log_level` のみ更新する短縮 API。 |
| `GET /api/config-log` | `.config_log` | なし | 壊れた行は無視する。 |
| `GET /api/api-rate-limit` | `.server_config`, `.api_rate_state` | なし | rate limit 設定と現在 window summary を返す。 |
| `POST /api/api-rate-limit` | `.server_config` | `.server_config`, `.api_rate_state`, `.config_log`, `.audit_log` | policy を保存し、window state を初期化する。 |
| `GET /api/repo-info` | `.repo_config` | なし | 不在時は定数値を返す。 |
| `POST /api/repo-config` | `.repo_config` | `.repo_config`, `.config_log` | 未指定キーは保持する。 |
| `GET /api/branch-config` | `.branch_config` | なし | 不在時は default。 |
| `POST /api/branch-config` | `.branch_config` | `.branch_config`, `.config_log` | 空配列は `.branch_config` 削除。 |
| `GET /api/sysinfo` | 出力サイト, process start time | なし | 状態ファイルは更新しない。 |
| `GET /api/health` | `.build_history`, `.pending_transfers` | なし | 認証不要。 |
| `GET /api/stats` | `.build_history`, `.build_logs/` | なし | `days` の範囲を集計する。 |
| `GET /api/stats/timeline` | `.build_history` | なし | 日別集計のみ。 |
| `GET /api/stats/build-duration` | `.build_logs/` | なし | duration 集計のみ。 |
| `GET /api/output-meta` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, 出力サイト | なし | 出力サイトと直近ログを集約する。通常 log 不在時は archive log を読む。 |
| `GET /api/pat-status` | `.github_token` | なし | 結果保存なし。 |
| `POST /api/pat-verify` | `.github_token` | なし | 結果保存なし。 |
| `POST /api/pat-update` | なし | `.github_token`, `.config_log` | token 値は `.config_log` でマスクする。 |
| `GET /api/backup` | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config` | なし | secrets は `"***"` へマスクする。 |
| `POST /api/restore` | request body | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config`, `.config_log` | restore 対象ファイルを検証後に表の順で置換する。 |
| `GET /api/schedule` | `.server_config` | なし | systemd 次回実行時刻は外部確認。 |
| `POST /api/schedule/interval` | `.server_config` | `.server_config`, `.config_log` | systemd timer 反映も行う。 |
| `POST /api/schedule/pause` | `.server_config` | `.server_config`, `.config_log` | 既に paused は `409`。 |
| `POST /api/schedule/resume` | `.server_config` | `.server_config`, `.config_log` | 稼働中は `409`。 |
| `POST /api/schedule/allowed-hours` | `.server_config` | `.server_config`, `.config_log` | `null` で解除。 |
| `POST /api/schedule/force-interval` | `.server_config` | `.server_config`, `.config_log` | `hours` を保存。 |
| `POST /api/schedule/cooldown` | `.server_config` | `.server_config`, `.config_log` | `seconds` を保存。 |
| `GET /api/dashboard` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock`, `.server_config`, `.alert_rules`, `.dashboard_layout`, 出力サイト, process start time | なし | 集約のみ。 |
| `GET /api/diagnostics` | `.github_token`, 出力サイト, systemd, `.notify_config` | なし | 診断結果は保存しない。 |
| `GET /api/rate-limit` | `.github_token` | なし | GitHub API 結果を返す。 |
| `GET /api/disk-usage` | `.build_logs/`, `.build_logs/archive/`, 出力サイト | なし | 集計のみ。 |
| `GET /api/webhook-events` | `.webhook_events.json` | なし | ページングして返す。 |
| `POST /api/webhook` | `.webhook_secret`, `.branch_config`, `.build_state`, `.maintenance`, `.build_circuit_state` | `.webhook_events.json`, `.build_state` または queue | 署名検証成功後のみイベント記録する。 |
| `GET /api/webhook-config` | `.webhook_secret` | なし | secret 本体は返さない。 |
| `POST /api/webhook-config` | なし | `.webhook_secret`, `.config_log` | secret 値は `.config_log` でマスクする。 |
| `POST /api/circuit-breaker/reset` | `.build_circuit_state` | `.build_circuit_state`, `.config_log` | 初期値へ戻す。冪等。 |
| `GET /api/snapshots` | `.snapshots/` | なし | 世代一覧を返す。 |
| `GET /api/snapshots/{id}/download` | `.snapshots/{id}/` | なし | バイナリを返す。 |
| `DELETE /api/snapshots/{id}` | `.snapshots/{id}/` | `.snapshots/`, `.config_log` | 対象 snapshot のみ削除する。 |
| `GET /api/maintenance` | `.maintenance` | なし | 不在時は disabled。 |
| `POST /api/maintenance/enable` | `.maintenance` | `.maintenance`, `.config_log` | `since` を現在時刻で保存する。 |
| `POST /api/maintenance/disable` | `.maintenance` | `.maintenance`, `.config_log` | disabled 状態を保存する。 |
| `GET /api/access-control` | `.access_control` | なし | 不在時は `allow: []`。 |
| `POST /api/access-control` | `.access_control` | `.access_control`, `.config_log` | allow 全体を置換する。 |
| `GET /api/hooks` | `.hooks` | なし | hook 一覧を返す。 |
| `POST /api/hooks` | `.hooks` | `.hooks`, `.config_log` | hook id を新規採番する。 |
| `DELETE /api/hooks/{id}` | `.hooks` | `.hooks`, `.config_log` | 対象 hook のみ削除する。 |
| `GET /api/hooks/{id}/log` | `.build_logs/{build_id}_hook_{id}.json` | なし | 直近 20 件を返す。 |
| `GET /api/alert-rules` | `.alert_rules` | なし | rule 一覧を返す。 |
| `POST /api/alert-rules` | `.alert_rules` | `.alert_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/alert-rules/{id}` | `.alert_rules` | `.alert_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `GET /api/tag-rules` | `.tag_rules` | なし | rule 一覧を返す。 |
| `POST /api/tag-rules` | `.tag_rules` | `.tag_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/tag-rules/{id}` | `.tag_rules` | `.tag_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `POST /api/verify-output` | `.build_history`, 出力サイト | なし | checksum 比較のみ。 |
| `GET /api/pipeline-config` | `.pipeline_config` | なし | 不在時は既定値。 |
| `POST /api/pipeline-config` | `.pipeline_config` | `.pipeline_config`, `.config_log` | config 全体を置換する。 |
| `GET /api/notes` | `.notes` | なし | 不在時は空文字。 |
| `POST /api/notes` | `.notes` | `.notes`, `.config_log` | content 全体を置換する。 |
| `GET /api/smtp-config` | `.smtp_config`, `.smtp_secret` | なし | password は返さず `password_set` だけ返す。 |
| `POST /api/smtp-config` | `.smtp_config`, `.smtp_secret` | `.smtp_config`, `.smtp_secret`, `.config_log` | password 指定時のみ `.smtp_secret` を更新する。 |
| `POST /api/smtp-test` | `.smtp_config`, `.smtp_secret` | `.notify_log` | 送信結果を記録する。 |
| `GET /api/queue` | `.build_state` | なし | queue 状態を返す。 |
| `DELETE /api/queue` | `.build_state` | `.build_state`, `.config_log` | 実行中 build は停止しない。 |
| `GET /api/dashboard-layout` | `.dashboard_layout` | なし | 不在時は既定 widget 順。 |
| `POST /api/dashboard-layout` | `.dashboard_layout` | `.dashboard_layout`, `.config_log` | widgets 全体を置換する。 |
| `GET /api/tokens` | `.api_tokens` | なし | token 本体は返さない。 |
| `POST /api/tokens` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | token 本体は作成時のみ返し、保存はハッシュのみ。 |
| `DELETE /api/tokens/{id}` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | 対象 token を失効する。 |

### 22.0e API 完全契約表

本表は API 実装、SDK 実装、標準管理ツール実装の契約インデックスである。実装者は endpoint を追加、削除、名称変更、body 変更、response 変更する前に本表を先に更新する。下表に存在しない endpoint は実装対象外とする。SHA reset 専用 endpoint とサマリー送信専用 endpoint は定義しない。

`Request` が `none` の場合、request body を受け付けない。空 JSON object `{}` も送信してはならない。`Response` は成功時 body の schema 名または最小 object を示す。詳細 schema は §22.0c、各 endpoint の個別例、§23 SDK 仕様、§24 UI 仕様を正とする。

`{message}` は `{"message": string}` を意味する。`{message,...}` 形式の response では `message` を必須キーとし、その他のキーも表記どおり必須とする。`?` が付いたキーだけを任意キーとする。成功時に空 body、`null` body、HTTP 204 は使用しない。

| Endpoint | Request | Response | Success | Errors | Read | Write | SDK | UI |
|----------|---------|----------|---------|--------|------|-------|-----|----|
| `POST /api/login` | `{password}` | `{token?,must_change,totp_required?,ticket?}` | `200` | `401`, `422`, `429`, `500` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.access_log`, `.audit_log` | `login(password)` | ログイン |
| `POST /api/login/totp` | `{ticket,code}` | `{token,must_change}` | `200` | `401`, `422`, `429`, `500` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.totp_secret`, `.access_log`, `.audit_log` | `loginTotp(ticket,code)` | ログイン |
| `POST /api/logout` | none | `{message}` | `200` | `401` | memory session | memory session | `logout()` | 全パネル共通 |
| `POST /api/change-password` | `{current_password,new_password}` | `{message}` | `200` | `401`, `422`, `500` | `.admin_credentials` | `.admin_credentials`, memory session | `changePassword()` | パスワード変更 |
| `GET /api/audit-log` | query `{limit,offset,actor?,action?,result?}` | `{log,total}` | `200` | `401`, `403`, `422`, `500` | `.audit_log` | none | `getAuditLog()` | 監査ログ |
| `GET /api/access-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.access_log` | none | `getAccessLog()` | アクセスログ |
| `GET /api/api-access-log` | query `{limit,offset,method?,path?,status?}` | `{log,total}` | `200` | `401`, `422`, `500` | `.api_access_log` | none | `getApiAccessLog()` | アクセスログ |
| `GET /api/sessions` | none | `{sessions}` | `200` | `401` | memory session | none | `getSessions()` | セッション管理 |
| `POST /api/sessions/revoke-all` | none | `{message,revoked_count}` | `200` | `401` | memory session | memory session, `.access_log` | `revokeAllSessions()` | セッション管理 |
| `GET /api/auth/totp-status` | none | `TotpStatus` | `200` | `401`, `500` | `.totp_secret` | none | `getTotpStatus()` | セキュリティ |
| `POST /api/auth/totp-setup` | none | `{secret,otpauth_uri}` | `200` | `401`, `403`, `409`, `500` | `.totp_secret` | none | `setupTotp()` | セキュリティ |
| `POST /api/auth/totp-confirm` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `.totp_secret` | `.totp_secret`, `.audit_log` | `confirmTotp(code)` | セキュリティ |
| `DELETE /api/auth/totp` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `.totp_secret` | `.totp_secret`, `.audit_log` | `disableTotp(code)` | セキュリティ |
| `GET /api/status` | none | `StatusObject` | `200` | `401`, `500` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock` | none | `getStatus()` | ステータス |
| `POST /api/build` | none | `{message,build_id?,queued?}` | `202` | `401`, `409`, `422`, `429`, `503` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state` or queue | `triggerBuild()` | 手動実行 |
| `POST /api/build/force` | none | `{message,build_id?,queued?}` | `202` | `401`, `409`, `422`, `429`, `503` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | `.build_state`, SHA cache or queue | `buildForce()` | 手動実行 |
| `POST /api/build/cancel` | none | `{message}` | `200` | `401`, `404`, `409` | `.build_lock` | `.build_state`, `.build_logs/{id}.json` | `cancelBuild()` | 手動実行 |
| `GET /api/build/stream` | none | SSE `log/end` events | `200` | `401`, `404` | `.build_logs/{id}.json`, `.build_state` | none | `streamBuild()` | 手動実行 |
| `GET /api/logs` | query `{n,q}` | `{lines}` | `200` | `401`, `422`, `500` | `.build_logs/` | none | `getLogs()` | ログビューア |
| `GET /api/logs/search` | query `{q,from,to,level}` | `SearchResult` | `200` | `401`, `422` | `.build_logs/` | none | `searchLogs()` | ログビューア |
| `GET /api/logs/export` | none | `{exported_at,lines}` | `200` | `401` | `.build_logs/` | none | `exportLogs()` | ログビューア |
| `POST /api/logs/cleanup` | none | `{message,deleted_count}` | `200` | `401`, `500` | `.server_config`, `.build_logs/` | `.build_logs/` | `cleanupLogs()` | ログビューア, 設定 |
| `POST /api/logs/archive` | none | `{message,archived_count}` | `200` | `401`, `500` | `.server_config`, `.build_logs/` | `.build_logs/archive/`, `.build_logs/` | `archiveLogs()` | ログビューア, 設定 |
| `GET /api/history` | query `{page,per_page,trigger?,tag?,flagged?}` | `HistoryPageObject` | `200` | `401`, `422` | `.build_history` | none | `getHistory()` | ビルド履歴 |
| `GET /api/history/export` | none | `ExportObject` | `200` | `401` | `.build_history` | none | `exportHistory()` | ビルド履歴 |
| `GET /api/history/{id}/log` | path `{id}` | `HistoryLogObject` | `200` | `401`, `404`, `500` | `.build_logs/{id}.json` | none | `getHistoryLog(id)` | ビルド履歴 |
| `GET /api/history/{id}/comment` | path `{id}` | `CommentObject` | `200` | `401`, `404`, `500` | `.build_logs/{id}.json` | none | `getHistoryComment(id)` | ビルド履歴 |
| `POST /api/history/{id}/comment` | `{comment}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | `setHistoryComment(id,comment)` | ビルド履歴 |
| `POST /api/history/{id}/flag` | `{flagged}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | `setHistoryFlag(id,flagged)` | ビルド履歴 |
| `POST /api/history/{id}/tags` | `{tags}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.build_history`, `.config_log` | `setHistoryTags(id,tags)` | ビルド履歴 |
| `POST /api/history/{id}/rollback` | path `{id}` | `{message,build_id}` | `202` | `401`, `404`, `409`, `500` | `.snapshots/{id}/`, `.server_config` | `.build_history`, `.build_logs/{new_id}.json` | `rollbackHistory(id)` | ビルド履歴, スナップショット |
| `GET /api/sysinfo` | none | `SysinfoObject` | `200` | `401`, `500` | output file, process start time | none | `getSysinfo()` | システム情報 |
| `GET /api/health` | none | `HealthObject` | `200` | `500` | `.build_history`, `.pending_transfers` | none | `health()` | 死活監視 |
| `GET /api/schedule` | none | `ScheduleObject` | `200` | `401`, `500` | `.server_config`, systemd | none | `getSchedule()` | リポジトリ情報 |
| `POST /api/schedule/interval` | `{interval_seconds}` | `{message,interval_seconds}` | `200` | `401`, `409`, `422`, `500` | `.server_config` | `.server_config`, `.config_log`, systemd timer | `setScheduleInterval(seconds)` | リポジトリ情報 |
| `POST /api/schedule/pause` | none | `{message}` | `200` | `401`, `409`, `500` | `.server_config` | `.server_config`, `.config_log` | `pauseSchedule()` | リポジトリ情報 |
| `POST /api/schedule/resume` | none | `{message}` | `200` | `401`, `409`, `500` | `.server_config` | `.server_config`, `.config_log` | `resumeSchedule()` | リポジトリ情報 |
| `POST /api/schedule/allowed-hours` | `{from,to}` | `{message,allowed_hours}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setAllowedHours()`, `clearAllowedHours()` | リポジトリ情報 |
| `POST /api/schedule/force-interval` | `{hours}` | `{message,hours}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setForceInterval(hours)` | リポジトリ情報 |
| `POST /api/schedule/cooldown` | `{seconds}` | `{message,seconds}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setBuildCooldown(seconds)` | リポジトリ情報 |
| `GET /api/notify-config` | none | `NotifyConfig` | `200` | `401`, `500` | `.notify_config`, `.smtp_config` | none | `getNotifyConfig()` | 通知設定 |
| `POST /api/notify-config` | `NotifyConfig` | `{message}` | `200` | `401`, `422`, `500` | `.notify_config` | `.notify_config`, `.config_log` | `setNotifyConfig(config)` | 通知設定 |
| `GET /api/notify-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.notify_log` | none | `getNotifyLog()` | 通知設定 |
| `POST /api/notify-test` | none | `{message,webhook_url}` | `200` | `401`, `422`, `500` | `.notify_config` | `.notify_log` | `notifyTest()` | 通知設定 |
| `POST /api/notify/weekly-summary` | none | `{message,period,success_count,failure_count,success_rate}` | `200` | `401`, `422`, `500` | `.notify_config`, `.build_logs/` | `.notify_log` | `notifyWeeklySummary()` | 通知設定 |
| `GET /api/config` | none | `ConfigObject` | `200` | `401`, `500` | `.server_config` | none | `getConfig()` | 設定 |
| `POST /api/config/validate` | partial `ConfigObject` | `ConfigValidationObject` | `200` | `401`, `422`, `500` | `.server_config` | none | `validateConfig(config)` | 設定 |
| `POST /api/config` | partial `ConfigObject` | `{message,config}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setConfig(config)` | 設定 |
| `POST /api/log-level` | `{level}` | `{message,level}` | `200` | `401`, `422`, `500` | `.server_config` | `.server_config`, `.config_log` | `setLogLevel(level)` | 設定 |
| `GET /api/config-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `.config_log` | none | `getConfigLog()` | 設定 |
| `GET /api/api-rate-limit` | none | `ApiRateLimitPolicy` | `200` | `401`, `403`, `500` | `.server_config`, `.api_rate_state` | none | `getApiRateLimit()` | セキュリティ |
| `POST /api/api-rate-limit` | `ApiRateLimitPolicy` | `ApiRateLimitPolicy` | `200` | `401`, `403`, `422`, `500` | `.server_config` | `.server_config`, `.api_rate_state`, `.config_log`, `.audit_log` | `setApiRateLimit(policy)` | セキュリティ |
| `GET /api/pat-status` | none | `PatStatusObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `getPatStatus()` | システム情報 |
| `POST /api/pat-verify` | none | `PatVerifyObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `patVerify()` | システム情報 |
| `POST /api/pat-update` | `{token}` | `{message}` | `200` | `401`, `422`, `500` | none | `.github_token`, `.config_log` | `updatePat(token)` | システム情報 |
| `GET /api/stats` | query `{days}` | `StatsObject` | `200` | `401`, `422` | `.build_history`, `.build_logs/` | none | `getStats(days)` | 統計 |
| `GET /api/stats/timeline` | query `{days}` | `TimelineObject` | `200` | `401`, `422` | `.build_history` | none | `getStatsTimeline(days)` | 統計 |
| `GET /api/stats/build-duration` | query `{n}` | `BuildDurationStats` | `200` | `401`, `422` | `.build_logs/` | none | `getStatsBuildDuration(n)` | 統計 |
| `GET /api/output-meta` | none | `OutputMetaObject` | `200` | `401`, `404`, `500` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, output file | none | `getOutputMeta()` | システム情報 |
| `GET /api/repo-info` | none | `RepoInfoObject` | `200` | `401`, `500` | `.repo_config` | none | `getRepoInfo()` | リポジトリ情報 |
| `POST /api/repo-config` | partial `RepoInfoObject` | `{message}` | `200` | `401`, `422`, `500` | `.repo_config` | `.repo_config`, `.config_log` | `setRepoConfig(config)` | リポジトリ情報 |
| `GET /api/branch-config` | none | `{source,branches}` | `200` | `401`, `500` | `.branch_config` | none | `getBranchConfig()` | リポジトリ情報 |
| `POST /api/branch-config` | `{branches}` | `{message,branches_count}` | `200` | `401`, `422`, `500` | `.branch_config` | `.branch_config`, `.config_log` | `setBranchConfig(branches)` | リポジトリ情報 |
| `GET /api/backup` | none | `BackupObject` | `200` | `401`, `500` | config state files | none | `backup()` | 設定 |
| `POST /api/restore` | `BackupObject` | `{message}` | `200` | `401`, `422`, `500` | request body | config state files, `.config_log` | `restore(config)` | 設定 |
| `GET /api/dashboard` | none | `DashboardObject` | `200` | `401`, `500` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock`, `.server_config`, `.alert_rules`, `.dashboard_layout`, output file, process start time | none | `getDashboard()` | ステータス, システム診断 |
| `GET /api/diagnostics` | none | `DiagnosticsObject` | `200` | `401`, `500` | `.github_token`, output file, systemd, `.notify_config` | none | `getDiagnostics()` | システム診断 |
| `GET /api/rate-limit` | none | `RateLimitObject` | `200` | `401`, `501`, `500` | `.github_token` | none | `getRateLimit()` | システム情報 |
| `GET /api/disk-usage` | none | `DiskUsageObject` | `200` | `401`, `500` | `.build_logs/`, `.build_logs/archive/`, output file | none | `getDiskUsage()` | システム情報 |
| `GET /api/webhook-events` | query `{limit,offset}` | `{events,total}` | `200` | `401`, `422`, `500` | `.webhook_events.json` | none | `getWebhookEvents(limit,offset)` | システム診断 |
| `POST /api/webhook` | GitHub webhook body | `{message,ref?}` | `200` | `400`, `403`, `409`, `501`, `503` | `.webhook_secret`, `.branch_config`, `.build_state`, `.maintenance`, `.build_circuit_state` | `.webhook_events.json`, `.build_state` or queue | none | 外部 Webhook |
| `GET /api/webhook-config` | none | `{configured}` | `200` | `401`, `500` | `.webhook_secret` | none | `getWebhookConfig()` | 通知設定 |
| `POST /api/webhook-config` | `{secret}` | `{message}` | `200` | `401`, `422`, `500` | none | `.webhook_secret`, `.config_log` | `setWebhookConfig(secret)` | 通知設定 |
| `POST /api/circuit-breaker/reset` | none | `{message,open,consecutive_failures}` | `200` | `401`, `500` | `.build_circuit_state` | `.build_circuit_state`, `.config_log` | `resetCircuitBreaker()` | 手動実行, システム診断 |
| `GET /api/snapshots` | none | `{snapshots}` | `200` | `401`, `500` | `.snapshots/` | none | `getSnapshots()` | スナップショット |
| `GET /api/snapshots/{id}/download` | path `{id}` | binary | `200` | `401`, `404`, `500` | `.snapshots/{id}/` | none | `downloadSnapshot(id)` | スナップショット |
| `DELETE /api/snapshots/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.snapshots/{id}/` | `.snapshots/`, `.config_log` | `deleteSnapshot(id)` | スナップショット |
| `GET /api/maintenance` | none | `MaintenanceObject` | `200` | `401`, `500` | `.maintenance` | none | `getMaintenance()` | メンテナンス |
| `POST /api/maintenance/enable` | `{reason}` | `{message,since}` | `200` | `401`, `422`, `500` | `.maintenance` | `.maintenance`, `.config_log` | `enableMaintenance(reason)` | メンテナンス |
| `POST /api/maintenance/disable` | none | `{message}` | `200` | `401`, `500` | `.maintenance` | `.maintenance`, `.config_log` | `disableMaintenance()` | メンテナンス |
| `GET /api/access-control` | none | `{allow}` | `200` | `401`, `500` | `.access_control` | none | `getAccessControl()` | アクセス制御 |
| `POST /api/access-control` | `{allow}` | `{message,allow}` | `200` | `401`, `422`, `500` | `.access_control` | `.access_control`, `.config_log` | `setAccessControl(allowList)` | アクセス制御 |
| `GET /api/hooks` | none | `{hooks}` | `200` | `401`, `500` | `.hooks` | none | `getHooks()` | フック |
| `POST /api/hooks` | `{phase,command_args,abort_on_failure,timeout_seconds?}` | `HookRecord` | `201` | `401`, `409`, `422`, `500` | `.hooks` | `.hooks`, `.config_log` | `addHook()` | フック |
| `DELETE /api/hooks/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.hooks` | `.hooks`, `.config_log` | `deleteHook(id)` | フック |
| `GET /api/hooks/{id}/log` | path `{id}` | `{id,runs}` | `200` | `401`, `404`, `500` | `.build_logs/{build_id}_hook_{id}.json` | none | `getHookLog(id)` | フック |
| `GET /api/alert-rules` | none | `{rules}` | `200` | `401`, `500` | `.alert_rules` | none | `getAlertRules()` | 設定 |
| `POST /api/alert-rules` | `{metric,operator,threshold,level,message}` | `AlertRule` | `201` | `401`, `409`, `422`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `addAlertRule()` | 設定 |
| `DELETE /api/alert-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.alert_rules` | `.alert_rules`, `.config_log` | `deleteAlertRule(id)` | 設定 |
| `GET /api/tag-rules` | none | `{rules}` | `200` | `401`, `500` | `.tag_rules` | none | `getTagRules()` | 設定 |
| `POST /api/tag-rules` | `{condition,tags}` | `TagRule` | `201` | `401`, `409`, `422`, `500` | `.tag_rules` | `.tag_rules`, `.config_log` | `addTagRule()` | 設定 |
| `DELETE /api/tag-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `.tag_rules` | `.tag_rules`, `.config_log` | `deleteTagRule(id)` | 設定 |
| `POST /api/verify-output` | none | `{match,expected,actual}` | `200` | `401`, `404`, `500` | `.build_history`, output file | none | `verifyOutput()` | システム診断 |
| `GET /api/pipeline-config` | none | `PipelineConfig` | `200` | `401`, `500` | `.pipeline_config` | none | `getPipelineConfig()` | 設定 |
| `POST /api/pipeline-config` | `PipelineConfig` | `{message}` | `200` | `401`, `422`, `500` | `.pipeline_config` | `.pipeline_config`, `.config_log` | `setPipelineConfig(config)` | 設定 |
| `GET /api/notes` | none | `{content,updated_at}` | `200` | `401`, `500` | `.notes` | none | `getNotes()` | 運用ノート |
| `POST /api/notes` | `{content}` | `{message,updated_at}` | `200` | `401`, `422`, `500` | `.notes` | `.notes`, `.config_log` | `setNotes(content)` | 運用ノート |
| `GET /api/smtp-config` | none | `SmtpConfig` | `200` | `401`, `500` | `.smtp_config`, `.smtp_secret` | none | `getSmtpConfig()` | 通知設定 |
| `POST /api/smtp-config` | partial `SmtpConfig` with optional password | `{message}` | `200` | `401`, `422`, `500` | `.smtp_config`, `.smtp_secret` | `.smtp_config`, `.smtp_secret`, `.config_log` | `setSmtpConfig(config)` | 通知設定 |
| `POST /api/smtp-test` | none | `{result,message}` | `200` | `401`, `422`, `500` | `.smtp_config`, `.smtp_secret` | `.notify_log` | `smtpTest()` | 通知設定 |
| `GET /api/queue` | none | `{queued,max_size}` | `200` | `401`, `500` | `.build_state` | none | `getQueue()` | 手動実行 |
| `DELETE /api/queue` | none | `{message,cleared_count}` | `200` | `401`, `500` | `.build_state` | `.build_state`, `.config_log` | `clearQueue()` | 手動実行 |
| `GET /api/dashboard-layout` | none | `{widgets}` | `200` | `401`, `500` | `.dashboard_layout` | none | `getDashboardLayout()` | ステータス |
| `POST /api/dashboard-layout` | `{widgets}` | `{message}` | `200` | `401`, `422`, `500` | `.dashboard_layout` | `.dashboard_layout`, `.config_log` | `setDashboardLayout(widgets)` | ステータス |
| `GET /api/tokens` | none | `{tokens}` | `200` | `401`, `403`, `500` | `.api_tokens` | none | `getTokens()` | API トークン管理 |
| `POST /api/tokens` | `{label,scopes,expires_at?}` | `TokenCreateResult` | `201` | `401`, `403`, `422`, `500` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | `createToken(label,scopes,expiresAt)` | API トークン管理 |
| `DELETE /api/tokens/{id}` | path `{id}` | `{message}` | `200` | `401`, `403`, `404`, `500` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | `revokeToken(id)` | API トークン管理 |

**成果物 / archive / backup API 副作用固定契約：**

| API | 処理順序 | 成功時副作用 | 失敗時副作用 |
|-----|----------|--------------|--------------|
| `POST /api/logs/cleanup` | `.server_config` 読込 → cleanup 対象 log 算出 → 対象 file 削除 → response。 | 対象 `.build_logs/{id}.json` だけを削除する。archive 済み `.json.gz` は削除しない。 | 算出前失敗は差分なし。途中削除失敗は `500` とし、削除済み file は戻さない。 |
| `POST /api/logs/archive` | `.server_config` 読込 → 対象 log 算出 → `.build_logs/archive/{id}.json.gz.tmp.{pid}` 作成 → gzip 書込 → fsync → rename → 元 log 削除 → response。 | archive 成功した log だけ元 `.json` を削除する。gzip は 1 log 1 file。 | gzip 作成または rename 失敗時は元 log を残す。元 log 削除失敗は `500` とし、archive 済み `.json.gz` は残す。 |
| `GET /api/backup` | 対象設定 file 読込 → 不在 file に既定値適用 → secret mask → response。 | 状態ファイルを更新しない。 | 読込不能な必須 file は `500`。任意 file 不在は既定値で返す。 |
| `POST /api/restore` | request 検証 → secret mask `"***"` の既存値補完 → 全対象 payload 生成 → §22.0d の順に atomic write → `.config_log` 追記 → response。 | 設定系状態 file だけを置換する。履歴、ログ、snapshot、session、token 本体は復元しない。 | 検証失敗は差分なし。途中 write 失敗は未処理 file を書かず `500`。処理済み file は戻さない。 |
| `GET /api/snapshots/{id}/download` | id 検証 → snapshot directory 検証 → tar.gz stream 生成 → response。 | 状態ファイルを更新しない。 | 不正 id は `422`、不在は `404`、stream 中の読込失敗は接続を終了し状態差分なし。 |
| `POST /api/history/{id}/rollback` | id 検証 → running 確認 → snapshot 検証 → 新 build id 採番 → deploy 転送 → rollback log/history 保存 → response。 | 新規 rollback build log/history だけを追加する。元 snapshot、元 history、state 全体は巻き戻さない。 | 転送失敗は rollback build log/history を failure として保存し、元 snapshot は削除しない。running 中は差分なし `409`。 |

archive 対象 id と snapshot id は build id 形式だけを許可する。API は request path の URL decode 後に `/`、`\`、`..`、空文字、NUL を含む id を `422` とする。tar.gz へ格納する path は snapshot directory からの相対 path とし、絶対 path、`..`、symlink entry、hardlink entry、device entry を含めてはならない。

backup response に secret 原文を含めてはならない。`password`、`token`、`secret`、`smtp_password`、`webhook_secret`、`.github_token`、`.smtp_secret`、`.api_tokens` の hash 元値は `"***"` または `*_set:boolean` で表現する。`POST /api/restore` で `"***"` を受け取った secret は既存値保持を意味し、既存値がない場合は未設定として扱う。

**backup / restore 固定契約：**

| 項目 | 仕様 |
|------|------|
| backup 対象 | `.server_config`、`.notify_config`、`.repo_config`、`.branch_config`、`.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.dashboard_layout`、`.smtp_config`、`.webhook_secret`、`.smtp_secret`。 |
| backup 対象外 | `.admin_credentials`、`.api_tokens`、`.totp_secret`、session、`.build_history`、`.build_logs/`、`.snapshots/`、`.notify_log`、`.notify_pending`、`.webhook_events.json`、`.approval_queue`、`.build_state`。 |
| backup 不在値 | 任意設定 file 不在は schema 既定値で返す。secret file 不在は `*_set:false`。 |
| backup secret | `.webhook_secret` と `.smtp_secret` は本体を返さず、`webhook_secret_set` / `smtp_password_set` boolean だけ返す。 |
| restore 検証 | すべての対象 payload を先に schema 検証し、1 件でも不正なら書込を開始せず `422`。 |
| restore `"***"` | 対応する既存 secret がある場合だけ既存値保持。既存 secret がない場合は未設定として扱い、新規 secret 文字列として保存しない。 |
| restore secret 削除 | secret key が `null` の場合は削除。key 省略は既存保持。 |
| restore 書込順 | `.server_config` → `.notify_config` → `.repo_config` → `.branch_config` → `.access_control` → `.hooks` → `.alert_rules` → `.tag_rules` → `.pipeline_config` → `.dashboard_layout` → `.smtp_config` → `.webhook_secret` → `.smtp_secret` → `.config_log`。 |
| restore 途中失敗 | 未処理 file は書かない。処理済み file は巻き戻さない。response は `500`。 |
| restore no-op | 全対象が既存値と同一の場合は file と `.config_log` を変更せず `{ "message":"No changes" }`。 |

restore の `.config_log` は対象 file ごとの差分を 1 record にまとめ、secret はすべて `"***"` とする。backup / restore の response、server log、fixture expected に secret 平文を含めてはならない。

**backup / restore fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| backup-mask | secret 設定済みで backup | secret 本体なし、`*_set:true`。 |
| restore-validate-fail | 1 file schema 不正 | `422`、全 file 差分なし。 |
| restore-secret-keep | `"***"` かつ既存 secret あり | 既存 secret 維持、平文出力なし。 |
| restore-secret-missing | `"***"` かつ既存 secret なし | secret 未設定のまま、`"***"` を保存しない。 |
| restore-secret-delete | secret `null` | secret file 削除。 |
| restore-write-failure | 中途 write 失敗 | 未処理 file は差分なし、処理済み file は維持、`500`。 |

**ビルド操作の競合優先順位：**

`POST /api/build`、`POST /api/build/force`、`POST /api/webhook`、`POST /api/history/{id}/rollback` は、以下の順に判定する。

1. 認証・権限を確認する。Webhook は署名検証を認証の代替とする。
2. `.maintenance.enabled == true` の場合は `503 {"error":"maintenance"}` を返す。
3. `.build_circuit_state.open == true` の場合は `409 {"error":"circuit_open"}` を返す。ただし `POST /api/circuit-breaker/reset` は対象外。
4. `.build_state.running == true` の場合:
   - queue 対応 API（`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`）は `queue_max_size > 0` かつ空きがあれば queue に追加し、`202 {"message":"Build queued","queued":true}` を返す。
   - queue が無効または満杯の場合は `429 {"error":"queue_full"}` を返す。
   - rollback は queue へ積まず、`409 {"error":"Build is running"}` を返す。
5. 実行中でなければ `.build_state.running=true` と新しい `build_id` を保存し、`202 {"message":"Build started","build_id":"...","queued":false}` を返す。

`POST /api/build/cancel` は queue を削除しない。実行中 build のみを cancel 対象とし、実行中でない場合は `409 {"error":"No build is running"}` を返す。`DELETE /api/queue` は実行中 build を停止せず、待機 queue のみ削除する。

### 22.0e.1 API 機能別処理契約

集約 API は、下表の読取元、算出方法、空状態の戻り値に従う。下表にない読取元や推測値を使用してレスポンスを補完してはならない。

| 機能 | Endpoint | 読取元 | 算出方法 | 空状態 / 不足時 |
|------|----------|--------|----------|-----------------|
| 現在状態 | `GET /api/status` | `.build_status.json`, `.build_history`, `.build_state`, `.build_lock` | `.build_status.json` から `last_sha`、`last_build_at`、`last_build_status`、`last_trigger`、pending 件数、circuit 状態を返す。不在時のみ `.build_history` の最新行と `.build_state` から算出する。`last_sha` は `.build_status.json.last_blob_sha` があればその値、なければ `last_commit_sha`、両方なければ `null` とする。 | 履歴なしは `last_sha:null`, `last_build_at:null`, `last_build_status:"none"`, `last_trigger:null`, `output_url:null`。 |
| 手動ビルド開始 | `POST /api/build` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | §22.0e の競合優先順位に従い、開始または queue 追加を行う。SHA cache は変更しない。 | queue 無効または満杯は `429 {"error":"queue_full"}`。 |
| 強制ビルド開始 | `POST /api/build/force` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state`, SHA cache | 開始可能な場合のみ SHA cache を空 SHA に更新し、同一状態更新内で build を開始する。queue 追加時は queue entry の `payload.force=true` を保存する。 | queue 無効または満杯は `429 {"error":"queue_full"}`。 |
| ログ一覧 | `GET /api/logs` | `.build_logs/` | 最新 build log の `stdout`、`stderr`、`warnings` を時系列順に連結し、`n` 件に丸める。`q` が空でない場合は部分一致行だけを返す。 | ログなしは `{"lines":[]}`。 |
| ログ検索 | `GET /api/logs/search` | `.build_logs/` | 全 build log を新しい順に読み、`q`、`from`、`to`、`level` で絞り込む。`level` は行内の `[INFO]`、`[WARNING]`、`[ERROR]`、`[DEBUG]` に一致させる。 | 一致なしは `results:[]`。 |
| 履歴一覧 | `GET /api/history` | `.build_history` | JSON Lines を新しい順で読み、`trigger`、`tag`、`flagged` を指定時のみ完全一致で絞り込み、ページングする。壊れた行は無視して ERROR ログに記録する。`total` と `pages` は絞り込み後の有効行だけで算出する。 | 履歴なしまたは一致なしは `total:0`, `pages:0`, `history:[]`。 |
| 出力メタ | `GET /api/output-meta` | `.build_history`, `.build_logs/`, `.build_logs/archive/`, 出力サイト | 出力サイトの現在サイズと mtime、直近成功履歴の `output_sha256`、直近ログの `report`、直近ログまたは HTML meta の `build_id` / `commit_sha` / `build_at` を返す。通常ログを先に読み、不在時だけ archive を読む。 | 出力サイト不在は `404`。report 不在の数値は `null`、warning は `[]`。build meta 不在は空文字。 |
| ダッシュボード | `GET /api/dashboard` | `.build_status.json`, `.build_history`, `.server_config`, `.alert_rules`, `.dashboard_layout`, `.build_state`, `.build_lock`, 出力サイト, process start time | `status`、`sysinfo`、`stats(days=7)`、`schedule`、`alerts` を同一リクエスト時点で算出し、widget 順序は `.dashboard_layout.widgets` を使用する。 | `.dashboard_layout` 不在は既定 widget 順。alerts なしは `[]`。 |
| 診断 | `GET /api/diagnostics` | `.github_token`, 出力サイト, systemd, `.notify_config`, `.webhook_secret` | PAT、GitHub API、出力サイト、systemd、Webhook 設定を個別 item として返す。診断結果は保存しない。 | 各項目は `ok`、`warn`、`error` のいずれかを返す。 |
| キュー | `GET /api/queue` | `.build_state`, `.server_config` | `.build_state.queued` と `.server_config.queue_max_size` を返す。 | `.build_state` 不在は初期値で `queued:[]`。 |
| バックアップ | `GET /api/backup` | §22.0d の backup 対象状態ファイル | 設定状態だけを export し、secret 値は `"***"` または boolean にマスクする。履歴、ログ、snapshot、session は含めない。 | 不在の任意設定ファイルは初期値で返す。 |
| リストア | `POST /api/restore` | request body | 対象 state schema をすべて検証してから §22.0d の write 順に保存する。secret が `"***"` の場合は既存 secret を保持する。 | 検証失敗は書き込み前に `422`。途中失敗は未処理ファイルを書かない。 |

### 22.0e.2 API ID 採番契約

API が新規 ID を生成する機能は、下表の形式に従う。既存 ID と衝突した場合は、同一時刻内で末尾に 3 桁の連番を付け、最大 999 まで試行する。999 回衝突した場合は `500 Internal Server Error` を返す。

| 対象 | 形式 | 例 | 衝突時 |
|------|------|----|--------|
| build id | `b{YYYYMMDDHHmmss}` | `b20260915100500` | `b20260915100500-001` |
| queue id | `q{YYYYMMDDHHmmss}` | `q20260915100500` | `q20260915100500-001` |
| token id | `tok` + 6 桁連番 | `tok000001` | 既存最大番号 + 1。999999 超過時は `500` |
| hook id | `h{YYYYMMDDHHmmss}` | `h20260915100500` | `h20260915100500-001` |
| approval id | `appr{YYYYMMDDHHmmss}` | `appr20260915100500` | `appr20260915100500-001` |
| alert rule id | `r{YYYYMMDDHHmmss}` | `r20260915100500` | `r20260915100500-001` |
| tag rule id | `t{YYYYMMDDHHmmss}` | `t20260915100500` | `t20260915100500-001` |
| snapshot id | `snap{YYYYMMDDHHmmss}` | `snap20260915100500` | `snap20260915100500-001` |

### 22.0e.3 API Endpoint 種別別実装契約

API handler は endpoint ごとの個別処理へ入る前に、§22.0 の判定順を必ず適用する。共通判定後の endpoint 種別別処理は下表に従う。

| 種別 | 対象 endpoint | 処理順序 | 成功時 | 失敗時 |
|------|---------------|----------|--------|--------|
| read-only list | `GET /api/history`, `GET /api/notify-log`, `GET /api/config-log`, `GET /api/webhook-events` | query 検証 → 対象ファイル読込 → 壊れた行を除外 → sort / paging → response 生成 | `total` がある endpoint は除外後件数を返す。 | query 不正は `422`。ファイル読込不能は `500`。 |
| read-only aggregate | `GET /api/status`, `GET /api/dashboard`, `GET /api/health`, `GET /api/output-meta` | 必要ファイルを read-only で読込 → 不在時初期値適用 → 算出 → response 生成 | GET の副作用なし。`.build_status.json` 破損時も API が自動修復しない。 | 必須ファイル破損は endpoint 固有の `500`。任意ファイル不在は初期値。 |
| single-file update | `POST /api/config`, `POST /api/repo-config`, `POST /api/maintenance/*`, `POST /api/access-control`, `POST /api/dashboard-layout` | body 検証 → 現在値読込 → 差分生成 → atomic write → `.config_log` 追記 | 更新後値または `{message}` を返す。 | body 検証失敗は書込前に `422`。`.config_log` 失敗時は対象更新済みのまま `500`。 |
| multi-file update | `POST /api/restore`, `POST /api/smtp-config`, `POST /api/history/{id}/tags` | 全入力検証 → 全対象読込 → 書込計画生成 → §22.0d の Write 順に atomic write | 全対象の更新完了後に response を返す。 | 検証失敗は書込なし `422`。途中失敗は未処理ファイルを書かず `500`。 |
| secret update | `POST /api/pat-update`, `POST /api/webhook-config`, `POST /api/smtp-config` password あり | secret 入力検証 → secret ファイル atomic write mode `0600` → `.config_log` へ `"***"` で記録 | secret 本体を response に含めない。 | secret 書込失敗は `500`。ログ、response、stdout へ平文を出さない。 |
| build command | `POST /api/build`, `POST /api/build/force`, `POST /api/history/{id}/rollback` | 認証 → maintenance → circuit → running / queue 判定 → `.build_state` 更新 | `202` と開始または queue 結果を返す。 | running 競合は `409` または `429`。状態更新失敗は `500`。 |
| destructive delete | `DELETE /api/queue`, `DELETE /api/snapshots/{id}`, `DELETE /api/hooks/{id}`, `DELETE /api/alert-rules/{id}`, `DELETE /api/tag-rules/{id}`, `DELETE /api/tokens/{id}` | path / auth 検証 → 対象存在確認 → 削除または失効 → 対象 endpoint の契約に従い `.config_log` または `.audit_log` 追記 | `{message}` と件数がある場合は件数を返す。 | 対象不在は `404`。部分削除は禁止し、失敗時は `500`。 |
| external check | `POST /api/pat-verify`, `GET /api/rate-limit`, `GET /api/diagnostics`, `POST /api/smtp-test`, `POST /api/notify-test` | 設定読込 → timeout 付き外部確認 → 結果 response → 必要時 log 追記 | 確認結果を保存しない。ただし test 送信 log は仕様どおり追記する。 | 未設定は `501` または endpoint 固有 `422`。timeout は `500`。 |
| binary response | `GET /api/snapshots/{id}/download` | path 検証 → snapshot 存在確認 → archive stream | `Content-Type` と `Content-Disposition` を付与する。 | 不在は `404`。読込失敗は `500`。 |
| stream response | `GET /api/build/stream` | 認証 → 最新 / 実行中 log 特定 → SSE header → frame 送信 | `log` frame 後、必ず `end` frame を送って close する。 | log 不在は `404`。送信中断時は状態ファイルを更新しない。 |

`.config_log`、`.access_log`、`.notify_log` への追記は JSON Lines 1 行単位で行う。追記失敗時は対象 endpoint の副作用が既に完了している場合でも、失敗を `500` として返し、次回 GET で破損行を無視できる形式を維持する。追記行の末尾改行を書けなかった場合は、その行を破損行として扱う。

**JSON Lines 読取・破損行契約：**

| ファイル | GET / 集計時 | 追記失敗時 | server log 固定コード |
|----------|--------------|------------|-----------------------|
| `.build_history` | 壊れた行を無視し、有効行だけで算出する。 | runner build 処理は失敗扱い。 | `BUILD_HISTORY_SKIP_CORRUPT` |
| `.config_log` | 壊れた行を無視する。 | 設定変更 API は `500`。対象状態ファイル更新済みの場合は戻さない。 | `CONFIG_LOG_WRITE_FAILED` |
| `.access_log` | 壊れた行を無視する。 | auth / token 操作は `500`。logout だけは token 破棄を優先する。 | `ACCESS_LOG_WRITE_FAILED` |
| `.api_access_log` | 壊れた行を無視する。 | 本来の API response を優先し、`500` へ変更しない。 | `API_ACCESS_LOG_WRITE_FAILED` |
| `.audit_log` | 壊れた行を無視する。 | 監査対象操作は `500`。既に状態更新済みの場合は戻さない。 | `AUDIT_LOG_WRITE_FAILED` |
| `.notify_log` | 壊れた行を無視する。 | 通知送信結果だけ失敗扱いにし、build 成否は反転しない。 | `NOTIFY_LOG_WRITE_FAILED` |
| `.webhook_events.json` | 壊れた行を無視する。 | queue 追加前なら `500`、queue 追加後なら response に `event_log_failed:true` を含める。 | `WEBHOOK_EVENT_LOG_WRITE_FAILED` |
| `.approval_queue` | 壊れた行を無視する。 | approval entry 作成は失敗扱い。 | `APPROVAL_QUEUE_WRITE_FAILED` |

JSON Lines の壊れた行は、空行、JSON parse 失敗、JSON object 以外、必須 key 不足、型不一致のいずれかとする。壊れた行を response に含めてはならない。壊れた行を検出しても GET API が対象ファイルを自動修復、削除、上書きしてはならない。

### 22.0e.4 API レスポンス正規化契約

API response は、§22.0e の Response 列、§22.0c の schema、§23 の SDK 型定義表に一致させる。実装者は endpoint ごとに以下の正規化を行う。

| 対象 | 仕様 |
|------|------|
| object response | 必須 key をすべて含める。値が存在しない場合は、schema で nullable の key だけ `null` を使用する。 |
| array response | 配列 key は未取得、対象なし、空状態のいずれでも `[]` を返す。`null`、key 省略は禁止する。 |
| boolean | `true` / `false` だけを返す。`"true"`、`1`、`0` は使用しない。 |
| integer | JSON number の整数として返す。文字列化しない。 |
| float | JSON number とし、小数第 2 位までに丸める指定がある値だけ `math.Round(x*100)/100` 相当にする。 |
| timestamp | UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ`。空状態は nullable key なら `null`、非 nullable key なら endpoint 固有の既定値。 |
| message | 成功 message は endpoint ごとの固定文言とし、入力値を連結しない。 |
| unknown key | response に schema 外 key を追加しない。互換目的の旧 key 追加も禁止する。 |

**ページング response：**

| Endpoint | 入力 | 出力 | 算出 |
|----------|------|------|------|
| `GET /api/history` | `page`, `per_page` | `total`, `page`, `per_page`, `pages`, `history` | `pages = ceil(total / per_page)`。`total=0` の場合 `pages=0`。 |
| `GET /api/webhook-events` | `limit`, `offset` | `events`, `total` | `total` は filter 後、limit/offset 前の件数。 |
| `GET /api/api-access-log` | `limit`, `offset` | `log`, `total` | `total` は filter 後、limit/offset 前の件数。 |
| `GET /api/audit-log` | `limit`, `offset` | `log`, `total` | `total` は filter 後、limit/offset 前の件数。 |

`offset >= total` の場合は空配列を返し、`404` にしない。`page > pages` の場合は `history:[]` を返し、`page` は request 値を保持する。`page < 1`、`per_page < 1`、`limit < 1`、`offset < 0` は `422`。

**部分更新 response：**

| Endpoint | 成功 response | 補足 |
|----------|---------------|------|
| `POST /api/config` | `{message:"Config updated",config}` または `{message:"No changes",config}` | `config` は更新後に既定値 merge 済みの `ConfigObject`。 |
| `POST /api/repo-config` | `{message:"Repo config updated"}` または `{message:"No changes"}` | 未指定 key は保持する。 |
| `POST /api/branch-config` | `{message:"Branch config updated",branches_count}` | 空配列で default 復帰した場合 `branches_count=0`。 |
| `POST /api/notify-config` | `{message:"Notify config updated"}` | secret は返さない。 |
| `POST /api/smtp-config` | `{message:"SMTP config updated"}` | password は返さない。 |
| `POST /api/dashboard-layout` | `{message:"Dashboard layout updated"}` | 保存後 widgets は再取得で確認する。 |

**no-op / 部分更新副作用契約：**

| Endpoint | no-op 判定 | no-op 時の副作用 | 更新時の副作用 |
|----------|------------|------------------|----------------|
| `POST /api/config` | 正規化後 config が既存値と一致 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 | `.server_config` → `.config_log` → `.audit_log`。 |
| `POST /api/repo-config` | 指定 key の正規化後値が既存値と一致 | `.repo_config`、`.config_log` を変更しない。 | `.repo_config` → `.config_log`。 |
| `POST /api/branch-config` | 正規化後 `branch_targets` が既存値と一致 | `.branch_config`、`.config_log` を変更しない。 | `.branch_config` 作成/置換/削除 → `.config_log`。 |
| `POST /api/notify-config` | secret mask 適用後の比較で既存値と一致 | `.notify_config`、`.config_log` を変更しない。 | `.notify_config` → `.config_log`。 |
| `POST /api/smtp-config` | config と password 更新有無が既存値と一致 | `.smtp_config`、`.smtp_secret`、`.config_log` を変更しない。 | `.smtp_config` → 必要時 `.smtp_secret` → `.config_log`。 |
| `POST /api/dashboard-layout` | widgets 配列が既存値と一致 | `.dashboard_layout`、`.config_log` を変更しない。 | `.dashboard_layout` → `.config_log`。 |

no-op response は endpoint 固有の `No changes` が定義されている場合はその文言を返す。定義がない endpoint は通常成功文言を返してよいが、状態ファイル、JSON Lines、監査ログ、通知ログに差分を作ってはならない。部分更新では未指定 key を保持し、`null` が削除を意味する key は個別節に明記された key だけとする。

**削除 / 失効 response：**

| Endpoint | 成功 response | 不在時 |
|----------|---------------|--------|
| `DELETE /api/snapshots/{id}` | `{message:"Snapshot deleted"}` | `404` |
| `DELETE /api/hooks/{id}` | `{message:"Hook deleted"}` | `404` |
| `DELETE /api/alert-rules/{id}` | `{message:"Alert rule deleted"}` | `404` |
| `DELETE /api/tag-rules/{id}` | `{message:"Tag rule deleted"}` | `404` |
| `DELETE /api/tokens/{id}` | `{message:"Token revoked"}` | `404` |
| `DELETE /api/queue` | `{message:"Queue cleared",cleared_count}` | queue 空でも `200`、`cleared_count=0` |

**SSE frame 契約：**

`GET /api/build/stream` は `text/event-stream; charset=utf-8` を返し、各 frame は以下の JSON を `data:` 行に 1 件ずつ出力する。

```json
{"type":"log","line":"[INFO] build started","at":"2026-09-16T10:00:00Z"}
{"type":"end","status":"success","duration_seconds":12}
```

`type` は `"log"` または `"end"` だけとする。`log.line` は最大 4000 文字とし、超過分は末尾を切り捨てる。`end` frame は接続終了前に 1 回だけ送る。送信中に client が切断した場合、状態ファイル、history、log を変更しない。

**バイナリ response 契約：**

`GET /api/snapshots/{id}/download` は JSON error 以外では binary response とし、成功時に JSON body を返さない。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"` を付与する。`id` に `"`、`\`、改行を含む値は path 検証で `422` とする。

**レスポンス検証条件：**

| ケース | 期待結果 |
|--------|----------|
| schema 必須 key | 全 endpoint が必須 key を返す。 |
| 空配列 | 配列 key は `[]`、`null` や key 省略なし。 |
| 変更なし config | `No changes`、不要な `.config_log` 追記なし。 |
| delete queue 空 | `200`、`cleared_count=0`。 |
| SSE 正常終了 | `end` frame が 1 回だけ送信される。 |
| snapshot download | binary body、固定 header、JSON success body なし。 |

**API / SDK / UI / 状態ファイル 横断固定契約：**

下表の機能群は、API endpoint、SDK method、UI 操作、状態ファイル副作用を同じ実装単位でそろえる。API だけ、SDK だけ、UI だけを先行して仕様外の仮実装にしてはならない。UI が未実装の Phase では、UI 列は fixture の期待操作として固定し、実装完了扱いには含めない。

| 機能群 | API | SDK | UI 操作 | 状態ファイル副作用 | 成功後再取得 | 失敗時固定 |
|--------|-----|-----|---------|--------------------|--------------|------------|
| login / session | `POST /api/login`, `POST /api/login/totp`, `POST /api/logout`, `GET /api/sessions`, `POST /api/sessions/revoke-all` | `login()`, `loginTotp()`, `logout()`, `getSessions()`, `revokeAllSessions()` | login、TOTP 確認、logout、session 一括失効。 | `.admin_credentials`、`.access_log`、`.audit_log`、memory session。token 本体は永続化しない。 | login 成功後は `getDashboard()` → `getStatus()` → `getQueue()`。session 失効後は `getSessions()`。 | `401` は SDK token / ticket / secret field を破棄し、UI は `panel-login` のみ表示する。 |
| build control | `POST /api/build`, `POST /api/build/force`, `POST /api/build/cancel`, `GET /api/build/stream`, `GET /api/queue`, `DELETE /api/queue` | `triggerBuild()`, `buildForce()`, `cancelBuild()`, `streamBuild()`, `getQueue()`, `clearQueue()` | manual build、force build、cancel、stream 開始/停止、queue clear。 | `.build_state`、`.build_lock`、`.build_logs/{id}.json`、queue entry。force build 時だけ SHA cache 更新。 | `getStatus()` → `getQueue()`、stream end 後は `getLogs()` も実行。 | running は `409`、queue full は `429`、maintenance / circuit は `503` または仕様済み `409`。UI は同じ build request を自動再送しない。 |
| logs / history | `GET /api/logs`, `GET /api/logs/search`, `GET /api/history`, `GET /api/history/{id}/log`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags` | `getLogs()`, `searchLogs()`, `getHistory()`, `getHistoryLog()`, `setHistoryComment()`, `setHistoryFlag()`, `setHistoryTags()` | log 表示、検索、history 表示、comment / flag / tag 保存。 | read-only GET は副作用なし。comment / flag / tag は `.build_logs/{id}.json` と必要時 `.build_history`、`.config_log`。 | 変更系は `getHistory()`、comment は `getHistoryComment(id)` も実行。 | `404` は選択解除または not found 表示。`422 details` は field error。壊れた JSON Lines は response に含めない。 |
| config / repo / branch / schedule | `GET/POST /api/config`, `POST /api/config/validate`, `GET/POST /api/repo-config`, `GET/POST /api/branch-config`, `GET /api/schedule`, `POST /api/schedule/*` | `getConfig()`, `setConfig()`, `validateConfig()`, `setRepoConfig()`, `getBranchConfig()`, `setBranchConfig()`, `getSchedule()`, schedule 系 method | config 保存、validate、repo 保存、branch 保存、schedule 変更。 | `.server_config`、`.repo_config`、`.branch_config`、`.config_log`。validate は保存なし。systemd 反映失敗時も保存済み値は戻さない。 | 保存系は対象 GET → `getConfigLog()`。validate は再取得なし。 | `422` は書込前停止。systemd 失敗 `500` は保存済み値を再取得して表示する。no-op は状態ファイルと log を変更しない。 |
| notify / SMTP / webhook | `GET/POST /api/notify-config`, `POST /api/notify-test`, `GET /api/notify-log`, `GET/POST /api/smtp-config`, `POST /api/smtp-test`, `GET/POST /api/webhook-config`, `GET /api/webhook-events`, `POST /api/notify/weekly-summary` | `getNotifyConfig()`, `setNotifyConfig()`, `notifyTest()`, `getNotifyLog()`, `getSmtpConfig()`, `setSmtpConfig()`, `smtpTest()`, `getWebhookConfig()`, `setWebhookConfig()`, `getWebhookEvents()`, `notifyWeeklySummary()` | notify 保存、test、SMTP 保存/test、webhook secret 保存、event 表示、weekly summary。 | `.notify_config`、`.smtp_config`、`.smtp_secret`、`.webhook_secret`、`.notify_log`、`.webhook_events.json`、`.config_log`。 | 保存系は対象 GET → `getConfigLog()`。test / summary は `getNotifyLog()`。 | secret は response / log / fixture へ平文出力しない。保存成功・失敗とも UI secret field を消去する。 |
| snapshots / rollback / maintenance | `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `DELETE /api/snapshots/{id}`, `POST /api/history/{id}/rollback`, `GET /api/maintenance`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable` | `getSnapshots()`, `downloadSnapshot()`, `deleteSnapshot()`, `rollbackHistory()`, `getMaintenance()`, `enableMaintenance()`, `disableMaintenance()` | snapshot list/download/delete、rollback、maintenance enable/disable。 | `.snapshots/`、`.build_history`、`.build_logs/{new_id}.json`、`.maintenance`、`.config_log`。download は副作用なし。 | delete は `getSnapshots()`。rollback は `getHistory()` → `getStatus()`。maintenance は `getMaintenance()`。 | delete / rollback は確認必須。running rollback は `409`。maintenance enabled 中は build / rollback / 設定変更系を disabled。 |
| access / hooks / rules / pipeline / notes / layout | access、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout の GET/POST/DELETE endpoint | 対応する §23 SDK method | 保存、追加、削除、notes 保存、dashboard layout 保存。 | `.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.notes`、`.dashboard_layout`、`.config_log`。 | 対象 GET → 必要時 `getConfigLog()`。layout は `getDashboardLayout()` → `getDashboard()`。 | duplicate `409` は競合表示。validation `422` は field error。削除対象不在は `404`。 |
| tokens / audit / access logs / rate limit | `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}`, `GET /api/audit-log`, `GET /api/access-log`, `GET /api/api-access-log`, `GET/POST /api/api-rate-limit` | `getTokens()`, `createToken()`, `revokeToken()`, `getAuditLog()`, `getAccessLog()`, `getApiAccessLog()`, `getApiRateLimit()`, `setApiRateLimit()` | token 発行/失効、audit / access log 表示、rate limit 保存。 | `.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.api_rate_state`、`.server_config`、`.config_log`。token 本体は作成時 response のみ。 | token 操作は `getTokens()` → `getAuditLog()`。rate limit は `getApiRateLimit()`。 | token 本体は再取得不可。`403` は logout しない。`429` は rate limit 表示し、同一操作を自動 retry しない。 |

**横断処理順固定：**

| 処理種別 | 固定順序 |
|----------|----------|
| 認証必須 JSON API | method / path 判定 → body 禁止判定 → JSON parse → 認証 / scope → rate limit → endpoint 固有 validation → read → write 計画 → atomic write → JSON Lines 追記 → response。 |
| read-only API | method / path 判定 → body 禁止判定 → 認証 / scope → query validation → read → 壊れた任意行除外 → response。read-only API は状態ファイルを書き換えない。 |
| UI 変更操作 | panel error / success 消去 → UI 入力検証 → 対象操作 disabled → SDK 呼び出し → 成功後再取得 → success 表示 → secret 消去 → disabled 再評価。 |
| UI 取得操作 | panel error 消去 → 対象操作 disabled → SDK 呼び出し → DOM 更新 → empty state 判定 → disabled 再評価。success 表示は行わない。 |
| SDK request | 引数検証 → path / query / body 生成 → Authorization 付与 → timeout 設定 → fetch → status 判定 → response parse → token 変化適用 → return / throw。 |
| multi-file write | 全入力検証 → 全対象 read → 全 write payload 生成 → §22.0d の Write 順に atomic write → JSON Lines 追記 → response。途中失敗時は未処理ファイルを書かない。 |

**横断 fixture 固定：**

| fixture | 入力 | 必須確認 |
|---------|------|----------|
| cross auth expired | 任意の認証必須 API が `401`。 | SDK は token を破棄し、UI は全 secret field を消去して `panel-login` だけを表示する。対象 API の状態ファイル副作用なし。 |
| cross config no-op | `POST /api/config` に既存値と同一の正規化済み body。 | API は `No changes`、SDK は response をそのまま返し、UI は成功表示する。`.server_config`、`.config_log`、`.audit_log` に差分なし。 |
| cross validation details | 任意の保存 API が `422 details`。 | SDK は `AdlaireCIError.details` を保持し、UI は該当 field と panel summary に表示する。状態ファイルを書かない。 |
| cross refresh failure | 変更 API は成功し、成功後再取得の 2 件目が `500`。 | 変更副作用は維持し、同じ変更 API を再実行しない。UI は操作成功を `global-success`、再取得失敗を panel error に分けて表示する。 |
| cross secret failure | secret 保存 API が `500`。 | SDK error に secret 原文を含めず、UI は secret field を消去する。状態ファイル、log、fixture に secret 原文が残らない。 |
| cross stream invalid frame | `GET /api/build/stream` が parse 不能 frame を返す。 | SDK は `AdlaireCIError(status=0,message="Invalid SSE frame")`、UI は stream error を表示し、status / queue を再取得する。状態ファイルは変更しない。 |
| cross binary snapshot | `downloadSnapshot(id)` が binary success。 | API は binary header、SDK は `Blob`、UI は download 開始表示。JSON parse、success JSON body、状態ファイル更新なし。 |
| cross destructive cancel | 削除 / rollback 確認 dialog を cancel。 | SDK method 呼び出し 0 回、状態ファイル副作用なし、success / error 表示差分なし。 |

### 22.0f 実装優先度

対象項目の実装時は、下表の順に進める。上位の完了条件を満たす前に下位へ進んではならない。同一優先度内では、API、SDK、UI、状態ファイル、検証手順を同じ Pull Request で同期する。

| 優先度 | 対象 | 完了条件 |
|--------|------|----------|
| P0 | 認証、セッション、共通エラー、状態ファイル読み書き、`.access_log`、`.config_log` | `POST /api/login` から認証必須 API の共通処理までが §22.0〜§22.0e と一致し、秘密情報がログとレスポンスに出ない。 |
| P1 | ビルド操作、status、logs、history、queue、circuit breaker | 手動ビルド、強制ビルド、キャンセル、キュー、履歴、ログ取得が同一状態ファイル契約で動作する。 |
| P2 | config、repo、branch、schedule、PAT、diagnostics、dashboard | 設定変更が `.config_log` に残り、GET 系集約 API が状態ファイルを更新しない。 |
| P3 | notify、SMTP、webhook、webhook config、weekly summary | 通知送信責務が `components/runner.go`、設定責務が `components/api.go` に分離され、secret はマスクされる。 |
| P4 | snapshots、rollback、maintenance、access control、hooks | 運用系 API が `409`、`422`、`503` を仕様どおり返し、ロールバックは履歴に `trigger: "rollback"` を残す。 |
| P5 | alert rules、tag rules、pipeline config、notes、dashboard layout、tokens | 拡張設定が schema どおり保存され、SDK と UI の操作名が §22.0e と一致する。 |

各優先度の検証条件は以下とする。

| 優先度 | 必須検証 |
|--------|----------|
| P0 | 認証成功、認証失敗、期限切れ session、`401` 時 SDK token 破棄、`.access_log` 追記、秘密情報マスクを確認する。 |
| P1 | 手動 build、force build、running 中の queue、cancel、history/log 取得、`409`、`429`、`503` を確認する。 |
| P2 | config/repo/branch/schedule の保存、`.config_log` 追記、GET 系 API が状態ファイルを書き換えないことを確認する。 |
| P3 | Webhook test、weekly summary、SMTP test、webhook secret 保存、secret mask、通知失敗ログを確認する。 |
| P4 | snapshot list/download/delete、rollback、maintenance enable/disable、access control block、hook success/failure を確認する。 |
| P5 | rule 追加/削除、pipeline config 保存、notes 保存、dashboard layout 保存、token 発行/失効、token 本体が再取得不可であることを確認する。 |

**API P0 / P1 fixture 固定：**

P0 / P1 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は実装言語の test case 名または subtest 名へそのまま写せる粒度とし、期待 HTTP status、期待 body、状態ファイル副作用を同時に確認する。

| Fixture | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------------------|---------------|--------------------|
| A1 common route errors | 未定義 `/api/unknown`、既存 path への未許可 method、body 禁止 endpoint への body 付き request、JSON 不正文を順に送る。 | `404 {"error":"Not found"}`、`405 {"error":"Method not allowed"}`、`400 {"error":"Request body is not allowed"}`、`400 {"error":"Invalid JSON"}`。 | 状態ファイルを作成、更新、削除しない。 |
| A2 auth errors | Bearer なし、無効 token、期限切れ session、scope 不足 API token で認証必須 endpoint を呼ぶ。 | `401 {"error":"Unauthorized"}` または `403 {"error":"Forbidden"}`。 | 秘密情報を response、`.access_log`、server log に出さない。 |
| A3 status primary | 正常な `.build_status.json`、正常な `.build_lock`、`.pending_transfers`、`.build_circuit_state` を置いて `GET /api/status`。 | `.build_status.json` の値を正とし、valid running lock がある場合だけ `running:true`。 | GET 副作用なし。mtime、mode、内容が変わらない。 |
| A4 status fallback | `.build_status.json` 不在、`.build_history` に成功履歴、`.build_state` に queue、`.build_lock` 不在で `GET /api/status`。 | `status`、`last_*`、`queued`、`running:false` を fallback 算出する。 | `.build_status.json` を生成しない。 |
| A5 status corrupted | `.build_status.json` を不正 JSON にして `GET /api/status`。 | `500 {"error":"State file is corrupted"}`。 | 退避ファイル作成、自動修復、初期値上書きを行わない。 |
| A6 history jsonl filtering | `.build_history` に有効行 2 件、空行 1 件、不正 JSON 1 件、必須 key 不足 1 件を置き `GET /api/history?page=1&per_page=10`。 | `total:2`、`pages:1`、`history` は有効行だけを新しい順で返す。 | 壊れた行を書き戻し削除しない。server log に `BUILD_HISTORY_SKIP_CORRUPT`。 |
| A7 build log lookup | `.build_logs/{id}.json` 正常、通常ログ不在で archive 正常、両方不在、破損ログをそれぞれ `GET /api/history/{id}/log`。 | 正常は log object、archive は gzip 展開結果、両方不在は `404 {"error":"Not found"}`、破損は `500 {"error":"State file is corrupted"}`。 | archive を通常ログへ復元しない。破損ログを上書きしない。 |
| A8 queue state | `.build_state` 不在、正常、破損の 3 状態で `GET /api/queue`。 | 不在は `queued:[]`、正常は保存値、破損は `500 {"error":"State file is corrupted"}`。 | 不在時も `.build_state` を生成しない。 |
| A9 build conflict | valid running `.build_lock`、形式不正 `.build_lock`、PID 判定不能 lock、queue 上限到達状態で `POST /api/build`。 | lock 系は `409 {"error":"Conflict"}`。queue 上限到達は endpoint 固有の `429`。 | 失敗時に `.build_state`、`.build_history`、`.build_status.json` を変更しない。 |
| A10 circuit breaker | `.build_circuit_state` が `open:true` の状態で `POST /api/build`、次に `POST /api/circuit-breaker/reset`。 | build は `503` または endpoint 固有 circuit open error。reset は `200` と固定 message。 | build 失敗では queue 追加なし。reset は `.build_circuit_state` だけを reset 後値へ atomic write。 |
| A11 write lock timeout | `.build_state.lock` を保持した状態で `.build_state` 更新 endpoint を呼ぶ。 | 10 秒経過後 `409 {"error":"Conflict"}`。 | tmp file を残さず、target を変更しない。 |
| A12 GET side-effect zero | `GET /api/status`、`GET /api/history`、`GET /api/logs`、`GET /api/queue` を連続実行する。 | 各 endpoint は入力状態に応じた正常 response または固定 error response。 | request 前後で対象状態ファイル一覧、mtime、mode、内容が一致する。 |

**API P2〜P5 fixture 固定：**

P2〜P5 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は既存 endpoint と既存状態ファイルだけを対象とし、§22.0e にない endpoint、§22.0a にない状態ファイル、§24 にない UI 操作を追加してはならない。

| Fixture | 優先度 | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------|--------------------|---------------|--------------------|
| B1 config no-op | P2 | 既存 `.server_config` と同じ body で `POST /api/config`。 | `200` と `{message:"No changes",config}`。 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 |
| B2 config validation failure | P2 | `queue_max_size=-1`、未知 enum、相対 path を含む `POST /api/config`。 | `422 {"error":"Validation failed","details":[...]}`。 | 状態ファイルを変更しない。 |
| B3 branch config save | P2 | 有効な branch target 2 件で `POST /api/branch-config`。 | `{message:"Branch config updated",branches_count:2}`。 | `.branch_config` を atomic write し、`.config_log` に差分を記録する。 |
| B4 schedule systemd failure | P2 | `.server_config` 保存成功後、systemd timer 更新を fake failure にする。 | `500`。 | `.server_config` は更新済み、`.config_log` に `systemd_update_failed` を記録し、未定義 rollback を行わない。 |
| B5 PAT update secret mask | P2 | `POST /api/pat-update` に token を送る。 | `{message:"PAT updated"}`。 | secret file は mode `600`。response、`.config_log`、`.audit_log`、server log に token 平文を出さない。 |
| B6 dashboard read-only | P2 | `.dashboard_layout`、`.build_state`、`.build_status.json`、`.alert_rules` を置き `GET /api/dashboard`。 | widget 順に dashboard object を返す。 | GET は対象状態ファイルを作成、修復、更新しない。 |
| C1 notify config mask | P3 | Webhook secret と SMTP password を含む通知設定保存後、GET / backup / log を確認する。 | secret は `"***"` または `*_set:true` だけを返す。 | secret 平文を状態表示、履歴、通知ログ、backup に残さない。 |
| C2 webhook receive signed | P3 | 正常署名の GitHub push payload を `POST /api/webhook`。 | `202` と queued 結果。 | `.webhook_events.json` 追記後、必要時 `.build_state.queued` へ `trigger:"webhook"` を追加する。 |
| C3 webhook invalid signature | P3 | 署名なし、不正 prefix、不一致署名。 | `401 {"error":"Unauthorized"}`。 | event log、queue、history を変更しない。 |
| C4 SMTP test disabled | P3 | SMTP disabled で `POST /api/smtp-test`。 | endpoint 固有の `422`。 | `.notify_log` へ成功扱いを残さず、secret を出力しない。 |
| D1 snapshot delete | P4 | 存在する snapshot id で `DELETE /api/snapshots/{id}`。 | `{message:"Snapshot deleted"}`。 | 対象 snapshot だけ削除し、`.config_log` または監査対象 log に削除を記録する。 |
| D2 rollback running conflict | P4 | `.build_state.running=true` で `POST /api/history/{id}/rollback`。 | `409 {"error":"Build is running"}`。 | queue、history、snapshot、deploy target を変更しない。 |
| D3 maintenance blocks build | P4 | maintenance enabled 状態で `POST /api/build`。 | `503` と endpoint 固有 maintenance error。 | build queue、history、log を変更しない。 |
| D4 access-control deny | P4 | allowlist に接続元が含まれない状態で任意認証必須 API。 | 認証判定前に `403 {"error":"Forbidden"}`。 | password / token 検証、access log 成功行、対象操作副作用を行わない。 |
| D5 hook timeout | P4 | `pre` hook が timeout。 | build status は `hook_error` または endpoint 固有の hook error。 | pipeline を実行せず、hook log と build log に timeout を固定値で記録する。 |
| E1 token issue once | P5 | `POST /api/tokens` で token 作成。 | token 本体を作成 response に 1 回だけ含める。 | `.api_tokens` には hash だけを保存し、再取得 API では token 本体を返さない。 |
| E2 token revoke missing | P5 | 存在しない token id を `DELETE /api/tokens/{id}`。 | `404 {"error":"Not found"}`。 | `.api_tokens`、`.audit_log` を変更しない。 |
| E3 alert/tag duplicate | P5 | 同一 alert rule または tag rule を 2 回作成。 | 2 回目は `409 {"error":"Conflict"}`。 | 2 回目は該当状態ファイルと `.config_log` を変更しない。 |
| E4 pipeline config reserved arg | P5 | `extra_args` に `--src`、`--out`、`--state-dir` を含める。 | `422 {"error":"Validation failed","details":[...]}`。 | `.pipeline_config` を変更しない。 |
| E5 notes same content | P5 | 同じ `content` を 2 回 `POST /api/notes`。 | 2 回目は `No changes`。 | 2 回目は `.notes`、`.config_log` を変更しない。 |
| E6 dashboard layout invalid | P5 | 重複 widget、未知 widget、空配列を `POST /api/dashboard-layout`。 | `422 {"error":"Validation failed","details":[...]}`。 | `.dashboard_layout` を変更しない。 |

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `POST` | `/api/login` | 不要 | ログイン（セッショントークン返却） |
| `POST` | `/api/logout` | 要 | ログアウト（セッション破棄） |
| `POST` | `/api/change-password` | 要 | パスワード変更 |
| `GET` | `/api/access-log` | 要 | ログイン、ログアウト、API token 監査ログを返す |
| `GET` | `/api/sessions` | 要 | 有効セッション一覧を返す |
| `POST` | `/api/sessions/revoke-all` | 要 | 現セッション以外の全セッションを強制無効化する |
| `GET` | `/api/status` | 要 | 最終ビルド時刻・SHA・成否・実行中フラグを返す |
| `POST` | `/api/build` | 要 | 手動ビルドトリガー（`components/runner.go` を即時起動） |
| `POST` | `/api/build/force` | 要 | SHA リセットとビルドをアトミックに実行する（強制ビルド） |
| `POST` | `/api/build/cancel` | 要 | 実行中のビルドを強制停止する（`running: true` のときのみ有効） |
| `GET` | `/api/build/stream` | 要 | 実行中または直近ビルドログを SSE で配信する |
| `GET` | `/api/logs?n=100&q=<keyword>` | 要 | 最新ビルドログを n 行返す（`q` 省略時は全行） |
| `GET` | `/api/logs/search?q=<keyword>&from=<date>&to=<date>&level=<warn\|error>` | 要 | 日付範囲・重大度を指定して過去ビルドログ（`.build_logs/`）を横断検索する |
| `GET` | `/api/logs/export` | 要 | ビルドログ全件を JSON 形式でエクスポートする |
| `POST` | `/api/logs/cleanup` | 要 | 保持期間（`log_retention_days`）を超えた `.build_logs/` エントリを削除する |
| `GET` | `/api/history?page=<n>&per_page=<n>` | 要 | 過去ビルド履歴一覧をページ指定で返す（省略時: `page=1`, `per_page=20`） |
| `GET` | `/api/history/export` | 要 | ビルド履歴一覧を JSON 形式でエクスポートする |
| `GET` | `/api/history/{id}/log` | 要 | 指定ビルド ID のログを取得する |
| `GET` | `/api/history/{id}/comment` | 要 | 指定ビルドのコメントを取得する |
| `POST` | `/api/history/{id}/comment` | 要 | 指定ビルドにコメントを付与・更新する |
| `POST` | `/api/history/{id}/flag` | 要 | 指定ビルドに重要フラグを設定・解除する |
| `POST` | `/api/history/{id}/tags` | 要 | 指定ビルドのタグを置換する |
| `POST` | `/api/history/{id}/rollback` | 要 | 指定ビルド ID のスナップショットから SSH 転送を再実行する（→ §14b） |
| `GET` | `/api/sysinfo` | 要 | 出力サイトサイズ・更新日時・稼働時間を返す |
| `GET` | `/api/health`        | 不要 | 死活監視用ヘルスチェック |
| `GET` | `/api/schedule` | 要 | systemd timer の次回実行予定時刻を返す |
| `POST` | `/api/schedule/interval` | 要 | systemd タイマーのポーリング間隔を動的変更する |
| `POST` | `/api/schedule/pause` | 要 | ポーリングを一時停止する |
| `POST` | `/api/schedule/resume` | 要 | ポーリングを再開する |
| `POST` | `/api/schedule/allowed-hours` | 要 | 自動ビルド許可時間帯を設定・解除する |
| `POST` | `/api/schedule/force-interval` | 要 | `FORCE_BUILD_INTERVAL`（強制再ビルド間隔）を動的変更する |
| `POST` | `/api/schedule/cooldown` | 要 | `BUILD_COOLDOWN_SECONDS`（ビルドクールダウン秒数）を動的変更する |
| `GET` | `/api/notify-config` | 要 | Webhook 通知設定を返す |
| `POST` | `/api/notify-config` | 要 | Webhook 通知設定を更新する |
| `GET` | `/api/notify-log` | 要 | Webhook 送信履歴（日時・イベント・HTTP ステータス・成否）を返す |
| `POST` | `/api/notify-test` | 要 | Webhook にテスト通知を送信し疎通を確認する |
| `POST` | `/api/notify/weekly-summary` | 要 | 週次サマリー Webhook を即時手動送信する（過去 7 日間の統計を集計して送信） |
| `GET` | `/api/config` | 要 | サーバー設定を返す |
| `POST` | `/api/config` | 要 | サーバー設定を更新する |
| `POST` | `/api/log-level` | 要 | `components/api.go` の `log_level` を変更する |
| `GET` | `/api/config-log` | 要 | 設定変更履歴（変更日時・種別・変更前後の値）を返す |
| `GET` | `/api/pat-status` | 要 | GitHub PAT の有効性確認 |
| `POST` | `/api/pat-verify` | 要 | GitHub API を呼び出し PAT の有効性をリアルタイム検証する |
| `POST` | `/api/pat-update` | 要 | `.github_token` ファイルを更新し PAT を差し替える |
| `GET` | `/api/stats?days=7` | 要 | ビルド統計（成功率・回数・平均間隔）を返す |
| `GET` | `/api/stats/timeline?days=30` | 要 | 日別ビルド成功/失敗件数の時系列配列を返す |
| `GET` | `/api/stats/build-duration?n=20` | 要 | 過去 N 件のビルド所要時間統計（平均・最小・最大・直近リスト）を返す |
| `GET` | `/api/stats/build-trends?n=100` | 要 | `.build_trends.json` から所要時間 trend、中央値、p95、異常件数を返す |
| `GET` | `/api/output-meta` | 要 | 出力サイトのサイズ・見出し数・生成日時・前回比サイズ差分を返す |
| `GET` | `/api/repo-info` | 要 | リポジトリ設定（OWNER/REPO/BRANCH/TARGET_FILE）を返す |
| `POST` | `/api/repo-config` | 要 | リポジトリ監視設定（OWNER / REPO / BRANCH / TARGET_FILE）を更新する |
| `GET` | `/api/branch-config` | 要 | 現在有効なブランチターゲット設定を返す |
| `POST` | `/api/branch-config` | 要 | ブランチターゲット設定を `.branch_config` へ書き込む |
| `GET` | `/api/backup` | 要 | 全設定（通知設定・サーバー設定）を JSON 形式でエクスポートする |
| `POST` | `/api/restore` | 要 | JSON 形式の設定をインポートし全設定を上書き復元する |
| `GET` | `/api/dashboard` | 要 | ステータス・システム情報・統計・スケジュールを一括返却する |
| `GET` | `/api/diagnostics` | 要 | PAT・GitHub API・出力サイト・systemd・Webhook の一括自己診断結果を返す |
| `GET` | `/api/rate-limit` | 要 | GitHub API のレート制限残量・上限・リセット時刻を返す |
| `GET` | `/api/disk-usage` | 要 | ビルドログ合計・出力サイトのディスク使用量を返す |
| `GET` | `/api/webhook-events?limit=50&offset=0` | 要 | 受信 Webhook イベント一覧を新しい順にページネーション付きで返す（→ `.webhook_events.json`） |
| `POST` | `/api/webhook` | 不要（Secret 検証） | GitHub push Webhook を受信し、署名検証後にビルドをトリガーする（→ §22 Webhook 受信仕様） |
| `GET` | `/api/webhook-config` | 要 | Webhook Secret 設定状態を返す |
| `POST` | `/api/webhook-config` | 要 | Webhook Secret を設定する |
| `POST` | `/api/circuit-breaker/reset` | 要 | サーキットブレーカーをリセットする（`open: false`・`consecutive_failures: 0` に戻しポーリングを再開） |
| `GET` | `/api/snapshots` | 要 | ビルド成果物スナップショット一覧を返す |
| `GET` | `/api/snapshots/{id}/download` | 要 | 指定スナップショットをダウンロードする |
| `DELETE` | `/api/snapshots/{id}` | 要 | 指定スナップショットを削除する |
| `GET` | `/api/maintenance` | 要 | メンテナンスモード状態を返す |
| `POST` | `/api/maintenance/enable` | 要 | メンテナンスモードを有効化する |
| `POST` | `/api/maintenance/disable` | 要 | メンテナンスモードを無効化する |
| `GET` | `/api/access-control` | 要 | IP / CIDR 許可リストを返す |
| `POST` | `/api/access-control` | 要 | IP / CIDR 許可リストを置換する |
| `GET` | `/api/hooks` | 要 | ビルドフック一覧を返す |
| `POST` | `/api/hooks` | 要 | ビルドフックを追加する |
| `DELETE` | `/api/hooks/{id}` | 要 | ビルドフックを削除する |
| `GET` | `/api/hooks/{id}/log` | 要 | ビルドフック実行ログを返す |
| `GET` | `/api/alert-rules` | 要 | アラートルール一覧を返す |
| `POST` | `/api/alert-rules` | 要 | アラートルールを追加する |
| `DELETE` | `/api/alert-rules/{id}` | 要 | アラートルールを削除する |
| `GET` | `/api/tag-rules` | 要 | 自動タグ付けルール一覧を返す |
| `POST` | `/api/tag-rules` | 要 | 自動タグ付けルールを追加する |
| `DELETE` | `/api/tag-rules/{id}` | 要 | 自動タグ付けルールを削除する |
| `POST` | `/api/verify-output` | 要 | 現在の出力サイト checksum を検証する |
| `GET` | `/api/pipeline-config` | 要 | ビルドパイプライン設定を返す |
| `POST` | `/api/pipeline-config` | 要 | ビルドパイプライン設定を置換する |
| `GET` | `/api/build-chain-config` | 要 | ビルド依存チェーン設定を返す |
| `POST` | `/api/build-chain-config` | 要 | ビルド依存チェーン設定を置換する |
| `GET` | `/api/notes` | 要 | 運用ノートを返す |
| `POST` | `/api/notes` | 要 | 運用ノートを保存する |
| `GET` | `/api/smtp-config` | 要 | SMTP 設定を返す |
| `POST` | `/api/smtp-config` | 要 | SMTP 設定を保存する |
| `POST` | `/api/smtp-test` | 要 | SMTP テスト送信を行う |
| `GET` | `/api/queue` | 要 | ビルドキュー状態を返す |
| `DELETE` | `/api/queue` | 要 | 待機中ビルドキューを削除する |
| `GET` | `/api/approvals` | 要 | ビルド承認待ち entry を返す |
| `POST` | `/api/approvals/{id}/approve` | 要 | 指定承認 entry を承認し build queue へ投入する |
| `POST` | `/api/approvals/{id}/reject` | 要 | 指定承認 entry を却下する |
| `GET` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を返す |
| `POST` | `/api/dashboard-layout` | 要 | ダッシュボードウィジェット設定を置換する |
| `GET` | `/api/tokens` | 要 | API token 一覧を返す。token 本体は返さない |
| `POST` | `/api/tokens` | 要 | API token を発行する。token 本体は作成時のみ返す |
| `DELETE` | `/api/tokens/{id}` | 要 | API token を失効する |

**`POST /api/login` リクエスト / レスポンス：**
```json
// リクエスト
{ "password": "admin" }

// レスポンス（TOTP 無効）
{ "token": "<session_token>", "must_change": "prompt" }

// レスポンス（TOTP 有効）
{ "totp_required": true, "ticket": "<login_ticket>", "must_change": "none" }
```

`must_change` の有効値：`"none"`（変更不要）| `"prompt"`（促す：初回ログイン時）| `"forced"`（強制：5 回目以降。変更完了まで管理画面の操作を制限）

**`POST /api/login/totp` リクエスト / レスポンス：**
```json
// リクエスト
{ "ticket": "<login_ticket>", "code": "123456" }

// レスポンス
{ "token": "<session_token>", "must_change": "none" }
```

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
  "last_build_at": "2026-09-14T10:00:00Z",
  "last_build_status": "success",
  "last_trigger": "polling",
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "output_url": "https://example.com/",
  "running": false
}
```

`last_build_status` の有効値：`"success"` | `"failure"` | `"success_deploy_pending"` | `"skipped_no_change"` | `"skipped_cooldown"` | `"circuit_open"` | `"config_recovered"` | `"config_error"` | `"lock_skipped"` | `"none"`（初回未実行時）
`last_trigger` の有効値：§13 の `trigger` 有効値または `null`
`running` の有効値：`true`（ビルド実行中）| `false`（待機中）
`running` の判定：`.build_state.running == true` または有効な `.build_lock` が存在する場合に `true` を返す。`.build_state.running == false` かつ `.build_lock` が存在しない場合は `false` を返す。形式不正または PID 判定不能な `.build_lock` が存在する場合は、状態競合として `running: true` を返し、API 側でロックを上書きしない。

**`GET /api/logs` レスポンス例：**
```json
{ "lines": ["2026-09-14T10:00:00Z [INFO] Build start", "..."] }
```

**`GET /api/logs/export` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00Z",
  "lines": ["2026-09-14T10:00:00Z [INFO] Build start", "..."]
}
```

**`GET /api/history` レスポンス例：**
```json
{
  "total": 42, "page": 1, "per_page": 20, "pages": 3,
  "history": [
    { "id": "b001", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success", "output_size_bytes": 2048576, "trigger": "polling", "duration_seconds": 42, "flagged": false, "tags": ["release"] },
    { "id": "b002", "build_at": "2026-09-14T18:30:00Z", "sha": "def456", "status": "failure", "output_size_bytes": null,    "trigger": "manual", "duration_seconds": 7,  "flagged": true,  "tags": [] }
  ]
}
```

`page` は 1 始まり。`per_page` の最大値は 100。`trigger` は §13 の固定値のみ許可する。`tag` はタグ検証と同じ文字列制約を適用する。`flagged` は `"true"` または `"false"` のみ許可する。未知 query、範囲外の `page` / `per_page`、不正な `trigger` / `tag` / `flagged` は `422` とする。範囲外ページを指定した場合は `history: []` を返す。

`id` はビルド実行時に生成するユニーク識別子（形式：`b{YYYYMMDDHHmmss}`）。`.build_logs/{id}.json` に対応するログファイルが保存される。

**`POST /api/build` / `POST /api/build/force` レスポンス例：**
```json
// 即時開始
{ "message": "Build started", "build_id": "b20260915100500", "queued": false }

// 実行中のため queue へ追加
{ "message": "Build queued", "queued": true }
```

SHA キャッシュのクリアだけを行う専用 API は定義しない。強制再ビルドは必ず `POST /api/build/force` を使用し、SHA reset と build trigger を同一ロック内で実行する。

**`GET /api/sysinfo` レスポンス例：**
```json
{
  "output_size_bytes": 2048576,
  "output_mtime": "2026-09-15T10:00:00Z",
  "uptime_seconds": 86400
}
```

**`GET /api/schedule` レスポンス例：**
```json
{ "next_run_at": "2026-09-15T10:05:00Z", "interval": "5min", "paused": false, "allowed_hours": { "from": 9, "to": 18 } }
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

> **責務分担：** Webhook 通知の**送信責務は `components/runner.go`** にある。`components/runner.go` はビルド完了時に `.notify_config` を読み込んで Webhook を送信する。`components/api.go`（通知 API）は設定の読み書きのみを担い、自身では通知を送信しない。

**`GET /api/notify-config` レスポンス例：**
```json
{
  "webhooks": [
    { "url": "https://hooks.example.com/...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": null }
  ],
  "channels": [
    { "id": "n001", "type": "webhook", "label": "メイン", "enabled": true, "on": ["failure"], "config": { "url": "https://hooks.example.com/..." }, "retry_count": 2, "retry_interval_seconds": 30 }
  ],
  "on": ["failure"],
  "summary": { "enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1 },
  "email": { "enabled": false, "to": [], "on": [] }
}
```

`secret`：Webhook 署名シークレット。未設定時は `null`、設定済み時は `"***"`（マスク）を返す（→ 16E 参照）。

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"deploy_failure"`（転送失敗時）| `"weekly_summary"`（定期サマリー送信時）| `"approval_required"`（承認待ち発生時）| `"duration_anomaly"`（所要時間異常時）| `"config_corrupt"`（設定破損復旧時）。複数指定可。

`summary`：週次サマリー通知の設定。`enabled: true` のとき指定曜日・時刻で統計サマリーを Webhook 送信する。`interval` の有効値は `"weekly"` 固定。`hour` は 0〜23（UTC）。`day_of_week` は 0 = 日曜〜6 = 土曜。

**`POST /api/notify/weekly-summary` レスポンス例：**
```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

即時週次サマリー送信。Webhook 未設定または無効時は `422` を返す。

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
{ "log_max_lines": 500, "history_max_count": 100, "build_timeout_seconds": 300, "log_retention_days": 30, "log_archive_after_days": 0, "log_level": "INFO", "pat_expires_at": null, "snapshots_keep": 5, "queue_max_size": 3, "build_retry_max": 0, "build_retry_base_seconds": 5, "commit_status_enabled": false, "commit_status_context": "Adlaire CI", "commit_status_target_url": null, "build_trend_keep_count": 1000, "duration_anomaly": { "enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5 } }
```

`pat_expires_at`：PAT の有効期限日（`YYYY-MM-DD` 形式）。`null` = 未設定。`GET /api/diagnostics` の `pat` 項目で 7 日以内なら `"warn"`、期限当日以前なら `"error"` に変更。

`POST /api/config` で更新可能なキーは `log_max_lines`、`history_max_count`、`build_timeout_seconds`、`log_retention_days`、`log_archive_after_days`、`log_level`、`pat_expires_at`、`snapshots_keep`、`queue_max_size`、`build_retry_max`、`build_retry_base_seconds`、`commit_status_enabled`、`commit_status_context`、`commit_status_target_url`、`build_trend_keep_count`、`duration_anomaly` に限定する。未知キーを含む場合は `422` を返し、既存設定を変更しない。

**`GET /api/health` レスポンス例：**
```json
{
  "status": "ok",
  "last_build_at": "2026-09-15T10:00:00Z",
  "last_build_status": "success",
  "last_deploy_at": "2026-09-15T10:01:00Z",
  "last_deploy_status": "success",
  "pending_transfers": 0,
  "uptime_seconds": 86400
}
```

- `status`：常に `"ok"`（サーバーが応答している限り）
- `last_build_at`：最終ビルド完了日時（未実行時 `null`）
- `last_build_status`：`GET /api/status` の `last_build_status` と同じ値
- `last_deploy_at`：最終 SSH 転送完了日時（未実行時 `null`）
- `last_deploy_status`：`"success"` | `"failure"` | `"skipped"` | `"none"`
- `pending_transfers`：ペンディングキューのエントリ数
- `uptime_seconds`：`components/api.go` 起動からの経過秒数

**`GET /api/pat-status` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:00:00Z" }
```

**`GET /api/access-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "result": "success" },
    { "at": "2026-09-15T09:00:00Z", "result": "failure" }
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
  "target_file": "docs"
}
```

**`GET /api/backup` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00Z",
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
{ "message": "Build started", "build_id": "b20260915100500", "queued": false }
```

**`POST /api/pat-verify` レスポンス例：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:05:00Z", "scopes": ["contents:read"] }
```

**`GET /api/history/{id}/log` レスポンス例：**
```json
{
  "id": "b001",
  "build_at": "2026-09-15T10:00:00Z",
  "sha": "abc123",
  "status": "success",
  "output_size_bytes": 2048576,
  "trigger": "polling",
  "comment": null,
  "flagged": false,
  "tags": ["release"],
  "lines": ["2026-09-15T10:00:00Z [INFO] Build start", "..."]
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
    { "at": "2026-09-15T10:00:00Z", "event": "failure", "http_status": 200, "result": "success", "attempt": 1, "error": null },
    { "at": "2026-09-14T18:30:00Z", "event": "start",   "http_status": 500, "result": "failure", "attempt": 3, "error": "HTTP 500" }
]}
```

`event` の有効値：`"start"` | `"success"` | `"failure"` | `"weekly_summary"`（`GET /api/notify-config` の `on` と同一）。`result` の有効値：`"success"` | `"failure"`（Webhook 送信の成否）。

**`GET /api/sessions` レスポンス例：**
```json
{ "sessions": [
    { "created_at": "2026-09-15T09:00:00Z", "expires_at": "2026-09-15T17:00:00Z", "current": true },
    { "created_at": "2026-09-15T08:00:00Z", "expires_at": "2026-09-15T16:00:00Z", "current": false }
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

`Content-Type: text/event-stream` で接続を維持し、ビルドログを逐次配信する。認証は他 API と同じ `Authorization: Bearer {SESSION_TOKEN}` ヘッダーで行う。セッショントークンを query parameter、Cookie、body で受け付けてはならない。

```
data: {"type": "log",  "line": "2026-09-15T10:00:01Z [INFO] Build start"}

data: {"type": "log",  "line": "2026-09-15T10:00:42Z [INFO] Build success"}

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
    "last_build_at": "2026-09-15T10:00:00Z",
    "last_build_status": "success",
    "output_url": "https://example.com/",
    "running": false
  },
  "sysinfo": {
    "output_size_bytes": 2048576,
    "output_mtime": "2026-09-15T10:00:00Z",
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
  "schedule": { "next_run_at": "2026-09-15T10:05:00Z", "interval": "5min", "paused": false },
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
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "lines": ["2026-09-15T10:00:01Z [ERROR] Build failed"] },
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
  "mtime": "2026-09-15T10:00:00Z",
  "heading_count": 342,
  "size_diff_bytes": 1024,
  "tables_count": 128,
  "code_blocks_count": 64,
  "build_warnings": ["未対応記法: admonition (3箇所)"],
  "size_warn": false,
  "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "build_id": "b20260915100000",
  "commit_sha": "abc1234",
  "build_at": "2026-09-15T10:00:00Z"
}
```

`size_diff_bytes`：前回ビルド時との差分（正＝増加、負＝減少、`null`＝比較不能）。
前回サイズは `.build_history` の直近エントリに記録された `output_size_bytes` フィールドから取得する。

`tables_count` / `code_blocks_count`：直近ビルドの変換レポート（§8）より取得。ビルド前は `null`。
`build_warnings`：直近ビルドで発生した警告メッセージの配列（§8 参照）。ビルド前は空配列 `[]`。
`build_id` / `commit_sha` / `build_at`：直近ビルドログの `build_meta` を優先し、不在の場合は出力 HTML の meta tag を読み取る。どちらにも存在しない場合は空文字を返す。
値は `components/runner.go` が `.build_logs/{id}.json` または `.build_logs/archive/{id}.json.gz` から最新エントリを読み取って返す。

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
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "duration_seconds": 42, "status": "success" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00Z", "duration_seconds": 7,  "status": "failure" }
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
      "timestamp": "2026-09-15T10:00:00Z",
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

`source` は設定の出所を示す。`.branch_config` ファイルが存在する場合は `"file"`、存在しない場合（`BRANCH_TARGETS` デフォルト値を使用）は `"default"` を返す。

**`POST /api/branch-config` リクエスト / レスポンス：**

```json
// リクエスト（GET /api/branch-config の branches と同一形式）
{
  "branches": [
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
  "started_at": "2026-09-15T09:59:18Z",
  "finished_at": "2026-09-15T10:00:00Z",
  "duration_seconds": 42,
  "status": "success",
  "commit_sha": "abc123def456",
  "commit_message": "fix: typo in §4.3 description",
  "commit_author": "Kazuhiro Kurata",
  "commit_at": "2026-09-15T09:58:00Z",
  "size_warn": false,
  "broken_links": 0,
  "heading_skips": 0,
  "reading_time": 87
}
```

`started_at` は `pipeline.sh` 実行開始時刻、`finished_at` は完了（または失敗）時刻。`duration_seconds` は整数（小数点以下切り捨て）。
`commit_sha` / `commit_message` / `commit_author` / `commit_at` はコミット情報取得 API（§13）の結果を記録する。API 失敗時は `null`。
`size_warn` は出力サイト合計サイズが `OUTPUT_SIZE_WARN_MB` 超過時 `true`、それ以外 `false`。`OUTPUT_SIZE_WARN_MB = 0` の場合は常に `false`。

**`GET /api/diagnostics` レスポンス例：**
```json
{
  "checked_at": "2026-09-15T10:00:00Z",
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
{ "limit": 5000, "remaining": 4823, "reset_at": "2026-09-15T11:00:00Z", "used": 177 }
```

**`GET /api/disk-usage` レスポンス例：**
```json
{
  "build_logs_bytes": 10485760,
  "build_logs_count": 42,
  "build_logs_archive_bytes": 2097152,
  "build_logs_archive_count": 18,
  "output_file_bytes": 2048576,
  "total_bytes": 14631448
}
```

`build_logs_bytes` / `build_logs_count` は `.build_logs/{id}.json` のみを集計する。`build_logs_archive_bytes` / `build_logs_archive_count` は `.build_logs/archive/{id}.json.gz` のみを集計する。`total_bytes` は通常 build log、archive build log、出力サイトの合計 bytes とする。

**`GET /api/config-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "type": "server_config", "diff": { "log_max_lines": [500, 1000] }, "diff_text": "- log_max_lines: 500\n+ log_max_lines: 1000" },
    { "at": "2026-09-14T18:00:00Z", "type": "notify_config", "diff": { "enabled": [false, true] },    "diff_text": "- enabled: false\n+ enabled: true" },
    { "at": "2026-09-13T12:00:00", "type": "repo_config",   "diff": { "branch": ["main", "develop"] }, "diff_text": "- branch: main\n+ branch: develop" }
]}
```

`type` の有効値：`"server_config"` | `"notify_config"` | `"repo_config"`
`diff` の形式：`{ フィールド名: [変更前, 変更後] }`。設定変更時に `.config_log` へ追記する。
`diff_text`：`diff` を `"- key: old_value\n+ key: new_value"` 形式の文字列に変換したフィールド。複数フィールド変更時は行を連結する。設定変更時に `diff` と同時に記録する。

**`GET /api/history/{id}/comment` レスポンス例：**
```json
{ "id": "b20260914183000", "comment": "ネットワーク障害による失敗。再ビルド済み。", "updated_at": "2026-09-14T19:00:00Z" }
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
    { "id": "tok000001", "label": "監視用", "scopes": ["read"], "created_at": "2026-09-15T10:00:00Z", "last_used_at": "2026-09-15T11:00:00Z", "expires_at": null, "revoked_at": null }
]}
```

**`POST /api/tokens` リクエスト / レスポンス：**
```json
// リクエスト
{ "label": "監視用", "scopes": ["read"], "expires_at": null }
// レスポンス: 201
{ "id": "tok000001", "token": "act_...", "label": "監視用", "scopes": ["read"], "created_at": "2026-09-15T10:00:00Z", "expires_at": null }
```

`token` はレスポンス時のみ返却し、以後は取得不可。`scopes` の有効値は `read`、`trigger`、`operate`、`config`、`admin` とする。管理 session は全 API 操作を許可し、API token は指定 scope の範囲だけを許可する。トークンは `Authorization: Bearer <token>` ヘッダーで送信する。

**`DELETE /api/tokens/{id}` レスポンス例：**
```json
{ "message": "Token revoked" }
```

---

### スナップショット（14A）

ビルド成功時に出力サイトを `.snapshots/` へ自動保存する。保持世代数は `GET /api/config` の `snapshots_keep`（デフォルト `5`、`0` = 機能無効）で制御し、超過した古い世代は自動削除する。

**`GET /api/snapshots` レスポンス例：**
```json
{ "snapshots": [
    { "id": "snap001", "build_id": "b20260915100000", "saved_at": "2026-09-15T10:00:00Z", "size_bytes": 2048576 },
    { "id": "snap002", "build_id": "b20260914183000", "saved_at": "2026-09-14T18:30:00Z", "size_bytes": 2031616 }
]}
```

**`GET /api/snapshots/{id}/download`**
バイナリレスポンス。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="site"` を付与する。

**`DELETE /api/snapshots/{id}` レスポンス例：**
```json
{ "message": "Snapshot deleted" }
```

**`POST /api/history/{id}/rollback` リクエスト / レスポンス：**
```json
// リクエスト: なし（パスパラメーターのみ）
// レスポンス: 202
{ "message": "Rollback started", "build_id": "b20260914183000" }
```

- `.snapshots/{id}/` が存在しない場合は `404 Not Found` を返す
- ロールバックは非同期で SSH 転送を実行する（`running: true` 中は `409 Conflict` を返す）
- 転送成功時は `.build_history` に `trigger: "rollback"` のエントリを追記する

**スナップショット保存・削除固定契約：**

| 処理 | 固定仕様 |
|------|----------|
| snapshot id | §22.0e.2 の `snap{YYYYMMDDHHmmss}` 形式。build id を使い回さない。 |
| 保存対象 | 出力サイトディレクトリ配下の通常ファイルだけ。symlink、socket、device、隠し一時ファイルは保存対象外。 |
| archive 形式 | `.snapshots/{snapshot_id}/site.tar.gz` と `.snapshots/{snapshot_id}/meta.json` を作成する。 |
| `meta.json` | `id`、`build_id`、`saved_at`、`size_bytes`、`file_count`、`output_sha256` を必須 key とする。 |
| 世代削除 | 新 snapshot 作成成功後に `saved_at` 昇順で超過分だけ削除する。削除失敗は build 成功を失敗へ反転しないが、WARN log に固定コード `SNAPSHOT_PRUNE_FAILED` を出す。 |
| rollback | 対象 snapshot の `site.tar.gz` を展開して転送し、新しい build id で `.build_logs/{id}.json` と `.build_history` を作成する。 |

rollback 開始時は `.build_lock` を取得し、取得できない場合は `409 {"error":"Build is running"}` を返す。`.build_lock` 取得後に `.build_state.running=true`、`current_build_id=<new_id>` を保存し、転送完了後に finalizer で `running=false` とする。rollback は queue に積まない。

---

### Webhook 受信仕様（22-W）

`POST /api/webhook` は GitHub からの push イベントを受信し、署名検証後にビルドをトリガーする。認証ヘッダー（`Authorization: Bearer`）は不要だが、`X-Hub-Signature-256` ヘッダーによる HMAC-SHA256 署名検証が必須である。

**署名検証：**
```go
// components/api.go の実装例
mac := hmac.New(sha256.New, []byte(secret))
mac.Write(body)
expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
if subtle.ConstantTimeCompare([]byte(expected), []byte(requestHeader.Get("X-Hub-Signature-256"))) != 1 {
    return 403
}
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
{"timestamp": "2026-09-15T10:00:00Z", "delivery_id": "abc-123-def", "event": "push", "ref": "refs/heads/main", "sha": "abc123def456", "build_triggered": true}
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

- `components/runner.go` 側の `FORCE_BUILD_INTERVAL` を動的変更する（`.server_config` に保存し、起動時に読み込む）
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

- `components/runner.go` 側の `BUILD_COOLDOWN_SECONDS` を動的変更する（`.server_config` に保存し、起動時に読み込む）
- `seconds` は 0 以上の整数。0 で機能無効化

---

### メンテナンスモード（14B）

メンテナンスモード有効中は手動ビルド（`POST /api/build`・`POST /api/build/force`）・スケジュールビルドの両方を拒否し、`503 Service Unavailable` + `{ "error": "maintenance" }` を返す。`GET /api/health` は制限対象外とする。

**`GET /api/maintenance` レスポンス例：**
```json
{ "enabled": false, "reason": null, "since": null }
```
メンテナンス中は `{ "enabled": true, "reason": "定期メンテナンス", "since": "2026-09-15T10:00:00Z" }`。

**`POST /api/maintenance/enable` リクエスト / レスポンス：**
```json
// リクエスト
{ "reason": "定期メンテナンス" }
// レスポンス: 200
{ "message": "Maintenance mode enabled", "since": "2026-09-15T10:00:00Z" }
```

**`POST /api/maintenance/disable` レスポンス：**
```json
{ "message": "Maintenance mode disabled" }
```

**メンテナンス更新固定契約：**

| API | 入力検証 | 保存値 | no-op | 失敗時 |
|-----|----------|--------|-------|--------|
| `POST /api/maintenance/enable` | `reason` 必須、1〜500 文字、前後空白除去後空は禁止 | `enabled:true`、`reason`、`since=now` | 既に同一 reason で enabled の場合は状態を変更せず `{ "message":"No changes","since":"<existing since>" }` | 検証失敗 `422`、保存失敗 `500` |
| `POST /api/maintenance/disable` | body 禁止 | `enabled:false`、`reason:null`、`since:null` | 既に disabled の場合は状態を変更せず `{ "message":"No changes" }` | 保存失敗 `500` |

メンテナンス判定は §22.0e の共通判定順に従い、`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`、`POST /api/history/{id}/rollback` を拒否対象とする。設定参照系 GET、認証、ログ参照、メンテナンス解除は拒否しない。

**メンテナンス判定・副作用固定契約：**

| 対象 | 判定 | 拒否時 |
|------|------|--------|
| `POST /api/build` | 認証、rate limit、入力検証後、`.build_state` 更新前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。queue 追加、build id 採番、`.build_state` 更新を行わない。 |
| `POST /api/build/force` | SHA reset 前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。SHA cache、queue、`.build_state` を変更しない。 |
| `POST /api/webhook` | 署名検証、payload 検証後、`.webhook_events.json` 追記前に確認する。 | `503 {"error":"maintenance"}`。event log と queue を変更しない。 |
| `POST /api/history/{id}/rollback` | snapshot 存在確認後、rollback build log 作成前に確認する。 | `503 {"error":"maintenance"}`。history、build log、pending transfer を変更しない。 |
| runner 定期起動 | lock 取得後、差分検出前に確認する。 | build せず `.build_status.json.status="skipped_maintenance"` を保存し、`.last_sha` を更新しない。 |

`POST /api/maintenance/enable` と `POST /api/maintenance/disable` の保存順は、`.maintenance` atomic write → `.config_log` 追記 → response とする。`.config_log` 追記失敗時は `500` を返し、保存済み `.maintenance` は巻き戻さない。同一状態 no-op では `.maintenance`、`.config_log`、`.audit_log` を変更しない。

**メンテナンス fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| maintenance-build-deny | enabled 中に `POST /api/build` | `503`、queue 差分なし、build id なし。 |
| maintenance-force-deny | enabled 中に `POST /api/build/force` | `503`、SHA cache 差分なし。 |
| maintenance-webhook-deny | enabled 中に署名済み webhook | `503`、event log と queue 差分なし。 |
| maintenance-disable-noop | disabled 中に disable | `200 No changes`、状態ファイル差分なし。 |

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

**アクセス制限判定固定契約：**

1. 接続元 IP は `X-Forwarded-For`、`X-Real-IP` を使わず、`net/http.Request.RemoteAddr` から取得する。
2. `RemoteAddr` が parse 不能な場合は `403 {"error":"Forbidden"}` を返す。
3. `.access_control` 不在、または `allow:[]` は全許可とする。
4. CIDR は `net.ParseCIDR`、単一 IPv4 は `net.ParseIP` で判定する。
5. 判定失敗時は認証処理より前に `403` を返し、password や token 検証を実行しない。

`POST /api/access-control` は正規化後の `allow` 配列が既存値と一致する場合、`.access_control` と `.config_log` を変更せず `{ "message":"No changes","allow":[...] }` を返す。

**アクセス制御更新・拒否固定契約：**

| 項目 | 仕様 |
|------|------|
| 判定順 | path / method 判定後、body parse 前、認証前に実行する。拒否時は password、session token、API token、rate limit state を検証または更新しない。 |
| 対象外 | `GET /api/health` だけを対象外とする。静的 admin UI、SDK JS、その他 `/api/` 以外の配信は本契約の対象外。 |
| allow 正規化 | 重複除去、辞書順 sort、単一 IPv4 は canonical 文字列、CIDR は `IP/mask` 表記へ正規化する。 |
| IPv6 | 初期実装では保存不可。IPv6 literal または IPv6 CIDR は `422`。 |
| private / public | private address に限定しない。入力が IPv4 または IPv4 CIDR として妥当なら保存可能。 |
| 保存順 | `.access_control` atomic write → `.config_log` 追記 → response。 |
| `.config_log` 失敗 | `500`。保存済み `.access_control` は巻き戻さない。 |
| 破損時 | §22.0a に従い初期値で再生成し、制限なしとして扱う。 |

**アクセス制御 fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| access-allow-empty | `.access_control.allow=[]` | 任意 IP の API が認証処理へ進む。 |
| access-deny-before-auth | allow 不一致 IP で `POST /api/login` | `403`、`.access_log`、`.audit_log`、rate state 差分なし。 |
| access-normalize | 重複 allow を保存 | sort / 重複除去後の配列を返し `.config_log` 記録。 |
| access-ipv6-reject | IPv6 literal を保存 | `422`、状態差分なし。 |

---

### ビルドフック（14E）

ビルド実行の直前（`pre`）・直後（`post`）に事前登録したコマンド引数配列を実行する。フック設定は `.hooks` に保存する。外部入力文字列をシェルへ渡す実装は禁止し、Go 標準ライブラリ `os/exec` の `exec.CommandContext(args[0], args[1:]...)` で実行する。

- `pre` フックが失敗（`exit_code != 0`）し `abort_on_failure: true` の場合、ビルドを中断しステータスを `hook_error` とする。
- `post` フックは `abort_on_failure` 設定に関わらずビルド結果（`success` / `failure`）を変更しない。
- フックの実行ログは `.build_logs/{build_id}_hook_{id}.json` に保存する。

**`GET /api/hooks` レスポンス例：**
```json
{ "hooks": [
    { "id": "h20260915100500", "phase": "pre",  "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true, "timeout_seconds": 300, "created_at": "2026-09-15T10:05:00Z" },
    { "id": "h20260915100600", "phase": "post", "command_args": ["echo", "build end"],   "enabled": true, "abort_on_failure": false, "timeout_seconds": 300, "created_at": "2026-09-15T10:06:00Z" }
]}
```

**`POST /api/hooks` リクエスト / レスポンス：**
```json
// リクエスト
{ "phase": "pre", "command_args": ["echo", "build start"], "abort_on_failure": true }
// レスポンス: 201
{ "id": "h20260915100500", "phase": "pre", "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true, "timeout_seconds": 300, "created_at": "2026-09-15T10:05:00Z" }
```

`phase` の有効値は `"pre"` または `"post"`。`command_args[0]` は絶対パス、または `PATH` 解決可能なコマンド名とする。`command_args` に空文字、NUL 文字、改行を含めてはならない。

**`DELETE /api/hooks/{id}` レスポンス：**
```json
{ "message": "Hook deleted" }
```

**`GET /api/hooks/{id}/log` レスポンス例：**
```json
{ "id": "h20260915100500", "runs": [
    { "build_id": "b20260915100000", "ran_at": "2026-09-15T10:00:00Z", "exit_code": 0, "output": "build start\n" },
    { "build_id": "b20260914183000", "ran_at": "2026-09-14T18:30:00Z", "exit_code": 1, "output": "Error: command not found\n" }
]}
```
`runs` は直近 20 件を返す（新しい順）。

**フック実行・保存固定契約：**

| 項目 | 仕様 |
|------|------|
| 実行順 | `phase` ごとに `.hooks.hooks` の配列順。`pre` は pipeline 前、`post` は pipeline/deploy/snapshot 後。 |
| disabled | `enabled:false` は読み飛ばし、hook log を作成しない。 |
| timeout | `timeout_seconds` 超過時は process を kill し、`exit_code:null`、`timed_out:true` として hook log を保存する。 |
| stdout/stderr | 最大各 10000 文字。超過分は末尾切り捨て、`truncated:true` を保存する。 |
| hook log | `.build_logs/{build_id}_hook_{hook_id}.json` に JSON object で保存し、同一 build/hook の再実行時は上書きせず `runs` へ追記する。 |
| pre abort | `pre` 失敗かつ `abort_on_failure:true` の場合、pipeline を実行せず build status を `hook_error` とする。 |
| secret mask | stdout、stderr、保存済み output、server log、通知 payload へ保存する前に runner の secret mask を適用する。 |
| log 保存失敗 | `pre` hook では build 本体を開始せず `hook_error`。`post` hook では build 結果を維持し、server log に `HOOK_LOG_WRITE_FAILED` を出す。 |

**`.hooks` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | §22.0e.2 の hook id。 |
| `phase` | string | 必須 | `"pre"` または `"post"`。 |
| `command_args` | string[] | 必須 | 1〜20 件。各値は NUL、LF、CR 禁止。 |
| `enabled` | boolean | 必須 | boolean。作成時 `true` 固定。 |
| `abort_on_failure` | boolean | 必須 | boolean。 |
| `timeout_seconds` | integer | 必須 | 1〜3600。省略時 300。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |

`.hooks` に未知 key、必須 key 不足、型不一致、不正 phase、不正 command、重複 id がある場合、`GET /api/hooks`、`POST /api/hooks`、`DELETE /api/hooks/{id}` は `500 {"error":"Internal server error"}` を返す。破損内容、command_args の secret らしき値、stdout/stderr は response と log に出さない。

**hooks API 更新順：**

| API | 更新順 | 失敗時 |
|-----|--------|--------|
| `POST /api/hooks` | 入力検証 → `.hooks` lock → id 採番 → record append → `.hooks` atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。追加済み record は巻き戻さない。 |
| `DELETE /api/hooks/{id}` | path id 検証 → `.hooks` lock → 対象存在確認 → record 削除 → `.hooks` atomic write → `.config_log` 追記 → response | 対象不在は `404`。`.config_log` 失敗時は `500`、削除済み record は巻き戻さない。 |
| `GET /api/hooks/{id}/log` | path id 検証 → `.hooks` で存在確認 → `.build_logs/*_hook_{id}.json` を新しい順で最大 20 件読込 → response | hook 不在は `404`。個別 hook log 破損はその file を除外し、server log に固定コードを出す。 |

hook log JSON は `{ "hook_id", "build_id", "phase", "started_at", "finished_at", "duration_seconds", "exit_code", "timed_out", "stdout", "stderr", "truncated" }` を必須 key とする。`GET /api/hooks/{id}/log` の `output` は `stdout + stderr` をこの順で連結した表示用互換値とし、保存時点で secret mask 済みの値だけを返す。

**hooks fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| hook-pre-success | pre hook exit 0 | pipeline 実行、hook log 保存、secret mask 済み。 |
| hook-pre-abort | pre hook exit 1 / abort true | pipeline 未実行、status `hook_error`、history に `failure_category:"hook_error"`。 |
| hook-pre-warn | pre hook exit 1 / abort false | build 継続、hook log に exit code。 |
| hook-post-failure | build success 後 post hook failure | build success 維持、hook log 保存。 |
| hook-timeout | timeout 超過 | process kill、`timed_out:true`、`exit_code:null`。 |
| hook-log-write-failure | pre hook log 保存失敗 | build 本体未実行、`hook_error`。 |
| post failure | build status を変更しない。WARN log と hook log だけを残す。 |

保存する hook log file は 1 実行 1 JSON object とし、`hook_id`、`build_id`、`phase`、`started_at`、`finished_at`、`duration_seconds`、`exit_code`、`timed_out`、`stdout`、`stderr`、`truncated` を必須 key とする。`GET /api/hooks/{id}/log` は複数 file を集約し、response の `runs[]` へ `build_id`、`ran_at`、`exit_code`、`output` を返す。

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

**アラート評価固定契約：**

`.alert_rules` 不在時は `rules:[]` と扱う。`GET /api/dashboard` は rule 配列順で評価し、条件一致した rule だけ `alerts` へ追加する。alert object は `id`、`level`、`message`、`metric`、`value`、`threshold` を必須 key とする。評価に必要な metric が算出不能な rule は alert 化せず、WARN log `ALERT_METRIC_UNAVAILABLE` を出す。

`POST /api/alert-rules` は `metric`、`operator`、`threshold`、`level`、`message` の正規化後値が既存 rule と一致する場合、`409 {"error":"Conflict"}` を返す。`DELETE` は対象 id 不在時 `404` とし、部分削除は行わない。

---

### 自動タグ付けルール（15B）

ビルド完了時に条件式を評価し、マッチしたルールのタグを `.build_history` のエントリへ自動追記する（手動タグと共存する）。`.tag_rules` に保存する。

条件式で使用可能な変数：`status`（`"success"` / `"failure"`）/ `duration_seconds`（整数）/ `trigger`（`"polling"` / `"force_interval"` / `"manual"` / `"webhook"` / `"retry_pending_transfer"` / `"startup_config_integrity"` / `"rollback"` / `"local_watch"` / `"approval"`）
演算子：`==`・`!=`・`>`・`<`・`>=`・`<=`

**`GET /api/tag-rules` レスポンス例：**
```json
{ "rules": [
    { "id": "t001", "condition": "status == 'failure'", "tags": ["要確認"] },
    { "id": "t002", "condition": "duration_seconds > 120", "tags": ["低速"] },
    { "id": "t003", "condition": "trigger == 'force_interval'", "tags": ["強制実行"] }
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

**自動タグ評価固定契約：**

タグ評価は build log の最終 status / duration / trigger が確定した後、`.build_history` 追記前に行う。手動タグと自動タグが重複した場合は 1 件に正規化し、既存順を保持したうえで自動タグを末尾へ追加する。条件式 parse 失敗を含む破損 rule がある場合、その rule を無視せず build を `failure` にし、ERROR log `TAG_RULE_INVALID` を出す。

`POST /api/tag-rules` は同一 `condition` と同一 `tags` の rule が存在する場合、`409 {"error":"Conflict"}` を返す。

---

### 出力サイトチェックサム（15C）

ビルド成功時に出力サイト配下の全通常ファイルから manifest SHA-256 を算出し、`.build_history` の該当エントリに `output_sha256` として記録する。manifest は `relative_path + "\n" + file_sha256 + "\n"` を相対パス昇順で連結した文字列とし、その SHA-256 hex を `output_sha256` とする。

**`GET /api/output-meta` レスポンス変更（`sha256` / build meta フィールド追加）：**
```json
{
  "size_bytes": 2048576,
  "mtime": "2026-09-15T10:00:00Z",
  "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "build_id": "b20260915100000",
  "commit_sha": "abc1234",
  "build_at": "2026-09-15T10:00:00Z"
}
```

**`GET /api/history/{id}/log` レスポンス変更（`output_sha256` フィールド追加）：**
```json
{ "id": "b001", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success",
  "output_size_bytes": 2048576, "output_sha256": "e3b0c44298fc1c149afbf4c8996fb924...",
  "trigger": "polling", "comment": null, "flagged": false, "tags": ["release"],
  "lines": ["2026-09-15T10:00:00Z [INFO] Build start", "..."] }
```

**`POST /api/verify-output` レスポンス例：**
```json
// 一致時
{ "match": true, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "e3b0c44298fc1c149afbf4c8996fb924..." }
// 不一致時
{ "match": false, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "f4a2d5591c8f3a742f902e3b6f7c1c3d..." }
```
`expected` は `.build_history` の最終成功エントリに記録された `output_sha256`。現在の出力サイトが存在しない場合は `404` を返す。

**checksum 算出固定契約：**

対象 path は `/` 区切りの相対 path とし、先頭 `/`、`..`、NUL、改行を含む path は manifest 算出前に build 失敗とする。ディレクトリ、symlink、device、socket は manifest に含めない。空ディレクトリの checksum は空文字列に対する SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` とする。

`POST /api/verify-output` は状態ファイルを変更しない。`.build_history` に成功履歴がない場合、`404 {"error":"Not Found"}` を返す。

---

### ビルドパイプライン設定（15D）

`components/runner.go` がビルド実行時に `.pipeline_config` を読み込み、`components/builder.go` の呼び出しに `extra_args`・`env` を適用する。`.pipeline_config` に保存する。

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

`POST /api/pipeline-config` は `.pipeline_config` 全体置換とし、部分更新を許可しない。`extra_args` または `env` のいずれかが欠ける場合は `422`。正規化後値が既存値と一致する場合は `.pipeline_config` と `.config_log` を変更せず `{ "message":"No changes" }` を返す。

runner は build 開始後、builder command 組み立て直前に `.pipeline_config` を 1 回だけ読む。読込不能または schema 不正は build を開始せず `failure` とし、SHA cache を更新しない。

---

### 運用ノート（15E）

システム全体の運用メモを Markdown テキストで保存・取得する。`.notes` に保存する。認証必須。

**`GET /api/notes` レスポンス例：**
```json
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00\nPAT 更新期限: 2026-12-01", "updated_at": "2026-09-15T10:00:00Z" }
```
初回（未作成）時：`{ "content": "", "updated_at": null }`

**`POST /api/notes` リクエスト / レスポンス：**
```json
// リクエスト
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00" }
// レスポンス: 200
{ "message": "Notes updated", "updated_at": "2026-09-15T10:00:00Z" }
```

`.notes` は UTF-8 text として保存し、JSON ではない。`POST /api/notes` の `content` は 0〜100000 文字、NUL 禁止とする。保存時は本文をそのまま `.notes` に atomic write し、更新時刻は `.config_log` の `at` を返す。既存本文と一致する場合は `.notes` と `.config_log` を変更せず `{ "message":"No changes","updated_at":null }` を返す。

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

**SMTP 更新・送信固定契約：**

| 処理 | 仕様 |
|------|------|
| 更新 | `.smtp_config` → 必要時 `.smtp_secret` → `.config_log` の順に書く。途中失敗時は未処理ファイルを書かない。 |
| 削除 | `password:null` は `.smtp_secret` 削除。削除対象が不在なら no-op。 |
| GET | `.smtp_secret` の存在だけを `password_set` で返し、password 本体は返さない。 |
| 送信 timeout | 接続、TLS、送信全体を合計 30 秒で timeout する。 |
| 送信 log | `.notify_log` に `type:"smtp_test"`、`result`、`message`、`at` を追記する。password は記録しない。 |

`enabled:true` にする場合は `host`、`port`、`from`、`to` 1 件以上を必須とする。`POST /api/smtp-test` は `enabled:false`、宛先なし、secret 必須構成で `.smtp_secret` 不在のいずれも `422 {"error":"SMTP not configured"}` を返す。

**SMTP 更新詳細：**

| ケース | `.smtp_config` | `.smtp_secret` | `.config_log` | response |
|--------|----------------|----------------|---------------|----------|
| config のみ変更 | 保存 | 変更なし | mask 済み diff 追記 | `200 {"message":"SMTP config updated"}` |
| password 追加 / 変更 | 保存 | mode `0600` で atomic write | password は `"***"` で追記 | `200` |
| `password:null` | 保存 | 存在すれば削除 | password は `"***"` で追記 | `200` |
| 完全 no-op | 変更なし | 変更なし | 追記なし | `200 {"message":"No changes"}` |
| `.smtp_secret` 書込失敗 | 必要なら `.smtp_config` 保存済み | 失敗 | 追記なし | `500` |
| `.config_log` 追記失敗 | 保存済み | 保存または削除済み | 失敗 | `500` |

`POST /api/smtp-test` は `.smtp_config` と `.smtp_secret` を読み、送信成功 / 失敗のどちらも `.notify_log` へ追記してから response を返す。`.notify_log` 追記失敗時は `500` を返す。SMTP password、認証失敗時の server response に含まれる credential 断片、接続 URL の userinfo は `message` と log に含めず固定文言へ置換する。

**SMTP fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| smtp-save-password | password 付き保存 | `.smtp_secret` mode `0600`、GET は `password_set:true`、log は `"***"`。 |
| smtp-delete-password | `password:null` | `.smtp_secret` 削除、password 平文なし。 |
| smtp-noop | 同一 config / password 未指定 | 状態差分なし、`.config_log` 追記なし。 |
| smtp-test-success | 設定済み test | `.notify_log` に success、response success。 |
| smtp-test-disabled | `enabled:false` | `422`、`.notify_log` 差分なし。 |
| smtp-log-failure | test 後 `.notify_log` 追記失敗 | `500`、password 平文なし。 |

**`GET /api/notify-config` への追加（`email` セクション）：**
```json
{
  "webhooks": [ { "url": "...", "label": "メイン", "enabled": true, "payload_template": null, "retry_count": 2, "retry_interval_seconds": 30, "secret": "***" } ],
  "channels": [ { "id": "n002", "type": "email", "label": "Ops", "enabled": true, "on": ["failure", "duration_anomaly"], "config": { "to": ["ops@example.com"] }, "retry_count": 2, "retry_interval_seconds": 30 } ],
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
    { "id": "q20260915100100", "trigger": "manual", "queued_at": "2026-09-15T10:01:00Z", "requested_by": "admin", "priority": "normal", "created_seq": 1, "payload": { "force": false } },
    { "id": "q20260915100200", "trigger": "webhook", "queued_at": "2026-09-15T10:02:00Z", "requested_by": "webhook", "priority": "normal", "created_seq": 2, "payload": { "delivery_id": "delivery-1", "branch": "main", "sha": "0123456789abcdef0123456789abcdef01234567" } }
  ], "max_size": 3 }
```
キューが空の場合：`{ "queued": [], "max_size": 3 }`

**`DELETE /api/queue` レスポンス：**
```json
{ "message": "Queue cleared", "cleared_count": 2 }
```
実行中のビルドは停止しない（`POST /api/build/cancel` を別途使用する）。

`POST /api/config` に `queue_max_size`（整数、`0` = キューなし）追加。

**queue 処理固定契約：**

queue 追加は `.build_state` の atomic write で行い、id は §22.0e.2 の queue id とする。queue entry は FIFO を標準とし、§27.35 の優先度キューが有効な場合だけ priority を使用する。`DELETE /api/queue` は `.build_state.queued` だけを空配列にし、実行中 build、lock、history、log を変更しない。

`queue_max_size=0` の場合、実行中に受けた queue 対応 API は `429 {"error":"queue_full"}` を返す。`.build_state.queued` の件数が `queue_max_size` 以上の場合も同じ body とする。

**queue entry schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | §22.0e.2 の queue id。 |
| `trigger` | string | 必須 | `"manual"`、`"webhook"`、`"approval"`。 |
| `queued_at` | string | 必須 | UTC ISO 8601。 |
| `requested_by` | string | 必須 | `"admin"`、`"webhook"`、`"approval"`、API token id。 |
| `priority` | string | 必須 | §27.35 の値。未指定作成時は `"normal"`。 |
| `created_seq` | integer | 必須 | 1 以上。既存最大 + 1。 |
| `payload` | object | 必須 | trigger ごとの固定 payload。未使用時は `{}`。 |

`payload` は trigger ごとに以下を許可する。未知 key は `422`、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

| trigger | payload |
|---------|---------|
| `manual` | `{ "force": boolean }`。 |
| `webhook` | `{ "delivery_id": string, "branch": string, "sha": string }`。 |
| `approval` | `{ "approval_id": string, "branch": string, "sha": string, "target": string }`。 |

**queue 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| API queue 追加 | `.build_state` lock → 最新 state 読込 → running / max_size / 重複確認 → id と created_seq 採番 → atomic write → response | write 失敗は `500`。queue 追加なし。 |
| webhook queue 追加 | event 検証 → `.webhook_events.json` 追記 → `.build_state` lock → queue append → response | event 追記前の失敗は queue なし。queue 追加失敗は event result を `error` にできる場合だけ追記し、response は `500`。 |
| approval approve queue 追加 | `.approval_queue` lock → pending 確認 → `.build_state` lock → queue append → approval status `approved` 追記 → response | queue full は `429`、approval は pending のまま。approval status 追記失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| runner 取り出し | `.build_state` lock → §27.35 の順で 1 件選択 → selected entry 削除 → `running=true` と `current_build_id` 設定 → atomic write | write 失敗は build 開始なし、lock を解放し終了コード `1`。 |
| queue clear | `.build_state` lock → `queued=[]` → atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。cleared queue は巻き戻さない。 |

重複判定は `trigger` と `payload` の正規化 JSON が一致する waiting entry を対象とする。重複時は新規 entry を追加せず `200 {"message":"Already queued","queued":true,"queue_id":"<existing>"}` を返す。`force=true` の manual entry は `force=false` と別 entry として扱う。

**queue fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| queue-add-running | running 中に manual build | queue append、created_seq 最大 + 1。 |
| queue-duplicate | 同一 manual payload を再投入 | 新規追加なし、既存 queue_id を返す。 |
| queue-full | max_size 到達 | `429 {"error":"queue_full"}`、差分なし。 |
| queue-clear | waiting 2 件で `DELETE /api/queue` | `cleared_count=2`、running/current_build_id 維持。 |
| queue-runner-take | urgent と normal が混在 | urgent を削除し running に設定、他 entry 維持。 |

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
`widgets` に未知の識別子が含まれる場合は `422` を返す。

**dashboard layout 固定契約：**

既定 widget 順は `["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]` とする。`POST /api/dashboard-layout` は `widgets` 全体置換のみ許可し、空配列、重複、未知 id は `422`。正規化後値が既存値と一致する場合は `.dashboard_layout` と `.config_log` を変更せず `{ "message":"No changes" }` を返す。

---

**`POST /api/logs/cleanup` レスポンス例：**
```json
{ "message": "Cleanup completed", "deleted_count": 12 }
```

`log_retention_days` が `0` の場合は削除せず `deleted_count: 0` を返す。

**`GET /api/history/export` レスポンス：**

`Content-Type: application/json` で返却される。

```json
{ "exported_at": "2026-09-15T10:00:00Z", "history": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "sha": "abc123", "status": "success", "trigger": "polling", "duration_seconds": 42, "flagged": false, "tags": ["release"], "comment": "" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00Z", "sha": "def456", "status": "failure", "trigger": "manual", "duration_seconds": 7, "flagged": true, "tags": [], "comment": "ネットワーク障害による失敗" }
]}
```

**`POST /api/repo-config` リクエスト / レスポンス：**
```json
// リクエスト（変更するフィールドのみ指定可）
{ "owner": "fqwink", "repo": "Adlaire-Design-System", "branch": "main", "target_file": "docs" }
// レスポンス: 200
{ "message": "Repo config updated" }
```

設定は `.repo_config` に保存し、`GET /api/repo-info` もこのファイルを参照する。

`trigger` の有効値：§13 の固定値（`"polling"`、`"force_interval"`、`"manual"`、`"webhook"`、`"retry_pending_transfer"`、`"startup_config_integrity"`、`"rollback"`、`"local_watch"`、`"approval"`）

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

## 25. 認証 実装仕様

**初期認証情報ファイル形式（JSON）：**
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

**`.admin_credentials` schema 固定：**

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

**ハッシュアルゴリズム：** Go 標準ライブラリのみで実装する `sha256_iter_v1`

| 項目 | 仕様 |
|------|------|
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| 初回 digest | `sha256(salt_bytes || password_utf8_bytes)` |
| 反復 | `iterations = 260000`。2 回目以降は `sha256(previous_digest || salt_bytes || password_utf8_bytes)` を繰り返す。 |
| 保存値 | 最終 digest を lowercase hex 文字列で `password_hash` に保存する。 |
| 比較 | 入力パスワードから同一手順で digest を生成し、`crypto/subtle.ConstantTimeCompare` で比較する。 |

`golang.org/x/crypto/pbkdf2` 等の外部パッケージは使用しない。PBKDF2、bcrypt、Argon2 等は本ファイルでは実装値を定義しない。password hash は `.admin_credentials.password_hash` に保存し、API response へは出さない。

**セッショントークン生成：**
`crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。

**セッション管理：** `components/api.go` 内のインメモリ辞書で管理。有効期限は新規発行時点の `.server_config.session_timeout_seconds` とする。設定不在時は 8 時間。再起動で全セッション破棄。単一 admin の複数同時セッションを許容する。辞書 key は token 本体ではなく `sha256(token)` の lowercase hex とし、API response、`.access_log`、`.audit_log`、サーバーログへ token 本体を出力してはならない。

**認証入力・保存禁止契約：**

| 項目 | 仕様 |
|------|------|
| 認証 header | `Authorization: Bearer {token}` だけを受け付ける。 |
| Cookie | session cookie、remember-me cookie、CSRF cookie は発行しない。受信しても認証に使わない。 |
| query token | `?token=`、`access_token`、`session` query は認証に使わず、存在しても無視する。 |
| body token | login / totp 以外の body token は認証に使わない。 |
| 永続化禁止 | session token、login ticket、setup 仮 secret、連続失敗回数はファイル保存しない。 |
| response 禁止 | password hash、salt、session token hash、ticket hash、TOTP secret 保存値は response に含めない。 |
| log 禁止 | password、current_password、new_password、token、ticket、hash、salt、TOTP code は `.access_log`、`.audit_log`、`.api_access_log`、journal に含めない。 |

**セッション期限切れ時：** `401 Unauthorized` を返す。クライアント（SDK）は `this._token` をクリアし、再ログインを促す。

**認証処理の副作用境界：**

| ケース | HTTP status | `.admin_credentials` | メモリ session / ticket | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|----------------------|-------------------------|---------------|--------------|------|
| password 不一致 | `401` | 変更なし | 変更なし | `login_failure` を追記 | `login_failure` を追記 | 連続失敗回数はメモリ上で +1。 |
| 連続失敗ロック | `429` | 変更なし | 変更なし | `login_locked` を追記 | `permission_denied` を追記 | password hash 検証は実行しない。 |
| password 成功 / TOTP 無効 | `200` | `login_count`、`last_login_at` 更新 | session 追加 | `login_success` を追記 | `login_success` を追記 | session token は全永続ログに保存しない。 |
| password 成功 / TOTP 有効 | `200` | 変更なし | ticket 追加 | `totp_required` を追記 | `totp_required` を追記 | `login_count` と `last_login_at` は更新しない。 |
| TOTP code 不一致 | `401` | 変更なし | ticket 削除 | `totp_failure` を追記 | `totp_failure` を追記 | ticket は再利用不可。 |
| TOTP 成功 | `200` | `login_count`、`last_login_at` 更新 | ticket 削除、session 追加 | `login_success` を追記 | `login_success` を追記 | session 期限は成功時点の設定で決める。 |
| session 期限切れ | `401` | 変更なし | 対象 session 削除 | 追記しない | 追記しない | `.api_access_log` は通常 API request として記録する。 |
| logout | `200` | 変更なし | 対象 session 削除 | `logout` を追記 | `logout` を追記 | token 本体と token hash は保存しない。 |
| password 変更成功 | `200` | 新 salt / hash、`must_change:false`、`updated_at` 更新 | 現 session 以外削除 | `password_change` を追記 | `password_change` を追記 | 新旧 password、hash、salt は保存しない。 |

上表で `.access_log` または `.audit_log` の追記が必要な処理は、response 返却前に追記を完了する。追記失敗時は、個別節で別指定がない限り `500 {"error":"Internal server error"}` を返す。session token、login ticket、TOTP setup secret は、必要なログ追記がすべて成功するまで response に含めてはならない。

**認証副作用順序固定契約：**

認証関連 endpoint は、以下の順序で副作用を確定する。response を返した後に session、ticket、credentials、ログを遅延更新してはならない。

| 処理 | 固定順序 | 失敗時 |
|------|----------|--------|
| `POST /api/login` password 不一致 | credentials 読込 → lock 判定 → hash 比較 → 失敗回数更新 → `.access_log` → `.audit_log` → `401` response | ログ追記失敗は `500`。失敗回数は戻さない。session/ticket は作成しない。 |
| `POST /api/login` password 成功 / TOTP 無効 | credentials 読込 → hash 比較 → 失敗回数 reset → `.admin_credentials` 更新 → session token 生成 → `.access_log` → `.audit_log` → token response | credentials またはログ追記失敗は `500`。token は response しない。 |
| `POST /api/login` password 成功 / TOTP 有効 | credentials 読込 → hash 比較 → 失敗回数 reset → ticket 生成 → `.access_log` → `.audit_log` → ticket response | ログ追記失敗は `500`。ticket は保存しない。 |
| `POST /api/login/totp` 成功 | ticket 検証 → TOTP secret 読込 → code 検証 → `.totp_secret.last_accepted_step` 更新 → `.admin_credentials` 更新 → session token 生成 → `.access_log` → `.audit_log` → token response | 途中失敗時は token を response しない。ticket は成功/失敗いずれも再利用不可にする。 |
| `POST /api/logout` | token 認証 → 対象 session 削除 → `.access_log` → `.audit_log` → response | ログ追記失敗は `500`。削除済み session は戻さない。 |
| `POST /api/change-password` | token 認証 → current password 検証 → new password 検証 → salt/hash 生成 → `.admin_credentials` 更新 → 現 session 以外削除 → `.access_log` → `.audit_log` → response | credentials 更新失敗は session を変更しない。ログ追記失敗時は `500` だが更新済み credentials は戻さない。 |
| `POST /api/sessions/revoke-all` | token 認証 → 現 session 以外を削除 → `.access_log` → `.audit_log` → response | ログ追記失敗時は `500`。削除済み session は戻さない。 |

session token と login ticket は `crypto/rand` 成功後にだけ生成し、生成した値はメモリ上で hash 化して保持する。response body に含める token / ticket は、その request の成功 response 1 回だけに含める。`403`、`429`、`500`、network 切断検出時に、未送信 token をログや状態ファイルへ退避してはならない。

**パスワード入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字 |
| 最大長 | 128 文字 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期パスワード | `--init-credentials` で `.admin_credentials` に生成する。初回ログイン時は `must_change: "prompt"` を返す。 |
| 変更時検証 | `new_password` が現在パスワードと同一の場合は `422` を返す。 |
| 失敗時応答 | パスワード不一致は `401` と `{"error":"Unauthorized"}` を返し、どの条件に失敗したかは返さない。 |
| 成功時保存 | `.admin_credentials` に新 salt、新 hash、`must_change:false`、`updated_at` を原子的に保存する。 |

**セッションレコード形式（メモリ上）：**

```json
{
  "token_hash": "<sha256_hex>",
  "session_id": "<16_byte_hex>",
  "created_at": "2026-09-15T10:00:00Z",
  "expires_at": "2026-09-15T18:00:00Z",
  "last_used_at": "2026-09-15T10:05:00Z"
}
```

`session_id` は `crypto/rand` で 16 bytes を生成し、lowercase hex とする。`GET /api/sessions` は `session_id` ではなく `current`、`created_at`、`expires_at`、`last_used_at` のみ返す。認証必須 API で有効 token を受信した場合、`last_used_at` を現在時刻へ更新する。期限切れ token は検出時にメモリから削除する。`POST /api/logout` は対象 token のみ削除する。`POST /api/sessions/revoke-all` は現在 token 以外を削除する。

**ログイン失敗制御：**

| 項目 | 仕様 |
|------|------|
| 失敗記録 | `components/api.go` はメモリ上で直近の連続ログイン失敗回数と最終失敗時刻を保持する。再起動で失敗回数はリセットされる。 |
| ロック条件 | 連続 10 回失敗した場合、最終失敗から 10 分間 `POST /api/login` を `429 Too Many Requests` と `{"error":"Too many attempts"}` で拒否する。 |
| 成功時 | ログイン成功時は連続失敗回数を 0 に戻す。 |
| 応答時間 | パスワード不一致、存在しない credentials、ロック中を除く検証失敗では、条件の詳細をレスポンスへ出さない。 |
| ログ | 成功、失敗、ロック拒否はいずれも `.access_log` へ追記する。password、token、hash、salt は記録しない。 |

`.access_log` または `.audit_log` 追記失敗時は、ログイン失敗では `500` を返し、失敗回数は増加済みのままとする。ログイン成功時は session token 発行前に `.admin_credentials`、`.access_log`、`.audit_log` を更新し、いずれかに失敗した場合は session token を発行しない。TOTP 有効時は ticket 発行前に `.access_log` と `.audit_log` を追記し、追記失敗時は ticket を発行しない。

**ログインフロー：**
```
POST /api/login
  ├─ 連続失敗ロック中 → 429
  └─ パスワードハッシュ検証
       ├─ 失敗 → 連続失敗回数 + 1 → .access_log 追記 → .audit_log 追記 → 401
       └─ 成功 → 連続失敗回数を 0 へリセット
                  ├─ must_change == true → must_change: "prompt"
                  ├─ TOTP 有効 → ticket 生成・返却
                  └─ TOTP 無効 → login_count / last_login_at 更新 → セッショントークン生成・返却
```

**パスワード変更時：** 新しい salt を生成しハッシュを更新し、`must_change:false` を保存する。変更完了後に現セッション以外のセッションを破棄。

**`--init-credentials` オプション：** `components/api.go` を `--init-credentials` 引数で起動した場合、初期パスワード `admin` で `.admin_credentials` を生成して終了する（HTTP サーバーは起動しない）。

`.admin_credentials` が既に存在する場合、`--init-credentials` は上書きせず `409` 相当の終了コード `2` で終了し、標準エラーへ `credentials already exist` を出力する。初期化成功時の終了コードは `0` とする。

**`--init-credentials` CLI 固定契約：**

| ケース | stdout | stderr | 終了コード | 副作用 |
|--------|--------|--------|------------|--------|
| 新規生成成功 | `credentials initialized` + LF | 空 | `0` | `.admin_credentials` を mode `0600` で作成する。 |
| 既存あり | 空 | `credentials already exist` + LF | `2` | 既存ファイルを変更しない。 |
| `--state-dir` 相対 path | 空 | `state directory must be absolute: {path}` + LF | `2` | ファイル作成なし。 |
| 書込失敗 | 空 | `credentials write failed` + LF | `1` | tmp を削除し、部分ファイルを残さない。 |
| rand 失敗 | 空 | `random source failed` + LF | `1` | ファイル作成なし。 |

生成手順は、state dir 検証 → 既存確認 → salt 生成 → hash 生成 → `{path}.tmp.{pid}` へ JSON + LF 書込 → mode `0600` → file sync → rename → parent directory sync の順に固定する。rename 後の sync に失敗した場合は `1` を返し、作成済みファイルは残る。実装者判断で初期パスワードを環境変数、対話入力、ランダム生成へ変更してはならない。

**認証 fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| auth-password-failure | 誤 password で `POST /api/login` | `401`、session/ticket なし、失敗回数 +1、`.access_log` と `.audit_log` に secret なし。 |
| auth-login-lock | 連続 10 回失敗後の `POST /api/login` | `429`、password hash 検証なし、`.access_log` に `login_locked`、`.audit_log` に `permission_denied`。 |
| auth-session-issued | TOTP 無効で password 成功 | token は response のみ、`.admin_credentials.login_count` +1、ログに token/hash/salt なし。 |
| auth-session-expired | 期限切れ session で保護 API | `401`、対象 session 削除、`.access_log` と `.audit_log` は追記しない。 |
| auth-password-change | password 変更成功 | 新 salt/hash、現 session 以外削除、`password_change` ログ、password/hash/salt 平文なし。 |
| auth-log-write-failure | login 成功時に `.audit_log` 追記失敗 | `500`、session token を response しない。 |

---

---

## 27. API owner 追加仕様化機能 詳細仕様

### 27.5 設定バリデーション API
owner component は `api` とする。collaborator component は `sdk`、`ui`、`statefile` とする。


`POST /api/config/validate` は、`POST /api/config` と同じ入力を受け取り、保存せずに検証結果を返す。

Request body は partial `ConfigObject` とする。未知 key を含む場合は `422 {"error":"Unknown config key","key":"..."}` を返す。型不一致、範囲外、URL 不正、commit status context 空文字、retry 設定不正は `200` で `valid=false` として返す。JSON 不正は `422` とする。

成功時 response:

```json
{
  "valid": false,
  "config": {"log_max_lines": 500},
  "errors": [
    {"field":"build_retry_max","code":"out_of_range","message":"build_retry_max must be 0..10"}
  ],
  "warnings": [
    {"field":"commit_status_target_url","code":"http_url","message":"https is recommended"}
  ]
}
```

`config` は既存 `.server_config` に request body を merge した正規化後の値を返す。ただし保存してはならない。`.config_log`、`.api_access_log` 以外の状態ファイルを更新してはならない。`.api_access_log` は通常 API 呼び出しとして記録する。

**validate 副作用固定契約：**

| 項目 | 仕様 |
|------|------|
| merge | request body は既存 config に shallow merge し、object 値は対象 object 単位で置換する。 |
| secret key | secret 風 key は未知 key と同じく `422`。値は response、access log、server log に出さない。 |
| `valid=false` | HTTP status は `200`。状態ファイルは更新しない。 |
| error 順 | request body の key 出現順で `errors[]` を返す。 |
| warning 順 | field 名昇順で返す。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 正常値 | `valid=true`、errors 空、状態差分なし。 |
| 範囲外 | `valid=false`、該当 field error。 |
| 未知 key | HTTP 422、状態差分なし。 |
| secret 風 key | HTTP 422、response と log に値を含めない。 |
| 複数 error | request key 出現順で返る。 |
| valid false | HTTP 200、保存差分なし。 |

### 27.6 API アクセスログ
owner component は `api` とする。collaborator component は `sdk`、`ui`、`statefile` とする。


`components/api.go` は全 `/api/` request について `.api_access_log` へ JSON Lines を追記する。`GET /api/health` も対象とする。静的 file 配信、admin HTML、SDK JS は対象外とする。

追記タイミングは response status 確定後とする。追記失敗時は、対象 API の本来の response を優先し、サーバーログに `API_ACCESS_LOG_WRITE_FAILED` を出す。access log 書き込み失敗を理由に API response を `500` へ変更してはならない。

`GET /api/api-access-log` は `limit`、`offset`、`method`、`path`、`status` query を受け付ける。`limit` は 1〜1000、既定値 100。`offset` は 0 以上、既定値 0。`method` は大文字 HTTP method 完全一致。`path` は prefix match。`status` は HTTP status 完全一致。壊れた行は無視し、新しい順で返す。

**access log record 固定契約：**

| key | 型 | 仕様 |
|-----|----|------|
| `timestamp` | string | response status 確定時刻。 |
| `request_id` | string | 16 byte hex。response header `X-Request-Id` と一致させる。 |
| `method` | string | 大文字 HTTP method。 |
| `path` | string | query を除いた path。 |
| `status` | integer | 実際に返した HTTP status。 |
| `actor` | string/null | admin、token id、webhook、または `null`。 |
| `auth_type` | string | `"none"`、`"session"`、`"api_token"`、`"webhook"`。 |

request body、query 全体、cookie、Authorization header、token、password、secret は保存しない。access log 追記失敗時も、既に確定した response status を変更しない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 認証成功 API | actor `admin` または token id で記録される。 |
| 認証失敗 API | status 401、actor null、secret なし。 |
| Webhook | auth_type `webhook`、actor `webhook`。 |
| filter | method/path/status が完全に効く。 |
| 書込失敗 | 本来 response を維持し、server log に WARN。 |
| request id | response header と log の `request_id` が一致する。 |
| secret query | query 全体を保存しない。 |

### 27.11 ポーリング間隔の動的変更
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は schedule 設定と systemd timer 更新の責務を持つ。`runner` は本機能で systemd timer を変更しない。

本機能の目的は、管理 API から systemd timer の実行間隔を変更し、次回以降の runner 起動間隔を固定仕様どおり反映することである。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/schedule/interval` |
| Request | `{ "interval_seconds": integer }` |
| 許容値 | 30〜86400 秒 |
| Response | `{ "message": "Schedule interval updated", "interval_seconds": N }` |
| 状態 | `.server_config.schedule_interval_seconds` を更新し、`.config_log` へ記録する。 |
| systemd | `/etc/systemd/system/adlaire-ci.timer` の `OnUnitActiveSec` を `N seconds` 相当へ更新する。 |

**処理順序：**

1. 認証、maintenance、入力型、範囲を検証する。
2. `.server_config` を atomic write で更新する。
3. `.config_log` に `schedule_interval_seconds` の diff を追記する。
4. systemd timer ファイルを書き換える。
5. `systemctl daemon-reload` を実行する。
6. `systemctl restart adlaire-ci.timer` を実行する。
7. `systemctl show adlaire-ci.timer -p OnUnitActiveSec` 相当で反映を確認する。
8. 成功 response を返す。

**systemd 更新固定契約：**

| 項目 | 仕様 |
|------|------|
| timer 書換 | 既存 unit を直接編集せず、管理対象 drop-in `/etc/systemd/system/adlaire-ci.timer.d/override.conf` を atomic write する。 |
| drop-in 内容 | `[Timer]`、`OnUnitActiveSec={N}s`、`Persistent=true` のみを書き込む。 |
| 反映確認 | `systemctl show` の値を秒へ正規化して request 値と一致確認する。 |
| rollback | systemd 失敗時に `.server_config` は巻き戻さない。失敗を `.config_log` に追加する。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 入力が integer でない、または範囲外 | `422`。状態ファイル、systemd は変更しない。 |
| `.server_config` 書き込み失敗 | `500`。systemd は変更しない。 |
| systemd timer 書き換え失敗 | `500`。`.server_config` は更新済みのまま残し、`.config_log` に `systemd_update_failed` を記録する。 |
| daemon-reload / restart / show 失敗 | `500`。server log に `SCHEDULE_INTERVAL_APPLY_FAILED` を出す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常値 300 | `.server_config`、timer、API response が 300 で一致する。 |
| 最小値 30 | 成功する。 |
| 範囲外 29 | `422`、差分なし。 |
| systemd 失敗 | `500`、`.config_log` に失敗記録。 |
| show 不一致 | `500`、server config は更新済み。 |

### 27.12 GitHub Webhook 受信
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は署名検証、イベント記録、queue 投入を担当し、`runner` は queue entry を処理する。

本機能の目的は、GitHub push event を HMAC-SHA256 署名検証したうえで受信し、定期 polling を待たずに build queue へ投入することである。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/webhook` |
| 必須 header | `X-GitHub-Event`, `X-GitHub-Delivery`, `X-Hub-Signature-256` |
| 対象 event | `push` のみ |
| Secret | `.webhook_secret` |
| 成功 response | `{ "message": "Webhook accepted", "queued": true, "event_id": "..." }` |
| 状態 | `.webhook_events.json` へ追記し、必要に応じて `.build_state.queued` へ `trigger="webhook"` entry を追加する。 |

**署名検証：**

署名は `sha256=` prefix を含む lowercase hex とする。検証は raw request body に対して `HMAC-SHA256(secret, body)` を計算し、定数時間比較で行う。secret 不在、header 不在、prefix 不正、hex 不正、署名不一致はすべて `401` とし、queue へ投入しない。

**処理順序：**

1. method と body size を検証する。
2. `.webhook_secret` を読み込む。
3. HMAC 署名を検証する。
4. JSON body を parse する。
5. event が `push` であることを確認する。
6. `ref`、`after`、`repository.owner.login`、`repository.name` を抽出する。
7. `.branch_config` と照合し、対象 branch がある場合だけ queue へ追加する。
8. `.webhook_events.json` にイベント結果を JSON Lines で追記する。
9. response を返す。

**Webhook 副作用固定契約：**

| ケース | `.webhook_events.json` | `.build_state.queued` | response |
|--------|------------------------|-----------------------|----------|
| 署名不正 / secret 不在 | 変更なし | 変更なし | `401 {"error":"Unauthorized"}` |
| JSON 不正 | 変更なし | 変更なし | `422 {"error":"Validation failed",...}` |
| event が `push` 以外 | `ignored_event` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| 対象 branch なし | `ignored_branch` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| queue 追加成功 | `queued` を追記 | queue entry を追加 | `202 {"message":"Webhook accepted","queued":true,"event_id":...,"queue_id":...}` |
| queue full | `queue_full` を追記 | 変更なし | `429 {"error":"queue_full"}` |
| event log 追記失敗 / queue 追加前 | 変更なし | 変更なし | `500 {"error":"Internal server error"}` |
| event log 追記失敗 / queue 追加後 | 変更なし | queue entry は残す | `202 {"message":"Webhook accepted","queued":true,"event_log_failed":true,"queue_id":...}` |

queue entry は §16C の queue entry schema を使用し、`trigger:"webhook"`、`requested_by:"webhook"`、`priority:"normal"`、`payload.delivery_id`、`payload.branch`、`payload.sha` を保存する。`X-GitHub-Delivery` が既に pending queue に存在し、同一 branch / sha の場合は重複投入せず、event log に `result:"duplicate"`、既存 `queued_id` を記録し、`202 {"message":"Webhook already queued","queued":true,"queue_id":"<existing>"}` を返す。

**Webhook validation 固定契約：**

| 項目 | 仕様 |
|------|------|
| body size | 最大 1 MiB。超過時は署名検証後に `413`。 |
| branch 抽出 | `ref` が `refs/heads/{branch}` 形式でない場合は `ignored_branch`。 |
| sha | `after` が 40 文字 lowercase hex でない場合は `422`。 |
| repository | `owner.login` と `name` から `owner/name` を作る。欠落時は `422`。 |
| event id | `wh{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| secret | signature、secret、raw payload は event log、access log、audit log に保存しない。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 署名不正 | `401`、イベントログ追記なし、queue なし。 |
| event が `push` 以外 | `202`、`queued=false`、イベントログには `ignored_event` として記録する。 |
| JSON 不正 | `422`、queue なし。 |
| 対象 branch なし | `202`、`queued=false`、イベントログに記録する。 |
| queue 上限 | `429`、イベントログに `queue_full` を記録する。 |
| イベントログ書き込み失敗 | `500`、queue 追加前なら queue しない。queue 追加後なら response に `queued=true` と `event_log_failed=true` を含める。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 push | 署名検証成功、イベントログ追記、queue entry `trigger="webhook"`。 |
| 署名不一致 | `401`、状態差分なし。 |
| 対象外 branch | `202 queued=false`、イベントログのみ。 |
| queue full | `429`、queue 差分なし。 |
| duplicate | queue 追加なし、既存 queue id を返す。 |
| payload oversized | `413`、状態差分なし。 |

### 27.13 Webhook イベントログ / 一覧取得 API
owner component は `api` とする。collaborator component は `sdk`、`ui`、`statefile` とする。

本機能の目的は、受信した GitHub Webhook の監査情報を `.webhook_events.json` に保存し、管理 API、SDK、UI からページング参照できるようにすることである。

**保存 schema：**

`.webhook_events.json` は JSON Lines とし、1 行 1 event を追記する。mode は `600` とする。

| key | 型 | 必須 | 説明 |
|-----|----|------|------|
| `timestamp` | string | 必須 | ISO 8601 UTC。 |
| `delivery_id` | string/null | 必須 | `X-GitHub-Delivery`。 |
| `event` | string | 必須 | GitHub event 名。 |
| `ref` | string/null | 必須 | push ref。 |
| `branch` | string/null | 必須 | `refs/heads/` を除いた branch。 |
| `sha` | string/null | 必須 | push `after`。 |
| `repository` | string/null | 必須 | `owner/repo`。 |
| `build_triggered` | boolean | 必須 | queue 追加済みなら `true`。 |
| `queued_id` | string/null | 必須 | queue id または `null`。 |
| `result` | string | 必須 | `"queued"`, `"duplicate"`, `"ignored_event"`, `"ignored_branch"`, `"queue_full"`, `"error"`。 |

`delivery_id` は 1〜200 文字、`event` は 1〜100 文字、`repository` は `owner/repo` 形式、`sha` は `null` または 40 文字 lowercase hex とする。保存時に request header 全体、署名値、secret、payload 全体を保存してはならない。

**一覧 API：**

`GET /api/webhook-events` は `limit` と `offset` query を受け付ける。`limit` は 1〜1000、既定値 50。`offset` は 0 以上、既定値 0。新しい順で返す。壊れた行は無視し、server log に `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` を出す。

Response は `{ "events": WebhookEventRecord[], "total": N }` とする。SDK `getWebhookEvents(limit,offset)` は `limit` と `offset` を常に query へ送信する。UI は件数、delivery id、event、branch、sha、result、queued id を表示する。

**webhook events 取得固定契約：**

| 項目 | 仕様 |
|------|------|
| 並び順 | `timestamp` 降順、同時刻は file 出現順の逆順。 |
| total | 壊れた行を除外した総件数。 |
| offset | filter 後、並び替え後に適用する。 |
| 壊れた行 | 内容を response、server log に含めない。固定コードだけ出す。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| event 追記 | JSON Lines へ schema 通り保存される。 |
| 一覧取得 | 新しい順、limit/offset が効く。 |
| 壊れた行 | API は継続し、壊れた行を返さない。 |
| duplicate delivery | queue 重複なし、event log は `duplicate`。 |
| total | 壊れた行を除外した件数。 |

### 27.16 ヘルスチェックエンドポイント
owner component は `api` とする。collaborator component は `statefile` とする。

本機能の目的は、認証不要の `GET /api/health` で、外部監視が Adlaire CI の最低限の稼働状態を確認できるようにすることである。

**Response：**

```json
{
  "status": "ok",
  "last_build_at": "2026-09-16T00:00:00Z",
  "last_build_status": "success",
  "last_deploy_at": "2026-09-16T00:01:00Z",
  "last_deploy_status": "success",
  "pending_transfers": 0,
  "uptime_seconds": 3600,
  "checks": []
}
```

`status` は `"ok"`、`"degraded"`、`"error"` のいずれかとする。必須状態ファイル破損がある場合は `degraded`、API process が応答できるが重大な read error がある場合は `error` とする。HTTP status は、API 自体が response を生成できる限り `200` とし、JSON 生成不能などの場合だけ `500` とする。

**読み取り元：**

`.build_status.json` を第一参照元とし、不在時は `.build_history` と `.pending_transfers` から算出する。`.build_status.json` 破損時は自動修復せず、`checks[]` に `build_status_corrupt` を含める。

**health checks 固定契約：**

| check | 条件 |
|-------|------|
| `build_status_missing` | `.build_status.json` 不在で fallback 算出した。 |
| `build_status_corrupt` | `.build_status.json` parse 失敗。 |
| `pending_transfers_read_error` | `.pending_transfers` 読取失敗。 |
| `notify_pending_read_error` | `.notify_pending` 読取失敗。 |
| `runner_stale` | `running=true` かつ `last_started_at` が 24 時間より古い。 |

`checks[]` は上表の順で返す。`status` は checks が空なら `"ok"`、read error または stale があれば `"degraded"`、response 生成不能だけ `"error"` とする。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 | HTTP 200、`status="ok"`。 |
| pending transfer あり | `pending_transfers` に件数、`status="degraded"`。 |
| build status 破損 | HTTP 200、`status="degraded"`、checks に記録。 |
| read error | HTTP 200 または 500 の条件が仕様通り。 |
| stale runner | checks に `runner_stale`。 |

### 27.17 ビルドログ重大度フィルター
owner component は `api` とする。collaborator component は `sdk`、`ui`、`archive`、`statefile` とする。

本機能の目的は、`GET /api/logs/search` と UI ログビューアで重大度別に build log を絞り込めるようにすることである。

**入力：**

`GET /api/logs/search` は既存 query に加えて `level` を受け付ける。`level` の許容値は `"info"`、`"warn"`、`"warning"`、`"error"`、`"debug"` とし、大文字小文字は区別しない。正規化後は `"INFO"`、`"WARNING"`、`"ERROR"`、`"DEBUG"` とする。`warn` は `"WARNING"` と同義とする。不正値は `422`。

**検索対象：**

`.build_logs/{id}.json.stdout`、`stderr`、`warnings`、`error`、archive log を対象とする。行頭が `[WARN]` または `[WARNING]` の行は WARNING、`[ERROR]` または stderr の非空行は ERROR、`[DEBUG]` は DEBUG、それ以外は INFO と分類する。

**ログ行分類固定契約：**

| 入力 | level |
|------|-------|
| stderr の非空行 | `ERROR` |
| stdout / warnings の `[ERROR]` prefix | `ERROR` |
| stdout / warnings の `[WARN]` または `[WARNING]` prefix | `WARNING` |
| stdout / warnings の `[DEBUG]` prefix | `DEBUG` |
| 上記以外 | `INFO` |

検索結果は `{build_id, level, source, line_number, message}` とし、`source` は `"stdout"`、`"stderr"`、`"warnings"`、`"error"` のいずれかとする。`line_number` は 1 始まり、配列項目は配列内 index + 1 とする。

**SDK / UI：**

SDK `searchLogs(q,from,to,level)` は `level` 指定時だけ query に送信する。UI は INFO / WARNING / ERROR / DEBUG の filter control を提供し、選択なしでは全件を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| `level=warn` | WARNING 行だけ返る。 |
| `level=error` | ERROR 行だけ返る。 |
| 不正 level | `422`。 |
| archive log | 通常 log と同じ分類で検索される。 |
| stderr | prefix なしでも ERROR。 |
| line number | 1 始まりで返る。 |

### 27.18 ブランチ設定の動的変更 API
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。

本機能の目的は、監視対象 branch / target / deploy target を `.branch_config` で管理し、API 経由で変更できるようにすることである。

**API：**

`GET /api/branch-config` は `.branch_config` が存在する場合 `{"source":"file","branches":[...]}`、不在の場合 `{"source":"default","branches":[...]}` を返す。`POST /api/branch-config` は `{ "branches": BranchTargetRecord[] }` を受け取る。

**保存仕様：**

永続ファイルの key は必ず `branch_targets` とする。API request / response で `branches` を使う場合も保存前に `branch_targets` へ変換する。空配列を受け取った場合は `.branch_config` を削除し、default 復帰とする。

**branch config 正規化固定契約：**

| 項目 | 仕様 |
|------|------|
| branch 重複 | 同一 `branch` の重複は `422`。 |
| deploy target id | branch 内で一意。重複は `422`。 |
| target_files | 存在する場合は §27.21 の正規化を適用する。 |
| 保存順 | branch 名昇順、同一 branch 内 deploy target id 昇順で保存する。 |
| 削除 | POST empty で `.branch_config` を削除した後、`.config_log` に default 復帰を記録する。 |

**検証：**

各 branch target は `branch`、`target_file`、`sha_file`、`src`、`out`、`deploy_targets` を必須とする。`target_file` は相対パスで `..` 禁止、`sha_file` / `src` / `out` / `dest_dir` は絶対パス、deploy target は最大 20 件、branch target は最大 50 件とする。

**runner 取り込み：**

`components/runner.go` は起動ごとに `.branch_config` を読む。API 更新後、runner 再起動は不要だが、既に実行中の runner へは反映しない。次回起動から反映する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| GET default | `.branch_config` 不在で `source="default"`。 |
| POST valid | `.branch_config.branch_targets` として保存、config log 追記。 |
| POST empty | `.branch_config` 削除、default 復帰。 |
| 相対 `sha_file` | `422`、状態差分なし。 |
| branch 重複 | `422`、状態差分なし。 |
| POST empty log failure | `.branch_config` は削除済み、response は `500`。 |

### 27.20 設定変更の詳細 diff 記録
owner component は `api` とする。collaborator component は `statefile` とする。

本機能の目的は、設定変更 API が何を変更したかを `.config_log` に機械可読 diff と人間可読 diff の両方で残すことである。

**対象 API：**

`.config_log` を write する全 API を対象とする。少なくとも `POST /api/config`、`POST /api/log-level`、`POST /api/notify-config`、`POST /api/repo-config`、`POST /api/branch-config`、`POST /api/webhook-config`、`POST /api/pat-update`、schedule 系 API、maintenance、access-control、hooks、alert-rules、tag-rules、pipeline-config、notes、smtp-config、dashboard-layout、snapshot delete を含む。

**ログ schema：**

各行は §22.0c `.config_log` schema に従う。`diff` は `{key:[before,after]}`、`diff_text` は 1 行以上の文字列とする。差分がない場合、対象 API は状態ファイルを書かず、`.config_log` も追記せず、response は `{ "message": "No changes" }` とする。

**マスク条件：**

キー名に `password`、`token`、`secret`、`pat`、`smtp_password` を含む値は before / after とも `"***"` に置換する。配列や object の内部 key も同じ規則で再帰的にマスクする。

**diff_text 形式：**

`{key}: {before} -> {after}` を key 名昇順で 1 行ずつ連結する。値は JSON 表現とし、secret は `"***"` とする。複数行値は `\n` escape した 1 行 JSON string とする。

**diff 生成固定契約：**

| 項目 | 仕様 |
|------|------|
| 比較対象 | 正規化後、既定値 merge 後の before / after object。保存対象外 key、response 専用 key、password 平文は比較対象に含めない。 |
| object | dot path で再帰比較する。例: `summary.hour`。 |
| array | index 比較ではなく配列全体を JSON 値として比較する。並び順が仕様上正規化される配列は正規化後に比較する。 |
| secret key | key path のいずれかに `password`、`token`、`secret`、`pat`、`smtp_password` を含む場合、before / after を `"***"` にする。 |
| no-op | diff が空の場合は状態ファイル、secret file、`.config_log` を変更しない。 |
| 追記順 | 対象状態ファイル保存後に `.config_log` を追記する。`.config_log` 失敗時は `500`、保存済み状態は巻き戻さない。 |

`.config_log` record の `target` は endpoint 固定名、`actor` は管理 session なら `"admin"`、API token なら token id とする。`request body` 全体、HTTP header、cookie、secret 平文を保存してはならない。

**diff 対象 endpoint 固定名：**

| endpoint | target |
|----------|--------|
| `POST /api/config` | `server_config` |
| `POST /api/notify-config` | `notify_config` |
| `POST /api/repo-config` | `repo_config` |
| `POST /api/branch-config` | `branch_config` |
| `POST /api/webhook-config` | `webhook_config` |
| `POST /api/smtp-config` | `smtp_config` |
| `POST /api/pipeline-config` | `pipeline_config` |
| その他 `.config_log` 対象 | method と path から `/api/` prefix を除き、`/` と `-` を `_` に置換した固定名。 |

diff 生成は状態保存前に memory 上で完了させる。diff 生成に失敗した場合は状態ファイルを書かない。`.config_log` 追記に失敗した場合は保存済み状態を巻き戻さず、response は `500` とする。

**異常系：**

| 条件 | 処理 |
|------|------|
| config log 追記失敗 | 対象状態ファイルの更新を失敗扱いにし、`500` を返す。 |
| diff 生成失敗 | 状態ファイルを書かず `500`。 |
| secret マスク漏れ検出 | 実装不合格。該当 API は完了扱いにしない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 通常変更 | `.config_log` に diff と diff_text が残る。 |
| secret 変更 | 値は `"***"` だけ保存される。 |
| 変更なし | 状態ファイルも config log も更新しない。 |
| 複数 key | key 昇順で diff_text を生成する。 |
| array 正規化 | 正規化後に同一なら no-op。 |
| config log 失敗 | 対象状態は保存済み、response は `500`。 |
| diff failure | 状態差分なしで `500`。 |
| target name | endpoint から固定 target 名が生成される。 |

**§27.21〜§27.38 / §27.42〜§27.47 機能別実装完全性固定契約：**

§27.21〜§27.38、§27.42〜§27.47 の各機能は、個別節の本文に加えて下表を満たした場合だけ実装完了とする。§27.38a は §27.21〜§27.38 の runner 拡張を横断検証する補足契約として扱う。

| 節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | 必須 fixture |
|----|------|------|------|---------------------------|--------------|--------------|
| §27.21 | 複数ファイル監視 | `.branch_config.branch_targets[].target_files`、GitHub content SHA または local SHA。 | `changed_targets[]`、target 単位 SHA cache、build log。 | 成功時だけ該当 target SHA cache を更新する。 | SHA 部分失敗では build を開始せず、成功取得済み cache も更新しない。 | 1 件変更、複数変更、変更なし、不正 path、force build、SHA 部分失敗。 |
| §27.22 | pipeline YAML | `.pipeline.yml` または `.pipeline_config.inline_yaml`。 | `pipeline_steps[]`、step stdout/stderr、build status。 | step を定義順に実行し、build 終了時に定義順で保存する。 | parse / command 不正では build を開始しない。required step 失敗で後続 required step を実行しない。 | 成功、required 失敗、optional 失敗、禁止 YAML、secret mask、timeout。 |
| §27.23 | local watch | `.server_config.watch_mode`、local `src` 配下 Markdown。 | `.local_watch_state.json`、trigger `local_watch`。 | local mode では GitHub API / PAT を呼ばず、成功時だけ state を置換する。 | file read 失敗は build なし。dry-run は state を作成 / 更新しない。 | 初回、変更なし、1 file 変更、token 不在、file 削除、out 除外。 |
| §27.24 | tag filter | `.server_config.tag_filter`、GitHub tags refs。 | `matched_tags[]`、skip status。 | tag 一致時だけ build。tag 不一致 skip は `.build_status.json` だけ更新する。 | tags API 最終失敗は build なし。local mode 併用は終了コード `2`。 | tag 一致、不一致、patterns 空、API 失敗、prefix pattern、local mode 併用。 |
| §27.25 | build cache | `--cache-dir`、cache 設定、input / deps SHA。 | cache entry、page cache、`[REPORT]` cache counts。 | hit 時は変換結果を再利用し、miss 成功時だけ entry を書く。 | cache read/write 失敗は build を成功可能にし、WARN / report に残す。 | 2 回目 hit、1 file miss、theme 変更、cache 破損、entry 不一致、write failure。 |
| §27.26 | parallel deploy | `.server_config.deploy_parallelism`、deploy targets。 | `target_results[]`、`.pending_transfers`。 | target ごとに並列転送し、result は設定順で保存する。 | worker 内部失敗は該当 target failure。他 target は継続する。 | 3 target/2 worker、1 target 失敗、parallelism 1、result order、pending duplicate。 |
| §27.27 | build hooks | `.hooks`、pre/post hook command_args。 | hook log、build log hook result。 | pre は build 前、post は build 後に id 昇順で実行する。 | pre abort で build 本体を開始しない。post 失敗は build status を変更しない。 | pre success、pre abort、post failure、timeout、delete hook、secret mask。 |
| §27.28 | dependency tracking | builder dependency manifest、Markdown link / asset reference。 | `.dependency_manifest.json`、affected target 判定。 | build 成功時だけ manifest を更新する。 | manifest 破損は full build 扱い。dry-run は更新しない。 | direct dep、shared asset、deleted dep、manifest corrupt、dry-run no write。 |
| §27.29 | remote build | remote build 設定、ssh target、artifact path。 | remote artifact metadata、build log remote section。 | remote command と artifact fetch を行い、検証成功時だけ deploy / history へ進む。 | remote timeout / checksum mismatch は failure とし、secret / command credential を保存しない。 | remote success、timeout、checksum mismatch、artifact missing、secret mask。 |
| §27.30 | approval flow | approval policy、pending request、approve/reject API。 | `.approval_queue`、history `approval_*`、通知。 | pending 作成後、承認時だけ queue / build へ進む。 | expired / rejected は build を開始しない。history append 失敗時も approval record は残す。 | pending、approve、reject、expired、duplicate approve、notify failure。 |
| §27.31 | branch env | branch env 設定、pipeline / hook env。 | merged env、masked log。 | runner env → branch env → step/hook env の順で上書きする。 | env key 不正は保存不可 / runner 設定エラー。secret はログに出さない。 | merge order、invalid key、secret stdout、branch missing、override。 |
| §27.32 | build notification | notify config、build event、channel 設定。 | notify payload、`.notify_log`、`.notify_pending`。 | event ごとに payload を生成し、送信結果を記録する。 | retry 対象失敗は pending。secret は payload/log に含めない。 | success notify、failure notify、webhook 5xx、command timeout、disabled、secret mask。 |
| §27.33 | build trends | build duration samples、history/log。 | `.build_trends.json`、trend stats response。 | build 完了時に sample を追加し、上限件数で trim する。 | trend 保存失敗は build 成否を反転しない。破損時は再集計契約に従う。 | sample append、trim、median/p95、corrupt rebuild、save failure。 |
| §27.34 | build chain | chain config、upstream build result。 | chain execution log、queued child build。 | chain 条件一致時だけ次 build を queue する。 | chain config 破損は chain 無効として通常 build は継続する。 | chain success、condition mismatch、loop detect、queue full、config corrupt。 |
| §27.35 | priority queue | queue request、priority、created_seq。 | sorted queue、queue API response。 | queue 保存時に priority / created_seq を固定する。 | queue full は `429`、既存 queue を変更しない。 | high priority、same priority FIFO、queue full、clear queue、invalid priority。 |
| §27.36 | failure classification | build failure evidence、logs/status。 | `failure_category`、`failure_evidence[]`、API filter。 | finalizer で分類し、log/history/status に保存する。 | 分類不能でも build 結果は保持し、category `unknown` とする。 | pipeline timeout、api failure、deploy failure、unknown、filter。 |
| §27.37 | execution environment | runner host/runtime info、builder version。 | environment record、build log/status。 | build 開始時に取得し、log に保存する。 | 取得失敗は `unknown` を保存し build は継続する。secret / env 全量は保存しない。 | normal、hostname unavailable、builder version timeout、secret env absent。 |
| §27.38 | duration anomaly | duration samples、anomaly config。 | anomaly flag、alert / notify、stats。 | build 完了時に閾値判定し、該当時だけ alert/notify を作る。 | samples 不足では判定しない。通知失敗は build 成否を反転しない。 | normal、avg anomaly、p95 anomaly、insufficient samples、notify failure。 |
| §27.42 | trigger API scope | API token scopes、endpoint group。 | allow/deny decision、audit。 | scope 一致時だけ endpoint 実行。 | scope 不足は `403`、endpoint 固有処理なし。 | read allowed、trigger denied、admin allowed、health no scope、audit deny。 |
| §27.43 | API key management | token label/scopes/expires_at。 | `.api_tokens` record、token 本体 1 回 response。 | token 作成 / 失効 / 認証成功時に状態と audit を更新する。 | audit 失敗時は token 本体を返さない。token 破損 state は自動再生成しない。 | create、revoke、expired auth、scope auth、audit failure、secret non-persistence。 |
| §27.44 | audit log | 認証、権限、設定、token、build trigger event。 | `.audit_log` JSON Lines、audit API response。 | 監査対象操作の成否確定後に追記する。 | 追記失敗は対象操作を `500` 扱い。ただし保存済み状態は個別契約どおり戻さない。 | password change、permission denied、secret input、corrupt line、filter、append failure。 |
| §27.45 | session timeout | `.server_config.session_timeout_seconds`。 | 新規 session `expires_at`。 | session 発行直前の設定で期限を計算する。 | 範囲外は `422`。既存 session の期限は変更しない。 | 300 秒、既存 session 維持、範囲外、key 不在、TOTP login、秒精度。 |
| §27.46 | TOTP | setup secret、ticket、TOTP code。 | `.totp_secret`、ticket/session response、audit/access log。 | secret / ticket はメモリと 1 回 response に限定し、成功時だけ永続状態を更新する。 | code 不一致 / replay は token を返さない。audit 失敗時も secret 平文を出さない。 | setup、login、replay、secret one-time、ticket reuse、disable、window、全角 code。 |
| §27.47 | API rate limit | rate policy、remote addr、actor key。 | `.api_rate_state`、`429`、state summary。 | key 群を同一 lock で判定 / 更新する。 | 上限超過では count を増やさず endpoint 固有処理を行わない。audit 失敗時は `500`。 | under limit、over IP、over token、disabled、policy update、audit failure、window reset。 |

**§27.21〜§27.38 / §27.42〜§27.47 API / SDK / UI 連動固定契約：**

| 節 | API | SDK | UI |
|----|-----|-----|----|
| §27.21 | branch config API と status/history/log に反映する。 | `getBranchConfig()` / `setBranchConfig()` は `target_files` を削除しない。 | リポジトリ情報 panel で複数 target を表示 / 保存する。 |
| §27.22 | pipeline config API。 | `getPipelineConfig()` / `setPipelineConfig()`。 | 設定 panel で pipeline config を表示 / 保存し、shell 文字列へ変換しない。 |
| §27.23 | `GET/POST /api/config` の `watch_mode`。 | `getConfig()` / `setConfig()`。 | 設定 panel で `github` / `local` を選択する。 |
| §27.24 | `GET/POST /api/config` の `tag_filter`。 | `getConfig()` / `setConfig()`。 | 設定 panel で tag filter を表示 / 保存する。 |
| §27.25 | `GET/POST /api/config` の cache key と builder report。 | `getConfig()` / `setConfig()`。 | 設定 panel と build result 表示で cache counts を表示する。 |
| §27.26 | `GET/POST /api/config` の `deploy_parallelism` と status/history。 | `getConfig()` / `setConfig()`。 | 設定 panel で parallelism を表示 / 保存し、結果は履歴/logで表示する。 |
| §27.27 | hooks API。 | hook methods。 | フック panel で command_args を 1 行 1 引数として表示 / 保存する。 |
| §27.28 | endpoint 追加なし。manifest は runner/builder 内部状態。 | SDK method 追加なし。 | UI 操作追加なし。build result で依存情報を表示してよい。 |
| §27.29 | remote build config は config / pipeline 系 API に含める。 | `getConfig()` / `setConfig()` または pipeline method。 | 設定 panel で remote build 設定を表示 / 保存する。 |
| §27.30 | approvals API。 | approval methods。 | 承認待ち panel で approve / reject を操作する。 |
| §27.31 | branch config または config API に含める。 | branch/config method。 | リポジトリ情報 panel で branch env を表示 / 保存する。 |
| §27.32 | notify config / notify log API。 | notify methods。 | 通知設定 panel で channel、payload、履歴を表示する。 |
| §27.33 | stats/build trends API。 | `getBuildTrends()`。 | 統計 panel で trend を表示する。 |
| §27.34 | build chain config API。 | `getBuildChainConfig()` / `setBuildChainConfig()`。 | 設定またはリポジトリ情報 panel で chain を表示 / 保存する。 |
| §27.35 | queue API。 | `getQueue()` / `clearQueue()` と build trigger methods。 | 手動実行 panel で priority 反映後 queue を表示する。 |
| §27.36 | history/stats API の filter と response に含める。 | history / stats methods は category を保持する。 | 履歴 panel で failure category filter を表示する。 |
| §27.37 | output meta / status / history log に含める。 | `getOutputMeta()` / `getStatus()` / `getHistoryLog()`。 | システム情報または履歴 detail に表示する。 |
| §27.38 | config / stats / dashboard alerts に含める。 | config / stats / dashboard methods。 | 統計 panel と dashboard alert で anomaly を表示する。 |
| §27.42 | endpoint ごとの scope 判定。 | SDK は token scope を推測せず API error を返す。 | UI は `403` を権限不足として表示し、logout しない。 |
| §27.43 | token API。 | token methods。 | API token 管理 panel で発行 token を 1 回だけ表示する。 |
| §27.44 | audit log API。 | `getAuditLog()`。 | 監査ログ panel で filter 表示し、secret は表示しない。 |
| §27.45 | config API。 | `setConfig({session_timeout_seconds})`。 | セキュリティ panel で session timeout を表示 / 保存する。 |
| §27.46 | TOTP/auth API。 | TOTP/auth methods。 | セキュリティ/login panel で one-time secret/ticket flow を扱う。 |
| §27.47 | rate limit API。 | `getApiRateLimit()` / `setApiRateLimit()`。 | セキュリティ panel で policy と state summary を表示する。 |

API / SDK / UI のいずれも、上表に存在しない補完 endpoint、補完 method、補完 UI 操作を追加してはならない。個別節が endpoint 追加なしとする機能は、runner / builder の内部挙動または既存 response field の範囲で実装する。

### 27.30 ビルド承認フロー
owner component は `api` とする。collaborator component は `runner`、`sdk`、`ui`、`statefile` とする。

本機能の目的は、本番向けなど approval_required な target の build / deploy を人間承認後にだけ実行することである。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 設定 key | `branch_targets[].approval_required` |
| 状態 | `.approval_queue` JSON Lines。mode `600`。 |
| API | `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject` |
| status | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| timeout | `.server_config.approval_timeout_seconds`。既定値 86400。 |

**正常系：**

1. runner は差分検出後、approval_required target について build を開始せず approval entry を作成する。
2. `.notify_config` に従い approval request 通知を送信する。
3. API approve 後、queue entry に `trigger="approval"` を追加する。
4. reject 後は build せず `.build_history.status="approval_rejected"` を記録する。
5. timeout 超過 entry は runner 起動時に `expired` へ更新する。

**approval queue / notify 固定契約：**

| 項目 | 仕様 |
|------|------|
| pending id | `appr{YYYYMMDDHHmmss}`、同秒衝突時は `-001` から連番。既存 id は再利用しない。 |
| queue payload | approve で追加する queue entry は `trigger:"approval"`、`priority:"normal"`、`requested_by` は actor id、`payload.approval_id` を含める。 |
| 通知 payload | `{event:"approval_required", approval_id, branch, sha, target, expires_at}`。secret、token、path secret は含めない。 |
| 重複 pending | 同一 branch / sha / target の pending がある場合、新規通知は送らず既存 id を返す。 |
| expired 更新 | runner 起動時に期限超過 pending をすべて処理し、created_at 昇順で expired record を追記する。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| entry 不在 | API は `404`。 |
| pending 以外への approve/reject | `409`。 |
| 通知失敗 | approval entry は残し、`.notify_pending` に追記する。 |
| queue full | approve API は `429`。approval status は pending のまま。 |

**SDK / UI：**

SDK は `getApprovals()`、`approveBuild(id)`、`rejectBuild(id)` を提供する。UI は pending 件数、branch、sha、target、created_at、expires_at、approve/reject 操作を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| approval required | build せず pending 作成。 |
| approve | queue 追加、trigger approval。 |
| reject | build なし、history 記録。 |
| timeout | expired、build なし。 |

**`.approval_queue` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | `appr{YYYYMMDDHHmmss}`、衝突時 `-001`。 |
| `status` | string | 必須 | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |
| `branch` | string | 必須 | branch target 名。 |
| `sha` | string | 必須 | 40 文字 lowercase hex。 |
| `target` | string | 必須 | branch target id または target file。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `expires_at` | string | 必須 | UTC ISO 8601。 |
| `decided_at` | string/null | 必須 | approve / reject / expire 時刻。 |
| `decided_by` | string/null | 必須 | 管理 session は `"admin"`、API token は token id。 |
| `queue_id` | string/null | 必須 | approve で追加した queue id。 |
| `reason` | string/null | 必須 | reject 理由または expire 理由。 |

`.approval_queue` は JSON Lines append-only とする。同一 id の最新 record を有効状態として扱い、古い record は監査履歴として残す。`GET /api/approvals` は id ごとに最新 record だけを返し、`created_at` 降順、同時刻は id 昇順で並べる。壊れた行は無視し、response に含めない。

**approval 状態遷移固定契約：**

| 現在 status | 操作 | 次 status | 副作用 |
|-------------|------|-----------|--------|
| `pending` | approve | `approved` | `.build_state.queued[]` に `trigger:"approval"` entry を追加し、`queue_id` を保存する。 |
| `pending` | reject | `rejected` | `.build_history` に `status:"approval_rejected"` を追記する。queue は追加しない。 |
| `pending` | timeout | `expired` | `.build_history` に `status:"approval_expired"` を追記する。queue は追加しない。 |
| `approved` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `rejected` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `expired` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |

**approval 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| pending 作成 | `.approval_queue` lock → 重複確認 → pending record append → 通知送信 → `.notify_log` / `.notify_pending` 更新 | 通知失敗でも pending は残す。pending append 失敗時は build を開始せず runner failure。 |
| approve | `.approval_queue` lock → 最新 pending 確認 → `.build_state` lock → queue append → approved record append → response | queue full は `429`、approval は pending のまま。approved append 失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| reject | `.approval_queue` lock → 最新 pending 確認 → rejected record append → `.build_history` append → response | history append 失敗時は `500`。rejected record は巻き戻さない。 |
| timeout | runner 起動時に `.approval_queue` lock → expires_at 超過 pending を expired record append → `.build_history` append | history append 失敗時も expired record は残し、runner は ERROR を出して継続する。 |

pending 作成時の重複判定は `branch`、`sha`、`target` が同一で、最新 status が `pending` の record とする。重複時は新規 record を作成せず、既存 pending id を使用する。approve / reject API は body を受け付けない。reject reason は初期実装では固定 `"rejected"` とする。

**approval fixture 固定：**

| fixture | 入力 | 期待結果 |
|---------|------|----------|
| approval-create | approval_required target に差分 | build なし、pending record、通知成功または pending。 |
| approval-duplicate | 同一 branch/sha/target を再検出 | pending 重複作成なし。 |
| approval-approve | pending approve | queue 追加、approved record、queue_id 保存。 |
| approval-reject | pending reject | rejected record、history `approval_rejected`。 |
| approval-timeout | expires_at 超過 | expired record、history `approval_expired`。 |
| approval-queue-full | max_size 到達時 approve | `429`、status pending 維持。 |
| approval duplicate notify | 重複時は通知を送らない。 |
| approved append failure | queue は残り、API は `500`。 |

### 27.42 ビルドトリガー専用 API スコープ
owner component は `api` とする。collaborator component は `sdk`、`ui`、`security`、`statefile` とする。


本機能の目的は、外部システムが最小権限で build を開始できる API token を発行できるようにすることである。

**scope：**

| scope | 許可 |
|-------|------|
| `trigger` | `POST /api/build`、`POST /api/build/force` のみ。 |
| `read` | 読み取り endpoint のみ。 |
| `operate` | 運用操作。 |
| `config` | 設定変更。 |
| `admin` | 管理操作。 |

`trigger` scope は cancel、queue clear、config、token、audit、PAT 更新を許可しない。

**endpoint group 対応：**

| scope | 許可 endpoint |
|-------|---------------|
| `trigger` | `POST /api/build`, `POST /api/build/force` |
| `read` | `GET /api/status`, `GET /api/logs`, `GET /api/logs/search`, `GET /api/logs/export`, `GET /api/history`, `GET /api/history/export`, `GET /api/history/{id}/log`, `GET /api/history/{id}/comment`, `GET /api/sysinfo`, `GET /api/health`, `GET /api/schedule`, `GET /api/notify-config`, `GET /api/notify-log`, `GET /api/config`, `GET /api/config-log`, `GET /api/pat-status`, `GET /api/stats`, `GET /api/stats/timeline`, `GET /api/stats/build-duration`, `GET /api/stats/build-trends`, `GET /api/output-meta`, `GET /api/repo-info`, `GET /api/branch-config`, `GET /api/backup`, `GET /api/dashboard`, `GET /api/diagnostics`, `GET /api/rate-limit`, `GET /api/disk-usage`, `GET /api/webhook-events`, `GET /api/webhook-config`, `GET /api/snapshots`, `GET /api/snapshots/{id}/download`, `GET /api/maintenance`, `GET /api/access-control`, `GET /api/hooks`, `GET /api/hooks/{id}/log`, `GET /api/alert-rules`, `GET /api/tag-rules`, `GET /api/pipeline-config`, `GET /api/build-chain-config`, `GET /api/notes`, `GET /api/smtp-config`, `GET /api/queue`, `GET /api/approvals`, `GET /api/dashboard-layout` |
| `operate` | `POST /api/build/cancel`, `GET /api/build/stream`, `POST /api/notify-test`, `POST /api/notify/weekly-summary`, `POST /api/pat-verify`, `POST /api/circuit-breaker/reset`, `DELETE /api/queue`, `POST /api/smtp-test`, `POST /api/verify-output`, `POST /api/history/{id}/rollback`, `POST /api/approvals/{id}/approve`, `POST /api/approvals/{id}/reject` |
| `config` | `POST /api/schedule/interval`, `POST /api/schedule/pause`, `POST /api/schedule/resume`, `POST /api/schedule/allowed-hours`, `POST /api/schedule/force-interval`, `POST /api/schedule/cooldown`, `POST /api/notify-config`, `POST /api/config/validate`, `POST /api/config`, `POST /api/log-level`, `POST /api/pat-update`, `POST /api/repo-config`, `POST /api/branch-config`, `POST /api/restore`, `POST /api/webhook-config`, `DELETE /api/snapshots/{id}`, `POST /api/maintenance/enable`, `POST /api/maintenance/disable`, `POST /api/access-control`, `POST /api/hooks`, `DELETE /api/hooks/{id}`, `POST /api/alert-rules`, `DELETE /api/alert-rules/{id}`, `POST /api/tag-rules`, `DELETE /api/tag-rules/{id}`, `POST /api/pipeline-config`, `POST /api/build-chain-config`, `POST /api/notes`, `POST /api/smtp-config`, `POST /api/dashboard-layout`, `POST /api/history/{id}/comment`, `POST /api/history/{id}/flag`, `POST /api/history/{id}/tags`, `POST /api/logs/cleanup`, `POST /api/logs/archive` |
| `admin` | `GET /api/access-log`, `GET /api/api-access-log`, `GET /api/audit-log`, `GET /api/api-rate-limit`, `POST /api/api-rate-limit`, `GET /api/sessions`, `POST /api/sessions/revoke-all`, `GET /api/auth/totp-status`, `POST /api/auth/totp-setup`, `POST /api/auth/totp-confirm`, `DELETE /api/auth/totp`, `GET /api/tokens`, `POST /api/tokens`, `DELETE /api/tokens/{id}` |

管理 session は上表に関係なく全 endpoint を許可する。API token が複数 scope を持つ場合は、いずれか 1 つの scope が endpoint に一致すれば許可する。`POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password` は API token scope 判定の対象外とし、API token では使用できない。`POST /api/webhook` は GitHub Webhook secret 検証専用であり、API token では使用できない。上表に存在しない endpoint は §22.0e と本表へ追加されるまで API token では許可してはならない。

**正常系：**

1. `POST /api/tokens` は `scopes:["trigger"]` を受け付ける。
2. 発行 token は `Authorization: Bearer <token>` で使用する。
3. `trigger` token による build は `.audit_log` と `.build_logs/{id}.json.trigger_actor` に token id を記録する。

**判定順：**

1. Bearer token の形式を検査する。`act_` prefix の場合は API token、64 文字 hex の場合は session token として扱う。どちらにも一致しない場合は `401`。
2. API token の hash、期限、失効状態を検証する。
3. endpoint group 対応表で scope を判定する。
4. 許可時だけ endpoint 処理へ進む。
5. 拒否時は `.audit_log` に `permission_denied`、`actor_type:"api_token"`、`actor_id:{token_id}`、`target_type:"endpoint"`、`target_id:{METHOD + " " + path}`、`result:"denied"` を追記してから `403` を返す。

`permission_denied` の監査ログ追記に失敗した場合は `500` を返し、対象 endpoint の処理は実行しない。

**認証・rate limit 組み合わせ順：**

| 段階 | 処理 |
|------|------|
| 1 | route match と method check を行う。未定義は `404` / `405` を返し、API token scope 判定は行わない。 |
| 2 | 認証不要 endpoint か判定する。`GET /api/health` と `POST /api/webhook` は API token scope 判定対象外。 |
| 3 | login endpoint は §27.47 の login rate limit を先に判定する。 |
| 4 | 認証必須 endpoint は Bearer token を検証し、actor を確定する。 |
| 5 | API token の場合だけ scope 判定を行う。管理 session は scope 表を参照しない。 |
| 6 | 認証後 endpoint の rate limit を §27.47 に従って判定する。 |
| 7 | endpoint 固有処理へ進む。 |

scope 判定前に endpoint 固有の request body parse、状態ファイル更新、外部送信を実行してはならない。

**scope 判定固定条件：**

| 条件 | 判定 |
|------|------|
| session token | scope 表を参照せず許可する。 |
| API token の `scopes` が空 | token record 破損として扱い `500`。 |
| API token の `scopes` に未知値 | token record 破損として扱い `500`。 |
| endpoint が複数 group に現れる | §27.42 の表で先に定義された group を採用する。 |
| path parameter 付き endpoint | 正規化 path pattern で scope 判定する。例: `DELETE /api/tokens/abc` は `DELETE /api/tokens/{id}`。 |
| query string | scope 判定には使用しない。method と path pattern だけで判定する。 |
| `GET /api/health` | 認証不要。API token scope 判定を行わない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| trigger token で `POST /api/build` | `202`。 |
| trigger token で `GET /api/status` | `403`。`read` を併用した場合のみ成功。 |
| trigger token で config | `403`。 |
| 未定義 endpoint | token scope に関係なく `404` または `405`。 |
| 複数 scope | いずれかに一致する endpoint だけ成功。 |
| path parameter endpoint | pattern 正規化後の scope で判定する。 |
| query 付き endpoint | query を除外して scope 判定する。 |
| scope 前 body | 権限不足時に request body validation や状態更新を行わない。 |

### 27.43 API キー管理
owner component は `api` とする。collaborator component は `sdk`、`ui`、`security`、`statefile` とする。


本機能の目的は、API key の発行、一覧、失効、期限、scope を実装し、key 本体を保存しないことである。

**API 権限：**

| endpoint | 必要認証 |
|----------|----------|
| `GET /api/tokens` | 管理 session または `admin` scope API token |
| `POST /api/tokens` | 管理 session または `admin` scope API token |
| `DELETE /api/tokens/{id}` | 管理 session または `admin` scope API token |

`admin` scope API token で新しい `admin` scope API token を発行してよい。発行者 token と発行対象 token は別 record とし、親子関係は保存しない。

**正常系：**

1. `POST /api/tokens` は `label`、`scopes`、`expires_at` を受け取る。
2. token 本体は `act_` + 32 byte 相当のランダム文字列とする。
3. `.api_tokens` には `sha256(token)` の lowercase hex だけを保存する。
4. response の `token` は作成時 1 回だけ返す。`GET /api/tokens` では返さない。
5. 認証時は hash 一致、`revoked_at == null`、`expires_at == null または now < expires_at` を満たす token だけ有効とする。
6. 作成、認証成功、失効、期限切れ拒否を `.audit_log` へ記録する。ただし token 本体は記録しない。

**`.api_tokens` record schema：**

| キー | 型 | 必須 | 許容値 |
|------|----|------|--------|
| `id` | string | 必須 | `tok` + 6 桁以上の数字。 |
| `label` | string | 必須 | 1〜64 文字。前後空白は保存前に除去し、空になれば `422`。 |
| `token_hash` | string | 必須 | `sha256(token)` の 64 文字 lowercase hex。 |
| `scopes` | array[string] | 必須 | `read`、`trigger`、`operate`、`config`、`admin` の 1〜5 件。保存時は重複除去し固定順へ正規化する。 |
| `created_at` | string | 必須 | UTC ISO 8601。 |
| `last_used_at` | string/null | 必須 | 認証成功後に更新する。 |
| `expires_at` | string/null | 必須 | `null` または UTC ISO 8601。 |
| `revoked_at` | string/null | 必須 | 失効時に UTC ISO 8601 を保存する。 |

`.api_tokens` に未知 key、必須 key 不足、型不一致、未知 scope、空 scopes、hash 形式不正、時刻形式不正がある場合は破損として扱い、token 認証、一覧、作成、失効をすべて `500 {"error":"Internal server error"}` で拒否する。ただし旧 `scope` 文字列から `scopes:[scope]` への正規化は §22.0c の例外に従う。破損内容、hash、token 本体は response、`.access_log`、`.audit_log`、journal に出力しない。

**API token 認証時の副作用境界：**

| ケース | HTTP status | `.api_tokens` | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|---------------|---------------|--------------|------|
| hash 不一致 | `401` | 変更なし | 追記しない | 追記しない | token id が特定できないため監査対象外。 |
| 期限切れ | `401` | `last_used_at` 更新なし | `token_expired` を追記 | `token_expired` を追記 | token id が特定できた場合だけ記録する。 |
| 失効済み | `401` | `last_used_at` 更新なし | `token_revoked_reject` を追記 | `token_revoked_reject` を追記 | token 本体と hash は保存しない。 |
| scope 不足 | `403` | `last_used_at` 更新なし | `permission_denied` を追記 | `permission_denied` を追記 | 対象 endpoint は実行しない。 |
| 認証成功 | endpoint 固有 | `last_used_at` を現在時刻へ更新 | `token_auth` を追記 | `token_auth` を追記 | endpoint 実行前に更新する。 |

上表で `.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、対象 API は `500` とし、endpoint 固有処理へ進まない。`last_used_at` 更新済み record は巻き戻さない。

**処理順：**

`POST /api/tokens` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. `label`、`scopes`、`expires_at` を検証する。
3. `.api_tokens` のファイルロックを取得する。
4. 既存 record を読み込み、id を採番する。
5. token 本体を生成し、hash を算出する。
6. record を append して `.api_tokens` を原子的に保存する。
7. `.access_log` に token 作成成功を追記する。
8. `.audit_log` に `token_create` を追記する。
9. response に token 本体を 1 回だけ含めて返す。

`.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、API は `500` を返す。作成済み token record は削除せず、再実行時は新しい token を発行する。

`DELETE /api/tokens/{id}` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. path `id` を検証する。
3. `.api_tokens` をロックして対象 record を検索する。
4. 対象が存在しない、または `revoked_at != null` の場合は `404` を返す。
5. `revoked_at` を現在時刻に設定して保存する。
6. `.access_log` と `.audit_log` に失効成功を追記する。
7. `{ "message": "Token revoked" }` を返す。

認証に使用中の API token 自身を失効してよい。その場合、当該リクエストは成功し、次リクエストから `401` になる。

**token ID 採番・返却固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.api_tokens.tokens[].id` の数値 suffix 最大値。存在しない場合は `tok000001`。 |
| 衝突時 | 最大 suffix + 1 を採用する。削除済みや失効済み id は再利用しない。 |
| token 本体 | `act_` + `crypto/rand` 32 bytes を base64url padding なしで encode した文字列。 |
| 作成 response | `{ "id", "label", "scopes", "created_at", "expires_at", "revoked_at", "last_used_at", "token" }`。`token` はこの response だけに含める。 |
| 一覧 response | `tokens` 配列に `token_hash` と `token` を含めない。並び順は `created_at` 降順、同時刻は `id` 昇順。 |
| label 正規化 | 前後空白を除去し、内部空白は保持する。64 文字判定は Unicode code point 数ではなく UTF-8 byte 数で行う。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| scopes 空 | `422`。 |
| 未知 scope | `422`。 |
| label 空 | `422`。 |
| expires_at が現在以前 | `422`。 |
| token 上限 100 件超過 | `422`。失効済み record も件数に含める。 |
| 期限切れ token | `401`。 |
| 失効済み token 再失効 | `404` または冪等成功にせず `404` 固定。 |
| `.api_tokens` 破損 | `500`。自動再生成しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 作成 | token 本体は response だけ。 |
| 一覧 | token hash と本体を返さない。 |
| 期限切れ | `401`、last_used_at 更新なし。 |
| 失効 | `revoked_at` 保存、以後 `401`。 |
| 自己失効 | 失効リクエストは `200`、以後同 token は `401`。 |
| log 追記失敗 | `500`。token record は残る。 |
| token record 破損 | `500`、自動再生成なし、secret 出力なし。 |
| scope 不足 | `403`、endpoint 実行なし、`permission_denied` 記録。 |
| 採番衝突 | 既存最大 suffix + 1 で作成される。 |
| 作成 response | token 本体は 1 回だけ含まれ、一覧では返らない。 |

### 27.44 監査ログ
owner component は `api` とする。collaborator component は `security`、`statefile` とする。


本機能の目的は、認証、権限拒否、token、設定、build trigger などの重要操作を追跡できる JSON Lines 監査ログとして保存することである。

**対象 action：**

| action | target_type |
|--------|-------------|
| `login_success`, `login_failure`, `logout`, `totp_required`, `totp_failure`, `totp_enabled`, `totp_disabled` | `auth` |
| `password_change` | `auth` |
| `token_create`, `token_revoke`, `token_auth`, `token_expired`, `token_revoked_reject` | `api_token` |
| `permission_denied` | `endpoint` または `auth` |
| `build_trigger`, `build_force_trigger` | `build` |
| `config_update`, `rate_limit_update` | `config` |

**record 生成規則：**

| 項目 | 仕様 |
|------|------|
| `timestamp` | 操作結果が確定した UTC 時刻。 |
| `request_id` | リクエスト受付時に生成した 16 byte hex。同一 API 処理中に複数 log を書く場合は同じ値を使う。 |
| `actor_type` | 未認証 login は `"anonymous"`、管理 session は `"admin"`、API token は `"api_token"`、内部処理は `"system"`。 |
| `actor_id` | 管理 session は `"admin"`、API token は token id、未認証は `null`、内部処理は `"system"`。 |
| `target_id` | 対象 id がある場合は id。endpoint 拒否は `"{METHOD} {path}"`。対象なしは `null`。 |
| `result` | 成功は `"success"`、認証失敗や検証失敗は `"failure"`、権限拒否は `"denied"`。 |
| `message` | 固定文言のみ。入力値を連結しない。最大 500 文字。 |

監査ログへ保存する object は `.audit_log` schema のキーだけとする。未知キー、request body、query 全体、header 全体、cookie、secret、token、password、hash、salt、TOTP secret を保存してはならない。

**action 生成固定値：**

| 操作 | action | actor_type | target_type | target_id | result |
|------|--------|------------|-------------|-----------|--------|
| password login 成功 | `login_success` | `anonymous` | `auth` | `login` | `success` |
| password login 失敗 | `login_failure` | `anonymous` | `auth` | `login` | `failure` |
| login ロック拒否 | `permission_denied` | `anonymous` | `auth` | `login` | `denied` |
| TOTP 必須 ticket 発行 | `totp_required` | `anonymous` | `auth` | `login` | `success` |
| TOTP login 失敗 | `totp_failure` | `anonymous` | `auth` | `login_totp` | `failure` |
| logout | `logout` | `admin` または `api_token` | `auth` | `logout` | `success` |
| password 変更 | `password_change` | `admin` | `auth` | `password` | `success` |
| API token 作成 | `token_create` | `admin` または `api_token` | `api_token` | 作成 token id | `success` |
| API token 認証成功 | `token_auth` | `api_token` | `api_token` | token id | `success` |
| API token 期限切れ拒否 | `token_expired` | `api_token` | `api_token` | token id | `failure` |
| API token 失効済み拒否 | `token_revoked_reject` | `api_token` | `api_token` | token id | `failure` |
| API token 失効 | `token_revoke` | `admin` または `api_token` | `api_token` | 失効 token id | `success` |
| scope 不足 | `permission_denied` | `api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| rate limit 超過 | `permission_denied` | `anonymous`、`admin`、`api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| TOTP 有効化 | `totp_enabled` | `admin` | `auth` | `totp` | `success` |
| TOTP 無効化 | `totp_disabled` | `admin` | `auth` | `totp` | `success` |
| 設定変更 | `config_update` | `admin` または `api_token` | `config` | 変更 key | `success` |
| rate limit 設定変更 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `success` |

同一操作で `.config_log` と `.audit_log` の両方を追記する場合、`.config_log` を先に追記する。`.config_log` 成功後に `.audit_log` が失敗した場合は `500` を返し、`.config_log` は巻き戻さない。`.audit_log` 追記失敗そのものを `.audit_log` に記録しようとしてはならない。

**監査 record 保存順・失敗契約：**

| 操作種別 | 保存順 | `.audit_log` 失敗時 |
|----------|--------|----------------------|
| 認証成功 | session または token 状態更新 → `.access_log` → `.audit_log` → response | token / session 状態は巻き戻さず `500`。response に token 本体を含めない。 |
| 認証失敗 | `.access_log` → `.audit_log` → response | `500`。失敗理由詳細は返さない。 |
| 権限拒否 | `.access_log` → `.audit_log` → `403` | `500`。対象 endpoint は実行しない。 |
| 設定変更 | 対象設定保存 → `.config_log` → `.audit_log` → response | 対象設定と `.config_log` は巻き戻さず `500`。 |
| token 作成 | `.api_tokens` 保存 → `.access_log` → `.audit_log` → response | 作成済み record は残し、token 本体は返さず `500`。 |

**正常系：**

1. 監査対象操作の成否が確定した後、`.audit_log` へ 1 行追記する。
2. 監査ログ追記に失敗した場合、対象操作は失敗扱いにし、`500` を返す。
3. secret、password、session token、API token 本体、hash、salt、TOTP secret は保存しない。
4. `GET /api/audit-log` は `limit`、`offset`、`actor`、`action`、`result` で絞り込み、新しい順で返す。

**取得仕様：**

| 項目 | 仕様 |
|------|------|
| 読み込み順 | ファイル先頭から全有効行を読み、フィルタ後に `timestamp` 降順、同時刻はファイル出現順の逆順で返す。 |
| `total` | フィルタ後、ページング前の有効 record 件数。 |
| `actor` filter | `actor_id` と完全一致。`anonymous` を指定した場合は `actor_type:"anonymous"` かつ `actor_id:null` に一致させる。 |
| `action` filter | `action` 完全一致。 |
| `result` filter | `success`、`failure`、`denied` の完全一致。 |
| 壊れた行 | 無視する。API response に壊れた行の内容を含めない。 |
| 返却上限 | `limit` は §22.0b の 1〜200。未指定時は 100。 |

`GET /api/audit-log` は監査ログ取得操作自体を `.audit_log` へ記録しない。通常の API request として `.api_access_log` には記録する。フィルタ値が未知 action、未知 result、200 文字超過 actor の場合は `422` とし、壊れた行の有無とは独立して判定する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| password 変更 | `password_change` が記録される。 |
| 権限拒否 | `permission_denied` が記録される。 |
| secret 入力 | 監査ログに平文がない。 |
| 壊れた行 | API は無視して返す。 |
| filter | `actor`、`action`、`result` が完全一致で絞り込まれる。 |
| 追記失敗 | 対象操作は `500`。 |
| 取得操作 | `GET /api/audit-log` 自身は監査ログへ追記されない。 |
| 未知 filter | `422`。 |
| audit failure token create | token record は残るが token 本体は返らない。 |
| body secret | request body 全体が保存されない。 |

### 27.45 セッションタイムアウト変更設定
owner component は `api` とする。collaborator component は `sdk`、`ui`、`security`、`statefile` とする。


本機能の目的は、新規 session の有効期限を管理 API から変更可能にし、既存 session への影響を明確にすることである。

**仕様：**

| 項目 | 値 |
|------|----|
| 設定 key | `.server_config.session_timeout_seconds` |
| 既定値 | `28800` |
| 最小値 | `300` |
| 最大値 | `2592000` |
| 更新 API | `POST /api/config` または `setConfig({session_timeout_seconds})` |

設定変更は新規 session にだけ適用する。既存 session の `expires_at` は延長も短縮もしない。

**処理契約：**

| 項目 | 仕様 |
|------|------|
| 保存先 | `.server_config.session_timeout_seconds`。他 key と同時更新された場合も §22.0a の単一ファイル更新手順で保存する。 |
| 既定値 merge | `.server_config` に key がない場合、`GET /api/config` は `28800` を返す。ファイルへ暗黙保存しない。 |
| login 時適用 | session token 発行直前に `.server_config` を読み、当該時点の値で `expires_at` を計算する。 |
| TOTP login | TOTP 有効時は `POST /api/login/totp` の成功時点で値を読む。`POST /api/login` の password 成功時点では session を発行しない。 |
| 変更監査 | `POST /api/config` で値が変わった場合は `.config_log` に差分を記録する。`.audit_log` は `config_update`、`target_type:"config"`、`target_id:"session_timeout_seconds"` を記録する。 |
| 同値更新 | 同じ値の更新は `200` とし、`.server_config` の再保存は行ってよいが、差分なしとして `.config_log` と `.audit_log` には記録しない。 |

session timeout の値は session 発行時に秒単位で加算する。`expires_at = issued_at + session_timeout_seconds` とし、計算後の時刻は UTC ISO 8601 秒精度で保存する。ミリ秒、ナノ秒、local timezone は保存しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 300 秒設定 | 新規 session が 300 秒後に期限切れ。 |
| 既存 session | 設定変更後も元の `expires_at`。 |
| 範囲外 | `422`。 |
| key 不在 | `GET /api/config` は `28800`。 |
| TOTP login | `POST /api/login/totp` 成功時点の値で session 期限を決める。 |
| 秒精度 | `expires_at` は UTC ISO 8601 秒精度。 |

### 27.46 TOTP 二要素認証
owner component は `api` とする。collaborator component は `sdk`、`ui`、`security`、`statefile` とする。


本機能の目的は、外部ライブラリなしで RFC 6238 TOTP を検証し、password 漏えい時の管理 API 不正利用を抑止することである。

**アルゴリズム：**

| 項目 | 仕様 |
|------|------|
| secret | `crypto/rand` 20 bytes を RFC 4648 base32 padding なしで表示する。 |
| HMAC | `crypto/hmac` + `crypto/sha1`。 |
| step | 30 秒。 |
| digits | 6 桁。 |
| window | 現在 step の前後 1 step を許可する。 |
| replay 防止 | `.totp_secret.last_accepted_step` 以下の step は拒否する。 |
| otpauth URI | `otpauth://totp/Adlaire%20CI:admin?secret={secret}&issuer=Adlaire%20CI&algorithm=SHA1&digits=6&period=30`。 |

QR code 生成は初期実装対象外とする。UI は secret と otpauth URI を一回表示し、ユーザーが認証アプリへ手入力またはURI貼り付けできるようにする。

**メモリ上状態：**

| 状態 | 保存場所 | 期限 | 内容 |
|------|----------|------|------|
| setup 仮 secret | `components/api.go` のメモリ | 10 分 | `secret_base32`, `created_at`。サーバー再起動で破棄する。 |
| login ticket | `components/api.go` のメモリ | 5 分 | `ticket_hash`, `created_at`, `password_verified_at`。ticket 本体は hash 化して保持する。 |

setup 仮 secret と login ticket は永続ファイルへ保存しない。API response、UI 一回表示、メモリ上状態以外に secret/ticket 本体を残してはならない。

**正常系：**

1. `POST /api/auth/totp-setup` は仮 secret と otpauth URI を返すが、永続化しない。
2. `POST /api/auth/totp-confirm` は仮 secret と code を検証し、成功時だけ `.totp_secret.enabled=true` と secret 情報を保存する。
3. TOTP 有効時の `POST /api/login` は password 成功後に ticket を返し、session token は返さない。
4. `POST /api/login/totp` は ticket と code を検証し、成功時に session token を返す。
5. `DELETE /api/auth/totp` は code を検証して TOTP を無効化する。

**endpoint 処理詳細：**

| endpoint | 処理 |
|----------|------|
| `GET /api/auth/totp-status` | `.totp_secret` を読み、`{enabled,confirmed_at}` だけを返す。`secret_base32` と `last_accepted_step` は返さない。 |
| `POST /api/auth/totp-setup` | TOTP 有効時は `409`。無効時は仮 secret を生成し、既存の未確認仮 secret を上書きする。`.totp_secret` は書き込まない。 |
| `POST /api/auth/totp-confirm` | 仮 secret がない、または期限切れなら `409`。code 成功時に `.totp_secret` を保存し、仮 secret をメモリから削除する。 |
| `DELETE /api/auth/totp` | `.totp_secret.enabled == false` は `409`。code 成功時に `enabled:false`, `secret_base32:null`, `confirmed_at:null`, `last_accepted_step:null` を保存する。 |
| `POST /api/login` | password 成功かつ TOTP 有効なら `ticket` を `crypto/rand` 32 bytes の lowercase hex で生成し、`{must_change,totp_required:true,ticket}` を返す。 |
| `POST /api/login/totp` | ticket hash と code を検証し、成功時に ticket を削除して session token を返す。失敗時も ticket は削除する。 |

TOTP code は 6 桁の ASCII 数字のみ受け付ける。空文字、全角数字、空白付き文字列、6 桁以外は `422` とする。検証は `window` 内の step を古い順に試し、最初に一致した step を採用する。採用 step が `.totp_secret.last_accepted_step` 以下の場合は `401` とする。

TOTP 関連の成功、失敗、無効化、ticket 発行は `.audit_log` へ記録する。code 不一致、replay、ticket 不正は `result:"failure"` とし、code、secret、ticket 本体は保存しない。

**TOTP 計算固定契約：**

| 項目 | 仕様 |
|------|------|
| counter | `floor(unix_seconds / 30)` を 8 byte big-endian unsigned integer として HMAC 入力にする。 |
| truncation | RFC 4226 dynamic truncation を使い、31 bit integer を `10^6` で剰余する。 |
| 表示 | 6 桁未満は左ゼロ埋めする。 |
| base32 decode | 大文字 ASCII のみ保存する。入力確認時は空白を除去せず、保存値と同じ RFC 4648 padding なし形式だけを扱う。 |
| clock source | API server の現在時刻だけを使う。client 時刻は受け取らない。 |

**TOTP 状態更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| setup 開始 | 仮 secret 生成 → メモリ保存 → response | メモリ保存失敗時は `500`、secret を返さない。 |
| confirm 成功 | `.totp_secret` 保存 → 仮 secret 削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み `.totp_secret` と仮 secret 削除は巻き戻さない。 |
| login password 成功 / TOTP 有効 | ticket 生成 → メモリ保存 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は ticket を削除し `500`。 |
| login TOTP 成功 | code step 採用 → `.totp_secret.last_accepted_step` 保存 → ticket 削除 → `.admin_credentials` 更新 → `.access_log` 追記 → `.audit_log` 追記 → session 追加 → response | session 追加前の失敗は token を返さない。ticket 削除後は同 ticket を再利用不可。 |
| login TOTP 失敗 | ticket 削除 → `.access_log` 追記 → `.audit_log` 追記 → response | log 失敗時は `500`。ticket は巻き戻さない。 |
| disable 成功 | `.totp_secret` を無効値で保存 → 未使用 ticket / 仮 secret 全削除 → `.audit_log` 追記 → response | `.audit_log` 失敗時は `500`。保存済み無効化は巻き戻さない。 |

`POST /api/login/totp` は、ticket が存在しない、期限切れ、hash 不一致、既に削除済みのいずれの場合も `401 {"error":"Unauthorized"}` を返す。ticket 不正の詳細、ticket hash、ticket 発行時刻は response と log に含めない。TOTP 有効化または無効化に成功した場合、既存 session は破棄しない。

**異常系：**

| 条件 | 処理 |
|------|------|
| code 不一致 | `401`。 |
| code 形式不正 | `422`。 |
| ticket 期限切れ | `401`。ticket 有効期限は 5 分。 |
| setup 未実行 confirm | `409`。 |
| TOTP 無効状態の disable | `409`。 |
| TOTP 有効状態の setup | `409`。 |
| `.totp_secret` 破損 | TOTP 有効 login、status、confirm、disable は `500`。自動無効化しない。 |
| `.totp_secret` 書き込み失敗 | `500`。response に secret、token、ticket を含めない。 |
| `.audit_log` 追記失敗 | `500`。保存済み `.totp_secret` は巻き戻さない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 setup | secret は確認まで保存されない。 |
| 正常 login | password 後 ticket、TOTP 後 token。 |
| replay | 同一 step の再利用は `401`。 |
| secret 表示 | response と UI の一回表示以外に平文が残らない。 |
| setup 期限切れ | confirm は `409`。 |
| login ticket 再利用 | 1 回成功後または失敗後の同一 ticket は `401`。 |
| confirm audit 失敗 | `500`、`.totp_secret` は保存済み、secret 平文は log なし。 |
| disable 成功 | `.totp_secret` は無効値、既存 session は維持、ticket と仮 secret は削除。 |
| window 前後 | 現在 step の前後 1 step が成功し、同 step 再利用は失敗する。 |
| 全角 code | `422`。 |

### 27.47 API レート制限
owner component は `api` とする。collaborator component は `sdk`、`ui`、`security`、`statefile` とする。


本機能の目的は、login 総当たり、API token 濫用、外部連携の暴走を Go 標準ライブラリだけで抑止することである。

**固定窓 policy：**

| group | 既定 window | 既定 max | 対象 |
|-------|-------------|----------|------|
| `login` | 60 秒 | 10 | `POST /api/login`、`POST /api/login/totp`。 |
| `read` | 60 秒 | 600 | 読み取り endpoint。 |
| `trigger` | 60 秒 | 60 | build trigger endpoint。 |
| `operate` | 60 秒 | 120 | 運用操作 endpoint。 |
| `config` | 60 秒 | 60 | 設定変更 endpoint。 |
| `admin` | 60 秒 | 60 | token、audit、rate limit endpoint。 |

**判定キー：**

| 認証状態 | key |
|----------|-----|
| 認証前 | `ip:{remote_addr}:{group}` |
| session | `session:admin:{group}` と `ip:{remote_addr}:{group}` の両方 |
| API token | `token:{token_id}:{group}` と `ip:{remote_addr}:{group}` の両方 |

どちらか一方でも上限を超えた場合は `429 Too Many Requests` と `{"error":"Too many requests"}` を返す。

`remote_addr` は `net/http.Request.RemoteAddr` の host 部分を使用する。`X-Forwarded-For`、`X-Real-IP`、`Forwarded` header は標準では信用せず、key 生成に使用しない。IPv6 は `net.SplitHostPort` で host を抽出し、正規化済み文字列をそのまま key に入れる。host 抽出に失敗した場合は `ip:unknown:{group}` を使用する。

**policy object：**

```json
{
  "enabled": true,
  "groups": {
    "login": { "window_seconds": 60, "max_requests": 10 },
    "read": { "window_seconds": 60, "max_requests": 600 },
    "trigger": { "window_seconds": 60, "max_requests": 60 },
    "operate": { "window_seconds": 60, "max_requests": 120 },
    "config": { "window_seconds": 60, "max_requests": 60 },
    "admin": { "window_seconds": 60, "max_requests": 60 }
  },
  "state_summary": []
}
```

永続化先は `.server_config.api_rate_limit` とし、保存時は `enabled` と `groups` だけを保存する。`state_summary` は `GET /api/api-rate-limit` と `POST /api/api-rate-limit` response 用の算出値であり、永続化しない。

**`state_summary` item：**

| キー | 型 | 説明 |
|------|----|------|
| `key` | string | `.api_rate_state.windows` の key。 |
| `group` | string | endpoint group。 |
| `window_start` | string | window 開始時刻。 |
| `count` | integer | 現在 count。 |
| `reset_at` | string | `window_start + window_seconds`。 |

`state_summary` は `reset_at` 降順、同時刻は `key` 昇順で最大 100 件返す。期限切れ window は response 算出前に `.api_rate_state` から削除してよい。

**正常系：**

1. path 解決後、認証前に IP key の login 制限を確認する。
2. 認証後、endpoint group を判定し、actor key と IP key の count を更新する。
3. `GET /api/api-rate-limit` は policy と window summary を返す。
4. `POST /api/api-rate-limit` は policy を検証して保存し、`.api_rate_state.windows` を空にする。

**endpoint group 判定：**

| 条件 | group |
|------|-------|
| `POST /api/login`、`POST /api/login/totp` | `login` |
| §27.42 の `trigger` 許可 endpoint | `trigger` |
| §27.42 の `operate` 許可 endpoint | `operate` |
| §27.42 の `config` 許可 endpoint | `config` |
| §27.42 の `admin` 許可 endpoint | `admin` |
| §27.42 の `read` 許可 endpoint | `read` |
| `GET /api/health` | rate limit 対象外 |
| 未知 path / method 不一致 | rate limit 判定前に `404` / `405` |

endpoint が複数 group に現れる場合は、`login`、`admin`、`config`、`operate`、`trigger`、`read` の順で最初に一致した group を採用する。

**判定・更新手順：**

1. `.server_config.api_rate_limit.enabled == false` の場合、`.api_rate_state` を読まずに対象 API 処理へ進む。
2. endpoint group を決める。
3. `.api_rate_state` のファイルロックを取得する。
4. 対象 key ごとに `.api_rate_state.windows[key]` を確認する。
5. window が存在しない、または `now >= window_start + window_seconds` の場合、`window_start=now`, `count=0` で初期化する。
6. `count >= max_requests` の key が 1 つでもあれば、count を増やさず `.audit_log` に `permission_denied` を追記し、`429` を返す。
7. 上限未満の場合、対象 key すべての `count` を 1 増やして `.api_rate_state` を保存し、対象 API 処理へ進む。

rate limit の `429` は `.audit_log` に `permission_denied` として記録する。監査ログ追記に失敗した場合は `500` を返す。login group の認証前 `429` は `actor_type:"anonymous"`、`actor_id:null` とする。

認証後 endpoint の rate limit では、actor key と IP key の両方を同じ lock 内で判定・更新する。片方だけの count 更新に成功した状態を残してはならない。`.api_rate_state` 保存失敗時は対象 API を実行せず `500` を返す。rate limit 判定で `429` になる request は count を増やさない。

**rate limit 副作用固定契約：**

| ケース | `.api_rate_state` | `.access_log` | `.audit_log` | endpoint 固有処理 |
|--------|-------------------|---------------|--------------|-------------------|
| 上限未満 | count を増やす | response 確定後に通常追記 | endpoint が監査対象の場合だけ追記 | 実行する |
| 上限超過 | count を増やさない | `429` として追記 | `permission_denied` を追記 | 実行しない |
| `.api_rate_state` 保存失敗 | 部分更新を残さない | `500` として追記を試行 | 追記しない | 実行しない |
| `.audit_log` 失敗 | count を増やさない | `500` として追記を試行 | 失敗 | 実行しない |

`429` 判定時は `.audit_log` 追記を `.api_rate_state` 保存前に行う。`.audit_log` 追記に成功した場合だけ `429` を返す。`.audit_log` 追記に失敗した場合は `.api_rate_state` を変更せず `500` を返す。

**設定更新契約：**

| 項目 | 仕様 |
|------|------|
| 必須 group | `login`、`read`、`trigger`、`operate`、`config`、`admin` をすべて含める。 |
| 余分な group | `422`。 |
| `enabled` | boolean 必須。 |
| `window_seconds` | integer 必須、1〜86400。 |
| `max_requests` | integer 必須、1〜100000。 |
| 保存順 | `.server_config` 保存 → `.api_rate_state.windows` 空保存 → `.config_log` 追記 → `.audit_log` に `rate_limit_update` 追記 → response。 |
| 同値更新 | `200` とし、`.server_config` と `.api_rate_state` は変更しない。`.config_log` と `.audit_log` に追記しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| policy group 不足 | `422`。 |
| window_seconds 範囲外 | `422`。許容値は 1〜86400。 |
| max_requests 範囲外 | `422`。許容値は 1〜100000。 |
| `.api_rate_state` 書き込み失敗 | `500`。対象 API は実行しない。 |
| unknown group | `422`。 |
| `enabled` 欠落 | `422`。 |
| `.api_rate_state` 破損 | §22.0a に従って退避し、空 window で再生成する。 |
| `.audit_log` 追記失敗 | `500`。`429` response は返さない。 |
| lock 取得 10 秒超過 | `409`。対象 API は実行しない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| login 11 回目 | `429`。 |
| window 経過 | count reset、成功。 |
| session と IP | どちらか超過で `429`。 |
| disabled | `enabled:false` の場合は判定せず成功。 |
| policy 更新 | `.api_rate_state.windows` が空になる。 |
| state_summary | 最大 100 件、`reset_at` 降順で返る。 |
| 429 count | `429` になった request では count が増えない。 |
| 同値更新 | state と log を変更せず `200`。 |
| proxy header | `X-Forwarded-For` ではなく `RemoteAddr` host で key を作る。 |
| audit failure on 429 | count は増えず `500`。 |
| state save failure | endpoint 固有処理なし、部分 count 更新なし。 |
