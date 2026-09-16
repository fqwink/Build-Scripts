# Build-Scripts Document Index

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

## Documents

| ファイル | 役割 |
|---------|------|
| `ADLAIRE_CI_SPEC.md` | Adlaire CI の方針、ポリシー、実装状態、正本関係を定めるマスター仕様書正本。 |
| `ADLAIRE_CI_DETAIL_SPEC.md` | `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様。実装対象コンポーネントの入出力、状態、処理順序、異常系、検証条件に関する正本。方針、ポリシー、実装状態、ロードマップ状態、実装可否は記載しない。 |
| `DESIGN.md` | 生成静的 Web サイトのデザイン仕様。レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を整理する。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| `DOCUMENT_INDEX.md` | 本索引。リポジトリ内の文書・実装ファイルの役割と正本関係を示す。仕様正本ではない。 |

## Specified Components

`ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` では、以下のコンポーネントも仕様化されている。

リポジトリ内ソース配置の標準構成は、`ADLAIRE_CI_SPEC.md` Part 1 §4.3 と `ADLAIRE_CI_DETAIL_SPEC.md` §0j を参照する。

下表は、現行リポジトリに存在する実装ファイルを示す。`components/` 標準配置への移行前は、現行ファイルを実装実体として扱う。

| パス | 状態 | 役割 |
|------|------|------|
| `build_spec.go` | 実装済み | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。標準移行後の配置は `components/builder.go`。 |
| `build_spec_test.go` | 実装済み | `build_spec.go` の Phase 1 fixture テスト。標準移行後のテスト配置は実装 PR で決定する。 |
| `go.mod` | 実装済み | Go module 定義。外部 module は追加しない。 |
| `testdata/build_spec/` | 実装済み | Phase 1 の受け入れ fixture 入力。標準移行後の配置は `testdata/builder/`。 |
| `runner.go` | 実装済み | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。Phase 2 完了判定パスを対象とする。標準移行後の配置は `components/runner.go`。 |
| `runner_test.go` | 実装済み | `runner.go` の Phase 2 fixture、hardening、完了判定パステスト。標準移行後のテスト配置は実装 PR で決定する。 |
| `components/api.go` | 仕様化済み・未実装 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `admin/adlaire-ci-sdk.js` | 仕様化済み・未実装 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | 仕様化済み・未実装 | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `components/mcp.go` | 将来計画 | MCP サーバー。将来追加コンポーネントとして追加予定。 |

## Source Of Truth

方針、ポリシー、実装状態、正本関係の判断では `ADLAIRE_CI_SPEC.md` を正とする。

実装の具体的詳細の判断では `ADLAIRE_CI_DETAIL_SPEC.md` を正とする。

デザイン判断では、`ADLAIRE_CI_SPEC.md` と矛盾しない範囲で `DESIGN.md` を参照する。

Go 版実装ファイルの挙動が `ADLAIRE_CI_SPEC.md` または `ADLAIRE_CI_DETAIL_SPEC.md` と矛盾する場合は、仕様と実装の不整合として扱う。

仕様を変更する場合は、先に該当する仕様書を更新し、その内容に基づいて実装ファイルを更新する。

`DESIGN.md` はデザイン仕様の補助文書であり、Adlaire CI 全体の機能仕様・運用仕様の正本ではない。

`DOCUMENT_INDEX.md` は索引であり、仕様・デザイン・実装判断の正本ではない。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

## Consistency Notes

現時点では、`ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` に記載された一部コンポーネントや機能は仕様化済みだが、リポジトリ内に実装ファイルが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。

現行配置では、`build_spec.go` は Phase 1 実装済みであり、`runner.go` は Phase 2 完了判定パス実装済みである。Go toolchain による `gofmt` と `go test` の検証を完了している。

Go 版コンポーネントの実装状態は、`ADLAIRE_CI_SPEC.md` を正とする。`ADLAIRE_CI_DETAIL_SPEC.md` は検証条件と実装詳細の参照先として扱う。
