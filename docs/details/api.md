# Adlaire CI — API 詳細仕様

本ファイルは `docs/DETAIL_INDEX.md` から分割した `api` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、正本関係、実装状態、ロードマップ状態、実装可否の上位判断を記載してはならない。方針、ポリシー、正本関係は `docs/SPEC.md`、実装状態、ロードマップ状態、実装可否は `docs/ROADMAP.md` を正とする。

本ファイルを読む前に、`docs/SPEC.md` で方針とポリシーを確認し、`docs/ROADMAP.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `api` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、security、表示、fixture、検証観点として参照する。

`docs/ROADMAP.md` §6.3 は runner / builder / api / sdk / ui / statefile / archive にまたがる横断補足契約であり、本ファイルへ移動しない。api 連動機能を実装する場合は、本ファイルの個別節を正本とし、横断処理順、成功後再取得、失敗時固定、api / sdk / ui / statefile 同期確認として `docs/ROADMAP.md` §6.3 を確認する。

`docs/details/security.md` §27.42〜§27.47 は security owner component の詳細仕様であり、本ファイルへ移動しない。api が security 機能に関わる場合、本ファイルは endpoint dispatch、request / response、状態ファイル read/write 呼び出し境界だけを担当し、scope、token、audit、session、TOTP、rate limit、漏えい禁止、security 横断順序の主本文は `docs/details/security.md` を正とする。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `api` |
| collaborator component | `statefile`、`sdk`、`ui`、`security`、`archive`、`runner` |
| 持つ内容 | `api` owner が主本文として定義する HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。 |
| 持たない内容 | SDK method 実装、UI DOM 詳細、runner の build 実行責務、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / release 手順、fixture / PR 証跡正本。 |

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

**`api` 設定値（スクリプト冒頭）：**

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
Owner             = "<GitHubオーナー名>"                       // 初期値。POST /api/repo-config の owner 更新成功時だけ .repo_config に保存
Repo              = "<リポジトリ名>"                           // 初期値。POST /api/repo-config の repo 更新成功時だけ .repo_config に保存
```

API service の systemd unit、配置、起動、更新、rollback は setup owner component の責務とし、`docs/details/setup.md` §26.3b、§26.4.2、§26.5 を正とする。

`api` は `adlaire-ci-api --addr 127.0.0.1:8765 --state-dir /opt/adlaire-builder` として起動された後の HTTP listener、request / response、状態ファイル read/write 呼び出し境界だけを定義する。

## 21a. 管理 API サーバー制限

本節は `api` の実行時制限を定義する。runner、setup、admin、sdk、ui は本節の制限を上書きしてはならない。

| 制限 | 詳細 | 実装時の禁止事項 |
|------|------|------------------|
| session はインメモリ管理 | 再起動で全 session を消去する。永続 session store は持たない。 | session を状態ファイル、cookie store、外部 DB、外部 cache に保存しない。 |
| HTTPS listener 非対応 | `api` は HTTP listener のみ起動する。標準 bind は `127.0.0.1:8765` とする。 | TLS listener、証明書読み込み、HTTPS redirect、外部公開 bind を実装しない。 |
| 外部認証非対応 | 認証は `.admin_credentials`、`.totp_secret`、session、API token で完結する。 | SSO、OAuth、LDAP、SAML、複数ユーザー管理を追加しない。 |
| 独自接続数制限なし | Go 標準ライブラリ `net/http` の標準 server で処理する。API rate limit は `docs/details/security.md` §27.47 の固定窓で行う。 | 独自 worker pool、connection pool、接続数上限、外部 queue を追加しない。 |
| runner 起動責務なし | api は HTTP endpoint の request / response と状態 read/write 呼び出し境界を担当する。 | runner の通常 polling loop、GitHub read、pipeline 実行、build log 確定処理を api 本文へ移動しない。 |

セッション、API token、TOTP、rate limit、audit log の security 主本文は `docs/details/security.md` §27.42〜§27.47 を正とし、本節は API server の実行時境界だけを定義する。

---

## 22. バックエンド API 仕様

**ベース URL：** `http://localhost:{PORT}/api`
**認証：** `Authorization: Bearer {SESSION_TOKEN}`（`/api/login` で取得したセッショントークン）
**レスポンス形式：** JSON

### 22.0 API 共通契約

本節の api は `api` の対象仕様である。実装時は、エンドポイント固有仕様より先に以下の共通契約を満たす。

| 項目 | 仕様 |
|------|------|
| Go バージョン | Go `1.22` 以上。HTTP 実装は Go 標準ライブラリ `net/http` を使用する。 |
| bind | 既定値は `127.0.0.1:8765`。`--addr <host:port>` が指定された場合は、起動中の listen address だけを置換する。`--addr 0.0.0.0:<port>` を指定しても、`api` は TLS listener、origin 制限、IP allowlist、reverse proxy 設定生成を追加実行しない。 |
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

状態ファイルのパス、形式、初期値、更新責務、破損時の扱い、更新手順、schema 厳格化、状態読取 adapter、状態読取 priority は `docs/details/statefile.md` §22.0a を正とする。`api` は同節の adapter と更新手順を利用し、endpoint 固有の request / response / validation は本ファイル §22.0b 以降を正とする。

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

複数エラーがある場合、body key は JSON object の出現順、query は URL query の出現順、path parameter は route 定義順で並べる。body、query、path にまたがる場合は body → query → path の順とする。sdk と ui は `details[].field` と `details[].message` をそのまま扱うため、実装者判断で文言を言い換えてはならない。

### 22.0c 主要状態ファイル schema

主要状態ファイル schema は `docs/details/statefile.md` §22.0c を正とする。本ファイルでは API endpoint と状態ファイルの read / write 対応を §22.0d 以降で定義する。

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

本表は API 実装、SDK 実装、標準管理ツール実装の契約インデックスである。endpoint を追加、削除、名称変更、body 変更、response 変更する場合は、本表、該当 endpoint 個別節、SDK method 表、UI 操作契約、fixture catalog を同じ仕様 PR で先に更新する。下表に存在しない endpoint は実装対象外とする。SHA reset 専用 endpoint とサマリー送信専用 endpoint は定義しない。

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
| `POST /api/logs/cleanup` | none | `{message,deleted_count,failed_count}` | `200` | `401`, `500` | `.server_config`, `.build_logs/` | `.build_logs/` | `cleanupLogs()` | ログビューア, 設定 |
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
| `POST /api/webhook` | GitHub webhook body | `{message,queued,event_id,queue_id?}` | `202` | `401`, `413`, `422`, `429`, `500`, `503` | `.webhook_secret`, `.branch_config`, `.build_state`, `.maintenance`, `.build_circuit_state` | `.webhook_events.json`, `.build_state` or queue | none | 外部 Webhook |
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

**archive / snapshot / rollback 処理参照：**

`POST /api/logs/cleanup`、`POST /api/logs/archive`、`GET /api/snapshots`、`GET /api/snapshots/{id}/download`、`DELETE /api/snapshots/{id}`、`POST /api/history/{id}/rollback` の保存、圧縮、削除、download 安全性、rollback 実体処理は `docs/details/archive.md` §27.7 および §27.15 を正とする。本ファイルでは API endpoint、request / response、HTTP status、read / write 境界だけを定義する。

**backup / restore API 副作用固定契約：**

| API | 処理順序 | 成功時副作用 | 失敗時副作用 |
|-----|----------|--------------|--------------|
| `GET /api/backup` | 対象設定 file 読込 → 不在 file に既定値適用 → secret mask → response。 | 状態ファイルを更新しない。 | 読込不能な必須 file は `500`。任意 file 不在は既定値で返す。 |
| `POST /api/restore` | request 検証 → secret mask `"***"` の既存 secret 再利用判定 → 全対象 payload 生成 → §22.0d の順に atomic write → `.config_log` 追記 → response。 | 設定系状態 file だけを置換する。履歴、ログ、snapshot、session、token 本体は復元しない。 | 検証失敗は差分なし。途中 write 失敗は未処理 file を書かず `500`。処理済み file は戻さない。 |

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

backup / restore fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは backup / restore fixture 本体を重複定義しない。

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
| snapshot id | build id と同一 | `b20260915100500` | build id 衝突時の値と同一。snapshot 専用 id は採番しない。 |

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
| external check | `POST /api/pat-verify`, `GET /api/rate-limit`, `GET /api/diagnostics`, `POST /api/smtp-test`, `POST /api/notify-test` | 設定読込 → timeout 付き外部確認 → 結果 response → endpoint 固有節で定義された log 追記 | 確認結果を保存しない。ただし test 送信 log は仕様どおり追記する。 | 未設定は `501` または endpoint 固有 `422`。timeout は `500`。 |
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
| `POST /api/smtp-config` | config と password 更新有無が既存値と一致 | `.smtp_config`、`.smtp_secret`、`.config_log` を変更しない。 | `.smtp_config` → password 更新がある場合だけ `.smtp_secret` → `.config_log`。 |
| `POST /api/dashboard-layout` | widgets 配列が既存値と一致 | `.dashboard_layout`、`.config_log` を変更しない。 | `.dashboard_layout` → `.config_log`。 |

no-op response は endpoint 固有の `No changes` が定義されている場合はその文言を返す。定義がない endpoint は通常成功文言を返す。no-op では、状態ファイル、JSON Lines、監査ログ、通知ログに差分を作ってはならない。部分更新では未指定 key を保持し、`null` が削除を意味する key は個別節に明記された key だけとする。

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

**api / sdk / ui / statefile 横断契約参照：**

API endpoint、SDK method、UI 操作、状態ファイル副作用の本文は各 owner component 別詳細仕様ファイルを正とする。成功後再取得、失敗時固定、横断処理順、api / sdk / ui / statefile の同期確認は `docs/ROADMAP.md` §6.3 を同時に確認する。本ファイルでは横断連動表と横断処理順表を重複定義しない。

**横断 fixture 参照：**

api / sdk / ui / statefile にまたがる横断 fixture の fixture 名、入力、必須確認は `docs/details/fixture.md` §22-F の cross fixture 固定契約を正とする。本ファイルでは横断 fixture 本体を重複定義しない。

### 22.0f Phase 3 / Phase 4 API fixture 参照

API の Phase 3 / Phase 4 必須検証、fixture 名、入力状態、期待 response、期待副作用は `docs/details/fixture.md` §22-F を正とする。

本ファイルでは、API endpoint の method、path、request、response、error、read / write 境界だけを定義する。fixture manifest、testdata 配置、期待副作用、PR 証跡、Phase 別の完了判定は本ファイルに重複定義しない。

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| `POST` | `/api/login` | 不要 | ログイン（セッショントークン返却） |
| `POST` | `/api/logout` | 要 | ログアウト（セッション破棄） |
| `POST` | `/api/change-password` | 要 | パスワード変更 |
| `GET` | `/api/access-log` | 要 | ログイン、ログアウト、API token 監査ログを返す |
| `GET` | `/api/sessions` | 要 | 有効セッション一覧を返す |
| `POST` | `/api/sessions/revoke-all` | 要 | 現セッション以外の全セッションを強制無効化する |
| `GET` | `/api/status` | 要 | 最終ビルド時刻・SHA・成否・実行中フラグを返す |
| `POST` | `/api/build` | 要 | 手動ビルドトリガー（`runner` を即時起動） |
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
| `POST` | `/api/schedule/force-interval` | 要 | `.server_config.force_build_interval_hours` を保存する |
| `POST` | `/api/schedule/cooldown` | 要 | `.server_config.build_cooldown_seconds` を保存する |
| `GET` | `/api/notify-config` | 要 | Webhook 通知設定を返す |
| `POST` | `/api/notify-config` | 要 | Webhook 通知設定を更新する |
| `GET` | `/api/notify-log` | 要 | Webhook 送信履歴（日時・イベント・HTTP ステータス・成否）を返す |
| `POST` | `/api/notify-test` | 要 | Webhook にテスト通知を送信し疎通を確認する |
| `POST` | `/api/notify/weekly-summary` | 要 | 週次サマリー Webhook を即時手動送信する（過去 7 日間の統計を集計して送信） |
| `GET` | `/api/config` | 要 | サーバー設定を返す |
| `POST` | `/api/config` | 要 | サーバー設定を更新する |
| `POST` | `/api/log-level` | 要 | `api` の `log_level` を変更する |
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
| `POST` | `/api/webhook` | 不要（Secret 検証） | GitHub push Webhook を受信し、署名検証後に build queue へ投入する（→ §22 Webhook 受信仕様 / §27.12） |
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

> **責務分担：** build lifecycle 通知、pending retry、自動 weekly summary の送信責務は `docs/details/runner.md` §27.32 および §27.19 を正とする。本節では通知 API の request / response、設定 read/write、履歴参照、手動送信 endpoint 境界だけを定義する。

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

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"deploy_failure"`（転送失敗時）| `"weekly_summary"`（定期サマリー送信時）| `"approval_required"`（承認待ち発生時）| `"duration_anomaly"`（所要時間異常時）| `"config_corrupt"`（設定破損復旧時）。複数指定は array 順を保持して保存する。

`summary`：週次サマリー通知の設定。`interval` の有効値は `"weekly"` 固定。`hour` は 0〜23（UTC）。`day_of_week` は 0 = 日曜〜6 = 土曜。自動送信条件、二重送信防止、集計、送信順序、`.build_state` 更新は `docs/details/runner.md` §27.19 を正とする。

**`POST /api/notify/weekly-summary` レスポンス例：**
```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

`POST /api/notify/weekly-summary` は手動送信 API とする。集計、対象 channel 抽出、通知送信、`.notify_log` 追記、失敗時 response、sent date を更新しない契約は `docs/details/runner.md` §27.19 の手動 weekly summary 契約を正とする。Webhook 未設定または無効時は `422` を返す。

`payload_template`：Webhook 送信 JSON ペイロードのテンプレート文字列。`null` = デフォルトペイロードを使用。build lifecycle 通知 payload、channel 選択、retry、pending 保存、secret mask は `docs/details/runner.md` §27.32 を正とする。テンプレート内で使用可能な変数は以下の通り。

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
- `uptime_seconds`：`api` 起動からの経過秒数

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
値は `runner` が `.build_logs/{id}.json` または `.build_logs/archive/{id}.json.gz` から最新エントリを読み取って返す。

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

クエリパラメータ `limit`（デフォルト 50、上限 1000）と `offset` でページネーションする。`.webhook_events.json` を逆順（新しい順）で返す。

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
      "branch": "main",
      "sha": "0123456789abcdef0123456789abcdef01234567",
      "repository": "fqwink/Build-Scripts",
      "build_triggered": true,
      "queued_id": "q20260915100200",
      "result": "queued"
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

手動 weekly summary の集計、送信、`.notify_log` 追記、失敗時 response は `docs/details/runner.md` §27.19 を正とする。`on: ["weekly_summary"]` 設定の Webhook 宛先がない場合は `422 Unprocessable Entity` を返す。

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

snapshot の保存、世代削除、download、delete、rollback 実体処理は `docs/details/archive.md` §27.15 を正とする。本節では API request / response と HTTP 境界だけを定義する。

**`GET /api/snapshots` レスポンス例：**
```json
{ "snapshots": [
    { "id": "b20260915100000", "build_id": "b20260915100000", "saved_at": "2026-09-15T10:00:00Z", "size_bytes": 2048576 },
    { "id": "b20260914183000", "build_id": "b20260914183000", "saved_at": "2026-09-14T18:30:00Z", "size_bytes": 2031616 }
]}
```

**`GET /api/snapshots/{id}/download`**
バイナリレスポンス。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"` を付与する。download tar.gz 生成と安全性検証は `docs/details/archive.md` §27.15 を正とする。

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

`.snapshots/{id}/` が存在しない場合は `404 Not Found` を返す。ロールバックは非同期で SSH 転送を実行し、`running: true` 中は `409 Conflict` を返す。rollback の状態更新、転送、履歴、ログ、pending transfer、禁止副作用は `docs/details/archive.md` §27.15 を正とする。

---

### Webhook 受信仕様（22-W）

`POST /api/webhook` は GitHub からの push event を受信し、HMAC-SHA256 署名検証後に build queue へ投入する。認証ヘッダー（`Authorization: Bearer`）は不要とし、Webhook 署名検証を認証代替として扱う。

本節は Webhook 受信 endpoint の概要と必須 header だけを定義する。署名検証、status code、response、event log、queue 投入、重複判定、異常系、検証条件は `docs/details/api.md` §27.12 を正とする。本節へ `POST /api/webhook` の response 例、event log schema、queue entry schema を重複定義してはならない。

**`POST /api/webhook` リクエストヘッダー：**
```
X-GitHub-Event: push
X-GitHub-Delivery: <delivery_id>
X-Hub-Signature-256: sha256=<hmac_hex>
Content-Type: application/json
```

Secret は `.webhook_secret` を正とする。secret 不在、header 不在、prefix 不正、hex 不正、署名不一致はいずれも §27.12 の固定契約どおり `401` とし、event log と queue を変更しない。

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

- `.server_config.force_build_interval_hours` を保存する。runner による読込タイミングと判定適用は `docs/details/runner.md` §13 の cooldown / force build 判定契約を正とする。
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

- `.server_config.build_cooldown_seconds` を保存する。runner による読込タイミングと判定適用は `docs/details/runner.md` §13 の cooldown / force build 判定契約を正とする。
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

メンテナンス fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルではメンテナンス fixture 本体を重複定義しない。

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

アクセス制御 fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルではアクセス制御 fixture 本体を重複定義しない。

---

### ビルドフック（14E）

本節は hooks API の request / response、`.hooks` 保存、`.config_log` 追記、hook log 参照境界だけを定義する。pre / post hook の実行順、timeout、process kill、hook log 保存、secret mask、pre abort、post failure、build status への影響は `docs/details/runner.md` §27.27 を正とする。

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

`.hooks` record schema は `docs/details/statefile.md` §22.0c `.hooks` schema を正とする。`.hooks` に未知 key、必須 key 不足、型不一致、不正 phase、不正 command、重複 id がある場合、`GET /api/hooks`、`POST /api/hooks`、`DELETE /api/hooks/{id}` は `500 {"error":"Internal server error"}` を返す。破損内容、command_args の secret らしき値、stdout/stderr は response と log に出さない。

**hooks API 更新順：**

| API | 更新順 | 失敗時 |
|-----|--------|--------|
| `POST /api/hooks` | 入力検証 → `.hooks` lock → id 採番 → record append → `.hooks` atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。追加済み record は巻き戻さない。 |
| `DELETE /api/hooks/{id}` | path id 検証 → `.hooks` lock → 対象存在確認 → record 削除 → `.hooks` atomic write → `.config_log` 追記 → response | 対象不在は `404`。`.config_log` 失敗時は `500`、削除済み record は巻き戻さない。 |
| `GET /api/hooks/{id}/log` | path id 検証 → `.hooks` で存在確認 → `.build_logs/*_hook_{id}.json` を新しい順で最大 20 件読込 → response | hook 不在は `404`。個別 hook log 破損はその file を除外し、server log に固定コードを出す。 |

hook log JSON の保存 schema、保存タイミング、失敗時の runner 挙動は `docs/details/runner.md` §27.27 を正とする。`GET /api/hooks/{id}/log` は保存済み hook log を読み取り、response の `runs[]` へ `build_id`、`ran_at`、`exit_code`、`output` を返す。`output` は保存済み `stdout + stderr` をこの順で連結した表示用互換値とし、保存時点で secret mask 済みの値だけを返す。

hooks fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは hooks fixture 本体を重複定義しない。

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

本節は、`GET /api/pipeline-config` と `POST /api/pipeline-config` の request / response、`.pipeline_config` read/write 境界だけを定義する。runner による `.pipeline_config` 読込タイミング、`extra_args` / `env` 適用、読込不能または schema 不正時の build 停止条件は `docs/details/runner.md` §27.22 を正とする。

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

Webhook に加えてメールでビルド結果を通知する機能。SMTP 接続設定は `.smtp_config` に、パスワードは `.smtp_secret`（パーミッション 600）に分離して保存する。`GET /api/notify-config` のレスポンスに `email` セクションを追加する。

**`GET /api/smtp-config` レスポンス例：**
```json
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true, "password_set": true }
```
`password_set`：`.smtp_secret` が存在するかを真偽値で返す。パスワード本体は返却しない。
未設定時：`{ "host": null, "port": 587, "user": null, "tls": true, "from": null, "to": [], "on": [], "enabled": false, "password_set": false }`

**`POST /api/smtp-config` リクエスト / レスポンス：**
```json
// リクエスト（下記キーだけを受け付ける。password フィールドは省略可能）
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
| 更新 | `.smtp_config` → password 更新がある場合だけ `.smtp_secret` → `.config_log` の順に書く。途中失敗時は未処理ファイルを書かない。 |
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

SMTP fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは SMTP fixture 本体を重複定義しない。

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

queue entry schema と trigger 別 payload schema は `docs/details/statefile.md` §22.0c `.build_state` schema を正とする。未知 key は `422`、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

**queue 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| API queue 追加 | `.build_state` lock → 最新 state 読込 → running / max_size / 重複確認 → id と created_seq 採番 → atomic write → response | write 失敗は `500`。queue 追加なし。 |
| webhook queue 追加 | event 検証 → `.webhook_events.json` 追記 → `.build_state` lock → queue append → response | event 追記前の失敗は queue なし。queue 追加失敗は event result を `error` にできる場合だけ追記し、response は `500`。 |
| approval approve queue 追加 | `.approval_queue` lock → pending 確認 → `.build_state` lock → queue append → approval status `approved` 追記 → response | queue full は `429`、approval は pending のまま。approval status 追記失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| runner 取り出し | `.build_state` lock → §27.35 の順で 1 件選択 → selected entry 削除 → `running=true` と `current_build_id` 設定 → atomic write | write 失敗は build 開始なし、lock を解放し終了コード `1`。 |
| queue clear | `.build_state` lock → `queued=[]` → atomic write → `.config_log` 追記 → response | `.config_log` 失敗時は `500`。cleared queue は巻き戻さない。 |

重複判定は `trigger` と `payload` の正規化 JSON が一致する waiting entry を対象とする。重複時は新規 entry を追加せず `200 {"message":"Already queued","queued":true,"queue_id":"<existing>"}` を返す。`force=true` の manual entry は `force=false` と別 entry として扱う。

queue fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは queue fixture 本体を重複定義しない。

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
{ "message": "Cleanup completed", "deleted_count": 12, "failed_count": 0 }
```

`log_retention_days` が `0` の場合は削除せず `deleted_count: 0`、`failed_count: 0` を返す。

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
// リクエスト（owner、repo、branch、target_file だけを受け付ける）
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

認証、session、password hash、login ticket、TOTP 連携、認証ログ、漏えい禁止、`--init-credentials` 生成手順の主本文は `docs/details/security.md` の認証共通詳細および §27.45〜§27.46 を正とする。`.admin_credentials` schema は `docs/details/statefile.md` §22.0c を正とする。

本ファイルでは、認証関連 API の endpoint、request / response、HTTP status、状態ファイル read / write 境界だけを定義する。

| API / CLI | API 側の担当 | 主本文 |
|-----------|--------------|--------|
| `POST /api/login` | route、body parse、response body、HTTP status、`.admin_credentials` / `.totp_secret` read/write 呼び出し境界。 | `docs/details/security.md` 認証共通詳細、§27.46 |
| `POST /api/login/totp` | route、body parse、response body、HTTP status、`.totp_secret` / `.admin_credentials` read/write 呼び出し境界。 | `docs/details/security.md` 認証共通詳細、§27.46 |
| `POST /api/logout` | route、body 禁止、response body、HTTP status。 | `docs/details/security.md` 認証共通詳細 |
| `POST /api/change-password` | route、body parse、response body、HTTP status、`.admin_credentials` write 呼び出し境界。 | `docs/details/security.md` 認証共通詳細 |
| `GET /api/sessions` / `POST /api/sessions/revoke-all` | route、response body、HTTP status、memory session 操作呼び出し境界。 | `docs/details/security.md` 認証共通詳細、§27.45 |
| `--init-credentials` | CLI option dispatch、stdout / stderr / exit code を security 契約どおり返す。 | `docs/details/security.md` 認証共通詳細、`docs/details/statefile.md` §22.0c |
| `GET /api/auth/totp-status` / `POST /api/auth/totp-setup` / `POST /api/auth/totp-confirm` / `DELETE /api/auth/totp` | route、body parse、response body、HTTP status、`.totp_secret` read/write 呼び出し境界。 | `docs/details/security.md` §27.46 |

認証 fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは認証 fixture 本体を重複定義しない。

---

## 27. api owner 追加仕様化機能 詳細仕様

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


`api` は全 `/api/` request について `.api_access_log` へ JSON Lines を追記する。`GET /api/health` も対象とする。静的 file 配信、admin HTML、SDK JS は対象外とする。

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
| 状態 | `.webhook_events.json` へ追記する。対象 push の after SHA が対象 branch の直近成功 SHA と異なり、maintenance が無効で、queue 上限に空きがあり、同一 delivery id の queued / running entry が存在しない場合だけ `.build_state.queued` へ `trigger="webhook"` entry を追加する。 |

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

queue entry は `docs/details/statefile.md` §22.0c `.build_state` schema の queue entry schema を使用し、`trigger:"webhook"`、`requested_by:"webhook"`、`priority:"normal"`、`payload.delivery_id`、`payload.branch`、`payload.sha` を保存する。`X-GitHub-Delivery` が既に pending queue に存在し、同一 branch / sha の場合は重複投入せず、event log に `result:"duplicate"`、既存 `queued_id` を記録し、`202 {"message":"Webhook already queued","queued":true,"queue_id":"<existing>"}` を返す。

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

本機能の目的は、受信した GitHub Webhook の監査情報を `.webhook_events.json` に保存し、管理 API、sdk、ui のページング参照対象にすることである。

`.webhook_events.json` の保存 schema は `docs/details/statefile.md` §22.0c `.webhook_events.json` JSON Lines schema を正とする。保存時に request header 全体、署名値、secret、payload 全体を保存してはならない。

**一覧 API：**

`GET /api/webhook-events` は `limit` と `offset` query を受け付ける。`limit` は 1〜1000、既定値 50。`offset` は 0 以上、既定値 0。新しい順で返す。壊れた行は無視し、server log に `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` を出す。

Response は `{ "events": WebhookEventRecord[], "total": N }` とする。SDK `getWebhookEvents(limit,offset)` は `limit` と `offset` を常に query へ送信する。ui は件数、delivery id、event、branch、sha、result、queued id を表示する。

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

本機能の目的は、認証不要の `GET /api/health` で、外部監視へ Adlaire CI の最低限の稼働状態を返すことである。

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

`status` は `"ok"`、`"degraded"`、`"error"` のいずれかとする。必須状態ファイル破損がある場合は `degraded`、API process が応答できるが重大な read error がある場合は `error` とする。HTTP status は、API 自体が response を生成できる限り `200` とする。response object 構築失敗、JSON encode 失敗、response 書き込み開始前の header 生成失敗の場合だけ `500` とする。

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

SDK `searchLogs(q,from,to,level)` は `level` 指定時だけ query に送信する。ui は INFO / WARNING / ERROR / DEBUG の filter control を提供し、選択なしでは全件を表示する。

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

本機能の目的は、監視対象 branch / target / deploy target を `.branch_config` で管理し、API 経由の変更対象を `.branch_config` に限定することである。

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

runner による `.branch_config` の読込、`RunnerConfig.BranchTargets` への正規化、起動中の反映タイミングは `docs/details/runner.md` §12〜§13 の runner 設定正規化契約を正とする。本節では API endpoint、request / response、保存、削除、検証条件だけを定義する。

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

`.config_log` のログ schema は `docs/details/statefile.md` §22.0c `.config_log` JSON Lines schema を正とする。差分がない場合、対象 API は状態ファイルを書かず、`.config_log` も追記せず、response は `{ "message": "No changes" }` とする。

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

§27.21〜§27.38、§27.42〜§27.47 の各機能は、owner component の個別節を主本文とし、下表を横断受け入れ確認として満たした場合だけ実装完了とする。下表は endpoint、SDK method、UI 操作、状態 schema、fixture を新規定義しない。`docs/ROADMAP.md` §6.3 は §27.21〜§27.38 の runner 拡張を横断検証する補足契約として扱う。

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
| §27.30 | approval API | approval 一覧取得、approve/reject API。 | `.approval_queue`、`.build_state.queued[]`、`.build_history`。 | pending entry だけ approve/reject できる。 | entry 不在は `404`、pending 以外は `409`、queue full は `429`。 | list、approve、reject、duplicate approve、queue full。 |
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

**§27.21〜§27.38 / §27.42〜§27.47 api / sdk / ui 連動固定契約：**

| 節 | API | SDK | UI |
|----|-----|-----|----|
| §27.21 | branch config API と status/history/log に反映する。 | `getBranchConfig()` / `setBranchConfig()` は `target_files` を削除しない。 | リポジトリ情報 panel で複数 target を表示 / 保存する。 |
| §27.22 | pipeline config API。 | `getPipelineConfig()` / `setPipelineConfig()`。 | 設定 panel で pipeline config を表示 / 保存し、shell 文字列へ変換しない。 |
| §27.23 | `GET/POST /api/config` の `watch_mode`。 | `getConfig()` / `setConfig()`。 | 設定 panel で `github` / `local` を選択する。 |
| §27.24 | `GET/POST /api/config` の `tag_filter`。 | `getConfig()` / `setConfig()`。 | 設定 panel で tag filter を表示 / 保存する。 |
| §27.25 | `GET/POST /api/config` の cache key と builder report。 | `getConfig()` / `setConfig()`。 | 設定 panel と build result 表示で cache counts を表示する。 |
| §27.26 | `GET/POST /api/config` の `deploy_parallelism` と status/history。 | `getConfig()` / `setConfig()`。 | 設定 panel で parallelism を表示 / 保存し、結果は履歴/logで表示する。 |
| §27.27 | hooks API。 | hook methods。 | フック panel で command_args を 1 行 1 引数として表示 / 保存する。 |
| §27.28 | endpoint 追加なし。manifest は runner/builder 内部状態。 | SDK method 追加なし。 | UI 操作追加なし。依存情報を表示する場合は、既存 build result 表示の範囲だけを使用する。 |
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
| §27.42 | endpoint ごとの scope 判定。 | sdk は token scope を推測せず API error を返す。 | ui は `403` を権限不足として表示し、logout しない。 |
| §27.43 | token API。 | token methods。 | API token 管理 panel で発行 token を 1 回だけ表示する。 |
| §27.44 | audit log API。 | `getAuditLog()`。 | 監査ログ panel で filter 表示し、secret は表示しない。 |
| §27.45 | config API。 | `setConfig({session_timeout_seconds})`。 | セキュリティ panel で session timeout を表示 / 保存する。 |
| §27.46 | TOTP/auth API。 | TOTP/auth methods。 | セキュリティ/login panel で one-time secret/ticket flow を扱う。 |
| §27.47 | rate limit API。 | `getApiRateLimit()` / `setApiRateLimit()`。 | セキュリティ panel で policy と state summary を表示する。 |

api / sdk / ui のいずれも、上表に存在しない endpoint、method、UI 操作を追加してはならない。追加が必要な場合は、本表、該当 endpoint 個別節、SDK method 表、UI 操作契約、fixture catalog を先に更新する。個別節が endpoint 追加なしとする機能は、runner / builder の内部挙動または既存 response field の範囲で実装する。

### 27.30 ビルド承認フロー
owner component は `api` とする。collaborator component は `runner`、`sdk`、`ui`、`statefile` とする。

本節は、ビルド承認フローにおける API endpoint、request / response、状態ファイル read/write 呼び出し境界だけを定義する。`approval_required` 検出、pending 作成、approval request 通知、timeout 処理、承認済み queue entry の実行は `docs/details/runner.md` §27.30 を正とする。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.approval_queue` JSON Lines。mode `600`。 |
| API | `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject` |
| status | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |

**正常系：**

1. `GET /api/approvals` は `.approval_queue` を読み取り、id ごとの最新 record だけを返す。
2. `POST /api/approvals/{id}/approve` は最新 status が `pending` の entry だけを承認し、`.build_state.queued[]` に `trigger:"approval"` の queue entry を追加する。
3. approve 成功後、`.approval_queue` に `approved` record を追記し、response を返す。
4. `POST /api/approvals/{id}/reject` は最新 status が `pending` の entry だけを却下し、`.approval_queue` に `rejected` record を追記する。
5. reject 成功後、`.build_history` に `status:"approval_rejected"` を追記し、response を返す。

**approval API 固定契約：**

| 項目 | 仕様 |
|------|------|
| queue payload | approve で追加する queue entry は `trigger:"approval"`、`priority:"normal"`、`requested_by` は actor id、`payload.approval_id` を含める。 |
| approval id | path parameter の `{id}` と `.approval_queue` の `id` が完全一致する entry だけを対象にする。 |
| list 並び順 | `created_at` 降順、同時刻は id 昇順。 |
| list 対象 | id ごとの最新 record だけを返す。古い record は監査履歴として返さない。 |
| body | approve / reject API は body を受け付けない。 |
| reject reason | 初期実装では固定 `"rejected"` とする。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| entry 不在 | API は `404`。 |
| pending 以外への approve/reject | `409`。 |
| queue full | approve API は `429`。approval status は pending のまま。 |
| `.approval_queue` の壊れた行 | `GET /api/approvals` response には含めない。 |

**SDK / UI：**

sdk は `getApprovals()`、`approveBuild(id)`、`rejectBuild(id)` を提供する。ui は pending 件数、branch、sha、target、created_at、expires_at、approve/reject 操作を表示する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| list | id ごとの最新 record だけを返す。 |
| approve | queue 追加、trigger approval。 |
| reject | build なし、history 記録。 |
| duplicate approve | `409`、状態差分なし。 |
| queue full | `429`、approval は pending のまま。 |

`.approval_queue` record schema は `docs/details/statefile.md` §22.0c `.approval_queue` JSON Lines schema を正とする。同一 id の最新 record を有効状態として扱い、古い record は監査履歴として残す。`GET /api/approvals` は id ごとに最新 record だけを返し、`created_at` 降順、同時刻は id 昇順で並べる。壊れた行は無視し、response に含めない。

**approval 状態遷移固定契約：**

| 現在 status | 操作 | 次 status | 副作用 |
|-------------|------|-----------|--------|
| `pending` | approve | `approved` | `.build_state.queued[]` に `trigger:"approval"` entry を追加し、`queue_id` を保存する。 |
| `pending` | reject | `rejected` | `.build_history` に `status:"approval_rejected"` を追記する。queue は追加しない。 |
| `approved` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `rejected` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `expired` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |

**approval 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| approve | `.approval_queue` lock → 最新 pending 確認 → `.build_state` lock → queue append → approved record append → response | queue full は `429`、approval は pending のまま。approved append 失敗時は `500`、queue 追加済み entry は巻き戻さない。 |
| reject | `.approval_queue` lock → 最新 pending 確認 → rejected record append → `.build_history` append → response | history append 失敗時は `500`。rejected record は巻き戻さない。 |

approve / reject API は body を受け付けない。reject reason は初期実装では固定 `"rejected"` とする。

approval fixture は `docs/details/fixture.md` §22-F の API 機能別 fixture 固定契約を正とする。本ファイルでは approval fixture 本体を重複定義しない。

### 27.42 ビルドトリガー専用 API スコープ

本節の主本文は `docs/details/security.md` §27.42 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

API 側は、route / method 確定、endpoint dispatch、HTTP status、request / response body、状態ファイル read/write 呼び出し境界だけを担当する。scope 判定順、許可 endpoint group、permission denied audit、body parse 前判定、漏えい禁止値は security owner の主本文を正とし、本ファイルへ重複定義しない。

### 27.43 API キー管理

本節の主本文は `docs/details/security.md` §27.43 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

API 側は、token API の route、HTTP method、request validation の入口、response schema、`.api_tokens` read/write 呼び出し境界だけを担当する。token 生成、hash 保存、scope 検証、作成時 1 回だけ token 本体を返す契約、認証成功時の `last_used_at` 更新、token 漏えい禁止は security owner の主本文を正とし、本ファイルへ重複定義しない。

### 27.44 監査ログ

本節の主本文は `docs/details/security.md` §27.44 を正とする。owner component は `security`、collaborator component は `api`、`statefile` とする。

API 側は、audit log API の route、query parameter、response schema、`.audit_log` read 呼び出し境界だけを担当する。audit record schema、action / actor / target / result、必須 audit 失敗時の `500`、secret / token / request body 保存禁止、壊れた行の扱いは security owner の主本文を正とし、本ファイルへ重複定義しない。

### 27.45 セッションタイムアウト変更設定

本節の主本文は `docs/details/security.md` §27.45 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

API 側は、session timeout config API の route、request body、response body、`.server_config.session_timeout_seconds` read/write 呼び出し境界だけを担当する。session の作成、期限判定、sliding update、期限切れ時 `401`、既存 session への反映条件、監査順序は security owner の主本文を正とし、本ファイルへ重複定義しない。

### 27.46 TOTP 二要素認証

本節の主本文は `docs/details/security.md` §27.46 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

API 側は、auth / TOTP API の route、request body、response body、`.totp_secret` read/write 呼び出し境界だけを担当する。TOTP secret 生成、setup 仮 secret、login ticket、code 検証、secret の一回表示、ticket 再利用禁止、TOTP 漏えい禁止、監査順序は security owner の主本文を正とし、本ファイルへ重複定義しない。

### 27.47 API レート制限

本節の主本文は `docs/details/security.md` §27.47 を正とする。owner component は `security`、collaborator component は `api`、`sdk`、`ui`、`statefile` とする。

API 側は、rate limit config API の route、request body、response body、`.server_config.api_rate_limit` と `.api_rate_state` の read/write 呼び出し境界だけを担当する。endpoint group 判定、window / count 更新、`429` 時に count を増やさない契約、actor key / IP key の同一 lock 更新、rate limit audit は security owner の主本文を正とし、本ファイルへ重複定義しない。
