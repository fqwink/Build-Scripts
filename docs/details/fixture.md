# Adlaire CI — Fixture 詳細仕様

[`docs/details/fixture.md`](fixture.md) は、[`docs/SPEC.md` 責務文書構成](../SPEC.md#document-responsibility-map) で割り当てられた fixture 証跡責務として、fixture 入力、expected、fake、effects、assertion、実装検証証跡、acceptance checklist、差し戻し条件だけを扱う。owner / collaborator 境界は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) を参照し、owner component 別詳細本文を再定義しない。

---

<a id="fixture-responsibility-boundary"></a>
**責務境界：**

| 項目 | 内容 |
|------|------|
| 証跡責務 | `fixture` |
| 対象 component | `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`security`、`setup`、`release`、`mcp` |
| 持つ内容 | fixture 証跡責務が本文として定義する fixture manifest、assertion、fake、testdata、expected / effects、受け入れ fixture 共通契約、実装検証証跡テンプレート、acceptance checklist、差し戻し条件。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) | fixture / testdata 配置、fake 実装、実装検証証跡。 |
| [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) | API の必須検証、API fixture、API / SDK / UI / 状態ファイル cross fixture 固定。 |
| [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の fixture 配置、fixture カタログ、manifest、assertion、expected/effects、相互整合、component 別検証責務。 |
| [`docs/details/fixture.md` fixture 証跡責務 §27-F-EVIDENCE](fixture.md#sec-27-f-20) | [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の実装検証証跡、受け入れゲート、差し戻し条件、部分失敗・再実行契約。 |

<a id="current-implementation-alignment-evidence"></a>
**現行実装整合証跡：**

現在状態の割当は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照する。該当する実装と必須証跡が契約に一致した場合は、同じ変更単位でこの証跡を更新または解消する。

<a id="phase-1-builder-implementation-evidence"></a>
**Phase 1 builder 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 1 `builder` | [PR #74](https://github.com/fqwink/Build-Scripts/pull/74) で [`components/builder/`](../../components/builder/)、[`components/builder/builder_test.go`](../../components/builder/builder_test.go)、[`docs/SPEC.md`](../SPEC.md) の Phase 1 builder 実装が merge 済みである。 |
| Docker 検証 | `golang:1.22` container で `builder` owner package の実装 artifact と対象 test artifact の `gofmt` 差分なし、`go test ./...` が成功、`git diff --check` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務で Phase 1、[`components/builder/`](../../components/builder/)、および `ビルドスクリプト` 担当機能を `実装済み` と扱う。 |

<a id="phase-2-runner-implementation-evidence"></a>
**Phase 2 runner 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 2 `runner` 実装進捗 | [PR #75](https://github.com/fqwink/Build-Scripts/pull/75) で [`components/runner/`](../../components/runner/) と [`components/runner/runner_test.go`](../../components/runner/runner_test.go) に、runner CLI parse / version、`.repo_config` / `.server_config` / `.branch_config` 統合、GitHub / local watch、複数 target files、dry-run、Commit Status、timeout 限定 retry、build status、起動時設定整合性 check、build trend、branch env、`.pipeline_config` の標準 builder command 拡張と単回読取、approval pending、output manifest、log archive、snapshot `site.tar.gz` / `meta.json`、通知 retry、通知署名、SMTP 未設定分類、command 通知、実行環境記録、build log の top-level commit fields と `remote_build` / `commit_status` / `build_meta` / `transfer_verified` 補助 key、parallel deploy `target_results[]`、dependency manifest、build chain、priority queue、build id / chain id 衝突回避、owner lock release、UTF-8 safe log trim、secret mask、exact builder version precheck の stderr / timeout / token 検証、FailureCategory 固定値と `failure_evidence[]` 保存の実装と Go test 証跡を追加した。本行は Phase 2 `runner` の実装検証証跡として扱う。 |
| Docker 検証 | `golang:1.23` container で `runner` owner package の実装 artifact と対象 test artifact の `gofmt` 差分なし、`go test ./...` が成功、`git diff --check` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 2 `runner` と [`components/runner/`](../../components/runner/) を `実装済み` と扱う。Phase 2 以外の owner、状態ファイル adapter 共通化、正式 fixture directory harness は将来計画または改訂予定の追加検証候補として扱い、Phase 2 runner 完了判定を取り消す根拠にしない。 |

<a id="phase-3-api-request-lifecycle-implementation-evidence"></a>
**Phase 3 api request lifecycle 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 3 `api` request lifecycle 実装進捗 | [`components/api/`](../../components/api/) と [`components/api/api_test.go`](../../components/api/api_test.go) に、固定 `http.Server`、request ID 生成と `X-Request-Id` / `.api_access_log.request_id` 伝播、request ID 生成失敗時の endpoint / auth / request log 停止、path / method / access control / auth の共通判定順、`.access_control` 破損時の fail-closed `503`、`.api_access_log` 必須 field、response status recorder の最初の status 固定、implicit `200`、`http.Flusher` 透過、SSE frame flush、JSON response の送信前 serialize、body 上限 `413` と trailing JSON rejection を追加した。本行は Phase 3 `api request lifecycle` の実装検証証跡として扱う。 |
| Docker 検証 | `golang:1.23` container で `api` owner package の実装 artifact と対象 test artifact を `gofmt` した後、`go test ./...` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 3 `api request lifecycle` を `実装済み` と扱う。Phase 3 以外の owner と正式 fixture directory harness は将来計画または改訂予定の追加検証候補として扱い、Phase 3 request lifecycle 完了判定を取り消す根拠にしない。 |

<a id="phase-4-api-operations-implementation-evidence"></a>
**Phase 4 api operations 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 4 `api` operations 実装検証 | [`main.go`](../../main.go)、[`components/api/`](../../components/api/)、[`main_test.go`](../../main_test.go)、[`components/api/api_test.go`](../../components/api/api_test.go) に、`adlaire-ci-api` の完全一致 basename dispatch、未知 basename 拒否、`--help` / `--version` 優先、`--state-dir` 検証、loopback `--addr` 検証、credentials 起動時検証、listener 起動前の未初期化 / 不正 credentials 拒否、`--init-credentials` の stdin 1 行入力・既存拒否・初期 credentials 作成、`SIGTERM` / `SIGINT` graceful shutdown / forced close 経路、`.admin_credentials` の 64 文字 salt・`login_count`・`last_login_at` schema、TOTP 無効 login の `login_count` 飽和加算と `must_change` `none` / `prompt` / `forced` 境界、`forced` session の通常 endpoint gate と password 変更後解除、同一 password 変更拒否、TOTP 有効 password 成功時の credentials 無変更 ticket、ticket credentials fingerprint、TOTP 成功時の最新 credentials 再読込・`login_count` 更新・`last_accepted_step` 保存・同一 step replay 拒否、IP key 別 memory-only login lock、auth transaction coordinator による login / TOTP / session 管理 endpoint と全 authenticated endpoint の session touch / expiry update の直列化と timeout `409`、credentials / TOTP statefile lock 内 read-modify-write adapter、password login 成功 / 失敗 / lock / TOTP required / TOTP setup / confirm 成功・失敗 / TOTP login 成功・失敗 / TOTP disable / logout / password change / session revoke-all の `.access_log` / `.audit_log` 追記と token / ticket / TOTP secret / otpauth URI 非漏えい、manual / force build の durable queue 投入、runner 非同期起動要求、active queue entry の読取・保持・重複照合、API token trigger 時の queue `requested_by` と build audit actor の token id 記録、maintenance fail-closed、Webhook の署名検証後 JSON parse、event / branch filter、40 桁 lowercase SHA 検証、branch 設定照合、event log、active / waiting queue 照合、delivery 衝突拒否、署名検証済み Webhook queue 成功時の `build_trigger` audit、Webhook audit 失敗時の dispatch 抑止、build stream の有限 SSE `log` / `end` frame、build trend 集計、build chain config GET / POST、approval list / approve / reject、health、PAT status / verify、GitHub rate-limit、status read model、sysinfo / stats / timeline / dashboard / output-meta / verify-output / disk-usage の選択 output target 参照、diagnostics、history `failure_category` filter と重複 id 除外、全 API endpoint の未知・重複 query 拒否、logs search の query / 日付 / level filter、snapshot list / download / delete と delete 成功時の `.config_log` / `.audit_log` 追記、API token の root object schema、`label` / `expires_at` / `last_used_at`、scope 検証・scope 認可・期限 / 失効拒否、token 作成・認証・scope 拒否・失効の `.access_log` 追記、API token event の `timestamp` / `request_id` / actor / target / result / `remote_addr` を持つ `.audit_log` 追記、`GET /api/audit-log` の `actor` filter・`timestamp` 降順・未知 action / result 拒否、API rate limit policy の strict schema、login / session / API token / 署名検証済み Webhook の固定窓判定、actor key と IP key の同一 lock 内更新、上限超過時の count 非増加と `permission_denied` audit、`state_summary` array、必須 token log 追記失敗時の token 本体非返却、`.api_tokens` の作成・認証時 `last_used_at` 更新・失効を statefile lock 内 read-modify-write で直列化し、lock conflict を `409` として返す処理、`POST /api/config/validate` の副作用なし検証・`valid=false`・未知 key / secret 風 key 拒否、top-level `BackupObject`、`RestoreObject` の必須 key 検証、secret file 保持・更新・削除、restore no-op、固定書込順の一部、`.hooks` / `.alert_rules` / `.tag_rules` の object wrapper 保存、schedule interval の `.server_config` 保存・systemd timer drop-in 反映・systemctl 呼出順・no-op 無副作用・partial failure log、circuit reset と maintenance mutation の log 追記の API 実装と Go test 証跡を追加した。本行は Phase 4 `api operations` の実装検証証跡として扱う。 |
| Docker 検証 | `golang:1.25.1` container で `main.go`、`main_test.go`、`api` owner package の `gofmt` 差分なし、Phase 4 API 対象 test selection と `go test ./... -count=1` が成功済みである。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 4 `api operations` を `実装済み` と扱う。追加管理 API の未実装項目、mcp、および正式 fixture directory harness は将来計画または改訂予定の追加検証候補として扱い、Phase 4 `api operations` 完了判定を取り消す根拠にしない。 |

<a id="phase-5-sdk-implementation-evidence"></a>
**Phase 5 sdk 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 5 `sdk` 実装検証 | [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) に、`baseUrl` の absolute `http` / `https` origin 固定、path / userinfo / query / fragment / `/api` 付き拒否、API 契約表と一致する public method、`encodeURIComponent(String(value))` による `%20` query 生成、`failure_category` query、access / notify / config log の `limit` / `offset` query、`createToken()` の `{label,scopes,expires_at}` body、`setRepoConfig()` の `{owner,repo}` 限定 patch、`setSmtpConfig()` の未定義 password 除外、30 秒 timeout の fetch / JSON text / Blob 読取全体適用、JSON media type 解析、binary `application/octet-stream` 検証、`^[A-Za-z0-9_-]{1,64}$` id 送信前検証、`streamBuild()` の `text/event-stream` 検証、fatal UTF-8、`data: ` 単一行 frame、`log` / `end` exact key、最終 `end` 後 EOF、`StreamHandle.close` / `closed` / `error` / `done` terminal state を実装した。 |
| 契約検証 | [`sdk_contract_test.go`](../../sdk_contract_test.go) で [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK column と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) の public method を照合し、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の method order、固定 private method、主要 request shape、引数検証 guard、`URLSearchParams`、`response.json()`、`EventSource`、browser storage、global export の不使用を検証する。 |
| Docker 検証 | `golang:1.22` container で `gofmt -w sdk_contract_test.go` 実行後、`go test ./...` が成功した。 |
| Deno 検証 | `denoland/deno:latest` container の `deno 2.9.7 stable` で `deno check admin/adlaire-ci-sdk.js` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 5 `sdk` と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) を `実装済み` と扱う。mcp、追加管理 API / SDK / UI の未実装項目、および正式 fixture directory harness は将来計画または改訂予定の追加検証候補として扱い、Phase 5 `sdk` 完了判定を取り消す根拠にしない。 |

<a id="phase-6-ui-implementation-evidence"></a>
**Phase 6 ui 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 6 `ui` 実装検証 | [`admin/index.html`](../../admin/index.html) に、same-origin absolute base URL 検証、`/api` 既定値排除、ES Module SDK import、ログイン後初期取得順、`must_change` / TOTP 分岐、`StreamHandle.done` 監視、one-time token / TOTP secret / otpauth URI の generation 一致消去、copy button、`409` / `422` / `429` / `503` 表示処理、field error focus、notify / repo / branch / duration anomaly request shape、backup / restore、Webhook events、output metadata、build trends、circuit reset、承認 approve / reject、履歴 comment / flag / tags / rollback、schedule interval / allowed hours / force interval / cooldown、pipeline config、build chain config、hook delete / log、alert rule、tag rule、token revoke、snapshot download / delete、focus-visible を実装した。 |
| 契約検証 | [`ui_contract_test.go`](../../ui_contract_test.go) で [`docs/details/ui.md`](ui.md) 詳細本文責務の DOM 固定値、SDK 境界、base URL、禁止 API、one-time secret、stream terminal、UI 操作表の主要 SDK method 接続、主要 request shape、status 別 error handling を静的照合する。 |
| Docker 検証 | `golang:1.22` container で `gofmt -w ui_contract_test.go` 実行後、`go test ./...` が成功した。 |
| Deno 検証 | `denoland/deno:latest` container の `deno 2.9.7 stable` で `deno check admin/adlaire-ci-sdk.js` が成功した。`admin/index.html` は module script を抽出し、SDK import を stub 化した `deno eval --allow-read` 構文検証が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 6 `ui`、[`admin/index.html`](../../admin/index.html)、および `標準管理ツール UI 契約` を `実装済み` と扱う。正式 `testdata/ui/` fixture directory harness、管理 UI 静的配布物構成、管理 UI 静的 HTTP 配信、mcp は将来計画または改訂予定の追加検証候補として扱い、Phase 6 `ui` 完了判定を取り消す根拠にしない。 |

<a id="phase-7-admin-cli-implementation-evidence"></a>
**Phase 7 admin CLI 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 7 `admin` CLI 実装検証 | [`main.go`](../../main.go)、[`main_test.go`](../../main_test.go)、[`components/admin/`](../../components/admin/)、[`components/admin/admin_test.go`](../../components/admin/admin_test.go) に、`adlaire-ci-admin` の完全一致 basename dispatch、`--help` / `--version` 優先、argv token safety、`--api-url` / `--token` / `--json` parse、URL 正規化、token 検証、7 command 固定表、request method / path / header / body、redirect 不追従、proxy 無効、retry なし、30 秒 timeout、1 MiB response body 上限、Content-Type 検証、JSON object 単一値検証、command 別 human stdout、`--json` wire body 出力、HTTP error / network error / invalid response、token 非表示を実装した。 |
| fixture 証跡 | [`testdata/admin/cli/partial-admin-cli-lifecycle/`](../../testdata/admin/cli/partial-admin-cli-lifecycle/)、[`testdata/admin/cli/success-admin-cli-transport/`](../../testdata/admin/cli/success-admin-cli-transport/)、[`testdata/admin/cli/failure-admin-cli-output-errors/`](../../testdata/admin/cli/failure-admin-cli-output-errors/)、[`testdata/admin/cli/security-admin-cli-secret-redaction/`](../../testdata/admin/cli/security-admin-cli-secret-redaction/) に Admin CLI fixture 固定契約の必須 fixture を配置し、[`components/admin/admin_test.go`](../../components/admin/admin_test.go) で必須 fixture file の存在と JSON 妥当性を検証する。 |
| Docker 検証 | `golang:1.22` container で `main.go`、`main_test.go`、`admin` owner package の `gofmt` 差分なし、`go test ./...` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 7 `admin` CLI 管理クライアント、[`components/admin/`](../../components/admin/)、および `CLI 管理クライアント` を `実装済み` と扱う。管理 UI 静的配布物構成、管理 UI 静的 HTTP 配信、mcp、および正式 fixture harness の追加拡張は将来計画または改訂予定の追加検証候補として扱い、Phase 7 `admin` CLI 完了判定を取り消す根拠にしない。 |

<a id="phase-8-setup-implementation-evidence"></a>
**Phase 8 setup 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 8 `setup` 実装検証 | [`main.go`](../../main.go)、[`main_test.go`](../../main_test.go)、[`components/setup/`](../../components/setup/)、[`components/setup/setup_test.go`](../../components/setup/setup_test.go) に、`adlaire-ci-setup` の完全一致 basename dispatch、`--help` / `--version` 優先、argv token safety、`install` / `install-api` / `update` mode、mode 別 option 許可、version / repository / OS arch / path / service user validation、Release asset 取得、redirect 境界、size 上限、`SHA256SUMS` 検証、binary 配置、secret / state create-only、admin archive 安全検証、API credentials 初期化、runner / API systemd unit 書込、activate / restart / `systemctl is-active` / `systemctl cat` 検証、Go `net/http` health check、API 導入済み update 判定、backup / rollback、stdout / stderr 固定契約を実装した。 |
| fixture 証跡 | [`testdata/setup/success-setup-admin-release-asset-layout/`](../../testdata/setup/success-setup-admin-release-asset-layout/)、[`testdata/setup/failure-setup-download-boundary/`](../../testdata/setup/failure-setup-download-boundary/)、[`testdata/setup/failure-setup-api-version-cohort/`](../../testdata/setup/failure-setup-api-version-cohort/)、[`testdata/setup/security-setup-admin-archive-boundary/`](../../testdata/setup/security-setup-admin-archive-boundary/)、[`testdata/setup/partial-setup-systemd-rollback-boundary/`](../../testdata/setup/partial-setup-systemd-rollback-boundary/)、[`testdata/setup/partial-setup-api-runner-dispatch/`](../../testdata/setup/partial-setup-api-runner-dispatch/)、[`testdata/setup/security-setup-secret-preservation/`](../../testdata/setup/security-setup-secret-preservation/) に setup / admin / Release asset 連動 fixture 固定契約の必須 fixture を配置し、[`components/setup/setup_test.go`](../../components/setup/setup_test.go) で必須 fixture file の存在と JSON 妥当性を検証する。 |
| Docker 検証 | `golang:1.22` container で `main.go`、`main_test.go`、`setup` owner package の実装 artifact と対象 test artifact を `gofmt` した後、`go test ./...` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 8 `setup`、[`components/setup/`](../../components/setup/)、および `初回セットアップ`、`管理 API 導入`、`バイナリアップデート`、`アップデート rollback` を `実装済み` と扱う。mcp、および正式 fixture harness の追加拡張は将来計画または改訂予定の追加検証候補として扱い、Phase 8 `setup` 完了判定を取り消す根拠にしない。 |

<a id="phase-9-release-implementation-evidence"></a>
**Phase 9 release 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 9 `release` 実装検証 | [`main.go`](../../main.go)、[`main_test.go`](../../main_test.go)、[`components/release/`](../../components/release/)、[`components/release/release_test.go`](../../components/release/release_test.go) に、`adlaire-ci-release` の完全一致 basename dispatch、`--help` / `--version` 優先、binary version 注入、Release asset basename の `-linux-amd64` 正規化、CLI parse / validation、clean checkout / tag / commit 検証、GitHub read-only 事前確認、immutable source snapshot A / B 展開、gofmt / go test / go build、admin archive 生成、checksum manifest、再現性確認、atomic output、GitHub draft 作成、asset upload、asset list / download digest 検証、publish、失敗時 draft cleanup、固定 stdout / stderr を実装した。 |
| fixture 証跡 | [`testdata/release/`](../../testdata/release/) に [`docs/details/fixture.md` fixture 証跡責務 Release fixture 固定契約](#release-fixture-contract) の 21 fixture を配置し、[`components/release/release_test.go`](../../components/release/release_test.go) で必須 fixture file の存在と JSON 妥当性を検証する。 |
| Docker 検証 | `golang:1.22` container で `main.go`、`main_test.go`、`builder` / `admin` / `setup` / `release` owner package の対象 Go artifact を `gofmt` した後、`go test ./...` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 9 `release`、[`components/release/`](../../components/release/)、および `GitHub Release 成果物生成・公開前検証・公開` を `実装済み` と扱う。Phase 10 `mcp` は [Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence) で扱う。 |

<a id="phase-10-mcp-implementation-evidence"></a>
**Phase 10 mcp 実装検証証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 10 `mcp` 実装検証 | [`main.go`](../../main.go)、[`main_test.go`](../../main_test.go)、[`components/mcp/`](../../components/mcp/)、[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go)、[`components/release/`](../../components/release/) に、`adlaire-ci-mcp` の完全一致 basename dispatch、`--help` / `--version` 優先、binary version 注入、CLI parse / validation、loopback listener、token auth、JSON-RPC lifecycle、tools / resources / prompts、sampling request、SSE、scope、read-only、confirmation、audit、metrics、timeout、statefile 連携を実装した。 |
| fixture 証跡 | [`testdata/mcp/`](../../testdata/mcp/) に [`docs/details/fixture.md` fixture 証跡責務 §30-F MCP fixture 固定契約](#mcp-fixture-contract) の 10 fixture を配置し、[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go) で必須 fixture file の存在と JSON 妥当性を検証する。 |
| Docker 検証 | `golang:1.22` container で `main.go`、`main_test.go`、`mcp` / `release` owner package の対象 Go artifact を `gofmt` した後、`go test ./...` が成功した。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 10 `mcp`、[`components/mcp/`](../../components/mcp/)、[`main.go`](../../main.go)、および MCP サーバー機能群を `実装済み` と扱う。 |

<a id="phase-11-quality-gate-implementation-evidence"></a>
**Phase 11 バグ修正ゼロ化 quality gate 証跡：**

| 対象 | 証跡 |
|------|------|
| Phase 11 `quality gate` 実装検証 | [`main_test.go`](../../main_test.go) に、標準実装 artifact inventory、`components/` の subdirectory 禁止、`cmd/` 非採用、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) 索引到達、fixture manifest の `name` / directory identity、component owner / collaborator 整合、参照先 Markdown file / anchor 到達、`fake_clock` UTC 秒精度、`not_applicable` / `missing_state` / `input_files` の安全相対 path、input / expected JSON 妥当性、assertion に対応する expected evidence、Phase 11 ROADMAP 状態、Phase 11 専用詳細ファイルの不在、`Part` 表現不在、未登録 fixture directory を検出する横断 gate を追加した。 |
| fixture 証跡 | [`testdata/builder/`](../../testdata/builder/)、[`testdata/admin/cli/`](../../testdata/admin/cli/)、[`testdata/setup/`](../../testdata/setup/)、[`testdata/release/`](../../testdata/release/)、[`testdata/mcp/`](../../testdata/mcp/)、および Phase 11 root coverage fixture 配下の正式 fixture manifest 56 件を [`main_test.go`](../../main_test.go) の Phase 11 manifest gate で横断確認する。 |
| Docker 検証 | `golang:1.22.12` container で `gofmt -w main_test.go components/builder/builder_test.go`、`gofmt -l main.go main_test.go components/*/*.go sdk_contract_test.go ui_contract_test.go`、`go test ./... -count=1`、`go test -race ./... -count=1` が成功することを Phase 11 完了検証とする。Phase 11 mutation 選択では、一時 copy 上で `main.go` の未知 basename 判定反転、fixture manifest `name` / directory identity 破壊、Phase 11 ROADMAP 状態反転、builder formal fixture expected 除去、statefile 共通永続化状態反転を実行し、対応 test が失敗することを確認対象に含める。集計は `survived=0` とする。 |
| 状態反映 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務では、Phase 11、状態ファイル共通永続化契約、検証基盤の Phase 11 対象機能群を `実装済み` と扱う。 |

<a id="phase-11-quality-gate-closure"></a>
**Phase 11 完了 closure：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 11 の完了証跡を記録する。Phase 11 の現在状態は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、closure record set の条件は [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) を正本とする。

Phase 11 closure record set は次の 18 record とする。

| closure_item | status | evidence | open_items |
|--------------|--------|----------|------------|
| `scope_inventory` | `closed` | [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan) | `[]` |
| `test_gap_inventory_closure` | `closed` | `test_gap_inventory open=0`、[`main_test.go`](../../main_test.go) Phase 11 manifest / ROADMAP / stale reference gate | `[]` |
| `test_improvement_batch_closure` | `closed` | builder formal fixture 6 件、Phase 11 root coverage fixture 8 件、既存 admin / setup / release / mcp fixture manifest 接続 | `[]` |
| `traceability_closure` | `closed` | [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) の実在所在、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) の owner 入口、owner 詳細本文、fixture 証跡の相互参照 | `[]` |
| `requirement_coverage_closure` | `closed` | Phase 11 対象機能 12 件と `状態ファイル共通永続化契約` を [`docs/ROADMAP.md`](../ROADMAP.md) で `実装済み` に統一 | `[]` |
| `fixture_root_closure` | `closed` | [fixture root coverage matrix 固定契約](#fixture-root-coverage-matrix-contract)、[未作成 fixture root closure record 固定契約](#fixture-root-missing-closure-record-contract)、builder formal fixture 6 件、root coverage fixture 8 件 | `[]` |
| `execution_evidence_closure` | `closed` | `gofmt -l main.go main_test.go components/*/*.go sdk_contract_test.go ui_contract_test.go`、`go test ./... -count=1`、`go test -race ./... -count=1` | `[]` |
| `oracle_closure` | `closed` | `expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`、既存 component test assertion | `[]` |
| `failure_diagnostics_closure` | `closed` | CLI exit code、stderr、stdout warning、strict publish prevention、manifest decode error、stale reference failure message | `[]` |
| `boundary_failure_matrix_closure` | `closed` | builder strict / URL safety / empty-dir / atomic output、runner / API / admin / setup / release / mcp 既存 fixture manifest | `[]` |
| `isolation_closure` | `closed` | `t.TempDir()` 出力、formal fixture read-only input、no external call expected、statefile lock / atomic write contract | `[]` |
| `determinism_closure` | `closed` | `fake_clock` UTC 秒精度、stable fixture manifest identity、idempotency assertion、生成物 byte comparison | `[]` |
| `concurrency_race_closure` | `closed` | `go test -race ./... -count=1`、statefile lock / append / atomic write 経路、MCP / API mutex 経路 | `[]` |
| `mutation_closure` | `closed` | mutation selection は `main dispatch exact basename`、`fixture manifest identity`、`Phase 11 state gate`、`builder formal expected`、`statefile common persistence state` を対象とし、集計を `killed=5`、`survived=0` とする。 | `[]` |
| `harness_self_verification_closure` | `closed` | [`main_test.go`](../../main_test.go) が missing manifest、重複 manifest name、unsupported assertion、missing expected、missing anchor、unexpected fixture directory、stale Phase 11 term を検出する。 | `[]` |
| `contract_drift_closure` | `closed` | test artifact traceability 固定契約、non-dedicated owner test routing 固定契約、ROADMAP 状態、DOCUMENT_INDEX 所在、owner 詳細本文の接続を確認する。 | `[]` |
| `cross_owner_closure` | `closed` | builder、runner、api、admin、setup、release、mcp、sdk、ui、statefile、archive、commitstatus、security の Phase 11 対象入口が DETAIL_INDEX から到達可能。 | `[]` |
| `final_open_item_count` | `closed` | `0` | `[]` |

<a id="phase-11-fixture-harness-reference"></a>
**Phase 11 fixture harness 参照：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 11 の fixture、expected、fake、実装検証証跡の一般形式と配置契約だけを固定する。Phase 11 の完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、現在状態と対象外理由は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、owner 割当入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) を参照する。

Phase 11 対象の fixture root、manifest、input、expected、effects、security、fake transcript、実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#sec-0g-8-f)、[`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](fixture.md#test-gap-inventory-record-contract)、[`docs/details/fixture.md` fixture 証跡責務 test improvement batch closure 固定契約](fixture.md#test-improvement-batch-closure-contract)、[`docs/details/fixture.md` fixture 証跡責務 test assertion identity / failure diagnostics evidence set 固定契約](fixture.md#test-assertion-failure-diagnostics-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](fixture.md#test-boundary-failure-matrix-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](fixture.md#test-concurrency-race-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 race trigger matrix 固定契約](fixture.md#race-trigger-matrix-contract)、[`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](fixture.md#mutation-test-evidence-contract)、[`docs/details/fixture.md` fixture 証跡責務 mutation selection ledger 固定契約](fixture.md#mutation-selection-ledger-contract)、[`docs/details/fixture.md` fixture 証跡責務 test / contract drift 証跡固定契約](fixture.md#test-contract-drift-evidence-contract)、[`docs/details/fixture.md` fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](fixture.md#sec-27-f-20)、および [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8) の Phase 11 完了判定を同時に満たす。Phase 11 対象 root の網羅判定は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) で割り当てた `testdata/builder/`、`testdata/runner/`、`testdata/api/`、`testdata/sdk/`、`testdata/ui/`、`testdata/statefile/`、`testdata/archive/`、`testdata/commitstatus/`、`testdata/security/`、`testdata/admin/cli/`、`testdata/setup/`、`testdata/release/`、`testdata/mcp/` の各 root について、存在、manifest、expected、fake、harness 接続、対象外理由のいずれかへ到達できることで判定する。

Phase 11 対象 root では、directory 名、`manifest.json.name`、fixture catalog 名が 1 対 1 に一致する場合だけ正式 fixture として扱う。`* 2` suffix 付き directory、同一 `manifest.json.name` を持つ複数 directory、fixture catalog 未登録 directory、harness から参照されない directory、expected だけを持つ directory は正式 fixture として扱わず、[`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8) の `fixture root identity zero duplicate` 判定で未完了として扱う。

<a id="phase-12-quality-gate-evidence"></a>
**Phase 12 実装品質ゲート再構築証跡：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 12 の fixture、expected、fake、実行型 harness、mutation、race、contract drift、closure record の証跡だけを固定する。Phase 12 の現在状態と実装割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](../DETAIL_INDEX.md#phase-12-quality-gate-entry)、Go 実装配置は [`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3)、Phase 完了単位は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、テスト方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を参照する。

Phase 12 の証跡 package は、下表の全対象を同一 closure record set へ接続する。1 件でも `open` が残る場合、Phase 12 を `実装済み` にしてはならない。

| 証跡対象 | 必須証跡 | 未完了条件 |
|----------|----------|------------|
| 契約不整合・状態安全性 | API credential 初期化と setup stdout、Admin CLI URL と API route 表、SDK method / UI 操作 / HTTP method / auth requirement / response schema、statefile read-modify-write lock、JSON Lines append lock / flush / 破損検出、strict JSON schema、secret redaction、atomic write / fsync / permission / crash recovery、冪等性、CLI exit code / stdout / stderr の照合 record。 | 照合対象の片側だけを確認している、route / method / auth / schema の差分が残る、statefile / JSON Lines / secret / atomic write の failure case が未記録、または差分を後続変更で閉じる説明がある。 |
| 実行型 fixture harness | fixture input を production entrypoint へ渡した実行記録、actual / expected / effects / security expected の比較記録、clock / filesystem / HTTP / command / systemd fake adapter の binding 記録、harness self-verification の positive / negative control。 | fixture directory の存在だけ、mock 結果だけ、production entrypoint を通らない実行、expected 比較なし、fake adapter 未接続、または harness が意図的な不正を fail にできない。 |
| mutation / race / queue state machine | production code への mutation selection ledger、mutation evidence set、`survived=0`、race detector または代替 interleaving 証跡、runner queue 遷移と finalizer の状態機械 coverage、goroutine / listener / shutdown / SSE / timer / queue / worker leak 証跡。 | mutation class 未定義、`survived>0`、race trigger 未判定、queue 状態遷移未網羅、goroutine / listener / timer leak 未検出、または race detector 未実行理由に正本 anchor がない。 |
| owner package 5 ファイル固定 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 Go owner package target owner](../DETAIL_INDEX.md#phase-12-go-owner-package-targets) の各 owner package が [`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3) の `<owner>.go`、`model.go`、`validate.go`、`execute.go`、`<owner>_test.go` だけを持つ検査結果、禁止ファイル名検出結果、巨大 component file 分割後の owner 責務接続、root 横断 test の除外理由。 | 6 ファイル目、禁止名、補助 package、owner 外逃がし、root 横断 test の誤分類、または 5 ファイルに収まらない責務を仕様分割せず残している。 |
| ROADMAP / fixture 分類 / MCP / release | ROADMAP 文言検査 test と実装品質 test の分離、root coverage fixture の `inventory fixture` 分類、MCP no-op の明示的未実装 error または実処理、version tag と release notes、release reproducibility の証跡。 | 文書文言 test と品質 test の責務混在、root coverage fixture の分類未定義、MCP no-op が黙って成功する、version tag / release notes / release reproducibility の証跡がない。 |
| CI required checks | GitHub Actions の通常 test、race、lint、fixture 実行 required check の結果、未実行時の skip / 未実行証跡、Pull Request 本文からの証跡到達。 | required check 未設定、未実行、skip を成功扱いにしている、または Pull Request 本文から対象証跡へ到達できない。 |
| final closure | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の 18 record へ Phase 12 対象を接続し、下表の全集計値を完了値で記録する。 | final open item が 0 でない、数値未記録、対象外理由 anchor 不足、または残件を後続 PR へ送っている。 |

Phase 12 closure record set は、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の 18 record を使用し、各 record の `scope` に `phase-12-quality-gate-reconstruction` を含める。Phase 12 では以下の集計値を同じ closure record set 内に記録する。

Phase 12 closure record set の記録先は、Phase 12 実装 PR 本文の `Verification` に置く implementation PR evidence package を必須とする。正式 fixture root を追加する場合は、同じ PR で `testdata/phase12/quality-gate-reconstruction/` を作成し、`manifest.json`、`input/`、`expected/`、closure record set への参照を置く。PR 本文または正式 fixture root のいずれにも closure record set の所在がない場合、Phase 12 を完了扱いにしてはならない。

Phase 12 の正式 fixture root は [`testdata/phase12/quality-gate-reconstruction/`](../../testdata/phase12/quality-gate-reconstruction/) とする。同 root の `manifest.json`、`input/scope.json`、`expected/effects.json` は Phase 12 closure counter を記録し、[`main_test.go`](../../main_test.go) の `TestPhase12QualityGateEvidence` と `TestPhase12StandardArtifactInventory` が `final_open_item_count=0`、owner package 5 ファイル固定、required check workflow 実在を検査する。

Phase 12 の CI required check は、GitHub workflow が作成されていない場合でも対象外にしてはならない。Phase 12 実装 PR は `.github/workflows/phase12-quality-gate.yml` を作成するか、同等の required check 名を GitHub 側で必須化した証跡を同じ PR 本文へ記録する。GitHub workflow は YAML 禁止の対象外であるが、workflow から実行する Adlaire CI 入出力、fixture manifest、expected、state、設定形式は JSON 契約に従う。

| required check name | 必須実行 | 完了条件 |
|---------------------|----------|----------|
| `phase12-go-format` | `rg --files -g '*.go'` が返す全 Go file を対象に `gofmt -l` を実行し、Phase 12 後の owner package file set を検査する。 | `gofmt -l` 出力空、owner package が 5 ファイル固定、`components/` 直下 `.go` file 0 件。 |
| `phase12-go-test` | Go stable toolchain で `go test ./... -count=1`。 | exit code `0`、skip / 未実行 record なし、対象 owner と fixture root が closure record set に接続済み。 |
| `phase12-go-race` | Go stable toolchain で `go test -race ./... -count=1`。 | exit code `0`、race trigger open `0`、goroutine / listener / timer / worker leak の対象外理由または証跡あり。 |
| `phase12-go-vet` | Go stable toolchain で `go vet ./...`。 | exit code `0`。外部 linter、npm、Node.js、Marketplace action で代替しない。 |
| `phase12-deno-check-sdk` | Deno stable runtime で `deno check admin/adlaire-ci-sdk.js`。 | SDK artifact が存在する限り exit code `0`。Deno 不在、JS 未変更、または Node.js 代替を完了扱いにしない。 |
| `phase12-fixture-harness` | Phase 12 で実在化した production entrypoint 実行型 fixture harness。 | `testdata/phase12/quality-gate-reconstruction/` または PR evidence package の closure record set に `phase12_fixture_harness_open_count=0` を記録する。 |
| `phase12-mutation` | production code、harness、assertion、expected 比較、security assertion、state diff assertion を対象にした mutation selection と mutation evidence。 | `phase12_mutation_survived_count=0`、`equivalent` / `invalid` は責務正本 anchor 付き。 |

| 集計値 | 完了値 | 未完了条件 |
|--------|--------|------------|
| `phase12_contract_mismatch_count` | `0` | API / SDK / UI / Admin CLI / route / method / auth / response schema / CLI 出力契約の差分が 1 件以上ある。 |
| `phase12_state_safety_open_count` | `0` | statefile read-modify-write、JSON Lines append、atomic write、fsync、permission、crash recovery、strict JSON schema の未解消項目が 1 件以上ある。 |
| `phase12_fixture_harness_open_count` | `0` | production entrypoint 実行、actual / expected 比較、fake adapter binding、harness self-verification の未接続が 1 件以上ある。 |
| `phase12_mutation_survived_count` | `0` | production code、harness、assertion、expected 比較、security assertion、state diff assertion の適用可能 mutation が 1 件以上 survived である。 |
| `phase12_race_trigger_open_count` | `0` | race trigger、race detector、代替 interleaving、goroutine / listener / timer / worker leak の未判定が 1 件以上ある。 |
| `phase12_owner_package_violation_count` | `0` | owner package が 5 ファイル固定に違反する、package 名が owner 名と一致しない、import path が `github.com/fqwink/build-scripts/components/<owner>` 以外である、`components/` 直下 `.go` file が残る、禁止名を持つ、補助 package へ逃がす、旧 `components/<owner>.go` / `components/<owner>_test.go` が残る、cross-cutting owner の duplicate 実装が呼び出し元 package に残る、または root 横断 test を owner package 内へ混在させる。 |
| `phase12_ci_required_check_open_count` | `0` | 通常 test、race、lint、fixture 実行 required check の未設定、未実行、skip 成功扱いが 1 件以上ある。 |
| `phase12_release_reproducibility_open_count` | `0` | version tag、release notes、checksum、archive、GitHub Release boundary、再取得検証の未解消項目が 1 件以上ある。 |
| `final_open_item_count` | `0` | 上記集計値または 18 record の `open_items` に残件がある。 |

Phase 12 の実装 PR 本文は、[`docs/details/fixture.md`](fixture.md) fixture 証跡責務 implementation PR evidence template 固定契約に加えて、上表の集計値、closure record set 所在、対象外理由 anchor、状態復帰が必要になった [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務の行、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の更新有無を記録する。集計値を口頭説明、検証コマンド名、または Pull Request の Summary だけで代替してはならない。

<a id="phase-13-implementation-alignment-quality-evidence"></a>
**Phase 13 実装整合・品質改善証跡：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 13 の fixture、expected、fake、fault injection、mutation、race、integration、E2E、CI、release evidence、closure record の証跡だけを固定する。Phase 13 の現在状態と実装割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry)、Go 実装配置は [`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3)、Phase 完了単位は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、テスト方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を参照する。

Phase 13 の証跡 package は、下表の全対象を同一 closure record set へ接続する。1 件でも `open` が残る場合、Phase 13 を `実装済み` にしてはならない。

| 証跡対象 | 必須証跡 | 未完了条件 |
|----------|----------|------------|
| owner package 責務純度 | `components/<owner>/` の 5 ファイル固定検査、file role 検査、package 名検査、import path 検査、6 ファイル目・サブディレクトリ・空ファイル・ダミー実装・旧 root path 不在検査。 | 5 ファイル固定違反、空 file、dummy success、helper / utils / common 追加、旧 root path 参照、または責務外処理を別 owner へ逃がす実装が 1 件以上ある。 |
| contract parity | API route 表、Admin CLI command 表、SDK public method 表、UI operation 表、MCP tool 表、setup stdout、credential 初期化、HTTP method、path、query、request、response、status、auth、token 取得元の照合 record。 | 片側だけの修正、route / method / auth / response shape 差分、stdout 契約差分、token argv 取得、未照合 endpoint が 1 件以上ある。 |
| state safety | statefile adapter 経由の read-modify-write lock、file / parent directory fsync、atomic rename、symlink 非追従、state directory `0700`、owner / mode 起動時検証、migration、stale lock、直接状態更新不在の証跡。 | API、runner、MCP、archive が状態 file を直接更新する、lock 範囲外 read-modify-write、fsync / rename / mode / owner 未検証、stale lock 復旧未定義が 1 件以上ある。 |
| JSON Lines safety | append lock、flush、fsync、破損行検出、破損 file 隔離、復旧対象登録、黙殺禁止、必須 log 書込み失敗時の失敗伝播。 | 破損行を読み飛ばして成功扱いにする、best-effort 扱いが禁止された audit / access / config log write failure を無視する、または復旧記録がない。 |
| security boundary | Go 標準ライブラリ内製 KDF、外部 password hash / KDF 依存不在、既存 SHA-256 反復 KDF 不在、token file / stdin、secret mask、署名、file safety、credential rotation、log / fixture / release notes の secret 不在証跡。 | 外部 KDF 依存、旧 SHA-256 反復 KDF 残存、token argv 残存、secret 出力、token hash / prefix / length 漏えい、credential rotation 未検証が 1 件以上ある。 |
| queue / recovery | `waiting`、`active`、`cancelling`、`succeeded`、`failed`、`cancelled`、`recovery_required` の状態機械、waiting to active atomic transition、ID 採番、queue 上限、panic / timeout / cancel / 保存失敗 finalizer、active 復旧、at-least-once 実行。 | 状態遷移表にない遷移、未原子的な active 化、finalizer 未実行経路、再起動後 active 放置、at-least-once 証跡欠落が 1 件以上ある。 |
| archive / commitstatus / release | archive owner の snapshot / compress / digest / verify / download / restore、commitstatus owner の retry / rate limit / timeout、release owner の tag / GitHub Release / SHA256SUMS / signature / SBOM / reproducible build evidence。 | 呼び出し元 owner に archive / commitstatus duplicate 実装が残る、release 証跡不足、SHA256SUMS / 署名 / SBOM / 再現ビルド証跡が未接続である。 |
| MCP real behavior | `resendWebhook`、`subscribe`、`unsubscribe`、`sampling` の実動作証跡、未実装機能の明示 error、HTTP timeout、body 上限、graceful shutdown、statefile adapter 接続。 | no-op success、未実装 success、statefile 直接更新、HTTP lifecycle 未検証、sampling の外部 AI API 直接呼出しが 1 件以上ある。 |
| external boundary | API / MCP HTTP timeout、body 上限、graceful shutdown、Webhook redirect / SSRF / private IP / DNS rebinding 防止、SSH host key / known_hosts / timeout / remote path、systemd 専用 user / 最小権限 / 書込み先限定。 | redirect 追従、private IP 許可、DNS rebinding 未検出、known_hosts 未検証、systemd root 前提、write path 無制限が 1 件以上ある。 |
| executable fixture | fixture input を production entrypoint へ渡し、actual と expected を比較した実行記録、fixture 数や JSON 妥当性だけを完了証拠にしない negative control。 | fixture count、JSON parse、file 存在だけの合格、production entrypoint 未通過、expected 比較なし、negative control が fail しない。 |
| mutation / race / fault | production code mutation selection、mutation evidence set、`survived=0`、race detector または代替 interleaving、disk full、permission denied、short write、fsync failure、rename failure、process kill。 | mutation survivor、race trigger 未判定、fault injection 未実行、skip に責務正本 anchor がない、障害を成功扱いする。 |
| integration / E2E | Admin to API、UI to SDK to API、runner to statefile、MCP to statefile の統合テスト、setup、install-api、update、rollback、Release、実インストール E2E。 | 単体 test だけで統合済み扱い、stub だけで E2E 扱い、実インストール未実行の未実行証跡欠落、連携先副作用未比較が 1 件以上ある。 |
| CI / GitHub Actions | format、test、race、vet、Go 標準 toolchain と内製検査による static analysis、dependency inventory、secret boundary scan、executable fixture、mutation、競合、fault、browser、setup、release の check、Actions の commit SHA pin、minimum permissions、timeout、required checks。 | required check 未設定、Actions tag pin、permissions 過大、timeout 欠落、外部解析 tool の導入、skip success、CI 外の口頭説明だけで合格にする。 |
| governance / recovery | LICENSE、SECURITY.md、CONTRIBUTING.md、CODEOWNERS、CHANGELOG、stale lock、状態破損、容量不足、credential rotation、rollback の復旧手順。 | governance file 未作成、復旧手順未記載、復旧手順が実装契約または fixture へ接続しない、document drift が 1 件以上ある。 |

Phase 13 closure record set は、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の 18 record を使用し、各 record の `scope` に `phase-13-implementation-alignment-quality` を含める。Phase 13 では以下の集計値を同じ closure record set 内に記録する。

| 集計値 | 完了値 | 未完了条件 |
|--------|--------|------------|
| `phase13_owner_file_violation_count` | `0` | 5 ファイル固定、file role、package 名、import path、root `.go` 残存、サブディレクトリ禁止に違反する。 |
| `phase13_dummy_or_empty_file_count` | `0` | 空 file、dummy success、no-op success、placeholder、未実装を成功にする実装がある。 |
| `phase13_direct_state_mutation_count` | `0` | API、runner、MCP、archive が statefile owner adapter を経由せず状態 file を変更する。 |
| `phase13_contract_mismatch_count` | `0` | API、Admin、SDK、UI、MCP、setup stdout、credential 初期化の契約差分がある。 |
| `phase13_state_safety_open_count` | `0` | lock、atomic write、fsync、rename、owner / mode、migration、recovery、stale lock の未解消項目がある。 |
| `phase13_jsonl_corruption_open_count` | `0` | JSON Lines 破損検出、隔離、復旧対象登録、黙殺禁止、append lock / flush / fsync に未解消項目がある。 |
| `phase13_security_kdf_open_count` | `0` | 外部 KDF 依存、旧 SHA-256 反復 KDF 残存、`pbkdf2_hmac_sha256_v1` fixture 未接続がある。 |
| `phase13_token_arg_open_count` | `0` | Admin または MCP が token を command argv から受け取る、token file / stdin の検証がない。 |
| `phase13_archive_commitstatus_ownership_open_count` | `0` | archive または commitstatus 責務が呼び出し元 owner に重複残存する。 |
| `phase13_queue_recovery_open_count` | `0` | queue 状態機械、finalizer、active recovery、at-least-once、backup / restore transaction の残件がある。 |
| `phase13_mcp_unimplemented_success_count` | `0` | MCP の未実装機能が success を返す、または no-op success が残る。 |
| `phase13_external_boundary_open_count` | `0` | HTTP、Webhook、SSH、systemd、state directory の外部境界 hardening 残件がある。 |
| `phase13_required_log_write_ignore_count` | `0` | 必須 audit、access、config log の書込み失敗を無視する経路がある。 |
| `phase13_fixture_execution_gap_count` | `0` | production entrypoint 実行、actual / expected 比較、negative control の未接続がある。 |
| `phase13_mutation_survived_count` | `0` | 適用可能な production code mutation が survived である。 |
| `phase13_race_or_concurrency_open_count` | `0` | race、multi-process 更新、interleaving、goroutine / listener / timer / worker leak の未判定がある。 |
| `phase13_fault_injection_open_count` | `0` | disk full、permission denied、short write、fsync failure、rename failure、process kill の未実行または未接続がある。 |
| `phase13_e2e_open_count` | `0` | setup、install-api、update、rollback、Release、実インストール E2E、Admin / UI / runner / MCP integration の未接続がある。 |
| `phase13_ci_required_check_open_count` | `0` | Phase 13 required check が未作成、未実行、required 化未証跡、skip success 扱いである。 |
| `phase13_action_pin_open_count` | `0` | GitHub Actions が commit SHA pin でない、minimum permissions または timeout が欠落する。 |
| `phase13_release_evidence_open_count` | `0` | Git tag、GitHub Release、SHA256SUMS、signature、SBOM、再現ビルド証跡に未解消項目がある。 |
| `phase13_recovery_procedure_open_count` | `0` | stale lock、状態破損、容量不足、credential rotation、rollback の復旧手順が未接続である。 |
| `phase13_document_drift_open_count` | `0` | 旧 path、重複仕様、実装済み表記、未作成 path、責務正本参照に drift がある。 |
| `final_open_item_count` | `0` | 上記集計値または 18 record の `open_items` に残件がある。 |

Phase 13 の正式 fixture root は `testdata/phase13/implementation-alignment-quality/` とする。同 root の実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 13 target path 所在](../DOCUMENT_INDEX.md#phase-13-target-paths) を参照する。

Phase 13 の正式 fixture root は以下の file を必須とする。下表の必須 file が 1 件でも欠ける場合、`phase13_fixture_execution_gap_count` に 1 件以上を計上し、Phase 13 を完了扱いにしてはならない。

| path | 内容 | 完了条件 |
|------|------|----------|
| `manifest.json` | Phase 13 evidence package の識別子、対象 owner、実行順序、required check、closure record set の所在。 | `name` が `phase-13-implementation-alignment-quality`、`scope` が同値、`owners` が Phase 13 対象 owner をすべて含む。 |
| `input/scope.json` | 実行対象の owner、file pattern、禁止 pattern、required check、source audit scope。 | `include`、`exclude`、`owner_packages`、`state_files`、`contract_sources`、`release_artifacts` を持つ。 |
| `input/owner_inventory.json` | `components/<owner>/` の 5 file 固定、package 名、import path、file role、旧 root path の棚卸し入力。 | 全 Go owner を 1 回だけ列挙し、同一 owner の重複、空 owner、未作成 path の実在扱いがない。 |
| `input/contract_inventory.json` | API route、Admin command、SDK method、UI operation、MCP tool、setup stdout、credential 初期化の照合入力。 | 各 entry は `contract_id`、`owner`、`contract_kind`、`method`、`path`、`query_schema_ref`、`request_schema_ref`、`response_schema_ref`、`error_schema_ref`、`auth`、`status_codes`、`state_effects`、`client_bindings`、`source_anchor` を持つ。 |
| `input/state_inventory.json` | statefile 経由に集約する state file、JSON Lines log、migration、recovery、直接状態更新禁止の棚卸し入力。 | 直接 open / truncate / append / rename の検査対象 file pattern と許可 owner `statefile` を固定する。 |
| `input/security_inventory.json` | password hash、token、secret mask、signature、file safety、credential rotation、required log の棚卸し入力。 | secret 値そのものを含めず、fixture secret は deterministic placeholder と hash / mask 判定だけを持つ。 |
| `input/faults.json` | disk full、permission denied、short write、fsync failure、rename failure、process kill、network timeout、DNS rebinding の障害注入入力。 | 各 fault は `fault_id`、`target_owner`、`trigger`、`expected_error`、`must_preserve_state` を持つ。 |
| `expected/effects.json` | Phase 13 実行後の状態差分、log、artifact、counter、error、recovery の期待結果。 | 全 `phase13_*` counter と `final_open_item_count` の期待値を持つ。 |
| `expected/counters.json` | closure counter の期待値だけを機械照合する正規化 JSON。 | 全 counter の値が `0` であり、対象外は `not_applicable_reason_anchor` を併記する。 |
| `records/closure.jsonl` | 18 record closure set と Phase 13 counter の実行記録。 | 1 行 1 record、UTF-8、LF 終端、JSON object、`scope` は `phase-13-implementation-alignment-quality`。 |
| `records/mutation.jsonl` | mutation selection、applied mutation、killed / survived、対象外理由。 | 適用可能 mutation の `survived` が `0`。 |
| `records/race.jsonl` | race detector または代替 interleaving、multi-process state update、queue transition の実行記録。 | 未判定 owner がない。 |
| `records/fault.jsonl` | fault injection の実行記録。 | `input/faults.json` の全 `fault_id` が 1 回以上実行され、期待 error と state preservation が一致する。 |
| `records/e2e.jsonl` | Admin to API、UI to SDK to API、runner to statefile、MCP to statefile、setup、release の E2E 記録。 | stub だけの成功を E2E と扱わない。 |
| `records/release.jsonl` | Git tag、GitHub Release、SHA256SUMS、signature、SBOM、reproducible build evidence の記録。 | release 証跡不足がない。 |

`manifest.json` は以下の key を必須とする。未定義 key を実装側の自由解釈にしてはならない。

| key | type | 値 |
|-----|------|----|
| `name` | string | `phase-13-implementation-alignment-quality` 固定。 |
| `scope` | string | `phase-13-implementation-alignment-quality` 固定。 |
| `owners` | array[string] | `builder`、`runner`、`api`、`admin`、`setup`、`release`、`statefile`、`security`、`archive`、`commitstatus`、`mcp`、`sdk`、`ui` を含む。 |
| `entrypoints` | array[string] | production entrypoint として実行する binary / browser / fixture harness 名。 |
| `required_checks` | array[string] | Phase 13 required check name をすべて含む。 |
| `closure_records` | array[string] | `records/*.jsonl` の相対 path。 |
| `negative_controls` | array[string] | 失敗しなければならない検査の id。 |
| `source_anchors` | array[string] | Phase 13 対象の責務正本 anchor。 |

`entrypoints` は下表の値だけを許可する。下表にない entrypoint を Phase 13 完了証跡へ追加する場合は、先に [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) と該当 owner 詳細本文を改訂する。

| entrypoint | owner / 責務 | 実行境界 |
|------------|--------------|----------|
| `adlaire-ci-build` | `builder` | [`main.go`](../../main.go) の dispatch を経由して builder production entrypoint を実行する。 |
| `adlaire-ci-runner` | `runner` | [`main.go`](../../main.go) の dispatch を経由して runner production entrypoint を実行する。 |
| `adlaire-ci-api` | `api` | [`main.go`](../../main.go) の dispatch を経由して API production entrypoint を実行する。 |
| `adlaire-ci-admin` | `admin` | [`main.go`](../../main.go) の dispatch を経由して Admin CLI production entrypoint を実行する。 |
| `adlaire-ci-setup` | `setup` | [`main.go`](../../main.go) の dispatch を経由して setup production entrypoint を実行する。 |
| `adlaire-ci-release` | `release` | [`main.go`](../../main.go) の dispatch を経由して release production entrypoint を実行する。 |
| `adlaire-ci-mcp` | `mcp` | [`main.go`](../../main.go) の dispatch を経由して MCP production entrypoint を実行する。 |
| `admin-ui-browser` | `ui` / `sdk` / `api` | browser runtime で [`admin/index.html`](../../admin/index.html) と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) を読み込み、UI to SDK to API fake の実行結果を照合する。 |
| `phase13-fixture-harness` | fixture 証跡責務 | `manifest.json`、`input/*.json`、`expected/*.json`、`records/*.jsonl` を読み込み、production entrypoint 実行、actual / expected 比較、negative control、counter 集計を実行する。 |

`required_checks` は [Phase 13 required check 固定表](#phase-13-required-checks) の `required check name` をすべて 1 回だけ含める。重複、欠落、表にない check name、skip success を許可する check name は `phase13_ci_required_check_open_count` に計上する。

`source_anchors` は本節 [`Phase 13 実装整合・品質改善証跡`](#phase-13-implementation-alignment-quality-evidence)、[Phase 13 required check 固定表](#phase-13-required-checks)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](../DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 13 target path 所在](../DOCUMENT_INDEX.md#phase-13-target-paths)、[`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3)、[`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、[`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、および Phase 13 対象 owner 詳細本文の各 `phase-13-*-alignment-contract` anchor を含める。anchor が存在しない、または owner だけが推測で補われる場合は `phase13_document_drift_open_count` に計上する。

`expected/counters.json` は上記 closure counter 表の全 `phase13_*` counter と `final_open_item_count` を 1 回だけ持つ。追加 counter、欠落 counter、`expected/effects.json` と異なる counter 値、`records/closure.jsonl` の集計と一致しない counter 値、対象外理由 anchor のない対象外 counter を禁止する。

Phase 13 fixture は negative control を必須とする。negative control は、少なくとも contract mismatch、direct state mutation、token argv、JSON Lines corruption、mutation survivor、race trigger、fault injection failure、GitHub Actions unpinned を 1 件ずつ含める。negative control が成功扱いになる場合、該当 checker 自体を未完成として `phase13_fixture_execution_gap_count` に計上する。

Phase 13 の document drift 判定では、Phase 1〜Phase 10 の過去実装検証証跡として旧 root artifact を説明する履歴本文を、現行実装 path drift として数えない。drift として数える対象は、現在状態、実装着手可否、標準配置、Phase 13 target path、owner package inventory、API / SDK / UI / MCP / setup / release の現行契約が旧 root artifact を実在または標準配置として扱う記載だけとする。履歴本文を残す場合も、現行実装の正本は [`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3) と [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 13 target path 所在](../DOCUMENT_INDEX.md#phase-13-target-paths) であることを Phase 13 closure record に記録する。

<a id="phase-13-required-checks"></a>
Phase 13 の CI required check は以下とする。GitHub workflow は YAML 禁止の対象外であるが、workflow から実行する Adlaire CI 入出力、fixture manifest、expected、state、設定形式は JSON 契約に従う。

| required check name | 必須実行 | 完了条件 |
|---------------------|----------|----------|
| `phase13-go-format` | 全 Go file に `gofmt -l`。 | 出力空、owner package 5 ファイル固定、旧 root `.go` 残存 0。 |
| `phase13-go-test` | Go stable toolchain で `go test ./... -count=1`。 | exit code `0`、skip / 未実行 record なし。 |
| `phase13-go-race` | Go stable toolchain で `go test -race ./... -count=1`。 | exit code `0`、race / concurrency open `0`。 |
| `phase13-go-vet` | Go stable toolchain で `go vet ./...`。 | exit code `0`。 |
| `phase13-stdlib-static-analysis` | Go 標準 toolchain と本リポジトリ内の内製検査だけで静的解析を実行する。 | exit code `0`、外部解析 tool 未使用、owner package 5 ファイル固定、unreachable fixture / obsolete path open `0`。 |
| `phase13-dependency-inventory` | `go.mod`、Go import、JavaScript / HTML artifact、workflow を検査し、production 外部依存が 0 であることを確認する。 | exit code `0`、`go.mod` の外部 `require`、`golang.org/x/*`、third party import、npm / bundler / polyfill 依存 open `0`。 |
| `phase13-secret-boundary-scan` | Go 標準 toolchain と内製検査だけで token argv、secret 出力、fixture / log / release notes secret 漏えい境界を検査する。 | exit code `0`、secret / file safety / command execution finding open `0`、外部 secret scan service 未使用。 |
| `phase13-deno-check-sdk` | Deno stable runtime で `deno check admin/adlaire-ci-sdk.js`。 | exit code `0`。Node.js 代替禁止。 |
| `phase13-owner-shape` | owner package inventory、5 ファイル固定、旧 root path、空 file、dummy / no-op success の検査。 | `phase13_owner_file_violation_count=0`、`phase13_dummy_or_empty_file_count=0`、`phase13_document_drift_open_count=0`。 |
| `phase13-contract-parity` | API、Admin、SDK、UI、MCP、setup stdout、credential 初期化の契約自動照合。 | `phase13_contract_mismatch_count=0`、`phase13_mcp_unimplemented_success_count=0`。 |
| `phase13-state-safety` | statefile 直接更新禁止、process lock、atomic write、JSON Lines safety、migration、recovery の検査。 | `phase13_direct_state_mutation_count=0`、`phase13_state_safety_open_count=0`、`phase13_jsonl_corruption_open_count=0`。 |
| `phase13-security-boundary` | Go 標準ライブラリ内製 KDF、token file / stdin、secret mask、file safety、required log write の検査。 | `phase13_security_kdf_open_count=0`、`phase13_token_arg_open_count=0`、`phase13_required_log_write_ignore_count=0`。 |
| `phase13-archive-commitstatus` | archive / commitstatus の owner 集約、呼び出し元 duplicate 排除、retry、rate limit、timeout の検査。 | `phase13_archive_commitstatus_ownership_open_count=0`。 |
| `phase13-runner-recovery` | queue 状態機械、finalizer、active recovery、at-least-once、backup / restore transaction の検査。 | `phase13_queue_recovery_open_count=0`、`phase13_race_or_concurrency_open_count=0`。 |
| `phase13-mcp-real-behavior` | MCP `resendWebhook`、`subscribe`、`unsubscribe`、`sampling` の実動作、未実装 error、statefile 接続の検査。 | `phase13_mcp_unimplemented_success_count=0`、`phase13_direct_state_mutation_count=0`。 |
| `phase13-external-boundary` | API / MCP HTTP lifecycle、Webhook SSRF、SSH strict、systemd 最小権限、state directory 権限の検査。 | `phase13_external_boundary_open_count=0`、`phase13_required_log_write_ignore_count=0`。 |
| `phase13-executable-fixture` | production entrypoint 実行型 fixture harness。 | `phase13_fixture_execution_gap_count=0`。 |
| `phase13-mutation` | production code mutation。 | `phase13_mutation_survived_count=0`。 |
| `phase13-concurrency` | multi-process state update、queue transition、race trigger。 | `phase13_race_or_concurrency_open_count=0`。 |
| `phase13-fault-injection` | disk full、permission denied、short write、fsync failure、rename failure、process kill。 | `phase13_fault_injection_open_count=0`。 |
| `phase13-browser` | UI to SDK to API browser fixture。 | UI が API response にない値を合成せず、secret を表示しない。 |
| `phase13-integration-e2e` | Admin to API、UI to SDK to API、runner to statefile、MCP to statefile の統合 E2E。 | `phase13_contract_mismatch_count=0`、`phase13_direct_state_mutation_count=0`、`phase13_e2e_open_count=0`。 |
| `phase13-setup-e2e` | setup、install-api、update、rollback、実インストール E2E。 | `phase13_e2e_open_count=0`。 |
| `phase13-release-e2e` | Git tag、GitHub Release、SHA256SUMS、signature、SBOM、再現ビルド証跡。 | `phase13_release_evidence_open_count=0`。 |
| `phase13-actions-pinning` | GitHub Actions commit SHA pin、minimum permissions、timeout、required checks の検査。 | `phase13_action_pin_open_count=0`、`phase13_ci_required_check_open_count=0`。 |
| `phase13-recovery-procedure` | stale lock、状態破損、容量不足、credential rotation、rollback の復旧手順と実装契約接続の検査。 | `phase13_recovery_procedure_open_count=0`。 |
| `phase13-document-drift` | 旧 path、重複仕様、実装済み表記、未作成 path、責務正本参照、ROADMAP / DOCUMENT_INDEX の drift 検査。 | `phase13_document_drift_open_count=0`、`final_open_item_count=0`。 |

<a id="phase-14-obsidian-vault-integration-evidence"></a>
**Phase 14 Obsidian Vault 連携証跡：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 14 の Obsidian local vault 入力、wikilink / embed / tag / asset 正規化、YAML frontmatter 拒否、builder handoff、closure record の証跡だけを固定する。Phase 14 の現在状態と実装割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 14 Obsidian Vault 連携参照](../DETAIL_INDEX.md#phase-14-obsidian-vault-integration-entry)、実装契約は [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian Vault 連携契約](obsidian.md#obsidian-phase14-vault-integration-contract) を参照する。

Phase 14 の正式 fixture root は `testdata/phase14/obsidian-vault-integration/` とする。同 root は Phase 14 実装 PR で実在化し、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 14 target path 所在](../DOCUMENT_INDEX.md#phase-14-target-paths) で `実在` として扱う。

| path | 内容 | 完了条件 |
|------|------|----------|
| `manifest.json` | Phase 14 evidence package の識別子、対象 owner、builder collaborator、closure record set の所在。 | `name` が `phase-14-obsidian-vault-integration`、`owners` が `obsidian`、`collaborators` が `builder` と `security` を含む。 |
| `input/options.json` | `adlaire-ci-build --input-mode obsidian-vault` の CLI option、strict、output mode、vault root、entry note、include / exclude。 | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 CLI 契約](obsidian.md#obsidian-phase14-cli-contract) の全 option を持ち、未知 option がない。 |
| `input/vault_tree.json` | vault root 配下の directory / file / symlink / hardlink / hidden path / `.obsidian/` の棚卸し入力。 | vault 外 path、symlink、device、`.obsidian/` 解釈、case-insensitive 解決を negative control として含む。 |
| `input/notes/` | Obsidian Markdown note 入力。 | wikilink、alias、heading link、asset embed、tag、code fence、inline code、YAML frontmatter 拒否 case を含む。 |
| `input/assets/` | 画像、PDF、その他許可 asset 入力。 | 許可拡張子、禁止拡張子、同名 asset、vault 外参照を含む。 |
| `expected/normalized.json` | note graph、resolved link、unresolved link、embed、tag、asset copy、slug、diagnostic の期待値。 | sort order、path、line / column、error code が固定される。 |
| `expected/output_tree.json` | builder handoff 後の中間 Markdown、asset、report、public output 期待値。 | 入力 vault への write が 0 件、handoff root 外 write が 0 件。 |
| `expected/counters.json` | Phase 14 closure counter の期待値。 | 全 counter が `0`、対象外理由は anchor 付き。 |
| `records/closure.jsonl` | Phase 14 closure record set。 | 1 行 1 record、UTF-8、LF 終端、`scope` は `phase-14-obsidian-vault-integration`。 |

`expected/normalized.json` は [`docs/details/obsidian.md` 詳細本文責務 Phase 14 output schema 契約](obsidian.md#obsidian-phase14-output-schema-contract) の `obsidian_map.json` root key、note object、link object、asset object、diagnostic object と同じ key 順、型、sort order を持つ。fixture harness は actual `obsidian_map.json` を JSON parse 後に比較し、unknown key、欠落 key、array sort 差分、absolute path、free-form diagnostic message、`null` を失敗にする。`expected/output_tree.json` は `created_paths`、`updated_paths`、`deleted_paths`、`unchanged_paths`、`forbidden_writes`、`stdout_json_keys`、`stderr_lines`、`normalized_note_sha256s`、`asset_sha256s` の root key だけを持ち、input vault 配下 write と handoff root 外 write を `forbidden_writes` で明示する。`normalized_note_sha256s` は normalized path の UTF-8 byte 昇順 object array とし、各 object は `path`、`sha256`、`size` を持つ。`asset_sha256s` は asset normalized path の UTF-8 byte 昇順 object array とし、各 object は `path`、`sha256`、`size`、`media_type` を持つ。fixture harness は note wikilink、asset image embed、asset PDF embed、tag 非変換、末尾 LF、relative link destination を byte 単位で比較し、`expected/normalized.json` の schema 比較だけで Phase 14 を合格にしてはならない。

| 集計値 | 完了値 | 未完了条件 |
|--------|--------|------------|
| `phase14_obsidian_vault_boundary_open_count` | `0` | vault root containment、symlink / hardlink / device 拒否、`.obsidian/` 非解釈、vault 外 path 拒否に残件がある。 |
| `phase14_obsidian_wikilink_open_count` | `0` | wikilink、alias、heading link、未解決 link、duplicate basename、case mismatch の deterministic 処理に残件がある。 |
| `phase14_obsidian_yaml_rejection_open_count` | `0` | YAML frontmatter、YAML metadata、YAML block を解釈または黙認する経路がある。 |
| `phase14_obsidian_asset_open_count` | `0` | asset embed、asset copy、禁止拡張子、vault 外 asset、asset digest、asset path safety に残件がある。 |
| `phase14_obsidian_builder_handoff_open_count` | `0` | builder handoff、stdout / stderr / report、公開出力維持、入力 vault no-write に残件がある。 |
| `phase14_fixture_execution_gap_count` | `0` | production entrypoint 実行、actual / expected 比較、negative control の未接続がある。 |
| `final_open_item_count` | `0` | 上記集計値または closure record の `open_items` に残件がある。 |

Phase 14 required check は以下に固定する。Phase 14 実装 PR は、下表の check 名、対象、完了条件を同一 Pull Request 本文へ記録する。

| check 名 | 対象 | 完了条件 |
|----------|------|----------|
| `phase14-go-format` | `components/obsidian/` と Obsidian 連携に触れた Go file | `gofmt -l` の差分が 0。 |
| `phase14-go-test` | Obsidian 連携に関係する production entrypoint と owner package | `go test ./...` が成功し、skip がある場合は対象外 anchor を記録する。 |
| `phase14-obsidian-vault-fixture` | `testdata/phase14/obsidian-vault-integration/` | production entrypoint を実行し、actual normalized graph、output tree、closure counters が expected と一致する。 |
| `phase14-builder-handoff` | Obsidian normalized root から builder への handoff | 入力 vault write が 0、handoff root 外 write が 0、builder stdout / stderr 契約差分が 0。 |
| `phase14-document-drift` | [`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md)、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md)、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md)、[`docs/details/obsidian.md`](obsidian.md)、[`docs/details/builder.md`](builder.md)、[`docs/details/security.md`](security.md)、[`docs/details/fixture.md`](fixture.md) | Phase 14 の状態、path、anchor、fixture root、builder handoff、security collaborator、未作成表記の drift が 0。 |

Phase 14 fixture は negative control を必須とする。negative control は、YAML frontmatter、vault escape、symlink、hardlink、invalid filter、invalid filter pattern、filter type mismatch、malformed wikilink、display text invalid、unresolved wikilink strict、duplicate tag strict、unused asset strict、duplicate basename、case mismatch、heading slug mismatch、relative link escape、unsupported `.canvas`、Dataview block、Templater block、external fetch attempt、note embed の禁止 case を 1 件以上含める。negative control が成功扱いになる場合、該当 checker 自体を未完成として `phase14_fixture_execution_gap_count` に計上する。

<a id="phase-15-obsidian-local-sync-evidence"></a>
**Phase 15 Obsidian local vault 同期証跡：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 15 の sync plan、apply、rollback、conflict、tombstone、atomic write、Obsidian Sync service 非依存、closure record の証跡だけを固定する。Phase 15 の現在状態と実装割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 15 Obsidian local vault 同期参照](../DETAIL_INDEX.md#phase-15-obsidian-local-sync-entry)、実装契約は [`docs/details/obsidian.md` 詳細本文責務 Phase 15 Obsidian local vault 同期契約](obsidian.md#obsidian-phase15-local-sync-contract) を参照する。

Phase 15 の正式 fixture root は [`testdata/phase15/obsidian-local-sync/`](../../testdata/phase15/obsidian-local-sync/) とする。同 root の所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 15 target path 所在](../DOCUMENT_INDEX.md#phase-15-target-paths) で `実在` として確認する。

| path | 内容 | 完了条件 |
|------|------|----------|
| `manifest.json` | Phase 15 evidence package の識別子、対象 owner、statefile / security / release / setup collaborator、closure record set の所在。 | `name` が `phase-15-obsidian-local-sync`、`owners` が `obsidian`、`collaborators` が `statefile`、`security`、`release`、`setup` を含む。 |
| `input/options.json` | `adlaire-ci-obsidian sync plan` / `sync apply` / `sync rollback` の CLI option。 | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 CLI 契約](obsidian.md#obsidian-phase15-cli-contract) の全 option を持ち、未知 option がない。invalid command、unknown option、duplicate option、必須 option 欠落、invalid `--direction`、invalid `--delete-policy`、invalid `--open-uri` を negative control として持つ。 |
| `input/sync_state.json` | 前回同期 state、digest、tombstone、conflict、rollback metadata。 | mtime だけの同一性判定を持たず、content digest と normalized path を持つ。 |
| `input/project_tree.json` | project 側 tree snapshot。 | create、modify、delete、rename、conflict、permission denied、read-only file を含む。 |
| `input/vault_tree.json` | vault 側 tree snapshot。 | create、modify、delete、rename、clock skew、partial write、official sync metadata 風 file を含む。 |
| `expected/plan.json` | project root、vault root、`sync_state.json` を変更しない plan の expected operations、conflict、tombstone、rollback precondition、conflict dir、tombstone dir。 | plan hash、operation order、`conflict_dir`、`tombstone_dir`、`--plan-file` 以外の created / updated / deleted path が 0 件で固定される。 |
| `expected/apply_state.json` | apply 後 state、staging cleanup、tombstone、conflict file、rollback record。 | atomicity、fsync、rollback、unchanged path が固定される。 |
| `expected/distribution.json` | Obsidian CLI 配布連携の expected asset、checksum、setup mode、version stdout。 | [`docs/details/release.md` 詳細本文責務 Phase 15 Obsidian Release 配布拡張契約](release.md#phase-15-obsidian-release-extension-contract) と [`docs/details/setup.md` 詳細本文責務 Phase 15 Obsidian CLI 導入手順](setup.md#phase-15-obsidian-setup-contract) の asset count、checksum 行数、asset 名、setup mode、version stdout と一致する。 |
| `expected/counters.json` | Phase 15 closure counter の期待値。 | 全 counter が `0`、対象外理由は anchor 付き。 |
| `records/closure.jsonl` | Phase 15 closure record set。 | 1 行 1 record、UTF-8、LF 終端、`scope` は `phase-15-obsidian-local-sync`。 |

`expected/plan.json` と `expected/apply_state.json` は [`docs/details/obsidian.md` 詳細本文責務 Phase 15 schema 契約](obsidian.md#obsidian-phase15-schema-contract) の `plan.json`、`sync_state.json`、conflict、tombstone、rollback record と同じ key 順、型、sort order を持つ。fixture harness は actual plan hash を canonical JSON byte から再計算し、`--plan-hash` 照合、`sync plan` 実行時の project root / vault root / `sync_state.json` 不変、apply 後 state digest、rollback record、stdout JSON key 順、stderr exact line を比較する。stdout の `plan_file` と `rollback_file` は state dir 相対 path として比較し、absolute path、home directory、temporary directory、host 固有 path が含まれる場合は失敗にする。`expected/plan.json` は `conflict_dir_default`、`tombstone_dir_default`、`conflict_dir`、`tombstone_dir` を含み、`sync apply` が CLI option ではなく plan file 内の `conflict_dir` と `tombstone_dir` だけを使用することを fixture で確認する。`expected/apply_state.json` は `delete_policy_default`、`open_uri_default`、`rollback_file_default`、`apply_id`、`conflict_ids`、`tombstone_ids`、`backup_paths` を含み、[`docs/details/obsidian.md` 詳細本文責務 Phase 15 schema 契約](obsidian.md#obsidian-phase15-schema-contract) の算出式と一致しなければならない。`expected/distribution.json` は `binary_asset`、`release_asset_count`、`checksum_line_count`、`setup_acceptance_asset`、`setup_mode`、`version_stdout`、`release_extension_anchor` の root key だけを持ち、[`docs/details/release.md` 詳細本文責務 Phase 15 Obsidian Release 配布拡張契約](release.md#phase-15-obsidian-release-extension-contract) と [`docs/details/setup.md` 詳細本文責務 Phase 15 Obsidian CLI 導入手順](setup.md#phase-15-obsidian-setup-contract) に一致しなければならない。

| 集計値 | 完了値 | 未完了条件 |
|--------|--------|------------|
| `phase15_sync_plan_mismatch_count` | `0` | plan hash、operation order、digest、`--plan-file` 以外の write 0 件、expected plan に差分がある。 |
| `phase15_sync_apply_atomicity_open_count` | `0` | lock、staging、atomic rename、fsync、partial write、process kill、state update order に残件がある。 |
| `phase15_sync_conflict_open_count` | `0` | both-side edit、rename collision、delete vs edit、clock skew、read-only conflict の処理に残件がある。 |
| `phase15_sync_tombstone_open_count` | `0` | delete policy、tombstone record、hard delete 禁止既定、restore path に残件がある。 |
| `phase15_sync_rollback_open_count` | `0` | rollback precondition、rollback record、staging cleanup、failed apply 後復旧に残件がある。 |
| `phase15_sync_service_dependency_open_count` | `0` | Obsidian Sync service、Obsidian cloud、remote vault API、plugin runtime、external watcher 依存が 1 件以上ある。 |
| `phase15_distribution_open_count` | `0` | `adlaire-ci-obsidian-linux-amd64` の Release asset、checksum、setup 取得対象、version 出力、fixture expected に残件がある。 |
| `phase15_fixture_execution_gap_count` | `0` | production entrypoint 実行、actual / expected 比較、negative control の未接続がある。 |
| `final_open_item_count` | `0` | 上記集計値または closure record の `open_items` に残件がある。 |

Phase 15 required check は以下に固定する。Phase 15 実装 PR は、下表の check 名、対象、完了条件を同一 Pull Request 本文へ記録する。

| check 名 | 対象 | 完了条件 |
|----------|------|----------|
| `phase15-go-format` | `components/obsidian/` と Obsidian 同期に触れた Go file | `gofmt -l` の差分が 0。 |
| `phase15-go-test` | Obsidian 同期に関係する production entrypoint と owner package | `go test ./...` が成功し、skip がある場合は対象外 anchor を記録する。 |
| `phase15-obsidian-sync-fixture` | `testdata/phase15/obsidian-local-sync/` | production entrypoint を実行し、actual plan、apply state、rollback state、closure counters が expected と一致する。 |
| `phase15-sync-atomicity` | lock、staging、fsync、atomic rename、state update、rollback record | partial write、process kill、rename failure、fsync failure、permission denied の failure injection が expected と一致する。 |
| `phase15-release-setup-distribution` | `adlaire-ci-obsidian-linux-amd64` の Release / setup 連携 | `expected/distribution.json` と Release / setup 詳細本文の asset 名、asset count、checksum line count、setup mode、version stdout が一致し、`phase15_distribution_open_count=0`。 |
| `phase15-document-drift` | [`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md)、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md)、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md)、[`docs/details/obsidian.md`](obsidian.md)、[`docs/details/release.md`](release.md)、[`docs/details/setup.md`](setup.md)、[`docs/details/fixture.md`](fixture.md) | Phase 15 の状態、path、anchor、fixture root、配布連携、未作成表記の drift が 0。 |

Phase 15 fixture は negative control を必須とする。negative control は、invalid command、unknown option、duplicate option、必須 option 欠落、invalid direction、invalid delete policy、invalid open-uri、apply 時の conflict / tombstone dir CLI 再指定、plan hash mismatch、both-side edit、opposite-side edit、delete vs edit、rename collision、clock skew、read-only file、permission denied、invalid plan file path、invalid rollback file path、invalid conflict / tombstone dir、plan file 内 conflict / tombstone dir path escape、default rollback file collision、stdout absolute path leak、rollback record write failure、rollback conflict、backup path escape、partial write、process kill、official Obsidian Sync service / cloud endpoint attempt、Obsidian URI 成功依存、vault outside write、Release asset 未追加、checksum 行不足、setup 取得対象不足を 1 件以上含める。negative control が成功扱いになる場合、該当 checker 自体を未完成として `phase15_fixture_execution_gap_count` または `phase15_distribution_open_count` に計上する。

<a id="phase-16-quality-evidence-closure-evidence"></a>
**Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務では、Phase 16 の宣言型証跡の実行型証跡化、`ALIGN-*` 追加検証候補分類、ignored error 分類、skip / 未実行 closure、determinism、filesystem durability parity、workflow hardening、Docker 検証、release rehearsal、巨大 owner risk ledger の証跡だけを固定する。Phase 16 の現在状態と実装割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、対象入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](../DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry)、Phase 完了単位は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、テスト方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を参照する。

Phase 16 の正式 fixture root は `testdata/phase16/quality-evidence-closure/` とする。同 root は Phase 16 実装 PR で実在化し、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](../DOCUMENT_INDEX.md#phase-16-target-paths) で `実在` として扱う。仕様策定時点で root が未作成の場合は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](../DOCUMENT_INDEX.md#phase-16-target-paths) に `未作成` として記録し、Phase 16 を `実装済み` に遷移させてはならない。

| path | 内容 | 完了条件 |
|------|------|----------|
| `manifest.json` | Phase 16 evidence package の識別子、対象 owner、work unit、required check、closure record set の所在。 | `name` が `phase-16-quality-evidence-closure`、`scope` が同値、対象 owner が Phase 16 work unit に必要な owner をすべて含む。 |
| `input/source_coverage.json` | live source coverage set から生成した Phase 16 対象 source と検出語の棚卸し。 | `rg --files --hidden -g '!.git/**'` の対象結果から Phase 16 source coverage set を再導出し、各 source の検出語、sha256、inventory 接続を記録する。 |
| `input/evidence_inventory.json` | Phase 12 から Phase 15 の manifest、records、expected、workflow、PR evidence、skip、未実行、対象外理由、counter の棚卸し。 | 宣言型証跡、record 存在確認だけの証跡、actual / expected 比較未接続、mutation 実行未接続、fault / race / E2E / release 実行未接続を分類する。 |
| `input/align_inventory.json` | [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](#current-implementation-alignment-evidence) の `ALIGN-*` 追加検証候補一覧。 | 各 `ALIGN-*` に `phase16_action`、`owner_component`、`closure_ref`、`future_ref`、`not_applicable_ref` のいずれかを持ち、未分類を残さない。 |
| `input/error_inventory.json` | ignored error、required write、audit / access / config / state / JSON Lines write、panic、skip / 未実行の棚卸し。 | safe ignore、must handle、must fail、cleanup only、hash writer、test skip、not applicable を区別し、根拠 anchor を持つ。 |
| `input/determinism_inventory.json` | clock、sleep、timeout、entropy、parallel worker、HTTP lifecycle、external I/O の棚卸し。 | fake adapter、固定時刻、固定 timeout、body limit、client disconnect、partial write、retry backoff の実行証跡接続先を持つ。 |
| `input/filesystem_inventory.json` | state、config、audit、history、snapshot、archive、release、setup、Obsidian sync の file operation 棚卸し。 | atomic write、fsync、parent directory fsync、rename、symlink 非追従、stale lock、recovery、permission failure を分類する。 |
| `input/workflow_inventory.json` | Phase 12 から Phase 16 の GitHub workflow と required check の棚卸し。 | action SHA pin、minimum permissions、timeout、skip success 禁止、required check 名、Docker 検証手順を分類する。 |
| `input/large_owner_inventory.json` | 巨大 owner file、5 ファイル原則、内部責務区画、test artifact 接続の棚卸し。 | 6 ファイル目、subdirectory、空 file、dummy 実装を 0 とし、巨大 file は risk ledger と内部責務区画検査で閉じる。 |
| `input/negative_controls.json` | checker、harness、counter、diagnostic、source coverage、workflow、mutation、fault injection の負例入力。 | 失敗すべき fixture 変異を固定し、checker が各負例を失敗として検出する。 |
| `expected/counters.json` | Phase 16 closure counter の期待値。 | 本節の全 `phase16_*` counter と `final_open_item_count` が `0`。 |
| `expected/actions.json` | Phase 16 で実施する修正・対象外・将来計画維持の期待分類。 | `ALIGN-*`、skip、ignored error、workflow、durability、determinism、large owner の各 item が 1 件ずつ closure action を持つ。 |
| `records/closure.jsonl` | Phase 16 closure record set。 | 1 行 1 record、UTF-8、LF 終端、`scope` は `phase-16-quality-evidence-closure`。 |
| `records/execution.jsonl` | production entrypoint 実行、actual / expected 比較、fixture harness、Docker 検証の実行証跡。 | 実行 command、input、actual、expected、diff、exit code、stdout / stderr、runtime、skip 有無を記録する。 |
| `records/mutation.jsonl` | production code mutation 実行証跡。 | mutation target、before、after、expected failure、actual failure、killed / survived / invalid / equivalent を記録し、`survived=0`。 |
| `records/fault.jsonl` | disk full、permission denied、short write、fsync failure、rename failure、process kill、partial write、client disconnect の障害注入証跡。 | failure injection が実行され、期待 failure code、rollback、recovery、no silent success が一致する。 |
| `records/workflow.jsonl` | workflow hardening と required check 実行証跡。 | action SHA pin、permissions、timeout、required check、Docker fallback、Deno stable、release rehearsal、Phase 12 workflow 差分が閉じている。 |

Phase 16 の棚卸し対象 source は、`rg --files --hidden -g '!.git/**'` で列挙できる実在 path と、必須所在確認対象の不在 record を基準にする。実装者は同列挙結果から `AGENTS.md`、`README.md`、`main.go`、`go.mod`、`components/*/*.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html`、`.github/workflows/*.yml`、`.github/workflows/*.yaml`、`testdata/phase12/quality-gate-reconstruction/**`、`testdata/phase13/implementation-alignment-quality/**`、`testdata/phase14/obsidian-vault-integration/**`、`testdata/phase15/obsidian-local-sync/**`、`testdata/phase16/quality-evidence-closure/**`、`docs/SPEC.md`、`docs/ROADMAP.md`、`docs/DETAIL_INDEX.md`、`docs/DOCUMENT_INDEX.md`、`docs/DESIGN.md`、`docs/details/fixture.md`、`docs/details/*.md` を Phase 16 source coverage set として扱う。`go.sum` は必須所在確認対象とし、実在する場合は source coverage set に含め、実在しない場合は `go.sum` が存在しないことを `input/source_coverage.json` の `not_applicable` record として扱い、存在しない file を作成して補ってはならない。source coverage set は UTF-8 byte 昇順の一意 path とし、重複 pattern で同一 path が複数回一致しても record は 1 件だけ作成する。source coverage set に含まれる path を inventory から省略してはならない。対象外とする場合も、該当 inventory record に `classification=not_applicable`、`not_applicable_ref`、`status=not_applicable` を記録する。

Phase 16 の source coverage set では、`_ =`、ignored return、`panic(`、`t.Skip`、`Skip(`、`time.Now`、`time.Sleep`、`time.After`、timer、timeout、`rand.`、`crypto/rand`、`http.Client`、`http.Server`、`http.NewRequest`、response `Write`、`Flush`、`exec.Command`、`exec.CommandContext`、`systemctl`、`os.Rename`、`os.Remove`、`os.OpenFile`、`os.WriteFile`、`Write`、`Sync`、`Close`、`filepath.WalkDir`、symlink 判定、`go.mod` の `require` / `replace` directive、external module path、external import path、JSONL record、closure counter、workflow action、workflow permissions、workflow timeout、Deno check、Docker fallback、Release rehearsal、5 ファイル原則、空 file、dummy 実装、owner file line count を検出対象に含める。owner file line count は LF 区切りの物理行数で数え、末尾 LF のない最終行も 1 行として数える。検出対象が source coverage set に存在するのに対応 inventory record がない場合、`phase16_execution_evidence_gap_count` または該当する `phase16_*_open_count` を `0` にしてはならない。

Phase 16 source coverage detection registry は以下に固定する。checker は対象 source を UTF-8 text として読み、コメント、文字列、Markdown code block、HTML comment、workflow comment の区別をせず、下表の `match` を `match_mode` に従って検出する。`match_mode=literal_substring` は case-sensitive UTF-8 substring、`match_mode=line_contains` は 1 行内に表記したすべての literal が現れる場合、`match_mode=go_mod_directive` は `go.mod` の logical line が表記した directive で始まる場合、`match_mode=path_literal` は path または import / module path に表記した literal が現れる場合を意味する。UTF-8 として読めない source は `phase16_execution_evidence_gap_count` に計上する。検出対象語が documentation だけに現れる場合も、該当 inventory record で `not_applicable`、`future_plan`、または `reject_invalid_evidence` に分類し、未記録のまま完了扱いにしてはならない。

| detected term id | match_mode | match | source scope | counter routing |
|------------------|------------|-------|--------------|-----------------|
| `ignored_assignment` | `literal_substring` | `_ =` | all text source | `phase16_ignored_error_unclassified_count` |
| `ignored_return_phrase` | `literal_substring` | `ignored return` | all text source | `phase16_ignored_error_unclassified_count` |
| `panic_call` | `literal_substring` | `panic(` | Go source | `phase16_panic_unclassified_count` |
| `test_skip_call` | `literal_substring` | `t.Skip`、`Skip(` | Go source | `phase16_skip_open_count` |
| `clock_now` | `literal_substring` | `time.Now` | Go source | `phase16_determinism_open_count` |
| `sleep_or_after` | `literal_substring` | `time.Sleep`、`time.After` | Go source | `phase16_determinism_open_count` |
| `timer_or_timeout` | `literal_substring` | `timer`、`timeout`、`Timeout` | all text source | `phase16_determinism_open_count` |
| `entropy` | `literal_substring` | `rand.`、`crypto/rand` | Go source | `phase16_determinism_open_count` |
| `http_boundary` | `literal_substring` | `http.Client`、`http.Server`、`http.NewRequest`、`Flush` | Go source | `phase16_http_boundary_open_count` |
| `http_response_write` | `line_contains` | `response` + `Write` | Go source | `phase16_http_boundary_open_count` |
| `external_command` | `literal_substring` | `exec.Command`、`exec.CommandContext`、`systemctl` | Go source | `phase16_validation_portability_open_count` |
| `filesystem_write` | `literal_substring` | `os.Rename`、`os.Remove`、`os.OpenFile`、`os.WriteFile`、`Write`、`Sync`、`Close`、`filepath.WalkDir` | Go source | `phase16_filesystem_durability_open_count` |
| `symlink_boundary` | `literal_substring` | `symlink`、`EvalSymlinks`、`Lstat` | all text source | `phase16_filesystem_durability_open_count` |
| `go_mod_require` | `go_mod_directive` | `require` | `go.mod` only | `phase16_execution_evidence_gap_count` |
| `go_mod_replace` | `go_mod_directive` | `replace` | `go.mod` only | `phase16_execution_evidence_gap_count` |
| `external_module_path` | `path_literal` | `golang.org/`、`github.com/`、`npm`、`node_modules` | `go.mod`、Go import、JavaScript、document source | `phase16_execution_evidence_gap_count` |
| `jsonl_or_counter` | `literal_substring` | `JSONL record`、`closure counter`、`final_open_item_count` | document / fixture source | `phase16_declarative_evidence_open_count` |
| `workflow_hardening` | `literal_substring` | `uses:`、`permissions:`、`timeout-minutes`、`required check` | workflow / document source | `phase16_workflow_hardening_open_count` |
| `runtime_portability` | `literal_substring` | `Deno check`、`Docker fallback`、`go test`、`gofmt` | workflow / document source | `phase16_validation_portability_open_count` |
| `release_rehearsal` | `literal_substring` | `Release rehearsal`、`SHA256SUMS`、`SBOM`、`signature`、`reproducible build` | release / document source | `phase16_release_rehearsal_open_count` |
| `owner_shape` | `literal_substring` | `5 ファイル原則`、`空 file`、`dummy 実装`、`owner file line count` | Go owner / document source | `phase16_large_owner_risk_open_count` |

Phase 16 `manifest.json` は以下の key だけを持つ JSON object とする。未知 key、欠落 key、型違い、空配列、空文字列、重複値は `phase16_execution_evidence_gap_count` に計上する。

| key | 型 | 固定値 / 条件 |
|-----|----|---------------|
| `name` | string | `phase-16-quality-evidence-closure` 固定。 |
| `scope` | string | `phase-16-quality-evidence-closure` 固定。 |
| `schema_version` | integer | `1` 固定。 |
| `fixture_root` | string | `testdata/phase16/quality-evidence-closure/` 固定。 |
| `owner_components` | array[string] | `admin`、`api`、`archive`、`builder`、`commitstatus`、`mcp`、`obsidian`、`release`、`runner`、`sdk`、`security`、`setup`、`statefile`、`ui` を含む。順序は ASCII 昇順。 |
| `work_units` | array[string] | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 16 work unit](../DETAIL_INDEX.md#phase-16-quality-evidence-closure-entry) の 7 work unit をすべて含む。順序は work unit 表の順序。 |
| `input_files` | array[string] | 本節 path 表の `input/*.json` をすべて含む。 |
| `expected_files` | array[string] | `expected/counters.json`、`expected/actions.json` を含む。 |
| `record_files` | array[string] | `records/closure.jsonl`、`records/execution.jsonl`、`records/mutation.jsonl`、`records/fault.jsonl`、`records/workflow.jsonl` を含む。 |
| `required_checks` | array[string] | 本節 required check 表の check 名をすべて含む。 |
| `counters` | object | 本節 closure counter 表の全 key を持ち、値はすべて `0`。 |

<a id="phase-16-inventory-common-schema"></a>
**Phase 16 inventory 共通 schema：**

Phase 16 の `input/*.json` は、root object に `schema_version`、`scope`、`records` だけを持つ。未知 key、欠落 key、型違い、空配列、空文字列、重複 key は `phase16_execution_evidence_gap_count` に計上する。`schema_version` は `1`、`scope` は `phase-16-quality-evidence-closure`、`records` は 1 件以上の array とする。全 inventory record は下表の key を必須とし、各 inventory 固有 key は後続表で追加する。

| key | 型 | 固定条件 |
|-----|----|----------|
| `id` | string | `phase16.<work_unit_slug>.<owner_or_area>.<slug>` 形式。`work_unit_slug` は `phase16-` prefix を除いた work unit 名を `.` 区切りへ変換した値とする。lowercase、ASCII、dot notation。重複禁止。 |
| `work_unit` | string | Phase 16 の 7 work unit のいずれか。 |
| `owner_component` | string | owner component 名、または workflow / document / cross のような責務領域名。空文字禁止。 |
| `source_ref` | string | 責務名付き Markdown link で到達できる正本 anchor、実装 path、workflow path、または fixture path。裸の説明文だけを禁止する。 |
| `source_locator` | string | `symbol:<name>`、`function:<name>`、`file:<basename>`、`fixture:<name>`、`workflow:<job>`、`align:<id>`、`counter:<key>`、`record:<file>` のいずれか。行番号に依存する値を禁止する。 |
| `classification` | string | `execute`、`fix`、`not_applicable`、`future_plan`、`reject_invalid_evidence` のいずれか。 |
| `required_action` | string | `execute_and_close`、`fix_and_close`、`classify_not_applicable`、`keep_future_plan`、`reject_as_invalid_evidence` のいずれか。 |
| `closure_ref` | string/null | 完了証跡へ到達できる責務名付き Markdown link、または Phase 16 実装前は `null`。`classification=execute` または `fix` の完了時は null 禁止。 |
| `not_applicable_ref` | string/null | 対象外理由の正本 anchor。`classification=not_applicable` 以外では `null`。 |
| `future_ref` | string/null | 将来計画維持の到達先。`classification=future_plan` 以外では `null`。 |
| `status` | string | `open`、`closed`、`not_applicable`、`future_plan`、`rejected` のいずれか。 |
| `counter_key` | string | 本節 closure counter 表のいずれか。 |

`closure_ref`、`not_applicable_ref`、`future_ref` のうち、完了時に値を持てる key は 1 つだけとする。`status=open` の record が 1 件でもある場合、対応する `phase16_*` counter と `final_open_item_count` は `0` にしてはならない。

| inventory | 追加必須 key | 固定条件 |
|-----------|--------------|----------|
| `input/source_coverage.json` | `path`、`path_kind`、`source_exists`、`sha256`、`size_bytes`、`absence_reason`、`detected_terms`、`inventory_refs` | `path` は source coverage set の UTF-8 byte 昇順 path、`path_kind` は `go`、`module`、`javascript`、`html`、`workflow`、`fixture`、`document` のいずれか、`source_exists` は boolean、`source_exists=true` では `sha256` が file content digest、`size_bytes` が 0 以上の integer、`absence_reason=null`、`source_exists=false` では `sha256=null`、`size_bytes=null`、`absence_reason` が責務正本 anchor を含む string とする。`source_exists=false` は `go.sum` の必須不在確認だけに許可し、その他 path の不在は失敗とする。`detected_terms` は下記の source coverage detected term object の一意 array、`inventory_refs` は該当 inventory record `id` の ASCII 昇順 1 件以上の array。検出対象語がない source でも `detected_terms=[]` とし、`inventory_refs` には対象外または drift 確認 record を接続する。 |
| `input/evidence_inventory.json` | `evidence_kind`、`phase_source`、`execution_requirement` | `evidence_kind` は `counter_only`、`record_exists_only`、`pr_text_only`、`skip_or_unexecuted`、`mutation_unexecuted`、`fault_unexecuted`、`race_unexecuted`、`release_unrehearsed`、`executable` のいずれか。`phase_source` は `phase12`、`phase13`、`phase14`、`phase15` のいずれか。`execution_requirement` は production entrypoint、actual / expected 比較、mutation、fault、race、release rehearsal のいずれかを含む。 |
| `input/align_inventory.json` | `align_id`、`align_title`、`phase16_action` | `align_id` は `ALIGN-` と 2 桁以上の数字。`phase16_action` は `execute`、`future_plan`、`not_applicable` のいずれか。 |
| `input/error_inventory.json` | `error_kind`、`must_fail`、`safe_ignore_reason` | `error_kind` は `ignored_return`、`required_write_failure`、`cleanup_failure`、`panic`、`skip`、`http_write_failure`、`hash_writer`、`best_effort` のいずれか。`must_fail` は boolean。`safe_ignore_reason` は safe ignore 以外では `null`。 |
| `input/determinism_inventory.json` | `source_kind`、`fake_adapter`、`repeat_count` | `source_kind` は `clock`、`timer`、`sleep`、`timeout`、`entropy`、`http_client`、`command`、`systemd`、`filesystem_order`、`parallel_worker` のいずれか。`fake_adapter` は責務正本 anchor または `null`。`repeat_count` は 2 以上。 |
| `input/filesystem_inventory.json` | `operation_kind`、`durability_requirement`、`failure_class` | `operation_kind` は `write`、`append`、`flush`、`fsync`、`parent_fsync`、`rename`、`remove`、`walk`、`lock`、`symlink_check`、`recovery` のいずれか。`durability_requirement` は statefile owner の手順または collaborator 固有手順への link。`failure_class` は 1 件以上の array。 |
| `input/workflow_inventory.json` | `workflow_path`、`job_name`、`hardening_items` | `hardening_items` は `action_sha_pin`、`permissions`、`timeout`、`required_check`、`skip_success_forbidden`、`docker_fallback`、`runtime_version`、`deno_check`、`release_rehearsal` のうち 1 件以上。Deno 不在時は Docker 上の Deno stable runtime 実行証跡だけを代替証跡として認める。 |
| `input/large_owner_inventory.json` | `owner_path`、`go_file_count`、`risk_items` | `go_file_count` は 5 以下。`risk_items` は `oversized_file`、`mixed_responsibility`、`missing_internal_section`、`test_gap`、`dummy_or_empty_file`、`subdirectory`、`none` のうち 1 件以上。production Go file が 2000 行を超える場合は `oversized_file` とする。test Go file が 2000 行を超える場合は `test_gap` として test artifact 分割または検証責務の見直し対象にする。`none` は `classification=not_applicable` の場合だけ許可する。 |
| `input/negative_controls.json` | `control_id`、`target_contract`、`mutation_kind`、`fixture_mutation`、`expected_failure_code`、`expected_counter_key`、`expected_exit_code`、`expected_diagnostic_ref`、`positive_control_ref`、`isolation_mode` | `control_id` は `phase16.negative.<slug>` 形式。`target_contract` は `manifest`、`source_coverage`、`inventory`、`expected`、`record`、`checker_stdout`、`diagnostic_jsonl`、`counter`、`runtime`、`workflow`、`mutation`、`fault`、`document_drift` のいずれか。`mutation_kind` は 失敗させる契約を 1 つだけ表す lowercase ASCII slug。`fixture_mutation` は対象 path、変更内容、変更前 sha256、変更後 sha256 を含む object。`expected_failure_code` は Phase 16 checker failure code 表のいずれか。`expected_counter_key` は本節 closure counter 表のいずれか。`expected_exit_code` は `1` 固定。`expected_diagnostic_ref` は期待 diagnostic の責務名付き Markdown link。`positive_control_ref` は同じ checker が正常 fixture で成功する証跡への責務名付き Markdown link。`isolation_mode` は `temporary_copy` 固定。live fixture root を直接破壊する負例実行を禁止する。 |

`input/source_coverage.json` の `detected_terms[]` は object とし、`term_id`、`match`、`match_mode`、`count`、`first_source_locator` を必須とする。`term_id`、`match`、`match_mode` は Phase 16 source coverage detection registry の行と完全一致する。`count` は 1 以上の integer、`first_source_locator` は `file:<path>#byte:<offset>` または `file:<path>#line:<line_number>` のいずれかとし、行番号を実装契約の正本として扱ってはならない。同一 `path` 内の `detected_terms[]` は `term_id`、`match`、`first_source_locator` の ASCII 昇順とし、registry に存在しない `term_id`、raw string だけの検出語、`count=0`、未接続 `inventory_refs` を禁止する。

Phase 16 の `expected/counters.json` は root object に `schema_version`、`scope`、`counters` だけを持つ。未知 key、欠落 key、型違い、重複 key を禁止する。`schema_version` は `1`、`scope` は `phase-16-quality-evidence-closure`、`counters` は本節 closure counter 表の全 key だけを 1 回ずつ持つ object とする。各 counter 値は integer `0` 固定であり、string `"0"`、boolean、null、負数、未記録、追加 counter、欠落 counter、同名 counter の重複、closure record set と異なる値を禁止する。

Phase 16 の `expected/actions.json` は root object に `schema_version`、`scope`、`actions` だけを持つ。未知 key、欠落 key、型違い、重複 key を禁止する。`actions[]` は `id`、`inventory_id`、`action`、`owner_component`、`expected_counter_key`、`expected_status`、`evidence_record` を必須とする。`inventory_id` はいずれかの inventory record `id` と一致し、`action` は `execute_and_close`、`fix_and_close`、`classify_not_applicable`、`keep_future_plan`、`reject_as_invalid_evidence` のいずれか、`evidence_record` は完了時に `records/closure.jsonl`、`records/execution.jsonl`、`records/mutation.jsonl`、`records/fault.jsonl`、`records/workflow.jsonl` のいずれかへ到達する。1 つの `inventory_id` に対する action は 1 件だけとし、複数 action、未接続 action、inventory に存在しない action、action のない inventory record を禁止する。

Phase 16 の `records/*.jsonl` は、1 行 1 JSON object、UTF-8、LF 終端、空行禁止、未知 key 禁止とする。全 record は `record_id`、`scope`、`work_unit`、`inventory_id`、`owner_component`、`source_ref`、`result`、`counter_key` を必須とする。`record_id` は `phase16.record.<record_file_slug>.<owner_or_area>.<slug>` 形式、`record_file_slug` は拡張子なしの record file 名とする。`inventory_id` はいずれかの inventory record `id` と一致する。`scope` は `phase-16-quality-evidence-closure`、`result` は `passed`、`failed`、`not_applicable`、`future_plan`、`rejected` のいずれかとし、完了時に `failed` は 0 件でなければならない。

| record file | 追加必須 key | 固定条件 |
|-------------|--------------|----------|
| `records/execution.jsonl` | `command`、`runtime`、`input_sha256`、`expected_sha256`、`actual_sha256`、`diff_result`、`exit_code`、`stdout_sha256`、`stderr_sha256` | `command` は実行した本番 entrypoint または checker。`diff_result` は `empty` または failure reason。exit code だけを合格根拠にしてはならない。 |
| `records/mutation.jsonl` | `mutation_id`、`target_ref`、`mutation_kind`、`expected_failure`、`actual_failure`、`mutation_result` | `mutation_result` は `killed`、`survived`、`invalid`、`equivalent` のいずれか。`survived` は完了時 0 件。`equivalent` は責務正本 anchor がある場合だけ許可する。 |
| `records/fault.jsonl` | `fault_id`、`fault_kind`、`injection_point`、`expected_recovery`、`actual_recovery`、`silent_success` | `silent_success` は完了時 `false`。`fault_kind` は disk full、permission denied、short write、fsync failure、rename failure、process kill、partial write、client disconnect のいずれか。 |
| `records/workflow.jsonl` | `workflow_path`、`job_name`、`action_pin_result`、`permissions_result`、`timeout_result`、`required_check_result`、`docker_fallback_result`、`runtime_version_result`、`deno_check_result`、`release_rehearsal_result` | 各 result は `passed`、`failed`、`not_applicable` のいずれか。`workflow_path` は `.github/workflows/*.yml` または `.github/workflows/*.yaml` に一致する。`deno_check_result=not_applicable` は Phase 16 で standalone JavaScript file を触らない場合だけ許可し、正本 anchor 必須。`release_rehearsal_result=not_applicable` は release / setup / distribution に影響がない場合だけ許可し、正本 anchor 必須。その他の `not_applicable` も正本 anchor 必須。 |
| `records/closure.jsonl` | `closure_status`、`open_items`、`closed_items`、`not_applicable_items`、`future_items` | `closure_status` は `closed`、`open` のいずれか。完了時は `closed` かつ `open_items=[]`。 |

<a id="phase-16-checker-contract"></a>
**Phase 16 checker 固定契約：**

Phase 16 checker は `testdata/phase16/quality-evidence-closure/` を入力に取り、live source coverage set、`manifest.json`、`input/*.json`、`expected/*.json`、`records/*.jsonl`、closure counter を相互照合する単一 gate とする。checker は JSON decode 前に duplicate key を検出し、duplicate key、unknown key、欠落 key、型違い、空文字列、不正 enum、不正 sort、LF 終端欠落、JSONL 空行、同一 `id` / `record_id` の重複をすべて失敗にする。

Phase 16 checker の標準実行入口は `go test ./... -run TestPhase16QualityEvidenceClosure -count=1` に固定する。required check `phase16-quality-evidence-fixture` はこの command を実行し、`records/execution.jsonl` の checker record は同じ command、runtime、stdout sha256、stderr sha256、actual sha256、expected sha256、exit code を記録する。Phase 16 checker のために新しい owner component、専用 CLI、別 executable、外部 tool、6 ファイル目、または owner subdirectory を作成してはならない。checker 実装は既存の repository-level Go test artifact または既存 owner の `<owner>_test.go` に置き、[`docs/SPEC.md` 方針責務 §4.3](../SPEC.md#sec-4-3) の 5 ファイル原則を崩してはならない。

checker は `rg --files --hidden -g '!.git/**'` の実行結果から Phase 16 source coverage set を再構築し、`input/source_coverage.json` と一致させる。live source coverage set に存在する path が `input/source_coverage.json` にない場合、または `input/source_coverage.json` の path が live source coverage set に存在しない場合、`phase16_document_drift_open_count` と `phase16_execution_evidence_gap_count` を `0` にしてはならない。checker は各 source の sha256、検出対象語、owner / area、inventory_refs を再計算し、対応 inventory record がない検出対象を失敗にする。

checker は各 inventory record から `expected/actions.json` と `records/*.jsonl` を再導出し、1 inventory record に 1 action、1 action に 1 件以上の evidence record、1 evidence record に 1 counter が接続していることを確認する。`closure_ref`、`not_applicable_ref`、`future_ref`、`counter_key`、`status`、`result` の組み合わせが矛盾する場合、手書き counter が `0` でも失敗にする。

checker は closure counter を `input/*.json` と `records/*.jsonl` から再集計する。`expected/counters.json`、`manifest.json.counters`、closure record set に記録された counter が再集計値と一致しない場合、または failed / open / survived / silent_success / unclassified / gap が 1 件以上ある場合、該当 counter と `final_open_item_count` は `0` にできない。

checker の診断は deterministic とし、`work_unit`、`owner_component`、`source_ref`、`source_locator`、`id`、`record_id` の ASCII 昇順で出力する。診断は free-form message だけにしてはならず、必ず diagnostic JSON Lines の固定 schema で出力する。checker が内部 error、panic、partial read、JSONL parse error、filesystem error、source enumeration error、hash mismatch を検出した場合は、成功扱いせず `phase16_execution_evidence_gap_count` に計上する。

checker の process 入出力は次のとおり固定する。成功時は exit code `0`、stdout は root object に `schema_version`、`scope`、`result`、`source_count`、`inventory_count`、`action_count`、`record_count`、`negative_control_count`、`negative_control_failed_count`、`negative_control_passed_count`、`diagnostic_count`、`error_count`、`warning_count`、`manifest_sha256`、`source_coverage_sha256`、`negative_controls_sha256`、`inventory_sha256`、`expected_sha256`、`counters_sha256`、`records_sha256`、`final_open_item_count` だけを持つ JSON object、stderr は空とする。`negative_control_failed_count` は失敗すべき negative control を checker が fail として検出した件数、`negative_control_passed_count` は失敗すべき negative control を誤って成功扱いした件数とする。成功時は `result=passed`、`diagnostic_count=0`、`error_count=0`、`warning_count=0`、`negative_control_passed_count=0`、`negative_control_failed_count` は `negative_control_count` と一致、`final_open_item_count=0` とする。失敗時は exit code `1`、stdout は成功時と同じ root key を持ち `result=failed`、`final_open_item_count` は 1 以上、stderr は 1 行 1 JSON object の diagnostic JSON Lines とする。checker の exit code だけ、stdout summary だけ、stderr message だけを合格根拠にしてはならず、`records/execution.jsonl` の `stdout_sha256`、`stderr_sha256`、`actual_sha256` と照合する。

diagnostic JSON Lines は `failure_code`、`severity`、`inventory_id`、`record_id`、`counter_key`、`source_ref`、`expected`、`actual`、`remediation_ref` だけを持つ。`severity` は `error`、`warning`、`info` のいずれかとし、Phase 16 完了時は `warning` と `info` も 0 件でなければならない。`failure_code` は uppercase ASCII の `PHASE16_<AREA>_<REASON>` 形式とし、下表のいずれかだけを許可する。`inventory_id` と `record_id` は string または null とし、JSON decode、duplicate key、source enumeration、manifest schema、source coverage schema、expected schema のように inventory / record へ到達する前に失敗した場合だけ null を許可する。`counter_key` は本節 closure counter 表の key、または counter 未確定の schema failure では `phase16_execution_evidence_gap_count` とする。`remediation_ref` は責務名付き Markdown link だけを許可し、裸 URL、裸 path、free-form 対応メモを禁止する。

| failure_code | 失敗条件 |
|--------------|----------|
| `PHASE16_JSON_DUPLICATE_KEY` | JSON object または JSONL record に duplicate key がある。 |
| `PHASE16_SCHEMA_UNKNOWN_KEY` | schema にない key がある。 |
| `PHASE16_SCHEMA_MISSING_KEY` | 必須 key が欠ける。 |
| `PHASE16_SCHEMA_TYPE_MISMATCH` | key の型、enum、null 可否、空文字列、空配列条件が一致しない。 |
| `PHASE16_SOURCE_MISSING` | source coverage set の必須 path が不在、または `source_exists=false` が `go.sum` 以外に使われた。 |
| `PHASE16_SOURCE_UNTRACKED` | live source coverage set の path が `input/source_coverage.json` に存在しない。 |
| `PHASE16_SOURCE_EXTRA` | `input/source_coverage.json` の path が live source coverage set に存在しない。 |
| `PHASE16_SOURCE_HASH_MISMATCH` | `sha256`、`size_bytes`、または `source_exists` が live source と一致しない。 |
| `PHASE16_TERM_UNCLASSIFIED` | detection registry で検出した term に対応 inventory record がない。 |
| `PHASE16_TERM_SCHEMA_INVALID` | `detected_terms[]` が object schema、registry id、sort、count、locator 条件を満たさない。 |
| `PHASE16_COUNTER_MISMATCH` | 再集計 counter、`expected/counters.json`、`manifest.json.counters`、closure record の値が一致しない。 |
| `PHASE16_RECORD_UNLINKED` | inventory、action、record、counter の 1 対 1 接続が欠ける。 |
| `PHASE16_RECORD_RESULT_OPEN` | `failed`、`open`、`survived`、`silent_success=true`、未分類、gap が残る。 |
| `PHASE16_NEGATIVE_CONTROL_SCHEMA_INVALID` | `input/negative_controls.json` の schema、期待 failure code、isolation、positive control 接続が不正である。 |
| `PHASE16_NEGATIVE_CONTROL_PASSED` | 失敗すべき negative control が成功扱いになった。 |
| `PHASE16_RUNTIME_UNEXECUTED` | 必須 runtime / production entrypoint / Deno / Docker / mutation / fault / race / release rehearsal が未実行である。 |
| `PHASE16_DOCUMENT_DRIFT` | path、anchor、状態、workflow、required check、対象外理由、将来計画維持に drift がある。 |
| `PHASE16_INTERNAL_ERROR` | checker の内部 error、panic、partial read、filesystem error、source enumeration error が発生した。 |

Phase 16 closure record set は、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の 18 record を使用し、各 record の `scope` に `phase-16-quality-evidence-closure` を含める。Phase 16 では以下の集計値を同じ closure record set 内に記録する。

| 集計値 | 完了値 | 未完了条件 |
|--------|--------|------------|
| `phase16_declarative_evidence_open_count` | `0` | JSONL record 存在、counter `0`、PR 本文記載、manifest 存在だけを実行証跡として扱う item が 1 件以上ある。 |
| `phase16_execution_evidence_gap_count` | `0` | production entrypoint 実行、actual / expected 比較、diff、stdout / stderr、exit code、fixture harness 接続が欠ける item がある。 |
| `phase16_mutation_survived_count` | `0` | production code mutation の survived が 1 件以上、または mutation を実際に加えた証跡がない。 |
| `phase16_fault_injection_open_count` | `0` | disk full、permission denied、short write、fsync failure、rename failure、process kill、partial write、client disconnect の注入・復旧証跡に残件がある。 |
| `phase16_align_unclassified_count` | `0` | `ALIGN-*` 追加検証候補に Phase 16 実施、将来計画維持、対象外 anchor のいずれもない。 |
| `phase16_align_open_count` | `0` | Phase 16 実施対象に分類した `ALIGN-*` に未実装、未検証、未証跡がある。 |
| `phase16_ignored_error_unclassified_count` | `0` | `_ =`、ignored write、cleanup、hash writer、best effort の分類がない箇所がある。 |
| `phase16_required_write_failure_open_count` | `0` | audit、access、config、state、JSON Lines、release、archive、rollback、credential の必須 write failure を無視する経路がある。 |
| `phase16_skip_open_count` | `0` | skip / 未実行が対象外 anchor、代替実行証跡、または失敗扱いに接続していない。 |
| `phase16_panic_unclassified_count` | `0` | panic が programmer error として明示許可されるか、error return に置換されるかの分類がない。 |
| `phase16_determinism_open_count` | `0` | clock、sleep、timeout、entropy、parallel worker、retry backoff が固定 fake または実行証跡に接続していない。 |
| `phase16_http_boundary_open_count` | `0` | response write failure、partial write、client disconnect、body limit、graceful shutdown、timeout、redirect / SSRF 境界の証跡に残件がある。 |
| `phase16_filesystem_durability_open_count` | `0` | atomic write、fsync、parent directory fsync、rename、symlink 非追従、stale lock、recovery、permission failure の owner 間水準差が残る。 |
| `phase16_workflow_hardening_open_count` | `0` | action SHA pin、minimum permissions、timeout、required check、skip success 禁止、Phase 12 workflow hardening に残件がある。 |
| `phase16_validation_portability_open_count` | `0` | host Go / Deno 不在時の Docker 検証手順、runtime version、実行 command、失敗時扱いが固定されていない。 |
| `phase16_release_rehearsal_open_count` | `0` | release rehearsal、tag / checksum / signature / SBOM / reproducible build 証跡、setup install rehearsal に残件がある。 |
| `phase16_large_owner_risk_open_count` | `0` | 巨大 owner file に risk ledger、内部責務区画、test artifact 接続、5 ファイル原則維持証跡がない。 |
| `phase16_document_drift_open_count` | `0` | Phase 16 の状態、path、anchor、fixture root、workflow、required check、対象外理由、将来計画維持の drift がある。 |
| `final_open_item_count` | `0` | 上記集計値または closure record の `open_items` に残件がある。 |

Phase 16 required check は以下に固定する。Phase 16 実装 PR は、下表の check 名、対象、完了条件を同一 Pull Request 本文へ記録する。

| check 名 | 対象 | 完了条件 |
|----------|------|----------|
| `phase16-go-format` | Phase 16 で触れた Go file | `gofmt -l` の差分が 0。 |
| `phase16-go-test` | Phase 16 の対象 owner と横断 gate | `go test ./...` が成功し、skip がある場合は [`docs/details/fixture.md` fixture 証跡責務 skip / 未実行証跡固定契約](#test-skip-evidence-contract) に接続する。 |
| `phase16-deno-check` | [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) と Phase 16 で触れた standalone JavaScript file | Deno stable runtime の `deno check` が成功する。Deno が host にない場合は Docker 上の Deno stable runtime 実行証跡を `phase16_validation_portability_open_count=0` に接続する。Node.js、npm、bundler、transpiler を代替完了根拠にしてはならない。[`admin/index.html`](../../admin/index.html) 内の inline JavaScript は `phase16-quality-evidence-fixture` で UI / SDK / API interaction として検証し、standalone JavaScript file として扱わない。 |
| `phase16-race` | statefile、runner、api、mcp、obsidian、archive、release、setup の並行更新または lock に関係する test | race detector または仕様化済み interleaving 証跡が成功し、`phase16_determinism_open_count=0`。 |
| `phase16-quality-evidence-fixture` | `testdata/phase16/quality-evidence-closure/` | `go test ./... -run TestPhase16QualityEvidenceClosure -count=1` が production entrypoint 実行、actual / expected 比較、checker 再導出、checker stdout schema、diagnostic JSON Lines schema、closure counters、negative control を照合し、成功時 `diagnostic_count=0`、`error_count=0`、`warning_count=0`、`negative_control_passed_count=0`、`final_open_item_count=0` である。 |
| `phase16-mutation` | production code mutation | `records/mutation.jsonl` の `survived=0`。 |
| `phase16-fault-injection` | filesystem、HTTP、process、permission、write failure | `records/fault.jsonl` が全 failure class を実行し、silent success が 0。 |
| `phase16-workflow-hardening` | `.github/workflows/*.yml`、`.github/workflows/*.yaml` と required check | action SHA pin、permissions、timeout、required check、Phase 12 workflow hardening がすべて閉じる。 |
| `phase16-document-drift` | [`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md)、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md)、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md)、[`docs/details/fixture.md`](fixture.md)、対象 owner 詳細本文 | Phase 16 の状態、path、anchor、fixture root、workflow、counter、対象外理由、将来計画維持の drift が 0。 |

Phase 16 fixture は `input/negative_controls.json` による negative control を必須とする。negative control は、JSON duplicate key を受理、diagnostic JSON Lines の必須 key 欠落を受理、`severity=warning` を成功扱い、`failure_code` registry 外値を受理、source coverage set の path 省略、`go.sum` 不在 record 欠落、`go.mod` の外部 `require` / `replace` を未検出、`detected_terms` を raw string だけで受理、detection registry の検出漏れ、checker stdout / stderr schema 不一致、JSONL record だけで成功、counter だけで成功、実行していない mutation を killed と扱う、survived mutation を 0 と誤集計、fault injection 未実行、skip を成功扱い、required write failure の無視、panic 未分類、clock / sleep / timeout の非決定、response write failure 無視、client disconnect 無視、fsync failure 無視、rename failure 無視、parent directory fsync 欠落、workflow tag pin、workflow permissions 欠落、workflow timeout 欠落、host tool 不在時の成功扱い、Deno 未実行の JavaScript 成功扱い、release rehearsal 未実行、巨大 owner risk ledger 欠落、`ALIGN-*` 未分類、将来計画対象を Phase 16 完了対象として誤計上する case を 1 件以上含める。checker は各 negative control を `isolation_mode=temporary_copy` で fixture root の一時 copy にだけ適用し、live fixture root、source coverage set、expected、record を直接破壊してはならない。negative control が成功扱いになる場合、該当 checker 自体を未完成として `phase16_execution_evidence_gap_count`、`phase16_align_unclassified_count`、または該当する `phase16_*_open_count` に計上する。

<a id="fixture-root-coverage-matrix-contract"></a>
**fixture root coverage matrix 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、Phase 11 対象 fixture root の網羅判定と未完了条件だけを固定する。実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](../DOCUMENT_INDEX.md#実装ファイル一覧)、現在状態と対象外理由は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8) を参照する。

Phase 11 対象 root は、正式 fixture directory、実装検証証跡、または [`docs/ROADMAP.md` 状態・計画責務](../ROADMAP.md) の対象外理由へ到達できる場合だけ閉じる。Phase 11 が割り当てた未作成 root は、正式 fixture、実装検証証跡、または対象外理由が追加されるまで未完了として扱う。`testdata/admin/` は fixture group root であり、正式 fixture directory として数えない。`testdata/admin/cli/` は Admin CLI の formal fixture root であり、その配下で directory 名、`manifest.json.name`、fixture catalog 名が 1 対 1 に一致する下位 directory だけを正式 fixture directory として扱う。下表の未作成 root 共通未完了条件は、正式 fixture、実装検証証跡、対象外理由のいずれにも到達できない状態とする。

| root / pattern | coverage status | closure source | incomplete condition |
|----------------|-----------------|----------------|----------------------|
| `testdata/builder/` | builder の正式 fixture root。`single/`、`site/`、`empty-dir/`、`strict/`、`safe/`、`url-safety/` は `manifest.json` と `expected/` を持つ正式 fixture directory である。 | [`§8a-F`](#8a-f-builder-初期受け入れ-fixture-契約)、[`§28-F`](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md#実装ファイル一覧) | catalog / manifest / expected / 実装検証証跡の対応が欠ける場合は未完了。 |
| `testdata/admin/cli/` | Admin CLI の正式 fixture root。 | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[`docs/DETAIL_INDEX.md` Phase 11 参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) | catalog / manifest / expected / effects / security の対応が欠ける場合は未完了。 |
| `testdata/setup/` | setup の正式 fixture root。 | [`§27-F setup / admin / Release asset 連動 fixture`](#sec-27-f-19)、[`docs/DETAIL_INDEX.md` Phase 11 参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) | catalog / manifest / expected / effects の対応が欠ける場合は未完了。 |
| `testdata/release/` | release の正式 fixture root。 | [Release fixture 固定契約](#release-fixture-contract)、[`docs/DETAIL_INDEX.md` Phase 11 参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) | catalog / manifest / expected / asset 証跡の対応が欠ける場合は未完了。 |
| `testdata/mcp/` | MCP の正式 fixture root。 | [MCP fixture 固定契約](#mcp-fixture-contract)、[`docs/DETAIL_INDEX.md` Phase 11 参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) | catalog / manifest / expected / fake transcript の対応が欠ける場合は未完了。 |
| `testdata/runner/` | Phase 11 root coverage formal fixture root。 | [`success-runner-phase11-root-coverage`](../../testdata/runner/success-runner-phase11-root-coverage/)、[`§15a-F`](#15a-f-runner-初期受け入れ-fixture-契約)、[`§27-F`](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。runner catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の runner 行で別管理する。 |
| `testdata/api/` | Phase 11 root coverage formal fixture root。 | [`success-api-phase11-root-coverage`](../../testdata/api/success-api-phase11-root-coverage/)、[`§22-F`](#22-f-api-fixture-契約)、[`docs/DETAIL_INDEX.md` Phase 11 参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。API catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の API 行で別管理する。 |
| `testdata/sdk/` | Phase 11 root coverage formal fixture root。 | [`success-sdk-phase11-root-coverage`](../../testdata/sdk/success-sdk-phase11-root-coverage/)、[`§22-F`](#22-f-api-fixture-契約)、[`docs/details/sdk.md` 詳細本文責務](sdk.md) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。SDK catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の SDK 行で別管理する。 |
| `testdata/ui/` | Phase 11 root coverage formal fixture root。 | [`success-ui-phase11-root-coverage`](../../testdata/ui/success-ui-phase11-root-coverage/)、[`§22-F`](#22-f-api-fixture-契約)、[`docs/details/ui.md` 詳細本文責務](ui.md) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。UI catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の UI 行で別管理する。 |
| `testdata/statefile/` | Phase 11 root coverage formal fixture root。 | [`success-statefile-phase11-root-coverage`](../../testdata/statefile/success-statefile-phase11-root-coverage/)、[`docs/details/statefile.md` 詳細本文責務](statefile.md)、[non-dedicated owner test routing 固定契約](#non-dedicated-owner-test-routing-contract) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。statefile catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の Statefile 行で別管理する。 |
| `testdata/archive/` | Phase 11 root coverage formal fixture root。 | [`success-archive-phase11-root-coverage`](../../testdata/archive/success-archive-phase11-root-coverage/)、[`docs/details/archive.md` 詳細本文責務](archive.md)、[non-dedicated owner test routing 固定契約](#non-dedicated-owner-test-routing-contract) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。Archive catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の Archive 行で別管理する。 |
| `testdata/commitstatus/` | Phase 11 root coverage formal fixture root。 | [`success-commitstatus-phase11-root-coverage`](../../testdata/commitstatus/success-commitstatus-phase11-root-coverage/)、[`docs/details/commitstatus.md` 詳細本文責務](commitstatus.md)、[non-dedicated owner test routing 固定契約](#non-dedicated-owner-test-routing-contract) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。Commit status catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の Commit status 行で別管理する。 |
| `testdata/security/` | Phase 11 root coverage formal fixture root。 | [`success-security-phase11-root-coverage`](../../testdata/security/success-security-phase11-root-coverage/)、[`docs/details/security.md` 詳細本文責務](security.md)、[non-dedicated owner test routing 固定契約](#non-dedicated-owner-test-routing-contract) | root coverage fixture directory、manifest、input、expected が欠ける場合は未完了。Security catalog 別 fixture 群は [`§0g.8-F`](#sec-0g-8-f) の Security 行で別管理する。 |
| `testdata/admin/archive/`、`testdata/admin/static-serving/`、`testdata/admin/security/` | Admin group root 配下の候補 root。現時点では Admin CLI の正式 fixture root ではない。 | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[`docs/ROADMAP.md` 状態・計画責務](../ROADMAP.md) | Admin CLI 完了証跡として数えた場合は未完了。作成する場合は該当 owner / state / fixture catalog を先に一致させる。 |

<a id="fixture-root-missing-closure-record-contract"></a>
**未作成 fixture root closure record 固定契約：**

Phase 11 対象 root のうち [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](../DOCUMENT_INDEX.md#実装ファイル一覧) で `未作成` とされている root は、正式 fixture directory、実装検証証跡、または [`docs/ROADMAP.md` 状態・計画責務](../ROADMAP.md) の対象外理由へ到達できるまで closure できない。未作成 root を hidden fixture、既存 Go test、または PR 本文の説明だけで閉じてはならない。

未作成 fixture root closure record は、root ごとに 1 件作成する。record は次の field を持つ。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `fixture_root` | `testdata/<owner>/` または固定表に記載された root / pattern。 | root が空、または [`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md#実装ファイル一覧) と一致しない。 |
| `owner_component` | 対象 owner component 名。 | owner 不明、または root から推測だけで決めている。 |
| `document_index_status` | `実在` または `未作成`。 | 索引責務の所在区分と一致しない。 |
| `closure_method` | `formal_fixture`、`implementation_evidence`、`not_applicable`、`open` のいずれか。 | 未登録値、または根拠なしで `formal_fixture` / `implementation_evidence` / `not_applicable` にしている。 |
| `closure_refs` | 正式 fixture directory、実装検証証跡、または対象外理由 anchor への責務名付き Markdown link 配列。 | 空、実在しない path、または責務正本 anchor なし。 |
| `open_items` | 作成すべき fixture、expected、fake、assertion、mutation class、drift record を列挙する。 | `closure_method=open` なのに不足内容が不明、または `closure_method` が open 以外なのに残件がある。 |

未作成 fixture root closure record の `closure_method=open` が 1 件でも残る場合、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `fixture_root_closure` と `final_open_item_count` は closed にできない。

Phase 11 root coverage closure record は次のとおりとする。

| fixture_root | owner_component | document_index_status | closure_method | closure_refs | open_items |
|--------------|-----------------|-----------------------|----------------|--------------|------------|
| [`testdata/builder/`](../../testdata/builder/) | `builder` | `実在` | `formal_fixture` | [`single`](../../testdata/builder/single/)、[`site`](../../testdata/builder/site/)、[`empty-dir`](../../testdata/builder/empty-dir/)、[`strict`](../../testdata/builder/strict/)、[`safe`](../../testdata/builder/safe/)、[`url-safety`](../../testdata/builder/url-safety/) | `[]` |
| [`testdata/runner/`](../../testdata/runner/) | `runner` | `実在` | `formal_fixture` | [`success-runner-phase11-root-coverage`](../../testdata/runner/success-runner-phase11-root-coverage/) | `[]` |
| [`testdata/api/`](../../testdata/api/) | `api` | `実在` | `formal_fixture` | [`success-api-phase11-root-coverage`](../../testdata/api/success-api-phase11-root-coverage/) | `[]` |
| [`testdata/sdk/`](../../testdata/sdk/) | `sdk` | `実在` | `formal_fixture` | [`success-sdk-phase11-root-coverage`](../../testdata/sdk/success-sdk-phase11-root-coverage/) | `[]` |
| [`testdata/ui/`](../../testdata/ui/) | `ui` | `実在` | `formal_fixture` | [`success-ui-phase11-root-coverage`](../../testdata/ui/success-ui-phase11-root-coverage/) | `[]` |
| [`testdata/statefile/`](../../testdata/statefile/) | `statefile` | `実在` | `formal_fixture` | [`success-statefile-phase11-root-coverage`](../../testdata/statefile/success-statefile-phase11-root-coverage/) | `[]` |
| [`testdata/archive/`](../../testdata/archive/) | `archive` | `実在` | `formal_fixture` | [`success-archive-phase11-root-coverage`](../../testdata/archive/success-archive-phase11-root-coverage/) | `[]` |
| [`testdata/commitstatus/`](../../testdata/commitstatus/) | `commitstatus` | `実在` | `formal_fixture` | [`success-commitstatus-phase11-root-coverage`](../../testdata/commitstatus/success-commitstatus-phase11-root-coverage/) | `[]` |
| [`testdata/security/`](../../testdata/security/) | `security` | `実在` | `formal_fixture` | [`success-security-phase11-root-coverage`](../../testdata/security/success-security-phase11-root-coverage/) | `[]` |

Phase 11 の test / contract drift 判定では、`main_test.go`、`components/*_test.go`、`sdk_contract_test.go`、`ui_contract_test.go` を fixture 証跡の補助 source として扱う。これらの test / contract assertion は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry)、owner 詳細本文、fixture catalog、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract)、または [`docs/ROADMAP.md` 状態・計画責務](../ROADMAP.md) の対象外理由へ到達できなければならない。test file 自体を仕様正本として扱ってはならず、fixture catalog と owner 詳細本文へ到達しない assertion は Phase 11 の未解消 drift とする。mutation test 証跡は [mutation test 証跡固定契約](#mutation-test-evidence-contract) に従い、Phase 11 では `survived=0` を満たすまで完了扱いにしない。

<a id="test-artifact-traceability-contract"></a>
**test artifact traceability 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実在 test artifact から owner 詳細本文、fixture root、boundary / failure matrix 証跡、concurrency / race 証跡、mutation test 証跡、test / contract drift 証跡へ到達するための接続だけを固定する。テスト方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](../DOCUMENT_INDEX.md#実装ファイル一覧)、owner 選択入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務](../DETAIL_INDEX.md#0b-詳細仕様参照表) を参照する。

| test artifact | owner / collaborator 接続 | fixture / 証跡接続 | drift 判定 |
|---------------|---------------------------|---------------------|------------|
| [`main_test.go`](../../main_test.go) | 起動入口 artifact と `builder`、`runner`、`api`、`admin`、`setup`、`release`、`mcp` の CLI 起動境界。owner 選択は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) の `main dispatch / binary version regression gate` を参照する。 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract)、[Phase 11 fixture harness 参照](#phase-11-fixture-harness-reference)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | 起動名、`--help`、`--version`、未知 basename、binary version、dispatch 副作用なしの assertion が owner 詳細本文または共通固定契約へ到達できない場合は [test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の孤立 test とする。 |
| [`components/builder/builder_test.go`](../../components/builder/builder_test.go) | `builder` owner。詳細本文は [`docs/details/builder.md` 詳細本文責務](builder.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §8a-F](fixture.md#8a-f-builder-初期受け入れ-fixture-契約)、[`§28-F`](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、[`testdata/builder/`](../../testdata/builder/)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | Markdown、HTML / CSS / JavaScript 生成、atomic output、URL safety、runtime fixture、manifest / expected 接続が `builder` 詳細本文または fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/runner/runner_test.go`](../../components/runner/runner_test.go) | `runner` owner、必要に応じて `statefile`、`archive`、`commitstatus`、`security` collaborator。詳細本文は [`docs/details/runner.md` 詳細本文責務](runner.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約)、[`§27-F`](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、`testdata/runner/`、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | GitHub / local target、pipeline、deploy、queue、lock、secret mask、state write、notification、Commit Status の assertion が owner 詳細本文、collaborator 詳細本文、または fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/api/api_test.go`](../../components/api/api_test.go) | `api` owner、必要に応じて `security`、`statefile`、`runner`、`archive` collaborator。詳細本文は [`docs/details/api.md` 詳細本文責務](api.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、`testdata/api/`、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | endpoint、HTTP status、response schema、認証・認可、rate limit、状態 read/write、audit / access log、副作用なしの assertion が API 詳細本文または fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/admin/admin_test.go`](../../components/admin/admin_test.go) | `admin` owner。CLI 管理クライアントの詳細本文は [`docs/details/admin.md` 詳細本文責務](admin.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 Admin CLI fixture 固定契約](fixture.md#admin-cli-fixture-contract)、[`testdata/admin/cli/`](../../testdata/admin/cli/)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | CLI request、stdout / stderr、exit code、secret redaction、transport 境界、fixture presence check が `admin` 詳細本文または Admin CLI fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/setup/setup_test.go`](../../components/setup/setup_test.go) | `setup` owner、必要に応じて `admin`、`release`、`security` collaborator。詳細本文は [`docs/details/setup.md` 詳細本文責務](setup.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture](fixture.md#sec-27-f-19)、[`testdata/setup/`](../../testdata/setup/)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | Release asset 取得、checksum、配置、systemd、rollback、secret preservation、admin archive 境界の assertion が `setup` 詳細本文または fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/release/release_test.go`](../../components/release/release_test.go) | `release` owner、必要に応じて `admin`、`security` collaborator。詳細本文は [`docs/details/release.md` 詳細本文責務](release.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 Release fixture 固定契約](fixture.md#release-fixture-contract)、[`testdata/release/`](../../testdata/release/)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | clean checkout、tag、binary build、admin archive、checksum、GitHub draft / asset / publish / cleanup、token mask の assertion が `release` 詳細本文または Release fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go) | `mcp` owner、必要に応じて `api`、`security`、`statefile` collaborator。詳細本文は [`docs/details/mcp.md` 詳細本文責務](mcp.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 MCP fixture 固定契約](fixture.md#mcp-fixture-contract)、[`testdata/mcp/`](../../testdata/mcp/)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | CLI、HTTP、JSON-RPC、tools、resources、prompts、sampling、SSE、auth、scope、confirmation、audit、metrics の assertion が `mcp` 詳細本文または MCP fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`sdk_contract_test.go`](../../sdk_contract_test.go) | `sdk` owner、`api` collaborator。詳細本文は [`docs/details/sdk.md` 詳細本文責務](sdk.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、`testdata/sdk/`、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | public method、request shape、query encoding、error class、timeout、stream、binary response、禁止 API の assertion が `sdk` 詳細本文、API 詳細本文、または fixture catalog へ到達できない場合は未解消 drift とする。 |
| [`ui_contract_test.go`](../../ui_contract_test.go) | `ui` owner、`sdk`、`api`、`security` collaborator。詳細本文は [`docs/details/ui.md` 詳細本文責務](ui.md) を参照する。 | [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、`testdata/ui/`、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract)、[mutation test 証跡固定契約](#mutation-test-evidence-contract)。 | DOM 固定値、SDK 境界、base URL、one-time secret 消去、stream terminal、主要操作、status 別 error handling の assertion が `ui` 詳細本文、SDK / API 詳細本文、または fixture catalog へ到達できない場合は未解消 drift とする。 |

<a id="non-dedicated-owner-test-routing-contract"></a>
**non-dedicated owner test routing 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、単独 test artifact を持たない owner component の検証接続だけを固定する。owner component の責務本文、状態、入出力、異常系、完了可否は各 owner 詳細本文を正とし、本項では再定義しない。

単独 Go artifact を持たない owner component について、専用 `components/<owner>_test.go` が存在しないことだけを未検証または仕様不足として扱ってはならない。ただし、その owner の runtime assertion は、下表の routing test artifact、owner 詳細本文、fixture root、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract)、[mutation test 証跡固定契約](#mutation-test-evidence-contract) のすべてへ到達できなければならない。下表の routing test artifact 以外に同 owner の assertion を追加する場合は、本項の routing 表を同時に更新する。routing 表に存在しない assertion は孤立 test として扱う。

| owner component | 実装包含 artifact | routing test artifact | fixture / 証跡接続 | drift 判定 |
|-----------------|--------------------|-----------------------|---------------------|------------|
| `statefile` | [`components/runner/`](../../components/runner/)、[`components/api/`](../../components/api/)。詳細本文は [`docs/details/statefile.md` 詳細本文責務](statefile.md) を参照する。 | [`components/runner/runner_test.go`](../../components/runner/runner_test.go)、[`components/api/api_test.go`](../../components/api/api_test.go)、[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go)。 | `testdata/statefile/`、`testdata/runner/`、`testdata/api/`、`testdata/mcp/`、[Phase 11 fixture harness 参照](#phase-11-fixture-harness-reference)。 | lock、atomic write、JSON object wrapper、direct runtime state write 禁止、state side effect order、read-only no-write、MCP state bridge の assertion が `statefile` 詳細本文または対象 fixture catalog へ到達できない場合は未解消 drift とする。 |
| `archive` | [`components/runner/`](../../components/runner/)、[`components/api/`](../../components/api/)。詳細本文は [`docs/details/archive.md` 詳細本文責務](archive.md) を参照する。 | [`components/runner/runner_test.go`](../../components/runner/runner_test.go)、[`components/api/api_test.go`](../../components/api/api_test.go)。 | `testdata/archive/`、`testdata/runner/`、`testdata/api/`、[Phase 11 fixture harness 参照](#phase-11-fixture-harness-reference)。 | snapshot 作成、保存済み archive 検証、download、delete、rollback / diff 連携、破損 archive、状態不変の assertion が `archive` 詳細本文または対象 fixture catalog へ到達できない場合は未解消 drift とする。 |
| `commitstatus` | [`components/runner/`](../../components/runner/)。詳細本文は [`docs/details/commitstatus.md` 詳細本文責務](commitstatus.md) を参照する。 | [`components/runner/runner_test.go`](../../components/runner/runner_test.go)。 | `testdata/commitstatus/`、`testdata/runner/`、[Phase 11 fixture harness 参照](#phase-11-fixture-harness-reference)。 | GitHub Commit Status request、state mapping、target URL、retry / rate-limit / failure handling、secret mask、副作用順序の assertion が `commitstatus` 詳細本文または対象 fixture catalog へ到達できない場合は未解消 drift とする。 |
| `security` | [`components/api/`](../../components/api/)、必要に応じて [`components/runner/`](../../components/runner/)、[`components/admin/`](../../components/admin/)、[`components/setup/`](../../components/setup/)、[`components/release/`](../../components/release/)、[`components/mcp/`](../../components/mcp/)、[`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js)、[`admin/index.html`](../../admin/index.html)。詳細本文は [`docs/details/security.md` 詳細本文責務](security.md) を参照する。 | [`components/api/api_test.go`](../../components/api/api_test.go)、[`components/runner/runner_test.go`](../../components/runner/runner_test.go)、[`components/admin/admin_test.go`](../../components/admin/admin_test.go)、[`components/setup/setup_test.go`](../../components/setup/setup_test.go)、[`components/release/release_test.go`](../../components/release/release_test.go)、[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go)、[`sdk_contract_test.go`](../../sdk_contract_test.go)、[`ui_contract_test.go`](../../ui_contract_test.go)。 | `testdata/security/`、`testdata/api/`、`testdata/runner/`、`testdata/admin/cli/`、`testdata/setup/`、`testdata/release/`、`testdata/mcp/`、`testdata/sdk/`、`testdata/ui/`、[Phase 11 fixture harness 参照](#phase-11-fixture-harness-reference)。 | authentication、authorization、session、TOTP、token、secret redaction、rate limit、audit / access log、credential file、external boundary、one-time secret の assertion が `security` 詳細本文または対象 fixture catalog へ到達できない場合は未解消 drift とする。 |

以下は実在する実装 artifact と owner component 詳細本文を照合した、現行実装整合の追加検証候補である。Phase 2 runner、Phase 3 api request lifecycle、Phase 4 api operations、Phase 5 sdk、Phase 6 ui、Phase 7 admin CLI、Phase 8 setup、Phase 9 release、および Phase 10 mcp の実装済み状態を取り消す一覧ではない。本表は [Phase 11 完了 closure](#phase-11-quality-gate-closure) の open item ではなく、将来計画または改訂予定へ分類する場合に使う追加検証候補だけを記録する。

| ID | 対象 | 現在確認できる実装証跡 | 追加検証候補の証跡 |
|----|------|------------------------------|----------------------------|
| <a id="align-01"></a>`ALIGN-01` | 起動入口 | [`main.go`](../../main.go) は `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci-api`、`adlaire-ci-admin`、`adlaire-ci-setup`、`adlaire-ci-release`、`adlaire-ci-mcp` を完全一致 basename で分岐し、未知 basename を拒否する。Release asset 名の `-linux-amd64` suffix は標準実行バイナリ名へ正規化し、`--version` は注入済み binary version を返す。[`components/api/`](../../components/api/) は `adlaire-ci-api` の起動契約、[`components/admin/`](../../components/admin/) は `adlaire-ci-admin` の CLI 管理契約、[`components/setup/`](../../components/setup/) は `adlaire-ci-setup` の setup / update 契約、[`components/release/`](../../components/release/) は `adlaire-ci-release` の Release 生成・公開契約、[`components/mcp/`](../../components/mcp/) は `adlaire-ci-mcp` の MCP server 契約を起動入口へ接続済みである。 | Phase 4 受け入れでは `adlaire-ci-api` CLI / lifecycle fixture を formal fixture harness へ昇格する。Phase 7 受け入れでは `adlaire-ci-admin` CLI fixture を [`testdata/admin/cli/`](../../testdata/admin/cli/) に配置済みである。Phase 8 受け入れでは `adlaire-ci-setup` fixture を [`testdata/setup/`](../../testdata/setup/) に配置済みである。Phase 9 受け入れでは `adlaire-ci-release` fixture を [`testdata/release/`](../../testdata/release/) に配置済みである。Phase 10 受け入れでは `adlaire-ci-mcp` fixture を [`testdata/mcp/`](../../testdata/mcp/) に配置済みである。 |
| <a id="align-02"></a>`ALIGN-02` | API / admin | API は [`admin/index.html`](../../admin/index.html) と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) を配信しない。 | [`docs/details/admin.md`](admin.md) の配信契約と API fixture。 |
| <a id="align-03"></a>`ALIGN-03` | API / SDK | [`components/api/`](../../components/api/) は `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject`、`GET /api/build-chain-config`、`POST /api/build-chain-config`、`GET /api/stats/build-trends` の 6 endpoint を実装済みである。`GET /api/build/stream` は保存済み log から `at` を持つ `log` frame と終端 `end` frame を有限 SSE として返す。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。SDK public method、SDK stream terminal、request / response / stream の Phase 5 実装証跡は [Phase 5 sdk 実装検証証跡](#phase-5-sdk-implementation-evidence) を参照する。一方で formal `testdata/api/operations/` / `testdata/sdk/operations/` cross fixture harness は追加検証候補として扱う。 | API / SDK cross fixture、正式 `testdata/api/operations/` と `testdata/sdk/operations/` の fixture。 |
| <a id="align-04"></a>`ALIGN-04` | API security | [`components/api/`](../../components/api/) は `.api_tokens` を top-level object `{"tokens":[]}` として読み書きし、record の `label`、`expires_at`、`last_used_at`、`revoked_at`、hash だけ保存、作成時 1 回だけ token 本体返却、token 一覧の secret 非表示、empty / unknown scope 拒否、endpoint ごとの scope 認可、期限切れ拒否、失効済み拒否、認証成功時の `last_used_at` 更新、token 作成・認証・scope 拒否・失効の `.access_log` 追記、API token event の `timestamp` / `request_id` / actor / target / result / `remote_addr` を持つ `.audit_log` 追記、必須 token log 追記失敗時の token 本体非返却、`.api_tokens` の作成・認証時 `last_used_at` 更新・失効を statefile lock 内 read-modify-write で直列化し、lock conflict を `409` として返す処理を実装済みである。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。一方で正式 fixture directory harness、access / audit log failure injection、secret 非表示 fixture は追加検証候補として扱う。 | [`docs/details/security.md`](security.md) と [`docs/details/statefile.md`](statefile.md) に一致する token 作成・認証・失効の formal fixture、access / audit log failure injection、secret 非表示 fixture。 |
| <a id="align-05"></a>`ALIGN-05` | UI / SDK / design | [`admin/index.html`](../../admin/index.html) は Phase 6 UI として、same-origin 初期化、`/api/api` 生成防止、外部 origin / path 付き override の SDK 生成前拒否、one-time token / TOTP secret / otpauth URI の generation 一致消去、copy button、`StreamHandle.done` 監視、`:focus-visible` indicator、Webhook events、output metadata、build trends、backup / restore、circuit reset、承認 approve / reject、履歴 comment / flag / tags / rollback、schedule / pipeline / build chain、token revoke、snapshot download / delete、hook delete / log、alert / tag rule 操作を SDK 経由で実装済みである。Go test 証跡は [Phase 6 ui 実装検証証跡](#phase-6-ui-implementation-evidence) を参照する。一方で正式 UI fixture directory harness と構造化 visual fixture は追加検証候補として扱う。 | `testdata/ui/login/`、`status/`、`build/`、`config/`、`secret/`、`stream/`、`operations/` の formal fixture、fake SDK transcript、DOM assertion、SDK call order、secret 非表示、[`docs/DESIGN.md`](../DESIGN.md) の selector / token / focus 固定値を照合する構造化 visual fixture。 |
| <a id="align-07"></a>`ALIGN-07` | fixture | [`testdata/admin/cli/`](../../testdata/admin/cli/) は Phase 7 Admin CLI fixture として作成済みである。[`testdata/setup/`](../../testdata/setup/) は Phase 8 setup / admin / Release asset 連動 fixture として作成済みである。[`testdata/release/`](../../testdata/release/) は Phase 9 Release fixture として作成済みである。[`testdata/mcp/`](../../testdata/mcp/) は Phase 10 MCP fixture として作成済みである。[`testdata/runner/`](../../testdata/runner/)、[`testdata/api/`](../../testdata/api/)、[`testdata/sdk/`](../../testdata/sdk/)、[`testdata/ui/`](../../testdata/ui/)、[`testdata/statefile/`](../../testdata/statefile/)、[`testdata/archive/`](../../testdata/archive/)、[`testdata/commitstatus/`](../../testdata/commitstatus/)、[`testdata/security/`](../../testdata/security/) は Phase 11 root coverage fixture として作成済みである。Admin CLI 以外の `testdata/admin/` fixture は未作成である。 | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#sec-0g-8-f) の必須 path、expected、fake、実行証跡。builder の Phase 1 証跡は [Phase 1 builder 実装検証証跡](#phase-1-builder-implementation-evidence)、Admin CLI の Phase 7 証跡は [Phase 7 admin CLI 実装検証証跡](#phase-7-admin-cli-implementation-evidence)、setup の Phase 8 証跡は [Phase 8 setup 実装検証証跡](#phase-8-setup-implementation-evidence)、release の Phase 9 証跡は [Phase 9 release 実装検証証跡](#phase-9-release-implementation-evidence)、mcp の Phase 10 証跡は [Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence) を参照する。 |
| <a id="align-09"></a>`ALIGN-09` | runner / state schema | [`components/runner/`](../../components/runner/) の build log / history は仕様外 field 構成であり、正規化 `status`、詳細 `target_status`、拡張 object、history の必須 key が [`docs/details/statefile.md`](statefile.md) の schema と一致しない。 | runner、API adapter、archive reader の同型 schema、未知 key、欠落 key、状態写像を検証する fixture。 |
| <a id="align-10"></a>`ALIGN-10` | API health / PAT | [`components/api/`](../../components/api/) は `GET /api/health` の `HealthObject`、固定順 `checks`、status summary fallback、read-only 副作用なし、`GET /api/pat-status` の期限応答、`POST /api/pat-verify` の GitHub `GET /user` 外部検証、scope 正規化、秘密情報非表示を実装済みである。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。一方で正式 `testdata/api/operations/` fixture harness、SDK / UI cross fixture、health response encode / write failure fake は追加検証候補として扱う。 | health、PAT status、PAT verify の SDK method / UI 操作接続、正式 `testdata/api/operations/` fixture、response encode / write failure fake、外部 GitHub 401 / 403 / 429 / 5xx / timeout fake、秘密情報非表示を照合する API / SDK / UI fixture。 |
| <a id="align-11"></a>`ALIGN-11` | API request / response / log schema | [`components/api/`](../../components/api/) は Phase 3 request lifecycle として request ID、`X-Request-Id`、`.api_access_log.request_id`、`.api_access_log` 必須 field、JSON response の送信前 serialize、最初の確定 status、implicit `200`、重複 `WriteHeader` 無視、`http.Flusher` 透過、SSE frame flush を実装済みである。一方で `.config_log` / `.audit_log` への request ID 伝播、`.api_access_log` 追記失敗時の固定 WARN、JSON encode / body write failure の server log、partial write / client disconnect の正式 fake fixture は追加検証候補として扱う。 | request ID、response header、最初の確定 status、implicit `200`、重複 `WriteHeader`、`http.Flusher` 透過、frame 単位 flush、response の事前 serialize、body write の全失敗 / partial write、client 切断、API access log、config diff log、partial failure、best-effort / 必須 log 追記失敗境界を検証する fixture。 |
| <a id="align-12"></a>`ALIGN-12` | queue dispatch / active state / cancel | [`components/api/`](../../components/api/) は manual / force request を `.build_state.queued[]` へ durable queue entry として保存し、保存後に runner 非同期起動要求を行う。API token trigger 時は `.build_state.queued[].requested_by` と build audit actor に token id を記録する。API は manual / force request で `running`、`current_build_id`、`last_started_at` を変更しない。API は `.build_state.active_queue_entry` を読取・返却・保持し、`DELETE /api/queue` では waiting だけを削除し、manual / force / webhook の重複判定では active と waiting を照合する。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。一方で runner の waiting-to-active atomic move、active entry の at-least-once retry / finalizer、cancel request、実行 process 停止、対象 build log の `cancelled` 最終化、history / status の正式 fixture は追加検証候補として扱う。 | [`docs/details/runner.md`](runner.md)、[`docs/details/statefile.md`](statefile.md) の atomic move、at-least-once retry、cancel request / process 停止 / `cancelled` 最終結果、runner finalizer、正式 queue / dispatch fixture。 |
| <a id="align-15"></a>`ALIGN-15` | SDK request / response / stream | [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) は Phase 5 SDK として、API 契約表と一致する public method、主要 request shape、`failure_category` query、`label` body、`URLSearchParams` 不使用の `%20` query wire encoding、header から body 読取完了までの timeout、body read error 固定変換、JSON / binary media type 検証、`StreamHandle.done` と terminal transition、fatal UTF-8 SSE parser、callback / frame error 固定化を実装済みである。詳細証跡は [Phase 5 sdk 実装検証証跡](#phase-5-sdk-implementation-evidence) と [ALIGN-15 SDK 契約解消証跡](#align-15-sdk-contract-gaps) を参照する。一方で formal SDK fixture directory harness と `expected/sdk_trace.json` / `expected/sdk_return.json` / `expected/sdk_error.json` による全 method の実行証跡は追加検証候補として扱う。 | `testdata/sdk/request-shape/`、`error-shape/`、`stream/`、`binary/`、`operations/` の formal fixture、全 public method の SDK trace / return / error expected、API / SDK / UI cross fixture。 |
| <a id="align-16"></a>`ALIGN-16` | runner config / execution | [`components/runner/`](../../components/runner/) は Phase 2 runner owner として、`.repo_config`、`.server_config`、`.branch_config`、`.notify_config` 統合、標準 builder command、`.pipeline_config` の標準 builder command 拡張と単回読取、設定 timeout、起動時設定整合性 check、branch environment、commit status、dry-run、timeout 限定 retry、multi-file、GitHub / local watch、approval pending、failure category / evidence、trigger / build status、trend、environment 証跡、build log の top-level commit fields と `remote_build` / `commit_status` / `build_meta` / `transfer_verified` 補助 key、exact builder version precheck の stderr / timeout / token 検証、remote build、tag filter、parallel deploy execution と `target_results[]`、duration anomaly、weekly summary、通知 channel schema、`.notify_log`、redirect 非追従、通知署名、SMTP 未設定分類、command 通知、priority queue の runner active 選択、build chain、dependency manifest、ビルド前後フックの runner 実行と hook log 保存を実装済みである。 | 将来計画または改訂予定では API owner に属する queue 追加 / full 判定、build chain 設定 endpoint、SDK / UI 接続、および [`docs/details/fixture.md`](fixture.md) の正式 `testdata/` fixture directory harness を整備する。Phase 2 runner 実装済み状態を取り消す根拠にしない。 |
| <a id="align-17"></a>`ALIGN-17` | statefile persistence / schema | [`components/runner/`](../../components/runner/) と [`components/api/`](../../components/api/) の JSON writer は `{name}.lock`、10 秒待機、排他的 temp 作成、file sync、mode 検証、rename 後 parent directory sync、失敗時 cleanup を共通化せず、runtime directory を `0755` で作成する。更新処理は file lock 外で read-modify-write するため、同一 target の並行更新で lost update を防止しない。JSON Lines は lock / sync なしで直接 append し、reader は破損行を証跡なしで無視する。一般 reader は未知 key と必須 key 欠落を拒否しない。`.hooks`、`.alert_rules`、`.tag_rules` は API の読み書きと backup / restore で object wrapper 保存へ整合済みである一方、共通 strict schema adapter と write 前 schema validation は追加検証候補である。UTF-8 text 契約の `.notes` は JSON object として扱う。`.api_tokens`、build log / history、`.build_state` の個別差分は [`ALIGN-04`](#align-04)、[`ALIGN-09`](#align-09)、[`ALIGN-12`](#align-12) を参照する。 | [`docs/details/statefile.md`](statefile.md) の共通 atomic write / append、lock 内 read-modify-write、1 file 1 lock、nested lock 禁止、`0700` directory、strict schema、root schema、lock timeout、sync、cleanup 契約を単一 adapter で満たし、正常系、並行更新、lost update、競合、未知 / 欠落 key、破損 JSON Lines、各 I/O failure injection を固定する statefile fixture。 |
| <a id="align-18"></a>`ALIGN-18` | API security lifecycle | [`components/api/`](../../components/api/) は request ID 生成後に path / method / access control の順で共通 gate を判定し、API token 認証では endpoint scope、expiry、`last_used_at`、token 必須 `.access_log` / `.audit_log` 追記成功前の response 禁止を実装済みである。API rate limit は strict policy schema、login pre-auth IP key、session / API token の actor key と IP key、署名検証済み Webhook の IP key、`.api_rate_state` lock 内 read-modify-write、上限超過時 count 非増加と `permission_denied` audit、policy 変更後の window 初期化、`state_summary` array を実装済みである。`.admin_credentials` は 64 文字 salt、`login_count`、`last_login_at`、`updated_at` を保存し、TOTP 無効 login と TOTP 成功時に `login_count` を飽和加算して `must_change` `none` / `prompt` / `forced` を返す。`--init-credentials` と `InitCredentials` は UTF-8、8〜128 code point、LF / CR / NUL 禁止の初期 password 制約を共有する。IP key 別 memory-only login lock は 10 回目の password 不一致で 10 分 lock し、lock 中の成功 password でも credentials を変更しない。auth transaction coordinator は login、login/totp、logout、password change、sessions、TOTP 管理 endpoint と全 authenticated endpoint の session touch / expiry update を直列化し、取得 timeout では endpoint 固有 body parse、credentials / TOTP read、entropy、memory mutation、`.access_log`、`.audit_log` を実行せず `409` を返す。credentials / TOTP / `.api_tokens` の read-modify-write は statefile lock 内で実行し、同一 credentials 更新と API token 並行更新の lost update を Go test で固定している。TOTP ticket は credentials fingerprint に結合し、TOTP 成功時に最新 credentials から再算出する。`forced` session gate と TOTP `last_accepted_step` replay 拒否も実装済みである。password login 成功、password 不一致、password login lock、TOTP required ticket 発行、TOTP setup、TOTP confirm 成功 / 失敗、TOTP login 成功 / 失敗、TOTP disable、logout、password change、session revoke-all は `.access_log` / `.audit_log` 成功後だけ one-time token / ticket / TOTP secret または成功 response を返す。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) と [`ALIGN-04`](#align-04) を参照する。一方で全 statefile の共通 atomic adapter 化、主状態保存後の追記失敗時 `500` 固定は追加検証候補として扱う。 | [`docs/details/security.md`](security.md)、[`docs/details/api.md`](api.md)、[`docs/details/statefile.md`](statefile.md) の同時 login / password change / TOTP / revoke の formal fixture、token 以外の必須 access / audit 追記成功前の ticket / secret response 禁止、主状態保存後の追記失敗における `500` と巻き戻し禁止を固定する security / API fixture。 |
| <a id="align-19"></a>`ALIGN-19` | log archive / snapshot | [`components/runner/`](../../components/runner/) は保持数超過の old build log を `.build_logs/archive/*.json.gz` へ圧縮保存してから削除する処理を持ち、Phase 2 runner snapshot として `.snapshots/{id}/site.tar.gz`、`meta.json`、checksum、USTAR、atomic publish、dot directory 除外を実装済みである。[`components/api/`](../../components/api/) は snapshot 一覧、保存済み archive download、snapshot delete と delete 成功時の `.config_log` / `.audit_log` 追記を実装済みである。 | 将来計画または改訂予定では archive owner / API owner に属する snapshot delete guard、rollback、snapshot diff、正式 fixture directory harness、各 I/O failure injection を整備する。Phase 2 runner snapshot 作成済み状態を取り消す根拠にしない。 |
| <a id="align-21"></a>`ALIGN-21` | API common HTTP validation | [`components/api/`](../../components/api/) の `decodeBody()` は 1 MiB 超過を `413 Payload too large` とし、最初の JSON value 後の trailing data を `400 Invalid JSON` として拒否する。`rejectBody()` は body 禁止 endpoint の body size 上限を扱い、history `{id}` は `^[A-Za-z0-9_-]{1,64}$` に一致しない場合 `422` を返す。整数 query helper は同一 key の重複を `422` にする。API 共通 query gate は path ごとの許可 key 以外と同一 query key 重複を `422` にし、query を持たない endpoint の query 付き request も `422` にする。一方で複数 validation detail の同時返却、全 endpoint の path / body / query 検証順序を正式 fixture として固定する証跡は追加検証候補として扱う。 | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の path / method / access / auth / body / query / path parameter 順序、未知・重複 query、複数 validation detail、状態不変を固定する API common lifecycle fixture。解消済み部分の Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。 |
| <a id="align-22"></a>`ALIGN-22` | API configuration / control | [`components/api/`](../../components/api/) の `POST /api/config/validate` は保存せずに正常・`valid=false`・warning・未知 key / secret 風 key 拒否を返す。backup は top-level `BackupObject` と必須設定の既定値を返し、restore は全必須 key、未知 key、主要 schema、secret file の保持・更新・削除、全体 no-op、固定書込順の一部、config / audit log 追記を実装済みである。`.hooks`、`.alert_rules`、`.tag_rules` は object wrapper 保存へ整合済みである。schedule interval は `.server_config` 保存後の systemd timer drop-in 作成、daemon-reload、restart、反映確認、no-op 無副作用、systemd 失敗時 partial failure log を実装済みである。一方で全 restore 対象の完全 schema validation、write failure injection、post-write partial failure、汎用設定 handler の strict schema と `.smtp_secret` 分離、SMTP test の実送信、Webhook 送信の redirect 非追従と canonical JSON byte 署名、hook / alert / tag rule record の owner schema 完全検証は追加検証候補として扱う。 | [`docs/details/api.md`](api.md)、[`docs/details/statefile.md`](statefile.md)、[`docs/details/security.md`](security.md) の backup / restore、設定 mutation、systemd、SMTP / notification、redirect 非追従、canonical payload と署名、rule / hook schema、no-op、partial failure、audit の処理順を固定する API configuration cross fixture。 |
| <a id="align-23"></a>`ALIGN-23` | API operational lifecycle | [`components/api/`](../../components/api/) は PAT verify、GitHub rate-limit、diagnostics の `pat` / `github_api` / `output_file` / `systemd` / `webhook` 固定 item、Webhook push event / branch filter、40 桁 lowercase SHA 検証、active / waiting queue の payload 重複排除、delivery 単位の衝突拒否、event log、署名検証済み Webhook queue 成功時の `build_trigger` audit、runner 非同期起動要求、manual / force build と署名・payload 検証済み webhook の maintenance fail-closed、circuit reset と maintenance mutation の config / audit log 追記を実装済みである。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。一方で正式 `testdata/api/operations/` fixture、各 I/O failure injection、SDK / UI cross fixture は追加検証候補として扱う。 | [`docs/details/api.md`](api.md)、[`docs/details/runner.md`](runner.md)、[`docs/details/security.md`](security.md)、[`docs/details/statefile.md`](statefile.md) の diagnostics、GitHub 外部確認、maintenance / circuit mutation、queue dispatch と必須副作用を fake 外部依存で固定する API operations fixture。 |
| <a id="align-24"></a>`ALIGN-24` | executable verification | [`components/runner/runner_test.go`](../../components/runner/runner_test.go)、[`components/api/api_test.go`](../../components/api/api_test.go)、[`components/admin/admin_test.go`](../../components/admin/admin_test.go)、[`components/setup/setup_test.go`](../../components/setup/setup_test.go)、[`components/release/release_test.go`](../../components/release/release_test.go)、[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go)、[`sdk_contract_test.go`](../../sdk_contract_test.go)、[`ui_contract_test.go`](../../ui_contract_test.go) は現行 Phase 実装を検証するが、[`docs/details/fixture.md`](fixture.md) の全 manifest、`input/`、`expected`、failure injection を実行する共通 harness には接続していない。[`components/admin/admin_test.go`](../../components/admin/admin_test.go) は Admin CLI fixture の必須 file 存在と JSON 妥当性を確認する。[`components/setup/setup_test.go`](../../components/setup/setup_test.go) は setup 連動 fixture の必須 file 存在と JSON 妥当性、fake download / checksum / systemd / rollback 境界を確認する。[`components/release/release_test.go`](../../components/release/release_test.go) は release fixture の必須 file 存在と JSON 妥当性、GitHub draft / upload / cleanup、checksum、dirty / version mismatch 境界を確認する。[`components/mcp/mcp_test.go`](../../components/mcp/mcp_test.go) は MCP fixture の必須 file 存在と JSON 妥当性、CLI / HTTP / JSON-RPC / tools / resources / prompts / sampling / timeout / auth / read-only / scope / confirmation / audit / metrics 境界を確認する。このため現行 test の成功だけでは正式 fixture directory harness 全体の適合を証明できない。Phase 2 `runner`、Phase 4 `api operations`、Phase 5 `sdk`、Phase 6 `ui`、Phase 7 `admin` CLI、Phase 8 `setup`、Phase 9 `release`、Phase 10 `mcp` の受け入れ証跡は、それぞれ [Phase 2 runner 実装検証証跡](#phase-2-runner-implementation-evidence)、[Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence)、[Phase 5 sdk 実装検証証跡](#phase-5-sdk-implementation-evidence)、[Phase 6 ui 実装検証証跡](#phase-6-ui-implementation-evidence)、[Phase 7 admin CLI 実装検証証跡](#phase-7-admin-cli-implementation-evidence)、[Phase 8 setup 実装検証証跡](#phase-8-setup-implementation-evidence)、[Phase 9 release 実装検証証跡](#phase-9-release-implementation-evidence)、[Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence) を正とする。 | [`docs/details/fixture.md`](fixture.md) の必須 fixture root、manifest、expected / effects、fake、証跡と [fake adapter 接続固定契約](#fixture-fake-adapter-binding-contract) を実行する Go / JavaScript 検証 harness、仕様外期待値の除去、formal fixture harness の `ALIGN-*` の正常・異常・partial / no-op / security assertion、および仕様非準拠動作を成功として受理しない regression test。builder の Phase 1 証跡と mcp の Phase 10 証跡は [Phase 1 builder 実装検証証跡](#phase-1-builder-implementation-evidence) と [Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence) を参照する。 |
| <a id="align-26"></a>`ALIGN-26` | runner boundary safety / external I/O | [`components/runner/`](../../components/runner/) は process lock の所有確認付き解放、`.repo_config` による repository 固定、GitHub owner / repo / branch / path の URL escape、rate-limit 待機、multi target SHA 解決、SSH deploy の remote path quote、`--` 境界、regular file / deploy path 検証を実装済みである。一方で Trees API の truncated / directory target response schema、fake clock を使う最大待機固定、source materialize の atomic publish、precheck no-write、SSH deploy timeout、stream copy、NUL 終端 checksum 完全一致、pending 保存失敗境界は formal fixture harness の追加検証候補として扱う。 | [`docs/details/runner.md` 詳細本文責務 lock ファイル契約](runner.md#lock-ファイル契約)、[§13](runner.md#13-処理フロー)、[§14a](runner.md#14a-ssh-サイト転送)、[`docs/details/statefile.md`](statefile.md) の owner / collaborator 契約に従い、残る GitHub schema / clock、atomic materialize、precheck no-write、SSH timeout / streaming / checksum、pending write failure を fake 外部依存で固定する runner fixture。 |
| <a id="align-27"></a>`ALIGN-27` | API read model / history consistency | [`components/api/`](../../components/api/) は `GET /api/status` の `StatusObject` 10 key、`last_target_status`、`output_url` / `queued` 非返却、`.build_status.json` 第一参照、history-only fallback 除外、sysinfo の固定 key、stats summary、timeline、build duration、dashboard、output-meta、log search、history `failure_category` filter、重複 id 除外、選択 output target の決定的参照、出力 target 不在時の固定応答、comment / flag / tags 更新の no-op 無副作用を実装済みである。Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。一方で正式 `testdata/api/operations/` fixture、archive fallback、HTML meta fallback、厳格 schema、read failure injection、複数状態の書込順、必須 config / audit log failure 境界は追加検証候補として扱う。 | [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1)、[§22.0e](api.md#sec-22-0e)、[`docs/details/statefile.md`](statefile.md)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](#sec-22-f) に従い、通常 log / archive、破損 / read failure、read-only no-write、history mutation の write order / partial failure を固定する API fixture。 |
| <a id="align-29"></a>`ALIGN-29` | output manifest / filesystem identity | [`components/runner/`](../../components/runner/) の `outputManifestSHA()` は output root symlink、symlink entry、hardlink、非通常 file、不正相対 path、読取中の size / mtime 変更を失敗として扱い、runner history / status は current target の `BranchTarget.Out` に帰属する manifest を保存する。一方で no-follow open 後の inode identity 比較、`outputSizeBytes()` の hardlink / identity 検証、走査 error の全 fixture、[`components/api/`](../../components/api/) の `inspectOutput()` / `fileTreeStats()`、runner と API の同一 manifest adapter は追加検証候補である。 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d 共通固定値](../DETAIL_INDEX.md#0d-共通固定値) の出力成果物 manifest SHA-256 契約と owner 詳細の通常 file 集計契約に従い、target ごとの output root 帰属、no-follow open、file type / link count / path / identity / size / mtime 検証、空 directory digest、runner / API 同一結果、走査中変更と走査 error の失敗を固定する cross fixture。 |
| <a id="align-30"></a>`ALIGN-30` | owner ID allocation | [`components/runner/`](../../components/runner/) の build ID は `.build_logs/{id}.json` 既存確認と `-001`〜`-999` suffix を実装済みである。[`components/api/`](../../components/api/) は queue id の既存 queue 照合と `-001`〜`-999` suffix、token id の `tok000001` 連番、hook id `h`、alert rule id `r`、tag rule id `t` の prefix を実装済みである。一方で API owner が rollback build ID を生成する残差、webhook event id の正式契約、suffix 上限到達時の no-write fixture は追加検証候補として扱う。 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d 共通固定値](../DETAIL_INDEX.md#0d-共通固定値) の時刻ベース ID 契約、[`docs/details/api.md` 詳細本文責務 §22.0e.2](api.md#sec-22-0e-2)、[`docs/details/runner.md` 詳細本文責務 build id 契約](runner.md#build-id-契約) に従い、rollback owner 境界、webhook event id、suffix 上限時 no-write を固定する API / runner cross fixture。解消済み API ID の Go test 証跡は [Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence) を参照する。 |
| <a id="align-31"></a>`ALIGN-31` | UI connected operation contract | [`admin/index.html`](../../admin/index.html) は Phase 6 UI として、notify form の schema 外 `webhooks` / `email` 送信を除去し、`channels` / `summary` を JSON 入力として扱う。repo 保存は `{owner,repo}` だけを送信し、branch config は別 request に分離する。config form は `duration_anomaly` object を送信する。UI 操作表の追加操作は SDK method だけを呼び、backup / export / snapshot は browser download を開始し、破壊的操作は確認後に実行し、`409` 後の再取得、`422` details の field focus、`429` の 10 秒操作無効化を実装済みである。Go test 証跡は [Phase 6 ui 実装検証証跡](#phase-6-ui-implementation-evidence) を参照する。一方で UI / SDK / API cross fixture directory harness は追加検証候補として扱う。 | [`docs/details/ui.md`](ui.md)、[`docs/details/api.md`](api.md)、[`docs/details/sdk.md`](sdk.md) の operation 対応、request shape、成功後再取得、status 別 UI 状態に従い、notify / repo / branch / config / approval / history / schedule / pipeline / token / snapshot / hook / rule の exact request、`409` / `422` / `429`、状態不変、再実行可否を固定する UI / SDK / API cross fixture。 |
| <a id="align-32"></a>`ALIGN-32` | runner startup integrity / secret safety | [`components/runner/`](../../components/runner/) は GitHub watch の `.github_token` について symlink / special file、mode、UTF-8、NUL / 制御文字、trim 後の内部空白を検証し、local watch では token を要求しない。一方で token 読取順序は local watch を許容するため契約との再整合が必要であり、読込成功直後の secret mask 登録、`.branch_config`、`.notify_config`、`.build_state`、`.build_circuit_state`、`.pending_transfers`、`.notify_pending`、`.maintenance` の固定順序検証、file / directory mode 確認、停止・復旧境界は追加検証候補である。 | [`docs/details/runner.md` GitHub token 読み込み契約](runner.md#github-token-読み込み契約)、[状態ファイル権限契約](runner.md#状態ファイル権限契約)、[設定ファイル起動時整合性チェック](runner.md#設定ファイル起動時整合性チェック)、[§13](runner.md#13-処理フロー) に従い、lock 前禁止副作用、token file type / mode / byte 契約、secret mask、startup 対象の順序・復旧・停止・無変更を固定する runner / security / statefile fixture。 |
| <a id="align-33"></a>`ALIGN-33` | runner terminal persistence / finalizer | [`components/runner/`](../../components/runner/) は build log 開始時保存、target 途中失敗の log / history / trend / status 保存、cooldown / circuit / schedule / no-change / allowed-hours / approval pending の `.build_status.json` 保存、所有確認付き lock 解放を実装済みである。一方で panic を含む全終了経路 finalizer、必須保存失敗後の補償、active queue entry の保持・解除境界、通常終了時の queue / active / `last_started_at` 維持、lock 解放失敗の永続記録は追加検証候補である。 | [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の終了結果保存集合、失敗境界、[runner finalizer 固定契約](runner.md#runner-finalizer-contract) に従い、正常・失敗・panic・timeout・skip、各保存失敗、queue 確定 / 未確定、所有確認付き lock 解放を固定する runner / statefile fixture。 |
| <a id="align-34"></a>`ALIGN-34` | runner target decision / output acceptance | [`components/runner/`](../../components/runner/) は `sha_file` の読取・schema 失敗を `failure_decode` とし、pipeline 成功後に output root と `index.html`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` の通常 file / symlink 非追従 / root 内包を検証し、検証成功後だけ SHA / deploy / snapshot へ進む。出力 manifest、size、REPORT 有無に依存しない output size warning も target の `BranchTarget.Out` に帰属する。 | [`docs/details/runner.md` SHA cache 読み書き契約](runner.md#sha-cache-読み書き契約)、[§13](runner.md#13-処理フロー)、[§14](runner.md#runner-build-pipeline-execution) に従い、SHA cache の厳格読取、必須出力 file set、file type / path 検証、検証成功後だけの SHA / deploy / snapshot、REPORT 有無に依存しない size 判定を formal fixture harness で固定する runner / builder cross fixture。 |
| <a id="align-35"></a>`ALIGN-35` | runner build identity / target attribution | [`components/runner/`](../../components/runner/) は実行する `BranchTarget` ごとに build id を採番し、同一 target 内 `target_files` の ID 共有、target ごとの log / history / snapshot / output SHA 帰属、全体終了コード集約を実装済みである。一方で `.build_state.current_build_id` の target 別遷移 fixture、active queue entry と組み合わせた attribution fixture は追加検証候補である。 | [`docs/details/runner.md` 詳細本文責務 build id 契約](runner.md#build-id-契約)、[§13](runner.md#13-処理フロー)、[§27.21](runner.md#sec-27-21) に従い、target ごとの current ID 遷移、queue 連動、log / history / snapshot / output SHA の非衝突と帰属、全体終了コード集約を固定する runner fixture。 |
| <a id="align-36"></a>`ALIGN-36` | runner output capture / report ingestion | [`components/runner/`](../../components/runner/) は stdout / stderr の CRLF / CR / NUL / invalid UTF-8 正規化、UTF-8 safe tail truncation、truncation flag、最初の `[REPORT]` 採用、重複 `[REPORT]` の `REPORT_DUPLICATE` warning、成功時 report 欠落の `REPORT_MISSING` warning を実装済みである。一方で secret mask、`[WARNING]` prefix 取り込み、13 key exact parse、拡張 key保持、不正 key 順 / 不足 key / 型不正の固定 warning 変換 fixture は追加検証候補である。 | [`docs/details/runner.md` 詳細本文責務 §14](runner.md#runner-build-pipeline-execution) と [§15](runner.md#15-ログ) に従い、secret mask、WARN 2 prefix、REPORT 0 / 1 / 複数件、13 key exact parse、拡張 key保持、pipeline exit code 非変更を固定する runner fixture。 |
| <a id="align-37"></a>`ALIGN-37` | binary release version contract | [`components/runner/`](../../components/runner/) は `--version` の `V.0.0-dev` default、exact 3 token、`--help` / `--version` 優先、runner と builder の exact version 一致 precheck を実装済みである。[`components/api/`](../../components/api/) と [`main.go`](../../main.go) は `adlaire-ci-api` の `--help` / `--version` 優先、exact 3 token、起動入口 dispatch を実装済みである。[`components/admin/`](../../components/admin/) と [`main.go`](../../main.go) は `adlaire-ci-admin` の `--help` / `--version` 優先、起動入口 dispatch を実装済みである。[`components/setup/`](../../components/setup/) と [`main.go`](../../main.go) は `adlaire-ci-setup` の `--help` / `--version` 優先、exact 3 token、Release asset 配置後の target version 照合を実装済みである。[`components/release/`](../../components/release/) と [`main.go`](../../main.go) は `adlaire-ci-release` の `--help` / `--version` 優先、Release tag との一致検証、`V.0.0-dev` Release 拒否、binary version 注入、Release asset の `--version` 検証を実装済みである。[`components/mcp/`](../../components/mcp/) と [`main.go`](../../main.go) は `adlaire-ci-mcp` の `--help` / `--version` 優先、exact 3 token、起動入口 dispatch を実装済みである。[`components/obsidian/`](../../components/obsidian/) と [`main.go`](../../main.go) は `adlaire-ci-obsidian` の `--help` / `--version` 優先、exact 3 token、起動入口 dispatch を実装済みである。 | [`docs/SPEC.md` ポリシー責務 §1](../SPEC.md#policy-versioning)、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](../DETAIL_INDEX.md#common-cli-contract)、runner / API / setup / release / admin / mcp / obsidian 詳細本文に従い、標準実行バイナリの開発 default、Release 注入、exact 3 token、tag 一致、`V.0.0-dev` Release 拒否、runner-builder 一致、不正出力の無副作用失敗を固定する cross fixture。builder の Phase 1 証跡、API CLI 証跡、Admin CLI 証跡、setup 証跡、release 証跡、mcp 証跡、Obsidian CLI 証跡は [Phase 1 builder 実装検証証跡](#phase-1-builder-implementation-evidence)、[Phase 4 api operations 実装検証証跡](#phase-4-api-operations-implementation-evidence)、[Phase 7 admin CLI 実装検証証跡](#phase-7-admin-cli-implementation-evidence)、[Phase 8 setup 実装検証証跡](#phase-8-setup-implementation-evidence)、[Phase 9 release 実装検証証跡](#phase-9-release-implementation-evidence)、[Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence)、[Phase 15 Obsidian local vault 同期証跡](#phase-15-obsidian-local-sync-evidence) を参照する。`adlaire-ci-release` を利用者向け Release asset に含めないことは、同バイナリの version 契約を免除しない。利用者向け Release asset は [`docs/details/release.md` 詳細本文責務 §R3](release.md#release-asset-contract) の 7 実行バイナリと管理 UI archive、checksum manifest に限定されるが、標準実行バイナリ名全体の version 契約は 8 件すべてへ適用する。 |

<a id="align-05-unconnected-sdk-operations"></a>
**ALIGN-05 UI 操作接続証跡：**

| 操作群 | Phase 6 接続状態 |
|--------|------------------------------------------------------------|
| 外部死活監視 | `health` は [`docs/details/ui.md`](ui.md) 詳細本文責務で標準管理 UI が呼び出さない method として定義済みであり、API / SDK fixture で扱う。 |
| 状態・診断 | `resetCircuitBreaker`、`getOutputMeta`、`getBuildTrends`、`getWebhookEvents` は Phase 6 UI で SDK 境界を通じて接続済みである。 |
| 承認・ビルド連鎖 | `approveBuild`、`rejectBuild`、`getBuildChainConfig`、`setBuildChainConfig` は Phase 6 UI で SDK 境界を通じて接続済みである。 |
| schedule・pipeline | `setAllowedHours`、`clearAllowedHours`、`setForceInterval`、`setBuildCooldown`、`setScheduleInterval`、`setPipelineConfig` は Phase 6 UI で接続済みである。 |
| history | `getHistoryComment`、`setHistoryComment`、`setHistoryFlag`、`setHistoryTags`、`rollbackHistory` は Phase 6 UI で接続済みである。 |
| snapshot | `deleteSnapshot`、`downloadSnapshot` は Phase 6 UI で接続済みである。 |
| token | `revokeToken` は Phase 6 UI で接続済みである。 |
| hook・rule | `getHookLog`、`deleteHook`、`addAlertRule`、`deleteAlertRule`、`addTagRule`、`deleteTagRule` は Phase 6 UI で接続済みである。 |

<a id="align-15-sdk-contract-gaps"></a>
**ALIGN-15 SDK 契約解消証跡：**

| 対象 | Phase 5 実装証跡 |
|------|------------------|
| `getAccessLog()` / `getNotifyLog()` / `getConfigLog()` | `{limit=100,offset=0}` を受け取り、`limit` / `offset` query を送信する。 |
| `getHistory()` | `failureCategory` を受け取り、`failure_category` query を送信する。 |
| `setRepoConfig()` | `{owner,repo}` の固定引数契約だけを受け取り、定義済み key だけを body に含める。 |
| `createToken()` | `{label,scopes,expires_at}` body を送信し、`name` key を送信しない。 |
| request validation | `_validateId()` は [`docs/details/api.md` 詳細本文責務 §22.0b](api.md#sec-22-0b) の `^[A-Za-z0-9_-]{1,64}$` に一致しない id を送信前に拒否する。主要変更 method は object / string / array / number / boolean の送信前検証を持つ。 |
| `streamBuild()` | `Content-Type: text/event-stream`、fatal UTF-8、frame の exact key / type / `at`、最終 `end`、EOF 後 buffer、`done` Promise、normal / error / user close の一回だけの terminal transition を実装する。 |
| JSON response | `Content-Type` を media type と parameter に分解し、`application/json` と `application/*+json` だけを JSON として扱う。 |
| query wire encoding | `_query()` は `URLSearchParams` を使わず、`encodeURIComponent(String(value))` で空白を `%20` として送信する。 |
| timeout / body read | `_request()` は `fetch()`、JSON text、Blob 読取の全体を 30 秒 timeout で保護し、timeout を `Request timeout`、その他の body 読取 reject を `Network error` へ固定変換する。 |
| binary response | `downloadSnapshot(id)` の `2xx` は `Content-Type: application/octet-stream` の場合だけ Blob として受理し、それ以外を `Invalid binary response` とする。 |

<a id="0g8-f-fixture--testdata--fake--実装検証証跡契約"></a>

**0g.8-F fixture / testdata / fake / 実装検証証跡契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装完了判定に必要な fixture、fake、testdata、expected / effects、実装検証証跡、acceptance checklist、差し戻し条件を扱う。[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0e](../DETAIL_INDEX.md#0e-完全実装検証マトリクス) と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) は詳細本文と証跡への入口、実装割当・順序・依存は [`docs/ROADMAP.md` 状態・計画責務 §4](../ROADMAP.md#roadmap-phase-plan)、実装変更単位と着手条件は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、setup / update の実行条件とRelease asset受け入れ条件は[`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順)、Release生成・公開契約は[`docs/details/release.md`](release.md)を参照する。fixture名、expected / effects、fake動作、実装検証証跡項目、不足時の扱い、差し戻し条件だけを[`docs/details/fixture.md`](fixture.md) fixture証跡責務で固定する。

実装検証証跡は、対象に応じて以下の 3 系統に分類する。複数系統にまたがる変更は、該当する全系統の証跡を実装検証証跡として記録する。

| 系統 | 対象 | 責務節 | 必須証跡 |
|------|------|--------|----------|
| component 実装 | builder、runner、api、admin、sdk、ui、statefile、archive、commitstatus、security、setup、release、mcp の実装。 | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) | owner component、collaborator component、変更ファイル、fixture / testdata path、fake、実行コマンド、期待結果、実結果、依存 component へ引き継ぐ contract、[mutation test 証跡固定契約](#mutation-test-evidence-contract) の対象有無。 |
| API 横断実装 | API と同期する SDK / UI / statefile の実装。 | [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) | endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、endpoint 固有の業務状態非変更、共通 security / observability 副作用、secret mask、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の対応。 |
| [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の追加仕様化機能 | runner / security 詳細本文責務。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | 対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様).x / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、関連 [`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) / [`docs/details/api.md` 詳細本文責務 §25](api.md#25-認証-実装仕様) / [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) / [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) / [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順)、owner / collaborator component、fixture 名、状態差分、外部副作用、partial failure、再実行、対象外確認、mutation test 証跡。 |

**不足時共通扱い：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務で必須とする fixture、manifest、expected、effects、security、実装検証証跡、対象外確認のいずれかが不足する場合、対象機能は未完了として扱う。fixture の pass だけでは完了証跡を満たさない。[`docs/details/fixture.md`](fixture.md) fixture 証跡責務内の対象別不足時表は、この不足時共通扱いに対する具体条件である。

component、API、[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) のいずれの実装検証証跡でも、記録形式は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務の表に従う。owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務、[`docs/details/setup.md`](setup.md) 詳細本文責務に同種の記録項目がある場合でも、[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は証跡分類、不足時の扱い、差し戻し条件だけを固定する。

<a id="test-evidence-package-record-location-contract"></a>
**test evidence package 記録先固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装変更、検証変更、fixture 変更、Phase 全体完了判定で提出する test evidence package の記録単位と到達条件だけを固定する。テスト方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

test evidence package は、単一の対象変更単位ごとに 1 組だけ作成する。対象変更単位は、1 つの実装 Pull Request、1 つの検証 Pull Request、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ test evidence package に複数 Phase、無関係な owner、後続 Phase、または完了判定しない将来作業を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| `evidence_package_id` | `phase-<number>-<scope>` または `<owner>-<feature>-<scope>` の lowercase kebab-case とし、同一 Pull Request 内で一意にする。 | id がない、空、重複、line number、commit hash だけ、または対象 owner / Phase を識別できない。 |
| `record_location` | Pull Request の `Verification` に記録する場合は `pr-verification`、fixture / expected / testdata に記録する場合は repository root 起点の path、文書内に記録する場合は責務名付き Markdown link を記録する。 | 記録先がない、口頭説明だけ、または実在しない path / anchor を記録している。 |
| `scope` | 対象 Phase、対象 owner、collaborator、対象 artifact、対象 test artifact、対象 fixture root、対象仕様 anchor、対象外 artifact を列挙する。 | 対象 owner、artifact、fixture root、仕様 anchor のいずれかが不明である。 |
| `required_evidence_sets` | 適用する test gap inventory、test improvement batch closure、coverage ledger、closure record set、oracle、failure diagnostics、boundary / failure matrix、isolation、determinism、race trigger matrix、concurrency / race、mutation selection ledger、mutation、harness self-verification、contract drift を列挙する。 | 適用判断がない、または適用対象の evidence set が記録先へ到達できない。 |
| `not_applicable_evidence_sets` | 適用外の evidence set は、対象外範囲、理由、責務正本 anchor、完了可否への影響を同じ package に記録する。 | 理由だけ、anchor だけ、将来対応、または対象外範囲不明で適用外にしている。 |
| `final_open_item_count` | 全 evidence set と closure record set の open item 合計を整数で記録し、完了扱いでは `0` に固定する。 | 件数未記録、`0` 以外、または残件を別変更で閉じる説明がある。 |

test evidence package は、[implementation PR evidence template 固定契約](#implementation-pr-evidence-template-contract)、[test gap inventory record 固定契約](#test-gap-inventory-record-contract)、[test improvement batch closure 固定契約](#test-improvement-batch-closure-contract)、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract)、[test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract)、該当 evidence set、[test / contract drift report schema 固定契約](#test-contract-drift-report-schema-contract) のいずれにも到達できない場合、完了証跡として扱わない。

<a id="test-gap-inventory-record-contract"></a>
**test gap inventory record 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、テスト関連改善作業で発見した問題点の棚卸し record schema だけを固定する。全件棚卸し義務、open item 残存時の完了禁止、対象外判断の方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を正本とする。

test gap inventory record は、発見した問題点 1 件につき 1 record 作成する。同じ原因に見える問題でも、検出元、対象 artifact、対象 owner、fixture root、または必要対応が異なる場合は別 record とする。重複統合する場合も、統合前の検出元を record 内にすべて残す。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `gap_id` | `gap-<owner>-<scope>-<number>` の lowercase kebab-case。 | 空、重複、または対象 owner / scope を識別できない。 |
| `source` | 発見元を `source-audit`、`test-drift`、`fixture-drift`、`mutation-survivor`、`race-trigger`、`review`、`validation`、`spec-gap` のいずれかで記録する。`spec-gap` は、テストまたは fixture の完了可否に影響する仕様不足だけに使用する。 | 発見元がない、自由記述だけで分類できない、またはテスト固有ではない仕様全般不備を `spec-gap` として混在している。 |
| `owner_component` | 対象 owner component。テスト固有の仕様不足で owner component が複数に見える場合も、検証不足を閉じる主 owner を 1 件だけ記録し、残りは collaborator として `spec_anchor` または `evidence_target` から到達させる。 | owner 未記録、または [docs/DETAIL_INDEX.md 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) に存在しない owner。 |
| `artifact` | 対象実装 artifact、test artifact、fixture、expected、fake、harness、checker、または文書 anchor。 | artifact が不明、または実在所在 / anchor へ到達できない。 |
| `spec_anchor` | 問題点を判定する責務正本への Markdown link。 | anchor なし、または説明文だけで仕様判断している。 |
| `spec_deficiency_record` | `source=spec-gap` の場合は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](../SPEC.md#spec-deficiency-inventory-record-contract) の `defect_id`。`source=spec-gap` 以外は `not_applicable`。 | `source=spec-gap` なのに `defect_id` がない、または仕様全般不備 inventory record と接続せず test gap だけで仕様不足を閉じている。 |
| `gap_type` | `missing-test`、`weak-oracle`、`missing-assertion-id`、`missing-fixture`、`fixture-duplicate`、`contract-drift`、`mutation-survived`、`race-unverified`、`non-deterministic`、`skip-without-anchor`、`unknown-side-effect`、`spec-missing` のいずれか。 | 分類なし、または複数分類を 1 record に混在している。 |
| `fixture_root` | 対象 fixture root、または `not_applicable` と対象外 anchor。 | fixture root が必要なのに空、または対象外理由がない。 |
| `evidence_target` | 作成または更新すべき evidence set、ledger、matrix、closure item。 | 対応する証跡種別が不明で closure へ接続できない。 |
| `required_action` | `specify`、`add-test`、`strengthen-oracle`、`add-fixture`、`dedupe-fixture`、`connect-anchor`、`add-mutation`、`add-race-evidence`、`mark-not-applicable` のいずれか。 | 必要対応が自由記述だけ、または実装判断で補完している。 |
| `closure_record` | 接続先 closure record、batch closure item、または Pull Request evidence label。 | closure へ到達できない。 |
| `status` | `open`、`closed`、`not_applicable` のいずれか。 | 完了時に `open` が残る、または `not_applicable` に責務正本 anchor がない。 |

test gap inventory record は、status が `closed` または `not_applicable` であり、`spec_anchor`、`evidence_target`、`closure_record` へ到達できる場合だけ閉じる。`source=spec-gap` の record は、`spec_deficiency_record` から [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](../SPEC.md#spec-deficiency-inventory-record-contract) へ到達でき、かつ [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](../SPEC.md#spec-deficiency-batch-closure-contract) の `test_gap_connection` で当該 `defect_id` と `gap_id` の 1 対 1 対応が確認できる場合だけ閉じる。`status=open` の record、分類不能 record、対象外理由 anchor のない record、または closure へ接続しない record が 1 件でも残る場合、test improvement batch は完了扱いにしない。

<a id="test-improvement-batch-closure-contract"></a>
**test improvement batch closure 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、テスト関連改善作業を一括で閉じる closure 条件だけを固定する。作業を一括 closure まで完了させる方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

test improvement batch closure は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの検証改善 Pull Request、1 つの fixture / expected 改善、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ batch closure に複数 Phase、無関係な owner、または完了判定しない後続作業を混在させてはならない。

| closure item | 必須記録 | 未完了条件 |
|--------------|----------|------------|
| batch identity | `batch_id`、対象 Phase、対象 owner、対象 artifact、対象 fixture root、対象 Pull Request。 | batch 範囲が不明、または複数変更単位を混在している。 |
| gap inventory | [test gap inventory record 固定契約](#test-gap-inventory-record-contract) の所在、record 件数、`open=0`、`not_applicable` 件数と anchor。 | inventory なし、record 件数未記録、または `open>0`。 |
| requirement ledger | [test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract) の所在、未検証契約 `0`。 | ledger なし、未検証契約残存、または対象外理由 anchor 不足。 |
| oracle / diagnostics / boundary | oracle、failure diagnostics、boundary / failure matrix の各 evidence set 所在、適用外理由。 | 弱い oracle、generic failure、境界未固定を残している。 |
| isolation / determinism | isolation、determinism の evidence set 所在、順序依存と変動要因の解消状態。 | 順序依存、flaky、実環境依存、cleanup 不明を残している。 |
| concurrency / race | [race trigger matrix 固定契約](#race-trigger-matrix-contract) と [test concurrency / race evidence set 固定契約](#test-concurrency-race-evidence-set-contract) の所在、open trigger `0`。 | race trigger 未判定、race detector / interleaving 証跡不足、または open trigger 残存。 |
| mutation | [mutation selection ledger 固定契約](#mutation-selection-ledger-contract) と [mutation test evidence set 固定契約](#mutation-test-evidence-set-contract) の所在、`survived=0`。 | mutation class 未判定、`survived>0`、または invalid / equivalent 根拠不足。 |
| harness self-verification | [test harness self-verification evidence set 固定契約](#test-harness-self-verification-evidence-set-contract) の所在、negative / positive control。 | harness が常時 pass / 常時 fail の可能性を排除できない。 |
| contract drift | [test / contract drift report schema 固定契約](#test-contract-drift-report-schema-contract) の所在、孤立 test `0`、未検証契約 `0`、期待値ドリフト `0`。 | drift 残存、または解消先 anchor 不足。 |
| fixture root closure | [fixture root coverage matrix 固定契約](#fixture-root-coverage-matrix-contract) と [未作成 fixture root closure record 固定契約](#fixture-root-missing-closure-record-contract) の所在。 | 未作成 root、重複 root、未接続 root を残している。 |
| final closure | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) への接続、`final_open_item_count=0`。 | final open item が 0 でない、または closure record set へ接続しない。 |

test improvement batch closure は、上表の closure item を表順に記録する。適用外の closure item は削除せず、対象外範囲、理由、責務正本 anchor、完了可否への影響を記録する。最終 record より後に残件、暫定対応、後続 PR 前提、または未確認事項を追記した batch closure は完了証跡として扱わない。

<a id="mutation-selection-ledger-contract"></a>
**mutation selection ledger 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、変更 artifact ごとの mutation class 選定記録だけを固定する。mutation class、判定語彙、必須条件は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、具体的な mutation 実行証跡は [mutation test evidence set 固定契約](#mutation-test-evidence-set-contract) を参照する。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `mutation_id` | `mut-<owner>-<artifact>-<number>` の lowercase kebab-case。 | 空、重複、または対象 artifact を識別できない。 |
| `owner_component` | 対象 owner component。 | owner 未記録、または対象変更単位と一致しない。 |
| `artifact` | 変異対象または対象外判断した実装、test、fixture、expected、security expected、state diff、harness、checker。 | 変更 artifact が ledger へ接続していない。 |
| `spec_anchor` | 変異で検出すべき仕様違反の責務正本 anchor。 | anchor なし、または実装挙動だけで判定している。 |
| `mutation_class` | [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の mutation class。 | class 未判定、または任意 class だけで標準 class を飛ばしている。 |
| `operation` | mutation operation、または `not_applicable`。 | applicable なのに operation がない。 |
| `applicability` | `applicable`、`not_applicable` のいずれか。 | 判定漏れ、対象外理由だけで anchor がない。 |
| `evidence_record` | [mutation test evidence set 固定契約](#mutation-test-evidence-set-contract) の record への参照、または対象外理由 anchor。 | evidence set へ到達できない。 |
| `result` | `killed`、`invalid`、`equivalent`、`not_applicable` のいずれか。`survived` は closure 前の未完了状態としてだけ記録できる。 | 完了時に `survived` が残る、または invalid / equivalent の anchor がない。 |

mutation selection ledger は、対象変更単位の全変更 artifact を少なくとも 1 件の `mutation_id` へ接続する。完了扱いでは `survived=0` を必須とし、`invalid` と `equivalent` は個別 record に責務正本 anchor と理由を持つ。

<a id="race-trigger-matrix-contract"></a>
**race trigger matrix 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、並行処理 / race 検証を要する trigger の判定表だけを固定する。data race、deadlock、goroutine leak、lost update、cleanup 漏れを完了不可とする方針は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、具体的な証跡 set は [test concurrency / race evidence set 固定契約](#test-concurrency-race-evidence-set-contract) を参照する。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `trigger_id` | `race-<owner>-<artifact>-<number>` の lowercase kebab-case。 | 空、重複、または対象 artifact を識別できない。 |
| `owner_component` | 対象 owner component。 | owner 未記録、または対象変更単位と一致しない。 |
| `artifact` | trigger を持つ実装、test、fixture、fake、harness、checker。 | artifact が不明、または実在所在へ到達できない。 |
| `trigger_kind` | `goroutine`、`channel`、`worker`、`lock`、`listener`、`timer`、`file-lock`、`queue`、`shutdown`、`shared-state`、`parallel-request`、`same-time-event`、`not_applicable` のいずれか。 | trigger 未判定、または自由記述だけで分類できない。 |
| `shared_resource` | 共有 resource、状態 path、listener、queue、lock、channel、file、clock、fake、または `none`。 | 競合対象が不明で evidence set へ接続できない。 |
| `required_evidence` | race detector、interleaving case、lifecycle assertion、cleanup assertion、atomicity assertion のうち必要な証跡。 | 必須証跡の判定漏れ。 |
| `race_detector` | `required`、`not_required`、`unavailable-open`、`not_applicable` のいずれか。 | race detector 対象なのに未実行理由だけで閉じている。 |
| `interleaving_case` | 固定 schedule、並行 request、shutdown 中操作、lost update case、または対象外理由 anchor。 | interleaving 未固定、または代替証跡がない。 |
| `status` | `open`、`closed`、`not_applicable` のいずれか。 | 完了時に `open` が残る、または `not_applicable` に責務正本 anchor がない。 |

race trigger matrix は、変更 artifact が並行処理へ影響しない場合でも `not_applicable` record で対象外理由 anchor を残す。trigger が 1 件でも `open` の場合、concurrency / race closure と test improvement batch closure は完了扱いにしない。

<a id="test-assertion-id-contract"></a>
**assertion id 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、test、subtest、contract assertion、fixture assertion、expected 比較、security assertion、state diff assertion を安定識別する assertion id の形式だけを固定する。assertion id がない test の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、failure diagnostics evidence set は [test assertion identity / failure diagnostics evidence set 固定契約](#test-assertion-failure-diagnostics-evidence-set-contract) を参照する。

assertion id は `owner.feature.case.assertion` を最小 4 segment とする lowercase dot notation に固定する。各 segment は `^[a-z][a-z0-9-]*$` とし、owner segment は `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`release`、`security`、`mcp` のいずれかにする。必要な場合だけ 5 segment 目以降に field、boundary、mutation class を追加できる。最大 segment 数は 7 とする。

| 対象 | assertion id の条件 | 禁止形式 |
|------|---------------------|----------|
| Go test / subtest | 対象 owner、feature、case、assertion を id に含め、subtest 名と 1 対 1 で対応できる。 | line number、`TestX/1`、実行順、table index、乱数、timestamp、commit hash。 |
| fixture assertion | `manifest.json.assertions` の値、expected file、比較対象 field、対象仕様 anchor へ接続する。 | fixture 名だけ、expected file 名だけ、`state` や `response` だけの汎用名。 |
| contract assertion | API endpoint、SDK method、UI operation、CLI command、state path、MCP method のいずれかを case または assertion に含める。 | HTTP status だけ、method 名だけ、DOM selector だけ、snapshot 名だけ。 |
| mutation assertion | mutation class、対象 file、期待 failure、実際の failure と接続する。 | mutation class だけ、対象 file だけ、`killed` 件数だけ。 |

assertion id を変更する場合は、変更前 id、変更後 id、変更理由、対象 requirement、対象 evidence set、contract drift への影響を同じ test evidence package に記録する。assertion id の rename だけで検証要求を covered と扱ってはならない。

<a id="test-execution-evidence-matrix-contract"></a>
**test execution evidence matrix 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装検証結果をどの証跡分類で記録し、何が不足すると未完了になるかだけを固定する。テスト方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、実行手順と Pull Request への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations)、実在 test artifact と fixture root の所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](../DOCUMENT_INDEX.md#実装ファイル一覧) を参照する。

実装検証証跡では、実行した command、対象 artifact、対象 owner、対象 fixture、終了 code、重要 stdout / stderr の要約、期待結果、実結果、未実行理由、完了可否への影響を記録する。必須検証を未実行にする場合は、対象外理由、再実行条件、未完了扱いかどうかを責務正本 anchor へ到達できる形で記録する。実行成功だけを記録し、対象 artifact、対象 fixture、assertion、mutation、contract drift との接続を記録しない証跡は完了証跡として扱わない。

| 証跡分類 | 対象 | 必須証跡 | 未完了条件 |
|----------|------|----------|------------|
| Go 構文・単体・contract 検証 | `main.go`、`components/<owner>/*.go`、`components/<owner>/*_test.go`、Go contract test。 | 対象 Go file、対象 test artifact、owner / collaborator、fixture root、実行 command、終了 code、差分有無、pass / fail、scope を限定した理由。 | Go 実装または Go test を変更したのに、対象 artifact と owner へ接続された Go 検証結果がない。 |
| JavaScript / UI 静的検証 | `admin/adlaire-ci-sdk.js`、`admin/index.html`、SDK / UI contract。 | 対象 artifact、Deno stable runtime での検証結果、SDK / UI owner、関連 API / security collaborator、実行 command、終了 code、pass / fail。 | JavaScript 系 artifact または UI contract を変更したのに、Deno 検証証跡がない、または Node.js / npm / bundler を標準検証の代替として扱っている。 |
| fixture schema / manifest 検証 | `testdata/`、`manifest.json`、`expected/`、`effects`、`security expected`。 | 対象 fixture root、fixture 名、directory 名、`manifest.json.name`、catalog 名、schema 検証結果、expected / effects / security expected の照合結果。 | manifest、directory、catalog、expected の対応が閉じていない、または [fixture root coverage matrix 固定契約](#fixture-root-coverage-matrix-contract) の未完了条件が残る。 |
| test requirement coverage ledger 検証 | owner 詳細本文の検証条件、fixture 証跡条件、Phase 対象、test / fixture / expected / assertion / mutation の接続。 | 対象 owner、対象 anchor、検証要求 id、test artifact、fixture root、expected、assertion、mutation class、closure item、covered / not_applicable / open。 | [test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract) の未完了条件が残る。 |
| test oracle 検証 | expected / actual 比較、negative case、状態差分、effects、security expected、禁止出力、禁止外部通信、失敗時 no mutation。 | 対象 owner、対象 anchor、入力、expected file、actual 取得元、比較単位、positive / negative case、禁止副作用、完了可否。 | [test oracle evidence set 固定契約](#test-oracle-evidence-set-contract) の未完了条件が残る。 |
| test assertion identity / failure diagnostics 検証 | assertion id、対象仕様 anchor、expected / actual / diff、failure reason、再現条件、secret-safe diagnostics。 | 対象 owner、対象 anchor、test / fixture / assertion、期待値、実値、差分、failure reason、再現 command、秘密情報非露出、完了可否。 | [test assertion identity / failure diagnostics evidence set 固定契約](#test-assertion-failure-diagnostics-evidence-set-contract) の未完了条件が残る。 |
| test boundary / failure matrix 検証 | 入力 class、limit、error taxonomy、partial failure、rollback、cleanup failure、retry / recovery、失敗時 no mutation。 | 対象 owner、対象 anchor、入力 class、境界値、失敗注入、期待 error、状態差分、許可副作用、禁止副作用、完了可否。 | [test boundary / failure matrix evidence set 固定契約](#test-boundary-failure-matrix-evidence-set-contract) の未完了条件が残る。 |
| test isolation 検証 | test order、共有状態、fixture / expected mutation、環境変数、working directory、temp root、state dir、listener、goroutine、process、timer、file lock、network fake。 | 対象 owner、対象 test / fixture、隔離境界、順序入替結果、共有状態初期化、cleanup 結果、残留 resource、parallel 可否。 | [test isolation evidence set 固定契約](#test-isolation-evidence-set-contract) の未完了条件が残る。 |
| test / contract drift 検証 | test、contract test、fixture assertion、owner 詳細本文の対応。 | test artifact、assertion 名または fixture 名、owner 詳細本文 anchor、fixture 証跡 anchor、drift 判定結果。 | [test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の孤立 test、未検証契約、期待値ドリフト、harness ドリフトが残る。 |
| test determinism 検証 | clock、timer、entropy、file order、map order、filesystem、network、GitHub API、systemd、process、browser runtime、parallel worker、retry 境界。 | 対象 artifact、対象 owner、変動要因、fake adapter、同一入力再実行結果、順序入替 case、禁止実環境依存、完了可否。 | [test determinism evidence set 固定契約](#test-determinism-evidence-set-contract) の未完了条件が残る。 |
| test concurrency / race 検証 | goroutine、channel、worker、lock、listener、timer、file lock、queue、shutdown、共有状態、並行 request、同時刻 event。 | 対象 artifact、対象 owner、共有 resource、race detector 結果、schedule / interleaving case、lock / channel / goroutine lifecycle、conflict outcome、atomicity、cleanup 結果。 | [test concurrency / race evidence set 固定契約](#test-concurrency-race-evidence-set-contract) の未完了条件が残る。 |
| mutation test 検証 | 実装コード、test harness、fixture assertion、expected 比較、security assertion、state diff assertion。 | mutation class、対象 file、対象 fixture、実行単位、判定、`killed` / `survived` / `invalid` / `equivalent` 件数。 | [mutation test 証跡固定契約](#mutation-test-evidence-contract) の `survived` が 1 件以上ある、または対象 mutation が未定義。 |
| cross-owner contract 検証 | API / SDK / UI、CLI / API、statefile / archive / security / runner、setup / release / admin の横断境界。 | 呼び出し元 owner、呼び出し先 owner、endpoint / method / command / state path、状態差分、security effect、成功後再取得、失敗時 no mutation、関連 fixture。 | 片側の契約だけを検証している、または collaborator の副作用、security、状態差分、失敗時固定が未確認。 |
| 未実行・対象外証跡 | 必須検証を実行できない場合、または仕様上対象外とする場合。 | 未実行 command、未実行理由、影響 owner、影響 fixture、再実行条件、対象外にする責務正本 anchor、完了可否への影響。 | 必須検証の未実行理由がない、対象外 anchor がない、または未実行のまま完了扱いにしている。 |

<a id="test-skip-evidence-contract"></a>
**skip / 未実行証跡固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、skip、未実行、環境機能不足、対象外判断を完了判定で扱うための証跡 schema だけを固定する。skip を成功として扱うことの禁止、対象外理由 anchor の必須条件は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を参照する。

skip / 未実行 record は、必須検証を実行しない command、test、subtest、fixture、assertion、evidence set ごとに 1 件作成する。複数の未実行理由を 1 record へまとめてはならない。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `skip_id` | `skip.<owner>.<feature>.<case>` の lowercase dot notation。 | id がない、重複、対象 owner 不明。 |
| `skip_type` | `environment_capability`、`tool_unavailable`、`runtime_unavailable`、`unsupported_platform`、`spec_not_applicable`、`blocked_by_open_item` のいずれか。 | 未登録値、または実行失敗を skip として分類している。 |
| `target` | 未実行 command、test、subtest、fixture、assertion、evidence set、または mutation class。 | 対象が空、または複数対象を 1 record に混在させている。 |
| `owner_component` | 対象 owner component 名。 | 対象 owner 不明。 |
| `requirement_refs` | 対象 owner 詳細本文、fixture 証跡、[`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md) の責務名付き Markdown link 配列。 | anchor なし、裸のファイル名、または責務正本へ到達できない。 |
| `reason` | 固定理由。環境機能不足の場合は不足した capability、確認方法、再実行条件を含める。 | `環境都合`、`手元で不可`、`不要` だけの説明。 |
| `substitute_evidence_refs` | 代替証跡がある場合は evidence package、test、fixture、expected、manual でない実行証跡への link。代替なしは空配列。 | 代替ありと書いているが参照先がない。 |
| `completion_impact` | `closed_by_alternative`、`not_applicable`、`open` のいずれか。 | 必須検証を未実行なのに `closed_by_alternative` または `not_applicable` の根拠がない。 |

`completion_impact=open` の skip / 未実行 record が 1 件でも残る場合、test evidence package、closure record set、Phase 全体完了、`実装済み` 判定を完了扱いにしてはならない。`completion_impact=not_applicable` は責務正本 anchor がある場合だけ使用できる。`completion_impact=closed_by_alternative` は代替証跡が同じ requirement を covered にできる場合だけ使用できる。

<a id="test-verification-closure-checklist-contract"></a>
**test verification closure checklist 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装検証証跡を完了扱いにする直前のクロージャ項目だけを固定する。テスト方針、完了可否、mutation test 必須条件は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Phase 11 仕様全般完了は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、現在状態と対象外理由は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、実在 test artifact と fixture root の所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](../DOCUMENT_INDEX.md#実装ファイル一覧) を参照する。

実装検証証跡のクロージャ記録は、下表の全項目を `closed` または `not_applicable` として根拠 anchor 付きで記録する。`not_applicable` は、対象機能、対象 owner、対象 fixture root、対象外理由の正本 anchor がある場合だけ使用できる。`open`、根拠 anchor なし、実行結果だけ、または pass 件数だけの記録は、完了証跡として扱わない。

| クロージャ項目 | 閉じる条件 | 未完了条件 |
|----------------|------------|------------|
| scope inventory | 対象 owner、collaborator、実装 artifact、test artifact、fixture root、owner 詳細本文 anchor、fixture 証跡 anchor、状態・計画 anchor が 1 件以上記録されている。 | 対象 artifact、対象 owner、または参照 anchor のいずれかが未記録。 |
| test gap inventory closure | 対象変更単位で発見したテスト関連問題点が [test gap inventory record 固定契約](#test-gap-inventory-record-contract) に従い、全 record が `closed` または `not_applicable` であり、`open=0` と closure record への接続を持つ。 | inventory record なし、発見元未分類、対象 owner / artifact / fixture root 未接続、`status=open`、対象外理由 anchor 不足、または closure record へ到達できない gap が残る。 |
| test improvement batch closure | 対象変更単位の一括改善が [test improvement batch closure 固定契約](#test-improvement-batch-closure-contract) に従い、gap inventory、requirement coverage、execution、oracle、failure diagnostics、boundary、isolation、determinism、concurrency / race、mutation、harness self-verification、contract drift、final open item count を同じ closure record set へ接続している。 | batch closure なし、必須 closure item 未接続、個別 gap の散発修正、対象外理由 anchor 不足、または `final_open_item_count` が `0` ではない。 |
| traceability closure | 変更または対象にした test、subtest、contract assertion、fixture assertion が [test artifact traceability 固定契約](#test-artifact-traceability-contract) へ接続されている。 | 孤立 test、孤立 assertion、仕様に存在しない期待値が残る。 |
| requirement coverage closure | 対象変更に関わる owner 詳細本文の検証条件、fixture 証跡条件、Phase 対象、既存 test / fixture / expected / assertion が [test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract) に従い、covered または not_applicable で閉じている。 | coverage record なし、open item 残存、owner 詳細本文 anchor 未接続、assertion / mutation 未接続、対象外理由 anchor 不足が残る。 |
| fixture root closure | 対象 fixture root が [fixture root coverage matrix 固定契約](#fixture-root-coverage-matrix-contract) の正式 fixture、実装検証証跡、または対象外理由へ到達できる。 | 未作成 root、未接続 root、重複 fixture、catalog 未登録 directory、harness 未参照 directory が残る。 |
| execution evidence closure | 必須検証ごとに、対象 artifact、対象 owner、対象 fixture、実行 command、終了 code、期待結果、実結果、未実行理由、完了可否への影響が [test execution evidence matrix 固定契約](#test-execution-evidence-matrix-contract) の分類で記録されている。 | 実行成功だけの記録、scope 不明、必須検証の未実行理由なし、または標準外 runtime を代替根拠にしている。 |
| oracle closure | 対象変更に関わる expected / actual、positive / negative case、状態差分、effects、security expected、禁止副作用、失敗時 no mutation が [test oracle evidence set 固定契約](#test-oracle-evidence-set-contract) に従い、責務正本 anchor と比較単位へ接続されている。 | 弱い oracle、実装結果の丸写し、snapshot 無条件受け入れ、status code だけ、fixture 存在だけ、禁止副作用未確認、失敗時 no mutation 未確認が残る。 |
| failure diagnostics closure | 対象変更に関わる test、subtest、contract assertion、fixture assertion、expected 比較、security assertion、state diff assertion が [test assertion identity / failure diagnostics evidence set 固定契約](#test-assertion-failure-diagnostics-evidence-set-contract) に従い、assertion id、対象仕様 anchor、expected / actual / diff、failure reason、再現条件、secret-safe diagnostics へ接続されている。 | assertion id なし、対象仕様 anchor なし、generic failure、差分不明、再現 command 不明、pass / fail 件数だけ、secret 露出、failure reason 未分類が残る。 |
| boundary / failure matrix closure | 対象変更に関わる入力 class、境界値、上限下限、error taxonomy、partial failure、rollback、cleanup failure、retry / recovery、失敗時 no mutation が [test boundary / failure matrix evidence set 固定契約](#test-boundary-failure-matrix-evidence-set-contract) に従い、責務正本 anchor と matrix 証跡へ接続されている。 | 正常系だけ、代表的異常系だけ、境界値未固定、error taxonomy 未固定、partial failure 未確認、rollback 未確認、cleanup failure 未確認、retry / recovery 境界未確認、失敗時 no mutation 未確認が残る。 |
| isolation closure | 対象変更に関わる test order、共有状態、fixture / expected mutation、環境変数、working directory、temp root、state dir、listener、goroutine、process、timer、file lock、network fake が [test isolation evidence set 固定契約](#test-isolation-evidence-set-contract) に従い、隔離境界と cleanup 証跡へ接続されている。 | 順序依存、共有状態汚染、fixture / expected 破壊、環境差分漏れ、残留 resource、cleanup failure 未確認、parallel 可否未記録が残る。 |
| determinism closure | 対象変更に関わる clock、timer、entropy、file order、map order、filesystem、network、process、parallel worker、retry 境界が [test determinism evidence set 固定契約](#test-determinism-evidence-set-contract) に従い、同一入力再実行と変動要因固定で同一結果を示す。 | flaky、retry pass、実時間、乱数、外部応答、OS 差分、順序差に依存する合格条件が残る。 |
| concurrency / race closure | 対象変更に関わる goroutine、channel、worker、lock、listener、timer、file lock、queue、shutdown、共有状態、並行 request、同時刻 event が [test concurrency / race evidence set 固定契約](#test-concurrency-race-evidence-set-contract) に従い、race detector、schedule / interleaving、atomicity、conflict outcome、cleanup 証跡へ接続されている。 | data race、goroutine leak、deadlock、lost update、二重 commit、二重 cleanup、順序依存、lock 競合、shutdown 中 mutation、parallel 結果未固定が残る。 |
| mutation closure | 対象変更に適用する mutation class、対象 file、対象 fixture、判定、集計が [mutation test 証跡固定契約](#mutation-test-evidence-contract) に従い、`survived=0` である。 | mutation class 未定義、判定不能、`survived` 残存、`invalid` / `equivalent` の根拠 anchor 不足。 |
| harness self-verification closure | 対象変更に関わる test harness、checker、contract drift checker、fixture assertion、expected 比較、security assertion、state diff assertion が [test harness self-verification evidence set 固定契約](#test-harness-self-verification-evidence-set-contract) に従い、検出すべき不正 fixture、欠損 expected、禁止副作用、secret leak、fake transcript 不一致、cleanup failure、無効化 mutation を fail として検出する。 | harness / checker / assertion の失効、負例なし、fail すべき fixture の pass、failure message 未固定、fake 未消費または過剰消費の未検出、cleanup failure の隠蔽が残る。 |
| contract drift closure | [test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の孤立 test、未検証契約、期待値ドリフト、harness ドリフトが 0 件である。 | いずれかの drift 種別が 1 件以上残る。 |
| cross-owner closure | API / SDK / UI、CLI / API、statefile / archive / security / runner、setup / release / admin の横断境界について、呼び出し元 owner、呼び出し先 owner、状態差分、security effect、成功後再取得、失敗時 no mutation が記録されている。 | 片側 owner のみの確認、collaborator 副作用未確認、security expected 未接続、失敗時固定なし。 |
| final open item count | 上記全項目の未完了条件が 0 件であり、残 open item が `0` として記録されている。 | open item が 1 件以上ある、件数が未記録、または残件を別変更で解消するとしている。 |

<a id="test-verification-closure-record-schema-contract"></a>
**test verification closure record schema 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[test verification closure checklist 固定契約](#test-verification-closure-checklist-contract) の判定結果を残す記録 schema だけを固定する。テスト方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Phase 11 仕様全般完了は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、現在状態と対象外理由は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan) を参照する。schema 記録は、方針、ポリシー、Phase 状態、owner 詳細本文、または実行手順を再定義してはならない。

クロージャ記録は、下表の field を持つ。`closure_item` ごとに 1 record を作成し、全 record の `open_items` が空であり、`status` が `closed` または `not_applicable` だけになった場合に限り、該当検証証跡を完了扱いにできる。`status=open` が 1 件以上ある場合、または `final_open_item_count` の `open_items` が空でない場合は、完了証跡として扱わない。

| field | 固定値 / 形式 | 必須条件 |
|-------|---------------|----------|
| `closure_item` | `scope_inventory`、`test_gap_inventory_closure`、`test_improvement_batch_closure`、`traceability_closure`、`requirement_coverage_closure`、`fixture_root_closure`、`execution_evidence_closure`、`oracle_closure`、`failure_diagnostics_closure`、`boundary_failure_matrix_closure`、`isolation_closure`、`determinism_closure`、`concurrency_race_closure`、`mutation_closure`、`harness_self_verification_closure`、`contract_drift_closure`、`cross_owner_closure`、`final_open_item_count` のいずれか。 | [test verification closure checklist 固定契約](#test-verification-closure-checklist-contract) のクロージャ項目と一致する。 |
| `status` | `closed`、`not_applicable`、`open` のいずれか。 | `closed` は未完了条件 0 件、`not_applicable` は対象外理由 anchor あり、`open` は未解消項目ありの場合だけ使用する。 |
| `owner_component` | 対象 owner component 名。 | 対象がある record では空にしてはならない。owner component の正本は [`docs/SPEC.md` 責務文書構成表](../SPEC.md#document-responsibility-map) と対象詳細本文を参照する。 |
| `collaborator_components` | collaborator component 名の配列。該当なしの場合は空配列。 | 横断境界、API / SDK / UI、statefile / archive / security / runner、setup / release / admin の接続がある場合は空配列にしてはならない。 |
| `artifacts` | 対象実装 artifact、文書 artifact、生成 artifact の配列。 | 対象 artifact がある record では空にしてはならない。実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](../DOCUMENT_INDEX.md#実装ファイル一覧) を参照する。 |
| `test_artifacts` | 対象 test、contract test、harness、checker の配列。 | `status=closed` では空にしてはならない。対象外の場合は `not_applicable_reason` と対象外 anchor を記録する。 |
| `fixture_roots` | 対象 fixture root の配列。 | fixture root を使う検証では空にしてはならない。fixture root が不要な検証では `evidence_refs` に不要根拠を記録する。 |
| `spec_refs` | 仕様正本への Markdown link 配列。 | 空配列禁止。owner 詳細本文、fixture 証跡、[`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md)、[`docs/DOCUMENT_INDEX.md`](../DOCUMENT_INDEX.md) のいずれかへ到達できる link を含める。 |
| `evidence_refs` | 実行証跡、fixture 証跡、生成物確認、差分確認、未実行理由への参照配列。 | 空配列禁止。実行成功件数だけ、または pass / fail だけの参照は不可とする。 |
| `open_items` | 未解消項目の配列。 | `status=closed` と `status=not_applicable` では空配列にする。1 件以上ある場合は `status=open` とする。 |
| `not_applicable_reason` | 対象外理由の本文または参照。 | `status=not_applicable` では必須。`status=closed` と `status=open` では空にする。 |

不正 record は完了証跡として扱わない。不正 record とは、`spec_refs` が空、`evidence_refs` が空、`status=closed` で `open_items` が空でない、`status=not_applicable` で `not_applicable_reason` または対象外理由 anchor がない、`status=open` を残したまま完了扱いにしている、または `closure_item` が [test verification closure checklist 固定契約](#test-verification-closure-checklist-contract) の項目へ対応しない record を指す。

<a id="test-verification-closure-record-set-contract"></a>
**test verification closure record set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[test verification closure record schema 固定契約](#test-verification-closure-record-schema-contract) の record を完了判定単位として束ねる条件だけを固定する。Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations)、完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) と [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit) を参照する。

closure record set は、単一の検証対象単位ごとに 1 組作成する。検証対象単位は、1 つの実装 PR、1 つの fixture harness closure、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ closure record set で複数の検証対象単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| record 数 | 1 つの closure record set は `scope_inventory`、`test_gap_inventory_closure`、`test_improvement_batch_closure`、`traceability_closure`、`requirement_coverage_closure`、`fixture_root_closure`、`execution_evidence_closure`、`oracle_closure`、`failure_diagnostics_closure`、`boundary_failure_matrix_closure`、`isolation_closure`、`determinism_closure`、`concurrency_race_closure`、`mutation_closure`、`harness_self_verification_closure`、`contract_drift_closure`、`cross_owner_closure`、`final_open_item_count` の 18 record だけを各 1 件持つ。 | 18 件未満、19 件以上、同じ `closure_item` の重複、未登録 `closure_item` がある。 |
| record 順序 | record を配列または箇条書きで記録する場合は、[test verification closure checklist 固定契約](#test-verification-closure-checklist-contract) の表順と同じ順序にする。 | 順序不一致により review 時に欠落または重複を判定できない。 |
| status closure | 完了扱いにできる closure record set は、全 record の `status` が `closed` または `not_applicable` であり、全 record の `open_items` が空であり、`final_open_item_count` が残件 `0` を示す。 | `status=open`、`open_items` 残存、残件数未記録、または残件を別変更で解消すると記録している。 |
| scope consistency | 全 record の `owner_component`、`collaborator_components`、`artifacts`、`test_artifacts`、`fixture_roots`、`spec_refs`、`evidence_refs` は同じ検証対象単位を指す。 | record 間で対象 owner、artifact、fixture root、または根拠 anchor が別範囲を指している。 |
| not applicable | `status=not_applicable` は、対象外理由、対象外にする owner / artifact / fixture root、責務正本 anchor、完了可否への影響を同じ record に持つ。 | 理由だけ、または anchor だけで対象外範囲と完了可否への影響が不明である。 |
| PR evidence | 実装変更、検証変更、fixture 変更、または意味のあるテスト / test gap inventory / batch closure / requirement coverage / oracle / failure diagnostics / boundary / failure matrix / isolation / determinism / race trigger / concurrency / race / mutation selection / mutation test / harness self-verification / contract drift の完了可否に関わる変更では、Pull Request の `Verification` に [implementation PR evidence template 固定契約](#implementation-pr-evidence-template-contract) に基づく implementation PR evidence package と closure record set の記録先、または対象外理由を記録する。 | Pull Request 上で implementation PR evidence package、closure record set の所在、対象外理由、または未完了扱いが確認できない。 |

<a id="implementation-pr-evidence-template-contract"></a>
**implementation PR evidence template 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装 PR、検証 PR、fixture PR、または Phase 全体完了 PR で Pull Request 本文に提出する証跡 package の最小構成だけを固定する。Pull Request 本文への記載義務は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations)、完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) と [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit)、対象 Phase と現在状態は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan) を参照する。

implementation PR evidence package は、単一の Pull Request ごとに 1 組作成する。同じ package に複数 Phase、複数の無関係な owner、または同一 Pull Request で完了判定しない後続作業を混在させてはならない。

| 提出物 | 最小記録形式 | 未完了条件 |
|--------|--------------|------------|
| submitted artifacts | 対象 Phase、対象 owner、collaborator、変更 artifact、test artifact、fixture root、expected / fake / effects / security expected、生成物、状態更新対象、対象外 artifact を列挙する。 | 対象 artifact、対象 owner、fixture root、対象外 artifact のいずれかが不明である。 |
| test gap inventory / batch closure | [test gap inventory record 固定契約](#test-gap-inventory-record-contract) と [test improvement batch closure 固定契約](#test-improvement-batch-closure-contract) の所在、gap record 件数、`open=0`、`final_open_item_count=0` を記録する。 | inventory 所在なし、record 件数未記録、`open>0`、または batch closure と closure record set が接続していない。 |
| requirement coverage ledger | [test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract) の所在、requirement id 数、covered 件数、not_applicable 件数、open 件数 `0` を記録する。 | ledger 所在なし、open 件数未記録、または `open>0`。 |
| 18 record closure set | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の所在、18 record の status 一覧、`final_open_item_count=0` を記録する。 | record 数不一致、`status=open`、`final_open_item_count` 未記録、または `open_items` 残存。 |
| oracle evidence | [test oracle evidence set 固定契約](#test-oracle-evidence-set-contract) の所在、positive / negative case、expected / actual 比較単位、状態差分、禁止副作用、失敗時 no mutation を記録する。 | status code だけ、fixture 存在だけ、expected / actual 不明、禁止副作用未確認。 |
| failure diagnostics evidence | [test assertion identity / failure diagnostics evidence set 固定契約](#test-assertion-failure-diagnostics-evidence-set-contract) の所在、assertion id、対象仕様 anchor、expected / actual / diff、failure reason、再現 command、secret-safe diagnostics を記録する。 | assertion id なし、差分不明、generic failure、secret-safe 診断未確認。 |
| boundary / failure matrix evidence | [test boundary / failure matrix evidence set 固定契約](#test-boundary-failure-matrix-evidence-set-contract) の所在、入力 class、境界値、error taxonomy、partial failure、rollback、cleanup、retry / recovery、失敗時 no mutation を記録する。 | 正常系だけ、境界未固定、partial failure / rollback / cleanup / retry 未確認。 |
| isolation evidence | [test isolation evidence set 固定契約](#test-isolation-evidence-set-contract) の所在、隔離境界、順序入替結果、共有状態初期化、cleanup、残留 resource、parallel 実行可否を記録する。 | 共有状態汚染、順序依存、cleanup 不明、残留 resource 未確認、parallel 可否未記録。 |
| determinism evidence | [test determinism evidence set 固定契約](#test-determinism-evidence-set-contract) の所在、変動要因、fake adapter、同一入力再実行結果、順序入替 case、禁止実環境依存を記録する。 | retry 成功だけ、実時間 / 乱数 / host 依存、再実行一致未確認、fake 境界不明。 |
| concurrency / race evidence | [race trigger matrix 固定契約](#race-trigger-matrix-contract) と [test concurrency / race evidence set 固定契約](#test-concurrency-race-evidence-set-contract) の所在、共有 resource、race trigger 件数、open trigger `0`、race detector、schedule / interleaving、lock / channel / goroutine lifecycle、conflict outcome、atomicity、cleanup を記録する。 | race trigger 未判定、data race 未確認、goroutine leak 未確認、lock / channel 終了条件不明、競合結果未固定。 |
| mutation evidence | [mutation selection ledger 固定契約](#mutation-selection-ledger-contract) と [mutation test evidence set 固定契約](#mutation-test-evidence-set-contract) の所在、mutation class decision、`killed` / `survived` / `invalid` / `equivalent` 件数、`survived=0` を記録する。 | mutation class 未定義、`survived>0`、または `equivalent` 根拠 anchor 不足。 |
| harness self-verification evidence | [test harness self-verification evidence set 固定契約](#test-harness-self-verification-evidence-set-contract) の所在、negative control、positive control、検出すべき不正、failure reason を記録する。 | harness が常に pass / 常に fail、negative / positive 片側だけ、failure reason 不明。 |
| contract drift evidence | [test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の所在、孤立 test、未検証契約、期待値ドリフト、harness ドリフトの件数、各 drift の解消状態を記録する。 | drift 件数未記録、孤立 test 残存、未検証契約残存、期待値または harness の正本不一致。 |
| execution evidence | [test execution evidence matrix 固定契約](#test-execution-evidence-matrix-contract) の分類、実行 command、runtime、終了 code、対象 artifact、対象 fixture、未実行理由を記録する。 | 実行成功件数だけ、対象 artifact 不明、必須検証の未実行理由なし。 |
| not applicable evidence | 対象外にした owner、artifact、fixture root、evidence set、理由、責務正本 anchor、完了可否への影響を記録する。 | 理由だけ、anchor だけ、対象外範囲不明、または将来対応を対象外理由にしている。 |
| completion declaration | `open item=0`、対象 Phase 全体完了可否、`docs/ROADMAP.md` 更新要否、`docs/DOCUMENT_INDEX.md` 更新要否、未実施確認の有無を記録する。 | 残件を別 PR に送る、状態更新要否不明、または完了可否が Phase 全体と一致しない。 |

Pull Request 本文で implementation PR evidence package を記録する場合は、上表の提出物名を label として表順に並べる。適用外の提出物は削除せず、`not applicable evidence` に対象外範囲、理由、責務正本 anchor、完了可否への影響を記録する。上表にない任意 label、表順と異なる記録順、または `completion declaration` より後に残件を追記する形式を完了証跡として扱ってはならない。

implementation PR evidence package は、上表の適用項目すべてが記録され、未完了条件が 0 件である場合だけ完了証跡として扱う。実装変更がない文書整理 PR では、対象外理由を Pull Request 本文へ記録すればよい。実装変更または検証変更があるのに package を作れない場合は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い検証不足として扱う。

<a id="test-requirement-coverage-ledger-contract"></a>
**test requirement coverage ledger 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、owner 詳細本文の検証条件、fixture 証跡条件、Phase 対象、既存 test / fixture / expected / assertion を、対象変更単位ごとに covered / not_applicable / open へ分類する ledger 条件だけを固定する。検証方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

test requirement coverage ledger は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ ledger で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| requirement inventory | owner 詳細本文の検証条件、fixture 証跡条件、Phase 対象、source-code audit 対象、既存 test / fixture / expected / assertion を requirement id 付きで列挙する。 | 検証条件、Phase 対象、既存 assertion のいずれかが ledger に存在しない。 |
| requirement identity | requirement id、owner component、collaborator component、対象 anchor、対象 artifact、検出したい仕様違反、必須 / 任意 / 対象外候補を記録する。 | requirement id が不安定、anchor なし、owner 不明、または検出したい仕様違反が空である。 |
| coverage link | 各 requirement は test artifact、fixture root、expected file、assertion id、mutation class、evidence set、closure item のいずれか 1 件以上へ接続する。 | covered と記録した requirement が test / fixture / expected / assertion / mutation / closure item のいずれにも接続していない。 |
| not applicable record | not_applicable は対象外 owner / artifact / fixture root、対象外理由、責務正本 anchor、完了可否への影響を同じ record に持つ。 | 理由だけ、anchor だけ、対象外範囲不明、または将来対応を対象外理由にしている。 |
| open item handling | open requirement は不足 test、追加すべき fixture、追加すべき expected / assertion / mutation class、完了不可理由を記録する。 | open item があるのに完了扱い、または不足内容が test / fixture / expected / mutation のどれか不明である。 |
| duplicate / conflict check | 同一 requirement が複数 owner に重複していないこと、または collaborator 境界として分離されていることを記録する。 | 同一検証条件を複数 owner が本文として持つ、または conflicting expected が残る。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `requirement_coverage_closure` record から ledger の所在へ到達できる。 | closure record set と ledger の対象 owner、artifact、fixture root、または scope が一致しない。 |

test requirement coverage ledger は、対象変更単位に含まれる検証条件、fixture 証跡条件、Phase 対象、既存 test / fixture / expected / assertion のいずれかが covered または not_applicable で閉じていない場合、完了証跡として扱わない。ledger を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-oracle-evidence-set-contract"></a>
**test oracle evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の oracle 条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。弱い oracle、実装結果の丸写し、snapshot 無条件受け入れ、status code だけの確認、fixture 存在だけの確認、禁止副作用未確認の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

oracle evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ oracle evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| oracle scope | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、対象仕様 anchor、検出したい仕様違反を記録する。 | 対象 artifact、assertion、または責務正本 anchor が空であり、対象外理由もない。 |
| input / expected / actual | 入力 fixture、fake input、expected file、actual 取得元、比較単位、比較方法、正規化有無を記録する。 | expected と actual の取得元が不明、実装出力をそのまま expected 化している、または比較単位が不明である。 |
| positive / negative case | 正常系だけでなく、入力不正、権限不足、欠損、境界値、外部失敗、timeout、partial failure、rollback、cleanup failure の該当 case を記録する。 | 正常系だけ、または対象機能が持つ失敗系を対象外 anchor なしで省略している。 |
| response / exit / error | HTTP status、JSON body、error code、CLI stdout、CLI stderr、exit code、SDK error、MCP JSON-RPC error の該当項目を固定する。 | status code だけ、stderr 空だけ、response body 不問、error class 不問、または終了 code 不問である。 |
| state / effects / security | 状態差分、write order、updated / unchanged / forbidden paths、外部通信、notification、log、audit、cache、secret 非露出を固定する。 | 状態差分なし、禁止副作用未確認、security expected 未接続、または secret 出力確認なしである。 |
| failure no mutation | 失敗 case では、業務状態、secret、audit、metrics、cache、snapshot、queue、external call の変更可否を固定し、禁止 mutation を `expected/effects.json` または同等の証跡へ接続する。 | 失敗時に何が変わらないか不明、または partial failure の許可副作用と禁止副作用が分離されていない。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `oracle_closure` record から oracle evidence set の所在へ到達できる。 | closure record set と oracle evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

oracle evidence set は、対象変更単位に含まれる expected、assertion、state diff、effects、security expected、error assertion、または failure case が 1 つでも責務正本 anchor と比較単位へ接続していない場合、完了証跡として扱わない。oracle evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-assertion-failure-diagnostics-evidence-set-contract"></a>
**test assertion identity / failure diagnostics evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の assertion identity と failure diagnostics 条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。assertion id 未固定、対象仕様 anchor 未接続、expected / actual / diff 不明、generic failure、pass / fail 件数だけの証跡、secret を含む診断出力の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

failure diagnostics evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ failure diagnostics evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| diagnostics scope | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、対象仕様 anchor、検出したい仕様違反を記録する。 | 対象 artifact、assertion、責務正本 anchor、または検出したい仕様違反が空であり、対象外理由もない。 |
| assertion identity | assertion id、subtest 名、fixture 名、manifest assertion、expected file、比較対象 field、関連 mutation class、関連 closure item を安定値として記録する。 | assertion id が不安定、fixture 名だけ、test 名だけ、line number だけ、または expected file だけで assertion を識別している。 |
| expected / actual / diff | expected 取得元、actual 取得元、比較単位、正規化有無、missing、extra、type mismatch、value mismatch、order mismatch、state diff、effects diff、security diff を記録する。 | 期待値と実値の片方だけ、差分単位不明、snapshot 全体差分だけ、または差分を人間の目視確認だけに依存している。 |
| failure reason | `mismatch`、`missing_expected`、`unexpected_side_effect`、`forbidden_output`、`security_leak`、`state_diff`、`harness_error`、`timeout`、`race`、`mutation_survived`、`contract_drift`、`unknown` のいずれかを記録し、`unknown` は未完了扱いにする。 | generic な `failed`、panic だけ、exit code だけ、件数だけ、または failure reason が未分類である。 |
| reproduction context | 再現 command、対象 fixture root、fake input、環境固定値、実行 owner、関連 expected file、必要な未実行理由を記録する。 | 再現 command なし、fixture root 不明、fake 不明、または local 環境依存の再現条件だけを記録している。 |
| secret-safe diagnostics | 診断出力、diff、failure message、stdout / stderr 要約、PR evidence に token、password、secret、Authorization header、Cookie、session id、TOTP secret の平文が出ないことを記録する。 | secret を含む diff、平文 token、hash 入力、復元可能な部分文字列、または secret 有無の未確認が残る。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `failure_diagnostics_closure` record から failure diagnostics evidence set の所在へ到達できる。 | closure record set と failure diagnostics evidence set の対象 owner、artifact、fixture root、assertion、または scope が一致しない。 |

failure diagnostics evidence set は、対象変更単位に含まれる test、subtest、contract assertion、fixture assertion、expected 比較、security assertion、state diff assertion のいずれかが assertion identity、expected / actual / diff、failure reason、reproduction context、secret-safe diagnostics へ接続していない場合、完了証跡として扱わない。failure diagnostics evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-boundary-failure-matrix-evidence-set-contract"></a>
**test boundary / failure matrix evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の正常系 / 異常系、境界値、外部境界失敗、partial failure、rollback、cleanup failure 条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。入力 class の判断漏れ、limit 境界未固定、error taxonomy 未固定、partial failure 未確認、rollback 未確認、cleanup failure 未確認、retry / recovery 境界未確認、失敗時 no mutation 未確認の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

boundary / failure matrix evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ boundary / failure matrix evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| matrix scope | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、対象仕様 anchor、入力境界、失敗境界、状態境界の該当有無を記録する。 | 対象 artifact、境界種別、または責務正本 anchor が空であり、対象外理由もない。 |
| input class matrix | valid、invalid、missing、empty、null、wrong type、unknown key、duplicate key、malformed encoding、path traversal、control character、oversized、under limit、over limit の該当 case と対象外理由を記録する。 | 入力 class の判断漏れ、正常系だけ、代表的 invalid だけ、または対象外 anchor なしで省略している。 |
| limit matrix | min、max、below min、above max、zero、one、empty、large、timeout、retry count、size、count、ID length、path depth、header length、body length、queue length、archive entry count の該当 case を記録する。 | off-by-one、上限超過、下限未満、timeout 境界、size / count 境界、または対象外理由が未記録である。 |
| error taxonomy | HTTP status、JSON error code、CLI stderr、exit code、SDK error class、MCP JSON-RPC error、statefile error、security error、external dependency error を対象 owner の契約へ接続する。 | status code だけ、error class だけ、stderr だけ、または owner 詳細本文の error 契約へ接続していない。 |
| partial failure matrix | write failure、read failure、decode failure、external call failure、fake failure、rollback failure、cleanup failure、audit failure、notification failure、archive failure、release upload failure、systemd failure の該当 case を記録する。 | partial failure の許可副作用と禁止副作用が分離されていない、または失敗注入が 1 種類だけで対象範囲を閉じている。 |
| rollback / cleanup | rollback 対象、rollback 順序、cleanup 対象、cleanup 順序、cleanup failure 時の最終状態、再実行可否、残留 resource、禁止復旧を記録する。 | rollback 不明、cleanup 成功だけ、cleanup failure 隠蔽、残留 resource 未確認、または再実行境界未記録である。 |
| retry / recovery | retry する条件、retry しない条件、retry 上限、backoff、timeout 優先順位、resume、idempotency、duplicate request、conflict response を記録する。 | retry 後成功だけ、retry 上限不明、timeout と成功の同時境界不明、resume / idempotency 未確認である。 |
| failure no mutation | 失敗 case ごとに、業務状態、secret、audit、metrics、cache、snapshot、queue、external call、file write の変更可否を固定し、許可副作用と禁止副作用を分離する。 | 失敗時に変わらない対象が不明、または許可副作用と禁止副作用が同じ expected に混在している。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `boundary_failure_matrix_closure` record から boundary / failure matrix evidence set の所在へ到達できる。 | closure record set と boundary / failure matrix evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

boundary / failure matrix evidence set は、対象変更単位に含まれる入力 class、limit、error taxonomy、partial failure、rollback、cleanup、retry / recovery、failure no mutation のいずれかが責務正本 anchor と matrix 証跡へ接続していない場合、完了証跡として扱わない。boundary / failure matrix evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-isolation-evidence-set-contract"></a>
**test isolation evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の isolation 条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。順序依存、共有状態汚染、fixture / expected 破壊、環境差分漏れ、残留 resource、cleanup failure 未確認の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

isolation evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ isolation evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| isolation scope | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、対象仕様 anchor、共有 resource の該当有無を記録する。 | 対象 artifact、共有 resource、または責務正本 anchor が空であり、対象外理由もない。 |
| test order independence | 対象 test / fixture は通常順、逆順、または fixture 名順以外の固定順で同一結果になることを記録する。順序入替が仕様上不要な場合は対象外理由 anchor を記録する。 | 1 順序だけの成功、順序入替で結果が変わる、または対象外理由がない。 |
| immutable fixture / expected | test 実行中に入力 fixture、expected、manifest、仕様文書を直接変更しないことを記録する。生成物は test ごとの temp root または state dir へ分離する。 | fixture / expected を直接 mutation する、または生成物と入力 fixture が同じ path を共有する。 |
| environment isolation | 環境変数、working directory、timezone、locale、home、temp root、state dir、process id 依存を test 単位で固定または復元する。 | 環境変更を復元しない、host 環境を合格条件にする、または test 間で state dir を共有する。 |
| runtime residue | listener、goroutine、process、timer、file lock、network fake、HTTP fake、browser fake、open file が test 後に残留しないことを記録する。 | port、goroutine、process、timer、lock、fake server、open file の残留確認がない。 |
| parallel boundary | parallel 実行可能な test は共有 state を持たず、parallel 実行不可の test は不可理由、影響 resource、直列化条件、責務正本 anchor を記録する。 | parallel 可否が不明、または不可理由なしで順序依存を許容している。 |
| cleanup failure | cleanup failure、temp root 削除失敗、lock 解放失敗、fake server close 失敗、state dir 復元失敗の扱いを記録する。 | cleanup 成功だけ、または cleanup failure が test 結果へ影響しないまま隠れる。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `isolation_closure` record から isolation evidence set の所在へ到達できる。 | closure record set と isolation evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

isolation evidence set は、対象変更単位に含まれる共有 resource、fixture、expected、環境差分、runtime resource、cleanup のいずれかが isolation scope、対象外理由、または closure record へ接続していない場合、完了証跡として扱わない。isolation evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-determinism-evidence-set-contract"></a>
**test determinism evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の決定性条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。flaky test、retry pass、実時間依存、乱数依存、外部応答依存の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

determinism evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ determinism evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| 変動要因棚卸し | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、clock、timer、entropy、file order、map order、filesystem、network、GitHub API、systemd、process、browser runtime、parallel worker、retry 境界の該当有無を記録する。 | 変動要因の判断漏れ、対象 artifact 未記録、または対象外理由と責務正本 anchor がない。 |
| fake adapter 接続 | 該当する変動要因は、[fake adapter 接続固定契約](#fixture-fake-adapter-binding-contract) または owner 詳細本文の固定 fake 境界へ接続する。 | 実 clock、実乱数、実 filesystem order、実 network、実 systemd、実 process、実 browser runtime、実 sleep を合格条件にしている。 |
| 同一入力再実行 | 同一 fixture、同一 fake input、同一 expected で 2 回以上の実行結果が一致することを記録する。 | 1 回の成功だけ、または再実行結果が異なるのに完了扱いにしている。 |
| 変動順序 case | file order、map order、parallel worker、timer、event、stream、retry に関わる変更は、順序入替 case または同時刻 case を記録する。 | 順序入替で結果が変わる、または順序差を検証していない。 |
| retry / rerun boundary | 失敗後の再実行、retry、backoff、timeout を扱う変更は、最初の失敗、retry 条件、retry 上限、最終結果、retry しない条件を同じ証跡に記録する。 | retry 後の成功だけを記録している、retry 上限が不明、または失敗を隠している。 |
| host independence | host 固有 path、OS path separator、timezone、locale、process id、file mtime、directory iteration order、環境変数差分が expected に影響しないことを記録する。 | host 固有値を expected に含める、または host 差分で結果が変わる。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `determinism_closure` record から determinism evidence set の所在へ到達できる。 | closure record set と determinism evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

determinism evidence set は、対象変更単位に含まれる変動要因が 1 つでも fake adapter、固定 fixture、対象外理由のいずれにも接続していない場合、完了証跡として扱わない。determinism evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-concurrency-race-evidence-set-contract"></a>
**test concurrency / race evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の並行処理 / race 条件を、対象変更単位ごとに完了判定できる証跡 set として記録する条件だけを固定する。data race、goroutine leak、deadlock、lost update、二重 commit、二重 cleanup、順序依存、lock 競合、shutdown 中 mutation、parallel 結果未固定の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

concurrency / race evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ concurrency / race evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| 並行対象棚卸し | 対象 owner、collaborator、対象 file、対象 test / fixture / expected / assertion、goroutine、channel、worker、lock、listener、timer、file lock、queue、shutdown、共有状態、並行 request、同時刻 event の該当有無を記録する。 | 並行対象の判断漏れ、対象 artifact 未記録、または対象外理由と責務正本 anchor がない。 |
| race detector | Go 実装または Go test が goroutine、channel、共有 memory、lock、timer、listener、worker を扱う場合は、対象 package または対象 test を `go test -race` 相当で検証した結果を記録する。実行不能な場合は、未実行理由、影響 owner、代替不可範囲、完了可否への影響を記録する。 | race detector 対象なのに結果がない、実行不能を対象外扱いにしている、または通常の `go test` 成功だけを race-free 根拠にしている。 |
| schedule / interleaving | 同時 start、同時 cancel、timeout 直前、shutdown 中 request、lock 待ち、queue 競合、worker 順序差、channel close 競合、timer 発火順序を対象機能に応じて固定する。 | 1 順序だけの成功、schedule 差分で結果が変わる、または同時刻 case が対象外 anchor なしで省略されている。 |
| lock / channel / goroutine lifecycle | lock 取得 / 解放、channel send / receive / close、goroutine start / exit、listener open / close、timer stop / drain、file lock cleanup の成功と失敗を記録する。 | lock 解放、goroutine 終了、timer cleanup、listener close、file lock cleanup の証跡がない。 |
| conflict outcome | 同じ状態 path、queue entry、snapshot、notification、audit、session、token、release tag、output path へ並行 mutation が到達する場合、勝者、敗者、HTTP status / exit code、状態差分、禁止副作用を固定する。 | lost update、二重 write、二重 notification、二重 audit、二重 cleanup、または敗者側の状態不変が未確認である。 |
| atomicity / visibility | 並行 reader / writer、partial failure、rollback、cleanup failure、rename、fsync、state reload、cache update を扱う変更は、中間状態の可視性、最終状態、read-your-write 可否、no partial publish を記録する。 | reader が中間状態を成功として読める、partial publish、rollback 不明、または visibility 条件が未記録である。 |
| parallel boundary | parallel 実行可能な test は共有 resource を持たず、parallel 実行不可の test は不可理由、影響 resource、直列化条件、責務正本 anchor を記録する。 | parallel 可否が不明、または不可理由なしで順序依存を許容している。 |
| cleanup after race | 競合、timeout、panic、cancel、shutdown、fake failure 後に残る goroutine、listener、timer、lock、temp root、state file、pending entry、fake event を記録する。 | 競合後 cleanup 未確認、残留 resource、fake event 未消費、または cleanup failure が隠れる。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `concurrency_race_closure` record から concurrency / race evidence set の所在へ到達できる。 | closure record set と concurrency / race evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

concurrency / race evidence set は、対象変更単位に含まれる並行対象が 1 つでも race detector、schedule / interleaving、lock / channel / goroutine lifecycle、conflict outcome、atomicity、cleanup、対象外理由のいずれにも接続していない場合、完了証跡として扱わない。concurrency / race evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="mutation-test-evidence-contract"></a>
**mutation test 証跡固定契約：**

mutation test の完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) を正本とし、fixture 証跡責務では実装検証証跡へ記録する項目だけを固定する。mutation test 証跡は、対象 test、fixture、expected、security expected、state diff、contract drift checker、または harness assertion が検出すべき仕様違反を、責務名付き Markdown link で正本へ到達できる形にする。

| 証跡項目 | 必須内容 | 不足時の扱い |
|----------|----------|--------------|
| 対象責務 | owner component、collaborator component、対象機能、対象 file、対象 test 名または fixture 名、対象 anchor。 | 仕様追跡性不足として未完了。 |
| mutation class | [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の mutation class 名、変異対象、検出したい仕様違反。 | mutation 対象未定義として未完了。 |
| mutation 実行単位 | 変更した実装コード、test harness、fixture assertion、expected 比較、security assertion、state diff assertion のいずれを変異対象にしたか。 | 実装変更と検証変更の対応不足として未完了。 |
| mutation operation | `operation_id`、対象 file、対象 function / assertion / expected、変異前、変異後、適用方法、temp copy 使用有無を記録する。 | 変異内容が再現できない、変異前後が不明、または正本 artifact を直接変更したまま残している。 |
| mutation patch 証跡 | unified diff、変換規則、または before / after expression のいずれかで mutation を再適用できる証跡を記録する。 | 文章説明だけ、目視確認だけ、または mutation を再実行できない。 |
| 実行方法 | 実行コマンド、固定入力、fake clock / fake entropy / fake filesystem / fake HTTP / fake process / fake browser runtime の使用有無、外部通信禁止確認。 | 再現不能として未完了。 |
| 判定 | `killed`、`survived`、`invalid`、`equivalent` のいずれか、期待 failure、実際の failure、終了コード、差分。 | 判定不能として未完了。 |
| survived 対応 | `survived` が 1 件以上ある場合の不足 test、未固定仕様、追加すべき fixture、完了不可理由。 | `survived` を残したまま完了扱い不可。 |
| invalid / equivalent 理由 | `invalid` は構文または観測不能理由、`equivalent` は観測可能挙動が同一である責務正本 anchor と理由。 | 理由または anchor 不足時は `survived` として扱う。 |
| 集計 | 対象変更単位ごとの `killed` / `survived` / `invalid` / `equivalent` 件数。 | 完了証跡不足として未完了。 |

mutation test 証跡の集計では、`survived=0` を完了条件とする。`invalid` と `equivalent` は件数を隠してはならず、対象 mutation class、対象 file、対象 fixture、対象 anchor を個別に記録する。`equivalent` を理由にする場合でも、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の `equivalent` 定義を満たす責務正本 anchor がないものは `survived` として扱う。

mutation test 証跡は、coverage 証跡、`go test` 成功、`deno check` 成功、fixture 存在確認、snapshot 一致、手作業の確認、または実装者の判断で代替してはならない。coverage を記録する場合でも、未検出 mutation、弱い assertion、未接続 fixture、仕様 anchor 不足がある場合は完了不可とする。

<a id="mutation-test-evidence-set-contract"></a>
**mutation test evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[mutation test 証跡固定契約](#mutation-test-evidence-contract) の証跡を対象変更単位ごとに束ね、`mutation_closure` の完了可否を判定できる条件だけを固定する。mutation test の必須条件、mutation class、判定語彙、完了可否は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

mutation evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ mutation evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| 対象棚卸し | `owner_component`、`collaborator_components`、対象 file、対象 test / fixture / expected / assertion、対象仕様 anchor を記録する。 | 対象 artifact、test artifact、fixture root、または仕様 anchor が空であり、対象外理由もない。 |
| mutation class decision | [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) の mutation class ごとに、`applicable` または `not_applicable` を記録する。`not_applicable` は対象外 owner / artifact / assertion、理由、責務正本 anchor、完了可否への影響を同じ record に持つ。 | mutation class の判断漏れ、理由だけの対象外、anchor だけの対象外、または対象外範囲が不明である。 |
| applicable record | `applicable` とした mutation class は、[mutation test 証跡固定契約](#mutation-test-evidence-contract) の証跡項目を満たす record を 1 件以上持つ。 | `applicable` なのに record がない、または record が対象 file、対象 fixture、期待 failure、実際の failure へ到達できない。 |
| verdict reconciliation | `killed` / `survived` / `invalid` / `equivalent` の集計は mutation evidence set 内の record 数と一致し、`survived=0` である。`invalid` と `equivalent` は個別 record に理由と責務正本 anchor を持つ。 | 集計不一致、`survived` 残存、`invalid` / `equivalent` の理由不足、または責務正本 anchor 不足がある。 |
| changed artifact coverage | 対象変更単位に含まれる実装、test、fixture、expected、security expected、state diff、contract drift checker、harness assertion は、少なくとも 1 つの mutation class decision へ接続する。 | 変更 artifact が mutation class decision へ接続していない。 |
| evidence refs | 実行 command、固定入力、fake clock / fake entropy / fake filesystem / fake HTTP / fake process / fake browser runtime、外部通信禁止確認、終了 code、重要 stdout / stderr 要約を記録する。 | 再実行できない、fake 境界が不明、外部通信禁止を確認できない、または実行結果だけで根拠が不足している。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `mutation_closure` record から mutation evidence set の所在へ到達できる。 | closure record set と mutation evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

mutation evidence set は、対象変更単位に含まれる変更 artifact が 1 つでも mutation class decision へ接続していない場合、完了証跡として扱わない。mutation evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-harness-self-verification-evidence-set-contract"></a>
**test harness self-verification evidence set 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、test harness、checker、contract drift checker、fixture assertion、expected 比較、security assertion、state diff assertion が、検出すべき不正を fail として検出できる証跡 set の条件だけを固定する。意味のあるテスト、完了可否、mutation test の必須条件は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Pull Request 本文への記録手順は [`AGENTS.md` Git 運用ルール](../../AGENTS.md#agents-git-operations) を参照する。

harness self-verification evidence set は、単一の対象変更単位ごとに 1 組作成する。対象変更単位は、1 つの実装変更、1 つの検証変更、1 つの fixture / expected 変更、1 つの owner artifact 検証、または 1 つの Phase 全体完了判定のいずれかに固定する。同じ harness self-verification evidence set で複数の対象変更単位を混在させてはならない。

| 固定項目 | 契約 | 未完了条件 |
|----------|------|------------|
| harness scope | `owner_component`、`collaborator_components`、対象 test、対象 harness、対象 checker、対象 assertion、対象 fixture root、対象 expected、対象仕様 anchor を記録する。 | harness、checker、assertion、fixture root、expected、仕様 anchor のいずれかが空であり、対象外理由もない。 |
| negative control set | expected 欠損、manifest / assertion mismatch、expected / actual mismatch、禁止副作用あり、secret leak あり、fake event 未消費、fake event 過剰、fake event 順序違反、cleanup failure、assertion 無効化、mutation disabled の該当 case を fail させる証跡を記録する。 | 該当する不正 case がない、または不正 case が pass している。 |
| expected failure identity | 各 negative control は、期待 exit code、期待 error class、期待 stderr / log 要約、期待 failure reason、対象 assertion 名を固定する。 | fail した事実だけで、どの assertion が何を検出したか判定できない。 |
| positive control boundary | 同じ harness / checker / assertion が、正しい fixture、正しい expected、正しい fake transcript、正しい cleanup では pass することを記録する。 | negative control だけで、harness 自体が常に fail する可能性を排除できない。 |
| fake transcript verification | fake clock / entropy / filesystem / HTTP / process / browser runtime / GitHub / systemd / notifier を使う場合は、期待 call、消費順、未消費 event、過剰 call、禁止 external call の検出結果を記録する。 | fake input の未消費、過剰消費、順序違反、禁止 external call を検出できない。 |
| self-test isolation | harness self-verification は canonical fixture / expected を直接破壊せず、copy、synthetic fixture、または temp root で不正 case を作る。 | canonical fixture / expected を mutation する、または自己検証の失敗が後続 test に影響する。 |
| closure connection | [test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `harness_self_verification_closure` record から harness self-verification evidence set の所在へ到達できる。 | closure record set と harness self-verification evidence set の対象 owner、artifact、fixture root、または scope が一致しない。 |

harness self-verification evidence set は、対象変更単位に含まれる harness、checker、assertion、expected 比較、security assertion、state diff assertion のいずれかが negative control set と positive control boundary の両方へ接続していない場合、完了証跡として扱わない。harness self-verification evidence set を作れない場合は、対象外として扱わず、[`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) に従い仕様不足または検証不足として扱う。

<a id="test-contract-drift-evidence-contract"></a>
**test / contract drift 証跡固定契約：**

test / contract drift 判定では、実装 test、contract test、fixture manifest、owner 詳細本文、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](../DETAIL_INDEX.md#phase-11-quality-gate-entry) の対応を照合する。test または assertion が存在するのに owner 詳細本文、fixture 証跡、または [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務の対象外理由へ到達できない場合は、孤立 test として未完了とする。owner 詳細本文で検証条件を持つのに test、fixture、expected、mutation 証跡へ到達できない場合は、未検証契約として未完了とする。

| drift 種別 | 未完了条件 | 解消方法 |
|------------|------------|----------|
| 孤立 test | test 名、subtest、contract assertion、fixture presence check が責務正本 anchor へ到達できない。 | 対象 owner 詳細本文または fixture 証跡へ参照を接続する。仕様に存在しない期待値なら test または期待値を削除する。 |
| 未検証契約 | owner 詳細本文に入力、出力、異常系、状態、副作用、security、境界値、mutation 条件があるが対応する test / fixture / expected / mutation 証跡がない。 | fixture と test を追加し、実装検証証跡へ記録する。 |
| 期待値ドリフト | expected、snapshot、security expected、effects が owner 詳細本文または [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test) と矛盾する。 | expected を責務正本に合わせる。実装挙動へ合わせるための期待値緩和は禁止する。 |
| harness ドリフト | fixture manifest、fixture catalog、directory 名、test 名、実装検証証跡の対象 owner または fixture 名が一致しない。 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) と対象 fixture catalog へ一致させる。 |

<a id="test-contract-drift-report-schema-contract"></a>
**test / contract drift report schema 固定契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、[test / contract drift 証跡固定契約](#test-contract-drift-evidence-contract) の判定結果を対象変更単位ごとに残す record schema だけを固定する。drift の完了禁止は [`docs/SPEC.md` ポリシー責務 §0g](../SPEC.md#policy-meaningful-test)、Phase 11 完了判定は [`docs/SPEC.md` 方針責務 §4.8](../SPEC.md#sec-4-8) を参照する。

test / contract drift report は、単一の対象変更単位ごとに 1 組作成する。同じ report に複数 Phase、無関係な owner、または完了判定しない後続作業を混在させてはならない。

| field | 固定値 / 形式 | 未完了条件 |
|-------|---------------|------------|
| `drift_id` | `drift.<owner>.<type>.<case>` の lowercase dot notation。 | id がない、重複、対象 owner 不明。 |
| `drift_type` | `orphan_test`、`unverified_contract`、`expected_drift`、`harness_drift` のいずれか。 | 未登録値、または drift 種別不明。 |
| `owner_component` | 対象 owner component 名。 | owner 不明、または collaborator のみで owner がない。 |
| `source_ref` | test、fixture、expected、manifest、harness、checker、owner 詳細本文 anchor のいずれかへの責務名付き Markdown link。 | 裸のファイル名、line number だけ、または実在しない anchor。 |
| `canonical_ref` | 正本とする owner 詳細本文、fixture 証跡、[`docs/SPEC.md`](../SPEC.md)、[`docs/ROADMAP.md`](../ROADMAP.md) の責務名付き Markdown link。 | 正本 anchor なし、または実装ファイルだけを正本にしている。 |
| `evidence_refs` | drift 判定に用いた test artifact、fixture root、expected、manifest、coverage ledger、mutation evidence、closure record への link 配列。 | 判定根拠が説明文だけ、または pass 件数だけである。 |
| `resolution` | `fixed`、`not_applicable`、`open` のいずれか。 | open なのに完了扱い、または `not_applicable` に責務正本 anchor がない。 |
| `open_item` | `resolution=open` の場合に必要な追加 test、fixture、expected、mutation、詳細仕様改訂を記録する。 | open なのに不足内容が不明。 |

test / contract drift report は、`resolution=open` が 1 件でも残る場合、[test verification closure record set 固定契約](#test-verification-closure-record-set-contract) の `contract_drift_closure` と `final_open_item_count` を closed にできない。report を作れない場合は、対象外として扱わず、検証不足として扱う。

<a id="sec-0g-8-f"></a>
**[fixture 証跡責務 §0g.8-F fixture / testdata 配置固定契約](fixture.md#sec-0g-8-f)：**

以下の配置はfixture証跡責務上の配置契約である。未作成pathは、該当componentまたは該当fixtureの実装検証変更で作成するまで現行実体として扱わない。`必須配置` 列は、検証群ごとの fixture group root または fixture directory pattern を示す。正式 fixture directory は、fixture catalog の fixture 名、実在 directory 名、`manifest.json.name` が一致した path だけとする。group root 自体、未作成 path、または catalog と `manifest.json.name` が一致しない directory を正式 fixture directory と扱ってはならない。`testdata/admin/cli/`、`testdata/setup/`、`testdata/release/`、`testdata/mcp/` 配下の実在 fixture directory は、各固定契約の catalog 行と実装検証証跡で確定する。`testdata/release/`はPhase 9 release実装検証変更で作成済みである。`mcp`のfixture契約は [`docs/details/fixture.md` fixture 証跡責務 §30-F](fixture.md#mcp-fixture-contract) を正本とする。`testdata/mcp/`はPhase 10 mcp実装検証変更で作成済みであり、実体と fixture 証跡は [Phase 10 mcp 実装検証証跡](#phase-10-mcp-implementation-evidence) を参照する。

| 検証群 | 必須配置（fixture group root / pattern） | 必須内容 | 禁止条件 |
|--------|------------------------------------------|----------|----------|
| builder | `testdata/builder/single/`、`testdata/builder/site/`、`testdata/builder/empty-dir/`、`testdata/builder/strict/`、`testdata/builder/safe/`、`testdata/builder/url-safety/`、各 fixture の `expected/`。 | 入力 Markdown、テーマ設定、asset 入力、期待 HTML / CSS / JS / search index、期待 stdout / stderr、期待終了コード。 | 実行環境ごとに変わる絶対 path、timestamp、乱数、外部 URL 取得結果を期待値へ含めてはならない。 |
| runner | `testdata/runner/r1/`〜`testdata/runner/r34/`、各 fixture の `state/`、`github/`、`pipeline/`、`ssh/`、`notify/`、`expected/`。 | GitHub fake response、状態ファイル初期値、lock 状態、pipeline fake 結果、deploy fake 結果、通知 fake 結果、期待 `.last_sha`、期待 queue / snapshot。 | 実 GitHub API、実 SSH、実通知先、実 remote branch 状態に依存して合否を決めてはならない。 |
| API request lifecycle | `testdata/api/request-lifecycle/auth/`、`status/`、`history/`、`logs/`、`queue/`、`stream/`、`errors/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | HTTP method / path / query / header / body、状態ファイル初期値、期待 response、期待 error body、SSE frame、状態 read/write 後の期待値。 | API 運用群 endpoint、外部公開設定、仕様未定義 endpoint を fixture に含めてはならない。 |
| API 運用 | `testdata/api/operations/config/`、`notify/`、`snapshots/`、`maintenance/`、`hooks/`、`tokens/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | config / notify / snapshot / rollback / maintenance / hook / token の正常系、validation error、secret mask、API request lifecycle 群の回帰確認。 | token 原文、secret 原文、mask 前 payload、再取得不可 token の復元値を fixture または expected に含めてはならない。 |
| Admin | `testdata/admin/archive/`、`static-serving/`、`cli/`、`security/`、各 fixture の `input/`、`expected/`。 | archive entry、配布 file set、HTTP method / path / header / body、CLI request / stdout / stderr / exit code、既存 admin directory の維持、secret path 非配信。 | UI / SDK 内容生成、未定義配布 file、unsafe archive entry、directory listing、CLI 未定義 endpoint passthrough を許可してはならない。 |
| SDK | `testdata/sdk/request-shape/`、`error-shape/`、`stream/`、`binary/`、`operations/`。 | fake fetch transcript、期待 request、期待 SDK return、期待 `AdlaireCIError`、期待 stream event、timeout / abort の期待結果。 | Node.js 専用 API、bundler、npm package、実 network、browser storage 依存を検証前提にしてはならない。 |
| UI | `testdata/ui/login/`、`status/`、`build/`、`config/`、`secret/`、`stream/`、`operations/`。 | fake SDK script、入力 DOM 状態、操作手順、期待 DOM assertion、期待 SDK call、期待 disabled / loading / error / success 表示。 | 直接 `fetch()`、CDN、外部 framework、画像 snapshot だけの合否判定、secret 表示を含めてはならない。 |
| Statefile | `testdata/statefile/read/`、`write/`、`lock/`、`json-lines/`、`corrupt/`、`partial/`、各 fixture の `input/`、`expected/`。 | schema、mode、mtime、atomic write、lock、破損時処理、write order、forbidden write。 | caller 固有の業務判断、暗黙の自動修復、未定義状態 file を含めてはならない。 |
| Archive | `testdata/archive/log/`、`snapshot/`、`download/`、`delete/`、`rollback/`、各 fixture の `input/`、`expected/`。 | archive entry、checksum、圧縮・展開結果、stream、削除・rollback 境界、元 file 維持。 | unsafe entry、未検証展開、build 成否反転、元 build log 改変を許可してはならない。 |
| Commit status | `testdata/commitstatus/pending/`、`final/`、`disabled/`、`failure/`、各 fixture の `input/`、`expected/`。 | GitHub Status request、送信順、payload、失敗理由、build 成否非反転、secret mask。 | 実 GitHub write、Authorization 値保存、status 失敗による build 成否反転を含めてはならない。 |
| Security | `testdata/security/auth/`、`session/`、`token/`、`totp/`、`audit/`、`rate-limit/`、各 fixture の `input/`、`expected/`。 | memory-only state、hash-only state、scope、rate count、audit、one-time response、forbidden leak / write / call。 | password、token、ticket、TOTP secret、Authorization header の平文を expected に保存してはならない。 |
| Setup | `testdata/setup/<fixture-name>/`、各 fixture の `input/`、`expected/`。具体 fixture 名は [fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture 固定契約](#sec-27-f-19) の catalog を正本とする。 | Release asset、checksum、binary / admin 配置、systemd 操作、Go `net/http` health、既存 state / secret 保持、rollback。 | 実Release、実systemd、実network、release生成・公開処理そのものに依存してはならない。 |
| Release | `testdata/release/<fixture-name>/`、各fixtureの`input/`、`expected/`。 | clean checkout、tag / commit、Go build、再現性、admin archive、checksum、GitHub draft / asset / publish / cleanup、token mask。 | 実GitHub write、実tag変更、checkout変更、secret平文、host固有pathを含めてはならない。 |

<a id="sec-0g-8-f-2"></a>
**[fixture 証跡責務 §0g.8-F fake 実装固定契約](fixture.md#sec-0g-8-f-2)：**

| fake | 対象責務 | 必須動作 | 必須記録 |
|------|----------|----------|----------|
| fake GitHub server | runner | fixture の JSON 応答だけを返す。未定義 method / path は `404` とする。rate limit、`304`、`409`、`500` は fixture で明示された場合のみ返す。 | method、path、query、request body、認証 header の有無、呼び出し順。token 値は `***` に置換する。 |
| fake ssh executable | runner | fixture 指定の stdout、stderr、終了コード、timeout を返す。実 shell、実 SSH、実 file 転送は実行しない。 | argv、stdin 有無、環境変数名、終了コード、timeout 発生有無。secret 値は記録しない。 |
| fake notifier | runner / API 運用 | fixture 指定の HTTP status、response body、timeout を返す。通知先へ送信しない。 | URL の host 部分、payload schema、mask 後 payload、retry 回数、最終結果。 |
| fake filesystem | builder / runner / API | atomic write 失敗、sync 失敗、lock 競合、JSON 破損、permission error を fixture 単位で再現する。通常 file I/O の代替にはしない。 | 対象 path、操作種別、注入した失敗、復旧後の状態。 |
| fake fetch | builder / SDK | `status`、`headers`、`body`、network error、timeout、abort、stream chunk を fixture どおり返す。実 network は使用しない。 | method、URL、query、headers、body 有無、abort 発生有無、呼び出し順。token 値は `***` に置換する。 |
| fake browser runtime | builder | fixture の DOM tree、event、`matchMedia`、localStorage、Clipboard API、`execCommand`、timer、IntersectionObserver、scroll geometry、console を決定的に再現する。実 browser、実 storage、実 clipboard を使用しない。fetch は `fake fetch` を使用する。 | listener の target / type / options / 登録順、API 呼出しと引数、DOM の class / attribute / text / tree、storage read / write、timer 登録 / 発火、console warning を呼出し順で記録する。 |
| fake SDK | UI | SDK method ごとに固定 return、固定 throw、固定 stream event を返す。UI からの直接 API 呼び出しは受け付けない。 | method 名、引数、呼び出し順、throw した error code、stream unsubscribe 実行有無。 |

<a id="sec-0g-8-f-3"></a>
**[fixture 証跡責務 §0g.8-F 実装検証証跡固定契約](fixture.md#sec-0g-8-f-3)：**

この契約は全 component 実装の検証証跡に適用する。API 実装は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の追加仕様化機能実装は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の acceptance checklist も同時に満たす。

| 証跡 | 必須記載 | 不足時の扱い |
|------|----------|--------------|
| 変更対象 | owner component、collaborator component、変更ファイル、追加 fixture / testdata path。 | 対象責務の成果物不足として未完了。 |
| 固定契約 | 追加または固定した CLI、状態 schema、HTTP API、SDK method、DOM id、fake 動作、終了コード、error body。 | 依存 component が参照できないため未完了。 |
| 検証 | 実行コマンド、fixture 名、期待結果、実結果、判定。 | 合否を再現できないため未完了。 |
| 未実装対象 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務で今回の実装対象外または将来計画に割り当てられた機能、外部公開設定、未定義 endpoint / UI / 状態ファイルを列挙する。 | 先取り実装または範囲不明として未完了。 |
| 依存 component への影響 | 依存 component が利用許可済みの contract と、利用禁止の未固定 contract を列挙する。 | 依存境界未確認として未完了。 |
| secret 確認 | log、fixture、snapshot、UI 表示、実装検証証跡に secret / token / password 原文がないこと。 | security 不合格として未完了。 |

<a id="8a-f-builder-初期受け入れ-fixture-契約"></a>
**8a-F builder 初期受け入れ fixture 契約：**

[`docs/details/fixture.md` fixture 証跡責務 §8a-F](fixture.md#8a-f-builder-初期受け入れ-fixture-契約) は、[`docs/details/builder.md` 詳細本文責務 §8a](builder.md#8a-builder-受け入れ検証条件) の builder 初期受け入れ fixture、入力、expected、fake、実装検証証跡を扱う正本である。[`docs/details/builder.md` 詳細本文責務 §8a](builder.md#8a-builder-受け入れ検証条件) は検証観点だけを持ち、fixture 本文を再定義しない。

| fixture | 入力 / 実行 | expected / effects |
|---------|-------------|--------------------|
| Fixture A: 単一 Markdown 入力 | `testdata/builder/single/source.md` を `adlaire-ci-build --src testdata/builder/single/source.md --out <tmp> --title "Fixture Site"` で実行する。入力は H1、self anchor link、bash fence、task list、table、footnote を含む。 | exit `0`。`index.html`、`assets/style.css`、`assets/app.js`、`assets/search-index.json` を作成し、`pages/` は作成しない。H1 id、self link、code block label、task checkbox、`[REPORT] pages=1`、`theme=adlaire-default` を固定する。 |
| Fixture B: ディレクトリ Markdown 入力 | `testdata/builder/site/docs/intro.md`、`guide/setup.md`、`guide/setup_copy.md` を入力し、`adlaire-ci-build --src testdata/builder/site/docs --out <tmp> --title "Docs"` で実行する。 | exit `0`。目次 `index.html`、`pages/guide-setup.html`、`pages/guide-setup-copy.html`、`pages/intro.html` を出力する。入力順、目次順、Markdown link 変換、page ごとの slug 空間、search index entry を固定する。 |
| Fixture C: 異常系 | unknown theme、source 不在、Markdown 不在 directory、空 title を入力する。 | exit `2`。stderr は固定 error。stdout に `[REPORT]` を出さず、既存正常出力を変更しない。 |
| Fixture D: 冪等性 | 同一入力、同一 CLI 引数と同一 `--build-id` / `--commit-sha` / `--build-at` で 2 回連続実行する。 | HTML、CSS、JavaScript、search index、`[REPORT]` が byte 単位で一致する。現在時刻、file mtime、process 起動時刻に由来する差分を許可しない。 |
| Fixture E: path 安全性と既存出力保護 | source 配下 out、source と out 同一、10 MiB 超 Markdown を入力する。 | exit `2`。stderr を固定し、出力作成なしまたは既存 `index.html` 維持。 |
| Fixture F: HTML escape と Markdown 境界 | raw `<script>`、先頭 h3、h3→h1、h1→h3、h3〜h6、`#no-space`、7 個の `#` で始まる行、escaped 内部 pipe、escaped 行末 pipe、不正 separator、不足 / 超過 cell を持つ table block、7 レベル以上 list nesting を含む Markdown を入力する。 | raw HTML は escape される。h1〜h6 と対応 class を出力し、先頭 h3 と h3→h1 では heading skip warning なし、h1→h3 だけ `[WARN] HEADING_SKIP` 1 件かつ `heading_skips=1` とする。`#no-space` と 7 個以上の `#` は通常段落にする。内部 escaped pipe は cell 文字、escaped 行末 pipe の行と不正 separator の候補行は各 1 件の段落、不足 cell は空 cell 補完、超過 cell は最終 cell へ ` | ` 連結とする。list nesting clamp と `[WARN] LIST_NESTING_CLAMPED` を固定する。 |
| Fixture G: search index / JavaScript contract | Fixture B と同じ directory 入力を使用する。 | `assets/search-index.json` の top-level array、entry key 順、body 長、HTML tag 除外、`assets/app.js` の localStorage guard、`search-results`、`data-search-hit`、外部 storage / network 不使用を固定する。 |
| Fixture H: strict warning before publish | strict 用 Markdown と既存正常出力を用意し、`adlaire-ci-build --src testdata/builder/strict/source.md --out <tmp> --strict` を実行する。 | exit `2`。stdout に `[WARN] BROKEN_LINK`、`[WARN] UNCLOSED_FENCE`、`[REPORT]` を出し、stderr は空。公開 rename を 0 回とし、既存 `index.html` を置換しない。 |
| Fixture I: atomic output compensation | 既存正常出力を用意し、fake filesystem で (1) 既存出力から `prev` への rename 後の `tmp` 公開 rename、(2) `tmp` 公開 rename 後の親 directory `Sync`、(3) 各 subcase の補償 rename または補償 `Sync` を個別に失敗させる。(3) の後は異なる current PID で再実行する。 | (1) と (2) で補償成功時は exit `1`、`[REPORT]` なし、更新前 `index.html` の byte 一致、`tmp` / `prev` 不在、補償 1 回を固定する。(3) は exit `1`、`cannot restore previous output directory` と最初の失敗文言、追加復旧 0 回、残存 `tmp` / `prev` / 公開 path の実状態を `expected/state/state-diff.json` と `expected/effects.json` に固定する。異なる PID の再実行も親 directory の残存 entry を削除せず、UTF-8 byte 列で最小の絶対 path を含む `output staging path already exists` で停止する。 |
| Fixture J: build metadata dataflow | 固定 `--build-id b20260926010203-001 --commit-sha abcdef0123456789abcdef0123456789abcdef01 --build-at 2026-09-26T01:02:03Z` と、各 metadata option を省略した実行を個別に行う。 | 指定時は `BuildConfig`、`SiteData.BuildMeta`、全 HTML meta、`[REPORT]` の 3 値が完全一致する。省略時は 3 値を空文字のまま出力する。7 文字 SHA は `invalid commit sha` で拒否する。独立した生成日時を追加せず、fake clock を変えても生成物が変化しない。 |
| Fixture K: sidebar runtime | desktop / mobile の fake `matchMedia`、`adb-sb` の不在・`"1"`・`"0"`・不正値、storage read/write throw、toggle click、breakpoint change、mobile `.tl` click を個別に実行する。 | [`docs/details/builder.md` 詳細本文責務 §7.2](builder.md#sec-7-2) の mode 別 `open` / `closed` class、`aria-expanded`、日本語 `aria-label`、保存値、既定状態を exact 比較する。storage failure でも DOM 状態を維持し、`#ct` inline style、cookie、IndexedDB を変更しない。 |
| Fixture L: TOC filter runtime | ASCII 大文字、非 ASCII、内部空白、group / leaf、0 件結果を含む TOC を用意し、空 query → 非空 query 1 → 非空 query 2 → 空 query を実行する。 | [`docs/details/builder.md` 詳細本文責務 §7.4](builder.md#sec-7-4) の正規化、最初の snapshot 1 回、leaf / group `hidden`、group 一時展開、`#sb-none`、`searchActive`、空 query での完全復元を固定する。検索中の localStorage write は 0 回とする。 |
| Fixture M: copy runtime | Clipboard API fulfill、同期 throw、Promise reject、fallback `true`、`false`、throw を、`.cb-copy` と `.hn-link` の両方で fake timer とともに実行する。 | [`docs/details/builder.md` 詳細本文責務 §7.6](builder.md#sec-7-6) と [§7.11](builder.md#sec-7-11) の呼出し順、fallback 最大 1 回、textarea 必須削除、成功 / 失敗表示、1,800 ms 復元を固定する。`false` / throw を成功表示せず、二重 copy を行わない。 |
| Fixture N: syntax highlight runtime | `python`、`bash`、`json`、`sql`、`ini`、`diff` の優先競合、未閉鎖 string、対応外言語、HTML 風文字列を含む `<code>` を入力する。 | [`docs/details/builder.md` 詳細本文責務 §7.8](builder.md#sec-7-8) の token class と優先順を exact DOM tree で比較し、処理前後の `textContent` を UTF-16 code unit 単位で一致させる。`innerHTML` 使用、HTML 解釈、token 入れ子、対応外言語の DOM 変更を 0 件とする。 |
| Fixture O: search runtime | 正常 index、HTTP failure、JSON failure、1 entry schema failure、同一 / 別 pathname、fetch 完了前の query 更新、1 / 2 code point query、先頭空白を持つ text node、連続検索と clear を fake fetch / DOM で実行する。 | [`docs/details/builder.md` 詳細本文責務 §7.9](builder.md#sec-7-9) の fetch 1 回、index 全体破棄、最新 query だけの結果、entry 順最大 20 件、same-origin / same-path 判定、mark 除外要素、UTF-16 index を保持する `normalizeSearchMark()`、左から非重複 mark、解除後 `normalize()`、警告最大 1 回を固定する。index 由来値を `innerHTML` へ渡さない。 |
| Fixture P: keyboard runtime | 通常 target、各 editable / interactive target とその子孫、`defaultPrevented`、composition、repeat、Ctrl / Meta / Alt、Shift、空 / 非空検索を組み合わせて `/`、`Escape`、`t`、未定義 key を dispatch する。 | [`docs/details/builder.md` 詳細本文責務 §7.12](builder.md#sec-7-12) の handler 1 個、処理 / no-op、`preventDefault()`、focus、input event、smooth scroll の呼出し回数を固定する。`keyCode` / `which` 参照は 0 回とする。 |
| Fixture Q: scroll runtime | scroll range 正、0、負相当、scrollTop が範囲外、`scrollY` 400 / 401、各 DOM 欠落、IntersectionObserver 不在を fake scroll source で実行する。 | [`docs/details/builder.md` 詳細本文責務 §7.7](builder.md#sec-7-7) と [§7.13](builder.md#sec-7-13) の passive listener 1 個、fallback → progress → top button の順、初期同期実行 1 回、進捗 `0`〜`100` clamp、`.visible` 境界、欠落機能だけの no-op を固定する。 |
| Fixture R: table sort runtime | 正負整数、小数、数学的同値、text、空、`-0`、指数表記、同値行、複数 `<tbody>`、対象 cell 欠落を含む table で各列を 2 回 click する。 | [`docs/details/builder.md` 詳細本文責務 §7.14](builder.md#sec-7-14) の number / text / empty 分類、昇順 / 降順、empty 常時末尾、stable 元 index、UTF-16 比較、`aria-sort`、同じ `<tbody>` への再配置を固定する。不正構造は DOM / aria 無変更、`Number()` / `localeCompare()` 使用は 0 回とする。 |
| Fixture S: URL 属性安全性 | `testdata/builder/url-safety/source.md` に `http` / `https` link、`mailto`、`tel`、同一 page fragment、Markdown 相対 link、`.md` 以外の相対 link、Markdown image、`javascript:`、`data:text/html`、credential 付き `https`、protocol-relative URL、absolute path、base 外 `..`、backslash、制御文字、ASCII space を含む URL を入れる。non-strict と `--strict` を個別に実行する。 | 許可 URL だけが attribute escape 済み `href` / `src` になる。Markdown 相対 link は対応 HTML path へ変換し、fragment は変換先 page の slug と照合する。拒否 link は label text だけ、拒否 image は alt text だけまたは空出力になり、拒否 URL 値、credential、query、secret 風文字列が stdout、stderr、`[REPORT]`、HTML、search index に残らない。non-strict は `[WARN] UNSAFE_URL` を拒否件数分出し exit `0`、`warnings` に加算する。strict は warning 出力後 exit `2`、stdout に `[WARN]` と `[REPORT]`、stderr 空、`Done` 行なし、公開出力維持にする。`expected/security.json` は許可 scheme、拒否 scheme、credential 非表示、raw HTML 不在、external call 0 件を固定する。 |
| Fixture T: 起動入口 / version | `adlaire-ci-build --version`、`adlaire-ci-build --help`、`adlaire-ci-build --version --src missing`、`adlaire-ci-build --help --theme invalid` を実行する。binary version は未注入の `V.0.0-dev` と具体値 `V.1.100` の 2 case を持つ。 | stdout は [`docs/details/builder.md` 詳細本文責務 §2](builder.md#固定出力) の 1 行固定出力だけ、stderr 空、exit `0`。`--src` 存在確認、`--theme` 検証、directory 作成、atomic writer、fake filesystem、外部通信、状態 read/write は 0 回。version は exact 3 token で、binary name は `adlaire-ci-build`、第 2 token は注入値と完全一致する。 |
| Fixture U: CLI parse / validation order | `--unknown`、位置引数、短縮 `-s`、未許可 `--src=...`、値欠落、unsafe argv token、重複 `--src` / `--out` / `--base-dir` / `--build-id` / `--commit-sha` / `--build-at` / `--title` / `--theme`、重複 `--strict` を個別 case として実行する。重複値 option は最後の値だけが存在確認対象になるよう、最後より前に存在しない path または secret 風文字列を置く。 | unsafe argv token は owner parse 前に `invalid command line token`、未知 option は `unknown option: <token>`、値欠落は `missing value: --name` を固定する。重複値 option は最後の値だけで成功または失敗し、最後より前の値を stdout、stderr、HTML、search index、`[REPORT]`、`expected/effects.json` へ出さない。`--strict --strict` は single strict と同じ結果にする。parse failure は stdout 空、exit `2`、file read/write、atomic writer、fake filesystem、外部通信 0 回。 |
| Fixture V: path / symlink / base-dir | `testdata/builder/safe/` 配下に通常 Markdown、hidden Markdown、symlink Markdown、symlink directory、大文字拡張子、10 MiB 超 file、CRLF / CR / UTF-8 BOM file、base 外 relative link を持つ fixture を置く。`--base-dir` 省略、明示空文字、存在しない path、通常 file、symlink directory、`--src` symlink、既存 `--out` symlink、`--out` 親 symlink、`--out` 親不存在、`--src` と `--out` の同一・相互包含を個別 case として実行する。 | hidden Markdown、symlink Markdown、symlink directory、大文字拡張子は入力収集対象外。CRLF / CR は LF、先頭 BOM は除去、本文途中 BOM は保持する。明示空 `--base-dir` は `base directory must not be empty`、`--src` symlink は `source is not markdown file or directory: <path>`、`--base-dir` symlink は `base path is not directory: <path>`、既存 `--out` symlink は `output path is not directory: <path>`、`--out` 親 symlink は `output parent is not directory: <path>`、親不存在は `output parent not found: <path>`、相互包含は `output path must be outside source: <path>` を固定する。failure case は公開出力を置換せず、`[REPORT]` を出さない。 |

`TOC ハイライト追従（アクティブ見出し追跡）` の runtime fixture は、同一挙動をここへ重複定義せず [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#sec-28-f-14) の [`docs/details/builder.md` 詳細本文責務 §28.16](builder.md#sec-28-16) 対応行を唯一の正本とする。

Fixture A〜V は fixture ごとに `manifest.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json` を持つ。HTML / CSS / JavaScript / search index を生成する成功 case と strict warning case は `expected/site/` を持つ。公開出力の置換、維持、staging cleanup、残存 path を判定する case は `expected/builder-output.json` を持つ。secret、credential、unsafe URL、argv token、metadata、外部通信禁止を扱う case は `expected/security.json` を持つ。対象外 expected は作成せず、`manifest.json.not_applicable` に対象外理由を固定する。

<a id="15a-f-runner-初期受け入れ-fixture-契約"></a>
**15a-F runner 初期受け入れ fixture 契約：**

[`docs/details/fixture.md` fixture 証跡責務 §15a-F](fixture.md#15a-f-runner-初期受け入れ-fixture-契約) は、[`docs/details/runner.md` 詳細本文責務 §15a](runner.md#15a-runner-受け入れ検証条件) の runner 初期受け入れ fixture、入力状態、expected、fake GitHub / fake ssh / fake notifier / fake filesystem、実装検証証跡を扱う正本である。[`docs/details/runner.md` 詳細本文責務 §15a](runner.md#15a-runner-受け入れ検証条件) は検証観点だけを持ち、fixture 本文を再定義しない。

**runner fixture 共通 expected：**

| 共通 expected | 固定内容 |
|---------------|----------|
| no external execution | 対象 fixture が禁止する GitHub API、pipeline、deploy、snapshot を実行しない。 |
| no log history creation | 対象 fixture が禁止する `.build_logs/` と `.build_history` を作成しない。 |
| clean final state | `.build_state.running=false`、`current_build_id=null`、`.build_lock` 不在で終了する。 |
| finalizer failure | ERROR ログ `BUILD_STATE_FINALIZE_FAILED` を出し、`.build_lock` は runner owner が所有確認後に 1 回だけ解放する。解放成功 / 失敗に関わらず、`.build_state.running=false` 保存失敗を正常扱いしない。 |

| fixture | 入力 / fake | expected / effects |
|---------|--------------|--------------------|
| Fixture R1 | CLI help、version、relative state-dir、unknown option、NUL / CR / LF / C0 制御文字 / DEL / invalid UTF-8 を含む argv token、runner 固有の重複 option case を実行する。重複 option case は `--state-dir /tmp/a --state-dir /tmp/b`、`--state-dir relative --state-dir /tmp/b`、`--state-dir /tmp/a --state-dir relative`、`--once --once`、`--dry-run --dry-run`、`--help --state-dir relative`、`--version --state-dir relative` を個別 case とする。 | help / version は exit `0`。不正 CLI は exit `2`、stdout 空、stderr 固定。unsafe argv token は owner 固有 parse より先に `invalid command line token` で拒否し、token 原文を expected、log、state へ保存しない。重複 `--state-dir` は最後の値だけを検証対象にし、最後より前の値を読取、存在確認、directory 作成、log、state、expected へ保存しない。最後が相対 path の case だけ exit `2` とし、最後より前の絶対 path を読まない。重複 `--once` は通常 oneshot、重複 `--dry-run` は通常 dry-run と同一扱いにする。help / version precedence case は `.github_token`、lock、状態ファイル、GitHub、pipeline、deploy、snapshot をすべて 0 call とする。CLI parse 失敗 case は lock、state、GitHub、pipeline、deploy、snapshot、`.build_logs/`、`.build_history` をすべて 0 write / 0 call とする。 |
| Fixture R2 | GitHub single file target の `.last_sha={"sha":"1111111111111111111111111111111111111111"}`、fake GitHub Trees API が同じ 40 文字 SHA を返す。 | exit `0`。`.last_sha` の 40 文字値を byte 不変で維持し、log/history 非作成、`NO_CHANGE` INFO。 |
| Fixture R3 | 旧 SHA、fake GitHub 変更あり、fake blob Markdown、deploy target なし、pipeline success、`snapshots_keep=0`。 | exit `0`。新 SHA 保存、build log/history success、archive owner 非呼出し、`snapshot_id=null`、clean final state。 |
| Fixture R4 | fake GitHub 変更あり、pipeline exit `7`。 | exit `1`。旧 SHA 維持、failure build log/history、deploy/snapshot 非実行。 |
| Fixture R5 | pipeline success、deploy target 1 件、引用対象文字を含む remote path、fake ssh transfer failure。 | local argv の `ssh`, `--`, `user@host`, remote command の順、`quoteRemoteArg` 適用後の `mkdir` / `tee`、stdin を固定する。exit `1`、新 SHA 保存、pending transfer 追加、build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、両方の `failure_category="deploy_failure"`、snapshot 非作成。 |
| Fixture R6 | 実行中 PID を指す `.build_lock`。 | exit `0`。状態、log、history を変更せず already running を出力する。 |
| Fixture R7 | `.notify_pending` が破損 JSON。 | corrupt backup を作成し、`.notify_pending=[]` を再生成して通常処理を継続する。timestamp は UTC 秒精度。 |
| Fixture R8 | build timeout `1`、pipeline が timeout まで終了しない。 | exit `1`。旧 SHA 維持、pipeline timeout log/history、clean final state。 |
| Fixture R9 | fake GitHub Trees API が retry 対象 `503` を返し続ける。 | exit `3`。旧 SHA 維持、failure_api log/history、pipeline/deploy/snapshot 非実行。 |
| Fixture R10 | pipeline success、notify webhook fake `500`。 | exit `0`。build success 維持、notify pending/log を保存し、history を failure にしない。 |
| Fixture R11 | pipeline stdout に `[REPORT]` 2 行。 | exit `0`。1 行目だけ report 保存、`REPORT_DUPLICATE` warning を保存する。 |
| Fixture R12 | finalizer 時だけ `.build_state` atomic write failure。 | exit `1`。成功 log/history/status は保存済み、finalizer failure を固定する。 |
| Fixture R13 | `.github_token` mode `0644`。 | exit `2`。`GITHUB_TOKEN_INSECURE_MODE`、running にせず、token 値 / 長さ / hash を出力しない。 |
| Fixture R14 | dry-run、repo/dist/log/snapshot directory 不在。 | exit `0`。directory 非作成、dry-run JSON だけ出力、外部副作用なし。 |
| Fixture R15 | `.last_sha` 破損、fake GitHub 変更あり。 | exit `1`。SHA cache 維持、failure_decode、pipeline/deploy/snapshot 非実行。 |
| Fixture R16 | fake GitHub `403` rate limit remaining `0`、reset header 不正。 | exit `3`。旧 SHA 維持、failure_api、長時間待機なし。 |
| Fixture R17 | cooldown 中、または manual force waiting entry。 | polling は skipped_cooldown。manual force は waiting から active へ atomic move し、事前に SHA cache を空値化せず cooldown と SHA 一致 skip を無視して build し、成功 log/history 保存後だけ新 SHA を保存して active を消去する。 |
| Fixture R18 | pipeline success、fake ssh transfer success、fake ssh checksum mismatch。 | `sha256sum --zero -- {quoted_remote_path}` の argv と NUL 終端 record を固定する。exit `1`、新 SHA 保存、pending transfer 保存、build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、両方の `failure_category="deploy_failure"`、snapshot 非作成。 |
| Fixture R19 | `source_kind="output"` の既存 pending transfer と同一転送先の deploy failure。 | entry 件数を増やさず、既存 entry の build id、source fields、`retry_count=0`、failed_at、last_error を更新し、投入順を保持する。 |
| Fixture R20 | pipeline/deploy success、`snapshots_keep=2`、`meta.json.saved_at` が異なる既存の schema-valid snapshot 2 件。 | archive owner 呼出し 1 回、`site.tar.gz` と `meta.json` を tmp 内で検証後に atomic publish、tmp 非残存、禁止 source 非含有、tar entry 辞書順、`size_bytes` / `file_count` / `output_sha256` 再計算一致、`saved_at` が最古の snapshot だけ tombstone 経由で prune、build log / history の `snapshot_id=build_id` を固定する。 |
| Fixture R21 | `.build_status.json` start write だけ fake failure。 | exit `1`。running にせず、外部副作用と log/history 作成なし、lock 削除。 |
| Fixture R22 | 外部副作用前の running build log atomic write だけ fake failure。 | exit `1`。history 非追記、旧 SHA 維持、GitHub API / pipeline / deploy / snapshot 非実行、status `failure_state_write`。 |
| Fixture R23 | 最終 build log 保存成功後の history append だけ fake failure。 | exit `1`。最終 build log、新 SHA、成功済み deploy / snapshot を維持し、history は追記せず自動 retry しない。status `failure_state_write`、clean final state。 |
| Fixture R24 | branch target 2 件、1 件目 GitHub API failure、2 件目 success。 | exit `1`。1 件目 failure_api、2 件目 success。1 件目失敗で全体中断しない。各 target に異なる build id を割り当て、log 2 件と history 2 行を同じ順序で保存し、成功 target の snapshot と output SHA をその target の `Out` だけから作る。 |
| Fixture R25 | branch target 2 件、両方 GitHub API failure。 | exit `3`。全 target failure_api、全 SHA 旧値維持、pipeline/deploy/snapshot 非実行。各 target に異なる build id を割り当て、log 2 件と重複 ID のない history 2 行を保存する。 |
| Fixture R26 | pipeline/deploy success、archive owner の snapshot save fake が `failed` を返す。 | exit `0`。build success 維持、build log / history の `snapshot_id=null`、`SNAPSHOT_SAVE_FAILED` warning、pending transfer 追加なし。 |
| Fixture R27 | success 保存後、finalizer `.build_state` atomic write だけ fake failure。 | exit `1`。success log/history/status は保持し、finalizer failure を固定する。 |
| Fixture R28 | manual active entry の build 中に process crash を発生させ、次回 runner を起動する。 | 初回は active を保持する。次回は waiting の urgent entry より先に同じ active payload を新しい build id で再実行し、結果確定後だけ active を消去する。 |
| Fixture R29 | queue build の final log atomic replace だけを失敗させる subcase と、final log 保存成功後の history append だけを失敗させる subcase を実行する。両 subcase で `.build_state` 読取と finalizer atomic write は成功させる。 | exit `1`。active と waiting を保持し、`running=false`、`current_build_id=null`、`last_finished_at={fake now}` を 1 回の atomic write で保存する。次回 runner は同じ active payload を新しい build id で waiting より先に再実行する。 |
| Fixture R30 | circuit open、active entry 1 件、waiting entry 1 件、stale `running=true` / `current_build_id`、pending transfer / notify 各 1 件で runner を起動する。 | pending transfer / notify retry だけを実行し、status は skipped / circuit_open。runtime flag は false / null へ保存し、active / waiting と最終時刻を保持する。queue 取得、build id 採番、GitHub、pipeline、deploy、snapshot、build log / history は実行せず、保存成功時 exit `0`、runtime flag 保存失敗時 exit `1`。 |
| Fixture R31 | SHA / deploy / snapshot 成功後の final build log atomic replace だけ fake failure。 | exit `1`。running build log を残し、history 非追記、新 SHA、deploy 結果、snapshot を巻き戻さず、status `failure_state_write`、running false、所有確認付き lock 解放を 1 回だけ実行することを固定する。 |
| Fixture R32 | `.repo_config` 不在、正常値、API 更新前に起動した process、更新後の次回 process、破損 JSON を個別に用意する。 | 不在時は `fqwink/Build-Scripts`、正常値は保存済み owner/repo を全 GitHub API path に使用する。起動済み process は更新前値を維持し、次回 process は更新後値を使う。破損時は exit `2`、`REPO_CONFIG_INVALID`、GitHub API / queue / build / pipeline / deploy / snapshot 呼出しと状態自動修復は 0 件。 |
| Fixture R33 | branch target 2 件を同一 runner process で成功させ、各 target は異なる `Out` と複数 `target_files` を持つ。fake clock は同一秒を返し、1 件目の build id と衝突させる。 | target ごとに suffix を含む一意 build id を 1 件ずつ採番する。同一 target の `target_files` は 1 ID に集約する。`.build_state.current_build_id` は target 順に遷移し、log、history、snapshot、output SHA は各 2 件かつ target の branch / output root と一致し、上書きと同一 history ID を 0 件とする。process 終了時は current ID を null にして exit `0`。 |
| Fixture R34 | stdout / stderr に CRLF、単独 CR、NUL、invalid UTF-8、既知 secret、1 MiB 超過、4000 文字超行を含める。REPORT なし、正常 1 行、重複、key 欠落、順序違い、型不正、必須 13 key 後の拡張 key を個別に実行し、`[WARN]` と `[WARNING]` も含める。 | LF 正規化、NUL 可視化、U+FFFD 置換、secret mask、UTF-8 rune 境界を保つ末尾 1 MiB 以下、truncation flag、行長制限を固定する。REPORT なしは `REPORT_MISSING`、重複は最初の 1 行と `REPORT_DUPLICATE`、不正は `report:null` と `REPORT_PARSE_FAILED`、正常と拡張 key は 13 key だけの Report object とし、いずれも元 pipeline exit code を変更しない。 |
| Fixture R35 | GitHub directory target 配下に path と blob SHA が異なる `.md` file 3 件を fake Trees API 順序を変えて返し、正規化済み path 順 payload から期待 64 文字 digest を事前算出する。`.last_sha` には同 digest を保存する。 | API 応答順に依存せず [`docs/details/runner.md` 詳細本文責務 §13](runner.md#13-処理フロー) の directory target digest が期待 64 文字値と一致する。exit `0`、`.last_sha` byte 不変、log/history/pipeline/deploy/snapshot 非実行、`NO_CHANGE` INFO。 |
| Fixture R36 | local mode の対象 file 2 件と空 target を個別に用意し、相対 path / file SHA-256 payload から期待 64 文字 digest を事前算出する。変更あり subcase は旧 64 文字値を保存する。 | [`docs/details/runner.md` 詳細本文責務 §27.23](runner.md#sec-27-23) の path 順 local target digest、空 target digest、dry-run `current_blob_sha`、build 成功後の 64 文字保存値を exact 比較する。GitHub API / Commit Status 呼出しは 0 件とする。 |
| Fixture R37 | `.server_config` sparse object に `build_timeout_seconds=1`、`snapshots_keep=0`、`schedule_paused=false`、`allowed_hours=null`、`build_cache_enabled=true` を保存し、標準 builder command が成功する target を実行する。別 subcase では `allowed_hours.from == allowed_hours.to` と `allowed_hours.from > allowed_hours.to` の保存済み config を個別に用意する。 | 正常 subcase は `RunnerConfig.BuildTimeoutSeconds=1`、`SnapshotsKeep=0`、`SchedulePaused=false`、`AllowedHours=null`、`BuildCacheEnabled=true` として 1 回だけ正規化し、builder command timeout、`--cache-dir` argv、archive owner 非呼出し、build log / history success を固定する。`from >= to` の各 invalid subcase は exit `2`、`CONFIG_ALLOWED_HOURS_INVALID`、GitHub API、queue、build id、pipeline、deploy、snapshot、build log / history を 0 件とする。 |
| Fixture R38 | queue entry なしの polling で `schedule_paused=true`、queue entry なしの polling で UTC hour が `allowed_hours` 外かつ target SHA が変更あり、manual queue entry ありで同じ pause / allowed hours 条件を個別に実行する。 | pause polling は `.build_status.json` に `status="skipped"`、`last_target_status="skipped_schedule_paused"`、`trigger="polling"` を保存し、GitHub API、pipeline、deploy、snapshot、build log / history は 0 件とする。allowed hours 外 polling は GitHub tree resolve と SHA decision で変更ありを確定した後、`skipped_allowed_hours` を保存し、blob API、pipeline、deploy、snapshot、build log / history、SHA 更新は 0 件とする。manual queue entry は schedule pause gate と allowed hours gate を無視して active move、build id 採番、SHA decision へ進み、pause / allowed hours を queue 完了条件に使用しない。 |

<a id="22-f-api-fixture-契約"></a>

**22-F API fixture 契約：**

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) は、API の必須検証、fixture 名、入力状態、期待 response、期待副作用を扱う fixture 証跡責務である。API endpoint の method、path、request、response、error、read / write 境界は [`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) を正本とする。

API 実装の検証証跡は、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに加えて、[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の検証群、endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、endpoint 固有の業務状態非変更、共通 security / observability 副作用、secret mask を記録する。

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) における「状態差分なし」「no-write」「read-only」は、endpoint 固有の業務状態に対する禁止を意味する。[`docs/details/api.md` 詳細本文責務 §22.0 GET の副作用](api.md#sec-22-0) が許可する 5 群の共通副作用は、該当する場合に `expected/state/`、`expected/logs/`、`expected/effects.json` へ明示し、省略または禁止扱いにしてはならない。

実装順序と現在の割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#roadmap-initial-phase-plan)、実装変更単位と完了判定単位は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#policy-phase-unit) を参照する。[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) は API fixture 証跡として記録する項目だけを固定する。

| 検証群 | 対象 | 完了条件 |
|--------|------|----------|
| API CLI dispatch | `--help`、`--version`、listener、signal、graceful shutdown、`--init-credentials`、入力異常、mode 競合 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](../DETAIL_INDEX.md#common-cli-contract)、[`docs/details/api.md` 詳細本文責務 `api` CLI 固定契約](api.md#api-cli-contract)、[API listener lifecycle 固定契約](api.md#api-server-lifecycle-contract) どおりに mode、stdout、stderr、終了コード、server 値、listener / shutdown / close / credentials 呼出回数、副作用が一致する。 |
| API request lifecycle | 認証、session、共通 error、状態ファイル読み書き、`.access_log`、`.config_log`、build 操作、status、logs、history、queue、circuit breaker | `POST /api/login` から認証必須 API の共通処理、手動 build、強制 build、cancel、queue、history、log 取得までが [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0)〜[§22.0e](api.md#sec-22-0e) と一致し、秘密情報が log と response に出ない。 |
| API 運用 | config、repo、branch、schedule、PAT、diagnostics、dashboard、notify、SMTP、webhook、snapshot、rollback、maintenance、access control、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout、tokens | 運用 API が schema どおり状態を保存し、secret mask、GET の endpoint 固有業務状態非変更、共通 security / observability 副作用、rollback / maintenance / hook / token の副作用が fixture と一致し、SDK と UI の操作名が [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と一致する。 |

各検証群の検証条件は以下とする。

| 検証群 | 必須検証 |
|--------|----------|
| API CLI dispatch | help/version優先、未知 option、禁止位置引数、値欠落、state-dir検証、loopback address検証、listener mode、server 固定値、listen failure、`SIGTERM` / `SIGINT`、shutdown success / timeout / failure / forced close、init-credentials mode、mode競合、password stdin、全失敗時副作用ゼロを確認する。 |
| API request lifecycle | 認証成功、認証失敗、期限切れ session、`401` 時 SDK token 破棄、`.access_log` 追記、秘密情報 mask、手動 build、force build、running 中の queue、cancel、history/log 取得、`409`、`429`、`503` を確認する。 |
| API 運用 | config/repo/branch/schedule の保存、`.config_log` と対象 audit の追記、GET 系 API の endpoint 固有業務状態非変更と共通副作用、Webhook test、weekly summary、SMTP test、webhook secret 保存、secret mask、snapshot list/download/delete、rollback、maintenance enable/disable、access control block、hook success/failure、rule 追加/削除、pipeline config 保存、notes 保存、dashboard layout 保存、token 発行/失効、token 本体が再取得不可であることを確認する。 |

<a id="sec-22-f"></a>
**[fixture 証跡責務 §22-F API request lifecycle fixture 固定契約](fixture.md#sec-22-f)：**

API request lifecycle 実装の完了証跡は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いと [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の固定表に従う。fixture は実装言語の test case 名または subtest 名へそのまま写せる粒度とし、期待 HTTP status、期待 body、状態ファイル副作用を同時に確認する。

| Fixture | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------------------|---------------|--------------------|
| A0 API CLI dispatch | fake listener、fake credentials initializer、既存 state directoryを用意する。help と未知 option、version と相対 state-dir、正常 listener、正常 init stdin、未知 option、位置引数、各値欠落、state-dir の未指定・空・相対・不在・通常 file・symlink、addr の hostname・IPv6・wildcard・port 0・65536・先頭 0、init と addr の同時指定、listener mode の credentials 不在・読取不能・schema 不正を個別に実行する。 | help/version は固定1行、stderr空、exit `0`。listenerは `127.0.0.1:8765` または指定した有効loopback addressでfake listener 1回、exit `0`。initはsecurity固定出力とexit。入力異常とcredentials起動時検証失敗は該当固定stderr 1行、stdout空、exit `2`。 | help/version/入力異常はstdin read、listener、credentials initializer、状態read/write、乱数、外部通信を0件とする。credentials起動時検証失敗は対象fileを修復、退避、上書きせずlistener 0件。listener modeはcredentials起動時readとlistener各1回、状態write 0件。init modeはlistener 0件、initializer 1回とし、password平文がargv、env、expected output、logに存在しない。 |
| A0a API listener lifecycle | fake listener へ `ReadHeaderTimeout=5s`、`ReadTimeout=30s`、`WriteTimeout=0`、`IdleTimeout=60s`、`MaxHeaderBytes=32768`、redacted `ErrorLog` を渡す。signal target は `api_process`、shutdown timer target は `api_shutdown` とする。listen failure、正常 serve 後の `SIGTERM`、正常 serve 後の `SIGINT`、shutdown success、elapsed `9999ms` の success、elapsed `10000ms` と同時の success、shutdown failure、close failure、shutdown 中の 2 回目 signal、raw server error log を個別に実行する。 | listen failure は stderr `listen failed` + LF、exit `1`。最初の signal 後に elapsed `<10000ms` で `Shutdown=nil`、server result `http.ErrServerClosed` となる success は stdout / stderr 空、exit `0`。elapsed `10000ms` 同時 success / timeout / shutdown failure / close failure / 2 回目 signal は stderr `shutdown failed` + LF、exit `1`。raw server error は server log の `API_HTTP_SERVER_ERROR` / `WARN` 1 件だけへ変換する。 | listen は各 case 1 回。最初の signal は shutdown 1 回、success は close 0 回、timeout / shutdown failure / close failure は close 1 回。2 回目 signal は新しい shutdown を開始せず、close 済みなら追加 close 0 回。raw message、client address、request byte、Go error、path、header、credential、新規 request、未定義 retry、業務状態 write、実 port bind、実 signal、実 sleep は 0 件。 |
| A1 route errors | 未定義 `/api/unknown`、既存 path への未許可 method を順に送る。 | `404 {"error":"Not found"}`、`405 {"error":"Method not allowed"}`。 | endpoint 固有の状態ファイルを作成、更新、削除しない。各 request の `.api_access_log` だけを共通処理順どおり 1 行追記し、access control、認証、rate limit、body 読取、外部呼び出しは開始しない。 |
| A1a authenticated body errors | 認証済み管理 session で body 禁止 endpoint への body 付き request、JSON 不正文、1 MiB 超過 body を個別に送る。 | `400 {"error":"Request body is not allowed"}`、`400 {"error":"Invalid JSON"}`、`413 {"error":"Payload too large"}`。 | access control、session、rate limit は [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2)、`.api_access_log` は [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6)、状態 read/write は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) どおり行う。`.access_log`、`.audit_log`、endpoint 固有の状態ファイルを変更せず、外部 API、command、通知を呼び出さない。 |
| A1b request read limits | header を 5 秒間完了しない接続、32768 bytes を超える header、認証済み管理 session の endpoint で body 読取を 30 秒超過させる request を個別に fake する。 | header timeout は handler / API response なしで connection close、header 超過は `431` と connection close、body timeout は `408 {"error":"Request timeout"}`。 | header 2 case は request ID、`.api_access_log`、認証、rate limit、状態差分なし。body timeout は完了済み session / rate limit 副作用を維持し、`.api_access_log.status=408`、endpoint 固有状態、`.access_log`、`.audit_log`、外部呼出しは差分なし。 |
| A2 auth errors | Bearer なし、無効 token、期限切れ session、scope 不足 API token で認証必須 endpoint を呼ぶ。 | `401 {"error":"Unauthorized"}` または `403 {"error":"Forbidden"}`。 | 秘密情報を response、`.access_log`、server log に出さない。 |
| A2a request id | fake entropy から固定 16 bytes を返す正常 request、未知 path、認証失敗 request を実行し、続けて entropy failure を発生させる。 | 正常系、`404`、`401` は 32 文字 lowercase hex の `X-Request-Id` を返し、対応する `.api_access_log` と同じ値を持つ。entropy failure は `500 {"error":"Internal server error"}` かつ `X-Request-Id` なし。 | 正常 request で audit / config log を作る場合は同じ request id を保存する。entropy failure は endpoint 状態、request log、audit log、config log を変更せず、server log は `API_REQUEST_ID_GENERATION_FAILED` 1 件だけとする。 |
| A3 status primary | 正常な `.build_status.json`と lock 不在、valid active lock、stale lock を個別に用意して `GET /api/status`。 | [StatusObject 固定契約](api.md#status-object-contract) の 10 key だけを返す。lock 不在は status の `running`、valid active lock は `running:true`、stale lock は status の `running` を維持する。`output_url`、`queued` は返さない。 | 参照対象の業務状態ファイルは mtime、mode、内容を変更しない。認証、rate limit、`.api_access_log` の共通副作用は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) どおり記録する。 |
| A3a status conflict lock | 形式不正 lock と PID 判定不能 lock を個別に用意し、正常な `.build_status.json` で `GET /api/status`。 | どちらも `200`、StatusObject の `running:true`。その他の 9 key は status file の写像値。 | lock の削除、上書き、退避、status / state の修復を行わない。 |
| A4 status fallback | `.build_status.json` 不在、`.build_history` に status summary 対象の成功履歴、pending 各 1 件、circuit open、`.build_state.running=false`、`.build_lock` 不在で `GET /api/status`。 | history から `last_sha`、`last_build_at`、`last_build_status`、`last_target_status`、`last_trigger`、状態から両 pending 件数、`circuit_open`、`running:false` を算出し、`last_deploy_status:null`。必須 10 key 以外は返さない。 | `.build_status.json` を生成しない。queue の有無を status response に反映しない。 |
| A4a status fallback history-only exclusion | `.build_status.json` 不在、`.build_history` の末尾時刻に `approval_rejected`、その直前に status summary 対象の成功履歴を置いて `GET /api/status` と `GET /api/health`。 | history-only 行を除外し、成功履歴から `last_build_status` と `last_target_status` を算出する。summary 対象行がない別ケースでは `last_build_status:"none"`。 | `.build_status.json`、`.build_history`、approval 状態を変更しない。 |
| A5 status corrupted | `.build_status.json` を不正 JSON にして `GET /api/status`。 | `500 {"error":"State file is corrupted"}`。 | 退避ファイル作成、自動修復、初期値上書きを行わない。 |
| A6 history jsonl filtering | `.build_history` に有効行 2 件、空行 1 件、不正 JSON 1 件、必須 key 不足 1 件を置き `GET /api/history?page=1&per_page=10`。 | `total:2`、`pages:1`、`history` は有効行だけを `finished_at` 降順、同時刻は id 降順で返す。 | 壊れた行を書き戻し削除しない。server log に `BUILD_HISTORY_SKIP_CORRUPT`。 |
| A6a paged log query | `.access_log`、`.api_access_log`、`.notify_log`、`.config_log` に同一 `at` を含む有効行 102 件と破損行 1 件ずつを置き、各 GET を query なし、`limit=1&offset=1`、`limit=0`、`limit=1001`、`offset=-1`、整数でない値で呼ぶ。 | query なしは `at` 降順、同時刻は物理行順の逆順の先頭 100 件、`limit=1&offset=1` は同固定順の 2 件目だけを返す。API access log の `total` は paging 前の有効行数。範囲外または整数でない値は `422 {"error":"Validation failed","details":[...]}`。 | 参照対象 log を作成、修復、更新、削除せず、mtime、mode、内容を変更しない。破損行ごとに対応する `*_SKIP_CORRUPT` 固定 code と行番号だけを server log へ記録する。 |
| A6b history order and duplicate id | `.build_history` に同じ `finished_at` で id が異なる有効行、同じ id の後続有効行、壊れた行を置いて `GET /api/history`。 | `finished_at` 降順、同時刻は id 降順。重複 id は最初の物理行だけを返し、`total` と paging は一意な有効行から算出する。 | file を書き換えず、壊れた行に `BUILD_HISTORY_SKIP_CORRUPT`、重複 id 行に `BUILD_HISTORY_DUPLICATE_ID` と line number を記録する。 |
| A6c history failure category | 既知 category、`failure_category:null`、列挙外値 `future_pipeline_error`、形式不正値 `Future-Error` の履歴を置き、filter 未指定、`failure_category=pipeline_timeout`、未知 query で `GET /api/history`。 | 列挙外値と形式不正値の行は破損行として除外する。既知 filter は完全一致行だけ、未知 query は `422` と `details[].field="failure_category"`。 | `.build_history`、build log、status を変更しない。除外行ごとに `BUILD_HISTORY_SKIP_CORRUPT` を記録する。 |
| A6d corrupt hook and search logs | 破損 hook log 1 件と正常 hook log 1 件で `GET /api/hooks/{id}/log` を実行し、破損通常 build log、gzip 展開失敗 archive、正常 build log で `GET /api/logs/search` を実行する。 | 両 endpoint とも破損対象を除外し、正常対象だけを固定順で返す。 | hook は `HOOK_LOG_SKIP_CORRUPT`、hook id、basename、検索は破損対象ごとに `LOG_SKIP_CORRUPT`、build id、basename だけを記録する。内容、絶対 path、Go error を記録せず、対象 file を更新・削除・修復しない。 |
| A7 build log lookup | `.build_logs/{id}.json` 正常、通常ログ不在で archive 正常、両方不在、破損ログをそれぞれ `GET /api/history/{id}/log`。 | 正常は log object、archive は gzip 展開結果、両方不在は `404 {"error":"Not found"}`、破損は `500 {"error":"State file is corrupted"}`。 | archive を通常ログへ復元しない。破損ログを上書きしない。 |
| A7a finite build stream | `.build_state` 不在、running/current id の `finished_at:null`, `duration_seconds:null` の途中保存 log あり・なし、running/current id 矛盾 state、idle で同時刻 id が異なる通常 / archive log、破損最新 log を用意し `GET /api/build/stream`。stdout に途中空行、末尾 LF、4001 code point の行を含める。 | state 不在は idle とする。running は current id だけを選び、`log.at=started_at`、end は `status:"running"`, `duration_seconds:null`。未保存は過去 log へ fallback せず `404`、矛盾 state は SSE header なしの JSON `500`。idle は `finished_at` 降順、id 降順、通常 log 優先で 1 件を選択する。各 frame は `data: {compact-json}\n\n`、stdout→stderr→warnings→error の log frame 後に end frame 1 件を送り close する。途中空行は保持、末尾空要素は除外、超過行は UTF-8 を壊さず 4000 code point へ切り詰める。 | wait、poll、tail、fallback write、archive 復元、状態更新を行わない。idle の破損候補は `BUILD_STREAM_SKIP_CORRUPT` を記録し次候補へ進む。 |
| A7b response writer semantics | JSON response で `Write` を先に呼ぶ case、`WriteHeader(202)` 後に `WriteHeader(500)` を呼ぶ case、SSE writer が `http.Flusher` を持つ / 持たない case、frame 全 write / partial write / client disconnect を個別に fake する。 | implicit status は `200`、明示 status は最初の `202` を維持する。Flusher ありは完全な各 frame 後に 1 回 flush、なしは実装済みと偽装しない。partial write / disconnect は追加 JSON response と後続 frameなし。 | `.api_access_log.status` は最初に確定した status と一致する。partial write / disconnect は `API_RESPONSE_WRITE_FAILED` を 1 件だけ記録し、業務状態、JSON Lines、body の再送、重複 flush を発生させない。 |
| A8 queue state | `.build_state` 不在、active / waiting がある正常状態、破損の 3 状態で `GET /api/queue`。 | 不在は `active:null,queued:[]`、正常は active と waiting を分離した保存値、破損は `500 {"error":"State file is corrupted"}`。 | 不在時も `.build_state` を生成しない。 |
| A9 build conflict / queue | valid running `.build_lock` と空き queue、形式不正 `.build_lock`、PID 判定不能 lock、queue 上限到達状態で `POST /api/build`。 | valid running かつ空きありは `202` queue 追加。判定不能 lock は `409 {"error":"Conflict"}`。queue 上限到達は `429`。 | `409` / `429` では `.build_state`、`.build_history`、`.build_status.json` を変更しない。valid running では queued だけを更新する。 |
| A10 circuit breaker | `.build_circuit_state` が `open:true` の状態で `POST /api/build`、続けて `POST /api/circuit-breaker/reset` を 2 回実行する。 | build は `409 {"error":"circuit_open"}`。reset は 2 回とも `200 {"message":"Circuit breaker reset","open":false,"consecutive_failures":0}`。 | build 拒否では queue、audit、runner 起動要求なし。1 回目の reset は初期値を atomic write し、`circuit_breaker_reset` config log と config update audit を順に 1 件ずつ追記する。2 回目は状態、config log、audit に差分なし。どちらも build / runner を起動しない。 |
| A11 write lock timeout | `.build_state.lock` を保持した状態で `.build_state` 更新 endpoint を呼ぶ。 | 10 秒経過後 `409 {"error":"Conflict"}`。 | tmp file を残さず、target を変更しない。 |
| A12 GET endpoint-state no-write | `GET /api/status`、`GET /api/history`、`GET /api/logs`、`GET /api/queue` を連続実行する。 | 各 endpoint は入力状態に応じた正常 response または固定 error response。 | 参照対象の業務状態ファイル一覧、mtime、mode、内容は request 前後で一致する。`.api_access_log` は request ごとに 1 行、認証 session と rate limit 状態は入力条件に応じて共通契約どおり更新し、それ以外の write / external call は行わない。 |
| A13 manual dispatch | active / waiting のない idle 状態で `POST /api/build`、fake systemctl exit `0`。 | `202 {message:"Build queued",queue_id,queued:true,dispatch:"requested"}`。`build_id` と `queued:false` は存在しない。 | manual entry、trigger audit を保存し、systemctl を固定引数で 1 回呼ぶ。running / current_build_id / active / SHA / status は変更しない。 |
| A14 dispatch fallback | idle 状態で `POST /api/build/force`、fake systemctl 非 `0` または 5 秒 timeout。 | `202`、`queued:true`、`dispatch:"timer_fallback"`。 | force entry と audit を保持し、`RUNNER_ACTIVATION_DEFERRED` を固定形式で 1 件記録する。retry、queue rollback、running 遷移、SHA 更新なし。 |
| A15 zero queue dispatch slot | `queue_max_size=0`、active なし、waiting なしで 1 回目と 2 回目の異なる manual request を送る。 | 1 回目は `202`、2 回目は `429 queue_full`。 | 1 件目だけ waiting に保存し、2 件目で state 差分なし。 |
| A16 health read / response output failures | `.build_status.json` permission failure、fallback history read failure、`.pending_transfers` permission failure を個別および同時に fake して `GET /api/health`。別 subcase で JSON serialize 失敗、body write の全失敗 / partial write、client 切断を fake する。 | 読取失敗は HTTP `200`、固定順の `build_status_read_error`、`build_history_read_error`、`pending_transfers_read_error`、`last_build_status:"none"`、`pending_transfers:0`。serialize 失敗は元 header / status / `HealthObject` なしの `500 {"error":"Internal server error"}`。送信開始後の write 失敗 / client 切断は追加 response なし。 | chmod、backup、修復、初期化、業務状態 write を行わない。serialize 失敗は `API_RESPONSE_ENCODE_FAILED` ERROR と最終 status `500` の `.api_access_log`、送信開始後失敗は `API_RESPONSE_WRITE_FAILED` WARN をそれぞれ 1 件だけ記録する。 |

<a id="sec-22-f-2"></a>
**[fixture 証跡責務 §22-F API 運用 fixture 固定契約](fixture.md#sec-22-f-2)：**

API 運用実装の完了証跡は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いと [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の固定表に従う。fixture は既存 endpoint と既存状態ファイルだけを対象とし、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) にない endpoint、[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) にない状態ファイル、[`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) にない UI 操作を追加してはならない。

| Fixture | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------------------|---------------|--------------------|
| B1 config no-op | 既存 `.server_config` と同じ body で `POST /api/config`。 | `200` と `{message:"No changes",config}`。 | `.server_config`、`.config_log`、`.audit_log` を変更しない。 |
| B1a config sparse normalization | `.server_config` 不在、`{}`、`{"log_level":"DEBUG"}`、未知 key 付き、型不一致を個別に用意し、`GET /api/config` を呼ぶ。 | 不在と `{}` は全 key が既定値の `ConfigObject`、partial object は `log_level:"DEBUG"` 以外が既定値の `ConfigObject`。未知 key と型不一致は `500 {"error":"State file is corrupted"}`。 | file の作成、完全形への書き戻し、退避、mtime / mode / key 順の変更を行わない。 |
| B1b config normalized write | partial `.server_config` に既定値と異なる許可 key を `POST /api/config`、専用 schedule API、rate-limit API でそれぞれ更新する。 | 更新 response は全 key を持つ正規化後 `ConfigObject`。 | 各成功更新は `.server_config` の全 top-level key を保存し、nested object の必須 key も省略しない。未指定値は読取時の正規化値を維持する。 |
| B2 config validation failure | `queue_max_size=-1`、未知 enum、相対 path を含む `POST /api/config`。 | `422 {"error":"Validation failed","details":[...]}`。 | 状態ファイルを変更しない。 |
| B2a repo config ownership | `.repo_config` 不在で `GET /api/repo-info`、続けて `{owner:"example",repo:"docs-site"}` で `POST /api/repo-config`、再度 GET。 | 初回 GET は `{owner:"fqwink",repo:"Build-Scripts",updated_at:null}`、POST は `{message:"Repo config updated"}`、再 GET は更新後 owner/repo と fake clock の `updated_at`。 | 初回 GET は file を作成しない。POST は 3 key の `.repo_config`、`.config_log`、`config_update` audit を固定順で保存し、`.branch_config` を変更しない。次回 runner 起動が新 owner/repo を使う。 |
| B2b repo config rejection / no-op | 現在値と同じ owner/repo、`branch`、`target_file`、`updated_at`、未知 key、空 object、不正 owner/repo を個別に `POST /api/repo-config`。 | 同値は `200 {message:"No changes"}`。それ以外は `422 {"error":"Validation failed","details":[...]}`。 | no-op と `422` は `.repo_config`、`.branch_config`、`.config_log`、`.audit_log` の byte、mtime、mode を変更しない。 |
| B2c repo config read failure | 破損 JSON と permission read failure の `.repo_config` を個別に用意し、`GET /api/repo-info` と `POST /api/repo-config` を呼ぶ。 | 破損は `500 {"error":"State file is corrupted"}`、読込不能は `500 {"error":"State file read failed"}`。 | 既定 owner/repo で補完せず、破損 file の退避、削除、上書き、`.config_log`、`.audit_log`、`.branch_config` の変更を行わない。 |
| B2d repo config backup / restore | `.repo_config` 不在と非既定値保存済みの 2 状態で backup し、非既定値への restore、固定既定値 `{owner:"fqwink",repo:"Build-Scripts",updated_at:null}` への restore、`updated_at:null` の非既定値、branch / target key 付きを個別に実行する。 | 不在時 backup は固定既定値、非既定値は 3 key を保持する。有効 restore は成功、固定既定値は file 削除または no-op、その他の不正入力は `422`。 | 非既定値 restore は 3 key object を atomic write し、固定既定値 restore は `.repo_config` だけを lock 下で削除する。`422` は全状態を不変に保ち、有効差分だけが [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の restore 保存順と [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の audit 契約に従う。 |
| B3 branch config save | 有効な branch target 2 件で `POST /api/branch-config`。 | `{message:"Branch config updated",branches_count:2}`。 | `.branch_config` を atomic write し、`.config_log` と `config_update` audit を 1 件ずつ記録する。 |
| B4 schedule systemd failure | `.server_config` 保存成功後、systemd timer 更新を fake failure にする。 | `500`。 | `.server_config` は更新済み、`.config_log` に `result="partial_failure"` / `error="systemd_update_failed"`、`.audit_log` に `config_update` / `result="failure"` を 1 件ずつ記録し、rollback しない。 |
| B5 PAT update secret mask | `POST /api/pat-update` に token を送る。 | `{message:"PAT updated"}`。 | secret file は mode `600`。response、`.config_log`、`.audit_log`、server log に token 平文を出さない。 |
| B6 dashboard read-only | `.dashboard_layout`、`.build_state`、`.build_status.json`、`.build_lock`、両 pending、circuit state、`.alert_rules` を置き `GET /api/dashboard`。 | widget 順に dashboard object を返す。`dashboard.status` は同じ入力の `GET /api/status` と key、型、値が完全一致し、`output_url`、`queued` を含まない。 | GET は参照対象の業務状態ファイルを作成、修復、更新しない。認証、rate limit、`.api_access_log` の共通副作用は別 path の expected に固定する。 |
| B7 diagnostics systemd | fake command で timer active / service unit installed、timer inactive、service unit missing、5 秒 timeout を順に返して `GET /api/diagnostics`。 | 成功時は `systemd` item が `ok` と固定 message、残りは `error` と固定 message。oneshot service の active 状態を要求しない。 | command 引数と 5 秒 timeout を固定し、shell 呼び出し、stdout / stderr 保存、診断結果保存、endpoint 固有業務状態の更新を行わない。共通副作用だけを expected に記録する。 |
| B8 output target exact match | branch 名順とは異なる `last_branch` / `last_target_file` を持つ `.build_status.json`、異なる `out` を持つ branch target 2 件、非選択 target のより新しい成功履歴、選択 target の現在と一つ前の成功履歴 / log / 出力 site を置き、全出力参照 endpoint を呼ぶ。 | `last_branch` と `last_target_file` の両方が一致する target の `out`、成功履歴、同じ id の log / archive だけから sysinfo、output-meta、dashboard、diagnostics、disk-usage、verify-output を算出し、`size_diff_bytes` は現在 build id の一つ前の同 target 成功履歴と比較する。 | 非選択 target の出力、履歴、log を fallback に使用せず、参照対象の作成、修復、更新、退避を行わない。 |
| B9 output target deterministic fallback | branch 名の入力順が逆の branch target 2 件と、`.build_status.json` 不在、破損、target 不一致の 3 状態を用意し、target 選択だけを必要とする出力参照 endpoint を呼ぶ。 | 3 状態とも branch 名の byte 昇順で先頭の target を request 中固定する。status 自体を応答する dashboard は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) の破損時契約に従う。 | `.build_status.json` と `.branch_config` を作成、修復、退避、書き戻しせず、request 中に target を切り替えない。 |
| B10 selected output missing | 有効な API 選択出力 target と履歴を置き、選択 target の `out` だけを不在にして各出力参照 endpoint を呼ぶ。 | sysinfo / dashboard は size `0` / mtime `null`、diagnostics は `output_file` error item、disk-usage は `output_file_bytes:0`、output-meta / verify-output は `404 {"error":"Not found"}`。 | 出力 directory を作成せず、非選択 target へ fallback せず、全業務状態ファイルを不変に保つ。 |
| C1 notify config mask | Webhook secret と SMTP password を含む通知設定保存後、GET / backup / log を確認する。 | secret は `"***"` または `*_set:true` だけを返す。 | secret 平文を状態表示、履歴、通知ログ、backup に残さない。 |
| C2 webhook receive signed | 正常署名の GitHub push payload を `POST /api/webhook`。 | `202` と queued 結果。 | queue 投入条件を満たす場合は `.build_state.queued` へ `trigger:"webhook"` を追加した後、`.webhook_events.json` に queued 結果を追記する。 |
| C3 webhook invalid signature | 署名なし、不正 prefix、不一致署名。 | `401 {"error":"Unauthorized"}`。 | event log、queue、history を変更しない。 |
| C4 SMTP test disabled | SMTP disabled で `POST /api/smtp-test`。 | endpoint 固有の `422`。 | `.notify_log` へ成功扱いを残さず、secret を出力しない。 |
| D1 snapshot delete | 存在する snapshot id で `DELETE /api/snapshots/{id}`。 | `{message:"Snapshot deleted"}`。 | 対象 snapshot だけ削除し、endpoint 契約どおり `.config_log` と対象 audit を記録する。 |
| D2 rollback running conflict | `.build_state.running=true` で `POST /api/history/{id}/rollback`。 | `409 {"error":"Build is running"}`。 | queue、history、snapshot、deploy target を変更しない。 |
| D3 maintenance blocks build | maintenance enabled 状態で `POST /api/build`。 | `503` と endpoint 固有 maintenance error。 | build queue、history、log を変更しない。 |
| D3a maintenance corrupt fail-closed | `.maintenance` 破損状態で `POST /api/build`。 | `503 {"error":"maintenance_unavailable"}`。 | `.maintenance` 自動修復 / backup、queue、build id、history、log、SHA cache を変更しない。 |
| D3b maintenance corrupt runner stop | `.maintenance` 破損状態、正常な state directory、書込可能な `.build_status.json` で runner 定期起動。 | `.build_status.json` の atomic write を 1 回実行して status failure / `config_error` / `startup_config_integrity` を保存し、終了コード `2`。 | `.maintenance` 自動修復 / backup、pending retry、build log、history、SHA cache、deploy、snapshot を変更しない。 |
| D4 access-control deny | allowlist に接続元が含まれない状態で任意認証必須 API。 | 認証判定前に `403 {"error":"Forbidden"}`。 | password / token 検証、access log 成功行、対象操作副作用を行わない。 |
| D4a access-control corrupt fail-closed | `.access_control` 破損状態で `GET /api/health` 以外の endpoint。 | body 読取と認証判定前に `503 {"error":"Access control unavailable"}`。 | `.access_control` 自動修復 / 上書き、password / token 検証、rate state、endpoint 固有副作用を発生させない。 |
| D5 hook timeout | `pre` hook が timeout。 | build status は `hook_error` または endpoint 固有の hook error。 | pipeline を実行せず、hook log と build log に timeout を固定値で記録する。 |
| E1 token issue once | `POST /api/tokens` で token 作成。 | token 本体を作成 response に 1 回だけ含める。 | `.api_tokens` には hash だけを保存し、再取得 API では token 本体を返さない。 |
| E2 token revoke missing | 存在しない token id を `DELETE /api/tokens/{id}`。 | `404 {"error":"Not found"}`。 | `.api_tokens`、`.audit_log` を変更しない。 |
| E3 alert/tag duplicate | 同一 alert rule または tag rule を 2 回作成。 | 2 回目は `409 {"error":"Conflict"}`。 | 2 回目は該当状態ファイル、`.config_log`、`.audit_log` を変更しない。 |
| E4 pipeline config reserved arg | `extra_args` に [`docs/details/statefile.md` 詳細本文責務 予約 builder option 固定契約](statefile.md#pipeline-config-reserved-builder-options) の各 exact option と各 `--name=value` 形式を個別に含める。 | 各 request は `422 {"error":"Validation failed","details":[...]}`。 | 各 request で `.pipeline_config`、`.config_log`、`.audit_log` を変更しない。 |
| E5 notes same content | 同じ `content` を 2 回 `POST /api/notes`。 | 2 回目は `No changes`。 | 2 回目は `.notes`、`.config_log`、`.audit_log` を変更しない。 |
| E6 dashboard layout invalid | 重複 widget、未知 widget、空配列を `POST /api/dashboard-layout`。 | `422 {"error":"Validation failed","details":[...]}`。 | `.dashboard_layout` を変更しない。 |

<a id="sec-22-f-3"></a>
**[fixture 証跡責務 §22-F API 機能別 fixture 固定契約](fixture.md#sec-22-f-3)：**

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の API fixture 固定表は、API component の機能別 fixture 名、入力、期待結果を固定する。API endpoint の処理順序、request / response、状態ファイル read / write 境界は [`docs/details/api.md`](api.md) 詳細本文責務を基準とし、API fixture 固定表では fixture 本体だけを定義する。

| 機能 | fixture | 入力 | 期待結果 |
|------|---------|------|----------|
| backup / restore | backup-mask | secret 設定済みで backup | secret 本体なし、`*_set:true`。 |
| backup / restore | server-config-sparse-roundtrip | `.server_config` 不在、`{}`、partial object、完全 object を個別に backup し、取得した `BackupObject` を restore する。 | 不在時の backup は `{}`。存在時は top-level key 集合と値を保持する。restore は正規化後値が同一なら no-op、差分があれば request の sparse / 完全形を保持し、再取得の正規化値が一致する。 |
| backup / restore | restore-validate-fail | 1 file schema 不正 | `422`、全 file 差分なし。 |
| backup / restore | restore-secret-keep | `"***"` かつ既存 secret あり | 既存 secret 維持、平文出力なし。 |
| backup / restore | restore-secret-missing | `"***"` かつ既存 secret なし | secret 未設定のまま、`"***"` を保存しない。 |
| backup / restore | restore-secret-delete | secret `null` | secret file 削除。 |
| backup / restore | restore-write-failure | 中途 write 失敗 | 未処理 file は差分なし、処理済み file は維持、`500`。 |
| メンテナンス | maintenance-build-deny | enabled 中に `POST /api/build` | `503`、queue 差分なし、build id なし。 |
| メンテナンス | maintenance-force-deny | enabled 中に `POST /api/build/force` | `503`、SHA cache 差分なし。 |
| メンテナンス | maintenance-webhook-deny | enabled 中に署名済み webhook | `503`、event log と queue 差分なし。 |
| メンテナンス | maintenance-disable-noop | disabled 中に disable | `200 No changes`、状態ファイル差分なし。 |
| メンテナンス | maintenance-corrupt-read | `.maintenance` 破損中に `GET /api/maintenance` | `500`、既存 file 上書きと corrupt backup なし。 |
| メンテナンス | maintenance-corrupt-build-deny | `.maintenance` 破損中に `POST /api/build` | `503 maintenance_unavailable`、queue / build id / SHA cache / log / history 差分なし。 |
| メンテナンス | maintenance-corrupt-runner-stop | `.maintenance` 破損中に runner 定期起動 | 終了コード `2`、可能な場合だけ status failure / `config_error`、pending retry / build / deploy なし。 |
| アクセス制御 | access-allow-empty | `.access_control.allow=[]` | 任意 IP の API が認証処理へ進む。 |
| アクセス制御 | access-deny-before-auth | allow 不一致 IP で `POST /api/login` | `403`、`.access_log`、`.audit_log`、rate state 差分なし。 |
| アクセス制御 | access-corrupt-fail-closed | `.access_control` 破損中に non-health API と `GET /api/health` | non-health は `503 Access control unavailable`、health は access control を読まず health 契約の status。自動修復、認証、rate state、endpoint 固有副作用なし。 |
| アクセス制御 | access-normalize | 重複 allow を保存 | sort / 重複除去後の配列を返し `.config_log` 記録。 |
| アクセス制御 | access-ipv6-reject | IPv6 literal を保存 | `422`、状態差分なし。 |
| hooks | hook-pre-success | pre hook exit 0 | pipeline 実行、hook log 保存、secret mask 済み。 |
| hooks | hook-pre-abort | pre hook exit 1 / abort true | pipeline 未実行、status `hook_error`、history に `failure_category:"hook_error"`。 |
| hooks | hook-pre-warn | pre hook exit 1 / abort false | build 継続、hook log に exit code。 |
| hooks | hook-post-failure | build success 後 post hook が非 0 で終了する。 | build status を変更せず、WARN log と hook log を保存する。 |
| hooks | hook-timeout | timeout 超過 | process kill、`timed_out:true`、`exit_code:null`。 |
| hooks | hook-log-write-failure | pre hook log 保存失敗 | build 本体未実行、`hook_error`。 |
| SMTP | smtp-save-password | password 付き保存 | `.smtp_secret` mode `0600`、GET は `password_set:true`、log は `"***"`。 |
| SMTP | smtp-delete-password | `password:null` | `.smtp_secret` 削除、password 平文なし。 |
| SMTP | smtp-noop | 同一 config / password 未指定 | 状態差分なし、`.config_log` 追記なし。 |
| SMTP | smtp-secret-write-failure | non-secret config 差分あり / なしの 2 subcase で `.smtp_secret` atomic write を fake failure | 差分ありでは先行する `.smtp_config` 保存を維持し、差分なしでは `.smtp_config` write 0 回。両方とも `.config_log` / `.audit_log` 0 回、response `500`、password 平文なし。 |
| 通知 test | notify-test-success | 有効な Webhook channel 複数、fake 2xx | channel id byte 昇順の先頭 1 件だけへ固定 payload を送信し、response は `{message:"Test notification sent",channel_id:<選択 channel id>}` と完全一致し、Webhook URL / secret を含まない。`.notify_log` に `event:"notify_test"`、対象 channel id、`result:"success"`、`http_status:200`、`error_code:null` を追記する。 |
| 通知 test | notify-test-failure | fake 1xx / 3xx / 4xx / 5xx / timeout / response 前接続失敗の各 subcase | `.notify_log` の `event:"notify_test"`、`result:"failure"`、`http_status`、`error_code` が statefile 正本に一致し、response は `500` かつ Webhook URL / secret を含まず、`.notify_pending` 差分なし。3xx は redirect 先への外部呼出し 0 件。 |
| SMTP | smtp-test-success | 設定済み test | `.notify_log` に `event:"smtp_test"`、`channel_id:"smtp-test"`、`channel_type:"email"`、`result:"success"`、`attempt:1`、`http_status:null`、`error_code:null` と固定 payload hash、response success。 |
| SMTP | smtp-test-failure | SMTP 接続・認証・送信失敗 / timeout の各 subcase | `.notify_log` の `result:"failure"`、`error_code:"smtp_error"` / `"timeout"`、mask 後固定 `error`、`.notify_pending` 差分なし。 |
| SMTP | smtp-test-disabled | `enabled:false` | `422`、`.notify_log` 差分なし。 |
| SMTP | smtp-log-failure | test 後 `.notify_log` 追記失敗 | `500`、password 平文なし。 |
| queue | queue-add-idle-dispatch | idle / active なしで manual build | waiting append、created_seq 最大 + 1、systemctl 固定引数 1 回、running / active / build id 差分なし。 |
| queue | queue-add-running | active / running 中に manual build、queue に空きあり | waiting append、created_seq 最大 + 1。active / running / current_build_id 維持。 |
| queue | queue-duplicate-waiting | waiting にある entry と同一の manual payload を再投入 | 新規追加なし、waiting の既存 queue_id と `queued:true` を返し、同じ id の runner 起動要求を 1 回再実行する。 |
| queue | queue-duplicate-active | active にある entry と同一の manual payload を再投入 | 新規追加なし、active の既存 queue_id と `queued:true` を返し、同じ id の runner 起動要求を 1 回再実行する。response から waiting と誤判定せず、続く `GET /api/queue` が active を返す。 |
| queue | queue-full | max_size 到達 | `429 {"error":"queue_full"}`、差分なし。 |
| queue | queue-clear | active 1 件、waiting 2 件で `DELETE /api/queue` | `cleared_count=2`、active / running / current_build_id 維持、waiting だけ空配列。 |
| queue | queue-runner-activate | urgent と normal が waiting に混在、active なし | urgent を waiting から active へ atomic move し running / current_build_id を設定、他 waiting 維持。 |
| queue | queue-active-retry | active normal、waiting urgent、valid running lock なし | active を先に新しい build id で再実行し、urgent は waiting に保持する。 |
| queue | queue-active-finalize | active build の log / history 最終保存成功 | 同じ atomic write で active null、running false、current_build_id null。waiting 維持。 |
| 認証 | auth-password-failure | 誤 password で `POST /api/login` | `401`、session/ticket なし、正規化 IP key の memory-only 失敗回数 +1、credentials 差分なし、`.access_log` と `.audit_log` に secret なし。 |
| 認証 | auth-login-lock | 同一 IP key の 10 回目 password 不一致、直後の request、rate window 超過後かつ lock 期限内の request、lock 期限一致の request | 10 回目は `401 Unauthorized`、直後は pre-auth rate limit の `429 Too many requests`、fake clock を 60 秒超 10 分未満へ進めた後は hash 検証なしの `429 Too many attempts`、`now >= locked_until` で entry 削除後に hash 検証を再開する。 |
| 認証 | auth-session-issued | TOTP 無効で password 成功、`must_change` と `login_count` 境界値 | token は response のみ、`.admin_credentials.login_count` は飽和 +1、`false` は `none`、`true` の増加後 1〜4 は `prompt`、5 以上は `forced`、ログに token/hash/salt なし。 |
| 認証 | auth-session-expired | 期限切れ session で保護 API | `401`、対象 session 削除、`.access_log` と `.audit_log` は追記しない。 |
| 認証 | auth-password-change | `forced` session で password 変更成功 | 新 salt/hash、`must_change:false`、現 session の `password_change_required=false`、現 session 以外削除、`password_change` ログ、password/hash/salt 平文なし。同じ現 session で通常 endpoint を使用できる。 |
| 認証 | auth-log-write-failure | login 成功時に `.audit_log` 追記失敗 | `500`、session token を response しない。 |
| 認証 | auth-concurrent-transactions | harness barrier で同時 password login 2 件、login と password change、旧 revision の TOTP ticket と password change、同時 TOTP setup 2 件、TOTP confirm と disable、revoke-all と通常認証、期限切れ session と logout の coordinator 取得順を各 subcase で固定する。 | request は [認証 transaction・並行更新固定契約](security.md#security-auth-concurrency-contract) の取得順に直列化する。login count の lost update 0、旧 password / ticket / setup secret の成功 0、revoke 後の次 request は `401`、二重 session 削除と二重 audit 0。 |
| 認証 | auth-transaction-timeout | 先行 request が auth transaction coordinator を保持し、後続 request の timer target `auth_transaction_wait` を elapsed `59999ms` と `60000ms` へ進め、各時刻の coordinator 解放を個別に実行する。 | `59999ms` 解放は後続 request が取得して処理を継続する。`60000ms` 同時解放は timeout を優先して `409 {"error":"Conflict"}`。timeout case の body parse、credentials / TOTP read、entropy、memory mutation、state write、`.access_log`、`.audit_log` は 0 件。共通 `.api_access_log` だけ最終 status `409` で記録する。 |
| approval | approval-create | approval_required target に差分 | build なし、`requested_force` を含む pending record、`approval_pending` audit、通知成功または pending。 |
| approval | approval-duplicate | branch/sha/target/requested_trigger/requested_force/delivery_id が同一の pending を再検出 | pending 重複作成なし。成功済み audit / channel 通知は重複しない。 |
| approval | approval-pending-audit-recovery | pending 作成後の `approval_pending` audit 追記を失敗させ、同一要求を再検出 | pending は 1 件のまま、不足 audit を 1 件追記し、証跡のない channel だけ通知する。 |
| approval | approval-approve | pending approve | queue 追加、approved record、queue_id 保存、`approval_approved` audit。 |
| approval | approval-reject | pending reject | rejected record、history `approval_rejected`、`approval_rejected` audit。 |
| approval | approval-timeout | expires_at 超過 | expired record、history `approval_expired`、`approval_expired` audit。 |
| approval | approval-queue-full | max_size 到達時 approve | `429`、status pending 維持。 |
| approval | approval-duplicate-notify | 同一 branch / SHA / target の pending approval を再検出する。 | pending record を重複作成せず、通知を送らない。 |
| approval | approval-approved-append-failure | approve による queue 追加後、approved record の保存を失敗させる。 | queue は残り、`approval_approved` audit は追記せず、API は `500` を返す。runner は queue を取り出さない。再 approve は既存 queue id を再利用し、queue 件数を増やさない。 |

<a id="sec-22-f-4"></a>
**[fixture 証跡責務 §22-F API / SDK / UI / 状態ファイル cross fixture 固定契約](fixture.md#sec-22-f-4)：**

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の固定表は、API endpoint、SDK method、UI 操作、状態ファイル副作用の横断整合を固定する。API / SDK / UI / 状態ファイルの本文は、それぞれ [`docs/details/api.md`](api.md) 詳細本文責務、[`docs/details/sdk.md`](sdk.md) 詳細本文責務、[`docs/details/ui.md`](ui.md) 詳細本文責務、[`docs/details/statefile.md`](statefile.md) 詳細本文責務を参照する。fixture 詳細本文では、横断 fixture の入力、副作用有無、期待結果だけを固定する。

| fixture | 入力 | 必須確認 |
|---------|------|----------|
| cross auth expired | 任意の認証必須 API が `401`。 | SDK は token を破棄し、UI は全 secret field を消去して `panel-login` だけを表示する。対象 API の状態ファイル副作用なし。 |
| cross config no-op | `POST /api/config` に既存値と同一の正規化済み body。 | API は `No changes`、SDK は response をそのまま返し、UI は成功表示する。`.server_config`、`.config_log`、`.audit_log` に差分なし。 |
| cross validation details | 任意の保存 API が `422 details`。 | SDK は `AdlaireCIError.details` を保持し、UI は該当 field と panel summary に表示する。状態ファイルを書かない。 |
| cross refresh failure | 変更 API は成功し、成功後再取得の 2 件目が `500`。 | 変更副作用は維持し、同じ変更 API を再実行しない。UI は操作成功を `global-success`、再取得失敗を panel error に分けて表示する。 |
| cross secret failure | secret 保存 API が `500`。 | SDK error に secret 原文を含めず、UI は secret field を消去する。状態ファイル、log、fixture に secret 原文が残らない。 |
| cross stream invalid frame | `GET /api/build/stream` が chunk 境界をまたぐ正常 frame の後、別ケースで parse 不能 frame、`end` なし EOF、`end` 後の追加 frame を返す。 | 正常 chunk は 1 frame に復元する。各不正ケースは `StreamHandle.error` に `AdlaireCIError(status=0,message="Invalid SSE frame")` を保持し、`done` を同じ error で reject する。UI は stream error を 1 回表示し、status / queue をこの順で再取得する。状態ファイルは変更しない。 |
| cross stream invalid response | `GET /api/build/stream` が `2xx` で media type 不一致、body 不在、readable body 不在を個別に返す。 | `streamBuild()` は `StreamHandle` を返す前に `AdlaireCIError(status=0,message="Invalid SSE response")` で reject する。callback 呼出しと UI の stream 開始表示は 0 回とし、UI は接続 error を 1 回表示する。 |
| cross stream callback failure | 正常 `log` frame で `onLine`、別ケースの正常 `end` frame で `onEnd` が例外を投げる。 | reader と connection を停止し、`StreamHandle.error` に `AdlaireCIError(status=0,message="Stream callback failed")` を保持し、`done` を同じ error で 1 回だけ reject する。callback 原始 error は SDK error、console、UI、fixture expected に転写しない。UI は stream error を 1 回表示し、status / queue をこの順で再取得する。 |
| cross binary snapshot | `downloadSnapshot(id)` が binary success。 | API は binary header、SDK は `Blob`、UI は download 開始表示。JSON parse、success JSON body、状態ファイル更新なし。 |
| cross destructive cancel | 削除 / rollback 確認 dialog を cancel。 | SDK method 呼び出し 0 回、状態ファイル副作用なし、success / error 表示差分なし。 |
| cross rollback conflict | `rollbackHistory(id)` が `409 {"error":"Build is running"}`。 | SDK は `AdlaireCIError.status=409` を保持し、自動 `getStatus()` と同一 rollback request の再送を行わない。 |
| cross token issue once | `createToken()` が token 本体を含む作成 response を返す。 | SDK は token を response として返すだけとし、内部保存、console 出力、token list 合成を行わない。UI は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の one-time display 契約に従う。 |
| cross duplicate rule | `addAlertRule()` または `addTagRule()` が `409 {"error":"Conflict"}`。 | SDK は `AdlaireCIError.status=409` を保持し、自動 retry を行わない。UI は既存 rule 表示を保持する。 |
| cross paged log query | `getAccessLog()`、`getNotifyLog()`、`getConfigLog()` を引数なしで呼び、続けて `{limit:1,offset:1}` で呼ぶ。 | SDK は引数なしでは `limit=100&offset=0`、指定時は `limit=1&offset=1` を各 endpoint へ送る。API response の配列順と内容を変更せず返し、UI は受信順を保持する。状態ファイル副作用なし。 |

<a id="27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約"></a>

**27-F fixture 証跡責務 / runner・security 実装検証証跡詳細契約：**

<a id="sec-27-f"></a>
**[fixture 証跡責務 §27-F 配置・命名固定契約](fixture.md#sec-27-f)：**

追加仕様化機能の fixture は、実装者が実行順や期待値を推測しないように、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の単位で配置する。実装ファイル名やテストフレームワークは fixture 証跡責務では固定しないが、fixture の入力、期待出力、期待副作用は機能単位で分離する。

| 機能範囲 | fixture 配置単位 | 必須内容 | 禁止条件 |
|----------|------------------|----------|----------|
| [§27.1〜§27.4](fixture.md#sec-27-f) | `status/`、`dry-run/`、`retry/`、`output-meta/` 相当の機能別単位。 | runner 入力、GitHub fake response、builder 入力、期待 log/history/status/report。 | 実際の GitHub API、実時刻依存の期待値、secret 平文。 |
| [§27.5〜§27.11](fixture.md#sec-27-f) | `config/`、`access-log/`、`archive/`、`schedule/` 相当の API 機能別単位。 | HTTP request、初期状態、期待 response、期待状態差分、期待 log。 | API response だけの検証、状態差分未確認。 |
| [§27.12〜§27.20](fixture.md#sec-27-f) | `webhook/`、`stats/`、`snapshot/`、`health/`、`branch-config/`、`summary/`、`config-log/` 相当の機能別単位。 | 署名、payload、query、snapshot 入力、破損行、期待 paging、期待 rollback。 | 署名検証省略、破損行の黙殺仕様未確認。 |
| [§27.21〜§27.38](fixture.md#sec-27-f) | `runner-extensions/` 配下の機能別単位。 | target、pipeline、cache、queue、hook、remote、approval、notification、trend の正常/異常/部分失敗。 | 実行完了順依存、外部 shell 展開、未定義状態ファイル。 |
| [§27.42〜§27.47](security.md#sec-27-42) | `security/` 配下の scope、token、audit、session、totp、rate-limit 単位。 | route 判定、body 未評価、token hash、audit failure、window reset、secret mask。 | token 本体保存、Authorization header 保存、監査なし権限拒否。 |
| [`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) | `builder-extensions/` 配下の機能別単位。 | Markdown 入力、CLI option、期待 HTML / CSS / JS / REPORT、strict / non-strict の終了コード。 | 外部 library、CDN、実 network、環境依存 timestamp、画像 snapshot だけの合否判定。 |

fixture 名は `success-*`、`failure-*`、`partial-*`、`noop-*`、`security-*` のいずれかで始める。fixture 名に実行時刻、乱数、環境依存 path、実 token 値を含めてはならない。期待時刻は固定値を使い、現在時刻依存の検証では `manifest.json.fake_clock` を fixture 入力に含める。

<a id="sec-27-f-2"></a>
**[fixture 証跡責務 §27-F カタログ固定契約](fixture.md#sec-27-f-2)：**

追加仕様化機能の実装変更は、対象機能について [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の fixture を作成する。複数の owner 機能契約を同一変更で実装する場合は、各対象行の fixture を省略せず作成する。fixture は入力、期待 response、期待 stdout/stderr、期待状態差分、期待 log/history/audit/notify、secret 非表示確認のうち該当するものを含める。

| 節 | fixture 名 | 入力 fixture | 期待出力 / 期待副作用 |
|----|------------|--------------|------------------------|
| [§27.1](commitstatus.md#sec-27-1) | `success-commit-status-pending-success`、`failure-commit-status-unavailable-sha`、`failure-commit-status-api-error`、`noop-commit-status-disabled`、`failure-commit-status-pending-error-final-success`、`failure-commit-status-invalid-payload` | server config、commit SHA、fake GitHub response、build result、fake clock、GitHub token 有無。 | GitHub request order、payload、build log `commit_status` schema keys、history state、fixed error reason、token / Authorization mask、disabled 時呼び出し 0。 |
| [§27.2](runner.md#sec-27-2) | `success-dry-run-changed`、`noop-dry-run-unchanged`、`failure-dry-run-github-error`、`security-dry-run-secret-mask` | CLI args、state dir、fake GitHub response。 | stdout JSON、終了コード、状態差分なし、lock なし、secret mask。 |
| [§27.3](runner.md#sec-27-3) | `success-retry-after-rate-limit`、`success-retry-after-timeout`、`failure-retry-limit-exceeded`、`noop-retry-nonretryable` | retry config、attempt sequence、pipeline / deploy fake result。 | attempts 配列、retry_count、最終 status、SHA 更新有無、未実行 attempt 不作成。 |
| [§27.4](builder.md#sec-27-4) | `success-output-meta-html-report-api`、`success-output-meta-empty-values`、`failure-output-meta-invalid-sha`、`failure-output-meta-invalid-time` | builder args、Markdown 入力、build meta 値。 | HTML meta、REPORT、output-meta response、終了コード、既存出力保護。 |
| [§27.5](api.md#sec-27-5) | `success-config-validate-valid`、`success-config-validate-invalid`、`failure-config-validate-unknown-key`、`security-config-validate-secret-mask` | config JSON、既存 `.server_config`。 | valid/errors/warnings、状態差分なし、unknown key `422`、secret 非表示。 |
| [§27.6](api.md#sec-27-6) | `success-api-access-log-authenticated`、`success-api-access-log-unauthorized`、`failure-api-access-log-append`、`security-api-access-log-body-mask` | HTTP request、actor、target endpoint result。 | `.api_access_log` 追記、status/duration、body 未保存、append failure 挙動。 |
| [§27.7](archive.md#sec-27-7) | `success-log-archive`、`noop-log-archive-empty`、`failure-log-archive-gzip`、`success-log-cleanup`、`partial-log-cleanup-delete-failure` | build logs、retention config、archive target。 | gzip archive、元 log 維持/削除条件、cleanup response、file 別 counter、後続処理、破損 log skip。 |
| [§27.8](runner.md#sec-27-8) | `success-build-status-running`、`success-build-status-final`、`failure-build-status-write`、`failure-build-status-corrupt-api` | runner state、build result、queue/pending/circuit state。 | `.build_status.json`、status API、finalizer、write failure 時 log/history 維持。 |
| [§27.9](runner.md#sec-27-9) | `success-trigger-manual`、`success-trigger-webhook`、`success-trigger-approval`、`failure-trigger-filter-invalid` | trigger source、history query、build log。 | trigger enum 保存、history filter、未知 trigger 不保存、`422`。 |
| [§27.10](runner.md#sec-27-10) | `success-startup-integrity-clean`、`success-startup-integrity-recovered`、`failure-startup-integrity-unrecoverable`、`security-startup-integrity-secret-mode` | startup state files、mode、corrupt files。 | recovery record、WARN/ERROR、build 開始可否、secret file mode 補正。 |
| [§27.11](api.md#sec-27-11) | `success-schedule-interval`、`success-schedule-pause-resume`、`failure-schedule-systemd-update`、`noop-schedule-same-value` | schedule request、server config、fake systemd。 | `.server_config`、systemd result、config log、保存済み config の再取得値。 |
| [§27.12](api.md#sec-27-12) | `success-webhook-push-queued`、`failure-webhook-invalid-signature`、`noop-webhook-duplicate-delivery`、`failure-webhook-queue-full` | GitHub headers、raw body、secret、queue state。 | webhook event、queue entry、状態差分なし条件、重複防止。 |
| [§27.13](api.md#sec-27-13) | `success-webhook-events-page`、`success-webhook-events-empty`、`partial-webhook-events-corrupt-line`、`failure-webhook-events-invalid-query` | `.webhook_events.json`、query。 | events/total、破損行除外、secret 非表示、`422`。 |
| [§27.14](runner.md#sec-27-14) | `success-stats-summary`、`success-stats-timeline`、`partial-stats-corrupt-log-skip`、`failure-stats-invalid-query` | history、build logs、archive logs、query。 | stats response、rounding、WARN、read-only 差分なし。 |
| [§27.15](archive.md#sec-27-15) | snapshot save / list / download / delete / rollback / pending retry の [§27.15 fixture 個別契約](#sec-27-f-4) | `.snapshots/`、output site、history id、running state、PendingTransfer object。 | 冪等 save、manifest 再計算、既存不一致の上書き禁止、保存済み tar.gz download、delete の log 順、rollback prepare / audit / worker / finalizer、snapshot 再展開による pending retry、不正 entry 防止。 |
| [§27.16](api.md#sec-27-16) | `success-health-ok`、`success-health-degraded`、`failure-health-read-error`、`noop-health-readonly` | status/pending/notify/process state。 | health JSON、HTTP status、degraded 判定、状態差分なし。 |
| [§27.17](api.md#sec-27-17) | `success-log-search-level`、`success-log-search-archive`、`partial-log-search-corrupt-skip`、`failure-log-search-invalid-level` | build logs、archive logs、query。 | build 単位の `SearchResult`、一致 `lines` の source / 行順、破損除外、`422`。 |
| [§27.18](api.md#sec-27-18) | `success-branch-config-get-default`、`success-branch-config-post`、`failure-branch-config-invalid-path`、`partial-branch-config-log-failure` | branch config request、existing config。 | branch_targets、config log、validation 差分なし、保存済み状態維持。 |
| [§27.19](runner.md#sec-27-19) | `success-weekly-summary-auto`、`success-weekly-summary-manual`、`noop-weekly-summary-same-day`、`failure-weekly-summary-send` | notify config、history、fake webhook。 | payload、notify log、sent date 更新条件、失敗時 pending。 |
| [§27.20](api.md#sec-27-20) | `success-config-diff-simple`、`success-config-diff-nested`、`noop-config-diff-same-value`、`security-config-diff-secret-mask` | before/after config、request。 | diff/diff_text、config log、no-op 差分なし、secret mask。 |
| [§27.21](runner.md#sec-27-21) | `success-multi-file-one-change`、`success-multi-file-many-change`、`noop-multi-file-all-skip`、`failure-multi-file-path-traversal` | branch target、target files、SHA cache。 | build 対象集合、対象別 SHA 更新、重複排除、path error。 |
| [§27.22](runner.md#sec-27-22) | `success-builder-command-extra-args`、`success-builder-command-env-merge`、`success-builder-command-cache`、`failure-builder-command-forbidden-key`、`failure-builder-command-reserved-env`、`security-builder-command-secret-mask` | pipeline config、branch env、builder fake result。 | 標準 builder argv、env merge、stdout/stderr、終了コード、deploy / SHA 更新禁止条件。 |
| [§27.23](runner.md#sec-27-23) | `success-local-watch-change`、`noop-local-watch-no-change`、`failure-local-watch-state-corrupt`、`failure-local-watch-tag-filter-conflict` | local files、watch state、server config。 | local watch state、trigger、GitHub call 0、終了コード `2`。 |
| [§27.24](runner.md#sec-27-24) | `success-tag-filter-match`、`noop-tag-filter-unmatched`、`failure-tag-filter-api`、`failure-tag-filter-pattern` | tag refs、patterns、SHA cache。 | matched tags、SHA 更新条件、skip 挙動、不正 pattern。 |
| [§27.25](builder.md#sec-27-25) | `success-build-cache-hit`、`success-build-cache-miss`、`partial-build-cache-byte-mismatch`、`failure-build-cache-save`、`failure-build-cache-cli-path` | Markdown、cache index/pages、theme/version、cache CLI path 条件。 | byte 一致、cache atomic save、hit 破棄、CLI 早期停止、build 成否維持。 |
| [§27.26](runner.md#sec-27-26) | `success-parallel-targets-all`、`partial-parallel-targets-some-fail`、`failure-parallel-targets-all-fail`、`success-parallel-targets-order-stable` | target list、parallel result sequence。 | 設定順 result、overall status、`deploy_timeout` / `deploy_ssh_error` / `deploy_checksum_error` / `deploy_internal_error`、history/status。 |
| [§27.27](runner.md#sec-27-27) | `success-hook-pre-post`、`failure-hook-pre-abort`、`partial-hook-post-fail`、`security-hook-shell-denied` | hook config、command_args、fake command result。 | hook log、secret mask、pre abort、shell 展開禁止。 |
| [§27.28](builder.md#sec-27-28) | `success-dependency-manifest`、`success-dependency-missing`、`failure-dependency-build-keeps-old`、`security-dependency-path-normalize` | Markdown refs、existing manifest、build result。 | dependency manifest、成功時置換、失敗時旧 manifest 維持。 |
| [§27.29](runner.md#sec-27-29) | `success-remote-build-artifact`、`failure-remote-build-auth`、`failure-remote-build-checksum`、`security-remote-build-argument-quoting`、`security-remote-build-unsafe-archive` | remote config、artifact archive、manifest/checksum、引用対象 argument。 | remote argument 引用、artifact atomic fetch、一時展開、検証後 deploy、既存出力保護、unsafe entry 拒否。 |
| [§27.30](runner.md#sec-27-30) | `success-approval-approve`、`noop-approval-reject`、`noop-approval-expire`、`failure-approval-double-approve`、`partial-approval-approved-append`、`partial-approval-pending-audit` | approval queue、API action、clock、state/audit fake failure。 | 状態 enum、requested_force、queue 連携、物理削除なし、audit、部分失敗後の冪等再開。 |
| [§27.31](runner.md#sec-27-31) | `success-branch-env-inject`、`security-branch-env-secret-mask`、`failure-branch-env-invalid-key`、`failure-branch-env-reserved-key`、`failure-branch-env-mask-failure` | branch config env、build command fake、reserved key、mask failure。 | child env、log key 名、secret mask、不正 key / value、reserved key 拒否、失敗後副作用禁止。 |
| [§27.32](runner.md#sec-27-32) | `success-notify-multi-channel`、`partial-notify-webhook-pending`、`failure-notify-webhook-nonretryable-http`、`failure-notify-webhook-network`、`failure-notify-email`、`failure-notify-command`、`partial-notify-retry-exhausted`、`failure-notify-log-write`、`noop-notify-disabled-event`、`security-notify-secret-mask` | notify config、build result、fake channel。 | notify log/pending、error code、channel 順、build 成否維持、payload mask、no-op。 |
| [§27.33](runner.md#sec-27-33) | `success-trend-summary-update`、`success-trend-replace-build-id`、`failure-trend-corrupt-rebuild`、`failure-trend-api-invalid-n` | build history/log、trend state、query。 | sample upsert、summary 再計算、retention、破損復旧、`422`。 |
| [§27.34](runner.md#sec-27-34) | `success-chain-dag-order`、`noop-chain-disabled-job`、`failure-chain-cycle`、`partial-chain-required-skip` | chain config、job result sequence。 | chain_run_id、DAG 順、disabled 除外、skipped 条件、validation `422`。 |
| [§27.35](runner.md#sec-27-35) | `success-priority-urgent-first`、`success-priority-active-first`、`failure-priority-created-seq-missing`、`failure-priority-invalid`、`failure-priority-queue-full` | active / waiting entries、runner lock、API action。 | active 優先、priority 順、created_seq 欠落拒否、waiting-to-active atomic move、`422` / `429`。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-timeout`、`success-failure-category-deploy`、`failure-failure-category-filter-invalid`、`security-failure-evidence-mask` | exit code、stderr、API error、status。 | category enum、filter `422`、純粋な success の `null`、deploy pending の `deploy_failure`、evidence 上限、secret mask。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-record`、`success-environment-builder-version-record`、`failure-environment-builder-version-contract`、`failure-environment-write`、`security-environment-secret-excluded` | OS/env/version/state path inputs。 | environment object、標準 builder 第 2 token、不正 version 時の builder 未起動、basename 保存、secret 非保存、write failure。 |
| [§27.38](runner.md#sec-27-38) | `success-duration-anomaly-avg`、`noop-duration-anomaly-insufficient-samples`、`partial-duration-anomaly-notify-failure`、`failure-duration-anomaly-invalid-config` | trend summary、duration config、build result。 | anomaly tag、history flag、通知 event、sample 不足、設定不正。 |
| [§27.42](security.md#sec-27-42) | `security-scope-trigger-allowed`、`security-scope-read-denied`、`security-scope-path-param`、`failure-scope-audit-failure` | route/method、token record、request body。 | body 未評価、`403`、audit、Authorization 非保存。 |
| [§27.43](security.md#sec-27-43) | `security-token-create-once`、`security-token-list-mask`、`success-token-revoke`、`failure-token-expired-auth`、`failure-token-random-source`、`failure-token-hash-collision` | token request、token state、clock、random source。 | token 形式、1 回表示、hash 保存、revoke、期限切れ `401`、random source / collision 失敗時の無副作用。 |
| [§27.44](security.md#sec-27-44) | `success-audit-operation`、`success-audit-denied`、`failure-audit-append`、`security-audit-secret-mask` | actor、target、result、operation input。 | JSON Lines、mask、必須 audit failure、paging。 |
| [§27.45](security.md#sec-27-45) | `success-session-active`、`failure-session-expired`、`success-session-timeout-update`、`success-session-revoke-all` | session state、clock、timeout config。 | 非 sliding `expires_at`、`last_used_at`、`401`、current 以外の revoke、token 非表示。 |
| [§27.46](security.md#sec-27-46) | `security-totp-setup-once`、`success-totp-confirm`、`failure-totp-code-reuse`、`success-totp-disable` | setup ticket、TOTP code、clock、secret state。 | secret 有効保存条件、ticket 一回使用、window、disable。 |
| [§27.47](security.md#sec-27-47) | `security-rate-limit-login`、`success-rate-limit-window-reset`、`security-rate-limit-ip-actor`、`failure-rate-limit-state-save` | policy、rate state、RemoteAddr、actor。 | count、`429`、audit 成功条件、部分 count 更新なし。 |

<a id="sec-27-f-3"></a>
**[fixture 証跡責務 §27-F security 詳細本文責務 §27.42〜§27.47 fixture 固定契約](fixture.md#sec-27-f-3)：**

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の fixture は、[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の認証、scope、token、audit、session、TOTP、rate limit の処理順、状態保存順、漏えい禁止、副作用境界を固定する。各 fixture は `manifest.json.owner_component` を `security`、`manifest.json.section` を対象 [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、`manifest.json.feature` を [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の feature 名に固定する。`components` には、HTTP request / response を検証する場合は `api`、SDK error 変換を検証する場合は `sdk`、UI 表示 / field 消去を検証する場合は `ui`、状態ファイルを検証する場合は `statefile` を含める。

| 節 | feature | fixture 名 | 固定する確認 |
|----|---------|------------|--------------|
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-trigger-allowed` | `trigger` scope token、`POST /api/build` / `POST /api/build/force` の許可、body parse 前 scope 判定、audit actor、build trigger actor、Authorization 非保存を固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-read-denied` | `trigger` scope token で read endpoint を呼んだ場合の `403`、endpoint 固有処理なし、body validation なし、permission_denied audit、token 本体非保存を固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-path-param` | `DELETE /api/tokens/{id}` の path parameter を正規化 path pattern で判定し、query string を scope 判定に使わないことを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `failure-scope-audit-failure` | scope 不足時の audit 追記失敗で `500`、対象 endpoint 未実行、request body 未評価、`.api_tokens` / endpoint 状態差分なしを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-multi-scope` | 複数 scope token はいずれか 1 scope が endpoint group に一致した場合だけ許可し、未知 scope / 空 scopes は token record 破損 `500` とすることを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `noop-scope-health-webhook-exempt` | `GET /api/health` と `POST /api/webhook` は API token scope 判定対象外とし、未定義 route は認証 / rate limit / body parse 前に `404` / `405` を返すことを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `security-rate-webhook-ip-only` | `POST /api/webhook` は body size と署名検証成功後、JSON parse 前に `trigger` group の IP key だけを判定・更新し、上限超過は `actor_type:"webhook"`、`actor_id:"webhook"` の audit 成功後の `429`、event log / queue / runner 起動要求の差分なしとすることを固定する。 |
| 認証共通 | `auth_common` | `security-auth-one-time-response` | session token、login ticket、API token 本体、TOTP setup secret、otpauth URI が許可された成功 response 1 回だけに出現し、以後の response / log / expected に残らないことを固定する。 |
| 認証共通 | `auth_common` | `failure-auth-log-before-token` | token / ticket / secret 返却前の access / audit log fake failure で `500`、one-time 値を response せず、平文保存と session record 追加は 0 件とする。login でログ失敗前に保存済みの `login_count` / `last_login_at`、TOTP で保存済みの `last_accepted_step` は戻さず、後続の未実行 write は 0 件とする。 |
| 認証共通 | `auth_common` | `security-auth-memory-only` | process restart 相当で session、ticket、TOTP setup 仮 secret、login 失敗回数が消え、未定義永続 state file が作成されないことを固定する。 |
| 認証共通 | `auth_common` | `security-init-credentials-atomic` | `--init-credentials` の非 terminal stdin 成功、terminal 拒否、末尾 LF 欠落、複数行、UTF-8 不正、8 文字未満、128 文字超過、NUL / CR、stdin read failure、予備確認時の既存あり、予備確認後かつ create-only lock 取得後の既存競合、相対 path、write failure、rand failure、rename 後 sync failure、create-only lock cleanup failureについて、stdout / stderr / exit code / stdin read 回数 / file mode / target 非上書き / partial file / password buffer 消去を固定する。予備確認時の既存ありでは stdin read を 0 件とする。 |
| 認証共通 | `auth_common` | `security-auth-forbidden-plaintexts` | password、session token、API token、ticket、hash 算出入力、salt、TOTP code、TOTP secret が response / stdout / stderr / journal / state / expected に平文で出ないことを固定する。 |
| 認証共通 | `auth_common` | `security-auth-credentials-schema` | `.admin_credentials` の必須 key、未知 key 拒否、64 文字 salt、`login_count` の 0〜9223372036854775807、上限での飽和加算、範囲外の破損扱いを固定する。 |
| 認証共通 | `auth_common` | `security-auth-must-change-boundary` | `must_change:false` の `none`、`true` かつ session 発行後 `login_count=1..4` の `prompt`、`5`、`6`、上限値の `forced`、各 session record の `password_change_required` 一致を固定する。 |
| 認証共通 | `auth_common` | `security-auth-totp-deferred-mode` | TOTP 有効の password 成功は exact `{totp_required:true,ticket}`、credentials 無変更、ticket に mode 非保持とする。ticket 発行後の別 session login による `login_count` / `last_login_at` 変化は fingerprint を変更せず、TOTP 成功時に最新 credentials から `login_count` 増加と最終 `must_change` を算出する。ticket 発行後の password 変更は fingerprint 不一致とし、ticket 消費、`401`、TOTP / credentials / session 無変更を固定する。 |
| 認証共通 | `auth_common` | `security-auth-forced-session-gate` | `forced` session の `POST /api/change-password` / `POST /api/logout` 許可、その他の認証必須 endpoint の body parse 前 `403 Password change required`、API token の gate 対象外、password 変更成功後の現 session gate 解除を固定する。 |
| 認証共通 | `auth_common` | `security-auth-login-lock-precedence` | IP key 別 10 回目 `401`、後続 lock の `429 Too many attempts`、pre-auth rate limit の `429 Too many requests` 優先、非 password 失敗の count 不変、成功 / 期限による entry 削除、restart 破棄を fake clock で固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `security-token-create-once` | 32 bytes random input、`base64.RawURLEncoding`、`act_` prefix、47 bytes 完成長、`[A-Za-z0-9_-]` 本体 alphabet、hash 保存、token response 一回表示、`.api_tokens` 保存 → `.access_log` の `token_create` → `.audit_log` の `token_create` → response の順序を固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `security-token-list-mask` | `GET /api/tokens` が token 本体、token hash、Authorization header を返さず、`created_at` 降順 / id 昇順で返すことを固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `success-token-revoke` | `DELETE /api/tokens/{id}` の id 検証、`revoked_at` 保存、自己失効、`.access_log` / `.audit_log` の `token_revoke` 追記、以後 `401` を固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `failure-token-expired-auth` | 期限切れ token の `401`、`last_used_at` 未更新、access / audit `token_expired`、endpoint 固有処理なしを固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `failure-token-record-corrupt` | `.api_tokens` の未知 key、必須 key 不足、hash 形式不正、未知 scope、空 scopes で token 認証 / 一覧 / 作成 / 失効を `500` にし、自動再生成しないことを固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `partial-token-create-audit-failure` | token record 保存後の access log または audit log 追記失敗で `500`、作成済み token record 維持、token 本体を response に含めないことを固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `failure-token-random-source` | random source がエラーまたは 32 bytes 未満を返す場合の `500`、`.api_tokens` / access log / audit log 無差分、lock 解放、token / hash 非出力を固定する。 |
| [§27.43](security.md#sec-27-43) | `api_key_management` | `failure-token-hash-collision` | 生成 token hash が既存 `token_hash` と一致する場合の `500`、再生成 0 回、既存 record / log 無差分、lock 解放、token / hash 非出力を固定する。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `success-audit-operation` | auth、session、TOTP、token、config、build trigger、approval の全対象 action について actor / target / result / request_id / timestamp を固定し、runner event の `request_id:null` を含める。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `success-audit-denied` | permission denied と rate limit denied の `actor_type`、`actor_id`、`target_type:"endpoint"`、`target_id:{METHOD path}`、`result:"denied"` を固定し、署名検証済み Webhook は `actor_type:"webhook"`、`actor_id:"webhook"` とする。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `failure-audit-append` | 必須 audit 追記失敗で対象操作を `500` とし、保存済み状態の巻き戻し有無を操作種別別保存順どおり固定する。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `security-audit-secret-mask` | request body、query 全体、header、cookie、secret、password、token、hash、salt、TOTP secret が `.audit_log`、response、server log に残らないことを固定する。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `success-audit-pagination-filter` | `GET /api/audit-log` の `limit`、`offset`、`actor`、`action`、`result`、timestamp 降順、壊れた行除外、取得操作自体を audit しないことを固定する。 |
| [§27.44](security.md#sec-27-44) | `audit_log` | `failure-audit-invalid-filter` | 未知 action、未知 result、200 Unicode scalar values 超過 actor を `422`、状態差分なし、壊れた行の有無と独立判定に固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `success-session-active` | login 成功時の `issued_at + session_timeout_seconds`、UTC ISO 8601 秒精度、非 sliding の `expires_at`、認証 request ごとの `last_used_at` 更新、session token 非表示を固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `failure-session-expired` | timeout session の `401`、endpoint 固有処理なし、session token 非保存、access / audit の結果を固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `success-session-timeout-update` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の `.server_config.session_timeout_seconds` 下限・上限、config log diff、audit `config_update`、新規 session だけへの適用を固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `success-session-revoke-all` | `POST /api/sessions/revoke-all` が current session 以外を失効し、実行中 request の response、`.access_log` / `.audit_log` の `session_revoke_all`、失効済み session の以後の `401`、current session の継続利用を固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `noop-session-timeout-same-value` | 同値更新は `200`、`.server_config` / `.config_log` / `.audit_log` / 既存 session 差分なしを固定する。 |
| [§27.45](security.md#sec-27-45) | `session_timeout` | `failure-session-timeout-invalid` | 範囲外、型不一致、不正 JSON を `422` / `400`、`.server_config` / `.sessions` / `.config_log` / `.audit_log` 差分なしに固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `security-totp-setup-once` | setup 仮 secret はメモリだけに保持し、response と UI 一回表示以外へ secret / otpauth URI を残さず、TOTP 有効時 setup `409` を固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `success-totp-confirm` | 仮 secret と code 検証、`.totp_secret` 保存、仮 secret 削除、audit 追記、status response の secret 非表示を固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `failure-totp-code-reuse` | `last_accepted_step` 以下の code replay を `401`、ticket 削除、secret / ticket 非保存、状態差分境界を固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `success-totp-disable` | code 検証後に `.totp_secret` を無効値保存、未使用 ticket / 仮 secret 削除、既存 session 維持、audit を固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `failure-totp-ticket-invalid` | ticket 不在、期限切れ、hash 不一致、削除済み ticket をすべて `401 {"error":"Unauthorized"}`、詳細非表示、ticket 巻き戻しなしに固定する。 |
| [§27.46](security.md#sec-27-46) | `totp` | `partial-totp-audit-failure` | confirm / disable の audit 失敗時 `500`、保存済み `.totp_secret` や仮 secret 削除は巻き戻さず、secret 平文を返さないことを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `security-rate-limit-login` | login group の認証前 IP key 判定、11 回目 `429`、count 非増加、permission_denied audit、request body 非保存を固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `success-rate-limit-window-reset` | `now >= window_start + window_seconds` で window reset、count 初期化、`reset_at`、state_summary 並び順を固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `security-rate-limit-ip-actor` | session / API token の actor key と IP key を同一 lock 内で判定 / 更新し、片方だけの count 更新を残さないことを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `failure-rate-limit-state-save` | `.api_rate_state` 保存失敗で endpoint 固有処理なし、部分 count 更新なし、`500`、audit 追記なしを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `noop-rate-limit-disabled` | `.server_config.api_rate_limit.enabled=false` では `.api_rate_state` を読まず、count / audit 差分なしで対象 endpoint へ進むことを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `partial-rate-limit-policy-update` | policy 保存 → windows 空保存 → config log → audit の順序、同値 no-op、audit 失敗時の保存済み状態維持を固定する。 |

[§27.42〜§27.47](security.md#sec-27-42) の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持つ。security fixture では、`.admin_credentials`、`.sessions`、`.api_tokens`、`.audit_log`、`.access_log`、`.api_access_log`、`.server_config`、`.totp_secret`、`.api_rate_state`、対象 endpoint 状態ファイルの forbidden side effect を必ず列挙する。request body、Authorization header、session token、API token、token hash、password hash、salt、TOTP secret、ticket、otpauth URI の forbidden leak は、[fixture 証跡責務 §27-F expected/security.json schema 固定契約](#sec-27-f-11-security) に従って列挙する。認証共通 fixture では、process restart 後に残ってはならない memory-only 値、未定義永続 state file、token 返却前 log failure 時の forbidden response field を必ず列挙する。

[`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の `expected/security.json` は、[fixture 証跡責務 §27-F expected/security.json schema 固定契約](#sec-27-f-11-security) の全 root key を持つ。token 本体、TOTP secret、ticket、otpauth URI、Authorization header、session token、API token は `allowed_one_time_response_fields` に明示された fixture の該当 response 以外では出現禁止とする。hash 値を検証する場合も、hash 算出入力の平文を expected file へ保存してはならない。`memory_only_values` には session、login ticket、TOTP setup 仮 secret、login 失敗回数を列挙し、restart 後に消えていることを expected に固定する。

<a id="sec-27-1"></a>
**[fixture 証跡責務 §27.1〜§27.11 feature fixture 固定契約](fixture.md#sec-27-1)：**

[`docs/details/fixture.md` fixture 証跡責務 §27.1〜§27.11](fixture.md#sec-27-1) の fixture は、Commit Status、dry-run、retry、output meta、config validation、access log、archive、build status、trigger、startup integrity、schedule の共通処理挙動を固定する。各 fixture は、owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務に定義された入力、状態、出力、外部呼び出し、副作用、secret mask を expected に固定し、実装検証証跡に対象 fixture と実行結果を列挙する。

| 節 | fixture | 固定する内容 |
|----|---------|--------------|
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `success-commit-status-pending-success` | fake GitHub server への `pending` → `success` 送信順、payload `state` / `context` / `description` / `target_url`、HTTP `201` success、`.build_logs.commit_status` の [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) schema keys、`.build_history.commit_status_state="success"` を固定する。 |
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `failure-commit-status-unavailable-sha` | commit SHA が取得できない場合に GitHub Status API 呼び出し 0 件、build 継続、`.build_logs.commit_status.state=null`、`error="commit sha unavailable"`、history の build status 非反転と `commit_status_state=null` を固定する。 |
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `failure-commit-status-api-error` | final の fake GitHub failure を table-driven subcase とし、HTTP `401` / `403` / `404` / `429` / `5xx` / その他の非 `201`、network / DNS / timeout の各固定 error reason、WARN、送信しようとした state、HTTP status、build 成否非反転、secret mask、Authorization header 非保存を固定する。`403` subcase は `Commit statuses: Write` 不足を入力条件に含め、`error="github status unauthorized"` を固定する。 |
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `noop-commit-status-disabled` | `.server_config.commit_status_enabled=false` で GitHub Status API 呼び出し 0 件、`commit_status.enabled=false`、`state=null`、`error=null`、payload / token / status side effect なしを固定する。 |
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `failure-commit-status-pending-error-final-success` | pending の fake GitHub failure 後も build を継続し、final success を 1 回送信し、build log は final summary だけを保存し、pending failure は WARN と `expected/effects.json.external_calls` で固定する。 |
| [`docs/details/commitstatus.md` 詳細本文責務 §27.1](commitstatus.md#sec-27-1) | `failure-commit-status-invalid-payload` | invalid owner / repo / context / target_url / token unavailable のいずれかで GitHub 呼び出し 0 件、`state="error"`、固定 error reason、secret 非表示、build result 非反転を固定する。 |
| [§27.2](runner.md#sec-27-2) | `success-dry-run-changed` | SHA 差分あり target の stdout JSON、`would_build=true`、`reason="sha_changed"`、`would_call`、`would_write`、終了コード `0`、状態差分なしを固定する。 |
| [§27.2](runner.md#sec-27-2) | `noop-dry-run-unchanged` | SHA 差分なし、cooldown、circuit open の `would_build=false` と reason、stdout JSON 1 件、lock / log / history / status / notification / deploy 差分なしを固定する。 |
| [§27.2](runner.md#sec-27-2) | `failure-dry-run-github-error` | fake GitHub read 最終失敗で終了コード `3`、`reason="github_error"`、`errors[]`、pipeline / deploy / commit status 呼び出し 0 件、状態差分なしを固定する。 |
| [§27.2](runner.md#sec-27-2) | `security-dry-run-secret-mask` | state dir、env、GitHub response、error message に secret 風値があっても stdout、stderr、effects、expected に平文を残さず、`secrets_masked=true` を固定する。 |
| [§27.3](runner.md#sec-27-3) | `success-retry-after-rate-limit` | GitHub API 429 後の retry、attempts 2 件、backoff fake clock、最終 success、`retry_count=1`、SHA 更新は最終成功後だけを固定する。 |
| [§27.3](runner.md#sec-27-3) | `success-retry-after-timeout` | pipeline timeout または deploy network timeout 後の retry success、attempt schema、未成功 attempt による deploy / snapshot / SHA 副作用なしを固定する。 |
| [§27.3](runner.md#sec-27-3) | `failure-retry-limit-exceeded` | `1 + build_retry_max` 件の attempts、最終詳細結果、未実行 attempt 不作成、history の同じ詳細結果、secret mask、終了コードを固定する。 |
| [§27.3](runner.md#sec-27-3) | `noop-retry-nonretryable` | pipeline exit code 非 0、checksum mismatch、validation failure の各 nonretryable 失敗で retry 0 件、pending transfer または failure 保存、SHA 更新なしを固定する。 |
| [§27.4](builder.md#sec-27-4) | `success-output-meta-html-report-api` | builder HTML meta、`[REPORT]`、build log、`GET /api/output-meta` response が同じ build id / sha / timestamp / title を返すことを固定する。 |
| [§27.4](builder.md#sec-27-4) | `success-output-meta-empty-values` | optional meta が空または null の場合の HTML 出力省略、REPORT 空値表現、API response の null / empty string 区別、secret 非表示を固定する。 |
| [§27.4](builder.md#sec-27-4) | `failure-output-meta-invalid-sha` | SHA 形式不正で終了コード `2` または API `422`、公開出力維持、REPORT なし、既存 output meta 非破壊を固定する。 |
| [§27.4](builder.md#sec-27-4) | `failure-output-meta-invalid-time` | timestamp parse 不能、範囲外、非 UTC 値で fixed error、公開出力維持、build log / API response に不正時刻を保存しないことを固定する。 |
| [§27.5](api.md#sec-27-5) | `success-config-validate-valid` | `POST /api/config/validate` が endpoint 固有の業務状態を変更せず、正規化後 config、`valid=true`、warnings/errors 空、共通 security / observability 副作用を固定する。 |
| [§27.5](api.md#sec-27-5) | `success-config-validate-invalid` | 型不一致、範囲外、相互排他違反で `valid=false`、`errors[]`、HTTP status、状態差分なし、secret mask を固定する。 |
| [§27.5](api.md#sec-27-5) | `failure-config-validate-unknown-key` | unknown root key / nested key を `422`、状態差分なし、`.config_log` 追記なし、response の key path 固定で返すことを固定する。 |
| [§27.5](api.md#sec-27-5) | `security-config-validate-secret-mask` | PAT、SMTP password、webhook secret、API token 風値を request / response / logs / effects に平文で残さず、key 名だけを返すことを固定する。 |
| [§27.6](api.md#sec-27-6) | `success-api-access-log-authenticated` | 認証済み request の `.api_access_log` JSON Lines 追記、request id、method、path、検証済み query の値 / mask、status、duration_ms、auth_type、actor、remote_addr、user_agent、body 非保存を固定する。 |
| [§27.6](api.md#sec-27-6) | `success-api-access-log-unauthorized` | 未認証 / 権限不足 request の log、`auth_type="none"`、`actor=null`、status `401` / `403`、query 検証前は `query={}`、Authorization header 非保存を固定する。 |
| [§27.6](api.md#sec-27-6) | `failure-api-access-log-append` | `.api_access_log` append failure で本来の HTTP status / body / header を維持し、対象 endpoint の状態と先行 audit / config / access / notify log を巻き戻さず、`API_ACCESS_LOG_WRITE_FAILED` WARN を 1 件だけ記録し、部分 JSON 行を残さないことを固定する。 |
| [§27.6](api.md#sec-27-6) | `security-api-access-log-body-mask` | request body、raw query、未知 / 検証失敗 query、free-form string query の実値、Authorization header、session token、API token が access log、response、effects に保存されないこと、許可済み free-form string は `"***"` で記録されることを固定する。 |
| [§27.7](archive.md#sec-27-7) | `success-log-archive` | 対象 `.build_logs/{id}.json` の gzip 作成、元 log 削除、archive からの API 参照、disk usage 反映、gzip path を固定する。 |
| [§27.7](archive.md#sec-27-7) | `noop-log-archive-empty` | archive 対象なし、実行中 build log、既存 archive の skip、件数 0、状態差分なしを固定する。 |
| [§27.7](archive.md#sec-27-7) | `failure-log-archive-gzip` | gzip write / JSON read failure で WARN、処理継続または fixed failure、元 log 維持、partial `.gz` 非公開を固定する。 |
| [§27.7](archive.md#sec-27-7) | `success-log-cleanup` | retention 対象の通常 log basename ASCII 昇順 → archive log basename ASCII 昇順の各 file `os.Remove` 1 回、各削除成功 / `os.IsNotExist` の `deleted_count+1`、`failed_count=0`、処理後再読込で空の archive directory だけ `os.Remove` 1 回を固定する。 |
| [§27.7](archive.md#sec-27-7) | `partial-log-cleanup-delete-failure` | 通常 log と archive log の削除失敗を個別に fake し、当該 file ごとの `failed_count+1`、`LOG_CLEANUP_DELETE_FAILED`、同一 file の再試行 0 回、後続 file 処理継続を固定する。archive directory 再読込 / 削除失敗は `LOG_ARCHIVE_DIRECTORY_CLEANUP_FAILED`、counter 不変、再試行 0 回とする。 |
| [§27.8](runner.md#sec-27-8) | `success-build-status-running` | build 開始前の `.build_status.json` atomic write、`status="running"`、`running=true`、`current_build_id`、pending 件数を固定する。 |
| [§27.8](runner.md#sec-27-8) | `success-build-status-final` | success / failure / skipped / deploy pending の finalizer、`running=false`、`current_build_id=null`、last fields 維持、API read 値を固定する。 |
| [§27.8](runner.md#sec-27-8) | `failure-build-status-write` | status write failure で runner 終了コード最低 `1`、ERROR log、build log / history 維持、部分 `.build_status.json` 非公開を固定する。 |
| [§27.8](runner.md#sec-27-8) | `failure-build-status-corrupt-api` | API が破損 `.build_status.json` を読んだ場合の `/api/status` / dashboard `500`、health degraded、自動修復なしを固定する。 |
| [§27.9](runner.md#sec-27-9) | `success-trigger-manual` | manual queue / force payload の trigger 保存、log / history / status / API filter / SDK / UI 表示の値一致を固定する。 |
| [§27.9](runner.md#sec-27-9) | `success-trigger-webhook` | 署名検証済み webhook queue entry の `trigger="webhook"`、target 限定、重複 delivery なし、history filter を固定する。 |
| [§27.9](runner.md#sec-27-9) | `success-trigger-approval` | approval queue entry 処理時の `trigger="approval"`、承認済み entry のみ処理、pending 以外は処理しないことを固定する。 |
| [§27.9](runner.md#sec-27-9) | `failure-trigger-filter-invalid` | unknown trigger query / queue trigger を `422` または処理中断にし、新規保存禁止、既存 unknown trigger の warning 表示を固定する。 |
| [§27.10](runner.md#sec-27-10) | `success-startup-integrity-clean` | 対象状態ファイルが正常な場合、backup / rewrite / notification / status warning なしで target 処理へ進むことを固定する。 |
| [§27.10](runner.md#sec-27-10) | `success-startup-integrity-recovered` | 破損、必須 key 不足、unknown key 正規化時の backup / 初期化 / atomic write / `.build_status.json.last_trigger` / config_corrupt 通知を固定する。 |
| [§27.10](runner.md#sec-27-10) | `failure-startup-integrity-unrecoverable` | permission / io error で自動復旧なし、終了コード `1` または `2`、`.build_state.running` 未変更、build log / history 非作成を固定する。 |
| [§27.10](runner.md#sec-27-10) | `security-startup-integrity-secret-mode` | secret file mode 補正、secret 平文非表示、backup byte 一致、dry-run で backup / rewrite なしを固定する。 |
| [§27.11](api.md#sec-27-11) | `success-schedule-interval` | schedule interval 更新 request、`.server_config` 保存、fake systemd update、`.config_log`、再取得 response の一致を固定する。 |
| [§27.11](api.md#sec-27-11) | `success-schedule-pause-resume` | pause / resume request、timer enable state、config 保存値、systemd fake 呼び出し順、UI / SDK response を固定する。 |
| [§27.11](api.md#sec-27-11) | `failure-schedule-systemd-update` | `.server_config` 保存後の fake systemd failure、HTTP `500`、config log `result="partial_failure"`、`error="systemd_update_failed"`、record 1 件、未定義 rollback なしを固定する。 |
| [§27.11](api.md#sec-27-11) | `noop-schedule-same-value` | 同一値更新時は `200 {"message":"No changes","interval_seconds":N}`、`.server_config` 差分なし、systemd 呼び出し 0 件、`.config_log` / `.audit_log` 追記 0 件、再取得値 `N` を固定する。 |

[§27.1〜§27.11](fixture.md#sec-27-1) の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持ち、未使用項目も空配列で明示する。dry-run、validation、noop、security fixture では、状態ファイル、lock、history、build log、archive、notification、deploy、commit status の forbidden side effect を必ず列挙する。

<a id="sec-27-f-4"></a>
**[fixture 証跡責務 §27-F 追加仕様化機能 §27.12〜§27.20 fixture 固定契約](fixture.md#sec-27-f-4)：**

[§27.12〜§27.20](fixture.md#sec-27-f-4) の fixture は、Webhook、Webhook event log、duration stats、snapshot artifact、health、log severity search、branch config、weekly summary、config diff の運用 API / runner / archive 連動を固定する。各 fixture は、HTTP response だけでなく、状態ファイル差分、外部呼び出し、保存順、失敗時に発生してはならない副作用、secret mask を expected に固定し、実装検証証跡に対象 fixture と実行結果を列挙する。

| 節 | fixture | 固定する内容 |
|----|---------|--------------|
| [§27.12](api.md#sec-27-12) | `success-webhook-push-queued` | raw body HMAC 検証、push payload parse、対象 branch 判定、`.webhook_events.json` `result="queued"`、`.build_state.queued[]` `trigger="webhook"`、event 保存後の systemctl 固定引数 1 回、response `202` の queue id / dispatch、delivery id 一致を固定する。 |
| [§27.12](api.md#sec-27-12) | `failure-webhook-invalid-signature` | secret 不在、署名 header 欠落、prefix 不正、hex 不正、署名不一致で `401`、event log / queue / build state / access log secret 値差分なしを固定する。 |
| [§27.12](api.md#sec-27-12) | `noop-webhook-duplicate-delivery` | 同一 delivery id、branch、sha の active または waiting entry がある場合に新規 queue 追加なし、event log `duplicate`、既存 queue id response、同じ id の起動要求再実行、idempotency を固定する。active / waiting 両方に同一 entry がある破損 fixture では active id を返し、waiting を自動削除せず固定 warning を確認する。 |
| [§27.12](api.md#sec-27-12) | `partial-webhook-dispatch-fallback` | queue / event 保存成功後の systemctl 非 `0` または timeout、`dispatch:"timer_fallback"`、queue / event 保持、固定 server log、追加 retry なしを固定する。 |
| [§27.12](api.md#sec-27-12) | `failure-webhook-queue-full` | queue 上限時に event log `queue_full` を追記し、queue 差分なし、response `429`、secret / raw payload 非保存を固定する。 |
| [§27.13](api.md#sec-27-13) | `success-webhook-events-page` | `.webhook_events.json` を timestamp 降順、同時刻 file 逆順で並べ、`limit` / `offset` 適用後の events、壊れていない行だけの `total`、secret 非表示を固定する。 |
| [§27.13](api.md#sec-27-13) | `success-webhook-events-empty` | event file 不在または空で `events=[]`、`total=0`、read-only no-write、server log なしを固定する。 |
| [§27.13](api.md#sec-27-13) | `partial-webhook-events-corrupt-line` | 破損 JSON Lines を response から除外し、固定 WARN code だけを server log に出し、破損行内容と secret を出さず、状態を修復しないことを固定する。 |
| [§27.13](api.md#sec-27-13) | `failure-webhook-events-invalid-query` | `limit` / `offset` 範囲外、未知 query、非整数 query で `422`、状態差分なし、server log に query 値の secret 風値を残さないことを固定する。 |
| [§27.14](runner.md#sec-27-14) | `success-stats-summary` | build history と通常 / archive build log から成功数、失敗数、成功率、平均 interval、平均 / 最大 duration を固定丸めで返し、read-only no-write を固定する。 |
| [§27.14](runner.md#sec-27-14) | `success-stats-timeline` | `days` 範囲、UTC 日付 bucket、日付降順、0 件日除外、status 分類、状態差分なしを固定する。 |
| [§27.14](runner.md#sec-27-14) | `partial-stats-corrupt-log-skip` | 通常 log 破損、archive gzip 展開失敗、duration 欠落を除外し、WARN code、集計継続、破損内容非表示を固定する。 |
| [§27.14](runner.md#sec-27-14) | `failure-stats-invalid-query` | `days` / `n` 範囲外、未知 query、非整数で `422`、通常 log / archive log / history 差分なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-save-publish` | `.snapshots` / tmp mode、USTAR + BestCompression、固定 tar / gzip header、metadata key 順、file sync / directory sync / validation / rename 順、public snapshot 2 file だけ、tmp 非残存を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-save-source` | output root symlink、source symlink / hardlink / special file、read 中変更、列挙後追加・削除・置換、UTF-8 / USTAR 非表現 path、USTAR 非表現 size を個別に fake し、`failed`、public path 非作成、既存 snapshot 不変を固定する。 |
| [§27.15](archive.md#sec-27-15) | `partial-snapshot-save-directory-sync` | publish rename 後の `.snapshots` sync だけを失敗させ、`saved`、`SNAPSHOT_DIRECTORY_SYNC_FAILED`、public snapshot 維持、prune 非実行、build success 維持を固定する。 |
| [§27.15](archive.md#sec-27-15) | `partial-snapshot-save-tmp-cleanup` | save 成否それぞれの tmp cleanup 失敗で `SNAPSHOT_TMP_CLEANUP_FAILED` が付加され、先行する save 結果と public snapshot を反転・巻戻しないことを固定する。 |
| [§27.15](archive.md#sec-27-15) | `partial-snapshot-prune-failure` | 超過分の `delete_partial` と `delete_failed` を個別に fake し、残りの prune 続行、save 結果 `saved`、warning `SNAPSHOT_PRUNE_FAILED` 1 回、partial public path 非復元、failed public path 維持を固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-save-existing` | 既存 `meta.json` と `site.tar.gz` が全検証に合格し、id / build_id / output_sha256 が呼出入力と一致する場合に `exists_valid`、`SNAPSHOT_EXISTS`、`snapshot_id=build_id`、archive byte / metadata / mtime 変化なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-save-existing-mismatch` | 既存 snapshot の archive 破損、metadata 破損、id / build_id / output_sha256 不一致の各ケースで `failed`、`SNAPSHOT_EXISTS_MISMATCH`、`snapshot_id=null`、既存 directory / byte / metadata 変化なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-list-download` | `meta.json` と保存済み `site.tar.gz` の読取、非圧縮通常 file の size / count / manifest SHA-256 再計算、tar entry 辞書順、固定 header、API download header、検証に使用した同一 file descriptor からの stream、stream byte と保存済み `site.tar.gz` の完全一致、状態差分なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-download-unsafe-entry` | unsafe path、symlink、secret file、meta mismatch のいずれかで stream 開始前 `500`、binary header なし、状態差分なし、固定 server log を固定する。 |
| [§27.15](archive.md#sec-27-15) | `partial-snapshot-download-stream-failure` | 事前検証に使用した同一 `site.tar.gz` file descriptor の stream 開始後 read error で stream 中断、JSON error 追加なし、状態差分なし、`SNAPSHOT_STREAM_FAILED` を固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-delete` | id validation、runner owner の delete guard 取得・state 再確認、archive owner の事前検証・tombstone rename・2 回の directory sync・cleanup・`deleted` 返却、api owner の success config log → success audit、guard 解放 → response 順、tombstone 非残存、削除対象以外の snapshot 維持を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-invalid-id` | 空、`..`、encoded path separator、build id 形式不一致で `422`、archive owner 非呼出し、snapshot / config log / audit / history / build log / pending / state 差分なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-state-error` | guard 取得後の `.build_state` 破損または読取失敗で `500`、archive owner 非呼出し、所有確認付き guard 解放、snapshot / config log / audit 差分なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-corrupt` | `meta.json` 不正、archive 破損、manifest 不一致で `500`、snapshot を削除せず log / audit を追記しないことを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-before-publish` | tombstone rename 前の permission / I/O failure で `delete_failed`、HTTP `500`、public snapshot 維持、config / audit 非追記、guard 解放を固定する。 |
| [§27.15](archive.md#sec-27-15) | `partial-snapshot-delete-cleanup` | tombstone rename 後の最初の directory sync、tombstone cleanup、最終 directory sync の各失敗で `delete_partial`、public snapshot 非復元、config log `partial_failure` / `snapshot_delete_cleanup_failed`、failure audit、guard 解放、HTTP `500`、warning `SNAPSHOT_DELETE_CLEANUP_FAILED` を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-config-log` | snapshot 削除後の `.config_log` failure で audit 非試行、response `500`、削除済み snapshot を巻き戻さず、history / build log / pending / state unchanged を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-audit` | snapshot 削除と `.config_log` 成功後の `.audit_log` failure で response `500`、削除と config log を巻き戻さず、history / build log / pending / state unchanged を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-delete-guard-release` | snapshot 削除と config / audit log 成功後の所有確認または lock 削除 failure で `500`、snapshot と log を巻き戻さず、他の lock を削除しないことを固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-rollback` | rollback lock、new build id、archive 検証・同一 descriptor からの tmp 展開、status → state → running log の prepare、build trigger audit、worker 開始、展開済み file のみ転送、tmp cleanup、final log → history → status → state → lock 解放、元 snapshot / 元 log / `.last_sha` unchanged を固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-rollback-pending` | deploy 再送可能失敗で `source_kind="snapshot"`、`out=null`、snapshot id / manifest SHA-256 付き pending を保存し、build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、両方の `failure_category="deploy_failure"`、tmp cleanup、元 snapshot unchanged を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-no-checksum` | `meta.json.output_sha256=null` の snapshot 転送失敗で pending を作成せず `failure_build`、final log / history / status / state / lock 解放、元 snapshot unchanged を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-audit` | prepared 後の build trigger audit 失敗で worker / SSH / deploy 非実行、handle 即時無効化、abort 1 回、final log → history → status → state → tmp cleanup →所有確認付き lock 解放を各 1 回、`error="rollback audit failed"`、HTTP `500` 維持、補償失敗時も後続手順継続と再試行禁止を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-start-state` | status、state、running log の各書込について rename 前 failure と post-rename partial failure を別 subcase で fake する。prepared / audit / worker / SSH / deploy は 0 回。同一 `new_build_id` の schema-valid running log / status / state だけを final log → history → status → state の順で各 1 回補償し、不在、破損、id / trigger / running 不一致の対象は作成・修復・上書きしない。最後に tmp cleanup と所有確認付き lock 解放を各 1 回実行し、固定 `ROLLBACK_PREWORKER_*` ERROR、`SNAPSHOT_ROLLBACK_TMP_CLEANUP_FAILED`、primary error 維持、各失敗後の後続手順継続、再試行禁止を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-finalizer` | final log atomic replace、history append、status write、state write を個別に fake し、先行成功済み deploy と元 snapshot を巻き戻さない。final log 失敗で history は 0 回、history 失敗で再追記は 0 回、status / state / 所有確認付き lock 解放はそれぞれ最大 1 回、失敗後も後続手順を継続し、同じ write の再試行がないことを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-stale-recovery` | stale lock、running state、同一 id の running rollback log がすべて一致する場合だけ `rollback interrupted` で final log → history → status → state → tmp cleanup →所有確認付き lock 解放を各 1 回実行し、失敗後も後続手順を継続し、同じ write を再試行しない。一要素でも不足または不一致なら推測修復も通常 build 開始もしないことを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-rollback-target-unavailable` | 元 history の branch / target_file に一致する現在 target 不在または deploy target 空で `409`、lock / build id / tmp / audit / log / history / pending 差分なしを固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-running-conflict` | running 中または valid lock 存在中の delete / rollback で `409`、delete guard 非返却、snapshot / history / log / pending / config log / audit 差分なしを固定する。delete guard が lock 取得後に state running を検出した場合は所有確認付き解放を固定する。 |
| [§27.15](archive.md#sec-27-15) | `success-snapshot-pending-retry` | `source_kind="snapshot"` entry から保存済み snapshot を再検証し、同一 descriptor から新規 pending tmp へ展開して保存済み 1 target へ転送・remote checksum 検証後に entry 削除、tmp cleanup、branch target `out` 非参照を固定する。 |
| [§27.15](archive.md#sec-27-15) | `failure-snapshot-pending-source` | snapshot 不在、archive 破損、manifest 不一致を個別に fake し、SSH 非実行、対応する固定 `last_error`、`retry_count+1`、`failed_at` 更新、entry / 元 snapshot 維持を固定する。 |
| [§27.16](api.md#sec-27-16) | `success-health-ok` | `.build_status.json` 正常時の HTTP `200`、`status="ok"`、checks 空、pending / notify count、uptime fake clock、read-only no-write を固定する。 |
| [§27.16](api.md#sec-27-16) | `success-health-degraded` | pending transfer、missing status fallback、runner stale、notify pending read warning で `status="degraded"`、checks 順序、HTTP `200`、状態差分なしを固定する。 |
| [§27.16](api.md#sec-27-16) | `failure-health-read-error` | response 生成不能は `500`、response を構築できる read error は `200` と `status="degraded"` および固定 `checks`、自動修復なし、状態差分なしを固定する。 |
| [§27.16](api.md#sec-27-16) | `noop-health-readonly` | health を複数回呼んでも状態ファイル、access 対象外ファイル、log archive、notification が変わらないことを固定する。 |
| [§27.17](api.md#sec-27-17) | `success-log-search-level` | `level` 正規化、stdout / stderr / warnings / error 分類、内部 line number、query filter、build 単位の grouped `lines` 固定順を固定する。HTTP response に内部 metadata を追加しない。 |
| [§27.17](api.md#sec-27-17) | `success-log-search-archive` | 通常 log と archive log を同一分類で検索し、通常 log 優先、archive gzip 展開順、状態差分なしを固定する。 |
| [§27.17](api.md#sec-27-17) | `partial-log-search-corrupt-skip` | 破損 build log と gzip 展開失敗を除外し、固定 WARN code、検索継続、破損内容非表示を固定する。 |
| [§27.17](api.md#sec-27-17) | `failure-log-search-invalid-level` | 不正 level、未知 query、date 範囲不正で `422`、read-only no-write、archive 展開呼び出し 0 件または固定中断位置を固定する。 |
| [§27.18](api.md#sec-27-18) | `success-branch-config-get-default` | `.branch_config` 不在時の default 正規化、`source="default"`、secret 非表示、状態差分なしを固定する。 |
| [§27.18](api.md#sec-27-18) | `success-branch-config-post` | request `branches` 検証、`branch_targets` 保存、sort、`.config_log` diff、`config_update` audit、response `branches_count`、runner が次回起動で読む状態を固定する。 |
| [§27.18](api.md#sec-27-18) | `failure-branch-config-invalid-path` | 相対禁止 path、`..`、制御文字、空 branch、branch 重複、deploy target id 重複、host / user の禁止文字、dest_dir 不正、上限超過で `422`、`.branch_config` / `.config_log` / `.audit_log` 差分なしを固定する。 |
| [§27.18](api.md#sec-27-18) | `partial-branch-config-log-failure` | `.branch_config` 保存または削除成功後の `.config_log` 追記失敗で response `500`、audit 未実行、保存済み状態を巻き戻さないことを固定する。 |
| [§27.19](runner.md#sec-27-19) | `success-weekly-summary-auto` | fake clock 条件一致、自動集計、対象 channel 抽出、通知 payload、`.notify_log`、`.build_state.weekly_summary_*` 更新順を固定する。 |
| [§27.19](runner.md#sec-27-19) | `success-weekly-summary-manual` | 手動 API の認証、runner 内部の channel 別送信結果、`.notify_log`、sent date 非更新、API の固定 response payload を検証し、公開 response に `channel_results` を含めないことを固定する。 |
| [§27.19](runner.md#sec-27-19) | `noop-weekly-summary-same-day` | 同日自動送信済みで通知 0 件、`.notify_log` / `.notify_pending` / `.build_state` 差分なし、idempotency を固定する。 |
| [§27.19](runner.md#sec-27-19) | `failure-weekly-summary-send` | 宛先なし `422` または送信失敗 `500`、sent date 非更新、retry 対象時だけ `.notify_pending` 追加、build status 非変更を固定する。 |
| [§27.20](api.md#sec-27-20) | `success-config-diff-simple` | 単一 key 更新の normalized before / after、machine diff、`diff_text`、type、action、actor、request_id、endpoint、保存後 `.config_log` と `config_update` audit の追記順を固定する。 |
| [§27.20](api.md#sec-27-20) | `success-config-diff-nested` | nested object の dot path diff、配列全体比較、key 昇順、JSON 値表現、複数行値 escape を固定する。 |
| [§27.20](api.md#sec-27-20) | `noop-config-diff-same-value` | 正規化後同一値で対象状態ファイル、secret file、`.config_log`、`.audit_log` 差分なし、endpoint 固有の no-op response、idempotency を固定する。 |
| [§27.20](api.md#sec-27-20) | `security-config-diff-secret-mask` | key path に password / token / secret / pat / smtp_password を含む値を before / after と `diff_text` で `"***"` にし、request body / header / cookie 非保存を固定する。 |
| [§27.20](api.md#sec-27-20) | `partial-config-log-failure` | 主状態保存後の `.config_log` 追記失敗で `500`、audit 未実行、保存済み主状態を巻き戻さないことを固定する。 |
| [§27.20](api.md#sec-27-20) | `partial-config-audit-failure` | 主状態と `.config_log` 保存後の `.audit_log` 追記失敗で `500`、保存済み主状態と `.config_log` を巻き戻さないことを固定する。 |

[§27.12〜§27.20](fixture.md#sec-27-f-4) の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持つ。read-only、noop、invalid query、invalid signature、running conflict fixture では、対象状態ファイル、queue、history、build log、snapshot、notification、`.config_log`、`.audit_log` の forbidden side effect を必ず列挙する。[§27.15](archive.md#sec-27-15) の download fixture では `downloads[]` に `content_type`、`content_disposition`、`entry_order`、`stream_started`、`stream_interrupted`、`error_after_stream_start` を固定し、delete / rollback fixture では `write_order` と `unchanged_paths` に元 snapshot、元 build log、`.last_sha`、対象外 history / pending を必ず列挙する。

<a id="sec-27-f-5"></a>
**[fixture 証跡責務 §27-F 追加仕様化機能 §27.21〜§27.30 fixture 固定契約](fixture.md#sec-27-f-5)：**

[§27.21〜§27.30](fixture.md#sec-27-f-5) の fixture は、runner 拡張、builder cache / dependency、remote artifact、approval API の詳細実装確認を固定する。各 fixture は、実行順、保存順、成功時だけ更新する状態、失敗時に絶対変更してはならない状態、外部 API / command / SSH / notification の呼び出し、secret mask を expected に固定し、実装検証証跡に対象 fixture と実行結果を列挙する。

| 節 | fixture | 固定する内容 |
|----|---------|--------------|
| [§27.21](runner.md#sec-27-21) | `success-multi-file-one-change` | target_files 正規化、1 target changed、build id 1 件、`ADLAIRE_CHANGED_TARGETS`、changed target だけの SHA cache 更新、build log `changed_targets[]` を固定する。 |
| [§27.21](runner.md#sec-27-21) | `success-multi-file-many-change` | 複数 target の辞書順、build 1 回、全 changed target の before / after SHA、force build 時の全 target SHA 更新を固定する。 |
| [§27.21](runner.md#sec-27-21) | `noop-multi-file-all-skip` | 全 target unchanged で build / deploy / snapshot / history / notify なし、`.build_status.json` skip、SHA cache 差分なし、idempotency を固定する。 |
| [§27.21](runner.md#sec-27-21) | `failure-multi-file-path-traversal` | target path の絶対 path、`..`、NUL、改行で validation failure、build なし、SHA cache / history / log 差分なしを固定する。 |
| [§27.22](runner.md#sec-27-22) | `success-builder-command-extra-args` | `.pipeline_config.extra_args` を固定 builder argv の末尾へ配列順で追加し、固定引数の順序、shell 起動 0 件、出力 file set 検証を固定する。 |
| [§27.22](runner.md#sec-27-22) | `success-builder-command-env-merge` | 親 process env、branch env、`.pipeline_config.env`、runner 固定値の上書き順、child env、stdout/stderr 保存を固定する。 |
| [§27.22](runner.md#sec-27-22) | `success-builder-command-cache` | 標準 builder command で `build_cache_enabled` false / true ごとの `--cache-dir` 不在 / 1 回、`extra_args` 後置、deploy / SHA 更新条件を固定する。 |
| [§27.22](runner.md#sec-27-22) | `failure-builder-command-forbidden-key` | `.pipeline_config` に未定義 key を指定した場合、API `422` または runner `failure_pipeline_config` / 終了コード `2`、child process 0 件、deploy / snapshot / SHA cache 更新なしを固定する。 |
| [§27.22](runner.md#sec-27-22) | `failure-builder-command-reserved-env` | `.pipeline_config.env` に [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Process environment entry 共通固定契約](../DETAIL_INDEX.md#process-environment-entry-contract) の reserved key と prefix 代表値 `ADLAIRE_CI_TEST` を個別に指定し、API `422` または runner `failure_pipeline_config` / 終了コード `2`、child process 0 件、deploy / snapshot / SHA cache 更新なしを固定する。branch env は [§27.31 fixture 固定契約](#sec-27-f-6) だけで検証する。 |
| [§27.22](runner.md#sec-27-22) | `security-builder-command-secret-mask` | `.pipeline_config.env` / branch env / stdout / stderr / command args の secret 風値が log、notify、effects、response に平文で残らないことを固定する。 |
| [§27.23](runner.md#sec-27-23) | `success-local-watch-change` | GitHub API 0 件、local scan 辞書順、changed file 差分、trigger `local_watch`、build success 後の `.local_watch_state.json` 置換を固定する。 |
| [§27.23](runner.md#sec-27-23) | `noop-local-watch-no-change` | GitHub API / PAT verify 0 件、build なし、state 差分なし、status `skipped_no_change`、idempotency を固定する。 |
| [§27.23](runner.md#sec-27-23) | `failure-local-watch-state-corrupt` | state 破損で full build 扱い、build 成功時だけ state 再作成、build failure 時は破損 state 維持を固定する。 |
| [§27.23](runner.md#sec-27-23) | `failure-local-watch-tag-filter-conflict` | `watch_mode="local"` と tag filter enabled の併用で終了コード `2`、GitHub API 0 件、build / state 更新なしを固定する。 |
| [§27.24](runner.md#sec-27-24) | `success-tag-filter-match` | tag refs API、pattern match、matched_tags 最大 100 件、build 実行、build log 保存、SHA cache 更新条件を固定する。 |
| [§27.24](runner.md#sec-27-24) | `noop-tag-filter-unmatched` | tag 不一致で build id / build log / history / deploy / snapshot / notify なし、`.build_status.json.status="skipped"`、`last_target_status="skipped_tag_filter"`、SHA cache 未更新を固定する。 |
| [§27.24](runner.md#sec-27-24) | `failure-tag-filter-api` | tags API retry と最終失敗、build なし、終了コード `3`、SHA cache / history / snapshot 差分なしを固定する。 |
| [§27.24](runner.md#sec-27-24) | `failure-tag-filter-pattern` | 不正 pattern で API `422` または runner 終了コード `2`、tag API / build / SHA cache 更新なしを固定する。 |
| [§27.25](builder.md#sec-27-25) | `success-build-cache-hit` | cache key 一致、dependency SHA 一致、通常変換 byte 等価、cache_hits REPORT、公開出力 staging → rename を固定する。 |
| [§27.25](builder.md#sec-27-25) | `success-build-cache-miss` | miss 時の通常変換、cache entry tmp write → rename、cache_misses REPORT、secret / absolute path 非保存を固定する。 |
| [§27.25](builder.md#sec-27-25) | `partial-build-cache-byte-mismatch` | cache entry byte mismatch を hit 破棄 / miss にし、`BUILD_CACHE_ENTRY_INVALID`、当該 page file だけの `os.Remove` 1 回、削除失敗時の `BUILD_CACHE_ENTRY_CLEANUP_FAILED`、再試行 0 回、build success、absolute path / cache 内容の非出力、公開出力保護を固定する。 |
| [§27.25](builder.md#sec-27-25) | `failure-build-cache-save` | cache write failure でも build success、REPORT `cache_write_failures`、公開出力 success、既存 cache index / page 維持を固定する。 |
| [§27.25](builder.md#sec-27-25) | `failure-build-cache-cli-path` | 値欠落、空文字、不在、非 directory、symlink、`Lstat` 失敗、source / output 配下を個別に与え、固定 stderr、終了コード `1` / `2`、Markdown 読込・cache 読取・staging 作成 0 件、公開出力と既存 cache 不変を固定する。 |
| [§27.26](runner.md#sec-27-26) | `success-parallel-targets-all` | worker 上限、target ごとの started / finished、target_results 設定順、pending なし、status success を固定する。 |
| [§27.26](runner.md#sec-27-26) | `partial-parallel-targets-some-fail` | 一部 target failure、成功 target pending なし、失敗 target だけ pending、SSH / checksum 失敗の `error_code` と固定 `error`、build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、両方の `failure_category="deploy_failure"`、notify 順を固定する。 |
| [§27.26](runner.md#sec-27-26) | `failure-parallel-targets-all-fail` | timeout / internal failure を含む全 target failure の `error_code` と固定 `error`、build 本体 success の場合の build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、両方の `failure_category="deploy_failure"`、pending 全件、SHA cache 更新可否を固定する。 |
| [§27.26](runner.md#sec-27-26) | `success-parallel-targets-order-stable` | 完了順が入れ替わる fake result でも target_results / pending / history が設定順で保存されることを固定する。 |
| [§27.27](runner.md#sec-27-27) | `success-hook-pre-post` | pre → build → post の順、hook log、build log warning なし、通知前実行、secret mask を固定する。 |
| [§27.27](runner.md#sec-27-27) | `failure-hook-pre-abort` | pre abort で builder / pipeline / remote / deploy / snapshot / SHA cache 更新なし、history `hook_error`、hook log 保存を固定する。 |
| [§27.27](runner.md#sec-27-27) | `partial-hook-post-fail` | build status 維持、post hook failure log、runner 終了コード最低 `1`、notification 順序、secret mask を固定する。 |
| [§27.27](runner.md#sec-27-27) | `security-hook-shell-denied` | shell metachar が展開されず argv として渡ること、glob / env 展開 0 件、stdout/stderr secret mask を固定する。 |
| [§27.28](builder.md#sec-27-28) | `success-dependency-manifest` | link / image / HTML img / include 抽出、dep path 正規化、manifest tmp → rename、REPORT counts、runner 逆引きを固定する。 |
| [§27.28](builder.md#sec-27-28) | `success-dependency-missing` | missing dependency の WARN、non-strict 継続、strict 終了コード `2`、broken_dependencies、manifest 保存条件を固定する。 |
| [§27.28](builder.md#sec-27-28) | `failure-dependency-build-keeps-old` | build failure / strict failure で既存 `.dependency_manifest.json` 維持、tmp 非公開、公開出力保護を固定する。 |
| [§27.28](builder.md#sec-27-28) | `security-dependency-path-normalize` | base 外、credential URL、query token、absolute path を manifest に保存せず、WARN / broken reason / secret mask を固定する。 |
| [§27.29](runner.md#sec-27-29) | `success-remote-build-artifact` | local builder 0 件、共通 deadline、`ssh -- user@host remoteCommand` argv、remote command、`cat --` artifact stream、stdout/stderr drain、exclusive tmp 作成、fsync / close / rename / directory fsync、unsafe entry 検査、manifest 検証、deploy 連携、`.remote_artifacts/{build_id}` の `os.RemoveAll` 1 回を固定する。 |
| [§27.29](runner.md#sec-27-29) | `failure-remote-build-auth` | SSH auth failure の retry、最終 `failure_remote_build`、artifact fetch / deploy / snapshot / SHA cache 更新なし、secret 非保存を固定する。 |
| [§27.29](runner.md#sec-27-29) | `failure-remote-build-checksum` | manifest checksum mismatch で deploy なし、公開 output / snapshot 保護、remote log mask、`.remote_artifacts/{build_id}` の `os.RemoveAll` 1 回、cleanup 失敗時の `REMOTE_ARTIFACT_CLEANUP_FAILED`、再試行 0 回、先行結果と終了コード維持を固定する。 |
| [§27.29](runner.md#sec-27-29) | `security-remote-build-argument-quoting` | `work_dir`、`command_args`、`artifact_path` に空白、single quote、`$`、backtick、semicolon、glob 文字を含め、論理 argv が remote 側で byte 一致し、追加 command、glob、変数展開、command 置換が 0 件であることを固定する。 |
| [§27.29](runner.md#sec-27-29) | `security-remote-build-unsafe-archive` | tar.gz の `..`、absolute path、symlink、device を拒否し、一時展開外書き込み 0 件、deploy なしを固定する。 |
| [§27.30](runner.md#sec-27-30) | `success-approval-approve` | pending list、approve body 禁止、queue id 採番、`trigger="approval"` queue 追加、requested_force 引継ぎ、approved record、`approval_approved` audit、その後の systemctl 固定引数 1 回、queue id / dispatch response、再取得 response を固定する。 |
| [§27.30](runner.md#sec-27-30) | `partial-approval-dispatch-fallback` | approved record / audit 確定後の systemctl failure、`dispatch:"timer_fallback"` 成功 response、approved / queue 保持、runner timer 処理可能を固定する。 |
| [§27.30](runner.md#sec-27-30) | `noop-approval-reject` | reject で queue 追加なし、rejected record、history `approval_rejected`、`approval_rejected` audit、runner build なし、UI / SDK 再取得を固定する。 |
| [§27.30](runner.md#sec-27-30) | `noop-approval-expire` | fake clock timeout、expired record、history `approval_expired`、`approval_expired` audit、queue 追加なし、期限後 approve `409` を固定する。 |
| [§27.30](runner.md#sec-27-30) | `failure-approval-double-approve` | approved / rejected / expired への二重 approve で `409`、queue / approval / history 差分なし、audit / secret mask を固定する。 |
| [§27.30](runner.md#sec-27-30) | `partial-approval-approved-append` | queue append 後の approved append 失敗、runner の実行拒否、同じ approve 再試行での queue id 再利用、queue 非重複、approved / audit 完了を固定する。 |
| [§27.30](runner.md#sec-27-30) | `partial-approval-pending-audit` | pending append 後の audit 失敗、同一 pending 再検出での audit 補完、channel 単位の通知証跡判定、pending 非重複、build 抑止を固定する。 |

[§27.21〜§27.30](fixture.md#sec-27-f-5) の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持つ。failure、noop、partial、security fixture では、SHA cache、build log、history、status、snapshot、dependency manifest、build cache、local watch state、approval queue、pending transfer、remote artifact tmp、public output の forbidden side effect を必ず列挙する。

<a id="sec-27-f-6"></a>
**[fixture 証跡責務 §27-F runner 詳細本文責務 §27.31〜§27.38 fixture 固定契約](fixture.md#sec-27-f-6)：**

[`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31)〜[§27.38](runner.md#sec-27-38) の fixture は、[`docs/details/runner.md` 詳細本文責務 §27.31](runner.md#sec-27-31)〜[§27.38](runner.md#sec-27-38) 実装確認固定契約に列挙された branch env、notification、trend、chain、priority queue、failure classification、environment record、duration anomaly の状態、log、API response、副作用、保存順、secret mask を固定する。各 fixture は `manifest.json.section` を対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様).x に固定し、`manifest.json.feature` を [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の feature 名と一致させる。

| 節 | fixture 名 | 固定する確認 |
|----|------------|--------------|
| [§27.31](runner.md#sec-27-31) | `success-branch-env-inject` | `.branch_config.branch_targets[].env` 正規化、ASCII 昇順保存、system env 上書き、builder / pipeline / hook / command notification への env 注入、`.build_logs/{id}.json.environment.env_keys` を固定する。 |
| [§27.31](runner.md#sec-27-31) | `security-branch-env-secret-mask` | 有効な大文字 key の `TOKEN` / `SECRET` / `PASSWORD` / `PAT` 値が response、stdout、stderr、build log、history、notify log、pending、effects に残らないことを固定する。 |
| [§27.31](runner.md#sec-27-31) | `failure-branch-env-invalid-key` | lowercase の `my_token`、先頭数字、制御文字、上限超過 key / value を API `422` または runner 終了コード `2`、`.branch_config` / build log / history / status 差分なしに固定する。 |
| [§27.31](runner.md#sec-27-31) | `failure-branch-env-reserved-key` | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Process environment entry 共通固定契約](../DETAIL_INDEX.md#process-environment-entry-contract) の各 exact reserved key と prefix 代表値 `ADLAIRE_CI_TEST` を個別に入力し、API `422`、保存済み不正状態の runner 終了コード `2`、child process 0 件、全業務状態不変を固定する。 |
| [§27.31](runner.md#sec-27-31) | `failure-branch-env-mask-failure` | mask failure 時に build 完了扱いにせず、平文保存なし、失敗地点以降の write / command / notification 禁止を固定する。 |
| [§27.32](runner.md#sec-27-32) | `success-notify-multi-channel` | 複数 channel の event 判定、channel id 昇順送信、`.notify_log` の `result:"success"` / `error_code:null`、build status 不変を固定する。 |
| [§27.32](runner.md#sec-27-32) | `partial-notify-webhook-pending` | webhook 5xx / timeout の `http_5xx` / `timeout` failure log、pending 追加、保存順、後続 channel 継続、build status 不変を固定する。 |
| [§27.32](runner.md#sec-27-32) | `failure-notify-webhook-nonretryable-http` | 1xx / 3xx / 4xx の `http_1xx` / `http_3xx` / `http_4xx`、HTTP status、pending 差分なし、後続 channel 継続を固定する。3xx は redirect 先への外部呼出し 0 件。 |
| [§27.32](runner.md#sec-27-32) | `failure-notify-webhook-network` | response 前接続失敗の `network_error`、`http_status:null`、pending 差分なしを固定する。 |
| [§27.32](runner.md#sec-27-32) | `failure-notify-email` | SMTP 未設定 / 送信失敗 / timeout の `smtp_not_configured` / `smtp_error` / `timeout`、statefile 正本の固定 `error`、`http_status:null`、pending 差分なしを固定する。 |
| [§27.32](runner.md#sec-27-32) | `failure-notify-command` | command 起動失敗 / 非 0 / timeout の `command_error` / `timeout`、statefile 正本の固定 `error`、`http_status:null`、stdout/stderr 非保存、pending 差分なしを固定する。 |
| [§27.32](runner.md#sec-27-32) | `partial-notify-retry-exhausted` | 初回失敗を `attempt:1`、最後に許可された retry を `attempt:1+retry_count` とし、最終 attempt では `result:"dropped"` / `error_code:"retry_exhausted"` / `error:"retry exhausted"` の 1 record だけを追記する。同じ attempt の failure record なし、pending 削除、build status 不変を固定する。 |
| [§27.32](runner.md#sec-27-32) | `failure-notify-log-write` | 初回送信後の notify log write failure では pending 追加判定を行わない。pending retry 後の log write failure では entry を byte 単位で保持し、attempt 番号を消費しない。両 subcase で build status を変更せず、固定 server log code を記録する。 |
| [§27.32](runner.md#sec-27-32) | `noop-notify-disabled-event` | enabled=false または event 不一致時に外部送信 0 件、`.notify_log` / `.notify_pending` 差分なし、idempotency を固定する。 |
| [§27.32](runner.md#sec-27-32) | `security-notify-secret-mask` | webhook secret、SMTP password、command env secret が GET、backup、notify log、pending、command stdout/stderr、UI 表示に残らないことを固定する。 |
| [§27.33](runner.md#sec-27-33) | `success-trend-summary-update` | sample 追加、保持件数 prune、avg / median / p95 / anomaly_count 再計算、atomic write を固定する。 |
| [§27.33](runner.md#sec-27-33) | `success-trend-replace-build-id` | 同一 `build_id` sample 置換、重複なし、`finished_at` 昇順再整列、summary 全再計算を固定する。 |
| [§27.33](runner.md#sec-27-33) | `failure-trend-corrupt-rebuild` | `.build_trends.json` 破損 backup、`.build_history` 有効行からの再集計、skip warning、再集計不能時初期化を固定する。 |
| [§27.33](runner.md#sec-27-33) | `failure-trend-api-invalid-n` | `GET /api/stats/build-trends?n=` 不正値を `422`、状態差分なし、status / log 更新なしに固定する。 |
| [§27.34](runner.md#sec-27-34) | `success-chain-dag-order` | DAG 検証、topological order、同順位 config 出現順、同一 `chain_run_id`、chain summary を固定する。 |
| [§27.34](runner.md#sec-27-34) | `noop-chain-disabled-job` | disabled job 除外、実行 command なし、history / build log 未作成、enabled job への影響なしを固定する。 |
| [§27.34](runner.md#sec-27-34) | `failure-chain-cycle` | 循環依存を API `422`、保存差分なし、runner では chain 無効化して通常 build へ戻す境界を固定する。 |
| [§27.34](runner.md#sec-27-34) | `partial-chain-required-skip` | required dependency failure 後の `skipped_dependency_failed` history、build log 未作成、summary skipped count、後続 write 禁止を固定する。 |
| [§27.35](runner.md#sec-27-35) | `success-priority-urgent-first` | active なしで urgent / high / normal / low の waiting 選択順、同一 priority FIFO、waiting から active への atomic move、`GET /api/queue` の active / waiting 分離表示を固定する。 |
| [§27.35](runner.md#sec-27-35) | `success-priority-active-first` | active normal と waiting urgent が共存する場合に active を先に再実行し、waiting を変更しないことを固定する。 |
| [§27.35](runner.md#sec-27-35) | `failure-priority-created-seq-missing` | `created_seq` 欠落 entry を queue 破損として拒否し、状態変更と build がないことを固定する。 |
| [§27.35](runner.md#sec-27-35) | `failure-priority-invalid` | 不正 priority を API `422`、`.build_state` / history / log 差分なしに固定する。 |
| [§27.35](runner.md#sec-27-35) | `failure-priority-queue-full` | queue full 時 `429`、urgent でも既存 low entry を削除しないこと、write / command なしを固定する。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-timeout` | pipeline timeout を `pipeline_timeout`、evidence source / code / message / at、build log / history 保存一致に固定する。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-deploy` | SSH / checksum / pending transfer failure の build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、API 正規化状態 `success`、両方の `failure_category="deploy_failure"`、分類優先順位、`failure_category=deploy_failure` の API filter 一致を固定する。 |
| [§27.36](runner.md#sec-27-36) | `failure-failure-category-filter-invalid` | 列挙外 `failure_category` query を `422`、状態差分なしとし、固定列挙値 `unknown` は有効値として区別することを固定する。 |
| [§27.36](runner.md#sec-27-36) | `security-failure-evidence-mask` | evidence 最大 10 件、分類 evidence 先頭、secret / token / path 全体 / 入力値連結なし、history `status="success"` の `failure_category:null`、history `status="success_deploy_pending"` の `failure_category:"deploy_failure"` を固定する。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-record` | build id 採番直後、builder 起動前の environment 保存、GOOS / GOARCH / Go version / hostname / state_dir / disk free / captured_at を固定する。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-builder-version-record` | 標準 builder version process を 1 回だけ起動し、検証済み第 2 token を `builder_version` へ保存し、別値へ fallback しないことを固定する。 |
| [§27.37](runner.md#sec-27-37) | `failure-environment-builder-version-contract` | 標準 builder の timeout、非 `0`、stderr 非空、不正 token 数、名前不一致、runner とのバージョン不一致、空 `go=`、追加行を個別に固定し、builder 未起動、environment 未作成とする。 |
| [§27.37](runner.md#sec-27-37) | `failure-environment-write` | environment 保存失敗時に pipeline / builder / deploy / notification を起動せず、`.build_status.json.status="failure"`、`last_target_status="failure_state_write"`、終了コード `1` を固定する。 |
| [§27.37](runner.md#sec-27-37) | `security-environment-secret-excluded` | 環境変数 value、token、secret、PATH 全体、VCS revision が build log、history、effects に保存されないことを固定する。 |
| [§27.38](runner.md#sec-27-38) | `success-duration-anomaly-avg` | trend 更新前 summary による avg 超過判定、WARN、`flagged=true`、tag 追加、notify payload `threshold_source` を固定する。 |
| [§27.38](runner.md#sec-27-38) | `noop-duration-anomaly-insufficient-samples` | sample 数不足、avg / p95 null、disabled 設定時に判定なし、通知なし、trend 更新だけ行う条件を固定する。 |
| [§27.38](runner.md#sec-27-38) | `partial-duration-anomaly-notify-failure` | anomaly 判定後、retry 対象 channel の Webhook `5xx` または timeout による通知失敗、build success 維持、当該 channel だけの notify pending 追加、trend 保存、history flag 維持を固定する。retry 対象外 error では pending を作成しない。 |
| [§27.38](runner.md#sec-27-38) | `failure-duration-anomaly-invalid-config` | API `422`、runner では既定値補正なしで機能無効、状態差分なし、通知なしを固定する。 |

[§27.31〜§27.38](fixture.md#sec-27-f-6) の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持つ。failure、noop、partial、security fixture では、`.branch_config`、`.notify_config`、`.notify_log`、`.notify_pending`、`.build_trends.json`、`.build_chain_config`、`.build_state`、`.build_logs/{id}.json`、`.build_history`、`.build_status.json`、外部 command、通知、public output の forbidden side effect を必ず列挙する。

<a id="sec-27-f-7"></a>
**[fixture 証跡責務 §27-F ファイルセット固定契約](fixture.md#sec-27-f-7)：**

各 fixture は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表のファイルセットを持つ。該当しない入出力は `not-applicable.txt` を置くのではなく、`manifest.json` の `not_applicable` 配列に理由付きで記録する。fixture ごとの必要ファイルは推測せず、`manifest.json` の定義で確定する。

| ファイル | 必須 | 内容 | 禁止条件 |
|----------|------|------|----------|
| `manifest.json` | 必須 | fixture 名、対象節、機能名、分類、fake clock、owner component、collaborator component、参照仕様節、not_applicable 理由。 | 実行環境依存 path、乱数、実 secret。 |
| `input/request.json` | fixture 固定表の必須 input が HTTP request、SDK invocation、UI action の 1 件以上を含む場合に必須 | primary HTTP request、SDK invocation、UI action。schema は [request input 固定契約](#fixture-request-input-contract) に従う。 | Authorization header の実 token、secret 平文、未知 key。 |
| `input/cli.json` | fixture 固定表の必須 input が CLI または setup script 実行を含む場合に必須 | binary 名、argv、env、cwd、stdin、expected exit code。schema は [CLI input 固定契約](#fixture-cli-input-contract) に従う。 | 実 home path、実 credential path、時刻上書き key、未知 key。 |
| `input/state/` | 状態参照 fixture で必須 | 実行前状態ファイル一式。存在しない状態は `manifest.json` の `missing_state` に列挙する。 | 期待状態を混ぜること、実 secret。 |
| `input/files/` | file byte 列を入力にする fixture で必須 | `manifest.json.input_files` に列挙した Markdown、設定 JSON、archive、snapshot、hook file、Release asset、既存公開出力、remote artifact だけを同じ相対 path で配置する。 | `manifest.json.input_files` にない file、実外部サービスから取得した未固定 file。 |
| `input/fakes.json` | fake を 1 件以上使う fixture で必須 | [fake input root 固定契約](#fixture-fake-input-root-contract) の全 root key と、使用する fake の応答順。 | 未知 root key、実ネットワーク呼び出し前提、実 command 実行前提。 |
| `expected/response.json` | `manifest.json.assertions` に `response` がある場合に必須 | HTTP response の status、headers、body。schema は [HTTP response expected 固定契約](#fixture-http-response-expected-contract) に従う。 | SDK return/error、UI state、未定義 key、順序非決定配列。 |
| `expected/sdk_trace.json` | `manifest.json.assertions` に `sdk-trace` がある場合に必須 | SDK method から HTTP request までの呼出順、引数、request、結果種別。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | Authorization 値、secret 平文、API 外部副作用。 |
| `expected/sdk_return.json` | `manifest.json.assertions` に `sdk-return` がある場合に必須 | SDK success return と return 後 token state。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | SDK error、UI 表示値、補完済み response。 |
| `expected/sdk_error.json` | `manifest.json.assertions` に `sdk-error` がある場合に必須 | SDK error と error 後 token state。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | SDK success return、UI error text、secret 平文。 |
| `expected/ui_trace.json` | `manifest.json.assertions` に `ui-trace` がある場合に必須 | UI action、SDK call、refresh 順、disabled 遷移、field 消去。schema は [UI expected 固定契約](#fixture-ui-expected-contract) に従う。 | HTTP request、Authorization 値、secret 平文。 |
| `expected/ui_dom.json` | `manifest.json.assertions` に `ui-dom` がある場合に必須 | 表示・非表示 panel、text、field error、control state、one-time 領域。schema は [UI expected 固定契約](#fixture-ui-expected-contract) に従う。 | 生 DOM snapshot だけによる判定、secret 平文。 |
| `expected/stdout.txt` | `manifest.json.assertions` に `stdout` がある場合に必須 | stdout 完全一致。stdout なしを検証する場合は空ファイル。 | 現在時刻、絶対環境 path。 |
| `expected/stderr.txt` | `manifest.json.assertions` に `stderr` がある場合に必須 | stderr 完全一致。stderr なしを検証する場合は空ファイル。 | secret、実 token、実 URL credential。 |
| `expected/state/` | [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) の `manifest.json.assertions` に `state` がある場合に必須 | 実行後状態ファイル一式、または [state diff expected 固定契約](#fixture-state-diff-expected-contract) に従う `expected/state/state-diff.json`。[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `state` は [§28-F ファイルセット固定契約](#sec-28-f-3) の `expected/site/` と `expected/builder-output.json` に割り当てる。 | 期待しないファイルの混入、実行後ファイル一式と `state-diff.json` の併用、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) での代用。 |
| `expected/logs/` | `manifest.json.assertions` に `logs` がある場合に必須 | build log、history、access log、audit log、notify log の期待値。 | secret 平文、実 Authorization header。 |
| `expected/events.json` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) または [MCP fixture 固定契約](#mcp-fixture-contract) が protocol event evidence を要求する場合に必須 | SSE frame、MCP notification、resource-updated、keepalive、disconnect の protocol event 期待値。fixture 名別 expected 固定表で要求し、`manifest.json.assertions` には追加しない。 | UI DOM、SDK return、状態 file、secret 平文。 |
| `expected/effects.json` | 必須 | 外部 API 呼び出し、command 実行、通知送信、download/stream 中断、呼び出し 0 件の期待値。 | 呼び出し順未指定、実外部送信。 |
| `expected/security.json` | `manifest.json.assertions` に `secret-mask` がある場合に必須 | secret 非表示確認対象、禁止文字列、token hash 検証、scope 判定、rate count。 | secret を検証用に平文保存すること。 |

`expected/state/` は、fixture が検証対象とする状態ファイルだけを含める。変更してはならない状態ファイルは `expected/effects.json` の `unchanged_paths` に列挙する。削除されるべきファイルは `expected/effects.json` の `deleted_paths` に列挙し、空 directory の存在可否も明記する。

<a id="fixture-request-input-contract"></a>
**request input 固定契約：**

`input/request.json` は `http`、`sdk`、`ui` の 3 root key だけをすべて持ち、未知 root key を禁止する。対象外 interface は `null` とする。fixture 固定表の必須 input が HTTP request を含む場合は `http`、SDK invocation を含む場合は `sdk`、UI action を含む場合は `ui` を `null` にしてはならない。`manifest.json.components` への component 追加だけを理由に、直接実行しない interface の object を作成してはならない。

```json
{
  "http": {
    "method": "POST",
    "path": "/api/build",
    "query": {},
    "headers": {
      "content-type": "application/json"
    },
    "body_present": false,
    "body": null
  },
  "sdk": {
    "method": "triggerBuild",
    "args": []
  },
  "ui": {
    "action": "click",
    "target": "btn-build",
    "value": null
  }
}
```

`http` object は `method`、`path`、`query`、`headers`、`body_present`、`body` の 6 key、`sdk` object は `method`、`args` の 2 key、`ui` object は `action`、`target`、`value` の 3 key だけを持つ。HTTP method、path、query、body 条件は対象 API request schema、SDK method と args は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、UI action と target は [UI expected 固定契約](#fixture-ui-expected-contract) の語彙に一致させる。header 名は lowercase ASCII、key は ASCII 昇順とする。synthetic secret を必要とする `headers`、`args`、`value` は `${secret:<source_id>}` を使い、`expected/security.json.forbidden_plaintexts[].id` を参照する。`body_present=false` は `body=null`、`body_present=true` は対象 API request schema と一致する JSON value とする。複数 request / invocation / action の後続入力は `input/fakes.json` に実行順で置き、本 file に array や追加 key を作成してはならない。

`expected/security.json` を除く fixture JSON に `${secret:<source_id>}` を 1 件以上置く fixture は、`manifest.json.assertions` に `secret-mask` を必ず含め、`expected/security.json` を必ず置く。各 `source_id` は `expected/security.json.forbidden_plaintexts[].id` の同名要素 1 件と完全一致させ、placeholder から参照されない同名要素、参照先のない placeholder、同じ `id` の重複を禁止する。

<a id="fixture-cli-input-contract"></a>
**CLI input 固定契約：**

`input/cli.json` は次の 6 root key だけを持ち、未知 key を禁止する。

```json
{
  "binary": "adlaire-ci-runner",
  "argv": ["--state-dir", "input/state"],
  "env": {},
  "cwd": ".",
  "stdin": null,
  "expected_exit_code": 0
}
```

`binary` は fixture が起動する repository 配布 executable または setup script の basename、`argv` は argv[0] を除く string array、`env` は明示的に渡す環境変数名と string 値だけの object、`cwd` は fixture root からの相対 directory、`stdin` は string または入力なしの `null`、`expected_exit_code` は `0`〜`255` の integer とする。`env` key は ASCII 昇順とし、host process から暗黙継承する key を expected に使用してはならない。secret 値は request input と同じ placeholder を使う。時刻は `manifest.json.fake_clock` だけを使用し、clock key、絶対 cwd、`..` segment、credential 付き argv を禁止する。

<a id="fixture-state-diff-expected-contract"></a>
**state diff expected 固定契約：**

`expected/state/state-diff.json` は、byte 比較用の実ファイル一式では固定できない file list、directory、mode、size、hash、mtime を検証する場合だけ使用する。同じ fixture の `expected/state/` に `state-diff.json` 以外の file または directory を置いてはならない。

```json
{
  "roots": ["install/admin"],
  "entries": [
    {
      "path": "install/admin",
      "type": "directory",
      "mode": "0755",
      "size": null,
      "sha256": null,
      "mtime": null
    },
    {
      "path": "install/admin/index.html",
      "type": "file",
      "mode": "0644",
      "size": 1234,
      "sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
      "mtime": "2026-09-16T00:00:00Z"
    }
  ]
}
```

root object は `roots`、`entries` の 2 key だけを持つ。`roots` は fixture 実行 root からの `/` 区切り相対 directory path を 1 件以上、ASCII 昇順、重複なしで持ち、相互の包含を禁止する。`entries` は各 `roots` 自身とその実行後の全 descendant を `path` の ASCII 昇順で過不足なく列挙し、各要素は `path`、`type`、`mode`、`size`、`sha256`、`mtime` の 6 key だけを持つ。各 `path` はいずれか 1 件の root と一致するか、その root の `/` 以下にある相対 path とする。絶対 path、空文字、`.`、`..` segment、backslash、symlink、socket、device、FIFO を禁止する。

`type` は `file` または `directory`、`mode` は特殊 bit を含まない正規表現 `^0[0-7]{3}$` の string とする。`type="file"` では `size` を 0 以上の integer、`sha256` を 64 桁 lowercase hexadecimal とする。`type="directory"` では `size=null`、`sha256=null` とする。`mtime` は比較しない場合の `null` または `manifest.json.fake_clock` と同じ UTC 秒精度形式の固定時刻とする。実行後に `roots` のいずれかの下へ未列挙 path が 1 件でもある場合、列挙 entry が不在、type / mode / size / sha256 / 比較対象 mtime が不一致、または path が複数 root に属する場合は `state` assertion を失敗とする。

<a id="fixture-http-response-expected-contract"></a>
**HTTP response expected 固定契約：**

`expected/response.json` は次の 3 root key だけを持ち、未知 key を禁止する。

```json
{
  "status": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "status": "ok"
  }
}
```

`status` は `100`〜`599` の integer、`headers` は lowercase ASCII header 名と string 値の object、`body` は JSON value または body なしを表す `null` とする。`headers` の key は ASCII 昇順とし、同名 header の複数値は HTTP 受信順に `, ` で連結した 1 string とする。`body` object の key、型、配列順は対象 API response schema と完全一致させ、SDK return、SDK error、UI 表示用既定値を追加してはならない。

<a id="fixture-sdk-expected-contract"></a>
**SDK expected 固定契約：**

`expected/sdk_trace.json` は `calls` だけを持ち、各要素は次の 11 key だけを持つ。未知 root key と未知要素 key を禁止する。

```json
{
  "calls": [
    {
      "order": 1,
      "method": "triggerBuild",
      "args": [],
      "request_method": "POST",
      "path": "/api/build",
      "query": {},
      "body_present": false,
      "body": null,
      "authorization_present": true,
      "result": "return",
      "status": 202
    }
  ]
}
```

`order` は 1 から始まる連続整数、`method` は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の SDK 列にある public method または `StreamHandle.close`、`args` は呼出時引数を順序どおり表す JSON array とする。secret 引数は平文を置かず、`expected/security.json.forbidden_plaintexts[].id` を参照する string `${secret:<source_id>}` に置換する。`request_method` は `GET`、`POST`、`PUT`、`PATCH`、`DELETE` のいずれかとする。HTTP 送信前 `TypeError` と `StreamHandle.close` だけは `request_method=null`、`path=null`、`query={}`、`body_present=false`、`body=null`、`authorization_present=false` とする。HTTP request を行う場合の `path` は query を含まない `/api/` 始まりの絶対 path、`query` は送信した string key / string value だけの object とし、送信なしは空 object とする。`body_present=false` では `body=null`、`body_present=true` では送信 JSON value を置く。`authorization_present` は header の有無だけを表し、token 値と `authorization_value` key を禁止する。`result` は `return`、`error`、`type_error`、`return_with_terminal_error` のいずれかとする。`return_with_terminal_error` は `streamBuild()` が `StreamHandle` を返した後に `done` が reject する場合だけ使用する。`status` は初期 HTTP status、network / timeout / protocol error の `0`、または HTTP request がない場合の `null` とする。配列順は実行順とし、SDK method または `StreamHandle.close` 1 回につき 1 要素を記録する。

`expected/sdk_return.json` は `returns` だけを持ち、各要素は `order`、`method`、`value_type`、`value`、`token_state` の 5 key だけを持つ。

```json
{
  "returns": [
    {
      "order": 1,
      "method": "triggerBuild",
      "value_type": "json",
      "value": {
        "message": "Build queued"
      },
      "token_state": "unchanged"
    }
  ]
}
```

`order` は対応する `sdk_trace.json.calls[].order`、`method` は同じ call の method と完全一致させる。`value_type` は `json`、`blob`、`stream_handle` のいずれかとする。`json` の `value` は API success response と key、型、値、配列順まで一致する JSON value とする。`blob` の `value` は `size`、`content_type`、`sha256` の 3 key だけを持ち、`size` は 0 以上の integer、`content_type` は string、`sha256` は 64 桁 lowercase hexadecimal とする。`stream_handle` の `value` は `closed`、`error`、`done` の 3 key だけを持ち、`closed` は boolean、`error` は `null` または `expected/sdk_error.json.errors[]` から `order`、`method`、`phase`、`token_state` を除いた exact object、`done` は `null` または `status` と `duration_seconds` だけを持つ `StreamEnd` とする。`token_state` は `set`、`cleared`、`unchanged` のいずれかとし、その call 完了直後の SDK token state transition を表す。`returns` は `order` 昇順、重複なしとする。

`expected/sdk_error.json` は `errors` だけを持ち、各要素は `order`、`method`、`phase`、`name`、`status`、`message`、`details`、`response_body`、`token_state` の 9 key だけを持つ。

```json
{
  "errors": [
    {
      "order": 1,
      "method": "getStatus",
      "phase": "call",
      "name": "AdlaireCIError",
      "status": 401,
      "message": "Unauthorized",
      "details": null,
      "response_body": {
        "error": "Unauthorized"
      },
      "token_state": "cleared"
    }
  ]
}
```

`order` と `method` は対応する SDK trace call と完全一致させる。`phase` は public method が reject / throw する `call`、または `StreamHandle` resolve 後に `done` が reject する `terminal` のいずれかとする。`terminal` は `streamBuild` だけに許可する。`name` は `AdlaireCIError` または HTTP 送信前引数不正の `TypeError` とする。`AdlaireCIError` の `status` は network / timeout / protocol error の `0`、または `100`〜`599` の integer、`details` は API の array または `null`、`response_body` は object、string、`null` のいずれかとする。`TypeError` は `phase="call"`、`status=null`、`details=null`、`response_body=null` とする。`message` は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) の固定文言または API error 値と完全一致させる。`token_state` は `set`、`cleared`、`unchanged` のいずれかとし、その error 確定直後の SDK token state transition を表す。`errors` は `order` 昇順、重複なしとする。

`sdk_trace.json.calls[].result="return"` は `sdk_return.json` の同一 `order` だけを、`result="error"` または `"type_error"` は `sdk_error.json` の同一 `order` かつ `phase="call"` だけを必須とする。`result="return_with_terminal_error"` は `sdk_return.json` と `sdk_error.json` の同一 `order` を 1 件ずつ必須とし、error の `phase` を `terminal` とする。この場合を除き、同じ call を return と error の両方へ記録してはならない。`sdk_trace.json` にない call、欠番、追加 return / error を禁止する。secret placeholder の解決値は比較時だけ使用し、差分、error、log、更新済み expected へ出力してはならない。

<a id="fixture-ui-expected-contract"></a>
**UI expected 固定契約：**

`expected/ui_trace.json` は次の 5 root key だけを持ち、未知 root key と各 array 要素の未知 keyを禁止する。

```json
{
  "actions": [
    {
      "order": 1,
      "action": "click",
      "target": "btn-build",
      "value": null
    }
  ],
  "sdk_calls": [
    {
      "order": 1,
      "method": "triggerBuild",
      "args": [],
      "result": "return"
    }
  ],
  "refresh_order": ["getStatus", "getQueue"],
  "disabled_transitions": [
    {
      "order": 1,
      "target": "btn-build",
      "disabled": true,
      "reason": "sending"
    }
  ],
  "cleared_fields": [
    {
      "order": 1,
      "target": "field-password",
      "reason": "failed"
    }
  ]
}
```

`actions[]` は `order`、`action`、`target`、`value`、`sdk_calls[]` は `order`、`method`、`args`、`result`、`disabled_transitions[]` は `order`、`target`、`disabled`、`reason`、`cleared_fields[]` は `order`、`target`、`reason` の key だけを持つ。各 `order` は array 内で 1 から始まる連続整数とする。`actions[].action` は `load`、`click`、`submit`、`input`、`change`、`copy`、`panel-transition`、`timer`、`stream-event`、`dialog-confirm`、`dialog-cancel` のいずれか、`target` は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の DOM id または `document`、`window`、`stream-handle`、`confirmation-dialog` のいずれかとする。`sdk_calls[].method` は SDK public method または `StreamHandle.close`、`result` は `return`、`error`、`return_with_terminal_error` のいずれかとし、SDK method 呼出し 1 回につき 1 要素を記録する。`refresh_order` は変更 API 成功後に実行した read SDK method 名を実行順で持ち、再取得なしは空配列とする。`reason` は [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の UI operation state、disabled 条件、secret 消去条件に記載された固定語を使用する。secret 値は SDK trace と同じ placeholder を使う。

`expected/ui_dom.json` は次の 7 root key だけを持ち、未知 root key と各 array 要素の未知 key を禁止する。

```json
{
  "visible_panels": ["panel-build"],
  "hidden_panels": [],
  "texts": [
    {
      "target": "build-success",
      "text": "Build queued"
    }
  ],
  "field_errors": [],
  "disabled_controls": [],
  "enabled_controls": ["btn-build"],
  "one_time_regions": []
}
```

`visible_panels`、`hidden_panels`、`disabled_controls`、`enabled_controls` は DOM id string の array、`texts[]` は `target` と `text`、`field_errors[]` は `field` と `message`、`one_time_regions[]` は `target`、`present`、`source_id` の key だけを持つ。各 array は最終 DOM の document order とし、同一 id の重複を禁止する。表示・非表示、disabled・enabled の同一 id 重複を禁止する。`source_id` は `expected/security.json.forbidden_plaintexts[].id` または secret を含まない領域の `null` とし、secret 平文を `text` または `message` に置いてはならない。生 DOM snapshot、CSS class の存在だけ、画面画像だけで `ui-dom` を合格にしてはならない。

<a id="fixture-fake-input-root-contract"></a>
**fake input root 固定契約：**

`input/fakes.json` は次の root object を使用し、未知 key を禁止する。全 root は呼び出し順の array とし、未使用 fake も空配列で残す。fixture 時刻は `manifest.json.fake_clock` だけを正本とし、`input/fakes.json` に clock root または時刻上書き値を置いてはならない。

```json
{
  "entropy": [],
  "filesystem": [],
  "timer": [],
  "listener": [],
  "signal": [],
  "github": [],
  "ssh": [],
  "smtp": [],
  "webhook": [],
  "systemd": [],
  "command": [],
  "fetch": [],
  "sdk": [],
  "archive": [],
  "download": [],
  "git": [],
  "mtime": [],
  "manifest": [],
  "cache": [],
  "clipboard": []
}
```

各 array 要素は次の共通 envelope だけを使用する。未知の envelope key を禁止する。`order` は root array ごとに 1 から始まる連続整数とする。使用しない fake kind に event を入れてはならない。

```json
{
  "order": 1,
  "operation": "call",
  "target": "github_api",
  "input": {},
  "output": {},
  "result": "success",
  "error_code": null
}
```

| key | 型 | 必須 | 固定契約 |
|-----|----|------|----------|
| `order` | integer | 必須 | root array 内で 1 から始まる連続整数。重複と欠番を禁止する。 |
| `operation` | string | 必須 | 対象 root の操作語彙固定表にある値。表にない値の任意追加を禁止する。 |
| `target` | string | 必須 | 安定した論理対象識別子または credential を含まない対象 path。環境固有の絶対 path、userinfo、token、secret を禁止する。 |
| `input` | object | 必須 | 対象 owner の入力 / request / payload 契約で定義した key だけを持つ。入力なしは空 object とする。 |
| `output` | object | 必須 | 対象 owner の response / result 契約で定義した key だけを持つ。出力なしまたは出力確定前の失敗は空 object とする。 |
| `result` | string | 必須 | `success`、`failure`、`timeout`、`cancelled`、`interrupted` のいずれか。 |
| `error_code` | string/null | 必須 | `result="success"` では `null`。それ以外では対象 owner 詳細本文の固定 error code。自由文の error message を禁止する。 |

| fake root | `operation` 許容値 |
|-----------|--------------------|
| `entropy` | `generate` |
| `filesystem` | `read`、`write`、`stat`、`mkdir`、`rename`、`remove`、`chmod`、`sync` |
| `timer` | `wait` |
| `listener` | `listen`、`serve`、`shutdown`、`close` |
| `signal` | `receive` |
| `github` | `call` |
| `ssh` | `execute` |
| `smtp` | `send` |
| `webhook` | `send` |
| `systemd` | `execute` |
| `command` | `execute` |
| `fetch` | `request`、`stream` |
| `sdk` | `call` |
| `archive` | `create`、`read`、`extract`、`validate` |
| `download` | `download` |
| `git` | `execute` |
| `mtime` | `read` |
| `manifest` | `read`、`write`、`validate` |
| `cache` | `read`、`write`、`remove` |
| `clipboard` | `write` |

`input` と `output` の key および値は、fixture カタログと対象 owner 詳細本文に一致させる。owner 詳細本文が固定していない場合は fixture で補完せず、先に対象 owner 詳細本文を改訂する。fake の secret は `input/` の synthetic 値だけを使用し、実 credential を含めてはならない。

<a id="fixture-fake-adapter-binding-contract"></a>
**fake adapter 接続固定契約：**

fixture harness は、対象 component の本番 entrypoint と同じ処理を、明示的に渡した runtime dependency set で構築する。runtime dependency set は少なくとも clock / timer、entropy、filesystem、listener、signal、外部 HTTP、systemd、command を個別 adapter として保持し、package global の差替え、環境変数による隠れた fake 選択、実装関数の直接 monkey patch、fixture 専用分岐を使用してはならない。本番 entrypoint は同じ dependency set に標準 clock / timer、`crypto/rand.Reader`、OS filesystem、`net.Listen` / `http.Server`、`os.Signal`、redirect 非追従を設定した HTTP client、引数配列で起動する systemd / command adapter を設定する。

fixture adapter は `manifest.json.fake_clock` と `input/fakes.json` だけから結果を返す。clock 呼出しは `manifest.json.fake_clock` の同一時刻を返す。timer の待機は実時間を使用せず harness が明示的に時刻を進めた時だけ発火し、同時刻に期限を迎える timer は作成順に発火する。時間経過を必要とする fixture は fixture step が明示する次の時刻へ harness が進める。entropy event は `input={"byte_count":<positive integer>}`、成功時 `output={"bytes_hex":"<byte_count*2 lowercase hex>"}`、失敗または短い読取時 `output={}` とする。要求 byte 数、hex 長、hex 形式が一致しない成功 event は fixture 定義不正として component を実行しない。production と fixture の両 adapter は要求 byte 数を完全に満たした場合だけ成功とし、短い読取を成功扱いにしない。

timer event の `target` は owner 詳細本文または fixture catalog が固定する stable timer id とし、`wait.input` は `{"duration_milliseconds":<non-negative integer>}`、正常終了の `output` は `{"elapsed_milliseconds":<same integer>}` とする。fake timer は実時間を待たず、manifest fake clock を同じ millisecond 数だけ進める。`cancelled` または `interrupted` は `output.elapsed_milliseconds` に `0` 以上かつ要求値未満の実経過値を返す。負数、小数、要求値を超える値を禁止する。

listener event の `target` は `api_listener` とする。`listen.input` は `{"addr":"127.0.0.1:<port>"}`、`serve.input` は `{"read_header_timeout_seconds":5,"read_timeout_seconds":30,"write_timeout_seconds":0,"idle_timeout_seconds":60,"max_header_bytes":32768,"error_log_mode":"redacted"}`、`shutdown.input` は `{"timeout_seconds":10}`、`close.input` は `{}` に固定する。`<port>` は API CLI 固定契約に合格した decimal port とし、fixture case が指定する実値へ置き換える。signal event の `target` は `api_process`、`receive.input` は `{}`、成功時 `receive.output` は `{"name":"SIGTERM"}` または `{"name":"SIGINT"}` だけを許可する。実 port bind、実 process signal、実 sleep を fixture で使用してはならない。

| root / operation | success の `output` / `error_code` | success 以外の固定結果 |
|------------------|-----------------------------------|--------------------------|
| `listener.listen` | `{}` / `null` | `result="failure"`、`output={}`、`error_code="LISTEN_FAILED"`。 |
| `listener.serve` | graceful shutdown では `{"server_error":"http.ErrServerClosed"}` / `null` | unexpected server error は `result="failure"`、`output={}`、`error_code="SERVE_FAILED"`。 |
| `listener.shutdown` | `{}` / `null` | timeout は `result="timeout"` / `SHUTDOWN_TIMEOUT`、failure は `result="failure"` / `SHUTDOWN_FAILED`。いずれも `output={}`。 |
| `listener.close` | `{}` / `null` | `result="failure"`、`output={}`、`error_code="CLOSE_FAILED"`。 |
| `signal.receive` | `{"name":"SIGTERM"}` または `{"name":"SIGINT"}` / `null` | fixture では success 以外を使用しない。 |
| `timer.wait` | `{"elapsed_milliseconds":<requested>}` / `null` | `cancelled` は `TIMER_CANCELLED`、`interrupted` は `TIMER_INTERRUPTED`。 |

各 adapter は対応 root の未消費先頭 event だけを 1 回消費する。adapter 呼出しがあるのに event がない、`operation` / `target` / `input` が一致しない、failure event を成功扱いする、または fixture 終了時に `manifest.json.not_applicable` で明示していない event が残る場合は fixture を失敗とする。並行 fixture は harness barrier で request の coordinator 取得順を固定し、その順序と fake 消費順を一致させる。実 filesystem、実 network、実 command、実 systemd、実 clipboard への fallback を禁止する。

`listener` と `signal` は process lifecycle を駆動する harness control event とし、`expected/effects.json` の process 外副作用 array へ重複記録しない。呼出回数、順序、入力、結果は当該 root の完全消費によって検証する。それ以外の process 外副作用 fake は [相互整合固定契約](#sec-27-f-13) に従って `expected/effects.json` の対応 array と一致させる。

<a id="sec-27-f-manifest-identity"></a>
**[fixture 証跡責務共通 manifest 識別子レジストリ固定契約](fixture.md#sec-27-f-manifest-identity)：**

`manifest.json` の `name`、`section`、`feature`、`owner_component` は、本節のレジストリと対象 fixture 固定表の組み合わせだけから決定する。節番号の範囲、fixture 名、配置 path、対象 component の列挙順から値を推測してはならない。同じ fixture 名が複数の固定表に現れる場合は、すべての出現箇所が同じ `section`、`feature`、`owner_component` を指すことを必須とし、異なる場合は仕様不整合として fixture 作成と実装を停止する。

[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の節別 fixture 固定表に記載された fixture は、次の識別子だけを使用する。fixture 名は該当節の固定表に列挙された値だけを許可する。

| `section` | `feature` | `owner_component` |
|-----------|-----------|-------------------|
| `27.1` | `commit_status` | `commitstatus` |
| `27.2` | `dry_run` | `runner` |
| `27.3` | `retry` | `runner` |
| `27.4` | `output_meta` | `builder` |
| `27.5` | `config_validate` | `api` |
| `27.6` | `api_access_log` | `api` |
| `27.7` | `log_archive` | `archive` |
| `27.8` | `build_status` | `runner` |
| `27.9` | `trigger` | `runner` |
| `27.10` | `startup_integrity` | `runner` |
| `27.11` | `schedule` | `api` |
| `27.12` | `webhook` | `api` |
| `27.13` | `webhook_events` | `api` |
| `27.14` | `stats` | `runner` |
| `27.15` | `snapshot` | `archive` |
| `27.16` | `health` | `api` |
| `27.17` | `log_search` | `api` |
| `27.18` | `branch_config` | `api` |
| `27.19` | `weekly_summary` | `runner` |
| `27.20` | `config_diff` | `api` |
| `27.21` | `multi_file` | `runner` |
| `27.22` | `builder_command_extension` | `runner` |
| `27.23` | `local_watch` | `runner` |
| `27.24` | `tag_filter` | `runner` |
| `27.25` | `build_cache` | `builder` |
| `27.26` | `parallel_targets` | `runner` |
| `27.27` | `hook` | `runner` |
| `27.28` | `dependency_manifest` | `builder` |
| `27.29` | `remote_build` | `runner` |
| `27.30` | `approval` | `runner` |
| `27.31` | `branch_env` | `runner` |
| `27.32` | `notification` | `runner` |
| `27.33` | `trend` | `runner` |
| `27.34` | `chain` | `runner` |
| `27.35` | `priority_queue` | `runner` |
| `27.36` | `failure_category` | `runner` |
| `27.37` | `environment` | `runner` |
| `27.38` | `duration_anomaly` | `runner` |
| `27.42` | `api_token_scope` | `security` |
| `27.43` | `api_key_management` | `security` |
| `27.44` | `audit_log` | `security` |
| `27.45` | `session_timeout` | `security` |
| `27.46` | `totp` | `security` |
| `27.47` | `api_rate_limit` | `security` |

節別 fixture 固定表以外の owner / 連動 fixture は、次の識別子だけを使用する。「対象 fixture 名」はリンク先固定表の第 1 列に列挙された fixture 名の完全一致を表し、行単位指定がある場合はその fixture 名だけに適用する。

| 対象 fixture 名 | `section` | `feature` | `owner_component` | fixture 名の正本 |
|-----------------|-----------|-----------|-------------------|------------------|
| security 詳細固定表の `認証共通` 行にある全 fixture | `security-auth-common` | `auth_common` | `security` | [security fixture 固定表](#sec-27-f-3) |
| runner / statefile 連動固定表の全 fixture | `22.0s` | `runner_state_integration` | `statefile` | [runner / statefile 連動 fixture 固定表](#sec-27-f-15) |
| statefile owner 固定表の全 fixture | `22.0s` | `statefile_contract` | `statefile` | [statefile owner fixture 固定表](#sec-27-f-16) |
| `success-api-sdk-ui-request-trace` | `22.0e` | `api_sdk_ui_integration` | `api` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| `failure-api-sdk-ui-error-propagation` | `22.0` | `api_sdk_ui_integration` | `api` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| `partial-api-sdk-ui-refresh-order` | `24` | `api_sdk_ui_integration` | `ui` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| `security-api-sdk-ui-secret-one-time` | `27.43` | `api_sdk_ui_integration` | `security` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| `security-api-sdk-ui-no-speculation` | `23` | `api_sdk_ui_integration` | `sdk` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| `security-api-sdk-ui-side-effect-boundary` | `22.0` | `api_sdk_ui_integration` | `api` | [API / SDK / UI 連動 fixture 固定表](#sec-27-f-17) |
| UI owner 固定表の全 fixture | `24` | `ui_contract` | `ui` | [UI owner fixture 固定表](#sec-27-f-18) |
| `success-setup-admin-release-asset-layout` | `26` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `failure-setup-download-boundary` | `26.2b` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `failure-setup-api-version-cohort` | `26.3b` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `security-setup-admin-archive-boundary` | `A2` | `setup_admin_integration` | `admin` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `partial-setup-systemd-rollback-boundary` | `26.5` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `partial-setup-api-runner-dispatch` | `26.4` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `security-admin-static-serving` | `A3` | `setup_admin_integration` | `admin` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `security-setup-secret-preservation` | `26.2b` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `partial-admin-cli-lifecycle` | `A7` | `admin_cli` | `admin` | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract) |
| `success-admin-cli-transport` | `A7` | `admin_cli` | `admin` | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract) |
| `failure-admin-cli-output-errors` | `A7` | `admin_cli` | `admin` | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract) |
| `security-admin-cli-secret-redaction` | `A7` | `admin_cli` | `admin` | [Admin CLI fixture 固定契約](#admin-cli-fixture-contract) |
| Release owner固定表の全fixture | `R1-R7` | `release_contract` | `release` | [Release fixture固定表](#release-fixture-contract) |
| `additional-management-users-roles-success` | `27.48-27.50` | `additional_management_users_roles` | `api` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-auth-permission-denied` | `27.58` | `additional_management_auth_permission` | `security` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-request-validation` | `27.48-27.70` | `additional_management_request_validation` | `api` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-secret-redaction` | `27.48-27.70` | `additional_management_secret_redaction` | `security` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-config-state-order` | `22.0d` | `additional_management_config_state_order` | `statefile` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-queue-retention-webhook` | `27.52,27.57,27.64,27.70` | `additional_management_queue_retention_webhook` | `api` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-sdk-transport` | `23.8` | `additional_management_sdk_transport` | `sdk` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-ui-flow` | `24.8` | `additional_management_ui_flow` | `ui` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `additional-management-cache-share-diff` | `27.60-27.70` | `additional_management_cache_share_diff` | `api` | [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) |
| `mcp-cli-lifecycle` | `29.0-29.3` | `mcp_cli_lifecycle` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-jsonrpc-errors` | `29.3` | `mcp_jsonrpc_errors` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-initialize-client-log` | `29.4` | `mcp_initialize_client_log` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-tools-list-call` | `29.5-29.6` | `mcp_tools_list_call` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-tools-confirmation` | `29.6,29.15` | `mcp_tools_confirmation` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-resources-subscription` | `29.7-29.8` | `mcp_resources_subscription` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-prompts` | `29.9` | `mcp_prompts` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-sampling` | `29.10` | `mcp_sampling` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-sse` | `29.12` | `mcp_sse` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |
| `mcp-state-metrics-audit` | `29.13-29.14` | `mcp_state_metrics_audit` | `mcp` | [MCP fixture 固定契約](#mcp-fixture-contract) |

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の fixture は、[§28-F カタログ固定契約](#sec-28-f-2) の同じ行にある fixture 名、節、feature slug を使用し、`owner_component` を `builder` に固定する。`§28` 共通行の `section` は `28-common` とし、`§28.1`〜`§28.25` 行は対応する `28.1`〜`28.25` とする。

<a id="sec-27-f-8"></a>
**[fixture 証跡責務共通 manifest schema 固定契約](fixture.md#sec-27-f-8)：**

[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、[Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) の `manifest.json` は次の共通 schema に従う。未知 key は禁止する。対象別の差分は各 fixture catalog と file set 契約で固定し、別 schema を作成してはならない。

<a id="fixture-component-identifier-contract"></a>
Fixture manifest の component 識別子は `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`release`、`security`、`mcp` の13件だけを許可する。`owner_component`、`collaborator_components`、`components` はこの識別子集合だけを使用する。

```json
{
  "name": "success-example",
  "section": "27.1",
  "feature": "commit_status",
  "category": "success",
  "owner_component": "commitstatus",
  "collaborator_components": ["runner", "statefile"],
  "components": ["commitstatus", "runner", "statefile"],
  "references": [
    "docs/details/commitstatus.md#sec-27-1",
    "docs/details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約",
    "docs/details/statefile.md#sec-22-0a"
  ],
  "fake_clock": "2026-09-16T00:00:00Z",
  "not_applicable": [
    { "path": "input/request.json", "reason": "CLI fixture" }
  ],
  "missing_state": [
    ".build_status.json"
  ],
  "input_files": [],
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
| `name` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) が指す fixture 固定表に記載された fixture 名との完全一致。固定表外の名前、prefix だけが一致する名前、別 fixture 名から推測した名前を禁止する。 |
| `section` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) で `name` に割り当てられた値との完全一致。範囲表記からの推測、主節の任意選択、別節の代用を禁止する。複数節を検証する fixture はレジストリの主節 1 件だけを `section` とし、残りを `references` に記録する。 |
| `feature` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) で `name` に割り当てられた値との完全一致。[fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、owner / 連動 fixture、[Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) はレジストリ記載の snake_case、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) はカタログ記載の kebab-case をそのまま使用し、相互変換、別名、case 変更を禁止する。 |
| `category` | string | 必須 | `success`、`failure`、`partial`、`noop`、`security`。[fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、owner / 連動 fixture は fixture 名 prefix と一致する。[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) と [MCP fixture 固定契約](#mcp-fixture-contract) は各契約の fixture 名別 category 固定表と一致する。 |
| `owner_component` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) で `name` に割り当てられた値との完全一致とし、[component 識別子固定契約](#fixture-component-identifier-contract) のいずれか 1 件を使用する。対象 component の列挙順や fixture 配置から推測してはならない。 |
| `collaborator_components` | array[string] | 必須 | [component 識別子固定契約](#fixture-component-identifier-contract) の値だけを使用する。`owner_component` を含めず、ASCII 昇順、重複なしとする。該当なしは空配列。 |
| `components` | array[string] | 必須 | `owner_component` 1 件と `collaborator_components` の全要素だけを ASCII 昇順、重複なしで含み、[component 識別子固定契約](#fixture-component-identifier-contract) の値だけを使用する。 |
| `references` | array[string] | 必須 | 各要素は repository root 起点の `docs/SPEC.md#<anchor>`、`docs/DESIGN.md#<anchor>`、`docs/DETAIL_INDEX.md#<anchor>`、`docs/details/<file>.md#<anchor>` のいずれかとし、実在 file と実在 anchor を指す。先頭は owner の主節、2 件目は対象の [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、[Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) の fixture カタログ、残りは collaborator と関連固定契約を UTF-8 byte 列の昇順で持つ。`docs/ROADMAP.md`、`docs/DOCUMENT_INDEX.md`、重複、fragment なし、絶対 path、`../` を禁止する。 |
| `fake_clock` | string/null | 必須 | UTC の `YYYY-MM-DDTHH:MM:SSZ` または `null`。カレンダー上無効な日時、offset、小数秒を禁止する。時刻依存 fixture は `null` 禁止。 |
| `not_applicable` | array[object] | 必須 | 各 object は `path` と `reason` の 2 key だけを持つ。`path` は対象 file set 契約の条件付き候補 file または directory の fixture root 相対 path、`reason` は空でない固定理由とする。`path` の ASCII 昇順、重複なしとし、必須 file、実在 path、契約外 path は列挙しない。該当なしは空配列とする。 |
| `missing_state` | array[string] | 必須 | fixture の仮想実行 cwd 起点の `/` 区切り相対 path で、実行前に存在しないことを期待する状態 file だけを ASCII 昇順、重複なしで持つ。directory、symlink、絶対 path、`.` / `..` segment を禁止する。該当なしは空配列とする。 |
| `input_files` | array[string] | 必須 | fixture directory からの `/` 区切り相対 path。対象 file set 契約が許可する `input/` 配下の通常 file だけを ASCII 昇順、重複なしで持ち、非空時は実在 input file set と完全一致させる。input file 不要時は空配列とする。directory、symlink、絶対 path、`.` / `..` segment を禁止する。 |
| `assertions` | array[string] | 必須 | `response`、`request`、`sdk-trace`、`sdk-return`、`sdk-error`、`ui-trace`、`ui-dom`、`stdout`、`stderr`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency`、`no-write` の 1 件以上。重複を禁止し、複数値はこの列挙順で記録する。 |

<a id="fixture-manifest-not-applicable-boundary-contract"></a>
**manifest `not_applicable` 境界固定契約：**

`manifest.json.not_applicable` は、fixture file set 内の条件付き候補 file または directory が当該 fixture で不要であることだけを示す。`manifest.json.not_applicable` を、意味のあるテスト、test gap inventory、test improvement batch closure、requirement coverage、mutation selection、mutation test、race trigger、concurrency / race、contract drift、fixture root coverage、Phase 全体完了の対象外理由として扱ってはならない。

完了判定上の対象外理由は、[test gap inventory record 固定契約](#test-gap-inventory-record-contract)、[test improvement batch closure 固定契約](#test-improvement-batch-closure-contract)、[test requirement coverage ledger 固定契約](#test-requirement-coverage-ledger-contract)、[test verification closure record schema 固定契約](#test-verification-closure-record-schema-contract)、[skip / 未実行証跡固定契約](#test-skip-evidence-contract)、[mutation selection ledger 固定契約](#mutation-selection-ledger-contract)、[mutation test evidence set 固定契約](#mutation-test-evidence-set-contract)、[race trigger matrix 固定契約](#race-trigger-matrix-contract)、または [test evidence package 記録先固定契約](#test-evidence-package-record-location-contract) の該当 record に、対象外範囲、理由、責務正本 anchor、完了可否への影響を記録した場合だけ成立する。

`manifest.json.not_applicable` の `reason` と完了判定 record の対象外理由が矛盾する場合は、fixture 証跡不整合として扱う。`manifest.json.not_applicable` だけが存在し、完了判定 record がない場合は、対象外が成立したものとして扱わない。

<a id="sec-27-f-9"></a>
**[fixture 証跡責務 §27-F 合否判定固定契約](fixture.md#sec-27-f-9)：**

| 判定 | 合格条件 |
|------|----------|
| response | status、headers、body、error details、request id が期待値と一致する。 |
| request | method、path、query、header、body byte、呼出順が期待値と一致する。secret は placeholder だけを許可し、禁止 request、redirect 追従、retry、proxy、未定義 external call が期待値に含まれる場合は不合格とする。 |
| sdk-trace | SDK method、引数、HTTP request、結果種別、status、呼出順が一致し、secret 値と Authorization 値を保持しない。 |
| sdk-return | SDK success return が API response を補完せず一致し、return 後 token state が一致する。 |
| sdk-error | error class、status、message、details、responseBody、error 後 token state が一致する。 |
| ui-trace | UI action、SDK method 呼出し、再取得、disabled 遷移、field 消去の順序が一致する。 |
| ui-dom | panel、text、field error、control state、one-time 領域が構造化期待値と一致する。 |
| stdout / stderr | 改行を含め完全一致する。時刻や path は fake 値だけを使う。 |
| state | 期待対象状態ファイルが byte 等価、または [state diff expected 固定契約](#fixture-state-diff-expected-contract) と一致する。未列挙状態ファイルに差分がない。 |
| logs | JSON Lines は行順、key、値、末尾改行が一致する。破損行 fixture では破損行を修復しない。 |
| effects | 外部 API、外部 command、通知、download、stream の呼び出し回数、順序、payload が一致する。 |
| secret-mask | 禁止文字列が response、stdout/stderr、state、logs、effects、UI DOM に存在しない。 |
| order | 複数状態更新は仕様の保存順と一致する。途中失敗 fixture は失敗地点以降の副作用がない。 |
| idempotency | 同一 fixture を 2 回適用した場合、2 回目の差分が no-op 仕様と一致する。 |
| no-write | `expected/effects.json.unchanged_paths` と `forbidden_writes` に列挙した対象へ差分がない。API fixture では、共通処理で許可された state / log / effect を禁止対象に含めず、期待差分として固定する。 |

対象 fixture 固定表のうち `manifest.json.assertions` に含まれる判定が 1 つでも失敗した場合、その fixture は失敗とする。対象機能の必須 fixture が 1 件でも存在しない、または skip された場合、その機能の実装変更は未完了とする。

<a id="sec-27-f-10"></a>
**[fixture 証跡責務 §27-F assertion 選択固定契約](fixture.md#sec-27-f-10)：**

fixture の `manifest.json.assertions` は、実装者が任意に減らしてはならない。fixture 名 prefix、owner component、collaborator component に応じて、対象 fixture 固定表の assertion を必ず含める。[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定した対象 owner 機能契約で追加検証が必要な場合は、該当する fixture 固定表へ追加してから fixture を作成する。

prefix 行と `components` 行の全条件を合成した集合を `manifest.json.assertions` に記録する。複数 assertion の選択肢がある行は、fixture が実際に生成する結果種別に該当する値を 1 件以上選ぶ。owner 詳細本文が log を生成しない fixture に `logs` を追加したり、expected file を減らすために該当 assertion を省略したりしてはならない。

| 条件 | 必須 assertion |
|------|----------------|
| `success-*` | `effects` と、`response`、`sdk-return`、`ui-dom`、`stdout`、`state`、`logs` のうち対象結果を表す 1 件以上。 |
| `failure-*` | `response`、`sdk-error`、`ui-dom`、`stdout`、`stderr` のうち対象失敗を表す 1 件以上、`effects`、`order`。状態差分なしを期待する場合は `no-write` も必須。 |
| `partial-*` | `state`、`effects`、`order`。log を成功地点まで生成する場合は `logs` も必須。どの副作用まで完了し、どこから未実行かを `expected/effects.json` で明示する。 |
| `noop-*` | `response`、`sdk-return`、`ui-dom`、`stdout` のうち対象結果を表す 1 件以上、`no-write`、`idempotency`、`effects`。外部呼び出し 0 件を `expected/effects.json` に明記する。 |
| `security-*` | `response`、`sdk-return`、`sdk-error`、`ui-dom`、`stdout`、`stderr` のうち対象結果を表す 1 件以上、`secret-mask`、`effects`。 |
| `components` に `api` を含み `input/request.json.http` が non-null | `response`、`state`、`effects`。read-only API は `no-write` も必須。 |
| `components` に `api` を含み HTTP request を実行しない | `state`、`effects`。API service file / process 境界を固定し、`response` と `expected/response.json` を作成しない。 |
| `components` に `admin` を含む | `response` または `stdout` / `stderr` の 1 件以上、`effects`、`secret-mask`。対象 fixture 固定表の必須 expected に `expected/state/` または `expected/state/state-diff.json` がある場合は `state` も必須とし、admin directory の `unchanged_paths` / `updated_paths` / `forbidden_writes` は `expected/effects.json` で併用する。状態比較と副作用境界の一方で他方を代用してはならない。 |
| `components` に `sdk` を含み `input/request.json.sdk` が non-null | `sdk-trace`、`effects` と、`sdk-return` / `sdk-error` の 1 件以上。success と error の両方を実行する fixture は両方を必須とする。`401`、logout、login の token 変化は該当する SDK expected の `token_state` に固定する。 |
| `components` に `ui` を含み `input/request.json.ui` が non-null | `ui-trace`、`ui-dom`、`effects`、`secret-mask`。SDK を実呼出しする fixture は `sdk-trace` と、`sdk-return` / `sdk-error` の該当値も必須とする。 |
| `components` に `runner` を含む | `state`、`logs`、`effects`、`order`。dry-run は `no-write` を必須とする。 |
| `components` に `builder` を含む | `stdout`、`state`、`effects`。[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) では `state` を `expected/site/` と `expected/builder-output.json` の比較に割り当てる。[fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) では `expected/state/` の実ファイル一式または `expected/state/state-diff.json` に出力 file set と byte / metadata 比較を記録する。 |
| `components` に `statefile` を含む | `state`、`effects`、`order`。read-only は `no-write`、JSON Lines を扱う場合は `logs` も必須。 |
| `components` に `archive` を含む | `state`、`effects`、`order`。API response または stream を扱う場合は `response` も必須。 |
| `components` に `commitstatus` を含む | `logs`、`effects`、`secret-mask`、`order`。GitHub Status API 呼び出しと build log 反映を固定する。 |
| `components` に `security` を含む | `effects`、`secret-mask`、`order` と、`response`、`sdk-return`、`sdk-error`、`ui-dom`、`stdout`、`stderr` のうち対象結果を表す 1 件以上。永続状態を作成、更新、削除する場合は `state`、audit / access / notify / server log を生成する場合は `logs` も必須とする。memory-only 値は `expected/security.json` で検証する。 |
| `components` に `setup` を含む | `stdout` または `stderr` の 1 件以上、`state`、`effects`、`secret-mask`、`order`。配置後 file set、mode、systemd 操作、rollback 境界を固定する。 |
| 外部 API / command / 通知を扱う | `effects`、`secret-mask`。呼び出し回数、順序、payload、mask 済み値を必須とする。 |

<a id="sec-27-f-11"></a>
**[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](fixture.md#sec-27-f-11)：**

`expected/effects.json` は次の schema に従う。未知 key は禁止する。

```json
{
  "external_calls": [
    {
      "order": 1,
      "operation": "call",
      "target": "github_api",
      "input": {
        "method": "POST",
        "path": "/repos/{owner}/{repo}/statuses/{sha}",
        "payload": {}
      },
      "output": {},
      "result": "success",
      "error_code": null
    }
  ],
  "commands": [],
  "notifications": [],
  "downloads": [],
  "streams": [],
  "read_api_calls": [
    {
      "order": 1,
      "method": "GET",
      "path": "/api/status",
      "query": {}
    }
  ],
  "unchanged_paths": [],
  "deleted_paths": [],
  "created_paths": [],
  "updated_paths": [],
  "write_order": [],
  "forbidden_created_paths": [],
  "forbidden_updated_paths": [],
  "forbidden_deleted_paths": [],
  "forbidden_writes": [],
  "forbidden_calls": []
}
```

| key | 型 | 必須 | 仕様 |
|-----|----|------|------|
| `external_calls` | array[object] | 必須 | GitHub API、SSH、remote build の process 外呼び出し。通知、download、stream はそれぞれ専用 array にだけ記録する。呼び出しなしは空配列。 |
| `commands` | array[object] | 必須 | pipeline、hook、systemd、archive、setup、update の local command 実行。command 通知は `notifications` にだけ記録する。実行なしは空配列。 |
| `notifications` | array[object] | 必須 | webhook、SMTP、command の各通知送信 attempt。pending 保存は path 副作用として別途記録する。通知なしは空配列。 |
| `downloads` | array[object] | 必須 | Release asset、snapshot、artifact の byte download。byte size、content type、checksum、中断有無は owner 出力契約に従って `output` へ記録する。該当なしは空配列。 |
| `streams` | array[object] | 必須 | SSE / fetch stream の event、close、error。通常の単発 read API と download は記録しない。該当なしは空配列。 |
| `read_api_calls` | array[object] | 必須 | API / SDK / UI fixture が明示的に行う read API call を実行順に記録する。各 object は `order`、`method`、`path`、`query` の 4 key だけを持つ。該当なしは空配列。 |
| `unchanged_paths` | array[string] | 必須 | 実行後に変更があってはならない状態ファイル、出力ファイル、log。 |
| `deleted_paths` | array[string] | 必須 | 実行後に削除される path。削除なしは空配列。 |
| `created_paths` | array[string] | 必須 | 実行後に新規作成される path。作成なしは空配列。 |
| `updated_paths` | array[string] | 必須 | 実行後に既存内容が更新される path。更新なしは空配列。 |
| `write_order` | array[string] | 必須 | 書き込み順。書き込みなしは空配列。複数状態更新 fixture では空配列禁止。 |
| `forbidden_created_paths` | array[string] | 必須 | 作成されてはならない path。作成禁止なしは空配列。 |
| `forbidden_updated_paths` | array[string] | 必須 | 更新されてはならない path。更新禁止なしは空配列。 |
| `forbidden_deleted_paths` | array[string] | 必須 | 削除されてはならない path。削除禁止なしは空配列。 |
| `forbidden_writes` | array[string] | 必須 | 書き込み禁止 path。read-only、dry-run、validation failure fixture では対象状態ファイルを必ず列挙する。 |
| `forbidden_calls` | array[object] | 必須 | 呼び出し禁止の effect 識別子。各 object は `category`、`operation`、`target` の 3 key だけを持つ。呼び出し禁止なしは空配列。 |

`external_calls[]`、`commands[]`、`notifications[]`、`downloads[]`、`streams[]` の各要素は [fake input root 固定契約](#fixture-fake-input-root-contract) の共通 envelope と同じ 7 key だけを持つ。`order` は 5 array をまたぐ実際の effect 開始順とし、1 から始まる重複と欠番のない連続整数とする。各 array 内は `order` 昇順とする。1 つの effect を複数 array へ二重計上してはならない。

| effect category | `operation` 固定値 | `target` 固定対象 |
|-----------------|------------------------|------------------------|
| `external_calls` | `call`、`execute` | 外部 service 識別子、credential なし API path、または SSH / remote build target 識別子。 |
| `commands` | `execute` | 起動する local executable の basename。argv は `input` に記録する。 |
| `notifications` | `send` | `.notify_config.channels[].id` または対象 owner 契約の固定 channel 識別子。URL、mail address、credential は記録しない。 |
| `downloads` | `download` | Release asset 名、snapshot id、artifact id のいずれか。credential 付き URL を記録しない。 |
| `streams` | `stream` | query と credential を含まない API path または対象 build id。 |

`operation`、`input`、`output`、`result`、`error_code` は対象 owner 詳細本文の入出力・異常系契約と一致させる。owner 詳細本文が固定していない値を fixture が独自定義することを禁止する。secret を含む `input` / `output` は平文を保存せず `"***"` を使う。

`read_api_calls[]` の `order` は同 array 内で 1 から始まる連続整数、`method` は `GET` 固定、`path` は query を含まない `/api/` 始まりの絶対 path、`query` は送信する query key と string value だけを持つ object とし、query なしは空 object とする。

`forbidden_calls[]` の `category` は `external_calls`、`commands`、`notifications`、`downloads`、`streams` のいずれか。`operation` と `target` は対応する effect envelope と同じ表記とする。配列は `category`、`operation`、`target` の順の ASCII 昇順、重複なしとする。

`unchanged_paths`、`deleted_paths`、`created_paths`、`updated_paths`、`forbidden_created_paths`、`forbidden_updated_paths`、`forbidden_deleted_paths`、`forbidden_writes` は fixture 実行 root からの `/` 区切り相対 path だけを持ち、絶対 path、空文字、`.`、`..` path segment を禁止する。各配列は ASCII 昇順、重複なしとする。`unchanged_paths`、`deleted_paths`、`created_paths`、`updated_paths` は相互に同一 path を含めない。`created_paths` は `forbidden_created_paths` と `forbidden_writes`、`updated_paths` は `forbidden_updated_paths` と `forbidden_writes`、`deleted_paths` は `forbidden_deleted_paths` と `forbidden_writes` の同一 path を禁止する。`write_order` は同じ path 表記を使い、実際の作成、更新、削除順を保持するため整列せず、同一 path の複数回操作は回数どおり重複記録する。`created_paths`、`updated_paths`、`deleted_paths` の各 path は `write_order` に 1 回以上存在し、`write_order` の path はこの 3 array のいずれかに存在しなければならない。並列処理 fixture で完了順が非決定の場合でも、仕様が要求する保存順を `write_order` に固定する。

<a id="sec-27-f-11-security"></a>
**[fixture 証跡責務 §27-F expected/security.json schema 固定契約](fixture.md#sec-27-f-11-security)：**

`expected/security.json` は次の root schema に従い、未知 root key と各 array object の未知 key を禁止する。secret の実 byte は `input/` 内の synthetic fixture 値から実行時に参照し、`expected/` 配下へ複製してはならない。

```json
{
  "forbidden_plaintexts": [
    {
      "id": "api_token",
      "source_file": "input/fakes.json",
      "source_pointer": "/entropy/0/output/value"
    }
  ],
  "forbidden_headers": ["authorization", "cookie", "set-cookie"],
  "forbidden_state_values": [
    {
      "state_path": ".api_tokens",
      "json_pointer": "/tokens/0/token",
      "source_id": "api_token"
    }
  ],
  "allowed_one_time_response_fields": [
    {
      "source_id": "api_token",
      "response_file": "expected/response.json",
      "json_pointer": "/body/token",
      "max_occurrences": 1
    }
  ],
  "hash_only_fields": [
    {
      "state_path": ".api_tokens",
      "json_pointer": "/tokens/0/token_hash",
      "source_id": "api_token"
    }
  ],
  "memory_only_values": [
    "login_failure_count",
    "login_ticket",
    "session",
    "totp_setup_secret"
  ],
  "one_time_response_assertions": [
    {
      "source_id": "api_token",
      "response_file": "expected/response.json",
      "json_pointer": "/body/token",
      "occurrences": 1
    }
  ],
  "scope_decisions": [
    {
      "method": "GET",
      "path": "/api/status",
      "token_scopes": ["read"],
      "expected_status": 200
    }
  ],
  "rate_limit_decisions": [
    {
      "group": "read",
      "key": "ip:127.0.0.1",
      "attempt": 1,
      "expected_status": 200,
      "expected_count": 1
    }
  ],
  "audit_required": [
    {
      "action": "token_create",
      "result": "success",
      "expected_records": 1
    }
  ]
}
```

| root key | 型 | 固定条件 |
|----------|----|----------|
| `forbidden_plaintexts` | array[object] | 各 object は `id`、`source_file`、`source_pointer` の 3 key だけを持つ。`id` は同一 file 内で一意な lowercase snake_case、`source_file` は fixture root からの `input/` 配下相対 path、`source_pointer` は RFC 6901 JSON Pointer とする。参照値は非空 string とし、実 credential を使用しない。 |
| `forbidden_headers` | array[string] | lowercase header 名を ASCII 昇順、重複なしで列挙する。値ではなく header 名を指定する。 |
| `forbidden_state_values` | array[object] | 各 object は `state_path`、`json_pointer`、`source_id` の 3 key だけを持つ。指定 state field に `source_id` の平文が存在しないことを確認する。 |
| `allowed_one_time_response_fields` | array[object] | 各 object は `source_id`、`response_file`、`json_pointer`、`max_occurrences` の 4 key だけを持つ。`response_file` は `expected/` 配下、`max_occurrences` は `1` 固定とする。ここにない response field へ secret を出現させてはならない。 |
| `hash_only_fields` | array[object] | 各 object は `state_path`、`json_pointer`、`source_id` の 3 key だけを持つ。field は source 平文と不一致であり、対象 owner 詳細本文が定める hash verifier で一致しなければならない。 |
| `memory_only_values` | array[string] | 許容値は `session`、`login_ticket`、`totp_setup_secret`、`login_failure_count`。ASCII 昇順、重複なしとし、restart 後は不在を期待する。 |
| `one_time_response_assertions` | array[object] | 各 object は `source_id`、`response_file`、`json_pointer`、`occurrences` の 4 key だけを持つ。`occurrences` は `1` 固定とし、同じ `source_id` の `allowed_one_time_response_fields` と完全一致する response location を指定する。 |
| `scope_decisions` | array[object] | 各 object は `method`、`path`、`token_scopes`、`expected_status` の 4 key だけを持つ。`token_scopes` は [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42) の固定 scope だけを ASCII 昇順、重複なしで持つ。 |
| `rate_limit_decisions` | array[object] | 各 object は `group`、`key`、`attempt`、`expected_status`、`expected_count` の 5 key だけを持つ。`attempt` は 1 以上、同一 group / key 内で昇順とする。 |
| `audit_required` | array[object] | 各 object は `action`、`result`、`expected_records` の 3 key だけを持つ。`expected_records` は 0 以上とし、fixture が要求する audit record 数を固定する。 |

`forbidden_plaintexts[].id` を参照する `source_id` は同じ file の `forbidden_plaintexts` に存在しなければならない。`expected/response.json` の one-time field には平文の代わりに string `"${secret:<source_id>}"` を置き、fixture runner は比較時だけ source を解決する。解決後の値を log、diff、error、更新済み expected file へ出力してはならない。`allowed_one_time_response_fields` と `one_time_response_assertions` に同じ location がない secret は response へ出現禁止とする。該当しない root key も空配列で残す。

<a id="sec-27-f-12"></a>
**[fixture 証跡責務 §27-F 不足時 未完了判定固定契約](fixture.md#sec-27-f-12)：**

| 不足 | 未完了理由 |
|------|------------|
| fixture 名が対象の [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) または [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) のカタログ固定表または owner fixture 固定表に存在しない。 | 固定表外 fixture のため未完了。 |
| 必須 fixture が存在しない。 | 機能の正常 / 異常 / no-op / security / partial coverage 不足。 |
| `manifest.json.assertions` が [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) assertion 選択固定契約を満たさない。 | 合否判定不足。 |
| `expected/effects.json` が存在しない、または必須 key が欠落する。 | 副作用検証不足。 |
| `unchanged_paths` または `forbidden_writes` が空で、fixture が failure / noop / read-only / dry-run / validation error のいずれかである。 | 無変更保証不足。 |
| 外部 API / command / notification を扱う fixture で `forbidden_calls` が空、かつ禁止副作用なしの理由が `manifest.json.not_applicable` にない。 | 外部副作用境界不足。 |
| `manifest.json.assertions` に `secret-mask` がある fixture で `expected/security.json` が存在しない。 | secret mask 検証不足。 |
| JSON Lines を扱う fixture で末尾改行、行順、破損行保持/除外条件を期待値に含めていない。 | log / audit / history 検証不足。 |
| idempotency fixture で 1 回目と 2 回目の期待差分を分離していない。 | 再実行検証不足。 |
| partial fixture で失敗地点より後の `forbidden_writes` / `forbidden_calls` を列挙していない。 | 部分失敗境界不足。 |

[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の不足は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。fixture が多くても、期待副作用、禁止副作用、secret mask、保存順、再実行差分が明示されていなければ、バグ修正ゼロ化の検証を満たさない。

<a id="sec-27-f-13"></a>
**[fixture 証跡責務 §27-F 相互整合固定契約](fixture.md#sec-27-f-13)：**

fixture 内の `manifest.json`、`input/*`、`expected/*` は相互に矛盾してはならない。[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の不整合が 1 件でもある場合、その fixture は失敗とし、対象機能の実装変更は未完了とする。

| 整合対象 | 固定条件 | 不整合時の扱い |
|----------|----------|----------------|
| `manifest.json.name` と directory 名 | directory 名は `manifest.json.name` と完全一致する。 | fixture 名不一致として失敗。 |
| `manifest.json.category` と fixture 固定表 | `success-*` は `success`、`failure-*` は `failure`、`partial-*` は `partial`、`noop-*` は `noop`、`security-*` は `security` とする。prefix を持たない [追加管理 API fixture 固定契約](#additional-management-api-fixture-contract) と [MCP fixture 固定契約](#mcp-fixture-contract) の fixture は、各契約の fixture 名別 category 固定表と完全一致させる。 | 分類不一致として失敗。 |
| `manifest.json.section` と fixture 固定表 | section、fixture 名、owner の組み合わせは対象の [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約)、[Admin CLI fixture 固定契約](#admin-cli-fixture-contract)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) のカタログ固定表または owner fixture 固定表に存在する組み合わせだけ許可する。 | 固定表外として失敗。 |
| `manifest.json.owner_component` と `components` | `owner_component` は `components` に必ず含める。 | owner 責務不一致として失敗。 |
| `manifest.json.collaborator_components` と `components` | `collaborator_components` は `components` にすべて含め、`owner_component` を含めてはならない。 | collaborator 責務不一致として失敗。 |
| fixture 固定表の必須 input と入力ファイル | CLI / setup script 実行は `input/cli.json`、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の builder は `input/options.json`、HTTP request / SDK invocation / UI action は `input/request.json` を置く。その他の builder fixture は対象 file set 契約に従う。複数条件に該当する場合は必要 file をすべて置き、該当する条件を `not_applicable` で除外してはならない。 | 入力責務不一致として失敗。 |
| `manifest.json.assertions` と期待値ファイル | `response` は `expected/response.json`、`request` は `expected/request.json`、`sdk-trace` は `expected/sdk_trace.json`、`sdk-return` は `expected/sdk_return.json`、`sdk-error` は `expected/sdk_error.json`、`ui-trace` は `expected/ui_trace.json`、`ui-dom` は `expected/ui_dom.json`、`stdout` は `expected/stdout.txt`、`stderr` は `expected/stderr.txt`、`logs` は `expected/logs/`、`effects` は `expected/effects.json`、`secret-mask` は `expected/security.json` を要求する。`state` は [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) では `expected/state/`、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) では `expected/site/` と `expected/builder-output.json` を要求する。 | 期待値不足として失敗。 |
| `expected/effects.json.write_order` と expected file set | `created_paths` / `updated_paths` として `write_order` に現れる path は [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[追加管理 API fixture 固定契約](#additional-management-api-fixture-contract)、[MCP fixture 固定契約](#mcp-fixture-contract) の `expected/state/` / `expected/logs/`、または [fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `expected/site/` / `expected/builder-output.json` で最終値を固定する。`deleted_paths` は対応する expected tree に存在させず、削除期待を `expected/effects.json` に固定する。 | 保存順だけの空検証として失敗。 |
| `expected/effects.json` の実行後 path 集合 | `unchanged_paths`、`created_paths`、`updated_paths`、`deleted_paths` は pairwise disjoint とする。 | 副作用分類矛盾として失敗。 |
| `expected/effects.json` の実行期待と禁止集合 | `created_paths` と `forbidden_created_paths` / `forbidden_writes`、`updated_paths` と `forbidden_updated_paths` / `forbidden_writes`、`deleted_paths` と `forbidden_deleted_paths` / `forbidden_writes` に同一 path を含めない。 | 禁止副作用矛盾として失敗。 |
| `expected/effects.json.write_order` と path 集合 | `write_order` に現れる各 path は `created_paths`、`updated_paths`、`deleted_paths` のいずれかに存在し、この 3 array の全 path は `write_order` に 1 回以上存在する。`unchanged_paths` と `forbidden_writes` の path を `write_order` に含めない。 | 保存順または禁止書込矛盾として失敗。 |
| `expected/effects.json.forbidden_calls` と effect 5 array | `category`、`operation`、`target` が一致する禁止 effect を実行期待に含めない。 | 禁止呼び出し矛盾として失敗。 |
| `expected/security.json` と期待値全体 | `expected/security.json.forbidden_plaintexts` が参照する source 値は、`allowed_one_time_response_fields` と `one_time_response_assertions` の同一 location 以外の response、stdout、stderr、state、logs、effects、UI DOM に出現してはならない。 | secret mask 不足として失敗。 |
| `input/fakes.json` と `expected/effects.json` | process 外副作用を表す fake は、本節の分類に対応する effect array に同数、同じ `operation` / `target` / `result` で記録する。fake root 内順序と対応 effect の相対順序を一致させる。未消費 fake がある場合は `manifest.json.not_applicable` に理由を置く。 | fake / effect 不一致として失敗。 |
| `missing_state` と `input/state/` | `manifest.json.missing_state` に列挙した path は `input/state/` に存在してはならない。 | 初期状態矛盾として失敗。 |

<a id="sec-27-f-14"></a>
**[fixture 証跡責務 §27-F component 別検証責務固定契約](fixture.md#sec-27-f-14)：**

実装変更は、対象 component ごとに [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の責務を満たす。複数 component を含む機能では、owner component と collaborator component を fixture manifest に分けて記録する。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

| component | owner 時の必須検証責務 | collaborator 時の必須検証責務 | 完了判定 | 禁止越境 |
|-----------|------------------------|-------------------------------|----------|----------|
| `builder` | CLI 入力、Markdown 入力、出力 site / HTML / REPORT、asset、cache、dependency、meta、終了コードを fixture で固定する。 | runner / API に渡す REPORT、meta、dependency manifest の key 名と nullable 条件を固定する。 | 出力 file の byte 比較または構造化 expected が存在し、既存出力保護と失敗時 no-write が検証済み。 | GitHub read、状態ファイル直接更新、通知送信。 |
| `runner` | CLI、設定、lock、SHA cache、build log、history、status、queue、external call、notification、終了コードを fixture で固定する。 | builder output、[`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の schema、commitstatus payload を仕様どおり消費し、未定義 key や未取得値を追加しない。 | 成功、失敗、skip、partial、dry-run の状態差分と write order が検証済み。 | API endpoint 追加、SDK method 追加、UI 操作追加。 |
| `api` | method/path/query/body/header、auth/scope/rate limit、response、状態 read/write、access/audit/config log を fixture で固定する。 | SDK / UI が追加 key を生成せずに扱える response schema、HTTP status、error body を返す。 | read-only は endpoint 固有の業務状態が no-write で、共通 security / observability 副作用、write API の保存順、validation failure の forbidden_writes が検証済み。 | runner CLI 処理、UI DOM 操作、未定義状態ファイル作成。 |
| `admin` | admin archive の file list、entry validation、配置 mode、static serving path/header/body、CLI request/stdout/stderr/exit code、secret 非配信、no mutation を fixture で固定する。 | setup の admin UI 展開、api の static serving / CLI fake response、security の secret mask、ui / sdk の配布物境界を壊さない expected を提供する。 | unsafe archive、未定義 file、directory listing、secret path、method denied、CLI 未定義 endpoint passthrough 禁止、no mutation が検証済み。 | UI / SDK 内容生成、API endpoint 実装、状態 schema 変更、systemd 操作、CLI からの任意 endpoint 実行。 |
| `sdk` | method、args、query 生成、error 変換、token 破棄、binary / stream handling を fixture で固定する。 | UI が API 詳細を知らずに扱える戻り値と error をそのまま伝播する。 | API response に存在しない key を生成せず、`401` / `403` / `429` / network error が区別される。 | API response 推測補完、未定義 endpoint 呼び出し、状態ファイル直接操作。 |
| `ui` | SDK method 呼び出し、DOM 表示、disabled/loading/error、secret field 消去、再取得順を fixture で固定する。 | SDK 戻り値だけを表示し、API / 状態ファイルの内部構造を再解釈しない。 | 直接 API 呼び出し、状態ファイル操作、secret DOM 残存がない。 | 直接 `fetch()`、状態ファイル操作、外部 command 実行。 |
| `statefile` | schema、atomic write、JSON Lines、lock、破損時処理、保存順、no-write / forbidden write を fixture で固定する。 | runner / API / archive の保存対象ごとに `write_order`、`unchanged_paths`、`forbidden_writes` を固定する。 | read-only、dry-run、validation failure、partial failure の副作用境界が検証済み。 | component 固有の業務判断、UI 表示判断、API response 補完。 |
| `archive` | log archive、snapshot、download、delete、rollback の保存 / 取得 / 削除境界を fixture で固定する。 | API download / rollback response と runner history への影響を [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の契約に合わせる。 | 元 log 保持、archive 破損時挙動、安全でない entry 拒否、rollback 失敗時 no-write が検証済み。 | build 成否の反転、元 build log の改変、未検証 archive 展開。 |
| `commitstatus` | GitHub Commit Status payload、pending / final 送信順、失敗時非反転、secret mask を fixture で固定する。 | runner の build 結果を受け取り、build 成否を変更せず送信結果だけを返す。 | pending 失敗、final 失敗、commit SHA なし、無効時呼び出し 0 が検証済み。 | build 実行判断、dry-run での GitHub write、status 失敗による build 成否反転。 |
| `security` | token hash、session、TOTP、scope、audit、rate limit、secret mask、forbidden call/write を fixture で固定する。 | API / SDK / UI / runner の secret 表示、認証失敗、副作用境界を検証する。 | 認証失敗、権限拒否、rate limit、audit failure の副作用境界が検証済み。 | 業務処理代行、認可前状態更新、secret 平文保存。 |
| `setup` | binary 配置、service 更新、rollback、secret 既存値保持、stdout/stderr mask、終了コードを fixture で固定する。 | runner / API の初期状態と既存 secret を壊さないことを effects で固定する。 | 部分失敗時の復元対象と復元禁止副作用が `expected/effects.json` に明記済み。 | runtime 機能追加、状態 schema 暗黙変更、外部依存追加。 |
| `release` | CLI、Git事前条件、build引数、asset byte、archive metadata、checksum、GitHub request順、draft cleanup、終了コードをfixtureで固定する。 | setupが受け入れるasset名、version、checksum、admin archiveを追加判断なしで提供する。 | 2回buildのbyte一致、upload後再取得digest、正式公開前全検証、失敗時cleanup、token maskが検証済み。 | tag変更、checkout変更、既存Release上書き、setup配置処理、secret出力。 |

component 責務を複数変更へ分ける場合でも、各変更が満たすべき owner component、collaborator component、fixture 名、期待ファイル、禁止副作用を実装検証証跡に明記する。責務の所在が不明な場合は、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

<a id="sec-27-f-15"></a>
**[fixture 証跡責務 §27-F runner / statefile 連動 fixture 固定契約](fixture.md#sec-27-f-15)：**

[§27.21〜§27.38](runner.md#sec-27-21) の runner owner 機能は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の対象機能行に定めた fixture に加えて、次の選択条件に該当する連動 fixture をすべて作成する。状態作成、更新、削除、複数保存順の変更は `success-runner-state-write-order`、途中の write / append / fsync failure を扱う変更は `partial-runner-state-write-failure`、同一入力または no-op の状態差分を扱う変更は `noop-runner-state-idempotency`、dry-run を扱う変更は `noop-runner-state-dry-run`、破損状態の read / recovery / stop を扱う変更は `failure-runner-state-corrupt-boundary`、secret を入力または状態に含む変更は `security-runner-state-secret-mask` を必須とする。複数条件に該当する場合は該当 fixture を省略せず、runner の業務判断と statefile の保存境界を分離して検証する。

| runner/statefile fixture 名 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|-----------------------------|----------------|------------|---------------|----------|
| `success-runner-state-write-order` | `runner`、`statefile` | build lifecycle、queue、history、status、log、対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 状態。 | `expected/effects.json.write_order`、`expected/state/`、`expected/logs/`。 | 対象の [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 機能契約の保存順と一致し、並列処理でも永続保存順が固定される。 |
| `partial-runner-state-write-failure` | `runner`、`statefile` | N 番目の state write / JSON Lines append / fsync fake failure。 | `expected/state/`、`expected/logs/`、`expected/effects.json` の `updated_paths`、`unchanged_paths`、`forbidden_writes`、`write_order`。 | 失敗地点前の成功済み状態は保持し、失敗地点以降は変更しない。未定義 rollback を行わない。 |
| `noop-runner-state-idempotency` | `runner`、`statefile` | 同一入力の 1 回目 / 2 回目、disabled、skip、duplicate、sample 不足。 | `expected/state/`、`expected/logs/`、`expected/effects.json` の 2 回目 `unchanged_paths` と空の実行 effect。 | 2 回目または no-op で不要な log / history / status / notify / audit 差分を作らない。 |
| `noop-runner-state-dry-run` | `runner`、`statefile` | dry-run option、変更あり target、外部 call fake。 | `expected/stdout.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json` の `forbidden_writes`、`forbidden_calls`、空の実行 effect。 | dry-run は状態ファイル、lock、external write、notification を一切変更しない。 |
| `failure-runner-state-corrupt-boundary` | `runner`、`statefile` | 破損 `.build_state`、`.build_status.json`、`.build_history`、対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) 状態、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json`。 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の退避 / 再生成 / 停止条件に従い、破損内容を stdout / stderr / log / expected に出さない。 |
| `security-runner-state-secret-mask` | `runner`、`statefile` | PAT、branch env secret、hook output secret、remote credential、notification secret。 | `expected/security.json`、`expected/logs/`、`expected/effects.json`。 | secret 平文、prefix、suffix、長さ、hash が stdout / stderr / state / log / fixture expected に残らない。 |

<a id="sec-27-f-16"></a>
**[fixture 証跡責務 §27-F statefile owner fixture 固定契約](fixture.md#sec-27-f-16)：**

[`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a)〜[§22.0s](statefile.md#sec-22-0s) の fixture は、状態ファイルの schema、atomic write、lock、JSON Lines、破損時処理、保存順、read-only no mutation を固定する。各 fixture の `manifest.json.name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) に従い、`expected/effects.json` に file list、mtime、mode、content、write order、forbidden writes を記録する。

| fixture | 初期状態 | 操作 | 合格条件 |
|---------|----------|------|----------|
| `success-state-read-missing` | target 不在 | 対応する read adapter 呼び出し | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) の不在時戻り値を返し、filesystem 差分なし。 |
| `success-state-text-payload` | `.notes`、`.webhook_secret`、`.smtp_secret` の各 target 不在または schema-valid な既存 byte 列 | [`docs/details/statefile.md` 詳細本文責務 UTF-8 text payload 固定契約](statefile.md#statefile-text-payload-contract) の path 別 valid payload を write adapter へ渡す | 保存後 byte 列が入力と完全一致し、暗黙の trim、Unicode normalization、末尾 LF の追加または削除がなく、mode `0600`、tmp / lock 残存なし。 |
| `failure-state-text-payload` | `.notes`、`.webhook_secret`、`.smtp_secret` の schema-valid な既存 byte 列 | 同固定契約の path 別拒否条件を 1 条件ずつ write adapter へ渡す | validation failure を返し、target content / mode / mtime、tmp、lock に差分なし。secret 平文を stdout、stderr、log、error、fixture expected へ出力しない。 |
| `failure-state-corrupt-object` | JSON parse 不能または未知 key あり | read adapter 呼び出し | `ErrStateCorrupted`。target 差分なし。API の公開応答は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。 |
| `success-state-corrupt-regenerates` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) で再生成指定済み file が破損 | write caller または再生成を伴う操作 | corrupt backup が 1 件作成され、初期値だけが保存される。 |
| `failure-state-lock-timeout` | `{name}.lock` が残る case、elapsed `9900ms` の試行前に解放する case、elapsed `10000ms` と同時に解放する case で、timer target `state_lock_wait` を `100ms` ずつ進める | write caller 呼び出し | elapsed `0` と `100ms`〜`9900ms` の合計 100 回だけ取得を試行する。`9900ms` 解放は取得成功、`10000ms` 同時解放は追加取得なしの conflict failure。timeout case は target / tmp 差分なし。 |
| `failure-state-create-only-existing` | create-only target として通常 file、directory、symlink、その他の file type を個別に配置 | create-only write caller 呼び出し | target 内容を読まず `ErrStateAlreadyExists`、tmp 作成なし、target 差分なし、取得した lock だけを 1 回削除する。lock 削除失敗 subcase は `STATE_LOCK_CLEANUP_FAILED`、残存 lock、cleanup failure を固定する。 |
| `success-state-locked-update` | schema-valid target と事前確定済み業務入力 | locked update adapter で更新値と result value を返す callback を 1 回実行 | lock 内で最新値を 1 回読み、callback 1 回、schema 検証、atomic write、lock 解放の順に実行する。更新後 target と result value が一致し、tmp / lock を残さない。 |
| `noop-state-locked-update` | schema-valid target | locked update adapter で `no-op` と result value を返す callback を 1 回実行 | current value を 1 回読み、callback 1 回、lock 解放だけを実行する。result value と `nil` error を返し、tmp、rename、target content / mode / mtime 差分を発生させない。 |
| `failure-state-mutation-callback` | schema-valid target | locked update adapter の callback が sentinel error を返す | callback error の `errors.Is` identity を維持し、tmp、rename、target 差分なし、lock 解放 1 回とする。lock 解放失敗 subcase は最初の callback error、`STATE_LOCK_CLEANUP_FAILED`、残存 lock を固定する。 |
| `failure-state-mutation-panic` | schema-valid target と secret を含む synthetic current value | locked update adapter の callback が synthetic panic を発生させる | process を panic させず `ErrStateMutationPanic` を返す。tmp、rename、target 差分なし、lock 解放 1 回、`STATE_MUTATION_PANIC` / `ERROR` 1 件とし、panic value、stack、path、current value、secret を出力しない。 |
| `failure-state-write-before-rename` | tmp create / full write / chmod / close のいずれかを fake failure | write caller 呼び出し | rename を実行せず旧 target を byte 単位で維持する。当該呼び出しが作成した tmp と取得した lock だけを tmp → lock の順で各 1 回削除し、write failure を返す。 |
| `failure-state-chmod` | rename 前の chmod を fake failure | write caller 呼び出し | 成功扱いにせず、旧 target を byte 単位で維持し、tmp → lock の順で各 1 回削除して chmod failure を返す。 |
| `failure-state-tmp-fsync` | rename 前の tmp file `Sync` を fake failure | write caller 呼び出し | rename を実行せず旧 target を byte 単位で維持し、tmp → lock の順で各 1 回削除して write failure を返す。 |
| `partial-state-tmp-cleanup` | rename 前 write failure 後の tmp `os.Remove` だけを fake failure | write caller 呼び出し | tmp 削除を再試行せず、lock は 1 回削除する。残存 tmp、`STATE_TMP_CLEANUP_FAILED` と basename だけの ERROR log、最初の write failure の返却を固定する。 |
| `partial-state-lock-cleanup-before-rename` | rename 前 write failure 後の lock `os.Remove` だけを fake failure | write caller 呼び出し | lock 削除を再試行せず、残存 lock、`STATE_LOCK_CLEANUP_FAILED` と basename だけの ERROR log、最初の write failure の返却を固定する。 |
| `partial-state-parent-fsync` | rename 成功後の parent directory `Sync` を fake failure | write caller 呼び出し | 新 target を維持し、post-rename partial write failure を返し、server log に `STATE_WRITE_AFTER_RENAME_FAILED`、secret 非表示を固定する。 |
| `partial-state-lock-remove` | rename / sync 成功後の lock 削除を fake failure | write caller 呼び出し | 新 target と残存 lock を確認し、post-rename partial write failure と `STATE_WRITE_AFTER_RENAME_FAILED` を固定する。 |
| `partial-json-lines-corrupt` | 有効行と破損行が混在 | list caller 呼び出し | 有効行だけ返し、server log に line number、呼び出し元の公開値に破損詳細なし。 |
| `noop-state-read-no-mutation` | 破損なし state 一式 | 全 read-only caller 呼び出し | state dir の file list、mtime、mode、content が変化しない。 |
| `partial-state-multi-write` | 2 file 目の rename 前 write を fake failure | 複数ファイル更新 caller 呼び出し | 1 file 目は保持、2 file 目以降は未変更。statefile owner は caller 固有 log を追記せず、caller 固有 fixture が詳細本文責務どおりの log 副作用を別途固定する。 |

statefile owner fixture が不足する場合、`statefile` は詳細実装確認を満たした扱いにしてはならない。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

<a id="sec-27-f-17"></a>
**[fixture 証跡責務 §27-F API / SDK / UI 連動 fixture 固定契約](fixture.md#sec-27-f-17)：**

[`docs/details/runner.md` 詳細本文責務 §27.21](runner.md#sec-27-21)〜[§27.38](runner.md#sec-27-38) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) のうち API、SDK、UI が連動する実装変更は、対象機能の owner fixture に加えて次の選択条件に該当する連動 fixture をすべて作成する。正常な UI → SDK → API request を追加または変更する場合は `success-api-sdk-ui-request-trace`、HTTP / SDK / UI error 伝播を追加または変更する場合は `failure-api-sdk-ui-error-propagation`、変更成功後の再取得を追加または変更する場合は `partial-api-sdk-ui-refresh-order`、one-time secret の発行または消去を扱う場合は `security-api-sdk-ui-secret-one-time`、不足 key、未知値、破損行除外済み値を扱う場合は `security-api-sdk-ui-no-speculation`、validation / authorization / rate limit / no-op / partial failure の副作用境界を扱う場合は `security-api-sdk-ui-side-effect-boundary` を必須とする。複数条件に該当する場合は該当 fixture を省略せず、owner component の本文を置き換えずに API response、SDK method、UI 表示の接続点を固定する。

| api/sdk/ui fixture 名 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|------------------------|----------------|------------|---------------|----------|
| `success-api-sdk-ui-request-trace` | `api`、`sdk`、`ui` | UI user action、same-origin の属性不在 / 空文字 / absolute origin override、空白と予約文字を含む query、SDK fake fetch trace、API request fixture。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | 最初の login URL を含む UI → SDK → API の method / path / query / body が [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) と一致する。`/api/api/` がなく、query 値は `encodeURIComponent` の 1 回適用と同じで空白を `%20` とする。UI → SDK → API の repository 内 call を `external_calls` へ記録せず、API 自身が外部 service を呼ばない場合は `external_calls` を空配列とする。 |
| `failure-api-sdk-ui-error-propagation` | `api`、`sdk`、`ui` | `401`、`403`、`409`、`422 details`、`429`、`500`、SDK header 待機 / JSON body / Blob body の timeout、timeout 以外の body read reject、binary 成功の不正 media type、SDK path `id` の空文字、`.`、`..`、`/`、`U+005C BACKSLASH`、空白、非 ASCII、65 文字以上、UI の外部 origin / path 付き `data-api-base-url` の fake response。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_error.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | HTTP error の status、message、details、token 破棄条件、panel error、field error、disabled が固定どおり。SDK timeout は body 読取完了まで有効で、その他の body read reject は `Network error`、binary media type 不一致は `Invalid binary response` とする。不正 `id` は `TypeError("Invalid argument: id")` として HTTP 送信を 0 回とする。不正 `baseUrl` は SDK 生成と HTTP 送信を 0 回とし、`UI initialization failed` だけを表示する。 |
| `partial-api-sdk-ui-refresh-order` | `api`、`sdk`、`ui` | 変更 API 成功、成功後再取得 1 件目成功、2 件目失敗。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/sdk_error.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/effects.json` の `read_api_calls`、`expected/security.json`。 | 再取得順を守り、変更成功は維持し、再取得失敗だけ panel error に表示する。変更 API を再送しない。 |
| `security-api-sdk-ui-secret-one-time` | `security`、`sdk`、`ui` | token 発行、TOTP setup、secret 保存失敗、trusted / synthetic event、clipboard success / failure / API 不在、新旧 generation、panel 遷移、logout、`401`。 | `expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/sdk_error.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/logs/`、`expected/security.json`、`expected/effects.json`。 | token / TOTP secret / otpauth URI は [UI one-time secret 消去契約](ui.md#ui-one-time-secret-contract) の 3 専用領域にだけ表示する。表示開始 event では消去せず、有効化後の trusted `click` / `submit` / `input` / `change` / `keydown` / `copy` をそれぞれ別 case で固定する。synthetic event と除外 callback では保持、copy は 1 回呼出し後の success / failure どちらでも消去、panel 遷移 / logout / `401` / revoke all は即時消去、旧 generation callback は新値を消去しない。消去後に DOM / SDK property / log / clipboard 再呼出しへ平文が残らない。 |
| `security-api-sdk-ui-no-speculation` | `api`、`sdk`、`ui` | 不足 key、未知 widget、unknown category、破損行除外済み response。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/security.json`、`expected/effects.json`。 | SDK は key を補完せず、UI は API 値だけ表示し、未知値は固定 error / warning / 空状態で扱う。 |
| `security-api-sdk-ui-side-effect-boundary` | `api`、`sdk`、`ui`、`statefile` | validation failure、認可失敗、rate limit、no-op、partial failure。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_error.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/effects.json`、`expected/state/`、`expected/logs/`、`expected/security.json`。 | 禁止 write / call が 0 件で、保存済み主状態、config log、audit、notify、UI 表示が [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定した対象 owner 機能契約の部分失敗契約と一致する。 |

SDK 連動 fixture の expected は [SDK expected 固定契約](#fixture-sdk-expected-contract)、UI 連動 fixture の expected は [UI expected 固定契約](#fixture-ui-expected-contract) に従う。SDK / UI fixture 証跡は [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) と [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の実装契約を検証する証跡であり、SDK method、戻り値、error class、token 破棄条件、DOM、表示順、disabled 条件、secret 消去条件を再定義しない。

<a id="sec-27-f-18"></a>
**[fixture 証跡責務 §27-F UI owner fixture 固定契約](fixture.md#sec-27-f-18)：**

[`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) の UI fixture は、SDK method 呼び出し、DOM 表示、disabled / loading / error、secret field 消去、再取得順、one-time 表示、no speculative state を固定する。各 fixture の `manifest.json.name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) に従い、対象 `§27.x` は `references` に記録する。`ui-trace` と `ui-dom` assertion、および対応する `expected/ui_trace.json` と `expected/ui_dom.json` を必須にする。API response と fake SDK 入力は `input/fakes.json`、実行後 SDK call は `expected/sdk_trace.json`、SDK result は `expected/sdk_return.json` / `expected/sdk_error.json` に分離し、expected file を入力として使用してはならない。

| fixture | 入力 / fake SDK 入力 | 合格条件 |
|---------|----------------------|----------|
| `success-ui-login-totp` | `login()` が exact `{totp_required:true,ticket}` を返す。 | password 消去、TOTP field 表示、ticket は DOM に表示せず、ログイン成功後の初期取得は 0 件。`must_change` を UI が推測しない。 |
| `success-ui-login-prompt` | `login()` または `loginTotp()` が `must_change:"prompt"` を返す。 | 初期取得 5 method を固定順で呼び、password panel と通常 panel を表示する。password 変更成功後は初期取得を再実行しない。 |
| `success-ui-forced-password` | `login()` または `loginTotp()` が `must_change:"forced"` を返す。 | session 発行直後の初期取得は 0 件、password panel 以外が操作不可。変更成功後に `getDashboard()` → `getStatus()` → `getQueue()` → `getMaintenance()` → `getDashboardLayout()` を 1 回ずつ呼び、その後に通常 panel を表示する。 |
| `partial-ui-refresh-failure` | 変更 API 成功後の再取得 2 件目が失敗。 | 変更成功は維持し、再取得失敗だけ panel error に表示する。 |
| `noop-ui-destructive-cancel` | 確認 dialog cancel。 | SDK method 呼び出し 0 回、表示差分なし。 |
| `security-ui-secret-clearing` | token 発行、TOTP setup、PAT 更新、Webhook secret 保存、[UI one-time secret 消去契約](ui.md#ui-one-time-secret-contract) の全消去 / 除外条件。 | 通常 secret field は成功、失敗、遷移で消去し、one-time 3 領域は trusted event、copy 完了、即時消去、generation guard を固定する。synthetic event と除外 callback では消去しない。 |
| `failure-ui-runtime-initial-status-error` | `getStatus()` が `AdlaireCIError(status=500,message="State file is corrupted")` を返す。 | status panel error に固定 message を表示し、build button を成功扱いにしない。 |
| `failure-ui-runtime-manual-build-conflict` | `triggerBuild()` が `409 Conflict` を返す。 | error 表示、`getStatus()` と `getQueue()` をこの順で再取得、同じ build request を再送しない。 |
| `failure-ui-runtime-queue-full` | `triggerBuild()` が `429 queue_full` を返す。 | build button を 10 秒 disabled、password や secret field は変更しない。 |
| `success-ui-runtime-build-dispatch` | `triggerBuild()` が queue id と `dispatch:"requested"` を返す。 | `Build queued` と queue id を表示し、`Build started` や build id を合成せず、status と queue を再取得する。 |
| `partial-ui-runtime-build-dispatch-fallback` | `buildForce()` が queue id と `dispatch:"timer_fallback"` を返す。 | `Force build queued` と runner activation deferred warning を表示し、失敗表示や自動 retry に変換しない。 |
| `success-ui-runtime-queue-active-waiting` | `getQueue()` が active 1 件、queued 2 件、max size を返す。 | active を専用行、waiting 2 件を API 順で表示し、clear 確認件数は 2。active を waiting 件数へ含めない。 |
| `success-ui-runtime-stream` | `streamBuild()` が log 2 件と完了済み end 1 件、および別ケースで `status:"running",duration_seconds:null` の end 1 件を返し、EOF で終わる。 | log 行 2 件を appendし、`onEnd` と `done` が同じ summary を 1 回ずつ受け、`done` resolve 後に status、queue、logs を順に再取得して stream indicator を消す。running end は途中保存 snapshot と表示し、完了表示へ変換しない。 |
| `noop-ui-runtime-stream-user-close` | ユーザーが `StreamHandle.close()` を 2 回押す。 | 例外と error 表示なし、`done` は `null` で 1 回 resolve、`onEnd` は 0 回、closed 表示、status/queue 再取得あり。 |
| `failure-ui-runtime-history-validation` | `getHistory()` が `422 details` を返す。 | 該当 filter field に message を紐付け、history rows を前回表示のまま維持する。 |
| `failure-ui-runtime-log-not-found` | `getHistoryLog(id)` が `404 Not found` を返す。 | detail panel に not found を表示し、履歴一覧は再取得しない。 |
| `success-ui-runtime-circuit-reset` | `resetCircuitBreaker()` 成功。 | circuit 表示を閉じ、status/queue を再取得し、build を自動開始しない。 |
| `security-ui-runtime-unauthorized` | 初期表示の `getStatus()` が `401` を返す。 | token/ticket/secret field を消去し、`panel-login` だけ表示する。 |
| `failure-ui-operations-config-validation` | `setConfig()` が `422 details` を返す。 | 該当 field に error、panel error summary 1 行、入力値保持、`getConfig()` を呼ばない。 |
| `failure-ui-operations-schedule-save` | `setScheduleInterval()` が `500` を返す。 | panel error 表示後に `getSchedule()` を 1 回呼び、保存済み値を表示する。 |
| `security-ui-operations-secret-save` | `setWebhookConfig()` または `setSmtpConfig()` が `500` を返す。 | secret field を消去し、secret 平文を error 表示しない。 |
| `success-ui-operations-notify-test` | `notifyTest()` 成功。 | 結果表示後に `getNotifyLog()` を呼び、通知設定を自動保存しない。 |
| `noop-ui-snapshot-delete` | delete 確認 dialog cancel。 | SDK method 呼び出し 0 回、success / error 表示差分なし。 |
| `failure-ui-rollback-conflict` | `rollbackHistory()` が `409 Build is running` を返す。 | error 表示、`getStatus()` を呼ぶ、rollback request を再送しない。 |
| `success-ui-maintenance-enabled` | `getMaintenance()` が enabled を返す。 | maintenance banner 表示、build / rollback / 設定変更系 disabled、disable maintenance は enabled。 |
| `failure-ui-alert-rule-duplicate` | `addAlertRule()` が `409 Conflict` を返す。 | 競合表示、rule list は前回表示を保持し、自動 retry しない。 |
| `failure-ui-dashboard-layout-invalid` | `setDashboardLayout()` が `422 details` を返す。 | 該当 widget field error、dashboard 表示順を変更しない。 |
| `partial-ui-dashboard-unknown-widget` | `getDashboardLayout()` が `["status","unknown","stats"]` を返す。 | `status`、`stats` だけ表示し、順序保持。`unknown` は panel error 1 行。layout 保存を自動実行しない。 |
| `success-ui-compare-two-builds` | history 2 件選択後、左右の `getHistoryLog()` が異なる stdout を返す。 | 左右ログを別 column で API 行順表示し、差分 class は DOM 一時表示だけ。状態保存 API を呼ばない。 |
| `partial-ui-compare-missing-build` | 右側 `getHistoryLog()` が `404` を返す。 | compare panel error、左側表示は維持、history 再取得なし、選択値は保持。 |
| `success-ui-approvals-expired` | `getApprovals()` が `status:"expired"` を含む。 | approve / reject button disabled、期限切れ表示、UI が pending へ戻さない。 |
| `success-ui-approval-force-visible` | `getApprovals()` が `requested_trigger:"manual"`, `requested_force:true`, `status:"pending"` を返す。 | approve 前に requested trigger と強制 build 表示を出し、requested_force を再計算または非表示にしない。 |
| `partial-ui-approval-dispatch-fallback` | `approveBuild(id)` が `queued:true`, queue id, `dispatch:"timer_fallback"` を返す。 | approve 成功と runner 待機 warning を表示し、approval / queue を再取得する。approve を再送しない。 |
| `failure-ui-approval-approve-conflict` | `approveBuild(id)` が `409` を返す。 | error 表示後に `getApprovals()` を 1 回呼び、同じ approve を再送しない。 |
| `success-ui-notes-preserve-content` | notes に前後空白と連続改行を含めて保存。 | `setNotes(content)` へ入力値そのまま送信し、trim しない。 |
| `success-ui-hook-command-args` | 3 行の command args を入力し、中央行が空。 | 空行を除いた配列を `addHook()` に渡し、shell 文字列を作らない。 |
| `failure-ui-pipeline-reserved-arg` | `setPipelineConfig()` が `422 details` を返す。 | field error を表示し、入力値を保持し、`getPipelineConfig()` を呼ばない。 |
| `security-ui-token-issued-clear` | `createToken()` が token 本体を返し、trusted event 全 6 種、synthetic event、clipboard success / failure / API 不在、新旧 generation を fake する。 | [UI one-time secret 消去契約](ui.md#ui-one-time-secret-contract) の有効化タイミング、消去 event、除外 event、copy 1 回、`Copied` / `Copy failed`、generation guard、3 領域の同時消去を固定する。`getTokens()` の一覧に token 本体を表示せず、UI は追加 SDK method を呼ばない。 |
| `success-ui-disabled-priority` | maintenance enabled 中に `429` が発生し 10 秒経過。 | maintenance が継続する限り build / rollback / 設定変更系は disabled のまま。 |

UI owner fixture が不足する場合、UI 実装変更は詳細実装確認を満たした扱いにしてはならない。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

<a id="sec-27-f-19"></a>
**[fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture 固定契約](fixture.md#sec-27-f-19)：**

[`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) の setup、admin UI 配布、API service 導入、update、rollback を含む実装変更は、対象機能の owner fixture に加えて次の選択条件に該当する連動 fixture をすべて作成する。Release asset または admin archive layout を追加・変更する場合は `success-setup-admin-release-asset-layout`、Release download、redirect、timeout、size上限を実装または変更する場合は`failure-setup-download-boundary`、管理API追加時の既存binary version cohortを実装または変更する場合は`failure-setup-api-version-cohort`、unsafe archive の拒否境界を変更する場合は `security-setup-admin-archive-boundary`、systemd 更新失敗時 rollback を変更する場合は `partial-setup-systemd-rollback-boundary`、API service と runner dispatch / timer fallback を変更する場合は `partial-setup-api-runner-dispatch`、admin static serving を変更する場合は `security-admin-static-serving`、既存 secret の保持境界を変更する場合は `security-setup-secret-preservation` を必須とする。複数条件に該当する場合は該当 fixture を省略せず、[`docs/details/setup.md` 詳細本文責務 §26.8](setup.md#sec-26-8) と [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界)〜[§A6](admin.md#a6-admin-fixture-参照契約) の合格条件を同じ expected で検証する。

| setup/admin/Release asset fixture 名 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|-----------------------------------|----------------|------------|---------------|----------|
| `success-setup-admin-release-asset-layout` | `setup`、`admin` | Release asset 一式、`SHA256SUMS`、`admin-ui.tar.gz`、fake download response、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 | asset 名、checksum 対象、admin archive root layout、`index.html` と `adlaire-ci-sdk.js` だけを含む file set、file mode、directory mode が [state diff expected 固定契約](#fixture-state-diff-expected-contract) と配布正本に一致し、未定義 file を拒否する。 |
| `failure-setup-download-boundary` | `setup` | HTTP scheme、userinfo、相対Location、redirect 4回、未知host、Content-Length 0 / 上限超過 / body不一致、timeout、partial write failureの各fake response、`input/cli.json.expected_exit_code=3`。 | mode別`download` stageまでの`expected/stdout.txt`、stderr exact `setup: error DOWNLOAD_FAILED stage=download rollback=none`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | 未許可URLへrequestせず、上限+1 byteで停止し、partial fileだけを削除する。既存asset、binary、admin、state、secret、unitを変更せず、Authorization、Cookie、Referer、URL、response bodyを出力しない。 |
| `failure-setup-api-version-cohort` | `setup` | build / runner / setup各binaryの不在、symlink、非0、stderr非空、名前 / version / go token不一致、timer inactive、runner unit不在の各case、`input/cli.json.expected_exit_code=2`。 | `validate` stage startまでの`expected/stdout.txt`、stderr exact `setup: error PRECONDITION_FAILED stage=validate rollback=none`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | 最初の不合格で停止し、directory作成、download、API / admin / credentials / unit変更を0件とする。3 binaryがtarget versionと一致するcaseだけ後続へ進む。 |
| `security-setup-admin-archive-boundary` | `setup`、`admin` | unsafe archive、既存 `$INSTALL_DIR/admin`、既存 API binary、API service fake、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/effects.json` の `unchanged_paths`、`forbidden_writes`、`forbidden_calls`、`expected/security.json`。 | unsafe archive では admin directory、API binary、credentials、runner state を変更せず、API service start / restart を呼ばない。 |
| `partial-setup-systemd-rollback-boundary` | `setup`、`runner`、`api` | systemd fake、全対象の旧 binary backup、旧 admin backup、API restart failure、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json` の `write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes`、`commands`、`expected/security.json`。 | API restart 失敗後の rollback を 1 回だけ実行し、旧 `admin/` と全対象 binary が更新前の version cohort に戻ることを固定する。`commands` は更新処理と rollback 処理を合わせ、runner timer restart 2 回、API restart 2 回、rollback 後の両 service の `is-active`、runner / API の 3 unit の `systemctl cat` を固定する。state、history、secret、credentials、systemd unit は変更せず、health HTTP request は実行しない。`response` assertion と `expected/response.json` を作成しない。 |
| `partial-setup-api-runner-dispatch` | `setup`、`runner`、`api` | API / runner / timer unit、systemd fake、manual queue request、`input/cli.json.expected_exit_code`。 | `expected/response.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json` の `commands` と `write_order`、`expected/security.json`。 | API unit が timer を Wants / After し、`KillSignal=SIGTERM`、`TimeoutStopSec=15s` を持ち、`SendSIGKILL=no`、`KillMode=none`、独自 `ExecStop`、shell wrapper を持たない。queue 保存後だけ `systemctl start --no-block adlaire-ci.service` を 1 回実行し、systemctl failure でも queue を保持して timer fallback を返す。 |
| `security-admin-static-serving` | `admin`、`api` | static request、secret/state/log/snapshot path、method variation。 | `expected/response.json`、`expected/state/`、`expected/security.json`、`expected/effects.json`。 | A3 の status、header、body 有無、method 制限、no mutation、secret 非表示が一致する。 |
| `security-setup-secret-preservation` | `setup`、`security`、`statefile` | 既存 `.github_token` / `.last_sha` / `.admin_credentials` / その他 secret files、fresh setup と update 入力、管理 API 新規導入では正常および path・type・symlink・mode・owner が不正な `ADMIN_INITIAL_PASSWORD_FILE`、create-only 競合、rename 前失敗、rename 後 partial failure、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/security.json`、`expected/effects.json` の `write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes`。 | fresh setup は `.github_token` → `.last_sha` の create-only 順序、新規 payload / mode / LF、競合後の read-only 再検証、既存有効 target の content / mode / mtime 保持、既存不正の無修復 exit `2`、statefile failure の exit `1`、rename 後 target 維持を固定する。管理 API 新規導入は、credentials 不在時だけ検証済み password file descriptor を child stdin へ1回渡し、入力fileを変更せず、不正fileではcredentials / systemdを変更しない。既存credentials時は password fileを要求、読取せず、mode / schema 合格時だけ content / mode / mtimeを保持して成功、不合格時は無修復 exit `2` とする。update は initializer を呼ばず、全 secret / state の content / mode / mtime を保持する。stdout、stderr、journal、argv、environment、expected に secret 原文を残さない。 |

setupをownerまたはcollaboratorに含む全fixtureの`expected/stdout.txt`と`expected/stderr.txt`は[`docs/details/setup.md` 詳細本文責務 §26.2d](setup.md#setup-output-contract)をbyte単位で検証する。成功caseはmode別stageのstart / ok全行と最終success行、stderr 0 byteを必須とする。失敗caseは失敗stageのstart行までのstdout、exact error code、stage、rollback-result、LF 1個のstderrを必須とし、child process出力、path、URL、journal本文、Go error、secret原文を含めてはならない。Release asset download失敗caseの`input/cli.json.expected_exit_code`は`3`とする。

setup / admin / Release asset 連動 fixture の `manifest.json.name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) の fixture 名別割り当てに従う。`setup` と `admin` のうち owner ではない component を `collaborator_components` に含める。実装変更が API service 起動、static serving、rollback、secret 保持を扱う場合は、`api`、`runner`、`security` のうち当該 fixture の対象 component を collaborator として追加し、`expected/effects.json` の `forbidden_calls` と `forbidden_writes` に禁止副作用を明記する。

<a id="admin-cli-fixture-contract"></a>
**[fixture 証跡責務 Admin CLI fixture 固定契約](fixture.md#admin-cli-fixture-contract)：**

Admin CLI fixture は [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) だけを確認する。API endpoint の response schema は [`docs/details/api.md`](api.md) 詳細本文責務、認証・secret の扱いは [`docs/details/security.md`](security.md) 詳細本文責務を参照し、本契約で再定義しない。

| fixture 名 | fixture directory | 必須確認 |
|------------|-------------------|----------|
| `partial-admin-cli-lifecycle` | `testdata/admin/cli/partial-admin-cli-lifecycle/` | `--help`、`--version`、argv token safety、option parse、必須 option、未知 command、終了 code、network / state no-write。 |
| `success-admin-cli-transport` | `testdata/admin/cli/success-admin-cli-transport/` | `--api-url` path prefix 連結、method、path、header、request body byte、redirect 不追従、retry なし、proxy なし、timeout。 |
| `failure-admin-cli-output-errors` | `testdata/admin/cli/failure-admin-cli-output-errors/` | command 別 human stdout、`--json` raw JSON、invalid JSON、Content-Type 不一致、body 上限超過、HTTP error、network error。 |
| `security-admin-cli-secret-redaction` | `testdata/admin/cli/security-admin-cli-secret-redaction/` | token、Authorization header、error body、URL、Location、server body 断片、fixture expected への secret 非出力。 |

Admin CLI fixture の expected file は fixture 名ごとに以下へ固定する。対象外の expected file は `manifest.json.not_applicable` に理由を記録する。

| fixture 名 | 必須 expected |
|------------|---------------|
| `partial-admin-cli-lifecycle` | `expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`。 |
| `success-admin-cli-transport` | `expected/request.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`。 |
| `failure-admin-cli-output-errors` | `expected/response.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`。 |
| `security-admin-cli-secret-redaction` | `expected/request.json`、`expected/response.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`。 |

Admin CLI fixture の `manifest.json.assertions` は次表に固定する。複数値は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) の列挙順で記録する。

| fixture 名 | 必須 assertions |
|------------|-----------------|
| `partial-admin-cli-lifecycle` | `stdout`、`stderr`、`effects`、`secret-mask`、`no-write` |
| `success-admin-cli-transport` | `request`、`stdout`、`stderr`、`effects`、`secret-mask`、`order` |
| `failure-admin-cli-output-errors` | `response`、`stdout`、`stderr`、`effects`、`secret-mask`、`no-write` |
| `security-admin-cli-secret-redaction` | `response`、`request`、`stdout`、`stderr`、`effects`、`secret-mask`、`no-write` |

Admin CLI fixture の `expected/request.json` は root object とし、root key は `requests`、`forbidden_requests` だけを許可する。`requests` は実際に送信する HTTP request を command 実行順で持つ配列、`forbidden_requests` は送信してはならない HTTP request または transport 挙動を列挙する配列とする。未知 root key、未知 request key、未知 forbidden request key を含む fixture は不合格とする。

| 対象 | key | 型 | 固定契約 |
|------|-----|----|----------|
| `requests[]` | `order` | integer | 1 始まりの送信順。重複、欠番、0 以下を禁止する。 |
| `requests[]` | `command` | string | [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) の command 固定表にある 7 command のいずれか。 |
| `requests[]` | `method` | string | 実送信 method の uppercase 文字列。 |
| `requests[]` | `path` | string | 正規化後の `--api-url` path prefix と command 固定 path を byte 連結した path。query、fragment、percent decode 後の値を含めてはならない。 |
| `requests[]` | `headers` | object | CLI が明示的に設定する header だけを持つ。`Authorization` は `${secret:admin_cli_token}` 固定、`Accept` は `application/json` 固定、`User-Agent` は `adlaire-ci-admin/<binary-version>` 固定、body を持つ command だけ `Content-Type: application/json` を持つ。 |
| `requests[]` | `body` | string または null | body なし command は null。`trigger-build` は `{}`。`config-snapshot` は未指定時 `{"label":null}`、指定時 `{"label":"<label>"}`。UTF-8、末尾 LF なし、余分な空白なしで固定する。 |
| `requests[]` | `redirects_followed` | integer | 常に 0。 |
| `requests[]` | `retry_count` | integer | 常に 0。 |
| `requests[]` | `proxy_used` | boolean | 常に false。 |
| `requests[]` | `cookies_sent` | integer | 常に 0。 |
| `forbidden_requests[]` | `method` | string | 送信禁止 request method。transport 挙動だけを禁止する場合は空文字を許可する。 |
| `forbidden_requests[]` | `path` | string | 送信禁止 request path。transport 挙動だけを禁止する場合は空文字を許可する。 |
| `forbidden_requests[]` | `reason` | string | `wrong-endpoint`、`redirect-follow`、`retry`、`proxy`、`cookie`、`secret-leak` のいずれか。 |

`success-admin-cli-transport` の `forbidden_requests` は、`trigger-build` に対する `POST /api/builds`、redirect 追従、retry、proxy、Cookie 送信を必ず含める。`security-admin-cli-secret-redaction` の `expected/request.json` は同じ schema を使用し、token、Authorization header 実値、URL query 内 secret、redirect `Location` 内 secret、server raw body 内 secret、Go error 内 secret、absolute path 内 secret を `requests`、`forbidden_requests`、`headers`、`body`、`reason` に平文または派生値として含めてはならない。

Admin CLI fixture は CLI 実行入力を `input/cli.json` に固定し、HTTP 応答 fake を使用する fixture は `input/fakes.json` の `fetch` root を使用する。Admin CLI fixture で `input/request.json` を作成してはならない。`partial-admin-cli-lifecycle` のうち API 呼び出しへ到達しない case は HTTP fake を使用せず、`input/fakes.json` を置く場合でも `fetch=[]` とする。`success-admin-cli-transport`、`failure-admin-cli-output-errors`、`security-admin-cli-secret-redaction` のうち API 呼び出しへ到達する case は、実 request 1 件につき `input/fakes.json.fetch[]` 1 件と `expected/effects.json.external_calls[]` 1 件を同じ順序で持つ。

Admin CLI 用 `input/fakes.json.fetch[]` は次表の固定契約に従う。未知 key、実 network、実 proxy、実 cookie jar、実 redirect、実 retry、環境変数由来 proxy を禁止する。

| 対象 | key | 型 | 固定契約 |
|------|-----|----|----------|
| `fetch[]` | `operation` | string | `request` 固定。 |
| `fetch[]` | `target` | string | `admin_cli_api` 固定。 |
| `fetch[].input` | `method` | string | Admin CLI が実送信する uppercase method。 |
| `fetch[].input` | `path` | string | 正規化後の `--api-url` path prefix と command 固定 path を byte 連結した path。query、fragment、scheme、host を含めない。 |
| `fetch[].input` | `headers` | object | CLI が明示的に送る header。`Authorization` は `${secret:admin_cli_token}`、`Accept`、`User-Agent`、body あり command の `Content-Type` だけを許可する。 |
| `fetch[].input` | `body` | string または null | 送信 body。body なしは null。byte 表現は [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) の request body 契約と完全一致させる。 |
| `fetch[].input` | `timeout_milliseconds` | integer | `30000` 固定。 |
| `fetch[].input` | `redirect_policy` | string | `manual` 固定。 |
| `fetch[].input` | `retry_count` | integer | `0` 固定。 |
| `fetch[].input` | `proxy_used` | boolean | false 固定。 |
| `fetch[].input` | `cookies_sent` | integer | `0` 固定。 |
| `fetch[].output` | `status` | integer または null | HTTP response を受け取った場合は `100`〜`599`。network error、TLS error、timeout、connection close before response は null。 |
| `fetch[].output` | `headers` | object | response header。header なしまたは response 確定前失敗は空 object。key は lowercase ASCII、値は受信順に `, ` で連結した string。 |
| `fetch[].output` | `body` | string または null | 1 MiB 以下で body read 成功時だけ raw body byte を UTF-8 string として置く。1 MiB 超過、body read timeout、connection close、response 確定前失敗は null。 |
| `fetch[].output` | `body_size_bytes` | integer | 実際に受信または fake が提示する body byte 数。body なしは 0。 |
| `fetch[].output` | `body_sha256` | string または null | `body` が null で `body_size_bytes>0` の場合だけ 64 文字 lowercase hex。通常 body を `body` に置く場合は null。 |
| `fetch[].output` | `body_read_result` | string | `success`、`over_limit`、`timeout`、`connection_closed`、`not_started` のいずれか。 |

Admin CLI 用 `fetch[]` の `result="success"` は HTTP response を受け取ったことだけを意味し、HTTP status `300`〜`599`、redirect response、invalid JSON、invalid Content-Type、1 MiB 超過 response も `result="success"` とする。1 MiB 超過 response は `error_code=null`、`output.body=null`、`output.body_read_result="over_limit"` とし、期待 stderr は `api error: invalid response` + LF に固定する。network error、TLS error、timeout、connection close before response、response body read timeout、response body read 中の connection close は `result="failure"` または `result="timeout"`、`output.body=null`、`error_code="connection_failed"` とし、期待 stderr は `api error: connection failed` + LF に固定する。

Admin CLI の `expected/effects.json.external_calls[]` は outbound HTTP 証跡として次を固定する。各要素は `operation="call"`、`target="admin_cli_api"`、`input` は対応する `fetch[].input` と同一 key / 同一値、`output` は `status`、`body_size_bytes`、`body_sha256`、`body_read_result` の 4 key、`result` と `error_code` は対応する `fetch[]` と同一値にする。response body 本文と secret を `expected/effects.json` に保存してはならない。`expected/request.json.requests[]`、`input/fakes.json.fetch[]`、`expected/effects.json.external_calls[]` は、実 request 件数、順序、method、path、header、body、redirect、retry、proxy、cookie が一致しなければならない。

Admin CLI fixture の stdout / response 検証は [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) の human stdout 写像契約に従う。`success-admin-cli-transport` は HTTP transport の成功に加えて、代表 success response から `expected/stdout.txt` を byte 単位で生成できることを固定する。`failure-admin-cli-output-errors` は HTTP status が `2xx` であっても、command ごとの stdout 写像に必要な key、型、許容値が不足または不一致の response を invalid response として扱う case を持つ。

| command | success stdout の source | invalid response 必須 case |
|---------|--------------------------|----------------------------|
| `status` | `last_build_status` string、`running` boolean。 | `last_build_status` 欠落、`running` 欠落、`running` string。 |
| `queue` | `active=null` または `active.id` string、`queued` array length。 | `queued` number、`queued` 欠落、`active` object の `id` 欠落、`active.id` number。 |
| `history` | `total` integer、`history[0].id` string。`history=[]` の場合だけ latest は `none`。 | `total` string、`history` 欠落、非空 `history[0].id` 欠落。 |
| `trigger-build` | `queue_id` string、`queued=true`、`dispatch` が `requested` または `timer_fallback`。 | `queue_id` 欠落、`queue_id` 空文字、`queued=false`、未知 `dispatch`。 |
| `cancel-queue` | response body の `message` に依存せず固定 `queue cancelled`。 | JSON object 以外、body 空。 |
| `config-snapshot` | `id` string。 | `id` 欠落、`id` 空文字、`id` number。 |
| `events` | `total` integer。`events` array は存在だけを検証し、stdout 件数には使わない。 | `total` 欠落、`total` string、`events` 欠落、`events` number。 |

`--json` success case は、command ごとの成功 response shape が有効な場合だけ API wire body から前後 ASCII whitespace を除去した byte 列に LF 1 個を付けた `expected/stdout.txt` を固定する。`--json` 指定時であっても、上表の invalid response 必須 case は stdout 空、stderr `api error: invalid response` + LF、終了 code `1` とする。`events.total` と `events.length` が異なる schema-valid response は valid case とし、human stdout は `total` だけを使用することを固定する。

`partial-admin-cli-lifecycle` は [`docs/details/admin.md` 詳細本文責務 §A7](admin.md#sec-a7) の parse / validation 順序を fixture で固定する。`--help` と `--version` は API URL、token、state directory、network fake、file read/write、乱数取得を一切参照しないことを `expected/effects.json.external_calls=[]`、`forbidden_writes`、`forbidden_reads` で固定する。parse error、未知 command、引数不足、引数過多、同一 option 重複、`--name=value`、短縮 option、禁止制御文字は個別 case とし、stdout 空、stderr 1 行、終了 code `2` を byte 単位で検証する。

`partial-admin-cli-lifecycle` の必須 case は次表に固定する。各 case は個別入力として実行し、前段の不合格が後段の不合格より優先されることを確認する。parse / validation 失敗 case の `expected/effects.json` は HTTP request 0 件、state read 0 件、state write 0 件、file read 0 件、file write 0 件、random 0 件、child process 0 件、external call 0 件を固定する。

| case | 入力条件 | 期待 stdout / stderr / exit |
|------|----------|-----------------------------|
| `help-precedence` | `--help` と、未指定の `--api-url`、未指定の `--token`、未知 command 相当 token、不正 URL 相当 token を同時に含める。 | stdout は help 固定 1 行、stderr 空、exit `0`。 |
| `version-precedence` | `--version` と、未指定の `--api-url`、未指定の `--token`、未知 command 相当 token、不正 URL 相当 token を同時に含める。 | stdout は version 固定 1 行、stderr 空、exit `0`。 |
| `unsafe-token-priority` | `--help` / `--version` を含まず、NUL、CR、LF、C0 制御文字、DEL、invalid UTF-8 のいずれかを含む argv token と、未知 option または不正 URL を同時に含める。 | stdout 空、stderr `invalid command line token` + LF、exit `2`。token 原文を expected に含めない。 |
| `missing-api-url-value` | `--api-url` の次 token がない、または次 token が `--` で始まる。 | stdout 空、stderr `missing value: --api-url` + LF、exit `2`。 |
| `missing-token-value` | `--token` の次 token がない、または次 token が `--` で始まる。 | stdout 空、stderr `missing value: --token` + LF、exit `2`。 |
| `unknown-option` | `--unknown` を command より前に含める。 | stdout 空、stderr `unknown option: --unknown` + LF、exit `2`。 |
| `equals-option` | `--api-url=http://127.0.0.1:8765` または `--token=value` を含める。 | stdout 空、stderr `unknown option: <token>` + LF、exit `2`。 |
| `short-option` | `-u`、`-t`、`-j` のいずれかを含める。 | stdout 空、stderr `unknown option: <token>` + LF、exit `2`。 |
| `duplicate-api-url` | `--api-url` を 2 回指定する。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `duplicate-token` | `--token` を 2 回指定する。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `duplicate-json` | `--json` を 2 回指定する。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `missing-required-api-url` | `--token` と command は有効、`--api-url` は未指定。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `missing-required-token` | `--api-url` と command は有効、`--token` は未指定。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `missing-command` | `--api-url` と `--token` は有効、command は未指定。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `invalid-api-url` | `--api-url` が scheme 不正、host 不在、userinfo あり、query あり、fragment あり、`..` segment、重複 slash、backslash のいずれか。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `invalid-token` | `--token` が空文字、4097 byte 以上、NUL、CR、LF のいずれか。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `unknown-command` | `--api-url` と `--token` は有効、command が固定 7 command 以外。 | stdout 空、stderr `unknown command: <command>` + LF、exit `2`。`<command>` は secret ではない固定入力値にする。 |
| `cancel-queue-arg-missing` | `cancel-queue` の `<queue_id>` 未指定。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `cancel-queue-arg-invalid` | `cancel-queue` の `<queue_id>` が空文字、`/`、`..`、NUL byte のいずれかを含む。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `cancel-queue-arg-extra` | `cancel-queue` に `<queue_id>` 以外の追加 token がある。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `config-snapshot-label-invalid` | `config-snapshot` の `label` が空文字、129 Unicode scalar values 以上、改行、NUL byte、BOM のいずれか。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `config-snapshot-arg-extra` | `config-snapshot` に `label` 以外の追加 token がある。 | stdout 空、stderr `usage error` + LF、exit `2`。 |
| `no-arg-command-extra` | `status`、`queue`、`history`、`trigger-build`、`events` のいずれかに追加 token がある。 | stdout 空、stderr `usage error` + LF、exit `2`。 |

`success-admin-cli-transport` は各 command について method、path、request body、header を `expected/request.json` に固定する。`Authorization` は placeholder `${secret:admin_cli_token}` だけを許可し、token 実値、token hash、部分文字列、長さから復元できる値を expected に置いてはならない。`config-snapshot` の JSON body、`cancel-queue` の 1 回だけの percent encode、`User-Agent`、`Accept`、`Content-Type`、redirect 不追従、retry 0 回、proxy 0 回、Cookie 0 件を検証する。

`success-admin-cli-transport` の必須 case は次表に固定する。各 case は `input/cli.json`、`input/fakes.json.fetch[]`、`expected/request.json.requests[]`、`expected/effects.json.external_calls[]`、`expected/stdout.txt`、`expected/stderr.txt` を byte 単位で照合する。

| case | 入力条件 | 必須検証 |
|------|----------|----------|
| `status-default-origin` | `--api-url http://127.0.0.1:8765 status`。 | `GET /api/status`、body null、`Authorization` / `Accept` / `User-Agent` だけ、stdout `status=<last_build_status> running=<running>`、stderr 空。 |
| `status-path-prefix` | `--api-url http://127.0.0.1:8765/adlaire/ status`。 | request path は `/adlaire/api/status`。末尾 `/` を 1 個だけ除去し、`/api` 重複除去、path clean、percent decode をしない。 |
| `queue-active-null` | `queue` が `active:null` と `queued` array を返す。 | `GET /api/queue`、stdout は `active=none queued=<array-length>`。 |
| `history-empty` | `history` が `total` integer と空 `history` array を返す。 | `GET /api/history`、stdout は `total=<total> latest=none`。 |
| `trigger-build-endpoint` | `trigger-build`。 | `POST /api/build`、body `{}`、`Content-Type: application/json`、stdout `queued=<queue_id>`、`forbidden_requests` に `POST /api/builds` を含める。 |
| `cancel-queue-percent-encode` | `cancel-queue` の `<queue_id>` に percent encode が必要な文字を含める。 | `DELETE /api/queue/{queue_id}` の path parameter を 1 回だけ encode する。slash 生成、2 重 encode、query 化を禁止する。 |
| `config-snapshot-null-body` | `config-snapshot` label 未指定。 | `POST /api/config-snapshots`、body `{"label":null}`、末尾 LF なし、stdout `snapshot=<id>`。 |
| `config-snapshot-label-body` | `config-snapshot` label 指定。 | body は `{"label":"<label>"}`。JSON string escape、UTF-8、余分な空白なしを固定する。 |
| `events-total` | `events` が `total` integer と `events` array を返す。 | `GET /api/events`、stdout は `events=<total>`。`events` array length を stdout に使用しない。 |
| `json-mode-wire-body` | `--json status`。 | response shape 合格後、API wire body の前後 ASCII whitespace だけを除去し、key order と number / string 表現を保持して LF 1 個を追加する。 |

`success-admin-cli-transport` の全 case は `expected/request.json.forbidden_requests` に redirect 追従、retry、proxy、Cookie 送信を含める。redirect 先 request、2 回目以降の同一 request、proxy 経由 request、Cookie 付き request が `input/fakes.json.fetch[]` または `expected/effects.json.external_calls[]` に 1 件でも存在する場合は不合格とする。

`failure-admin-cli-output-errors` は success body と failure body を混在させない。invalid JSON、複数 JSON value、body 空、Content-Type 不一致、複数 Content-Type、1 MiB 超過、HTTP `300`〜`599`、network error、TLS error、timeout、connection close before response を別 case とし、stdout 空、stderr 固定 1 行、終了 code `1` を byte 単位で検証する。`--json` success case では API wire body の key order と number / string 表現を保持し、末尾 LF 1 個だけを追加することを確認する。

`security-admin-cli-secret-redaction` は token、Authorization header、error response body 内 secret、redirect `Location`、URL query、server raw body、Go error、absolute path が stdout、stderr、`expected/request.json`、`expected/response.json`、`expected/effects.json`、`expected/security.json` に平文で出現しないことを列挙する。placeholder は `${secret:<source_id>}` だけを使用し、実値、hash 入力、prefix、suffix、長さ、base64、percent encode された派生値を禁止する。

Admin CLI fixture の `manifest.json.name`、`section="A7"`、`feature="admin_cli"`、`owner_component="admin"` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) に従う。HTTP fake を使用する fixture は `api` を collaborator、secret / auth 境界を検証する fixture は `security` を collaborator に含める。`components` は `admin` と collaborator を ASCII 昇順で持つ。`references` は先頭に `docs/details/admin.md#sec-a7`、次に `docs/details/fixture.md#admin-cli-fixture-contract`、以後 collaborator 詳細本文を UTF-8 byte 列の昇順で持つ。

<a id="release-fixture-contract"></a>
**[fixture 証跡責務 Release fixture 固定契約](fixture.md#release-fixture-contract)：**

[`docs/details/release.md`](release.md)の実装変更は次の全fixtureを必須とする。全caseで実GitHub、実tag変更、実branch変更、実Release作成を禁止し、`input/fakes.json`の`git`、`command`、`filesystem`、`archive`、`github`、`download` eventだけを使用する。

| fixture名 | 必須input | 必須expected | 合格条件 |
|-----------|-----------|--------------|----------|
| `success-release-assets-reproducible` | clean Git状態、tag / commit、2個のtemporary root、同一commitから生成した独立`git archive` stream 2件、未作成out path、Go command fake、admin配布物。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/release-assets.json`、`expected/security.json`。 | 2つのsource snapshotのfile set / mode / size / digest、7 binaryとadmin archiveのbyte、固定build引数、`GOENV=off` / `GOWORK=off` / `GOTOOLCHAIN=local` / network offを含むenvironment、version出力、commit timestamp、USTAR / gzip level 9を含むarchive metadata、asset mode、checksum 8行が一致する。保持済みout parent descriptor相対のsibling staging作成、file sync / directory sync / atomic rename / parent sync順、local Git再検証、temporary cleanupを固定し、live checkout file readとGitHub writeは0件。 |
| `failure-release-output-parent-race` | 初回検証時のparent symlink、検証後のsymlink差替え、parent directory置換、`<out>.tmp`競合、異なるmountへの誘導を個別に発生させるfake。 | 初回不正はstderr`release: INVALID_INPUT`、検証後の差替え / 競合はstderr`release: OUTPUT_FAILED`、全caseでstdout空、`expected/effects.json`、`expected/security.json`。 | 最終parent descriptorを1回だけ保持し、staging作成前とrename前にdevice / inodeを再検証する。通常path操作、descriptor外write、既存path上書き、checkout内write、GitHub writeを0件とし、staging作成後の失敗は保持済みdescriptor相対でstagingだけをcleanupする。 |
| `success-release-draft-publish` | 検証済み9 asset、memory保持済みnotes、token placeholder、local Git、repository / default branch ref / compare / tag / Release read、create / upload / download / patch fake。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/response.json`、`expected/effects.json`、`expected/release-assets.json`、`expected/security.json`。 | write前local / remote再検証、draft作成、draft直後remote再検証、9件順次upload、一覧再取得、9件download digest、publish直前remote / metadata再検証、publish、公開後再取得の順序が一致する。 |
| `failure-release-dirty-worktree` | tracked変更またはuntracked fileを含むGit status。 | stdout空、stderr`release: DIRTY_WORKTREE`、`expected/effects.json`、`expected/security.json`。 | directory作成、Go command、archive、GitHub callが0件である。 |
| `failure-release-version-mismatch` | HEAD、local tag、remote tag、`--commit`、remote default branch ref / compare、binary `--version`の各不一致case。 | Git ref / commit / ancestor不一致caseはstderr`release: TAG_MISMATCH`、binary出力不一致caseはstderr`release: VERSION_MISMATCH`、全caseで`expected/effects.json`、`expected/security.json`。 | local / remote Git・GitHub検証またはversion検証の最初の不一致で停止し、GitHub writeが0件、checkoutとtagが不変である。各stderrはLF 1個で終わる。 |
| `failure-release-source-snapshot-boundary` | 128 MiB超過、絶対path、`.` / `..`、backslash、symlink、hardlink、device、FIFO、socket、重複entry、root脱出、A / B不一致の各`git archive` fake。 | stdout空、stderr`release: BUILD_FAILED`、`expected/effects.json`、`expected/security.json`。 | unsafe entryへのwrite、existing path上書き、Go command、`--out`、GitHub writeが0件で、A / B temporary rootをcleanupする。 |
| `failure-release-source-race` | snapshot生成後にtracked変更、untracked追加、HEAD移動、local tag移動を個別に発生させるfake。 | dirtyはstderr`release: DIRTY_WORKTREE`、HEAD / tagはstderr`release: TAG_MISMATCH`、`expected/effects.json`、`expected/security.json`。 | build inputは先に固定したsnapshotだけとし、local再検証で失敗し、`<out>.tmp`を削除して`--out`とGitHub writeを作成しない。 |
| `failure-release-notes-identity` | notes fileのsymlink、open後のidentity / size / mtime変更、262145 bytes、invalid UTF-8、NUL / CR、末尾LF不正の各case。 | stdout空、stderr`release: INVALID_INPUT`、`expected/effects.json`、`expected/security.json`。 | no-follow openは1回、path再openは0回、directory作成、source snapshot、Go command、GitHub call、`--out`が0件である。 |
| `failure-release-token-identity` | token path segmentまたは対象fileのsymlink、path segment差替え、open後のidentity / mode / owner / size / mtime変更、4097 bytes、invalid UTF-8、trim後の空白・NUL・CR / LF・制御文字の各case。 | stdout空、stderr`release: INVALID_INPUT`、`expected/effects.json`、`expected/security.json`。 | descriptor相対no-follow open、対象file open 1回、path再open 0回、token原文の記録、directory作成、source snapshot、Go command、GitHub call、`--out`が0件である。 |
| `failure-release-go-environment` | hostに未定義`GO*`、workspace、auto toolchain、proxy / sumdb、A / B異値を注入する。 | stdout空、固定environmentの`expected/release-assets.json`、`expected/effects.json`、`expected/security.json`。 | childは固定対象だけを固定値で受け、その他の`GO*`を継承せず、toolchain / module networkを0件にする。fakeが固定値不一致を返すcaseは対応する`FORMAT_FAILED` / `TEST_FAILED` / `BUILD_FAILED`で停止する。 |
| `failure-release-gofmt-order` | A snapshot に `.go` regular file 3 件を ASCII 順外で配置し、2 件目だけ `gofmt -d` が差分 stdout を返す fake。 | stdout空、stderr`release: FORMAT_FAILED`、`expected/effects.json`、`expected/security.json`。 | `gofmt -d <path>` は relative path の ASCII 昇順で 1 file 1 process とし、2 件目の差分検出後に3件目を実行しない。`go test`、build、`--out`、GitHub writeは0件。 |
| `failure-release-permission` | read-only APIの`401` / `403`、draft createの`401` / `403` / permission-masked `404`、upload / PATCH / DELETEの権限不足、workflow変更時の`Workflows: write`不足を個別に返すGitHub fake。 | 操作別exact stderr、`expected/effects.json`、`expected/security.json`。 | 無認証 / 別token fallback、既存Release再利用、write retryを0件とし、draft id確定後だけcleanupを1回実行する。 |
| `partial-release-remote-race-cleanup` | draft作成直後またはpublish直前にdefault branchが`--commit`非包含へ変化、remote tag移動、repository / default branch不一致、Release metadata不一致を発生させるfake。 | stdout空、stderr`release: PUBLISH_FAILED`、`expected/effects.json`、`expected/security.json`。 | draft直後caseはasset upload 0件、publish直前caseはPATCH 0件とし、確定済みdraftだけをDELETE 1回でcleanupする。 |
| `failure-release-non-reproducible` | 2回目のbinaryまたはarchiveを1 byteだけ変更するfake。 | stdout空、stderr`release: NON_REPRODUCIBLE`、`expected/effects.json`、`expected/security.json`。 | `--out`公開成果物とGitHub writeが0件で、host temporary pathを出力しない。 |
| `failure-release-output-atomicity` | out path既存、`<out>.tmp`既存、staging作成、copy、chmod、file sync、directory sync、rename、parent sync、cleanupの各失敗fake。 | 入力不正caseはstderr`release: INVALID_INPUT`、出力処理caseはstderr`release: OUTPUT_FAILED`、`expected/effects.json`、`expected/security.json`。 | rename前失敗では`--out`を作らず、cleanup成功caseは`<out>.tmp`を削除する。parent sync失敗では検証済み6fileだけを持つ`--out`を維持し、GitHub writeを0件とする。cleanup失敗caseは元errorを維持してretryせず、残存`<out>.tmp`を再利用せず、pathを出力しない。checkoutと事前存在pathを変更しない。 |
| `failure-release-existing-release` | 同一tagのdraftまたはpublished Release read response。 | stdout空、stderr`release: RELEASE_EXISTS`、`expected/effects.json`、`expected/security.json`。 | 既存Releaseのupdate/delete、asset upload、local buildが0件である。 |
| `partial-release-upload-cleanup` | draft作成成功、asset 1〜5件目の任意upload failure、DELETE成功。 | stdout空、stderr`release: ASSET_UPLOAD_FAILED`、`expected/effects.json`、`expected/security.json`。 | failure後に追加uploadとpublishをせず、作成したdraftだけをDELETE 1回で削除する。 |
| `partial-release-cleanup-failure` | draft作成成功、upload failure、DELETE failure。 | stdout空、stderr`release: DRAFT_CLEANUP_FAILED`、`expected/effects.json`、`expected/security.json`。 | DELETEを再送せず、URL、draft id、response body、元errorをstdout / stderrへ出さない。fake call記録だけで対象draftへのDELETE 1回を検証し、tagとlocal assetを削除しない。 |
| `failure-release-timeout-boundaries` | checkout root、status、HEAD / local tag、commit timestamp、`git archive`、gofmt、test、build、binary version、GitHub metadata、upload、downloadの各timeout fake。 | 順に`INVALID_INPUT`、`DIRTY_WORKTREE`、`TAG_MISMATCH`、`BUILD_FAILED`、`BUILD_FAILED`、`FORMAT_FAILED`、`TEST_FAILED`、`BUILD_FAILED`、`BUILD_FAILED`、`GITHUB_READ_FAILED`、`ASSET_UPLOAD_FAILED`、`ASSET_VERIFY_FAILED`のexact stderrと、`expected/effects.json`、`expected/security.json`。 | 各timeout値でcontextをcancelし、childまたはrequestの終了後に後続を停止する。draft作成後caseだけcleanupを1回行い、その他のGitHub writeは0件とする。 |
| `security-release-http-boundary` | userinfo、HTTP scheme、redirect、未知host、過大metadata/error body、assetのContent-Length不在・不一致・超過case。 | create responseのupload / html URL不正は`DRAFT_CREATE_FAILED`、read-only metadata上限超過は`GITHUB_READ_FAILED`、asset redirect / length / body不正は`ASSET_VERIFY_FAILED`のexact stderrと、`expected/effects.json`、`expected/security.json`。 | 未許可URLへrequestせず、上限+1 byteで読取を停止し、Authorization、token、response body、URLを出力しない。 |
| `security-release-token-mask` | token原文、userinfo付き拒否URL、API error body、upload / download failure。 | stdout、stderr、effects、security、全asset。 | token、Authorization値、token file内容、credential付きURLが全expectedと生成物に存在せず、header有無だけをeffectsへ記録する。 |

`expected/release-assets.json`は`assets`、`checksum_lines`、`builds`、`archive`の4 root keyだけを持つ。`assets`はname、size、sha256、modeをasset名ASCII昇順で持ち、`builds`はbinary名、argv、environment、version_stdout、first_sha256、second_sha256を持つ。`archive`はentry名、typeflag、format、mode、uid、gid、uname、gname、mtime、atime、ctime、pax_records、gzip_level、gzip_mtime、gzip_name、gzip_comment、gzip_extra、gzip_osを持つ。unknown key、実token、host絶対pathを禁止する。

Release fixtureの`manifest.json.name`、`section="R1-R7"`、`feature="release_contract"`、`owner_component="release"`は[manifest識別子レジストリ固定契約](#sec-27-f-manifest-identity)に従う。asset受け入れの照合が必要なcaseだけ`setup`、admin archive内容の照合が必要なcaseだけ`admin`を`collaborator_components`へ追加する。

<a id="sec-27-f-20"></a>
**[fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](fixture.md#sec-27-f-20)：**

正式 fixture directory harness の実装検証証跡を追加または更新する変更は、実装検証証跡に [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表を記録する。記録がない formal fixture harness 変更は、コードと fixture が存在しても未完了とする。Phase 単位の実装検証は、同じ変更単位で記録する Phase 実装検証証跡を正とする。

| 記録項目 | 必須内容 |
|----------|----------|
| wave | 対象 wave、対象 [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、先行 wave 完了 commit または変更識別子。 |
| 実装対象 | 実装する機能名、owner component、collaborator component、変更ファイル、追加 fixture path。 |
| 実装対象外 | 同じ wave 内で今回実装しない [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、後続 wave、MCP、外部公開構成、未定義 endpoint / UI / 状態ファイルを実装検証証跡に列挙する。 |
| fixture | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) カタログの fixture 名、manifest / effects / security の検証結果。 |
| acceptance | [`docs/details/fixture.md`](fixture.md) fixture 証跡責務 runner / security 実装 acceptance checklist の各項目の pass / fail / 未実行。 |
| 後続影響 | 後続変更が利用許可済みの contract、利用禁止の未固定 contract を実装検証証跡に列挙する。 |

**runner / security 実装最終受け入れゲート：**

Formal fixture harness の runner / security 実装検証証跡は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の全 gate を満たした場合だけ「完了」と判定する。1 件でも未達がある場合は「未完了」、仕様逸脱または secret 漏えいリスクがある場合は「差し戻し」とする。Phase 単位の現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務と、該当する Phase 実装検証証跡によって判定する。

| gate | 完了条件 | 未完了条件 | 差し戻し条件 |
|------|----------|------------|--------------|
| scope | 実装対象が [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) カタログ固定契約に存在する機能だけである。 | 実装対象節の記載が実装検証証跡にない。 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加している。 |
| fixture | 対象 [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の必須 fixture がすべて存在し、skip されていない。 | 必須 fixture が不足、または fixture 名が不一致。 | fixture が実装挙動に合わせて期待値を緩めている。 |
| manifest | 全 fixture の `manifest.json` が schema、assertion 選択、owner / collaborator component 責務を満たす。 | assertion、references、owner_component、collaborator_components、components、fake_clock のいずれかが不足。 | unknown key、実 secret、実環境 path、乱数依存を含む。 |
| expected | `response`、SDK、UI、stdout/stderr、state、logs、effects、security の各 assertion に対応する [ファイルセット固定契約](#sec-27-f-7) の expected file が存在し、各固定 schema と一致する。 | assertion に対応する expected file が不足。 | expected と manifest / effects / security が矛盾する。 |
| side effect | `write_order`、`unchanged_paths`、`forbidden_writes`、`forbidden_calls` が対象機能の成功 / 失敗 / no-op / partial を説明できる。 | 禁止副作用または無変更保証が不足。 | 失敗時に未許可状態を書き換える、外部呼び出しを行う。 |
| secret | secret 平文が expected、logs、effects、UI DOM、stdout/stderr に存在しない。 | secret 検証対象が不足。 | token、password、TOTP secret、PAT、Authorization header が平文で残る。 |
| component | builder / runner / api / admin / sdk / ui / statefile / archive / commitstatus / security / setup / releaseの該当責務が全てfixtureに紐づく。 | owner / collaborator componentの所在が不明。 | SDK / UIがAPI responseを推測補完、またはUIが直接API / 状態ファイルを操作する。 |
| repeatability | fake clock、fake external response、固定 path により、同じ fixture が同じ結果を再現する。 | idempotency / no-op の 2 回目期待値が不足。 | 現在時刻、実ネットワーク、実 OS 差分に依存する。 |

**runner / security 実装 acceptance checklist：**

実装検証証跡には、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の項目を記録する。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

| 項目 | 記録内容 |
|------|----------|
| 対象責務参照 | 実装した [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様).x / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、関連 [`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) / [`docs/details/api.md` 詳細本文責務 §25](api.md#25-認証-実装仕様) / [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) / [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) / [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順)、owner component、collaborator component。 |
| 対象 fixture | 作成または更新した fixture 名一覧。fixture 名は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) カタログ固定契約と一致させる。 |
| 実行結果 | fixture ごとの pass / fail、実行コマンド、終了コード。 |
| 状態差分 | 作成、更新、削除、変更禁止の path。`expected/effects.json` と一致させる。 |
| 外部副作用 | GitHub、SSH、SMTP、webhook、systemd、hook、remote build、notification の呼び出し回数と順序。 |
| secret 検証 | 禁止文字列、mask 対象、平文が残らないことを確認した出力範囲。 |
| 部分失敗 | partial / failure fixture の失敗地点、完了済み副作用、禁止副作用。 |
| 再実行 | idempotency / no-op fixture の 1 回目と 2 回目の差分。 |
| 対象外確認 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加していないことを、差分対象 file と fixture manifest の両方で確認する。 |

**runner / security 差し戻し固定条件：**

次のいずれかに該当する変更は、fixture が pass していても差し戻しとする。

| 条件 | 理由 |
|------|------|
| 仕様にない endpoint、SDK method、UI 操作、状態ファイルを追加している。 | 仕様外実装。 |
| fixture の期待値が実装都合に合わせて仕様より弱い。 | 検証の形骸化。 |
| secret 平文、Authorization header、TOTP secret、PAT、password が expected または log に残る。 | secret 漏えい。 |
| read-only / dry-run / validation failure で、個別仕様または API 共通処理順に許可されていない状態、log、外部 call が変化する。 | 副作用違反。 |
| partial failure で失敗地点以降の write / call が発生する。 | 部分失敗境界違反。 |
| SDK が API response を補完し、UI が SDK を迂回し、runner / builder が未定義状態ファイルを作成する。 | component 責務違反。 |
| 実ネットワーク、実時刻、実ユーザー環境、実 secret に依存する fixture だけで合格している。 | 再現性不足。 |

<a id="sec-27-f-21"></a>
**[fixture 証跡責務 §27-F runner / security 部分失敗・再実行固定契約](fixture.md#sec-27-f-21)：**

| ケース | 固定挙動 |
|--------|----------|
| 状態保存前の validation 失敗 | 状態ファイル、JSON Lines、外部 API、外部 command、通知を実行しない。 |
| 主状態保存後の log 追記失敗 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定した対象 owner 機能契約が rollback を明記しない限り主状態は戻さず、response は `500` とし、次回 GET は保存済み主状態を返す。 |
| log 追記後の audit 失敗 | [`docs/details/security.md` 詳細本文責務 §27.44](security.md#sec-27-44) の対象 action は API では `500`、runner 内部 event では runner failure とする。保存済み主状態と先行 log は巻き戻さない。対象 action でない操作は audit を試行しない。 |
| 外部 API 送信成功後の状態保存失敗 | 外部送信の再送を自動実行しない。状態保存失敗を `500` または runner failure として記録し、再実行時は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定した対象 owner 機能契約の重複防止 key で判定する。 |
| 通知送信失敗 | build / config / security の主結果を反転しない。`.notify_pending` または [`docs/details/runner.md` 詳細本文責務 §27.32](runner.md#sec-27-32) の失敗記録だけを更新する。 |
| download / stream 中断 | サーバー側状態を成功/失敗へ変更しない。access log は [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint で中断 status 記録が定義されている場合だけ追記し、history と build log は変更しない。 |
| 再実行 no-op | 同一入力で差分がない保存 API は、[`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) の対象 endpoint に定義された no-op response を返し、状態、config log、audit log、notify log に新規差分を作らない。 |
| 破損 JSON Lines | read API は破損行を除外し、破損内容を response に出さない。write API は既存破損行を修復、削除、並べ替えしない。 |

<a id="28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約"></a>

**28-F fixture 証跡責務 / builder 拡張実装検証証跡詳細契約：**

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) は、[`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) の fixture、fake、expected、effects、実装検証証跡を扱う fixture 証跡責務である。各 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 機能は、Markdown 入力、CLI option、期待 HTML、期待 CSS / JS、`[REPORT]`、終了コード、strict / non-strict の差分を fixture で固定する。外部 library、CDN、実 network、現在時刻、実 git repository、実画像取得、画像 snapshot だけの合否判定を fixture の前提にしてはならない。

<a id="sec-28-f"></a>
**[fixture 証跡責務 §28-F 配置固定契約](fixture.md#sec-28-f)：**

| 節 | fixture 配置単位 | 必須 fixture |
|----|------------------|--------------|
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/config-resolution/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `config-resolution` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/determinism/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `determinism` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/atomicity/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `atomicity` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/parser-precedence/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `parser-precedence` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/browser-runtime/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `browser-runtime` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `builder-extensions/visual-layout/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の `visual-layout` fixture をすべて作成する。 |
| [`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) | `builder-extensions/<feature-slug>/` | 次の [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約に列挙した fixture をすべて作成する。 |

<a id="sec-28-f-2"></a>
**[fixture 証跡責務 §28-F カタログ固定契約](fixture.md#sec-28-f-2)：**

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の固定表の fixture 名は固定値である。実装変更では、対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の全 fixture を追加または更新し、`expected/stdout.txt` の `[REPORT]` key が [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の固定契約と一致することを証跡に含める。

| 節 | feature slug | 必須 fixture |
|----|--------------|--------------|
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `config-resolution` | `success-config-defaults-only`、`success-config-file-values`、`success-config-env-over-file`、`success-config-cli-over-env-over-file`、`failure-config-json-corrupt`、`failure-config-unknown-key`、`failure-config-invalid-type`、`failure-config-duplicate-nonrepeatable-cli`、`security-config-secret-not-echoed` |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `determinism` | `success-slug-duplicates`、`success-search-index-text-sources`、`success-local-storage-payload`、`success-hash-targets`、`security-deterministic-no-runtime-variance` |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `atomicity` | `success-atomic-write-all-files`、`success-incremental-reuse-byte-identical`、`success-incremental-delete-stale-page`、`failure-strict-warning-no-replace`、`failure-write-error-no-partial-update`、`failure-changed-manifest-invalid-no-output`、`success-dependency-manifest-corrupt-full-build`、`security-atomic-no-stale-temp-promoted` |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `parser-precedence` | `success-block-precedence-code-math-heading`、`success-inline-precedence-code-image-link`、`success-admonition-inline-composition`、`success-heading-inline-slug-source`、`success-list-definition-task-boundary`、`failure-unclosed-math-strict`、`noop-code-fence-protects-extensions`、`security-parser-raw-html-escaped` |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `browser-runtime` | `success-runtime-init-order`、`success-section-collapse-storage-print`、`success-runtime-light-mode-print`、`success-toc-active-observer-fallback`、`success-hash-history-focus-navigation`、`success-lightbox-focus-trap-close`、`success-keyboard-scope-skip-link`、`security-runtime-no-storage-leak` |
| [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 共通 | `visual-layout` | `success-css-output-order`、`success-responsive-320-layout`、`success-print-layout`、`success-light-mode-variables`、`success-minify-visual-preservation`、`success-component-overflow-boundaries`、`security-visual-no-external-assets`、`security-focus-visible-no-overlap` |
| [`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1) | `incremental` | `success-one-page-change`、`success-dependency-change`、`success-reuse-page-copied-to-staging`、`success-stale-page-delete-on-success`、`success-incremental-dependency-manifest-corrupt-full-build`、`failure-changed-manifest-corrupt`、`failure-changed-manifest-base-escape`、`noop-unchanged-pages-kept`、`security-incremental-no-public-write-on-failure` |
| [`docs/details/builder.md` 詳細本文責務 §28.2](builder.md#sec-28-2) | `formats` | `success-html-default`、`success-html-explicit`、`failure-pdf-reserved`、`failure-epub-reserved`、`failure-unknown-format`、`failure-multiple-format`、`security-format-no-reserved-output-files` |
| [`docs/details/builder.md` 詳細本文責務 §28.3](builder.md#sec-28-3) | `markdown-extensions` | `success-admonition-note-warn-tip`、`success-badge-color`、`success-extension-csv-normalization`、`failure-unknown-extension`、`failure-badge-invalid-text-strict`、`noop-badge-invalid-text-nonstrict`、`security-extension-escape`、`noop-extension-disabled` |
| [`docs/details/builder.md` 詳細本文責務 §28.4](builder.md#sec-28-4) | `code-line-numbers` | `success-line-numbers-fence`、`success-line-numbers-cli`、`success-line-numbers-diff-composition`、`noop-line-numbers-empty-code`、`noop-line-numbers-disabled`、`security-line-numbers-copy-clean` |
| [`docs/details/builder.md` 詳細本文責務 §28.5](builder.md#sec-28-5) | `heading-numbering` | `success-heading-numbering-h2-h3`、`success-heading-numbering-implicit-h2`、`success-heading-numbering-toc-search`、`failure-heading-numbering-unknown-mode`、`noop-heading-numbering-none`、`security-heading-slug-unchanged` |
| [`docs/details/builder.md` 詳細本文責務 §28.6](builder.md#sec-28-6) | `section-collapse` | `success-collapse-h2-h3`、`success-collapse-local-storage`、`success-collapse-print-search-hash`、`failure-collapse-duplicate-target`、`noop-collapse-no-heading`、`noop-collapse-disabled`、`security-collapse-state-parse-guard` |
| [`docs/details/builder.md` 詳細本文責務 §28.7](builder.md#sec-28-7) | `toc-depth` | `success-toc-depth-h2-h3`、`success-toc-depth-h1-h6`、`success-toc-depth-heading-numbering-sync`、`failure-toc-depth-invalid-range`、`failure-toc-depth-invalid-format`、`security-toc-depth-active-sync` |
| [`docs/details/builder.md` 詳細本文責務 §28.8](builder.md#sec-28-8) | `updated-at` | `success-updated-at-git`、`success-updated-at-file`、`success-updated-at-fallback`、`success-updated-at-none`、`failure-updated-at-unknown-source`、`failure-updated-at-unavailable`、`security-updated-at-no-search-index` |
| [`docs/details/builder.md` 詳細本文責務 §28.9](builder.md#sec-28-9) | `diff-highlight` | `success-diff-insert-delete-context`、`success-diff-header`、`success-diff-line-number-composition`、`noop-diff-non-diff-language`、`security-diff-escape`、`security-diff-copy-text-clean` |
| [`docs/details/builder.md` 詳細本文責務 §28.10](builder.md#sec-28-10) | `lazy-images` | `success-lazy-relative-image`、`success-lazy-external-image-no-fetch`、`success-lazy-data-uri-no-fetch`、`noop-lazy-base-outside-nonstrict`、`failure-lazy-base-outside-strict`、`noop-lazy-disabled`、`security-lazy-alt-escape`、`security-lazy-invalid-scheme-strict` |
| [`docs/details/builder.md` 詳細本文責務 §28.11](builder.md#sec-28-11) | `custom-meta` | `success-meta-name-property-order`、`success-meta-og-twitter`、`success-meta-duplicate-last-wins`、`failure-meta-forbidden-key`、`failure-meta-invalid-type`、`security-meta-escape`、`security-meta-secret-not-reported` |
| [`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) | `light-mode-fixed` | `success-light-mode-fixed`、`success-light-mode-print`、`security-light-mode-no-theme-toggle`、`security-light-mode-no-storage`、`failure-light-mode-dark-output` |
| [`docs/details/builder.md` 詳細本文責務 §28.13](builder.md#sec-28-13) | `code-title` | `success-code-title-colon`、`success-code-title-key-value`、`success-code-title-title-only`、`noop-code-title-empty`、`security-code-title-escape`、`security-code-title-copy-search-excluded` |
| [`docs/details/builder.md` 詳細本文責務 §28.14](builder.md#sec-28-14) | `template-vars` | `success-template-var-replace`、`success-template-var-multiple-sources`、`noop-template-var-code-fence-span`、`noop-template-var-invalid-syntax`、`failure-template-var-missing-strict`、`failure-template-var-key-validation`、`security-template-var-secret-not-reported` |
| [`docs/details/builder.md` 詳細本文責務 §28.15](builder.md#sec-28-15) | `minify-html` | `success-minify-html`、`success-minify-preserve-code`、`success-minify-attribute-order`、`failure-minify-structure-broken`、`failure-minify-marker-missing`、`noop-minify-disabled`、`security-minify-no-script-style-inline` |
| [`docs/details/builder.md` 詳細本文責務 §28.16](builder.md#sec-28-16) | `toc-active` | `success-toc-active-scroll`、`success-toc-active-fallback`、`noop-toc-active-disabled`、`security-toc-active-depth-sync` |
| [`docs/details/builder.md` 詳細本文責務 §28.17](builder.md#sec-28-17) | `mermaid` | `success-mermaid-graph-td`、`failure-mermaid-unsupported-strict`、`noop-mermaid-disabled`、`security-mermaid-no-external-script` |
| [`docs/details/builder.md` 詳細本文責務 §28.18](builder.md#sec-28-18) | `footnotes` | `success-footnotes-multiple`、`success-footnotes-backlink`、`failure-footnote-undefined-strict`、`security-footnote-escape` |
| [`docs/details/builder.md` 詳細本文責務 §28.19](builder.md#sec-28-19) | `math` | `success-math-inline-block`、`failure-math-unclosed-strict`、`noop-math-code-fence`、`security-math-escape` |
| [`docs/details/builder.md` 詳細本文責務 §28.20](builder.md#sec-28-20) | `hash-history` | `success-hash-history-click`、`success-hash-history-back-forward`、`noop-hash-history-disabled`、`security-hash-history-missing-target` |
| [`docs/details/builder.md` 詳細本文責務 §28.21](builder.md#sec-28-21) | `a11y` | `success-a11y-landmarks-labels`、`success-a11y-skip-link-tab-order`、`failure-a11y-duplicate-id-strict`、`security-a11y-no-keyboard-trap` |
| [`docs/details/builder.md` 詳細本文責務 §28.22](builder.md#sec-28-22) | `image-lightbox` | `success-lightbox-open-close`、`success-lightbox-escape-backdrop`、`failure-lightbox-alt-missing-strict`、`security-lightbox-focus-trap` |
| [`docs/details/builder.md` 詳細本文責務 §28.23](builder.md#sec-28-23) | `print-qr` | `success-print-qr-url`、`noop-print-qr-empty-url`、`failure-print-qr-url-too-long`、`security-print-qr-svg-escape` |
| [`docs/details/builder.md` 詳細本文責務 §28.24](builder.md#sec-28-24) | `definition-lists` | `success-definition-list-single`、`success-definition-list-multiple`、`noop-definition-list-empty-term`、`security-definition-list-inline-escape` |
| [`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) | `task-lists` | `success-task-list-unchecked`、`success-task-list-checked-nested`、`noop-task-list-non-target`、`security-task-list-disabled-aria` |

<a id="sec-28-f-3"></a>
**[fixture 証跡責務 §28-F ファイルセット固定契約](fixture.md#sec-28-f-3)：**

| ファイル | 必須 | 内容 |
|----------|------|------|
| `manifest.json` | 必須 | [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) に従う fixture 名、対象節、feature slug、分類、owner `builder`、collaborator、参照仕様節、fake clock、not_applicable 理由。strict と期待終了コードは含めない。 |
| `input/source.md` または `input/site/` | 必須 | Markdown 入力。site fixture は複数 Markdown、asset、dependency を含める。 |
| `input/options.json` | 必須 | CLI option、env key、expected exit code、strict / non-strict。 |
| `input/adlaire-ci-build.json` | 条件付き | 設定ファイル fixture で使用する。未使用 fixture では存在させない。 |
| `input/fakes.json` | 条件付き | [fake input root 固定契約](#fixture-fake-input-root-contract) の全 root key。builder fixture は `entropy`、`filesystem`、`git`、`mtime`、`manifest`、`cache`、`clipboard` のうち使用する root だけを非空とし、その他を空配列にする。実外部呼び出しは禁止。 |
| `input/existing-site/` | 条件付き | atomicity、incremental、failure fixture で既存公開出力を表す。成功 fixture では置換前状態、failure fixture では維持されるべき状態を置く。 |
| `expected/site/` | 必須 | 期待 HTML、`assets/style.css`、`assets/app.js`、`assets/search-index.json` のうち対象機能が変更する file。 |
| `expected/stdout.txt` | 必須 | 進捗、`[WARN]`、`[REPORT]` を含む stdout 完全一致。fatal failure は空 file。 |
| `expected/stderr.txt` | 必須 | fatal failure の `[ERROR]` 完全一致。stderr なしは空 file。 |
| `expected/effects.json` | 必須 | [fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) に従う作成、更新、維持、削除禁止 path、外部 call 0 件、既存出力保護。 |
| `expected/builder-output.json` | 必須 | staging cleanup、公開出力置換、dependency manifest 公開、search index 再生成の transaction flag。 |
| `expected/security.json` | `manifest.json.assertions` に `secret-mask` がある場合に必須 | HTML escape、attribute escape、外部 library 不使用、secret / URL credential 非表示、base 外 path 拒否。 |
| `expected/visual.json` | visual layout fixture で必須 | viewport、selector、media query、declaration、overflow、visibility、focus、禁止 asset。[fixture 証跡責務 §28-F visual layout 固定契約](#sec-28-f-10) の schema に従う。 |

<a id="sec-28-f-4"></a>
**[fixture 証跡責務 §28-F 判定粒度固定契約](fixture.md#sec-28-f-4)：**

[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、目視確認、画像 snapshot、現在時刻、実 git repository、実 network、ブラウザ環境だけに依存して合否判定してはならない。期待値は file 内容、stdout、stderr、終了コード、副作用、禁止出力のいずれかで固定する。

`input/options.json` は次の 4 key だけをすべて必須とし、未知 key を禁止する。

```json
{
  "argv": ["--strict", "--src", "input/source.md", "--out", "output"],
  "env": {},
  "strict": true,
  "expected_exit_code": 0
}
```

| key | 型 | 固定契約 |
|-----|----|----------|
| `argv` | array[string] | binary 名を含まない引数列。対象 fixture が検証する引数だけを実行順で持ち、shell 展開を適用しない。 |
| `env` | object[string,string] | 対象 [`docs/details/builder.md`](builder.md) 詳細本文責務が定義する環境変数だけを key の ASCII 昇順で持つ。未列挙の host 環境変数は実行環境へ渡さない。実 secret を禁止する。 |
| `strict` | boolean | CLI、env、config、default 解決後の有効値。`argv` / `env` / `input/adlaire-ci-build.json` から得られる値と完全一致させる。 |
| `expected_exit_code` | integer | 成功は `0`、内部エラーは `1`、入力 / 設定 / 契約エラーは `2` のいずれか。strict 固有 error に限定しない。`expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json` と矛盾させない。 |

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

builder 拡張 fixture の `expected/effects.json` は、[fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key を持つ。byte 単位で維持する path は `unchanged_paths`、新規作成、更新、削除する path はそれぞれ `created_paths`、`updated_paths`、`deleted_paths` に ASCII 昇順、重複なしで列挙する。`external_calls`、`commands`、`notifications`、`downloads`、`streams`、`read_api_calls` は空配列とし、builder が process 外副作用を行う期待値を記載してはならない。tmp、未定義 asset、reserved format 出力は `forbidden_created_paths`、failure 時に維持する既存 HTML、manifest、search index、asset は `forbidden_updated_paths` と `forbidden_deleted_paths` に列挙する。

`expected/builder-output.json` は次の 4 key だけを持ち、未知 key を禁止する。

```json
{
  "staging_cleaned": true,
  "public_output_replaced": true,
  "manifest_written": true,
  "search_index_regenerated": true
}
```

| key | 固定条件 |
|-----|----------|
| `staging_cleaned` | staging directory が残らない場合 `true`。staging cleanup 失敗 fixture だけ `false` を許可し、その場合は終了コード `1` とする。 |
| `public_output_replaced` | 成功 transaction で公開 `--out` を置換した場合だけ `true`。failure、strict failure、no-op は `false`。 |
| `manifest_written` | 成功 transaction で `.dependency_manifest.json` を公開した場合だけ `true`。 |
| `search_index_regenerated` | 成功 transaction で最終 page set から search index を再生成した場合だけ `true`。 |

<a id="sec-28-f-5"></a>
**[fixture 証跡責務 §28-F 設定解決固定契約](fixture.md#sec-28-f-5)：**

`builder-extensions/config-resolution/` は [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 全機能の前提 fixture である。[`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) の個別 fixture は、この共通 fixture と矛盾する CLI / env / config / default 解決をしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-config-defaults-only` | CLI、環境変数、`adlaire-ci-build.json` がない場合、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の全設定 key が既定値から採用され、`config_default_keys` に全 key、その他 source key に空 array が出る。 |
| `success-config-file-values` | `input/adlaire-ci-build.json` の `builder_extensions` 値が採用され、`config_file_used=true`、`config_file_path="adlaire-ci-build.json"`、`config_file_keys` が ASCII 昇順で出る。 |
| `success-config-env-over-file` | 同一 key が環境変数と設定ファイルに存在する場合、環境変数を採用し、設定ファイル source は `config_overridden_keys` に記録する。 |
| `success-config-cli-over-env-over-file` | 同一 key が CLI、環境変数、設定ファイルに存在する場合、CLI を採用し、環境変数 / 設定ファイル source は `config_overridden_keys` に記録する。 |
| `failure-config-json-corrupt` | 設定ファイルが JSON として parse できない場合、終了コード `2`、stdout 空、stderr `BUILDER28_INVALID_OPTION`、出力作成なし、既存出力維持。 |
| `failure-config-unknown-key` | root unknown key、`builder_extensions` unknown key のいずれも終了コード `2`、stdout 空、stderr に key 名だけを記録し、`[REPORT]` は出力しない。 |
| `failure-config-invalid-type` | boolean / string / array / object の型不一致、JSON object 値が string 以外、空 array 要素を終了コード `2`、stdout 空にする。 |
| `failure-config-duplicate-nonrepeatable-cli` | repeatable ではない CLI option の重複指定を終了コード `2`、stdout 空にし、Markdown 読込前に停止する。 |
| `security-config-secret-not-echoed` | `meta` / `template_vars` / 環境変数値に secret 風文字列、credential 付き URL、raw HTML が含まれても、stderr、stdout、REPORT へ値を出さない。key 名だけを出す。 |

設定解決成功 fixture の `expected/stdout.txt` は、`config_file_used`、`config_file_path`、`config_cli_keys`、`config_env_keys`、`config_file_keys`、`config_default_keys`、`config_overridden_keys`、`config_rejected_keys` を必ず含める。設定解決失敗 fixture の `expected/stdout.txt` は空 file とし、`expected/stderr.txt` に `[ERROR] BUILDER28_INVALID_OPTION -:0 28 ...` を固定する。設定解決失敗 fixture の `expected/effects.json` は、HTML、CSS、JS、search index、manifest、設定ファイルが新規作成、更新、削除されないことを固定する。

<a id="sec-28-f-6"></a>
**[fixture 証跡責務 §28-F 決定性固定契約](fixture.md#sec-28-f-6)：**

`builder-extensions/determinism/` は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の HTML identity、slug、search index、hash、localStorage の共通 fixture である。個別 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、本 fixture と異なる slug、id、search text、storage key、hash target を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-slug-duplicates` | ASCII、非 ASCII、記号、空 heading、同名 heading、`foo` と `foo-2` の衝突を含む入力で、heading id、TOC href、section wrapper id、hash target が [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の slug 規則と一致する。 |
| `success-search-index-text-sources` | 採番 heading、paragraph、list、code、admonition、badge、image alt、footnote、definition list、task list を含む入力で、search index に含める text と除外する UI text が完全一致する。 |
| `success-local-storage-payload` | section collapse を有効にし、`adlaire:section-state` の key、payload、未知値無視、JSON parse failure fallback を `expected/site/assets/app.js` と `expected/effects.json` で固定する。color scheme 用 localStorage key は存在しないことを固定する。 |
| `success-hash-targets` | hash history と TOC active tracking を有効にし、heading id だけを target にすること、TOC depth 外 heading を active 化しないこと、存在しない hash を no-op にすることを固定する。 |
| `security-deterministic-no-runtime-variance` | 同一入力を fake clock、fake git、異なる OS path separator 相当入力、異なる map order 相当 config で実行しても、HTML、search index、stdout、stderr が同一になることを固定する。 |

決定性 fixture の `expected/site/*.html` は、heading id、TOC href、collapse wrapper id、`aria-controls`、`data-section-id` を完全一致で確認する。`expected/site/assets/search-index.json` は page key、heading text、body text、除外 text の不在を JSON parse 後完全一致で確認する。`expected/site/assets/app.js` は localStorage key、payload schema、unknown value guard、parse failure guard、hash no-op guard を文字列または構造で確認する。

<a id="sec-28-f-7"></a>
**[fixture 証跡責務 §28-F atomicity 固定契約](fixture.md#sec-28-f-7)：**

`builder-extensions/atomicity/` は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の出力副作用、公開置換、manifest、search index、stale 削除を固定する共通 fixture である。個別 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、本 fixture と異なる失敗時副作用、manifest 更新条件、search index 更新条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-atomic-write-all-files` | HTML、CSS、JS、search index、`.dependency_manifest.json` が staging にそろってから公開 `--out` へ置換され、`expected/builder-output.json` の `public_output_replaced=true`、`manifest_written=true`、`search_index_regenerated=true` になる。 |
| `success-incremental-reuse-byte-identical` | 未変更 page の既存 HTML が byte 単位で維持され、changed page、manifest、search index だけが成功 transaction として更新される。 |
| `success-incremental-delete-stale-page` | 入力 source から削除された Markdown に対応する HTML、search index entry、manifest entry が成功時だけ削除される。 |
| `failure-strict-warning-no-replace` | non-strict なら fallback 出力できる警告を strict で実行し、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力、manifest、search index 維持を固定する。 |
| `failure-write-error-no-partial-update` | staging 書き込みまたは validation 失敗を fake し、終了コード `1`、stderr `[ERROR] BUILDER28_INTERNAL_IO` または `BUILDER28_OUTPUT_VALIDATION_FAILED`、公開出力、manifest、search index 維持を固定する。 |
| `failure-changed-manifest-invalid-no-output` | `--changed-manifest` が base 外 path、絶対 path、URL scheme、JSON 破損のいずれかの場合、終了コード `2`、stdout 空、公開出力維持を固定する。 |
| `success-dependency-manifest-corrupt-full-build` | 既存 `.dependency_manifest.json` が破損または schema 不一致の場合、warning なし full build とし、成功時だけ新 manifest と search index を公開する。 |
| `security-atomic-no-stale-temp-promoted` | staging path、absolute path、host user path、tmp path が HTML、CSS、JS、search index、manifest、stdout、stderr、REPORT に混入しないことを固定する。 |

atomicity fixture の `input/existing-site/` は、既存 HTML、既存 `assets/search-index.json`、既存 `.dependency_manifest.json`、stale HTML、既存 asset を含める。failure fixture の `expected/site/` は `input/existing-site/` と byte 単位で一致させる。success fixture の `expected/effects.json` は、`created_paths`、`updated_paths`、`unchanged_paths`、`deleted_paths` をすべて明示する。

<a id="sec-28-f-8"></a>
**[fixture 証跡責務 §28-F parser precedence 固定契約](fixture.md#sec-28-f-8)：**

`builder-extensions/parser-precedence/` は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の Markdown parser 優先順位、構文 grammar、曖昧構文、機能併用順を固定する共通 fixture である。個別 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、本 fixture と異なる token 解釈、別順序の inline 変換、code fence / code span 内変換を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-block-precedence-code-math-heading` | code fence 継続中の heading / footnote / badge / math が code text のまま残り、math block と heading の判定順が [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) Markdown parser 優先順位固定契約と一致する。 |
| `success-inline-precedence-code-image-link` | code span 内の badge / footnote / math が変換されず、image が link より優先され、link text 内の badge / math だけが inline 変換される。 |
| `success-admonition-inline-composition` | admonition body 内の badge、footnote、math、fenced code の併用で、body inline 変換と fenced code 保護が両立する。 |
| `success-heading-inline-slug-source` | heading 内の badge、footnote、math 表示変換と、slug source text から UI text を除外する規則が同時に成立する。 |
| `success-list-definition-task-boundary` | task list、通常 list、definition list、list 内 `: definition` の境界が固定どおりに分かれる。 |
| `failure-unclosed-math-strict` | 未閉鎖 math inline / math block が non-strict では通常 text、strict では `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持になる。 |
| `noop-code-fence-protects-extensions` | code fence 内の template var、badge、footnote、math、definition marker、task marker が一切変換されない。 |
| `security-parser-raw-html-escaped` | raw HTML、event handler、`javascript:` URL、HTML comment 指示が parser 段階で実行可能要素にならず、expected HTML と security.json で escape を確認する。 |

parser precedence fixture の `expected/site/*.html` は、対象 token の tag、text node、未変換 text、変換済み node、属性順を完全一致で確認する。`expected/stdout.txt` は warning の有無、warning code、line、section を完全一致で確認する。`expected/security.json` は raw HTML、script、event handler、credential URL、CDN、外部 library が出力に存在しないことを固定する。

<a id="sec-28-f-9"></a>
**[fixture 証跡責務 §28-F browser runtime 固定契約](fixture.md#sec-28-f-9)：**

`builder-extensions/browser-runtime/` は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) のブラウザ JS 初期化順、状態復元、event handler、focus、keyboard、print、fallback を固定する共通 fixture である。個別 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、本 fixture と異なる localStorage key、focus 移動、keyboard scope、lightbox close 条件、TOC active 条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-runtime-init-order` | `assets/app.js` 内で static guard、storage guard、section collapse、hash history、TOC active、lightbox、accessibility guard の初期化順が固定どおりである。 |
| `success-section-collapse-storage-print` | section collapse の既定展開、保存値復元、toggle、`aria-expanded`、`adlaire-section-collapsed`、search hit 一時展開、beforeprint / afterprint 復元が一致する。 |
| `success-runtime-light-mode-print` | light 固定の CSS variables、theme toggle 不在、color scheme 永続化不在、print light が一致する。 |
| `success-toc-active-observer-fallback` | IntersectionObserver 使用時と fallback scroll 時の active link 1 件化、`.is-active`、`aria-current="location"`、TOC depth 外除外が一致する。 |
| `success-hash-history-focus-navigation` | heading / TOC click、`history.pushState`、`tabindex="-1"`、focus、back / forward、missing hash no-op が一致する。 |
| `success-lightbox-focus-trap-close` | trigger click、`Enter` / `Space`、dialog open、Escape、backdrop、close button、opener focus return、Tab / Shift+Tab focus trap が一致する。 |
| `success-keyboard-scope-skip-link` | [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) keyboard handler が対象 UI focus 中だけ有効で、[`docs/details/builder.md` 詳細本文責務 §7.12](builder.md#sec-7-12) shortcut を上書きせず、skip link が main content へ移動する。 |
| `security-runtime-no-storage-leak` | cookie、sessionStorage、IndexedDB、runtime network fetch、external script、secret / credential の storage 書込が 0 件である。 |

browser runtime fixture の `expected/site/assets/app.js` は、初期化関数名または固定 marker、localStorage key、event 名、guard、fallback 分岐、focus trap 分岐を文字列または構造で確認する。`expected/security.json` は、cookie、sessionStorage、IndexedDB、fetch、XMLHttpRequest、external script、secret / credential storage が存在しないことを固定する。ブラウザ実行がない fixture でも、期待 JS 構造と expected HTML / CSS / security を組み合わせて合否判定する。

<a id="sec-28-f-10"></a>
**[fixture 証跡責務 §28-F visual layout 固定契約](fixture.md#sec-28-f-10)：**

`builder-extensions/visual-layout/` は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の CSS 出力順、responsive layout、print layout、light 固定変数、minify 後の視覚維持、component overflow、外部 visual asset 禁止、focus 表示を固定する共通 fixture である。個別 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、本 fixture と異なる selector 順、media query 境界、print visibility、外部 asset 許可、focus 表示条件を期待値にしてはならない。

| fixture | 固定する内容 |
|---------|--------------|
| `success-css-output-order` | `expected/site/assets/style.css` で、既存 base、light visual baseline、typography / block、code extension、navigation runtime UI、media UI、responsive、print の順序が固定どおりである。 |
| `success-responsive-320-layout` | 幅 `320px` 相当の fixture metadata と expected CSS / HTML で、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) UI が text overlap、text clipping、不可視 overflow を発生させず、table と code block だけが scroll wrapper 内で横 overflow を持つ。 |
| `success-print-layout` | `@media print` で interactive controls を非表示、本文要素を表示、collapsed section を展開、light 固定表示、QR を print 専用表示にする。 |
| `success-light-mode-variables` | `:root` と print の custom property 名と既定値が固定どおりであり、dark / auto selector が存在しない。 |
| `success-minify-visual-preservation` | minify 有効時も required selector、custom property、[`docs/DESIGN.md` デザイン責務 Builder 拡張コンポーネント視覚契約](../DESIGN.md#builder-拡張コンポーネント視覚契約) が定義する responsive breakpoint の media query、typography stability、focus / active 時の layout 寸法維持、`@media print`、`pre` / `code` の空白保持 property が削除、改名、結合破壊されない。 |
| `success-component-overflow-boundaries` | admonition、badge、definition list、task list、footnote、math、code title、line numbers、diff、lightbox、Mermaid、print QR が [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) CSS / layout / print / visual 固定契約の境界どおりに出力される。 |
| `security-visual-no-external-assets` | `@import`、remote `url()`、external font、CDN、追加 asset file、runtime CSS fetch、inline style が出力されない。 |
| `security-focus-visible-no-overlap` | skip link、TOC active、hash target、lightbox control、collapse toggle の focus / active style が可視で、hover / focus により寸法が変わらず、text を隠さない。 |

visual layout fixture の viewport 条件は `expected/visual.json.viewport_width` に固定し、`manifest.json` に追加 key を置かない。判定を画像 snapshot だけに依存させてはならない。visual layout fixture は `expected/site/assets/style.css`、`expected/site/*.html`、`expected/security.json`、`expected/visual.json` を必須とし、selector、media query、declaration、overflow、visibility、focus、禁止 asset を構造化して固定する。ブラウザ実行がない fixture でも、期待 HTML / CSS / security / visual の組み合わせで合否判定できなければならない。

`expected/visual.json` は次の root key だけをすべて必須とし、未知 key を禁止する。string array は完全文字列の ASCII 昇順、object array は各行が定める複合 key 順の ASCII 昇順とし、いずれも重複なしとする。該当なしは空配列とする。

| key | 型 | 固定契約 |
|-----|----|----------|
| `viewport_width` | integer/null | viewport 固定 fixture は `320` 以上の正整数、viewport 非依存 fixture は `null`。単位は CSS pixel。 |
| `required_selectors` | array[string] | 期待 CSS と HTML の両方で照合する selector。空文字を禁止する。 |
| `required_media_queries` | array[string] | `expected/site/assets/style.css` に必要な media query 文字列の完全一致値。 |
| `required_declarations` | array[object] | 各 object は `selector`、`property`、`value` の 3 key だけを持ち、期待 CSS の正規化後宣言と一致する。順序は `selector`、`property`、`value` の複合 key とする。 |
| `overflow_expectations` | array[object] | 各 object は `selector`、`axis`、`mode` の 3 key だけを持つ。`axis` は `x`、`y`、`both`、`mode` は `visible`、`clip`、`scroll`、`auto` のいずれかとする。 |
| `visibility_expectations` | array[object] | 各 object は `selector`、`context`、`visible` の 3 key だけを持つ。`context` は `screen`、`print`、`focus`、`active` のいずれか、`visible` は boolean とする。 |
| `focus_expectations` | array[object] | 各 object は `selector`、`indicator_visible`、`layout_shift` の 3 key だけを持ち、後ろ 2 key は boolean とする。 |
| `forbidden_assets` | array[string] | 出力が含んではならない external URL、`@import`、remote font、CDN、追加 asset path の固定文字列。 |

<a id="sec-28-f-11"></a>
**[fixture 証跡責務 §28-F builder 詳細本文責務 §28.1〜§28.5 fixture 固定契約](fixture.md#sec-28-f-11)：**

[`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) の fixture は、対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の実装詳細固定契約に列挙された HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。各 fixture は `manifest.json.section` を対象 `§28.x`、`manifest.json.feature` を [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) カタログ固定契約の feature slug と一致させる。各ブロックの `expected/effects.json` は [fixture 証跡責務 §27-F expected/effects.json schema 固定契約](#sec-27-f-11) の全 root key、`expected/builder-output.json` は本節の 4 transaction flag を持つ。

[`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1)〜[`docs/details/builder.md` 詳細本文責務 §28.5](builder.md#sec-28-5) の fixture は、中間状態、HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。

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
| `markdown-extensions` | `failure-badge-invalid-text-strict` | strict で不正 badge label / color を `BUILDER28_INVALID_CONTENT`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `markdown-extensions` | `noop-badge-invalid-text-nonstrict` | non-strict で不正 badge を元 text のまま出力し、`markdown_extension_warnings` を加算する。 |
| `markdown-extensions` | `security-extension-escape` | admonition body、badge label、attribute、raw HTML、危険 URL、event handler が escape される。 |
| `markdown-extensions` | `noop-extension-disabled` | extension 未指定時に admonition / badge 構文を特別扱いせず、既存 Markdown 変換結果を維持する。 |
| `code-line-numbers` | `success-line-numbers-fence` | fence option `line-numbers` のある code block だけに `.code-lines`、`.line-no`、`data-line` を出す。 |
| `code-line-numbers` | `success-line-numbers-cli` | CLI 有効時に全 code fence へ line number を出し、`code_line_number_blocks` と `code_line_number_lines` が一致する。 |
| `code-line-numbers` | `success-line-numbers-diff-composition` | [`docs/details/builder.md` 詳細本文責務 §28.9](builder.md#sec-28-9) diff class と line number が同時に存在し、line number node に diff class が付かない。 |
| `code-line-numbers` | `noop-line-numbers-empty-code` | 空 code block には line number を出さず、REPORT count に含めない。 |
| `code-line-numbers` | `noop-line-numbers-disabled` | CLI 無効かつ fence option なしの code block は既存出力と一致する。 |
| `code-line-numbers` | `security-line-numbers-copy-clean` | `expected/site/assets/search-index.json`、copy text fixture、minify 後 HTML のいずれにも line number text が混入しない。 |
| `heading-numbering` | `success-heading-numbering-h2-h3` | h2 / h3 に `span.heading-number` を出し、`1.`、`1.1.` 形式、ASCII space 1 個、`numbered_headings` を固定する。 |
| `heading-numbering` | `success-heading-numbering-implicit-h2` | h2 がない page の h3 で暗黙 h2 counter `1` を使い、`1.1.` から開始する。 |
| `heading-numbering` | `success-heading-numbering-toc-search` | 本文 heading、TOC 表示 text、search index 表示 text に番号を含め、search index 検索対象正規化 text には番号を含めない。 |
| `heading-numbering` | `failure-heading-numbering-unknown-mode` | 未知 mode を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `heading-numbering` | `noop-heading-numbering-none` | `none` で heading、TOC、search index 表示 text を変更せず、`numbered_headings=0` にする。 |
| `heading-numbering` | `security-heading-slug-unchanged` | 採番有無で heading id、anchor href、collapse target、hash history target が byte 単位で一致する。 |

[§28.1〜§28.5](builder.md#sec-28-group-1-5) の failure / security fixture では、`forbidden_updated_paths` と `forbidden_deleted_paths` に公開 `--out`、既存 `.dependency_manifest.json`、既存 `assets/search-index.json` を必ず含める。

<a id="sec-28-f-12"></a>
**[fixture 証跡責務 §28-F builder 詳細本文責務 §28.6〜§28.10 fixture 固定契約](fixture.md#sec-28-f-12)：**

[`docs/details/builder.md` 詳細本文責務 §28.6](builder.md#sec-28-6)〜[`docs/details/builder.md` 詳細本文責務 §28.10](builder.md#sec-28-10) の fixture は、UI 状態、TOC、timestamp、code token、image token、HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。manifest と共通 effects key は [§28.1〜§28.5](builder.md#sec-28-group-1-5) の共通契約に従う。

| feature slug | fixture | 固定する内容 |
|--------------|---------|--------------|
| `section-collapse` | `success-collapse-h2-h3` | h2 / h3 の section 範囲、`button.adlaire-section-toggle`、`aria-controls`、`aria-expanded`、`data-section-id`、`section-<slug>` wrapper、`collapsible_sections` を固定する。 |
| `section-collapse` | `success-collapse-local-storage` | `adlaire:section-state` の `<page_key>#<slug>` key、boolean payload、ASCII key order、unknown key 無視、JSON parse failure fallback を `expected/site/assets/app.js` で確認する。 |
| `section-collapse` | `success-collapse-print-search-hash` | print 全展開、検索 hit 一時展開、hash target 一時展開、localStorage 保存値非変更を HTML / CSS / JS expected で確認する。 |
| `section-collapse` | `failure-collapse-duplicate-target` | `section-<slug>` wrapper id が既存 id と衝突した場合に `BUILDER28_OUTPUT_VALIDATION_FAILED`、終了コード `1`、公開出力維持になる。 |
| `section-collapse` | `noop-collapse-no-heading` | h2 / h3 がない page では toggle、wrapper、JS state、REPORT count を増やさない。 |
| `section-collapse` | `noop-collapse-disabled` | option 無効時に toggle、wrapper、collapse JS、localStorage key を出力せず、既存 heading HTML と一致する。 |
| `section-collapse` | `security-collapse-state-parse-guard` | localStorage に JSON 破損、boolean 以外、未知 page / section key があっても例外化せず、静的 HTML、TOC、本文を壊さない。 |
| `toc-depth` | `success-toc-depth-h2-h3` | `--toc-depth 2:3` で TOC link が h2 / h3 だけになり、本文 heading、heading id、search index heading source が変化しない。 |
| `toc-depth` | `success-toc-depth-h1-h6` | `--toc-depth 1:6` で全 heading level の TOC link を本文出現順に出力する。 |
| `toc-depth` | `success-toc-depth-heading-numbering-sync` | [`docs/details/builder.md` 詳細本文責務 §28.5](builder.md#sec-28-5) と併用し、TOC 表示 text だけに numbering を含め、href と heading id が採番で変わらない。 |
| `toc-depth` | `failure-toc-depth-invalid-range` | `0:6`、`1:7`、`4:2` を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `toc-depth` | `failure-toc-depth-invalid-format` | 空値、整数以外、separator 不一致、余分な値を `BUILDER28_INVALID_OPTION`、終了コード `2` にする。 |
| `toc-depth` | `security-toc-depth-active-sync` | [`docs/details/builder.md` 詳細本文責務 §28.16](builder.md#sec-28-16) 有効時の active tracking 対象が TOC 出力 link と一致し、depth 外 heading を active 化しない。 |
| `updated-at` | `success-updated-at-git` | fake git timestamp を UTC RFC3339 秒精度へ正規化し、`time.page-updated-at`、表示 text、`updated_at_source="git"`、`updated_at_fallback=0` を固定する。 |
| `updated-at` | `success-updated-at-file` | fake file mtime を UTC RFC3339 秒精度へ正規化し、`updated_at_source="file"`、`updated_at_fallback=0` を固定する。 |
| `updated-at` | `success-updated-at-fallback` | git 取得不能かつ file mtime 取得可能時に file へ fallback し、`updated_at_source="file"`、`updated_at_fallback=1` になる。 |
| `updated-at` | `success-updated-at-none` | `none` で timestamp 取得なし、`.page-updated-at` 出力なし、`updated_at=""`、`updated_at_source="none"`、`updated_at_fallback=0` になる。 |
| `updated-at` | `failure-updated-at-unknown-source` | 未知 source を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `updated-at` | `failure-updated-at-unavailable` | git / file timestamp とも取得不能、または timestamp parse 不能を `BUILDER28_INTERNAL_IO`、終了コード `1`、公開出力維持にする。 |
| `updated-at` | `security-updated-at-no-search-index` | search index に `.page-updated-at` の label、timestamp、UI text が混入しない。 |
| `diff-highlight` | `success-diff-insert-delete-context` | diff / patch fence の inserted、deleted、context 行へ `.tok-inserted`、`.tok-deleted`、`.tok-context` を付与し、REPORT count を固定する。 |
| `diff-highlight` | `success-diff-header` | `+++` / `---` header 行を `.tok-diff-header` とし、insertions / deletions に加算しない。 |
| `diff-highlight` | `success-diff-line-number-composition` | [`docs/details/builder.md` 詳細本文責務 §28.4](builder.md#sec-28-4) と併用し、line number node に diff class が付かず、code text 側だけに diff class が付く。 |
| `diff-highlight` | `noop-diff-non-diff-language` | 通常 code fence では行頭 `+` / `-` / space があっても diff class を付けない。 |
| `diff-highlight` | `security-diff-escape` | diff 行内の raw HTML、event handler、`javascript:` URL が escape され、class 付与後も実行可能にならない。 |
| `diff-highlight` | `security-diff-copy-text-clean` | copy text と search index に diff class、line number、UI label が混入せず、元の diff 記号と code text だけを含む。 |
| `lazy-images` | `success-lazy-relative-image` | base 内相対 image path を正規化し、`loading="lazy"`、`decoding="async"`、escaped alt を出力する。 |
| `lazy-images` | `success-lazy-external-image-no-fetch` | `http` / `https` URL に lazy 属性を付けるが、external call は 0 件である。 |
| `lazy-images` | `success-lazy-data-uri-no-fetch` | `data:` URL に lazy 属性を付けるが、decode、MIME 判定、external call を行わない。 |
| `lazy-images` | `noop-lazy-base-outside-nonstrict` | non-strict で base 外相対 path を `BUILDER28_PATH_OUTSIDE_BASE` warning、exit `0` にし、該当 image は alt text だけまたは空出力、`lazy_images` 加算なし、拒否 URL 値の HTML / search index / stdout / stderr 不在にする。 |
| `lazy-images` | `failure-lazy-base-outside-strict` | strict で base 外相対 path を `BUILDER28_PATH_OUTSIDE_BASE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `lazy-images` | `noop-lazy-disabled` | option 無効時に `loading`、`decoding` を追加せず、既存 img 出力と一致する。 |
| `lazy-images` | `security-lazy-alt-escape` | alt、src、title 相当の attribute に raw HTML、quote、event handler が混入しても attribute escape される。 |
| `lazy-images` | `security-lazy-invalid-scheme-strict` | `javascript:`、`file:`、その他未許可 scheme を strict で `[WARN] UNSAFE_URL`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |

[§28.6〜§28.10](builder.md#sec-28-group-6-10) で browser runtime、visual layout、parser precedence と併用する fixture では、該当共通 fixture と同じ localStorage key、media query、parser 保護、external call 0 件を再確認する。

<a id="sec-28-f-13"></a>
**[fixture 証跡責務 §28-F builder 詳細本文責務 §28.11〜§28.15 fixture 固定契約](fixture.md#sec-28-f-13)：**

[`docs/details/builder.md` 詳細本文責務 §28.11](builder.md#sec-28-11)〜[`docs/details/builder.md` 詳細本文責務 §28.15](builder.md#sec-28-15) の fixture は、head meta、theme state、code title、template var、minify byte、HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。manifest と共通 effects key は [§28.1〜§28.5](builder.md#sec-28-group-1-5) の共通契約に従う。

| feature slug | fixture | 固定する内容 |
|--------------|---------|--------------|
| `custom-meta` | `success-meta-name-property-order` | `name:*`、bare key、`property:og:*`、`property:twitter:*` の正規化、ASCII key order、head 内の既存 meta 後 / stylesheet 前の出力順を固定する。 |
| `custom-meta` | `success-meta-og-twitter` | OGP と Twitter meta を `property` attribute で出力し、`content` attribute を escape 済みで固定する。 |
| `custom-meta` | `success-meta-duplicate-last-wins` | CLI、env、config、同一 source 内重複の last wins と、`custom_meta_count` / `custom_meta_rejected` を固定する。 |
| `custom-meta` | `failure-meta-forbidden-key` | `script`、`http-equiv`、`charset`、`refresh`、`set-cookie`、`content-security-policy` を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空にする。 |
| `custom-meta` | `failure-meta-invalid-type` | `ADLAIRE_META_JSON` または config meta が object 以外、value string 以外、空 key、制御文字 key の場合に `BUILDER28_INVALID_OPTION` になる。 |
| `custom-meta` | `security-meta-escape` | meta key / value の quote、raw HTML、event handler、credential URL を attribute escape し、実行可能 HTML を出力しない。 |
| `custom-meta` | `security-meta-secret-not-reported` | secret 風 value と credential 付き URL value が stdout、stderr、REPORT、manifest に平文出力されない。 |
| `light-mode-fixed` | `success-light-mode-fixed` | `:root` の light CSS variables、REPORT `color_scheme_fixed=true`、dark / auto selector 不在を固定する。 |
| `light-mode-fixed` | `success-light-mode-print` | `@media print` で light 固定の背景 / 文字色になり、dark background を印刷しない。 |
| `light-mode-fixed` | `security-light-mode-no-theme-toggle` | [`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) の theme toggle 禁止識別子が HTML / CSS / JS に存在しない。 |
| `light-mode-fixed` | `security-light-mode-no-storage` | [`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) の入力、storage、REPORT 禁止識別子が存在しない。 |
| `light-mode-fixed` | `failure-light-mode-dark-output` | [`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) の dark / auto / 永続化禁止識別子が出力された場合は `BUILDER28_OUTPUT_VALIDATION_FAILED`、終了コード `1`、公開出力維持にする。 |
| `code-title` | `success-code-title-colon` | `go:main.go` 形式で language と title を分離し、`.code-block-header` 内 `.code-title` を出力する。 |
| `code-title` | `success-code-title-key-value` | `bash:title=deploy.sh` 形式で title を出力し、language は `bash` として code block に残す。 |
| `code-title` | `success-code-title-title-only` | `title=README.md` 形式で language 空、title ありの code block を固定する。 |
| `code-title` | `noop-code-title-empty` | 空 title、空白 title、`title=`、`lang:` の値なしを no-op にし、warning と REPORT count を増やさない。 |
| `code-title` | `security-code-title-escape` | title 内 raw HTML、quote、event handler、`javascript:` URL が escape される。 |
| `code-title` | `security-code-title-copy-search-excluded` | copy text、search index、line number count、diff count に code title text が混入しない。 |
| `template-vars` | `success-template-var-replace` | `{{ KEY }}` を Markdown parse 前に通常 text だけ置換し、置換後 text が Markdown 処理へ渡る。 |
| `template-vars` | `success-template-var-multiple-sources` | CLI、env、config の source 優先順位、repeatable CLI、key count、replacement count、missing array を固定する。 |
| `template-vars` | `noop-template-var-code-fence-span` | code fence と code span 内の `{{ KEY }}` が置換されない。 |
| `template-vars` | `noop-template-var-invalid-syntax` | `{{KEY}}`、`{{ key }}`、<code>{{ KEY &#124; filter }}</code> が通常 text として残る。 |
| `template-vars` | `failure-template-var-missing-strict` | strict で未定義 var を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `template-vars` | `failure-template-var-key-validation` | key 不正、object 以外、value string 以外を `BUILDER28_INVALID_OPTION`、終了コード `2` にする。 |
| `template-vars` | `security-template-var-secret-not-reported` | secret 風 value と credential URL value が stdout、stderr、REPORT、manifest に平文出力されない。 |
| `minify-html` | `success-minify-html` | tag 間 whitespace と HTML comment の安全な削減、byte before / after / saved の REPORT を固定する。 |
| `minify-html` | `success-minify-preserve-code` | `pre` / `code` 内 whitespace、改行、escape 済み text が byte 単位で保持される。 |
| `minify-html` | `success-minify-attribute-order` | attribute order、quote、escape、URL、data / aria attribute が minify 前後で保持される。 |
| `minify-html` | `failure-minify-structure-broken` | minify 後に doctype / html / head / body、必須 id / class / attribute、search index 対象 text が壊れる場合に `BUILDER28_OUTPUT_VALIDATION_FAILED`、終了コード `1` になる。 |
| `minify-html` | `failure-minify-marker-missing` | 必須 marker が定義されている fixture で marker 消失を `BUILDER28_OUTPUT_VALIDATION_FAILED`、終了コード `1`、公開出力維持にする。 |
| `minify-html` | `noop-minify-disabled` | minify 無効時に HTML byte を変更せず、minify REPORT byte count を 0 にする。 |
| `minify-html` | `security-minify-no-script-style-inline` | minify 実装が新規 inline script / style を追加せず、既存 script / style 相当領域の内部 byte を変更しない。 |

[§28.11〜§28.15](builder.md#sec-28-group-11-15) の security fixture では secret / credential が stdout、stderr、REPORT、manifest、HTML attribute、search index のいずれにも平文で残らないことを `expected/security.json` に固定する。

<a id="sec-28-f-14"></a>
**[fixture 証跡責務 §28-F builder 詳細本文責務 §28.16〜§28.20 fixture 固定契約](fixture.md#sec-28-f-14)：**

[`docs/details/builder.md` 詳細本文責務 §28.16](builder.md#sec-28-16)〜[`docs/details/builder.md` 詳細本文責務 §28.20](builder.md#sec-28-20) の fixture は、TOC active、Mermaid、footnote、math、hash history の HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。manifest と共通 effects key は [§28.1〜§28.5](builder.md#sec-28-group-1-5) の共通契約に従う。

| feature slug | fixture | 固定する内容 |
|--------------|---------|--------------|
| `toc-active` | `success-toc-active-scroll` | TOC link 集合、初期 active、scroll 時の `.is-active` 1 件化、`aria-current="location"` の付与 / 削除、REPORT `toc_active_tracking=true`、`toc_active_items` を固定する。 |
| `toc-active` | `success-toc-active-fallback` | IntersectionObserver が使えない前提で fallback scroll handler が同じ active 候補集合を使い、例外時 no-break になることを `expected/site/assets/app.js` で確認する。 |
| `toc-active` | `noop-toc-active-disabled` | `--toc-active=false` で active handler、`.is-active` 初期 class、`aria-current`、未使用 JS branch を出力せず、既存 TOC HTML と一致する。 |
| `toc-active` | `security-toc-active-depth-sync` | [`docs/details/builder.md` 詳細本文責務 §28.7](builder.md#sec-28-7) と併用し、TOC depth 外 heading、footnote backlink、collapse wrapper、lightbox target を active 対象にしない。 |
| `mermaid` | `success-mermaid-graph-td` | `graph TD`、node 定義、edge 定義を deterministic SVG へ変換し、`.mermaid-diagram`、`.mermaid-node`、`.mermaid-edge`、viewBox、REPORT rendered count を固定する。 |
| `mermaid` | `failure-mermaid-unsupported-strict` | strict で未対応 Mermaid 構文を `BUILDER28_UNSUPPORTED_RESERVED`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `mermaid` | `noop-mermaid-disabled` | `--mermaid=false` で `mermaid` fence を通常 code block として出力し、SVG、Mermaid class、external script、REPORT rendered count を増やさない。 |
| `mermaid` | `security-mermaid-no-external-script` | SVG 内に `script`、`foreignObject`、event handler、external href、CDN、runtime fetch が存在しないことを `expected/security.json` で固定する。 |
| `footnotes` | `success-footnotes-multiple` | 複数 definition / reference、同一 id 複数参照、参照順番号、`sup.footnote-ref`、末尾 `section.footnotes`、REPORT count を固定する。 |
| `footnotes` | `success-footnotes-backlink` | 各 footnote item の `.footnote-backref`、本文 reference への backlink target、一意 id、search index から UI label を除外することを固定する。 |
| `footnotes` | `failure-footnote-undefined-strict` | strict で未定義 reference、未参照 definition、重複 definition を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `footnotes` | `security-footnote-escape` | definition text、reference 周辺 text、id、backlink label の raw HTML、quote、event handler、`javascript:` が escape される。 |
| `math` | `success-math-inline-block` | `$...$` と `$$...$$` を `span.math-inline` / `div.math-block` へ変換し、delimiter 除去、escape、REPORT inline / block count を固定する。 |
| `math` | `failure-math-unclosed-strict` | strict で未閉鎖 inline delimiter、未閉鎖 block delimiter、長さ超過を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `math` | `noop-math-code-fence` | code fence、code span、link destination、image src、escaped dollar、通貨表現では math 変換しない。 |
| `math` | `security-math-escape` | math content 内 raw HTML、script 風 text、event handler 風 text が escape 済み text として残り、外部 renderer / SVG / canvas / image を出力しない。 |
| `hash-history` | `success-hash-history-click` | heading / TOC link click、`history.pushState`、`tabindex="-1"`、focus、scroll、REPORT `hash_history_enabled=true` / target count を固定する。 |
| `hash-history` | `success-hash-history-back-forward` | popstate / hashchange で back / forward 時に既存 heading target へ focus し、TOC active と衝突しない handler 順を固定する。 |
| `hash-history` | `noop-hash-history-disabled` | `--hash-history=false` で pushState handler、popstate handler、tabindex 補助を出力せず、通常 anchor fallback だけを残す。 |
| `hash-history` | `security-hash-history-missing-target` | 存在しない hash、heading 以外の id、external URL、footnote backlink、empty hash を no-op にし、runtime 例外、build 時公開出力破壊、search index 混入を発生させない。 |

[§28.16〜§28.20](builder.md#sec-28-group-16-20) の browser runtime と parser precedence に関わる fixture では、`expected/site/assets/app.js` に handler 登録順、fallback 分岐、保護対象 token、`expected/effects.json.external_calls` に external call 0 件を固定する。security fixture では external script、CDN、runtime network fetch、raw HTML、event handler、credential、secret が HTML、CSS、JS、search index、stdout、stderr、REPORT、manifest に残らないことを `expected/security.json` に固定する。

<a id="sec-28-f-15"></a>
**[fixture 証跡責務 §28-F builder 詳細本文責務 §28.21〜§28.25 fixture 固定契約](fixture.md#sec-28-f-15)：**

[`docs/details/builder.md` 詳細本文責務 §28.21](builder.md#sec-28-21)〜[`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) の fixture は、accessibility、lightbox、print QR、definition list、task list の HTML / CSS / JS / search index、stdout、stderr、REPORT、副作用を固定する。manifest と共通 effects key は [§28.1〜§28.5](builder.md#sec-28-group-1-5) の共通契約に従う。

| feature slug | fixture | 固定する内容 |
|--------------|---------|--------------|
| `a11y` | `success-a11y-landmarks-labels` | `.skip-link`、`#main-content`、landmark role、TOC / search / icon button の `aria-label`、`:focus-visible`、REPORT `a11y_*` を固定する。 |
| `a11y` | `success-a11y-skip-link-tab-order` | skip link が最初の focus target になり、main content へ移動し、既存 keyboard shortcut と衝突しない focus 順を `expected/site/assets/app.js` と HTML で固定する。 |
| `a11y` | `failure-a11y-duplicate-id-strict` | 重複 id、空 label、focus 不能 skip target、keyboard trap を `BUILDER28_OUTPUT_VALIDATION_FAILED`、終了コード `1`、stdout 空、stderr 固定 error、公開出力維持にする。 |
| `a11y` | `security-a11y-no-keyboard-trap` | section collapse、TOC active、hash target、lightbox、skip link を併用しても Tab / Shift+Tab が閉じ込められず、focus outline が text を隠さない。 |
| `image-lightbox` | `success-lightbox-open-close` | trigger 数、page 1 個の dialog、open / close button、`aria-modal`、`aria-hidden`、opener focus return、REPORT `lightbox_images` を固定する。 |
| `image-lightbox` | `success-lightbox-escape-backdrop` | Escape、backdrop click、close button、Enter / Space activation、dialog hidden state、body scroll への副作用なしを `expected/site/assets/app.js` で固定する。 |
| `image-lightbox` | `failure-lightbox-alt-missing-strict` | strict で alt なし / 空 alt image を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持にする。 |
| `image-lightbox` | `security-lightbox-focus-trap` | Tab / Shift+Tab focus trap、external image no-fetch、escaped `data-lightbox-src`、external script / asset 不在を `expected/security.json` で固定する。 |
| `print-qr` | `success-print-qr-url` | `http` / `https` URL から `.print-qr`、`.print-qr-svg`、viewBox、rect order、print CSS、REPORT `print_qr=true` を固定する。 |
| `print-qr` | `noop-print-qr-empty-url` | URL 空値で QR SVG、print QR CSS、REPORT URL、search index text を出力せず、`print_qr=false`、`print_qr_url=""` にする。 |
| `print-qr` | `failure-print-qr-url-too-long` | 512 byte 超過、scheme 不正、credential 付き URL、制御文字入り URL を `BUILDER28_INVALID_OPTION`、終了コード `2`、stdout 空、公開出力維持にする。 |
| `print-qr` | `security-print-qr-svg-escape` | SVG 内 `script`、event handler、external href、foreignObject がなく、credential / secret 風 query が stdout、stderr、REPORT、manifest、search index に平文で残らない。 |
| `definition-lists` | `success-definition-list-single` | 単一 term / definition を `dl.definition-list`、`dt`、`dd` へ変換し、inline escape、REPORT `definition_lists=1` / `definition_terms=1` を固定する。 |
| `definition-lists` | `success-definition-list-multiple` | 複数 term、複数 definition、空行境界、paragraph 復帰、search index text order を固定する。 |
| `definition-lists` | `noop-definition-list-empty-term` | 空 term、空 definition、blockquote 内、list item 内、disabled option では通常 paragraph / list として扱い、warning と REPORT count を増やさない。 |
| `definition-lists` | `security-definition-list-inline-escape` | term / definition 内の raw HTML、quote、event handler、`javascript:` が escape され、link / badge / footnote / math inline 併用順が parser precedence と一致する。 |
| `task-lists` | `success-task-list-unchecked` | `[ ]` marker を disabled unchecked checkbox、`.task-list-item`、`.task-list-checkbox`、aria label、REPORT item count へ変換する。 |
| `task-lists` | `success-task-list-checked-nested` | `[x]` / `[X]` checked、nested list 階層維持、checked count、通常 list との混在を固定する。 |
| `task-lists` | `noop-task-list-non-target` | `[-]`、`[o]`、`[]`、`[xx]`、文中 marker、disabled option では通常 list text として扱い、checkbox を出力しない。 |
| `task-lists` | `security-task-list-disabled-aria` | checkbox が常に disabled、click で状態変更不可、aria label 非空、search index から checkbox label / marker text を除外する。 |

[§28.21〜§28.25](builder.md#sec-28-group-21-25) の accessibility と lightbox fixture は browser runtime fixture と同じ focus / keyboard / no-break 条件を再確認する。print QR、definition list、task list の security fixture は external call 0 件、外部 library 不使用、raw HTML 不在、credential / secret 非表示、search index 除外対象を `expected/security.json` に固定する。

<a id="sec-28-f-16"></a>
**[fixture 証跡責務 §28-F expected 比較方式固定契約](fixture.md#sec-28-f-16)：**

expected 比較は、実装環境差分で揺れないように以下の正規化だけを許可する。[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の固定表にない正規化、部分一致、snapshot 差し替え、目視承認は合格条件にしてはならない。

| ファイル | 比較方式 | 許可する正規化 | 禁止 |
|----------|----------|----------------|------|
| `expected/site/*.html` | DOM 構造、tag、属性順、text node の完全一致。 | 改行コードを LF に統一。末尾改行 1 個を許可。 | 属性順の無視、class subset 比較、画像 snapshot だけの比較。 |
| `expected/site/assets/style.css` | selector、property、値、media query の完全一致。 | 空行の連続を 1 行へ正規化可。 | selector の部分一致、未使用 selector の黙認。 |
| `expected/site/assets/app.js` | 対象 handler、storage key、guard、fallback 分岐を含む文字列完全一致。 | 改行コードを LF に統一。 | minify 差分の黙認、外部 script 参照の黙認。 |
| `expected/site/assets/search-index.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | HTML tag 混入、line number 混入、未定義 key の黙認。 |
| `expected/stdout.txt` | 行完全一致。`[REPORT]` は 1 行 key=value 形式。[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 追加 key は既存 [`docs/details/builder.md` 詳細本文責務 §8](builder.md#8-実行方法) key の後ろに ASCII 昇順で並べる。 | 末尾改行 1 個を許可。 | `[REPORT]` key 省略、型違い、件数差分、追加 key 順序違い、compact JSON 内空白の黙認。 |
| `expected/stderr.txt` | 行完全一致。fatal failure 以外は空 file。 | 末尾改行 1 個を許可。 | error code 差分、line 差分、message 差分、stdout warning 混入の黙認。 |
| `expected/effects.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | 外部 call、削除、既存出力破壊の黙認。 |
| `expected/builder-output.json` | 4 transaction flag の key、boolean 型、値完全一致。 | object key 順だけ無視可。 | staging 残存、公開置換、manifest 公開、search index 再生成の黙認。 |
| `expected/security.json` | JSON parse 後の key、型、値完全一致。 | object key 順だけ無視可。array 順は固定。 | CDN、credential、raw HTML、secret 残存の黙認。 |

<a id="sec-28-f-17"></a>
**[fixture 証跡責務 §28-F stdout / stderr / REPORT 固定契約](fixture.md#sec-28-f-17)：**

stdout、stderr、`[REPORT]` は、同じ入力から常に同じ順序で出力する。順序は、入力 file path 昇順、line 昇順、section 昇順、code 昇順とする。fatal failure の場合、`[REPORT]` は出力せず、stdout は空 file とする。

| 対象 | 固定内容 |
|------|----------|
| stdout warning | `[WARN] CODE file:line section message` の形式で完全一致。warning は stdout だけに出す。 |
| stderr error | `[ERROR] CODE file:line section message` の形式で完全一致。error は stderr だけに出す。 |
| message | 句点ありの日本語または ASCII 英文に統一し、secret、credential、raw HTML を含めない。 |
| REPORT boolean | `true` / `false` 小文字。 |
| REPORT integer | 0 以上の 10 進数。 |
| REPORT string | double quote 付き JSON string。[`docs/details/builder.md` 詳細本文責務 §8](builder.md#8-実行方法) 既存 key は既存形式を維持する。 |
| REPORT array | compact JSON array。要素順は発生順ではなく sorted string 昇順。ただし page order を意味する配列は input path 昇順。 |

<a id="sec-28-f-18"></a>
**[fixture 証跡責務 §28-F 基準出力安定性固定契約](fixture.md#sec-28-f-18)：**

[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の実装変更では、対象機能の fixture だけでなく、基準出力安定性 fixture を 1 件以上含める。既定有効機能は、有効化前後ではなく「機能対象入力なし」の基準 fixture を含める。

| 機能種別 | 基準出力安定性 fixture |
|----------|--------------|
| 明示有効化機能 | option 未指定時に既存 HTML / CSS / JS / search index / REPORT が変わらない fixture。 |
| 既定有効機能 | 対象 Markdown 記法や対象 DOM が存在しない入力で既存出力が変わらない fixture。 |
| reserved feature | 予約値指定時に出力が作られず、既存出力も破壊しない fixture。 |
| security failure | strict / non-strict の差分と、拒否対象が出力に残らない fixture。 |

<a id="sec-28-f-19"></a>
**[fixture 証跡責務 §28-F 機能別最低確認項目固定契約](fixture.md#sec-28-f-19)：**

各 [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の fixture は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 個別固定補足契約の validation、HTML / asset 固定、warning / error、REPORT count を最低 1 件以上の expected で確認する。[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の固定表の項目を fixture から省略してはならない。

| 節 | 最低確認項目 |
|----|--------------|
| [`docs/details/builder.md` 詳細本文責務 §28.1](builder.md#sec-28-1) | changed / reused page count、未変更 HTML byte 維持、manifest path validation、search index 全体再生成。 |
| [`docs/details/builder.md` 詳細本文責務 §28.2](builder.md#sec-28-2) | `html` 成功、`pdf` / `epub` 拒否、unknown format 拒否、出力破壊なし。 |
| [`docs/details/builder.md` 詳細本文責務 §28.3](builder.md#sec-28-3) | admonition type 正規化、badge color validation、disabled 時の基準出力安定性、escape。 |
| [`docs/details/builder.md` 詳細本文責務 §28.4](builder.md#sec-28-4) | line number node、copy 対象除外、空 code、line count。 |
| [`docs/details/builder.md` 詳細本文責務 §28.5](builder.md#sec-28-5) | slug 不変、表示番号、TOC / search index 番号、unknown mode 拒否。 |
| [`docs/details/builder.md` 詳細本文責務 §28.6](builder.md#sec-28-6) | h2 / h3 section 範囲、`aria-controls`、`data-section-id`、`adlaire:section-state` payload、print / search / hash 一時展開、重複 target 検出。 |
| [`docs/details/builder.md` 詳細本文責務 §28.7](builder.md#sec-28-7) | min/max validation、invalid format、TOC filter、heading numbering 併用、active tracking 対象一致。 |
| [`docs/details/builder.md` 詳細本文責務 §28.8](builder.md#sec-28-8) | fake git、fake file mtime、fallback、none、取得不能 failure、RFC3339 UTC 秒精度、search index 除外。 |
| [`docs/details/builder.md` 詳細本文責務 §28.9](builder.md#sec-28-9) | inserted / deleted / context / header class、line number 併用、escape、copy text / search index 清浄性。 |
| [`docs/details/builder.md` 詳細本文責務 §28.10](builder.md#sec-28-10) | lazy 属性、外部 URL no-fetch、data URI no-fetch、base 外 path strict、disabled no-op、alt escape、invalid scheme strict。 |
| [`docs/details/builder.md` 詳細本文責務 §28.11](builder.md#sec-28-11) | name / property key 正規化、head 内順序、重複 last wins、禁止 key、型 validation、attribute escape、secret 非表示。 |
| [`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) | light 固定、[`docs/details/builder.md` 詳細本文責務 §28.12](builder.md#sec-28-12) の禁止識別子不在、print light、禁止出力拒否。 |
| [`docs/details/builder.md` 詳細本文責務 §28.13](builder.md#sec-28-13) | colon / key-value / title-only title、copy / search 除外、empty title no-op、escape。 |
| [`docs/details/builder.md` 詳細本文責務 §28.14](builder.md#sec-28-14) | key validation、source 優先順位、code fence / span 非置換、invalid syntax no-op、missing var strict、secret 非表示、replacement count。 |
| [`docs/details/builder.md` 詳細本文責務 §28.15](builder.md#sec-28-15) | byte count、pre/code 保持、attribute order 保持、structure validation、marker validation、disabled 時の基準出力安定性、inline script / style 非追加。 |
| [`docs/details/builder.md` 詳細本文責務 §28.16](builder.md#sec-28-16) | active link 1 件化、aria-current、fallback scroll、depth sync。 |
| [`docs/details/builder.md` 詳細本文責務 §28.17](builder.md#sec-28-17) | graph TD SVG、unsupported source fallback、external script 不在。 |
| [`docs/details/builder.md` 詳細本文責務 §28.18](builder.md#sec-28-18) | reference order、backlink、duplicate definition warning、undefined strict。 |
| [`docs/details/builder.md` 詳細本文責務 §28.19](builder.md#sec-28-19) | inline / block math、code 内非変換、unclosed delimiter、escape。 |
| [`docs/details/builder.md` 詳細本文責務 §28.20](builder.md#sec-28-20) | pushState、focus、back / forward、missing target no-op。 |
| [`docs/details/builder.md` 詳細本文責務 §28.21](builder.md#sec-28-21) | skip link、landmark、button label、duplicate id strict、keyboard trap 不在。 |
| [`docs/details/builder.md` 詳細本文責務 §28.22](builder.md#sec-28-22) | trigger count、dialog 1 個、Escape / backdrop close、focus trap、alt warning。 |
| [`docs/details/builder.md` 詳細本文責務 §28.23](builder.md#sec-28-23) | URL validation、512 byte 制限、print-only SVG、通常表示非表示。 |
| [`docs/details/builder.md` 詳細本文責務 §28.24](builder.md#sec-28-24) | dl / dt / dd 構造、paragraph 境界、empty term no-op、inline escape。 |
| [`docs/details/builder.md` 詳細本文責務 §28.25](builder.md#sec-28-25) | checked / unchecked、disabled checkbox、nested list、aria、通常 list 非変換。 |

<a id="sec-28-f-20"></a>
**[fixture 証跡責務 §28-F strict / non-strict 固定契約](fixture.md#sec-28-f-20)：**

| ケース | non-strict fixture | strict fixture | 固定する差分 |
|--------|--------------------|----------------|--------------|
| warning で継続できる構文不正 | 終了コード `0`、warning count 増加、fallback 出力あり。 | 終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、`Done` 行なし、公開出力維持。 | stdout / stderr / effects。 |
| base 外 path | 対象参照を無効化し warning。 | 終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持。 | 参照先 file が作成されず、公開出力が置換されないこと。 |
| 未定義参照 | 通常 text または非表示 fallback。 | 終了コード `2`、stdout `[WARN]` と `[REPORT]`、stderr 空、公開出力維持。 | HTML fallback と strict 停止、`[REPORT]` の warning count。 |
| reserved feature | 終了コード `2`。 | 終了コード `2`。 | strict 差分なし。 |
| 内部エラー fixture | 終了コード `1`。 | 終了コード `1`。 | 既存出力維持。 |

<a id="sec-28-f-21"></a>
**[`docs/details/fixture.md` fixture 証跡責務 §28-F manifest 固定 schema](fixture.md#sec-28-f-21)：**

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `manifest.json` は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) をそのまま使用し、builder 専用の別 schema を定義しない。`name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) と [§28-F カタログ固定契約](#sec-28-f-2) の完全一致とする。`components` は `builder` と `collaborator_components` だけを ASCII 昇順、重複なしで持つ。`references` は対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) と [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) を含める。

`strict` と期待終了コードは `input/options.json` の正本値とし、`manifest.json` に `strict`、`expected_exit_code`、`feature_slug`、`owner`、`collaborators`、`spec_refs` を追加してはならない。

<a id="sec-28-f-22"></a>
**[fixture 証跡責務 §28-F builder 拡張実装検証証跡固定契約](fixture.md#sec-28-f-22)：**

実装検証証跡には、対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様)、追加 fixture 名、変更した HTML / CSS / JS / REPORT key、strict / non-strict 結果、外部依存なし確認、基準出力安定性確認、未実装の [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 機能を列挙する。対象外の [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 機能を先取り実装した場合、または [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) に存在しない Markdown 記法、CLI option、CSS class、JS 挙動を追加した場合は未完了として扱う。

<a id="sec-28-f-23"></a>
**[fixture 証跡責務 §28-F builder 拡張実装受け入れゲート固定契約](fixture.md#sec-28-f-23)：**

[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の実装変更は、[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の固定表を実装検証証跡で確認できる場合だけ受け入れ可能とする。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

| ゲート | 実装検証証跡 | 不足時の扱い |
|--------|---------|--------------|
| 対象範囲 | 対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) と対象外 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の一覧。 | 先取り実装または範囲不明として未完了。 |
| fixture catalog | 追加 / 更新した fixture 名の一覧と [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) catalog との対応。 | fixture 不足として未完了。 |
| manifest schema | 各 fixture の `manifest.json` が必須 key を満たす確認。 | fixture schema 不足として未完了。 |
| expected files | HTML、CSS、JS、search index、stdout、stderr、effects、builder output、security の該当 expected 更新。 | expected 不足として未完了。 |
| strict / non-strict | strict と non-strict の終了コード、stdout、stderr、effects の差分。 | 異常系未固定として未完了。 |
| REPORT | 追加 / 更新した REPORT key、型、count 単位、既定値。 | runner / API 連携不能として未完了。 |
| baseline stability | 対象機能無効時または対象入力なし時の基準出力安定性確認。 | 基準出力破壊リスクとして未完了。 |
| security | external call 0 件、CDN / external library 不使用、secret / credential / raw HTML 不在。 | security 不足として未完了。 |
| atomicity | failure fixture で既存出力、manifest、search index が維持される確認。 | 部分更新リスクとして未完了。 |
| not run | 未実施確認がある場合の理由と影響範囲。 | 検証不足として未完了。 |

<a id="sec-28-f-24"></a>
**[`docs/details/fixture.md` fixture 証跡責務 §28-F 不足時の固定扱い](fixture.md#sec-28-f-24)：**

[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の実装中に fixture 不足を発見した場合、実装判断で対象 fixture を省略してはならない。不足時は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。各機能の現在状態は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照する。

| 不足 | 扱い |
|------|------|
| fixture catalog の必須 fixture がない。 | 実装未完了。 |
| strict / non-strict の片方がない。 | 異常系未完了。 |
| `expected/builder-output.json` がない、または 4 件の transaction flag が不足する。 | builder output transaction 検証未完了。 |
| expected/security.json が必要なのにない。 | security 検証未完了。 |
| expected/effects.json に既存出力維持がない。 | atomicity 検証未完了。 |
| REPORT key の型または件数が expected にない。 | REPORT 検証未完了。 |
| 比較除外が `not_applicable` に理由付きで記載されていない。 | fixture schema 不備。 |

<a id="additional-management-api-fixture-contract"></a>
**[fixture 証跡責務 §29-F 追加管理 API fixture 固定契約](fixture.md#additional-management-api-fixture-contract)：**

追加管理 API fixture は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70)、[`docs/details/sdk.md` 詳細本文責務 §23.8](sdk.md#sec-23-8)、[`docs/details/ui.md` 詳細本文責務 §24.8](ui.md#sec-24-8)、[`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) を確認する。

| fixture 名 | fixture group root | 必須確認 |
|------------|--------------------|----------|
| `additional-management-users-roles-success` | `testdata/api/additional-management/` | users / roles の create、update、disable、delete、最後の admin 保護、audit、admin event。 |
| `additional-management-auth-permission-denied` | `testdata/api/additional-management/` | [`docs/details/security.md` 詳細本文責務 §27.58](security.md#sec-27-58) の permission 不足で `403`、対象 state 差分なし。 |
| `additional-management-request-validation` | `testdata/api/additional-management/` | body 禁止、未知 key、型不一致、path parameter 不正、query 不正が `400` または `422` で no-write。 |
| `additional-management-secret-redaction` | `testdata/api/additional-management/` | external auth secret、share token、Authorization header、Webhook secret が response、log、expected に出ない。share token は作成 response だけ許可。 |
| `additional-management-config-state-order` | `testdata/statefile/additional-management/` | config snapshot、template、restore、apply の write order、partial failure、secret mask。 |
| `additional-management-queue-retention-webhook` | `testdata/api/additional-management/` | queue cancel / reorder、retention run、webhook resend の conflict、confirmation、runner owner 境界。 |
| `additional-management-sdk-transport` | `testdata/sdk/additional-management/` | [`docs/details/sdk.md` 詳細本文責務 §23.8](sdk.md#sec-23-8) の全 method、path encode、query 省略、body 禁止、Blob / text / StreamHandle、admin event stream callback。 |
| `additional-management-ui-flow` | `testdata/ui/additional-management/` | [`docs/details/ui.md` 詳細本文責務 §24.8](ui.md#sec-24-8) の panel id、SDK only、成功後再取得順、confirmation、one-time secret 消去、SSE close、event feed callback 表示。 |
| `additional-management-cache-share-diff` | `testdata/api/additional-management/` | response cache の hit / miss / policy / purge、share link、snapshot diff、badge、metrics、OpenAPI、version、admin event stream の response、状態差分、副作用境界。 |

追加管理 API 実装検証証跡には、対象機能名、owner component、API endpoint、必要 permission、SDK method、UI panel、state path、fixture 名、実行結果、未実施項目を列挙する。未実施項目が 1 件でもある場合、対象機能を実装済みへ遷移してはならない。

追加管理 API fixture の必須 expected file は fixture 名ごとに以下へ固定する。対象外の expected file は `manifest.json.not_applicable` に理由を記録する。

| fixture 名 | 必須 expected |
|------------|---------------|
| `additional-management-users-roles-success` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-auth-permission-denied` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-request-validation` | `expected/response.json`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-secret-redaction` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-config-state-order` | `expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-queue-retention-webhook` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `additional-management-sdk-transport` | `expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/sdk_error.json`、`expected/security.json`。 |
| `additional-management-ui-flow` | `expected/sdk_trace.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/security.json`。 |
| `additional-management-cache-share-diff` | `expected/response.json`、`expected/events.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |

追加管理 API fixture の `manifest.json.category` は fixture 名ごとに以下へ固定する。fixture 名から category を推測してはならない。

| fixture 名 | `category` |
|------------|------------|
| `additional-management-users-roles-success` | `success` |
| `additional-management-auth-permission-denied` | `failure` |
| `additional-management-request-validation` | `failure` |
| `additional-management-secret-redaction` | `security` |
| `additional-management-config-state-order` | `partial` |
| `additional-management-queue-retention-webhook` | `partial` |
| `additional-management-sdk-transport` | `partial` |
| `additional-management-ui-flow` | `success` |
| `additional-management-cache-share-diff` | `partial` |

追加管理 API fixture は [`docs/details/api.md` 詳細本文責務 追加管理 API validation / error 優先順位固定契約](api.md#additional-management-validation-order) の順序を検証する。validation failure fixture は、後続の状態 read/write、外部通信、audit、admin event が発生しないことを `expected/effects.json.forbidden_writes` と `expected/effects.json.external_calls=[]` で示す。

追加管理 API validation failure fixture は、次の確認をすべて持つ。HTTP status と error body は [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70) の対象 endpoint 契約と完全一致させる。`expected/state/state-diff.json` は対象状態ファイルが byte 不変であることを示し、`expected/effects.json` は `created_paths`、`updated_paths`、`deleted_paths`、`external_calls`、`commands`、`notifications`、`streams`、`read_api_calls` を空配列にする。認証または認可まで到達する failure だけは、[`docs/details/security.md` 詳細本文責務](security.md) が許可する共通 audit / access log 副作用を `expected/effects.json` に明示する。`expected/security.json` は request body、Authorization header、share token、external auth secret、Webhook secret、SMTP password、raw restore template secret の平文が response、stdout、stderr、log、expected file に出現しないことを列挙する。

`additional-management-cache-share-diff` は `GET /api/version`、`GET /api/openapi.json`、`GET /api/events`、`GET /api/events/stream`、`GET /api/share/{token}/status`、`GET /api/cache-policy`、`POST /api/cache-policy`、`DELETE /api/response-cache` を同一 `additional-management-cache-share-diff` fixture 内の別 case として持つ。`GET /api/version` は `api_version="1"`、`compatible_versions=["1"]`、`deprecated_versions=[]`、`spec_version="V.N"`、`binary_version` が注入値または `V.0.0-dev` であること、`Cache-Control: no-store`、`.response_cache` no-read/no-write を `expected/response.json` と `expected/effects.json` で固定する。

`additional-management-cache-share-diff` の OpenAPI case は `openapi="3.1.0"`、`info.title="Adlaire CI API"`、`servers=[{"url":"/"}]`、path ASCII 昇順、method 固定順、`operationId` の重複なし、未実装 endpoint 不在、secret / token example 不在、`Cache-Control: no-store`、`.response_cache` no-read/no-write を固定する。OpenAPI の JSON object key 順を比較対象にする場合は fixture runner が canonical JSON へ正規化してから比較し、array 順は byte 単位で固定する。

`additional-management-cache-share-diff` の admin event stream case は、`input/events.jsonl` に schema-valid event、壊れた行、別 `type` event を含める。`expected/events.json` は `id:`、`event: admin-event`、`data:`、`: keepalive`、`event: error` の frame byte 列を LF 区切りで固定し、CRLF、`retry:`、複数 `data:` 行、secret 値が存在しないことを `expected/security.json` で確認する。client disconnect case は `expected/effects.json` の `streams` に close reason を記録し、状態、audit、admin event、cache が不変であることを固定する。

`additional-management-cache-share-diff` の share link case は token 作成 response にだけ token 本体が存在すること、`.share_links` に token hash だけが保存されること、list response に token / token_hash が出ないこと、作成と失効の `.share_links` → `.audit_log` → `.admin_events` 順、不正 token `401`、不在 hash / revoke 済み `404`、期限切れ `410`、`GET /api/share/{token}/status` の `.audit_log` / `.admin_events` no-write、scope 別 `SharedStatusResponse` key 1 件だけ、snapshot 2 件未満の `404` を固定する。`expected/security.json` は token 本体、token hash、token prefix、Authorization header、Cookie が response 許可箇所以外の expected file、log、effects、UI trace に存在しないことを列挙する。

`additional-management-cache-share-diff` の response cache case は許可 endpoint、禁止 endpoint、cache key canonical JSON、hit、miss、expired、corrupt、body hash mismatch、forbidden header、policy disabled、purge success、purge failure、業務状態 write 成功後の全 entry 削除を個別 case として固定する。hit case は endpoint 固有 state read が発生しないこと、miss case は handler 実行後に `200` JSON object だけを保存すること、write failure case は original response を維持し cache 保存だけを諦めることを `expected/effects.json.write_order` と `forbidden_writes` で示す。

`additional-management-config-state-order` は `config-snapshot-create-list-get`、`config-snapshot-corrupt-list-restore-denied`、`config-snapshot-restore-partial-failure`、`config-template-create-apply-delete`、`config-template-apply-partial-failure` を必須 case とする。snapshot case は [`docs/details/api.md` 詳細本文責務 §27.55](api.md#sec-27-55) の backup 対象、secret marker、id 採番、restore 書込順、`restored_paths`、corrupted list entry、delete 後 audit / admin event を固定する。template case は [`docs/details/api.md` 詳細本文責務 §27.65](api.md#sec-27-65) の name 重複、secret 平文拒否、subset restore schema、apply 書込順、`updated_paths`、partial failure no rollback を固定する。

`additional-management-queue-retention-webhook` は `queue-cancel-success-running-conflict`、`queue-reorder-success-invalid-set`、`history-retention-policy-run-success-disabled-partial`、`webhook-resend-accepted-not-resendable-conflict` を必須 case とする。queue case は `.build_state.queued` の exact 集合、active entry 不変、entry payload / priority / created_seq 不変、audit / admin event 順を固定する。retention case は `.server_config.history_retention` の既定値、policy no-op、run disabled no-write、削除順、log 削除失敗時の partial state、成功時の `.build_history` 差分を固定する。webhook resend case は `.notify_log` 対象条件、`.notify_config` 再取得、secret 非露出、accepted response、同一 delivery guard、runner notification 契約への引渡しを固定する。

`additional-management-cache-share-diff` は `status-badge-svg-contract` と `snapshot-site-diff-contract` を必須 case として持つ。status badge case は SVG root、固定 text、固定色、禁止要素、`Cache-Control: no-store`、`.response_cache` no-read/no-write、target / branch 不在時 `404` を `expected/response.json` と `expected/security.json` で固定する。snapshot site diff case は archive 検証済み 2 snapshot、added / removed / modified / unchanged_count、path sort、hash だけの response、同一 id `422`、破損 snapshot `500`、状態差分なしを固定する。

`additional-management-sdk-transport` は `streamAdminEvents(query,onEvent)` の `type` query だけ送信、`limit` / `offset` / `after` 不送信、`onEvent` TypeError、admin event frame callback、keepalive 破棄、error frame reject、user close `done=null`、server EOF reject、invalid frame reject を `expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/sdk_error.json` で固定する。

`additional-management-ui-flow` は event feed panel で `streamAdminEvents(query,onEvent)` だけを使用し、UI が `EventSource`、`fetch`、`ReadableStream` reader を直接生成しないことを `expected/ui_trace.json` で固定する。`onEvent` 受信時は API record の値だけを表示へ挿入し、stream error では error 表示 1 回と `getAdminEvents` 1 回、user stop では error 表示 0 回と `getAdminEvents` 1 回を `expected/ui_dom.json` と `expected/sdk_trace.json` で固定する。

追加管理 API fixture の `manifest.json` は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) に従う。`section` は対象 [`docs/details/api.md` 詳細本文責務 §27.48](api.md#sec-27-48)〜[§27.70](api.md#sec-27-70)、[`docs/details/sdk.md` 詳細本文責務 §23.8](sdk.md#sec-23-8)、[`docs/details/ui.md` 詳細本文責務 §24.8](ui.md#sec-24-8)、または [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) のいずれかを固定文字列で持つ。複数責務を横断する fixture は主 owner を 1 件だけ `owner_component` に置き、残りを `collaborator_components`、`components`、`references` に記録する。`assertions` は実在する expected file と 1 対 1 で対応させ、fixture 実行時に使わない expected file は作成せず、`not_applicable` に理由を置く。ただし `expected/events.json` は protocol event evidence として fixture 名別 expected 固定表で要求し、`manifest.json.assertions` の値としては追加しない。

追加管理 API fixture の `manifest.json.assertions` は次表に固定する。複数値は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) の列挙順で記録する。

| fixture 名 | 必須 assertions |
|------------|-----------------|
| `additional-management-users-roles-success` | `response`、`state`、`logs`、`effects`、`secret-mask`、`order` |
| `additional-management-auth-permission-denied` | `response`、`state`、`logs`、`effects`、`secret-mask`、`no-write` |
| `additional-management-request-validation` | `response`、`state`、`effects`、`secret-mask`、`no-write` |
| `additional-management-secret-redaction` | `response`、`state`、`logs`、`effects`、`secret-mask` |
| `additional-management-config-state-order` | `state`、`effects`、`secret-mask`、`order`、`idempotency` |
| `additional-management-queue-retention-webhook` | `response`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency` |
| `additional-management-sdk-transport` | `sdk-trace`、`sdk-return`、`sdk-error`、`secret-mask` |
| `additional-management-ui-flow` | `sdk-trace`、`ui-trace`、`ui-dom`、`secret-mask`、`order` |
| `additional-management-cache-share-diff` | `response`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency` |

<a id="mcp-fixture-contract"></a>
**[fixture 証跡責務 §30-F MCP fixture 固定契約](fixture.md#mcp-fixture-contract)：**

MCP fixture は [`docs/details/mcp.md`](mcp.md) 詳細本文責務を確認する。次表の `fixture group root` は分類用の親 directory であり、正式 fixture directory ではない。MCP の正式 fixture directory は `fixture group root` 直下の `fixture 名` と同名 directory とし、当該 directory の `manifest.json.name` と一致しなければならない。

| fixture 名 | fixture group root | 必須確認 |
|------------|--------------------|----------|
| `mcp-cli-lifecycle` | `testdata/mcp/cli/` | `--help`、`--version`、`--state-dir`、重複 option 拒否、loopback bind、non-loopback 拒否、禁止 host 拒否、read-only 起動。 |
| `mcp-jsonrpc-errors` | `testdata/mcp/jsonrpc/` | parse error、invalid request、batch 拒否、method not found、invalid params、unauthorized、forbidden、timeout。 |
| `mcp-initialize-client-log` | `testdata/mcp/jsonrpc/` | initialize、notifications/initialized、`.mcp_client_log` append、initialize 前 method 拒否。 |
| `mcp-tools-list-call` | `testdata/mcp/tools/` | `tools/list` descriptor、ToolDescriptor schema、read-only tool 除外、`tools/call` success/error、scope 不足。 |
| `mcp-tools-confirmation` | `testdata/mcp/tools/` | 副作用 tool の confirmation required、params hash mismatch、期限切れ、confirmation 付き成功。 |
| `mcp-resources-subscription` | `testdata/mcp/resources/` | `resources/list`、`resources/read`、subscribe、unsubscribe、resource-updated notification。 |
| `mcp-prompts` | `testdata/mcp/prompts/` | `prompts/list`、`prompts/get`、prompt params validation、prompt が副作用を実行しないこと。 |
| `mcp-sampling` | `testdata/mcp/sampling/` | client sampling request、timeout、外部 AI API direct call 不在、prompt hash、本文非保存。 |
| `mcp-sse` | `testdata/mcp/sse/` | keepalive、resource-updated、shutdown、client disconnect、未接続時 notification 破棄。 |
| `mcp-state-metrics-audit` | `testdata/mcp/state/` | `.mcp_config`、`.mcp_audit_log`、`.mcp_client_log`、`.mcp_metrics` の schema、append、破損時処理。 |

`mcp-state-metrics-audit` は [`docs/details/statefile.md` 詳細本文責務 §22.0d](statefile.md#sec-22-0d) の `McpConfig`、`McpScopeRecord`、`McpAuditRecord`、`McpClientRecord`、`McpMetrics` を byte 単位で検証する。`.mcp_config.scopes[].token_hash` は SHA-256 lowercase hex だけを許可し、token 本体、Authorization header 値、token の部分文字列を `expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json` に含めてはならない。

`mcp-cli-lifecycle` は [`docs/details/mcp.md` 詳細本文責務 §29.1](mcp.md#sec-29-1) の CLI 検証順を固定する。少なくとも `--state-dir` 重複、`--addr` 重複、`--client-token` 重複、`--read-only` 重複、`--allow-non-loopback` 重複、`--addr 0.0.0.0:8766`、`--addr 255.255.255.255:8766`、`--addr 224.0.0.1:8766`、`--addr [::1]:8766`、`--addr example.local:8766`、`--addr 192.168.1.10:8766` かつ `--allow-non-loopback` なし、`--addr 192.168.1.10:8766` かつ `--allow-non-loopback` ありを別 case とする。失敗 case は stdout 空、stderr 固定 1 行、終了 code `2`、listener 起動 0 件、state read/write 0 件、`.mcp_client_log` / `.mcp_audit_log` / `.mcp_metrics` 更新 0 件を `expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json` で固定する。成功 case は bind 対象、read-only flag、token 要否、状態副作用なしを `expected/effects.json` に固定する。

`mcp-tools-list-call` は [`docs/details/mcp.md` 詳細本文責務 §29.6](mcp.md#sec-29-6) の `ToolDescriptor` を tool 名ごとに byte 等価で検証する。`name`、`description`、`inputSchema`、`annotations` 以外の key、固定 description 不一致、`required` の順序不一致、未知 `properties`、`additionalProperties:false` 欠落、`readOnlyHint` 不一致、[`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) と異なる tool order は不合格とする。read-only 起動 fixture は副作用 `yes` の tool が `tools/list` から除外され、残る tool の相対順序が維持されることを `expected/response.json` で固定する。read-only 起動中に副作用 tool を直接 `tools/call` する case は JSON-RPC `-32002 Forbidden`、対象 owner state 変更 0 件、`.mcp_audit_log` 追記 0 件、`.mcp_metrics` 更新 0 件を `expected/response.json` と `expected/effects.json` で固定する。

`mcp-state-metrics-audit` は `.mcp_metrics.tools["<tool_name>"]` の `success_count`、`error_count`、`timeout_count`、`last_status`、`last_duration_ms`、`last_at` を検証する。`tool_name` は [`docs/details/mcp.md` 詳細本文責務 §29.5](mcp.md#sec-29-5) の tool name と完全一致させる。`tool_name`、`status`、`count` を sibling key として保存する旧形式、`forbidden_count`、未知 metrics key は不合格とする。

`mcp-state-metrics-audit` は `adlaire.getMetrics` の返却 snapshot が現在の `adlaire.getMetrics` 呼び出し分を含まず、response data 確定後に `.mcp_metrics.tools["adlaire.getMetrics"]` を更新する順序を `expected/effects.json.write_order` で固定する。

`mcp-state-metrics-audit` は副作用 tool の実行前 audit と sampling result audit を分けて検証する。副作用 tool の audit は対象 owner component 呼び出し前に `status:"accepted"`、`duration_ms:0`、`jsonrpc_error_code:null`、`build_id:null`、`prompt_hash:null` で `.mcp_audit_log` へ追記する。sampling result audit は `tool:"adlaire.analyzeBuildError"`、`build_id`、`prompt_hash` を持ち、sampling response 本文を保存しない。sampling result audit の追記失敗では JSON-RPC `-32603 Internal error` を返し、`.mcp_metrics` を更新しないことを固定する。

`mcp-state-metrics-audit` は metrics 更新失敗の partial case を持つ。対象 owner component または sampling request が完了した後に `.mcp_metrics` 更新が失敗した場合は、完了済み owner state を rollback せず、JSON-RPC `-32603 Internal error` を返し、失敗地点以降の追加 write を行わないことを `expected/effects.json` の `write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes` で固定する。

MCP 実装検証証跡には、tool 名、resource URI、prompt 名、scope、confirmation、timeout、audit、metrics、client log、SSE event、未実施項目を含める。副作用 tool を確認する場合、confirmation なし実行拒否と confirmation 付き実行成功の両方を必須とする。secret、token、Authorization header、raw params の secret 値が stdout、stderr、audit、client log、metrics、expected に出現した場合は不合格とする。

MCP fixture の expected file は fixture 名ごとに以下へ固定する。対象外の expected file は `manifest.json.not_applicable` に理由を記録する。

| fixture 名 | 必須 expected |
|------------|---------------|
| `mcp-cli-lifecycle` | `expected/stdout.txt`、`expected/stderr.txt`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-jsonrpc-errors` | `expected/response.json`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-initialize-client-log` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-tools-list-call` | `expected/response.json`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-tools-confirmation` | `expected/response.json`、`expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-resources-subscription` | `expected/response.json`、`expected/events.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-prompts` | `expected/response.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-sampling` | `expected/request.json`、`expected/response.json`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-sse` | `expected/events.json`、`expected/effects.json`、`expected/security.json`。 |
| `mcp-state-metrics-audit` | `expected/state/state-diff.json`、`expected/logs/`、`expected/effects.json`、`expected/security.json`。 |

MCP fixture の `manifest.json.category` は fixture 名ごとに以下へ固定する。fixture 名から category を推測してはならない。

| fixture 名 | `category` |
|------------|------------|
| `mcp-cli-lifecycle` | `partial` |
| `mcp-jsonrpc-errors` | `failure` |
| `mcp-initialize-client-log` | `success` |
| `mcp-tools-list-call` | `partial` |
| `mcp-tools-confirmation` | `partial` |
| `mcp-resources-subscription` | `success` |
| `mcp-prompts` | `noop` |
| `mcp-sampling` | `partial` |
| `mcp-sse` | `success` |
| `mcp-state-metrics-audit` | `partial` |

MCP fixture は [`docs/details/mcp.md` 詳細本文責務 §29.6 Tool schema](mcp.md#sec-29-6) の未知 key、必須 key 不足、型不一致、範囲外を個別に検証する。params validation failure では `.mcp_audit_log`、`.mcp_metrics`、対象 owner state を更新しないことを固定する。

MCP fixture の `manifest.json` は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) に従う。`owner_component` は `mcp` 固定、`collaborator_components` は fixture が呼び出す owner component だけを持ち、`components` は `mcp` と `collaborator_components` を ASCII 昇順で持つ。`section` は [`docs/details/mcp.md` 詳細本文責務 §29.0](mcp.md#sec-29-0)〜[§29.17](mcp.md#sec-29-17) の対象節を固定文字列で持つ。`references` は対象 MCP 節、呼び出す owner component 詳細本文、必要な [`docs/details/fixture.md` fixture 証跡責務 §30-F](fixture.md#mcp-fixture-contract) を含める。read-only 起動 fixture では副作用 tool を `tools/list` から除外する expected を必須とし、除外した tool 名を `not_applicable` ではなく `expected/response.json` に記録する。

MCP fixture の `manifest.json.assertions` は次表に固定する。複数値は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) の列挙順で記録する。`expected/events.json` と `expected/request.json` は MCP protocol evidence として fixture 名別 expected 固定表で要求し、`manifest.json.assertions` の値としては追加しない。

| fixture 名 | 必須 assertions |
|------------|-----------------|
| `mcp-cli-lifecycle` | `stdout`、`stderr`、`effects`、`secret-mask` |
| `mcp-jsonrpc-errors` | `response`、`state`、`effects`、`secret-mask`、`no-write` |
| `mcp-initialize-client-log` | `response`、`state`、`logs`、`effects`、`secret-mask`、`order` |
| `mcp-tools-list-call` | `response`、`state`、`effects`、`secret-mask` |
| `mcp-tools-confirmation` | `response`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency` |
| `mcp-resources-subscription` | `response`、`effects`、`secret-mask`、`order` |
| `mcp-prompts` | `response`、`effects`、`secret-mask`、`no-write` |
| `mcp-sampling` | `response`、`state`、`effects`、`secret-mask`、`order` |
| `mcp-sse` | `effects`、`secret-mask`、`order` |
| `mcp-state-metrics-audit` | `state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency` |

MCP JSON-RPC failure fixture は、`id` の保持、`jsonrpc:"2.0"`、error code、message、data の有無を `expected/response.json` に固定する。parse error、invalid request、batch request、method not found、invalid params、unauthorized、forbidden、timeout は別 fixture とし、1 fixture で複数の失敗分類を兼用してはならない。initialize 前に許可されない method を呼んだ場合、`.mcp_client_log`、`.mcp_audit_log`、`.mcp_metrics`、対象 owner state の各副作用有無を `expected/effects.json` に明示する。

MCP tool call fixture は、`tools/list` の descriptor と `tools/call` の request / response を別 assertion として記録する。副作用 `yes` の tool は、confirmation なし拒否、confirmation params hash 不一致、confirmation 期限切れ、confirmation 付き成功を別 fixture とする。confirmation なし拒否、不一致、期限切れでは対象 owner state、runner queue、config、archive、notification、`.mcp_audit_log`、`.mcp_metrics` を更新してはならない。confirmation 付き成功では、`status:"accepted"` の audit 追記、対象 owner 呼び出し、metrics 更新、JSON-RPC success response の順序を `expected/effects.json` に記録する。

MCP security fixture は、Authorization header、API token、session token、tool params 内 secret、prompt 本文に含まれる secret、sampling request 本文の secret が stdout、stderr、`.mcp_client_log`、`.mcp_audit_log`、`.mcp_metrics`、`expected/response.json`、`expected/events.json` に平文で出現しないことを `expected/security.json` に列挙する。secret を placeholder として表す場合は `${secret:<source_id>}` だけを使用し、実値、hash 入力、部分文字列、長さから復元可能な値を expected file に置いてはならない。
