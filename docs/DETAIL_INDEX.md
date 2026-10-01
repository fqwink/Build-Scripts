# Adlaire CI — 詳細仕様入口

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は、対象機能から唯一の owner component を確定する対応表、owner 詳細本文が定義する collaborator 境界への参照入口、fixture 証跡への参照入口、共通固定値だけを管理する。方針・ポリシー・状態定義・着手可否は [`docs/SPEC.md` 方針責務・ポリシー責務](SPEC.md)、現在状態と実装計画は [`docs/ROADMAP.md` 状態・計画責務](ROADMAP.md)、実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md) を正本とする。

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は owner component の処理本文、HTTP body、状態 schema、SDK method、UI DOM、fixture assertion、現在状態を再定義しない。

<a id="詳細仕様参照入口"></a>
**対象機能から詳細本文への選択手順：**

1. [`docs/ROADMAP.md` 状態・計画責務](ROADMAP.md) で対象機能の現在状態と実装計画上の割当を確認する。
2. [詳細節対応表](#0i-詳細節対応表) で機能の owner component を一件に確定し、対応する詳細節を開く。
3. [詳細仕様参照表](#0b-詳細仕様参照表) で owner component の主本文が正しいことを確認する。
4. collaborator がある場合だけ、その component の詳細本文を境界確認として読む。
5. fixture と実装検証証跡は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) を読む。

<a id="0b-詳細仕様参照表"></a>
**owner component 別詳細本文・fixture 参照表：**

| owner component | 主本文 | fixture / 証跡 |
|-----------------|--------|----------------|
| `builder` | [`docs/details/builder.md` 詳細本文責務](details/builder.md) | [`docs/details/fixture.md` fixture 証跡責務 §8a-F](details/fixture.md#8a-f-builder-初期受け入れ-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §28-F](details/fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) |
| `runner` | [`docs/details/runner.md` 詳細本文責務](details/runner.md) | [`docs/details/fixture.md` fixture 証跡責務 §15a-F](details/fixture.md#15a-f-runner-初期受け入れ-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `api` | [`docs/details/api.md` 詳細本文責務](details/api.md) | [`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `admin` | [`docs/details/admin.md` 詳細本文責務](details/admin.md) | [`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約)、[`docs/details/fixture.md` fixture 証跡責務 Admin CLI fixture 固定契約](details/fixture.md#admin-cli-fixture-contract) |
| `sdk` | [`docs/details/sdk.md` 詳細本文責務](details/sdk.md) | [`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `ui` | [`docs/details/ui.md` 詳細本文責務](details/ui.md) | [`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `setup` | [`docs/details/setup.md` 詳細本文責務](details/setup.md) | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](details/fixture.md#0g8-f-fixture--testdata--fake--実装検証証跡契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture](details/fixture.md#sec-27-f-19) |
| `release` | [`docs/details/release.md` 詳細本文責務](details/release.md) | [`docs/details/fixture.md` fixture 証跡責務 Release fixture 固定契約](details/fixture.md#release-fixture-contract) |
| `statefile` | [`docs/details/statefile.md` 詳細本文責務](details/statefile.md) | [`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `archive` | [`docs/details/archive.md` 詳細本文責務](details/archive.md) | [`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `commitstatus` | [`docs/details/commitstatus.md` 詳細本文責務](details/commitstatus.md) | [`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `security` | [`docs/details/security.md` 詳細本文責務](details/security.md) | [`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| `mcp` | [`docs/details/mcp.md` 詳細本文責務](details/mcp.md) | [`docs/details/fixture.md` fixture 証跡責務 MCP fixture 固定契約](details/fixture.md#mcp-fixture-contract) |

<a id="0b1-owner-component-別-owner-collaborator-境界管理"></a>
**owner / collaborator 境界参照：**

owner / collaborator 境界の規則は [`docs/SPEC.md` 方針責務 §4.2a](SPEC.md#sec-4-2a) を参照する。[詳細節対応表](#0i-詳細節対応表) の `owner` 列は機能から owner component を特定する入口、[詳細仕様参照表](#0b-詳細仕様参照表) は owner component から主本文を特定する入口とする。

<a id="0c-実装前確認項目"></a>
**実装前参照：**

実装着手可否は、[`docs/SPEC.md` 方針責務 §4.7](SPEC.md#sec-4-7) の着手ゲート、[`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) の実装可否、[`docs/SPEC.md` ポリシー責務 §0d](SPEC.md#policy-spec-freeze) の凍結条件、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) の active Phase 条件、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の現在状態と Phase 割当てをすべて使用して判定する。この入口では、対象機能が [詳細節対応表](#0i-詳細節対応表) に存在し、owner 詳細本文と fixture 証跡へ到達できることだけを確認する。

<a id="0d-共通固定値"></a>
**共通固定値：**

| 項目 | 固定値 |
|------|--------|
| Go 最小バージョン | Go `1.22` 以上。 |
| 文字コード | JSON、JSON Lines、および owner 詳細本文が text と定義する入力・出力・HTTP body は UTF-8 とする。binary response、圧縮 archive、出力成果物内の通常 file は byte 列として扱い、文字コードを適用しない。状態ファイルの text / binary 区分は [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a) を正本とする。 |
| 改行 | 新規 JSON object / array と JSON Lines は LF とする。text の改行正規化、末尾 LF、byte 保持、CR / CRLF の許否は各 owner 詳細本文を正本とし、共通処理で変換しない。binary response、圧縮 archive、出力成果物内の通常 file には改行規則を適用しない。状態ファイルは [`docs/details/statefile.md` 詳細本文責務 UTF-8 text payload 固定契約](details/statefile.md#statefile-text-payload-contract) を参照する。 |
| <a id="common-machine-time"></a>機械処理時刻 | UTC の ISO 8601 秒精度 `YYYY-MM-DDTHH:MM:SSZ`。ミリ秒、ナノ秒、UTC 以外の offset、local timezone の保存を禁止する。ローカル時刻は UI 表示だけで使用する。 |
| CLI 終了コード | `0` 成功、`1` 一般エラー、`2` 入力・設定エラー、`3` 外部サービス・ネットワークエラー、`4` `.build_lock` の schema / PID 解析不正、PID 実行中判定不能、または lock 作成失敗。実行中 lock による通常 skip は `0`。 |
| CLI 共通 option | `--help` と `--version`。parse、優先順位、出力、副作用は [CLI 共通固定契約](#common-cli-contract) に従う。 |
| バイナリバージョン | 未注入の local / 検証用ビルドは `V.0.0-dev`。バージョン付き開発ビルドと Release 用実行バイナリ asset は割当済みバージョンと exact 一致する `^V\.[1-9][0-9]*\.[0-9]+$` 形式。Release 用実行バイナリ asset では tag とも exact 一致させる。その他の値、空文字、前後空白、改行を禁止する。 |
| 時刻ベース ID | prefix と UTC `YYYYMMDDHHmmss` を連結する。未衝突 ID に suffix は付けない。衝突時は `-001` から `-999` まで 3 桁連番を順に試し、上限到達時は既存 ID を上書きせず失敗とする。 |
| <a id="common-target-digest"></a>target SHA / digest | GitHub の単一 file target は Git object SHA の 40 文字 lowercase hex、GitHub の directory target と local target は SHA-256 の 64 文字 lowercase hex とする。算出 byte 列は GitHub target が [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、local target が [`docs/details/runner.md` 詳細本文責務 §27.23](details/runner.md#sec-27-23) を正本とする。空文字または `null` は owner schema が個別に許可する場合だけ使用する。commit SHA、GitHub Webhook `after`、commit status SHA はこの契約に含めず、40 文字 lowercase hex に固定する。互換性のため既存の `blob_sha`、`last_blob_sha`、`previous_blob_sha`、`current_blob_sha`、`before_sha`、`after_sha` の key 名は変更しない。 |
| 出力成果物 manifest SHA-256 | 出力 root 配下の通常 file だけを entry とし、`/` 区切りの相対 path を UTF-8 byte 辞書順に並べる。各 file の SHA-256 を lowercase hex で算出し、各 entry の `relative_path + "\n" + file_sha256 + "\n"` を順に連結した byte 列全体の SHA-256 lowercase hex を `output_sha256` とする。directory は走査だけに使用し entry に含めない。symlink、link count 2 以上の hardlink、device、socket、FIFO を 1 件でも検出した場合は除外継続せず算出失敗とする。出力 root は symlink でない directory、全 path は valid UTF-8 とし、先頭 `/`、空 segment、`.`、`..`、backslash、NUL、CR、LF を含む相対 path は算出失敗とする。file は no-follow open 後の identity / type と読取前後の size / mtime が列挙時から不変の場合だけ採用し、走査中の追加・削除・置換・変更は算出失敗とする。通常 file が 0 件の出力 root は空 byte 列の SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` とする。 |

<a id="common-cli-contract"></a>
**CLI 共通固定契約：**

| 項目 | 固定契約 |
|------|----------|
| option token | option は owner 詳細本文で列挙した `--[a-z][a-z0-9-]*` 形式だけを許可する。短縮 option、列挙外 option、owner 詳細本文が明示的に許可していない位置引数を禁止する。 |
| 値 option | 値を取る option は、owner 詳細本文が当該 option に `--name=value` を明示的に許可した場合を除き、`--name value` の 2 token 形式だけを許可する。値 token がない、または次 token が `--` で始まる場合は `missing value: --name` とする。未許可の `--name=value` は token 全体を未知 option とする。 |
| 共通 option 優先順位 | argv に exact `--help` が 1 件以上あれば `--help`、それ以外で exact `--version` が 1 件以上あれば `--version` を、残りの argv の parse、必須値検証、path / file / state 検証より先に確定する。両方がある場合は `--help` を採用する。 |
| argv token safety | `--help` / `--version` の優先順位で成功終了しない場合、owner 固有 parse の前に全 argv token を検証する。各 token は有効な UTF-8 文字列とし、NUL、CR、LF、C0 制御文字、DEL を含んではならない。不合格時は stdout 空、stderr `invalid command line token` + LF、終了コード `2` とし、token 原文を stdout、stderr、log、状態ファイルへ出力せず、状態変更と外部副作用を開始しない。 |
| help / version 結果 | owner 詳細本文の固定文字列 1 行と LF だけを stdout へ出力し、stderr は空、終了コードは `0` とする。file read/write、directory 作成、lock、listener、外部通信、child process、乱数取得を行わない。 |
| version 出力 | exact 3 token の `<binary-name> <binary-version> go=<runtime.Version()>` + LF とする。token 間は ASCII space 1 文字、前後空白、追加行、空 token を禁止する。`binary-name` は owner 詳細本文の実行ファイル名、`binary-version` は共通固定値のバイナリバージョン、`runtime.Version()` は空文字禁止とする。 |
| version 注入 | 実装はビルド時に不変の `binary-version` を受け取る。未指定時は `V.0.0-dev` とする。Release build は tag と同じ具体値を全バイナリへ注入し、出力が不一致または `V.0.0-dev` の場合は成果物作成を失敗させる。 |
| parse / 入力検証失敗 | stdout は空、stderr は owner 詳細本文で固定した最初のエラー 1 行と LF だけ、終了コードは `2` とする。owner 詳細本文で固定した検証順に最初の 1 件を選び、状態変更と外部副作用を開始しない。未知 option と禁止位置引数は、argv token safety 合格後の token だけを対象に `unknown option: <token>` とする。 |
| owner 固有契約 | 許可 option、固定 help / version 文字列、重複指定、値正規化、検証順、固有エラー、実行 mode は owner 詳細本文を正本とする。owner 詳細本文はこの共通契約を暗黙に上書きせず、異なる parse 形式を許可する option を個別に明示する。 |

<a id="common-state-dir-contract"></a>
**CLI state directory 共通固定契約：**

`--state-dir` を持つ owner CLI は、option の必須性、既定値、symlink、mode、owner、検証順、検証後の処理を各 owner 詳細本文で定義する。以下の共通失敗の stdout、stderr、終了コード、副作用はこの表だけを正本とし、owner 詳細本文で再定義しない。

| 条件 | stdout | stderr | 終了コード | 副作用 |
|------|--------|--------|------------|--------|
| 値が空文字 | 空 | `state directory must not be empty` + LF | `2` | state read/write、directory 作成、listener、外部通信、child process を開始しない。 |
| 相対 path | 空 | `state directory must be absolute: <path>` + LF | `2` | 同上。`<path>` は argv token safety 合格後の入力値を使用する。 |
| path 不在 | 空 | `state directory not found: <path>` + LF | `2` | 同上。 |
| directory でない | 空 | `state path is not directory: <path>` + LF | `2` | 同上。 |

<a id="process-environment-entry-contract"></a>
**Process environment entry 共通固定契約：**

environment object は 0〜100 key とする。各 key は `^[A-Z_][A-Z0-9_]{0,63}$`、各 value は UTF-8 の 0〜4096 bytes とし、NUL、LF、CR を禁止する。利用者入力または保存対象の environment object では、`PATH`、`HOME`、`SHELL`、`USER`、`GITHUB_TOKEN`、`ADLAIRE_TOKEN`、`ADLAIRE_CHANGED_TARGETS` と `ADLAIRE_CI_` prefix を reserved とし、保存、利用者指定値としての process 注入を禁止する。runner が所有する固定 process environment は [`docs/details/runner.md` 詳細本文責務 §27.31](details/runner.md#sec-27-31) の列挙値だけを例外とし、利用者入力による上書きを禁止する。永続化する environment object は key を ASCII 昇順で保存する。

状態ファイルの lock、atomic write、権限、JSON 処理は [`docs/details/statefile.md` 詳細本文責務](details/statefile.md)、秘密情報は [`docs/details/security.md` 詳細本文責務](details/security.md)、外部依存とデータ交換形式は [`docs/SPEC.md` 方針責務・ポリシー責務](SPEC.md) を正本とする。

<a id="0e-完全実装検証マトリクス"></a>
**完全実装検証参照：**

対象 owner の詳細本文と fixture 証跡は [owner component 別詳細本文・fixture 参照表](#0b-詳細仕様参照表) の同一行を使用する。完了判定と状態遷移は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity)、現在状態は [`docs/ROADMAP.md` 状態・計画責務](ROADMAP.md) を参照する。

<a id="phase-11-quality-gate-entry"></a>
**Phase 11 バグ修正ゼロ化参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 11 のバグ修正ゼロ化対象から該当する詳細本文と fixture 証跡への入口だけを固定する。Phase 11 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、完了判定方針と Phase 専用詳細仕様ファイルの扱いは [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) と [`docs/SPEC.md` 方針責務 §4.4](SPEC.md#sec-4-4)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit) を参照する。

| 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|------|--------------|--------------------------|
| Phase 11 仕様全般完了 gate | 方針責務 / 状態・計画責務 / 文書・実装ファイル所在の索引責務 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |
| 全標準実装 artifact source-code audit | 方針責務 / 詳細仕様入口責務 / 文書・実装ファイル所在の索引責務 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |
| 意味のあるテスト / boundary / failure matrix / concurrency / race / mutation zero survivor / harness self-verification gate | ポリシー責務 / fixture 証跡責務 | [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/details/fixture.md` fixture 証跡責務 test oracle evidence set 固定契約](details/fixture.md#test-oracle-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test isolation evidence set 固定契約](details/fixture.md#test-isolation-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract)、[`docs/details/fixture.md` fixture 証跡責務 mutation test evidence set 固定契約](details/fixture.md#mutation-test-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test / contract drift 証跡固定契約](details/fixture.md#test-contract-drift-evidence-contract) |
| test artifact traceability / drift source routing | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [`docs/details/fixture.md` fixture 証跡責務 test artifact traceability 固定契約](details/fixture.md#test-artifact-traceability-contract)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| non-dedicated owner test routing / artifact coverage zero gap | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [`docs/details/fixture.md` fixture 証跡責務 non-dedicated owner test routing 固定契約](details/fixture.md#non-dedicated-owner-test-routing-contract)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| fixture root coverage matrix / missing-root closure gate | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 / 状態・計画責務 | [`docs/details/fixture.md` fixture 証跡責務 fixture root coverage matrix 固定契約](details/fixture.md#fixture-root-coverage-matrix-contract)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| test execution evidence matrix / verification closure gate | fixture 証跡責務 / AGENTS.md 作業ルール / 文書・実装ファイル所在の索引責務 | [`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract)、[`docs/details/fixture.md` fixture 証跡責務 test oracle evidence set 固定契約](details/fixture.md#test-oracle-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test isolation evidence set 固定契約](details/fixture.md#test-isolation-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract)、[`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| test verification closure checklist / open item zero gate | fixture 証跡責務 / ポリシー責務 / 状態・計画責務 | [`docs/details/fixture.md` fixture 証跡責務 test verification closure checklist 固定契約](details/fixture.md#test-verification-closure-checklist-contract)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| test verification closure record schema / closure record set | fixture 証跡責務 / ポリシー責務 / 状態・計画責務 | [`docs/details/fixture.md` fixture 証跡責務 test verification closure record schema 固定契約](details/fixture.md#test-verification-closure-record-schema-contract)、[`docs/details/fixture.md` fixture 証跡責務 test verification closure record set 固定契約](details/fixture.md#test-verification-closure-record-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test verification closure checklist 固定契約](details/fixture.md#test-verification-closure-checklist-contract)、[`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| main dispatch / binary version regression gate | 起動入口 artifact / 方針責務 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](DETAIL_INDEX.md#0d-共通固定値)、[`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`main.go`](../main.go) |
| builder parser / config / output / generated site regression gate | `builder` / fixture 証跡責務 | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`§5`](details/builder.md#5-静的-web-サイト出力構造)、[`§8`](details/builder.md#8-実行方法)、[`docs/details/fixture.md` fixture 証跡責務 ALIGN-02](details/fixture.md#align-02)、[`ALIGN-07`](details/fixture.md#align-07)、[`ALIGN-29`](details/fixture.md#align-29)、[`docs/details/fixture.md` fixture 証跡責務 §28-F](details/fixture.md#28-f-fixture-証跡責務--builder-拡張実装検証証跡詳細契約) |
| 状態ファイル共通永続化 / direct runtime state write / JSON Lines elimination gate | `statefile` / fixture 証跡責務 | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`§22.0c`](details/statefile.md#sec-22-0c)、[`§22.0s`](details/statefile.md#sec-22-0s)、[`Phase 11 statefile バグ修正ゼロ化対象`](details/statefile.md#phase-11-statefile-quality-gate)、[`docs/details/fixture.md` fixture 証跡責務 ALIGN-17](details/fixture.md#align-17) |
| runner corrupt state recovery / startup repair gate | `runner` / `statefile` / fixture 証跡責務 | [`docs/details/runner.md` 詳細本文責務 §15a](details/runner.md#15a-runner-受け入れ検証条件)、[`docs/details/runner.md` 詳細本文責務 §27](details/runner.md#27-runner-owner-追加仕様化機能-詳細仕様)、[`docs/details/statefile.md` 詳細本文責務 Phase 11 statefile バグ修正ゼロ化対象](details/statefile.md#phase-11-statefile-quality-gate)、[`docs/details/fixture.md` fixture 証跡責務 §15a-F](details/fixture.md#15a-f-runner-初期受け入れ-fixture-契約) |
| runner source / pipeline / deploy / notification regression gate | `runner` / `security` / `archive` / fixture 証跡責務 | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`§13`](details/runner.md#13-処理フロー)、[`§14`](details/runner.md#runner-build-pipeline-execution)、[`§14b`](details/runner.md#14b-スナップショット管理)、[`docs/details/archive.md` 詳細本文責務](details/archive.md)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/fixture.md` fixture 証跡責務 ALIGN-33](details/fixture.md#align-33)、[`ALIGN-36`](details/fixture.md#align-36) |
| output manifest / file tree consistency gate | `runner` / `api` / `archive` | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](DETAIL_INDEX.md#0d-共通固定値)、[`docs/details/api.md` 詳細本文責務 出力サイトチェックサム](details/api.md#output-site-checksum)、[`docs/details/runner.md` 詳細本文責務 .build_history JSON Lines 追記契約](details/runner.md#build-history-json-lines-contract)、[`docs/details/archive.md` 詳細本文責務 保存済み archive 検証・配信契約](details/archive.md#snapshot-archive-validation) |
| secret / auth / one-time response regression gate | `security` / `api` / `sdk` / `ui` | [`docs/details/security.md` 詳細本文責務 認証共通詳細](details/security.md#認証共通詳細)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| queue / cancel / finalizer regression gate | `runner` / `api` / `statefile` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 runner finalizer 固定契約](details/runner.md#runner-finalizer-contract)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c) |
| API config / backup / restore / read model / filesystem / external side effect regression gate | `api` / `statefile` / `archive` / `security` / `runner` | [`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`§22.0d`](details/statefile.md#sec-22-0d)、[`docs/details/archive.md` 詳細本文責務](details/archive.md)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| API / MCP listener lifecycle and HTTP boundary gate | `api` / `mcp` / `security` | [`docs/details/api.md` 詳細本文責務 API listener lifecycle 固定契約](details/api.md#21a-管理-api-サーバー制限)、[`docs/details/api.md` 詳細本文責務 §22.0](details/api.md#sec-22-0)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)、[`docs/details/security.md` 詳細本文責務](details/security.md) |
| runtime clock / entropy / parallel worker determinism gate | `api` / `runner` / `mcp` / `setup` / `release` / fixture 証跡責務 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0d](DETAIL_INDEX.md#0d-共通固定値)、[`docs/details/fixture.md` fixture 証跡責務 fake adapter 接続固定契約](details/fixture.md#fixture-fake-adapter-binding-contract)、[`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract)、[`docs/details/api.md` 詳細本文責務 request ID 固定契約](details/api.md#api-request-id-contract)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0) |
| SDK / UI client runtime regression gate | `sdk` / `ui` / `api` / `security` | [`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/security.md` 詳細本文責務 認証共通詳細](details/security.md#認証共通詳細) |
| setup / release distribution / filesystem / GitHub release boundary gate | `setup` / `release` / `admin` / `security` / fixture 証跡責務 | [`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[`docs/details/admin.md` 詳細本文責務 §A1](details/admin.md#a1-管理-ui-静的ファイル境界)、[`§A2`](details/admin.md#a2-管理-ui-archive-検証)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/fixture.md` fixture 証跡責務 §27-F setup / admin / Release asset 連動 fixture](details/fixture.md#sec-27-f-19) |
| admin CLI / MCP bridge regression gate | `admin` / `sdk` / `ui` / `mcp` / `api` / `security` / `statefile` | [`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/security.md` 詳細本文責務 認証共通詳細](details/security.md#認証共通詳細)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c) |
| admin CLI / API / SDK / UI route parity gate | `admin` / `api` / `sdk` / `ui` | [`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/api.md` 詳細本文責務 §27.48〜§27.70](details/api.md#additional-management-api-contract)、[`docs/details/sdk.md` 詳細本文責務 §23.8](details/sdk.md#sec-23-8)、[`docs/details/ui.md` 詳細本文責務 §24.8](details/ui.md#sec-24-8)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |
| MCP protocol / SSE / sampling regression gate | `mcp` / `security` / fixture 証跡責務 | [`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/fixture.md` MCP fixture 固定契約](details/fixture.md#mcp-fixture-contract) |
| 正式 fixture directory harness / fixture root identity / manifest closure / test contract drift zero | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](details/fixture.md#sec-0g-8-f)、[`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract)、[`docs/details/fixture.md` fixture 証跡責務 mutation test evidence set 固定契約](details/fixture.md#mutation-test-evidence-set-contract)、[`docs/details/fixture.md` fixture 証跡責務 test / contract drift 証跡固定契約](details/fixture.md#test-contract-drift-evidence-contract)、[`docs/details/fixture.md` fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](details/fixture.md#sec-27-f-20)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md) |
| cross-owner regression gate | 方針責務 / fixture 証跡責務 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/details/fixture.md` fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](details/fixture.md#sec-27-f-20)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |

<a id="0f-仕様策定完了チェック"></a>
**仕様策定完了条件の参照：**

仕様策定または仕様改訂の完了条件は [`docs/SPEC.md` ポリシー責務 §0b](SPEC.md#policy-spec-pr-completion) を正本とする。この入口では、[詳細仕様参照表](#0b-詳細仕様参照表)、[`docs/SPEC.md` ポリシー責務 §0 詳細仕様必須項目](SPEC.md#detail-contract-required-fields)、[詳細節対応表](#0i-詳細節対応表)、[完全実装検証マトリクス](#0e-完全実装検証マトリクス) の参照が揃っていることだけを確認する。実装着手の凍結判定は [実装前参照](#0c-実装前確認項目) で別途行う。

<a id="0i-詳細節対応表"></a>
**機能・owner component・詳細本文対応表：**

[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i 詳細節対応表](DETAIL_INDEX.md#0i-詳細節対応表) は対象機能から唯一の owner と関連詳細本文へ移動するための対応表である。`owner` 列だけが機能 owner の正本であり、collaborator の接続境界と担当処理は owner の該当詳細節を参照する。現在状態は [`docs/ROADMAP.md` 状態・計画責務](ROADMAP.md)、受け入れ assertion は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) を正本とし、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 §0i 詳細節対応表](DETAIL_INDEX.md#0i-詳細節対応表) では再掲しない。

<a id="0i1-builder--静的-web-サイト出力"></a>
**0i.1 Builder / 静的 Web サイト出力：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| 出力サイトサイズ警告閾値 | `builder` | [`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 出力サイトへのビルドメタ埋め込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/builder.md` 詳細本文責務 §27.4](details/builder.md#sec-27-4) |
| builder 起動入口 / version | `builder` | [`docs/details/builder.md` 詳細本文責務 §0](details/builder.md#0-責務境界)、[`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| builder CLI parse / validation order | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| builder 入力 path / symlink / base-dir | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §2a](details/builder.md#2a-入力収集出力パス決定)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| 変換レポート出力 | `builder` | [`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| サイドバー開閉 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.2](details/builder.md#sec-7-2) |
| TOC 検索フィルター | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.4](details/builder.md#sec-7-4) |
| トップへ戻るボタン | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.7](details/builder.md#sec-7-7) |
| シンタックスハイライト | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.8](details/builder.md#sec-7-8) |
| 本文内全文検索 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.9](details/builder.md#sec-7-9) |
| アンカーリンク自動検証 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| コードブロックの折りたたみ | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.10](details/builder.md#sec-7-10) |
| 印刷スタイル（`@media print`） | `builder` | [`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/DESIGN.md` デザイン責務 Builder 拡張コンポーネント視覚契約](DESIGN.md#builder-拡張コンポーネント視覚契約) |
| 静的 Web サイト出力 | `builder` | [`docs/details/builder.md` 詳細本文責務 §2](details/builder.md#2-ファイルパス設定)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §7](details/builder.md#7-javascript-機能)、[`docs/DESIGN.md` デザイン責務](DESIGN.md) |
| テーマコンポーネント | `builder` | [`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §7](details/builder.md#7-javascript-機能)、[`docs/DESIGN.md` デザイン責務](DESIGN.md) |
| 外部リンクの自動処理 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3) |
| 読み取り進捗バー | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.13](details/builder.md#sec-7-13) |
| コードブロックのコピーボタン | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.6](details/builder.md#sec-7-6) |
| 見出しアンカーリンクコピー | `builder` | [`docs/details/builder.md` 詳細本文責務 §3](details/builder.md#3-処理パイプライン)、[`docs/details/builder.md` 詳細本文責務 §7.11](details/builder.md#sec-7-11) |
| TOC 開閉状態の永続化 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.3](details/builder.md#sec-7-3) |
| 見出しスラグ重複解決 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5) |
| 前後章ナビゲーションボタン | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §7.15](details/builder.md#sec-7-15) |
| 内部リンク整合性チェック | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| 見出し階層スキップ警告 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| 読了時間推計と表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.5](details/builder.md#sec-4-5)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §6](details/builder.md#6-css-クラス一覧)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法) |
| テーブルのソート機能 | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.14](details/builder.md#sec-7-14) |
| キーボードショートカット | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.12](details/builder.md#sec-7-12) |
| ビルドキャッシュ | `builder` | [`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14](details/runner.md#runner-build-pipeline-execution)、[`docs/details/builder.md` 詳細本文責務 §27.25](details/builder.md#sec-27-25) |
| 依存ファイルトラッキング | `builder` | [`docs/details/builder.md` 詳細本文責務 §4.3](details/builder.md#sec-4-3)、[`docs/details/builder.md` 詳細本文責務 §5](details/builder.md#5-静的-web-サイト出力構造)、[`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/builder.md` 詳細本文責務 §27.28](details/builder.md#sec-27-28) |
| 差分ビルド | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.1](details/builder.md#sec-28-1) |
| 複数出力形式 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.2](details/builder.md#sec-28-2) |
| Markdown 拡張記法サポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.3](details/builder.md#sec-28-3) |
| コードブロック行番号表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.4](details/builder.md#sec-28-4) |
| 見出しの自動採番 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.5](details/builder.md#sec-28-5) |
| セクション折りたたみ | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.6](details/builder.md#sec-28-6) |
| TOC 深さ制御 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.7](details/builder.md#sec-28-7) |
| 最終更新日の自動埋め込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.8](details/builder.md#sec-28-8) |
| diff ハイライト | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.9](details/builder.md#sec-28-9) |
| 画像の遅延読み込み | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.10](details/builder.md#sec-28-10) |
| カスタムメタタグ注入 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.11](details/builder.md#sec-28-11) |
| ライトモード固定 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.12](details/builder.md#sec-28-12) |
| コードブロックのファイル名表示 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.13](details/builder.md#sec-28-13) |
| テンプレート変数展開 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.14](details/builder.md#sec-28-14) |
| HTML ミニファイ | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.15](details/builder.md#sec-28-15) |
| TOC ハイライト追従（アクティブ見出し追跡） | `builder` | [`docs/details/builder.md` 詳細本文責務 §7.5](details/builder.md#sec-7-5)、[`docs/details/builder.md` 詳細本文責務 §28.16](details/builder.md#sec-28-16) |
| Mermaid ダイアグラム描画 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.17](details/builder.md#sec-28-17) |
| 脚注サポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.18](details/builder.md#sec-28-18) |
| インライン数式レンダリング | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.19](details/builder.md#sec-28-19) |
| ページ内ナビゲーション履歴 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.20](details/builder.md#sec-28-20) |
| 読み上げ対応（アクセシビリティ） | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.21](details/builder.md#sec-28-21) |
| 画像ライトボックス | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.22](details/builder.md#sec-28-22) |
| 印刷時 QR コード挿入 | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.23](details/builder.md#sec-28-23) |
| 定義リストサポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.24](details/builder.md#sec-28-24) |
| タスクリストサポート | `builder` | [`docs/details/builder.md` 詳細本文責務 §28.25](details/builder.md#sec-28-25) |

<a id="0i2-runner--ci-実行"></a>
**0i.2 Runner / CI 実行：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| ビルドタイムアウト | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドログのファイル保存 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| ネットワーク断時の再試行 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| GitHub API レート制限自動待機 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 転送後リモート整合性検証 | `runner` | [`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| マルチブランチビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドログ世代管理 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| ビルド出力の外部転送 | `runner` | [`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| ビルドクールダウン | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| ビルド前の事前チェック | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順) |
| 定期強制ビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルド中重複スキップ | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー) |
| GitHub PAT 有効期限の事前警告 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| コミット情報のビルドログ記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ) |
| GitHub API 連続失敗によるサーキットブレーカー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| 設定ファイル起動時整合性チェック | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/runner.md` 詳細本文責務 §27.10](details/runner.md#sec-27-10) |
| ビルドステータスファイル出力 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.8](details/runner.md#sec-27-8) |
| ビルドトリガー種別の記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/runner.md` 詳細本文責務 §27.9](details/runner.md#sec-27-9) |
| GitHub Commit Status API | `commitstatus` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/commitstatus.md` 詳細本文責務 §27.1](details/commitstatus.md#sec-27-1) |
| ドライラン実行モード | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.2](details/runner.md#sec-27-2) |
| ビルド失敗時の自動リトライ | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/runner.md` 詳細本文責務 §27.3](details/runner.md#sec-27-3) |
| Webhook 通知失敗リトライキュー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照) |
| 複数ファイル監視 | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.21](details/runner.md#sec-27-21) |
| 標準 builder command 拡張設定 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14](details/runner.md#runner-build-pipeline-execution)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.22](details/runner.md#sec-27-22) |
| ローカルファイル監視モード | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §27.23](details/runner.md#sec-27-23) |
| タグ付きコミットのみビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.24](details/runner.md#sec-27-24) |
| 並列マルチターゲットビルド | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.26](details/runner.md#sec-27-26) |
| ビルド前後フック | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.27](details/runner.md#sec-27-27) |
| リモートビルド対応 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §14a](details/runner.md#14a-ssh-サイト転送)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.29](details/runner.md#sec-27-29) |
| ブランチ別環境変数 | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.31](details/runner.md#sec-27-31) |
| ビルド通知連携 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.32](details/runner.md#sec-27-32) |
| ビルド時間トレンド記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.33](details/runner.md#sec-27-33) |
| ビルド依存チェーン | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.34](details/runner.md#sec-27-34) |
| ビルド優先度キュー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.35](details/runner.md#sec-27-35) |
| 失敗原因の自動分類 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.36](details/runner.md#sec-27-36) |
| ビルド実行環境の記録 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §27.37](details/runner.md#sec-27-37) |
| ビルド所要時間の異常検知 | `runner` | [`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.38](details/runner.md#sec-27-38) |

<a id="0i3-api--sdk--ui"></a>
**0i.3 API / SDK / UI：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| JavaScript SDK 公開契約 | `sdk` | [`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| 標準管理ツール UI 契約 | `ui` | [`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/DESIGN.md` デザイン責務 標準管理 UI 視覚契約](DESIGN.md#admin-ui-visual-contract)、[`docs/details/fixture.md` fixture 証跡責務 §22-F](details/fixture.md#22-f-api-fixture-契約)、[`docs/details/fixture.md` fixture 証跡責務 §27-F](details/fixture.md#27-f-fixture-証跡責務--runnersecurity-実装検証証跡詳細契約) |
| 管理 UI 静的配布物構成・Archive 検証 | `admin` | [`docs/details/admin.md` 詳細本文責務 §A1](details/admin.md#a1-管理-ui-静的ファイル境界)〜[§A2](details/admin.md#a2-管理-ui-archive-検証)、[`docs/details/admin.md` 詳細本文責務 §A4](details/admin.md#a4-setup-連携境界)〜[§A6](details/admin.md#a6-admin-fixture-参照契約)、[`docs/details/fixture.md` fixture 証跡責務 setup / admin / Release asset 連動契約](details/fixture.md#sec-27-f-19) |
| 管理 UI 静的 HTTP 配信 | `admin` | [`docs/details/admin.md` 詳細本文責務 §A3](details/admin.md#a3-静的配信契約)、[`docs/details/admin.md` 詳細本文責務 §A5](details/admin.md#a5-受け入れ条件)〜[§A6](details/admin.md#a6-admin-fixture-参照契約)、[`docs/details/fixture.md` fixture 証跡責務 setup / admin / Release asset 連動契約](details/fixture.md#sec-27-f-19) |
| ポーリング間隔の動的変更 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/api.md` 詳細本文責務 §27.11](details/api.md#sec-27-11) |
| GitHub Webhook 受信 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 Webhook 受信境界](details/api.md#webhook-receive-overview)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §27.12](details/api.md#sec-27-12) |
| 設定バリデーション API | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.5](details/api.md#sec-27-5) |
| API アクセスログ | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.6](details/api.md#sec-27-6) |
| Webhook イベントログ | `api` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 Webhook 受信境界](details/api.md#webhook-receive-overview)、[`docs/details/api.md` 詳細本文責務 §27.13](details/api.md#sec-27-13) |
| ヘルスチェックエンドポイント | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.16](details/api.md#sec-27-16) |
| Webhook イベント一覧取得 API | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.13](details/api.md#sec-27-13) |
| ビルドログ重大度フィルター | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §27.17](details/api.md#sec-27-17) |
| ブランチ設定の動的変更 API | `api` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.18](details/api.md#sec-27-18) |
| 週次ビルドサマリー Webhook | `runner` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.19](details/runner.md#sec-27-19) |
| 設定変更の詳細 diff 記録 | `api` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.20](details/api.md#sec-27-20) |
| ビルド承認フロー | `runner` | [`docs/details/runner.md` 詳細本文責務 §11](details/runner.md#11-ci-ランナー-ファイル構成)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/runner.md` 詳細本文責務 §16](details/runner.md#16-systemd-タイマー参照)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/api.md` 詳細本文責務 §27.30](details/api.md#sec-27-30) |
| 通知チャンネル管理 | `runner` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.32](details/runner.md#sec-27-32)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルドログの保存済み有限 SSE 配信 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルド統計ダッシュボード | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| IP アドレス制限 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| ビルドキュー可視化 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| メンテナンスモード | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| アラート閾値設定 | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |
| バックアップ／リストア | `api` | [`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様) |

<a id="0i4-statefile--archive--security"></a>
**0i.4 Statefile / Archive / Security：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| 状態ファイル共通永続化契約 | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/fixture.md` fixture 証跡責務 statefile owner fixture 固定契約](details/fixture.md#sec-27-f-16) |
| ビルドログのアーカイブ圧縮 | `archive` | [`docs/details/runner.md` 詳細本文責務 §12](details/runner.md#12-設定値runner)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/archive.md` 詳細本文責務 §27.7](details/archive.md#sec-27-7) |
| ビルド所要時間の記録と統計入力 | `runner` | [`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.14](details/runner.md#sec-27-14) |
| ビルドアーティファクト世代管理 | `archive` | [`docs/details/runner.md` 詳細本文責務 §14b](details/runner.md#14b-スナップショット管理)、[`docs/details/archive.md` 詳細本文責務 §27.15](details/archive.md#sec-27-15)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e) |
| ビルドアーティファクト管理 | `archive` | [`docs/details/runner.md` 詳細本文責務 §14b](details/runner.md#14b-スナップショット管理)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/archive.md` 詳細本文責務 §27.15](details/archive.md#sec-27-15) |
| ビルドトリガー専用 API スコープ | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.42](details/security.md#sec-27-42) |
| API キー管理 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.43](details/security.md#sec-27-43) |
| 監査ログ | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §27.30](details/runner.md#sec-27-30)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.44](details/security.md#sec-27-44) |
| セッションタイムアウト変更設定 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §25](details/api.md#25-認証-実装仕様)、[`docs/details/security.md` 詳細本文責務 §27.45](details/security.md#sec-27-45) |
| TOTP 二要素認証 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/api.md` 詳細本文責務 §25](details/api.md#25-認証-実装仕様)、[`docs/details/security.md` 詳細本文責務 §27.46](details/security.md#sec-27-46) |
| API レート制限 | `security` | [`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/security.md` 詳細本文責務 §27.47](details/security.md#sec-27-47) |

<a id="0i5-setup--release"></a>
**0i.5 Setup / Release：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| 初回セットアップ | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.1](details/setup.md#sec-26-1)〜[§26.3](details/setup.md#sec-26-3)、[§26.4.1](details/setup.md#sec-26-4-1)、[§26.7](details/setup.md#sec-26-7)〜[§26.8](details/setup.md#sec-26-8) |
| 管理 API 導入 | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.3b](details/setup.md#sec-26-3b)、[§26.4.2](details/setup.md#sec-26-4-2)、[§26.7](details/setup.md#sec-26-7)〜[§26.8](details/setup.md#sec-26-8) |
| バイナリアップデート | `setup` | [`docs/details/setup.md` 詳細本文責務 §26.5](details/setup.md#sec-26-5)、[§26.7](details/setup.md#sec-26-7)〜[§26.8](details/setup.md#sec-26-8) |
| アップデート rollback | `setup` | [`docs/details/setup.md` 詳細本文責務 アップデート rollback 固定契約](details/setup.md#setup-update-rollback-contract)、[§26.7](details/setup.md#sec-26-7)〜[§26.8](details/setup.md#sec-26-8) |
| GitHub Release 成果物生成・公開前検証・公開 | `release` | [`docs/details/release.md` 詳細本文責務 §R1](details/release.md#release-cli-contract)〜[§R7](details/release.md#release-acceptance-contract)。成果物の受け入れ側契約は[`docs/details/setup.md` 詳細本文責務 §26.2a](details/setup.md#sec-26-2a)と[§26.8](details/setup.md#sec-26-8)を参照する。 |

<a id="0i6-追加管理api機能"></a>
**0i.6 追加管理 API 機能：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| マルチユーザー対応 | `security` | [`docs/details/security.md` 詳細本文責務 §27.48](details/security.md#sec-27-48)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.48](details/api.md#sec-27-48) |
| データストア切り替え | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0d](details/statefile.md#sec-22-0d)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.49](details/api.md#sec-27-49) |
| 外部認証連携 | `security` | [`docs/details/security.md` 詳細本文責務 §27.50](details/security.md#sec-27-50)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.50](details/api.md#sec-27-50) |
| 統計データの JSON エクスポート | `api` | [`docs/details/api.md` 詳細本文責務 §27.51](details/api.md#sec-27-51)、SDK は [`docs/details/sdk.md` 詳細本文責務 §23.8](details/sdk.md#sec-23-8)、UI は [`docs/details/ui.md` 詳細本文責務 §24.8](details/ui.md#sec-24-8) |
| キュー内個別エントリのキャンセル | `runner` | [`docs/details/runner.md` 詳細本文責務 §27.35 queue 契約](details/runner.md#sec-27-35)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.52](details/api.md#sec-27-52) |
| Prometheus メトリクスエンドポイント | `api` | [`docs/details/api.md` 詳細本文責務 §27.53](details/api.md#sec-27-53) |
| CLI 管理クライアント | `admin` | [`docs/details/admin.md` 詳細本文責務 §A7](details/admin.md#sec-a7)、API 対応は [`docs/details/api.md` 詳細本文責務 §27.54](details/api.md#sec-27-54)、fixture 証跡は [`docs/details/fixture.md` Admin CLI fixture 固定契約](details/fixture.md#admin-cli-fixture-contract) |
| 設定の自動スナップショット | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0d](details/statefile.md#sec-22-0d)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.55](details/api.md#sec-27-55) |
| ステータスバッジ生成 | `api` | [`docs/details/api.md` 詳細本文責務 §27.56](details/api.md#sec-27-56) |
| ビルド履歴の自動削除設定 | `runner` | 状態 schema は [`docs/details/statefile.md` 詳細本文責務 HistoryRetentionPolicy object](details/statefile.md#history-retention-policy-object)、API 境界と削除実行順は [`docs/details/api.md` 詳細本文責務 §27.57](details/api.md#sec-27-57) |
| ロールベースアクセス制御 | `security` | [`docs/details/security.md` 詳細本文責務 §27.58](details/security.md#sec-27-58)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.58](details/api.md#sec-27-58) |
| 設定スナップショット差分表示 | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0d](details/statefile.md#sec-22-0d)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.59](details/api.md#sec-27-59) |
| 複数プロジェクト管理 | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0d](details/statefile.md#sec-22-0d)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.60](details/api.md#sec-27-60) |
| ユーザー管理 API | `security` | [`docs/details/security.md` 詳細本文責務 §27.61](details/security.md#sec-27-61)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.61](details/api.md#sec-27-61) |
| API バージョニング | `api` | [`docs/details/api.md` 詳細本文責務 §27.62](details/api.md#sec-27-62) |
| API ドキュメント自動生成 | `api` | [`docs/details/api.md` 詳細本文責務 §27.63](details/api.md#sec-27-63) |
| ビルドキューの手動並び替え | `runner` | [`docs/details/runner.md` 詳細本文責務 §27.35 queue 契約](details/runner.md#sec-27-35)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.64](details/api.md#sec-27-64) |
| 設定テンプレート | `statefile` | [`docs/details/statefile.md` 詳細本文責務 §22.0d](details/statefile.md#sec-22-0d)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.65](details/api.md#sec-27-65) |
| 管理者向けイベントフィード | `api` | [`docs/details/api.md` 詳細本文責務 §27.66](details/api.md#sec-27-66) |
| 読み取り専用共有リンク | `security` | [`docs/details/security.md` 詳細本文責務 §27.67](details/security.md#sec-27-67)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.67](details/api.md#sec-27-67) |
| API レスポンスキャッシュ制御 | `api` | [`docs/details/api.md` 詳細本文責務 §27.68](details/api.md#sec-27-68) |
| スナップショット間サイト差分 API | `archive` | [`docs/details/archive.md` 詳細本文責務 snapshot 契約](details/archive.md#sec-27-15)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.69](details/api.md#sec-27-69) |
| Webhook 送信履歴の手動再送 API | `runner` | [`docs/details/runner.md` 詳細本文責務 notification 契約](details/runner.md#sec-27-38)、API 境界は [`docs/details/api.md` 詳細本文責務 §27.70](details/api.md#sec-27-70) |

<a id="0i7-mcp"></a>
**0i.7 MCP：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| MCP サーバー実装 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)〜[§29.4](details/mcp.md#sec-29-4) |
| MCP ツール・リソース公開 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.5](details/mcp.md#sec-29-5)〜[§29.8](details/mcp.md#sec-29-8) |
| AI 支援ビルドエラー分析 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.10](details/mcp.md#sec-29-10) |
| MCP Prompts 定義 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.9](details/mcp.md#sec-29-9) |
| MCP Sampling によるビルドログ自動分析 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.10](details/mcp.md#sec-29-10) |
| MCP Notifications（イベントプッシュ） | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.11](details/mcp.md#sec-29-11) |
| MCP HTTP SSE transport 対応 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.12](details/mcp.md#sec-29-12) |
| MCP ツールスコープ細分化 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.13](details/mcp.md#sec-29-13) |
| MCP ツール呼び出し監査ログ | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.14](details/mcp.md#sec-29-14) |
| MCP リソース購読（Resource Subscriptions） | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.8](details/mcp.md#sec-29-8) |
| MCP クライアント情報ログ | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.14](details/mcp.md#sec-29-14) |
| MCP ツール実行統計 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.14](details/mcp.md#sec-29-14) |
| MCP ツール実行タイムアウト設定 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.15](details/mcp.md#sec-29-15) |
| MCP 設定 CRUD ツール | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.15](details/mcp.md#sec-29-15) |
| MCP Elicitation による副作用操作の確認 | `mcp` | [`docs/details/mcp.md` 詳細本文責務 §29.15](details/mcp.md#sec-29-15) |
