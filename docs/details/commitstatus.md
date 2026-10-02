# Adlaire CI — Commit Status 詳細仕様

---

<a id="0-責務境界"></a>

**0. 責務境界：**

| 項目 | 内容 |
|------|------|
| owner component | `commitstatus` |
| 実装主体 | [`components/commitstatus/`](../../components/commitstatus/)。GitHub Commit Status の payload 生成、送信、結果保存境界は [`components/runner/`](../../components/runner/) と接続する。 |
| 持つ内容 | `commitstatus` owner が主本文として定義する GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| 検証接続 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 owner 詳細本文 検証接続共通入口](../DETAIL_INDEX.md#owner-detail-verification-route) を参照する。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | GitHub Commit Status API。 |

---

<a id="sec-27-1"></a>
**27.1 GitHub Commit Status API：**

owner component は `commitstatus` とする。collaborator component は `runner`、`statefile` とする。

runner は `.server_config.commit_status_enabled == true` の場合、commitstatus component を呼び出し、対象 commit に GitHub Commit Status API を送信する。送信先は `POST /repos/{owner}/{repo}/statuses/{sha}` とし、`sha` は対象 target の commit SHA を使用する。blob SHA しか取得できない場合は status を送信せず、`.build_logs/{id}.json.commit_status.error="commit sha unavailable"` を保存する。

送信 payload は次に固定する。

| key | 値 |
|-----|----|
| `state` | 開始時 `"pending"`、成功時 `"success"`、失敗時 `"failure"`、GitHub API 認証失敗、commit status 設定不正、状態ファイル読込失敗、状態ファイル書込失敗、payload 生成不能のいずれかの CI 自体の異常時 `"error"`。 |
| `context` | `.server_config.commit_status_context`。型、既定値、許容値は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config` schema を正とする。 |
| `description` | 140 文字以内。開始時 `Build started`、成功時 `Build succeeded`、失敗時 `Build failed: <target_status>`、pending deploy 時 `Build succeeded with deploy pending`。 |
| `target_url` | `.server_config.commit_status_target_url` が `null` でなければ送信する。 |

GitHub request は以下に固定する。retry、GraphQL API、Check Runs API、任意 header、任意 payload key は [`docs/details/commitstatus.md`](commitstatus.md) 詳細本文責務の対象外とする。

| 項目 | 固定値 |
|------|--------|
| method | `POST` |
| path | `/repos/{owner}/{repo}/statuses/{sha}` |
| header | `Accept: application/vnd.github+json`、`Authorization: Bearer {GITHUB_TOKEN}`、`X-GitHub-Api-Version: 2022-11-28` |
| 呼出条件 | `.server_config.commit_status_enabled=true` の場合だけ呼び出す。GitHub permission は [`docs/SPEC.md` ポリシー責務 §5](../SPEC.md#policy-runner-secrets) を正本とする。 |
| body | `state`、`context`、`description`、任意の `target_url` だけを持つ JSON object。 |
| timeout | 10 秒。timeout は network error と同じ送信失敗扱い。 |
| success status | HTTP `201`。`200`、`202`、`204` は失敗扱い。 |
| retry | 行わない。pending / final の各送信は 1 回だけ。 |

処理順序は以下とする。

1. runner が commit SHA を確定する。
2. runner が build id を採番する。
3. `.build_status.json status="running"` を書く前に、runner が commitstatus component を呼び出し、GitHub status `"pending"` を送信する。
4. pipeline / deploy / snapshot / history の最終結果確定後、runner が commitstatus component を呼び出し、GitHub status の最終 state を送信する。
5. `.build_logs/{id}.json.commit_status` に最終送信結果を保存する。
6. `.build_history.commit_status_state` に最終 state を保存する。

Commit Status 送信失敗は build 成否を反転させない。送信失敗時は WARN ログ `COMMIT_STATUS_FAILED: status=<http_status> error=<reason>` を出し、`.build_logs/{id}.json.commit_status.state` を最後に送信しようとした state、`error` を固定文言で保存する。GitHub API 認証失敗、403、404、5xx、network error はすべて送信失敗として扱う。

**Commit Status 呼び出し入力固定契約：**

| 入力 | 取得元 | validation | 不正時 |
|------|--------|------------|--------|
| `owner` | runner repository 設定 | 1〜100 文字、`/`、空白、制御文字禁止。 | 送信せず、`state="error"`、`error="invalid repository owner"`。 |
| `repo` | runner repository 設定 | 1〜100 文字、`/`、空白、制御文字禁止。 | 送信せず、`state="error"`、`error="invalid repository name"`。 |
| `sha` | runner が取得した commit SHA | 40 文字 hex。blob SHA は不可。 | 送信せず、`state=null`、`error="commit sha unavailable"`。 |
| `context` | `.server_config.commit_status_context` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config` schema に一致すること。 | config validation で `422`。runner 実行時に schema 不一致を検出した場合は送信せず `state="error"`。 |
| `target_url` | `.server_config.commit_status_target_url` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config` schema に一致すること。 | config validation で `422`。runner 実行時に schema 不一致を検出した場合は送信せず `state="error"`。 |
| GitHub token | runner secret | 空でない文字列。 | 送信せず、`state="error"`、`error="github token unavailable"`。secret 値は出力しない。 |
| build result | runner final result | `"success"`、`"failure"`、`"success_deploy_pending"`、CI internal error のいずれかへ正規化する。 | 正規化不能なら `state="error"`。 |

**state 算出固定契約：**

| runner 状態 | Commit Status `state` | description |
|-------------|-----------------------|-------------|
| build 開始前 | `"pending"` | `Build started` |
| pipeline / build / deploy 成功 | `"success"` | `Build succeeded` |
| build 成功、deploy pending | `"success"` | `Build succeeded with deploy pending` |
| pipeline / build / deploy 失敗 | `"failure"` | `Build failed: <target_status>` |
| GitHub token 不在、payload 生成不能、設定不正、状態 read/write 失敗 | `"error"` | `Build status error: <reason>` |

`description` の `<target_status>` と `<reason>` は英数字、空白、`_`、`-`、`:`、`.` だけに正規化する。その他の文字、改行、tab は空白へ置換し、連続空白は 1 個へ畳み込む。140 文字を超える場合は 137 文字 + `...` に切り詰める。

**Commit Status 保存固定契約：**

`.build_logs/{id}.json.commit_status` は [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の CommitStatus object に定義された key だけを保存する。`pending_sent`、`pending_error`、`final_sent` など未定義 key を保存してはならない。pending / final の送信順、HTTP request、response、失敗有無は fixture の `expected/effects.json` と server WARN log で検証し、build log schema へ未定義 key を追加しない。

| ケース | `enabled` | `state` | `sent_at` | `http_status` | `error` |
|--------|-----------|---------|-----------|---------------|---------|
| disabled | `false` | `null` | `null` | `null` | `null` |
| commit SHA なし | `true` | `null` | `null` | `null` | `"commit sha unavailable"` |
| pending 成功、final 成功 | `true` | final state | final 送信時刻 | final HTTP status | `null` |
| pending 失敗、final 成功 | `true` | final state | final 送信時刻 | final HTTP status | `null`。pending failure は WARN log と `expected/effects.json` で検証する。 |
| pending 成功、final 失敗 | `true` | final に送信しようとした state | final 試行時刻 | final HTTP status または `null` | 固定 error reason。 |
| pending 失敗、final 失敗 | `true` | final に送信しようとした state | final 試行時刻 | final HTTP status または `null` | final failure の固定 error reason。pending failure は WARN log と `expected/effects.json` で検証する。 |
| payload 生成不能 | `true` | `"error"` | `null` | `null` | 固定 error reason。 |
| GitHub token 不在 | `true` | `"error"` | `null` | `null` | `"github token unavailable"` |

`.build_history.commit_status_state` は `.build_logs/{id}.json.commit_status.state` と同じ値を保存する。`commit_status` が `null` の場合、または disabled の場合は `null` を保存する。commit SHA なしの場合も `null` とし、build result 自体を failure へ反転してはならない。

**Commit Status error reason 固定契約：**

| 失敗 | 保存する `error` | WARN log |
|------|------------------|----------|
| HTTP `401` / `403` | `"github status unauthorized"` | `COMMIT_STATUS_FAILED: status={status} error=github status unauthorized` |
| HTTP `404` | `"github status target not found"` | `COMMIT_STATUS_FAILED: status=404 error=github status target not found` |
| HTTP `429` | `"github status rate limited"` | `COMMIT_STATUS_FAILED: status=429 error=github status rate limited` |
| HTTP `5xx` | `"github status server error"` | `COMMIT_STATUS_FAILED: status={status} error=github status server error` |
| HTTP `201` 以外 | `"github status unexpected response"` | `COMMIT_STATUS_FAILED: status={status} error=github status unexpected response` |
| network / DNS / timeout | `"github status network error"` | `COMMIT_STATUS_FAILED: status=0 error=github status network error` |
| invalid owner / repo / context / target_url | [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) の固定表の固定 validation error | `COMMIT_STATUS_FAILED: status=0 error={reason}` |
| state write failure | `"commit status state write failed"` | `COMMIT_STATUS_FAILED: status=0 error=commit status state write failed` |

GitHub response body 全体、Authorization header、GitHub token、credential 付き URL は build log、history、server log、fixture expected に保存しない。HTTP status、固定 error reason、request path、payload state/context/description/target_url だけを保存・検証対象にする。

**副作用境界固定契約：**

| ケース | 許可する副作用 | 禁止する副作用 |
|--------|----------------|----------------|
| disabled | `.build_logs/{id}.json.commit_status.enabled=false`、history `commit_status_state=null`。 | GitHub Status API 呼び出し、WARN log、payload 生成。 |
| dry-run | stdout の `would_write` に `"commit_status"` を含めることだけ。 | GitHub write API 呼び出し、build log / history / status / lock / SHA cache 更新。 |
| commit SHA なし | build log / history への summary 保存。 | GitHub Status API 呼び出し、build 成否反転、SHA cache 更新。 |
| pending 送信失敗 | WARN log、fixture effects 上の failed call 記録。 | build 停止、pipeline / deploy / snapshot / history の中断。 |
| final 送信失敗 | WARN log、commit_status summary の error 保存。 | build 成否反転、deploy result 反転、snapshot / history の削除。 |
| payload validation failure | commit_status summary の error 保存。 | GitHub Status API 呼び出し、未定義 key 保存、secret 出力。 |
| state write failure | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の state write failure 契約に従う。 | Commit Status 失敗を理由に未定義 rollback を実行すること。 |

**Commit Status fixture 参照：**

Commit Status fixture の fake GitHub Status API、expected/effects、expected logs、expected security、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) を正本とする。[`docs/details/commitstatus.md`](commitstatus.md) 詳細本文責務では、GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証観点だけを扱う。

検証観点:

| ケース | 期待結果 |
|--------|----------|
| 有効・成功 build | fake GitHub server に `pending` → `success` の順で 2 回送信される。 |
| 有効・pipeline 失敗 | `pending` → `failure` の順で送信される。 |
| deploy pending | 最終 state は `success`、description は deploy pending を含む。 |
| commit SHA なし | status 送信なし、build は継続、log に `commit sha unavailable`、history `commit_status_state=null`。 |
| GitHub status 送信失敗 | build 成否は維持、WARN と固定 error reason を保存する。 |
| 無効 | status API 呼び出し 0 回、`commit_status.enabled=false`。 |
| description 長文 | 140 文字以内に切り詰められる。 |
| pending 失敗 | build 継続、pending failure は WARN と effects で検証し、未定義 build log key を保存しない。 |
| final 失敗 | final に送信しようとした state と固定 error reason を保存し、build result は反転しない。 |

---

<a id="phase-13-commitstatus-alignment-contract"></a>
**Phase 13 commitstatus 実装整合契約：**

[`docs/details/commitstatus.md`](commitstatus.md) 詳細本文責務は、Phase 13 で GitHub Commit Status payload、送信、retry、rate limit、timeout、結果正規化、保存境界を所有する。Phase 13 の対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry)、完了証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) を参照する。

| 対象 | 固定契約 |
|------|----------|
| 責務集約 | GitHub Commit Status の HTTP request、payload 生成、response 判定、retry 可否、rate limit 判定、timeout、保存値正規化は `commitstatus` owner だけが実装する。`runner` は build lifecycle 上の呼出時点と結果受け取りだけを扱う。 |
| retry / rate limit | `401`、`403`、`404`、`429`、`5xx`、network、timeout の分類は [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) の error reason 固定契約へ一致させる。Phase 13 実装で retry を追加する場合は、retry 対象、回数、待機、最終保存値を同節へ先に追加する。 |
| timeout | request 全体 timeout、response body 上限、body discard 条件、HTTP client redirect 禁止を `commitstatus` owner の実装内で固定し、呼び出し元が個別 timeout を注入して挙動を変えない。 |
| secret safety | Authorization header、GitHub token、credential 付き URL、GitHub response body 全体を build log、history、server log、fixture expected に保存しない。secret mask の本文は [`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](security.md#phase-13-security-alignment-contract) を参照する。 |
| 状態保存 | `.build_logs`、`.build_history` への保存は [`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](statefile.md#phase-13-statefile-alignment-contract) の statefile 経路だけを使用する。 |
| duplicate 排除 | `components/runner/`、`components/api/`、`components/mcp/`、`components/archive/` に Commit Status HTTP payload 生成、GitHub status endpoint 直呼び、同等 retry、同等 error reason table を残さない。 |

Phase 13 の Commit Status result record は `commit_sha`、`context`、`state`、`target_url_hash`、`description`、`attempt_count`、`final_http_status`、`final_error_reason`、`rate_limit_reset_at`、`timeout_ms`、`secret_masked` を持つ。`state` は GitHub Commit Status API が許可する状態だけを保存し、runner 独自状態を保存しない。`target_url_hash` は URL raw value を保存せず canonical URL の SHA-256 lowercase hex とする。

Commit Status 送信失敗は runner の build result を成功へ反転してはならない。送信失敗が build failure を意味するか、status notification failure を意味するかは `final_error_reason` と runner finalizer record の両方で区別する。

Phase 13 の `commitstatus` 実装は、[`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](fixture.md#phase-13-implementation-alignment-quality-evidence) の `phase13_archive_commitstatus_ownership_open_count=0`、`phase13_external_boundary_open_count=0` を満たすまで完了扱いにしてはならない。
