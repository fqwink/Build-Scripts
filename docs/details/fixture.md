# Adlaire CI — Fixture 詳細仕様

本ファイルは `docs/DETAIL_INDEX.md` から分割した fixture / fake / testdata / assertion / PR 証跡の詳細仕様である。

本ファイルに、方針、ポリシー、正本関係、実装状態、ロードマップ状態、実装可否の上位判断を記載してはならない。方針、ポリシー、正本関係は `docs/SPEC.md`、実装状態、ロードマップ状態、実装可否は `docs/ROADMAP.md` を正とする。

本ファイルを読む前に、`docs/SPEC.md` で方針とポリシーを確認し、`docs/ROADMAP.md` で実装状態と実装可否を確認し、`docs/DETAIL_INDEX.md` §0〜§0j で共通固定値、責務 component、詳細節対応表、リポジトリ内ソース配置を確認する。本ファイルは `fixture` owner component の主本文であり、collaborator component の仕様は fixture 入力、expected、effects、assertion、PR 証跡、検証観点として参照する。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| owner component | `fixture` |
| collaborator component | `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`security`、`setup` |
| 持つ内容 | `fixture` owner が主本文として定義する fixture manifest、assertion、fake、testdata、expected / effects、受け入れ fixture 共通契約、PR 証跡テンプレート、acceptance checklist、差し戻し条件。 |
| 持たない内容 | 個別 component の通常処理本文、API endpoint 詳細、SDK method 実装、UI DOM 詳細、状態 schema、setup / release 実行手順。 |

---

## 対象範囲

| 範囲 | 内容 |
|------|------|
| §0g.8-F | Phase fixture / testdata 配置、fake 実装、実装 PR 証跡。 |
| §22-F | Phase 3 / Phase 4 API の必須検証、API fixture、API / SDK / UI / 状態ファイル cross fixture 固定。 |
| §27-F | §27 fixture 配置、fixture カタログ、manifest、assertion、expected/effects、相互整合、component 別検証責務。 |
| §27-F-PR | §27 実装 PR 証跡、受け入れゲート、差し戻し条件、部分失敗・再実行契約。 |

## 0g.8-F Phase fixture / testdata / fake / PR 証跡契約

本ファイルは、実装完了判定に必要な fixture、fake、testdata、expected / effects、PR 証跡、acceptance checklist、差し戻し条件の正本である。`docs/DETAIL_INDEX.md` §0e、§0g、§0i は完了判定の入口を示すだけとし、`docs/details/setup.md` §26 は setup / release / Phase 判定の実行条件を示すだけとする。fixture 名、expected / effects、fake 動作、PR 証跡項目、不足時の扱い、差し戻し条件は本ファイルを正とする。

実装 PR の完了証跡は、対象に応じて以下の 3 系統に分類する。複数系統にまたがる PR は、該当する全系統の証跡を PR 本文または検証ログに記録する。

| 系統 | 対象 | 正本節 | 必須証跡 |
|------|------|--------|----------|
| Phase 実装 | Phase 1〜Phase 6 の初期実装。 | §0g.8-F | 対象 Phase、owner component、collaborator component、変更ファイル、fixture / testdata path、fake、実行コマンド、期待結果、実結果、後続 Phase へ引き継ぐ contract。 |
| API 実装 | Phase 3 / Phase 4 の API、SDK / UI / statefile と同期する API 実装。 | §22-F | Phase、endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、状態副作用、secret mask、GET 副作用なし確認。 |
| §27 実装 | §27.1〜§27.47 の追加仕様化機能。 | §27-F | 対象 §27.x、関連 §22 / §23 / §24 / §25 / §26、owner / collaborator component、fixture 名、状態差分、外部副作用、partial failure、再実行、対象外確認。 |

上表の証跡が不足する場合、対象機能は未完了として扱う。実装者は fixture が pass したことだけを完了証跡として扱ってはならない。

Phase、API、§27 のいずれの実装 PR でも、証跡の記録形式は本ファイルの表に従う。component 別詳細仕様ファイル、`docs/DETAIL_INDEX.md`、`docs/details/setup.md` に同種の記録項目がある場合は、本ファイルの証跡分類、不足時の扱い、差し戻し条件を優先する。

**Phase fixture / testdata 配置固定契約：**

| Phase | 必須配置 | 必須内容 | 禁止事項 |
|-------|----------|----------|----------|
| Phase 1 | `testdata/builder/single/`、`testdata/builder/site/`、`testdata/builder/empty-dir/`、`testdata/builder/strict/`、`testdata/builder/safe/`、各 fixture の `expected/`。 | 入力 Markdown、テーマ設定、asset 入力、期待 HTML / CSS / JS / search index、期待 stdout / stderr、期待終了コード。 | 実行環境ごとに変わる絶対 path、timestamp、乱数、外部 URL 取得結果を期待値へ含めてはならない。 |
| Phase 2 | `testdata/runner/r1/`〜`testdata/runner/r20/`、各 fixture の `state/`、`github/`、`pipeline/`、`ssh/`、`notify/`、`expected/`。 | GitHub fake response、状態ファイル初期値、lock 状態、pipeline fake 結果、deploy fake 結果、通知 fake 結果、期待 `.last_sha`、期待 queue / snapshot。 | 実 GitHub API、実 SSH、実通知先、実 remote branch 状態に依存して合否を決めてはならない。 |
| Phase 3 | `testdata/api/phase3/auth/`、`status/`、`history/`、`logs/`、`queue/`、`stream/`、`errors/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | HTTP method / path / query / header / body、状態ファイル初期値、期待 response、期待 error body、SSE frame、状態 read/write 後の期待値。 | 仕様未定義 endpoint、Phase 4 endpoint、外部公開設定を fixture に含めてはならない。 |
| Phase 4 | `testdata/api/phase4/config/`、`notify/`、`snapshots/`、`maintenance/`、`hooks/`、`tokens/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | config / notify / snapshot / rollback / maintenance / hook / token の正常系、validation error、secret mask、Phase 3 互換確認の期待値。 | token 原文、secret 原文、mask 前 payload、再取得不可 token の復元値を fixture または expected に含めてはならない。 |
| Phase 5 | `testdata/sdk/request-shape/`、`error-shape/`、`stream/`、`binary/`、`phase4/`。 | fake fetch transcript、期待 request、期待 SDK return、期待 `AdlaireCIError`、期待 stream event、timeout / abort の期待結果。 | Node.js 専用 API、bundler、npm package、実 network、browser storage 依存を検証前提にしてはならない。 |
| Phase 6 | `testdata/ui/login/`、`status/`、`build/`、`config/`、`secret/`、`stream/`、`phase4/`。 | fake SDK script、入力 DOM 状態、操作手順、期待 DOM assertion、期待 SDK call、期待 disabled / loading / error / success 表示。 | 直接 `fetch()`、CDN、外部 framework、画像 snapshot だけの合否判定、secret 表示を含めてはならない。 |

**fake 実装固定契約：**

| fake | 対象 Phase | 必須動作 | 必須記録 |
|------|------------|----------|----------|
| fake GitHub server | Phase 2 | fixture の JSON 応答だけを返す。未定義 method / path は `404` とする。rate limit、`304`、`409`、`500` は fixture で明示された場合のみ返す。 | method、path、query、request body、認証 header の有無、呼び出し順。token 値は `***` に置換する。 |
| fake ssh executable | Phase 2 | fixture 指定の stdout、stderr、終了コード、timeout を返す。実 shell、実 SSH、実 file 転送は実行しない。 | argv、stdin 有無、環境変数名、終了コード、timeout 発生有無。secret 値は記録しない。 |
| fake notifier | Phase 2 / Phase 4 | fixture 指定の HTTP status、response body、timeout を返す。通知先へ送信しない。 | URL の host 部分、payload schema、mask 後 payload、retry 回数、最終結果。 |
| fake filesystem | Phase 1〜Phase 4 | atomic write 失敗、sync 失敗、lock 競合、JSON 破損、permission error を fixture 単位で再現する。通常 file I/O の代替にはしない。 | 対象 path、操作種別、注入した失敗、復旧後の状態。 |
| fake fetch | Phase 5 | `status`、`headers`、`body`、network error、timeout、abort、stream chunk を fixture どおり返す。実 network は使用しない。 | method、URL、query、headers、body 有無、abort 発生有無、呼び出し順。token 値は `***` に置換する。 |
| fake SDK | Phase 6 | SDK method ごとに固定 return、固定 throw、固定 stream event を返す。UI からの直接 API 呼び出しは受け付けない。 | method 名、引数、呼び出し順、throw した error code、stream unsubscribe 実行有無。 |

**実装 PR 証跡固定契約：**

本契約は Phase 実装の PR 証跡に適用する。API 実装は §22-F、§27 実装は §27-F の acceptance checklist も同時に満たす。

| 証跡 | 必須記載 | 不足時の扱い |
|------|----------|--------------|
| 変更対象 | 対象 Phase、owner component、collaborator component、変更ファイル、追加 fixture / testdata path。 | 対象 Phase の成果物不足として未完了。 |
| 固定契約 | 追加または固定した CLI、状態 schema、HTTP API、SDK method、DOM id、fake 動作、終了コード、error body。 | 後続 Phase が参照できないため未完了。 |
| 検証 | 実行コマンド、fixture 名、期待結果、実結果、判定。 | 合否を再現できないため未完了。 |
| 未実装対象 | 対象 Phase 外の機能、将来計画、MCP、外部公開設定、未定義 endpoint / UI / 状態ファイルのうち今回実装しない範囲を PR 証跡に列挙する。 | 先取り実装または範囲不明として未完了。 |
| 後続 Phase への影響 | 後続 Phase が利用許可済みの contract と、利用禁止の未固定 contract を PR 証跡に列挙する。 | 次 Phase 着手条件未充足として未完了。 |
| secret 確認 | log、fixture、snapshot、UI 表示、PR 本文に secret / token / password 原文がないこと。 | security 不合格として未完了。 |

## 22-F Phase 3 / Phase 4 API fixture 契約

本節は、Phase 3 / Phase 4 API の必須検証、fixture 名、入力状態、期待 response、期待副作用の正本である。API endpoint の method、path、request、response、error、read / write 境界は `docs/details/api.md` §22 を正とし、本節では再定義しない。

API 実装 PR は、§0g.8-F の PR 証跡固定契約に加えて、本節の Phase、endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、状態副作用、secret mask、GET 副作用なし確認を記録する。これらの記録が不足する場合、API 実装は完了扱いにしない。

対象項目の実装時は、Phase 3 を完了してから Phase 4 へ進める。同一 Phase 内では、API、SDK、UI、状態ファイル、検証手順を同じ Pull Request で同期する。

| Phase | 対象 | 完了条件 |
|-------|------|----------|
| Phase 3 | 認証、セッション、共通エラー、状態ファイル読み書き、`.access_log`、`.config_log`、ビルド操作、status、logs、history、queue、circuit breaker | `POST /api/login` から認証必須 API の共通処理、手動ビルド、強制ビルド、キャンセル、キュー、履歴、ログ取得までが `docs/details/api.md` §22.0〜§22.0e と一致し、秘密情報がログとレスポンスに出ない。 |
| Phase 4 | config、repo、branch、schedule、PAT、diagnostics、dashboard、notify、SMTP、webhook、snapshot、rollback、maintenance、access control、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout、tokens | 拡張運用 API が schema どおり状態を保存し、secret mask、GET 副作用なし、rollback / maintenance / hook / token の副作用が fixture と一致し、SDK と UI の操作名が `docs/details/api.md` §22.0e と一致する。 |

各 Phase の検証条件は以下とする。

| Phase | 必須検証 |
|-------|----------|
| Phase 3 | 認証成功、認証失敗、期限切れ session、`401` 時 SDK token 破棄、`.access_log` 追記、秘密情報マスク、手動 build、force build、running 中の queue、cancel、history/log 取得、`409`、`429`、`503` を確認する。 |
| Phase 4 | config/repo/branch/schedule の保存、`.config_log` 追記、GET 系 API の無副作用、Webhook test、weekly summary、SMTP test、webhook secret 保存、secret mask、snapshot list/download/delete、rollback、maintenance enable/disable、access control block、hook success/failure、rule 追加/削除、pipeline config 保存、notes 保存、dashboard layout 保存、token 発行/失効、token 本体が再取得不可であることを確認する。 |

**Phase 3 API fixture 固定：**

Phase 3 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は実装言語の test case 名または subtest 名へそのまま写せる粒度とし、期待 HTTP status、期待 body、状態ファイル副作用を同時に確認する。

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

**Phase 4 API fixture 固定：**

Phase 4 実装は、下表の fixture をすべて満たした場合だけ完了扱いにする。fixture は既存 endpoint と既存状態ファイルだけを対象とし、`docs/details/api.md` §22.0e にない endpoint、`docs/details/statefile.md` §22.0a にない状態ファイル、`docs/details/ui.md` §24 にない UI 操作を追加してはならない。

| Fixture | Phase | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------|--------------------|---------------|--------------------|
| B1 config no-op | Phase 4 | 既存 `.server_config` と同じ body で `POST /api/config`。 | `200` と `{message:"No changes",config}`。 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 |
| B2 config validation failure | Phase 4 | `queue_max_size=-1`、未知 enum、相対 path を含む `POST /api/config`。 | `422 {"error":"Validation failed","details":[...]}`。 | 状態ファイルを変更しない。 |
| B3 branch config save | Phase 4 | 有効な branch target 2 件で `POST /api/branch-config`。 | `{message:"Branch config updated",branches_count:2}`。 | `.branch_config` を atomic write し、`.config_log` に差分を記録する。 |
| B4 schedule systemd failure | Phase 4 | `.server_config` 保存成功後、systemd timer 更新を fake failure にする。 | `500`。 | `.server_config` は更新済み、`.config_log` に `systemd_update_failed` を記録し、未定義 rollback を行わない。 |
| B5 PAT update secret mask | Phase 4 | `POST /api/pat-update` に token を送る。 | `{message:"PAT updated"}`。 | secret file は mode `600`。response、`.config_log`、`.audit_log`、server log に token 平文を出さない。 |
| B6 dashboard read-only | Phase 4 | `.dashboard_layout`、`.build_state`、`.build_status.json`、`.alert_rules` を置き `GET /api/dashboard`。 | widget 順に dashboard object を返す。 | GET は対象状態ファイルを作成、修復、更新しない。 |
| C1 notify config mask | Phase 4 | Webhook secret と SMTP password を含む通知設定保存後、GET / backup / log を確認する。 | secret は `"***"` または `*_set:true` だけを返す。 | secret 平文を状態表示、履歴、通知ログ、backup に残さない。 |
| C2 webhook receive signed | Phase 4 | 正常署名の GitHub push payload を `POST /api/webhook`。 | `202` と queued 結果。 | `.webhook_events.json` 追記後、queue 投入条件を満たす場合だけ `.build_state.queued` へ `trigger:"webhook"` を追加する。 |
| C3 webhook invalid signature | Phase 4 | 署名なし、不正 prefix、不一致署名。 | `401 {"error":"Unauthorized"}`。 | event log、queue、history を変更しない。 |
| C4 SMTP test disabled | Phase 4 | SMTP disabled で `POST /api/smtp-test`。 | endpoint 固有の `422`。 | `.notify_log` へ成功扱いを残さず、secret を出力しない。 |
| D1 snapshot delete | Phase 4 | 存在する snapshot id で `DELETE /api/snapshots/{id}`。 | `{message:"Snapshot deleted"}`。 | 対象 snapshot だけ削除し、`.config_log` または監査対象 log に削除を記録する。 |
| D2 rollback running conflict | Phase 4 | `.build_state.running=true` で `POST /api/history/{id}/rollback`。 | `409 {"error":"Build is running"}`。 | queue、history、snapshot、deploy target を変更しない。 |
| D3 maintenance blocks build | Phase 4 | maintenance enabled 状態で `POST /api/build`。 | `503` と endpoint 固有 maintenance error。 | build queue、history、log を変更しない。 |
| D4 access-control deny | Phase 4 | allowlist に接続元が含まれない状態で任意認証必須 API。 | 認証判定前に `403 {"error":"Forbidden"}`。 | password / token 検証、access log 成功行、対象操作副作用を行わない。 |
| D5 hook timeout | Phase 4 | `pre` hook が timeout。 | build status は `hook_error` または endpoint 固有の hook error。 | pipeline を実行せず、hook log と build log に timeout を固定値で記録する。 |
| E1 token issue once | Phase 4 | `POST /api/tokens` で token 作成。 | token 本体を作成 response に 1 回だけ含める。 | `.api_tokens` には hash だけを保存し、再取得 API では token 本体を返さない。 |
| E2 token revoke missing | Phase 4 | 存在しない token id を `DELETE /api/tokens/{id}`。 | `404 {"error":"Not found"}`。 | `.api_tokens`、`.audit_log` を変更しない。 |
| E3 alert/tag duplicate | Phase 4 | 同一 alert rule または tag rule を 2 回作成。 | 2 回目は `409 {"error":"Conflict"}`。 | 2 回目は該当状態ファイルと `.config_log` を変更しない。 |
| E4 pipeline config reserved arg | Phase 4 | `extra_args` に `--src`、`--out`、`--state-dir` を含める。 | `422 {"error":"Validation failed","details":[...]}`。 | `.pipeline_config` を変更しない。 |
| E5 notes same content | Phase 4 | 同じ `content` を 2 回 `POST /api/notes`。 | 2 回目は `No changes`。 | 2 回目は `.notes`、`.config_log` を変更しない。 |
| E6 dashboard layout invalid | Phase 4 | 重複 widget、未知 widget、空配列を `POST /api/dashboard-layout`。 | `422 {"error":"Validation failed","details":[...]}`。 | `.dashboard_layout` を変更しない。 |

**API 機能別 fixture 固定契約：**

下表は、API component の機能別 fixture 名、入力、期待結果を固定する。API endpoint の処理順序、request / response、状態ファイル read / write 境界は `docs/details/api.md` を正とし、本表では fixture 本体だけを定義する。

| 機能 | fixture | 入力 | 期待結果 |
|------|---------|------|----------|
| backup / restore | backup-mask | secret 設定済みで backup | secret 本体なし、`*_set:true`。 |
| backup / restore | restore-validate-fail | 1 file schema 不正 | `422`、全 file 差分なし。 |
| backup / restore | restore-secret-keep | `"***"` かつ既存 secret あり | 既存 secret 維持、平文出力なし。 |
| backup / restore | restore-secret-missing | `"***"` かつ既存 secret なし | secret 未設定のまま、`"***"` を保存しない。 |
| backup / restore | restore-secret-delete | secret `null` | secret file 削除。 |
| backup / restore | restore-write-failure | 中途 write 失敗 | 未処理 file は差分なし、処理済み file は維持、`500`。 |
| メンテナンス | maintenance-build-deny | enabled 中に `POST /api/build` | `503`、queue 差分なし、build id なし。 |
| メンテナンス | maintenance-force-deny | enabled 中に `POST /api/build/force` | `503`、SHA cache 差分なし。 |
| メンテナンス | maintenance-webhook-deny | enabled 中に署名済み webhook | `503`、event log と queue 差分なし。 |
| メンテナンス | maintenance-disable-noop | disabled 中に disable | `200 No changes`、状態ファイル差分なし。 |
| アクセス制御 | access-allow-empty | `.access_control.allow=[]` | 任意 IP の API が認証処理へ進む。 |
| アクセス制御 | access-deny-before-auth | allow 不一致 IP で `POST /api/login` | `403`、`.access_log`、`.audit_log`、rate state 差分なし。 |
| アクセス制御 | access-normalize | 重複 allow を保存 | sort / 重複除去後の配列を返し `.config_log` 記録。 |
| アクセス制御 | access-ipv6-reject | IPv6 literal を保存 | `422`、状態差分なし。 |
| hooks | hook-pre-success | pre hook exit 0 | pipeline 実行、hook log 保存、secret mask 済み。 |
| hooks | hook-pre-abort | pre hook exit 1 / abort true | pipeline 未実行、status `hook_error`、history に `failure_category:"hook_error"`。 |
| hooks | hook-pre-warn | pre hook exit 1 / abort false | build 継続、hook log に exit code。 |
| hooks | hook-post-failure | build success 後 post hook failure | build success 維持、hook log 保存。 |
| hooks | hook-timeout | timeout 超過 | process kill、`timed_out:true`、`exit_code:null`。 |
| hooks | hook-log-write-failure | pre hook log 保存失敗 | build 本体未実行、`hook_error`。 |
| hooks | post failure | build status を変更しない。WARN log と hook log だけを残す。 |
| SMTP | smtp-save-password | password 付き保存 | `.smtp_secret` mode `0600`、GET は `password_set:true`、log は `"***"`。 |
| SMTP | smtp-delete-password | `password:null` | `.smtp_secret` 削除、password 平文なし。 |
| SMTP | smtp-noop | 同一 config / password 未指定 | 状態差分なし、`.config_log` 追記なし。 |
| SMTP | smtp-test-success | 設定済み test | `.notify_log` に success、response success。 |
| SMTP | smtp-test-disabled | `enabled:false` | `422`、`.notify_log` 差分なし。 |
| SMTP | smtp-log-failure | test 後 `.notify_log` 追記失敗 | `500`、password 平文なし。 |
| queue | queue-add-running | running 中に manual build | queue append、created_seq 最大 + 1。 |
| queue | queue-duplicate | 同一 manual payload を再投入 | 新規追加なし、既存 queue_id を返す。 |
| queue | queue-full | max_size 到達 | `429 {"error":"queue_full"}`、差分なし。 |
| queue | queue-clear | waiting 2 件で `DELETE /api/queue` | `cleared_count=2`、running/current_build_id 維持。 |
| queue | queue-runner-take | urgent と normal が混在 | urgent を削除し running に設定、他 entry 維持。 |
| 認証 | auth-password-failure | 誤 password で `POST /api/login` | `401`、session/ticket なし、失敗回数 +1、`.access_log` と `.audit_log` に secret なし。 |
| 認証 | auth-login-lock | 連続 10 回失敗後の `POST /api/login` | `429`、password hash 検証なし、`.access_log` に `login_locked`、`.audit_log` に `permission_denied`。 |
| 認証 | auth-session-issued | TOTP 無効で password 成功 | token は response のみ、`.admin_credentials.login_count` +1、ログに token/hash/salt なし。 |
| 認証 | auth-session-expired | 期限切れ session で保護 API | `401`、対象 session 削除、`.access_log` と `.audit_log` は追記しない。 |
| 認証 | auth-password-change | password 変更成功 | 新 salt/hash、現 session 以外削除、`password_change` ログ、password/hash/salt 平文なし。 |
| 認証 | auth-log-write-failure | login 成功時に `.audit_log` 追記失敗 | `500`、session token を response しない。 |
| approval | approval-create | approval_required target に差分 | build なし、pending record、通知成功または pending。 |
| approval | approval-duplicate | 同一 branch/sha/target を再検出 | pending 重複作成なし。 |
| approval | approval-approve | pending approve | queue 追加、approved record、queue_id 保存。 |
| approval | approval-reject | pending reject | rejected record、history `approval_rejected`。 |
| approval | approval-timeout | expires_at 超過 | expired record、history `approval_expired`。 |
| approval | approval-queue-full | max_size 到達時 approve | `429`、status pending 維持。 |
| approval | approval duplicate notify | 重複時は通知を送らない。 |
| approval | approved append failure | queue は残り、API は `500`。 |

**API / SDK / UI / 状態ファイル cross fixture 固定：**

下表の fixture は、API endpoint、SDK method、UI 操作、状態ファイル副作用の横断整合を固定する。API endpoint の method、path、request、response、error、read / write 境界は `docs/details/api.md`、SDK method と error 変換は `docs/details/sdk.md`、UI DOM と表示状態は `docs/details/ui.md`、状態ファイル schema と保存手順は `docs/details/statefile.md` を正とする。

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

## 27-F §27 fixture / PR 証跡詳細契約

**§27 fixture 配置・命名固定契約：**

§27 の fixture は、実装者が実行順や期待値を推測しないように、下表の単位で配置する。実装ファイル名やテストフレームワークは本節では固定しないが、fixture の入力、期待出力、期待副作用は機能単位で分離する。

| 機能範囲 | fixture 配置単位 | 必須内容 | 禁止事項 |
|----------|------------------|----------|----------|
| §27.1〜§27.4 | `status/`、`dry-run/`、`retry/`、`output-meta/` 相当の機能別単位。 | runner 入力、GitHub fake response、builder 入力、期待 log/history/status/report。 | 実際の GitHub API、実時刻依存の期待値、secret 平文。 |
| §27.5〜§27.11 | `config/`、`access-log/`、`archive/`、`schedule/` 相当の API 機能別単位。 | HTTP request、初期状態、期待 response、期待状態差分、期待 log。 | API response だけの検証、状態差分未確認。 |
| §27.12〜§27.20 | `webhook/`、`stats/`、`snapshot/`、`health/`、`branch-config/`、`summary/`、`config-log/` 相当の機能別単位。 | 署名、payload、query、snapshot 入力、破損行、期待 paging、期待 rollback。 | 署名検証省略、破損行の黙殺仕様未確認。 |
| §27.21〜§27.38 | `runner-extensions/` 配下の機能別単位。 | target、pipeline、cache、queue、hook、remote、approval、notification、trend の正常/異常/部分失敗。 | 実行完了順依存、外部 shell 展開、未定義状態ファイル。 |
| §27.42〜§27.47 | `security/` 配下の scope、token、audit、session、totp、rate-limit 単位。 | route 判定、body 未評価、token hash、audit failure、window reset、secret mask。 | token 本体保存、Authorization header 保存、監査なし権限拒否。 |
| §28.1〜§28.25 | `builder-extensions/` 配下の機能別単位。 | Markdown 入力、CLI option、期待 HTML / CSS / JS / REPORT、strict / non-strict の終了コード。 | 外部 library、CDN、実 network、環境依存 timestamp、画像 snapshot だけの合否判定。 |

fixture 名は `success-*`、`failure-*`、`partial-*`、`noop-*`、`security-*` のいずれかで始める。fixture 名に実行時刻、乱数、環境依存 path、実 token 値を含めてはならない。期待時刻は固定値を使い、現在時刻依存の検証では fake clock を fixture 入力に含める。

**§27 fixture カタログ固定契約：**

§27 の実装 PR は、対象機能について下表の fixture を作成する。複数節を同一 PR で実装する場合は、各節の fixture を省略せず作成する。fixture は入力、期待 response、期待 stdout/stderr、期待状態差分、期待 log/history/audit/notify、secret 非表示確認のうち該当するものを含める。

| 節 | fixture 名 | 入力 fixture | 期待出力 / 期待副作用 |
|----|------------|--------------|------------------------|
| §27.1 | `success-commit-status-pending-success`、`failure-commit-status-unavailable-sha`、`failure-commit-status-api-error`、`noop-commit-status-disabled` | server config、commit SHA、fake GitHub response、build result。 | GitHub payload、build log `commit_status`、history state、token mask、disabled 時呼び出し 0。 |
| §27.2 | `success-dry-run-changed`、`noop-dry-run-unchanged`、`failure-dry-run-github-error`、`security-dry-run-secret-mask` | CLI args、state dir、fake GitHub response。 | stdout JSON、終了コード、状態差分なし、lock なし、secret mask。 |
| §27.3 | `success-retry-after-rate-limit`、`success-retry-after-timeout`、`failure-retry-limit-exceeded`、`noop-retry-nonretryable` | retry config、attempt sequence、pipeline / deploy fake result。 | attempts 配列、retry_count、最終 status、SHA 更新有無、未実行 attempt 不作成。 |
| §27.4 | `success-output-meta-html-report-api`、`success-output-meta-empty-values`、`failure-output-meta-invalid-sha`、`failure-output-meta-invalid-time` | builder args、Markdown 入力、build meta 値。 | HTML meta、REPORT、output-meta response、終了コード、既存出力保護。 |
| §27.5 | `success-config-validate-valid`、`success-config-validate-invalid`、`failure-config-validate-unknown-key`、`security-config-validate-secret-mask` | config JSON、既存 `.server_config`。 | valid/errors/warnings、状態差分なし、unknown key `422`、secret 非表示。 |
| §27.6 | `success-api-access-log-authenticated`、`success-api-access-log-unauthorized`、`failure-api-access-log-append`、`security-api-access-log-body-mask` | HTTP request、actor、target endpoint result。 | `.api_access_log` 追記、status/duration、body 未保存、append failure 挙動。 |
| §27.7 | `success-log-archive`、`noop-log-archive-empty`、`failure-log-archive-gzip`、`success-log-cleanup` | build logs、retention config、archive target。 | gzip archive、元 log 維持/削除条件、cleanup response、破損 log skip。 |
| §27.8 | `success-build-status-running`、`success-build-status-final`、`failure-build-status-write`、`failure-build-status-corrupt-api` | runner state、build result、queue/pending/circuit state。 | `.build_status.json`、status API、finalizer、write failure 時 log/history 維持。 |
| §27.9 | `success-trigger-manual`、`success-trigger-webhook`、`success-trigger-approval`、`failure-trigger-filter-invalid` | trigger source、history query、build log。 | trigger enum 保存、history filter、未知 trigger 不保存、`422`。 |
| §27.10 | `success-startup-integrity-clean`、`success-startup-integrity-recovered`、`failure-startup-integrity-unrecoverable`、`security-startup-integrity-secret-mode` | startup state files、mode、corrupt files。 | recovery record、WARN/ERROR、build 開始可否、secret file mode 補正。 |
| §27.11 | `success-schedule-interval`、`success-schedule-pause-resume`、`failure-schedule-systemd-update`、`noop-schedule-same-value` | schedule request、server config、fake systemd。 | `.server_config`、systemd result、config log、保存済み config の再取得値。 |
| §27.12 | `success-webhook-push-queued`、`failure-webhook-invalid-signature`、`noop-webhook-duplicate-delivery`、`failure-webhook-queue-full` | GitHub headers、raw body、secret、queue state。 | webhook event、queue entry、状態差分なし条件、重複防止。 |
| §27.13 | `success-webhook-events-page`、`success-webhook-events-empty`、`partial-webhook-events-corrupt-line`、`failure-webhook-events-invalid-query` | `.webhook_events.json`、query。 | events/total、破損行除外、secret 非表示、`422`。 |
| §27.14 | `success-stats-summary`、`success-stats-timeline`、`partial-stats-corrupt-log-skip`、`failure-stats-invalid-query` | history、build logs、archive logs、query。 | stats response、rounding、WARN、read-only 差分なし。 |
| §27.15 | `success-snapshot-list-download`、`success-snapshot-delete`、`success-snapshot-rollback`、`failure-snapshot-running-conflict` | `.snapshots/`、history id、running state。 | tar.gz download、delete 差分、rollback history/log、不正 entry 防止。 |
| §27.16 | `success-health-ok`、`success-health-degraded`、`failure-health-read-error`、`noop-health-readonly` | status/pending/notify/process state。 | health JSON、HTTP status、degraded 判定、状態差分なし。 |
| §27.17 | `success-log-search-level`、`success-log-search-archive`、`partial-log-search-corrupt-skip`、`failure-log-search-invalid-level` | build logs、archive logs、query。 | search results、line_number、破損除外、`422`。 |
| §27.18 | `success-branch-config-get-default`、`success-branch-config-post`、`failure-branch-config-invalid-path`、`partial-branch-config-log-failure` | branch config request、existing config。 | branch_targets、config log、validation 差分なし、保存済み状態維持。 |
| §27.19 | `success-weekly-summary-auto`、`success-weekly-summary-manual`、`noop-weekly-summary-same-day`、`failure-weekly-summary-send` | notify config、history、fake webhook。 | payload、notify log、sent date 更新条件、失敗時 pending。 |
| §27.20 | `success-config-diff-simple`、`success-config-diff-nested`、`noop-config-diff-same-value`、`security-config-diff-secret-mask` | before/after config、request。 | diff/diff_text、config log、no-op 差分なし、secret mask。 |
| §27.21 | `success-multi-file-one-change`、`success-multi-file-many-change`、`noop-multi-file-all-skip`、`failure-multi-file-path-traversal` | branch target、target files、SHA cache。 | build 対象集合、対象別 SHA 更新、重複排除、path error。 |
| §27.22 | `success-yaml-pipeline-file-priority`、`success-yaml-pipeline-inline`、`failure-yaml-pipeline-parse`、`security-yaml-pipeline-secret-mask` | `.pipeline.yml`、inline YAML、step fake result。 | pipeline_source、step log、終了コード、deploy/SHA 更新なし条件。 |
| §27.23 | `success-local-watch-change`、`noop-local-watch-no-change`、`failure-local-watch-state-corrupt`、`failure-local-watch-tag-filter-conflict` | local files、watch state、server config。 | local watch state、trigger、GitHub call 0、終了コード `2`。 |
| §27.24 | `success-tag-filter-match`、`noop-tag-filter-unmatched`、`failure-tag-filter-api`、`failure-tag-filter-pattern` | tag refs、patterns、SHA cache。 | matched tags、SHA 更新条件、skip 挙動、不正 pattern。 |
| §27.25 | `success-build-cache-hit`、`success-build-cache-miss`、`partial-build-cache-byte-mismatch`、`failure-build-cache-save` | Markdown、cache index/pages、theme/version。 | byte 一致、cache atomic save、hit 破棄、build 成否維持。 |
| §27.26 | `success-parallel-targets-all`、`partial-parallel-targets-some-fail`、`failure-parallel-targets-all-fail`、`success-parallel-targets-order-stable` | target list、parallel result sequence。 | 設定順 result、overall status、timeout、history/status。 |
| §27.27 | `success-hook-pre-post`、`failure-hook-pre-abort`、`partial-hook-post-fail`、`security-hook-shell-denied` | hook config、command_args、fake command result。 | hook log、secret mask、pre abort、shell 展開禁止。 |
| §27.28 | `success-dependency-manifest`、`success-dependency-missing`、`failure-dependency-build-keeps-old`、`security-dependency-path-normalize` | Markdown refs、existing manifest、build result。 | dependency manifest、成功時置換、失敗時旧 manifest 維持。 |
| §27.29 | `success-remote-build-artifact`、`failure-remote-build-auth`、`failure-remote-build-checksum`、`security-remote-build-unsafe-archive` | remote config、artifact archive、manifest/checksum。 | 一時展開、検証後 deploy、既存出力保護、unsafe entry 拒否。 |
| §27.30 | `success-approval-approve`、`noop-approval-reject`、`noop-approval-expire`、`failure-approval-double-approve` | approval queue、API action、clock。 | 状態 enum、queue 連携、物理削除なし、audit。 |
| §27.31 | `success-branch-env-match`、`noop-branch-env-default`、`security-branch-env-secret-mask`、`failure-branch-env-unknown-key` | branch config env、build command fake。 | child env、log key 名、secret mask、未定義 key 除外。 |
| §27.32 | `success-notify-build-success`、`success-notify-build-failure`、`partial-notify-send-failure`、`security-notify-secret-mask` | notify config、build result、fake channel。 | notify log/pending、build 成否維持、payload mask。 |
| §27.33 | `success-trend-append`、`success-trend-replace-same-id`、`partial-trend-retention`、`failure-trend-save` | build history/log、trend state。 | sample upsert、summary、retention、終了コード最低 `1`。 |
| §27.34 | `success-chain-linear`、`success-chain-parallel`、`failure-chain-cycle`、`partial-chain-required-fail` | chain config、job result sequence。 | chain_run_id、job log、skipped 条件、validation `422`。 |
| §27.35 | `success-priority-queue-order`、`success-priority-queue-fifo-tie`、`failure-priority-queue-full`、`partial-priority-queue-dequeue-failure` | queued entries、runner lock、API action。 | dequeue atomic、running 設定、entry 残存、`429`。 |
| §27.36 | `success-failure-category-timeout`、`success-failure-category-network`、`success-failure-category-deploy`、`noop-failure-category-success-null` | exit code、stderr、API error、status。 | category enum、unknown fallback、success null、history/log。 |
| §27.37 | `success-environment-full`、`partial-environment-null-fields`、`security-environment-secret-path`、`failure-environment-version-missing` | OS/env/version/state path inputs。 | environment object、null field、basename 保存、secret path 非表示。 |
| §27.38 | `success-duration-anomaly-avg`、`success-duration-anomaly-p95`、`noop-duration-anomaly-sample-shortage`、`partial-duration-anomaly-notify-failure` | trend summary、duration config、build result。 | anomaly tag、history flag、通知 event、重複なし。 |
| §27.42 | `security-scope-trigger-allowed`、`security-scope-read-denied`、`security-scope-path-param`、`failure-scope-audit-failure` | route/method、token record、request body。 | body 未評価、`403`、audit、Authorization 非保存。 |
| §27.43 | `security-token-create-once`、`security-token-list-mask`、`success-token-revoke`、`failure-token-expired-auth` | token request、token state、clock。 | token 1 回表示、hash 保存、revoke、期限切れ `401`。 |
| §27.44 | `success-audit-operation`、`success-audit-denied`、`failure-audit-append`、`security-audit-secret-mask` | actor、target、result、operation input。 | JSON Lines、mask、必須 audit failure、paging。 |
| §27.45 | `success-session-active`、`failure-session-expired`、`success-session-timeout-update`、`success-session-revoke-all` | session state、clock、timeout config。 | last_seen/expires_at、`401`、revoke all、token 非表示。 |
| §27.46 | `security-totp-setup-once`、`success-totp-confirm`、`failure-totp-code-reuse`、`success-totp-disable` | setup ticket、TOTP code、clock、secret state。 | secret 有効保存条件、ticket 一回使用、window、disable。 |
| §27.47 | `security-rate-limit-login`、`success-rate-limit-window-reset`、`security-rate-limit-ip-actor`、`failure-rate-limit-state-save` | policy、rate state、RemoteAddr、actor。 | count、`429`、audit 成功条件、部分 count 更新なし。 |

**§27 fixture ファイルセット固定契約：**

各 fixture は、下表のファイルセットを持つ。該当しない入出力は `not-applicable.txt` を置くのではなく、`manifest.json` の `not_applicable` 配列に理由付きで記録する。実装者は fixture ごとに必要ファイルを推測してはならない。

| ファイル | 必須 | 内容 | 禁止事項 |
|----------|------|------|----------|
| `manifest.json` | 必須 | fixture 名、対象節、機能名、分類、fake clock、owner component、collaborator component、参照仕様節、not_applicable 理由。 | 実行環境依存 path、乱数、実 secret。 |
| `input/request.json` | API / SDK / UI fixture で必須 | HTTP method、path、query、headers、body、SDK method、SDK args、UI action。 | Authorization header の実 token、secret 平文。 |
| `input/cli.json` | runner / builder fixture で必須 | binary 名、argv、env key、cwd、stdin、fake clock。 | 実 home path、実 credential path。 |
| `input/state/` | 状態参照 fixture で必須 | 実行前状態ファイル一式。存在しない状態は `manifest.json` の `missing_state` に列挙する。 | 期待状態を混ぜること、実 secret。 |
| `input/files/` | builder / artifact / hook / remote fixture で必須 | Markdown、YAML、archive、snapshot、hook 入力などの対象ファイル。 | 実外部サービスから取得した未固定ファイル。 |
| `input/fakes.json` | 外部 API / command fixture で必須 | fake GitHub、fake SSH、fake SMTP、fake webhook、fake systemd、fake command の応答順。 | 実ネットワーク呼び出し前提。 |
| `expected/response.json` | HTTP / SDK fixture で必須 | HTTP status、headers、JSON body、SDK return/error。 | 未定義 key、順序非決定配列。 |
| `expected/stdout.txt` | CLI stdout がある fixture で必須 | stdout 完全一致。stdout なしは空ファイル。 | 現在時刻、絶対環境 path。 |
| `expected/stderr.txt` | CLI stderr がある fixture で必須 | stderr 完全一致。stderr なしは空ファイル。 | secret、実 token、実 URL credential。 |
| `expected/state/` | 状態差分がある fixture で必須 | 実行後状態ファイル一式、または `state-diff.json`。 | 期待しないファイルの混入。 |
| `expected/logs/` | log / history / audit / notify fixture で必須 | build log、history、access log、audit log、notify log の期待値。 | secret 平文、実 Authorization header。 |
| `expected/effects.json` | 必須 | 外部 API 呼び出し、command 実行、通知送信、download/stream 中断、呼び出し 0 件の期待値。 | 呼び出し順未指定、実外部送信。 |
| `expected/security.json` | secret / auth / rate limit fixture で必須 | secret 非表示確認対象、禁止文字列、token hash 検証、scope 判定、rate count。 | secret を検証用に平文保存すること。 |

`expected/state/` は、fixture が検証対象とする状態ファイルだけを含める。変更してはならない状態ファイルは `expected/effects.json` の `unchanged_paths` に列挙する。削除されるべきファイルは `expected/effects.json` の `deleted_paths` に列挙し、空 directory の存在可否も明記する。

**§27 fixture manifest schema 固定契約：**

`manifest.json` は次の schema に従う。未知 key は禁止する。

```json
{
  "name": "success-example",
  "section": "27.1",
  "feature": "commit_status",
  "category": "success",
  "owner_component": "commitstatus",
  "collaborator_components": ["runner", "statefile"],
  "components": ["commitstatus", "runner", "statefile"],
  "references": ["§27.1", "§22.0a"],
  "fake_clock": "2026-09-16T00:00:00Z",
  "not_applicable": [
    { "path": "input/request.json", "reason": "CLI fixture" }
  ],
  "missing_state": [
    ".build_status.json"
  ],
  "assertions": [
    "state",
    "logs",
    "effects",
    "secret-mask"
  ]
}
```

| key | 型 | 必須 | 許容値 |
|-----|----|------|--------|
| `name` | string | 必須 | §27 fixture カタログ固定契約に記載された fixture 名。 |
| `section` | string | 必須 | `27.1`〜`27.38`、`27.42`〜`27.47`。 |
| `feature` | string | 必須 | lowercase snake_case。 |
| `category` | string | 必須 | `success`、`failure`、`partial`、`noop`、`security`。fixture 名 prefix と一致する。 |
| `owner_component` | string | 必須 | `builder`、`runner`、`api`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`security` のいずれか 1 件。 |
| `collaborator_components` | array[string] | 必須 | owner 以外の component。該当なしは空配列。 |
| `components` | array[string] | 必須 | `owner_component` と `collaborator_components` を重複なしで含む配列。許容値は `builder`、`runner`、`api`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`security`。 |
| `references` | array[string] | 必須 | 参照仕様節。対象 §27.x と関連 §22 / §23 / §24 / §25 / §26 を含める。 |
| `fake_clock` | string/null | 必須 | UTC ISO 8601 または `null`。時刻依存 fixture は `null` 禁止。 |
| `not_applicable` | array[object] | 必須 | 該当しない必須候補ファイルと理由。空配列可。 |
| `missing_state` | array[string] | 必須 | 実行前に存在しないことを期待する状態ファイル。空配列可。 |
| `assertions` | array[string] | 必須 | `response`、`stdout`、`stderr`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency`、`no-write` の 1 件以上。 |

**§27 fixture 合否判定固定契約：**

| 判定 | 合格条件 |
|------|----------|
| response | status、headers、body、error details、request id が期待値と一致する。 |
| stdout / stderr | 改行を含め完全一致する。時刻や path は fake 値だけを使う。 |
| state | 期待対象状態ファイルが byte 等価または `state-diff.json` と一致する。未列挙状態ファイルに差分がない。 |
| logs | JSON Lines は行順、key、値、末尾改行が一致する。破損行 fixture では破損行を修復しない。 |
| effects | 外部 API、外部 command、通知、download、stream の呼び出し回数、順序、payload が一致する。 |
| secret-mask | 禁止文字列が response、stdout/stderr、state、logs、effects、UI DOM に存在しない。 |
| order | 複数状態更新は仕様の保存順と一致する。途中失敗 fixture は失敗地点以降の副作用がない。 |
| idempotency | 同一 fixture を 2 回適用した場合、2 回目の差分が no-op 仕様と一致する。 |
| no-write | read-only / dry-run fixture で状態、logs、effects の差分がない。 |

上表のうち `manifest.json.assertions` に含まれる判定が 1 つでも失敗した場合、その fixture は失敗とする。対象機能の必須 fixture が 1 件でも存在しない、または skip された場合、その機能の実装 PR は未完了とする。

**§27 fixture assertion 選択固定契約：**

fixture の `manifest.json.assertions` は、実装者が任意に減らしてはならない。fixture 名 prefix、owner component、collaborator component に応じて、下表の assertion を必ず含める。個別節で追加検証が必要な場合は、下表へ追加してから fixture を作成する。

| 条件 | 必須 assertion |
|------|----------------|
| `success-*` | `state`、`logs`、`effects`。HTTP / SDK fixture では `response` も必須。CLI fixture では `stdout` または `stderr` の少なくとも一方を必須とする。 |
| `failure-*` | `response` または `stdout` / `stderr`、`state`、`effects`、`order`。状態差分なしを期待する場合は `no-write` も必須。 |
| `partial-*` | `state`、`logs`、`effects`、`order`。どの副作用まで完了し、どこから未実行かを `expected/effects.json` で明示する。 |
| `noop-*` | `response` または `stdout`、`no-write`、`idempotency`、`effects`。外部呼び出し 0 件を `expected/effects.json` に明記する。 |
| `security-*` | `response` または `stdout` / `stderr`、`secret-mask`、`effects`。認証 / scope / rate limit / TOTP / token fixture では `state` も必須。 |
| `components` に `api` を含む | `response`、`state`、`effects`。read-only API は `no-write` も必須。 |
| `components` に `sdk` を含む | `response`、`effects`。`401` fixture では token 破棄の期待値を `expected/state/` または `expected/effects.json` に含める。 |
| `components` に `ui` を含む | `response`、`effects`、`secret-mask`。DOM 期待値または UI action 後の field 消去期待を含める。 |
| `components` に `runner` を含む | `state`、`logs`、`effects`、`order`。dry-run は `no-write` を必須とする。 |
| `components` に `builder` を含む | `stdout` または `state`、`effects`。出力ファイル byte 比較または `state-diff.json` を必須とする。 |
| 外部 API / command / 通知を扱う | `effects`、`secret-mask`。呼び出し回数、順序、payload、mask 済み値を必須とする。 |

**§27 expected/effects.json schema 固定契約：**

`expected/effects.json` は次の schema に従う。未知 key は禁止する。

```json
{
  "external_calls": [
    {
      "type": "github_status",
      "order": 1,
      "method": "POST",
      "target": "/repos/{owner}/{repo}/statuses/{sha}",
      "payload": {},
      "result": "success"
    }
  ],
  "commands": [],
  "notifications": [],
  "downloads": [],
  "streams": [],
  "unchanged_paths": [],
  "deleted_paths": [],
  "created_paths": [],
  "write_order": [],
  "forbidden_writes": [],
  "forbidden_calls": []
}
```

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `external_calls` | array[object] | 必須 | GitHub API、SMTP、webhook、SSH、remote build の process 外呼び出し。呼び出しなしは空配列。 |
| `commands` | array[object] | 必須 | pipeline、hook、systemd、archive、setup、update の local command 実行。実行なしは空配列。 |
| `notifications` | array[object] | 必須 | 通知送信、pending 化、retry 対象。通知なしは空配列。 |
| `downloads` | array[object] | 必須 | snapshot / artifact download の byte size、content type、中断有無。該当なしは空配列。 |
| `streams` | array[object] | 必須 | SSE / fetch stream の event、close、error。該当なしは空配列。 |
| `unchanged_paths` | array[string] | 必須 | 実行後に変更があってはならない状態ファイル、出力ファイル、log。 |
| `deleted_paths` | array[string] | 必須 | 実行後に削除される path。削除なしは空配列。 |
| `created_paths` | array[string] | 必須 | 実行後に新規作成される path。作成なしは空配列。 |
| `write_order` | array[string] | 必須 | 書き込み順。書き込みなしは空配列。複数状態更新 fixture では空配列禁止。 |
| `forbidden_writes` | array[string] | 必須 | 書き込み禁止 path。read-only、dry-run、validation failure fixture では対象状態ファイルを必ず列挙する。 |
| `forbidden_calls` | array[string] | 必須 | 呼び出し禁止の外部 API / command / notification。呼び出し禁止なしは空配列。 |

`external_calls[]`、`commands[]`、`notifications[]` は `order` を持つ。並列処理 fixture で完了順が非決定の場合でも、期待保存順は `write_order` に固定する。secret を含む payload は、fixture 内でも平文を保存せず `"***"` を使う。

**§27 fixture 不足時 未完了判定固定契約：**

| 不足 | 未完了理由 |
|------|------------|
| fixture 名が §27 fixture カタログに存在しない。 | カタログ外 fixture のため未完了。 |
| 必須 fixture が存在しない。 | 機能の正常 / 異常 / no-op / security / partial coverage 不足。 |
| `manifest.json.assertions` が §27 fixture assertion 選択固定契約を満たさない。 | 合否判定不足。 |
| `expected/effects.json` が存在しない、または必須 key が欠落する。 | 副作用検証不足。 |
| `unchanged_paths` または `forbidden_writes` が空で、fixture が failure / noop / read-only / dry-run / validation error のいずれかである。 | 無変更保証不足。 |
| 外部 API / command / notification を扱う fixture で `forbidden_calls` が空、かつ禁止対象なしの理由が `manifest.json.not_applicable` にない。 | 外部副作用境界不足。 |
| secret fixture で `expected/security.json` が存在しない。 | secret mask 検証不足。 |
| JSON Lines を扱う fixture で末尾改行、行順、破損行保持/除外条件を期待値に含めていない。 | log / audit / history 検証不足。 |
| idempotency fixture で 1 回目と 2 回目の期待差分を分離していない。 | 再実行検証不足。 |
| partial fixture で失敗地点より後の `forbidden_writes` / `forbidden_calls` を列挙していない。 | 部分失敗境界不足。 |

上表に 1 件でも該当する場合、対象機能は実装済みとして扱わない。fixture が多くても、期待副作用、禁止副作用、secret mask、保存順、再実行差分が明示されていなければ、バグ修正ゼロ化の検証を満たさない。

**§27 fixture 相互整合固定契約：**

fixture 内の `manifest.json`、`input/*`、`expected/*` は相互に矛盾してはならない。下表の不整合が 1 件でもある場合、その fixture は失敗とし、対象機能の実装 PR は未完了とする。

| 整合対象 | 固定条件 | 不整合時の扱い |
|----------|----------|----------------|
| `manifest.json.name` と directory 名 | directory 名は `manifest.json.name` と完全一致する。 | fixture 名不一致として失敗。 |
| `manifest.json.category` と fixture 名 prefix | `success-*` は `success`、`failure-*` は `failure`、`partial-*` は `partial`、`noop-*` は `noop`、`security-*` は `security` とする。 | 分類不一致として失敗。 |
| `manifest.json.section` と fixture カタログ | section と fixture 名の組み合わせは §27 fixture カタログ固定契約に存在する組み合わせだけ許可する。 | カタログ外として失敗。 |
| `manifest.json.owner_component` と `components` | `owner_component` は `components` に必ず含める。 | owner 責務不一致として失敗。 |
| `manifest.json.collaborator_components` と `components` | `collaborator_components` は `components` にすべて含め、`owner_component` を含めてはならない。 | collaborator 責務不一致として失敗。 |
| `manifest.json.components` と入力ファイル | `runner` / `builder` を含む場合は `input/cli.json`、`api` / `sdk` / `ui` を含む場合は `input/request.json` を置く。該当しない場合は `not_applicable` に理由を置く。 | 入力責務不一致として失敗。 |
| `manifest.json.assertions` と期待値ファイル | `response` は `expected/response.json`、`stdout` は `expected/stdout.txt`、`stderr` は `expected/stderr.txt`、`state` は `expected/state/`、`logs` は `expected/logs/`、`effects` は `expected/effects.json`、`secret-mask` は `expected/security.json` を要求する。 | 期待値不足として失敗。 |
| `expected/effects.json.write_order` と `expected/state/` | `write_order` に列挙された path は `expected/state/` または `expected/logs/` に期待値を持つ。 | 保存順だけの空検証として失敗。 |
| `expected/effects.json.unchanged_paths` と `created_paths` / `deleted_paths` | 同一 path を `unchanged_paths` と `created_paths` または `deleted_paths` に同時に含めない。 | 副作用矛盾として失敗。 |
| `expected/effects.json.forbidden_writes` と `write_order` | 同一 path を `forbidden_writes` と `write_order` に同時に含めない。 | 禁止書込矛盾として失敗。 |
| `expected/effects.json.forbidden_calls` と `external_calls` / `commands` / `notifications` | 禁止対象と同じ call target を実行期待に含めない。 | 禁止呼び出し矛盾として失敗。 |
| `expected/security.json` と期待値全体 | `expected/security.json.forbidden_values` に含まれる文字列は、fixture directory 内の `input/` を除く全 expected file に出現してはならない。 | secret mask 不足として失敗。 |
| `input/fakes.json` と `expected/effects.json` | fake response を消費する外部 call / command は `expected/effects.json` に同数、同順で記録する。未消費 fake がある場合は `manifest.json.not_applicable` に理由を置く。 | fake / effect 不一致として失敗。 |
| `missing_state` と `input/state/` | `manifest.json.missing_state` に列挙した path は `input/state/` に存在してはならない。 | 初期状態矛盾として失敗。 |

**§27 component 別検証責務固定契約：**

実装 PR は、対象 component ごとに下表の責務を満たす。複数 component を含む機能では、owner component と collaborator component を fixture manifest に分けて記録し、各 component の責務をすべて満たすまで完了扱いにしてはならない。

| component | owner 時の必須検証責務 | collaborator 時の必須検証責務 | 完了判定 | 禁止越境 |
|-----------|------------------------|-------------------------------|----------|----------|
| `builder` | CLI 入力、Markdown 入力、出力 site / HTML / REPORT、asset、cache、dependency、meta、終了コードを fixture で固定する。 | runner / API に渡す REPORT、meta、dependency manifest の key 名と nullable 条件を固定する。 | 出力 file の byte 比較または構造化 expected が存在し、既存出力保護と失敗時 no-write が検証済み。 | GitHub read、状態ファイル直接更新、通知送信。 |
| `runner` | CLI、設定、lock、SHA cache、build log、history、status、queue、external call、notification、終了コードを fixture で固定する。 | builder output、statefile schema、commitstatus payload を仕様どおり消費し、未定義 key や未取得値を追加しない。 | 成功、失敗、skip、partial、dry-run の状態差分と write order が検証済み。 | API endpoint 追加、SDK method 追加、UI 操作追加。 |
| `api` | method/path/query/body/header、auth/scope/rate limit、response、状態 read/write、access/audit/config log を fixture で固定する。 | SDK / UI が追加 key を生成せずに扱える response schema、HTTP status、error body を返す。 | read-only は no-write、write API は保存順、validation failure は forbidden_writes が検証済み。 | runner CLI 処理、UI DOM 操作、未定義状態ファイル作成。 |
| `sdk` | method、args、query 生成、error 変換、token 破棄、binary / stream handling を fixture で固定する。 | UI が API 詳細を知らずに扱える戻り値と error をそのまま伝播する。 | API response に存在しない key を生成せず、`401` / `403` / `429` / network error が区別される。 | API response 推測補完、未定義 endpoint 呼び出し、状態ファイル直接操作。 |
| `ui` | SDK method 呼び出し、DOM 表示、disabled/loading/error、secret field 消去、再取得順を fixture で固定する。 | SDK 戻り値だけを表示し、API / 状態ファイルの内部構造を再解釈しない。 | 直接 API 呼び出し、状態ファイル操作、secret DOM 残存がない。 | 直接 `fetch()`、状態ファイル操作、外部 command 実行。 |
| `statefile` | schema、atomic write、JSON Lines、lock、破損時処理、保存順、no-write / forbidden write を fixture で固定する。 | runner / API / archive の保存対象ごとに `write_order`、`unchanged_paths`、`forbidden_writes` を固定する。 | read-only、dry-run、validation failure、partial failure の副作用境界が検証済み。 | component 固有の業務判断、UI 表示判断、API response 補完。 |
| `archive` | log archive、snapshot、download、delete、rollback の保存 / 取得 / 削除境界を fixture で固定する。 | API download / rollback response と runner history への影響を statefile 契約に合わせる。 | 元 log 保持、archive 破損時挙動、安全でない entry 拒否、rollback 失敗時 no-write が検証済み。 | build 成否の反転、元 build log の改変、未検証 archive 展開。 |
| `commitstatus` | GitHub Commit Status payload、pending / final 送信順、失敗時非反転、secret mask を fixture で固定する。 | runner の build 結果を受け取り、build 成否を変更せず送信結果だけを返す。 | pending 失敗、final 失敗、commit SHA なし、無効時呼び出し 0 が検証済み。 | build 実行判断、dry-run での GitHub write、status 失敗による build 成否反転。 |
| `security` | token hash、session、TOTP、scope、audit、rate limit、secret mask、forbidden call/write を fixture で固定する。 | API / SDK / UI / runner の secret 表示、認証失敗、副作用境界を検証する。 | 認証失敗、権限拒否、rate limit、audit failure の副作用境界が検証済み。 | 業務処理代行、認可前状態更新、secret 平文保存。 |
| `setup` | binary 配置、service 更新、rollback、secret 既存値保持、stdout/stderr mask、終了コードを fixture で固定する。 | runner / API の初期状態と既存 secret を壊さないことを effects で固定する。 | 部分失敗時の復元対象と復元禁止対象が `expected/effects.json` に明記済み。 | runtime 機能追加、状態 schema 暗黙変更、外部依存追加。 |

component 責務を別 PR へ分割する場合でも、分割先 PR が満たすべき owner component、collaborator component、fixture 名、期待ファイル、禁止副作用を PR 本文に明記する。責務の所在が不明な場合は、その機能を実装完了扱いにしてはならない。

**§27 PR 別必須記録固定契約：**

各 §27 実装 PR は、本文に下表を記録する。記録がない PR は、コードと fixture が存在しても未完了とする。

| 記録項目 | 必須内容 |
|----------|----------|
| wave | 対象 wave、対象 §27.x、先行 wave 完了 commit または PR 番号。 |
| 実装対象 | 実装する機能名、owner component、collaborator component、変更ファイル、追加 fixture path。 |
| 実装対象外 | 同じ wave 内で今回実装しない §27.x、後続 wave、MCP、外部公開構成、未定義 endpoint / UI / 状態ファイルを PR 証跡に列挙する。 |
| fixture | §27 fixture カタログの fixture 名、manifest / effects / security の検証結果。 |
| acceptance | §27 実装 PR acceptance checklist の各項目の pass / fail / 未実行。 |
| 後続影響 | 後続 PR が利用許可済みの contract、利用禁止の未固定 contract を PR 証跡に列挙する。 |

**§27 実装 PR 最終受け入れゲート：**

§27 の機能実装 PR は、下表の全 gate を満たした場合だけ「完了」と判定する。1 件でも未達がある場合は「未完了」、仕様逸脱または secret 漏えいリスクがある場合は「差し戻し」とする。

| gate | 完了条件 | 未完了条件 | 差し戻し条件 |
|------|----------|------------|--------------|
| scope | 実装対象が §27 fixture カタログ固定契約に存在する機能だけである。 | 実装対象節の記載が PR 本文にない。 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加している。 |
| fixture | 対象 §27.x の必須 fixture がすべて存在し、skip されていない。 | 必須 fixture が不足、または fixture 名が不一致。 | fixture が実装挙動に合わせて期待値を緩めている。 |
| manifest | 全 fixture の `manifest.json` が schema、assertion 選択、owner / collaborator component 責務を満たす。 | assertion、references、owner_component、collaborator_components、components、fake_clock のいずれかが不足。 | unknown key、実 secret、実環境 path、乱数依存を含む。 |
| expected | `expected/response.json`、`expected/state/`、`expected/logs/`、`expected/effects.json`、`expected/security.json` が assertion と一致する。 | assertion に対応する expected file が不足。 | expected と manifest / effects / security が矛盾する。 |
| side effect | `write_order`、`unchanged_paths`、`forbidden_writes`、`forbidden_calls` が対象機能の成功 / 失敗 / no-op / partial を説明できる。 | 禁止副作用または無変更保証が不足。 | 失敗時に未許可状態を書き換える、外部呼び出しを行う。 |
| secret | secret 平文が expected、logs、effects、UI DOM、stdout/stderr に存在しない。 | secret 検証対象が不足。 | token、password、TOTP secret、PAT、Authorization header が平文で残る。 |
| component | builder / runner / api / sdk / ui / statefile / archive / commitstatus / security / setup の該当責務が全て fixture に紐づく。 | owner / collaborator component の所在が不明。 | SDK / UI が API response を推測補完、または UI が直接 API / 状態ファイルを操作する。 |
| repeatability | fake clock、fake external response、固定 path により、同じ fixture が同じ結果を再現する。 | idempotency / no-op の 2 回目期待値が不足。 | 現在時刻、実ネットワーク、実 OS 差分に依存する。 |

**§27 実装 PR acceptance checklist：**

実装 PR 本文または検証ログには、下表の項目を記録する。記録がない項目は未検証として扱い、対象機能を完了扱いにしてはならない。

| 項目 | 記録内容 |
|------|----------|
| 対象仕様 | 実装した §27.x、関連 §22 / §23 / §24 / §25 / §26、owner component、collaborator component。 |
| 対象 fixture | 作成または更新した fixture 名一覧。fixture 名は §27 fixture カタログ固定契約と一致させる。 |
| 実行結果 | fixture ごとの pass / fail、実行コマンド、終了コード。 |
| 状態差分 | 作成、更新、削除、変更禁止の path。`expected/effects.json` と一致させる。 |
| 外部副作用 | GitHub、SSH、SMTP、webhook、systemd、hook、remote build、notification の呼び出し回数と順序。 |
| secret 検証 | 禁止文字列、mask 対象、平文が残らないことを確認した出力範囲。 |
| 部分失敗 | partial / failure fixture の失敗地点、完了済み副作用、禁止副作用。 |
| 再実行 | idempotency / no-op fixture の 1 回目と 2 回目の差分。 |
| 対象外確認 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加していないことを、差分対象 file と fixture manifest の両方で確認する。 |

**§27 差し戻し固定条件：**

次のいずれかに該当する PR は、fixture が pass していても差し戻しとする。

| 条件 | 理由 |
|------|------|
| 仕様にない endpoint、SDK method、UI 操作、状態ファイルを追加している。 | 仕様外実装。 |
| fixture の期待値が実装都合に合わせて仕様より弱い。 | 検証の形骸化。 |
| secret 平文、Authorization header、TOTP secret、PAT、password が expected または log に残る。 | secret 漏えい。 |
| read-only / dry-run / validation failure で状態、log、外部 call が変化する。 | 副作用違反。 |
| partial failure で失敗地点以降の write / call が発生する。 | 部分失敗境界違反。 |
| SDK が API response を補完し、UI が SDK を迂回し、runner / builder が未定義状態ファイルを作成する。 | component 責務違反。 |
| 実ネットワーク、実時刻、実ユーザー環境、実 secret に依存する fixture だけで合格している。 | 再現性不足。 |

**§27 部分失敗・再実行固定契約：**

| ケース | 固定挙動 |
|--------|----------|
| 状態保存前の validation 失敗 | 状態ファイル、JSON Lines、外部 API、外部 command、通知を実行しない。 |
| 主状態保存後の log 追記失敗 | 個別節が rollback を明記しない限り主状態は戻さず、response は `500` とし、次回 GET は保存済み主状態を返す。 |
| log 追記後の audit 失敗 | 個別節で必須 audit と定義された操作は `500`。任意 audit の場合は対象操作成功を維持し、WARN を保存する。 |
| 外部 API 送信成功後の状態保存失敗 | 外部送信の再送を自動実行しない。状態保存失敗を `500` または runner failure として記録し、再実行時は個別節の重複防止 key で判定する。 |
| 通知送信失敗 | build / config / security の主結果を反転しない。`.notify_pending` または個別節の失敗記録だけを更新する。 |
| download / stream 中断 | サーバー側状態を成功/失敗へ変更しない。access log は endpoint 固有節で中断 status 記録が定義されている場合だけ追記し、history と build log は変更しない。 |
| 再実行 no-op | 同一入力で差分がない保存 API は、個別節の no-op response を返し、状態、config log、audit log、notify log に新規差分を作らない。 |
| 破損 JSON Lines | read API は破損行を除外し、破損内容を response に出さない。write API は既存破損行を修復、削除、並べ替えしない。 |

## 28-F §28 builder 拡張 fixture / PR 証跡詳細契約

本節は、`docs/details/builder.md` §28.1〜§28.25 の fixture、fake、expected、effects、PR 証跡の正本である。各 §28 機能は、Markdown 入力、CLI option、期待 HTML、期待 CSS / JS、`[REPORT]`、終了コード、strict / non-strict の差分を fixture で固定する。外部 library、CDN、実 network、現在時刻、実 git repository、実画像取得、画像 snapshot だけの合否判定を fixture の前提にしてはならない。

**§28 fixture 配置固定契約：**

| 節 | fixture 配置単位 | 必須 fixture |
|----|------------------|--------------|
| §28 共通 | `builder-extensions/config-resolution/` | 下記 §28 fixture カタログ固定契約の `config-resolution` fixture をすべて作成する。 |
| §28 共通 | `builder-extensions/determinism/` | 下記 §28 fixture カタログ固定契約の `determinism` fixture をすべて作成する。 |
| §28 共通 | `builder-extensions/atomicity/` | 下記 §28 fixture カタログ固定契約の `atomicity` fixture をすべて作成する。 |
| §28 共通 | `builder-extensions/parser-precedence/` | 下記 §28 fixture カタログ固定契約の `parser-precedence` fixture をすべて作成する。 |
| §28 共通 | `builder-extensions/browser-runtime/` | 下記 §28 fixture カタログ固定契約の `browser-runtime` fixture をすべて作成する。 |
| §28 共通 | `builder-extensions/visual-layout/` | 下記 §28 fixture カタログ固定契約の `visual-layout` fixture をすべて作成する。 |
| §28.1〜§28.25 | `builder-extensions/<feature-slug>/` | 下記 §28 fixture カタログ固定契約に列挙した fixture をすべて作成する。 |

**§28 fixture カタログ固定契約：**

下表の fixture 名は固定値である。実装 PR では、対象 §28.x の全 fixture を追加または更新し、`expected/stdout.txt` の `[REPORT]` key が `docs/details/builder.md` §28 の固定契約と一致することを証跡に含める。

| 節 | feature slug | 必須 fixture |
|----|--------------|--------------|
| §28 共通 | `config-resolution` | `success-config-defaults-only`、`success-config-file-values`、`success-config-env-over-file`、`success-config-cli-over-env-over-file`、`failure-config-json-corrupt`、`failure-config-unknown-key`、`failure-config-invalid-type`、`failure-config-duplicate-nonrepeatable-cli`、`security-config-secret-not-echoed` |
| §28 共通 | `determinism` | `success-slug-duplicates`、`success-search-index-text-sources`、`success-local-storage-payload`、`success-hash-targets`、`security-deterministic-no-runtime-variance` |
| §28 共通 | `atomicity` | `success-atomic-write-all-files`、`success-incremental-reuse-byte-identical`、`success-incremental-delete-stale-page`、`failure-strict-warning-no-replace`、`failure-write-error-no-partial-update`、`failure-changed-manifest-invalid-no-output`、`success-dependency-manifest-corrupt-full-build`、`security-atomic-no-stale-temp-promoted` |
| §28 共通 | `parser-precedence` | `success-block-precedence-code-math-heading`、`success-inline-precedence-code-image-link`、`success-admonition-inline-composition`、`success-heading-inline-slug-source`、`success-list-definition-task-boundary`、`failure-unclosed-math-strict`、`noop-code-fence-protects-extensions`、`security-parser-raw-html-escaped` |
| §28 共通 | `browser-runtime` | `success-runtime-init-order`、`success-section-collapse-storage-print`、`success-color-scheme-cycle-print`、`success-toc-active-observer-fallback`、`success-hash-history-focus-navigation`、`success-lightbox-focus-trap-close`、`success-keyboard-scope-skip-link`、`security-runtime-no-storage-leak` |
| §28 共通 | `visual-layout` | `success-css-output-order`、`success-responsive-320-layout`、`success-print-layout`、`success-color-scheme-variables`、`success-minify-visual-preservation`、`success-component-overflow-boundaries`、`security-visual-no-external-assets`、`security-focus-visible-no-overlap` |
| §28.1 | `incremental` | `success-one-page-change`、`success-dependency-change`、`success-reuse-page-copied-to-staging`、`success-stale-page-delete-on-success`、`success-incremental-dependency-manifest-corrupt-full-build`、`failure-changed-manifest-corrupt`、`failure-changed-manifest-base-escape`、`noop-unchanged-pages-kept`、`security-incremental-no-public-write-on-failure` |
| §28.2 | `formats` | `success-html-default`、`success-html-explicit`、`failure-pdf-reserved`、`failure-epub-reserved`、`failure-unknown-format`、`failure-multiple-format`、`security-format-no-reserved-output-files` |
| §28.3 | `markdown-extensions` | `success-admonition-note-warn-tip`、`success-badge-color`、`success-extension-csv-normalization`、`failure-unknown-extension`、`failure-badge-invalid-text-strict`、`noop-badge-invalid-text-nonstrict`、`security-extension-escape`、`noop-extension-disabled` |
| §28.4 | `code-line-numbers` | `success-line-numbers-fence`、`success-line-numbers-cli`、`success-line-numbers-diff-composition`、`noop-line-numbers-empty-code`、`noop-line-numbers-disabled`、`security-line-numbers-copy-clean` |
| §28.5 | `heading-numbering` | `success-heading-numbering-h2-h3`、`success-heading-numbering-implicit-h2`、`success-heading-numbering-toc-search`、`failure-heading-numbering-unknown-mode`、`noop-heading-numbering-none`、`security-heading-slug-unchanged` |
| §28.6 | `section-collapse` | `success-collapse-h2-h3`、`success-collapse-local-storage`、`noop-collapse-no-heading`、`security-collapse-duplicate-target-strict` |
| §28.7 | `toc-depth` | `success-toc-depth-h2-h3`、`success-toc-depth-h1-h6`、`failure-toc-depth-invalid-range`、`security-toc-depth-active-sync` |
| §28.8 | `updated-at` | `success-updated-at-git`、`success-updated-at-file`、`success-updated-at-fallback`、`failure-updated-at-unknown-source` |
| §28.9 | `diff-highlight` | `success-diff-insert-delete-context`、`success-diff-header`、`noop-diff-non-diff-language`、`security-diff-escape` |
| §28.10 | `lazy-images` | `success-lazy-relative-image`、`success-lazy-external-image-no-fetch`、`failure-lazy-base-outside-strict`、`security-lazy-alt-escape` |
| §28.11 | `custom-meta` | `success-meta-og-twitter`、`success-meta-duplicate-last-wins`、`failure-meta-forbidden-key`、`security-meta-escape` |
| §28.12 | `color-scheme` | `success-color-scheme-light`、`success-color-scheme-dark`、`success-color-scheme-auto`、`failure-color-scheme-unknown` |
| §28.13 | `code-title` | `success-code-title-colon`、`success-code-title-key-value`、`noop-code-title-empty`、`security-code-title-escape` |
| §28.14 | `template-vars` | `success-template-var-replace`、`noop-template-var-code-fence`、`failure-template-var-missing-strict`、`security-template-var-key-validation` |
| §28.15 | `minify-html` | `success-minify-html`、`success-minify-preserve-code`、`failure-minify-marker-missing`、`noop-minify-disabled` |
| §28.16 | `toc-active` | `success-toc-active-scroll`、`success-toc-active-fallback`、`noop-toc-active-disabled`、`security-toc-active-depth-sync` |
| §28.17 | `mermaid` | `success-mermaid-graph-td`、`failure-mermaid-unsupported-strict`、`noop-mermaid-disabled`、`security-mermaid-no-external-script` |
| §28.18 | `footnotes` | `success-footnotes-multiple`、`success-footnotes-backlink`、`failure-footnote-undefined-strict`、`security-footnote-escape` |
| §28.19 | `math` | `success-math-inline-block`、`failure-math-unclosed-strict`、`noop-math-code-fence`、`security-math-escape` |
| §28.20 | `hash-history` | `success-hash-history-click`、`success-hash-history-back-forward`、`noop-hash-history-disabled`、`security-hash-history-missing-target` |
| §28.21 | `a11y` | `success-a11y-landmarks-labels`、`success-a11y-skip-link-tab-order`、`failure-a11y-duplicate-id-strict`、`security-a11y-no-keyboard-trap` |
| §28.22 | `image-lightbox` | `success-lightbox-open-close`、`success-lightbox-escape-backdrop`、`failure-lightbox-alt-missing-strict`、`security-lightbox-focus-trap` |
| §28.23 | `print-qr` | `success-print-qr-url`、`noop-print-qr-empty-url`、`failure-print-qr-url-too-long`、`security-print-qr-svg-escape` |
| §28.24 | `definition-lists` | `success-definition-list-single`、`success-definition-list-multiple`、`noop-definition-list-empty-term`、`security-definition-list-inline-escape` |
| §28.25 | `task-lists` | `success-task-list-unchecked`、`success-task-list-checked-nested`、`noop-task-list-non-target`、`security-task-list-disabled-aria` |

**§28 fixture ファイルセット固定契約：**

| ファイル | 必須 | 内容 |
|----------|------|------|
| `manifest.json` | 必須 | fixture 名、対象 §28 共通または §28.x、feature slug、owner `builder`、collaborator、参照仕様節、strict 有無、fake clock、not_applicable 理由。 |
| `input/source.md` または `input/site/` | 必須 | Markdown 入力。site fixture は複数 Markdown、asset、dependency を含める。 |
| `input/options.json` | 必須 | CLI option、env key、expected exit code、strict / non-strict。 |
| `input/adlaire-ci-build.json` | 条件付き | 設定ファイル fixture で使用する。未使用 fixture では存在させない。 |
| `input/fakes.json` | 条件付き | fake git timestamp、fake file mtime、fake manifest、fake cache、fake clipboard など。実外部呼び出しは禁止。 |
| `input/existing-site/` | 条件付き | atomicity、incremental、failure fixture で既存公開出力を表す。成功 fixture では置換前状態、failure fixture では維持されるべき状態を置く。 |
| `expected/site/` | 必須 | 期待 HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` のうち対象機能が変更する file。 |
| `expected/stdout.txt` | 必須 | 進捗、`[WARN]`、`[REPORT]` を含む stdout 完全一致。fatal failure は空 file。 |
| `expected/stderr.txt` | 必須 | fatal failure の `[ERROR]` 完全一致。stderr なしは空 file。 |
| `expected/effects.json` | 必須 | 作成、更新、維持、削除禁止 path、外部 call 0 件、既存出力保護、strict 昇格条件。 |
| `expected/security.json` | 条件付き | HTML escape、attribute escape、外部 library 不使用、secret / URL credential 非表示、base 外 path 拒否。 |

**§28 fixture 判定粒度固定契約：**

§28 fixture は、目視確認、画像 snapshot、現在時刻、実 git repository、実 network、ブラウザ環境だけに依存して合否判定してはならない。期待値は file 内容、stdout、stderr、終了コード、副作用、禁止出力のいずれかで固定する。

| 判定対象 | 固定内容 |
|----------|----------|
| HTML | 対象機能が生成する tag、attribute、class、data attribute、aria attribute、escape 済み text を完全一致または正規化済み比較で確認する。 |
| CSS | 対象機能が追加する selector、custom property、`@media print`、focus style を確認する。未使用機能の selector が出ないことも確認する。 |
| JS | 対象機能が追加する event handler、localStorage key、history handler、dialog handler、fallback 分岐の文字列または構造を確認する。外部 script 参照がないことを確認する。 |
| search index | 採番表示、line number 除外、HTML tag 除外、対象 page path、updated time の有無、UI text 除外を確認する。 |
| stdout | 既存進捗行、`[WARN]`、`[REPORT]` の有無、`[REPORT]` key 順、値型、既定値、件数、warning count を完全一致で確認する。fatal failure は空 file にする。 |
| stderr | fatal failure の `[ERROR]` code、対象 file、対象 section、strict 昇格有無を完全一致で確認する。warning 継続 case と strict warning case は空 file にする。 |
| effects | 作成、更新、維持、削除禁止、既存出力維持、manifest 上書き有無、外部 call 0 件を JSON で確認する。 |
| security | HTML escape、attribute escape、base 外 path、URL credential 非表示、secret 非表示、CDN / external library 不使用を確認する。 |
| visual layout | CSS selector 順、custom property、media query、responsive overflow、print visibility、minify preserve、focus visibility、layout shift 禁止を確認する。 |
| parser precedence | block token 優先順位、inline token 優先順位、code fence / code span 保護、曖昧構文、機能併用順を確認する。 |
| browser runtime | JS 初期化順、event handler、focus、keyboard、localStorage、print、fallback、例外時 no-break を確認する。 |

`expected/effects.json` は、§28 fixture では以下の key を固定する。未使用 key も省略せず、空配列、空 object、または `false` で明示する。

| key | 型 | 固定 |
|-----|----|------|
| `created_paths` | array | 新規作成される公開出力 path。ASCII 昇順。 |
| `updated_paths` | array | 既存から内容が変わる公開出力 path。ASCII 昇順。 |
| `preserved_paths` | array | failure または reuse により byte 単位で維持される path。ASCII 昇順。 |
| `deleted_paths` | array | 成功時に削除される stale path。ASCII 昇順。failure fixture では空配列。 |
| `forbidden_created_paths` | array | 作成してはならない path。tmp、未定義 asset、reserved format 出力を含める。 |
| `forbidden_updated_paths` | array | 更新してはならない path。failure fixture では既存 HTML、manifest、search index、asset を含める。 |
| `forbidden_deleted_paths` | array | 削除してはならない path。failure fixture では既存公開出力を含める。 |
| `external_calls` | integer | 常に `0`。 |
| `staging_cleaned` | boolean | staging directory が残らない場合 `true`。staging cleanup 失敗 fixture では `false` を許可し、その場合は終了コード `1` を期待する。 |
| `public_output_replaced` | boolean | 成功 transaction で公開 `--out` が置換される場合だけ `true`。 |
| `manifest_written` | boolean | 成功 transaction で `.dependency_manifest.json` が公開される場合だけ `true`。 |
| `search_index_regenerated` | boolean | 成功 transaction で最終 page set から search index を再生成する場合だけ `true`。 |

**§28 設定解決 fixture 固定契約：**

`builder-extensions/config-resolution/` は §28 全機能の前提 fixture である。§28.1〜§28.25 の個別 fixture は、この共通 fixture と矛盾する CLI / env / config / default 解決をしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-config-defaults-only` | CLI、環境変数、`adlaire-ci-build.json` がない場合、§28 の全設定 key が既定値から採用され、`config_default_keys` に全 key、その他 source key に空 array が出る。 |
| `success-config-file-values` | `input/adlaire-ci-build.json` の `builder_extensions` 値が採用され、`config_file_used=true`、`config_file_path="adlaire-ci-build.json"`、`config_file_keys` が ASCII 昇順で出る。 |
| `success-config-env-over-file` | 同一 key が環境変数と設定ファイルに存在する場合、環境変数を採用し、設定ファイル source は `config_overridden_keys` に記録する。 |
| `success-config-cli-over-env-over-file` | 同一 key が CLI、環境変数、設定ファイルに存在する場合、CLI を採用し、環境変数 / 設定ファイル source は `config_overridden_keys` に記録する。 |
| `failure-config-json-corrupt` | 設定ファイルが JSON として parse できない場合、終了コード `2`、stdout 空、stderr `BUILDER28_INVALID_OPTION`、出力作成なし、既存出力維持。 |
| `failure-config-unknown-key` | root unknown key、`builder_extensions` unknown key のいずれも終了コード `2`、stdout 空、stderr に key 名だけを記録し、`[REPORT]` は出力しない。 |
| `failure-config-invalid-type` | boolean / string / array / object の型不一致、JSON object 値が string 以外、空 array 要素を終了コード `2`、stdout 空にする。 |
| `failure-config-duplicate-nonrepeatable-cli` | repeatable ではない CLI option の重複指定を終了コード `2`、stdout 空にし、Markdown 読込前に停止する。 |
| `security-config-secret-not-echoed` | `meta` / `template_vars` / 環境変数値に secret 風文字列、credential 付き URL、raw HTML が含まれても、stderr、stdout、REPORT へ値を出さない。key 名だけを出す。 |

設定解決成功 fixture の `expected/stdout.txt` は、`config_file_used`、`config_file_path`、`config_cli_keys`、`config_env_keys`、`config_file_keys`、`config_default_keys`、`config_overridden_keys`、`config_rejected_keys` を必ず含める。設定解決失敗 fixture の `expected/stdout.txt` は空 file とし、`expected/stderr.txt` に `[ERROR] BUILDER28_INVALID_OPTION -:0 28 ...` を固定する。設定解決失敗 fixture の `expected/effects.json` は、HTML、CSS、JS、search index、manifest、設定ファイルが新規作成、更新、削除されないことを固定する。

**§28 決定性 fixture 固定契約：**

`builder-extensions/determinism/` は、§28 の HTML identity、slug、search index、hash、localStorage の共通 fixture である。個別 §28 fixture は、本 fixture と異なる slug、id、search text、storage key、hash target を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-slug-duplicates` | ASCII、非 ASCII、記号、空 heading、同名 heading、`foo` と `foo-2` の衝突を含む入力で、heading id、TOC href、section wrapper id、hash target が `docs/details/builder.md` §28 の slug 規則と一致する。 |
| `success-search-index-text-sources` | 採番 heading、paragraph、list、code、admonition、badge、image alt、footnote、definition list、task list を含む入力で、search index に含める text と除外する UI text が完全一致する。 |
| `success-local-storage-payload` | section collapse と color scheme を有効にし、`adlaire:section-state`、`adlaire:color-scheme` の key、payload、未知値無視、JSON parse failure fallback を `expected/site/assets/app.js` と `expected/effects.json` で固定する。 |
| `success-hash-targets` | hash history と TOC active tracking を有効にし、heading id だけを target にすること、TOC depth 外 heading を active 化しないこと、存在しない hash を no-op にすることを固定する。 |
| `security-deterministic-no-runtime-variance` | 同一入力を fake clock、fake git、異なる OS path separator 相当入力、異なる map order 相当 config で実行しても、HTML、search index、stdout、stderr が同一になることを固定する。 |

決定性 fixture の `expected/site/*.html` は、heading id、TOC href、collapse wrapper id、`aria-controls`、`data-section-id` を完全一致で確認する。`expected/site/assets/search-index.json` は page key、heading text、body text、除外 text の不在を JSON parse 後完全一致で確認する。`expected/site/assets/app.js` は localStorage key、payload schema、unknown value guard、parse failure guard、hash no-op guard を文字列または構造で確認する。

**§28 atomicity fixture 固定契約：**

`builder-extensions/atomicity/` は、§28 の出力副作用、公開置換、manifest、search index、stale 削除を固定する共通 fixture である。個別 §28 fixture は、本 fixture と異なる失敗時副作用、manifest 更新条件、search index 更新条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-atomic-write-all-files` | HTML、CSS、JS、search index、`.dependency_manifest.json` が staging にそろってから公開 `--out` へ置換され、`expected/effects.json` の `public_output_replaced=true`、`manifest_written=true`、`search_index_regenerated=true` になる。 |
| `success-incremental-reuse-byte-identical` | 未変更 page の既存 HTML が byte 単位で維持され、changed page、manifest、search index だけが成功 transaction として更新される。 |
| `success-incremental-delete-stale-page` | 入力 source から削除された Markdown に対応する HTML、search index entry、manifest entry が成功時だけ削除される。 |
| `failure-strict-warning-no-replace` | non-strict なら fallback 出力できる警告を strict で実行し、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力、manifest、search index 維持を固定する。 |
| `failure-write-error-no-partial-update` | staging 書き込みまたは validation 失敗を fake し、終了コード `1`、stderr `[ERROR] BUILDER28_INTERNAL_IO` または `BUILDER28_OUTPUT_VALIDATION_FAILED`、公開出力、manifest、search index 維持を固定する。 |
| `failure-changed-manifest-invalid-no-output` | `--changed-manifest` が base 外 path、絶対 path、URL scheme、JSON 破損のいずれかの場合、終了コード `2`、stdout 空、公開出力維持を固定する。 |
| `success-dependency-manifest-corrupt-full-build` | 既存 `.dependency_manifest.json` が破損または schema 不一致の場合、warning なし full build とし、成功時だけ新 manifest と search index を公開する。 |
| `security-atomic-no-stale-temp-promoted` | staging path、absolute path、host user path、tmp path が HTML、CSS、JS、search index、manifest、stdout、stderr、REPORT に混入しないことを固定する。 |

atomicity fixture の `input/existing-site/` は、既存 HTML、既存 `assets/search-index.json`、既存 `.dependency_manifest.json`、stale HTML、既存 asset を含める。failure fixture の `expected/site/` は `input/existing-site/` と byte 単位で一致させる。success fixture の `expected/effects.json` は、`created_paths`、`updated_paths`、`preserved_paths`、`deleted_paths` をすべて明示する。

**§28 parser precedence fixture 固定契約：**

`builder-extensions/parser-precedence/` は、§28 の Markdown parser 優先順位、構文 grammar、曖昧構文、機能併用順を固定する共通 fixture である。個別 §28 fixture は、本 fixture と異なる token 解釈、別順序の inline 変換、code fence / code span 内変換を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-block-precedence-code-math-heading` | code fence 継続中の heading / footnote / badge / math が code text のまま残り、math block と heading の判定順が `docs/details/builder.md` §28 Markdown parser 優先順位固定契約と一致する。 |
| `success-inline-precedence-code-image-link` | code span 内の badge / footnote / math が変換されず、image が link より優先され、link text 内の badge / math だけが inline 変換される。 |
| `success-admonition-inline-composition` | admonition body 内の badge、footnote、math、fenced code の併用で、body inline 変換と fenced code 保護が両立する。 |
| `success-heading-inline-slug-source` | heading 内の badge、footnote、math 表示変換と、slug source text から UI text を除外する規則が同時に成立する。 |
| `success-list-definition-task-boundary` | task list、通常 list、definition list、list 内 `: definition` の境界が固定どおりに分かれる。 |
| `failure-unclosed-math-strict` | 未閉鎖 math inline / math block が non-strict では通常 text、strict では終了コード `2` と `BUILDER28_UNRESOLVED_REFERENCE` になる。 |
| `noop-code-fence-protects-extensions` | code fence 内の template var、badge、footnote、math、definition marker、task marker が一切変換されない。 |
| `security-parser-raw-html-escaped` | raw HTML、event handler、`javascript:` URL、HTML comment 指示が parser 段階で実行可能要素にならず、expected HTML と security.json で escape を確認する。 |

parser precedence fixture の `expected/site/*.html` は、対象 token の tag、text node、未変換 text、変換済み node、属性順を完全一致で確認する。`expected/stdout.txt` は warning の有無、warning code、line、section を完全一致で確認する。`expected/security.json` は raw HTML、script、event handler、credential URL、CDN、外部 library が出力に存在しないことを固定する。

**§28 browser runtime fixture 固定契約：**

`builder-extensions/browser-runtime/` は、§28 のブラウザ JS 初期化順、状態復元、event handler、focus、keyboard、print、fallback を固定する共通 fixture である。個別 §28 fixture は、本 fixture と異なる localStorage key、focus 移動、keyboard scope、lightbox close 条件、TOC active 条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-runtime-init-order` | `assets/app.js` 内で static guard、storage guard、color scheme、section collapse、hash history、TOC active、lightbox、accessibility guard の初期化順が固定どおりである。 |
| `success-section-collapse-storage-print` | section collapse の既定展開、保存値復元、toggle、`aria-expanded`、`adlaire-section-collapsed`、search hit 一時展開、beforeprint / afterprint 復元が一致する。 |
| `success-color-scheme-cycle-print` | `light → dark → auto → light` の toggle、root `data-color-scheme`、toggle `aria-label`、`adlaire:color-scheme`、print light が一致する。 |
| `success-toc-active-observer-fallback` | IntersectionObserver 使用時と fallback scroll 時の active link 1 件化、`.is-active`、`aria-current="location"`、TOC depth 外除外が一致する。 |
| `success-hash-history-focus-navigation` | heading / TOC click、`history.pushState`、`tabindex="-1"`、focus、back / forward、missing hash no-op が一致する。 |
| `success-lightbox-focus-trap-close` | trigger click、`Enter` / `Space`、dialog open、Escape、backdrop、close button、opener focus return、Tab / Shift+Tab focus trap が一致する。 |
| `success-keyboard-scope-skip-link` | §28 keyboard handler が対象 UI focus 中だけ有効で、既存 §7.12 shortcut を上書きせず、skip link が main content へ移動する。 |
| `security-runtime-no-storage-leak` | cookie、sessionStorage、IndexedDB、runtime network fetch、external script、secret / credential の storage 書込が 0 件である。 |

browser runtime fixture の `expected/site/assets/app.js` は、初期化関数名または固定 marker、localStorage key、event 名、guard、fallback 分岐、focus trap 分岐を文字列または構造で確認する。`expected/security.json` は、cookie、sessionStorage、IndexedDB、fetch、XMLHttpRequest、external script、secret / credential storage が存在しないことを固定する。ブラウザ実行がない fixture でも、期待 JS 構造と expected HTML / CSS / security を組み合わせて合否判定する。

**§28 visual layout fixture 固定契約：**

`builder-extensions/visual-layout/` は、§28 の CSS 出力順、responsive layout、print layout、color scheme 変数、minify 後の視覚維持、component overflow、外部 visual asset 禁止、focus 表示を固定する共通 fixture である。個別 §28 fixture は、本 fixture と異なる selector 順、media query 境界、print visibility、外部 asset 許可、focus 表示条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-css-output-order` | `expected/site/assets/style.css` で、既存 base、color scheme、typography / block、code extension、navigation runtime UI、media UI、responsive、print の順序が固定どおりである。 |
| `success-responsive-320-layout` | 幅 `320px` 相当の fixture metadata と expected CSS / HTML で、§28 UI が text overlap、text clipping、不可視 overflow を発生させず、table と code block だけが scroll wrapper 内で横 overflow を持つ。 |
| `success-print-layout` | `@media print` で interactive controls を非表示、本文要素を表示、collapsed section を展開、color scheme を light 相当、QR を print 専用表示にする。 |
| `success-color-scheme-variables` | `:root`、`[data-color-scheme="light"]`、`[data-color-scheme="dark"]`、`[data-color-scheme="auto"]`、print の custom property 名と既定値が固定どおりである。 |
| `success-minify-visual-preservation` | minify 有効時も required selector、custom property、`@media (max-width: 768px)`、`@media print`、`pre` / `code` の空白保持 property が削除、改名、結合破壊されない。 |
| `success-component-overflow-boundaries` | admonition、badge、definition list、task list、footnote、math、code title、line numbers、diff、lightbox、Mermaid、print QR が `docs/details/builder.md` §28 CSS / layout / print / visual 固定契約の境界どおりに出力される。 |
| `security-visual-no-external-assets` | `@import`、remote `url()`、external font、CDN、追加 asset file、runtime CSS fetch、inline style が出力されない。 |
| `security-focus-visible-no-overlap` | theme toggle、skip link、TOC active、hash target、lightbox control、collapse toggle の focus / active style が可視で、hover / focus により寸法が変わらず、text を隠さない。 |

visual layout fixture の `manifest.json` は、`viewport_width` を使う場合でも判定を画像 snapshot だけに依存させてはならない。`expected/site/assets/style.css`、`expected/site/*.html`、`expected/security.json`、必要に応じて `expected/visual.json` に、selector、media query、display、overflow、visibility、focus、禁止 asset を構造化して固定する。ブラウザ実行がない fixture でも、期待 HTML / CSS / security の組み合わせで合否判定できなければならない。

**§28.1〜§28.5 feature fixture 固定契約：**

§28.1〜§28.5 の fixture は、`docs/details/builder.md` §28.1〜§28.5 実装詳細固定契約に列挙された中間状態、HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。各 fixture は `manifest.json.section` を対象 §28.x に固定し、`manifest.json.feature_slug` を §28 fixture カタログ固定契約の feature slug と一致させる。

| feature slug | fixture | 固定する内容 |
|--------------|---------|--------------|
| `incremental` | `success-one-page-change` | 複数 page のうち 1 page だけが changed となり、changed page は再生成、reused page は staging へ byte copy、search index と manifest は全 page から再生成される。 |
| `incremental` | `success-dependency-change` | `dependencies_changed` に含まれる asset を参照する page だけが changed になり、参照しない page は reused になる。 |
| `incremental` | `success-reuse-page-copied-to-staging` | reused page が公開出力から staging へ複写され、公開出力を直接参照したまま成功扱いにしないことを `expected/effects.json` で確認する。 |
| `incremental` | `success-stale-page-delete-on-success` | 入力 source から消えた page の HTML、manifest entry、search index entry が成功置換時だけ削除される。 |
| `incremental` | `success-incremental-dependency-manifest-corrupt-full-build` | 既存 `.dependency_manifest.json` の JSON 破損または schema 不一致を warning なし full build とし、終了コード `0`、`incremental_reason=["manifest-invalid"]` になる。 |
| `incremental` | `failure-changed-manifest-corrupt` | `--changed-manifest` の JSON 破損で stdout 空、stderr `BUILDER28_INVALID_OPTION`、終了コード `2`、公開出力維持になる。 |
| `incremental` | `failure-changed-manifest-base-escape` | `--changed-manifest` 内 path が base 外、絶対 path、または `..` escape を含む場合に `BUILDER28_PATH_OUTSIDE_BASE`、終了コード `2`、公開出力維持になる。 |
| `incremental` | `noop-unchanged-pages-kept` | changed page が 0 件の場合でも search index と manifest を最終 page set から再生成し、HTML byte は既存と一致する。 |
| `incremental` | `security-incremental-no-public-write-on-failure` | build 途中 failure、strict warning、changed manifest fatal failure のすべてで公開 `--out`、既存 manifest、既存 search index が変更されない。 |
| `formats` | `success-html-default` | format 未指定で `html` として実行し、REPORT は `output_format="html"`、`output_format_supported=true` になる。 |
| `formats` | `success-html-explicit` | `--format html`、env、設定 file のいずれでも HTML 出力だけを作り、format 由来の追加 file を作らない。 |
| `formats` | `failure-pdf-reserved` | `pdf` を予約値として `BUILDER28_UNSUPPORTED_RESERVED` で拒否し、stdout 空、REPORT なし、staging なし、既存出力維持になる。 |
| `formats` | `failure-epub-reserved` | `epub` を予約値として `BUILDER28_UNSUPPORTED_RESERVED` で拒否し、stdout 空、REPORT なし、staging なし、既存出力維持になる。 |
| `formats` | `failure-unknown-format` | 未知 format を `BUILDER28_INVALID_OPTION` で拒否し、予約値用 error code を使わない。 |
| `formats` | `failure-multiple-format` | 同一 source 内の複数 format 指定を `BUILDER28_INVALID_OPTION` で拒否する。source 優先順位による上書きとは区別する。 |
| `formats` | `security-format-no-reserved-output-files` | `pdf`、`epub`、未知 format、複数指定のすべてで `.pdf`、`.epub`、追加 directory、追加 asset file が生成されない。 |
| `markdown-extensions` | `success-admonition-note-warn-tip` | NOTE / WARN / TIP を `section.adlaire-admonition`、`data-adlaire-admonition`、`div.adlaire-admonition-title` として出力し、title text を固定する。 |
| `markdown-extensions` | `success-badge-color` | `gray`、`blue`、`green`、`yellow`、`red` の badge を `span.adlaire-badge`、`data-adlaire-badge-color` として出力し、label を escape 済み text にする。 |
| `markdown-extensions` | `success-extension-csv-normalization` | csv の trim、空要素無視、重複除去、許可値順序の正規化を確認する。 |
| `markdown-extensions` | `failure-unknown-extension` | 未知 extension を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `markdown-extensions` | `failure-badge-invalid-text-strict` | strict で不正 badge label / color を `BUILDER28_INVALID_OPTION`、終了コード `2` にする。 |
| `markdown-extensions` | `noop-badge-invalid-text-nonstrict` | non-strict で不正 badge を元 text のまま出力し、`markdown_extension_warnings` を加算する。 |
| `markdown-extensions` | `security-extension-escape` | admonition body、badge label、attribute、raw HTML、危険 URL、event handler が escape される。 |
| `markdown-extensions` | `noop-extension-disabled` | extension 未指定時に admonition / badge 構文を特別扱いせず、既存 Markdown 変換結果を維持する。 |
| `code-line-numbers` | `success-line-numbers-fence` | fence option `line-numbers` のある code block だけに `.code-lines`、`.line-no`、`data-line` を出す。 |
| `code-line-numbers` | `success-line-numbers-cli` | CLI 有効時に全 code fence へ line number を出し、`code_line_number_blocks` と `code_line_number_lines` が一致する。 |
| `code-line-numbers` | `success-line-numbers-diff-composition` | §28.9 diff class と line number が同時に存在し、line number node に diff class が付かない。 |
| `code-line-numbers` | `noop-line-numbers-empty-code` | 空 code block には line number を出さず、REPORT count に含めない。 |
| `code-line-numbers` | `noop-line-numbers-disabled` | CLI 無効かつ fence option なしの code block は既存出力と一致する。 |
| `code-line-numbers` | `security-line-numbers-copy-clean` | `expected/site/assets/search-index.json`、copy text fixture、minify 後 HTML のいずれにも line number text が混入しない。 |
| `heading-numbering` | `success-heading-numbering-h2-h3` | h2 / h3 に `span.heading-number` を出し、`1.`、`1.1.` 形式、ASCII space 1 個、`numbered_headings` を固定する。 |
| `heading-numbering` | `success-heading-numbering-implicit-h2` | h2 がない page の h3 で暗黙 h2 counter `1` を使い、`1.1.` から開始する。 |
| `heading-numbering` | `success-heading-numbering-toc-search` | 本文 heading、TOC 表示 text、search index 表示 text に番号を含め、search index 検索対象正規化 text には番号を含めない。 |
| `heading-numbering` | `failure-heading-numbering-unknown-mode` | 未知 mode を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `heading-numbering` | `noop-heading-numbering-none` | `none` で heading、TOC、search index 表示 text を変更せず、`numbered_headings=0` にする。 |
| `heading-numbering` | `security-heading-slug-unchanged` | 採番有無で heading id、anchor href、collapse target、hash history target が byte 単位で一致する。 |

§28.1〜§28.5 の `expected/effects.json` は、少なくとも `created_paths`、`updated_paths`、`preserved_paths`、`deleted_paths`、`forbidden_created_paths`、`forbidden_updated_paths`、`forbidden_deleted_paths`、`external_calls` を持つ。failure / security fixture では、`forbidden_updated_paths` と `forbidden_deleted_paths` に公開 `--out`、既存 `.dependency_manifest.json`、既存 `assets/search-index.json` を必ず含める。

**§28 expected 比較方式固定契約：**

expected 比較は、実装環境差分で揺れないように以下の正規化だけを許可する。下表にない正規化、部分一致、snapshot 差し替え、目視承認は合格条件にしてはならない。

| ファイル | 比較方式 | 許可する正規化 | 禁止 |
|----------|----------|----------------|------|
| `expected/site/*.html` | DOM 構造、tag、属性順、text node の完全一致。 | 改行コードを LF に統一。末尾改行 1 個を許可。 | 属性順の無視、class subset 比較、画像 snapshot だけの比較。 |
| `expected/site/assets/style.css` | selector、property、値、media query の完全一致。 | 空行の連続を 1 行へ正規化可。 | selector の部分一致、未使用 selector の黙認。 |
| `expected/site/assets/app.js` | 対象 handler、storage key、guard、fallback 分岐を含む文字列完全一致。 | 改行コードを LF に統一。 | minify 差分の黙認、外部 script 参照の黙認。 |
| `expected/site/assets/search-index.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | HTML tag 混入、line number 混入、未定義 key の黙認。 |
| `expected/stdout.txt` | 行完全一致。`[REPORT]` は 1 行 key=value 形式。§28 追加 key は既存 §8 key の後ろに ASCII 昇順で並べる。 | 末尾改行 1 個を許可。 | `[REPORT]` key 省略、型違い、件数差分、追加 key 順序違い、compact JSON 内空白の黙認。 |
| `expected/stderr.txt` | 行完全一致。fatal failure 以外は空 file。 | 末尾改行 1 個を許可。 | error code 差分、line 差分、message 差分、stdout warning 混入の黙認。 |
| `expected/effects.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | 外部 call、削除、既存出力破壊の黙認。 |
| `expected/security.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | CDN、credential、raw HTML、secret 残存の黙認。 |

**§28 stdout / stderr / REPORT fixture 固定契約：**

stdout、stderr、`[REPORT]` は、同じ入力から常に同じ順序で出力する。順序は、入力 file path 昇順、line 昇順、section 昇順、code 昇順とする。fatal failure の場合、`[REPORT]` は出力せず、stdout は空 file とする。

| 対象 | 固定内容 |
|------|----------|
| stdout warning | `[WARN] CODE file:line section message` の形式で完全一致。warning は stdout だけに出す。 |
| stderr error | `[ERROR] CODE file:line section message` の形式で完全一致。error は stderr だけに出す。 |
| message | 句点ありの日本語または ASCII 英文に統一し、secret、credential、raw HTML を含めない。 |
| REPORT boolean | `true` / `false` 小文字。 |
| REPORT integer | 0 以上の 10 進数。 |
| REPORT string | double quote 付き JSON string。§8 既存 key は既存形式を維持する。 |
| REPORT array | compact JSON array。要素順は発生順ではなく sorted string 昇順。ただし page order を意味する配列は input path 昇順。 |

**§28 既存出力互換 fixture 固定契約：**

§28 実装 PR では、対象機能の fixture だけでなく、機能を無効化した互換 fixture を 1 件以上含める。既定有効機能は、有効化前後ではなく「機能対象入力なし」の互換 fixture を含める。

| 機能種別 | 互換 fixture |
|----------|--------------|
| 明示有効化機能 | option 未指定時に既存 HTML / CSS / JS / search index / REPORT が変わらない fixture。 |
| 既定有効機能 | 対象 Markdown 記法や対象 DOM が存在しない入力で既存出力が変わらない fixture。 |
| reserved feature | 予約値指定時に出力が作られず、既存出力も破壊しない fixture。 |
| security failure | strict / non-strict の差分と、拒否対象が出力に残らない fixture。 |

**§28 機能別 fixture 最低確認項目固定契約：**

各 §28 fixture は、`docs/details/builder.md` §28 個別固定補足契約の validation、HTML / asset 固定、warning / error、REPORT count を最低 1 件以上の expected で確認する。下表の項目を fixture から省略してはならない。

| 節 | 最低確認項目 |
|----|--------------|
| §28.1 | changed / reused page count、未変更 HTML byte 維持、manifest path validation、search index 全体再生成。 |
| §28.2 | `html` 成功、`pdf` / `epub` 拒否、unknown format 拒否、出力破壊なし。 |
| §28.3 | admonition type 正規化、badge color validation、disabled 時互換、escape。 |
| §28.4 | line number node、copy 対象除外、空 code、line count。 |
| §28.5 | slug 不変、表示番号、TOC / search index 番号、unknown mode 拒否。 |
| §28.6 | toggle target、localStorage key、print 展開、重複 target 検出。 |
| §28.7 | min/max validation、TOC filter、active tracking 対象一致。 |
| §28.8 | fake git、fake file mtime、fallback、RFC3339 UTC 秒精度。 |
| §28.9 | inserted / deleted / context / header class、escape、copy 本文維持。 |
| §28.10 | lazy 属性、外部 URL no-fetch、base 外 path warning、alt escape。 |
| §28.11 | meta head 内順序、禁止 key、重複 last wins、attribute escape。 |
| §28.12 | light / dark / auto、toggle、localStorage、print light。 |
| §28.13 | colon / key-value title、copy 除外、empty title no-op、escape。 |
| §28.14 | key validation、code fence 非置換、missing var、replacement count。 |
| §28.15 | byte count、pre/code 保持、marker validation、disabled 互換。 |
| §28.16 | active link 1 件化、aria-current、fallback scroll、depth sync。 |
| §28.17 | graph TD SVG、unsupported source fallback、external script 不在。 |
| §28.18 | reference order、backlink、duplicate definition warning、undefined strict。 |
| §28.19 | inline / block math、code 内非変換、unclosed delimiter、escape。 |
| §28.20 | pushState、focus、back / forward、missing target no-op。 |
| §28.21 | skip link、landmark、button label、duplicate id strict、keyboard trap 不在。 |
| §28.22 | trigger count、dialog 1 個、Escape / backdrop close、focus trap、alt warning。 |
| §28.23 | URL validation、512 byte 制限、print-only SVG、通常表示非表示。 |
| §28.24 | dl / dt / dd 構造、paragraph 境界、empty term no-op、inline escape。 |
| §28.25 | checked / unchecked、disabled checkbox、nested list、aria、通常 list 非変換。 |

**§28 strict / non-strict fixture 固定契約：**

| ケース | non-strict fixture | strict fixture | 固定する差分 |
|--------|--------------------|----------------|--------------|
| warning で継続できる構文不正 | 終了コード `0`、warning count 増加、fallback 出力あり。 | 終了コード `2`、出力なし。 | stdout / stderr / effects。 |
| base 外 path | 対象参照を無効化し warning。 | 終了コード `2`。 | 参照先 file が作成されないこと。 |
| 未定義参照 | 通常 text または非表示 fallback。 | 終了コード `2`。 | HTML fallback と strict 停止。 |
| reserved feature | 終了コード `2`。 | 終了コード `2`。 | strict 差分なし。 |
| 内部エラー fixture | 終了コード `1`。 | 終了コード `1`。 | 既存出力維持。 |

**§28 fixture manifest 固定 schema：**

`manifest.json` は以下の key を必須とする。未使用 key も省略せず、空配列または空文字で明示する。

| key | 型 | 固定内容 |
|-----|----|----------|
| `name` | string | fixture 名。§28 fixture カタログ固定契約の値と完全一致。 |
| `section` | string | `28-common`、または `28.1`〜`28.25` のいずれか。 |
| `feature_slug` | string | §28 fixture カタログ固定契約の feature slug。 |
| `owner` | string | 常に `builder`。 |
| `collaborators` | array | `runner`、`statefile` など該当する補助 component。該当なしは空配列。 |
| `spec_refs` | array | `docs/details/builder.md §28.x` と `docs/details/fixture.md §28-F` を含める。 |
| `strict` | boolean | strict mode fixture なら `true`。 |
| `expected_exit_code` | integer | `0`、`1`、`2` のいずれか。 |
| `fakes` | object | fake clock、fake git、fake mtime、fake manifest、fake clipboard。不要なら空 object。 |
| `not_applicable` | array | 比較対象外にする file や観点。理由なしの除外は禁止。 |

**§28 PR 証跡固定契約：**

実装 PR 本文には、対象 §28.x、追加 fixture 名、変更した HTML / CSS / JS / REPORT key、strict / non-strict 結果、外部依存なし確認、既存出力互換確認、未実装の §28 機能を列挙する。対象外の §28 機能を先取り実装した場合、または `docs/details/builder.md` §28 に存在しない Markdown 記法、CLI option、CSS class、JS 挙動を追加した場合は未完了として扱う。

**§28 PR 受け入れゲート固定契約：**

§28 実装 PR は、下表を PR 本文または検証ログで確認できる場合だけ受け入れ可能とする。確認できない項目は、実装漏れではなく証跡不足として未完了扱いにする。

| ゲート | PR 証跡 | 不足時の扱い |
|--------|---------|--------------|
| 対象範囲 | 対象 §28.x と対象外 §28.x の一覧。 | 先取り実装または範囲不明として未完了。 |
| fixture catalog | 追加 / 更新した fixture 名の一覧と §28 fixture catalog との対応。 | fixture 不足として未完了。 |
| manifest schema | 各 fixture の `manifest.json` が必須 key を満たす確認。 | fixture schema 不足として未完了。 |
| expected files | HTML、CSS、JS、search index、stdout、stderr、effects、security の該当 expected 更新。 | expected 不足として未完了。 |
| strict / non-strict | strict と non-strict の終了コード、stdout、stderr、effects の差分。 | 異常系未固定として未完了。 |
| REPORT | 追加 / 更新した REPORT key、型、count 単位、既定値。 | runner / API 連携不能として未完了。 |
| compatibility | 対象機能無効時または対象入力なし時の既存出力互換確認。 | 既存出力破壊リスクとして未完了。 |
| security | external call 0 件、CDN / external library 不使用、secret / credential / raw HTML 不在。 | security 不足として未完了。 |
| atomicity | failure fixture で既存出力、manifest、search index が維持される確認。 | 部分更新リスクとして未完了。 |
| not run | 未実施確認がある場合の理由と影響範囲。 | 検証不足として未完了。 |

**§28 fixture 不足時の固定扱い：**

§28 実装中に fixture 不足を発見した場合、実装判断で対象 fixture を省略してはならない。fixture が不足している機能は、実装済みではなく `仕様化済み・未実装` のまま扱う。

| 不足 | 扱い |
|------|------|
| fixture catalog の必須 fixture がない。 | 実装未完了。 |
| strict / non-strict の片方がない。 | 異常系未完了。 |
| expected/security.json が必要なのにない。 | security 検証未完了。 |
| expected/effects.json に既存出力維持がない。 | atomicity 検証未完了。 |
| REPORT key の型または件数が expected にない。 | REPORT 検証未完了。 |
| 比較除外が `not_applicable` に理由付きで記載されていない。 | fixture schema 不備。 |
