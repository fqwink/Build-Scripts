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
| Phase 11 仕様全般完了判定の所在 | [`docs/SPEC.md` 方針責務 §4.8](SPEC.md#sec-4-8) |
| fixture、expected、fake、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) |
| test artifact と fixture 証跡の接続 | [`docs/details/fixture.md` fixture 証跡責務 test artifact traceability 固定契約](details/fixture.md#test-artifact-traceability-contract) |
| 単独 test artifact を持たない owner の検証接続 | [`docs/details/fixture.md` fixture 証跡責務 non-dedicated owner test routing 固定契約](details/fixture.md#non-dedicated-owner-test-routing-contract) |
| Phase 11 fixture root 網羅判定 | [`docs/details/fixture.md` fixture 証跡責務 fixture root coverage matrix 固定契約](details/fixture.md#fixture-root-coverage-matrix-contract) |
| 検証実行証跡の分類と未完了条件 | [`docs/details/fixture.md` fixture 証跡責務 test execution evidence matrix 固定契約](details/fixture.md#test-execution-evidence-matrix-contract) |
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
| [`components/builder.go`](../components/builder.go) | `builder` | 実在 |
| [`components/builder_test.go`](../components/builder_test.go) | `builder` test | 実在 |
| [`components/runner.go`](../components/runner.go) | `runner` | 実在 |
| [`components/runner_test.go`](../components/runner_test.go) | `runner` test | 実在 |
| [`components/api.go`](../components/api.go) | `api` | 実在 |
| [`components/api_test.go`](../components/api_test.go) | `api` test | 実在 |
| [`components/admin.go`](../components/admin.go) | `admin` CLI | 実在 |
| [`components/admin_test.go`](../components/admin_test.go) | `admin` CLI test | 実在 |
| [`components/setup.go`](../components/setup.go) | `setup` | 実在 |
| [`components/setup_test.go`](../components/setup_test.go) | `setup` test | 実在 |
| [`components/release.go`](../components/release.go) | `release` | 実在 |
| [`components/release_test.go`](../components/release_test.go) | `release` test | 実在 |
| [`components/mcp.go`](../components/mcp.go) | `mcp` | 実在 |
| [`components/mcp_test.go`](../components/mcp_test.go) | `mcp` test | 実在 |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | `sdk` | 実在 |
| [`admin/index.html`](../admin/index.html) | `ui` | 実在 |
| [`testdata/builder/`](../testdata/builder/) | `builder` fixture root | 実在 |
| [`testdata/builder/single/source.md`](../testdata/builder/single/source.md) | `builder` fixture | 実在 |
| [`testdata/builder/site/`](../testdata/builder/site/) | `builder` fixture | 実在 |
| `testdata/builder/empty-dir/.keep` | `builder` fixture marker | 実在 |
| `testdata/builder/strict/` | `builder` fixture | 未作成 |
| `testdata/builder/safe/` | `builder` fixture | 未作成 |
| `testdata/builder/url-safety/` | `builder` fixture | 未作成 |
| `testdata/builder/**/expected/` | `builder` expected | 未作成 |
| `testdata/runner/` | `runner` fixture root | 未作成 |
| `testdata/api/` | `api` fixture root | 未作成 |
| `testdata/api/additional-management/` | 追加管理 API fixture root | 未作成 |
| [`testdata/admin/`](../testdata/admin/) | `admin` fixture root | 実在 |
| `testdata/sdk/` | `sdk` fixture root | 未作成 |
| `testdata/sdk/additional-management/` | 追加管理 SDK fixture root | 未作成 |
| `testdata/ui/` | `ui` fixture root | 未作成 |
| `testdata/ui/additional-management/` | 追加管理 UI fixture root | 未作成 |
| `testdata/statefile/` | `statefile` fixture root | 未作成 |
| `testdata/statefile/additional-management/` | 追加管理 statefile fixture root | 未作成 |
| `testdata/archive/` | `archive` fixture root | 未作成 |
| `testdata/commitstatus/` | `commitstatus` fixture root | 未作成 |
| `testdata/security/` | `security` fixture root | 未作成 |
| [`testdata/setup/`](../testdata/setup/) | `setup` fixture root | 実在 |
| [`testdata/release/`](../testdata/release/) | `release` fixture root | 実在 |
| [`testdata/mcp/`](../testdata/mcp/) | `mcp` fixture root | 実在 |

所在区分はファイルまたは path の存在だけを示す。現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、状態語彙と実装可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#policy-spec-maturity) を参照する。
