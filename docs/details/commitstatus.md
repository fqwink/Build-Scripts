# Adlaire CI — Commit Status 詳細仕様

本ファイルは `commitstatus` owner component の詳細本文責務の正本である。

本ファイルの詳細本文境界管理条件は [`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 に従う。本ファイルは `commitstatus` owner component の主本文であり、collaborator component の仕様は呼び出し境界、状態、fixture、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `commitstatus` |
| collaborator component | `runner`、`statefile` |
| 持つ内容 | `commitstatus` owner が主本文として定義する GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| 持たない内容 | runner の build 実行判断、GitHub read、API endpoint、SDK method、UI DOM 詳細、状態 schema、setup / release 手順、fixture / PR 証跡責務。 |

---

## 対象範囲

| 範囲 | 内容 |
|------|------|
| §27.1 | GitHub Commit Status API。 |

---

### 27.1 GitHub Commit Status API

owner component は `commitstatus` とする。collaborator component は `runner`、`statefile` とする。

runner は `.server_config.commit_status_enabled == true` の場合、commitstatus component を呼び出し、対象 commit に GitHub Commit Status API を送信する。送信先は `POST /repos/{owner}/{repo}/statuses/{sha}` とし、`sha` は対象 target の commit SHA を使用する。blob SHA しか取得できない場合は status を送信せず、`.build_logs/{id}.json.commit_status.error="commit sha unavailable"` を保存する。

送信 payload は次に固定する。

| key | 値 |
|-----|----|
| `state` | 開始時 `"pending"`、成功時 `"success"`、失敗時 `"failure"`、GitHub API 認証失敗、commit status 設定不正、状態ファイル読込失敗、状態ファイル書込失敗、payload 生成不能のいずれかの CI 自体の異常時 `"error"`。 |
| `context` | `.server_config.commit_status_context`。既定値 `"Adlaire CI"`。 |
| `description` | 140 文字以内。開始時 `Build started`、成功時 `Build succeeded`、失敗時 `Build failed: <target_status>`、pending deploy 時 `Build succeeded with deploy pending`。 |
| `target_url` | `.server_config.commit_status_target_url` が `null` でなければ送信する。 |

GitHub request は以下に固定する。retry、GraphQL API、Check Runs API、任意 header、任意 payload key は本詳細仕様の対象外とする。

| 項目 | 固定値 |
|------|--------|
| method | `POST` |
| path | `/repos/{owner}/{repo}/statuses/{sha}` |
| header | `Accept: application/vnd.github+json`、`Authorization: Bearer {GITHUB_TOKEN}`、`X-GitHub-Api-Version: 2022-11-28` |
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
| `context` | `.server_config.commit_status_context` | 1〜100 文字、前後空白は validation error。 | config validation で `422`。runner 実行時に不正を検出した場合は送信せず `state="error"`。 |
| `target_url` | `.server_config.commit_status_target_url` | `null`、空文字、`http://`、`https://` のいずれか。空文字は `null` と同じ。 | config validation で `422`。runner 実行時に不正を検出した場合は送信せず `state="error"`。 |
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

`.build_logs/{id}.json.commit_status` は [`docs/details/statefile.md`](statefile.md) §22.0c の CommitStatus object に定義された key だけを保存する。`pending_sent`、`pending_error`、`final_sent` など未定義 key を保存してはならない。pending / final の送信順、HTTP request、response、失敗有無は fixture の `expected/effects.json` と server WARN log で検証し、build log schema へ未定義 key を追加しない。

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
| invalid owner / repo / context / target_url | 上表の固定 validation error | `COMMIT_STATUS_FAILED: status=0 error={reason}` |
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
| state write failure | runner の既存 state write failure 契約に従う。 | Commit Status 失敗を理由に未定義 rollback を実行すること。 |

**Commit Status fixture expected 固定契約：**

| expected file | 必須内容 |
|---------------|----------|
| `input/fakes.json` | fake GitHub Status API の応答順、HTTP status、response body mask、network failure 指定。 |
| `expected/effects.json.external_calls` | `method`、`path`、`headers_present`、`body`、`order`、`result_status`、`error_reason`。Authorization 値は保存しない。 |
| `expected/logs/build-log.json` | `commit_status` object。key は `enabled`、`state`、`context`、`target_url`、`sent_at`、`http_status`、`error` だけ。 |
| `expected/logs/history.jsonl` | `commit_status_state` と build result が非反転であること。 |
| `expected/security.json` | `forbidden_plaintexts` に GitHub token、Authorization header、credential 付き URL、GitHub response body を含める。 |

検証条件:

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
