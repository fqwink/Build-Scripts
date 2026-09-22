# Build-Scripts — 文書索引

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。責務分離と重複禁止の方針は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a を正本とする。

仕様・詳細仕様・デザイン正本・[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は `docs/` 配下に集約する。ルールブック [`AGENTS.md`](../AGENTS.md) と入口文書 [`README.md`](../README.md) はリポジトリ root に置く。

## 索引責務の使い方

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は、読みたい目的から参照先を選ぶための索引である。仕様本文、判定本文、実装可否本文は、以下の表に示す各責務正本を参照する。

| 読みたいもの | 参照先 | 参照理由 |
|--------------|--------|----------|
| 作業してよい条件 | [`AGENTS.md`](../AGENTS.md) | 承認、Git、PR、文書整合の最上位ルールを確認する。 |
| 初見向け概要 | [`README.md`](../README.md) | 最小限の入口と主要参照先を確認する。 |
| 文書構造 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書の役割、所在、実装ファイル所在を確認する。 |
| 方針・ポリシー | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、正本参照先、禁止事項、リリース判断を確認する。 |
| 状態・計画 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態分類、実装可否、Phase、機能インベントリ、将来計画、追加仕様化機能参照、横断補足契約を確認する。 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、対応表、ソース配置を確認する。 |
| 詳細仕様本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン関係を確認する。 |
| 実装所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の [実装ファイル一覧](#実装ファイル一覧) | 実装ファイル、テスト、fixture の所在を確認する。 |

## 目的別参照先

| 目的 | 最初に読む文書 | 次に確認する文書 | 確認内容 |
|------|----------------|------------------|----------|
| 作業ルール、承認、Git 運用を確認したい | [`AGENTS.md`](../AGENTS.md) | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 作業開始可否、変更承認、PR 作成、文書整合手順の確認先。 |
| リポジトリ全体の文書構造を把握したい | [`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 各文書の責務範囲と所在の確認先。 |
| 方針、ポリシー、禁止事項を確認したい | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 対象領域に適用する原則と制約の確認先。 |
| 状態分類、実装可否、Phase、将来計画を確認したい | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 状態分類、実装可否、Phase、将来計画の確認先。 |
| 詳細仕様本文を探したい | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象機能の owner component、参照節、受け入れ条件の確認先。 |
| 実装ファイル、テスト、fixture の所在を確認したい | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j | 現行実装ファイル、将来追加予定 path、標準配置の確認先。 |
| component の入出力、状態、処理順序、異常系、検証条件を確認したい | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | collaborator component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 実装時に従う具体仕様と collaborator 境界の確認先。 |
| 生成静的 Web サイトの見た目を確認したい | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | [`docs/details/builder.md`](details/builder.md) | 生成 HTML のデザイン関係の確認先。 |

## リポジトリ文書索引

以降は、リポジトリ内文書と実装ファイル所在の索引である。

| 索引領域 | 確認する節 | 確認内容 |
|----------|------------|--------------|
| 参照順序 | [参照順序](#参照順序) | 作業開始から詳細仕様本文までの確認先。 |
| 文書一覧 | [文書一覧](#文書一覧) | 各文書の役割と所在。 |
| 詳細仕様管理 | [詳細仕様管理](#詳細仕様管理) | owner component 別詳細本文責務の配置。 |
| 実装ファイル一覧 | [実装ファイル一覧](#実装ファイル一覧) | 実装ファイル、テスト、fixture の所在。 |
| 正本参照先 | [正本参照先](#正本参照先) | 対象ごとの正本参照先。 |

## 参照順序

| 順序 | 参照先 | 目的 |
|------|----------|------|
| 1 | [`AGENTS.md`](../AGENTS.md) | 作業ルール、承認、Git 運用、文書整合ルールを確認する。 |
| 2 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、禁止事項、リリース判断、実装着手可否を確認する。 |
| 3 | [`README.md`](../README.md) | 利用入口と主要参照先を確認する。 |
| 4 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書と実装ファイルの役割、正本参照先、配置を確認する。 |
| 5 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成 HTML のデザイン関係を確認する。 |
| 6 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態分類、実装可否、Phase、将来計画、[`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 の追加仕様化機能参照、横断補足契約を確認する。 |
| 7 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、対応表、ソース配置を確認する。 |
| 8 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |

文書一覧表は参照先索引である。作業開始前の最上位確認は [`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務の両方を先に読むことを必須とする。

## 文書一覧

| ファイル | 役割 |
|---------|------|
| [`README.md`](../README.md) | 初見向け入口。概要と主要参照先を示す。 |
| [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針、ポリシー、正本参照先、禁止事項、リリース判断を示す。 |
| [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態分類、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約を示す。 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、対応表、リポジトリ内ソース配置を示す。 |
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
| [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン関係の正本。 |
| [`AGENTS.md`](../AGENTS.md) | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書・実装ファイルの役割と所在を示す。 |

## 詳細仕様管理

詳細仕様管理表は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の配置先を示す索引である。詳細仕様参照入口、共通固定値、対応表、リポジトリ内ソース配置は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、追加仕様化機能参照と横断補足契約は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6、各ファイルの責務境界と owner / collaborator の扱いは [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b および [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0b.1 を参照する。

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

## 実装ファイル一覧

実装ファイル一覧は、現行リポジトリに存在する実装ファイル、テスト、fixture、および将来追加予定 path の所在を示す索引である。状態分類と実装可否は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。

リポジトリ内ソース配置は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 のディレクトリ構成と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j を参照する。

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の実装ファイル一覧表は、現行リポジトリに存在する実装ファイルと、将来追加予定 path を区別して示す。`main.go`、`components/*.go`、`admin/` 配下の静的 UI ファイル、`testdata/<component>/` を現行配置として扱う。

標準配置図に含まれる未作成 path は、将来追加予定 path として扱い、該当 owner component が実装対象になった変更で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成しない。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は標準外配置である。標準配置と現行実体の所在確認は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j、実装ファイル一覧を同時に確認する。

| パス | component | 所在区分 | 役割 |
|------|-----------|----------|------|
| `main.go` | `-` | 現行実体 | 起動入口。現時点では実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | `builder` | 現行実体 | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `components/builder_test.go` | `builder` | 現行実体 | `components/builder.go` の Go 検証ファイル。 |
| `go.mod` | `-` | 現行実体 | Go module 定義。外部 module は追加しない。 |
| `testdata/builder/` | `builder` | 現行実体 | builder 用 fixture 入力。 |
| `components/runner.go` | `runner` | 現行実体 | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。 |
| `components/runner_test.go` | `runner` | 現行実体 | `components/runner.go` の Go 検証ファイル。 |
| `components/api.go` | `api` | 現行実体 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `components/api_test.go` | `api` | 現行実体 | `components/api.go` の API endpoint、認証、状態ファイル、管理操作の検証テスト。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 現行実体 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | `ui` | 現行実体 | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | `mcp` | 将来追加予定 | MCP サーバー。現時点では未作成であり、実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## 正本参照先

| 正本対象 | 参照先 |
|----------|------|
| 作業ルール、承認、Git 運用、文書整合 | [`AGENTS.md`](../AGENTS.md) |
| 方針、ポリシー、正本参照先、禁止事項、リリース判断 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 |
| 状態分類、実装可否、Phase、機能インベントリ、将来計画、昇格手順、追加仕様化機能参照、横断補足契約 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 |
| 詳細仕様参照入口、索引、共通固定値、実装前確認項目、検証マトリクス、詳細節対応表、リポジトリ内ソース配置 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 |
| 各 component の入出力、状態、処理順序、異常系、検証条件の本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 |
| 生成静的 Web サイトのデザイン関係 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 |
| 文書・実装ファイルの参照先と役割 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 |
| 実装ファイル、テスト、fixture の所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の [実装ファイル一覧](#実装ファイル一覧) |

文書整合の禁止事項、責務分離、参照リンク化、重複禁止は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a、作業手順は [`AGENTS.md`](../AGENTS.md) を参照する。
