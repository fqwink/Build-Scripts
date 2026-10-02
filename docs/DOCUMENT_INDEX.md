# Build-Scripts — 文書・実装ファイル所在索引

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は、Build-Scripts リポジトリ内の文書、実装ファイル、テスト、fixture、未作成 path の所在だけを示す索引である。方針、状態定義、現在状態、実装詳細、デザイン、検証契約は本文として定義しない。

<a id="索引責務の使い方"></a>
**所在確認対象と参照先：**

| 確認対象 | 参照先 |
|----------|--------|
| 作業ルール | [`AGENTS.md`](../AGENTS.md) |
| 方針、ポリシー、状態語彙、状態遷移条件 | [`docs/SPEC.md`](SPEC.md) |
| 実装 artifact と各機能の現在状態、Phase、将来計画 | [`docs/ROADMAP.md`](ROADMAP.md) |
| 詳細仕様入口、共通固定値、owner 対応表、collaborator 境界参照入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) |
| 詳細仕様本文・証跡ディレクトリ | [`docs/details/`](details/) |
| Phase 11 バグ修正ゼロ化の参照入口 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 11 バグ修正ゼロ化参照](DETAIL_INDEX.md#phase-11-quality-gate-entry) |
| Phase 12 実装品質ゲート再構築の参照入口 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 実装品質ゲート再構築参照](DETAIL_INDEX.md#phase-12-quality-gate-entry) |
| Phase 12 正式 fixture root | [`testdata/phase12/quality-gate-reconstruction/`](../testdata/phase12/quality-gate-reconstruction/) |
| Phase 12 required check workflow | [`.github/workflows/phase12-quality-gate.yml`](../.github/workflows/phase12-quality-gate.yml) |
| Phase 13 実装整合・品質改善の参照入口 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) |
| Phase 13 正式 fixture root | `testdata/phase13/implementation-alignment-quality/` |
| Phase 13 required check workflow | `.github/workflows/phase13-implementation-alignment-quality.yml` |
| 横断テスト証跡の共通入口 | [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route) |
| Phase 11 仕様全般完了判定の所在 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) |
| fixture、expected、fake、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) |
| test artifact と fixture 証跡の接続 | [`docs/details/fixture.md` fixture 証跡責務 test artifact traceability 固定契約](details/fixture.md#test-artifact-traceability-contract) |
| 単独 test artifact を持たない owner の検証接続 | [`docs/details/fixture.md` fixture 証跡責務 non-dedicated owner test routing 固定契約](details/fixture.md#non-dedicated-owner-test-routing-contract) |
| Phase 11 fixture root 網羅判定 | [`docs/details/fixture.md` fixture 証跡責務 fixture root coverage matrix 固定契約](details/fixture.md#fixture-root-coverage-matrix-contract) |
| fixture group root と formal fixture root の境界 | [`docs/details/fixture.md` fixture 証跡責務 fixture root coverage matrix 固定契約](details/fixture.md#fixture-root-coverage-matrix-contract) |
| 未作成 fixture root closure record | [`docs/details/fixture.md` fixture 証跡責務 未作成 fixture root closure record 固定契約](details/fixture.md#fixture-root-missing-closure-record-contract) |
| 検証実行証跡の分類と未完了条件 | [`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract) |
| test evidence package 記録先 | [`docs/details/fixture.md` fixture 証跡責務 test evidence package 記録先固定契約](details/fixture.md#test-evidence-package-record-location-contract) |
| test gap inventory 記録 | [`docs/details/fixture.md` fixture 証跡責務 test gap inventory record 固定契約](details/fixture.md#test-gap-inventory-record-contract)。テスト固有の仕様不足と仕様全般不備の接続入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route)。 |
| test improvement batch closure | [`docs/details/fixture.md` fixture 証跡責務 test improvement batch closure 固定契約](details/fixture.md#test-improvement-batch-closure-contract)。 |
| 仕様全般不備 inventory record | [`docs/SPEC.md` ポリシー責務 仕様全般不備 inventory record 固定契約](SPEC.md#spec-deficiency-inventory-record-contract)。テスト固有の仕様不足と仕様全般不備の接続入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route)。 |
| 仕様全般不備 batch closure | [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract)。`test_gap_connection` の接続入口は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 横断テスト証跡共通入口](DETAIL_INDEX.md#cross-test-evidence-route)。 |
| 仕様全般不備 inspection pass 証跡 | [`docs/SPEC.md` ポリシー責務 仕様全般不備 batch closure 固定契約](SPEC.md#spec-deficiency-batch-closure-contract) の `inspection_scope`、`inspection_pass_count`、`inspection_pass_summary`。 |
| assertion id 形式 | [`docs/details/fixture.md` fixture 証跡責務 assertion id 固定契約](details/fixture.md#test-assertion-id-contract) |
| skip / 未実行証跡 | [`docs/details/fixture.md` fixture 証跡責務 skip / 未実行証跡固定契約](details/fixture.md#test-skip-evidence-contract) |
| 検証要求網羅 ledger | [`docs/details/fixture.md` fixture 証跡責務 test requirement coverage ledger 固定契約](details/fixture.md#test-requirement-coverage-ledger-contract) |
| 検証証跡クロージャと open item 0 判定 | [`docs/details/fixture.md` fixture 証跡責務 test verification closure checklist 固定契約](details/fixture.md#test-verification-closure-checklist-contract) |
| 検証証跡クロージャ記録 schema | [`docs/details/fixture.md` fixture 証跡責務 test verification closure record schema 固定契約](details/fixture.md#test-verification-closure-record-schema-contract) |
| 検証証跡クロージャ記録 set | [`docs/details/fixture.md` fixture 証跡責務 test verification closure record set 固定契約](details/fixture.md#test-verification-closure-record-set-contract) |
| 実装 PR 証跡 package | [`docs/details/fixture.md` fixture 証跡責務 implementation PR evidence template 固定契約](details/fixture.md#implementation-pr-evidence-template-contract) |
| test oracle 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test oracle evidence set 固定契約](details/fixture.md#test-oracle-evidence-set-contract) |
| test assertion identity / failure diagnostics 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test assertion identity / failure diagnostics evidence set 固定契約](details/fixture.md#test-assertion-failure-diagnostics-evidence-set-contract) |
| test boundary / failure matrix 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test boundary / failure matrix evidence set 固定契約](details/fixture.md#test-boundary-failure-matrix-evidence-set-contract) |
| test isolation 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test isolation evidence set 固定契約](details/fixture.md#test-isolation-evidence-set-contract) |
| test determinism 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test determinism evidence set 固定契約](details/fixture.md#test-determinism-evidence-set-contract) |
| test concurrency / race 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test concurrency / race evidence set 固定契約](details/fixture.md#test-concurrency-race-evidence-set-contract) |
| race trigger matrix | [`docs/details/fixture.md` fixture 証跡責務 race trigger matrix 固定契約](details/fixture.md#race-trigger-matrix-contract) |
| mutation test 証跡固定契約 | [`docs/details/fixture.md` fixture 証跡責務 mutation test 証跡固定契約](details/fixture.md#mutation-test-evidence-contract) |
| mutation test 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 mutation test evidence set 固定契約](details/fixture.md#mutation-test-evidence-set-contract) |
| mutation selection ledger | [`docs/details/fixture.md` fixture 証跡責務 mutation selection ledger 固定契約](details/fixture.md#mutation-selection-ledger-contract) |
| test harness self-verification 証跡 set | [`docs/details/fixture.md` fixture 証跡責務 test harness self-verification evidence set 固定契約](details/fixture.md#test-harness-self-verification-evidence-set-contract) |
| test / contract drift 証跡 | [`docs/details/fixture.md` fixture 証跡責務 test / contract drift 証跡固定契約](details/fixture.md#test-contract-drift-evidence-contract) |
| test / contract drift report schema | [`docs/details/fixture.md` fixture 証跡責務 test / contract drift report schema 固定契約](details/fixture.md#test-contract-drift-report-schema-contract) |
| manifest not_applicable 境界 | [`docs/details/fixture.md` fixture 証跡責務 manifest not_applicable 境界固定契約](details/fixture.md#fixture-manifest-not-applicable-boundary-contract) |
| 生成静的 Web サイトと標準管理 UI のデザイン | [`docs/DESIGN.md`](DESIGN.md) |
| 利用入口 | [`README.md`](../README.md) |

<a id="文書一覧"></a>
**文書所在：**

| パス | 所在区分 |
|------|----------|
| [`AGENTS.md`](../AGENTS.md) | 実在 |
| [`README.md`](../README.md) | 実在 |
| [`docs/SPEC.md`](SPEC.md) | 実在 |
| [`docs/ROADMAP.md`](ROADMAP.md) | 実在 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) | 実在 |
| [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) | 実在 |
| [`docs/DESIGN.md`](DESIGN.md) | 実在 |
| [`docs/details/`](details/) | 実在 |
| `docs/examples/` | 未作成 |
| `LICENSE` | 未作成 |
| `SECURITY.md` | 未作成 |
| `CONTRIBUTING.md` | 未作成 |
| `CODEOWNERS` | 未作成 |
| `CHANGELOG.md` | 未作成 |

<a id="詳細仕様本文の所在"></a>
**詳細仕様本文・証跡所在：**

| パス | owner component / 責務 |
|------|-------------------------|
| [`docs/details/builder.md`](details/builder.md) | `builder` |
| [`docs/details/runner.md`](details/runner.md) | `runner` |
| [`docs/details/api.md`](details/api.md) | `api` |
| [`docs/details/admin.md`](details/admin.md) | `admin` |
| [`docs/details/sdk.md`](details/sdk.md) | `sdk` |
| [`docs/details/ui.md`](details/ui.md) | `ui` |
| [`docs/details/setup.md`](details/setup.md) | `setup` |
| [`docs/details/release.md`](details/release.md) | `release` |
| [`docs/details/statefile.md`](details/statefile.md) | `statefile` |
| [`docs/details/archive.md`](details/archive.md) | `archive` |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | `commitstatus` |
| [`docs/details/security.md`](details/security.md) | `security` |
| [`docs/details/mcp.md`](details/mcp.md) | `mcp` |
| [`docs/details/fixture.md`](details/fixture.md) | fixture 証跡責務 |

<a id="実装ファイル一覧"></a>
**実装・テスト・fixture 所在：**

| パス | 実装上の役割 | 所在区分 |
|------|-----------|----------|
| [`main.go`](../main.go) | 起動入口 | 実在 |
| [`main_test.go`](../main_test.go) | 起動入口 test | 実在 |
| [`go.mod`](../go.mod) | Go module | 実在 |
| [`sdk_contract_test.go`](../sdk_contract_test.go) | `sdk` contract test | 実在 |
| [`ui_contract_test.go`](../ui_contract_test.go) | `ui` contract test | 実在 |
| [`components/admin/admin.go`](../components/admin/admin.go) | `admin` owner package 公開実行境界 | 実在 |
| [`components/admin/model.go`](../components/admin/model.go) | `admin` owner package model | 実在 |
| [`components/admin/validate.go`](../components/admin/validate.go) | `admin` owner package validation | 実在 |
| [`components/admin/execute.go`](../components/admin/execute.go) | `admin` owner package execution | 実在 |
| [`components/admin/admin_test.go`](../components/admin/admin_test.go) | `admin` owner package test | 実在 |
| [`components/api/api.go`](../components/api/api.go) | `api` owner package 公開実行境界 | 実在 |
| [`components/api/model.go`](../components/api/model.go) | `api` owner package model | 実在 |
| [`components/api/validate.go`](../components/api/validate.go) | `api` owner package validation | 実在 |
| [`components/api/execute.go`](../components/api/execute.go) | `api` owner package execution | 実在 |
| [`components/api/api_test.go`](../components/api/api_test.go) | `api` owner package test | 実在 |
| [`components/archive/archive.go`](../components/archive/archive.go) | `archive` owner package 公開実行境界 | 実在 |
| [`components/archive/model.go`](../components/archive/model.go) | `archive` owner package model | 実在 |
| [`components/archive/validate.go`](../components/archive/validate.go) | `archive` owner package validation | 実在 |
| [`components/archive/execute.go`](../components/archive/execute.go) | `archive` owner package execution | 実在 |
| [`components/archive/archive_test.go`](../components/archive/archive_test.go) | `archive` owner package test | 実在 |
| [`components/builder/builder.go`](../components/builder/builder.go) | `builder` owner package 公開実行境界 | 実在 |
| [`components/builder/model.go`](../components/builder/model.go) | `builder` owner package model | 実在 |
| [`components/builder/validate.go`](../components/builder/validate.go) | `builder` owner package validation | 実在 |
| [`components/builder/execute.go`](../components/builder/execute.go) | `builder` owner package execution | 実在 |
| [`components/builder/builder_test.go`](../components/builder/builder_test.go) | `builder` owner package test | 実在 |
| [`components/commitstatus/commitstatus.go`](../components/commitstatus/commitstatus.go) | `commitstatus` owner package 公開実行境界 | 実在 |
| [`components/commitstatus/model.go`](../components/commitstatus/model.go) | `commitstatus` owner package model | 実在 |
| [`components/commitstatus/validate.go`](../components/commitstatus/validate.go) | `commitstatus` owner package validation | 実在 |
| [`components/commitstatus/execute.go`](../components/commitstatus/execute.go) | `commitstatus` owner package execution | 実在 |
| [`components/commitstatus/commitstatus_test.go`](../components/commitstatus/commitstatus_test.go) | `commitstatus` owner package test | 実在 |
| [`components/mcp/mcp.go`](../components/mcp/mcp.go) | `mcp` owner package 公開実行境界 | 実在 |
| [`components/mcp/model.go`](../components/mcp/model.go) | `mcp` owner package model | 実在 |
| [`components/mcp/validate.go`](../components/mcp/validate.go) | `mcp` owner package validation | 実在 |
| [`components/mcp/execute.go`](../components/mcp/execute.go) | `mcp` owner package execution | 実在 |
| [`components/mcp/mcp_test.go`](../components/mcp/mcp_test.go) | `mcp` owner package test | 実在 |
| [`components/release/release.go`](../components/release/release.go) | `release` owner package 公開実行境界 | 実在 |
| [`components/release/model.go`](../components/release/model.go) | `release` owner package model | 実在 |
| [`components/release/validate.go`](../components/release/validate.go) | `release` owner package validation | 実在 |
| [`components/release/execute.go`](../components/release/execute.go) | `release` owner package execution | 実在 |
| [`components/release/release_test.go`](../components/release/release_test.go) | `release` owner package test | 実在 |
| [`components/runner/runner.go`](../components/runner/runner.go) | `runner` owner package 公開実行境界 | 実在 |
| [`components/runner/model.go`](../components/runner/model.go) | `runner` owner package model | 実在 |
| [`components/runner/validate.go`](../components/runner/validate.go) | `runner` owner package validation | 実在 |
| [`components/runner/execute.go`](../components/runner/execute.go) | `runner` owner package execution | 実在 |
| [`components/runner/runner_test.go`](../components/runner/runner_test.go) | `runner` owner package test | 実在 |
| [`components/security/security.go`](../components/security/security.go) | `security` owner package 公開実行境界 | 実在 |
| [`components/security/model.go`](../components/security/model.go) | `security` owner package model | 実在 |
| [`components/security/validate.go`](../components/security/validate.go) | `security` owner package validation | 実在 |
| [`components/security/execute.go`](../components/security/execute.go) | `security` owner package execution | 実在 |
| [`components/security/security_test.go`](../components/security/security_test.go) | `security` owner package test | 実在 |
| [`components/setup/setup.go`](../components/setup/setup.go) | `setup` owner package 公開実行境界 | 実在 |
| [`components/setup/model.go`](../components/setup/model.go) | `setup` owner package model | 実在 |
| [`components/setup/validate.go`](../components/setup/validate.go) | `setup` owner package validation | 実在 |
| [`components/setup/execute.go`](../components/setup/execute.go) | `setup` owner package execution | 実在 |
| [`components/setup/setup_test.go`](../components/setup/setup_test.go) | `setup` owner package test | 実在 |
| [`components/statefile/statefile.go`](../components/statefile/statefile.go) | `statefile` owner package 公開実行境界 | 実在 |
| [`components/statefile/model.go`](../components/statefile/model.go) | `statefile` owner package model | 実在 |
| [`components/statefile/validate.go`](../components/statefile/validate.go) | `statefile` owner package validation | 実在 |
| [`components/statefile/execute.go`](../components/statefile/execute.go) | `statefile` owner package execution | 実在 |
| [`components/statefile/statefile_test.go`](../components/statefile/statefile_test.go) | `statefile` owner package test | 実在 |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | `sdk` | 実在 |
| [`admin/index.html`](../admin/index.html) | `ui` | 実在 |
| [`testdata/builder/`](../testdata/builder/) | `builder` fixture root | 実在 |
| [`testdata/builder/single/`](../testdata/builder/single/) | `builder` fixture | 実在 |
| [`testdata/builder/single/source.md`](../testdata/builder/single/source.md) | `builder` fixture input | 実在 |
| [`testdata/builder/single/expected/`](../testdata/builder/single/expected/) | `builder` expected | 実在 |
| [`testdata/builder/site/`](../testdata/builder/site/) | `builder` fixture | 実在 |
| [`testdata/builder/site/expected/`](../testdata/builder/site/expected/) | `builder` expected | 実在 |
| [`testdata/builder/empty-dir/`](../testdata/builder/empty-dir/) | `builder` fixture | 実在 |
| [`testdata/builder/empty-dir/.keep`](../testdata/builder/empty-dir/.keep) | `builder` fixture marker | 実在 |
| [`testdata/builder/empty-dir/expected/`](../testdata/builder/empty-dir/expected/) | `builder` expected | 実在 |
| [`testdata/builder/strict/`](../testdata/builder/strict/) | `builder` fixture | 実在 |
| [`testdata/builder/strict/expected/`](../testdata/builder/strict/expected/) | `builder` expected | 実在 |
| [`testdata/builder/safe/`](../testdata/builder/safe/) | `builder` fixture | 実在 |
| [`testdata/builder/safe/expected/`](../testdata/builder/safe/expected/) | `builder` expected | 実在 |
| [`testdata/builder/url-safety/`](../testdata/builder/url-safety/) | `builder` fixture | 実在 |
| [`testdata/builder/url-safety/expected/`](../testdata/builder/url-safety/expected/) | `builder` expected | 実在 |
| [`testdata/runner/`](../testdata/runner/) | `runner` fixture root | 実在 |
| [`testdata/runner/success-runner-phase11-root-coverage/`](../testdata/runner/success-runner-phase11-root-coverage/) | `runner` Phase 11 root coverage fixture | 実在 |
| [`testdata/api/`](../testdata/api/) | `api` fixture root | 実在 |
| [`testdata/api/success-api-phase11-root-coverage/`](../testdata/api/success-api-phase11-root-coverage/) | `api` Phase 11 root coverage fixture | 実在 |
| `testdata/api/additional-management/` | 追加管理 API fixture root | 未作成 |
| [`testdata/admin/`](../testdata/admin/) | `admin` fixture group root | 実在 |
| [`testdata/admin/cli/`](../testdata/admin/cli/) | Admin CLI formal fixture root | 実在 |
| [`testdata/sdk/`](../testdata/sdk/) | `sdk` fixture root | 実在 |
| [`testdata/sdk/success-sdk-phase11-root-coverage/`](../testdata/sdk/success-sdk-phase11-root-coverage/) | `sdk` Phase 11 root coverage fixture | 実在 |
| `testdata/sdk/additional-management/` | 追加管理 SDK fixture root | 未作成 |
| [`testdata/ui/`](../testdata/ui/) | `ui` fixture root | 実在 |
| [`testdata/ui/success-ui-phase11-root-coverage/`](../testdata/ui/success-ui-phase11-root-coverage/) | `ui` Phase 11 root coverage fixture | 実在 |
| `testdata/ui/additional-management/` | 追加管理 UI fixture root | 未作成 |
| [`testdata/statefile/`](../testdata/statefile/) | `statefile` fixture root | 実在 |
| [`testdata/statefile/success-statefile-phase11-root-coverage/`](../testdata/statefile/success-statefile-phase11-root-coverage/) | `statefile` Phase 11 root coverage fixture | 実在 |
| `testdata/statefile/additional-management/` | 追加管理 statefile fixture root | 未作成 |
| [`testdata/archive/`](../testdata/archive/) | `archive` fixture root | 実在 |
| [`testdata/archive/success-archive-phase11-root-coverage/`](../testdata/archive/success-archive-phase11-root-coverage/) | `archive` Phase 11 root coverage fixture | 実在 |
| [`testdata/commitstatus/`](../testdata/commitstatus/) | `commitstatus` fixture root | 実在 |
| [`testdata/commitstatus/success-commitstatus-phase11-root-coverage/`](../testdata/commitstatus/success-commitstatus-phase11-root-coverage/) | `commitstatus` Phase 11 root coverage fixture | 実在 |
| [`testdata/security/`](../testdata/security/) | `security` fixture root | 実在 |
| [`testdata/security/success-security-phase11-root-coverage/`](../testdata/security/success-security-phase11-root-coverage/) | `security` Phase 11 root coverage fixture | 実在 |
| [`testdata/setup/`](../testdata/setup/) | `setup` fixture root | 実在 |
| [`testdata/release/`](../testdata/release/) | `release` fixture root | 実在 |
| [`testdata/mcp/`](../testdata/mcp/) | `mcp` fixture root | 実在 |
| [`testdata/phase12/quality-gate-reconstruction/`](../testdata/phase12/quality-gate-reconstruction/) | Phase 12 quality gate fixture root | 実在 |
| [`.github/workflows/phase12-quality-gate.yml`](../.github/workflows/phase12-quality-gate.yml) | Phase 12 required check workflow | 実在 |
| `testdata/phase13/implementation-alignment-quality/` | Phase 13 implementation alignment quality fixture root | 未作成 |
| `.github/workflows/phase13-implementation-alignment-quality.yml` | Phase 13 required check workflow | 未作成 |

所在区分はファイルまたは path の存在だけを示す。現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、状態語彙と実装可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) を参照する。

<a id="phase-12-target-paths"></a>
**Phase 12 target path 所在：**

以下は [`docs/SPEC.md` 方針責務 §4.3](SPEC.md#sec-4-3) が定義する Phase 12 完了後の Go owner package 到達形であり、Phase 12 実装により実在化済みである。

| path pattern | 対象 owner | 所在区分 |
|--------------|------------|----------|
| `components/<owner>/` | Phase 12 対象 Go owner package directory | 実在 |
| `components/<owner>/<owner>.go` | owner の公開実行境界 | 実在 |
| `components/<owner>/model.go` | owner の入力・出力・状態 model | 実在 |
| `components/<owner>/validate.go` | owner の入力・状態・設定・権限検証 | 実在 |
| `components/<owner>/execute.go` | owner の正常系・異常系実行順序 | 実在 |
| `components/<owner>/<owner>_test.go` | owner package の仕様契約検証 | 実在 |
| `testdata/phase12/quality-gate-reconstruction/` | Phase 12 正式 fixture root。closure record set 所在、manifest、input、expected を持つ。 | 実在 |
| `.github/workflows/phase12-quality-gate.yml` | Phase 12 required check workflow。[`docs/details/fixture.md` fixture 証跡責務 Phase 12 実装品質ゲート再構築証跡](details/fixture.md#phase-12-quality-gate-evidence) の required check name を実行する。 | 実在 |

`<owner>` に入る値は、[`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 12 Go owner package target owner](DETAIL_INDEX.md#phase-12-go-owner-package-targets) に存在する owner だけとする。JavaScript / HTML artifact である `sdk` と `ui` は、Go owner package target path として扱わず、[`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) と [`admin/index.html`](../admin/index.html) の実在行を正とする。

<a id="phase-13-target-paths"></a>
**Phase 13 target path 所在：**

以下は [`docs/DETAIL_INDEX.md` 詳細仕様入口責務 Phase 13 実装整合・品質改善参照](DETAIL_INDEX.md#phase-13-implementation-alignment-quality-entry) から参照される所在である。現在状態と完了可否は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。

| path pattern | 対象 | 所在区分 |
|--------------|------|----------|
| `components/<owner>/<owner>.go` | owner の公開実行境界 | 実在 |
| `components/<owner>/model.go` | owner の入力・出力・状態 model | 実在 |
| `components/<owner>/validate.go` | owner の入力・状態・設定・権限検証 | 実在 |
| `components/<owner>/execute.go` | owner の正常系・異常系実行順序 | 実在 |
| `components/<owner>/<owner>_test.go` | owner package の仕様契約検証 | 実在 |
| `testdata/phase13/implementation-alignment-quality/` | Phase 13 正式 fixture root | 未作成 |
| `testdata/phase13/implementation-alignment-quality/manifest.json` | Phase 13 evidence package manifest | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/scope.json` | Phase 13 実行 scope 入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/owner_inventory.json` | Phase 13 owner package inventory 入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/contract_inventory.json` | Phase 13 API / Admin / SDK / UI / MCP / setup stdout / credential 初期化契約照合入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/state_inventory.json` | Phase 13 statefile 経路棚卸し入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/security_inventory.json` | Phase 13 security 境界棚卸し入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/input/faults.json` | Phase 13 fault injection 入力 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/expected/effects.json` | Phase 13 期待効果 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/expected/counters.json` | Phase 13 closure counter 期待値 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/closure.jsonl` | Phase 13 closure record set | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/mutation.jsonl` | Phase 13 mutation 証跡 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/race.jsonl` | Phase 13 race / concurrency 証跡 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/fault.jsonl` | Phase 13 fault injection 証跡 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/e2e.jsonl` | Phase 13 integration / E2E 証跡 | 未作成 |
| `testdata/phase13/implementation-alignment-quality/records/release.jsonl` | Phase 13 release 証跡 | 未作成 |
| `.github/workflows/phase13-implementation-alignment-quality.yml` | Phase 13 required check workflow | 未作成 |
| `LICENSE` | release governance artifact | 未作成 |
| `SECURITY.md` | release governance artifact | 未作成 |
| `CONTRIBUTING.md` | release governance artifact | 未作成 |
| `CODEOWNERS` | release governance artifact | 未作成 |
| `CHANGELOG.md` | release governance artifact | 未作成 |
