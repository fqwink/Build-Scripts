# Adlaire CI — API 詳細仕様

[`docs/details/api.md`](api.md) は `api` owner component の詳細本文責務として、`api` が主本文として持つ実装契約だけを扱う。

owner / collaborator 境界管理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) に従う。`api` owner component の主本文であり、collaborator component の仕様は呼び出し境界、schema、security、表示、検証観点として参照する。fixture、expected、fake、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務を参照する。

api 連動機能の owner / collaborator は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.3](../DETAIL_INDEX.md#0i3-api--sdk--ui) を入口とする。endpoint、成功時・失敗時動作、状態副作用は [`docs/details/api.md`](api.md)、SDK method は [`docs/details/sdk.md`](sdk.md)、UI 操作は [`docs/details/ui.md`](ui.md)、状態 schema は [`docs/details/statefile.md`](statefile.md)、検証証跡は [`docs/details/fixture.md`](fixture.md) を正本とする。

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) は security 領域の詳細本文責務である。api が security 機能に関わる場合、[`docs/details/api.md`](api.md) 詳細本文責務は endpoint dispatch、request / response、状態ファイル read/write 呼び出し境界だけを担当し、scope、token、audit、session、TOTP、rate limit、漏えい禁止、security 横断順序の主本文は [`docs/details/security.md`](security.md) 詳細本文責務を参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `api` |
| 実装主体 | [`components/api.go`](../../components/api.go)。起動入口は [`main.go`](../../main.go)、実行バイナリ名は `adlaire-ci-api` とする。 |
| 持つ内容 | `api` owner が主本文として定義する HTTP 共通契約、endpoint、request / response、状態ファイル read/write 呼び出し境界、認証連携、api owner 追加機能。 |
| 持たない内容 | SDK method 実装、UI DOM 詳細、runner の build 実行責務、builder の変換処理、admin 静的配信、security 主本文、状態 schema、setup / update 手順、release 生成・公開手順、fixture 証跡責務。 |

---

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

**`api` 実行時設定：**

| 項目 | 入力元 | 固定契約 |
|------|--------|----------|
| listen address | `--addr` | 省略時 `127.0.0.1:8765`。[`docs/details/api.md` 詳細本文責務 §21a](api.md#21a-管理-api-サーバー制限) に従い外部公開 bind を使用しない。 |
| state root | `--state-dir` | 必須の絶対 path。API が参照するすべての runtime 状態 path は、この directory と [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル名から決定する。固定絶対 path を別途持たない。 |
| log level | `readServerConfig().log_level` | request 処理開始時の正規化済み値を使用する。不在時の既定値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema を正本とする。 |
| GitHub repository identity | `readRepoConfig()` | `owner` / `repo` だけを取得する。既定値、schema、破損時の扱いは [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) を正本とする。branch と監視対象は `readBranchConfig()` から取得する。 |

API service の systemd unit、配置、起動、更新、rollback は setup owner component の責務とし、[`docs/details/setup.md` 詳細本文責務 §26.3b](setup.md#sec-26-3b)、[`docs/details/setup.md` 詳細本文責務 §26.4.2](setup.md#sec-26-4-2)、[`docs/details/setup.md` 詳細本文責務 §26.5](setup.md#sec-26-5) を参照する。

`api` は `adlaire-ci-api --addr 127.0.0.1:8765 --state-dir /opt/adlaire-builder` として起動された後の HTTP listener、request / response、状態ファイル read/write 呼び出し境界だけを定義する。

## 21a. 管理 API サーバー制限

[`docs/details/api.md` 詳細本文責務 §21a](api.md#21a-管理-api-サーバー制限) は `api` の実行時制限を定義する。runner、setup、admin、sdk、ui は [`docs/details/api.md` 詳細本文責務 §21a](api.md#21a-管理-api-サーバー制限) の制限を上書きしてはならない。

| 制限 | 詳細 | 実装時の禁止条件 |
|------|------|------------------|
| session はインメモリ管理 | 再起動で全 session を消去する。永続 session store は持たない。 | session を状態ファイル、cookie store、外部 DB、外部 cache に保存しない。 |
| HTTPS listener 非対応 | `api` は HTTP listener のみ起動する。標準 bind は `127.0.0.1:8765` とする。 | TLS listener、証明書読み込み、HTTPS redirect、外部公開 bind を実装しない。 |
| 外部認証非対応 | 認証は `.admin_credentials`、`.totp_secret`、session、API token で完結する。 | SSO、OAuth、LDAP、SAML、複数ユーザー管理を追加しない。 |
| 独自接続数制限なし | Go 標準ライブラリ `net/http` の標準 server で処理する。API rate limit は [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の固定窓で行う。 | 独自 worker pool、connection pool、接続数上限、外部 queue を追加しない。 |
| runner 処理責務なし | api は HTTP endpoint、durable queue 投入、systemd への非同期起動要求、response だけを担当する。 | build id 採番、`running=true` 遷移、runner の polling loop、GitHub read、pipeline 実行、build log 確定処理を api で行わない。 |

セッション、API token、TOTP、rate limit、audit log の security 主本文は [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) を基準とし、[`docs/details/api.md` 詳細本文責務 §21a](api.md#21a-管理-api-サーバー制限) は API server の実行時境界だけを定義する。

---

## 22. バックエンド API 仕様

**ベース URL：** `http://localhost:{PORT}/api`
**認証：** `Authorization: Bearer {SESSION_TOKEN}`（`/api/login` で取得したセッショントークン）
**レスポンス形式：** JSON

<a id="sec-22-0"></a>
**22.0 API 共通契約：**

[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) は、`api` owner の API 共通契約を定義する。実装時は、エンドポイント固有契約より先に以下の共通契約を満たす。

| 項目 | 仕様 |
|------|------|
| 実装前提 | Go 最小バージョンは [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値)、HTTP 技術選定は [`docs/SPEC.md` 方針責務 §4](../SPEC.md#4-技術方針) を参照する。 |
| bind | 既定値は `127.0.0.1:8765`。`--addr <host:port>` が指定された場合は、起動中の listen address だけを置換する。`--addr 0.0.0.0:<port>` を指定しても、`api` は TLS listener、origin 制限、IP allowlist、reverse proxy 設定生成を追加実行しない。 |
| 文字コード | リクエストボディ、レスポンスボディ、状態ファイルは [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の文字コード契約を使用する。 |
| JSON レスポンス | JSON レスポンスには `Content-Type: application/json; charset=utf-8` を付与する。 |
| request ID | 全 `/api/` request の受付時に `crypto/rand` で 16 bytes を生成し、32 文字 lowercase hex として扱う。全 response の `X-Request-Id`、`.api_access_log.request_id`、同一 request で作成する `.audit_log.request_id` と `.config_log.request_id` は同じ値を使用する。生成失敗時は endpoint 処理、認証、状態更新、各 request log 追記を行わず `500 {"error":"Internal server error"}` を返し、`X-Request-Id` は付与しない。server log には secret や乱数値を含まない固定エラーを記録する。 |
| リクエスト body 上限 | JSON body は 1 MiB を上限とする。超過時は `413 Payload Too Large` と `{"error": "Payload too large"}` を返す。 |
| request body 禁止 | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) で `Request` が `none` の endpoint に body がある場合は `400 Bad Request` と `{"error": "Request body is not allowed"}` を返す。 |
| 成功レスポンス | 各エンドポイント例に記載した JSON オブジェクトを返す。空レスポンスは使用しない。 |
| エラーレスポンス | エラー時は `{"error":"<message>"}` を返す。入力検証失敗時のみ `details` を配列 `[{ "field": "<field>", "message": "<reason>" }]` とし、複数エラーがある場合はリクエスト JSON の出現順、query、path parameter の順で並べる。入力検証以外の補足は `details` を使わず、`error` を実装者向けではない固定文言にする。 |
| 未知のパス | 定義されていない `/api/...` は `404 Not Found` と `{"error": "Not found"}` を返す。 |
| 未対応メソッド | パスは存在するがメソッドが異なる場合は `405 Method Not Allowed` と `{"error": "Method not allowed"}` を返す。 |
| JSON 不正 | JSON ボディのパースに失敗した場合は `400 Bad Request` と `{"error": "Invalid JSON"}` を返す。 |
| 入力検証失敗 | 型、必須キー、範囲、有効値が仕様と異なる場合は `422 Unprocessable Entity` と `{"error":"Validation failed","details":[...]}` を返す。`field` は JSON body key、query key、または path parameter 名とし、body 全体の形式不正は `field` を `"$"` とする。 |
| 認証なし | 認証必須エンドポイントで Bearer トークンがない、または無効な場合は `401 Unauthorized` と `{"error": "Unauthorized"}` を返す。 |
| 権限不足 | 認証済み API token の scope が不足する場合は `403 Forbidden` と `{"error":"Forbidden"}` を返す。管理 session は全 API 操作を許可する。API token は `read`、`trigger`、`operate`、`config`、`admin` の scope だけを許可し、token 作成時に指定された scope 外の endpoint は拒否する。 |
| 競合 | 現在状態と要求操作が両立しない場合は `409 Conflict` を返す。対象は、ビルド未実行時の cancel、停止済みスケジュールへの pause、稼働中スケジュールへの resume、lock 取得 10 秒超過、stale 判定不能な `.build_lock` である。実行中の manual / force / webhook build request は queue 上限内なら waiting entry として受理し、上限到達時は `429 queue_full` とする。 |
| 未設定機能 | endpoint の必須 secret、必須外部設定、必須状態ファイルが未設定で処理を開始できない場合は `501 Not Implemented` と `{"error":"Not configured"}` を返す。エンドポイント固有仕様で `422`、`503`、`500` を明記している場合のみ個別指定を優先する。 |
| 時刻形式 | API レスポンスと状態ファイルの機械処理用時刻は UTC ISO 8601 `YYYY-MM-DDTHH:MM:SSZ` とする。明示オフセット、timezone なし文字列、ミリ秒付き文字列は保存しない。外部 API から取得した時刻も保存前に UTC `Z` へ正規化する。 |
| GET の副作用 | `GET` endpoint 固有処理は業務状態を変更しない。許可する共通副作用は、(1) 有効な管理 session の memory-only `last_used_at` 更新または期限切れ session の memory からの削除、(2) API token 認証成功時の `.api_tokens.last_used_at` 更新、(3) [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) が要求する `.api_rate_state` 更新、(4) [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2) が対象 event に要求する `.access_log` と `.audit_log` の追記、(5) 全 `/api/` request に対する `.api_access_log` の best-effort 追記の 5 群だけとする。endpoint 固有契約が method、対象、timeout を固定した外部 read または read-only command は観測処理として実行できるが、結果を業務状態へ保存しない。これら以外の file 作成・更新・削除、queue 変更、外部 write、変更 command、未定義 command、通知送信を禁止する。 |
| 状態ファイル更新 | JSON 状態ファイルの atomic write、lock、mode、file sync、親ディレクトリ sync は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に従う。永続化失敗時は `500` を返す。 |
| 秘密情報 | PAT、Webhook Secret、セッショントークン、API トークンはログ、バックアップ、GET レスポンスへ平文出力しない。設定済み表示は `"***"` または boolean で返す。 |
| 並列更新 | 同一状態ファイルを更新する API は、ファイル単位のロックを取得してから読み込み、検証、書き込みを行う。ロック取得待ちは最大 10 秒とし、超過時は `409 Conflict` を返す。 |
| 設定変更ログ | 設定変更 API は、変更前後の値を `.config_log` に追記する。ただし秘密情報の値は変更前後とも `"***"` にマスクする。監査 event は [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) を参照する。 |
| CORS | 既定では CORS ヘッダーを付与しない。標準管理ツールは同一 origin から配信する。`OPTIONS` preflight は定義しない。CORS を有効化する拡張は [`docs/details/api.md`](api.md) 詳細本文責務で未定義とし、実装してはならない。 |
| セキュリティヘッダー | すべての API レスポンスに `Cache-Control: no-store`、`X-Content-Type-Options: nosniff` を付与する。SSE は `Cache-Control: no-store` と `X-Accel-Buffering: no` を付与する。 |
| 判定順 | request ID 生成 → path 解決 → method 検証 → access control → [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2) の認証・scope・rate limit → body / query / path 入力検証 → endpoint 固有 gate → 状態競合 → 処理実行の順に判定する。`GET /api/health` だけは access control と maintenance を読まない。`POST /api/webhook` の body size 判定と署名検証、login endpoint の pre-auth rate limit と credential body parse は security 横断順序と各個別契約に従う。 |

エンドポイント例に記載されたフィールド名、型、有効値、HTTP ステータスは規範とする。API、SDK、標準管理ツールのいずれかを変更する場合は、[`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様)、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の対応関係を参照する。

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
| access control 破損 / 読取不能 | `503` | `{"error":"Access control unavailable"}` |
| maintenance 有効中 | `503` | `{"error":"maintenance"}` |
| maintenance 破損 / 読取不能 | `503` | `{"error":"maintenance_unavailable"}` |
| rate limit 超過 | `429` | `{"error":"Too many requests"}` |
| 入力検証失敗 | `422` | `{"error":"Validation failed","details":[...]}` |
| 状態競合 | `409` | endpoint 固有文言。未定義の場合は `{"error":"Conflict"}` |
| 必須設定なし | `501` | `{"error":"Not configured"}` |
| 内部処理失敗 | `500` | `{"error":"Internal server error"}` |

[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の固定表の body は空白差分を除いて固定とする。`500` の response body に Go error、path、secret、状態ファイル内容、外部 API response body を含めてはならない。内部原因は server log にだけ固定コード付きで出力する。

**API 実行順・副作用境界固定契約：**

全 endpoint は、[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の API 共通処理順序固定表の順序と副作用境界に従う。[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約で異なる順序を明記していない限り、実装者判断で検証順、状態読取、状態書込、外部呼び出し、ログ追記の順序を入れ替えてはならない。

| 段階 | 処理 | 失敗時 status | 失敗時副作用 |
|------|------|---------------|--------------|
| 0 | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の request ID を生成する。 | `500` | endpoint 処理、認証、状態更新、request 単位の log 追記を行わず、`X-Request-Id` を付与しない。 |
| 1 | path 解決。未知 path を判定する。 | `404` | 状態ファイル、外部 API、監査ログ、設定ログを変更しない。`.api_access_log` だけ [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) に従って記録する。 |
| 2 | method 検証。path が存在し method が不一致か判定する。 | `405` | endpoint 固有処理を開始しない。`.api_access_log` 以外を変更しない。 |
| 3 | `GET /api/health` 以外で `.access_control` を読み、接続元 IP を判定する。 | `403` / `503` | body 読取、認証、rate limit 更新、endpoint 状態読取、状態書込、外部呼び出しを行わない。 |
| 4 | [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2) の認証不要判定、pre-auth rate limit、認証、scope、認証後 rate limit を順番どおり実行する。 | `401` / `403` / `413` / `429` / `500` | [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2) の副作用境界に従う。scope 不足時は endpoint 固有 body を parse しない。Webhook は body size 判定と署名検証、login は login rate limit 後の credential body parse をこの段階で行う。 |
| 5 | body 禁止、body size、JSON parse、query、path parameter、body schema、enum、範囲を検証する。Webhook と login で段階 4 に実施済みの検証は再実行しない。 | `400` / `413` / `422` | endpoint 固有の業務状態を変更しない。段階 4 までに完了した security / observability 副作用は保持する。外部 API、systemd、runner、hook、通知を呼び出さない。 |
| 6 | endpoint 固有 gate と read adapter を呼び、maintenance、状態破損、状態競合を判定する。 | `409` / `500` / `503` | write lock を取得していても target を変更しない。tmp file があれば削除する。 |
| 7 | endpoint 固有処理を実行し、必要な状態ファイルを [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の Write 列順に更新する。 | endpoint 固有 | 途中失敗時の巻き戻しは、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約または [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) に明記された範囲だけ行う。 |
| 8 | endpoint 固有契約が response 確定前に要求する `.config_log`、`.audit_log`、`.access_log`、`.notify_log`、event log を仕様順に追記する。security 段階で追記済みの log は再追記しない。 | endpoint 固有 | 必須 log の失敗時挙動は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の該当 endpoint 固有契約に従う。既に確定済みの endpoint 状態は自動推測で再変更しない。 |
| 9 | response status、body、header を確定する。 | - | response 生成時に追加の状態読取、状態書込、外部呼び出しを行わない。 |
| 10 | 確定した response の status と request ID を用いて `.api_access_log` を [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) に従い追記する。 | 元の status を維持 | 追記失敗で response を `500` へ変更しない。server log に `API_ACCESS_LOG_WRITE_FAILED` を WARN で記録する。 |
| 11 | 確定済み response を送信する。 | - | 送信開始後に状態や log を追加更新しない。 |

`GET` endpoint は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の API 共通処理順序固定表の段階 7 で業務状態ファイルを書き換えない。`POST`、`DELETE` endpoint でも、段階 6 までに失敗した場合は endpoint 固有の状態書込を一切行わない。外部 API 送信、systemd 操作、hook 実行、通知送信、snapshot 操作は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約の処理順に現れる場合だけ実行する。実装者判断で「先に外部確認してから validation error を返す」処理にしてはならない。

<a id="sec-22-0a"></a>
**22.0a API から参照する状態ファイル共通仕様：**

状態ファイルのパス、形式、初期値、更新責務、破損時の扱い、更新手順、schema 厳格化、状態読取 adapter は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) を参照する。`api` は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の adapter と更新手順を利用し、endpoint 固有の request / response / validation は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) 以降を参照する。

<a id="sec-22-0b"></a>
**22.0b 入力検証共通仕様：**

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

共通検証値と [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約の値が異なる場合は、[`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の固定表を優先する。[`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の固定表にない query は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の共通検証を使用する。

| Endpoint | query | 既定値 | 許容値 | 補足 |
|----------|-------|--------|--------|------|
| `GET /api/history` | `page` | `1` | 1 以上 | 整数文字列だけ許可する。 |
| `GET /api/history` | `per_page` | `20` | 1〜100 | `0`、負数、小数、指数表記は禁止。 |
| `GET /api/history` | `failure_category` | `""` | `github_api`, `pipeline_timeout`, `pipeline_exit`, `deploy_failure`, `hook_error`, `config_error`, `resource_error`, `unknown`, 空文字 | 空文字は未指定。大文字小文字を区別し、未知値は `422`。 |
| `GET /api/logs` | `n` | `100` | 1〜1000 | `q` は空文字を許可する。 |
| `GET /api/logs/search` | `q` | `""` | 0〜500 文字 | 空文字は全件検索ではなく level/from/to のみ検索として扱う。 |
| `GET /api/logs/search` | `from`, `to` | `""` | 空文字または `YYYY-MM-DD` | `from > to` は `422`。 |
| `GET /api/logs/search` | `level` | `""` | `info`, `warn`, `warning`, `error`, `debug`, 空文字 | 大文字小文字は区別しない。 |
| `GET /api/access-log` | `limit` | `100` | 1〜1000 | 整数文字列だけ許可する。 |
| `GET /api/access-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/api-access-log` | `limit` | `100` | 1〜1000 | [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) を優先する。 |
| `GET /api/api-access-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/notify-log` | `limit` | `100` | 1〜1000 | 整数文字列だけ許可する。 |
| `GET /api/notify-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/config-log` | `limit` | `100` | 1〜1000 | 整数文字列だけ許可する。 |
| `GET /api/config-log` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/webhook-events` | `limit` | `50` | 1〜1000 | [`docs/details/api.md` 詳細本文責務 §27.13](api.md#sec-27-13) を優先する。 |
| `GET /api/webhook-events` | `offset` | `0` | 0 以上 | 整数文字列だけ許可する。 |
| `GET /api/audit-log` | `limit` | `100` | 1〜200 | [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) を優先する。 |
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
| 未知 query key | query 名 | `unknown field` |
| path parameter 不正 | parameter 名 | `invalid path parameter` |
| body 全体が object でない | `$` | `object required` |
| query が整数でない | query 名 | `integer required` |
| 日付の暦日不正 | query 名または key 名 | `invalid date` |

複数エラーがある場合、body key は JSON object の出現順、query は URL query の出現順、path parameter は route 定義順で並べる。body、query、path にまたがる場合は body → query → path の順とする。API は `details[].field` と `details[].message` の順序、文言、大小文字を変更しない。SDK の error 伝播は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI の field error 表示は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) を参照する。

<a id="sec-22-0c"></a>
**22.0c API から参照する主要状態ファイル schema：**

主要状態ファイル schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。`api` 詳細では API endpoint と状態ファイルの read / write 対応を [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) 以降で定義する。

<a id="sec-22-0c-1"></a>
**22.0c.1 API 状態読取優先順：**

[`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) は、API endpoint ごとの状態読取順、response 算出、不在時 response、破損時 response を定義する。状態ファイルの schema、adapter 戻り値、atomic write、破損退避は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) と [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。

<a id="api-state-read-priority-table"></a>
[API 状態読取優先順表](#api-state-read-priority-table) または [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の Read 列で `.server_config` を参照するすべての endpoint は、raw JSON を endpoint 内で直接 parse せず、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の `readServerConfig()` を 1 request につき 1 回呼び出し、その戻り値を request 中固定する。不在時は既定値、破損時は `500 {"error":"State file is corrupted"}`、読込不能時は `500 {"error":"State file read failed"}` とし、read-only endpoint は file を作成または修復しない。

| Endpoint | 読取順 | 正常時 response 算出 | 不在時 | 破損時 / 読込不能時 |
|----------|--------|----------------------|--------|---------------------|
| `GET /api/status` | `readBuildStatus()` と `readBuildLock()` → status 不在時だけ `readBuildHistory()`、`readBuildState()`、`readPendingTransfers()`、`readNotifyPending()`、`readCircuitState()` | `.build_status.json` がある場合は同ファイルだけから `last_sha`、`last_build_at`、`last_build_status`、`last_target_status`、`last_trigger`、`last_deploy_status`、両 pending 件数、`circuit_open` を写像する。`running` は同ファイルの値とするが、`readBuildLock().running=true` の場合は valid active lock と形式不正または PID 判定不能の conflict lock のどちらでも `true` へ上書きする。stale lock は上書きしない。正確な写像は [`StatusObject` 固定契約](#status-object-contract) を使用する。 | `.build_status.json` 不在時は最新の status summary 対象 history 行から `last_sha`、`last_build_at`、正規化 / 詳細 status、`last_trigger` を導出し、`last_deploy_status:null`、pending 配列件数、circuit state を返す。history-only 行は対象外とし、対象履歴なしは `last_sha:null`、`last_build_at:null`、`last_build_status:"none"`、`last_target_status:null`、`last_trigger:null`。`running` は `.build_state.running \|\| readBuildLock().running`。 | `.build_status.json` 破損は `500 {"error":"State file is corrupted"}`。fallback 中の必須読取破損も `500`。形式不正または PID 判定不能の lock は response を失敗させず `running:true`、stale lock は status または state の値を維持し、API は lock を変更しない。 |
| `GET /api/health` | `readBuildStatus()` →不在、破損、読込不能時だけ `readBuildHistory()`、常に `readPendingTransfers()`、`readNotifyPending()` | `.build_status.json` が有効な場合は `last_build_at=last_finished_at`、`last_build_status=status`、`last_deploy_at`、`last_deploy_status` を同名または対応 field から返す。[`docs/details/api.md` 詳細本文責務 §27.16](api.md#sec-27-16) の固定順で `checks` を作り、`checks` が空なら `status:"ok"`、1 件以上なら `status:"degraded"` とする。 | `.build_status.json` 不在は history fallback を行い、`last_deploy_at:null`、`last_deploy_status:null`、`build_status_missing` check を返す。 | 読取または parse 異常を対応する `checks` 値へ写像し、response 自体を構築できる場合は `200`。response 構築不能だけ `500`。 |
| `GET /api/history` | `readBuildHistory()` | schema-valid かつ id 重複除外済みの行を `finished_at` 降順、同時刻は `id` 降順に sort し、`trigger`、`tag`、`flagged`、`failure_category` filter 後に paging する。response の `warnings` は返却対象に未知の保存済み `failure_category` が 1 件以上ある場合だけ `["unknown_failure_category"]`、それ以外は `[]`。 | `total:0`、検証済み query の `page` と `per_page`、`pages:0`、`history:[]`、`warnings:[]`。 | 行単位破損は除外し、`BUILD_HISTORY_SKIP_CORRUPT`、id 重複は `BUILD_HISTORY_DUPLICATE_ID` を server log へ記録する。ファイル読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/history/{id}/log` | `readBuildLog(id)` | 対象 ID の log object を返す。 | 通常ログと archive の両方が不在なら `404 {"error":"Not found"}`。 | 対象 ID の log 破損は `500 {"error":"State file is corrupted"}`。 |
| `GET /api/logs` | `readLatestBuildLogs(n,q)` | 最新 log line を時系列順へ正規化し、`q` 指定時は部分一致で絞り込む。 | `{"lines":[]}`。 | 個別 log 破損は除外する。ディレクトリ読込不能は `500 {"error":"State file read failed"}`。 |
| `GET /api/queue` | `readBuildState()`、queue 上限算出時に `.server_config` | `active=.build_state.active_queue_entry`、priority 順の `queued`、`max_size=queue_max_size` を返す。 | `.build_state` 不在は `active:null`、`queued:[]`。`.server_config` 不在は既定 queue 上限。 | `.build_state` 破損は `500 {"error":"State file is corrupted"}`。 |
| `GET /api/repo-info`、`POST /api/repo-config` | request 中 1 回の `readRepoConfig()` | GET は `RepoInfoObject`、POST は正規化済み現在値へ `RepoConfigPatch` の指定 key だけを適用する。 | `{owner:"fqwink",repo:"Build-Scripts",updated_at:null}` を memory 上で使用し、GET は file を作成しない。POST は差分がある場合だけ非 `null` の `updated_at` を付与した 3 key object を保存する。 | 破損は `500 {"error":"State file is corrupted"}`、読込不能は `500 {"error":"State file read failed"}`。既定値で補完せず、既存 file、`.config_log`、`.audit_log` を変更しない。 |
| API 選択出力 target を使う endpoint | `readBranchConfig()` →不在時は default 正規化→ `readBuildStatus()` → endpoint 固有読取 | [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) で target を 1 件に固定してから endpoint 固有値を算出する。 | `.branch_config` 不在は default。`.build_status.json` 不在または target 不一致は正規化済み先頭 target。 | `.branch_config` 破損 / 読取不能は `500`。`.build_status.json` 破損 / 読取不能は target 選択では先頭 target へ fallback し、status 自体を応答する endpoint は個別契約に従う。 |
| `POST /api/build` / `POST /api/build/force` | `readMaintenance()` → `readCircuitState()` → `readBuildLock()` → `readBuildState()` | maintenance disabled、circuit closed、queue 容量内なら manual queue entry を保存し、systemd へ非同期起動を要求する。API は `running`、`current_build_id`、`last_started_at` を変更しない。 | `.maintenance` / `.build_state` / `.build_circuit_state` 不在は各初期値。 | maintenance 破損 / 読取不能は `503`、circuit / state 破損は `500`。circuit open は `409 {"error":"circuit_open"}`。`.build_lock` が形式不正または PID 判定不能の場合は `409 {"error":"Conflict"}`。queue 上限超過は `429`。 |
| `POST /api/build/cancel` | `readBuildLock()`、`readBuildState()` | 実行中 build を cancel request 状態へ更新する。 | lock 不在かつ running false は `409 {"error":"Conflict"}`。 | `.build_state` 破損は `500`。`.build_lock` 形式不正または PID 判定不能は `409`。 |
| `POST /api/circuit-breaker/reset` | `readCircuitState()` | 初期値と異なる場合だけ `open:false`、`consecutive_failures:0`、`opened_at:null`、`last_failure_at:null`、`last_error:null` を atomic write し、`.config_log` と `config_update` audit を順に追記する。既に初期値なら no-op とし、同じ response を返す。 | 不在は初期値とみなし、file を作成しない no-op。 | 読取破損は `500 {"error":"State file is corrupted"}` とし、上書きしない。主状態保存後の log 失敗は `500` とし、保存済み状態と先行 log を巻き戻さない。 |
| `GET /api/maintenance`、`POST /api/maintenance/enable`、`POST /api/maintenance/disable` | `readMaintenance()` | schema 準拠値を返却または更新する。 | `.maintenance` 不在は `enabled:false`、`reason:null`、`since:null` とし、GET では file を作成しない。 | 破損または読取不能は `500 {"error":"State file is corrupted"}` または `500 {"error":"State file read failed"}`。既存 file を上書きしない。 |
| `GET /api/health` 以外の全 endpoint | `readAccessControl()` | `allow` による接続元 IP 判定後だけ後続処理へ進む。 | `.access_control` 不在または `allow:[]` は全許可とし、file を作成しない。 | 破損または読取不能は `503 {"error":"Access control unavailable"}`。body 読取、認証、rate limit 更新、endpoint 固有処理を行わない。 |

<a id="sec-22-0d"></a>
**22.0d API と状態ファイル対応表：**

API 実装では、[`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の固定表の endpoint 固有 read/write と、[`docs/details/security.md` 詳細本文責務](security.md) が明記する認証、scope、rate limit、監査の共通 read/write 以外を実行してはならない。固定表の Read / Write 列は endpoint 固有処理を示し、API token 認証時の `.api_tokens` / `.access_log` / `.audit_log`、rate limit の `.api_rate_state`、[`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の監査対象 action は重複記載しない。複数ファイルを write する API は、表と参照先が定める順序で検証、バックアップ、書き込みを行い、途中失敗時は後続ファイルを書き込まない。

| API | Read | Write | 補足 |
|-----|------|-------|------|
| `POST /api/login` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.access_log`, `.audit_log` | TOTP 有効時は session を発行せず ticket を返す。 |
| `POST /api/login/totp` | `.admin_credentials`, `.totp_secret` | `.admin_credentials`, `.totp_secret`, `.access_log`, `.audit_log` | password verified ticket と TOTP code を検証し、成功時に session を発行する。 |
| `POST /api/logout` | メモリ上 session | メモリ上 session, `.access_log`, `.audit_log` | 対象 session を削除する。 |
| `POST /api/change-password` | `.admin_credentials` | `.admin_credentials`, メモリ上 session, `.access_log`, `.audit_log` | 現 session 以外をメモリから削除する。 |
| `GET /api/audit-log` | `.audit_log` | なし | 読取、並び順、filter、paging は [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) を正本とする。 |
| `GET /api/access-log` | `.access_log` | なし | [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の JSON Lines 一覧取得契約を使用する。 |
| `GET /api/api-access-log` | `.api_access_log` | なし | [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の JSON Lines 一覧取得契約と [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の filter 契約を使用する。 |
| `GET /api/sessions` | メモリ上 session | なし | token 本体は返さない。 |
| `POST /api/sessions/revoke-all` | メモリ上 session | メモリ上 session, `.access_log`, `.audit_log` | 現 session 以外を削除する。 |
| `GET /api/auth/totp-status` | `.totp_secret` | なし | 単一 admin の TOTP 状態を返す。 |
| `POST /api/auth/totp-setup` | `.totp_secret` | メモリ上 setup secret, `.audit_log` | 仮 secret と otpauth URI をメモリで生成し、未確認のまま永続化しない。 |
| `POST /api/auth/totp-confirm` | `.totp_secret` | `.totp_secret`, `.audit_log` | 仮 secret を code 検証後に保存し、TOTP を有効化する。 |
| `DELETE /api/auth/totp` | `.totp_secret` | `.totp_secret`, `.audit_log` | code 検証後に TOTP を無効化する。 |
| `GET /api/status` | `.build_status.json`, `.build_lock`, `.build_history`, `.build_state`, `.pending_transfers`, `.notify_pending`, `.build_circuit_state` | なし | `.build_status.json` を第一参照元とする。不在時のみ残りの状態から status 不在時集約値を算出する。 |
| `POST /api/build` / `POST /api/build/force` | `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state.queued`, `.audit_log` | idle / running にかかわらず manual queue entry を保存し、保存後に runner 起動を要求する。後者は `payload.force=true` とし、API は SHA cache を読み書きしない。 |
| `POST /api/build/cancel` | `.build_lock`, `.build_state` | `.build_state`, `.build_logs/{id}.json` | 実行中でない場合は `409`。 |
| `GET /api/build/stream` | `.build_logs/{id}.json`, `.build_logs/archive/`, `.build_state` | なし | 保存済み log 1 件の有限 SSE 配信だけを行い、通常 log と archive を更新しない。 |
| `GET /api/logs` | `.build_logs/`, `.build_logs/archive/` | なし | 通常 log を優先し、同一 build id の通常 log 不在時だけ archive を読む。 |
| `GET /api/logs/search` | `.build_logs/`, `.build_logs/archive/` | なし | 通常 log と archive を build id 重複なしで横断検索する。 |
| `GET /api/logs/export` | `.build_logs/`, `.build_logs/archive/` | なし | 通常 log 優先で JSON export する。 |
| `POST /api/logs/cleanup` | `.server_config`, `.build_logs/` | `.build_logs/` | 削除対象のみ削除する。 |
| `POST /api/logs/archive` | `.server_config`, `.build_logs/` | `.build_logs/archive/`, `.build_logs/` | archive 対象を gzip 圧縮し、成功後に通常 log を削除する。 |
| `GET /api/history` | `.build_history` | なし | `page`、`per_page`、`trigger`、`tag`、`flagged`、`failure_category` で絞り込み、ページングして返す。 |
| `GET /api/history/export` | `.build_history` | なし | 全件 export。 |
| `GET /api/history/{id}/log` | `.build_logs/{id}.json`, `.build_logs/archive/{id}.json.gz` | なし | 通常 log を優先し、不在時だけ archive を読む。両方不在は `404`、選択した log の破損は `500`。 |
| `GET /api/history/{id}/comment` | `.build_logs/{id}.json` | なし | なし。 |
| `POST /api/history/{id}/comment` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | コメントだけ更新する。 |
| `POST /api/history/{id}/flag` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.config_log` | flag だけ更新する。 |
| `POST /api/history/{id}/tags` | `.build_logs/{id}.json` | `.build_logs/{id}.json`, `.build_history`, `.config_log` | `.build_history` の同一 id にも反映する。 |
| `POST /api/history/{id}/rollback` | `.build_history`、`.snapshots/{id}/`、`.server_config`、`.branch_config`、`.maintenance`、`.build_circuit_state`、`.build_lock`、`.build_state` | `.audit_log`。runner owner が `.build_lock`、`.build_status.json`、`.build_state`、`.build_logs/{new_id}.json`、`.build_history`、必要時の `.pending_transfers` を更新する。 | maintenance、circuit、実行中判定後に runner owner の rollback coordinator を prepare し、`build_trigger` audit 成功後だけ worker 開始を確定する。API owner は rollback build log / history / pending を直接書き込まない。 |
| `GET /api/notify-config` | `.notify_config`, `.smtp_config` | なし | secrets はマスクする。 |
| `POST /api/notify-config` | `.notify_config` | `.notify_config`, `.config_log` | secret は GET で返さない。 |
| `GET /api/notify-log` | `.notify_log` | なし | [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の JSON Lines 一覧取得契約を使用する。 |
| `POST /api/notify-test` | `.notify_config` | `.notify_log` | 送信結果を追記する。 |
| `POST /api/notify/weekly-summary` | `.notify_config`, `.build_logs/` | `.notify_log` | 宛先なしは `422`。 |
| `GET /api/config` | `.server_config` | なし | 既定値を merge して返す。 |
| `POST /api/config/validate` | `.server_config`, request body | なし | 保存せず検証結果だけ返す。`.config_log` も更新しない。 |
| `POST /api/config` | `.server_config` | `.server_config`, `.config_log` | 許可キーのみ更新する。 |
| `POST /api/log-level` | `.server_config` | `.server_config`, `.config_log` | `log_level` のみ更新する短縮 API。 |
| `GET /api/config-log` | `.config_log` | なし | [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の JSON Lines 一覧取得契約を使用する。 |
| `GET /api/api-rate-limit` | `.server_config`, `.api_rate_state` | なし | rate limit 設定と現在 window summary を返す。 |
| `POST /api/api-rate-limit` | `.server_config` | `.server_config`, `.api_rate_state`, `.config_log`, `.audit_log` | policy を保存し、window state を初期化する。 |
| `GET /api/repo-info` | `.repo_config` | なし | `readRepoConfig()` の `owner`、`repo`、`updated_at` だけを返す。不在時は [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) の固定既定値を返し、file を作成しない。 |
| `POST /api/repo-config` | `.repo_config` | `.repo_config`, `.config_log` | `owner` / `repo` だけを部分更新する。branch / target は受け付けず、未指定の repository key は正規化済み現在値を保持する。 |
| `GET /api/branch-config` | `.branch_config` | なし | 不在時は default。 |
| `POST /api/branch-config` | `.branch_config` | `.branch_config`, `.config_log` | 空配列は `.branch_config` 削除。 |
| `GET /api/sysinfo` | `.branch_config`, `.build_status.json`, API 選択出力 target, process start time | なし | [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) で target を 1 件に固定し、状態ファイルは更新しない。 |
| `GET /api/health` | `.build_status.json`, `.build_history`, `.pending_transfers`, `.notify_pending` | なし | 認証不要。`.build_status.json` 不在、破損、または読込不能時に history fallback を行う。 |
| `GET /api/stats` | `.build_history`, `.build_logs/` | なし | `days` の範囲を集計する。 |
| `GET /api/stats/timeline` | `.build_history` | なし | 日別集計のみ。 |
| `GET /api/stats/build-duration` | `.build_logs/` | なし | duration 集計のみ。 |
| `GET /api/stats/build-trends` | `.build_trends.json` | なし | 最新 `n` 件を返し、返却対象 sample だけから response summary を再計算する。 |
| `GET /api/output-meta` | `.branch_config`, `.build_status.json`, `.build_history`, `.build_logs/`, `.build_logs/archive/`, API 選択出力 target | なし | 同じ target の出力サイト、成功履歴、対応 log だけを集約する。通常 log 不在時は同じ id の archive log を読む。 |
| `GET /api/pat-status` | `.github_token`, `.server_config` | なし | `.server_config.pat_expires_at` を参照する。結果保存なし。 |
| `POST /api/pat-verify` | `.github_token` | なし | 結果保存なし。 |
| `POST /api/pat-update` | なし | `.github_token`, `.config_log` | token 値は `.config_log` でマスクする。 |
| `GET /api/backup` | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config`, `.webhook_secret`, `.smtp_secret` | なし | secret 本体は返さず `*_set` boolean に変換する。 |
| `POST /api/restore` | request body | `.server_config`, `.notify_config`, `.repo_config`, `.branch_config`, `.access_control`, `.hooks`, `.alert_rules`, `.tag_rules`, `.pipeline_config`, `.dashboard_layout`, `.smtp_config`, `.webhook_secret`, `.smtp_secret`, `.config_log` | restore 対象ファイルを検証後、Write 列に記載された左から右の順で置換する。 |
| `GET /api/schedule` | `.server_config` | なし | systemd 次回実行時刻は外部確認。 |
| `POST /api/schedule/interval` | `.server_config` | `.server_config`, `.config_log` | systemd timer 反映も行う。 |
| `POST /api/schedule/pause` | `.server_config` | `.server_config`, `.config_log` | 既に paused は `409`。 |
| `POST /api/schedule/resume` | `.server_config` | `.server_config`, `.config_log` | 稼働中は `409`。 |
| `POST /api/schedule/allowed-hours` | `.server_config` | `.server_config`, `.config_log` | `null` で解除。 |
| `POST /api/schedule/force-interval` | `.server_config` | `.server_config`, `.config_log` | `hours` を保存。 |
| `POST /api/schedule/cooldown` | `.server_config` | `.server_config`, `.config_log` | `seconds` を保存。 |
| `GET /api/dashboard` | `.branch_config`, `.build_status.json`, `.build_history`, `.build_state`, `.build_lock`, `.pending_transfers`, `.notify_pending`, `.build_circuit_state`, `.server_config`, `.alert_rules`, `.dashboard_layout`, API 選択出力 target, process start time | なし | `status` は [StatusObject 固定契約](#status-object-contract) と完全一致させ、その後に target 選択後の集約値を組み立てる。 |
| `GET /api/diagnostics` | `.branch_config`, `.build_status.json`, `.github_token`, API 選択出力 target, systemd, `.notify_config`, `.webhook_secret` | なし | 診断結果は保存しない。 |
| `GET /api/rate-limit` | `.github_token` | なし | GitHub API 結果を返す。 |
| `GET /api/disk-usage` | `.branch_config`, `.build_status.json`, `.build_logs/`, `.build_logs/archive/`, API 選択出力 target | なし | log 全体と選択 target の出力サイトだけを集計する。 |
| `GET /api/webhook-events` | `.webhook_events.json` | なし | ページングして返す。 |
| `POST /api/webhook` | `.webhook_secret`, `.branch_config`, `.server_config`, `.build_state`, `.build_lock`, `.maintenance`, `.build_circuit_state` | `.build_state`, `.webhook_events.json` | 署名検証成功後のみイベント記録する。queue 追加時は queue 保存後にイベントを記録する。 |
| `GET /api/webhook-config` | `.webhook_secret` | なし | secret 本体は返さない。 |
| `POST /api/webhook-config` | なし | `.webhook_secret`, `.config_log` | secret 値は `.config_log` でマスクする。 |
| `POST /api/circuit-breaker/reset` | `.build_circuit_state` | `.build_circuit_state`, `.config_log` | 初期値へ戻す。冪等。 |
| `GET /api/snapshots` | `.snapshots/` | なし | 世代一覧を返す。 |
| `GET /api/snapshots/{id}/download` | `.snapshots/{id}/` | なし | バイナリを返す。 |
| `DELETE /api/snapshots/{id}` | `.snapshots/{id}/`、`.build_lock`、`.build_state` | `.snapshots/{id}/`、`.config_log`、`.audit_log`。runner owner が排他中だけ `.build_lock` を作成・削除する。 | runner owner の snapshot delete guard で idle を確定した後だけ archive owner の delete 実体を呼び出す。`deleted` で success config / audit、`delete_partial` で partial-failure config / failure audit、`delete_failed` で log 非追記とする。 |
| `GET /api/maintenance` | `.maintenance` | なし | 不在時は disabled。 |
| `POST /api/maintenance/enable` | `.maintenance` | `.maintenance`, `.config_log` | `since` を現在時刻で保存する。 |
| `POST /api/maintenance/disable` | `.maintenance` | `.maintenance`, `.config_log` | disabled 状態を保存する。 |
| `GET /api/access-control` | `.access_control` | なし | 不在時は `allow: []`。 |
| `POST /api/access-control` | `.access_control` | `.access_control`, `.config_log` | allow 全体を置換する。 |
| `GET /api/hooks` | `.hooks` | なし | hook 一覧を返す。 |
| `POST /api/hooks` | `.hooks` | `.hooks`, `.config_log` | hook id を新規採番する。 |
| `DELETE /api/hooks/{id}` | `.hooks` | `.hooks`, `.config_log` | 対象 hook のみ削除する。 |
| `GET /api/hooks/{id}/log` | `.build_logs/{build_id}_hook_{id}.json` | なし | `ran_at` 降順、同時刻は `build_id` 降順の先頭 20 件を返す。 |
| `GET /api/alert-rules` | `.alert_rules` | なし | rule 一覧を返す。 |
| `POST /api/alert-rules` | `.alert_rules` | `.alert_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/alert-rules/{id}` | `.alert_rules` | `.alert_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `GET /api/tag-rules` | `.tag_rules` | なし | rule 一覧を返す。 |
| `POST /api/tag-rules` | `.tag_rules` | `.tag_rules`, `.config_log` | rule id を新規採番する。 |
| `DELETE /api/tag-rules/{id}` | `.tag_rules` | `.tag_rules`, `.config_log` | 対象 rule のみ削除する。 |
| `POST /api/verify-output` | `.branch_config`, `.build_status.json`, `.build_history`, API 選択出力 target | なし | 同じ target の直近成功履歴と現在の出力 checksum 比較のみ。 |
| `GET /api/pipeline-config` | `.pipeline_config` | なし | 不在時は既定値。 |
| `POST /api/pipeline-config` | `.pipeline_config` | `.pipeline_config`, `.config_log` | config 全体を置換する。 |
| `GET /api/build-chain-config` | `.build_chain_config` | なし | 不在時は `{"chains":[]}`。 |
| `POST /api/build-chain-config` | `.build_chain_config` | `.build_chain_config`, `.config_log` | schema、参照整合、循環を検証後に config 全体を置換する。 |
| `GET /api/notes` | `.notes` | なし | 不在時は空文字。 |
| `POST /api/notes` | `.notes` | `.notes`, `.config_log` | content 全体を置換する。 |
| `GET /api/smtp-config` | `.smtp_config`, `.smtp_secret` | なし | password は返さず `password_set` だけ返す。 |
| `POST /api/smtp-config` | `.smtp_config`, `.smtp_secret` | `.smtp_config`, `.smtp_secret`, `.config_log` | password 指定時のみ `.smtp_secret` を更新する。 |
| `POST /api/smtp-test` | `.smtp_config`, `.smtp_secret` | `.notify_log` | 送信結果を記録する。 |
| `GET /api/queue` | `.build_state`, `.server_config` | なし | active entry と waiting queue を分離して返す。 |
| `DELETE /api/queue` | `.build_state` | `.build_state`, `.config_log` | waiting queue だけを削除し、active entry と実行中 build は変更しない。 |
| `GET /api/approvals` | `.approval_queue` | なし | id ごとの最新 record だけを返す。 |
| `POST /api/approvals/{id}/approve` | `.approval_queue`, `.server_config`, `.build_state`, `.build_lock` | `.build_state`, `.approval_queue`, `.audit_log` | pending entry だけを approval queue と build queue の固定更新順で承認する。 |
| `POST /api/approvals/{id}/reject` | `.approval_queue` | `.approval_queue`, `.build_history`, `.audit_log` | pending entry だけを却下し、rejected history を追記する。 |
| `GET /api/dashboard-layout` | `.dashboard_layout` | なし | 不在時は既定 widget 順。 |
| `POST /api/dashboard-layout` | `.dashboard_layout` | `.dashboard_layout`, `.config_log` | widgets 全体を置換する。 |
| `GET /api/tokens` | `.api_tokens` | なし | token 本体は返さない。 |
| `POST /api/tokens` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | token 本体は作成時のみ返し、保存はハッシュのみ。 |
| `DELETE /api/tokens/{id}` | `.api_tokens` | `.api_tokens`, `.access_log`, `.audit_log` | 対象 token を失効する。 |

<a id="sec-22-0e"></a>
**22.0e API 完全契約表：**

[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の固定表は、HTTP request / response / status、SDK method、UI 操作先の契約インデックスである。endpoint 固有の状態 read / write と処理補足は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) を唯一の正本とし、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) には再掲しない。endpoint を追加、削除、名称変更、body 変更、response 変更する場合は、[`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の状態アクセス表、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約、[`docs/details/sdk.md` SDK 引数変換契約](sdk.md#sdk-argument-contract)、[`docs/details/ui.md` UI 操作契約表](ui.md#ui-operation-contract)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) を同じ仕様変更範囲で先に更新する。[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の固定表に存在しない endpoint は実装対象外とする。SHA cache clear 専用 endpoint は定義しない。週次サマリーの手動送信は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `POST /api/notify/weekly-summary` だけを使用する。

`Request` が `none` の場合、request body を受け付けない。空 JSON object `{}` も送信してはならない。`Response` は成功時 body の schema 名または正確な top-level object を示す。schema 名の全 key は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約、状態由来の nested record の key は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を正本とする。`SDK` 列と `UI` 列は endpoint と collaborator の対応索引であり、response schema の正本ではない。SDK method の本文は [`docs/details/sdk.md`](sdk.md)、UI 操作の本文は [`docs/details/ui.md`](ui.md) を参照する。

`{message}` は `{"message": string}` を意味する。`{message,...}` 形式の response では `message` を必須キーとし、その他のキーも表記どおり必須とする。`?` が付いたキーだけを任意キーとする。成功時に空 body、`null` body、HTTP 204 は使用しない。

`Errors` 列は、[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の共通エラーに加算する endpoint 固有 status と認証境界を示す。共通処理で発生する `400`、`403`、`405`、`413`、`429`、`500`、`503` は、`Errors` 列に再掲していない場合も共通処理順序に従って適用する。`401` の記載は session または API token 認証を必須とすることを示す。`GET /api/health` は認証不要とし、`POST /api/webhook` の `401` は Webhook 署名検証失敗を示す。共通契約と endpoint 固有契約が同じ status を持つ場合、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約が条件と body を明記したときだけ endpoint 固有契約を優先する。

| Endpoint | Request | Response | Success | Errors | SDK | UI |
| ---------- | --------- | ---------- | --------- | -------- | ----- | ---- |
| `POST /api/login` | `{password}` | `{token?,must_change,totp_required?,ticket?}` | `200` | `401`, `422`, `429`, `500` | `login(password)` | ログイン |
| `POST /api/login/totp` | `{ticket,code}` | `{token,must_change}` | `200` | `401`, `422`, `429`, `500` | `loginTotp(ticket,code)` | ログイン |
| `POST /api/logout` | none | `{message}` | `200` | `401`, `500` | `logout()` | 全パネル共通 |
| `POST /api/change-password` | `{current_password,new_password}` | `{message}` | `200` | `401`, `422`, `500` | `changePassword()` | パスワード変更 |
| `GET /api/audit-log` | query `{limit,offset,actor?,action?,result?}` | `{log,total}` | `200` | `401`, `403`, `422`, `500` | `getAuditLog()` | 監査ログ |
| `GET /api/access-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `getAccessLog()` | アクセスログ |
| `GET /api/api-access-log` | query `{limit,offset,method?,path?,status?}` | `{log,total}` | `200` | `401`, `422`, `500` | `getApiAccessLog()` | アクセスログ |
| `GET /api/sessions` | none | `{sessions}` | `200` | `401` | `getSessions()` | セッション管理 |
| `POST /api/sessions/revoke-all` | none | `{message,revoked_count}` | `200` | `401`, `500` | `revokeAllSessions()` | セッション管理 |
| `GET /api/auth/totp-status` | none | `TotpStatus` | `200` | `401`, `500` | `getTotpStatus()` | セキュリティ |
| `POST /api/auth/totp-setup` | none | `{secret,otpauth_uri}` | `200` | `401`, `403`, `409`, `500` | `setupTotp()` | セキュリティ |
| `POST /api/auth/totp-confirm` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `confirmTotp(code)` | セキュリティ |
| `DELETE /api/auth/totp` | `{code}` | `TotpStatus` | `200` | `401`, `403`, `409`, `422`, `500` | `disableTotp(code)` | セキュリティ |
| `GET /api/status` | none | `StatusObject` | `200` | `401`, `500` | `getStatus()` | ステータス |
| `POST /api/build` | none | `{message,queue_id,queued,dispatch}` | `200`, `202` | `401`, `409`, `429`, `500`, `503` | `triggerBuild()` | 手動実行 |
| `POST /api/build/force` | none | `{message,queue_id,queued,dispatch}` | `200`, `202` | `401`, `409`, `429`, `500`, `503` | `buildForce()` | 手動実行 |
| `POST /api/build/cancel` | none | `{message}` | `200` | `401`, `409` | `cancelBuild()` | 手動実行 |
| `GET /api/build/stream` | none | SSE `log/end` events | `200` | `401`, `404`, `500` | `streamBuild()` | 手動実行 |
| `GET /api/logs` | query `{n,q}` | `{lines}` | `200` | `401`, `422`, `500` | `getLogs()` | ログビューア |
| `GET /api/logs/search` | query `{q,from,to,level}` | `SearchResult` | `200` | `401`, `422`, `500` | `searchLogs()` | ログビューア |
| `GET /api/logs/export` | none | `{exported_at,lines}` | `200` | `401`, `500` | `exportLogs()` | ログビューア |
| `POST /api/logs/cleanup` | none | `{message,deleted_count,failed_count}` | `200` | `401`, `500` | `cleanupLogs()` | ログビューア, 設定 |
| `POST /api/logs/archive` | none | `{message,archived_count}` | `200` | `401`, `500` | `archiveLogs()` | ログビューア, 設定 |
| `GET /api/history` | query `{page,per_page,trigger?,tag?,flagged?,failure_category?}` | `HistoryPageObject` | `200` | `401`, `422`, `500` | `getHistory()` | ビルド履歴 |
| `GET /api/history/export` | none | `ExportObject` | `200` | `401`, `500` | `exportHistory()` | ビルド履歴 |
| `GET /api/history/{id}/log` | path `{id}` | `HistoryLogObject` | `200` | `401`, `404`, `500` | `getHistoryLog(id)` | ビルド履歴 |
| `GET /api/history/{id}/comment` | path `{id}` | `CommentObject` | `200` | `401`, `404`, `500` | `getHistoryComment(id)` | ビルド履歴 |
| `POST /api/history/{id}/comment` | `{comment}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `setHistoryComment(id,comment)` | ビルド履歴 |
| `POST /api/history/{id}/flag` | `{flagged}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `setHistoryFlag(id,flagged)` | ビルド履歴 |
| `POST /api/history/{id}/tags` | `{tags}` | `{message}` | `200` | `401`, `404`, `422`, `500` | `setHistoryTags(id,tags)` | ビルド履歴 |
| `POST /api/history/{id}/rollback` | path `{id}` | `{message,build_id}` | `202` | `401`, `404`, `409`, `422`, `500`, `503` | `rollbackHistory(id)` | ビルド履歴, スナップショット |
| `GET /api/sysinfo` | none | `SysinfoObject` | `200` | `401`, `500` | `getSysinfo()` | システム情報 |
| `GET /api/health` | none | `HealthObject` | `200` | `500` | `health()` | なし（外部死活監視） |
| `GET /api/schedule` | none | `ScheduleObject` | `200` | `401`, `500` | `getSchedule()` | リポジトリ情報 |
| `POST /api/schedule/interval` | `{interval_seconds}` | `{message,interval_seconds}` | `200` | `401`, `409`, `422`, `500` | `setScheduleInterval(seconds)` | リポジトリ情報 |
| `POST /api/schedule/pause` | none | `{message}` | `200` | `401`, `409`, `500` | `pauseSchedule()` | リポジトリ情報 |
| `POST /api/schedule/resume` | none | `{message}` | `200` | `401`, `409`, `500` | `resumeSchedule()` | リポジトリ情報 |
| `POST /api/schedule/allowed-hours` | `{from,to}` | `{message,allowed_hours}` | `200` | `401`, `422`, `500` | `setAllowedHours()`, `clearAllowedHours()` | リポジトリ情報 |
| `POST /api/schedule/force-interval` | `{hours}` | `{message,hours}` | `200` | `401`, `422`, `500` | `setForceInterval(hours)` | リポジトリ情報 |
| `POST /api/schedule/cooldown` | `{seconds}` | `{message,seconds}` | `200` | `401`, `422`, `500` | `setBuildCooldown(seconds)` | リポジトリ情報 |
| `GET /api/notify-config` | none | `NotifyConfig` | `200` | `401`, `500` | `getNotifyConfig()` | 通知設定 |
| `POST /api/notify-config` | `NotifyConfig` | `{message}` | `200` | `401`, `422`, `500` | `setNotifyConfig(config)` | 通知設定 |
| `GET /api/notify-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `getNotifyLog()` | 通知設定 |
| `POST /api/notify-test` | none | `{message,channel_id}` | `200` | `401`, `422`, `500` | `notifyTest()` | 通知設定 |
| `POST /api/notify/weekly-summary` | none | `{message,period,success_count,failure_count,success_rate}` | `200` | `401`, `422`, `500` | `notifyWeeklySummary()` | 通知設定 |
| `GET /api/config` | none | `ConfigObject` | `200` | `401`, `500` | `getConfig()` | 設定 |
| `POST /api/config/validate` | partial `ConfigObject` | `ConfigValidationObject` | `200` | `401`, `422`, `500` | `validateConfig(config)` | 設定 |
| `POST /api/config` | partial `ConfigObject` | `{message,config}` | `200` | `401`, `422`, `500` | `setConfig(config)` | 設定 |
| `POST /api/log-level` | `{level}` | `{message,level}` | `200` | `401`, `422`, `500` | `setLogLevel(level)` | 設定 |
| `GET /api/config-log` | query `{limit,offset}` | `{log}` | `200` | `401`, `422` | `getConfigLog()` | 設定 |
| `GET /api/api-rate-limit` | none | `ApiRateLimitPolicyResponse` | `200` | `401`, `403`, `500` | `getApiRateLimit()` | セキュリティ |
| `POST /api/api-rate-limit` | `ApiRateLimitPolicyInput` | `ApiRateLimitPolicyResponse` | `200` | `401`, `403`, `422`, `500` | `setApiRateLimit(policy)` | セキュリティ |
| `GET /api/pat-status` | none | `PatStatusObject` | `200` | `401`, `500` | `getPatStatus()` | システム情報 |
| `POST /api/pat-verify` | none | `PatVerifyObject` | `200` | `401`, `501`, `500` | `patVerify()` | システム情報 |
| `POST /api/pat-update` | `{token}` | `{message}` | `200` | `401`, `422`, `500` | `updatePat(token)` | システム情報 |
| `GET /api/stats` | query `{days}` | `StatsObject` | `200` | `401`, `422` | `getStats(days)` | 統計 |
| `GET /api/stats/timeline` | query `{days}` | `TimelineObject` | `200` | `401`, `422` | `getStatsTimeline(days)` | 統計 |
| `GET /api/stats/build-duration` | query `{n}` | `BuildDurationStats` | `200` | `401`, `422` | `getStatsBuildDuration(n)` | 統計 |
| `GET /api/stats/build-trends` | query `{n}` | `BuildTrendStats` | `200` | `401`, `422`, `500` | `getBuildTrends(n)` | 統計 |
| `GET /api/output-meta` | none | `OutputMetaObject` | `200` | `401`, `404`, `500` | `getOutputMeta()` | システム情報 |
| `GET /api/repo-info` | none | `RepoInfoObject` | `200` | `401`, `500` | `getRepoInfo()` | リポジトリ情報 |
| `POST /api/repo-config` | `RepoConfigPatch` | `{message}` | `200` | `401`, `422`, `500` | `setRepoConfig(config)` | リポジトリ情報 |
| `GET /api/branch-config` | none | `{source,branches}` | `200` | `401`, `500` | `getBranchConfig()` | リポジトリ情報 |
| `POST /api/branch-config` | `{branches}` | `{message,branches_count}` | `200` | `401`, `422`, `500` | `setBranchConfig(branches)` | リポジトリ情報 |
| `GET /api/backup` | none | `BackupObject` | `200` | `401`, `500` | `backup()` | 設定 |
| `POST /api/restore` | `RestoreObject` | `{message}` | `200` | `401`, `422`, `500` | `restore(config)` | 設定 |
| `GET /api/dashboard` | none | `DashboardObject` | `200` | `401`, `500` | `getDashboard()` | ステータス, システム診断 |
| `GET /api/diagnostics` | none | `DiagnosticsObject` | `200` | `401`, `500` | `getDiagnostics()` | システム診断 |
| `GET /api/rate-limit` | none | `RateLimitObject` | `200` | `401`, `501`, `500` | `getRateLimit()` | システム情報 |
| `GET /api/disk-usage` | none | `DiskUsageObject` | `200` | `401`, `500` | `getDiskUsage()` | システム情報 |
| `GET /api/webhook-events` | query `{limit,offset}` | `{events,total}` | `200` | `401`, `422`, `500` | `getWebhookEvents(limit,offset)` | システム診断 |
| `POST /api/webhook` | GitHub webhook body | `{message,queued,event_id,queue_id?,dispatch?}` | `202` | `401`, `413`, `422`, `429`, `500`, `503` | none | 外部 Webhook |
| `GET /api/webhook-config` | none | `{configured}` | `200` | `401`, `500` | `getWebhookConfig()` | 通知設定 |
| `POST /api/webhook-config` | `{secret}` | `{message}` | `200` | `401`, `422`, `500` | `setWebhookConfig(secret)` | 通知設定 |
| `POST /api/circuit-breaker/reset` | none | `{message,open,consecutive_failures}` | `200` | `401`, `500` | `resetCircuitBreaker()` | 手動実行, システム診断 |
| `GET /api/snapshots` | none | `{snapshots}` | `200` | `401`, `500` | `getSnapshots()` | スナップショット |
| `GET /api/snapshots/{id}/download` | path `{id}` | binary | `200` | `401`, `404`, `422`, `500` | `downloadSnapshot(id)` | スナップショット |
| `DELETE /api/snapshots/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `409`, `422`, `500` | `deleteSnapshot(id)` | スナップショット |
| `GET /api/maintenance` | none | `MaintenanceObject` | `200` | `401`, `500` | `getMaintenance()` | メンテナンス |
| `POST /api/maintenance/enable` | `{reason}` | `{message,since}` | `200` | `401`, `422`, `500` | `enableMaintenance(reason)` | メンテナンス |
| `POST /api/maintenance/disable` | none | `{message}` | `200` | `401`, `500` | `disableMaintenance()` | メンテナンス |
| `GET /api/access-control` | none | `{allow}` | `200` | `401`, `500` | `getAccessControl()` | アクセス制御 |
| `POST /api/access-control` | `{allow}` | `{message,allow}` | `200` | `401`, `422`, `500` | `setAccessControl(allowList)` | アクセス制御 |
| `GET /api/hooks` | none | `{hooks}` | `200` | `401`, `500` | `getHooks()` | フック |
| `POST /api/hooks` | `{phase,command_args,abort_on_failure,timeout_seconds?}` | `HookRecord` | `201` | `401`, `409`, `422`, `500` | `addHook()` | フック |
| `DELETE /api/hooks/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `deleteHook(id)` | フック |
| `GET /api/hooks/{id}/log` | path `{id}` | `{id,runs}` | `200` | `401`, `404`, `500` | `getHookLog(id)` | フック |
| `GET /api/alert-rules` | none | `{rules}` | `200` | `401`, `500` | `getAlertRules()` | 設定 |
| `POST /api/alert-rules` | `{metric,operator,threshold,level,message}` | `AlertRule` | `201` | `401`, `409`, `422`, `500` | `addAlertRule()` | 設定 |
| `DELETE /api/alert-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `deleteAlertRule(id)` | 設定 |
| `GET /api/tag-rules` | none | `{rules}` | `200` | `401`, `500` | `getTagRules()` | 設定 |
| `POST /api/tag-rules` | `{condition,tags}` | `TagRule` | `201` | `401`, `409`, `422`, `500` | `addTagRule()` | 設定 |
| `DELETE /api/tag-rules/{id}` | path `{id}` | `{message}` | `200` | `401`, `404`, `500` | `deleteTagRule(id)` | 設定 |
| `POST /api/verify-output` | none | `{match,expected,actual}` | `200` | `401`, `404`, `500` | `verifyOutput()` | システム診断 |
| `GET /api/pipeline-config` | none | `PipelineConfig` | `200` | `401`, `500` | `getPipelineConfig()` | 設定 |
| `POST /api/pipeline-config` | `PipelineConfig` | `{message}` | `200` | `401`, `422`, `500` | `setPipelineConfig(config)` | 設定 |
| `GET /api/build-chain-config` | none | `BuildChainConfig` | `200` | `401`, `500` | `getBuildChainConfig()` | 設定, リポジトリ情報 |
| `POST /api/build-chain-config` | `BuildChainConfig` | `{message,chains_count}` | `200` | `401`, `422`, `500` | `setBuildChainConfig(chains)` | 設定, リポジトリ情報 |
| `GET /api/notes` | none | `{content,updated_at}` | `200` | `401`, `500` | `getNotes()` | 運用ノート |
| `POST /api/notes` | `{content}` | `{message,updated_at}` | `200` | `401`, `422`, `500` | `setNotes(content)` | 運用ノート |
| `GET /api/smtp-config` | none | `SmtpConfig` | `200` | `401`, `500` | `getSmtpConfig()` | 通知設定 |
| `POST /api/smtp-config` | partial `SmtpConfig` with optional password | `{message}` | `200` | `401`, `422`, `500` | `setSmtpConfig(config)` | 通知設定 |
| `POST /api/smtp-test` | none | `{result,message}` | `200` | `401`, `422`, `500` | `smtpTest()` | 通知設定 |
| `GET /api/queue` | none | `{active,queued,max_size}` | `200` | `401`, `500` | `getQueue()` | 手動実行 |
| `DELETE /api/queue` | none | `{message,cleared_count}` | `200` | `401`, `500` | `clearQueue()` | 手動実行 |
| `GET /api/approvals` | none | `{approvals}` | `200` | `401`, `500` | `getApprovals()` | 承認待ち |
| `POST /api/approvals/{id}/approve` | path `{id}`; body none | `{message,queued,queue_id,dispatch}` | `200` | `401`, `404`, `409`, `429`, `500` | `approveBuild(id)` | 承認待ち, 手動実行 |
| `POST /api/approvals/{id}/reject` | path `{id}`; body none | `{message}` | `200` | `401`, `404`, `409`, `500` | `rejectBuild(id)` | 承認待ち |
| `GET /api/dashboard-layout` | none | `{widgets}` | `200` | `401`, `500` | `getDashboardLayout()` | ステータス |
| `POST /api/dashboard-layout` | `{widgets}` | `{message}` | `200` | `401`, `422`, `500` | `setDashboardLayout(widgets)` | ステータス |
| `GET /api/tokens` | none | `{tokens}` | `200` | `401`, `403`, `500` | `getTokens()` | API トークン管理 |
| `POST /api/tokens` | `{label,scopes,expires_at?}` | `TokenCreateResult` | `201` | `401`, `403`, `422`, `500` | `createToken(label,scopes,expiresAt)` | API トークン管理 |
| `DELETE /api/tokens/{id}` | path `{id}` | `{message}` | `200` | `401`, `403`, `404`, `500` | `revokeToken(id)` | API トークン管理 |

**archive / snapshot / rollback 処理参照：**

`POST /api/logs/cleanup`、`POST /api/logs/archive`、`GET /api/snapshots`、`GET /api/snapshots/{id}/download`、`DELETE /api/snapshots/{id}` の archive 実体処理、download 安全性、rollback 用 artifact の検証・展開・転送は [`docs/details/archive.md` 詳細本文責務 §27.7](archive.md#sec-27-7) および [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を参照する。rollback build の lock、ID、log、history、pending、status、state は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。`api` 詳細では API endpoint、request / response、HTTP status、`.config_log` / `.audit_log` 境界だけを定義する。

**backup / restore API 副作用固定契約：**

| API | 処理順序 | 成功時副作用 | 失敗時副作用 |
|-----|----------|--------------|--------------|
| `GET /api/backup` | 対象設定 file 読込 → 不在 file に既定値適用 → secret mask → response。 | 状態ファイルを更新しない。 | 読込不能な必須 file は `500`。任意 file 不在は既定値で返す。 |
| `POST /api/restore` | request 検証 → secret mask `"***"` の既存 secret 再利用判定 → 全対象 payload 生成 → 定義済み Write 順に atomic write → `.config_log` 追記 → [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の `config_update` 追記 → response。 | 設定系状態 file だけを置換する。履歴、ログ、snapshot、session、token 本体は復元しない。 | 検証失敗は差分なし。途中 write 失敗は未処理 file を書かず `500`。処理済み file は戻さない。 |

backup response に secret 原文を含めてはならない。`password`、`token`、`secret`、`smtp_password`、`webhook_secret`、`.github_token`、`.smtp_secret`、`.api_tokens` の hash 元値は `"***"` または `*_set:boolean` で表現する。`POST /api/restore` で `"***"` を受け取った secret は既存値保持を意味し、既存値がない場合は未設定として扱う。

`BackupObject` は `exported_at`、`server_config`、`notify_config`、`repo_config`、`branch_config`、`access_control`、`hooks`、`alert_rules`、`tag_rules`、`pipeline_config`、`dashboard_layout`、`smtp_config`、`webhook_secret_set`、`smtp_password_set` の全 key を必須とし、未知 key を返さない。各設定 object の schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の同名状態ファイル schema を参照する。

`RestoreObject` は `server_config`、`notify_config`、`repo_config`、`branch_config`、`access_control`、`hooks`、`alert_rules`、`tag_rules`、`pipeline_config`、`dashboard_layout`、`smtp_config` を必須 key とし、`exported_at`、`webhook_secret_set`、`smtp_password_set`、`webhook_secret`、`smtp_password` だけを任意 key とする。このため `BackupObject` は変換なしで restore request として使用できる。`exported_at` は UTC ISO 8601 として検証するが永続化しない。`webhook_secret_set` と `smtp_password_set` は backup 時の参照値であり、restore の secret 書込を発生させない。`webhook_secret` と `smtp_password` は、key 省略で既存値保持、`"***"` で既存値保持、`null` で削除、それ以外の空でない string で新規保存とする。未知 key、空文字の secret、型不一致は `422`とし、どの file も変更しない。

**backup / restore 固定契約：**

| 項目 | 仕様 |
|------|------|
| backup 対象 | `.server_config`、`.notify_config`、`.repo_config`、`.branch_config`、`.access_control`、`.hooks`、`.alert_rules`、`.tag_rules`、`.pipeline_config`、`.dashboard_layout`、`.smtp_config`、`.webhook_secret`、`.smtp_secret`。 |
| backup 対象外 | `.admin_credentials`、`.api_tokens`、`.totp_secret`、session、`.build_history`、`.build_logs/`、`.snapshots/`、`.notify_log`、`.notify_pending`、`.webhook_events.json`、`.approval_queue`、`.build_state`。 |
| backup 不在値 | 任意設定 file 不在は schema 既定値で返す。secret file 不在は `*_set:false`。 |
| `.server_config` backup | 存在する場合は schema 検証済みの on-disk sparse object を返し、不在時は `{}` を返す。読取時に既定値の key を追加せず、破損または読込不能は `500`。 |
| backup secret | `.webhook_secret` と `.smtp_secret` は本体を返さず、`webhook_secret_set` / `smtp_password_set` boolean だけ返す。 |
| restore 検証 | すべての対象 payload を先に schema 検証し、1 件でも不正なら書込を開始せず `422`。 |
| restore `"***"` | 対応する既存 secret がある場合だけ既存値保持。既存 secret がない場合は未設定として扱い、新規 secret 文字列として保存しない。 |
| restore secret 削除 | secret key が `null` の場合は削除。key 省略は既存保持。 |
| restore 書込順 | `.server_config` → `.notify_config` → `.repo_config` → `.branch_config` → `.access_control` → `.hooks` → `.alert_rules` → `.tag_rules` → `.pipeline_config` → `.dashboard_layout` → `.smtp_config` → `.webhook_secret` → `.smtp_secret` → `.config_log` → `.audit_log`。 |
| restore 途中失敗 | 未処理 file は書かない。処理済み file は巻き戻さない。response は `500`。 |
| restore no-op | 全対象が既存値と同一の場合は file、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes" }`。 |
| `.server_config` restore | sparse object として検証し、正規化後の全 key 値で現在値と no-op 比較する。差分がある場合は request に含まれた top-level key だけを保って atomic write し、省略 key を追加しない。nested object は部分 merge しない。 |
| `.repo_config` backup / restore | backup は `readRepoConfig()` の `owner`、`repo`、`updated_at` を返す。restore で `{owner:"fqwink",repo:"Build-Scripts",updated_at:null}` を受けた場合は `.repo_config` 不在の既定状態を表すため、既存 file があれば lock 下で削除し、不在なら no-op とする。それ以外は `updated_at` が非 `null` の 3 key object だけを atomic write する。branch / target key は `422`。 |

restore の `.config_log` は対象 file ごとの差分を 1 record にまとめ、secret はすべて `"***"` とする。backup / restore の response、server log、fixture expected に secret 平文を含めてはならない。

**ビルド操作の競合優先順位：**

`POST /api/build`、`POST /api/build/force`、`POST /api/webhook`、`POST /api/history/{id}/rollback` は、以下の順に判定する。queue 対応 API と rollback API は、maintenance / circuit 判定後の実行境界を分離する。

1. [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の path、method、access control、[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2) の認証・scope・rate limit、body / query / path 入力検証をその順で完了する。Webhook は [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) の body size 判定と署名検証を認証の代替とする。
2. `.maintenance` が破損または読取不能の場合は `503 {"error":"maintenance_unavailable"}` を返す。
3. `.maintenance.enabled == true` の場合は `503 {"error":"maintenance"}` を返す。
4. `.build_circuit_state.open == true` の場合は `409 {"error":"circuit_open"}` を返す。ただし `POST /api/circuit-breaker/reset` は対象外。
5. `.build_lock` の形式が不正、PID 判定不能、または OS 判定失敗の場合は `409 {"error":"Conflict"}` を返し、queue と rollback を変更しない。
6. rollback は valid running lock または `.build_state.running == true` の場合に queue へ積まず、`409 {"error":"Build is running"}` を返す。idle の場合は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) の rollback coordinator を呼び出し、manual / webhook queue 処理を実行しない。
7. queue 対応 API（`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`）は idle / running にかかわらず `.build_state.queued[]` へ新しい queue id の entry を追加する。valid running lock、`.build_state.running == true`、または `.build_state.active_queue_entry != null` のいずれかが成立するときは `queue_max_size` 件、すべて成立しないときは `max(1, queue_max_size)` 件を waiting queue の保存上限とする。上限到達は `429 {"error":"queue_full"}` とする。
8. `POST /api/build` は `build_trigger`、`POST /api/build/force` は `build_force_trigger` を `.audit_log` へ追記する。`target_id` は queue id とする。audit 失敗時は保存済み queue entry を巻き戻さず `500` を返し、API からの runner 起動要求は行わない。
9. queue entry と endpoint 固有の必須永続化手順が完了した後、API は `systemctl start --no-block adlaire-ci.service` をシェルを介さず 5 秒 timeout で 1 回だけ実行する。Webhook で queue 追加後の event log 追記だけが失敗した場合は、保存済み queue entry を失わないための部分成功として起動要求へ進み、response に `event_log_failed:true` を含める。終了コード `0` は `dispatch:"requested"`、非 `0`、timeout、起動不能は `dispatch:"timer_fallback"` とする。
10. manual / force は `202 {"message":"Build queued","queue_id":"...","queued":true,"dispatch":"requested|timer_fallback"}` を返す。API は build id を採番せず、`running`、`current_build_id`、`last_started_at`、SHA cache を更新しない。

`dispatch:"requested"` は systemd が起動 job を受理したことだけを示し、runner の lock 取得または build 開始を意味しない。`dispatch:"timer_fallback"` の場合も durable queue entry は受理済みであり、API は server log に `RUNNER_ACTIVATION_DEFERRED: queue_id={id}` だけを記録する。次回 `adlaire-ci.timer` 起動が同じ entry を処理する。systemctl の stdout / stderr と OS error 詳細は HTTP response、audit、queue payload に含めない。

`POST /api/build/cancel` は queue を削除しない。実行中 build のみを cancel 対象とし、実行中でない場合は `409 {"error":"No build is running"}` を返す。`DELETE /api/queue` は実行中 build を停止せず、待機 queue のみ削除する。

<a id="sec-22-0e-1"></a>
**22.0e.1 API 機能別処理契約：**

集約 API の読取元と write 境界は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の endpoint 別状態アクセス表を正本とする。[`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の固定表は算出方法と空状態の戻り値だけを定義する。両表にない読取元や推測値を使用してレスポンスを補完してはならない。

`GET /api/sysinfo`、`GET /api/output-meta`、`GET /api/dashboard`、`GET /api/diagnostics`、`GET /api/disk-usage`、`POST /api/verify-output` の単数形の出力値は、以下の固定順で 1 件に決定した API 選択出力 target だけを対象とする。複数 target の集約値、任意の最新 directory、ファイル更新時刻の最大値から target を推測してはならない。

1. `readBranchConfig()` で `.branch_config` を読む。不在時は [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner) の `BRANCH_TARGETS` 既定値を同じ schema で正規化する。破損は `500 {"error":"State file is corrupted"}`、読取不能は `500 {"error":"State file read failed"}` とし、read-only request で退避、削除、default 書込みを行わない。
2. 正規化後の branch target が 0 件の場合は `500 {"error":"Output target unavailable"}` を返す。正常な `.branch_config` の空配列は default 復帰として保存ファイルを削除するため、保存済み空配列を有効状態として扱わない。
3. `readBuildStatus()` が有効で、`last_branch` と `last_target_file` が同じ 1 件の branch target の `branch` と `target_file` に完全一致する場合はその target を選ぶ。status 不在、破損、読取不能、nullable 値、または設定変更後の不一致では、branch 名の byte 昇順で先頭の target を選ぶ。target 選択だけを目的とする caller は `.build_status.json` を修復せず、status 自体を応答する endpoint の破損時処理は個別契約を優先する。
4. 選択値は `branch`、`target_file`、`out` の 3 値とし、1 request 内で固定する。応答生成中に `.branch_config` または `.build_status.json` を再読込して target を切り替えてはならない。
5. 選択 target の履歴は `.build_history` の `branch` と `target_file` が両方完全一致し、`status` が `success` または `success_deploy_pending` の行だけを対象とする。`finished_at` 降順、同時刻は `id` 降順で並べ、対応 log は同じ `id` の通常 log、archive log の順で読む。他 target の履歴、log、report、checksum を fallback に使用してはならない。

| 機能 | Endpoint | 算出方法 | 空状態 / 不足時 |
|------|----------|----------|-----------------|
| 現在状態 | `GET /api/status` | [StatusObject 固定契約](#status-object-contract) の 10 key だけを返す。`.build_status.json` がある場合は同 file と lock だけから算出し、不在時だけ履歴、state、pending、circuit へ fallback する。 | 履歴なしは `last_sha:null`, `last_build_at:null`, `last_build_status:"none"`, `last_target_status:null`, `last_trigger:null`, `last_deploy_status:null`, 両 pending 件数 `0`, `circuit_open:false`, `running:false`。 |
| 手動ビルド要求 | `POST /api/build` | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の競合優先順位に従い、manual queue entry の保存と runner 起動要求を行う。build 開始状態と SHA cache は変更しない。 | queue 上限到達は `429 {"error":"queue_full"}`。 |
| 強制ビルド要求 | `POST /api/build/force` | `payload.force=true` の manual queue entry を保存し、runner 起動を要求する。API は SHA cache を空値化せず、runner は SHA 一致 skip だけを bypass し、build 成功後に限り新 SHA を保存する。 | queue 上限到達は `429 {"error":"queue_full"}`。 |
| ログ一覧 | `GET /api/logs` | 同一 build id は通常 log を優先し、通常 log 不在時だけ archive を展開する。対象 build log の `pipeline.stdout`、`pipeline.stderr`、`warnings` をこの順で行配列へ連結し、build log の `finished_at` 昇順、同時刻は build id 昇順で並べた後に末尾 `n` 行を返す。`q` が空でない場合は大文字小文字を区別する部分一致行だけを絞り込んでから `n` を適用する。 | ログなしは `{"lines":[]}`。 |
| ログ検索 | `GET /api/logs/search` | [`docs/details/api.md` 詳細本文責務 §27.17](api.md#sec-27-17) の固定契約に従い、通常 log 優先で、`finished_at` 降順、build id 降順に検索し、一致 message を build 単位の `lines` へ固定 source 順でまとめる。 | 一致なしは `results:[]`。 |
| 履歴一覧 | `GET /api/history` | schema-valid かつ id 重複除外済みの行を `finished_at` 降順、同時刻は `id` 降順に並べ、`trigger`、`tag`、`flagged`、`failure_category` を指定時のみ完全一致で絞り込み、ページングする。壊れた行と重複 id 行は無視して固定 ERROR code を記録する。`total` と `pages` は絞り込み後の一意な有効行だけで算出する。`warnings` は返却 page の保存済み未知 category に限定する。 | 履歴なしまたは一致なしは `total:0`、検証済み query の `page` と `per_page`、`pages:0`, `history:[]`, `warnings:[]`。 |
| システム情報 | `GET /api/sysinfo` | API 選択出力 target の `out` 配下の通常 file 合計 bytes、directory mtime、process uptime を返す。 | 出力 directory 不在は `output_size_bytes:0`, `output_mtime:null`。 |
| 出力メタ | `GET /api/output-meta` | 選択 target の出力サイトの現在サイズと mtime、同じ target の直近成功履歴の `output_sha256`、対応 log の `report`、対応 log または HTML meta の `build_id` / `commit_sha` / `build_at` を返す。 | 出力サイト不在は `404`。同じ target の成功履歴がない場合は `sha256:""`、report 数値は `null`、warning は `[]`、build meta は HTML meta 不在時に空文字。 |
| ダッシュボード | `GET /api/dashboard` | `status`、API 選択出力 target の `sysinfo`、`stats(days=7)`、`schedule`、`alerts` を同一リクエスト時点で算出し、widget 順序は `.dashboard_layout.widgets` を使用する。 | `.dashboard_layout` 不在は既定 widget 順。出力 directory 不在の `sysinfo` は size `0` / mtime `null`。alerts なしは `[]`。 |
| 診断 | `GET /api/diagnostics` | PAT、GitHub API、API 選択出力 target の出力サイト、systemd、Webhook 設定を個別 item として返す。診断結果は保存しない。 | 出力 directory 不在は `output_file` item を `error` とする。各項目は `ok`、`warn`、`error` のいずれかを返す。 |
| ディスク使用量 | `GET /api/disk-usage` | 通常 build log、archive build log、API 選択出力 target の `out` 配下の通常 file を分類ごとに集計する。 | 出力 directory 不在は `output_file_bytes:0`。log directory 不在は対応 bytes / count を `0`。 |
| 出力検証 | `POST /api/verify-output` | 選択 target の直近成功履歴 `output_sha256` と、同 target の現在の出力サイト manifest SHA-256 を比較する。 | 同 target の成功履歴なし、`output_sha256:null`、または出力 directory 不在は `404 {"error":"Not found"}`。 |
| キュー | `GET /api/queue` | `.build_state.active_queue_entry` を `active`、`.build_state.queued` を priority 順の `queued`、`.server_config.queue_max_size` を `max_size` として返す。 | `.build_state` 不在は `active:null`, `queued:[]`。 |
| バックアップ | `GET /api/backup` | 設定状態だけを export し、secret 値は返さず `webhook_secret_set` / `smtp_password_set` boolean に変換する。履歴、ログ、snapshot、session は含めない。 | 不在の任意設定ファイルは初期値で返す。 |
| リストア | `POST /api/restore` | 対象 state schema をすべて検証してから [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の write 順に保存する。secret が `"***"` の場合は既存 secret を保持する。 | 検証失敗は書き込み前に `422`。途中失敗は未処理ファイルを書かない。 |

<a id="sec-22-0e-2"></a>
**22.0e.2 API ID 採番契約：**

API owner が新規 ID を生成する機能は、[`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の固定表の prefix と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。時刻ベース ID が suffix 上限へ到達した場合は `500 Internal Server Error` を返し、既存 ID を上書きしない。build id は runner owner が [`docs/details/runner.md` 詳細本文責務 build id 契約](runner.md#build-id-契約) で採番し、snapshot id は archive owner が build id と同一値を使用する。API owner がどちらも採番しない。

| 対象 | 形式 | 例 | 衝突時 |
|------|------|----|--------|
| queue id | `q{YYYYMMDDHHmmss}` | `q20260915100500` | `q20260915100500-001` |
| token id | `tok` + 6 桁連番 | `tok000001` | 既存最大番号 + 1。999999 超過時は `500` |
| hook id | `h{YYYYMMDDHHmmss}` | `h20260915100500` | `h20260915100500-001` |
| approval id | `appr{YYYYMMDDHHmmss}` | `appr20260915100500` | `appr20260915100500-001` |
| alert rule id | `r{YYYYMMDDHHmmss}` | `r20260915100500` | `r20260915100500-001` |
| tag rule id | `t{YYYYMMDDHHmmss}` | `t20260915100500` | `t20260915100500-001` |

<a id="sec-22-0e-3"></a>
**22.0e.3 API Endpoint 種別別実装契約：**

API handler は endpoint ごとの個別処理へ入る前に、[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の判定順を必ず適用する。共通判定後の endpoint 種別別処理は [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の固定表に従う。

| 種別 | 対象 endpoint | 処理順序 | 成功時 | 失敗時 |
|------|---------------|----------|--------|--------|
| read-only list | `GET /api/history`, `GET /api/access-log`, `GET /api/api-access-log`, `GET /api/notify-log`, `GET /api/config-log`, `GET /api/webhook-events` | query 検証 → 対象ファイル読込 → 壊れた行を除外 → filter → sort → paging → response 生成 | `total` がある endpoint は paging 前の除外・filter 後件数を返す。 | query 不正は `422`。ファイル読込不能は `500`。 |
| read-only aggregate | `GET /api/status`, `GET /api/dashboard`, `GET /api/health`, `GET /api/sysinfo`, `GET /api/output-meta`, `GET /api/disk-usage` | 必要ファイルを read-only で読込 → 不在時初期値適用 → 算出 → response 生成 | endpoint 固有の業務状態は変更しない。security / observability 共通副作用は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の GET 副作用契約に限定する。`.build_status.json` 破損時も API が自動修復しない。 | 必須ファイル破損は endpoint 固有の `500`。任意ファイル不在は初期値。 |
| read-only verification | `POST /api/verify-output` | API 選択出力 target 決定 →同 target の成功履歴読込 →現在 manifest 算出 → checksum 比較 → response 生成 | 出力、履歴、log、SHA cache、状態ファイルを変更しない。 | 出力または検証可能な成功履歴不在は `404`。読取不能は `500`。 |
| single-file update | `POST /api/config`, `POST /api/repo-config`, `POST /api/maintenance/*`, `POST /api/access-control`, `POST /api/dashboard-layout` | body 検証 → 現在値読込 → 差分生成 → atomic write → `.config_log` → `config_update` audit | 更新後値または `{message}` を返す。 | body 検証失敗は書込前に `422`。`.config_log` または audit 失敗時は対象更新済みのまま `500`。 |
| multi-file update | `POST /api/restore`, `POST /api/smtp-config`, `POST /api/history/{id}/tags` | 全入力検証 → 全対象読込 → 書込計画生成 → 定義済み Write 順に atomic write → `.config_log` → `config_update` audit | 全対象の更新完了後に response を返す。 | 検証失敗は書込なし `422`。途中失敗は未処理ファイルを書かず `500`。 |
| secret update | `POST /api/pat-update`, `POST /api/webhook-config`, `POST /api/smtp-config` password あり | secret 入力検証 → [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の atomic write と mode を適用 → `.config_log` へ `"***"` で記録 → `config_update` audit | secret 本体を response に含めない。 | secret 書込失敗は `500`。ログ、response、stdout へ平文を出さない。 |
| build dispatch | `POST /api/build`, `POST /api/build/force` | 共通判定 → maintenance → circuit → active / waiting 重複判定 → lock / queue 容量 → queue entry atomic write → trigger audit → `systemctl start --no-block adlaire-ci.service` | 新規 entry は `202` と queue id、`queued:true`、`dispatch` を返す。同一 active / waiting entry は再起動要求後に `200`。 | maintenance 破損 / 読取不能は `503 maintenance_unavailable`、enabled は `503 maintenance`、circuit open は `409 circuit_open`、lock 判定不能は `409 Conflict`、queue full は `429`。queue 書込みまたは audit 失敗は `500`。起動要求失敗は queue を残し `dispatch:"timer_fallback"`。 |
| rollback command | `POST /api/history/{id}/rollback` | 共通判定 → id 検証 → maintenance → circuit → runner rollback prepare → `build_trigger` audit → prepared worker 開始 | worker 開始境界が `accepted` を返した後だけ `202` と rollback build id を返す。 | 不正 id は `422`、history / snapshot 不在は `404`、maintenance は `503`、circuit open、running、現在の deploy target 不在は `409`、snapshot 破損、状態書込失敗は `500`。audit 失敗は prepared rollback を abort して `500`。manual queue entry と systemd runner 起動要求を作成しない。 |
| snapshot delete | `DELETE /api/snapshots/{id}` | 共通判定 → id 検証 → runner owner の snapshot delete guard 取得・state 再確認 → archive owner の事前検証・削除 → 三値結果別 log 処理 → guard 解放 → response | `deleted` かつ success config / audit と guard 解放の成功後だけ `{ "message":"Snapshot deleted" }` を返す。 | 不正 id は `422`、不在は `404`、running / lock conflict は `409`、state 破損・読取失敗、archive 破損、`delete_partial`、`delete_failed`、log 失敗、guard 解放失敗は `500`。public snapshot 削除後の失敗で snapshot と先行 log を巻き戻さない。guard 解放は guard 取得後の全終了経路で 1 回試行する。 |
| destructive delete | `DELETE /api/queue`, `DELETE /api/hooks/{id}`, `DELETE /api/alert-rules/{id}`, `DELETE /api/tag-rules/{id}`, `DELETE /api/tokens/{id}` | path / auth 検証 → 対象存在確認 → 削除または失効 → [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約に従い `.config_log` または `.audit_log` 追記 | `{message}` と件数がある場合は件数を返す。 | 対象不在は `404`。部分削除は禁止し、失敗時は `500`。 |
| external check | `POST /api/pat-verify`, `GET /api/rate-limit`, `GET /api/diagnostics`, `POST /api/smtp-test`, `POST /api/notify-test` | 設定読込 → timeout 付き外部確認 → [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約が要求する log 追記 → 結果 response | 確認結果そのものを状態へ保存しない。test 送信 log は response 確定前に仕様どおり追記する。 | 未設定は `501` または endpoint 固有 `422`。timeout は `500`。 |
| binary response | `GET /api/snapshots/{id}/download` | path 検証 → snapshot 存在確認 → archive owner による保存済み archive 事前検証 → header 確定 →検証に使用した同一 file descriptor から stream | 事前検証成功後だけ `Content-Type` と `Content-Disposition` を付与する。 | 不正 id は `422`、不在は `404`、事前検証または読込失敗は `500`。stream 開始後は JSON error を追加しない。 |
| stream response | `GET /api/build/stream` | 認証 → `.build_state` 読取 → 保存済み log を固定規則で 1 件選択 → SSE header → 有限 frame 送信 | `log` frame 後、必ず `end` frame を 1 回送って close する。wait、poll、tail、自動 reconnect は行わない。 | 選択可能 log 不在は `404`。選択対象または必須状態の読込不能・破損は `500`。送信中断時は状態ファイルを更新しない。 |

`.config_log`、`.access_log`、`.audit_log`、`.notify_log` への追記は JSON Lines 1 行単位で行う。追記失敗時は対象 endpoint の副作用が既に完了している場合でも、失敗を `500` として返し、次回 GET で破損行を無視できる形式を維持する。追記行の末尾改行を書けなかった場合は、その行を破損行として扱う。

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

<a id="api-json-lines-list-contract"></a>
**JSON Lines 一覧取得固定契約：**

| Endpoint | 対象 / 配列 key | 並び順 | 破損行 server log | 件数 key |
|----------|---------------------|----------|----------------------|----------|
| `GET /api/access-log` | `.access_log` / `log` | `at` 降順。同時刻は物理行順の逆順。 | `ACCESS_LOG_SKIP_CORRUPT` | なし |
| `GET /api/api-access-log` | `.api_access_log` / `log` | `at` 降順。同時刻は物理行順の逆順。 | `API_ACCESS_LOG_SKIP_CORRUPT` | `total` |
| `GET /api/notify-log` | `.notify_log` / `log` | `at` 降順。同時刻は物理行順の逆順。 | `NOTIFY_LOG_SKIP_CORRUPT` | なし |
| `GET /api/config-log` | `.config_log` / `log` | `at` 降順。同時刻は物理行順の逆順。 | `CONFIG_LOG_SKIP_CORRUPT` | なし |

各 endpoint は schema-valid 行だけを候補とし、endpoint 固有 filter を適用した後に固定順へ sort し、その後に `offset`、`limit` を適用する。`total` を持つ endpoint は paging 前の候補数を返す。対象ファイル不在は空配列とし、ファイル読込不能は `500 {"error":"State file read failed"}`とする。破損行は行番号と [JSON Lines 一覧取得固定契約](#api-json-lines-list-contract) の対象 endpoint 行に記載された破損行 server log code だけを server log へ 1 件記録し、行内容と secret 風値を記録しない。

<a id="sec-22-0e-4"></a>
**22.0e.4 API レスポンス正規化契約：**

API response は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の Response 列と対象 endpoint 契約、状態由来 record では [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema に一致させる。Response 列で `{...}` として示した key と対象 endpoint 契約で定義した schema key は、`?` または nullable 指定がある場合を除き必須である。SDK と UI は API response schema の正本ではなく、response を補完、再計算、別名化しない。endpoint ごとの正規化は以下とする。

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
| `GET /api/history` | `page`, `per_page` | `total`, `page`, `per_page`, `pages`, `history`, `warnings` | `pages = ceil(total / per_page)`。`total=0` の場合 `pages=0`。`warnings` は常に配列。 |
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
| `POST /api/repo-config` | 指定 key の正規化後値が既存値と一致 | `.repo_config`、`.config_log`、`.audit_log` を変更しない。 | `.repo_config` → `.config_log` → `.audit_log`。 |
| `POST /api/branch-config` | 正規化後 `branch_targets` が既存値と一致 | `.branch_config`、`.config_log`、`.audit_log` を変更しない。 | `.branch_config` 作成/置換/削除 → `.config_log` → `.audit_log`。 |
| `POST /api/notify-config` | secret mask 適用後の比較で既存値と一致 | `.notify_config`、`.config_log`、`.audit_log` を変更しない。 | `.notify_config` → `.config_log` → `.audit_log`。 |
| `POST /api/smtp-config` | config と password 更新有無が既存値と一致 | `.smtp_config`、`.smtp_secret`、`.config_log`、`.audit_log` を変更しない。 | `.smtp_config` → password 更新がある場合だけ `.smtp_secret` → `.config_log` → `.audit_log`。 |
| `POST /api/dashboard-layout` | widgets 配列が既存値と一致 | `.dashboard_layout`、`.config_log`、`.audit_log` を変更しない。 | `.dashboard_layout` → `.config_log` → `.audit_log`。 |

no-op response は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に固有の `No changes` が定義されている場合はその文言を返す。定義がない endpoint は通常成功文言を返す。no-op では、状態ファイル、JSON Lines、監査ログ、通知ログに差分を作ってはならない。部分更新では未指定 key を保持し、`null` が削除を意味する key は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に明記された key だけとする。

`.server_config` の no-op 比較と保存形式は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config` sparse object 契約を適用する。通常の設定更新 API は差分がある場合だけ全 top-level key を持つ正規化後 object を保存する。

**削除 / 失効 response：**

| Endpoint | 成功 response | 不在時 |
|----------|---------------|--------|
| `DELETE /api/snapshots/{id}` | `{message:"Snapshot deleted"}` | `404` |
| `DELETE /api/hooks/{id}` | `{message:"Hook deleted"}` | `404` |
| `DELETE /api/alert-rules/{id}` | `{message:"Alert rule deleted"}` | `404` |
| `DELETE /api/tag-rules/{id}` | `{message:"Tag rule deleted"}` | `404` |
| `DELETE /api/tokens/{id}` | `{message:"Token revoked"}` | `404` |
| `DELETE /api/queue` | `{message:"Queue cleared",cleared_count}` | queue 空でも `200`、`cleared_count=0` |

<a id="build-stream-sse-contract"></a>
**SSE frame 契約：**

`GET /api/build/stream` は request 開始時点で保存済みの build log 1 件を有限 SSE response として配信する。実行中 process の stdout / stderr を追尾する live tail ではない。`text/event-stream; charset=utf-8`、`Cache-Control: no-store`、`X-Accel-Buffering: no` を返す。各 frame は compact JSON 1 object を `data: {json}\n\n` として出力し、`event:`、`id:`、`retry:`、comment frame、複数 `data:` 行への分割を使用しない。

```jsonl
{"type":"log","line":"[INFO] build started","at":"2026-09-16T10:00:00Z"}
{"type":"end","status":"success","duration_seconds":12}
```

`type` は `"log"` または `"end"` だけとする。`log` frame の必須 key は `type`、`line`、`at`、`end` frame の必須 key は `type`、`status`、`duration_seconds` とし、他の key を追加しない。`log.line` は最大 4000 Unicode code point とし、JSON encode 前に先頭 4000 code point を残して末尾を切り捨てる。UTF-8 byte 列の途中で切断してはならない。`end.status="running"` は `end.duration_seconds:null` とだけ組み合わせる。`end.status` が `"success"`、`"failure"`、`"cancelled"` のいずれかの場合、`end.duration_seconds` は 0 以上の整数とする。`end` frame は接続終了前に 1 回だけ送る。送信中に client が切断した場合、状態ファイル、history、log を変更しない。

**SSE 対象 log 選択・送信固定契約：**

| 項目 | 仕様 |
|------|------|
| state 読取 | request 単位で `readBuildState()` を 1 回呼び、不在時は初期値を idle として使用する。破損は `500 {"error":"State file is corrupted"}`、読込不能は `500 {"error":"State file read failed"}` とし、log 選択を開始しない。SSE header は state と対象 log の読取成功後にだけ確定し、`404` / `500` は共通 JSON error response で返す。 |
| 実行中 | `.build_state.running=true` かつ `current_build_id` が非 `null` なら、その id の通常 log を選択する。途中保存 log の `status="running"`、`finished_at:null`、`duration_seconds:null` は有効とする。保存済み log がまだない場合は `404 {"error":"Not found"}` とし、過去 log へ fallback しない。 |
| 非実行中 | 通常 log と archive の schema-valid 完了 log から `finished_at` 降順、同時刻は build id 降順の先頭 1 件を選択する。同一 id は通常 log を優先する。 |
| 破損 log | 実行中に id で指定された log の破損・読込不能は `500`。非実行中の最新探索では破損 log を除外し、`BUILD_STREAM_SKIP_CORRUPT` と build id だけを server log に記録して次候補を選択する。有効候補がなければ `404`。 |
| log frame 順 | `pipeline.stdout` の各行 → `pipeline.stderr` の各行 → `warnings[]` の各要素 → `error` の非空文字の順。stdout / stderr は LF で分割し、空文字は 0 行、末尾 LF で生じる最後の空要素だけを除外し、途中の空行は `line:""` として保持する。warnings は配列の保存順を維持する。log 対象が 0 行でも `end` は送る。 |
| `log.at` | 選択 log の `finished_at` が非 `null` ならその値、実行中の途中保存 log だけは `started_at`。両方を欠く場合は schema 破損とし、現在時刻で補完しない。 |
| `end` | 選択 log の `status` と `duration_seconds` をそのまま使用する。実行中の途中保存 log は `status:"running"`、`duration_seconds:null` とする。全 `log` frame 後に 1 回だけ送る。 |
| 終了 | `end` 送信後に response body を close する。server 側の wait、poll、follow、keepalive と SDK 側の自動 reconnect を行わない。 |
| 副作用 | 通常 log、archive、`.build_state`、history、status を作成、更新、復元、削除しない。 |

**バイナリ response 契約：**

`GET /api/snapshots/{id}/download` は JSON error 以外では binary response とし、成功時に JSON body を返さない。archive owner の事前検証が成功した後だけ `Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"` を付与し、検証済みの保存済み `site.tar.gz` を配信する。`id` に `"`、`\`、改行を含む値は path 検証で `422` とする。

**レスポンス検証条件：**

| ケース | 期待結果 |
|--------|----------|
| schema 必須 key | 全 endpoint が必須 key を返す。 |
| 空配列 | 配列 key は `[]`、`null` や key 省略なし。 |
| 変更なし config | `No changes`、状態、`.config_log`、`.audit_log` に差分なし。 |
| delete queue 空 | `200`、`cleared_count=0`。 |
| SSE 正常終了 | `end` frame が 1 回だけ送信される。 |
| snapshot download | binary body、固定 header、JSON success body なし。 |

**api / sdk / ui / statefile 横断契約参照：**

API endpoint、SDK method、UI 操作、状態ファイル副作用の本文は各 owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務を参照する。owner / collaborator の対応は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.3](../DETAIL_INDEX.md#0i3-api--sdk--ui)、横断検証証跡は [`docs/details/fixture.md`](fixture.md) を参照する。

**横断 fixture 参照：**

api / sdk / ui / statefile にまたがる横断 fixture の fixture 名、入力、必須確認は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) API / SDK / UI / 状態ファイル cross fixture 固定契約を参照する。

<a id="sec-22-0f"></a>
**22.0f API fixture 参照：**

API の必須検証、fixture 名、入力状態、期待 response、期待副作用は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) を参照する。個別 API 節では fixture 名、fixture manifest、testdata 配置、期待副作用、実装検証証跡を再定義しない。

`api` 詳細では、API endpoint の method、path、request、response、error、read / write 境界だけを定義する。fixture manifest、testdata 配置、期待副作用、実装検証証跡は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務、実装計画上の割当と現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務、実装可否と完了判定は [`docs/SPEC.md` ポリシー責務 §0a](../SPEC.md#0a-仕様成熟度ポリシー) を参照する。

endpoint の method、path、認証境界、request、response、error、read / write 境界の一覧は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の API 完全契約表を唯一の正本とする。以降は endpoint ごとの追加詳細だけを記載する。

**`POST /api/login` リクエスト / レスポンス：**
```text
// リクエスト
{ "password": "admin" }

// レスポンス（TOTP 無効）
{ "token": "<session_token>", "must_change": "prompt" }

// レスポンス（TOTP 有効）
{ "totp_required": true, "ticket": "<login_ticket>", "must_change": "none" }
```

`must_change` の有効値：`"none"`（変更不要）| `"prompt"`（促す：初回ログイン時）| `"forced"`（強制：5 回目以降。変更完了まで管理画面の操作を制限）

**`POST /api/login/totp` リクエスト / レスポンス：**
```text
// リクエスト
{ "ticket": "<login_ticket>", "code": "123456" }

// レスポンス
{ "token": "<session_token>", "must_change": "none" }
```

**`POST /api/logout` リクエスト / レスポンス：**
```text
// リクエスト: なし（Bearer トークンのみ）
// レスポンス: 200
{ "message": "Logged out" }
```

**`POST /api/change-password` リクエスト / レスポンス：**
```text
// リクエスト
{ "current_password": "...", "new_password": "..." }
// レスポンス: 200
{ "message": "Password changed" }
```

**`GET /api/status` レスポンス例：**
```json
{
  "last_sha": "0123456789abcdef0123456789abcdef01234567",
  "last_build_at": "2026-09-14T10:00:00Z",
  "last_build_status": "success",
  "last_target_status": "success",
  "last_trigger": "polling",
  "last_deploy_status": "success",
  "pending_transfers_count": 0,
  "notify_pending_count": 0,
  "circuit_open": false,
  "running": false
}
```

<a id="status-object-contract"></a>
**`StatusObject` 固定契約：**

`StatusObject` は以下の 10 key をすべて必須とし、未知 key、`output_url`、`queued`、`current_build_id` を含めない。queue 状態は `GET /api/queue`、出力成果物情報は `GET /api/output-meta` をそれぞれ唯一の API 正本とする。

| key | 型 | `.build_status.json` あり | `.build_status.json` 不在時 fallback |
|-----|----|----------------------------|---------------------------------------|
| `last_sha` | string/null | `last_blob_sha` が非 `null` ならその値、それ以外は `last_commit_sha`。 | 最新 status summary 対象 history の `blob_sha` が非 `null` ならその値、それ以外は `commit_sha`。履歴なしは `null`。 |
| `last_build_at` | string/null | `last_finished_at`。 | 最新 status summary 対象 history の `finished_at`。履歴なしは `null`。 |
| `last_build_status` | string | `status` をそのまま返す。 | history の詳細結果を [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の runner 結果値写像で正規化する。履歴なしは `"none"`。 |
| `last_target_status` | string/null | `last_target_status`。 | 最新 status summary 対象 history の `status`。履歴なしは `null`。 |
| `last_trigger` | string/null | `last_trigger`。 | 最新 status summary 対象 history の `trigger`。履歴なしは `null`。 |
| `last_deploy_status` | string/null | `last_deploy_status`。 | history schema に同値がないため `null`。推測しない。 |
| `pending_transfers_count` | integer | `pending_transfers_count`。 | `readPendingTransfers()` の配列件数。 |
| `notify_pending_count` | integer | `notify_pending_count`。 | `readNotifyPending()` の配列件数。 |
| `circuit_open` | boolean | `circuit_open`。 | `readCircuitState().open`。 |
| `running` | boolean | `running`。ただし `readBuildLock().running=true` なら `true`。 | `readBuildState().running \|\| readBuildLock().running`。 |

`last_build_status` は [`docs/details/statefile.md` 詳細本文責務 `.build_status.json` schema](statefile.md#build-status-schema) の正規化値、`last_target_status` は同 schema の詳細結果値とする。`last_trigger` は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の trigger 有効値または `null`、`last_deploy_status` は同じ statefile schema の許容値または `null` とする。`last_sha` は 40 または 64 文字 lowercase hex または `null`、pending 件数は 0 以上の整数とする。

`running` の読取順、stale lock、形式不正または PID 判定不能の lock の扱いは、[`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) の `GET /api/status` 行を正本とする。

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

**`GET /api/history` 空履歴レスポンス例：**
```json
{
  "total": 0, "page": 1, "per_page": 20, "pages": 0,
  "warnings": [],
  "history": []
}
```

`page` は 1 始まり。`per_page` の最大値は 100。`trigger` は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の固定値のみ許可する。`tag` はタグ検証と同じ文字列制約を適用する。`flagged` は `"true"` または `"false"` のみ許可する。`failure_category` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の FailureCategory 固定値だけを許可し、指定時は保存値の完全一致で絞り込む。`.build_history.status="success"` の純粋な成功履歴は `failure_category:null` のため category filter 指定時に一致しない。`.build_history.status="success_deploy_pending"` の履歴は正規化状態が `success` でも `failure_category:"deploy_failure"` を保存し、`failure_category=deploy_failure` に一致する。未知 query、1 未満の `page`、範囲外の `per_page`、不正な `trigger` / `tag` / `flagged` / `failure_category` は `422 {"error":"Validation failed","details":[...]}` とし、`details[]` は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の固定規則に従って実際に不正な query 名と理由を返す。例えば未知の `failure_category` は `{"field":"failure_category","message":"invalid value"}`、整数でない `page` は `{"field":"page","message":"integer required"}` とする。有効な正整数だが最終ページを超える `page` は `422` にせず、その `page` と `per_page`、絞り込み後の `total` と `pages`、`history:[]`、`warnings:[]` を返す。各 history object は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の有効な `.build_history` record の全必須 key を保持し、response 専用 `build_at=finished_at` と `sha=commit_sha ?? blob_sha` を追加する。保存ファイルへ `build_at` と `sha` を逆書きしてはならない。

`HistoryPageObject.warnings` は常に string array とする。filter と paging 後の返却対象に、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の前方互換 category id が 1 件以上ある場合だけ `"unknown_failure_category"` を 1 回含める。前方互換 category id を書き換えず、当該 history object の原値を返す。前方互換 category id が paging 対象外の行だけにある場合は `warnings:[]` とする。この top-level `warnings` は個々の history object にある integer の `warnings` と別の field であり、互いに上書き、合算、型変換しない。

history `id` と対応 log の有無は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_history` schema を正とする。build record は build id を使い同名 build log を持つ。approval event と dependency skip は正本で定義された id を使い、対応 build log を持たない。

**`POST /api/build` / `POST /api/build/force` レスポンス例：**
```json
{ "message": "Build queued", "queue_id": "q20260915100500", "queued": true, "dispatch": "requested" }
```

systemd への非同期起動要求が失敗した場合は、同じ `202` で `dispatch:"timer_fallback"` を返す。`dispatch` の許容値は `"requested"` と `"timer_fallback"` だけとする。manual / force response に `build_id`、`queued:false`、`Build started` を返してはならない。build id と実行中状態は runner が queue entry を取得した後に確定する。

SHA キャッシュのクリアだけを行う専用 API は定義しない。強制再ビルドは必ず `POST /api/build/force` を使用する。API は SHA cache をクリアせず、force 実行要求を通常 build 要求と区別して runner へ渡す。runner は SHA 一致による skip だけを無効化し、成功後に限り SHA cache を新 SHA で更新する。

**`GET /api/sysinfo` レスポンス例：**
```json
{
  "output_size_bytes": 2048576,
  "output_mtime": "2026-09-15T10:00:00Z",
  "uptime_seconds": 86400
}
```

`output_size_bytes` と `output_mtime` は [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target だけから算出する。選択 target の `out` directory が不在の場合は HTTP `200`、`output_size_bytes:0`、`output_mtime:null` とし、directory の作成や他 target への fallback を行わない。

**`GET /api/schedule` レスポンス例：**
```json
{ "next_run_at": "2026-09-15T10:05:00Z", "interval": "5min", "paused": false, "allowed_hours": { "from": 9, "to": 18 } }
```

`paused` が `true` のとき、ポーリングは停止中で `next_run_at` は `null` を返す。

`allowed_hours`：自動ビルドを許可する時間帯（UTC）。`null` = 無制限。`from` 以上 `to` 未満の時刻のみビルドを実行する。許可時間帯外のポーリングでは変更を検出しても実行を保留し、次の許可時間帯に入った時点で実行する。

**`POST /api/schedule/allowed-hours` リクエスト / レスポンス：**
```text
// 設定
{ "from": 9, "to": 18 }
// 解除（無制限に戻す）
{ "from": null, "to": null }
// レスポンス: 200
{ "message": "Allowed hours updated", "allowed_hours": { "from": 9, "to": 18 } }
```

**`POST /api/schedule/pause` / `POST /api/schedule/resume` レスポンス例：**
```jsonl
{ "message": "Schedule paused" }
{ "message": "Schedule resumed" }
```

既に一時停止中に `pause`、または稼働中に `resume` を呼び出した場合は `409 Conflict` を返す。

> **責務分担：** build lifecycle 通知、pending retry、自動 weekly summary の送信責務は [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) および [`docs/details/runner.md` 詳細本文責務 §27.19](runner.md#sec-27-19) を参照する。`api` 詳細では通知 API の request / response、設定 read/write、履歴参照、手動送信 endpoint 境界だけを定義する。

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

`secret`：Webhook 署名シークレット。未設定時は `null`、設定済み時は `"***"`（マスク）を返す。

`on` の有効値：`"start"`（ビルド開始時）| `"success"`（ビルド成功時）| `"failure"`（ビルド失敗時）| `"deploy_failure"`（転送失敗時）| `"weekly_summary"`（定期サマリー送信時）| `"approval_required"`（承認待ち発生時）| `"duration_anomaly"`（所要時間異常時）| `"config_corrupt"`（設定破損復旧時）。複数指定は array 順を保持して保存する。

`summary`：週次サマリー通知の設定。`interval` の有効値は `"weekly"` 固定。`hour` は 0〜23（UTC）。`day_of_week` は 0 = 日曜〜6 = 土曜。自動送信条件、二重送信防止、集計、送信順序、`.build_state` 更新は [`docs/details/runner.md` 詳細本文責務 §27.19](runner.md#sec-27-19) を参照する。

**`POST /api/notify/weekly-summary` レスポンス例：**
```json
{ "message": "Weekly summary sent", "period": "2026-09-08/2026-09-14", "success_count": 12, "failure_count": 1, "success_rate": 92.3 }
```

`POST /api/notify/weekly-summary` は手動送信 API とする。集計対象、対象 channel 抽出、通知送信、`.notify_log` 追記は [`docs/details/runner.md` 詳細本文責務 §27.19](runner.md#sec-27-19) の共通集計・送信契約を使用する。HTTP response、status、validation details と、手動送信で sent date を更新しない境界は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/notify/weekly-summary` レスポンス参照を正本とする。

`payload_template`：Webhook 送信 JSON ペイロードのテンプレート文字列。`null` = デフォルトペイロードを使用。build lifecycle 通知 payload、channel 選択、retry、pending 保存、secret mask は [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) を参照する。テンプレート内で使用可能な変数は以下の通り。

| 変数 | 内容 |
|------|------|
| `{{id}}` | ビルド ID |
| `{{status}}` | ビルド結果（`success` / `failure`） |
| `{{sha}}` | 対象コミット SHA |
| `{{duration_seconds}}` | ビルド所要時間（秒） |
| `{{build_at}}` | ビルド実行日時（ISO 8601） |

`secret`：Webhook 送信時の HMAC-SHA256 署名用シークレット文字列。設定時はリクエストヘッダーに `X-Adlaire-Signature: sha256=<hmac>` を付与する。`null` = 署名なし。`GET /api/notify-config` で返却する際、設定済みの場合は `"***"` でマスクし、未設定の場合は `null` を返す。`POST /api/notify-config` で更新可能。

**`GET /api/config` レスポンス例：**
```json
{
  "log_max_lines": 500,
  "history_max_count": 100,
  "build_timeout_seconds": 300,
  "log_retention_days": 30,
  "log_archive_after_days": 0,
  "log_level": "INFO",
  "pat_expires_at": null,
  "snapshots_keep": 5,
  "queue_max_size": 3,
  "build_retry_max": 0,
  "build_retry_base_seconds": 5,
  "commit_status_enabled": false,
  "commit_status_context": "Adlaire CI",
  "commit_status_target_url": null,
  "build_trend_keep_count": 1000,
  "duration_anomaly": { "enabled": false, "min_samples": 20, "avg_multiplier": 2.0, "p95_multiplier": 1.5 },
  "session_timeout_seconds": 28800,
  "watch_mode": "github",
  "tag_filter": { "enabled": false, "patterns": [] },
  "build_cache_enabled": false,
  "deploy_parallelism": 1,
  "remote_build": { "enabled": false, "host": null, "user": null, "work_dir": null, "command_args": [], "artifact_path": null },
  "approval_timeout_seconds": 86400
}
```

`pat_expires_at`：PAT の有効期限日（`YYYY-MM-DD` 形式）。`null` = 未設定。`GET /api/diagnostics` の `pat` 項目で 7 日以内なら `"warn"`、期限当日以前なら `"error"` に変更。

`POST /api/config` で更新可能なキーは `log_max_lines`、`history_max_count`、`build_timeout_seconds`、`log_retention_days`、`log_archive_after_days`、`log_level`、`pat_expires_at`、`snapshots_keep`、`queue_max_size`、`build_retry_max`、`build_retry_base_seconds`、`commit_status_enabled`、`commit_status_context`、`commit_status_target_url`、`build_trend_keep_count`、`duration_anomaly`、`session_timeout_seconds`、`watch_mode`、`tag_filter`、`build_cache_enabled`、`deploy_parallelism`、`remote_build`、`approval_timeout_seconds` に限定する。各 key の型、既定値、許容値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config` schema を正本とする。未知キー、専用 schedule API の key、`api_rate_limit` を含む場合は `422` を返し、既存設定を変更しない。

**`GET /api/health` レスポンス：**

response schema、状態判定、`checks` の順序、不在・破損時の扱いは [`docs/details/api.md` 詳細本文責務 §27.16](api.md#sec-27-16) を唯一の正本とする。

**`GET /api/pat-status` レスポンス：**
```json
{ "configured": true, "expires_at": "2026-10-01", "expires_in_days": 16 }
```

`configured` は `.github_token` が通常ファイルとして存在し、trim 後の token が空でない場合だけ `true` とする。`expires_at` は `.server_config.pat_expires_at` を返し、未設定は `null` とする。`expires_in_days` は API 実行時の UTC 暦日から `expires_at` までの UTC 暦日差とし、当日は `0`、期限切れは負数、`expires_at:null` は `null` とする。この endpoint は GitHub API を呼び出さず、token 本体、prefix、長さ、file mode を返さない。

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
  "repo": "Build-Scripts",
  "updated_at": null
}
```

`RepoInfoObject` は `owner`、`repo`、`updated_at` の 3 key を必須とし、未知 key、`branch`、`target_file`、`branches` を含めない。値の型、既定値、許容値は [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) を正本とする。

**`GET /api/backup` レスポンス例：**
```json
{
  "exported_at": "2026-09-15T10:00:00Z",
  "server_config": {},
  "notify_config": {
    "webhooks": [],
    "channels": [],
    "on": [],
    "summary": {"enabled": false, "interval": "weekly", "hour": 9, "day_of_week": 1},
    "email": {"enabled": false, "to": [], "on": []}
  },
  "repo_config": {
    "owner": "fqwink",
    "repo": "Build-Scripts",
    "updated_at": "2026-09-15T10:00:00Z"
  },
  "branch_config": {"branch_targets": []},
  "access_control": {"allow": []},
  "hooks": {"hooks": []},
  "alert_rules": {"rules": []},
  "tag_rules": {"rules": []},
  "pipeline_config": {"extra_args": [], "env": {}},
  "dashboard_layout": {"widgets": ["status", "stats", "schedule", "alerts", "disk", "rate_limit", "snapshots", "maintenance", "queue"]},
  "smtp_config": {"host": null, "port": 587, "user": null, "tls": true, "from": null, "to": [], "on": [], "enabled": false},
  "webhook_secret_set": false,
  "smtp_password_set": false
}
```

**`POST /api/notify-test` レスポンス例：**
```json
{ "message": "Test notification sent", "channel_id": "n001" }
```

`.notify_config.channels[]` が 1 件以上ある場合は `type="webhook"` かつ `enabled=true` の channel、空の場合は有効な legacy `webhooks[]` を runner と同じ規則で channel へ正規化し、channel id の byte 昇順で先頭 1 件だけを test 対象とする。test は `on[]` による event 選択を適用しない。有効な Webhook 通知先がない場合は `422 {"error":"Webhook not configured"}` を返し、送信と `.notify_log` 追記を行わない。

test payload は未知 key を含まない `{"event":"notify_test","message":"Test notification"}` 固定とし、key を ASCII 昇順にした canonical JSON byte を送信と `payload_sha256` の両方に使う。HTTP redirect は追従しない。送信 attempt の成功・失敗のどちらも、response 確定前に [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.notify_log` schema どおり `event="notify_test"`、選択 channel の `channel_id`、`channel_type="webhook"`、`attempt=1` で 1 record を追記する。HTTP 応答を得た場合は `http_status` を記録し、2xx だけを `result="success"` とする。1xx、3xx、4xx、5xx、timeout、response 受信前の接続失敗は同 schema の対応する `error_code` と secret mask 後 `error` で `result="failure"` とする。test 失敗は `.notify_pending` を作成しない。

送信失敗または `.notify_log` 追記失敗は `500 {"error":"Internal server error"}` を返す。送信に成功しても log 追記に失敗した場合は success response を返さない。成功 response の `channel_id` は選択した channel id と完全一致させる。response、log、server log に Webhook URL または Webhook secret を含めない。

**`POST /api/build/force` レスポンス例：**

`POST /api/build/force` の即時開始 / queue 追加 response 例は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/build` / `POST /api/build/force` レスポンス例を参照する。

**`POST /api/pat-verify` レスポンス：**
```json
{ "valid": true, "checked_at": "2026-09-15T10:05:00Z", "scopes": ["contents:read"] }
```

`.github_token` が未設定または空の場合は `501 {"error":"Not configured"}` とする。設定済みの場合は GitHub `GET /user` を 10 秒 timeout で 1 回だけ呼び出す。GitHub が 2xx を返した場合は `valid:true`、401 または 403 を返した場合は `valid:false` とし、いずれも HTTP `200` で `checked_at` に検証完了時の UTC ISO 8601 秒精度を返す。`scopes` は 2xx response の `X-OAuth-Scopes` を comma で分割し、前後空白除去、空要素除外、byte 昇順、重複除去した配列とする。header 不在または `valid:false` は `[]` とする。network error、timeout、429、5xx は `500 {"error":"PAT verification failed"}` とし、検証結果を状態ファイルへ保存しない。

`valid:true` は `GET /user` が 2xx を返したことだけを表し、対象リポジトリの `Contents: Read` または `Commit statuses: Write` を保証しない。`scopes` は response header の観測値だけであり、repository permission の代替判定に使用してはならない。必要権限は [`docs/details/runner.md` 詳細本文責務 §17](runner.md#17-github-連携前提) を正とし、不足時は実際の repository API response に従って runner または commitstatus component が失敗を処理する。

**`GET /api/history/{id}/log` レスポンス：**

response object は [`docs/details/api.md` 詳細本文責務の `.build_logs/{id}.json` API 応答境界](api.md#22-バックエンド-api-仕様) に従う。`approval_rejected`、`approval_expired`、`skipped_dependency_failed` の history id は正常に build log を持たないため `404 {"error":"Not found"}` を返す。

**`POST /api/restore` リクエスト / レスポンス：**

リクエストは `GET /api/backup` の `BackupObject` をそのまま受け取り、`webhook_secret` と `smtp_password` だけを任意で追加できる `RestoreObject` とする。成功レスポンスは次のとおり。

```json
{ "message": "Restored" }
```

**`GET /api/notify-log` レスポンス例：**
```json
{ "log": [
    { "id": "ntfy20260915100000", "at": "2026-09-15T10:00:00Z", "event": "failure", "channel_id": "n001", "channel_type": "webhook", "payload_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "result": "success", "attempt": 1, "http_status": 200, "error_code": null, "error": null },
    { "id": "ntfy20260914183000", "at": "2026-09-14T18:30:00Z", "event": "start", "channel_id": "n001", "channel_type": "webhook", "payload_sha256": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "result": "failure", "attempt": 3, "http_status": 500, "error_code": "http_5xx", "error": "HTTP 500" }
]}
```

record key、型、`event`、`result`、nullable 条件は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.notify_log` JSON Lines schema を参照する。API は payload 本文、送信先、secret を補完しない。

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
```text
// リクエスト
{ "interval_seconds": 300 }
// レスポンス: 200
{ "message": "Schedule interval updated", "interval_seconds": 300 }
```

**`POST /api/build/cancel` レスポンス例：**
```json
{ "message": "Build cancelled" }
```

`running: false` のときに呼び出した場合は `409 Conflict` → `{"error": "No build is running"}` を返す。

**`GET /api/build/stream` — SSE ストリーミング：**

wire format、対象 log 選択、途中保存 log、frame 順、有限 close、`404`、破損時処理、副作用禁止は [`docs/details/api.md` 詳細本文責務 SSE frame 契約](api.md#build-stream-sse-contract) を唯一の正本とする。認証は他 API と同じ `Authorization: Bearer {SESSION_TOKEN}` header だけを受け付け、session token を query parameter、Cookie、body で受け付けてはならない。

**`POST /api/log-level` リクエスト / レスポンス：**
```text
// リクエスト
{ "level": "DEBUG" }
// レスポンス: 200
{ "message": "Log level changed", "level": "DEBUG" }
```

`level` の有効値：`"DEBUG"` | `"INFO"` | `"WARNING"` | `"ERROR"`

**`POST /api/pat-update` リクエスト / レスポンス：**
```text
// リクエスト
{ "token": "github_pat_..." }
// レスポンス: 200
{ "message": "PAT updated" }
```

**`GET /api/dashboard` レスポンス例：**
```json
{
  "status": {
    "last_sha": "0123456789abcdef0123456789abcdef01234567",
    "last_build_at": "2026-09-15T10:00:00Z",
    "last_build_status": "success",
    "last_target_status": "success",
    "last_trigger": "polling",
    "last_deploy_status": "success",
    "pending_transfers_count": 0,
    "notify_pending_count": 0,
    "circuit_open": false,
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

`dashboard.status` は同一 request 内で [StatusObject 固定契約](#status-object-contract) に従って 1 回だけ算出した完全な `StatusObject` とし、key の省略、追加、別名化、queue 情報の混入を行わない。`dashboard.sysinfo` は同じ request で一度だけ決定した [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target を使う。status と `sysinfo` の算出中に target を再選択してはならない。出力 directory 不在時は `sysinfo.output_size_bytes:0`、`sysinfo.output_mtime:null` とする。

**`GET /api/logs/search` レスポンス例：**
```json
{
  "query": "ERROR",
  "from": "2026-09-10",
  "to": "2026-09-15",
  "results": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "lines": ["2026-09-15T10:00:01Z [ERROR] Build failed"] },
    { "id": "b20260912183000", "build_at": "2026-09-12T18:30:00Z", "lines": ["2026-09-12T18:30:05Z [ERROR] Timeout"] }
  ]
}
```

`from` / `to` は `YYYY-MM-DD` 形式。省略時は全期間。`q` 省略時は全行返却。
`level=warn` で `[WARN]` 行のみ、`level=error` で `[ERROR]` 行のみを絞り込む。省略時は全レベルを返却する。
`SearchResult` は `query`、`from`、`to`、`level`、`results` を必須キーとする。`query` は未指定時 `""`、`from`、`to`、`level` は未指定時 `null` とする。`results` は build log の `finished_at` 降順、同時刻は build id 降順とし、各要素は `id`、`build_at`、`lines` だけを持つ。`id` は build id、`build_at` は選択 log の `finished_at` を UTC ISO 8601 秒精度で返す。`lines` は一致した message を source 順 `stdout` → `stderr` → `warnings` → `error`、各 source 内の保存順で保持する string 配列とする。一致なしは `results:[]` とする。

<a id="output-meta-response"></a>
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
  "commit_sha": "0123456789abcdef0123456789abcdef01234567",
  "build_at": "2026-09-15T10:00:00Z"
}
```

<a id="output-meta-size-diff-contract"></a>
`size_diff_bytes`：前回の同一 target の成功ビルド時との差分（正＝増加、負＝減少、`null`＝比較不能）。出力 HTML meta または対応 log の `build_id` が同じ target の成功履歴と一致する場合は、その履歴より 1 件前の成功履歴 `output_size_bytes` を比較元とする。`build_id` が履歴と一致しない場合は、同じ target の直近成功履歴 `output_size_bytes` を比較元とする。該当履歴がない場合または比較元が `null` の場合は `size_diff_bytes:null` とする。

`OutputMetaObject` は以下の 12 key をすべて必須とし、未知 key を返さない。同一 target の成功履歴を `finished_at` 降順、同時刻は `id` 降順に並べた先頭を基準履歴とし、通常 build log、通常 log 不在時の archive log の順で基準履歴と同じ `id` の log だけを対応 build log とする。`sha256` は基準履歴、`heading_count`、`tables_count`、`code_blocks_count`、`build_warnings`、`size_warn` は対応 build log だけから算出し、別 target または別 build id の log へ fallback してはならない。`build_id`、`commit_sha`、`build_at` は対応 build log の `build_meta` を優先し、対応 build log が不在または `build_meta:null` の場合だけ現在の同一 target 出力 HTML meta へ fallback する。HTML meta の build id が基準履歴と異なっても基準履歴、対応 build log、`sha256` の選択を変更しない。

| key | 型 | 生成規則 |
|-----|----|----------|
| `size_bytes` | integer | API 選択出力 target の `out` 配下にある通常 file の byte 数を再帰的に合計する。0 以上。 |
| `mtime` | string | API 選択出力 target の `out` directory 自体の mtime を [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の機械処理時刻で返す。 |
| `heading_count` | integer/null | 対応 build log の `report.headings`。`report:null` または対応 log 不在は `null`。 |
| `size_diff_bytes` | integer/null | [size difference 選択契約](#output-meta-size-diff-contract) で確定した同一 target 成功履歴との差分。増加は正、減少は負、比較不能は `null`。 |
| `tables_count` | integer/null | 対応 build log の `report.tables_count`。`report:null` または対応 log 不在は `null`。 |
| `code_blocks_count` | integer/null | 対応 build log の `report.code_blocks_count`。`report:null` または対応 log 不在は `null`。 |
| `build_warnings` | string[] | 対応 build log の `warnings` を保存順で返す。対応 log 不在は `[]`。 |
| `size_warn` | boolean | 対応 build log の `report.size_warn`。`report:null` または対応 log 不在は `false`。 |
| `sha256` | string | 同一 target の直近成功履歴 `output_sha256`。履歴不在または `output_sha256:null` は空文字。それ以外は 64 文字 lowercase hex。 |
| `build_id` | string | 対応 build log の `build_meta.build_id` を優先し、log 不在または `build_meta:null` では同 target 出力 HTML の `adlaire-build-id` meta を使う。どちらも不在は空文字。 |
| `commit_sha` | string | 対応 build log の `build_meta.commit_sha` を優先し、log 不在または `build_meta:null` では同 target 出力 HTML の `adlaire-commit-sha` meta を使う。どちらも不在は空文字。非空文字は 7〜40 文字 lowercase hex。 |
| `build_at` | string | 対応 build log の `build_meta.build_at` を優先し、log 不在または `build_meta:null` では同 target 出力 HTML の `adlaire-build-at` meta を使う。どちらも不在は空文字。非空文字は機械処理時刻形式。 |

`runner` は参照元 build log の生成だけを担当し、API response の不在値や target 選択を定義しない。

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

クエリパラメータ `n`（既定値 10）で対象件数を指定する。このレスポンス例は `n=20` に固定する。`.build_logs/{id}.json` と `.build_logs/archive/{id}.json.gz` の `duration_seconds` フィールドを集計する。

```json
{
  "n": 20,
  "count": 18,
  "avg_seconds": 38.5,
  "min_seconds": 22,
  "max_seconds": 67,
  "recent": [
    { "id": "b20260915100000", "build_at": "2026-09-15T10:00:00Z", "duration_seconds": 42, "status": "success", "target_status": "success" },
    { "id": "b20260914183000", "build_at": "2026-09-14T18:30:00Z", "duration_seconds": 7, "status": "failure", "target_status": "failure_build" }
  ]
}
```

`count` は、schema-valid で `duration_seconds` が 0 以上の log を `finished_at` 降順、同時刻は build id 降順に並べた先頭 `n` 件の件数とする。同一 build id が通常 log と archive にある場合は通常 log だけを採用する。`recent` はこの選択順のまま返し、`recent[].id=log.id`、`recent[].build_at=log.finished_at`、`recent[].status=log.status`、`recent[].target_status=log.target_status` とする。`avg_seconds`は小数第 3 位を四捨五入し小数第 2 位まで、`min_seconds` と `max_seconds` は選択対象から算出する。対象 0 件は `count:0`、`avg_seconds:null`、`min_seconds:null`、`max_seconds:null`、`recent:[]`とする。個別 log 破損と archive 展開失敗は除外し固定 WARN code だけを記録し、log directory 自体の読込不能は `500`とする。

**`GET /api/stats/build-trends` レスポンス固定契約：**

```json
{
  "n": 100,
  "count": 2,
  "samples": [
    { "build_id": "b20260914183000", "finished_at": "2026-09-14T18:30:00Z", "branch": "main", "trigger": "manual", "duration_seconds": 7, "status": "failure", "target_status": "failure_build", "anomaly": false },
    { "build_id": "b20260915100000", "finished_at": "2026-09-15T10:00:00Z", "branch": "main", "trigger": "polling", "duration_seconds": 42, "status": "success", "target_status": "success", "anomaly": false }
  ],
  "summary": { "count": 2, "avg_seconds": 24.5, "median_seconds": 24.5, "p95_seconds": 42, "anomaly_count": 0 }
}
```

`BuildTrendStats` は `n`、`count`、`samples`、`summary` を必須 key とする。API は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_trends.json.samples` の末尾から最新 `n` 件を選択し、選択後の sample を `finished_at` 昇順、同時刻は `build_id` 昇順で返す。`count` は返却 sample 数、`summary` は返却 sample だけから [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) の算出規則で再計算する。保存済み全件 summary を部分集合の response として流用しない。ファイル不在は `{"n":N,"count":0,"samples":[],"summary":{"count":0,"avg_seconds":null,"median_seconds":null,"p95_seconds":null,"anomaly_count":0}}`、破損または読込不能は `500`とする。GET は trend file と他の状態を更新しない。`warnings` key は `BuildTrendStats` に定義しない。

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

サーキットブレーカーが既に初期値なら同じ response を返し、`.build_circuit_state`、`.config_log`、`.audit_log` を変更しない。初期値と異なる場合は `.build_circuit_state` を atomic write した後、`type="circuit_breaker_reset"` の `.config_log`、`action="config_update"` / `target_id="circuit_breaker_reset"` の `.audit_log` をこの順で 1 件ずつ追記する。log 追記失敗は `500` とし、保存済み主状態と先行 log を巻き戻さない。reset 成功後に API が build または runner を起動してはならない。

**`GET /api/branch-config` response：**

`branches[]` の各 object は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.branch_config` schema の `branch_targets[]` と同じ key を返す。`deploy_targets[]` の具体構造と固定サンプル値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を正本とし、api 詳細本文責務では response wrapper と `source` の値だけを固定する。

| key | type | 必須 | 値 |
|-----|------|------|----|
| `source` | string | 必須 | `"file"` または `"default"`。 |
| `branches` | object[] | 必須 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `branch_targets[]`。 |

`source` は設定の出所を示す。`.branch_config` ファイルが存在する場合は `"file"`、存在しない場合（`BRANCH_TARGETS` デフォルト値を使用）は `"default"` を返す。

**`POST /api/branch-config` リクエスト / レスポンス：**

| 区分 | key | type | 必須 | 値 |
|------|-----|------|------|----|
| request | `branches` | object[] | 必須 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `branch_targets[]`。 |
| response `200` | `message` | string | 必須 | `"Branch config updated"`。 |
| response `200` | `branches_count` | integer | 必須 | 正規化後の `branches` 件数。 |

`branches` が空配列 `[]` の場合は `.branch_config` ファイルを削除し、`BRANCH_TARGETS` のデフォルト値に戻す（`source: "default"` に戻る）。変更は次回ポーリング周回から反映される。

**`POST /api/notify/weekly-summary` レスポンス参照：**

レスポンス例は [`docs/details/api.md`](api.md) 詳細本文責務の `POST /api/notify/weekly-summary` レスポンス例を参照する。

手動 weekly summary の集計対象、channel 抽出、送信、`.notify_log` 追記は [`docs/details/runner.md` 詳細本文責務 §27.19](runner.md#sec-27-19) の共通集計・送信契約を使用する。API はその結果を `{message:"Weekly summary sent",period,success_count,failure_count,success_rate}` へ写像する。`period` は `{period_from}/{period_to}`、両方とも UTC `YYYY-MM-DD`。`on: ["weekly_summary"]` 設定の有効な宛先がない場合は `422 {"error":"Validation failed","details":[{"field":"channels","message":"required"}]}`。1 channel 以上の送信または `.notify_log` 追記が失敗した場合は `500 {"error":"Notification failed"}` とし、success response を返さない。手動送信で `.build_state.weekly_summary_last_sent_at` と `weekly_summary_sent_date` を更新しない。

**`.build_logs/{id}.json` API 応答境界：**

`GET /api/history/{id}/log` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_logs/{id}.json` object を key の削除、別名化、状態値の再分類なしで返し、表示用 `lines` だけを追加する。`lines` は `pipeline.stdout`、`pipeline.stderr`、`warnings` を保存順のまま行配列へ変換した response 専用 field とする。API は `broken_links`、`heading_skips`、`reading_time` を top-level へ複製せず `report` 内の正本値を返す。通常 log 不在時は同 schema の archive log を読み、response 生成によって通常 log または archive を変更してはならない。

**`GET /api/diagnostics` レスポンス例：**
```json
{
  "checked_at": "2026-09-15T10:00:00Z",
  "items": [
    { "name": "pat",         "status": "ok",   "message": "PAT is valid" },
    { "name": "github_api",  "status": "ok",   "message": "GitHub API reachable" },
    { "name": "output_file", "status": "ok",   "message": "Output site exists (2.0 MB)" },
    { "name": "systemd",     "status": "ok",   "message": "adlaire-ci.timer is active and adlaire-ci.service is installed" },
    { "name": "webhook",     "status": "warn", "message": "Webhook URL not configured" }
  ]
}
```

`status` の有効値：`"ok"` | `"warn"` | `"error"`。

`output_file` は API response の固定 item 名であり、単一 HTML file ではなく [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target に対応する出力サイト directory 全体を検査する。出力 directory が不在の場合は HTTP `200` を維持し、`{"name":"output_file","status":"error","message":"Output site is missing"}` を返す。`output_file_bytes` も同じ API 選択出力 target の directory 配下の通常 file の合計 bytes を表す。

`systemd` item は shell を介さず、各 command を最大 5 秒で実行する。`systemctl is-active adlaire-ci.timer` が exit `0` かつ stdout `active`、および `systemctl cat adlaire-ci.service` が exit `0` の両方を満たす場合だけ `ok` とする。timer が inactive の場合、service unit が不在の場合、command timeout、command 起動失敗はいずれも `error` とし、oneshot の `adlaire-ci.service` 自体へ `is-active` を要求しない。診断 command の stdout / stderr は response と状態ファイルへ保存せず、固定 message だけを返す。

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

`build_logs_bytes` / `build_logs_count` は `.build_logs/{id}.json` のみを集計する。`build_logs_archive_bytes` / `build_logs_archive_count` は `.build_logs/archive/{id}.json.gz` のみを集計する。`output_file_bytes` は [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target だけを集計し、出力 directory 不在時は `0` とする。`total_bytes` は通常 build log、archive build log、`output_file_bytes` の合計とし、非選択 target の出力 directory は含めない。

**`GET /api/config-log` レスポンス例：**
```json
{ "log": [
    { "at": "2026-09-15T10:00:00Z", "type": "server_config", "action": "update", "actor": "admin", "request_id": "00112233445566778899aabbccddeeff", "endpoint": "POST /api/config", "result": "success", "error": null, "diff": { "log_max_lines": [500, 1000] }, "diff_text": "log_max_lines: 500 -> 1000" },
    { "at": "2026-09-14T18:00:00Z", "type": "notify_config", "action": "update", "actor": "admin", "request_id": "112233445566778899aabbccddeeff00", "endpoint": "POST /api/notify-config", "result": "success", "error": null, "diff": { "enabled": [false, true] }, "diff_text": "enabled: false -> true" },
    { "at": "2026-09-13T12:00:00Z", "type": "repo_config", "action": "update", "actor": "tok000001", "request_id": "2233445566778899aabbccddeeff0011", "endpoint": "POST /api/repo-config", "result": "success", "error": null, "diff": { "repo": ["Build-Scripts", "Adlaire-Design-System"] }, "diff_text": "repo: \"Build-Scripts\" -> \"Adlaire-Design-System\"" }
]}
```

record schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.config_log` JSON Lines schema、`type`、`diff`、`diff_text` の生成規則は [`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) を参照する。

**`GET /api/history/{id}/comment` レスポンス例：**
```json
{ "id": "b20260914183000", "comment": "ネットワーク障害による失敗。再ビルド済み。", "updated_at": "2026-09-14T19:00:00Z" }
```

コメント未設定時は `"comment": null`。コメントは `.build_logs/{id}.json` の `comment` フィールドに保存する。

**`POST /api/history/{id}/comment` リクエスト / レスポンス：**
```text
// リクエスト
{ "comment": "ネットワーク障害による失敗。再ビルド済み。" }
// レスポンス: 200
{ "message": "Comment saved" }
```

**`POST /api/history/{id}/flag` リクエスト / レスポンス：**
```text
// リクエスト
{ "flagged": true }
// レスポンス: 200
{ "message": "Flag updated" }
```

フラグは `.build_logs/{id}.json` の `flagged` フィールドに保存する。

**`POST /api/history/{id}/tags` リクエスト / レスポンス：**
```text
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
```text
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

**スナップショット：**

snapshot の保存、世代削除、download、delete 実体、rollback artifact の検証・展開・転送は [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を参照する。rollback build の状態書込と worker lifecycle は [`docs/details/runner.md` 詳細本文責務 §14b](runner.md#14b-スナップショット管理) を参照する。`api` 詳細では API request / response、HTTP status、config / audit log 境界だけを定義する。

**`GET /api/snapshots` レスポンス例：**
```json
{ "snapshots": [
    { "id": "b20260915100000", "build_id": "b20260915100000", "saved_at": "2026-09-15T10:00:00Z", "size_bytes": 2048576 },
    { "id": "b20260914183000", "build_id": "b20260914183000", "saved_at": "2026-09-14T18:30:00Z", "size_bytes": 2031616 }
]}
```

**`GET /api/snapshots/{id}/download`**
バイナリレスポンス。`Content-Type: application/octet-stream`、`Content-Disposition: attachment; filename="{id}.tar.gz"` を付与する。保存済み `site.tar.gz` の事前検証と stream 引渡しは [`docs/details/archive.md` 詳細本文責務 §27.15](archive.md#sec-27-15) を参照する。

<a id="snapshot-delete-api-contract"></a>
**`DELETE /api/snapshots/{id}` レスポンス例：**
```json
{ "message": "Snapshot deleted" }
```

`{id}` は build id 形式に一致し、URL decode 後に `/`、`\\`、`..`、NUL、空文字を含まない場合だけ受理する。不正は `422`。validation 後、api owner は [`docs/details/runner.md` 詳細本文責務 lock ファイル契約](runner.md#lock-ファイル契約) の snapshot delete guard を呼び出す。runner owner は `.build_lock` 取得後に `.build_state` を再読取し、idle を確定した場合だけ guard を返す。valid running lock または `.build_state.running=true` / `current_build_id!=null` は `409 {"error":"Build is running"}`、state 破損または読取失敗は `500` とし、archive owner を呼び出さない。snapshot 不在は `404`、archive 事前検証失敗は `500` とする。

archive owner の返却値ごとの api owner 処理は次に固定する。どの結果でも guard を保持したまま必要な log 処理を実行する。

| archive 結果 | `.config_log` | `.audit_log` | HTTP |
|----------------|---------------|--------------|------|
| `deleted` | `type="snapshot_delete"`、`action="delete"`、`endpoint="DELETE /api/snapshots/{id}"`、`diff={"snapshot_id":["<id>",null]}`、`result="success"`、`error=null` を追記する。 | `action="config_update"`、request actor、`target_type="config"`、`target_id="snapshot_delete"`、`result="success"`、`message="snapshot deleted"` を追記する。 | 両 log と guard 解放成功後だけ `200` と成功 response。 |
| `delete_partial` | `type="snapshot_delete"`、`action="delete"`、`endpoint="DELETE /api/snapshots/{id}"`、`diff={"snapshot_id":["<id>",null]}`、`result="partial_failure"`、`error="snapshot_delete_cleanup_failed"` を追記する。 | config log 成功後に `action="config_update"`、request actor、`target_type="config"`、`target_id="snapshot_delete"`、`result="failure"`、`message="snapshot delete cleanup failed"` を追記する。 | public snapshot は削除済みのまま `500 {"error":"Internal server error"}`。 |
| `delete_failed` | 追記しない。 | 追記しない。 | public snapshot を維持し `500 {"error":"Internal server error"}`。 |

`.config_log` 失敗時は `.audit_log` を追記せず `500`、`.audit_log` 失敗時も `500` とし、削除済み snapshot と先行 log を巻き戻さない。api owner は guard 取得後の success / failure 全終了経路で runner owner に guard 解放を 1 回要求する。解放失敗は先行結果にかかわらず `500` とし、snapshot と追記済み log を巻き戻さない。

**`POST /api/history/{id}/rollback` リクエスト / レスポンス：**

リクエスト body はなく、path parameter の `{id}` だけを使用する。成功時は HTTP `202` と次の JSON を返す。

```json
{ "message": "Rollback started", "build_id": "b20260925123000" }
```

response の `build_id` は runner owner が新規採番した rollback build id であり、path parameter の元 snapshot id と同じ値を再利用しない。

`{id}` は snapshot delete と同じ build id 検証を適用し、不正は `422` とする。元 history または `.snapshots/{id}/` が存在しない場合は `404 Not Found`、snapshot 破損は `500`、`running: true`、valid running lock、circuit open、現在の branch target または deploy target 不在は `409 Conflict` を返す。maintenance enabled は `503 {"error":"maintenance"}`、maintenance 破損または読取不能は `503 {"error":"maintenance_unavailable"}` とする。

api owner は runner owner の rollback coordinator から `prepared` と `new_build_id` を受け取った後、`.audit_log` に `action="build_trigger"`、request actor、`target_type="build"`、`target_id=<new_build_id>`、`result="success"`、固定 `message="rollback prepared"` を追記する。audit 成功後だけ `StartPreparedRollback` を呼び出し、`accepted` 取得後に `202` を返す。audit 失敗では `AbortPreparedRollback` を呼び出して `500`、worker 開始境界が accepted を返せない場合も `500` とする。rollback は manual queue entry と systemd runner 起動要求を作成しない。worker 開始後に client 切断または response 書込失敗が起きても rollback を取り消さない。

---

<a id="webhook-receive-overview"></a>

**Webhook 受信境界：**

`POST /api/webhook` は GitHub からの push event を受信し、HMAC-SHA256 署名検証後に build queue へ投入する。認証ヘッダー（`Authorization: Bearer`）は不要とし、Webhook 署名検証を認証代替として扱う。

[`docs/details/api.md` 詳細本文責務 Webhook 受信境界](api.md#webhook-receive-overview) は Webhook 受信 endpoint の概要と必須 header だけを定義する。署名検証、status code、response、event log、queue 投入、重複判定、異常系、検証条件は [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) を参照する。Webhook 受信境界へ `POST /api/webhook` の response 例、event log schema、queue entry schema を重複定義してはならない。

**`POST /api/webhook` リクエストヘッダー：**
```
X-GitHub-Event: push
X-GitHub-Delivery: <delivery_id>
X-Hub-Signature-256: sha256=<hmac_hex>
Content-Type: application/json
```

Secret は `.webhook_secret` を基準とする。secret 不在、header 不在、prefix 不正、hex 不正、署名不一致はいずれも [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) の固定契約どおり `401` とし、event log と queue を変更しない。

**Webhook Secret 設定 API：**

method、path、認証境界、request、response、error、read / write 境界は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `GET /api/webhook-config` と `POST /api/webhook-config` を参照する。

**`GET /api/webhook-config` レスポンス例：**
```json
{ "configured": true }
```

**`POST /api/webhook-config` リクエスト / レスポンス：**
```text
// リクエスト
{ "secret": "<新しいSecret文字列>" }
// レスポンス: 200
{ "message": "Webhook secret updated" }
```

---

<a id="api-cooldown-common-reference"></a>
**スケジュール強制再ビルド：**

cooldown 共通参照は [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の cooldown / force build 判定契約とする。

**`POST /api/schedule/force-interval` リクエスト / レスポンス：**
```text
// リクエスト
{ "hours": 24 }
// 無効化
{ "hours": 0 }
// レスポンス: 200
{ "message": "Force build interval updated", "hours": 24 }
```

- `.server_config.force_build_interval_hours` を保存する。runner による読込タイミングと判定適用は [cooldown 共通参照](#api-cooldown-common-reference) を使用する。
- `hours` は 0 以上の整数。0 で機能無効化

**`POST /api/schedule/cooldown` リクエスト / レスポンス：**
```text
// リクエスト
{ "seconds": 120 }
// 無効化
{ "seconds": 0 }
// レスポンス: 200
{ "message": "Build cooldown updated", "seconds": 120 }
```

- `.server_config.build_cooldown_seconds` を保存する。runner による読込タイミングと判定適用は [cooldown 共通参照](#api-cooldown-common-reference) を使用する。
- `seconds` は 0 以上の整数。0 で機能無効化

---

**メンテナンスモード：**

メンテナンスモード有効中は手動ビルド（`POST /api/build`・`POST /api/build/force`）・スケジュールビルドの両方を拒否し、`503 Service Unavailable` + `{ "error": "maintenance" }` を返す。`GET /api/health` は制限対象外とする。

**`GET /api/maintenance` レスポンス例：**
```json
{ "enabled": false, "reason": null, "since": null }
```
メンテナンス中は `{ "enabled": true, "reason": "定期メンテナンス", "since": "2026-09-15T10:00:00Z" }`。

**`POST /api/maintenance/enable` リクエスト / レスポンス：**
```text
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

メンテナンス判定は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の共通判定順に従い、`POST /api/build`、`POST /api/build/force`、署名検証済み `POST /api/webhook`、`POST /api/history/{id}/rollback` を拒否対象とする。設定参照系 GET、認証、ログ参照、メンテナンス解除は拒否しない。

`.maintenance` が破損または読取不能の場合、`GET /api/maintenance`、`POST /api/maintenance/enable`、`POST /api/maintenance/disable` は `500` を返し、既存 file を上書きしない。`POST /api/build`、`POST /api/build/force`、署名と payload の検証に成功した `POST /api/webhook`、`POST /api/history/{id}/rollback` は `503 {"error":"maintenance_unavailable"}` を返し、maintenance 無効と推測して実行してはならない。どちらの失敗でも自動修復、自動初期化、corrupt backup 作成、queue 追加、build id 採番、SHA cache 更新、event log 追記、history 追記、build log 作成、snapshot 変更を行わない。

**メンテナンス判定・副作用固定契約：**

| 対象 | 判定 | 拒否時 |
|------|------|--------|
| `POST /api/build` | 認証、rate limit、入力検証後、`.build_state` 更新前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。queue 追加、build id 採番、`.build_state` 更新を行わない。 |
| `POST /api/build/force` | force 実行要求の保存前に `.maintenance.enabled` を確認する。 | `503 {"error":"maintenance"}`。SHA cache、queue、`.build_state` を変更しない。 |
| `POST /api/webhook` | 署名検証、payload 検証後、`.webhook_events.json` 追記前に確認する。 | `503 {"error":"maintenance"}`。event log と queue を変更しない。 |
| `POST /api/history/{id}/rollback` | snapshot 存在確認後、rollback build log 作成前に確認する。 | `503 {"error":"maintenance"}`。history、build log、pending transfer を変更しない。 |
| runner 定期起動 | lock 取得後、差分検出前に確認する。 | build せず `.build_status.json.status="skipped"`、`last_target_status="skipped_maintenance"` を保存し、`.last_sha` を更新しない。 |
| runner 定期起動で `.maintenance` 破損 / 読取不能 | lock 取得後、pending retry、差分検出、build id 採番前に判定する。 | [`docs/details/runner.md` 詳細本文責務 §27.10](runner.md#sec-27-10) の status failure 固定契約を 1 回だけ適用し、終了コード `2`。build log、history、SHA cache、deploy、snapshot、pending retry を変更しない。 |

`POST /api/maintenance/enable` と `POST /api/maintenance/disable` の保存順は、`.maintenance` atomic write → `.config_log` 追記 → `config_update` audit → response とする。`.config_log` または `.audit_log` 追記失敗時は `500` を返し、保存済み `.maintenance` と先行 log は巻き戻さない。同一状態 no-op では `.maintenance`、`.config_log`、`.audit_log` を変更しない。

---

**IP アクセス制限：**

`allow` に IPv4 アドレスまたは CIDR 表記のリストを設定する。空リストは制限なし（全接続許可）を意味する。制限に一致しない接続元からのリクエストは `403 Forbidden` を返す。`GET /api/health` は制限対象外とする。設定は `.access_control` に保存する。

**`GET /api/access-control` レスポンス例：**
```json
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
```
制限なしの場合: `{ "allow": [] }`

**`POST /api/access-control` リクエスト / レスポンス：**
```text
// リクエスト
{ "allow": ["192.168.1.0/24", "10.0.0.1"] }
// レスポンス: 200
{ "message": "Access control updated", "allow": ["192.168.1.0/24", "10.0.0.1"] }
```

**アクセス制限判定固定契約：**

1. 接続元 IP は `X-Forwarded-For`、`X-Real-IP` を使わず、`net/http.Request.RemoteAddr` から取得する。
2. `RemoteAddr` が parse 不能な場合は `403 {"error":"Forbidden"}` を返す。
3. `.access_control` 不在、または `allow:[]` は全許可とする。
4. `.access_control` 破損または読取不能は全許可と推測せず、`GET /api/health` 以外を `503 {"error":"Access control unavailable"}` で拒否する。
5. CIDR は `net.ParseCIDR`、単一 IPv4 は `net.ParseIP` で判定する。
6. 判定失敗時は body 読取と認証処理より前に `403` を返し、password や token 検証を実行しない。

`POST /api/access-control` は正規化後の `allow` 配列が既存値と一致する場合、`.access_control`、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes","allow":[...] }` を返す。

**アクセス制御更新・拒否固定契約：**

| 項目 | 仕様 |
|------|------|
| 判定順 | path / method 判定後、body parse 前、認証前に実行する。拒否時は password、session token、API token、rate limit state を検証または更新しない。 |
| 対象外 | `GET /api/health` だけを対象外とする。静的 admin UI、SDK JS、その他 `/api/` 以外の配信はこの契約の対象外。 |
| allow 正規化 | 重複除去、辞書順 sort、単一 IPv4 は canonical 文字列、CIDR は `IP/mask` 表記へ正規化する。 |
| IPv6 | 初期実装では保存不可。IPv6 literal または IPv6 CIDR は `422`。 |
| private / public | private address に限定しない。入力が IPv4 または IPv4 CIDR として妥当なら保存可能。 |
| 保存順 | `.access_control` atomic write → `.config_log` 追記 → `config_update` audit → response。 |
| `.config_log` 失敗 | audit を試行せず `500`。保存済み `.access_control` は巻き戻さない。 |
| `.audit_log` 失敗 | `500`。保存済み `.access_control` と `.config_log` は巻き戻さない。 |
| 破損 / 読取不能時 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に従い自動修復と上書きを行わず、`GET /api/health` 以外の全 request を `503 {"error":"Access control unavailable"}` で拒否する。body 読取、認証、rate limit 更新、endpoint 固有処理は行わない。 |

---

**ビルドフック：**

[`docs/details/api.md`](api.md) 詳細本文責務のビルドフックは hooks API の request / response、`.hooks` 保存、`.config_log` 追記、hook log 参照境界だけを定義する。pre / post hook の実行順、timeout、process kill、hook log 保存、secret mask、pre abort、post failure、build status への影響は [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) を参照する。

**`GET /api/hooks` レスポンス例：**
```json
{ "hooks": [
    { "id": "h20260915100500", "phase": "pre",  "command_args": ["echo", "build start"], "enabled": true, "abort_on_failure": true, "timeout_seconds": 300, "created_at": "2026-09-15T10:05:00Z" },
    { "id": "h20260915100600", "phase": "post", "command_args": ["echo", "build end"],   "enabled": true, "abort_on_failure": false, "timeout_seconds": 300, "created_at": "2026-09-15T10:06:00Z" }
]}
```

**`POST /api/hooks` リクエスト / レスポンス：**
```text
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
`runs` は `ran_at` 降順、同時刻は `build_id` 降順で直近 20 件を返す。

`.hooks` record schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.hooks` schema を参照する。`.hooks` に未知 key、必須 key 不足、型不一致、不正 phase、不正 command、重複 id がある場合、`GET /api/hooks`、`POST /api/hooks`、`DELETE /api/hooks/{id}` は `500 {"error":"Internal server error"}` を返す。破損内容、command_args の secret らしき値、stdout/stderr は response と log に出さない。

**hooks API 更新順：**

| API | 更新順 | 失敗時 |
|-----|--------|--------|
| `POST /api/hooks` | 入力検証 → `.hooks` lock → id 採番 → record append → `.hooks` atomic write → `.config_log` 追記 → `config_update` audit → response | `.config_log` または audit 失敗時は `500`。追加済み record と先行 log は巻き戻さない。 |
| `DELETE /api/hooks/{id}` | path id 検証 → `.hooks` lock → 対象存在確認 → record 削除 → `.hooks` atomic write → `.config_log` 追記 → `config_update` audit → response | 対象不在は `404`。`.config_log` または audit 失敗時は `500`、削除済み record と先行 log は巻き戻さない。 |
| `GET /api/hooks/{id}/log` | path id 検証 → `.hooks` で存在確認 → `.build_logs/*_hook_{id}.json` を `ran_at` 降順、同時刻は `build_id` 降順で最大 20 件読込 → response | hook 不在は `404`。個別 hook log 破損はその file を除外し、server log に固定コードを出す。 |

hook log JSON の保存 schema、保存タイミング、失敗時の runner 挙動は [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) を参照する。`GET /api/hooks/{id}/log` は保存済み hook log を読み取り、response の `runs[]` へ `build_id`、`ran_at`、`exit_code`、`output` を返す。`output` は保存済み `stdout + stderr` をこの順で連結した表示用互換値とし、保存時点で secret mask 済みの値だけを返す。

---

**カスタムアラートルール：**

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
```text
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

<a id="tag-rule-api"></a>
**自動タグ付けルール：**

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
```text
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

タグ評価は build log の最終 status / duration / trigger が確定した後、`.build_history` 追記前に行う。手動タグと自動タグが重複した場合は 1 件に正規化し、既存順を保持したうえで自動タグを末尾へ追加する。条件式 parse 失敗を含む破損 rule がある場合、その rule を無視せず build log の `status="failure"`、`target_status="failure_tag_rule"`、history の `status="failure_tag_rule"` とし、ERROR log `TAG_RULE_INVALID` を出す。

`POST /api/tag-rules` は同一 `condition` と同一 `tags` の rule が存在する場合、`409 {"error":"Conflict"}` を返す。

---

**出力サイトチェックサム：**

build 完了時の `output_sha256` 算出と history 保存は runner owner の責務とする。API owner は `GET /api/output-meta`、`GET /api/history/{id}/log`、`POST /api/verify-output` で保存値の読取り、現在出力の再計算、比較、HTTP response だけを担当する。runner と API の両方は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の出力成果物 manifest SHA-256 を使用する。

`GET /api/output-meta` の response key、型、不在時値、参照順は [`GET /api/output-meta` レスポンス例](#output-meta-response) を唯一の正本とする。API owner は subset response や別名 field を定義しない。

`GET /api/history/{id}/log` の `output_sha256` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の build log schema にある同名 field をそのまま返す。API 固有の追加 field または別名 field は作らない。

**`POST /api/verify-output` レスポンス例：**
```text
// 一致時
{ "match": true, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "e3b0c44298fc1c149afbf4c8996fb924..." }
// 不一致時
{ "match": false, "expected": "e3b0c44298fc1c149afbf4c8996fb924...", "actual": "f4a2d5591c8f3a742f902e3b6f7c1c3d..." }
```
`expected` は [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target と `branch` / `target_file` が一致する直近成功履歴の `output_sha256` とする。同 target の成功履歴がない場合、`output_sha256` が `null` の場合、または現在の選択 target の出力サイトが存在しない場合は `404 {"error":"Not found"}` を返す。別 target の成功履歴へ fallback してはならない。

`POST /api/verify-output` は状態ファイルを変更しない。対象 path と checksum の算出は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値)、履歴と出力の選択、不在時の応答は [`docs/details/api.md` 詳細本文責務 §22.0e.1](api.md#sec-22-0e-1) の API 選択出力 target 契約に従う。path 不正、読取不能、manifest 算出不能は `500` とし、部分結果や推測値を返さない。

---

<a id="pipeline-config-api"></a>
**ビルドパイプライン設定：**

[`docs/details/api.md` 詳細本文責務 ビルドパイプライン設定](api.md#pipeline-config-api) は、`GET /api/pipeline-config` と `POST /api/pipeline-config` の request / response、`.pipeline_config` read/write 境界だけを定義する。runner による `.pipeline_config` 読込タイミング、`extra_args` / `env` 適用、読込不能または schema 不正時の build 停止条件は [`docs/details/runner.md` 詳細本文責務 §27.22](runner.md#sec-27-22) を参照する。

**`GET /api/pipeline-config` レスポンス例：**
```json
{ "extra_args": ["--verbose"], "env": { "DEBUG": "1" }, "inline_yaml": null }
```
初期値（未設定時）：`{ "extra_args": [], "env": {}, "inline_yaml": null }`

**`POST /api/pipeline-config` リクエスト / レスポンス：**
```text
// リクエスト
{ "extra_args": "<GET /api/pipeline-config response extra_args>", "env": "<GET /api/pipeline-config response env>", "inline_yaml": "<GET /api/pipeline-config response inline_yaml>" }
// レスポンス: 200
{ "message": "Pipeline config updated" }
```

`POST /api/pipeline-config` は `.pipeline_config` 全体置換とし、部分更新を許可しない。`extra_args`、`env`、`inline_yaml` のいずれかが欠ける場合は `422`。`inline_yaml` の空文字は保存前に `null` へ正規化し、空白だけの文字列は保持して runner の YAML 検証対象とする。正規化後値が既存値と一致する場合は `.pipeline_config`、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes" }` を返す。

---

**運用ノート：**

システム全体の運用メモを Markdown テキストで保存・取得する。`.notes` に保存する。認証必須。

**`GET /api/notes` レスポンス例：**
```json
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00\nPAT 更新期限: 2026-12-01", "updated_at": "2026-09-15T10:00:00Z" }
```
初回（未作成）時：`{ "content": "", "updated_at": null }`

**`POST /api/notes` リクエスト / レスポンス：**
```text
// リクエスト
{ "content": "# 運用メモ\n定期メンテナンス: 毎週日曜 2:00〜4:00" }
// レスポンス: 200
{ "message": "Notes updated", "updated_at": "2026-09-15T10:00:00Z" }
```

`.notes` は UTF-8 text として保存し、JSON ではない。`POST /api/notes` の `content` は 0〜100000 文字、NUL 禁止とする。保存時は本文をそのまま `.notes` に atomic write し、更新時刻は `.config_log` の `at` を返す。既存本文と一致する場合は `.notes`、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes","updated_at":null }` を返す。

---

**メール通知（SMTP）：**

Webhook に加えてメールでビルド結果を通知する機能。SMTP 接続設定は `.smtp_config` に、パスワードは `.smtp_secret`（パーミッション 600）に分離して保存する。`GET /api/notify-config` のレスポンスに `email` セクションを追加する。

**`GET /api/smtp-config` レスポンス例：**
```json
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true, "password_set": true }
```
`password_set`：`.smtp_secret` が存在するかを真偽値で返す。パスワード本体は返却しない。
未設定時：`{ "host": null, "port": 587, "user": null, "tls": true, "from": null, "to": [], "on": [], "enabled": false, "password_set": false }`

**`POST /api/smtp-config` リクエスト / レスポンス：**
```text
// リクエスト（次のキーだけを受け付ける。password フィールドは省略可能）
{ "host": "smtp.example.com", "port": 587, "user": "notify@example.com", "password": "s3cr3t", "tls": true, "from": "notify@example.com", "to": ["ops@example.com"], "on": ["failure"], "enabled": true }
// レスポンス: 200
{ "message": "SMTP config updated" }
```
`on` の有効値：`"start"` / `"success"` / `"failure"`

**`POST /api/smtp-test` レスポンス例：**
```text
// 成功時
{ "result": "success", "message": "Test email sent to ops@example.com" }
// 失敗時
{ "result": "failure", "message": "Connection refused: smtp.example.com:587" }
```
SMTP 未設定または `enabled: false` の場合は `422` を返す。

**SMTP 更新・送信固定契約：**

| 処理 | 仕様 |
|------|------|
| 更新 | request 正規化後に non-secret SMTP config と現在値を比較する。差分がある場合だけ `.smtp_config` を 1 回 atomic write し、差分がなければ書かない。次に password 追加 / 変更 / 削除が指定された場合だけ `.smtp_secret` を 1 回更新し、その後 `.config_log` → `.audit_log` の `config_update` の順に各 1 回書く。途中失敗時は未処理ファイルと後続 log を書かず、先行成功済み状態を巻き戻さない。 |
| 削除 | `password:null` は `.smtp_secret` 削除。削除対象が不在なら no-op。 |
| GET | `.smtp_secret` の存在だけを `password_set` で返し、password 本体は返さない。 |
| 送信 timeout | 接続、TLS、送信全体を合計 30 秒で timeout する。 |
| 送信 log | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.notify_log` schema で `event="smtp_test"`、`channel_id="smtp-test"`、`channel_type="email"`、`attempt=1`、`http_status=null` を記録する。password、宛先、SMTP response 本文は記録しない。 |

`enabled:true` にする場合は `host`、`port`、`from`、`to` 1 件以上を必須とする。`POST /api/smtp-test` は `enabled:false`、宛先なし、secret 必須構成で `.smtp_secret` 不在のいずれも `422 {"error":"SMTP not configured"}` を返す。

**SMTP 更新詳細：**

| ケース | `.smtp_config` | `.smtp_secret` | `.config_log` | `.audit_log` | response |
|--------|----------------|----------------|---------------|--------------|----------|
| config のみ変更 | 保存 | 変更なし | mask 済み diff 追記 | `config_update` / `success` | `200 {"message":"SMTP config updated"}` |
| password 追加 / 変更 | 保存 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の mode で atomic write | password は `"***"` で追記 | `config_update` / `success` | `200` |
| `password:null` | 保存 | 存在すれば削除 | password は `"***"` で追記 | `config_update` / `success` | `200` |
| 完全 no-op | 変更なし | 変更なし | 追記なし | 追記なし | `200 {"message":"No changes"}` |
| `.smtp_secret` 書込失敗 | non-secret config に差分があれば保存済み、差分がなければ未更新 | 失敗 | 追記なし | 追記なし | `500`。保存済み `.smtp_config` を巻き戻さない。 |
| `.config_log` 追記失敗 | 保存済み | 保存または削除済み | 失敗 | 追記なし | `500` |
| `.audit_log` 追記失敗 | 保存済み | 保存または削除済み | 追記済み | 失敗 | `500`。保存済み状態と `.config_log` は巻き戻さない。 |

`POST /api/smtp-test` は `.smtp_config` と `.smtp_secret` を読み、test payload を未知 key のない `{"event":"smtp_test","message":"Test email"}` 固定とする。key を ASCII 昇順にした canonical JSON byte の SHA-256 を `.notify_log.payload_sha256` とする。送信成功は `result="success"`、SMTP 接続・認証・送信失敗は `result="failure"` / `error_code="smtp_error"`、timeout は `result="failure"` / `error_code="timeout"` とし、成功 / 失敗のどちらも `.notify_log` へ 1 record を追記してから response を返す。test 送信は `.notify_pending` を作成しない。SMTP 未設定の `422` は送信 attempt ではないため `.notify_log` を追記しない。`.notify_log` 追記失敗時は `500` を返す。SMTP password、認証失敗時の server response に含まれる credential 断片、接続 URL の userinfo は response `message`、`.notify_log.error`、server log のいずれにも含めず固定文言へ置換する。

**`GET /api/notify-config` への追加（`email` セクション）：**
```json
{
  "channels": [ { "id": "n002", "type": "email", "label": "Ops", "enabled": true, "on": ["failure", "duration_anomaly"], "config": { "to": ["ops@example.com"] }, "retry_count": 2, "retry_interval_seconds": 30 } ],
  "email": { "enabled": true, "to": ["ops@example.com"], "on": ["failure"] }
}
```

---

<a id="build-queueing"></a>
**ビルドキューイング：**

`POST /api/build`・`POST /api/build/force`・署名検証済み `POST /api/webhook`・承認確定済み approval は、idle / running にかかわらず durable queue へ追加し、runner だけが queue を取り出して build を開始する。`queue_max_size`（既定 `3`）は実行中に保持できる待機 entry 数とする。`0` は実行中の追加待機を無効にするが、idle で queue が空の場合の最初の dispatch entry 1 件は受理する。上限到達時は `429 Too Many Requests` + `{ "error": "queue_full" }` を返す。

**`GET /api/queue` レスポンス例：**
```json
{ "active": null, "queued": [
    { "id": "q20260915100100", "trigger": "manual", "queued_at": "2026-09-15T10:01:00Z", "requested_by": "admin", "priority": "normal", "created_seq": 1, "payload": { "force": false } },
    { "id": "q20260915100200", "trigger": "webhook", "queued_at": "2026-09-15T10:02:00Z", "requested_by": "webhook", "priority": "normal", "created_seq": 2, "payload": { "delivery_id": "delivery-1", "branch": "main", "sha": "0123456789abcdef0123456789abcdef01234567" } }
  ], "max_size": 3 }
```
active も waiting queue も空の場合：`{ "active": null, "queued": [], "max_size": 3 }`。runner が queue entry を処理中または再実行待ちの場合は、その entry を `active` に Queue entry schema のまま返し、同じ id を `queued` に含めない。

**`DELETE /api/queue` レスポンス：**
```json
{ "message": "Queue cleared", "cleared_count": 2 }
```
active entry と実行中のビルドは停止または削除しない（`POST /api/build/cancel` を別途使用する）。`cleared_count` は削除した waiting entry 件数だけを数える。

`POST /api/config` の `queue_max_size` は整数 `0`〜`100` とする。`0` は実行中の待機 queue を無効にする値であり、idle 時の最初の dispatch entry を禁止する値ではない。

**queue 処理固定契約：**

queue 追加は `.build_state` の atomic write で行い、id は [`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2) の queue id とする。queue entry の保存値と表示順は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の Queue entry schema、および [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) の priority、created_seq、id 順を使用する。`DELETE /api/queue` は `.build_state.queued` だけを空配列にし、`active_queue_entry`、実行中 build、lock、history、log を変更しない。

valid running lock、`.build_state.running=true`、または `.build_state.active_queue_entry != null` の場合、waiting queue の保存上限は `queue_max_size` 件とし、`queue_max_size=0` では queue 対応 API を `429 {"error":"queue_full"}` とする。すべて成立しない idle 状態の保存上限は `max(1, queue_max_size)` 件とする。上限判定と append は同じ `.build_state` lock 内で行い、判定後に lock を解放してから append してはならない。

queue entry schema と trigger 別 payload schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.build_state` schema を参照する。未知 key は `422`、runner 読込時は queue entry 破損として当該 entry を処理せず ERROR ログに記録する。

**queue 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| API queue 追加 | `.build_state` lock → 最新 state 読込 → active / waiting 重複確認 → lock / running に応じた上限確認 → id と created_seq 採番 → atomic write → trigger audit → runner 起動要求 → response | write 失敗は `500` で queue 追加なし。audit 失敗は queue を残して `500`、起動要求失敗は queue を残して `202` と `dispatch:"timer_fallback"`。 |
| webhook queue 追加 | event 検証 → `.build_state` lock → queue append → `.webhook_events.json` 追記 → runner 起動要求 → response | queue 追加失敗時は event log を追記せず `500`。queue 追加後の event 追記失敗は queue を残し、`event_log_failed:true` を含む `202` とした後に起動要求を行う。 |
| approval approve queue 追加 | `.approval_queue` lock → 最新 record 確認 → `.build_state` lock → 同一 `approval_id` entry 確認 → 未存在時だけ queue append → approval status `approved` 追記 → audit → runner 起動要求 → response | queue full は `429`、approval は pending のまま。approval status または audit 失敗時は `500`、queue / approved record を巻き戻さない。確定前は起動要求を行わず、再試行時は既存 entry を再利用する。 |
| runner active 移動 | runner process lock 取得 → pending retry → circuit closed 判定 → `.build_state` lock → active があれば再利用、なければ優先度順で waiting 1 件を選択 → waiting から削除して `active_queue_entry` へ移動 → `running=true` と `current_build_id` 設定 → atomic write → build | circuit open では active / waiting を変更せず build を開始しない。start state write 失敗は build 開始なしで既存 state を保持し、process lock を解放して終了コード `1`。結果 log / history 確定後の finalizer だけが active を消去する。 |
| queue clear | `.build_state` lock → `active_queue_entry` を保持して `queued=[]` → atomic write → `.config_log` 追記 → `config_update` audit → response | `.config_log` または audit 失敗時は `500`。cleared waiting queue と先行 log は巻き戻さず、active entry は全段階で変更しない。 |

重複判定は `trigger` と `payload` の正規化 JSON が一致する `active_queue_entry` と waiting entry を、active、priority、created_seq、id の順で検索する。重複時は新規 entry を追加せず、同じ queue id に対して runner 起動要求を再実行し、`200 {"message":"Already queued","queued":true,"queue_id":"<existing>","dispatch":"requested|timer_fallback"}` を返す。ここで `queued:true` は同じ durable queue request が受理済みであることを示し、waiting 状態を意味しない。active / waiting の判定は続く `GET /api/queue` の応答を正とする。`force=true` の manual entry は `force=false` と別 entry として扱う。

**runner 起動要求固定契約：**

| 項目 | 仕様 |
|------|------|
| 実行コマンド | `exec.CommandContext(ctx, "systemctl", "start", "--no-block", "adlaire-ci.service")` と同値の引数配列。`/bin/sh -c`、文字列連結、追加引数は禁止する。 |
| timeout | 5 秒。timeout 時は child process を終了し、queue entry は保持する。 |
| 実行条件 | queue entry と endpoint 固有の audit / event / approved record が必要な順序まで確定した後だけ 1 回実行する。失敗途中の approval entry には実行しない。 |
| 成功 | systemctl 終了コード `0`。response の `dispatch` は `"requested"`。runner 起動完了や build 開始を推測しない。 |
| 失敗 | 非 `0`、timeout、binary 不在、実行不能。response の `dispatch` は `"timer_fallback"`。固定 server log `RUNNER_ACTIVATION_DEFERRED: queue_id={id}` を出し、queue、audit、event、approval を巻き戻さない。 |
| fallback | [`docs/details/setup.md` 詳細本文責務 §26.4.1](setup.md#sec-26-4-1) の `adlaire-ci.timer` による次回起動。API は timer の次回時刻を response に推測表示しない。 |
| 状態所有 | API は起動要求の前後で `.build_state.running`、`current_build_id`、`last_started_at`、`.build_status.json`、`.build_lock` を書き換えない。runner が lock 取得後にのみ更新する。 |

---

<a id="dashboard-layout-api"></a>

**ダッシュボードウィジェットカスタマイズ：**

ダッシュボードに表示するウィジェットの種類・順序を設定できる機能。`.dashboard_layout` に保存する。

有効なウィジェット識別子：`status`（ビルド状態）/ `stats`（統計サマリー）/ `schedule`（次回実行）/ `alerts`（アラート）/ `disk`（ディスク使用量）/ `rate_limit`（GitHub API レート制限）/ `snapshots`（スナップショット件数）/ `maintenance`（メンテナンス状態）/ `queue`（キュー状態）

**`GET /api/dashboard-layout` レスポンス例：**
```json
{ "widgets": ["status", "alerts", "stats", "schedule", "disk"] }
```
未設定時はデフォルト順（全ウィジェット）を返す。

**`POST /api/dashboard-layout` リクエスト / レスポンス：**
```text
// リクエスト
{ "widgets": ["status", "alerts", "schedule", "stats"] }
// レスポンス: 200
{ "message": "Dashboard layout updated" }
```
`widgets` に未知の識別子が含まれる場合は `422` を返す。

**dashboard layout 固定契約：**

既定 widget 順は `["status","stats","schedule","alerts","disk","rate_limit","snapshots","maintenance","queue"]` とする。`POST /api/dashboard-layout` は `widgets` 全体置換のみ許可し、空配列、重複、未知 id は `422`。正規化後値が既存値と一致する場合は `.dashboard_layout`、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes" }` を返す。

---

**`POST /api/logs/cleanup` レスポンス例：**
```json
{ "message": "Cleanup completed", "deleted_count": 12, "failed_count": 0 }
```

`log_retention_days` が `0` の場合は削除せず `deleted_count: 0`、`failed_count: 0` を返す。

**`GET /api/history/export` レスポンス：**

`Content-Type: application/json` で返却される。

```json
{ "exported_at": "2026-09-15T10:00:00Z", "history": [] }
```

`ExportObject` は `exported_at` と `history` だけを必須 key とし、未知 key を返さない。`exported_at` は response 作成完了時の UTC ISO 8601 秒精度とする。`history` は filter と paging を適用しない全有効行を `GET /api/history` と同じ重複 id 除外、`finished_at` 降順、同時刻 id 降順で返し、各 item は `HistoryPageObject.history[]` と同じ全保存 key、`build_at`、`sha` を持つ。行単位破損、重複 id、ファイル不在、読込不能は `GET /api/history` と同じ扱いとし、export によって `.build_history` を変更しない。前方互換 category id は history item 内に原値のまま含める。`ExportObject` に paging key と top-level `warnings` は追加しない。

**`POST /api/repo-config` リクエスト / レスポンス：**
```text
// リクエスト（owner と repo のうち 1 key 以上）
{ "owner": "fqwink", "repo": "Adlaire-Design-System" }
// レスポンス: 200
{ "message": "Repo config updated" }
```

`RepoConfigPatch` は `owner` と `repo` だけを許可する JSON object とし、少なくとも 1 key を必須とする。型と許容値は [`docs/details/statefile.md` 詳細本文責務 `.repo_config` schema](statefile.md#repo-config-schema) を使用する。`updated_at`、`branch`、`target_file`、`branches` を含む未知 key は `422 {"error":"Validation failed","details":[...]}` とし、`.repo_config`、`.branch_config`、`.config_log`、`.audit_log` を変更しない。

API は `readRepoConfig()` の正規化済み現在値へ request の指定 key を適用する。差分がない場合は `{ "message":"No changes" }` を返し、差分がある場合は server が `updated_at` を request 完了時刻に設定して `.repo_config` を atomic write する。`GET /api/repo-info` と次回起動する runner は同じ `readRepoConfig()` の値を使用する。実行中 runner の repository 識別子を途中で差し替えない。

branch、監視対象、source / output path、deploy target の表示と更新は `GET /api/branch-config` と `POST /api/branch-config` の責務とする。`POST /api/repo-config` が `.branch_config` を読み書きしたり、`POST /api/branch-config` が `.repo_config` を読み書きしたりしてはならない。

history response の `trigger` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.build_history.trigger` をそのまま返す。status summary 専用の `startup_config_integrity` を history へ追加してはならない。

共通エラーの HTTP status と response body は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の API 共通エラー固定文言を正本とする。状態競合、通知未設定、queue 上限、maintenance などの endpoint 固有エラーは [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の固定表と対象 endpoint 契約を正本とし、ここでは再定義しない。

---

## 25. 認証 実装仕様

認証、session、password hash、login ticket、TOTP 連携、認証ログ、漏えい禁止、`--init-credentials` 生成手順の主本文は [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) および [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45)〜[§27.46](security.md#sec-27-46) を参照する。`.admin_credentials` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。

`api` 詳細では、認証関連 API の endpoint、request / response、HTTP status、状態ファイル read / write 境界だけを定義する。

| API / CLI | API 側の担当 | 主本文 |
|-----------|--------------|--------|
| `POST /api/login` | route、body parse、response body、HTTP status、`.admin_credentials` / `.totp_secret` read/write 呼び出し境界。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細)、[`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) |
| `POST /api/login/totp` | route、body parse、response body、HTTP status、`.totp_secret` / `.admin_credentials` read/write 呼び出し境界。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細)、[`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) |
| `POST /api/logout` | route、body 禁止、response body、HTTP status。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) |
| `POST /api/change-password` | route、body parse、response body、HTTP status、`.admin_credentials` write 呼び出し境界。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) |
| `GET /api/sessions` / `POST /api/sessions/revoke-all` | route、response body、HTTP status、memory session 操作呼び出し境界。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細)、[`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) |
| `--init-credentials` | CLI option dispatch、stdout / stderr / exit code を security 契約どおり返す。 | [`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) |
| `GET /api/auth/totp-status` / `POST /api/auth/totp-setup` / `POST /api/auth/totp-confirm` / `DELETE /api/auth/totp` | route、body parse、response body、HTTP status、`.totp_secret` read/write 呼び出し境界。 | [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) |

---

## 27. api owner 追加仕様化機能 詳細仕様

<a id="sec-27-5"></a>
**27.5 設定バリデーション API：**
[`docs/details/api.md` 詳細本文責務 §27.5](api.md#sec-27-5) の境界は owner component `api`、collaborator component `sdk`、`ui`、`statefile` とする。


`POST /api/config/validate` は、`POST /api/config` と同じ入力を受け取り、保存せずに検証結果を返す。

Request body は partial `ConfigObject` とする。未知 key を含む場合は `422 {"error":"Validation failed","details":[{"field":"<unknown-key>","message":"Unknown config key"}]}` を返す。型不一致、範囲外、URL 不正、commit status context 空文字、retry 設定不正は `200` で `valid=false` として返す。JSON parse 失敗は API 共通契約どおり `400 {"error":"Invalid JSON"}` とする。

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

`config` は既存 `.server_config` に request body を merge した正規化後の値を返す。ただし保存してはならない。通常 API 呼び出しとして追記する `.api_access_log` 以外の状態ファイルまたはログを更新してはならない。

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

<a id="sec-27-6"></a>
**27.6 API アクセスログ：**
[`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の境界は owner component `api`、collaborator component `sdk`、`ui`、`statefile` とする。


`api` は全 `/api/` request について `.api_access_log` へ JSON Lines を追記する。`GET /api/health` も対象とする。静的 file 配信、admin HTML、SDK JS は対象外とする。

追記タイミングは response status 確定後とする。追記失敗時は、対象 API の本来の response を優先し、サーバーログに `API_ACCESS_LOG_WRITE_FAILED` を出す。access log 書き込み失敗を理由に API response を `500` へ変更してはならない。

`GET /api/api-access-log` は `limit`、`offset`、`method`、`path`、`status` query を受け付ける。`limit` は 1〜1000、既定値 100。`offset` は 0 以上、既定値 0。`method` は大文字 HTTP method 完全一致。`path` は prefix match。`status` は HTTP status 完全一致。filter、並び順、paging、破損行、`total` は [`docs/details/api.md` 詳細本文責務 §22.0e.3](api.md#sec-22-0e-3) の JSON Lines 一覧取得固定契約に従う。

**access log field 導出固定契約：**

record の key、型、必須性は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.api_access_log` JSON Lines schema を正本とする。`api` owner は各 field の導出方法を以下に固定する。

| key | 導出仕様 |
|-----|----------|
| `at` | response status 確定後、`.api_access_log` 追記直前の UTC 秒精度時刻。 |
| `request_id` | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の 32 文字 lowercase hex。response header `X-Request-Id` と一致させる。 |
| `method` | request の HTTP method を大文字化した値。 |
| `path` | URL decode や path 正規化で別名化せず、query を除いた request path。 |
| `query` | route の query 契約で許可され、[`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の検証に成功した明示指定 key だけを保存する。number、boolean、固定 enum、`YYYY-MM-DD` は正規化後の値を保存し、それ以外の string 値は常に `"***"` へ置換する。既定値、未指定 key、未知 key、検証失敗値、raw query string は保存しない。path / method / query 検証完了前の response は `{}` とする。 |
| `status` | 送信対象として確定した HTTP status。access log 追記失敗で変更しない。 |
| `duration_ms` | handler 受付開始直後の monotonic clock から response status 確定までの経過をミリ秒で切り捨て、0 未満は 0 とする。 |
| `auth_type` | 終了時の認証経路を `"none"`、`"session"`、`"api_token"`、`"webhook"` のいずれかで記録する。認証前失敗と認証失敗は `"none"` とする。 |
| `actor` | session 認証成功は `"admin"`、API token 認証成功は token id、Webhook 署名検証成功は `"webhook"`、それ以外は `null` とする。 |
| `remote_addr` | `http.Request.RemoteAddr` を `net.SplitHostPort` で分解した host を IP 文字列として保存する。分解または IP parse に失敗した場合は `null`。`X-Forwarded-For`、`Forwarded`、`X-Real-IP` は参照しない。 |
| `user_agent` | `User-Agent` を `strings.TrimSpace` し、改行と C0 / DEL 制御文字を除去した先頭 512 Unicode scalar values。空文字は `null`。 |

request body、raw query string、未知 query key、検証失敗 query 値、cookie、Authorization header、token、password、secret は保存しない。access log 追記失敗時も、既に確定した response status を変更しない。

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 認証成功 API | actor `admin` または token id で記録される。 |
| 認証失敗 API | status 401、actor null、secret なし。 |
| Webhook | auth_type `webhook`、actor `webhook`。 |
| filter | method/path/status が完全に効く。 |
| 書込失敗 | 本来 response を維持し、server log に WARN。 |
| request id | response header と log の `request_id` が一致する。 |
| secret query | raw query string と secret / free-form query の実値を保存しない。許可済み query key だけを固定規則で保存する。 |

<a id="sec-27-11"></a>
**27.11 ポーリング間隔の動的変更：**
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は schedule 設定と systemd timer 更新の責務を持つ。`runner` は本機能で systemd timer を変更しない。

本機能の目的は、管理 API から systemd timer の実行間隔を変更し、次回以降の runner 起動間隔を固定仕様どおり反映することである。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/schedule/interval` |
| Request | `{ "interval_seconds": integer }` |
| 許容値 | 30〜86400 秒 |
| Response | `{ "message": "Schedule interval updated", "interval_seconds": N }` |
| 状態 | `.server_config.schedule_interval_seconds` を更新し、`.config_log` と `.audit_log` の `config_update` へ記録する。 |
| systemd | `/etc/systemd/system/adlaire-ci.timer` の `OnUnitActiveSec` を `N seconds` 相当へ更新する。 |

**処理順序：**

1. 認証、maintenance、入力型、範囲を検証する。
2. memory 上で `schedule_interval_seconds` の diff を生成する。diff が空の場合は `.server_config`、systemd、`.config_log`、`.audit_log` を変更せず、`200 {"message":"No changes","interval_seconds":N}` を返して終了する。
3. diff が空でない場合は `.server_config` を atomic write で更新する。rename 前の write failure では処理を終了する。rename 成功後の file / parent directory Sync または lock 削除失敗では手順 4〜7 を実行せず、手順 8 で `result="partial_failure"`、`error="state_write_after_rename_failed"` を追記する。
4. systemd timer drop-in を書き換える。
5. `systemctl daemon-reload` を実行する。
6. `systemctl restart adlaire-ci.timer` を実行する。
7. `systemctl show adlaire-ci.timer -p OnUnitActiveSec` 相当で反映を確認する。
8. 手順 4〜7 がすべて成功した場合は `.config_log` に `result="success"`、`error=null` の record を 1 件追記する。手順 4〜7 のいずれかが失敗した場合は、後続の systemd 操作を中止し、`result="partial_failure"`、`error="systemd_update_failed"` の record を 1 件追記する。
9. 手順 8 の `.config_log` 追記成功後、`config_update` を `.audit_log` へ追記する。`.config_log.result="success"` なら audit result は `success`、`partial_failure` なら `failure` とする。
10. 手順 4〜7 と手順 9 がすべて成功した場合は成功 response、systemd または post-rename partial failure がある場合は audit 追記後に `500` を返す。audit 追記失敗時も `500` とし、保存済み状態、systemd 操作、`.config_log` は巻き戻さない。

**systemd 更新固定契約：**

| 項目 | 仕様 |
|------|------|
| timer 書換 | 既存 unit を直接編集せず、管理対象 drop-in `/etc/systemd/system/adlaire-ci.timer.d/override.conf` を atomic write する。 |
| drop-in 内容 | `[Timer]`、`OnUnitActiveSec={N}s`、`Persistent=true` のみを書き込む。 |
| 反映確認 | `systemctl show` の値を秒へ正規化して request 値と一致確認する。 |
| rollback | systemd 失敗時に `.server_config` は巻き戻さない。`.config_log` に `result="partial_failure"`、`error="systemd_update_failed"`、`.audit_log` に `config_update` / `result="failure"` を 1 件ずつ追記する。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 入力が integer でない、または範囲外 | `422`。状態ファイル、systemd は変更しない。 |
| 既存値と同一 | `200 {"message":"No changes","interval_seconds":N}`。`.server_config`、systemd、`.config_log`、`.audit_log` は変更しない。 |
| `.server_config` の rename 前書き込み失敗 | `500`。`.server_config`、systemd、`.config_log`、`.audit_log` は変更しない。 |
| `.server_config` の rename 成功後の Sync または lock 削除失敗 | `500`。systemd は変更しない。新しい `.server_config` は維持し、`.config_log` に `result="partial_failure"`、`error="state_write_after_rename_failed"`、`.audit_log` に `config_update` / `result="failure"` を記録する。 |
| systemd timer 書き換え / daemon-reload / restart / show 失敗 | `500`。`.server_config` は更新済みのまま残し、server log に `SCHEDULE_INTERVAL_APPLY_FAILED`、`.config_log` に `result="partial_failure"`、`error="systemd_update_failed"`、`.audit_log` に `config_update` / `result="failure"` を記録する。 |
| `.config_log` 追記失敗 | `500`。`.server_config` と実行済み systemd 操作は巻き戻さず、server log に `CONFIG_LOG_WRITE_FAILED` を追記する。失敗した config log record と後続 audit record は作成しない。 |
| `.audit_log` 追記失敗 | `500`。`.server_config`、実行済み systemd 操作、`.config_log` は巻き戻さず、server log に `AUDIT_LOG_WRITE_FAILED` を追記する。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常値 300 | `.server_config`、timer、API response が 300 で一致する。 |
| 最小値 30 | 成功する。 |
| 範囲外 29 | `422`、差分なし。 |
| 同一値 | `200` と `{ "message":"No changes","interval_seconds":N }`。状態、systemd call、config log、audit log に差分なし。 |
| systemd 失敗 | `500`、`.config_log` に `result="partial_failure"` / `error="systemd_update_failed"`、`.audit_log` に `config_update` / `result="failure"` を 1 件ずつ記録する。 |
| show 不一致 | `500`、server config は更新済み。 |

<a id="sec-27-12"></a>
**27.12 GitHub Webhook 受信：**
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。`api` は署名検証、イベント記録、queue 投入を担当し、`runner` は queue entry を処理する。

本機能の目的は、GitHub push event を HMAC-SHA256 署名検証したうえで受信し、build queue へ投入して runner の非同期起動を要求することである。起動要求を systemd が受理できない場合も queue を保持し、定期 timer を fallback とする。

**入力 / 出力：**

| 項目 | 仕様 |
|------|------|
| API | `POST /api/webhook` |
| 必須 header | `X-GitHub-Event`, `X-GitHub-Delivery`, `X-Hub-Signature-256` |
| 対象 event | `push` のみ |
| Secret | `.webhook_secret` |
| 成功 response | `{ "message": "Webhook accepted", "queued": true, "event_id": "...", "queue_id": "...", "dispatch": "requested\|timer_fallback" }` |
| 状態 | `.webhook_events.json` へ追記する。対象 push の after SHA が対象 branch の直近成功 SHA と異なり、maintenance が無効で、queue 上限に空きがあり、同一 delivery id、branch、SHA の active / waiting entry が存在しない場合だけ `.build_state.queued` へ `trigger="webhook"` entry を追加する。 |

**署名検証：**

署名は `sha256=` prefix を含む lowercase hex とする。検証は raw request body に対して `HMAC-SHA256(secret, body)` を計算し、定数時間比較で行う。secret 不在、header 不在、prefix 不正、hex 不正、署名不一致はすべて `401` とし、queue へ投入しない。

**処理順序：**

1. method と body size を検証する。1 MiB 超過時は secret と署名を検証せず `413` を返す。
2. `.webhook_secret` を読み込む。
3. HMAC 署名を検証する。
4. [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) に従い、`trigger` group の IP key だけで rate limit を判定・更新する。上限超過時は JSON parse 前に `429 {"error":"Too many requests"}` を返す。
5. JSON body を parse する。
6. event が `push` であることを確認する。
7. `ref`、`after`、`repository.owner.login`、`repository.name` を抽出する。
8. `.branch_config` と照合する。対象 branch がある場合は `.build_state` lock 内で active / waiting の重複、queue 容量を判定し、重複がなければ waiting queue へ entry を atomic write する。
9. queue 追加後、または重複する既存 queue id の確定後に、`.webhook_events.json` へ `queued` または `duplicate` のイベント結果を JSON Lines で追記する。対象 event / branch でない場合は queue を変更せず `ignored_event` または `ignored_branch` を追記する。
10. queue id を確定した request は、event log 追記の成否を確定した後、[`docs/details/api.md` 詳細本文責務 ビルドキューイング](api.md#build-queueing) の runner 起動要求を 1 回実行する。queue 追加後の event log 失敗でも起動要求へ進み、`event_log_failed:true` を返す。queue 追加前の失敗では起動要求を行わない。
11. queue を使用した response に `queue_id` と `dispatch` を含めて返す。

**Webhook 副作用固定契約：**

| ケース | `.webhook_events.json` | `.build_state.queued` | response |
|--------|------------------------|-----------------------|----------|
| 署名不正 / secret 不在 | 変更なし | 変更なし | `401 {"error":"Unauthorized"}` |
| rate limit 上限超過 | 変更なし | 変更なし | `429 {"error":"Too many requests"}` |
| JSON parse 失敗 | 変更なし | 変更なし | `400 {"error":"Invalid JSON"}` |
| event が `push` 以外 | `ignored_event` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| 対象 branch なし | `ignored_branch` を追記 | 変更なし | `202 {"message":"Webhook ignored","queued":false,"event_id":...}` |
| queue 追加成功 | `queued` を追記 | queue entry を追加 | `202 {"message":"Webhook accepted","queued":true,"event_id":...,"queue_id":...,"dispatch":"requested\|timer_fallback"}` |
| queue full | `queue_full` を追記 | 変更なし | `429 {"error":"queue_full"}` |
| event log 追記失敗 / queue 追加前 | 変更なし | 変更なし | `500 {"error":"Internal server error"}` |
| event log 追記失敗 / queue 追加後 | 変更なし | queue entry は残す | runner 起動要求後に `202 {"message":"Webhook accepted","queued":true,"event_log_failed":true,"queue_id":...,"dispatch":"requested\|timer_fallback"}` |

queue entry は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.build_state` schema の queue entry schema を使用し、`trigger:"webhook"`、`requested_by:"webhook"`、`priority:"normal"`、`payload.delivery_id`、`payload.branch`、`payload.sha` を保存する。`X-GitHub-Delivery` が `active_queue_entry` または waiting queue に存在し、同一 branch / sha の場合は重複投入せず、event log に `result:"duplicate"`、既存 `queued_id` を記録し、同じ queue id の runner 起動要求を再実行して `202 {"message":"Webhook already queued","queued":true,"queue_id":"<existing>","dispatch":"requested|timer_fallback"}` を返す。active と waiting の両方に同じ delivery が存在する破損状態では active の queue id を使用し、waiting duplicate を自動削除せず server log に `QUEUE_DUPLICATE_STATE: queue_id={id}` を記録する。

**Webhook validation 固定契約：**

| 項目 | 仕様 |
|------|------|
| body size | 最大 1 MiB。超過時は署名検証前に `413`。 |
| branch 抽出 | `ref` が `refs/heads/{branch}` 形式でない場合は `ignored_branch`。 |
| sha | `after` が 40 文字 lowercase hex でない場合は `422`。 |
| repository | `owner.login` と `name` から `owner/name` を作る。欠落時は `422`。 |
| event id | base は `wh{YYYYMMDDHHmmss}` とし、衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| secret | signature、secret、raw payload は event log、access log、audit log に保存しない。 |

**Webhook 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| 署名前処理 | raw body は署名検証用 byte slice として保持し、JSON parse、文字コード変換、空白整形、改行変換を署名前に行わない。 |
| body size | body は 1 MiB + 1 byte まで読み、1 MiB 超過を検出した時点で `413` とする。secret 読取、署名検証、JSON parse、event log、queue 操作は行わない。上限以内の場合だけ署名検証へ進む。 |
| rate limit | 署名検証成功後、JSON parse 前に `trigger` group の IP key だけを判定・更新する。`429` 時は event log、queue、runner 起動要求を変更または実行しない。 |
| 保存順 | queue 追加が必要な場合は `.build_state` lock 取得 → queue 追加 atomic write → `.webhook_events.json` 追記 → response の順とする。 |
| queue 追加後の event log 失敗 | queue entry は巻き戻さず、response に `event_log_failed:true` を含める。server log には固定 code だけを出し、payload 内容は出さない。 |
| event id 採番 | fake clock の UTC 秒を使う。prefix、suffix、上限到達時の失敗は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。 |
| duplicate 判定 | `X-GitHub-Delivery`、branch、after SHA が一致する active / waiting entry だけを duplicate とする。delivery id だけ一致して branch または sha が異なる場合は `409 {"error":"Conflicting delivery"}` とし、queue を変更しない。 |
| 対象 branch | `.branch_config` が存在する場合は `branch_targets[].branch`、不在時は default branch 設定と照合する。照合できない branch は `ignored_branch` として event log だけ残す。 |
| maintenance | maintenance 有効時は署名、JSON、branch 検証後、event log / queue 追記前に `503` とし、状態を変更しない。 |

**異常系：**

| 条件 | 応答 / 処理 |
|------|-------------|
| 署名不正 | `401`、イベントログ追記なし、queue なし。 |
| rate limit 上限超過 | `429 {"error":"Too many requests"}`、`.api_rate_state` の count は増やさず、[`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の audit 成功後だけ返す。event log、queue、runner 起動要求は変更または実行しない。 |
| event が `push` 以外 | `202`、`queued=false`、イベントログには `ignored_event` として記録する。 |
| JSON parse 失敗 | `400 {"error":"Invalid JSON"}`、event log と queue の変更なし。 |
| 対象 branch なし | `202`、`queued=false`、イベントログに記録する。 |
| queue 上限 | `429`、イベントログに `queue_full` を記録する。 |
| イベントログ書き込み失敗 | `500`、queue 追加前なら queue しない。queue 追加後なら response に `queued=true` と `event_log_failed=true` を含める。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 push | 署名検証成功、イベントログ追記、queue entry `trigger="webhook"`。 |
| 署名不一致 | `401`、状態差分なし。 |
| Webhook rate limit | 署名成功後・JSON parse 前に `trigger` group の IP key だけを使い、上限未満は count +1、上限超過は `429` かつ event log / queue 差分なし。 |
| 対象外 branch | `202 queued=false`、イベントログのみ。 |
| queue full | `429`、queue 差分なし。 |
| duplicate | queue 追加なし、既存 queue id を返す。 |
| payload oversized | `413`、状態差分なし。 |

<a id="sec-27-13"></a>
**27.13 Webhook イベントログ / 一覧取得 API：**
[`docs/details/api.md` 詳細本文責務 §27.13](api.md#sec-27-13) の境界は owner component `api`、collaborator component `sdk`、`ui`、`statefile` とする。

本機能の目的は、受信した GitHub Webhook の監査情報を `.webhook_events.json` に保存し、管理 API、sdk、ui のページング参照対象にすることである。

`.webhook_events.json` の保存 schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.webhook_events.json` JSON Lines schema を参照する。保存時に request header 全体、署名値、secret、payload 全体を保存してはならない。

**WebhookEventRecord 生成固定契約：**

キー、型、必須性、許容値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.webhook_events.json` schema を唯一の正本とする。`api` owner は各 field の導出だけを以下に固定する。

| key | 導出元 | API 生成規則 |
|-----|--------|--------------|
| `id` | event id 採番 | [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) の event id。 |
| `timestamp` | API clock | event 結果確定後、追記直前の UTC ISO 8601 秒精度。 |
| `delivery_id` | `X-GitHub-Delivery` | 空、NUL、改行を拒否した header 値。 |
| `event` | `X-GitHub-Event` | 受信した GitHub event 名。 |
| `ref` | payload `ref` | 欠落時は `null`。 |
| `branch` | payload `ref` | `refs/heads/{branch}` から抽出し、抽出不能時は `null`。 |
| `sha` | payload `after` | 40 文字 lowercase hex。`ignored_event` / `ignored_branch` で取得不能な場合だけ `null`。 |
| `repository` | payload repository | `owner.login` と `name` から `owner/name` を生成し、取得不能時は `null`。 |
| `result` | Webhook 副作用固定契約 | `queued`、`duplicate`、`ignored_event`、`ignored_branch`、`queue_full` の該当結果。maintenance、署名失敗、JSON parse 失敗、payload validation 失敗は event を保存しない。 |
| `build_triggered` | queue write 結果 | queue entry を新規追加した場合だけ `true`。 |
| `queued_id` | queue entry | 新規または既存 queue id。queue に関係しない結果では `null`。 |
| `error_code` | Webhook 処理結果 | `queue_full` では `"queue_full"`、その他は `null`。secret、payload 断片、signature を含めない。 |

未知 key は保存しない。API response では [`docs/details/api.md` 詳細本文責務 §27.13](api.md#sec-27-13) の固定表の key だけを返す。

**一覧 API：**

`GET /api/webhook-events` は `limit` と `offset` query を受け付ける。`limit` は 1〜1000、既定値 50。`offset` は 0 以上、既定値 0。`timestamp` 降順、同時刻は file 出現順の逆順で返す。壊れた行は無視し、server log に `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` を出す。

Response は `{ "events": WebhookEventRecord[], "total": N }` の 2 key だけとし、両 key を必須とする。`events` の各要素は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の WebhookEventRecord、`total` は paging 前の schema-valid record 件数とし、未知 key を返さない。SDK 引数変換は [`docs/details/sdk.md` SDK 引数変換契約](sdk.md#sdk-argument-contract)、UI 表示と再取得は [`docs/details/ui.md` UI 操作契約表](ui.md#ui-operation-contract) を参照する。

**webhook events 取得固定契約：**

| 項目 | 仕様 |
|------|------|
| 並び順 | `timestamp` 降順、同時刻は file 出現順の逆順。 |
| total | 壊れた行を除外した総件数。 |
| offset | filter 後、並び替え後に適用する。 |
| 壊れた行 | 内容を response、server log に含めない。固定コードだけ出す。 |

**Webhook events API 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| query | 許可 query は `limit` と `offset` だけ。未知 query は `422`。 |
| read-only | `.webhook_events.json` の破損行、未知 key、古い schema を API 取得時に修復しない。 |
| collaborator 参照 | SDK request は [`docs/details/sdk.md` SDK 引数変換契約](sdk.md#sdk-argument-contract)、UI 操作と表示は [`docs/details/ui.md` UI 操作契約表](ui.md#ui-operation-contract) を正本とする。API は SDK の既定値や UI の表示列を定義しない。 |
| server log | corrupt line skip は `WEBHOOK_EVENT_LOG_SKIP_CORRUPT` の固定 code と file path だけを出す。行本文は出さない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| event 追記 | JSON Lines へ schema 通り保存される。 |
| 一覧取得 | `timestamp` 降順、同時刻は file 出現順の逆順で limit/offset が効く。 |
| 壊れた行 | API は継続し、壊れた行を返さない。 |
| duplicate delivery | queue 重複なし、event log は `duplicate`。 |
| total | 壊れた行を除外した件数。 |

<a id="sec-27-16"></a>
**27.16 ヘルスチェックエンドポイント：**
[`docs/details/api.md` 詳細本文責務 §27.16](api.md#sec-27-16) の境界は owner component `api`、collaborator component `statefile` とする。

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

`status` は `"ok"` または `"degraded"` とする。`checks` が空の場合だけ `"ok"`、1 件以上の場合は `"degraded"` とする。response object 構築失敗、JSON encode 失敗、response 書き込み開始前の header 生成失敗では `500` とし、`HealthObject` を返さない。response body の `status:"error"` は使用しない。

**読み取り元：**

`.build_status.json` を第一参照元とし、不在、破損、読込不能時は `.build_history` から build 関連値を fallback 算出する。fallback 対象は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の共通順序で選んだ最新の status summary 対象行だけとし、history-only 行を採用しない。対象行がなければ `last_build_at:null`、`last_build_status:"none"` とする。`pending_transfers` は `.pending_transfers`、notify pending の check は `.notify_pending`、`uptime_seconds` は API process 起動時刻から算出する。`.build_status.json` 破損または読込不能時は自動修復せず、`checks[]` に対応する check を含める。fallback では `last_deploy_at` と `last_deploy_status` を推測せず `null` とする。

**health checks 固定契約：**

| check | 条件 |
|-------|------|
| `build_status_missing` | `.build_status.json` 不在で fallback 算出した。 |
| `build_status_corrupt` | `.build_status.json` parse 失敗。 |
| `build_status_read_error` | `.build_status.json` が permission または I/O error で読み取れない。 |
| `build_history_read_error` | status fallback に必要な `.build_history` を読み取れない。 |
| `pending_transfers_present` | `.pending_transfers` に 1 件以上の entry がある。 |
| `pending_transfers_read_error` | `.pending_transfers` 読取失敗。 |
| `notify_pending_read_error` | `.notify_pending` 読取失敗。 |
| `runner_stale` | `running=true` かつ `last_started_at` が 24 時間より古い。 |

`checks[]` は [`docs/details/api.md` 詳細本文責務 §27.16](api.md#sec-27-16) の固定表の順で返す。同じ check は 1 回だけ含める。`.build_status.json` 不在、破損、読込不能時は read-only で `.build_history` の最新の status summary 対象行へ fallback し、状態ファイルの作成、修復、退避、上書きを行わない。history-only 行は fallback から除外する。

**health 実装確認固定契約：**

| 条件 | HTTP status | `status` | `checks` | 副作用 |
|------|-------------|----------|----------|--------|
| 全参照成功 | `200` | `ok` | `[]` | 状態差分なし。 |
| `.build_status.json` 不在 | `200` | `degraded` | `build_status_missing` | fallback 算出のみ。作成しない。 |
| `.build_status.json` 破損 | `200` | `degraded` | `build_status_corrupt` | 修復、backup、削除を行わない。 |
| `.build_status.json` 読込不能 | `200` | `degraded` | `build_status_read_error` | fallback 算出のみ。chmod、修復、backup、削除を行わない。 |
| history fallback 読取失敗 | `200` | `degraded` | `build_history_read_error` | 履歴を変更せず、build 関連 nullable 値を `null`、status を `none` とする。 |
| pending transfer あり | `200` | `degraded` | `pending_transfers_present` | entry を消費または変更しない。 |
| `.pending_transfers` 読取失敗 | `200` | `degraded` | `pending_transfers_read_error` | `pending_transfers:0` を返し、pending を初期化しない。 |
| `.notify_pending` 読取失敗 | `200` | `degraded` | `notify_pending_read_error` | notify pending を初期化しない。 |
| JSON encode 不能 | `500` | response なし | response なし | 書き込み済み header がなければ `500`。状態差分なし。 |

`GET /api/health` は認証、access control、maintenance、rate limit、session timeout の拒否対象にしない。ただし path / method 判定は通常どおり行い、`POST /api/health` は `405` とする。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 正常 | HTTP 200、`status="ok"`。 |
| pending transfer あり | `pending_transfers` に件数、`status="degraded"`。 |
| build status 破損 | HTTP 200、`status="degraded"`、checks に記録。 |
| build status 読込不能 | HTTP 200、history fallback、`build_status_read_error`、状態差分なし。 |
| history-only 最新行 | その行を fallback から除外し、直前の status summary 対象行を返す。対象行がなければ `last_build_status="none"`。 |
| pending 読込不能 | HTTP 200、`pending_transfers:0`、`pending_transfers_read_error`、状態差分なし。 |
| response 構築不能 | HTTP 500、`HealthObject` なし、状態差分なし。 |
| stale runner | checks に `runner_stale`。 |

<a id="sec-27-17"></a>
**27.17 ビルドログ重大度フィルター：**
owner component は `api` とする。collaborator component は `sdk`、`ui`、`archive`、`statefile` とする。

本機能の目的は、`GET /api/logs/search` と UI ログビューアで重大度別に build log を絞り込めるようにすることである。

**入力：**

`GET /api/logs/search` は既存 query に加えて `level` を受け付ける。`level` の許容値は `"info"`、`"warn"`、`"warning"`、`"error"`、`"debug"` とし、大文字小文字は区別しない。正規化後は `"INFO"`、`"WARNING"`、`"ERROR"`、`"DEBUG"` とする。`warn` は `"WARNING"` と同義とする。不正値は `422`。

**検索対象：**

`.build_logs/{id}.json.pipeline.stdout`、`pipeline.stderr`、`warnings`、`error`、archive log を対象とする。行頭が `[WARN]` または `[WARNING]` の行は WARNING、`[ERROR]` または stderr の非空行は ERROR、`[DEBUG]` は DEBUG、それ以外は INFO と分類する。

**ログ行分類固定契約：**

| 入力 | level |
|------|-------|
| stderr の非空行 | `ERROR` |
| stdout / warnings の `[ERROR]` prefix | `ERROR` |
| stdout / warnings の `[WARN]` または `[WARNING]` prefix | `WARNING` |
| stdout / warnings の `[DEBUG]` prefix | `DEBUG` |
| stderr 非空行でも `[ERROR]`、`[WARN]`、`[WARNING]`、`[DEBUG]` prefix 付き行でもない入力 | `INFO` |

検索中の各行は `{build_id, level, source, line_number, message}` として分類する。`source` は `"stdout"`、`"stderr"`、`"warnings"`、`"error"` のいずれか、`line_number` は 1 始まりとする。これらは filter と固定順の決定に使う内部表現であり、HTTP response に `level`、`source`、`line_number`、`message` を追加しない。response は [`docs/details/api.md`](api.md) 詳細本文責務の `SearchResult` 固定契約に従い、build ごとに一致 message を `lines` へグループ化する。

**ログ検索実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| query | 許可 query は `q`、`from`、`to`、`level` だけ。未知 query は `422`。paging は定義せず、`limit`、`offset` を受け付けない。 |
| 日付 | `from` / `to` は `YYYY-MM-DD`。UTC 日付として扱い、`from` は 00:00:00 以上、`to` は 23:59:59 以下を含む。`from > to` は `422`。 |
| q | 部分一致。大文字小文字は区別しない。空文字は全件。NUL、改行を含む q は `422`。 |
| 並び順 | `results` は build log の `finished_at` 降順、同時刻は build id 降順。`lines` は source 順 `stdout` → `stderr` → `warnings` → `error`、各 source 内の line number 昇順。 |
| archive | 同一 build id が通常 log と archive にある場合は通常 log を優先し、二重に返さない。 |
| 破損 log | 破損通常 log または gzip 展開失敗は除外し、固定 WARN code だけを server log に出す。本文は出さない。 |
| read-only | 検索 API は log、archive、history、status、config を変更しない。 |

SDK request の引数変換は [`docs/details/sdk.md` SDK 引数変換契約](sdk.md#sdk-argument-contract)、UI の filter 操作と表示は [`docs/details/ui.md` UI 操作契約表](ui.md#ui-operation-contract) を正本とする。API は `q`、`from`、`to`、`level` の受信、検証、検索、response 生成だけを担当する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| `level=warn` | WARNING 行だけ返る。 |
| `level=error` | ERROR 行だけ返る。 |
| 不正 level | `422`。 |
| archive log | 通常 log と同じ分類で検索される。 |
| stderr | prefix なしでも ERROR。 |
| grouped response | `level`、`source`、`line_number`、`message` は HTTP response に返さず、一致 message だけを build ごとの `lines` に固定順で返す。 |

<a id="sec-27-18"></a>
**27.18 ブランチ設定の動的変更 API：**
owner component は `api` とする。collaborator component は `runner`、`statefile` とする。

本機能の目的は、監視対象 branch / target / deploy target を `.branch_config` で管理し、API 経由の変更対象を `.branch_config` に限定することである。

**API：**

`GET /api/branch-config` は `.branch_config` が存在する場合 `{"source":"file","branches":[...]}`、不在の場合 `{"source":"default","branches":[...]}` を返す。`POST /api/branch-config` は `{ "branches": BranchTargetRecord[] }` を受け取る。

**保存仕様：**

永続ファイルの key は必ず `branch_targets` とする。API request / response で `branches` を使う場合も保存前に `branch_targets` へ変換する。空配列を受け取った場合は `.branch_config` を削除し、default 復帰とする。

**branch config 正規化固定契約：**

| 項目 | 仕様 |
|------|------|
| branch 一意性 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config` schema が定める一意性に違反する request は `422` とする。 |
| deploy target id | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) に従い、省略値を request 配列順から導出した後の一意性違反を `422` とする。 |
| target_files / approval_required / env | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の型、既定値、正規化、検証を適用する。 |
| 保存順 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config` 正規化順をそのまま使用する。 |
| 削除 | POST empty で `.branch_config` を削除した後、`.config_log` に default 復帰、`.audit_log` に `config_update` を記録する。 |

**branch config 保存・失敗時固定契約：**

| 操作 | 保存順 | 失敗時 |
|------|--------|--------|
| GET default | default branch target を response へ正規化する。 | 状態ファイルを作成しない。default 算出不能なら `500`。 |
| GET file | `readBranchConfig()` が返す正規化済み target を `source:"file"` で返す。 | 破損は `500 {"error":"State file is corrupted"}`、読取不能は `500 {"error":"State file read failed"}`。退避、削除、default 書込みを行わない。 |
| POST valid | request 検証 → 正規化 → `.branch_config` atomic write → `.config_log` 追記 → `config_update` audit → response。 | write 失敗は `500`、`.config_log` / audit なし。`.config_log` または audit 失敗は `500`、保存済み `.branch_config` と先行 log は巻き戻さない。 |
| POST empty | `.branch_config` 存在確認 → 削除 → `.config_log` 追記 → `config_update` audit → response。 | 削除失敗は `500`。`.config_log` または audit 失敗は `500`、削除済み状態と先行 log は巻き戻さない。 |
| POST no-op | 正規化後の `branch_targets` が既存値と一致する。 | `.branch_config`、`.config_log`、`.audit_log` を変更せず `{ "message":"No changes","branches_count":N }` を返す。 |

`.config_log` の diff target は `branch_config` とする。deploy target の `host`、`user`、`dest_dir` は secret として扱わないが、値に token / password / secret 風 key が含まれる object を追加した場合は [`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) の mask 規則を適用する。

**検証：**

request の必須 key、任意 key、型、件数上限、path 制約、env 制約、deploy target id は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.branch_config` schema を正本とする。API は任意 key の既定値を正規化した後、全 key を永続 object に明示保存する。

**runner 取り込み：**

runner による `.branch_config` の読込、`RunnerConfig.BranchTargets` への正規化、起動中の反映タイミングは [`docs/details/runner.md` 詳細本文責務 §12](runner.md#12-設定値runner)〜[§13](runner.md#13-処理フロー) の runner 設定正規化契約を参照する。`api` 詳細では API endpoint、request / response、保存、削除、検証条件だけを定義する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| GET default | `.branch_config` 不在で `source="default"`。 |
| POST valid | `.branch_config.branch_targets` として保存、`.config_log` と `config_update` audit を 1 件ずつ追記する。 |
| POST empty | `.branch_config` 削除、default 復帰。 |
| 相対 `sha_file` | `422`、状態差分なし。 |
| branch 重複 | `422`、状態差分なし。 |
| POST empty log failure | `.branch_config` は削除済み、response は `500`。 |

<a id="sec-27-20"></a>
**27.20 設定変更の詳細 diff 記録：**
[`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) の境界は owner component `api`、collaborator component `statefile` とする。

本機能の目的は、設定変更 API が何を変更したかを `.config_log` に機械可読 diff と人間可読 diff の両方で残すことである。

**対象 API：**

対象は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の Write 列に `.config_log` を持つ全 endpoint とする。対象 endpoint の追加・削除は [`docs/details/api.md` 詳細本文責務 §22.0d](api.md#sec-22-0d) の状態アクセス固定表と [`docs/details/api.md` 詳細本文責務 §27.20](api.md#sec-27-20) を同一変更で整合させ、本文中に別の手書き endpoint 一覧を作成しない。

`.config_log` のログ schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.config_log` JSON Lines schema を参照する。差分がない場合、対象 API は状態ファイル、secret file、`.config_log`、`.audit_log` を変更せず、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に定義された no-op response を返す。対象 endpoint 契約が no-op response を定めていない endpoint だけ `{ "message": "No changes" }` を返す。

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
| no-op | diff が空の場合は状態ファイル、secret file、`.config_log`、`.audit_log` を変更しない。 |
| 追記順 | 対象状態ファイル保存後に `.config_log`、続いて `.audit_log` の `config_update` を追記する。`.config_log` 失敗時は audit を試行しない。いずれかの追記失敗時は `500`、保存済み状態と先行 log は巻き戻さない。 |

`.config_log` record の `type` は endpoint 固定名、`actor` は管理 session なら `"admin"`、API token なら token id とする。`request body` 全体、HTTP header、cookie、secret 平文を保存してはならない。

**diff 対象 endpoint 固定名：**

| endpoint | type |
|----------|--------|
| `POST /api/config` | `server_config` |
| `POST /api/notify-config` | `notify_config` |
| `POST /api/repo-config` | `repo_config` |
| `POST /api/branch-config` | `branch_config` |
| `POST /api/webhook-config` | `webhook_config` |
| `POST /api/smtp-config` | `smtp_config` |
| `POST /api/pipeline-config` | `pipeline_config` |
| `DELETE /api/snapshots/{id}` | `snapshot_delete` |
| その他 `.config_log` 対象 | path template から `/api/` prefix を除き、`{name}` は `name` へ変換し、`/` と `-` を `_` に置換した固定名。先頭または末尾の `_` は除去する。HTTP method と path parameter の実値は使用しない。 |

**config diff 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| actor | 管理 session は `"admin"`、API token は token id、未認証で許可される設定変更 API は存在しない。 |
| request id | access log と response header の request id と同じ値を `request_id` として保存する。request ID 生成失敗時は endpoint 処理へ進まないため `.config_log` を追記しない。 |
| endpoint | `endpoint` は HTTP method と path template を保存する。path param の実値が secret 風値でも path template だけを保存する。 |
| action | 対象状態が不在から作成された場合は `create`、既存状態を変更した場合は `update`、対象を削除した場合は `delete` とする。複数状態を変更する restore 等は `update` とする。 |
| result | 主状態変更と必須後続処理が成功した record は `"success"`。主状態変更後に systemd、sync、外部適用等の必須後続処理が失敗した record は `"partial_failure"`。主状態変更前の失敗では record を作成しない。 |
| error | `result="success"` では `null`。`result="partial_failure"` では [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約が定める lowercase snake case の固定 error code を保存し、Go error、command output、path、secret を保存しない。 |
| config log write failure | `.config_log` 追記失敗時は `CONFIG_LOG_WRITE_FAILED` を server log に記録し、audit は試行せず response は `500`。追記に失敗した record は存在しない。保存済みの主状態は巻き戻さない。 |
| audit | `.config_log.type` を `target_id` とする `config_update` を 1 件追記する。`.config_log.result="success"` は audit `result="success"`、`partial_failure` は audit `result="failure"` とする。`POST /api/api-rate-limit` だけ action を `rate_limit_update` とする。 |
| audit write failure | `.audit_log` 追記失敗時は `AUDIT_LOG_WRITE_FAILED` を server log に記録して `500`。保存済みの主状態と `.config_log` は巻き戻さない。 |
| rollback 禁止 | `.config_log` または `.audit_log` の追記失敗時に保存済み設定と先行 log を巻き戻さない。巻き戻さないことを fixture で固定する。 |
| no-op | diff が空の場合は `.config_log`、対象状態ファイル、secret file、`.audit_log` を変更しない。 |
| secret | secret 平文は diff、diff_text、server log、access log、`.audit_log`、response に出さない。 |

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| 単一 key | key 昇順 diff、diff_text、type、action、actor、request_id、endpoint、`result="success"`、`error=null` が保存される。 |
| nested object | dot path で差分が保存される。 |
| array | 正規化後の配列全体が JSON 値として diff に入る。 |
| 複数 key | key 昇順で `diff_text` を生成する。 |
| type 名 | path template から固定 type 名が生成される。 |
| no-op | 状態、secret、`.config_log`、`.audit_log` に差分なし。endpoint 固有の no-op response。 |
| secret | before / after と diff_text が `"***"` になる。 |
| config log failure | 対象設定は保存済み、audit は未実行、response は `500`、未定義 rollback なし。 |
| audit failure | 対象設定と `.config_log` は保存済み、response は `500`、未定義 rollback なし。 |
| diff 生成失敗 | 状態差分なしで `500`。 |

diff 生成は状態保存前に memory 上で完了させる。diff 生成に失敗した場合は状態ファイルを書かない。`.config_log` または `.audit_log` 追記に失敗した場合は保存済み状態と先行 log を巻き戻さず、response は `500` とする。

**異常系：**

| 条件 | 処理 |
|------|------|
| config log 追記失敗 | 対象状態は保存済みのまま維持し、API response を `500` とする。 |
| audit log 追記失敗 | 対象状態と `.config_log` は保存済みのまま維持し、API response を `500` とする。 |
| diff 生成失敗 | 状態ファイルを書かず `500`。 |
| secret マスク漏れ検出 | 実装不合格。該当 API は詳細実装確認を満たした扱いにしない。 |

<a id="sec-27-21"></a>
**[`docs/details/api.md` 詳細本文責務 API 連動境界確認表](api.md#sec-27-21)：**

[`docs/details/api.md` 詳細本文責務 §27.21〜§27.38 / §27.42〜§27.47](api.md#sec-27-21) の固定表は、api owner が関与する場合の入力境界、出力境界、状態 read/write 呼び出し境界、失敗時副作用、fixture 参照を確認する表である。runner / builder / security / fixture の主本文は、表の節リンク先と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定する。

詳細実装確認では、対象機能の owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、必要な collaborator component 別の詳細本文責務を参照する。api が owner ではない行は、api が受け渡す endpoint、response、状態 read/write 境界の確認だけに使用する。fixture 証跡参照列は fixture 正本への参照入口であり、fixture 名、input、expected、fake、合格条件の正本ではない。

| API 連動境界確認節 | 機能 | 入力 | 出力 | 状態ファイル / 外部副作用 | 失敗時副作用 | fixture 証跡参照 |
|--------------------|------|------|------|---------------------------|--------------|----------------|
| [`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21) | 複数ファイル監視 | `.branch_config.branch_targets[].target_files`、GitHub content SHA または local SHA。 | `changed_targets[]`、target 単位 SHA cache、build log。 | 成功時だけ該当 target SHA cache を更新する。 | SHA 部分失敗では build を開始せず、成功取得済み cache も更新しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.22](runner.md#sec-27-22) | pipeline YAML | `.pipeline.yml` または `.pipeline_config.inline_yaml`。 | `pipeline_steps[]`、step stdout/stderr、build status。 | step を定義順に実行し、build 終了時に定義順で保存する。 | parse / command 不正では build を開始しない。required step 失敗で後続 required step を実行しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) | local watch | `.server_config.watch_mode`、local `src` 配下 Markdown。 | `.local_watch_state.json`、trigger `local_watch`。 | local mode では GitHub API / PAT を呼ばず、成功時だけ state を置換する。 | file read 失敗は build なし。dry-run は state を作成 / 更新しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.24](runner.md#sec-27-24) | tag filter | `.server_config.tag_filter`、GitHub tags refs。 | `matched_tags[]`、skip status。 | tag 一致時だけ build。tag 不一致 skip は `.build_status.json` だけ更新する。 | tags API 最終失敗は build なし。local mode 併用は終了コード `2`。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/builder.md` 詳細本文責務 §27.25](builder.md#sec-27-25) | build cache | `--cache-dir`、cache 設定、input / deps SHA。 | cache entry、page cache、`[REPORT]` cache counts。 | hit 時は変換結果を再利用し、miss 成功時だけ entry を書く。 | cache read/write 失敗は build を成功可能にし、WARN / report に残す。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.26](runner.md#sec-27-26) | parallel deploy | `.server_config.deploy_parallelism`、deploy targets。 | `target_results[]`、`.pending_transfers`。 | target ごとに並列転送し、result は設定順で保存する。 | worker 内部失敗は該当 target failure。他 target は継続する。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) | build hooks | `.hooks`、pre/post hook command_args。 | hook log、build log hook result。 | pre は build 前、post は build 後に id 昇順で実行する。 | pre abort で build 本体を開始しない。post 失敗は build status を変更しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/builder.md` 詳細本文責務 §27.28](builder.md#sec-27-28) | dependency tracking | builder dependency manifest、Markdown link / asset reference。 | `.dependency_manifest.json`、affected target 判定。 | build 成功時だけ manifest を更新する。 | manifest 破損は full build 扱い。dry-run は更新しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.29](runner.md#sec-27-29) | remote build | remote build 設定、ssh target、artifact path。 | remote artifact metadata、build log remote section。 | remote command と artifact fetch を行い、検証成功時だけ deploy / history へ進む。 | remote timeout / checksum mismatch は failure とし、secret / command credential を保存しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) | approval API | approval 一覧取得、approve/reject API。 | `.approval_queue`、`.build_state.queued[]`、`.build_history`。 | pending entry だけ approve/reject できる。 | entry 不在は `404`、pending 以外は `409`、queue full は `429`。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31) | branch env | branch env 設定、pipeline / hook env。 | merged env、masked log。 | runner env → branch env → step/hook env の順で上書きする。 | env key 不正は保存不可 / runner 設定エラー。secret はログに出さない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) | build notification | notify config、build event、channel 設定。 | notify payload、`.notify_log`、`.notify_pending`。 | event ごとに payload を生成し、送信結果を記録する。 | retry 対象失敗は pending。secret は payload/log に含めない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) | build trends | build duration samples、history/log。 | `.build_trends.json`、trend stats response。 | build 完了時に sample を追加し、上限件数で trim する。 | trend 保存失敗は build 成否を反転しない。破損時は再集計契約に従う。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.34](runner.md#sec-27-34) | build chain | chain config、upstream build result。 | chain execution log、queued child build。 | chain 条件一致時だけ次 build を queue する。 | chain config 破損は chain 無効として通常 build は継続する。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) | priority queue | queue request、priority、created_seq。 | sorted queue、queue API response。 | queue 保存時に priority / created_seq を固定する。 | queue full は `429`、既存 queue を変更しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.36](runner.md#sec-27-36) | failure classification | build / deploy 結果、failure evidence。 | build log の `failure_category` / `failure_evidence[]`、history の `failure_category`、API filter。 | runner finalizer で分類し、build log と history に保存する。`.build_status.json` に category / evidence を追加しない。 | 分類不能でも build 結果は保持し、category `unknown` とする。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.37](runner.md#sec-27-37) | execution environment | runner host/runtime info、builder version。 | build log の `environment` object、`GET /api/history/{id}/log` response。 | build id 採番直後に取得し、build log に保存する。`.build_status.json` に environment を追加しない。 | 取得失敗は `unknown` を保存し build は継続する。secret / env 全量は保存しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) | duration anomaly | duration samples、anomaly config。 | anomaly flag、alert / notify、stats。 | build 完了時に閾値判定し、該当時だけ alert/notify を作る。 | samples 不足では判定しない。通知失敗は build 成否を反転しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | trigger API scope | API token scopes、endpoint group。 | allow/deny decision、audit。 | scope 一致時だけ endpoint 実行。 | scope 不足は `403`、endpoint 固有処理なし。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | API key management | token label/scopes/expires_at。 | `.api_tokens` record、token 本体 1 回 response。 | token 作成 / 失効 / 認証成功時に状態と audit を更新する。 | audit 失敗時は token 本体を返さない。token 破損 state は自動再生成しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | audit log | 認証、権限、設定、token、build trigger event。 | `.audit_log` JSON Lines、audit API response。 | 監査対象操作の成否確定後に追記する。 | 追記失敗は対象操作を `500` 扱い。ただし保存済み状態は個別契約どおり戻さない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | session timeout | `.server_config.session_timeout_seconds`。 | 新規 session `expires_at`。 | session 発行直前の設定で期限を計算する。 | 範囲外は `422`。既存 session の期限は変更しない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP | setup secret、ticket、TOTP code。 | `.totp_secret`、ticket/session response、audit/access log。 | secret / ticket はメモリと 1 回 response に限定し、成功時だけ永続状態を更新する。 | code 不一致 / replay は token を返さない。audit 失敗時も secret 平文を出さない。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | API rate limit | rate policy、remote addr、actor key。 | `.api_rate_state`、`429`、state summary。 | key 群を同一 lock で判定 / 更新する。 | 上限超過では count を増やさず endpoint 固有処理を行わない。audit 失敗時は `500`。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |

<a id="sec-27-21-2"></a>
**[`docs/details/api.md` 詳細本文責務 API / SDK / UI 連動参照表](api.md#sec-27-21-2)：**

[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表は api、sdk、ui の接続点をそろえるための参照表である。API endpoint は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint 契約、SDK method は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI 操作は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様)、fixture と実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を参照する。[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表だけを根拠に endpoint、method、DOM、状態ファイル、fixture を追加してはならない。

| 節 | API | SDK | UI |
|----|-----|-----|----|
| [`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21) | branch config API と status/history/log に反映する。 | `getBranchConfig()` / `setBranchConfig()` は `target_files` を削除しない。 | リポジトリ情報 panel で複数 target を表示 / 保存する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.22](runner.md#sec-27-22) | pipeline config API。 | `getPipelineConfig()` / `setPipelineConfig()`。 | 設定 panel で pipeline config を表示 / 保存し、shell 文字列へ変換しない。 |
| [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) | `GET/POST /api/config` の `watch_mode`。 | `getConfig()` / `setConfig()`。 | 設定 panel で `github` / `local` を選択する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.24](runner.md#sec-27-24) | `GET/POST /api/config` の `tag_filter`。 | `getConfig()` / `setConfig()`。 | 設定 panel で tag filter を表示 / 保存する。 |
| [`docs/details/builder.md` 詳細本文責務 §27.25](builder.md#sec-27-25) | `GET/POST /api/config` の cache key と builder report。 | `getConfig()` / `setConfig()`。 | 設定 panel と build result 表示で cache counts を表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.26](runner.md#sec-27-26) | `GET/POST /api/config` の `deploy_parallelism` と status/history。 | `getConfig()` / `setConfig()`。 | 設定 panel で parallelism を表示 / 保存し、結果は履歴 / log で表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.27](runner.md#sec-27-27) | hooks API。 | hook methods。 | フック panel で command_args を 1 行 1 引数として表示 / 保存する。 |
| [`docs/details/builder.md` 詳細本文責務 §27.28](builder.md#sec-27-28) | endpoint 追加なし。manifest は runner/builder 内部状態。 | SDK method 追加なし。 | UI 操作追加なし。依存情報を表示する場合は、既存 build result 表示の範囲だけを使用する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.29](runner.md#sec-27-29) | `GET/POST /api/config` の `remote_build`。 | `getConfig()` / `setConfig()`。 | 設定 panel で remote build 設定を表示 / 保存する。 |
| [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) | approvals API。 | approval methods。 | 承認待ち panel で approve / reject を操作する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31) | `GET/POST /api/branch-config` の `branches[].env`。 | `getBranchConfig()` / `setBranchConfig()`。 | リポジトリ情報 panel で branch env を表示 / 保存する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) | notify config / notify log API。 | notify methods。 | 通知設定 panel で channel、payload、履歴を表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.33](runner.md#sec-27-33) | stats/build trends API。 | `getBuildTrends()`。 | 統計 panel で trend を表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.34](runner.md#sec-27-34) | build chain config API。 | `getBuildChainConfig()` / `setBuildChainConfig()`。 | 設定またはリポジトリ情報 panel で chain を表示 / 保存する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.35](runner.md#sec-27-35) | queue API。 | `getQueue()` / `clearQueue()` と build trigger methods。 | 手動実行 panel で priority 反映後 queue を表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.36](runner.md#sec-27-36) | history/stats API の filter と response に含める。 | history / stats methods は category を保持する。 | 履歴 panel で failure category filter を表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.37](runner.md#sec-27-37) | output meta / status / history log に含める。 | `getOutputMeta()` / `getStatus()` / `getHistoryLog()`。 | システム情報または履歴 detail に表示する。 |
| [`docs/details/runner.md` 詳細本文責務 §27.38](runner.md#sec-27-38) | config / stats / dashboard alerts に含める。 | config / stats / dashboard methods。 | 統計 panel と dashboard alert で anomaly を表示する。 |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | endpoint ごとの scope 判定。 | sdk は token scope を推測せず API error を返す。 | ui は `403` を権限不足として表示し、logout しない。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | token API。 | token methods。 | API token 管理 panel で発行 token を 1 回だけ表示する。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | audit log API。 | `getAuditLog()`。 | 監査ログ panel で filter 表示し、secret は表示しない。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | config API。 | `setConfig({session_timeout_seconds})`。 | セキュリティ panel で session timeout を表示 / 保存する。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP/auth API。 | TOTP/auth methods。 | セキュリティ/login panel で one-time secret/ticket flow を扱う。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | rate limit API。 | `getApiRateLimit()` / `setApiRateLimit()`。 | セキュリティ panel で policy と state summary を表示する。 |

api / sdk / ui のいずれも、[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表に存在しない endpoint、method、UI 操作を追加してはならない。追加が必要な場合は、[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint、[`docs/details/sdk.md` SDK 引数変換契約](sdk.md#sdk-argument-contract)、[`docs/details/ui.md` UI 操作契約表](ui.md#ui-operation-contract)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を先に更新する。対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) または [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約が endpoint 追加なしと定める機能は、runner / builder の内部挙動または既存 response field の範囲で実装する。

<a id="sec-27"></a>
**[`docs/details/api.md` 詳細本文責務 API / SDK / UI 連動実装確認ゲート](api.md#sec-27)：**

[`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21)〜[§27.38](runner.md#sec-27-38) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) のうち API、SDK、UI が連動する機能の詳細実装確認では、[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表の全条件を満たす。owner component が `api` ではない機能でも、API response を SDK / UI が利用する場合は [`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) api / sdk / ui 連動実装確認ゲートを満たす。

| ゲート | API 側の合格条件 | SDK / UI への固定契約 | 禁止条件 |
|--------|------------------|------------------------|----------|
| endpoint 対応 | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) 機能契約、[`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表の API 列が同じ method / path / query / body / response を示す。 | SDK method と UI 操作は [`docs/details/api.md` 詳細本文責務 §27](api.md#27-api-owner-追加仕様化機能-詳細仕様) API / SDK / UI 接続固定表の SDK / UI 列だけを使用する。 | 表にない endpoint、method、UI 操作を実装都合で追加すること。 |
| request 正規化 | path parameter、query、body key、nullable、既定値を API 側で検証し、不正値は `422` とする。 | SDK は指定値を表どおり送信し、UI は入力正規化だけを行う。 | SDK が未知 key を削除する、UI が API 既定値を保存前に補完すること。 |
| response 透過 | 成功 response は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に定義された key だけを返し、array 順、nullable、mask 値を固定する。 | SDK は response を補完・再計算せず返し、UI は API 順序で表示する。 | SDK / UI が存在しない key、集計値、状態名、token list を合成すること。 |
| error 伝播 | `401` / `403` / `409` / `422` / `429` / `500` の status と error body を固定する。 | SDK は `AdlaireCIError` として保持し、UI は status 別表示と仕様上の再取得だけを行う。 | 自動 retry、自動 refresh、自動 logout、同一変更 API の再送を仕様外で行うこと。 |
| side effect | validation 失敗、認可失敗、rate limit、no-op、partial failure の write / call / log 差分を [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint と [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) で固定する。 | SDK / UI は副作用完了を推測せず、成功後再取得で確認する。 | read-only、dry-run、validation failure、`403`、`429` で、[`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の共通処理順序または対象 endpoint に許可されていない状態、log、外部副作用を発生させること。 |
| secret / one-time | token、PAT、Webhook secret、SMTP password、TOTP secret、ticket、Authorization header は response・log・state の許可箇所以外へ出さない。 | SDK は内部保存せず、UI は [`docs/details/ui.md` 詳細本文責務 one-time secret 消去契約](ui.md#ui-one-time-secret-contract) の専用領域、trusted event、copy、即時消去、generation 一致条件に従う。 | token 本体の再表示、token list への合成、secret の error message / 許可外 DOM / expected への残存。 |
| fixture 証跡 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の API / SDK / UI 連動 fixture で request、response、error、side effect、secret、refresh order を確認する。 | [`docs/details/fixture.md`](fixture.md) fixture 証跡責務の実装検証証跡に対象 fixture、未実装対象、未定義 endpoint / UI / 状態ファイル不追加を記録する。 | fixture なし、または実装挙動に合わせて期待値を弱めること。 |

<a id="sec-27-30"></a>
**27.30 ビルド承認フロー：**
機能 owner は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) を参照する。`api` は collaborator として API endpoint、request / response、状態ファイル read/write 呼び出し境界を担当し、`sdk`、`ui`、`statefile` との接続境界だけを定義する。

[`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) は、ビルド承認フローにおける API endpoint、request / response、状態ファイル read/write 呼び出し境界だけを定義する。`approval_required` 検出、pending 作成、approval request 通知、timeout 処理、承認済み queue entry の実行は [`docs/details/runner.md` 詳細本文責務 §27.30](runner.md#sec-27-30) を参照する。

**入力 / 状態：**

| 項目 | 仕様 |
|------|------|
| 状態 | `.approval_queue` JSON Lines。mode は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の状態ファイル固定値。 |
| API | `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject` |
| status | `"pending"`、`"approved"`、`"rejected"`、`"expired"`。 |

**正常系：**

1. `GET /api/approvals` は `.approval_queue` を読み取り、id ごとの最新 record だけを返す。
2. `POST /api/approvals/{id}/approve` は最新 status が `pending` の entry だけを承認し、`.build_state.queued[]` に `trigger:"approval"` の queue entry を追加する。
3. approve 成功後、`.approval_queue` に `approved` record、`.audit_log` に `approval_approved` を追記し、同じ queue id の runner 起動を要求してから response を返す。
4. `POST /api/approvals/{id}/reject` は最新 status が `pending` の entry だけを却下し、`.approval_queue` に `rejected` record を追記する。
5. reject 成功後、`.build_history` に `status:"approval_rejected"`、`.audit_log` に `approval_rejected` を追記し、response を返す。

**approval API 固定契約：**

| 項目 | 仕様 |
|------|------|
| queue payload | approve で追加する queue entry は `trigger:"approval"`、`priority:"normal"`、`requested_by` は approve actor id とする。`payload` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の approval payload とし、pending record から `approval_id`、`branch`、`sha`、`target`、`requested_trigger`、`requested_force`、`delivery_id` をコピーする。 |
| response | `{ "message":"Build approved", "queued":true, "queue_id":"...", "dispatch":"requested\|timer_fallback" }`。`dispatch` は systemd 起動要求の受付結果であり、approval status や build 開始結果ではない。 |
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
| pending に同一 `approval_id` の queue entry が 1 件存在 | 前回 approve の部分失敗として既存 `queue_id` を再利用し、approved record だけを追記する。queue 上限は再判定しない。 |
| pending に同一 `approval_id` の queue entry が複数または payload 不一致で存在 | `500`。approval / queue を変更せず、server log `APPROVAL_QUEUE_INCONSISTENT` を出す。 |

SDK method、request 変換、error 伝播は [`docs/details/sdk.md` 詳細本文責務 §27.21〜§27.47](sdk.md#sec-27-21)、UI 表示、操作、再取得、disabled 判定は [`docs/details/ui.md` 詳細本文責務 §27.21〜§27.47](ui.md#sec-27-21) を正本とする。API は approval endpoint の request / response、状態遷移呼出境界、[approval 更新順](#api-approval-update-order) だけを担当する。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| list | id ごとの最新 record だけを返す。 |
| approve | queue 追加、trigger approval。 |
| reject | build なし、history 記録。 |
| duplicate approve | `409`、状態差分なし。 |
| queue full | `429`、approval は pending のまま。 |
| approve partial retry | 既存 queue entry を再利用し、queue 件数を増やさず approved record と audit を完了する。 |
| runner activation | approved record と audit の確定後だけ起動要求を行う。起動要求失敗は approved / queue を保持し、`dispatch:"timer_fallback"` で成功 response を返す。 |

`.approval_queue` record schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) `.approval_queue` JSON Lines schema を参照する。同一 id の最新 record を有効状態として扱い、古い record は監査履歴として残す。`GET /api/approvals` は id ごとに最新 record だけを返し、`created_at` 降順、同時刻は id 昇順で並べる。壊れた行は無視し、response に含めない。

**approval 状態遷移固定契約：**

| 現在 status | 操作 | 次 status | 副作用 |
|-------------|------|-----------|--------|
| `pending` | approve | `approved` | `.build_state.queued[]` に `trigger:"approval"` entry を追加し、`queue_id` を保存する。 |
| `pending` | reject | `rejected` | `.build_history` に `status:"approval_rejected"` を追記する。queue は追加しない。 |
| `approved` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `rejected` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |
| `expired` | approve / reject | 変更なし | `409 {"error":"Conflict"}`。 |

<a id="api-approval-update-order"></a>
**approval 更新順：**

| 操作 | 更新順 | 失敗時 |
|------|--------|--------|
| approve | `.approval_queue` lock → 最新 pending 確認 → `.build_state` lock → 同一 `approval_id` の queue entry 検索 → 存在しない場合だけ queue append → approved record append → `approval_approved` audit → runner 起動要求 → response | queue full は新規 queue append 時だけ `429`、approval は pending のまま。approved append または audit 失敗時は `500`、queue / approved record と先行 log は巻き戻さず起動要求を行わない。起動要求失敗だけは `dispatch:"timer_fallback"` で `200`。再試行は同一 queue entry を再利用する。 |
| reject | `.approval_queue` lock → 最新 pending 確認 → rejected record append → `.build_history` append → `approval_rejected` audit → response | history append または audit 失敗時は `500`。rejected record / history と先行 log は巻き戻さない。 |

**approval API 実装確認固定契約：**

| 項目 | 仕様 |
|------|------|
| list corrupt line | `.approval_queue` の破損行は response から除外し、server log に固定 code `APPROVAL_QUEUE_CORRUPT_LINE` と file path だけを出す。行本文、token、payload は出さない。 |
| approve queue id | approve で作成する queue id の base は `q{YYYYMMDDHHmmss}` とし、衝突処理は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約に従う。approval id とは別 id とする。 |
| approve partial failure | queue append 成功後に approved record append が失敗した場合、queue entry は残し、response `500`、server log `APPROVAL_APPROVED_RECORD_FAILED` を出す。runner は latest approval status が `approved` で同じ `queue_id` を持つことを確認できるまでその entry を取り出さない。同じ approve 再試行は既存 entry を再利用して approved record と audit を完了し、queue を重複追加しない。 |
| reject partial failure | rejected record append 成功後に history append が失敗した場合、rejected は残し、response `500`、server log `APPROVAL_REJECT_HISTORY_FAILED` を出す。 |
| audit | list は read audit 対象外。approve は `approval_approved`、reject は `approval_rejected` とし、`target_type:"approval"`、`target_id` は path の approval id とする。audit 失敗時は [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) に従い、保存済み状態を巻き戻さない。 |
| body 禁止 | `Content-Length > 0` または JSON body がある approve / reject は `400` と `{"error":"Request body is not allowed"}` を返し、状態差分なし。 |
| secret | approval payload、queue payload、history、audit、server log、SDK error、UI 表示に Authorization header、session token、API token、repository token を保存しない。 |
| collaborator 参照 | SDK の response / error 伝播は [`docs/details/sdk.md` 詳細本文責務 §27.21〜§27.47](sdk.md#sec-27-21)、UI の表示、再取得、自動再送禁止は [`docs/details/ui.md` 詳細本文責務 §27.21〜§27.47](ui.md#sec-27-21) を正本とする。API は collaborator 内部状態を定義しない。 |

<a id="sec-security-owner-api-boundary"></a>
**security owner 機能 API 側境界：**

security owner 機能に対し、[`docs/details/api.md`](api.md) 詳細本文責務は route、HTTP request / response、状態 read/write 呼び出し境界だけを持つ。owner component、security 主本文、漏えい禁止、認証・認可・監査・rate limit の判定本文は [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) を正本とする。

**API rate limit request / response 契約：**

`ApiRateLimitPolicyInput` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ApiRateLimitPolicy object と同じ `enabled`、`groups` だけを持つ。未知 key、response 専用 `state_summary`、group の不足または追加、型不一致、範囲外は `422` とし、状態を変更しない。

`ApiRateLimitPolicyResponse` は保存済みまたは既定の `enabled`、`groups` と、response 専用 `state_summary` を持つ。`groups` は `login`、`read`、`trigger`、`operate`、`config`、`admin` の順で生成する。`state_summary` は次の item を `reset_at` 降順、同時刻は `key` 昇順で最大 100 件返す。

| キー | 型 | 説明 |
|------|----|------|
| `key` | string | `.api_rate_state.windows` の key。 |
| `group` | string | key に含まれる endpoint group。 |
| `window_start` | string | 保存済み window 開始時刻。 |
| `count` | integer | 保存済み現在 count。 |
| `reset_at` | string | `window_start` に当該 group の `window_seconds` を加算した UTC ISO 8601。 |

response 算出時点で `reset_at` が現在時刻以下の window は `state_summary` から除外する。`GET /api/api-rate-limit` は除外した window を `.api_rate_state` から削除せず、endpoint 固有の write を行わない。`POST /api/api-rate-limit` で policy が変わった場合は [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の保存順で window を空にして `state_summary:[]` を返す。同値更新の場合は状態を変更せず、現在の有効 window から `state_summary` を算出する。

**security owner 機能の API 側共通境界：**

| security owner 機能 | api 側が持つ内容 | security 正本参照 |
|----|------------------|------------------|
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | route / method、endpoint dispatch、HTTP status、request / response body、状態ファイル read/write 呼び出し境界。 | scope 判定順、許可 endpoint group、permission denied audit、body parse 前判定、漏えい禁止値。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | token API の route、HTTP method、request validation 入口、response schema、`.api_tokens` read/write 呼び出し境界。 | token 生成、hash 保存、scope 検証、作成時 1 回だけ token 本体を返す契約、認証成功時の `last_used_at` 更新、token 漏えい禁止。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | audit log API の route、query parameter、response envelope、`.audit_log` read 呼び出し境界。 | action / actor / target / result の導出、必須 audit 失敗時の `500`、secret / token / request body 保存禁止、壊れた行の扱い。record schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | session timeout config API の route、request body、response body、`.server_config.session_timeout_seconds` read/write 呼び出し境界。 | session の作成、非 sliding の期限判定、`last_used_at` 更新、期限切れ時 `401`、既存 session への反映条件、監査順序。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | auth / TOTP API の route、request body、response body、`.totp_secret` read/write 呼び出し境界。 | TOTP secret 生成、setup 仮 secret、login ticket、code 検証、secret の一回表示、ticket 再利用禁止、TOTP 漏えい禁止、監査順序。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | rate limit config API の route、request body、response body、`.server_config.api_rate_limit` と `.api_rate_state` の read/write 呼び出し境界。 | endpoint group 判定、window / count 更新、`429` 時に count を増やさない契約、actor key / IP key の同一 lock 更新、rate limit audit。 |
