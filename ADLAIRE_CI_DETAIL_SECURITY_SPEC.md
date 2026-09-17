# Adlaire CI — Security 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` および `ADLAIRE_CI_DETAIL_API_SPEC.md` から分割した `security` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

本ファイルを読む前に、`ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `security` owner component の主本文であり、collaborator component の仕様は呼び出し境界、endpoint、SDK、UI、状態 schema、fixture、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `security` |
| collaborator component | `api`、`sdk`、`ui`、`statefile` |
| 持つ内容 | API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| 持たない内容 | API endpoint 共通処理、SDK method 実装、UI DOM 詳細、runner / builder の業務処理。 |

---

## 対象範囲

| 範囲 | 内容 |
|------|------|
| §27.42 | ビルドトリガー専用 API スコープ。 |
| §27.43 | API キー管理。 |
| §27.44 | 監査ログ。 |
| §27.45 | セッションタイムアウト変更設定。 |
| §27.46 | TOTP 二要素認証。 |
| §27.47 | API レート制限。 |

---

## 横断固定契約

**§27.42〜§27.47 認証・監査・制限機能 実装完全性固定契約：**

§27.42〜§27.47 は、API token、監査、session、TOTP、rate limit に関する安全機能である。各機能は個別節に加えて下表を満たした場合だけ実装完了とする。

| 節 | 機能 | 判定入口 | 成功時副作用 | 失敗時副作用 | 漏えい禁止値 | 必須 fixture |
|----|------|----------|--------------|--------------|--------------|--------------|
| §27.42 | API token scope | route / method 確定後、body parse 前。 | 許可 endpoint だけ処理し、§27.44 の監査対象 event に該当する場合は audit に actor を残す。 | 権限不足は対象処理を実行せず `403`。audit 失敗時は `500`。 | token 本体、Authorization header。 | trigger allowed、read denied、multi scope、path param、body 未評価、audit failure。 |
| §27.43 | API key 管理 | admin session または admin scope。 | token hash だけ保存し、作成時だけ token 本体を返す。 | validation 失敗は保存差分なし。失効済み token は再有効化しない。 | token 本体、token hash の不要露出。 | create、list mask、revoke、expired、duplicate label、admin token create。 |
| §27.44 | 監査ログ | security / config / operation event 確定時。 | 1 event 1 JSON Lines で追記し、actor / target / result を保存する。 | 必須 audit 失敗は対象処理を `500` にする。任意 audit は個別節優先。 | secret、password、token、TOTP secret、raw request body。 | success、denied、failure、mask、append failure、pagination。 |
| §27.45 | session timeout | login、authenticated request、timeout config API。 | session の last_seen / expires_at を固定規則で更新する。 | timeout session は `401`、対象 endpoint は実行しない。 | session token。 | active、expired、sliding update、config update、revoke all、clock boundary。 |
| §27.46 | TOTP | setup、confirm、login/totp、disable。 | secret は confirm 成功後だけ有効保存し、ticket は一回だけ使う。 | ticket 再利用、期限切れ、code 不正は対象状態を変更しない。 | TOTP secret、backup code 相当値、ticket token。 | setup、confirm、login success、code reuse、disable、audit failure。 |
| §27.47 | API rate limit | route / auth / scope 判定の定義済み位置。 | 上限未満だけ count を増やし endpoint 処理へ進む。 | `429` は count を増やさず、audit 成功時だけ返す。 | API token、session token、request body。 | login 11 回目、window reset、IP+actor、disabled、policy update、state save failure。 |

**§27.42〜§27.47 セキュリティ機能 横断順序固定契約：**

| 順序 | 処理 | 固定条件 |
|------|------|----------|
| 1 | route / method を確定する。 | 未定義 route は認証、rate limit、body parse より前に `404` / `405`。 |
| 2 | 認証不要 endpoint を判定する。 | `GET /api/health` と `POST /api/webhook` は個別契約を優先する。 |
| 3 | login rate limit を判定する。 | login group は認証前 IP key で判定する。 |
| 4 | 認証情報を検証する。 | session と API token を混同しない。形式不一致は `401`。 |
| 5 | API token scope を判定する。 | scope 不足時は body validation と状態更新を行わない。 |
| 6 | 認証後 rate limit を判定する。 | actor key と IP key を同一 lock 内で判定・更新する。 |
| 7 | endpoint 固有 validation を行う。 | 失敗時は対象状態、外部 API、外部 command を変更しない。 |
| 8 | endpoint 固有処理を実行する。 | 成功時だけ個別節の保存順で状態、access log、audit log を確定する。 |

上表の順序を変更してはならない。個別節が別順序を明記する場合は、セキュリティ上の漏えいを増やさない範囲で個別節を優先する。順序変更が必要な場合は、先に本表と該当個別節を同時に改訂する。

---

## 認証共通詳細

本節は、password 認証、session、login ticket、認証ログ、`--init-credentials` の security owner 詳細仕様である。HTTP endpoint の method、path、request、response、status は `ADLAIRE_CI_DETAIL_API_SPEC.md` §22.0e および §25 を正とする。`.admin_credentials` schema は `ADLAIRE_CI_DETAIL_STATEFILE_SPEC.md` §22.0c を正とする。

**password hash 固定契約：**

| 項目 | 仕様 |
|------|------|
| algorithm | Go 標準ライブラリのみで実装する `sha256_iter_v1`。 |
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| 初回 digest | `sha256(salt_bytes || password_utf8_bytes)`。 |
| 反復 | `iterations = 260000`。2 回目以降は `sha256(previous_digest || salt_bytes || password_utf8_bytes)` を繰り返す。 |
| 保存値 | 最終 digest を lowercase hex 文字列で `.admin_credentials.password_hash` に保存する。 |
| 比較 | 入力 password から同一手順で digest を生成し、`crypto/subtle.ConstantTimeCompare` で比較する。 |
| 外部依存 | `golang.org/x/crypto/pbkdf2` 等の外部パッケージは使用しない。 |

**password 入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字。 |
| 最大長 | 128 文字。 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期 password | `--init-credentials` で `.admin_credentials` に生成する。初回 login 時は `must_change:"prompt"` を返す。 |
| 変更時検証 | `new_password` が現在 password と同一の場合は `422`。 |
| 失敗時応答 | password 不一致は `401 {"error":"Unauthorized"}`。どの条件に失敗したかは返さない。 |
| 成功時保存 | `.admin_credentials` に新 salt、新 hash、`must_change:false`、`updated_at` を原子的に保存する。 |

**session / ticket 固定契約：**

| 項目 | 仕様 |
|------|------|
| session token 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。 |
| session 保存場所 | `api` 内のインメモリ辞書。再起動で全 session を破棄する。 |
| session key | token 本体ではなく `sha256(token)` の lowercase hex。 |
| session 期限 | 新規発行時点の `.server_config.session_timeout_seconds`。設定不在時は 8 時間。 |
| session_id | `crypto/rand` で 16 bytes を生成し、lowercase hex とする。 |
| login ticket | TOTP 有効 login の password 成功時だけ生成し、メモリ上で hash 化して保持する。 |
| 永続化禁止 | session token、login ticket、setup 仮 secret、連続失敗回数はファイル保存しない。 |
| response 禁止 | password hash、salt、session token hash、ticket hash、TOTP secret 保存値は response に含めない。 |
| log 禁止 | password、current_password、new_password、token、ticket、hash、salt、TOTP code は `.access_log`、`.audit_log`、`.api_access_log`、journal に含めない。 |

`GET /api/sessions` は `session_id` ではなく `current`、`created_at`、`expires_at`、`last_used_at` のみ返す。認証必須 API で有効 token を受信した場合、`last_used_at` を現在時刻へ更新する。期限切れ token は検出時にメモリから削除する。`POST /api/logout` は対象 token のみ削除する。`POST /api/sessions/revoke-all` は現在 token 以外を削除する。

**login 失敗制御：**

| 項目 | 仕様 |
|------|------|
| 失敗記録 | `api` はメモリ上で直近の連続 login 失敗回数と最終失敗時刻を保持する。再起動で失敗回数はリセットされる。 |
| lock 条件 | 連続 10 回失敗した場合、最終失敗から 10 分間 `POST /api/login` を `429 {"error":"Too many attempts"}` で拒否する。 |
| 成功時 | login 成功時は連続失敗回数を 0 に戻す。 |
| 応答時間 | password 不一致、存在しない credentials、lock 中を除く検証失敗では、条件の詳細を response へ出さない。 |
| log | 成功、失敗、lock 拒否はいずれも `.access_log` へ追記する。password、token、hash、salt は記録しない。 |

**認証処理の副作用境界：**

| ケース | HTTP status | `.admin_credentials` | メモリ session / ticket | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|----------------------|-------------------------|---------------|--------------|------|
| password 不一致 | `401` | 変更なし | 変更なし | `login_failure` を追記 | `login_failure` を追記 | 連続失敗回数はメモリ上で +1。 |
| 連続失敗 lock | `429` | 変更なし | 変更なし | `login_locked` を追記 | `permission_denied` を追記 | password hash 検証は実行しない。 |
| password 成功 / TOTP 無効 | `200` | `login_count`、`last_login_at` 更新 | session 追加 | `login_success` を追記 | `login_success` を追記 | session token は全永続ログに保存しない。 |
| password 成功 / TOTP 有効 | `200` | 変更なし | ticket 追加 | `totp_required` を追記 | `totp_required` を追記 | `login_count` と `last_login_at` は更新しない。 |
| TOTP code 不一致 | `401` | 変更なし | ticket 削除 | `totp_failure` を追記 | `totp_failure` を追記 | ticket は再利用不可。 |
| TOTP 成功 | `200` | `login_count`、`last_login_at` 更新 | ticket 削除、session 追加 | `login_success` を追記 | `login_success` を追記 | session 期限は成功時点の設定で決める。 |
| session 期限切れ | `401` | 変更なし | 対象 session 削除 | 追記しない | 追記しない | `.api_access_log` は通常 API request として記録する。 |
| logout | `200` | 変更なし | 対象 session 削除 | `logout` を追記 | `logout` を追記 | token 本体と token hash は保存しない。 |
| password 変更成功 | `200` | 新 salt / hash、`must_change:false`、`updated_at` 更新 | 現 session 以外削除 | `password_change` を追記 | `password_change` を追記 | 新旧 password、hash、salt は保存しない。 |

上表で `.access_log` または `.audit_log` の追記が必要な処理は、response 返却前に追記を完了する。追記失敗時は、個別節で別指定がない限り `500 {"error":"Internal server error"}` を返す。session token、login ticket、TOTP setup secret は、必要なログ追記がすべて成功するまで response に含めてはならない。

**認証副作用順序固定契約：**

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

**`--init-credentials` CLI 固定契約：**

| ケース | stdout | stderr | 終了コード | 副作用 |
|--------|--------|--------|------------|--------|
| 新規生成成功 | `credentials initialized` + LF | 空 | `0` | `.admin_credentials` を mode `0600` で作成する。 |
| 既存あり | 空 | `credentials already exist` + LF | `2` | 既存ファイルを変更しない。 |
| `--state-dir` 相対 path | 空 | `state directory must be absolute: {path}` + LF | `2` | ファイル作成なし。 |
| 書込失敗 | 空 | `credentials write failed` + LF | `1` | tmp を削除し、部分ファイルを残さない。 |
| rand 失敗 | 空 | `random source failed` + LF | `1` | ファイル作成なし。 |

生成手順は、state dir 検証 → 既存確認 → salt 生成 → hash 生成 → `{path}.tmp.{pid}` へ JSON + LF 書込 → mode `0600` → file sync → rename → parent directory sync の順に固定する。rename 後の sync に失敗した場合は `1` を返し、作成済みファイルは残る。実装者判断で初期 password を環境変数、対話入力、ランダム生成へ変更してはならない。

---

### 27.42 ビルドトリガー専用 API スコープ
owner component は `security` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。


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
owner component は `security` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、API key の発行、一覧、失効、期限、scope を実装し、key 本体を保存しないことである。

**API 権限：**

| endpoint | 必要認証 |
|----------|----------|
| `GET /api/tokens` | 管理 session または `admin` scope API token |
| `POST /api/tokens` | 管理 session または `admin` scope API token |
| `DELETE /api/tokens/{id}` | 管理 session または `admin` scope API token |

`admin` scope API token は新しい `admin` scope API token を発行できる。発行者 token と発行対象 token は別 record とし、親子関係は保存しない。

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

認証に使用中の API token 自身を失効対象にできる。その場合、当該リクエストは成功し、次リクエストから `401` になる。

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
owner component は `security` とする。collaborator component は `api`、`statefile` とする。


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
owner component は `security` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。


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
owner component は `security` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。


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
| setup 仮 secret | `api` のメモリ | 10 分 | `secret_base32`, `created_at`。サーバー再起動で破棄する。 |
| login ticket | `api` のメモリ | 5 分 | `ticket_hash`, `created_at`, `password_verified_at`。ticket 本体は hash 化して保持する。 |

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
owner component は `security` とする。collaborator component は `api`、`sdk`、`ui`、`statefile` とする。


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

`state_summary` は `reset_at` 降順、同時刻は `key` 昇順で最大 100 件返す。期限切れ window は response 算出前に `.api_rate_state` から削除する。

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
