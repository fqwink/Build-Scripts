# Build-Scripts — 文書索引

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。仕様本文、方針本文、状態本文、詳細本文、デザイン本文を持たない。責務分離と重複禁止の方針は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a、正本所在は [正本参照先](#正本参照先) を正本とする。

方針・ポリシー、状態・計画、詳細仕様、デザイン、[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は `docs/` 配下に集約する。ルールブック [`AGENTS.md`](../AGENTS.md) と入口文書 [`README.md`](../README.md) はリポジトリ root に置く。

索引内で方針、状態、詳細本文、デザイン本文に触れる場合は、対象文書の内容を再掲せず、責務名付きリンクで参照する。

## 索引責務の使い方

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務は、読みたい目的から参照先を選ぶための索引である。仕様本文、判定本文、実装可否本文は持たず、正本所在は [正本参照先](#正本参照先) を参照する。

| 読みたいもの | 参照先 | 参照理由 |
|--------------|--------|----------|
| 作業してよい条件 | [`AGENTS.md`](../AGENTS.md) | 作業ルール正本。 |
| 初見向け概要 | [`README.md`](../README.md) | 利用入口。 |
| 文書構造 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書・実装所在の索引。 |
| 方針・ポリシー | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針・ポリシー正本。 |
| 状態・計画 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態・計画正本。 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様入口。 |
| 詳細仕様本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | owner component 別詳細本文。 |
| デザイン責務 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン正本。 |
| 実装所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の [実装ファイル一覧](#実装ファイル一覧) | 実装ファイル、テスト、fixture の所在。 |

## 目的別参照先

| 目的 | 最初に読む文書 | 次に確認する文書 | 確認内容 |
|------|----------------|------------------|----------|
| 作業ルール、承認、Git 運用を確認したい | [`AGENTS.md`](../AGENTS.md) | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 作業開始可否、変更承認、PR 作成、文書整合手順の確認先。 |
| リポジトリ全体の文書構造を把握したい | [`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 各文書の責務範囲と所在の確認先。 |
| 方針、ポリシー、禁止事項を確認したい | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | - | 方針・ポリシー・禁止事項は [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務だけを正本とする。 |
| 状態分類、実装可否、Phase、将来計画を確認したい | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 状態・計画正本への入口。 |
| 詳細仕様本文を探したい | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象機能の owner component、参照節、受け入れ条件の確認先。 |
| 実装ファイル、テスト、fixture の所在を確認したい | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j | 現行実装ファイル、将来追加予定 path、標準配置の確認先。 |
| component の入出力、状態、処理順序、異常系、検証条件を確認したい | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 対象 owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | owner component を確定してから、実装時に従う具体仕様と collaborator 境界を確認する。 |
| 生成静的 Web サイトの見た目を確認したい | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | - | 生成 HTML のデザイン関係は [`docs/DESIGN.md`](DESIGN.md) デザイン責務だけを正本とする。 |

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
| 2 | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針・ポリシー正本を確認する。 |
| 3 | [`README.md`](../README.md) | 利用入口と主要参照先を確認する。 |
| 4 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書と実装ファイルの役割、正本参照先、配置を確認する。 |
| 5 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態・計画正本を確認する。 |
| 6 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、対応表、ソース配置を確認する。 |
| 7 | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成 HTML のデザイン関係を確認する。 |
| 8 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 | 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。 |

文書一覧表は参照先索引である。作業開始前の最上位確認は [`AGENTS.md`](../AGENTS.md) と [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務の両方を先に読むことを必須とする。

## 文書一覧

| ファイル | 役割 |
|---------|------|
| [`README.md`](../README.md) | 初見向け入口。概要と主要参照先を示す。 |
| [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 | 方針・ポリシー正本。 |
| [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 | 状態・計画正本。 |
| [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 | 詳細仕様参照入口、共通固定値、対応表、リポジトリ内ソース配置を示す。 |
| [`docs/details/builder.md`](details/builder.md) 詳細本文責務 | `builder` owner component の詳細本文責務。 |
| [`docs/details/runner.md`](details/runner.md) 詳細本文責務 | `runner` owner component の詳細本文責務。 |
| [`docs/details/api.md`](details/api.md) 詳細本文責務 | `api` owner component の詳細本文責務。 |
| [`docs/details/admin.md`](details/admin.md) 詳細本文責務 | `admin` owner component の詳細本文責務。 |
| [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務 | `sdk` owner component の詳細本文責務。 |
| [`docs/details/ui.md`](details/ui.md) 詳細本文責務 | `ui` owner component の詳細本文責務。 |
| [`docs/details/setup.md`](details/setup.md) 詳細本文責務 | `setup` owner component の詳細本文責務。 |
| [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務 | `statefile` owner component の詳細本文責務。 |
| [`docs/details/archive.md`](details/archive.md) 詳細本文責務 | `archive` owner component の詳細本文責務。 |
| [`docs/details/commitstatus.md`](details/commitstatus.md) 詳細本文責務 | `commitstatus` owner component の詳細本文責務。 |
| [`docs/details/security.md`](details/security.md) 詳細本文責務 | `security` owner component の詳細本文責務。 |
| [`docs/details/fixture.md`](details/fixture.md) 詳細本文責務 | `fixture` owner component の詳細本文責務。 |
| [`docs/DESIGN.md`](DESIGN.md) デザイン責務 | 生成静的 Web サイトのデザイン関係の正本。 |
| [`AGENTS.md`](../AGENTS.md) | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 | 文書・実装ファイルの役割と所在を示す。 |

## 詳細仕様管理

詳細仕様管理表は、owner component 別の [`docs/details/*.md`](details/) 詳細本文責務の配置先だけを示す索引である。詳細仕様入口は [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務、追加仕様化機能参照と横断補足契約は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 §6 を参照する。

| ファイル | owner component |
|----------|-----------------|
| [`docs/details/builder.md`](details/builder.md) 詳細本文責務 | `builder` |
| [`docs/details/runner.md`](details/runner.md) 詳細本文責務 | `runner` |
| [`docs/details/api.md`](details/api.md) 詳細本文責務 | `api` |
| [`docs/details/admin.md`](details/admin.md) 詳細本文責務 | `admin` |
| [`docs/details/sdk.md`](details/sdk.md) 詳細本文責務 | `sdk` |
| [`docs/details/ui.md`](details/ui.md) 詳細本文責務 | `ui` |
| [`docs/details/setup.md`](details/setup.md) 詳細本文責務 | `setup` |
| [`docs/details/statefile.md`](details/statefile.md) 詳細本文責務 | `statefile` |
| [`docs/details/archive.md`](details/archive.md) 詳細本文責務 | `archive` |
| [`docs/details/commitstatus.md`](details/commitstatus.md) 詳細本文責務 | `commitstatus` |
| [`docs/details/security.md`](details/security.md) 詳細本文責務 | `security` |
| [`docs/details/fixture.md`](details/fixture.md) 詳細本文責務 | `fixture` |

## 実装ファイル一覧

実装ファイル一覧は、現行リポジトリに存在する実装ファイル、テスト、fixture、および将来追加予定 path の所在を示す索引である。状態分類と実装可否は [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務を参照する。

リポジトリ内ソース配置は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3 のディレクトリ構成と [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j を参照する。

[`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務の実装ファイル一覧表は、現行リポジトリに存在する実装ファイルと、将来追加予定 path を区別して示す。`main.go`、実在する `components/*.go`、`admin/` 配下の静的 UI ファイル、実在する `testdata/<component>/` だけを現行実体として扱う。

標準配置図に含まれる未作成 path は、将来追加予定 path として扱い、該当 owner component が実装対象になった変更で追加する。標準配置図に含まれていることだけを理由に、未実装ファイル、将来追加予定 path、空ディレクトリ、placeholder を作成しない。

`build_spec.go`、`runner.go`、`build_spec_test.go`、`runner_test.go`、`testdata/build_spec/` は標準外配置である。標準配置と現行実体の所在確認は、[`docs/SPEC.md`](SPEC.md) 方針責務 §4.3、[`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 §0j、実装ファイル一覧を参照する。

| パス | component | 所在区分 | 役割 |
|------|-----------|----------|------|
| `main.go` | `-` | 現行実体 | 起動入口。現時点では実行ファイル名に応じて `builder` または `runner` component を呼び出す。 |
| `components/builder.go` | `builder` | 現行実体 | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `components/builder_test.go` | `builder` | 現行実体 | `components/builder.go` の Go 検証ファイル。 |
| `go.mod` | `-` | 現行実体 | Go module 定義。外部 module は追加しない。 |
| `testdata/builder/` | `builder` | 現行実体 | builder 用 fixture 入力。 |
| `components/runner.go` | `runner` | 現行実体 | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。 |
| `components/runner_test.go` | `runner` | 現行実体 | `components/runner.go` の Go 検証ファイル。 |
| `testdata/runner/` | `runner` | 未作成 path | runner fixture の配置予定 path。runner fixture を追加する実装検証変更で作成し、未作成の間は現行実体として扱わない。 |
| `components/api.go` | `api` | 現行実体 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `components/api_test.go` | `api` | 現行実体 | `components/api.go` の API endpoint、認証、状態ファイル、管理操作の検証テスト。 |
| `testdata/api/` | `api` | 未作成 path | API fixture の配置予定 path。API fixture を追加する実装検証変更で作成し、未作成の間は現行実体として扱わない。 |
| `admin/adlaire-ci-sdk.js` | `sdk` | 現行実体 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `testdata/sdk/` | `sdk` | 未作成 path | SDK fixture の配置予定 path。SDK fixture を追加する実装検証変更で作成し、未作成の間は現行実体として扱わない。 |
| `admin/index.html` | `ui` | 現行実体 | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `testdata/ui/` | `ui` | 未作成 path | UI fixture の配置予定 path。UI fixture を追加する実装検証変更で作成し、未作成の間は現行実体として扱わない。 |
| `docs/examples/` | `-` | 未作成 path | 利用例、設定例、サンプル構成を追加する場合の配置先。実装コンポーネントではなく、空ディレクトリや placeholder は作成しない。 |
| `components/mcp.go` | `mcp` | 将来追加予定 | MCP サーバー。現時点では未作成であり、実装可能な詳細仕様を持たず、MCP 専用詳細仕様が新設されるまで実装対象ではない。 |

## 正本参照先

以下は正本所在の一覧である。各正本の内容説明は再掲せず、詳細な読み分けは [索引責務の使い方](#索引責務の使い方) と [目的別参照先](#目的別参照先) を参照する。

| 正本対象 | 参照先 |
|----------|------|
| 作業ルール | [`AGENTS.md`](../AGENTS.md) |
| 方針・ポリシー | [`docs/SPEC.md`](SPEC.md) 方針責務・ポリシー責務 |
| 状態・計画 | [`docs/ROADMAP.md`](ROADMAP.md) 状態・計画責務 |
| 詳細仕様入口 | [`docs/DETAIL_INDEX.md`](DETAIL_INDEX.md) 詳細仕様入口責務 |
| 詳細本文 | owner component 別の [`docs/details/*.md`](details/) 詳細本文責務 |
| デザイン | [`docs/DESIGN.md`](DESIGN.md) デザイン責務 |
| 文書・実装ファイル所在 | [`docs/DOCUMENT_INDEX.md`](DOCUMENT_INDEX.md) 文書・実装ファイル所在の索引責務 |

文書整合の禁止事項、責務分離、参照リンク化、重複禁止は [`docs/SPEC.md`](SPEC.md) 方針責務 §4.2a、作業手順は [`AGENTS.md`](../AGENTS.md) を参照する。
