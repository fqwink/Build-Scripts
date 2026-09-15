# Build-Scripts Document Index

このファイルは、Build-Scripts リポジトリ内の文書・実装ファイルの参照先と役割を整理する索引である。

## Documents

| ファイル | 役割 |
|---------|------|
| `build_spec_v3_spec.md` | Adlaire CI の仕様正本。`build_spec_v3.py`、`runner.py`、将来コンポーネントである `api_server.py`、`adlaire-ci-sdk.js`、`admin/index.html` の仕様判断の最上位基準。 |
| `DESIGN.md` | `Adlaire-db-spec.html` のデザイン仕様。レイアウト、色、タイポグラフィ、TOC、コードブロック等の視覚仕様を整理する。 |
| `build_spec_v3.py` | Adlaire DB 仕様書 Markdown を単一 HTML へ変換するビルドスクリプト。仕様改訂の入力元ではなく、`build_spec_v3_spec.md` に基づいて更新する実装ファイル。 |
| `runner.py` | GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する CI ランナー。仕様改訂の入力元ではなく、`build_spec_v3_spec.md` に基づいて更新する実装ファイル。 |
| `AGENTS.md` | エージェント作業ルールブック。承認、仕様書管理、実装管理、Git 運用、文書整合の最上位ルール。 |
| `DOCUMENT_INDEX.md` | 本索引。リポジトリ内の文書・実装ファイルの役割と正本関係を示す。仕様正本ではない。 |

## Planned Components

`build_spec_v3_spec.md` では、以下のコンポーネントも仕様化されている。

| パス | 状態 | 役割 |
|------|------|------|
| `api_server.py` | 未実装 | 管理 API サーバー。常駐 HTTP サーバーとして Adlaire CI の状態確認・操作 API を提供する予定。 |
| `adlaire-ci-sdk.js` | 未実装 | 管理ツール用 JavaScript SDK。管理 API 通信を抽象化する予定。 |
| `admin/index.html` | 未実装 | 標準管理ツール UI。SDK 経由で API と通信する予定。 |
| `mcp_server.py` | 将来計画 | MCP サーバー。将来の 4 コンポーネント構成で追加予定。 |

## Source Of Truth

仕様判断では `build_spec_v3_spec.md` を正とする。

デザイン判断では、`build_spec_v3_spec.md` と矛盾しない範囲で `DESIGN.md` を参照する。

`build_spec_v3.py` または `runner.py` の挙動が `build_spec_v3_spec.md` と矛盾する場合は、仕様と実装の不整合として扱う。

仕様を変更する場合は、先に `build_spec_v3_spec.md` を更新し、その内容に基づいて実装ファイルを更新する。

`DESIGN.md` はデザイン仕様の補助文書であり、Adlaire CI 全体の機能仕様・運用仕様の正本ではない。

`DOCUMENT_INDEX.md` は索引であり、仕様・デザイン・実装判断の正本ではない。

`AGENTS.md` と他ファイルが作業ルール上矛盾する場合は、`AGENTS.md` を正とする。

## Consistency Notes

現時点では、`build_spec_v3_spec.md` に記載された一部コンポーネントや機能は仕様化済みだが、リポジトリ内に実装ファイルが存在しない。

仕様化済みだが未実装の内容は、実装済み機能として扱わない。
