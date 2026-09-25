# Build-Scripts — 文書・実装ファイル所在索引

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) は、Build-Scripts リポジトリ内の文書、実装ファイル、テスト、fixture、未作成 path の所在だけを示す索引である。方針、状態定義、現在状態、実装詳細、デザイン、検証契約は本文として定義しない。

## 索引責務の使い方

| 確認対象 | 参照先 |
|----------|--------|
| 作業ルール | [`AGENTS.md`](../AGENTS.md) |
| 方針、ポリシー、状態語彙、状態遷移条件 | [`docs/SPEC.md`](SPEC.md) |
| 各 component と各機能の現在状態、Phase、将来計画 | [`docs/ROADMAP.md`](ROADMAP.md) |
| 詳細仕様入口、共通固定値、owner / collaborator 対応 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) |
| owner component 別詳細本文 | [`docs/details/`](details/) |
| fixture、expected、fake、実装検証証跡 | [`docs/details/fixture.md`](details/fixture.md) |
| 生成 HTML のデザイン | [`docs/DESIGN.md`](DESIGN.md) |
| 利用入口 | [`README.md`](../README.md) |

## 目的別参照先

目的別の正本所在は [索引責務の使い方](#索引責務の使い方) を参照する。対象機能の詳細本文を探す場合は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) から owner component を特定する。

## リポジトリ文書索引

リポジトリ文書の実在所在は [文書一覧](#文書一覧)、owner component 別詳細本文の所在は [詳細仕様本文の所在](#詳細仕様本文の所在)、実装・テスト・fixture の所在は [実装ファイル一覧](#実装ファイル一覧) を参照する。

## 参照順序

作業開始時の必須読了順序は [`AGENTS.md`](../AGENTS.md) を正とする。仕様を追跡する場合は [`docs/SPEC.md`](SPEC.md)、[`docs/ROADMAP.md`](ROADMAP.md)、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md)、対象 owner component の詳細本文、[`docs/details/fixture.md`](details/fixture.md) の順に所在を確認する。

## 文書一覧

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

## 詳細仕様本文の所在

| パス | owner component / 責務 |
|------|-------------------------|
| [`docs/details/builder.md`](details/builder.md) | `builder` |
| [`docs/details/runner.md`](details/runner.md) | `runner` |
| [`docs/details/api.md`](details/api.md) | `api` |
| [`docs/details/admin.md`](details/admin.md) | `admin` |
| [`docs/details/sdk.md`](details/sdk.md) | `sdk` |
| [`docs/details/ui.md`](details/ui.md) | `ui` |
| [`docs/details/setup.md`](details/setup.md) | `setup` |
| [`docs/details/statefile.md`](details/statefile.md) | `statefile` |
| [`docs/details/archive.md`](details/archive.md) | `archive` |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | `commitstatus` |
| [`docs/details/security.md`](details/security.md) | `security` |
| [`docs/details/fixture.md`](details/fixture.md) | fixture 証跡責務 |

## 実装ファイル一覧

| パス | component | 所在区分 |
|------|-----------|----------|
| [`main.go`](../main.go) | 起動入口 | 実在 |
| [`go.mod`](../go.mod) | Go module | 実在 |
| [`components/builder.go`](../components/builder.go) | `builder` | 実在 |
| [`components/builder_test.go`](../components/builder_test.go) | `builder` test | 実在 |
| [`components/runner.go`](../components/runner.go) | `runner` | 実在 |
| [`components/runner_test.go`](../components/runner_test.go) | `runner` test | 実在 |
| [`components/api.go`](../components/api.go) | `api` | 実在 |
| [`components/api_test.go`](../components/api_test.go) | `api` test | 実在 |
| [`admin/adlaire-ci-sdk.js`](../admin/adlaire-ci-sdk.js) | `sdk` | 実在 |
| [`admin/index.html`](../admin/index.html) | `ui` | 実在 |
| [`testdata/builder/`](../testdata/builder/) | `builder` fixture root | 実在 |
| [`testdata/builder/single/source.md`](../testdata/builder/single/source.md) | `builder` fixture | 実在 |
| [`testdata/builder/site/`](../testdata/builder/site/) | `builder` fixture | 実在 |
| `testdata/builder/empty-dir/.keep` | `builder` fixture marker | 実在 |
| `testdata/builder/strict/` | `builder` fixture | 未作成 |
| `testdata/builder/safe/` | `builder` fixture | 未作成 |
| `testdata/builder/**/expected/` | `builder` expected | 未作成 |
| `testdata/runner/` | `runner` fixture root | 未作成 |
| `testdata/api/` | `api` fixture root | 未作成 |
| `testdata/sdk/` | `sdk` fixture root | 未作成 |
| `testdata/ui/` | `ui` fixture root | 未作成 |
| `components/mcp.go` | `mcp` | 未作成 |

所在区分はファイルまたは path の存在だけを示す。現在状態は [`docs/ROADMAP.md`](ROADMAP.md)、状態語彙と実装可否は [`docs/SPEC.md` ポリシー責務 §0a](SPEC.md#0a-仕様成熟度ポリシー) を参照する。

## 正本参照先

正本所在は [索引責務の使い方](#索引責務の使い方) の表を唯一の一覧とする。
