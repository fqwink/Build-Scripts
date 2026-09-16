# Build-Scripts

Build-Scripts は、Adlaire DB 仕様書を HTML 化し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Adlaire DB 仕様書 Markdown から単一 HTML を生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

方針・ポリシー・実装状態の正本は `ADLAIRE_CI_SPEC.md`、実装詳細の正本は `ADLAIRE_CI_DETAIL_SPEC.md` です。

## 主要ファイル

- `ADLAIRE_CI_SPEC.md`: Adlaire CI のマスター仕様書。
- `ADLAIRE_CI_DETAIL_SPEC.md`: Adlaire CI の Part 3 詳細仕様。
- `build_spec.go`: Adlaire DB 仕様書 Markdown を単一 HTML へ変換する Go 版ビルドスクリプト。
- `runner.go`: GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する Go 版 CI ランナー。
- `DESIGN.md`: 生成 HTML のデザイン仕様を整理する補助文書。
- `DOCUMENT_INDEX.md`: 文書・実装ファイルの役割を整理する索引。
- `AGENTS.md`: 本リポジトリにおけるエージェント作業ルール。

## 仕様化済みコンポーネント

Adlaire CI は、最初から Go を前提として仕様策定します。

`build_spec.go`、`runner.go`、`api_server.go`、`adlaire-ci-sdk.js`、`admin/index.html` は仕様化済み・未実装のコンポーネントです。`mcp_server.go` は将来計画のコンポーネントです。

## 注意

仕様や挙動を変更する場合は、先に `ADLAIRE_CI_SPEC.md` または `ADLAIRE_CI_DETAIL_SPEC.md` を改訂し、その内容に基づいて実装ファイルと関連文書を整合させます。
