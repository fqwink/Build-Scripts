# Adlaire CI — Commit Status 詳細仕様

本ファイルは `ADLAIRE_CI_DETAIL_SPEC.md` および `ADLAIRE_CI_DETAIL_RUNNER_SPEC.md` から分割した `commitstatus` owner component の詳細仕様である。

本ファイルに、方針、ポリシー、実装状態、正本関係、ロードマップ状態、実装可否の上位判断を記載してはならない。これらは `ADLAIRE_CI_SPEC.md` を正とする。

本ファイルを読む前に、`ADLAIRE_CI_SPEC.md` で実装状態と実装可否を確認し、`ADLAIRE_CI_DETAIL_SPEC.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `commitstatus` owner component の主本文であり、collaborator component の仕様は呼び出し境界、状態、fixture、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `commitstatus` |
| collaborator component | `runner`、`statefile` |
| 持つ内容 | GitHub Commit Status API payload、送信順、失敗時非反転、保存値、secret mask、検証条件。 |
| 持たない内容 | runner の build 実行判断、GitHub read、API endpoint、SDK method、UI DOM 詳細。 |

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
| `state` | 開始時 `"pending"`、成功時 `"success"`、失敗時 `"failure"`、設定/状態書込など CI 自体の異常時 `"error"`。 |
| `context` | `.server_config.commit_status_context`。既定値 `"Adlaire CI"`。 |
| `description` | 140 文字以内。開始時 `Build started`、成功時 `Build succeeded`、失敗時 `Build failed: <target_status>`、pending deploy 時 `Build succeeded with deploy pending`。 |
| `target_url` | `.server_config.commit_status_target_url` が `null` でなければ送信する。 |

処理順序は以下とする。

1. runner が commit SHA を確定する。
2. runner が build id を採番する。
3. `.build_status.json status="running"` を書く前に、runner が commitstatus component を呼び出し、GitHub status `"pending"` を送信する。
4. pipeline / deploy / snapshot / history の最終結果確定後、runner が commitstatus component を呼び出し、GitHub status の最終 state を送信する。
5. `.build_logs/{id}.json.commit_status` に最終送信結果を保存する。
6. `.build_history.commit_status_state` に最終 state を保存する。

Commit Status 送信失敗は build 成否を反転させない。送信失敗時は WARN ログ `COMMIT_STATUS_FAILED: status=<http_status> error=<reason>` を出し、`.build_logs/{id}.json.commit_status.state` を最後に送信しようとした state、`error` を固定文言で保存する。GitHub API 認証失敗、403、404、5xx、network error はすべて送信失敗として扱う。

**Commit Status 保存固定契約：**

| 項目 | 仕様 |
|------|------|
| pending 送信成功 | `.build_logs/{id}.json.commit_status.pending_sent=true` を保存する。 |
| pending 送信失敗 | `pending_sent=false`、`pending_error` に固定文言を保存し、build は継続する。 |
| final 送信成功 | `final_sent=true`、`state` に最終 state を保存する。 |
| final 送信失敗 | `final_sent=false`、`state` に送信しようとした最終 state、`error` に固定文言を保存する。 |
| description | 140 文字を超える場合は 137 文字 + `...` に切り詰める。改行は空白へ置換する。 |
| target_url | `null` または空文字の場合は payload から省略する。`http` / `https` 以外は設定 validation で `422`。 |
| secret | GitHub token、Authorization header、response body 全体は build log に保存しない。 |

検証条件:

| ケース | 期待結果 |
|--------|----------|
| 有効・成功 build | fake GitHub server に `pending` → `success` の順で 2 回送信される。 |
| 有効・pipeline 失敗 | `pending` → `failure` の順で送信される。 |
| deploy pending | 最終 state は `success`、description は deploy pending を含む。 |
| commit SHA なし | status 送信なし、build は継続、log に `commit sha unavailable`。 |
| GitHub status 送信失敗 | build 成否は維持、WARN と `commit_status.error` を保存する。 |
| 無効 | status API 呼び出し 0 回、`commit_status.enabled=false`。 |
| description 長文 | 140 文字以内に切り詰められる。 |
| pending 失敗 | build 継続、`pending_sent=false`。 |
