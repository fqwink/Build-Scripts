# Build-Scripts Document Index

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

## Documents

| ファイル | 役割 |
|---------|------|
| `ADLAIRE_CI_SPEC.md` | Adlaire CI の方針、ポリシー、実装状態、正本関係を定めるマスター仕様書正本。 |
| `ADLAIRE_CI_DETAIL_SPEC.md` | `ADLAIRE_CI_SPEC.md` の Part 3 詳細仕様。Go 版コンポーネント、標準管理ツール、将来計画コンポーネントの実装詳細に関する正本。 |
| `DESIGN.md` | 生成静的 Web サイトのデザイン仕様。レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を整理する。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| `DOCUMENT_INDEX.md` | 本索引。リポジトリ内の文書・実装ファイルの役割と正本関係を示す。仕様正本ではない。 |

## Specified Components

`ADLAIRE_CI_SPEC.md` と `ADLAIRE_CI_DETAIL_SPEC.md` では、以下のコンポーネントも仕様化されている。

| パス | 状態 | 役割 |
|------|------|------|
| `build_spec.go` | 仕様化済み・未実装 | Go 版静的 Web サイトビルドスクリプト。`adlaire-ci-build` バイナリとして実行する。 |
| `runner.go` | 仕様化済み・未実装 | Go 版 CI ランナー。`adlaire-ci-runner` バイナリとして実行する。 |
| `api_server.go` | 仕様化済み・未実装 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する。 |
| `adlaire-ci-sdk.js` | 仕様化済み・未実装 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する。 |
| `admin/index.html` | 仕様化済み・未実装 | 標準管理ツール UI。SDK 経由で API と通信する。 |
| `mcp_server.go` | 将来計画 | MCP サーバー。将来の 4 コンポーネント構成で追加予定。 |

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

リポジトリ内に既存の `build_spec.py` や `runner.py` が存在する場合でも、`build_spec.go`、`runner.go`、`api_server.go`、`adlaire-ci-sdk.js`、`admin/index.html`、`mcp_server.go` の状態判定には使用しない。

Go 版コンポーネントの実装状態は、`ADLAIRE_CI_SPEC.md`、`ADLAIRE_CI_DETAIL_SPEC.md`、実ファイルの存在、検証結果の整合によって判定する。
