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
| `obsidian` | [`docs/details/obsidian.md` 詳細本文責務](details/obsidian.md) | [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence)、[`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) |

<a id="owner-detail-verification-route"></a>
**owner 詳細本文 検証接続共通入口：**

owner component 別詳細本文が検証接続を示す場合の参照入口である。検証方針と完了可否は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、owner から fixture への入口は [詳細仕様参照表](#0b-詳細仕様参照表)、横断テスト証跡は [横断テスト証跡共通入口](#cross-test-evidence-route)、Phase 11 横断入口は [Phase 11 バグ修正ゼロ化参照](#phase-11-quality-gate-entry)、Phase 12 横断入口は [Phase 12 実装品質ゲート再構築参照](#phase-12-quality-gate-entry)、Phase 13 横断入口は [Phase 13 実装整合・品質改善参照](#phase-13-implementation-alignment-quality-entry)、Phase 14 入口は [Phase 14 Obsidian Vault 連携参照](#phase-14-obsidian-vault-integration-entry)、Phase 15 入口は [Phase 15 Obsidian local vault 同期参照](#phase-15-obsidian-local-sync-entry)、Phase 16 入口は [Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照](#phase-16-quality-evidence-closure-entry)、Phase 17 入口は [Phase 17 ConoHa VPS 試験本番運用参照](#phase-17-production-validation-entry)、Phase 18 入口は [Phase 18 試験本番VPS 実運用接続・運用証跡参照](#phase-18-trial-production-operation-entry)、仕様全般不備は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract) と [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) を参照する。

<a id="cross-test-evidence-route"></a>
**横断テスト証跡共通入口：**

owner component 別詳細本文、作業ルール、または Pull Request 本文要件から横断テスト証跡へ到達する場合は、この入口を使用する。各 fixture 証跡の schema、必須 key、記録先、例外条件は [`docs/details/fixture.md` fixture 証跡責務](details/fixture.md) を正本とする。

| 対象 | 正本参照 |
|------|----------|
| テスト方針 / 完了可否 | [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) |
| traceability / owner routing / fixture coverage | [`docs/details/fixture.md` fixture 証跡責務 test artifact traceability 固定契約](details/fixture.md#test-artifact-traceability-contract)、[`non-dedicated owner test routing 固定契約`](details/fixture.md#non-dedicated-owner-test-routing-contract)、[`fixture root coverage matrix 固定契約`](details/fixture.md#fixture-root-coverage-matrix-contract) |
| test execution / closure | [`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract)、[`test verification closure checklist 固定契約`](details/fixture.md#test-verification-closure-checklist-contract)、[`test verification closure record schema 固定契約`](details/fixture.md#test-verification-closure-record-schema-contract)、[`test verification closure record set 固定契約`](details/fixture.md#test-verification-closure-record-set-contract) |
| test gap inventory / batch closure | [`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract)、[`test improvement batch closure 固定契約`](details/fixture.md#test-improvement-batch-closure-contract) |
| 仕様全般不備 inventory / batch closure | [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract)、[`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) |
| テスト固有 spec-gap closure | テストまたは fixture の完了可否に影響する仕様不足の記録先は [`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract) の `source=spec-gap`。`source=spec-gap` は [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) の `test_gap_connection` で `defect_id` と `gap_id` の 1 対 1 対応を閉じる。非テストの仕様全般不備の記録先は [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract) と [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract)。 |
| requirement / oracle / assertion / boundary | [`docs/details/fixture.md` fixture 証跡責務 test requirement coverage ledger 固定契約](details/fixture.md#test-requirement-coverage-ledger-contract)、[`test oracle evidence set 固定契約`](details/fixture.md#test-oracle-evidence-set-contract)、[`test assertion identity / failure diagnostics evidence set 固定契約`](details/fixture.md#test-assertion-failure-diagnostics-evidence-set-contract)、[`test assertion id 固定契約`](details/fixture.md#test-assertion-id-contract)、[`test boundary / failure matrix evidence set 固定契約`](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract) |
| isolation / determinism / concurrency | [`docs/details/fixture.md` fixture 証跡責務 test isolation evidence set 固定契約](details/fixture.md#test-isolation-evidence-set-contract)、[`test determinism evidence set 固定契約`](details/fixture.md#test-determinism-evidence-set-contract)、[`test concurrency / race evidence set 固定契約`](details/fixture.md#test-concurrency-race-evidence-set-contract)、[`race trigger matrix 固定契約`](details/fixture.md#race-trigger-matrix-contract) |
| mutation | [`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract)、[`mutation test evidence set 固定契約`](details/fixture.md#mutation-test-evidence-set-contract)、[`mutation selection ledger 固定契約`](details/fixture.md#mutation-selection-ledger-contract) |
| harness self-verification | [`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract) |
| contract drift | [`docs/details/fixture.md` fixture 証跡責務 test / contract drift 証跡固定契約](details/fixture.md#test-contract-drift-evidence-contract)、[`test / contract drift report schema 固定契約`](details/fixture.md#test-contract-drift-report-schema-contract) |
| test evidence package | [`docs/details/fixture.md` fixture 証跡責務 test evidence package 記録先固定契約](details/fixture.md#test-evidence-package-record-location-contract) |
| skip / 未実行 | [`docs/details/fixture.md` fixture 証跡責務 skip / 未実行証跡固定契約](details/fixture.md#test-skip-evidence-contract) |
| fixture root closure | [`docs/details/fixture.md` fixture 証跡責務 未作成 fixture root closure record 固定契約](details/fixture.md#fixture-root-missing-closure-record-contract) |
| manifest not_applicable boundary | [`docs/details/fixture.md` fixture 証跡責務 manifest not_applicable 境界固定契約](details/fixture.md#fixture-manifest-not-applicable-boundary-contract) |
| PR 証跡 | [`docs/details/fixture.md` fixture 証跡責務 implementation PR evidence template 固定契約](details/fixture.md#implementation-pr-evidence-template-contract)、[`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations) |

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
| 意味のあるテスト / test evidence package / test gap inventory / batch closure / requirement coverage / failure diagnostics / boundary / failure matrix / concurrency / race trigger / mutation selection / mutation zero survivor / harness self-verification gate | ポリシー責務 / fixture 証跡責務 | [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| test artifact traceability / drift source routing | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| non-dedicated owner test routing / artifact coverage zero gap | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| fixture root coverage matrix / missing-root closure gate | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 / 状態・計画責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| test execution evidence matrix / verification closure gate | fixture 証跡責務 / AGENTS.md 作業ルール / 文書・実装ファイル所在の索引責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| test verification closure checklist / open item zero gate | fixture 証跡責務 / ポリシー責務 / 状態・計画責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| test verification closure record schema / 18 record closure set | fixture 証跡責務 / ポリシー責務 / 状態・計画責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| 仕様全般不備一括棚卸し / 重複ゼロ / inspection pass gate | ポリシー責務 / 詳細仕様入口責務 / fixture 証跡責務 | [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract)、[`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract)、[横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract) |
| implementation PR evidence package / PR 提出証跡 gate | fixture 証跡責務 / AGENTS.md 作業ルール / ポリシー責務 / 状態・計画責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) |
| test evidence package / assertion id / skip evidence gate | fixture 証跡責務 / AGENTS.md 作業ルール / ポリシー責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations)、[`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test) |
| missing fixture root / manifest not applicable / drift report gate | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
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
| 正式 fixture directory harness / fixture root identity / manifest closure / test contract drift zero | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](details/fixture.md#sec-0g-8-f)、[`§27-F runner / security 実装検証証跡 必須記録固定契約`](details/fixture.md#sec-27-f-20)、[`Phase 11 fixture harness 参照`](details/fixture.md#phase-11-fixture-harness-reference)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md) |
| cross-owner regression gate | 方針責務 / fixture 証跡責務 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/details/fixture.md` fixture 証跡責務 §27-F runner / security 実装検証証跡 必須記録固定契約](details/fixture.md#sec-27-f-20)、[`docs/details/fixture.md` fixture 証跡責務 Phase 11 fixture harness 参照](details/fixture.md#phase-11-fixture-harness-reference) |

<a id="phase-12-quality-gate-entry"></a>
**Phase 12 実装品質ゲート再構築参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 12 の対象から該当する owner 詳細本文、fixture 証跡、方針責務、状態・計画責務への入口だけを固定する。Phase 12 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、Go 実装構造方針は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、テスト方針は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) を参照する。

<a id="phase-12-implementation-order"></a>
Phase 12 は以下の内部順序で実装、検証、証跡記録、状態更新を進める。前の順序に open item が残る場合、後続順序を完了扱いにしてはならない。

| 順序 | 固定する gate | 完了時の必須到達先 |
|------|---------------|--------------------|
| 1 | Phase 11 と状態ファイル共通永続化の再評価、必要な状態復帰 | [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) |
| 2 | API / SDK / UI / Admin CLI / setup / security の契約不整合、状態安全性、出力契約 | 下表の契約不整合・状態安全性対象、`phase12_contract_mismatch_count=0`、`phase12_state_safety_open_count=0` |
| 3 | production entrypoint 実行型 fixture harness、fake adapter、actual / expected 比較 | [`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence)、`phase12_fixture_harness_open_count=0` |
| 4 | production code mutation、race trigger、runner queue / finalizer 状態機械 | `phase12_mutation_survived_count=0`、`phase12_race_trigger_open_count=0` |
| 5 | Go owner package 5 ファイル固定、巨大 component file 分割、target path 整合 | [Phase 12 Go owner package target owner](#phase-12-go-owner-package-targets)、[`docs/DOCUMENT_INDEX.md` Phase 12 target path 所在](DOCUMENT_INDEX.md#phase-12-target-paths)、`phase12_owner_package_violation_count=0` |
| 6 | CI required checks、release reproducibility、version tag、release notes | `phase12_ci_required_check_open_count=0`、`phase12_release_reproducibility_open_count=0` |
| 7 | closure record set、対象外理由 anchor、状態復帰完了、PR 証跡 | `final_open_item_count=0`、[`docs/details/fixture.md` fixture 証跡責務 implementation PR evidence template 固定契約](details/fixture.md#implementation-pr-evidence-template-contract) |

<a id="phase-12-go-owner-package-targets"></a>
Phase 12 で Go owner package 5 ファイル固定の対象にする owner は下表だけとする。`sdk` と `ui` は JavaScript / HTML artifact を正本とするため、Go owner package 化の対象に含めない。

下表の `Phase 12 開始時の実装 artifact` は移行元を示す履歴情報であり、現行標準配置、実在ファイル、実装着手可否、Phase 13 target path の正本として扱ってはならない。現行標準配置は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、実在所在は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md) を正本とする。

| Go owner package target | Phase 12 開始時の実装 artifact | Phase 12 完了時の artifact | 主本文 |
|-------------------------|-------------------------------|---------------------------|--------|
| `builder` | `components/builder.go` | `components/builder/` の 5 ファイルだけ。`components/builder.go` と `components/builder_test.go` は残さない。 | [`docs/details/builder.md` 詳細本文責務](details/builder.md) |
| `runner` | `components/runner.go` | `components/runner/` の 5 ファイルだけ。`components/runner.go` と `components/runner_test.go` は残さない。 | [`docs/details/runner.md` 詳細本文責務](details/runner.md) |
| `api` | `components/api.go` | `components/api/` の 5 ファイルだけ。`components/api.go` と `components/api_test.go` は残さない。 | [`docs/details/api.md` 詳細本文責務](details/api.md) |
| `admin` | `components/admin.go` | `components/admin/` の 5 ファイルだけ。`components/admin.go` と `components/admin_test.go` は残さない。 | [`docs/details/admin.md` 詳細本文責務](details/admin.md) |
| `setup` | `components/setup.go` | `components/setup/` の 5 ファイルだけ。`components/setup.go` と `components/setup_test.go` は残さない。 | [`docs/details/setup.md` 詳細本文責務](details/setup.md) |
| `release` | `components/release.go` | `components/release/` の 5 ファイルだけ。`components/release.go` と `components/release_test.go` は残さない。 | [`docs/details/release.md` 詳細本文責務](details/release.md) |
| `mcp` | `components/mcp.go` | `components/mcp/` の 5 ファイルだけ。`components/mcp.go` と `components/mcp_test.go` は残さない。 | [`docs/details/mcp.md` 詳細本文責務](details/mcp.md) |
| `statefile` | `components/api.go` / `components/runner.go` / `components/mcp.go` 内の状態ファイル共通責務 | `components/statefile/` の 5 ファイルだけ。呼び出し元 owner package に状態ファイル共通責務の実装詳細を残さない。 | [`docs/details/statefile.md` 詳細本文責務](details/statefile.md) |
| `archive` | `components/api.go` / `components/admin.go` / `components/setup.go` / `components/release.go` 内の archive 責務 | `components/archive/` の 5 ファイルだけ。呼び出し元 owner package に archive 実装詳細を残さない。 | [`docs/details/archive.md` 詳細本文責務](details/archive.md) |
| `commitstatus` | `components/runner.go` 内の GitHub Commit Status 責務 | `components/commitstatus/` の 5 ファイルだけ。`runner` owner package に Commit Status 実装詳細を残さない。 | [`docs/details/commitstatus.md` 詳細本文責務](details/commitstatus.md) |
| `security` | `components/builder.go` / `components/runner.go` / `components/api.go` / `components/admin.go` / `components/setup.go` / `components/release.go` / `components/mcp.go` 内の security 責務 | `components/security/` の 5 ファイルだけ。呼び出し元 owner package に security 実装詳細を残さない。 | [`docs/details/security.md` 詳細本文責務](details/security.md) |

Go owner package の package 名、import path、`components/` 直下 `.go` file 残存禁止、`main.go` からの参照境界は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) を正本とする。Phase 12 の owner package 検査は、上表の各 owner について directory、5 file set、package declaration、import graph、旧 root artifact 不在、呼び出し元 package への duplicate 実装不在をすべて確認する。

| Phase 12 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|---------------|--------------|--------------------------|
| Phase 11 と状態ファイル共通永続化の再評価 | 状態・計画責務 / fixture 証跡責務 / `statefile` | [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) |
| API credential 初期化と setup stdout 契約統一 | `api` / `setup` / `security` / fixture 証跡責務 | [`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/security.md` 詳細本文責務 認証共通詳細](details/security.md#認証共通詳細)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| Admin CLI から呼ぶ全 URL と API route 表の自動照合 | `admin` / `api` / `sdk` / `ui` / fixture 証跡責務 | [`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/api.md` 詳細本文責務 追加管理 API 固定契約](details/api.md#additional-management-api-contract)、[`docs/details/sdk.md` 詳細本文責務 §23.8](details/sdk.md#sec-23-8)、[`docs/details/ui.md` 詳細本文責務 §24.8](details/ui.md#sec-24-8)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| SDK method、UI 操作、Admin CLI、API route、HTTP method、auth requirement、response schema の横断照合 | `api` / `sdk` / `ui` / `admin` / `security` / fixture 証跡責務 | [`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/sdk.md` 詳細本文責務 §23](details/sdk.md#23-javascript-sdk-仕様)、[`docs/details/ui.md` 詳細本文責務 §24](details/ui.md#24-標準管理ツール-仕様)、[`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| Statefile read-modify-write 全体のプロセス間 lock | `statefile` / `api` / `runner` / `mcp` / fixture 証跡責務 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| JSON Lines append の lock、flush、破損検出 | `statefile` / `runner` / `api` / fixture 証跡責務 | [`docs/details/statefile.md` 詳細本文責務 §22.0s](details/statefile.md#sec-22-0s)、[`docs/details/runner.md` 詳細本文責務 .build_history JSON Lines 追記契約](details/runner.md#build-history-json-lines-contract)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| API request、statefile、config、fixture manifest、expected の strict JSON schema | `api` / `statefile` / fixture 証跡責務 | [`docs/details/api.md` 詳細本文責務 §22.0](details/api.md#sec-22-0)、[`docs/details/statefile.md` 詳細本文責務 §22.0a](details/statefile.md#sec-22-0a)、[`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](details/fixture.md#sec-0g-8-f)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| secret redaction の API response、CLI stdout / stderr、log、fixture、release notes、audit log 横断照合 | `security` / `api` / `admin` / `runner` / `release` / fixture 証跡責務 | [`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/runner.md` 詳細本文責務 §15](details/runner.md#15-ログ)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| state、config、audit、history、snapshot、release、archive の atomic write、fsync、permission、crash recovery | `statefile` / `archive` / `runner` / `api` / `setup` / `release` / fixture 証跡責務 | [`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[`docs/details/archive.md` 詳細本文責務](details/archive.md)、[`docs/details/runner.md` 詳細本文責務 §14b](details/runner.md#14b-スナップショット管理)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| setup、release、runner、API mutation endpoint の冪等性 | `setup` / `release` / `runner` / `api` / fixture 証跡責務 | [`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| CLI exit code、stdout、stderr 統一 | 方針責務 / `builder` / `runner` / `admin` / `setup` / `release` / `mcp` | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 CLI 共通固定契約](#common-cli-contract)、[`docs/details/builder.md` 詳細本文責務 §8](details/builder.md#8-実行方法)、[`docs/details/runner.md` 詳細本文責務 §15a](details/runner.md#15a-runner-受け入れ検証条件)、[`docs/details/admin.md` 詳細本文責務](details/admin.md)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0) |
| fixture を実際に実行する共通 harness | fixture 証跡責務 / 全 owner | [`docs/details/fixture.md` fixture 証跡責務 §0g.8-F](details/fixture.md#sec-0g-8-f)、[横断テスト証跡共通入口](#cross-test-evidence-route)、[`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) |
| clock、filesystem、HTTP、command、systemd の依存注入 | fixture 証跡責務 / `runner` / `api` / `setup` / `release` / `mcp` | [`docs/details/fixture.md` fixture 証跡責務 fake adapter 接続固定契約](details/fixture.md#fixture-fake-adapter-binding-contract)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 request ID 固定契約](details/api.md#api-request-id-contract)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0) |
| fixture input を production entrypoint へ渡し、actual と expected を比較する gate | fixture 証跡責務 / 全 owner | [`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract)、[`test oracle evidence set 固定契約`](details/fixture.md#test-oracle-evidence-set-contract)、[`test harness self-verification evidence set 固定契約`](details/fixture.md#test-harness-self-verification-evidence-set-contract) |
| production code 対象 mutation test と `survived=0` | ポリシー責務 / fixture 証跡責務 / 全 owner | [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、[`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract)、[`mutation selection ledger 固定契約`](details/fixture.md#mutation-selection-ledger-contract) |
| Runner queue 遷移と finalizer の状態機械 | `runner` / `api` / `statefile` / fixture 証跡責務 | [`docs/details/runner.md` 詳細本文責務 runner finalizer 固定契約](details/runner.md#runner-finalizer-contract)、[`docs/details/runner.md` 詳細本文責務 §13](details/runner.md#13-処理フロー)、[`docs/details/api.md` 詳細本文責務 §22.0e](details/api.md#sec-22-0e)、[`docs/details/statefile.md` 詳細本文責務 §22.0c](details/statefile.md#sec-22-0c)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| GitHub Actions の通常 test、race、lint、fixture 実行 required check | AGENTS.md 作業ルール / fixture 証跡責務 | [`AGENTS.md` Git 運用ルール](../AGENTS.md#agents-git-operations)、[`docs/details/fixture.md` fixture 証跡責務 implementation PR evidence template 固定契約](details/fixture.md#implementation-pr-evidence-template-contract)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| 巨大 component file の owner package 5 ファイル固定分割 | 方針責務 / 文書・実装ファイル所在の索引責務 / 各 owner | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md)、各 owner の [`docs/details/*.md` 詳細本文責務](details/) |
| ROADMAP 文言検査 test と実装品質 test の分離 | 状態・計画責務 / fixture 証跡責務 | [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/details/fixture.md` fixture 証跡責務 test artifact traceability 固定契約](details/fixture.md#test-artifact-traceability-contract)、[横断テスト証跡共通入口](#cross-test-evidence-route) |
| root coverage fixture の inventory fixture 分類 | fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [`docs/details/fixture.md` fixture 証跡責務 fixture root coverage matrix 固定契約](details/fixture.md#fixture-root-coverage-matrix-contract)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 実装・テスト・fixture 所在](DOCUMENT_INDEX.md#実装ファイル一覧) |
| MCP no-op 実装の明示的未実装 error または実処理化 | `mcp` / `api` / `security` / fixture 証跡責務 | [`docs/details/mcp.md` 詳細本文責務 §29.0](details/mcp.md#sec-29-0)、[`docs/details/api.md` 詳細本文責務 §22](details/api.md#22-バックエンド-api-仕様)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/fixture.md` MCP fixture 固定契約](details/fixture.md#mcp-fixture-contract) |
| version tag と release notes 整備、release 再現性 gate | `release` / `setup` / `security` / fixture 証跡責務 | [`docs/SPEC.md` ポリシー責務 §1](SPEC.md#policy-versioning)、[`docs/details/release.md` 詳細本文責務 §R3](details/release.md#release-asset-contract)、[`docs/details/setup.md` 詳細本文責務 §26](details/setup.md#26-セットアップアップデート手順)、[`docs/details/security.md` 詳細本文責務](details/security.md)、[`docs/details/fixture.md` fixture 証跡責務 Release fixture 固定契約](details/fixture.md#release-fixture-contract) |
| Phase 12 完了 closure / 状態復帰 gate | 方針責務 / 状態・計画責務 / fixture 証跡責務 / 文書・実装ファイル所在の索引責務 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md) |

<a id="phase-13-implementation-alignment-quality-entry"></a>
**Phase 13 実装整合・品質改善参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 13 の対象から該当する owner 詳細本文、fixture 証跡、状態・計画責務への入口だけを固定する。Phase 13 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、Go 実装構造方針は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、テスト方針は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) を参照する。

Phase 13 は以下の内部順序で実装、検証、証跡記録、状態更新を進める。前の順序に open item が残る場合、後続順序を完了扱いにしてはならない。

| 順序 | 固定する gate | 完了時の必須到達先 |
|------|---------------|--------------------|
| 1 | owner package 5 ファイル責務純度、6 ファイル目、サブディレクトリ、空ファイル、ダミー実装、旧 path drift の解消 | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 13 target path 所在](DOCUMENT_INDEX.md#phase-13-target-paths)、`phase13_owner_file_violation_count=0`、`phase13_dummy_or_empty_file_count=0`、`phase13_document_drift_open_count=0` |
| 2 | statefile への状態 schema、process lock、read-modify-write、atomic write、JSON Lines、migration、復旧、直接状態更新禁止の集約 | [`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract)、`phase13_direct_state_mutation_count=0`、`phase13_state_safety_open_count=0`、`phase13_jsonl_corruption_open_count=0` |
| 3 | security、archive、commitstatus への横断責務集約と呼び出し元 duplicate 排除 | [`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract)、[`docs/details/archive.md` 詳細本文責務 Phase 13 archive 実装整合契約](details/archive.md#phase-13-archive-alignment-contract)、[`docs/details/commitstatus.md` 詳細本文責務 Phase 13 commitstatus 実装整合契約](details/commitstatus.md#phase-13-commitstatus-alignment-contract)、`phase13_security_kdf_open_count=0`、`phase13_token_arg_open_count=0`、`phase13_archive_commitstatus_ownership_open_count=0` |
| 4 | API、Admin、SDK、UI、MCP、setup の method、path、query、request、response、stdout、credential 初期化、token 取得、HTTP lifecycle の契約一致 | [`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/admin.md` 詳細本文責務 Phase 13 admin 実装整合契約](details/admin.md#phase-13-admin-alignment-contract)、[`docs/details/sdk.md` 詳細本文責務 Phase 13 SDK 実装整合契約](details/sdk.md#phase-13-sdk-alignment-contract)、[`docs/details/ui.md` 詳細本文責務 Phase 13 UI 実装整合契約](details/ui.md#phase-13-ui-alignment-contract)、[`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、`phase13_contract_mismatch_count=0`、`phase13_mcp_unimplemented_success_count=0` |
| 5 | runner queue、finalizer、active recovery、at-least-once、backup / restore transaction、systemd、SSH、Webhook external boundary を実装契約へ一致させる | [`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、`phase13_queue_recovery_open_count=0`、`phase13_external_boundary_open_count=0`、`phase13_required_log_write_ignore_count=0` |
| 6 | fixture execution、mutation、race、fault injection、integration、E2E、CI required checks、GitHub Actions pinning、release evidence、recovery procedure を closure へ接続する | [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence)、[`docs/details/release.md` 詳細本文責務 Phase 13 release 実装整合契約](details/release.md#phase-13-release-alignment-contract)、`phase13_fixture_execution_gap_count=0`、`phase13_mutation_survived_count=0`、`phase13_race_or_concurrency_open_count=0`、`phase13_fault_injection_open_count=0`、`phase13_e2e_open_count=0`、`phase13_ci_required_check_open_count=0`、`phase13_action_pin_open_count=0`、`phase13_release_evidence_open_count=0`、`phase13_recovery_procedure_open_count=0` |
| 7 | closure record set、対象外理由 anchor、ROADMAP 状態、DOCUMENT_INDEX 所在、PR 証跡 | [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務](DOCUMENT_INDEX.md)、`final_open_item_count=0` |

| Phase 13 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|---------------|--------------|--------------------------|
| 5 ファイル責務純度と旧 root path drift | 方針責務 / 文書・実装ファイル所在の索引責務 / 全 Go owner | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 13 target path 所在](DOCUMENT_INDEX.md#phase-13-target-paths)、各 owner の [`docs/details/*.md` 詳細本文責務](details/) |
| `<owner>.go`、`model.go`、`validate.go`、`execute.go`、`<owner>_test.go` の責務純度 | 全 Go owner | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) |
| builder production entrypoint / generated site output parity | `builder` / fixture 証跡責務 | [`docs/details/builder.md` 詳細本文責務 Phase 13 builder 実装整合契約](details/builder.md#phase-13-builder-alignment-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) |
| API credential 初期化と setup stdout 契約統一 | `api` / `setup` / `security` | [`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| API route / method / query / request / response、Admin / SDK / UI client binding、MCP bridge、setup stdout の契約自動照合 | `admin` / `api` / `sdk` / `ui` / `mcp` / `setup` / `security` | [`docs/details/admin.md` 詳細本文責務 Phase 13 admin 実装整合契約](details/admin.md#phase-13-admin-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/sdk.md` 詳細本文責務 Phase 13 SDK 実装整合契約](details/sdk.md#phase-13-sdk-alignment-contract)、[`docs/details/ui.md` 詳細本文責務 Phase 13 UI 実装整合契約](details/ui.md#phase-13-ui-alignment-contract)、[`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| API、runner、MCP、archive の状態 schema 共通化と直接状態更新禁止 | `statefile` / `api` / `runner` / `mcp` / `archive` | [`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/archive.md` 詳細本文責務 Phase 13 archive 実装整合契約](details/archive.md#phase-13-archive-alignment-contract) |
| lock、atomic write、JSON Lines、migration、復旧、fsync、rename、symlink 非追従、stale lock | `statefile` | [`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract) |
| password hash、Go 標準ライブラリ内製 KDF、token、署名、secret mask、file safety、token file / stdin | `security` / `admin` / `mcp` | [`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract)、[`docs/details/admin.md` 詳細本文責務 Phase 13 admin 実装整合契約](details/admin.md#phase-13-admin-alignment-contract)、[`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract) |
| snapshot、圧縮、digest、検証、download、restore、backup / restore transaction | `archive` / `api` / `runner` | [`docs/details/archive.md` 詳細本文責務 Phase 13 archive 実装整合契約](details/archive.md#phase-13-archive-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract) |
| GitHub Commit Status、retry、rate limit、timeout | `commitstatus` / `runner` | [`docs/details/commitstatus.md` 詳細本文責務 Phase 13 commitstatus 実装整合契約](details/commitstatus.md#phase-13-commitstatus-alignment-contract)、[`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract) |
| queue 状態機械、waiting to active、ID 採番、queue 上限、finalizer、active 復旧、at-least-once | `runner` / `statefile` / `api` | [`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract) |
| MCP resendWebhook、subscribe、unsubscribe、sampling 実動作化、未実装 error、HTTP lifecycle | `mcp` / `api` / `statefile` / `security` | [`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| Webhook SSRF / redirect / private IP / DNS rebinding、防御的 HTTP timeout / body 上限 / graceful shutdown | `api` / `mcp` / `runner` / `security` | [`docs/details/api.md` 詳細本文責務 Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/mcp.md` 詳細本文責務 Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| systemd 専用 user、最小権限、書込み先限定、SSH host key / timeout / known_hosts / remote path | `setup` / `runner` / `security` | [`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、[`docs/details/runner.md` 詳細本文責務 Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| fixture、mutation、race、fault injection、integration、E2E、CI required checks | fixture 証跡責務 / 全 owner | [`docs/details/fixture.md` fixture 証跡責務 Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) |
| Git tag、GitHub Release、SHA256SUMS、署名、SBOM、再現ビルド証跡、復旧手順 | `release` / `setup` / `security` / `archive` | [`docs/details/release.md` 詳細本文責務 Phase 13 release 実装整合契約](details/release.md#phase-13-release-alignment-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract)、[`docs/details/archive.md` 詳細本文責務 Phase 13 archive 実装整合契約](details/archive.md#phase-13-archive-alignment-contract) |

Phase 13 の実装作業は下表の work unit に分けて進める。work unit は実装順序と evidence package の記録単位であり、Phase 13 Pull Request を分割する根拠ではない。すべての work unit が `closed` にならない限り Phase 13 全体を完了扱いにしてはならない。

| work unit | 対象 | 実装詳細の正本 | 必須 evidence |
|-----------|------|----------------|---------------|
| `phase13-owner-shape` | 5 file 固定、旧 root path、空 file、dummy / no-op success | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` Phase 13 target path 所在](DOCUMENT_INDEX.md#phase-13-target-paths) | [`docs/details/fixture.md` Phase 13 証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence) の `input/owner_inventory.json` と `records/closure.jsonl` |
| `phase13-state-safety` | statefile 経由、process lock、atomic write、JSON Lines、migration、recovery | [`docs/details/statefile.md` Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract) | `input/state_inventory.json`、`records/fault.jsonl`、`records/race.jsonl` |
| `phase13-security-boundary` | Go 標準ライブラリ内製 KDF、token file / stdin、secret mask、required log | [`docs/details/security.md` Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract)、[`docs/details/admin.md` Phase 13 admin 実装整合契約](details/admin.md#phase-13-admin-alignment-contract)、[`docs/details/mcp.md` Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract) | `input/security_inventory.json`、`records/closure.jsonl` |
| `phase13-contract-parity` | API、Admin、SDK、UI、MCP、setup stdout、credential 初期化 | [`docs/details/api.md` Phase 13 API 実装整合契約](details/api.md#phase-13-api-alignment-contract)、[`docs/details/admin.md` Phase 13 admin 実装整合契約](details/admin.md#phase-13-admin-alignment-contract)、[`docs/details/sdk.md` Phase 13 SDK 実装整合契約](details/sdk.md#phase-13-sdk-alignment-contract)、[`docs/details/ui.md` Phase 13 UI 実装整合契約](details/ui.md#phase-13-ui-alignment-contract)、[`docs/details/mcp.md` Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/setup.md` Phase 13 setup 実装整合契約](details/setup.md#phase-13-setup-alignment-contract) | `input/contract_inventory.json`、`records/e2e.jsonl` |
| `phase13-runner-recovery` | queue 状態機械、finalizer、active recovery、at-least-once、backup / restore transaction | [`docs/details/runner.md` Phase 13 runner 実装整合契約](details/runner.md#phase-13-runner-alignment-contract)、[`docs/details/archive.md` Phase 13 archive 実装整合契約](details/archive.md#phase-13-archive-alignment-contract) | `records/race.jsonl`、`records/fault.jsonl`、`records/e2e.jsonl` |
| `phase13-mcp-release` | MCP 実動作化、commitstatus、release governance、GitHub Actions pinning | [`docs/details/mcp.md` Phase 13 MCP 実装整合契約](details/mcp.md#phase-13-mcp-alignment-contract)、[`docs/details/commitstatus.md` Phase 13 commitstatus 実装整合契約](details/commitstatus.md#phase-13-commitstatus-alignment-contract)、[`docs/details/release.md` Phase 13 release 実装整合契約](details/release.md#phase-13-release-alignment-contract) | `records/e2e.jsonl`、`records/release.jsonl` |
| `phase13-final-closure` | mutation、race、fault、CI required checks、document drift、ROADMAP / DOCUMENT_INDEX 更新 | [`docs/details/fixture.md` Phase 13 実装整合・品質改善証跡](details/fixture.md#phase-13-implementation-alignment-quality-evidence)、[`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan) | `expected/counters.json`、全 `records/*.jsonl`、Pull Request 本文の closure summary |

<a id="phase-14-obsidian-vault-integration-entry"></a>
**Phase 14 Obsidian Vault 連携参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 14 の対象から `obsidian` owner 詳細本文、builder collaborator 境界、fixture 証跡、状態・計画責務への入口だけを固定する。Phase 14 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) を参照する。

Phase 14 の Go owner package target は `obsidian` だけとする。Phase 14 実装で `obsidian` 以外の新規 owner package を追加してはならない。`builder` と `security` は collaborator として参照し、Phase 14 の vault 読取、wikilink / embed / tag / asset 正規化、builder handoff 判断を所有してはならない。

| 順序 | 固定する gate | 完了時の必須到達先 |
|------|---------------|--------------------|
| 1 | local vault 入力境界、vault root containment、symlink / hardlink / device 拒否、`.obsidian/` 非解釈 | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian Vault 連携契約](details/obsidian.md#obsidian-phase14-vault-integration-contract)、`phase14_obsidian_vault_boundary_open_count=0` |
| 2 | Markdown / wikilink / embed / tag / asset の deterministic 正規化、未解決 link の失敗、YAML frontmatter 拒否、`obsidian_map.json` schema | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian 正規化契約](details/obsidian.md#obsidian-phase14-normalization-contract)、[`docs/details/obsidian.md` 詳細本文責務 Phase 14 output schema 契約](details/obsidian.md#obsidian-phase14-output-schema-contract)、`phase14_obsidian_wikilink_open_count=0`、`phase14_obsidian_yaml_rejection_open_count=0`、`phase14_obsidian_asset_open_count=0` |
| 3 | builder handoff、stdout / stderr / report、公開出力への副作用境界、canonical JSON schema 比較 | [`docs/details/builder.md` 詳細本文責務](details/builder.md)、[`docs/details/obsidian.md` 詳細本文責務 Phase 14 builder handoff 契約](details/obsidian.md#obsidian-phase14-builder-handoff-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence)、`phase14_obsidian_builder_handoff_open_count=0` |
| 4 | fixture、negative control、closure record、DOCUMENT_INDEX 未作成 path | [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence)、[`docs/DOCUMENT_INDEX.md` Phase 14 target path 所在](DOCUMENT_INDEX.md#phase-14-target-paths)、`phase14_fixture_execution_gap_count=0`、`final_open_item_count=0` |

| Phase 14 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|---------------|--------------|--------------------------|
| Obsidian local vault 読取と path safety | `obsidian` / `security` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian Vault 連携契約](details/obsidian.md#obsidian-phase14-vault-integration-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| wikilink / embed / tag / asset 正規化 | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian 正規化契約](details/obsidian.md#obsidian-phase14-normalization-contract) |
| builder handoff | `obsidian` / `builder` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 builder handoff 契約](details/obsidian.md#obsidian-phase14-builder-handoff-contract)、[`docs/details/builder.md` 詳細本文責務](details/builder.md) |
| Phase 14 fixture / closure | fixture 証跡責務 | [`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) |

<a id="phase-15-obsidian-local-sync-entry"></a>
**Phase 15 Obsidian local vault 同期参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 15 の対象から `obsidian` owner 詳細本文、statefile / security / release / setup collaborator 境界、fixture 証跡、状態・計画責務への入口だけを固定する。Phase 15 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、完了判定方針は [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) を参照する。

Phase 15 の Go owner package target は `obsidian` だけとする。Phase 15 実装で `obsidian` 以外の新規 owner package を追加してはならない。`statefile`、`security`、`release`、`setup` は collaborator として参照し、Phase 15 の sync plan、sync apply、sync rollback、conflict、tombstone 判断を所有してはならない。

| 順序 | 固定する gate | 完了時の必須到達先 |
|------|---------------|--------------------|
| 1 | `sync plan` の input snapshot、digest、change set、conflict 事前検出、`--plan-file` 以外への書込み禁止、plan / state schema | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 Obsidian local vault 同期契約](details/obsidian.md#obsidian-phase15-local-sync-contract)、[`docs/details/obsidian.md` 詳細本文責務 Phase 15 schema 契約](details/obsidian.md#obsidian-phase15-schema-contract)、`phase15_sync_plan_mismatch_count=0` |
| 2 | `sync apply` の plan hash 照合、process lock、staging、atomic rename、fsync、rollback record schema | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 apply / rollback 契約](details/obsidian.md#obsidian-phase15-apply-rollback-contract)、[`docs/details/obsidian.md` 詳細本文責務 Phase 15 schema 契約](details/obsidian.md#obsidian-phase15-schema-contract)、`phase15_sync_apply_atomicity_open_count=0`、`phase15_sync_rollback_open_count=0` |
| 3 | conflict、tombstone、delete policy、clock skew、read-only file、partial write / kill 復旧 | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 conflict / tombstone 契約](details/obsidian.md#obsidian-phase15-conflict-tombstone-contract)、`phase15_sync_conflict_open_count=0`、`phase15_sync_tombstone_open_count=0` |
| 4 | Obsidian Sync service / cloud / plugin 非依存、URI open の任意境界、Release / setup 配布連携、fixture closure | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 外部境界契約](details/obsidian.md#obsidian-phase15-external-boundary-contract)、[`docs/details/obsidian.md` 詳細本文責務 Phase 15 配布連携契約](details/obsidian.md#obsidian-phase15-distribution-contract)、[`docs/details/release.md` 詳細本文責務 Phase 15 Obsidian Release 配布拡張契約](details/release.md#phase-15-obsidian-release-extension-contract)、[`docs/details/setup.md` 詳細本文責務 Phase 15 Obsidian CLI 導入手順](details/setup.md#phase-15-obsidian-setup-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence)、`phase15_sync_service_dependency_open_count=0`、`phase15_distribution_open_count=0`、`phase15_fixture_execution_gap_count=0`、`final_open_item_count=0` |

| Phase 15 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|---------------|--------------|--------------------------|
| Obsidian local vault 同期 plan / apply / rollback | `obsidian` / `statefile` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 Obsidian local vault 同期契約](details/obsidian.md#obsidian-phase15-local-sync-contract)、[`docs/details/statefile.md` 詳細本文責務 Phase 13 statefile 実装整合契約](details/statefile.md#phase-13-statefile-alignment-contract) |
| conflict / tombstone / delete policy | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 conflict / tombstone 契約](details/obsidian.md#obsidian-phase15-conflict-tombstone-contract) |
| local filesystem safety / secret 非保持 | `obsidian` / `security` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 外部境界契約](details/obsidian.md#obsidian-phase15-external-boundary-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
| Obsidian CLI 配布 | `obsidian` / `release` / `setup` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 配布連携契約](details/obsidian.md#obsidian-phase15-distribution-contract)、[`docs/details/release.md` 詳細本文責務 Phase 15 Obsidian Release 配布拡張契約](details/release.md#phase-15-obsidian-release-extension-contract)、[`docs/details/setup.md` 詳細本文責務 §26.2a](details/setup.md#sec-26-2a)、[`docs/details/setup.md` 詳細本文責務 Phase 15 Obsidian CLI 導入手順](details/setup.md#phase-15-obsidian-setup-contract) |
| Phase 15 fixture / closure | fixture 証跡責務 | [`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) |

<a id="phase-16-quality-evidence-closure-entry"></a>
**Phase 16 実装済み品質証跡実体化・追加検証候補 closure 参照：**

[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、Phase 16 の対象から該当する owner 詳細本文、fixture 証跡、状態・計画責務、checker implementation artifact、実装順序固定契約への入口だけを固定する。Phase 16 の現在状態、依存 Phase、active Phase は [`docs/ROADMAP.md` 状態・計画責務 §4.1](ROADMAP.md#roadmap-initial-phase-plan)、Phase 単位の完了条件は [`docs/SPEC.md` ポリシー責務 §0f](SPEC.md#policy-phase-unit)、テスト方針と mutation 必須条件は [`docs/SPEC.md` ポリシー責務 §0g](SPEC.md#policy-meaningful-test)、fixture 証跡は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence)、Phase 16 実装順序は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 implementation sequence 固定契約](details/fixture.md#phase-16-implementation-sequence-contract)、checker 実装 artifact は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](DOCUMENT_INDEX.md#phase-16-target-paths) の [`main_test.go`](../main_test.go) を参照する。

Phase 16 は以下の内部順序で仕様確認、実装、検証、証跡記録、状態更新を進める。前の順序に open item が残る場合、後続順序を完了扱いにしてはならない。実装 artifact 作成時の細分化順序、途中完了禁止、checker skeleton 先行条件は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 implementation sequence 固定契約](details/fixture.md#phase-16-implementation-sequence-contract) を参照する。

| 順序 | 固定する gate | 完了時の必須到達先 |
|------|---------------|--------------------|
| 1 | Phase 12 から Phase 15 の宣言型証跡、JSONL counter、record 存在確認、skip、未実行、対象外理由を棚卸しし、実行型証跡へ置換する対象を確定する | [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence)、`phase16_declarative_evidence_open_count=0` |
| 2 | [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence) の `ALIGN-*` 追加検証候補を、Phase 16 実施、将来計画維持、または対象外 anchor へ分類し、未分類を 0 にする | `phase16_align_unclassified_count=0`、`phase16_align_open_count=0` |
| 3 | ignored error、required write failure、audit / access / config / state / JSON Lines write、panic、skip / 未実行を全 owner で分類し、成功証跡として扱えない経路を排除する | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](#cross-test-evidence-route)、`phase16_ignored_error_unclassified_count=0`、`phase16_required_write_failure_open_count=0`、`phase16_skip_open_count=0`、`phase16_panic_unclassified_count=0` |
| 4 | clock、sleep、timeout、entropy、parallel worker、HTTP lifecycle、partial write、client disconnect、external I/O を deterministic fake または実行証跡へ接続する | [`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract)、`phase16_determinism_open_count=0`、`phase16_http_boundary_open_count=0` |
| 5 | state、config、audit、history、snapshot、archive、release、setup、Obsidian sync の atomic write、fsync、parent directory fsync、rename、symlink 非追従、recovery を owner 間で同じ証跡水準へ揃える | [`docs/details/statefile.md` 詳細本文責務](details/statefile.md)、[`docs/details/archive.md` 詳細本文責務](details/archive.md)、`phase16_filesystem_durability_open_count=0` |
| 6 | `.github/workflows/phase12-quality-gate.yml` を含む required workflow hardening、Docker 検証手順、Deno stable runtime による JavaScript 検証、release rehearsal、巨大 owner risk ledger、5 ファイル原則維持を closure へ接続する | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](DOCUMENT_INDEX.md#phase-16-target-paths)、`phase16_workflow_hardening_open_count=0`、`phase16_validation_portability_open_count=0`、`phase16_large_owner_risk_open_count=0` |
| 7 | closure record set、正式 fixture root、`input/negative_controls.json`、checker implementation artifact、checker 再導出、checker 実行入口、checker 入出力、required check、ROADMAP 状態、DOCUMENT_INDEX 所在、PR 証跡を一致させる | [`docs/details/fixture.md` fixture 証跡責務 implementation PR evidence template 固定契約](details/fixture.md#implementation-pr-evidence-template-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 checker 固定契約](details/fixture.md#phase-16-checker-contract)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](DOCUMENT_INDEX.md#phase-16-target-paths)、`final_open_item_count=0` |

| Phase 16 対象 | owner / 責務 | 詳細本文 / fixture 証跡 |
|---------------|--------------|--------------------------|
| 宣言型証跡から実行型証跡への置換 | fixture 証跡責務 / 全 owner | [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence)、[`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract) |
| `ALIGN-*` 追加検証候補分類 | fixture 証跡責務 / 状態・計画責務 | [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence)、[`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) |
| ignored error / required write failure / panic / skip | 全 Go owner / fixture 証跡責務 | 各 owner の [`docs/details/*.md` 詳細本文責務](details/)、[`docs/details/fixture.md` fixture 証跡責務 skip / 未実行証跡固定契約](details/fixture.md#test-skip-evidence-contract) |
| determinism / HTTP boundary / external I/O | `api` / `runner` / `mcp` / `setup` / `release` / `obsidian` | [`docs/details/api.md` 詳細本文責務](details/api.md)、[`docs/details/runner.md` 詳細本文責務](details/runner.md)、[`docs/details/mcp.md` 詳細本文責務](details/mcp.md)、[`docs/details/setup.md` 詳細本文責務](details/setup.md)、[`docs/details/release.md` 詳細本文責務](details/release.md)、[`docs/details/obsidian.md` 詳細本文責務](details/obsidian.md) |
| filesystem durability parity | `statefile` / `archive` / `runner` / `api` / `setup` / `release` / `obsidian` | [`docs/details/statefile.md` 詳細本文責務](details/statefile.md)、[`docs/details/archive.md` 詳細本文責務](details/archive.md)、各 collaborator owner 詳細本文 |
| workflow hardening / Docker 検証 / release rehearsal | fixture 証跡責務 / `release` / `setup` | [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence)、[`docs/details/release.md` 詳細本文責務](details/release.md)、[`docs/details/setup.md` 詳細本文責務](details/setup.md) |
| 巨大 owner risk ledger / 5 ファイル原則維持 | 方針責務 / 文書・実装ファイル所在の索引責務 / 全 Go owner | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](DOCUMENT_INDEX.md#phase-16-target-paths)、各 owner の [`docs/details/*.md` 詳細本文責務](details/) |

Phase 16 の実装作業は下表の work unit に分けて進める。work unit は実装順序と evidence package の記録単位であり、Phase 16 Pull Request を分割する根拠ではない。すべての work unit が `closed` にならない限り Phase 16 全体を完了扱いにしてはならない。

| work unit | 対象 | 正本 / 証跡 |
|-----------|------|-------------|
| `phase16-evidence-execution` | Phase 12 から Phase 15 の宣言型証跡、mutation、fault、race、E2E、release evidence の実行化 | [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence) |
| `phase16-align-closure` | `ALIGN-*` 追加検証候補の分類、実施、対象外 anchor、将来計画維持 | [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence)、[`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) |
| `phase16-error-boundary` | ignored error、required write、panic、skip / 未実行 | [横断テスト証跡共通入口](#cross-test-evidence-route)、各 owner 詳細本文 |
| `phase16-determinism-boundary` | clock、sleep、timeout、entropy、HTTP lifecycle、external I/O | [`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract) |
| `phase16-filesystem-durability` | atomic write、fsync、parent directory fsync、rename、symlink 非追従、recovery | [`docs/details/statefile.md` 詳細本文責務](details/statefile.md)、[`docs/details/archive.md` 詳細本文責務](details/archive.md) |
| `phase16-ci-release-portability` | Phase 12 workflow hardening、Docker 検証、Deno stable runtime による JavaScript 検証、release rehearsal | [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence) |
| `phase16-maintainability` | 巨大 owner risk ledger、5 ファイル原則維持、内部責務区画検査 | [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3)、[`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 16 target path 所在](DOCUMENT_INDEX.md#phase-16-target-paths) |

Phase 16 の各 work unit は、下表の inventory と owner 範囲をすべて処理対象に含める。対象 source が存在しない owner は、存在しないことを理由に省略せず、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 inventory 共通 schema](details/fixture.md#phase-16-inventory-common-schema) の `classification=not_applicable` とし、責務正本 anchor を記録する。棚卸し対象 source と検出語は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 実装済み品質証跡実体化・追加検証候補 closure 証跡](details/fixture.md#phase-16-quality-evidence-closure-evidence) の source coverage set を正本とし、下表の対象 owner / artifact から存在する source を任意に除外してはならない。

| work unit | 必須 inventory | 対象 owner / artifact | 必須処理 |
|-----------|----------------|------------------------|----------|
| `phase16-evidence-execution` | `input/evidence_inventory.json` | Phase 12、Phase 13、Phase 14、Phase 15 の fixture root、required check、PR evidence、closure counter、mutation、fault、race、E2E、release 証跡 | 宣言型証跡、存在確認だけの証跡、未実行証跡、counter だけの完了表記を `execute_and_close`、`reject_as_invalid_evidence`、`not_applicable` のいずれかへ分類する。`expected/counters.json` の全 counter は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 expected/counters schema](details/fixture.md#phase-16-quality-evidence-closure-evidence) と一致させる。 |
| `phase16-align-closure` | `input/align_inventory.json` | [`docs/details/fixture.md` fixture 証跡責務 現行実装整合証跡](details/fixture.md#current-implementation-alignment-evidence) の全 `ALIGN-*` | 各 `ALIGN-*` を Phase 16 実施、将来計画維持、対象外に 1 つだけ分類し、未分類を残さない。 |
| `phase16-error-boundary` | `input/error_inventory.json` | 全 Go owner package、[`main.go`](../main.go)、[`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js)、[`admin/index.html`](../admin/index.html) | `_ =`、ignored return、`panic(`、`t.Skip` / skip、cleanup failure、required write failure、SSE / HTTP write failure を分類し、成功証跡として扱えない経路を closure へ接続する。 |
| `phase16-determinism-boundary` | `input/determinism_inventory.json` | `api`、`runner`、`mcp`、`setup`、`release`、`obsidian`、`builder`、`admin`、`sdk`、`ui` | `time.Now`、sleep / timer、timeout、entropy、HTTP client、external command、systemd、filesystem order、parallel worker を fake adapter または実行証跡へ接続する。 |
| `phase16-filesystem-durability` | `input/filesystem_inventory.json` | `statefile`、`archive`、`runner`、`api`、`setup`、`release`、`obsidian`、`builder` | write、append、flush、fsync、parent directory fsync、rename、remove、walk、symlink 非追従、stale lock、recovery、permission failure を owner 間で同じ証跡水準へ揃える。 |
| `phase16-ci-release-portability` | `input/workflow_inventory.json` | `.github/workflows/*.yml`、`.github/workflows/*.yaml`、Docker fallback、Release / setup rehearsal、Go / Deno runtime 検証 | action SHA pin、minimum permissions、timeout、required check、host tool 不在時の Docker 検証、release rehearsal を closure へ接続する。 |
| `phase16-maintainability` | `input/large_owner_inventory.json` | 全 Go owner package、[`main.go`](../main.go)、root contract test | 5 ファイル固定違反、6 ファイル目、subdirectory、空 file、dummy 実装、巨大 file risk、内部責務区画、test artifact 接続を分類し、`phase16_large_owner_risk_open_count=0` へ接続する。 |

`input/negative_controls.json` は 8 番目の work unit ではなく、7 work unit 全体に対する checker self-verification inventory とする。各 negative control record は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 inventory 共通 schema](details/fixture.md#phase-16-inventory-common-schema) の `work_unit` で上表のいずれか 1 つへ所属し、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 required check](details/fixture.md#phase-16-quality-evidence-closure-evidence) の check 名のいずれか 1 つへ接続し、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 negative control coverage matrix 固定契約](details/fixture.md#phase-16-negative-control-coverage-matrix) の `target_contract` をすべて網羅する。7 work unit のいずれにも negative control がない場合、required check 表のいずれにも negative control がない場合、または coverage matrix の未充足 `target_contract` が 1 件でもある場合、Phase 16 の仕様策定完了および実装完了を認めない。

Phase 16 の実装者は、各 inventory の `required_action` が `future_plan` の場合でも、[`docs/ROADMAP.md` 状態・計画責務 §5](ROADMAP.md#522-統合ロードマップ表) への到達可能な `future_ref` を記録する。`future_ref` を持たない将来計画維持、`closure_ref` を持たない完了、`not_applicable_ref` を持たない対象外、または同一 item の複数 action は未完了として扱う。

Phase 16 は後続 Phase 文書を Phase 16 の未完了項目として吸収してはならない。Phase 17 以降の詳細本文、target path、fixture root、required check、provider validation、production simulation は、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 future Phase document boundary 固定契約](details/fixture.md#phase-16-future-phase-document-boundary) に従い、document drift と将来計画維持の境界確認だけに使用する。

Phase 16 の完了判定では、[`docs/details/fixture.md` fixture 証跡責務 Phase 16 checker 固定契約](details/fixture.md#phase-16-checker-contract) を必ず通過する。checker が live source coverage set から検出した対象を `input/source_coverage.json`、各 inventory、`expected/actions.json`、`records/*.jsonl`、closure counter に再導出できない場合、該当 work unit は `closed` にしてはならない。手書き counter、手書き closure record、または PR 本文だけを checker の代替証跡にしてはならない。checker の stdout / stderr、exit code、diagnostic schema、source coverage detection registry は [`docs/details/fixture.md` fixture 証跡責務 Phase 16 checker 固定契約](details/fixture.md#phase-16-checker-contract) のみを正本とする。

<a id="phase-17-production-validation-entry"></a>
**Phase 17 ConoHa VPS 試験本番運用参照：**

Phase 17 の現在状態と依存順序は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、ConoHa VPS 試験本番運用の本文と証跡接続条件は [`docs/details/production-validation.md`](details/production-validation.md) 本番検証詳細本文責務、fixture と required check は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](details/fixture.md#phase-17-production-validation-evidence)、実在 path は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 17 target path 所在](DOCUMENT_INDEX.md#phase-17-target-paths) を参照する。

Phase 17 は owner component を追加しない。ConoHa VPS 試験本番運用、provider target、OS / plan 前提、opaque environment identity、検証 mode、試験本番運用窓・停止条件、本番環境同等 simulation、destructive operation 禁止境界、credential / SSH 入力境界、運用中バグ修正順序、エックスサーバ VPS 将来判断の本文は [`docs/details/production-validation.md`](details/production-validation.md) 本番検証詳細本文責務だけに置く。fixture path、input / expected / record schema、closure counter、required check、required check workflow、checker 実行入口、negative boundary、document drift record の本文は [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](details/fixture.md#phase-17-production-validation-evidence) だけに置く。API、Admin、SDK、UI、runner、statefile、setup、release、security、archive、commitstatus、mcp の処理仕様は各 owner component 別詳細本文を参照し、Phase 17 本文では再定義しない。

| 判断対象 | 正本参照 | DETAIL_INDEX での扱い |
|----------|----------|-----------------------|
| provider target / OS / plan / opaque environment identity / 検証 mode / 構築プラン / 試験本番運用窓・停止条件 / ConoHa VPS 試験本番運用 / credential / SSH 入力境界 / エックスサーバ VPS 将来判断 | [`docs/details/production-validation.md`](details/production-validation.md) 本番検証詳細本文責務 | 固定値、順序、禁止境界を再掲しない。 |
| fixture path / input schema / expected schema / record schema / closure counter / required check / workflow / checker / negative boundary / document drift record | [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](details/fixture.md#phase-17-production-validation-evidence) | schema、check 名、完了条件を再掲しない。 |
| 実在 path / 未作成 path | [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 17 target path 所在](DOCUMENT_INDEX.md#phase-17-target-paths) | 所在確認だけを行う。 |

Phase 17 実装 PR は、required check 名、完了条件、required check workflow 契約、checker 実行入口、document drift record を [`docs/details/fixture.md` fixture 証跡責務 Phase 17 ConoHa VPS 試験本番運用証跡](details/fixture.md#phase-17-production-validation-evidence) から参照し、同一 PR 本文へ記録する。

<a id="phase-18-trial-production-operation-entry"></a>
**Phase 18 試験本番VPS 実運用接続・運用証跡参照：**

Phase 18 の現在状態と依存順序は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、試験本番VPS 実運用接続・運用証跡の本文は [`docs/details/production-validation.md` 本番検証詳細本文責務 Phase 18 試験本番VPS 実運用接続・運用証跡契約](details/production-validation.md#phase-18-trial-production-operation-contract)、fixture と required check は [`docs/details/fixture.md` fixture 証跡責務 Phase 18 試験本番VPS 実運用接続・運用証跡](details/fixture.md#phase-18-trial-production-operation-evidence)、実在 path は [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 18 target path 所在](DOCUMENT_INDEX.md#phase-18-target-paths) を参照する。

Phase 18 は owner component を追加しない。試験本番VPS の実運用接続、Phase 17 からの引継ぎ、運用 session、service health、runtime flow、update / rollback drill、運用中 issue、仕様先行バグ修正、再検証、secret / destructive operation 境界、document drift の本文は [`docs/details/production-validation.md`](details/production-validation.md) 本番検証詳細本文責務だけに置く。fixture path、input / expected / record schema、closure counter、required check、required check workflow、checker 実行入口、negative boundary、document drift record の本文は [`docs/details/fixture.md` fixture 証跡責務 Phase 18 試験本番VPS 実運用接続・運用証跡](details/fixture.md#phase-18-trial-production-operation-evidence) だけに置く。API、Admin、SDK、UI、runner、statefile、setup、release、security、archive、commitstatus、mcp の処理仕様は各 owner component 別詳細本文を参照し、Phase 18 本文では再定義しない。

| 判断対象 | 正本参照 | DETAIL_INDEX での扱い |
|----------|----------|-----------------------|
| 実運用接続 / Phase 17 引継ぎ / 運用 session / service health / runtime flow / update / rollback drill / issue / bugfix loop / secret 境界 / destructive operation 境界 / document drift 境界 | [`docs/details/production-validation.md` 本番検証詳細本文責務 Phase 18 試験本番VPS 実運用接続・運用証跡契約](details/production-validation.md#phase-18-trial-production-operation-contract) | 固定値、順序、禁止境界を再掲しない。 |
| fixture path / input schema / expected schema / record schema / closure counter / required check / workflow / checker / negative boundary / document drift record | [`docs/details/fixture.md` fixture 証跡責務 Phase 18 試験本番VPS 実運用接続・運用証跡](details/fixture.md#phase-18-trial-production-operation-evidence) | schema、check 名、完了条件を再掲しない。 |
| 実在 path / 未作成 path | [`docs/DOCUMENT_INDEX.md` 文書・実装ファイル所在の索引責務 Phase 18 target path 所在](DOCUMENT_INDEX.md#phase-18-target-paths) | 所在確認だけを行う。 |

Phase 18 実装 PR は、required check 名、完了条件、required check workflow 契約、checker 実行入口、document drift record を [`docs/details/fixture.md` fixture 証跡責務 Phase 18 試験本番VPS 実運用接続・運用証跡](details/fixture.md#phase-18-trial-production-operation-evidence) から参照し、同一 PR 本文へ記録する。

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

<a id="0i8-obsidian"></a>
**0i.8 Obsidian：**

| 機能 | owner | 詳細本文 |
|------|-------|----------|
| Obsidian Vault 連携 | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian Vault 連携契約](details/obsidian.md#obsidian-phase14-vault-integration-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) |
| Obsidian wikilink / embed / tag / asset 正規化 | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 Obsidian 正規化契約](details/obsidian.md#obsidian-phase14-normalization-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 14 Obsidian Vault 連携証跡](details/fixture.md#phase-14-obsidian-vault-integration-evidence) |
| Obsidian builder handoff | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 14 builder handoff 契約](details/obsidian.md#obsidian-phase14-builder-handoff-contract)、[`docs/details/builder.md` 詳細本文責務](details/builder.md) |
| Obsidian local vault 同期 plan / apply / rollback | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 Obsidian local vault 同期契約](details/obsidian.md#obsidian-phase15-local-sync-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) |
| Obsidian conflict / tombstone / delete policy | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 conflict / tombstone 契約](details/obsidian.md#obsidian-phase15-conflict-tombstone-contract)、[`docs/details/fixture.md` fixture 証跡責務 Phase 15 Obsidian local vault 同期証跡](details/fixture.md#phase-15-obsidian-local-sync-evidence) |
| Obsidian 外部境界 | `obsidian` | [`docs/details/obsidian.md` 詳細本文責務 Phase 15 外部境界契約](details/obsidian.md#obsidian-phase15-external-boundary-contract)、[`docs/details/security.md` 詳細本文責務 Phase 13 security 実装整合契約](details/security.md#phase-13-security-alignment-contract) |
