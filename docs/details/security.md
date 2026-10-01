# Adlaire CI — Security 詳細仕様

---

<a id="0-責務境界"></a>

**0. 責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `security` |
| 実装主体 | 単独の Go artifact は持たない。認証、認可、token、session、TOTP、rate limit、audit の実装は [`components/api.go`](../../components/api.go) に内包する。 |
| 持つ内容 | `security` owner が主本文として定義する API token scope、API key、audit、session timeout、TOTP、rate limit、漏えい禁止、security 横断順序。 |
| 検証接続 | 検証方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、owner から fixture への入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b](../DETAIL_INDEX.md#0b-詳細仕様参照表)、Phase 11 横断入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry)、fixture / expected / fake / 実行証跡 / PR 証跡 package / クロージャ記録 / mutation / contract drift の記録先は [`docs/details/fixture.md` fixture 証跡責務](fixture.md) と [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](../DOCUMENT_INDEX.md) を参照する。この owner 詳細本文では owner 固有の入出力、状態、処理順序、異常系、検証条件だけを記載する。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | ビルドトリガー専用 API スコープ。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | API キー管理。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 監査ログ。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | セッションタイムアウト変更設定。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP 二要素認証。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | API レート制限。 |

---

<a id="横断固定契約"></a>
**横断固定契約：**

<a id="sec-27-42"></a>
**[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 認証・監査・制限機能 実装確認固定契約](security.md#sec-27-42)：**

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) は、API token、監査、session、TOTP、rate limit に関する安全機能である。各機能の詳細実装確認では、対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約に加えて、この認証・監査・制限機能実装確認固定表を満たす。

| 節 | 機能 | 判定入口 | 成功時副作用 | 失敗時副作用 | 漏えい禁止値 | fixture 証跡入口 |
|----|------|----------|--------------|--------------|--------------|----------------|
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) | API token scope | route / method 確定後、body parse 前。 | 許可 endpoint だけ処理し、[`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の監査対象 event に該当する場合は audit に actor を残す。 | 権限不足は対象処理を実行せず `403`。audit 失敗時は `500`。 | token 本体、Authorization header。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) | API key 管理 | admin session または admin scope。 | token hash だけ保存し、作成時だけ token 本体を返す。 | validation 失敗は保存差分なし。失効済み token は再有効化しない。 | token 本体、token hash の不要露出。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) | 監査ログ | security / config / operation event 確定時。 | 1 event 1 JSON Lines で追記し、actor / target / result を保存する。 | 対象 action の audit 失敗は API では `500`、runner 内部 event では runner failure とする。 | secret、password、token、TOTP secret、raw request body。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) | session timeout | login、authenticated request、timeout config API。 | session 発行時に `expires_at` を固定し、有効な認証 request ごとに `last_used_at` だけを更新する。 | timeout session は `401`、対象 endpoint は実行しない。 | session token。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) | TOTP | setup、confirm、login/totp、disable。 | secret は confirm 成功後だけ有効保存し、ticket は一回だけ使う。 | ticket 再利用、期限切れ、code 不正は対象状態を変更しない。 | TOTP secret、backup code 相当値、ticket token。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |
| [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) | API rate limit | route / auth / scope 判定の定義済み位置。 | 上限未満だけ count を増やし endpoint 処理へ進む。 | `429` は count を増やさず、audit 成功時だけ返す。 | API token、session token、request body。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)。 |

<a id="sec-27-42-2"></a>
**[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2)：**

| 順序 | 処理 | 固定条件 |
|------|------|----------|
| 1 | route / method を確定する。 | 未定義 route は認証、rate limit、body parse より前に `404` / `405`。 |
| 2 | IP access control を判定する。 | `GET /api/health` だけは省略する。読取、拒否、破損時の副作用は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。 |
| 3 | 認証不要 endpoint を判定する。 | `GET /api/health` と `POST /api/webhook` は個別契約を優先する。Webhook は [`docs/details/api.md` 詳細本文責務 §27.12](api.md#sec-27-12) の body size 判定と署名検証を認証代替とする。 |
| 4 | login rate limit を判定する。 | login group は credential body の parse と認証より前に IP key で判定する。login 以外ではこの段階を省略する。 |
| 5 | 認証情報を検証する。 | 認証必須 endpoint は session と API token を混同しない。形式不一致は `401`。login は検証済み body 上限内で credential body を parse し、個別認証契約を適用する。 |
| 6 | 認可 gate を判定する。 | API token は endpoint scope を判定し、scope 不足は `403 {"error":"Forbidden"}` とする。管理 session は session 発行時の `password_change_required` を判定し、`true` の場合は `POST /api/change-password` と `POST /api/logout` 以外を `403 {"error":"Password change required"}` とする。拒否時は endpoint 固有 body の parse / validation、認証後 rate limit 更新、endpoint 固有状態更新を行わない。認証不要 endpoint では省略する。 |
| 7 | 認証後 rate limit を判定する。 | actor key と IP key を同一 lock 内で判定・更新する。`POST /api/webhook` は署名検証成功後に `trigger` group の IP key だけを判定・更新する。その他の認証不要 endpoint は個別契約で対象 group が明記された場合だけ実行する。 |
| 8 | endpoint 固有 body / query / path validation を行う。 | 失敗時は対象状態、外部 API、外部 command を変更しない。前段階までに完了した security 副作用は保持する。 |
| 9 | endpoint 固有処理を実行する。 | 成功時だけ [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に固定された保存順で状態、access log、audit log を確定する。 |

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) 認証・監査・制限機能実装確認固定表の順序を変更してはならない。対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約はこの順序を上書きできず、異なる順序の記載がある場合は仕様不整合として実装しない。例外が必要な場合は、先に横断順序固定表と対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能契約を同一仕様変更で改訂し、例外の endpoint、挿入位置、前後の副作用、失敗時の応答を固定表に明記する。その改訂が完了するまで例外は存在しない。

---

<a id="認証共通詳細"></a>
**認証共通詳細：**

[`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) は、password 認証、session、login ticket、認証ログ、`--init-credentials` の [`docs/details/security.md` 詳細本文責務](security.md) である。HTTP endpoint の method、path、request、response、status は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) および [`docs/details/api.md` 詳細本文責務 §25](api.md#25-認証-実装仕様) を参照する。`.admin_credentials` schema は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) を参照する。

**password hash 固定契約：**

| 項目 | 仕様 |
|------|------|
| algorithm | Go 標準ライブラリのみで実装する `sha256_iter_v1`。 |
| salt 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の hex 文字列として保存する。 |
| 初回 digest | `salt_bytes` と `password_utf8_bytes` をこの順で byte 連結し、SHA-256 digest を生成する。 |
| 反復 | `iterations = 260000`。2 回目以降は `previous_digest`、`salt_bytes`、`password_utf8_bytes` をこの順で byte 連結して SHA-256 digest を生成する処理を繰り返す。 |
| 保存値 | 最終 digest を lowercase hex 文字列で `.admin_credentials.password_hash` に保存する。 |
| 比較 | 入力 password から同一手順で digest を生成し、`crypto/subtle.ConstantTimeCompare` で比較する。 |
| 実装境界 | `crypto/sha256`、`crypto/rand`、`crypto/subtle`、`encoding/hex` だけで成立する。外部依存の採否は [`docs/SPEC.md` ポリシー責務 §4](../SPEC.md#policy-dependencies) を参照する。 |

**password 入力制約：**

| 項目 | 仕様 |
|------|------|
| 最小長 | 8 文字。 |
| 最大長 | 128 文字。 |
| 許可文字 | UTF-8 文字列。NUL 文字は禁止。前後空白はトリムせず、入力値そのものを検証・ハッシュ化する。 |
| 初期 password | `--init-credentials` の非 terminal 標準入力から取得し、平文を保存または出力せず `.admin_credentials` に hash、`must_change:true`、`login_count:0` を生成する。session 発行時の response は session / ticket 固定契約の `must_change` 算出を使用する。 |
| 変更時検証 | `new_password` が現在 password と同一の場合は `422`。 |
| 失敗時応答 | password 不一致は `401 {"error":"Unauthorized"}`。どの条件に失敗したかは返さない。 |
| 成功時保存 | `.admin_credentials` に新 salt、新 hash、`must_change:false`、`updated_at` を原子的に保存する。 |

**session / ticket 固定契約：**

| 項目 | 仕様 |
|------|------|
| session token 生成 | `crypto/rand` で 32 bytes を生成し、`encoding/hex` で 64 文字の lowercase hex 文字列へ変換する。 |
| session 保存場所 | `api` 内のインメモリ辞書。再起動で全 session を破棄する。 |
| session key | token 本体ではなく `sha256(token)` の lowercase hex。 |
| session 期限 | 新規発行時点の `.server_config.session_timeout_seconds`。設定不在時は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の既定値を使用する。 |
| session_id | `crypto/rand` で 16 bytes を生成し、lowercase hex とする。 |
| `must_change` 算出 | session token 発行と同じ credentials 更新で `login_count` を飽和加算した後に算出する。`.admin_credentials.must_change=false` は `"none"`、`true` かつ増加後 `login_count` が 1〜4 は `"prompt"`、`true` かつ増加後 `login_count` が 5 以上は `"forced"` とする。 |
| password 変更 gate | session record は発行 response の `must_change` が `"forced"` の場合だけ `password_change_required:true` をメモリ上に保持する。password 変更成功時に現 session の値を `false` にし、その他の session は失効する。 |
| login ticket | TOTP 有効 login の password 成功時だけ生成し、メモリ上で hash 化して保持する。 |
| login ticket 保持値 | password 検証済みであること、ticket hash、発行時刻、期限、credentials fingerprint だけを保持する。`must_change`、`login_count`、`last_login_at`、`password_change_required` は保持しない。 |
| credentials fingerprint | ticket 発行時の `password_hash`、`salt`、`algorithm`、`iterations` の 10 進 ASCII、`updated_at` をこの順で LF 1 byte で連結した UTF-8 byte 列の SHA-256 lowercase hex。TOTP 成功処理で最新 credentials から再算出し、`crypto/subtle.ConstantTimeCompare` で一致しない場合は ticket を消費して `401 {"error":"Unauthorized"}` とする。fingerprint は response と log に出力しない。 |
| 永続化禁止 | session token、login ticket、setup 仮 secret、連続失敗回数はファイル保存しない。 |
| response 禁止 | password hash、salt、session token hash、ticket hash、TOTP secret 保存値は response に含めない。 |
| log 禁止 | password、current_password、new_password、token、ticket、hash、salt、TOTP code は `.access_log`、`.audit_log`、`.api_access_log`、journal に含めない。 |

`GET /api/sessions` は `session_id` ではなく `current`、`created_at`、`expires_at`、`last_used_at` のみ返す。認証必須 API で有効 token を受信した場合、`last_used_at` を現在時刻へ更新する。期限切れ token は検出時にメモリから削除する。`POST /api/logout` は対象 token のみ削除する。`POST /api/sessions/revoke-all` は現在 token 以外を削除する。

**login 失敗制御：**

| 項目 | 仕様 |
|------|------|
| key | [`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の IP key と同じ正規化を使用する。`RemoteAddr` を `net.SplitHostPort` で host に正規化し、parse 失敗時は同節の固定 fallback key を使用する。 |
| 失敗記録 | `POST /api/login` の password hash 不一致だけが、IP key ごとの連続失敗回数と最終失敗時刻をメモリ上で更新する。JSON / body 検証失敗、pre-auth rate limit 拒否、TOTP 失敗はこの回数を更新しない。再起動で全 entry を破棄する。 |
| lock 成立 | 10 回目の連続 password 不一致は `401 {"error":"Unauthorized"}` を返し、`locked_until = last_failure_at + 10分` を記録する。その後 `now < locked_until` の request は password hash 検証を行わず `429 {"error":"Too many attempts"}` で拒否する。 |
| lock 解除 | `now >= locked_until` であれば対象 entry を削除してから password hash 検証を行う。password hash 一致時も、TOTP の有効・無効にかかわらず対象 entry を削除する。 |
| rate limit 優先順位 | pre-auth login rate limit を login lock より前に判定する。両方が拒否条件を満たす場合は `429 {"error":"Too many requests"}` を返し、login lock 判定、password hash 検証、失敗回数更新を行わない。 |
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
| password 変更成功 | `200` | 新 salt / hash、`must_change:false`、`updated_at` 更新 | 現 session の `password_change_required=false`、現 session 以外削除 | `password_change` を追記 | `password_change` を追記 | 新旧 password、hash、salt は保存しない。 |
| 他 session 一括失効 | `200` | 変更なし | 現 session 以外削除 | `session_revoke_all` を追記 | `session_revoke_all` を追記 | 現 session と実行中 response は維持する。 |

[`docs/details/security.md` 詳細本文責務 認証共通詳細](security.md#認証共通詳細) の固定表で `.access_log` または `.audit_log` の追記が必要な処理は、response 返却前に追記を完了する。追記失敗時は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint または対象の [§27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) 機能契約で別の応答が定義されていない限り `500 {"error":"Internal server error"}` を返す。session token、login ticket、TOTP setup secret は、必要なログ追記がすべて成功するまで response に含めてはならない。

**認証副作用順序固定契約：**

| 処理 | 固定順序 | 失敗時 |
|------|----------|--------|
| `POST /api/login` password 不一致 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials 読込 → hash 比較 → 失敗回数更新 → `.access_log` → `.audit_log` → `401` response | 10 回目も `401`。ログ追記失敗は `500`。失敗回数は戻さない。session/ticket は作成しない。 |
| `POST /api/login` password 成功 / TOTP 無効 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials lock 取得 → credentials 読込 → hash 比較 → IP key の失敗 entry 削除 → session token 生成 → `login_count=min(login_count+1,9223372036854775807)` と `last_login_at` 更新 → `must_change` 算出 → credentials 保存 → `.access_log` → `.audit_log` →同値の `password_change_required` を持つ session record 追加 → token response | token 生成失敗は credentials 更新前の `500`。credentials またはログ追記失敗も `500` とし、token を response せず session record を追加しない。ログ失敗前に保存済みの credentials は戻さない。 |
| `POST /api/login` password 成功 / TOTP 有効 | pre-auth rate limit → IP key の lock entry 期限判定 → credentials 読込 → hash 比較 → IP key の失敗 entry 削除 → password 検証済み ticket 生成 → `.access_log` → `.audit_log` → ticket response | `.admin_credentials`、`login_count`、`last_login_at` は更新せず、ticket に `must_change` を保存しない。ログ追記失敗は `500`。ticket は保存しない。 |
| `POST /api/login/totp` 成功 | ticket 検証・即時無効化 → TOTP secret 読込 → code 検証 → credentials lock 取得 → credentials 再読込 → credentials fingerprint 一致判定 → session token 生成 → `.totp_secret.last_accepted_step` 更新 → `login_count=min(login_count+1,9223372036854775807)` と `last_login_at` 更新 → `must_change` 算出 → credentials 保存 → `.access_log` → `.audit_log` →同値の `password_change_required` を持つ session record 追加 → token response | fingerprint 不一致は状態更新前の `401`。その他の途中失敗時は token を response せず session record を追加しない。ticket は成功/失敗いずれも再利用不可とし、ticket 発行時の credentials 値から `must_change` を推測しない。失敗地点より前に永続化済みの `.totp_secret` または credentials は戻さない。 |
| `POST /api/logout` | token 認証 → 対象 session 削除 → `.access_log` → `.audit_log` → response | ログ追記失敗は `500`。削除済み session は戻さない。 |
| `POST /api/change-password` | token 認証 → current password 検証 → new password 検証 → salt/hash 生成 → `.admin_credentials` 更新 → 現 session の `password_change_required=false` → 現 session 以外削除 → `.access_log` → `.audit_log` → response | credentials 更新失敗は session を変更しない。ログ追記失敗時は `500` だが更新済み credentials と session 変更は戻さない。 |
| `POST /api/sessions/revoke-all` | token 認証 → 現 session 以外を削除 → `.access_log` → `.audit_log` → response | ログ追記失敗時は `500`。削除済み session は戻さない。 |

<a id="security-auth-concurrency-contract"></a>
**認証 transaction・並行更新固定契約：**

security owner は、credentials、TOTP、session、login ticket、TOTP setup 仮 secret、login 失敗 entry をまたぐ認証状態遷移を process 内の単一 auth transaction coordinator で直列化する。対象は `POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password`、`GET /api/sessions`、`POST /api/sessions/revoke-all`、`GET /api/auth/totp-status`、`POST /api/auth/totp-setup`、`POST /api/auth/totp-confirm`、`DELETE /api/auth/totp`、session 認証時の期限切れ削除と `last_used_at` 更新である。同じ request から coordinator を再帰取得してはならない。

auth transaction coordinator の取得待ちは最大 10 秒とする。待機開始を elapsed `0` とし、elapsed `<10s` でだけ取得成功を許可する。elapsed `10s` 到達時は timeout を優先し、同時刻の解放を取得成功扱いにしない。timeout 時は endpoint 固有 body parse、credentials / TOTP read、entropy 取得、memory state 変更、状態 write、`.access_log`、`.audit_log` を開始せず、`409 {"error":"Conflict"}` を返す。request context が取得前に cancel された場合も副作用なしで終了する。取得後は、対象処理の認証副作用順序固定契約が完了するまで coordinator を保持し、response 値を確定した後に解放する。

auth transaction coordinator は、statefile locked update adapter 内の file lock と memory auth lock のどちらよりも先に取得し、どちらを解放した後にも最後まで保持する。file lock と memory auth lock を同時に保持してはならず、一方を保持中に他方を取得してはならない。同じ transaction で両方が必要な場合は、memory auth lock で必要値を read-copy して解放 → statefile locked update adapter を完了して file lock を解放 → memory auth lock で確定済み mutation を適用して解放、の順に固定する。coordinator が他 request の認証状態遷移を排他するため、各区間の間に同じ auth 状態を別 request が変更することはない。filesystem、network、command、log write、password hash、TOTP 計算、entropy 読取は memory auth lock 外で実行する。memory auth lock 内では session、ticket、pending TOTP、login failure / rate entry の read-copy または確定済み mutation だけを行う。auth transaction coordinator は statefile I/O と log write をまたいで保持できるが、別 request の response writer、network client、systemd command を待ってはならない。

| 競合処理 | 直列化後の固定結果 |
|----------|--------------------|
| 同時 password login 成功 | 先に coordinator を取得した request から 1 件ずつ最新 credentials を読み、各成功ごとに `login_count` を飽和加算する。lost update を禁止する。 |
| login と password change | password change の credentials rename が先なら旧 password login は `401`。login の session 発行が先なら password change はその session を使用でき、変更成功時に他 session を失効する。 |
| TOTP ticket と credentials 変更 | ticket の credentials fingerprint と成功時再読込値が一致しない場合は `401`。旧 revision の ticket から session を発行しない。 |
| 同時 TOTP setup | 後に coordinator を取得した setup が pending secret を置換する。confirm は request 開始時ではなく coordinator 取得後の最新 pending secret だけを検証し、置換前 secret の code は `409`。 |
| TOTP confirm と disable | 先に完了した transaction の永続状態を後続 request が再読込する。後続 request は stale memory snapshot を使用しない。 |
| revoke-all と通常認証 | revoke-all が先なら削除対象 token の後続認証は `401`。通常認証が先ならその request だけ完了でき、revoke-all 完了後の次 request は `401`。 |
| session 期限切れと logout | 最初の transaction が対象 session を削除し、後続 request は `401`。同じ session を二重削除、二重 audit しない。 |

credentials、TOTP、API token の read-modify-write は [`docs/details/statefile.md` 詳細本文責務 lock 内 read-modify-write adapter 固定契約](statefile.md#statefile-locked-update-adapter-contract) を使用する。transaction 開始前の read-only snapshot に基づく保存、file lock 外の `login_count` 加算、TOTP replay step 更新、token `last_used_at` 更新を禁止する。

session token と login ticket は `crypto/rand` 成功後にだけ生成し、生成した値はメモリ上で hash 化して保持する。response body に含める token / ticket は、その request の成功 response 1 回だけに含める。`403`、`429`、`500`、network 切断検出時に、未送信 token をログや状態ファイルへ退避してはならない。

<a id="init-credentials-cli-contract"></a>
**`--init-credentials` CLI 固定契約：**

初期 password は非 terminal の標準入力から exactly 1 行で受け取る。入力 byte 列は UTF-8 password byte 列と末尾 LF 1 byte だけで構成し、末尾 LF より後の byte、途中の LF、CR、NUL、UTF-8 不正を禁止する。末尾 LF は password に含めず、前後空白を trim しない。UTF-8 code point 数は 8〜128 とする。実装は最大 514 bytes まで読み、513 bytes 以下で exactly 1 個の末尾 LFを確認してから UTF-8 と code point 数を検証する。標準入力が terminal の場合は password を読み取らない。password byte 列は hash 入力以外に複製、保存、log 出力せず、処理終了前に保持 buffer を上書きする。

`--state-dir` の検証は、security 処理を開始する前に [`docs/details/api.md` 詳細本文責務 `api` CLI 固定契約](api.md#api-cli-contract) と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI state directory 共通固定契約](../DETAIL_INDEX.md#common-state-dir-contract) に従って完了させる。security owner は state directory の parse、検証順、エラー出力を再定義しない。

| init-credentials CLI ケース | stdout | stderr | 終了コード | 副作用 |
|----------------------------|--------|--------|------------|--------|
| 新規生成成功 | `credentials initialized` + LF | 空 | `0` | `.admin_credentials` を mode `0600` で作成する。 |
| 既存あり | 空 | `credentials already exist` + LF | `2` | 既存ファイルを変更しない。 |
| 標準入力が terminal | 空 | `initial password stdin must not be terminal` + LF | `2` | password 読取とファイル作成なし。 |
| password 入力が形式不正 | 空 | `invalid initial password` + LF | `2` | salt / hash 生成とファイル作成なし。 |
| password 読取失敗 | 空 | `password input failed` + LF | `1` | salt / hash 生成とファイル作成なし。 |
| 書込失敗 | 空 | `credentials write failed` + LF | `1` | tmp を削除し、部分ファイルを残さない。 |
| rand 失敗 | 空 | `random source failed` + LF | `1` | ファイル作成なし。 |

生成手順は、API CLI の state dir 検証 → `.admin_credentials` の予備存在確認 → 標準入力の terminal 判定 → password 読取・検証 → salt 生成 → hash 生成 → [`docs/details/statefile.md` 詳細本文責務 §22.0a 状態ファイル更新手順](statefile.md#statefile-update-procedure) の create-only mode 呼出し、の順に固定する。予備確認で既存 file を検出した場合は標準入力を読まず「既存あり」で終了する。create-only mode は lock 取得後に存在を再確認するため、予備確認後の競合でも既存 `.admin_credentials` を上書きしない。`ErrStateAlreadyExists` は「既存あり」の終了コード `2`、lock / tmp / write / chmod / file sync / rename / parent sync / cleanup failure は「書込失敗」の終了コード `1` へ写像する。rename 後の partial failure では作成済みファイルを残す。security owner は tmp 名、lock 名、rename、cleanup を再定義しない。初期 password を process 引数、環境変数、設定ファイル、固定値、terminal 対話、ランダム生成、stdout / stderr 出力から取得してはならない。

**認証共通実装確認ゲート：**

| 観点 | 合格条件 | 禁止条件 |
|------|----------|----------|
| one-time response | session token、login ticket、API token 本体、TOTP setup secret、otpauth URI は、該当成功 response 1 回だけに含める。 | `500`、`401`、`403`、`409`、`422`、`429` response、log、状態ファイル、fixture expected へ平文を残すこと。 |
| memory-only state | session、login ticket、TOTP setup 仮 secret、login 失敗回数は process memory だけに保持し、再起動で破棄する。 | session、login ticket、TOTP setup 仮 secret、login 失敗回数を保存する未定義永続ファイルの作成、ticket / session / 仮 secret の backup / restore 対象化。 |
| required security log before token | token / ticket / secret を response に含める前に、endpoint 固有契約が必須とする `.access_log` と `.audit_log` の追記を完了する。`.api_access_log` は [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の best-effort 記録とし、one-time 値の返却 gate にしない。 | 必須 `.access_log` または `.audit_log` の追記失敗時に token / ticket / secret を response へ含めること。 |
| hash-only storage | password、session token、login ticket、API token は保存時に hash 化し、平文を保存しない。 | hash 算出入力の平文、token 本体、ticket 本体、password 本体を expected / log に保存すること。 |
| fixed error body | 認証失敗、権限不足、rate limit、validation failure は固定 error body だけを返す。 | password 不一致理由、token record 詳細、scope 一覧、TOTP step、rate limit key を response に出すこと。 |
| no endpoint side effect | 認証、scope、rate limit、body validation のいずれかで失敗した場合、endpoint 固有処理を開始しない。 | 状態ファイル更新、外部 API / command 実行、通知送信、snapshot / archive / deploy 操作。 |
| audit dependency | 必須 audit が定義された操作は、audit 追記成功まで完了 response を返さない。 | 必須 audit 失敗時に成功 response を返すこと、audit failure を audit に再帰記録すること。 |
| secret scan | response、stdout、stderr、journal、`.access_log`、`.audit_log`、`.api_access_log`、`.config_log`、fixture expected に禁止値が存在しないことを検証する。 | 目視確認だけで secret 非表示を合格扱いにすること。 |

**one-time response 固定契約：**

| 値 | 返却 endpoint | 保存先 | 成功 response 以外の扱い |
|----|---------------|--------|--------------------------|
| session token | `POST /api/login`、`POST /api/login/totp` | memory に `sha256(token)`。 | response 前の必須 `.access_log` / `.audit_log` 失敗時は token を破棄し、`500`。 |
| login ticket | TOTP 有効時の `POST /api/login` | memory に `sha256(ticket)`。 | response 前の必須 `.access_log` / `.audit_log` 失敗時は ticket を破棄し、`500`。期限切れ / 失敗後は再利用不可。 |
| API token 本体 | `POST /api/tokens` | `.api_tokens.token_hash` のみ。 | `.api_tokens` 保存後の必須 `.access_log` / `.audit_log` 失敗時は record を残し、token 本体は返さず `500`。 |
| TOTP setup secret | `POST /api/auth/totp-setup` | memory の仮 secret。confirm 成功後だけ `.totp_secret`。 | setup response 以外へ出さない。audit / access / server log / status response へ出さない。 |
| otpauth URI | `POST /api/auth/totp-setup` | 永続保存しない。 | UI 一回表示以外へ出さない。fixture expected では secret 部分を `***` として扱う。 |

認証共通 fixture 名、入力、expected、合格条件、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を唯一の正本とする。[`docs/details/security.md`](security.md) 詳細本文責務では、one-time response、memory-only state、hash-only storage、audit dependency、secret scan の実装契約だけを扱う。

---

<a id="sec-27-42-3"></a>
**27.42 ビルドトリガー専用 API スコープ：**
[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42-3) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、外部システムによる build 開始操作を `trigger` scope の API token と build 開始 endpoint だけに限定することである。

<a id="api-token-scope-contract"></a>
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

管理 session は API token scope 固定表の対象外とするが、[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) の `password_change_required` 判定を常に適用する。API token が複数 scope を持つ場合は、いずれか 1 つの scope が endpoint に一致すれば許可する。`POST /api/login`、`POST /api/login/totp`、`POST /api/logout`、`POST /api/change-password` は API token scope 判定の対象外とし、API token では使用できない。`POST /api/webhook` は GitHub Webhook secret 検証専用であり、API token では使用できない。[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の scope 固定表に存在しない endpoint は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の scope 固定表へ追加されるまで API token では許可してはならない。

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

認証、API token scope、login rate limit、認証後 rate limit、endpoint 固有 validation の実行順は [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) を唯一の正本とする。scope 判定前に endpoint 固有の request body parse、状態ファイル更新、外部送信を実行してはならない。

**scope 判定固定条件：**

| 条件 | 判定 |
|------|------|
| session token | API token scope 表は参照せず、[`docs/details/security.md` 詳細本文責務 §27.42〜§27.47 セキュリティ機能 横断順序固定契約](security.md#sec-27-42-2) の `password_change_required` 判定後に許可する。 |
| API token の `scopes` が空 | token record 破損として扱い `500`。 |
| API token の `scopes` に未知値 | token record 破損として扱い `500`。 |
| endpoint が複数 group に現れる | 仕様不整合とし、group の優先順を推測せず実装しない。[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の endpoint 割当ては method と正規化 path pattern の組で正確に 1 group だけに属しなければならない。 |
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

<a id="sec-27-43"></a>
**27.43 API キー管理：**
[`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、API key の発行、一覧、失効、期限、scope を実装し、key 本体を保存しないことである。

**API 権限：**

| endpoint | 必要認証 |
|----------|----------|
| `GET /api/tokens` | 管理 session または `admin` scope API token |
| `POST /api/tokens` | 管理 session または `admin` scope API token |
| `DELETE /api/tokens/{id}` | 管理 session または `admin` scope API token |

`admin` scope API token は、`POST /api/tokens` の認証と scope 判定を通過した場合だけ新しい `admin` scope API token を発行する。発行者 token と発行対象 token は別 record とし、親子関係は保存しない。

**正常系：**

1. `POST /api/tokens` は `label`、`scopes`、`expires_at` を受け取る。
2. token 本体は `io.ReadFull(crypto/rand.Reader, buffer)` で取得した 32 bytes を `base64.RawURLEncoding.EncodeToString` で符号化し、`act_` を前置した文字列とする。
3. `.api_tokens` には `sha256(token)` の lowercase hex だけを保存する。
4. response の `token` は作成時 1 回だけ返す。`GET /api/tokens` では返さない。
5. 認証時は hash 一致、`revoked_at == null`、`expires_at == null または now < expires_at` を満たす token だけ有効とする。
6. 作成、認証成功、失効、期限切れ拒否を `.audit_log` へ記録する。ただし token 本体は記録しない。

**`.api_tokens` schema 参照：**

`.api_tokens` の top-level object、record のキー、型、必須性、許容値、正規化例外、採番固定値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.api_tokens` schema を唯一の正本とする。[`docs/details/security.md`](security.md) 詳細本文責務は token の生成、hash、認証、scope 判定、作成・失効・認証時の副作用と漏えい禁止だけを定義し、schema を再掲しない。

`.api_tokens` が同 schema に違反する場合は破損として扱い、token 認証、一覧、作成、失効をすべて `500 {"error":"Internal server error"}` で拒否する。破損内容、hash、token 本体は response、`.access_log`、`.audit_log`、journal に出力しない。

**API token 認証時の副作用境界：**

| ケース | HTTP status | `.api_tokens` | `.access_log` | `.audit_log` | 備考 |
|--------|-------------|---------------|---------------|--------------|------|
| hash 不一致 | `401` | 変更なし | 追記しない | 追記しない | token id が特定できないため監査対象外。 |
| 期限切れ | `401` | `last_used_at` 更新なし | `token_expired` を追記 | `token_expired` を追記 | token id が特定できた場合だけ記録する。 |
| 失効済み | `401` | `last_used_at` 更新なし | `token_revoked_reject` を追記 | `token_revoked_reject` を追記 | token 本体と hash は保存しない。 |
| scope 不足 | `403` | `last_used_at` 更新なし | `permission_denied` を追記 | `permission_denied` を追記 | 対象 endpoint は実行しない。 |
| 認証成功 | endpoint 固有 | `last_used_at` を現在時刻へ更新 | `token_auth` を追記 | `token_auth` を追記 | endpoint 実行前に更新する。 |

[`docs/details/security.md` 詳細本文責務 §27.43](security.md#sec-27-43) の固定表で `.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、対象 API は `500` とし、endpoint 固有処理へ進まない。`last_used_at` 更新済み record は巻き戻さない。

**処理順：**

`POST /api/tokens` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. `label`、`scopes`、`expires_at` を検証する。
3. `.api_tokens` のファイルロックを取得する。
4. 既存 record を読み込み、id を採番する。
5. token 本体を生成し、hash を算出する。
6. record を append して `.api_tokens` を原子的に保存する。
7. `.access_log` に `token_create` を追記する。
8. `.audit_log` に `token_create` を追記する。
9. response に token 本体を 1 回だけ含めて返す。

`.api_tokens` 保存後に `.access_log` または `.audit_log` 追記へ失敗した場合、API は `500` を返す。作成済み token record は削除せず、再実行時は新しい token を発行する。

`DELETE /api/tokens/{id}` は以下の順で処理する。

1. 認証と `admin` scope を確認する。
2. path `id` を検証する。
3. `.api_tokens` をロックして対象 record を検索する。
4. 対象が存在しない、または `revoked_at != null` の場合は `404` を返す。
5. `revoked_at` を現在時刻に設定して保存する。
6. `.access_log` と `.audit_log` に `token_revoke` を追記する。
7. `{ "message": "Token revoked" }` を返す。

認証に使用中の API token 自身を path `id` に指定した場合、その token を失効対象にする。当該リクエストは成功し、次リクエストから `401` になる。

**token ID 採番・返却固定契約：**

| 項目 | 仕様 |
|------|------|
| 採番元 | `.api_tokens.tokens[].id` の数値 suffix 最大値。存在しない場合は `tok000001`。 |
| 衝突時 | 最大 suffix + 1 を採用する。削除済みや失効済み id は再利用しない。 |
| token 本体 | `io.ReadFull(crypto/rand.Reader, buffer)` で 32 bytes を完全に取得し、`base64.RawURLEncoding.EncodeToString` で padding なしの 43 ASCII 文字へ変換し、`act_` を前置する。完成 token は合計 47 bytes、本体 alphabet は `[A-Za-z0-9_-]` に固定する。 |
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
| `crypto/rand.Reader` が 32 bytes を取得できない | `500 {"error":"Internal server error"}`。token、hash、record、access log、audit log を作成せず、lock を解放する。短い読み取りを成功扱いしない。 |
| 生成した `sha256(token)` が既存 record の `token_hash` と一致 | `500 {"error":"Internal server error"}`。自動再生成せず、既存 record と全 log を変更せず lock を解放する。token または hash の値を出力しない。 |

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

<a id="sec-27-44"></a>
**27.44 監査ログ：**
owner component は `security` とする。collaborator component は `api`、`runner`、`statefile` とする。


本機能の目的は、認証、権限拒否、token 作成・失効、設定変更、build trigger、session revoke、TOTP、approval の状態変更を追跡できる JSON Lines 監査ログとして保存することである。

**record 生成規則：**

| 項目 | 仕様 |
|------|------|
| `timestamp` | 操作結果が確定した UTC 時刻。 |
| `request_id` | API event は [`docs/details/api.md` 詳細本文責務 request ID 固定契約](api.md#api-request-id-contract) で生成した値を使用し、同一 API 処理中に複数 log を書く場合は同じ値を使う。`runner` 内部 event は `null`。 |
| `actor_type` | 未認証 login は `"anonymous"`、管理 session は `"admin"`、API token は `"api_token"`、署名検証済み Webhook は `"webhook"`、内部処理は `"system"`。 |
| `actor_id` | 管理 session は `"admin"`、API token は token id、署名検証済み Webhook は `"webhook"`、未認証は `null`、内部処理は `"system"`。 |
| `target_id` | 対象 id がある場合は id。endpoint 拒否は `"{METHOD} {path}"`。対象なしは `null`。 |
| `result` | 成功は `"success"`、認証失敗や検証失敗は `"failure"`、権限拒否は `"denied"`。 |
| `remote_addr` | [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6) の `remote_addr` 導出規則に従う。API request に紐づかない内部処理は `null`。 |
| `message` | 固定文言のみ。入力値を連結しない。最大 500 文字。 |

監査 record のキー、型、必須性、許容値は [`docs/details/statefile.md` 詳細本文責務 `.audit_log` schema](statefile.md#audit-log-schema) を唯一の正本とする。監査ログへ保存する object はその schema のキーだけとし、未知キー、request body、query 全体、header 全体、cookie、secret、token、password、hash、salt、TOTP secret を保存してはならない。

**action 生成固定値：**

| 操作 | action | actor_type | target_type | target_id | result |
|------|--------|------------|-------------|-----------|--------|
| password login 成功 | `login_success` | `anonymous` | `auth` | `login` | `success` |
| password login 失敗 | `login_failure` | `anonymous` | `auth` | `login` | `failure` |
| login ロック拒否 | `permission_denied` | `anonymous` | `auth` | `login` | `denied` |
| TOTP 必須 ticket 発行 | `totp_required` | `anonymous` | `auth` | `login` | `success` |
| TOTP login 失敗 | `totp_failure` | `anonymous` | `auth` | `login_totp` | `failure` |
| logout | `logout` | `admin` または `api_token` | `auth` | `logout` | `success` |
| 他 session 一括失効 | `session_revoke_all` | `admin` または `api_token` | `auth` | `sessions` | `success` |
| password 変更 | `password_change` | `admin` | `auth` | `password` | `success` |
| API token 作成 | `token_create` | `admin` または `api_token` | `api_token` | 作成 token id | `success` |
| API token 認証成功 | `token_auth` | `api_token` | `api_token` | token id | `success` |
| API token 期限切れ拒否 | `token_expired` | `api_token` | `api_token` | token id | `failure` |
| API token 失効済み拒否 | `token_revoked_reject` | `api_token` | `api_token` | token id | `failure` |
| API token 失効 | `token_revoke` | `admin` または `api_token` | `api_token` | 失効 token id | `success` |
| scope 不足 | `permission_denied` | `api_token` | `endpoint` | `{METHOD} {path}` | `denied` |
| rate limit 超過 | `permission_denied` | `anonymous`、`admin`、`api_token`、`webhook` | `endpoint` | `{METHOD} {path}` | `denied` |
| TOTP setup 開始 | `totp_setup` | `admin` または `api_token` | `auth` | `totp_setup` | `success` |
| TOTP confirm / disable code 失敗 | `totp_failure` | `admin` または `api_token` | `auth` | `totp_confirm` または `totp_disable` | `failure` |
| TOTP 有効化 | `totp_enabled` | `admin` または `api_token` | `auth` | `totp` | `success` |
| TOTP 無効化 | `totp_disabled` | `admin` または `api_token` | `auth` | `totp` | `success` |
| 通常 build 開始または queue 追加 | `build_trigger` | `admin` または `api_token` | `build` | build id または queue id | `success` |
| force build 開始または queue 追加 | `build_force_trigger` | `admin` または `api_token` | `build` | build id または queue id | `success` |
| `.config_log` 対象設定変更成功 | `config_update` | `admin` または `api_token` | `config` | `.config_log.type` と同じ固定名 | `success` |
| `.config_log` 対象設定変更部分失敗 | `config_update` | `admin` または `api_token` | `config` | `.config_log.type` と同じ固定名 | `failure` |
| rate limit 設定変更成功 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `success` |
| rate limit 設定変更部分失敗 | `rate_limit_update` | `admin` または `api_token` | `config` | `api_rate_limit` | `failure` |
| approval pending 作成 | `approval_pending` | `system` | `approval` | approval id | `success` |
| approval 承認 | `approval_approved` | `admin` または `api_token` | `approval` | approval id | `success` |
| approval 却下 | `approval_rejected` | `admin` または `api_token` | `approval` | approval id | `success` |
| approval 期限切れ | `approval_expired` | `system` | `approval` | approval id | `success` |
| share link 作成 | `share_link_create` | `admin` または `api_token` | `share_link` | 作成 share link id | `success` |
| share link 失効 | `share_link_revoke` | `admin` または `api_token` | `share_link` | 失効 share link id | `success` |

<a id="sec-27-44-config-audit-order"></a>
**設定変更の保存順・失敗固定契約：**

同一操作で `.config_log` と `.audit_log` の両方を追記する場合、対象状態保存 → `.config_log` → `.audit_log` → response の順に固定する。`.config_log` 追記失敗時は server log に `CONFIG_LOG_WRITE_FAILED` だけを記録し、`.audit_log` を試行せず `500` を返す。`.audit_log` 追記失敗時は server log に `AUDIT_LOG_WRITE_FAILED` だけを記録し、`500` を返す。いずれの失敗でも保存済み対象状態と追記済み先行 log を巻き戻さず、入力値、diff、Go error、path、secret を server log に含めない。`.audit_log` 追記失敗そのものを `.audit_log` に再帰記録してはならない。

**監査 record 保存順・失敗契約：**

| 操作種別 | 保存順 | `.audit_log` 失敗時 |
|----------|--------|----------------------|
| 認証成功 | session または token 状態更新 → `.access_log` → `.audit_log` → response | token / session 状態は巻き戻さず `500`。response に token 本体を含めない。 |
| 認証失敗 | `.access_log` → `.audit_log` → response | `500`。失敗理由詳細は返さない。 |
| 権限拒否 | `.access_log` → `.audit_log` → `403` | `500`。対象 endpoint は実行しない。 |
| 設定変更 | [設定変更の保存順・失敗固定契約](#sec-27-44-config-audit-order) を適用する。 | 同契約を適用する。 |
| token 作成 | `.api_tokens` 保存 → `.access_log` → `.audit_log` → response | 作成済み record は残し、token 本体は返さず `500`。 |
| build trigger | build / queue 状態保存 → `.audit_log` → response | 保存済み build / queue 状態は巻き戻さず `500`。 |
| session revoke / TOTP | 対象状態更新 → `.access_log` が必要な場合は追記 → `.audit_log` → response | 保存済み状態は巻き戻さず `500`。one-time secret または token は返さない。 |
| approval API | queue / approval / history を [`docs/details/api.md` 詳細本文責務 §27.30](api.md#sec-27-30) の順で更新 → `.audit_log` → response | 保存済み状態は巻き戻さず `500`。 |
| approval runner event | approval / history を [`docs/details/runner.md` 詳細本文責務 §27.30](runner.md#sec-27-30) の順で更新 → `.audit_log` | 保存済み状態は巻き戻さず runner failure とし、後続通知は実行しない。 |
| share link 作成 | `.share_links` 保存 → `.audit_log` → `.admin_events` → token response | 保存済み share link record は巻き戻さず、token 本体は返さず `500`。 |
| share link 失効 | `.share_links` 保存 → `.audit_log` → `.admin_events` → response | 保存済み失効状態は巻き戻さず `500`。 |

**取得仕様：**

| 項目 | 仕様 |
|------|------|
| 読み込み順 | ファイル先頭から全有効行を読み、フィルタ後に `timestamp` 降順、同時刻はファイル出現順の逆順で返す。 |
| `total` | フィルタ後、ページング前の有効 record 件数。 |
| `actor` filter | `actor_id` と完全一致。`anonymous` を指定した場合は `actor_type:"anonymous"` かつ `actor_id:null` に一致させる。 |
| `action` filter | `action` 完全一致。 |
| `result` filter | `success`、`failure`、`denied` の完全一致。 |
| 壊れた行 | 無視する。API response に壊れた行の内容を含めない。 |
| 返却上限 | `limit` は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の 1〜200。未指定時は 100。 |

`GET /api/audit-log` は監査ログ取得操作自体を `.audit_log` へ記録しない。通常の API request として `.api_access_log` には記録する。フィルタ値が未知 action、未知 result、200 Unicode scalar values 超過の actor の場合は `422` とし、壊れた行の有無とは独立して判定する。

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

<a id="sec-27-45"></a>
**27.45 セッションタイムアウト変更設定：**
[`docs/details/security.md` 詳細本文責務 §27.45](security.md#sec-27-45) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、新規 session の有効期限を管理 API から更新し、既存 session への影響を明確にすることである。

**設定参照：**

| 項目 | 値 |
|------|----|
| 設定 schema | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.session_timeout_seconds`。型、既定値、許容範囲は同 schema を正とする。 |
| 更新 API | [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の `POST /api/config`。SDK 呼出は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の `setConfig(config)` を正とする。 |

設定変更は新規 session にだけ適用する。既存 session の `expires_at` は延長も短縮もしない。

**処理契約：**

| 項目 | 仕様 |
|------|------|
| 保存先 | `.server_config.session_timeout_seconds`。他 key と同時更新された場合も [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の単一ファイル更新手順で保存する。 |
| 既定値 merge | `.server_config` に key がない場合、`GET /api/config` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 既定値を返す。ファイルへ暗黙保存しない。 |
| login 時適用 | session token 発行直前に `.server_config` を読み、当該時点の値で `expires_at` を計算する。 |
| TOTP login | TOTP 有効時は `POST /api/login/totp` の成功時点で値を読む。`POST /api/login` の password 成功時点では session を発行しない。 |
| 変更監査 | `POST /api/config` で値が変わった場合は `.config_log` に差分を記録する。`.audit_log` は `config_update`、`target_type:"config"`、`target_id:"session_timeout_seconds"` を記録する。 |
| 同値更新 | 同じ値の更新は `200` とし、`.server_config`、`.config_log`、`.audit_log`、既存 session を変更しない。 |

session timeout の値は session 発行時に秒単位で加算する。`expires_at = issued_at + session_timeout_seconds` とし、計算後の時刻は UTC ISO 8601 秒精度で保存する。ミリ秒、ナノ秒、local timezone は保存しない。

**検証条件：**

| ケース | 期待結果 |
|--------|----------|
| schema 最小値の設定 | 新規 session が設定値の秒数後に期限切れ。 |
| 既存 session | 設定変更後も元の `expires_at`。 |
| 範囲外 | `422`。 |
| key 不在 | `GET /api/config` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema 既定値を返す。 |
| TOTP login | `POST /api/login/totp` 成功時点の値で session 期限を決める。 |
| 秒精度 | `expires_at` は UTC ISO 8601 秒精度。 |

<a id="sec-27-46"></a>
**27.46 TOTP 二要素認証：**
[`docs/details/security.md` 詳細本文責務 §27.46](security.md#sec-27-46) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


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

QR code 生成はこの詳細本文責務では仕様化しない。UI は secret と otpauth URI を一回表示し、ユーザーの認証アプリ登録手段は手入力または URI 貼り付けに限定する。

**メモリ上状態：**

| 状態 | 保存場所 | 期限 | 内容 |
|------|----------|------|------|
| setup 仮 secret | `api` のメモリ | 10 分 | `secret_base32`, `created_at`。サーバー再起動で破棄する。 |
| login ticket | `api` のメモリ | 5 分 | `ticket_hash`, `created_at`, `password_verified_at`, `credentials_fingerprint`。ticket 本体は hash 化して保持する。 |

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
| `POST /api/login` | password 成功かつ TOTP 有効なら `ticket` を `crypto/rand` 32 bytes の lowercase hex で生成する。credentials を更新せず、ticket に `must_change` を保持しない。HTTP response は [`docs/details/api.md` 詳細本文責務 `POST /api/login`](api.md#25-認証-実装仕様) の `LoginResult` を参照する。 |
| `POST /api/login/totp` | ticket hash と code を検証し、成功時に ticket を削除する。最新 credentials を再読込して session 発行時の `must_change` を算出する。失敗時も ticket は削除する。HTTP response は [`docs/details/api.md` 詳細本文責務 `POST /api/login/totp`](api.md#25-認証-実装仕様) を参照する。 |

TOTP code は 6 桁の ASCII 数字のみ受け付ける。空文字、全角数字、空白付き文字列、6 桁以外は `422` とする。検証は `window` 内の step を古い順に試し、最初に一致した step を採用する。採用 step が `.totp_secret.last_accepted_step` 以下の場合は `401` とする。

TOTP setup、ticket 発行、login TOTP 成功 / 失敗、confirm 成功 / 失敗、disable 成功 / 失敗は `.audit_log` へ記録する。code 不一致、replay、ticket 不正は `result:"failure"` とし、code、secret、ticket 本体は保存しない。

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
| setup 開始 | 仮 secret 生成 → メモリ保存 → `.audit_log` 追記 → response | メモリ保存または audit 追記失敗時は仮 secret を削除して `500`、secret を返さない。 |
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

<a id="sec-27-47"></a>
**27.47 API レート制限：**
[`docs/details/security.md` 詳細本文責務 §27.47](security.md#sec-27-47) の境界は owner component `security`、collaborator component `api`、`sdk`、`ui`、`statefile` とする。


本機能の目的は、login 総当たり、API token 濫用、外部連携の暴走を Go 標準ライブラリだけで抑止することである。

**固定窓 policy：**

| group | 既定 window | 既定 max | 対象 |
|-------|-------------|----------|------|
| `login` | 60 秒 | 10 | `POST /api/login`、`POST /api/login/totp`。 |
| `read` | 60 秒 | 600 | 読み取り endpoint。 |
| `trigger` | 60 秒 | 60 | build trigger endpoint と署名検証済み `POST /api/webhook`。 |
| `operate` | 60 秒 | 120 | 運用操作 endpoint。 |
| `config` | 60 秒 | 60 | 設定変更 endpoint。 |
| `admin` | 60 秒 | 60 | token、audit、rate limit endpoint。 |

**判定キー：**

| 認証状態 | key |
|----------|-----|
| 認証前 | `ip:{remote_addr}:{group}` |
| session | `session:admin:{group}` と `ip:{remote_addr}:{group}` の両方 |
| API token | `token:{token_id}:{group}` と `ip:{remote_addr}:{group}` の両方 |
| 署名検証済み Webhook | `ip:{remote_addr}:trigger` だけ |

どちらか一方でも上限を超えた場合は `429 Too Many Requests` と `{"error":"Too many requests"}` を返す。

`remote_addr` は `net/http.Request.RemoteAddr` の host 部分を使用する。`X-Forwarded-For`、`X-Real-IP`、`Forwarded` header は標準では信用せず、key 生成に使用しない。IPv6 は `net.SplitHostPort` で host を抽出し、正規化済み文字列をそのまま key に入れる。host 抽出に失敗した場合は `ip:unknown:{group}` を使用する。

**policy / response schema 参照：**

永続化する `.server_config.api_rate_limit` のキー、型、必須 group、範囲、既定値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ApiRateLimitPolicy object を唯一の正本とする。`GET /api/api-rate-limit` と `POST /api/api-rate-limit` の request / response、response 専用 `state_summary` のキー、並び順、件数は [`docs/details/api.md` 詳細本文責務 security owner 機能 API 側境界](api.md#sec-security-owner-api-boundary) を正とする。`state_summary` は永続化しない。

**正常系：**

1. path 解決後、認証前に IP key の login 制限を確認する。
2. 認証後、endpoint group を判定し、actor key と IP key の count を更新する。署名検証済み Webhook は `trigger` group の IP key だけを更新する。
3. `GET /api/api-rate-limit` は policy と window summary を返す。
4. `POST /api/api-rate-limit` は policy を検証して保存し、`.api_rate_state.windows` を空にする。

<a id="security-rate-limit-group-map"></a>
**endpoint group 判定：**

| 条件 | group |
|------|-------|
| `POST /api/login`、`POST /api/login/totp` | `login` |
| 署名検証済み `POST /api/webhook` | `trigger` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `trigger` 許可 endpoint | `trigger` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `operate` 許可 endpoint | `operate` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `config` 許可 endpoint | `config` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `admin` 許可 endpoint | `admin` |
| [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の `read` 許可 endpoint | `read` |
| `GET /api/health` | rate limit 対象外 |
| 未知 path / method 不一致 | rate limit 判定前に `404` / `405` |

rate limit group は [endpoint group 判定表](#security-rate-limit-group-map) と [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の一意な endpoint 割当てから決定する。同じ method と正規化 path pattern が複数 group に存在する場合は仕様不整合であり、実装者判断の優先順で解決してはならない。

**判定・更新手順：**

1. `.server_config.api_rate_limit.enabled == false` の場合、`.api_rate_state` を読まずに対象 API 処理へ進む。
2. endpoint group を決める。
3. `.api_rate_state` のファイルロックを取得する。
4. 対象 key ごとに `.api_rate_state.windows[key]` を確認する。
5. window が存在しない、または `now >= window_start + window_seconds` の場合、`window_start=now`, `count=0` で初期化する。
6. `count >= max_requests` の key が 1 つでもあれば、count を増やさず `.audit_log` に `permission_denied` を追記し、`429` を返す。
7. 上限未満の場合、対象 key すべての `count` を 1 増やして `.api_rate_state` を保存し、対象 API 処理へ進む。

rate limit の `429` は `.audit_log` に `permission_denied` として記録する。監査ログ追記に失敗した場合は `500` を返す。login group の認証前 `429` は `actor_type:"anonymous"`、`actor_id:null`、署名検証済み Webhook の `429` は `actor_type:"webhook"`、`actor_id:"webhook"` とする。

認証後 endpoint の rate limit では、actor key と IP key の両方を同じ lock 内で判定・更新する。片方だけの count 更新に成功した状態を残してはならない。`.api_rate_state` 保存失敗時は対象 API を実行せず `500` を返す。rate limit 判定で `429` になる request は count を増やさない。

**rate limit 副作用固定契約：**

| ケース | `.api_rate_state` | `.api_access_log` | `.audit_log` | endpoint 固有処理 |
|--------|-------------------|---------------|--------------|-------------------|
| 上限未満 | count を増やす | response 確定後に通常追記 | endpoint が監査対象の場合だけ追記 | 実行する |
| 上限超過 | count を増やさない | `429` として追記 | `permission_denied` を追記 | 実行しない |
| `.api_rate_state` 保存失敗 | 部分更新を残さない | `500` として追記を試行 | 追記しない | 実行しない |
| `.audit_log` 失敗 | count を増やさない | `500` として追記を試行 | 失敗 | 実行しない |

`429` 判定時は `.audit_log` 追記を `.api_rate_state` 保存前に行う。`.audit_log` 追記に成功した場合だけ `429` を返す。`.audit_log` 追記に失敗した場合は `.api_rate_state` を変更せず `500` を返す。

**設定更新契約：**

| 項目 | 仕様 |
|------|------|
| policy validation | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の ApiRateLimitPolicy object schema に完全一致する場合だけ保存する。 |
| 保存順 | `.server_config` 保存 → `.api_rate_state.windows` 空保存 → `.config_log` 追記 → `.audit_log` に `rate_limit_update` 追記 → response。 |
| 同値更新 | `200` とし、`.server_config` と `.api_rate_state` は変更しない。`.config_log` と `.audit_log` に追記しない。 |

**異常系：**

| 条件 | 処理 |
|------|------|
| policy schema 不一致 | `422`。`.server_config`、`.api_rate_state`、`.config_log`、`.audit_log` を変更しない。 |
| `.api_rate_state` 書き込み失敗 | `500`。対象 API は実行しない。 |
| `.api_rate_state` 破損 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) に従って退避し、空 window で再生成する。 |
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

<a id="sec-27-48"></a>
**27.48 マルチユーザー対応：**

マルチユーザー対応の owner は `security` とする。

`.users` と `.roles` の schema は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) を正本とする。

既定 user は初期セットアップ時に 1 件だけ作成する。既定 user は `admin` role を持つ。

session には `user_id`、`role_ids`、`permissions_hash`、`issued_at`、`expires_at` を保存する。

role または permission が変更された場合、既存 session の `permissions_hash` が不一致になった request は `401` とし、再 login を要求する。

最後の active admin user を disabled、locked、role 剥奪してはならない。

**監査対象：**

| 操作 | audit action |
|------|--------------|
| user create | `user_create` |
| user update | `user_update` |
| user disable | `user_disable` |
| role assignment change | `user_role_update` |

<a id="sec-27-50"></a>
**27.50 外部認証連携：**

外部認証連携の owner は `security` とする。

対応 provider type は `oidc` だけとする。

OIDC discovery URL は `issuer + "/.well-known/openid-configuration"` とする。

外部認証連携は local password 認証を暗黙に無効化しない。

provider secret は `.external_auth_config` に保存せず、`client_secret_ref` だけを保存する。

callback 検証では `state`、`nonce`、`issuer`、`audience`、`exp`、`iat`、`sub` を検証する。

検証失敗は `401` とし、失敗理由の詳細、token、claim 全文を response、server log、audit log に出してはならない。

<a id="sec-27-58"></a>
**27.58 ロールベースアクセス制御：**

ロールベースアクセス制御の owner は `security` とする。

permission 名は以下に固定する。

| permission | 許可範囲 |
|------------|----------|
| `*` | 全操作。system admin role だけが保持できる。 |
| `status:read` | status、history、metrics、badge の read。 |
| `build:write` | build trigger、queue reorder、queue cancel。 |
| `config:read` | config、snapshot、template の read。 |
| `config:write` | config、snapshot create / restore、template apply。 |
| `user:admin` | user、role、external auth、share link 管理。 |
| `system:admin` | datastore、cache、retention、event stream、webhook resend。 |

API endpoint group と必要 permission は以下に固定する。ここにない追加管理 API endpoint は `403` とし、実装判断で既存 permission に割り当ててはならない。route、request、response は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70) を参照する。

| endpoint 群 | 必要 permission |
|--------------|----------------|
| `/api/users`, `/api/roles`, `/api/external-auth-*`, `/api/share-links`, `/api/share/{token}/status` | `user:admin`。ただし `/api/share/{token}/status` は share token 検証成功時だけ session permission 不要。 |
| `/api/datastore`, `/api/datastore/switch`, `/api/cache-policy`, `/api/response-cache`, `/api/events`, `/api/events/stream`, `/api/notify-log/{delivery_id}/resend` | `system:admin` |
| `/api/config-snapshots`, `/api/config-snapshots/{id}`, `/api/config-snapshots/{id}/restore`, `/api/config-snapshots/{left_id}/diff/{right_id}`, `/api/config-templates`, `/api/config-templates/{id}/apply` | read は `config:read`、write / restore / apply / delete は `config:write` |
| `/api/projects` | read は `config:read`、write / archive は `config:write` |
| `/api/queue/{queue_id}`, `/api/queue/reorder` | `build:write` |
| `/api/history/retention`, `/api/history/retention/run` | read は `status:read`、write / run は `system:admin` |
| `/api/stats/export`, `/api/metrics`, `/api/badge/status.svg`, `/api/snapshots/{left_id}/diff/{right_id}` | `status:read` |
| `/api/version`, `/api/openapi.json` | 認証不要。個別 permission 不要。 |

`*` 以外の permission は暗黙に他 permission を含めない。

追加管理 API の認証・認可判定は以下の順に固定する。route、request、response、HTTP status の一覧は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70) を参照する。

| 順序 | 判定 | 失敗時 |
|------|------|--------|
| 1 | `/api/version` と `/api/openapi.json` は認証判定を行わない。 | security 副作用なし。 |
| 2 | `/api/share/{token}/status` は share token を検証する。 | token 不正は `401`、不在または revoke 済みは `404`、期限切れは `410`。session を作成しない。 |
| 3 | その他の追加管理 API は session または API token を検証する。 | 不在、不正、期限切れは `401`。対象状態は読まない。 |
| 4 | user status を確認する。 | `disabled` または `locked` は `401`。対象状態は読まない。 |
| 5 | permission を判定する。 | 不足は `403`。対象状態は読まず、permission denied audit だけを許可する。 |
| 6 | endpoint 固有の user / role / share link 制約を確認する。 | 最後の admin 保護、system role 変更、role 使用中、self lock は `409`。 |

初期 system role は以下に固定する。

| role id | name | permissions | system |
|---------|------|-------------|--------|
| `admin` | `Administrator` | `*` | `true` |
| `operator` | `Operator` | `status:read`, `build:write`, `config:read` | `true` |
| `viewer` | `Viewer` | `status:read` | `true` |

system role は削除禁止、`id` 変更禁止、`system:false` への変更禁止とする。system role の `name` と `permissions` を変更する request は `409` とし、状態を変更しない。

role 削除時、対象 role を持つ active user が 1 件でも存在する場合は `409` とする。

<a id="sec-27-61"></a>
**27.61 ユーザー管理 API：**

ユーザー管理 API の owner は `security` とする。

username は case-sensitive ではなく、保存前に lowercase に正規化する。

同一 username の active または disabled user が存在する場合、作成は `409` とする。

password 更新時は平文を保存せず、既存 `.admin_credentials` と同じ hash 方式を使用する。

password 未設定 user は local password login を禁止し、外部認証または API token だけを許可する。

disabled user の session は次 request で `401` とする。

locked user の password login は `401` とし、API token と外部認証も拒否する。

user 更新時の最後の admin 保護は、変更適用後に `*` permission を持つ active user が 1 件以上残ることを条件とする。自分自身の `role_ids` 変更、`status` 変更、password 削除、external subject 削除がこの条件を満たさない場合は `409` とし、session、user、role、audit、admin event を変更しない。

<a id="sec-27-67"></a>
**27.67 読み取り専用共有リンク：**

読み取り専用共有リンクの owner は `security` とする。

share token は 32 byte 以上の乱数を URL-safe Base64 で表現する。

token 本体は作成 response で 1 回だけ返し、`.share_links` には SHA-256 hash だけを保存する。

share link scope は `status`、`history`、`snapshot_diff` のいずれかとする。

share link は read-only であり、build trigger、config update、queue 操作、user 操作、webhook resend を実行してはならない。

期限切れは `410`、revoke 済みは `404` とする。

share link request は通常 session を作成しない。

share link 作成時の `expires_at` は `null` または現在時刻より後の UTC ISO 8601 秒精度とする。過去時刻、現在時刻と同一秒、local timezone、offset 付き時刻は `422` とする。revoke は `revoked_at` が `null` の record だけを対象とし、revoke 済み record への再 revoke は `404` とする。

share token 生成は CSPRNG から 32 byte を取得し、padding なし URL-safe Base64 で表現する。生成 token の文字集合は `A-Z`、`a-z`、`0-9`、`-`、`_` だけとし、`=`、`+`、`/`、空白、改行を含めてはならない。CSPRNG 失敗時は `.share_links`、audit、admin event を変更せず `500` とする。token hash は token 文字列の UTF-8 byte 列ではなく、Base64 decode 後の 32 byte 以上の raw token byte に対する SHA-256 lowercase hex とする。

share token 検証は、(1) path parameter を percent decode する、(2) URL-safe Base64 padding なしとして decode する、(3) decode 後 byte length が 32 以上であることを確認する、(4) SHA-256 lowercase hex を算出する、(5) `.share_links.links[].token_hash` と constant-time 比較する、(6) `revoked_at`、(7) `expires_at` の順に行う。decode 不能または 32 byte 未満では `.share_links` を読まない。hash 不一致、不在、revoke 済みでは同じ `404` 結果とし、どの条件だったかを response、audit、admin event、access log で区別できる値として出さない。

share link 管理 API の `.audit_log` は [`docs/details/statefile.md` 詳細本文責務 `.audit_log` schema](statefile.md#audit-log-schema) の key だけを使用し、`action` は `share_link_create` または `share_link_revoke`、`target_type` は `share_link`、`target_id` は share link id、`actor_type` と `actor_id` は管理 session または API token とする。share link 管理 API の `.admin_events` は [`docs/details/statefile.md` 詳細本文責務 AdminEventRecord](statefile.md#sec-22-0d) の key だけを使用し、`type` は `security`、`target_type` は `share_link`、`target_id` は share link id、`actor` は管理 session user id または API token id、`message` は固定文言だけとする。`GET /api/share/{token}/status` は `.audit_log` と `.admin_events` を書き込まない。`.audit_log`、`.admin_events`、`.access_log`、`.api_access_log`、server log には token 本体、token hash、token prefix、token length、Authorization header、Cookie、raw path を保存してはならない。作成 response 返却後の token 再表示、list response への token 追加、SDK / UI による token 復元、log からの token 復元を禁止する。

share link の `scope` は `status`、`history`、`snapshot_diff` の単一値だけを許可する。複数 scope、wildcard、空配列、将来 scope 名、permission 名、API path を受け付けてはならない。`scope` ごとの response composition は [`docs/details/api.md` 詳細本文責務 §27.67](api.md#sec-27-67) を参照し、security owner は token、hash、期限、revoke、漏えい禁止だけを正本として持つ。
