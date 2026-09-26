# Adlaire CI — Fixture 詳細仕様

[`docs/details/fixture.md`](fixture.md) は fixture 証跡責務として、fixture、expected、fake、実装検証証跡、acceptance checklist、差し戻し条件だけを扱う。

component 境界管理の参照先は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0b.1](../DETAIL_INDEX.md#0b1-owner-component-別-owner-collaborator-境界管理) とする。[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は owner component 別詳細本文責務ではなく、fixture 入力、expected、effects、assertion、実装検証証跡、検証観点の証跡本文を持つ。

---

## 0. 責務境界

| 項目 | 内容 |
|------|------|
| 証跡責務 | `fixture` |
| 対象 component | `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`security`、`setup` |
| 持つ内容 | fixture 証跡責務が本文として定義する fixture manifest、assertion、fake、testdata、expected / effects、受け入れ fixture 共通契約、実装検証証跡テンプレート、acceptance checklist、差し戻し条件。 |
| 持たない内容 | 個別 component の通常処理本文、API endpoint 詳細、SDK method 実装、UI DOM 詳細、状態 schema、setup / update 実行手順、release 生成・公開手順。 |

---

<a id="対象範囲"></a>
**対象範囲：**

| 範囲 | 内容 |
|------|------|
| [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) | fixture / testdata 配置、fake 実装、実装検証証跡。 |
| [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) | API の必須検証、API fixture、API / SDK / UI / 状態ファイル cross fixture 固定。 |
| [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の fixture 配置、fixture カタログ、manifest、assertion、expected/effects、相互整合、component 別検証責務。 |
| [`docs/details/fixture.md` fixture 証跡責務 §27-F-EVIDENCE](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の実装検証証跡、受け入れゲート、差し戻し条件、部分失敗・再実行契約。 |

<a id="current-implementation-alignment-evidence"></a>
**現行実装整合証跡：**

以下は実在する実装 artifact と owner component 詳細本文を照合した未完了証跡である。現在状態の割当は [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務を参照する。該当する実装と必須証跡が契約に一致した場合は、同じ変更単位でこの証跡を更新または解消する。

| ID | 対象 | 現在確認できる実装証跡 | 受け入れに必要な証跡 |
|----|------|------------------------------|----------------------------|
| <a id="align-01"></a>`ALIGN-01` | 起動入口 | [`main.go`](../../main.go) は実行ファイル basename に `adlaire-ci-runner` を含む場合だけ runner を起動し、それ以外は `adlaire-ci-api` を含めて builder を起動する。標準 3 バイナリ名の完全一致判定、API 用 CLI の `--addr` / `--state-dir` / `--init-credentials` / `--help` / `--version`、API listener 起動経路が存在しない。[`components/api.go`](../../components/api.go) の `NewAPIServer`、`ListenAndServe`、`InitCredentials` は起動入口へ接続されていない。 | `adlaire-ci-build`、`adlaire-ci-runner`、`adlaire-ci-api` がそれぞれ対応 owner component を起動する [`main.go`](../../main.go) の完全一致分岐、API CLI parser、listener 実行、`--init-credentials` 実行経路、および 3 バイナリ別 CLI fixture。 |
| <a id="align-02"></a>`ALIGN-02` | API / admin | API は [`admin/index.html`](../../admin/index.html) と [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) を配信しない。 | [`docs/details/admin.md`](admin.md) の配信契約と API fixture。 |
| <a id="align-03"></a>`ALIGN-03` | API / SDK | [`components/api.go`](../../components/api.go) に `GET /api/approvals`、`POST /api/approvals/{id}/approve`、`POST /api/approvals/{id}/reject`、`GET /api/build-chain-config`、`POST /api/build-chain-config`、`GET /api/stats/build-trends` の 6 endpoint が存在しない。`GET /api/build/stream` は `log` frame の `at` と終端 `end` frame を送信しない。 | [`docs/details/api.md`](api.md) と [`docs/details/sdk.md`](sdk.md) の公開契約に対応する 6 endpoint、`log` / `end` の exact SSE frame、および API / SDK cross fixture。 |
| <a id="align-04"></a>`ALIGN-04` | API security | [`components/api.go`](../../components/api.go) は `.api_tokens` を root array として読み書きし、record に `name` を保存して `label`、`expires_at`、`last_used_at` を保存しない。作成 request の旧 `name` key と GET / POST response の追加 `name` key を許可し、scope 省略時に未定義の `write` を含む `read,write` を保存する。認証時は token の有効性だけを確認し、endpoint ごとの scope と `expires_at` を強制せず、認証成功時の `last_used_at` も更新しない。 | [`docs/details/security.md`](security.md) と [`docs/details/statefile.md`](statefile.md) に一致する `.api_tokens` object schema、`label` / `expires_at` / `last_used_at`、許可 scope、expiry、scope 認可、認証成功更新、および secret 非表示 fixture。 |
| <a id="align-05"></a>`ALIGN-05` | UI / SDK / design | [`admin/index.html`](../../admin/index.html) は SDK 113 public method のうち 31 method を UI 操作へ接続せず、承認一覧に approve / reject 操作がない。一時秘密情報は token と TOTP secret の 2 領域だけを一般 click で消去し、otpauth URI 専用領域、trusted event、copy 完了、generation 一致、全終了経路の消去を実装しない。build stream は `StreamHandle.done` を監視しない。inline CSS は [`docs/DESIGN.md` デザイン責務 標準管理 UI 視覚契約](../DESIGN.md#admin-ui-visual-contract) の `:focus-visible` indicator を実装していない。 | [ALIGN-05 UI 未接続 SDK 操作](#align-05-unconnected-sdk-operations)、[`docs/details/ui.md`](ui.md) の全必須操作、one-time secret、stream terminal、SDK 呼出しを照合する UI fixture、および [`docs/DESIGN.md`](../DESIGN.md) の selector / token / focus 固定値を照合する構造化 visual fixture。 |
| <a id="align-06"></a>`ALIGN-06` | builder / design | builder の生成テンプレートが [`docs/DESIGN.md`](../DESIGN.md) の token、寸法、sidebar、Markdown 画像 selector、トップへ戻る表示と一致しない。 | デザイン正本に一致する生成物と builder fixture。 |
| <a id="align-07"></a>`ALIGN-07` | fixture | builder の必須 fixture 一部と `testdata/runner/`、`testdata/api/`、`testdata/admin/`、`testdata/sdk/`、`testdata/ui/`、`testdata/statefile/`、`testdata/archive/`、`testdata/commitstatus/`、`testdata/security/`、`testdata/setup/` の fixture root が未作成である。 | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#sec-0g-8-f) の必須 path、expected、fake、実行証跡。 |
| <a id="align-08"></a>`ALIGN-08` | empty directory fixture | `testdata/builder/empty-dir/` は `.keep` を含む marker-based fixture であり、文字どおりの空 directory ではない。 | [`docs/details/fixture.md` fixture 証跡責務 §8a-F](fixture.md#8a-f-builder-初期受け入れ-fixture-契約) と検証処理の marker 除外条件が一致する証跡。 |
| <a id="align-09"></a>`ALIGN-09` | runner / state schema | [`components/runner.go`](../../components/runner.go) の build log / history は旧 field 構成であり、正規化 `status`、詳細 `target_status`、拡張 object、history の必須 key が [`docs/details/statefile.md`](statefile.md) の schema と一致しない。 | runner、API adapter、archive reader の同型 schema、未知 key、欠落 key、状態写像を検証する fixture。 |
| <a id="align-10"></a>`ALIGN-10` | API health / PAT | [`components/api.go`](../../components/api.go) の health response は field 名、正規化状態、`checks` が [`docs/details/api.md`](api.md) と一致しない。PAT status は期限を返さず、PAT verify は GitHub 外部検証と scope 取得を行わない。 | health、PAT status、PAT verify の request、response、外部呼出し、秘密情報非表示を照合する API / SDK / UI fixture。 |
| <a id="align-11"></a>`ALIGN-11` | API request / log schema | [`components/api.go`](../../components/api.go) は request ID を生成・伝播せず、`.api_access_log` と `.config_log` の record field、必須 log の失敗境界が [`docs/details/api.md`](api.md) と [`docs/details/statefile.md`](statefile.md) の契約に一致しない。 | request ID、response header、API access log、config diff log、partial failure、log 追記失敗境界を検証する fixture。 |
| <a id="align-12"></a>`ALIGN-12` | queue dispatch / active state / cancel | [`components/api.go`](../../components/api.go) は idle の manual request で build id と `running=true` を直接保存し、systemd への runner 起動要求を行わない。API と [`components/runner.go`](../../components/runner.go) の `.build_state` schema は `active_queue_entry` を持たず、runner は waiting-to-active 遷移と未確定 entry の再実行を行わない。cancel は `.build_state.running=false` と `current_build_id=null` を直接保存するだけで、実行 process への cancel request、対象 build log の `cancelled` 最終化、history / status、active entry を保持した finalizer を実行しない。 | [`docs/details/api.md`](api.md)、[`docs/details/runner.md`](runner.md)、[`docs/details/statefile.md`](statefile.md) の durable queue、非同期 runner 起動、atomic move、at-least-once retry、cancel request / process 停止 / `cancelled` 最終結果、response 契約を検証する queue / dispatch fixture。 |
| <a id="align-13"></a>`ALIGN-13` | builder Markdown / report | [`components/builder.go`](../../components/builder.go) は h5 / h6 と対応 CSS、table の固定 grammar、未閉鎖 fence の fallback warning、7 段以上 list の clamp warning、再帰 blockquote、definition list、`build_id` / `commit_sha` / `build_at` を含む REPORT 固定順の全部を実装していない。 | [`docs/details/builder.md` 詳細本文責務 §4.4〜§4.5](builder.md#sec-4-4)、[§6](builder.md#6-css-クラス一覧)、[§8](builder.md#8-実行方法) と一致する HTML / CSS / warning / REPORT、および各正常・fallback・strict fixture。 |
| <a id="align-14"></a>`ALIGN-14` | builder CLI / atomic publication | [`components/builder.go`](../../components/builder.go) の `BuildConfig`、CLI parser、help、生成 HTML は `build_id`、`commit_sha`、`build_at` を受け渡さない。strict warning は warning と REPORT を stdout へ出す前に stderr error で終了する。atomic writer は出力親 directory を暗黙作成し、既存 staging path を削除し、公開前検証、file sync、directory sync、parent sync を行わず、cleanup failure を報告しない。 | [`docs/details/builder.md`](builder.md) の CLI metadata、strict stdout / stderr / exit、更新前出力維持、staging 衝突、公開前検証、sync、rename、復元、cleanup 契約を正常系と failure injection fixture で確認した証跡。 |
| <a id="align-15"></a>`ALIGN-15` | SDK request / response / stream | [`admin/adlaire-ci-sdk.js`](../../admin/adlaire-ci-sdk.js) は 5 public method の引数契約と `createToken()` の body mapping が詳細仕様と一致せず、SDK 送信前の型・範囲・空配列検証が不足する。JSON media type を parse せず部分文字列と末尾文字列で判定するため、不正 media type を JSON として受理し、parameter 付き `application/*+json` を JSON として受理しない。`streamBuild()` は response media type、fatal UTF-8、SSE key / type / `at` / 最終 `end` / EOF を検証せず、`StreamHandle.done` と terminal state を持たず、接続後の frame / callback error を固定 error として通知しない。 | [ALIGN-15 SDK 契約差分](#align-15-sdk-contract-gaps)、[`docs/details/sdk.md`](sdk.md) の request shape、送信前 validation、JSON media type / response error、SSE parser、`StreamHandle` terminal transition、および error shape fixture。 |
| <a id="align-16"></a>`ALIGN-16` | runner config / execution | [`components/runner.go`](../../components/runner.go) は `.branch_config` の `branch_targets` だけを読み、`.server_config`、`.repo_config`、`.pipeline_config`、`.notify_config` の owner 契約を統合しない。pipeline は固定 `.ci/pipeline.sh` と固定 300 秒 timeout だけを実行し、pipeline YAML step、設定 timeout、hook、branch environment、remote build、commit status、dry-run、retry、multi-file / watch / tag / parallel trigger、approval、dependency chain、priority queue、failure category、trigger / build status、trend / anomaly / environment 証跡を実行しない。通知は旧 `webhooks` array の webhook 送信だけで、channel schema、署名、SMTP / command、`.notify_log`、週次 summary、規定の retry 分類を満たさず、送信 error を一律 pending とする。 | [`docs/details/runner.md`](runner.md)、[`docs/details/commitstatus.md`](commitstatus.md)、[`docs/details/statefile.md`](statefile.md) の owner / collaborator 契約に従う設定統合、pipeline 実行、通知、承認、連鎖、優先度、観測値の実装と、成功、部分失敗、再実行、dry-run、外部 fake を固定する runner fixture。 |
| <a id="align-17"></a>`ALIGN-17` | statefile persistence / schema | [`components/runner.go`](../../components/runner.go) と [`components/api.go`](../../components/api.go) の JSON writer は `{name}.lock`、10 秒待機、排他的 temp 作成、file sync、mode 検証、rename 後 parent directory sync、失敗時 cleanup を共通化せず、runtime directory を `0755` で作成する。JSON Lines は lock / sync なしで直接 append し、reader は破損行を証跡なしで無視する。一般 reader は未知 key と必須 key 欠落を拒否しない。`.hooks`、`.alert_rules`、`.tag_rules` は object wrapper 契約に対して root array を保存し、UTF-8 text 契約の `.notes` を JSON object として扱う。`.api_tokens`、build log / history、`.build_state` の個別差分は [`ALIGN-04`](#align-04)、[`ALIGN-09`](#align-09)、[`ALIGN-12`](#align-12) を参照する。 | [`docs/details/statefile.md`](statefile.md) の共通 atomic write / append、`0700` directory、strict schema、root schema、lock timeout、sync、cleanup 契約を単一 adapter で満たし、正常系、競合、未知 / 欠落 key、破損 JSON Lines、各 I/O failure injection を固定する statefile fixture。 |
| <a id="align-18"></a>`ALIGN-18` | API security lifecycle | [`components/api.go`](../../components/api.go) は route / method 判定より前に IP access 判定を行い、session middleware で API token を endpoint scope、expiry、`last_used_at` なしに受理する。API 全体と login の規定 rate limit / lock、必須 audit / access log の成功保証を実装しない。credential 初期化の salt は 16 bytes、password 比較は通常文字列比較で、初期化と変更の長さ、NUL、同一 password 条件が不一致であり、login 失敗回数を永続化する。TOTP は空白を除去して受理し、`last_accepted_step` の比較 / 更新による replay 防止を行わない。token schema 固有差分は [`ALIGN-04`](#align-04) を参照する。 | [`docs/details/security.md`](security.md)、[`docs/details/api.md`](api.md)、[`docs/details/statefile.md`](statefile.md) の順序付き security middleware、password / hash、memory-only login lock、TOTP replay 防止、token scope / expiry / 利用時更新、rate window、必須 audit / access log 契約を固定する security / API fixture。 |
| <a id="align-19"></a>`ALIGN-19` | log archive / snapshot | [`components/runner.go`](../../components/runner.go) の log archive は build log の `finished_at` ではなく file mtime を使用し、current build 除外、決定的順序、temp / rename / sync、gzip 検証後の source 削除、失敗時 cleanup を行わない。snapshot は出力 directory を `.snapshots/{id}/site` へ直接 copy し、`site.tar.gz`、`meta.json`、checksum、USTAR、atomic publish、公開前検証を実装しない。[`components/api.go`](../../components/api.go) は metadata 不足 snapshot を一覧へ含め、download artifact を検証せず、削除を runner guard / tombstone / sync なしの `RemoveAll` で行い、rollback は archive の検証、展開、deploy、coordinator を行わず成功 history を追記する。 | [`docs/details/archive.md`](archive.md)、[`docs/details/runner.md`](runner.md)、[`docs/details/api.md`](api.md)、[`docs/details/statefile.md`](statefile.md) の log archive、snapshot 作成 / 一覧 / download / 三状態削除 / rollback 契約を満たし、時刻境界、current build、決定的 archive、checksum、破損、競合、各 I/O failure injection を固定する archive cross fixture。 |
| <a id="align-20"></a>`ALIGN-20` | builder extension catalog | [`components/builder.go`](../../components/builder.go) の CLI、`BuildConfig`、生成 HTML / CSS / JavaScript は、[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の差分ビルド、出力形式、行番号、見出し番号、section 折りたたみ、TOC depth、更新日時、diff highlight、画像 lazy load、custom meta、出力名、template 変数、minify、scrollspy、Mermaid、数式、履歴、accessibility、lightbox、print QR の option、設定解決、生成物、runtime 挙動を実装しない。現在時刻を直接取得する箇所があり、§28 の決定性入力にも接続していない。 | [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) と [`docs/details/fixture.md` fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の全機能を、CLI / env / config precedence、決定性、atomic publication、browser runtime、visual / accessibility assertion まで確認する builder extension fixture。 |
| <a id="align-21"></a>`ALIGN-21` | API common HTTP validation | [`components/api.go`](../../components/api.go) の `decodeBody()` は 1 MiB 超過を `400 Invalid JSON` として扱い、最初の JSON value 後の trailing data を拒否しない。`rejectBody()` は size 上限と read error を扱わず、query parser は未知 key と同一 key の重複を拒否しない。method helper は endpoint ごとの許可 method 集合を統合せず、validation は最初の 1 件だけを返すため、共通処理順序、status、error detail 順序が [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) と一致しない。 | [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) の path / method / access / auth / body / query / path parameter 順序、`413`、単一 JSON document、未知・重複 query、複数 validation detail、状態不変を固定する API common lifecycle fixture。 |
| <a id="align-22"></a>`ALIGN-22` | API configuration / control | [`components/api.go`](../../components/api.go) の backup は必須設定の既定値を補わず `config` 配下へ格納し、restore は全必須 key と全 schema を先行検証せず、secret file の保持・削除、全体 no-op、固定書込順、audit を実装しない。汎用設定 handler は pipeline、dashboard、SMTP の strict schema と `.smtp_secret` 分離を満たさず、SMTP test は送信せず成功 log を作成する。schedule 更新は systemd timer を変更せず、hook / alert / tag rule は owner schema を検証しない。 | [`docs/details/api.md`](api.md)、[`docs/details/statefile.md`](statefile.md)、[`docs/details/security.md`](security.md) の backup / restore、設定 mutation、systemd、SMTP / notification、rule / hook schema、no-op、partial failure、audit の処理順を固定する API configuration cross fixture。 |
| <a id="align-23"></a>`ALIGN-23` | API operational lifecycle | [`components/api.go`](../../components/api.go) の diagnostics は GitHub token と出力の 2 item だけを返し、GitHub API、systemd、Webhook 設定を検査しない。GitHub rate-limit と PAT verify は外部確認を行わず `remote_checked:false` の固定値を返す。Webhook は branch / event filter、delivery 重複排除、active queue 照合、audit、runner 非同期起動要求を実装しない。maintenance 判定は破損時に disabled と扱い、circuit reset と maintenance mutation は規定 audit / config log 境界を満たさない。 | [`docs/details/api.md`](api.md)、[`docs/details/runner.md`](runner.md)、[`docs/details/security.md`](security.md)、[`docs/details/statefile.md`](statefile.md) の diagnostics、GitHub 外部確認、Webhook、maintenance fail-closed、circuit / operation mutation、queue dispatch と必須副作用を fake 外部依存で固定する API operations fixture。 |
| <a id="align-24"></a>`ALIGN-24` | executable verification | [`components/builder_test.go`](../../components/builder_test.go)、[`components/runner_test.go`](../../components/runner_test.go)、[`components/api_test.go`](../../components/api_test.go) は現行の部分実装と旧 schema を直接組み立てる test が中心で、[`docs/details/fixture.md`](fixture.md) の manifest、`input/`、`expected/`、failure injection を実行する共通 harness に接続していない。SDK / UI の executable test harness も存在しない。このため現行 test の成功だけでは責務正本への適合を証明できない。 | [`docs/details/fixture.md`](fixture.md) の必須 fixture root、manifest、expected / effects、fake、証跡を実行する Go / JavaScript 検証 harness、旧期待値の除去、各 `ALIGN-*` の正常・異常・partial / no-op / security assertion、および仕様非準拠動作を成功として受理しない regression test。 |
| <a id="align-25"></a>`ALIGN-25` | builder input / path safety | [`components/builder.go`](../../components/builder.go) は `--src` と `--out` の同一・相互包含を拒否せず、入力列挙で hidden Markdown file と symlink Markdown file を含め、大文字拡張子を lowercase 化して受理する。入力 file の 10 MiB 上限、UTF-8 BOM 除去、CRLF / CR の LF 正規化も行わない。 | [`docs/details/builder.md` 詳細本文責務 §2〜§2a](builder.md#2-ファイルパス設定)、[§8a](builder.md#8a-builder-受け入れ検証条件) と [`docs/details/fixture.md` fixture 証跡責務 §8a-F](#8a-f-builder-初期受け入れ-fixture-契約) に従い、path 重複、hidden / symlink、大小文字拡張子、file size、BOM、改行、既存出力不変を固定する builder fixture。 |
| <a id="align-26"></a>`ALIGN-26` | runner boundary safety / external I/O | [`components/runner.go`](../../components/runner.go) の process lock 解放は自身が書いた `pid` / `started_at` の再確認なしに path を削除し、stale 判定も Linux `/proc/{pid}` 契約を使用しない。GitHub API は `.repo_config` を使わず環境変数から repository を決定し、path / branch を URL encode せず、Trees API の truncated / directory target と response schema、rate-limit 最大待機、fake clock を固定契約どおり扱わない。source materialize は既存 `src` を `RemoveAll` して直接再作成し、precheck は出力親 directory を作成する。SSH deploy は remote 引数引用、`--`、分離 command、timeout、stream copy、file type / path 検証、NUL 終端 checksum 完全一致を実装せず、pending 保存失敗を無視する。 | [`docs/details/runner.md` 詳細本文責務 lock ファイル契約](runner.md#lock-ファイル契約)、[§13](runner.md#13-処理フロー)、[§14a](runner.md#14a-ssh-サイト転送)、[`docs/details/statefile.md`](statefile.md) の owner / collaborator 契約に従い、所有確認付き解放、repository 固定、GitHub retry / schema / clock、atomic materialize、precheck no-write、SSH argv / timeout / streaming / checksum、pending write failure を fake 外部依存で固定する runner fixture。 |
| <a id="align-27"></a>`ALIGN-27` | API read model / history consistency | [`components/api.go`](../../components/api.go) の status は `last_target_status` を欠き `output_url` と `queued` を追加し、状態読取優先順と fallback が固定契約と一致しない。sysinfo は `state_dir` / `output_exists` を追加し、stats、timeline、dashboard、output-meta、log search は key、集計、target / log 選択、日付・level filter、archive fallback、順序が各 response 契約と一致しない。history は top-level `warnings`、`failure_category` filter、重複 id 判定、厳格 schema、読込失敗を実装せず、comment / flag / tags 更新では no-op、複数状態の書込順、必須 config / audit log failure を固定境界で扱わない。 | [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1)、[§22.0e](api.md#sec-22-0e)、[`docs/details/statefile.md`](statefile.md)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](#sec-22-f) に従い、exact response、決定的選択 / 集計 / filter、通常 log / archive、破損 / 重複 / read failure、read-only no-write、history mutation の no-op / write order / partial failure を固定する API fixture。 |

<a id="align-05-unconnected-sdk-operations"></a>
**ALIGN-05 UI 未接続 SDK 操作：**

| 操作群 | [`admin/index.html`](../../admin/index.html) から呼び出されていない SDK public method |
|--------|------------------------------------------------------------|
| 状態・診断 | `health`、`getOutputMeta`、`getBuildTrends`、`getWebhookEvents`、`resetCircuitBreaker` |
| 承認・ビルド連鎖 | `approveBuild`、`rejectBuild`、`getBuildChainConfig`、`setBuildChainConfig` |
| schedule・pipeline | `setAllowedHours`、`clearAllowedHours`、`setForceInterval`、`setBuildCooldown`、`setScheduleInterval`、`setPipelineConfig` |
| history | `getHistoryComment`、`setHistoryComment`、`setHistoryFlag`、`setHistoryTags`、`rollbackHistory` |
| backup・restore | `backup`、`restore` |
| snapshot | `deleteSnapshot`、`downloadSnapshot` |
| token | `revokeToken` |
| hook・rule | `getHookLog`、`deleteHook`、`addAlertRule`、`deleteAlertRule`、`addTagRule`、`deleteTagRule` |

<a id="align-15-sdk-contract-gaps"></a>
**ALIGN-15 SDK 契約差分：**

| 対象 | 現在確認できる差分 |
|------|--------------------|
| `getAccessLog()` | `{limit=100,offset=0}` を受け取らず、必須 query を送信しない。 |
| `getNotifyLog()` | `{limit=100,offset=0}` を受け取らず、必須 query を送信しない。 |
| `getConfigLog()` | `{limit=100,offset=0}` を受け取らず、必須 query を送信しない。 |
| `getHistory()` | `failureCategory` を受け取らず、`failure_category` query を送信しない。 |
| `setRepoConfig()` | `{owner,repo}` の固定引数契約ではなく任意 object をそのまま body として送信する。 |
| `createToken()` | `{label,scopes,expires_at}` ではなく `{name,scopes,expires_at}` を送信する。 |
| request validation | 数値範囲、列挙値、空配列、object key の送信前検証が public method ごとの固定契約を満たさない。 |
| `streamBuild()` | `Content-Type: text/event-stream`、fatal UTF-8、frame の exact key / type / `at`、最終 `end`、EOF 後 buffer を検証せず、`done` Promise と normal / error / user close の一回だけの terminal transition を実装しない。 |
| JSON response | `Content-Type` を media type と parameter に分解せず、`includes("application/json")` または raw header の `endsWith("+json")` で判定する。 |

<a id="0g8-f-fixture--testdata--fake--実装検証証跡契約"></a>

**0g.8-F fixture / testdata / fake / 実装検証証跡契約：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は、実装完了判定に必要な fixture、fake、testdata、expected / effects、実装検証証跡、acceptance checklist、差し戻し条件を扱う。[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0e](../DETAIL_INDEX.md#0e-完全実装検証マトリクス) と [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) は詳細本文と証跡への入口、実装割当・順序・依存は [`docs/ROADMAP.md` 状態・計画責務 §4](../ROADMAP.md#4-phase-実装計画)、実装変更単位と着手条件は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#0f-phase-実装単位ポリシー)、setup / update の実行条件と Release asset 受け入れ条件は [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) を参照する。`release` owner component の詳細本文は未作成であり、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i.5](../DETAIL_INDEX.md#0i5-setup--release) を参照する。fixture 名、expected / effects、fake 動作、実装検証証跡項目、不足時の扱い、差し戻し条件だけを [`docs/details/fixture.md`](fixture.md) fixture 証跡責務で固定する。

実装検証証跡は、対象に応じて以下の 3 系統に分類する。複数系統にまたがる変更は、該当する全系統の証跡を実装検証証跡として記録する。

| 系統 | 対象 | 責務節 | 必須証跡 |
|------|------|--------|----------|
| component 実装 | builder、runner、api、admin、sdk、ui、statefile、archive、commitstatus、security、setup の実装。 | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) | owner component、collaborator component、変更ファイル、fixture / testdata path、fake、実行コマンド、期待結果、実結果、依存 component へ引き継ぐ contract。 |
| API 横断実装 | API と同期する SDK / UI / statefile の実装。 | [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) | endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、endpoint 固有の業務状態非変更、共通 security / observability 副作用、secret mask。 |
| [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の追加仕様化機能 | runner / security 詳細本文責務。 | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) | 対象 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様).x / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、関連 [`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) / [`docs/details/api.md` 詳細本文責務 §25](api.md#25-認証-実装仕様) / [`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様) / [`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) / [`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順)、owner / collaborator component、fixture 名、状態差分、外部副作用、partial failure、再実行、対象外確認。 |

**不足時共通扱い：**

[`docs/details/fixture.md`](fixture.md) fixture 証跡責務で必須とする fixture、manifest、expected、effects、security、実装検証証跡、対象外確認のいずれかが不足する場合、対象機能は未完了として扱う。fixture の pass だけでは完了証跡を満たさない。[`docs/details/fixture.md`](fixture.md) fixture 証跡責務内の対象別不足時表は、この不足時共通扱いに対する具体条件である。

component、API、[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) のいずれの実装検証証跡でも、記録形式は [`docs/details/fixture.md`](fixture.md) fixture 証跡責務の表に従う。owner component 別の [`docs/details/*.md`](../details/) 詳細本文責務、[`docs/DETAIL_INDEX.md`](../DETAIL_INDEX.md) 詳細仕様入口責務、[`docs/details/setup.md`](setup.md) 詳細本文責務に同種の記録項目がある場合でも、[`docs/details/fixture.md`](fixture.md) fixture 証跡責務は証跡分類、不足時の扱い、差し戻し条件だけを固定する。

<a id="sec-0g-8-f"></a>
**[fixture 証跡責務 §0g.8-F fixture / testdata 配置固定契約](fixture.md#sec-0g-8-f)：**

以下の配置は fixture 証跡責務上の配置契約である。未作成 path は、該当 component または該当 fixture の実装検証変更で作成するまで現行実体として扱わない。`release` と `mcp` は owner 詳細本文と fixture 契約が未作成であるため、対応する `testdata/` root を推測して作成してはならない。

| 検証群 | 必須配置 | 必須内容 | 禁止条件 |
|--------|----------|----------|----------|
| builder | `testdata/builder/single/`、`testdata/builder/site/`、`testdata/builder/empty-dir/`、`testdata/builder/strict/`、`testdata/builder/safe/`、各 fixture の `expected/`。 | 入力 Markdown、テーマ設定、asset 入力、期待 HTML / CSS / JS / search index、期待 stdout / stderr、期待終了コード。 | 実行環境ごとに変わる絶対 path、timestamp、乱数、外部 URL 取得結果を期待値へ含めてはならない。 |
| runner | `testdata/runner/r1/`〜`testdata/runner/r31/`、各 fixture の `state/`、`github/`、`pipeline/`、`ssh/`、`notify/`、`expected/`。 | GitHub fake response、状態ファイル初期値、lock 状態、pipeline fake 結果、deploy fake 結果、通知 fake 結果、期待 `.last_sha`、期待 queue / snapshot。 | 実 GitHub API、実 SSH、実通知先、実 remote branch 状態に依存して合否を決めてはならない。 |
| API request lifecycle | `testdata/api/request-lifecycle/auth/`、`status/`、`history/`、`logs/`、`queue/`、`stream/`、`errors/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | HTTP method / path / query / header / body、状態ファイル初期値、期待 response、期待 error body、SSE frame、状態 read/write 後の期待値。 | API 運用群 endpoint、外部公開設定、仕様未定義 endpoint を fixture に含めてはならない。 |
| API 運用 | `testdata/api/operations/config/`、`notify/`、`snapshots/`、`maintenance/`、`hooks/`、`tokens/`、各 fixture の `state/`、`requests/`、`responses/`、`expected/`。 | config / notify / snapshot / rollback / maintenance / hook / token の正常系、validation error、secret mask、API request lifecycle 群の回帰確認。 | token 原文、secret 原文、mask 前 payload、再取得不可 token の復元値を fixture または expected に含めてはならない。 |
| Admin | `testdata/admin/archive/`、`static-serving/`、`security/`、各 fixture の `input/`、`expected/`。 | archive entry、配布 file set、HTTP method / path / header / body、既存 admin directory の維持、secret path 非配信。 | UI / SDK 内容生成、未定義配布 file、unsafe archive entry、directory listing を許可してはならない。 |
| SDK | `testdata/sdk/request-shape/`、`error-shape/`、`stream/`、`binary/`、`operations/`。 | fake fetch transcript、期待 request、期待 SDK return、期待 `AdlaireCIError`、期待 stream event、timeout / abort の期待結果。 | Node.js 専用 API、bundler、npm package、実 network、browser storage 依存を検証前提にしてはならない。 |
| UI | `testdata/ui/login/`、`status/`、`build/`、`config/`、`secret/`、`stream/`、`operations/`。 | fake SDK script、入力 DOM 状態、操作手順、期待 DOM assertion、期待 SDK call、期待 disabled / loading / error / success 表示。 | 直接 `fetch()`、CDN、外部 framework、画像 snapshot だけの合否判定、secret 表示を含めてはならない。 |
| Statefile | `testdata/statefile/read/`、`write/`、`lock/`、`json-lines/`、`corrupt/`、`partial/`、各 fixture の `input/`、`expected/`。 | schema、mode、mtime、atomic write、lock、破損時処理、write order、forbidden write。 | caller 固有の業務判断、暗黙の自動修復、未定義状態 file を含めてはならない。 |
| Archive | `testdata/archive/log/`、`snapshot/`、`download/`、`delete/`、`rollback/`、各 fixture の `input/`、`expected/`。 | archive entry、checksum、圧縮・展開結果、stream、削除・rollback 境界、元 file 維持。 | unsafe entry、未検証展開、build 成否反転、元 build log 改変を許可してはならない。 |
| Commit status | `testdata/commitstatus/pending/`、`final/`、`disabled/`、`failure/`、各 fixture の `input/`、`expected/`。 | GitHub Status request、送信順、payload、失敗理由、build 成否非反転、secret mask。 | 実 GitHub write、Authorization 値保存、status 失敗による build 成否反転を含めてはならない。 |
| Security | `testdata/security/auth/`、`session/`、`token/`、`totp/`、`audit/`、`rate-limit/`、各 fixture の `input/`、`expected/`。 | memory-only state、hash-only state、scope、rate count、audit、one-time response、forbidden leak / write / call。 | password、token、ticket、TOTP secret、Authorization header の平文を expected に保存してはならない。 |
| Setup | `testdata/setup/install/`、`update/`、`rollback/`、`admin/`、`systemd/`、`health/`、各 fixture の `input/`、`expected/`。 | Release asset、checksum、binary / admin 配置、systemd 操作、Go `net/http` health、既存 state / secret 保持、rollback。 | 実 Release、実 systemd、実 network、未定義 release 生成・公開処理に依存してはならない。 |

<a id="sec-0g-8-f-2"></a>
**[fixture 証跡責務 §0g.8-F fake 実装固定契約](fixture.md#sec-0g-8-f-2)：**

| fake | 対象責務 | 必須動作 | 必須記録 |
|------|----------|----------|----------|
| fake GitHub server | runner | fixture の JSON 応答だけを返す。未定義 method / path は `404` とする。rate limit、`304`、`409`、`500` は fixture で明示された場合のみ返す。 | method、path、query、request body、認証 header の有無、呼び出し順。token 値は `***` に置換する。 |
| fake ssh executable | runner | fixture 指定の stdout、stderr、終了コード、timeout を返す。実 shell、実 SSH、実 file 転送は実行しない。 | argv、stdin 有無、環境変数名、終了コード、timeout 発生有無。secret 値は記録しない。 |
| fake notifier | runner / API 運用 | fixture 指定の HTTP status、response body、timeout を返す。通知先へ送信しない。 | URL の host 部分、payload schema、mask 後 payload、retry 回数、最終結果。 |
| fake filesystem | builder / runner / API | atomic write 失敗、sync 失敗、lock 競合、JSON 破損、permission error を fixture 単位で再現する。通常 file I/O の代替にはしない。 | 対象 path、操作種別、注入した失敗、復旧後の状態。 |
| fake fetch | SDK | `status`、`headers`、`body`、network error、timeout、abort、stream chunk を fixture どおり返す。実 network は使用しない。 | method、URL、query、headers、body 有無、abort 発生有無、呼び出し順。token 値は `***` に置換する。 |
| fake SDK | UI | SDK method ごとに固定 return、固定 throw、固定 stream event を返す。UI からの直接 API 呼び出しは受け付けない。 | method 名、引数、呼び出し順、throw した error code、stream unsubscribe 実行有無。 |

<a id="sec-0g-8-f-3"></a>
**[fixture 証跡責務 §0g.8-F 実装検証証跡固定契約](fixture.md#sec-0g-8-f-3)：**

この契約は全 component 実装の検証証跡に適用する。API 実装は [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約)、[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の追加仕様化機能実装は [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の acceptance checklist も同時に満たす。

| 証跡 | 必須記載 | 不足時の扱い |
|------|----------|--------------|
| 変更対象 | owner component、collaborator component、変更ファイル、追加 fixture / testdata path。 | 対象責務の成果物不足として未完了。 |
| 固定契約 | 追加または固定した CLI、状態 schema、HTTP API、SDK method、DOM id、fake 動作、終了コード、error body。 | 依存 component が参照できないため未完了。 |
| 検証 | 実行コマンド、fixture 名、期待結果、実結果、判定。 | 合否を再現できないため未完了。 |
| 未実装対象 | [`docs/ROADMAP.md`](../ROADMAP.md) 状態・計画責務で今回の実装対象外または将来計画に割り当てられた機能、MCP、外部公開設定、未定義 endpoint / UI / 状態ファイルを列挙する。 | 先取り実装または範囲不明として未完了。 |
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
| Fixture D: 冪等性 | 同一入力、同一 CLI 引数で 2 回連続実行する。 | 生成時刻 meta と footer 生成時刻以外の HTML、CSS、JavaScript、search index、`[REPORT]` が一致する。 |
| Fixture E: path 安全性と既存出力保護 | source 配下 out、source と out 同一、10 MiB 超 Markdown を入力する。 | exit `2`。stderr を固定し、出力作成なしまたは既存 `index.html` 維持。 |
| Fixture F: HTML escape と Markdown 境界 | raw `<script>`、先頭 h3、h3→h1、h1→h3、h3〜h6、`#no-space`、7 個の `#` で始まる行、escaped 内部 pipe、escaped 行末 pipe、不正 separator、不足 / 超過 cell を持つ table block、7 レベル以上 list nesting を含む Markdown を入力する。 | raw HTML は escape される。h1〜h6 と対応 class を出力し、先頭 h3 と h3→h1 では heading skip warning なし、h1→h3 だけ `[WARN] HEADING_SKIP` 1 件かつ `heading_skips=1` とする。`#no-space` と 7 個以上の `#` は通常段落にする。内部 escaped pipe は cell 文字、escaped 行末 pipe の行と不正 separator の候補行は各 1 件の段落、不足 cell は空 cell 補完、超過 cell は最終 cell へ ` | ` 連結とする。list nesting clamp と `[WARN] LIST_NESTING_CLAMPED` を固定する。 |
| Fixture G: search index / JavaScript contract | Fixture B と同じ directory 入力を使用する。 | `assets/search-index.json` の top-level array、entry key 順、body 長、HTML tag 除外、`assets/app.js` の localStorage guard、`search-results`、`data-search-hit`、外部 storage / network 不使用を固定する。 |
| Fixture H: strict warning before publish | strict 用 Markdown と既存正常出力を用意し、`adlaire-ci-build --src testdata/builder/strict/source.md --out <tmp> --strict` を実行する。 | exit `2`。stdout に `[WARN] BROKEN_LINK`、`[WARN] UNCLOSED_FENCE`、`[REPORT]` を出し、stderr は空。公開 rename を 0 回とし、既存 `index.html` を置換しない。 |
| Fixture I: atomic output compensation | 既存正常出力を用意し、fake filesystem で (1) 既存出力から `prev` への rename 後の `tmp` 公開 rename、(2) `tmp` 公開 rename 後の親 directory `Sync`、(3) 各 subcase の補償 rename または補償 `Sync` を個別に失敗させる。(3) の後は異なる current PID で再実行する。 | (1) と (2) で補償成功時は exit `1`、`[REPORT]` なし、更新前 `index.html` の byte 一致、`tmp` / `prev` 不在、補償 1 回を固定する。(3) は exit `1`、`cannot restore previous output directory` と最初の失敗文言、追加復旧 0 回、残存 `tmp` / `prev` / 公開 path の実状態を `expected/state/state-diff.json` と `expected/effects.json` に固定する。異なる PID の再実行も親 directory の残存 entry を削除せず、UTF-8 byte 列で最小の絶対 path を含む `output staging path already exists` で停止する。 |

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
| Fixture R1 | CLI help、relative state-dir、unknown option。 | help は exit `0`。不正 CLI は exit `2`、stdout 空、stderr 固定。 |
| Fixture R2 | `.last_sha={"sha":"blob-1"}`、fake GitHub Trees API が同一 SHA を返す。 | exit `0`。`.last_sha` 維持、log/history 非作成、`NO_CHANGE` INFO。 |
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
| Fixture R24 | branch target 2 件、1 件目 GitHub API failure、2 件目 success。 | exit `1`。1 件目 failure_api、2 件目 success。1 件目失敗で全体中断しない。 |
| Fixture R25 | branch target 2 件、両方 GitHub API failure。 | exit `3`。全 target failure_api、全 SHA 旧値維持、pipeline/deploy/snapshot 非実行。 |
| Fixture R26 | pipeline/deploy success、archive owner の snapshot save fake が `failed` を返す。 | exit `0`。build success 維持、build log / history の `snapshot_id=null`、`SNAPSHOT_SAVE_FAILED` warning、pending transfer 追加なし。 |
| Fixture R27 | success 保存後、finalizer `.build_state` atomic write だけ fake failure。 | exit `1`。success log/history/status は保持し、finalizer failure を固定する。 |
| Fixture R28 | manual active entry の build 中に process crash を発生させ、次回 runner を起動する。 | 初回は active を保持する。次回は waiting の urgent entry より先に同じ active payload を新しい build id で再実行し、結果確定後だけ active を消去する。 |
| Fixture R29 | queue build の final log atomic replace だけを失敗させる subcase と、final log 保存成功後の history append だけを失敗させる subcase を実行する。両 subcase で `.build_state` 読取と finalizer atomic write は成功させる。 | exit `1`。active と waiting を保持し、`running=false`、`current_build_id=null`、`last_finished_at={fake now}` を 1 回の atomic write で保存する。次回 runner は同じ active payload を新しい build id で waiting より先に再実行する。 |
| Fixture R30 | circuit open、active entry 1 件、waiting entry 1 件、stale `running=true` / `current_build_id`、pending transfer / notify 各 1 件で runner を起動する。 | pending transfer / notify retry だけを実行し、status は skipped / circuit_open。runtime flag は false / null へ保存し、active / waiting と最終時刻を保持する。queue 取得、build id 採番、GitHub、pipeline、deploy、snapshot、build log / history は実行せず、保存成功時 exit `0`、runtime flag 保存失敗時 exit `1`。 |
| Fixture R31 | SHA / deploy / snapshot 成功後の final build log atomic replace だけ fake failure。 | exit `1`。running build log を残し、history 非追記、新 SHA、deploy 結果、snapshot を巻き戻さず、status `failure_state_write`、running false、所有確認付き lock 解放を 1 回だけ実行することを固定する。 |
| Fixture R32 | `.repo_config` 不在、正常値、API 更新前に起動した process、更新後の次回 process、破損 JSON を個別に用意する。 | 不在時は `fqwink/Build-Scripts`、正常値は保存済み owner/repo を全 GitHub API path に使用する。起動済み process は更新前値を維持し、次回 process は更新後値を使う。破損時は exit `2`、`REPO_CONFIG_INVALID`、GitHub API / queue / build / pipeline / deploy / snapshot 呼出しと状態自動修復は 0 件。 |

<a id="22-f-api-fixture-契約"></a>

**22-F API fixture 契約：**

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) は、API の必須検証、fixture 名、入力状態、期待 response、期待副作用を扱う fixture 証跡責務である。API endpoint の method、path、request、response、error、read / write 境界は [`docs/details/api.md` 詳細本文責務 §22](api.md#22-バックエンド-api-仕様) を正本とする。

API 実装の検証証跡は、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに加えて、[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の検証群、endpoint、SDK method、UI 操作、状態 read/write、fixture 名、HTTP status、response、endpoint 固有の業務状態非変更、共通 security / observability 副作用、secret mask を記録する。

[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) における「状態差分なし」「no-write」「read-only」は、endpoint 固有の業務状態に対する禁止を意味する。[`docs/details/api.md` 詳細本文責務 §22.0 GET の副作用](api.md#sec-22-0) が許可する 5 群の共通副作用は、該当する場合に `expected/state/`、`expected/logs/`、`expected/effects.json` へ明示し、省略または禁止扱いにしてはならない。

実装順序と現在の割当は [`docs/ROADMAP.md` 状態・計画責務 §4.1](../ROADMAP.md#41-初期実装-phase-単位)、実装変更単位と完了判定単位は [`docs/SPEC.md` ポリシー責務 §0f](../SPEC.md#0f-phase-実装単位ポリシー) を参照する。[`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) は API fixture 証跡として記録する項目だけを固定する。

| 検証群 | 対象 | 完了条件 |
|--------|------|----------|
| API request lifecycle | 認証、session、共通 error、状態ファイル読み書き、`.access_log`、`.config_log`、build 操作、status、logs、history、queue、circuit breaker | `POST /api/login` から認証必須 API の共通処理、手動 build、強制 build、cancel、queue、history、log 取得までが [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0)〜[§22.0e](api.md#sec-22-0e) と一致し、秘密情報が log と response に出ない。 |
| API 運用 | config、repo、branch、schedule、PAT、diagnostics、dashboard、notify、SMTP、webhook、snapshot、rollback、maintenance、access control、hooks、alert rules、tag rules、pipeline config、notes、dashboard layout、tokens | 運用 API が schema どおり状態を保存し、secret mask、GET の endpoint 固有業務状態非変更、共通 security / observability 副作用、rollback / maintenance / hook / token の副作用が fixture と一致し、SDK と UI の操作名が [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e) と一致する。 |

各検証群の検証条件は以下とする。

| 検証群 | 必須検証 |
|--------|----------|
| API request lifecycle | 認証成功、認証失敗、期限切れ session、`401` 時 SDK token 破棄、`.access_log` 追記、秘密情報 mask、手動 build、force build、running 中の queue、cancel、history/log 取得、`409`、`429`、`503` を確認する。 |
| API 運用 | config/repo/branch/schedule の保存、`.config_log` と対象 audit の追記、GET 系 API の endpoint 固有業務状態非変更と共通副作用、Webhook test、weekly summary、SMTP test、webhook secret 保存、secret mask、snapshot list/download/delete、rollback、maintenance enable/disable、access control block、hook success/failure、rule 追加/削除、pipeline config 保存、notes 保存、dashboard layout 保存、token 発行/失効、token 本体が再取得不可であることを確認する。 |

<a id="sec-22-f"></a>
**[fixture 証跡責務 §22-F API request lifecycle fixture 固定契約](fixture.md#sec-22-f)：**

API request lifecycle 実装の完了証跡は [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いと [`docs/details/fixture.md` fixture 証跡責務 §22-F](fixture.md#22-f-api-fixture-契約) の固定表に従う。fixture は実装言語の test case 名または subtest 名へそのまま写せる粒度とし、期待 HTTP status、期待 body、状態ファイル副作用を同時に確認する。

| Fixture | 入力状態 / Request | 期待 response | 状態ファイル副作用 |
|---------|--------------------|---------------|--------------------|
| A1 route errors | 未定義 `/api/unknown`、既存 path への未許可 method を順に送る。 | `404 {"error":"Not found"}`、`405 {"error":"Method not allowed"}`。 | endpoint 固有の状態ファイルを作成、更新、削除しない。各 request の `.api_access_log` だけを共通処理順どおり 1 行追記し、access control、認証、rate limit、body 読取、外部呼び出しは開始しない。 |
| A1a authenticated body errors | 認証済み管理 session で body 禁止 endpoint への body 付き request、JSON 不正文、1 MiB 超過 body を個別に送る。 | `400 {"error":"Request body is not allowed"}`、`400 {"error":"Invalid JSON"}`、`413 {"error":"Payload too large"}`。 | access control、session、rate limit は [`docs/details/security.md` 詳細本文責務 §27.42〜§27.47](security.md#sec-27-42-2)、`.api_access_log` は [`docs/details/api.md` 詳細本文責務 §27.6](api.md#sec-27-6)、状態 read/write は [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) どおり行う。`.access_log`、`.audit_log`、endpoint 固有の状態ファイルを変更せず、外部 API、command、通知を呼び出さない。 |
| A2 auth errors | Bearer なし、無効 token、期限切れ session、scope 不足 API token で認証必須 endpoint を呼ぶ。 | `401 {"error":"Unauthorized"}` または `403 {"error":"Forbidden"}`。 | 秘密情報を response、`.access_log`、server log に出さない。 |
| A2a request id | fake entropy から固定 16 bytes を返す正常 request、未知 path、認証失敗 request を実行し、続けて entropy failure を発生させる。 | 正常系、`404`、`401` は 32 文字 lowercase hex の `X-Request-Id` を返し、対応する `.api_access_log` と同じ値を持つ。entropy failure は `500 {"error":"Internal server error"}` かつ `X-Request-Id` なし。 | 正常 request で audit / config log を作る場合は同じ request id を保存する。entropy failure は endpoint 状態、request log、audit log、config log を変更しない。 |
| A3 status primary | 正常な `.build_status.json`と lock 不在、valid active lock、stale lock を個別に用意して `GET /api/status`。 | [StatusObject 固定契約](api.md#status-object-contract) の 10 key だけを返す。lock 不在は status の `running`、valid active lock は `running:true`、stale lock は status の `running` を維持する。`output_url`、`queued` は返さない。 | 参照対象の業務状態ファイルは mtime、mode、内容を変更しない。認証、rate limit、`.api_access_log` の共通副作用は [`docs/details/api.md` 詳細本文責務 §22.0](api.md#sec-22-0) どおり記録する。 |
| A3a status conflict lock | 形式不正 lock と PID 判定不能 lock を個別に用意し、正常な `.build_status.json` で `GET /api/status`。 | どちらも `200`、StatusObject の `running:true`。その他の 9 key は status file の写像値。 | lock の削除、上書き、退避、status / state の修復を行わない。 |
| A4 status fallback | `.build_status.json` 不在、`.build_history` に status summary 対象の成功履歴、pending 各 1 件、circuit open、`.build_state.running=false`、`.build_lock` 不在で `GET /api/status`。 | history から `last_sha`、`last_build_at`、`last_build_status`、`last_target_status`、`last_trigger`、状態から両 pending 件数、`circuit_open`、`running:false` を算出し、`last_deploy_status:null`。必須 10 key 以外は返さない。 | `.build_status.json` を生成しない。queue の有無を status response に反映しない。 |
| A4a status fallback history-only exclusion | `.build_status.json` 不在、`.build_history` の末尾時刻に `approval_rejected`、その直前に status summary 対象の成功履歴を置いて `GET /api/status` と `GET /api/health`。 | history-only 行を除外し、成功履歴から `last_build_status` と `last_target_status` を算出する。summary 対象行がない別ケースでは `last_build_status:"none"`。 | `.build_status.json`、`.build_history`、approval 状態を変更しない。 |
| A5 status corrupted | `.build_status.json` を不正 JSON にして `GET /api/status`。 | `500 {"error":"State file is corrupted"}`。 | 退避ファイル作成、自動修復、初期値上書きを行わない。 |
| A6 history jsonl filtering | `.build_history` に有効行 2 件、空行 1 件、不正 JSON 1 件、必須 key 不足 1 件を置き `GET /api/history?page=1&per_page=10`。 | `total:2`、`pages:1`、`history` は有効行だけを `finished_at` 降順、同時刻は id 降順で返す。 | 壊れた行を書き戻し削除しない。server log に `BUILD_HISTORY_SKIP_CORRUPT`。 |
| A6a paged log query | `.access_log`、`.api_access_log`、`.notify_log`、`.config_log` に同一 `at` を含む有効行 102 件と破損行 1 件ずつを置き、各 GET を query なし、`limit=1&offset=1`、`limit=0`、`limit=1001`、`offset=-1`、整数でない値で呼ぶ。 | query なしは `at` 降順、同時刻は物理行順の逆順の先頭 100 件、`limit=1&offset=1` は同固定順の 2 件目だけを返す。API access log の `total` は paging 前の有効行数。範囲外または整数でない値は `422 {"error":"Validation failed","details":[...]}`。 | 参照対象 log を作成、修復、更新、削除せず、mtime、mode、内容を変更しない。破損行ごとに対応する `*_SKIP_CORRUPT` 固定 code と行番号だけを server log へ記録する。 |
| A6b history order and duplicate id | `.build_history` に同じ `finished_at` で id が異なる有効行、同じ id の後続有効行、壊れた行を置いて `GET /api/history`。 | `finished_at` 降順、同時刻は id 降順。重複 id は最初の物理行だけを返し、`total` と paging は一意な有効行から算出する。 | file を書き換えず、壊れた行に `BUILD_HISTORY_SKIP_CORRUPT`、重複 id 行に `BUILD_HISTORY_DUPLICATE_ID` と line number を記録する。 |
| A6c history failure category | 既知 category、`failure_category:null`、前方互換 category id `future_pipeline_error`、正規表現不一致値 `Future-Error` の履歴を置き、filter 未指定、`failure_category=pipeline_timeout`、未知 query で `GET /api/history`。 | 正規表現不一致行は破損行として除外する。未指定は前方互換 category id の原値を保持し、返却 page にその値があるときだけ top-level `warnings:["unknown_failure_category"]`。history item の integer `warnings` は変更しない。既知 filter は完全一致行だけ、未知 query は `422` と `details[].field="failure_category"`。 | `.build_history`、build log、status を変更しない。正規表現不一致行に `BUILD_HISTORY_SKIP_CORRUPT` を記録する。 |
| A7 build log lookup | `.build_logs/{id}.json` 正常、通常ログ不在で archive 正常、両方不在、破損ログをそれぞれ `GET /api/history/{id}/log`。 | 正常は log object、archive は gzip 展開結果、両方不在は `404 {"error":"Not found"}`、破損は `500 {"error":"State file is corrupted"}`。 | archive を通常ログへ復元しない。破損ログを上書きしない。 |
| A7a finite build stream | `.build_state` 不在、running/current id の `finished_at:null`, `duration_seconds:null` の途中保存 log あり・なし、running/current id 矛盾 state、idle で同時刻 id が異なる通常 / archive log、破損最新 log を用意し `GET /api/build/stream`。stdout に途中空行、末尾 LF、4001 code point の行を含める。 | state 不在は idle とする。running は current id だけを選び、`log.at=started_at`、end は `status:"running"`, `duration_seconds:null`。未保存は過去 log へ fallback せず `404`、矛盾 state は SSE header なしの JSON `500`。idle は `finished_at` 降順、id 降順、通常 log 優先で 1 件を選択する。各 frame は `data: {compact-json}\n\n`、stdout→stderr→warnings→error の log frame 後に end frame 1 件を送り close する。途中空行は保持、末尾空要素は除外、超過行は UTF-8 を壊さず 4000 code point へ切り詰める。 | wait、poll、tail、fallback write、archive 復元、状態更新を行わない。idle の破損候補は `BUILD_STREAM_SKIP_CORRUPT` を記録し次候補へ進む。 |
| A8 queue state | `.build_state` 不在、active / waiting がある正常状態、破損の 3 状態で `GET /api/queue`。 | 不在は `active:null,queued:[]`、正常は active と waiting を分離した保存値、破損は `500 {"error":"State file is corrupted"}`。 | 不在時も `.build_state` を生成しない。 |
| A9 build conflict / queue | valid running `.build_lock` と空き queue、形式不正 `.build_lock`、PID 判定不能 lock、queue 上限到達状態で `POST /api/build`。 | valid running かつ空きありは `202` queue 追加。判定不能 lock は `409 {"error":"Conflict"}`。queue 上限到達は `429`。 | `409` / `429` では `.build_state`、`.build_history`、`.build_status.json` を変更しない。valid running では queued だけを更新する。 |
| A10 circuit breaker | `.build_circuit_state` が `open:true` の状態で `POST /api/build`、続けて `POST /api/circuit-breaker/reset` を 2 回実行する。 | build は `409 {"error":"circuit_open"}`。reset は 2 回とも `200 {"message":"Circuit breaker reset","open":false,"consecutive_failures":0}`。 | build 拒否では queue、audit、runner 起動要求なし。1 回目の reset は初期値を atomic write し、`circuit_breaker_reset` config log と config update audit を順に 1 件ずつ追記する。2 回目は状態、config log、audit に差分なし。どちらも build / runner を起動しない。 |
| A11 write lock timeout | `.build_state.lock` を保持した状態で `.build_state` 更新 endpoint を呼ぶ。 | 10 秒経過後 `409 {"error":"Conflict"}`。 | tmp file を残さず、target を変更しない。 |
| A12 GET endpoint-state no-write | `GET /api/status`、`GET /api/history`、`GET /api/logs`、`GET /api/queue` を連続実行する。 | 各 endpoint は入力状態に応じた正常 response または固定 error response。 | 参照対象の業務状態ファイル一覧、mtime、mode、内容は request 前後で一致する。`.api_access_log` は request ごとに 1 行、認証 session と rate limit 状態は入力条件に応じて共通契約どおり更新し、それ以外の write / external call は行わない。 |
| A13 manual dispatch | active / waiting のない idle 状態で `POST /api/build`、fake systemctl exit `0`。 | `202 {message:"Build queued",queue_id,queued:true,dispatch:"requested"}`。`build_id` と `queued:false` は存在しない。 | manual entry、trigger audit を保存し、systemctl を固定引数で 1 回呼ぶ。running / current_build_id / active / SHA / status は変更しない。 |
| A14 dispatch fallback | idle 状態で `POST /api/build/force`、fake systemctl 非 `0` または 5 秒 timeout。 | `202`、`queued:true`、`dispatch:"timer_fallback"`。 | force entry と audit を保持し、`RUNNER_ACTIVATION_DEFERRED` を固定形式で 1 件記録する。retry、queue rollback、running 遷移、SHA 更新なし。 |
| A15 zero queue dispatch slot | `queue_max_size=0`、active なし、waiting なしで 1 回目と 2 回目の異なる manual request を送る。 | 1 回目は `202`、2 回目は `429 queue_full`。 | 1 件目だけ waiting に保存し、2 件目で state 差分なし。 |
| A16 health read failures | `.build_status.json` permission failure、fallback history read failure、`.pending_transfers` permission failure を個別および同時に fake して `GET /api/health`。 | HTTP `200`、固定順の `build_status_read_error`、`build_history_read_error`、`pending_transfers_read_error`、`last_build_status:"none"`、`pending_transfers:0`。response encode fake failure だけは `500` かつ `HealthObject` なし。 | chmod、backup、修復、初期化、状態 write を行わない。 |

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
| 認証 | auth-password-failure | 誤 password で `POST /api/login` | `401`、session/ticket なし、失敗回数 +1、`.access_log` と `.audit_log` に secret なし。 |
| 認証 | auth-login-lock | 連続 10 回失敗後の `POST /api/login` | `429`、password hash 検証なし、`.access_log` に `login_locked`、`.audit_log` に `permission_denied`。 |
| 認証 | auth-session-issued | TOTP 無効で password 成功 | token は response のみ、`.admin_credentials.login_count` +1、ログに token/hash/salt なし。 |
| 認証 | auth-session-expired | 期限切れ session で保護 API | `401`、対象 session 削除、`.access_log` と `.audit_log` は追記しない。 |
| 認証 | auth-password-change | password 変更成功 | 新 salt/hash、現 session 以外削除、`password_change` ログ、password/hash/salt 平文なし。 |
| 認証 | auth-log-write-failure | login 成功時に `.audit_log` 追記失敗 | `500`、session token を response しない。 |
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
| [§27.22](runner.md#sec-27-22) | `success-yaml-pipeline-file-priority`、`success-yaml-pipeline-inline`、`success-standard-builder-fallback-cache`、`failure-yaml-pipeline-parse`、`failure-yaml-pipeline-reserved-env`、`security-yaml-pipeline-secret-mask` | `.pipeline.yml`、inline YAML、pipeline config、branch env、step fake result。 | source 優先順、標準 builder argv、env merge、step log、終了コード、deploy / SHA 更新禁止条件。 |
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
| [§27.35](runner.md#sec-27-35) | `success-priority-urgent-first`、`success-priority-active-first`、`success-priority-created-seq-normalize`、`failure-priority-invalid`、`failure-priority-queue-full` | active / waiting entries、runner lock、API action。 | active 優先、priority 順、created_seq 正規化、waiting-to-active atomic move、`422` / `429`。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-timeout`、`success-failure-category-deploy`、`failure-failure-category-filter-invalid`、`security-failure-evidence-mask` | exit code、stderr、API error、status。 | category enum、filter `422`、純粋な success の `null`、deploy pending の `deploy_failure`、evidence 上限、secret mask。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-record`、`success-environment-builder-version-timeout`、`failure-environment-write`、`security-environment-secret-excluded` | OS/env/version/state path inputs。 | environment object、unknown fallback、basename 保存、secret 非保存、write failure。 |
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
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-path-param` | `DELETE /api/tokens/{id}` など path parameter 付き endpoint を正規化 path pattern で判定し、query string を scope 判定に使わないことを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `failure-scope-audit-failure` | scope 不足時の audit 追記失敗で `500`、対象 endpoint 未実行、request body 未評価、`.api_tokens` / endpoint 状態差分なしを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `security-scope-multi-scope` | 複数 scope token はいずれか 1 scope が endpoint group に一致した場合だけ許可し、未知 scope / 空 scopes は token record 破損 `500` とすることを固定する。 |
| [§27.42](security.md#sec-27-42) | `api_token_scope` | `noop-scope-health-webhook-exempt` | `GET /api/health` と `POST /api/webhook` は API token scope 判定対象外とし、未定義 route は認証 / rate limit / body parse 前に `404` / `405` を返すことを固定する。 |
| [§27.47](security.md#sec-27-47) | `api_rate_limit` | `security-rate-webhook-ip-only` | `POST /api/webhook` は body size と署名検証成功後、JSON parse 前に `trigger` group の IP key だけを判定・更新し、上限超過は `actor_type:"webhook"`、`actor_id:"webhook"` の audit 成功後の `429`、event log / queue / runner 起動要求の差分なしとすることを固定する。 |
| 認証共通 | `auth_common` | `security-auth-one-time-response` | session token、login ticket、API token 本体、TOTP setup secret、otpauth URI が許可された成功 response 1 回だけに出現し、以後の response / log / expected に残らないことを固定する。 |
| 認証共通 | `auth_common` | `failure-auth-log-before-token` | token / ticket / secret 返却前の access / audit log fake failure で `500`、one-time 値を response せず、平文保存なし、保存済み hash-only state の巻き戻し有無を固定する。 |
| 認証共通 | `auth_common` | `security-auth-memory-only` | process restart 相当で session、ticket、TOTP setup 仮 secret、login 失敗回数が消え、未定義永続 state file が作成されないことを固定する。 |
| 認証共通 | `auth_common` | `security-init-credentials-atomic` | `--init-credentials` の成功、予備確認時の既存あり、予備確認後かつ create-only lock 取得後の既存競合、相対 path、write failure、rand failure、rename 後 sync failure、create-only lock cleanup failure の stdout / stderr / exit code / file mode / target 非上書き / partial file を固定する。 |
| 認証共通 | `auth_common` | `security-auth-forbidden-plaintexts` | password、session token、API token、ticket、hash 算出入力、salt、TOTP code、TOTP secret が response / stdout / stderr / journal / state / expected に平文で出ないことを固定する。 |
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
| [§27.3](runner.md#sec-27-3) | `noop-retry-nonretryable` | pipeline exit code 非 0、checksum mismatch、validation failure 等の nonretryable 失敗で retry 0 件、pending transfer または failure 保存、SHA 更新なしを固定する。 |
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
| [§27.6](api.md#sec-27-6) | `failure-api-access-log-append` | access log append failure 時の HTTP response、対象 endpoint の状態変更有無、audit / server log、部分 JSON 行禁止を固定する。 |
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
| [§27.22](runner.md#sec-27-22) | `success-yaml-pipeline-file-priority` | `.pipeline.yml` 優先、inline YAML 未読、step 定義順実行、step log 定義順保存、標準 builder command 追加引数非適用を固定する。 |
| [§27.22](runner.md#sec-27-22) | `success-yaml-pipeline-inline` | file 不在時の inline YAML 採用、env merge、optional failure 継続、`PIPELINE_OPTIONAL_STEP_FAILED` の `warnings[]` と stdout 固定行、REPORT / build log / status success を固定する。 |
| [§27.22](runner.md#sec-27-22) | `success-standard-builder-fallback-cache` | repository file 不在、`inline_yaml:null` で標準 builder command を選択し、固定 argv 順、`build_cache_enabled` false / true ごとの `--cache-dir` 不在 / 1 回、`extra_args` 後置、shell 起動 0 件、出力 file set 検証を固定する。 |
| [§27.22](runner.md#sec-27-22) | `failure-yaml-pipeline-parse` | 禁止 YAML 構文で build 本体、deploy、snapshot、SHA cache 更新を開始せず、`failure_pipeline_config`、未実行 step `not_run` を固定する。 |
| [§27.22](runner.md#sec-27-22) | `failure-yaml-pipeline-reserved-env` | `.pipeline_config.env` と step env それぞれに [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Process environment entry 共通固定契約](../DETAIL_INDEX.md#process-environment-entry-contract) の各 exact reserved key と prefix 代表値 `ADLAIRE_CI_TEST` を個別に指定し、API `422` または runner `failure_pipeline_config` / 終了コード `2`、child process 0 件、deploy / snapshot / SHA cache 更新なしを固定する。branch env は [§27.31 fixture 固定契約](#sec-27-f-6) だけで検証する。 |
| [§27.22](runner.md#sec-27-22) | `security-yaml-pipeline-secret-mask` | step env / branch env / stdout / stderr / command args の secret 風値が log、notify、effects、response に平文で残らないことを固定する。 |
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
| [§27.35](runner.md#sec-27-35) | `success-priority-created-seq-normalize` | `created_seq` 欠落旧 entry の lock 内正規化保存、正規化後取り出し、正規化失敗時 build なしを固定する。 |
| [§27.35](runner.md#sec-27-35) | `failure-priority-invalid` | 不正 priority を API `422`、`.build_state` / history / log 差分なしに固定する。 |
| [§27.35](runner.md#sec-27-35) | `failure-priority-queue-full` | queue full 時 `429`、urgent でも既存 low entry を削除しないこと、write / command なしを固定する。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-timeout` | pipeline timeout を `pipeline_timeout`、evidence source / code / message / at、build log / history 保存一致に固定する。 |
| [§27.36](runner.md#sec-27-36) | `success-failure-category-deploy` | SSH / checksum / pending transfer failure の build log `target_status="success_deploy_pending"`、history `status="success_deploy_pending"`、API 正規化状態 `success`、両方の `failure_category="deploy_failure"`、分類優先順位、`failure_category=deploy_failure` の API filter 一致を固定する。 |
| [§27.36](runner.md#sec-27-36) | `failure-failure-category-filter-invalid` | 未知 `failure_category` query を `422`、状態差分なし、既存未知値 warning と区別することを固定する。 |
| [§27.36](runner.md#sec-27-36) | `security-failure-evidence-mask` | evidence 最大 10 件、分類 evidence 先頭、secret / token / path 全体 / 入力値連結なし、history `status="success"` の `failure_category:null`、history `status="success_deploy_pending"` の `failure_category:"deploy_failure"` を固定する。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-record` | build id 採番直後、builder 起動前の environment 保存、GOOS / GOARCH / Go version / hostname / state_dir / disk free / captured_at を固定する。 |
| [§27.37](runner.md#sec-27-37) | `success-environment-builder-version-timeout` | builder version 2 秒 timeout 時 `"unknown"`、stderr 非保存、build 継続を固定する。 |
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
| `input/files/` | file byte 列を入力にする fixture で必須 | `manifest.json.input_files` に列挙した Markdown、設定 JSON / YAML、archive、snapshot、hook file、Release asset、既存公開出力、remote artifact だけを同じ相対 path で配置する。 | `manifest.json.input_files` にない file、実外部サービスから取得した未固定 file。 |
| `input/fakes.json` | fake を 1 件以上使う fixture で必須 | [fake input root 固定契約](#fixture-fake-input-root-contract) の全 root key と、使用する fake の応答順。 | 未知 root key、実ネットワーク呼び出し前提、実 command 実行前提。 |
| `expected/response.json` | `manifest.json.assertions` に `response` がある場合に必須 | HTTP response の status、headers、body。schema は [HTTP response expected 固定契約](#fixture-http-response-expected-contract) に従う。 | SDK return/error、UI state、未定義 key、順序非決定配列。 |
| `expected/sdk_trace.json` | `manifest.json.assertions` に `sdk-trace` がある場合に必須 | SDK method から HTTP request までの呼出順、引数、request、結果種別。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | Authorization 値、secret 平文、API 外部副作用。 |
| `expected/sdk_return.json` | `manifest.json.assertions` に `sdk-return` がある場合に必須 | SDK success return と return 後 token state。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | SDK error、UI 表示値、補完済み response。 |
| `expected/sdk_error.json` | `manifest.json.assertions` に `sdk-error` がある場合に必須 | SDK error と error 後 token state。schema は [SDK expected 固定契約](#fixture-sdk-expected-contract) に従う。 | SDK success return、UI error text、secret 平文。 |
| `expected/ui_trace.json` | `manifest.json.assertions` に `ui-trace` がある場合に必須 | UI action、SDK call、refresh 順、disabled 遷移、field 消去。schema は [UI expected 固定契約](#fixture-ui-expected-contract) に従う。 | HTTP request、Authorization 値、secret 平文。 |
| `expected/ui_dom.json` | `manifest.json.assertions` に `ui-dom` がある場合に必須 | 表示・非表示 panel、text、field error、control state、one-time 領域。schema は [UI expected 固定契約](#fixture-ui-expected-contract) に従う。 | 生 DOM snapshot だけによる判定、secret 平文。 |
| `expected/stdout.txt` | `manifest.json.assertions` に `stdout` がある場合に必須 | stdout 完全一致。stdout なしを検証する場合は空ファイル。 | 現在時刻、絶対環境 path。 |
| `expected/stderr.txt` | `manifest.json.assertions` に `stderr` がある場合に必須 | stderr 完全一致。stderr なしを検証する場合は空ファイル。 | secret、実 token、実 URL credential。 |
| `expected/state/` | [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の `manifest.json.assertions` に `state` がある場合に必須 | 実行後状態ファイル一式、または [state diff expected 固定契約](#fixture-state-diff-expected-contract) に従う `expected/state/state-diff.json`。[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `state` は [§28-F ファイルセット固定契約](#sec-28-f-3) の `expected/site/` と `expected/builder-output.json` に割り当てる。 | 期待しないファイルの混入、実行後ファイル一式と `state-diff.json` の併用、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) での代用。 |
| `expected/logs/` | `manifest.json.assertions` に `logs` がある場合に必須 | build log、history、access log、audit log、notify log の期待値。 | secret 平文、実 Authorization header。 |
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
| `27.22` | `yaml_pipeline` | `runner` |
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
| `security-setup-admin-archive-boundary` | `A2` | `setup_admin_integration` | `admin` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `partial-setup-systemd-rollback-boundary` | `26.5` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `partial-setup-api-runner-dispatch` | `26.4` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `security-admin-static-serving` | `A3` | `setup_admin_integration` | `admin` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |
| `security-setup-secret-preservation` | `26.2b` | `setup_admin_integration` | `setup` | [setup / admin / Release asset 連動 fixture 固定表](#sec-27-f-19) |

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の fixture は、[§28-F カタログ固定契約](#sec-28-f-2) の同じ行にある fixture 名、節、feature slug を使用し、`owner_component` を `builder` に固定する。`§28` 共通行の `section` は `28-common` とし、`§28.1`〜`§28.25` 行は対応する `28.1`〜`28.25` とする。

<a id="sec-27-f-8"></a>
**[fixture 証跡責務共通 manifest schema 固定契約](fixture.md#sec-27-f-8)：**

[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) と [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `manifest.json` は次の共通 schema に従う。未知 key は禁止する。対象別の差分は各 fixture catalog と file set 契約で固定し、別 schema を作成してはならない。

<a id="fixture-component-identifier-contract"></a>
Fixture manifest の component 識別子は `builder`、`runner`、`api`、`admin`、`sdk`、`ui`、`statefile`、`archive`、`commitstatus`、`setup`、`security` の 11 件だけを許可する。`owner_component`、`collaborator_components`、`components` はこの識別子集合だけを使用する。

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
| `feature` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) で `name` に割り当てられた値との完全一致。[fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) と owner / 連動 fixture はレジストリ記載の snake_case、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) はカタログ記載の kebab-case をそのまま使用し、相互変換、別名、case 変更を禁止する。 |
| `category` | string | 必須 | `success`、`failure`、`partial`、`noop`、`security`。fixture 名 prefix と一致する。 |
| `owner_component` | string | 必須 | [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) で `name` に割り当てられた値との完全一致とし、[component 識別子固定契約](#fixture-component-identifier-contract) のいずれか 1 件を使用する。対象 component の列挙順や fixture 配置から推測してはならない。 |
| `collaborator_components` | array[string] | 必須 | [component 識別子固定契約](#fixture-component-identifier-contract) の値だけを使用する。`owner_component` を含めず、ASCII 昇順、重複なしとする。該当なしは空配列。 |
| `components` | array[string] | 必須 | `owner_component` 1 件と `collaborator_components` の全要素だけを ASCII 昇順、重複なしで含み、[component 識別子固定契約](#fixture-component-identifier-contract) の値だけを使用する。 |
| `references` | array[string] | 必須 | 各要素は repository root 起点の `docs/SPEC.md#<anchor>`、`docs/DESIGN.md#<anchor>`、`docs/DETAIL_INDEX.md#<anchor>`、`docs/details/<file>.md#<anchor>` のいずれかとし、実在 file と実在 anchor を指す。先頭は owner の主節、2 件目は対象の [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) または [fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の fixture カタログ、残りは collaborator と関連固定契約を UTF-8 byte 列の昇順で持つ。`docs/ROADMAP.md`、`docs/DOCUMENT_INDEX.md`、重複、fragment なし、絶対 path、`../` を禁止する。 |
| `fake_clock` | string/null | 必須 | UTC の `YYYY-MM-DDTHH:MM:SSZ` または `null`。カレンダー上無効な日時、offset、小数秒を禁止する。時刻依存 fixture は `null` 禁止。 |
| `not_applicable` | array[object] | 必須 | 各 object は `path` と `reason` の 2 key だけを持つ。`path` は対象 file set 契約の条件付き候補 file または directory の fixture root 相対 path、`reason` は空でない固定理由とする。`path` の ASCII 昇順、重複なしとし、必須 file、実在 path、契約外 path は列挙しない。該当なしは空配列とする。 |
| `missing_state` | array[string] | 必須 | fixture の仮想実行 cwd 起点の `/` 区切り相対 path で、実行前に存在しないことを期待する状態 file だけを ASCII 昇順、重複なしで持つ。directory、symlink、絶対 path、`.` / `..` segment を禁止する。該当なしは空配列とする。 |
| `input_files` | array[string] | 必須 | fixture directory からの `/` 区切り相対 path。対象 file set 契約が許可する `input/` 配下の通常 file だけを ASCII 昇順、重複なしで持ち、非空時は実在 input file set と完全一致させる。input file 不要時は空配列とする。directory、symlink、絶対 path、`.` / `..` segment を禁止する。 |
| `assertions` | array[string] | 必須 | `response`、`sdk-trace`、`sdk-return`、`sdk-error`、`ui-trace`、`ui-dom`、`stdout`、`stderr`、`state`、`logs`、`effects`、`secret-mask`、`order`、`idempotency`、`no-write` の 1 件以上。重複を禁止し、複数値はこの列挙順で記録する。 |

<a id="sec-27-f-9"></a>
**[fixture 証跡責務 §27-F 合否判定固定契約](fixture.md#sec-27-f-9)：**

| 判定 | 合格条件 |
|------|----------|
| response | status、headers、body、error details、request id が期待値と一致する。 |
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

[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表のうち `manifest.json.assertions` に含まれる判定が 1 つでも失敗した場合、その fixture は失敗とする。対象機能の必須 fixture が 1 件でも存在しない、または skip された場合、その機能の実装変更は未完了とする。

<a id="sec-27-f-10"></a>
**[fixture 証跡責務 §27-F assertion 選択固定契約](fixture.md#sec-27-f-10)：**

fixture の `manifest.json.assertions` は、実装者が任意に減らしてはならない。fixture 名 prefix、owner component、collaborator component に応じて、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の assertion を必ず含める。[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i](../DETAIL_INDEX.md#0i-詳細節対応表) から特定した対象 owner 機能契約で追加検証が必要な場合は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表へ追加してから fixture を作成する。

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
| `manifest.json.category` と fixture 名 prefix | `success-*` は `success`、`failure-*` は `failure`、`partial-*` は `partial`、`noop-*` は `noop`、`security-*` は `security` とする。 | 分類不一致として失敗。 |
| `manifest.json.section` と fixture 固定表 | section、fixture 名、owner の組み合わせは対象の [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) または [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) のカタログ固定表または owner fixture 固定表に存在する組み合わせだけ許可する。 | 固定表外として失敗。 |
| `manifest.json.owner_component` と `components` | `owner_component` は `components` に必ず含める。 | owner 責務不一致として失敗。 |
| `manifest.json.collaborator_components` と `components` | `collaborator_components` は `components` にすべて含め、`owner_component` を含めてはならない。 | collaborator 責務不一致として失敗。 |
| fixture 固定表の必須 input と入力ファイル | CLI / setup script 実行は `input/cli.json`、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の builder は `input/options.json`、HTTP request / SDK invocation / UI action は `input/request.json` を置く。その他の builder fixture は対象 file set 契約に従う。複数条件に該当する場合は必要 file をすべて置き、該当する条件を `not_applicable` で除外してはならない。 | 入力責務不一致として失敗。 |
| `manifest.json.assertions` と期待値ファイル | `response` は `expected/response.json`、`sdk-trace` は `expected/sdk_trace.json`、`sdk-return` は `expected/sdk_return.json`、`sdk-error` は `expected/sdk_error.json`、`ui-trace` は `expected/ui_trace.json`、`ui-dom` は `expected/ui_dom.json`、`stdout` は `expected/stdout.txt`、`stderr` は `expected/stderr.txt`、`logs` は `expected/logs/`、`effects` は `expected/effects.json`、`secret-mask` は `expected/security.json` を要求する。`state` は [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) では `expected/state/`、[fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) では `expected/site/` と `expected/builder-output.json` を要求する。 | 期待値不足として失敗。 |
| `expected/effects.json.write_order` と expected file set | `created_paths` / `updated_paths` として `write_order` に現れる path は [fixture 証跡責務 §27-F](#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の `expected/state/` / `expected/logs/`、または [fixture 証跡責務 §28-F](#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `expected/site/` / `expected/builder-output.json` で最終値を固定する。`deleted_paths` は対応する expected tree に存在させず、削除期待を `expected/effects.json` に固定する。 | 保存順だけの空検証として失敗。 |
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
| `admin` | admin archive の file list、entry validation、配置 mode、static serving path/header/body、secret 非配信、no mutation を fixture で固定する。 | setup の admin UI 展開、api の static serving、ui / sdk の配布物境界を壊さない expected を提供する。 | unsafe archive、未定義 file、directory listing、secret path、method denied、no mutation が検証済み。 | UI / SDK 内容生成、API endpoint 実装、状態 schema 変更、systemd 操作。 |
| `sdk` | method、args、query 生成、error 変換、token 破棄、binary / stream handling を fixture で固定する。 | UI が API 詳細を知らずに扱える戻り値と error をそのまま伝播する。 | API response に存在しない key を生成せず、`401` / `403` / `429` / network error が区別される。 | API response 推測補完、未定義 endpoint 呼び出し、状態ファイル直接操作。 |
| `ui` | SDK method 呼び出し、DOM 表示、disabled/loading/error、secret field 消去、再取得順を fixture で固定する。 | SDK 戻り値だけを表示し、API / 状態ファイルの内部構造を再解釈しない。 | 直接 API 呼び出し、状態ファイル操作、secret DOM 残存がない。 | 直接 `fetch()`、状態ファイル操作、外部 command 実行。 |
| `statefile` | schema、atomic write、JSON Lines、lock、破損時処理、保存順、no-write / forbidden write を fixture で固定する。 | runner / API / archive の保存対象ごとに `write_order`、`unchanged_paths`、`forbidden_writes` を固定する。 | read-only、dry-run、validation failure、partial failure の副作用境界が検証済み。 | component 固有の業務判断、UI 表示判断、API response 補完。 |
| `archive` | log archive、snapshot、download、delete、rollback の保存 / 取得 / 削除境界を fixture で固定する。 | API download / rollback response と runner history への影響を [`docs/details/statefile.md` 詳細本文責務 §22.0c](statefile.md#sec-22-0c) の契約に合わせる。 | 元 log 保持、archive 破損時挙動、安全でない entry 拒否、rollback 失敗時 no-write が検証済み。 | build 成否の反転、元 build log の改変、未検証 archive 展開。 |
| `commitstatus` | GitHub Commit Status payload、pending / final 送信順、失敗時非反転、secret mask を fixture で固定する。 | runner の build 結果を受け取り、build 成否を変更せず送信結果だけを返す。 | pending 失敗、final 失敗、commit SHA なし、無効時呼び出し 0 が検証済み。 | build 実行判断、dry-run での GitHub write、status 失敗による build 成否反転。 |
| `security` | token hash、session、TOTP、scope、audit、rate limit、secret mask、forbidden call/write を fixture で固定する。 | API / SDK / UI / runner の secret 表示、認証失敗、副作用境界を検証する。 | 認証失敗、権限拒否、rate limit、audit failure の副作用境界が検証済み。 | 業務処理代行、認可前状態更新、secret 平文保存。 |
| `setup` | binary 配置、service 更新、rollback、secret 既存値保持、stdout/stderr mask、終了コードを fixture で固定する。 | runner / API の初期状態と既存 secret を壊さないことを effects で固定する。 | 部分失敗時の復元対象と復元禁止副作用が `expected/effects.json` に明記済み。 | runtime 機能追加、状態 schema 暗黙変更、外部依存追加。 |

component 責務を複数変更へ分ける場合でも、各変更が満たすべき owner component、collaborator component、fixture 名、期待ファイル、禁止副作用を実装検証証跡に明記する。責務の所在が不明な場合は、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約) の不足時共通扱いに従う。

<a id="sec-27-f-15"></a>
**[fixture 証跡責務 §27-F runner / statefile 連動 fixture 固定契約](fixture.md#sec-27-f-15)：**

[§27.21〜§27.38](runner.md#sec-27-21) の runner owner 機能は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の対象機能行に定めた fixture に加えて、次の選択条件に該当する連動 fixture をすべて作成する。状態作成、更新、削除、複数保存順の変更は `success-runner-state-write-order`、途中の write / append / fsync failure を扱う変更は `partial-runner-state-write-failure`、同一入力または no-op の状態差分を扱う変更は `noop-runner-state-idempotency`、dry-run を扱う変更は `noop-runner-state-dry-run`、破損状態の read / recovery / stop を扱う変更は `failure-runner-state-corrupt-boundary`、secret を入力または状態に含む変更は `security-runner-state-secret-mask` を必須とする。複数条件に該当する場合は該当 fixture を省略せず、runner の業務判断と statefile の保存境界を分離して検証する。

| runner/statefile fixture 群 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|----------------------------|----------------|------------|---------------|----------|
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
| `failure-state-corrupt-object` | JSON parse 不能または未知 key あり | read adapter 呼び出し | `ErrStateCorrupted`。target 差分なし。API の公開応答は [`docs/details/api.md` 詳細本文責務 §22.0c.1](api.md#sec-22-0c-1) を参照する。 |
| `success-state-corrupt-regenerates` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](statefile.md#sec-22-0a) で再生成指定済み file が破損 | write caller または再生成を伴う操作 | corrupt backup が 1 件作成され、初期値だけが保存される。 |
| `failure-state-lock-timeout` | `{name}.lock` が 10 秒以上残る | write caller 呼び出し | conflict failure、target/tmp 差分なし。 |
| `failure-state-create-only-existing` | create-only target として通常 file、directory、symlink、その他の file type を個別に配置 | create-only write caller 呼び出し | target 内容を読まず `ErrStateAlreadyExists`、tmp 作成なし、target 差分なし、取得した lock だけを 1 回削除する。lock 削除失敗 subcase は `STATE_LOCK_CLEANUP_FAILED`、残存 lock、cleanup failure を固定する。 |
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

| api/sdk/ui fixture 群 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|------------------------|----------------|------------|---------------|----------|
| `success-api-sdk-ui-request-trace` | `api`、`sdk`、`ui` | UI user action、SDK fake fetch trace、API request fixture。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_return.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | UI → SDK → API の method / path / query / body が [`docs/details/api.md` 詳細本文責務 §22.0e](api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](ui.md#24-標準管理ツール-仕様) と一致する。UI → SDK → API の repository 内 call を `external_calls` へ記録せず、API 自身が外部 service を呼ばない場合は `external_calls` を空配列とする。 |
| `failure-api-sdk-ui-error-propagation` | `api`、`sdk`、`ui` | `401`、`403`、`409`、`422 details`、`429`、`500` の fake response。 | `expected/response.json`、`expected/sdk_trace.json`、`expected/sdk_error.json`、`expected/ui_trace.json`、`expected/ui_dom.json`、`expected/state/`、`expected/effects.json`、`expected/security.json`。 | status、message、details、token 破棄条件、panel error、field error、disabled が固定どおり。 |
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
| `success-ui-login-totp` | `login()` が `totp_required:true` を返す。 | password 消去、TOTP field 表示、ticket は DOM に表示しない。 |
| `success-ui-forced-password` | `must_change:"forced"`。 | password panel 以外が操作不可。変更成功後に通常初期取得を行う。 |
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

[`docs/details/setup.md` 詳細本文責務 §26](setup.md#26-セットアップアップデート手順) の setup、admin UI 配布、API service 導入、update、rollback を含む実装変更は、対象機能の owner fixture に加えて次の選択条件に該当する連動 fixture をすべて作成する。Release asset または admin archive layout を追加・変更する場合は `success-setup-admin-release-asset-layout`、unsafe archive の拒否境界を変更する場合は `security-setup-admin-archive-boundary`、systemd 更新失敗時 rollback を変更する場合は `partial-setup-systemd-rollback-boundary`、API service と runner dispatch / timer fallback を変更する場合は `partial-setup-api-runner-dispatch`、admin static serving を変更する場合は `security-admin-static-serving`、既存 secret の保持境界を変更する場合は `security-setup-secret-preservation` を必須とする。複数条件に該当する場合は該当 fixture を省略せず、[`docs/details/setup.md` 詳細本文責務 §26.8](setup.md#sec-26-8) と [`docs/details/admin.md` 詳細本文責務 §A1](admin.md#a1-管理-ui-静的ファイル境界)〜[§A6](admin.md#a6-admin-fixture-参照契約) の合格条件を同じ expected で検証する。

| setup/admin/Release asset fixture 群 | 対象 component | 必須 input | 必須 expected | 合格条件 |
|---------------------------------|----------------|------------|---------------|----------|
| `success-setup-admin-release-asset-layout` | `setup`、`admin` | Release asset 一式、`SHA256SUMS`、`admin-ui.tar.gz`、fake download response、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/state-diff.json`、`expected/effects.json`、`expected/security.json`。 | asset 名、checksum 対象、admin archive root layout、`index.html` と `adlaire-ci-sdk.js` だけを含む file set、file mode、directory mode が [state diff expected 固定契約](#fixture-state-diff-expected-contract) と配布正本に一致し、未定義 file を拒否する。 |
| `security-setup-admin-archive-boundary` | `setup`、`admin` | unsafe archive、既存 `$INSTALL_DIR/admin`、既存 API binary、API service fake、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/effects.json` の `unchanged_paths`、`forbidden_writes`、`forbidden_calls`、`expected/security.json`。 | unsafe archive では admin directory、API binary、credentials、runner state を変更せず、API service start / restart を呼ばない。 |
| `partial-setup-systemd-rollback-boundary` | `setup`、`runner`、`api` | systemd fake、全対象の旧 binary backup、旧 admin backup、API restart failure、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json` の `write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes`、`commands`、`expected/security.json`。 | API restart 失敗後の rollback を 1 回だけ実行し、旧 `admin/` と全対象 binary が更新前の version cohort に戻ることを固定する。`commands` は更新処理と rollback 処理を合わせ、runner timer restart 2 回、API restart 2 回、rollback 後の両 service の `is-active`、runner / API の 3 unit の `systemctl cat` を固定する。state、history、secret、credentials、systemd unit は変更せず、health HTTP request は実行しない。`response` assertion と `expected/response.json` を作成しない。 |
| `partial-setup-api-runner-dispatch` | `setup`、`runner`、`api` | API / runner / timer unit、systemd fake、manual queue request、`input/cli.json.expected_exit_code`。 | `expected/response.json`、`expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/effects.json` の `commands` と `write_order`、`expected/security.json`。 | API unit が timer を Wants / After し、queue 保存後だけ `systemctl start --no-block adlaire-ci.service` を 1 回実行する。systemctl failure でも queue を保持して timer fallback を返す。 |
| `security-admin-static-serving` | `admin`、`api` | static request、secret/state/log/snapshot path、method variation。 | `expected/response.json`、`expected/state/`、`expected/security.json`、`expected/effects.json`。 | A3 の status、header、body 有無、method 制限、no mutation、secret 非表示が一致する。 |
| `security-setup-secret-preservation` | `setup`、`security`、`statefile` | 既存 `.github_token` / `.last_sha` / その他 secret files、fresh setup と update 入力、create-only 競合、rename 前失敗、rename 後 partial failure、`input/cli.json.expected_exit_code`。 | `expected/stdout.txt`、`expected/stderr.txt`、`expected/state/`、`expected/logs/`、`expected/security.json`、`expected/effects.json` の `write_order`、`updated_paths`、`unchanged_paths`、`forbidden_writes`。 | fresh setup は `.github_token` → `.last_sha` の create-only 順序、新規 payload / mode / LF、競合後の read-only 再検証、既存有効 target の content / mode / mtime 保持、既存不正の無修復 exit `2`、statefile failure の exit `1`、rename 後 target 維持を固定する。update は initializer を呼ばず、全 secret / state の content / mode / mtime を保持する。stdout、stderr、journal、expected に secret 原文を残さない。 |

setup / admin / Release asset 連動 fixture の `manifest.json.name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) の fixture 名別割り当てに従う。`setup` と `admin` のうち owner ではない component を `collaborator_components` に含める。実装変更が API service 起動、static serving、rollback、secret 保持を扱う場合は、`api`、`runner`、`security` のうち当該 fixture の対象 component を collaborator として追加し、`expected/effects.json` の `forbidden_calls` と `forbidden_writes` に禁止副作用を明記する。

<a id="sec-27-f-20"></a>
**[fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](fixture.md#sec-27-f-20)：**

各 [`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の実装変更は、実装検証証跡に [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表を記録する。記録がない変更は、コードと fixture が存在しても未完了とする。

| 記録項目 | 必須内容 |
|----------|----------|
| wave | 対象 wave、対象 [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、先行 wave 完了 commit または変更識別子。 |
| 実装対象 | 実装する機能名、owner component、collaborator component、変更ファイル、追加 fixture path。 |
| 実装対象外 | 同じ wave 内で今回実装しない [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47)、後続 wave、MCP、外部公開構成、未定義 endpoint / UI / 状態ファイルを実装検証証跡に列挙する。 |
| fixture | [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) カタログの fixture 名、manifest / effects / security の検証結果。 |
| acceptance | [`docs/details/fixture.md`](fixture.md) fixture 証跡責務 runner / security 実装 acceptance checklist の各項目の pass / fail / 未実行。 |
| 後続影響 | 後続変更が利用許可済みの contract、利用禁止の未固定 contract を実装検証証跡に列挙する。 |

**runner / security 実装最終受け入れゲート：**

[`docs/details/runner.md` 詳細本文責務 §27](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の機能実装は、[`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) の固定表の全 gate を満たした場合だけ「完了」と判定する。1 件でも未達がある場合は「未完了」、仕様逸脱または secret 漏えいリスクがある場合は「差し戻し」とする。

| gate | 完了条件 | 未完了条件 | 差し戻し条件 |
|------|----------|------------|--------------|
| scope | 実装対象が [`docs/details/fixture.md` fixture 証跡責務 §27-F](fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) カタログ固定契約に存在する機能だけである。 | 実装対象節の記載が実装検証証跡にない。 | 未定義 endpoint、未定義 UI、未定義状態ファイル、MCP、外部公開構成を追加している。 |
| fixture | 対象 [`docs/details/runner.md` 詳細本文責務 §27.x](runner.md#27-runner-owner-追加仕様化機能-詳細仕様) / [`docs/details/security.md` 詳細本文責務 §27.42](security.md#sec-27-42)〜[§27.47](security.md#sec-27-47) の必須 fixture がすべて存在し、skip されていない。 | 必須 fixture が不足、または fixture 名が不一致。 | fixture が実装挙動に合わせて期待値を緩めている。 |
| manifest | 全 fixture の `manifest.json` が schema、assertion 選択、owner / collaborator component 責務を満たす。 | assertion、references、owner_component、collaborator_components、components、fake_clock のいずれかが不足。 | unknown key、実 secret、実環境 path、乱数依存を含む。 |
| expected | `response`、SDK、UI、stdout/stderr、state、logs、effects、security の各 assertion に対応する [ファイルセット固定契約](#sec-27-f-7) の expected file が存在し、各固定 schema と一致する。 | assertion に対応する expected file が不足。 | expected と manifest / effects / security が矛盾する。 |
| side effect | `write_order`、`unchanged_paths`、`forbidden_writes`、`forbidden_calls` が対象機能の成功 / 失敗 / no-op / partial を説明できる。 | 禁止副作用または無変更保証が不足。 | 失敗時に未許可状態を書き換える、外部呼び出しを行う。 |
| secret | secret 平文が expected、logs、effects、UI DOM、stdout/stderr に存在しない。 | secret 検証対象が不足。 | token、password、TOTP secret、PAT、Authorization header が平文で残る。 |
| component | builder / runner / api / admin / sdk / ui / statefile / archive / commitstatus / security / setup の該当責務が全て fixture に紐づく。 | owner / collaborator component の所在が不明。 | SDK / UI が API response を推測補完、または UI が直接 API / 状態ファイルを操作する。 |
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
| [`docs/details/builder.md` 詳細本文責務 §28.10](builder.md#sec-28-10) | `lazy-images` | `success-lazy-relative-image`、`success-lazy-external-image-no-fetch`、`success-lazy-data-uri-no-fetch`、`failure-lazy-base-outside-strict`、`noop-lazy-disabled`、`security-lazy-alt-escape`、`security-lazy-invalid-scheme-strict` |
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
| `failure-unclosed-math-strict` | 未閉鎖 math inline / math block が non-strict では通常 text、strict では終了コード `2` と `BUILDER28_UNRESOLVED_REFERENCE` になる。 |
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
| `markdown-extensions` | `failure-badge-invalid-text-strict` | strict で不正 badge label / color を `BUILDER28_INVALID_OPTION`、終了コード `2` にする。 |
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
| `lazy-images` | `failure-lazy-base-outside-strict` | strict で base 外相対 path を `BUILDER28_PATH_OUTSIDE_BASE`、終了コード `2`、公開出力維持にする。 |
| `lazy-images` | `noop-lazy-disabled` | option 無効時に `loading`、`decoding` を追加せず、既存 img 出力と一致する。 |
| `lazy-images` | `security-lazy-alt-escape` | alt、src、title 相当の attribute に raw HTML、quote、event handler が混入しても attribute escape される。 |
| `lazy-images` | `security-lazy-invalid-scheme-strict` | `javascript:`、`file:`、その他未許可 scheme を strict で `BUILDER28_INVALID_OPTION`、終了コード `2`、公開出力維持にする。 |

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
| `template-vars` | `failure-template-var-missing-strict` | strict で未定義 var を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、公開出力維持にする。 |
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
| `mermaid` | `failure-mermaid-unsupported-strict` | strict で未対応 Mermaid 構文を `BUILDER28_UNSUPPORTED_RESERVED`、終了コード `2`、stdout 空、stderr 固定 error、公開出力維持にする。 |
| `mermaid` | `noop-mermaid-disabled` | `--mermaid=false` で `mermaid` fence を通常 code block として出力し、SVG、Mermaid class、external script、REPORT rendered count を増やさない。 |
| `mermaid` | `security-mermaid-no-external-script` | SVG 内に `script`、`foreignObject`、event handler、external href、CDN、runtime fetch が存在しないことを `expected/security.json` で固定する。 |
| `footnotes` | `success-footnotes-multiple` | 複数 definition / reference、同一 id 複数参照、参照順番号、`sup.footnote-ref`、末尾 `section.footnotes`、REPORT count を固定する。 |
| `footnotes` | `success-footnotes-backlink` | 各 footnote item の `.footnote-backref`、本文 reference への backlink target、一意 id、search index から UI label を除外することを固定する。 |
| `footnotes` | `failure-footnote-undefined-strict` | strict で未定義 reference、未参照 definition、重複 definition を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout 空、stderr 固定 error、公開出力維持にする。 |
| `footnotes` | `security-footnote-escape` | definition text、reference 周辺 text、id、backlink label の raw HTML、quote、event handler、`javascript:` が escape される。 |
| `math` | `success-math-inline-block` | `$...$` と `$$...$$` を `span.math-inline` / `div.math-block` へ変換し、delimiter 除去、escape、REPORT inline / block count を固定する。 |
| `math` | `failure-math-unclosed-strict` | strict で未閉鎖 inline delimiter、未閉鎖 block delimiter、長さ超過を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、公開出力維持にする。 |
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
| `image-lightbox` | `failure-lightbox-alt-missing-strict` | strict で alt なし / 空 alt image を `BUILDER28_UNRESOLVED_REFERENCE`、終了コード `2`、stdout 空、stderr 固定 error、公開出力維持にする。 |
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
**[fixture 証跡責務 §28-F 既存出力互換固定契約](fixture.md#sec-28-f-18)：**

[`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) の実装変更では、対象機能の fixture だけでなく、機能を無効化した互換 fixture を 1 件以上含める。既定有効機能は、有効化前後ではなく「機能対象入力なし」の互換 fixture を含める。

| 機能種別 | 互換 fixture |
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
| [`docs/details/builder.md` 詳細本文責務 §28.3](builder.md#sec-28-3) | admonition type 正規化、badge color validation、disabled 時互換、escape。 |
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
| [`docs/details/builder.md` 詳細本文責務 §28.15](builder.md#sec-28-15) | byte count、pre/code 保持、attribute order 保持、structure validation、marker validation、disabled 互換、inline script / style 非追加。 |
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
| warning で継続できる構文不正 | 終了コード `0`、warning count 増加、fallback 出力あり。 | 終了コード `2`、出力なし。 | stdout / stderr / effects。 |
| base 外 path | 対象参照を無効化し warning。 | 終了コード `2`。 | 参照先 file が作成されないこと。 |
| 未定義参照 | 通常 text または非表示 fallback。 | 終了コード `2`。 | HTML fallback と strict 停止。 |
| reserved feature | 終了コード `2`。 | 終了コード `2`。 | strict 差分なし。 |
| 内部エラー fixture | 終了コード `1`。 | 終了コード `1`。 | 既存出力維持。 |

<a id="sec-28-f-21"></a>
**[`docs/details/fixture.md` fixture 証跡責務 §28-F manifest 固定 schema](fixture.md#sec-28-f-21)：**

[`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) の `manifest.json` は [fixture 証跡責務共通 manifest schema 固定契約](#sec-27-f-8) をそのまま使用し、builder 専用の別 schema を定義しない。`name`、`section`、`feature`、`owner_component` は [manifest 識別子レジストリ固定契約](#sec-27-f-manifest-identity) と [§28-F カタログ固定契約](#sec-28-f-2) の完全一致とする。`components` は `builder` と `collaborator_components` だけを ASCII 昇順、重複なしで持つ。`references` は対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) と [`docs/details/fixture.md` fixture 証跡責務 §28-F](fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) を含める。

`strict` と期待終了コードは `input/options.json` の正本値とし、`manifest.json` に `strict`、`expected_exit_code`、`feature_slug`、`owner`、`collaborators`、`spec_refs` を追加してはならない。

<a id="sec-28-f-22"></a>
**[fixture 証跡責務 §28-F builder 拡張実装検証証跡固定契約](fixture.md#sec-28-f-22)：**

実装検証証跡には、対象 [`docs/details/builder.md` 詳細本文責務 §28.x](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様)、追加 fixture 名、変更した HTML / CSS / JS / REPORT key、strict / non-strict 結果、外部依存なし確認、既存出力互換確認、未実装の [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 機能を列挙する。対象外の [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) 機能を先取り実装した場合、または [`docs/details/builder.md` 詳細本文責務 §28](builder.md#28-builder-owner-静的サイト出力拡張追加仕様化機能-詳細仕様) に存在しない Markdown 記法、CLI option、CSS class、JS 挙動を追加した場合は未完了として扱う。

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
| compatibility | 対象機能無効時または対象入力なし時の既存出力互換確認。 | 既存出力破壊リスクとして未完了。 |
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
