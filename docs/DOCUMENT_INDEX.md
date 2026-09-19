# Build-Scripts — 文書索引

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

仕様・詳細仕様・デザイン正本・本索引は `docs/` 配下に集約する。ルールブック [`AGENTS.md`](../AGENTS.md) と入口文書 [`README.md`](../README.md) はリポジトリ root に置く。

## 索引責務の使い方

本ファイルは、読みたい目的から参照先を選ぶための索引である。方針、ポリシー、正本参照先、禁止事項、責務ベース明示的原則は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務を参照する。実装状態、実装可否、Phase、将来計画は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。詳細仕様入口と owner component 参照表は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務を参照する。

| 読みたいもの | 参照先 | 参照理由 |
|--------------|--------|----------|
| 作業してよい条件 | [`AGENTS.md`](../AGENTS.md) | 承認、Git、PR、文書整合の最上位ルールを確認する。 |
| 初見向け概要 | [`README.md`](../README.md) | 最小限の入口と主要参照先を確認する。 |
| 文書構造 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書の役割、所在、実装ファイル所在を確認する。 |
| 方針・ポリシー | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、正本参照先、禁止事項、リリース判断を確認する。 |
| 状態・計画 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 実装状態、実装可否、Phase、機能インベントリ、将来計画、追加仕様化機能参照、横断補足契約を確認する。 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様の読み方、共通固定値、対応表、ソース配置を確認する。 |
| 詳細仕様本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン方針と視覚仕様を確認する。 |
| 実装所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の [仕様化済みコンポーネント](#仕様化済みコンポーネント) | 実装ファイル、テスト、fixture の所在と状態を確認する。 |

## 目的別参照先

| 目的 | 最初に読む文書 | 次に確認する文書 | 判断内容 |
|------|----------------|------------------|----------|
| 作業ルール、承認、Git 運用を確認したい | [`AGENTS.md`](../AGENTS.md) | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 作業開始可否、変更承認、PR 作成、文書整合の手順を判断する。 |
| リポジトリ全体の文書構造を把握したい | [`README.md`](../README.md) | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | どの文書が何を持つか、どの順番で読むかを判断する。 |
| 方針、ポリシー、禁止事項を判断したい | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 対象領域に適用する原則と制約を判断する。 |
| 実装状態、実装可否、Phase、将来計画を判断したい | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象が未仕様化、将来計画、改訂予定、仕様化済み・未実装、実装中・検証未完了、実装済みのどれかを判断する。 |
| 詳細仕様本文を探したい | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象機能の owner component、参照節、受け入れ条件を判断する。 |
| 実装ファイル、テスト、fixture の所在を確認したい | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j | 現行実装ファイル、将来追加予定 path、標準配置の扱いを判断する。 |
| component の入出力、状態、処理順序、異常系、検証条件を確認したい | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 実装時に従う具体仕様と collaborator 境界を判断する。 |
| 生成静的 Web サイトの見た目を確認したい | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | [`docs/details/builder.md`](details/builder.md) | レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を判断する。 |

## リポジトリ文書索引

本節以降は、リポジトリ内文書と実装ファイル所在の索引である。仕様判断の本文は各正本を参照する。

| 索引領域 | 確認する節 | 判断する内容 |
|----------|------------|--------------|
| 読む順番 | [読む順番](#読む順番) | 作業開始から詳細仕様本文までの確認順序。 |
| 文書一覧 | [文書一覧](#文書一覧) | 各文書の役割と所在。 |
| 実装所在入口 | [実装ファイル索引](#実装ファイル索引) | 実装ファイル所在の判断原則。 |
| 詳細仕様管理 | [詳細仕様管理](#詳細仕様管理) | owner component 別詳細本文責務の配置と状態。 |
| 実装ファイル一覧 | [仕様化済みコンポーネント](#仕様化済みコンポーネント) | 実装ファイル、テスト、fixture の所在と状態。 |
| 正本参照先 | [正本参照先](#正本参照先) | 判断対象ごとの正本参照先。 |
| 整合確認先 | [整合確認先](#整合確認先) | 文書整合で確認する正本への参照。 |
| 整合メモ | [整合メモ](#整合メモ) | 現行状態に関する注意点。 |

## 読む順番

| 順序 | ファイル | 目的 |
|------|----------|------|
| 1 | [`AGENTS.md`](../AGENTS.md) | 作業ルール、承認、Git 運用、文書整合ルールを確認する。 |
| 2 | [`README.md`](../README.md) | 初見向けの概要、読む順番、主要ファイルを確認する。 |
| 3 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書と実装ファイルの役割、正本参照先、配置を確認する。 |
| 4 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、禁止事項、リリース判断を確認する。 |
| 5 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成 HTML のデザイン関係を確認する。 |
| 6 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 実装状態、実装可否、Phase、将来計画、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能参照、横断補足契約を確認する。 |
| 7 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様の入口、読み順、共通固定値、対応表、ソース配置を確認する。 |
| 8 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |

上記の順序は、文書整理、仕様改訂、実装、検証、PR 作成のすべてで共通とする。[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は索引であり、仕様判断の正本ではない。

## 文書一覧

| ファイル | 役割 |
|---------|------|
| [`README.md`](../README.md) | 初見向け入口。概要、読む順番、主要ファイルを示す。 |
| [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、正本参照先、禁止事項、リリース判断を示す。 |
| [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約を示す。 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様の入口、読み順、共通固定値、対応表、リポジトリ内ソース配置を示す。 |
| [`docs/details/builder.md`](details/builder.md) | `builder` owner component の詳細本文責務。 |
| [`docs/details/runner.md`](details/runner.md) | `runner` owner component の詳細本文責務。 |
| [`docs/details/api.md`](details/api.md) | `api` owner component の詳細本文責務。 |
| [`docs/details/admin.md`](details/admin.md) | `admin` owner component の詳細本文責務。 |
| [`docs/details/sdk.md`](details/sdk.md) | `sdk` owner component の詳細本文責務。 |
| [`docs/details/ui.md`](details/ui.md) | `ui` owner component の詳細本文責務。 |
| [`docs/details/setup.md`](details/setup.md) | `setup` owner component の詳細本文責務。 |
| [`docs/details/statefile.md`](details/statefile.md) | `statefile` owner component の詳細本文責務。 |
| [`docs/details/archive.md`](details/archive.md) | `archive` owner component の詳細本文責務。 |
| [`docs/details/commitstatus.md`](details/commitstatus.md) | `commitstatus` owner component の詳細本文責務。 |
| [`docs/details/security.md`](details/security.md) | `security` owner component の詳細本文責務。 |
| [`docs/details/fixture.md`](details/fixture.md) | `fixture` owner component の詳細本文責務。 |
| [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン関係の正本。デザイン方針、レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を示す。 |
| [`AGENTS.md`](../AGENTS.md) | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 本索引。文書・実装ファイルの役割と所在を示す。仕様本文を定義しない。 |

## 実装ファイル索引

実装ファイルの所在と状態は [仕様化済みコンポーネント](#仕様化済みコンポーネント) を参照する。実装ファイルが存在することだけで、仕様化済み、実装可、完了済みとは判断しない。

## 詳細仕様管理

詳細仕様本文の所在は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務を参照する。[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務は、入口、読み順、共通固定値、対応表、リポジトリ内ソース配置を示す。[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 は、追加仕様化機能参照と横断補足契約を示す。

下表は、詳細仕様本文の配置先を示す索引である。各ファイルの責務境界、持つ内容、持たない内容、owner / collaborator の扱いは [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b および [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 を参照する。

| ファイル | owner component |
|----------|-----------------|
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
| [`docs/details/fixture.md`](details/fixture.md) | `fixture` |

詳細仕様本文の配置先は上表のとおりである。実装状態と実装可否は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、owner component の特定は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務を参照する。

## 仕様化済みコンポーネント

[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務では、以下のコンポーネントも仕様化されている。

リポジトリ内ソース配置は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 のディレクトリ構成と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j を参照する。

下表は、現行リポジトリに存在する実装ファイルと、将来追加予定 path を区別して示す。`main.go`、`components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` を現行配置として扱う。

標準配置図に含まれる未作成 path は、将来追加予定 path として扱い、該当 owner component が実装対象になった PR で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成しない。

| パス | component | 状態 | 役割 |
|------|-----------|------|------|
| `main.go` | `-` | 実装済み | 起動入口。現時点では実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | `builder` | 実装済み | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `components/builder_test.go` | `builder` | 実装済み | `components/builder.go` の Phase 1 fixture テスト。 |
| `go.mod` | `-` | 実装済み | Go module 定義。外部 module は追加しない。 |
| `testdata/builder/` | `builder` | 実装済み | Phase 1 の受け入れ fixture 入力。 |
| `components/runner.go` | `runner` | 実装済み | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。Phase 2 完了判定パスを対象とする。 |
| `components/runner_test.go` | `runner` | 実装済み | `components/runner.go` の Phase 2 fixture、hardening、完了判定パステスト。 |
| `components/api.go` | `api` | 実装済み | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `components/api_test.go` | `api` | 実装済み | `components/api.go` の API endpoint、認証、状態ファイル、管理操作の検証テスト。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 実装済み | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | `ui` | 実装済み | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | `mcp` | 将来追加予定 path | MCP サーバー。現時点では未作成であり、実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## 正本参照先

| 判断対象 | 参照先 |
|----------|------|
| 作業ルール、承認、Git 運用、文書整合 | [`AGENTS.md`](../AGENTS.md) |
| 方針、ポリシー、正本参照先、禁止事項、リリース判断 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 |
| 実装状態、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 |
| 詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 |
| 各 component の入出力、状態、処理順序、異常系、検証条件の本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 |
| 生成静的 Web サイトのデザイン関係 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 |
| 文書・実装ファイルの参照先と役割 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 |
| 実装ファイル、テスト、fixture の所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の [仕様化済みコンポーネント](#仕様化済みコンポーネント) |

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は索引であり、仕様本文、デザイン本文、実装可否判断の正本ではない。

## 整合確認先

文書整合の禁止事項、責務分離、参照リンク化、重複禁止は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a を参照する。作業ルール、承認、Git 操作、PR 作成、検証手順は [`AGENTS.md`](../AGENTS.md) を参照する。本ファイルでは確認先だけを示す。

## 整合メモ

現時点では、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務に記載された一部機能は仕様化済みだが、リポジトリ内に実装コードが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。

現行実装実体は `main.go`、`components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` である。`components/builder.go`、`components/runner.go`、`components/api.go`、`admin/adlaire-ci-sdk.js`、`admin/index.html` は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務の実装状態と本ファイルの [仕様化済みコンポーネント](#仕様化済みコンポーネント) に従って実装済みとして扱う。Go toolchain による `gofmt` と `go test` の検証対象は Go ファイルとする。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は標準外配置であり、現行実装実体として扱わない。標準配置と現行実体の判断は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j、本ファイルの [仕様化済みコンポーネント](#仕様化済みコンポーネント) を同時に確認する。
