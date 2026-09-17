# Build-Scripts

Build-Scripts は、Markdown から静的 Web サイトを生成し、Go で実装する Adlaire CI の仕様と運用を管理するためのリポジトリです。

## 概要

本リポジトリでは、Markdown ファイルまたは Markdown ディレクトリから静的 Web サイトを生成する Go 版ビルドスクリプトと、GitHub 上の対象 Markdown 変更を検出してビルドを実行する Go 版 CI ランナーの仕様を管理します。

方針・ポリシー・実装状態・実装可否・ロードマップの正本は `docs/ADLAIRE_CI_SPEC.md` です。詳細仕様の入口、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、詳細節対応表、リポジトリ内ソース配置、横断補足契約は `docs/ADLAIRE_CI_DETAIL_SPEC.md` を参照します。各 owner component の入出力、状態、処理順序、異常系、検証条件の本文は `docs/ADLAIRE_CI_DETAIL_*_SPEC.md` を正本として扱います。

## 読む順番

1. `AGENTS.md`: 作業ルール、承認、Git 運用を確認する。
2. `docs/DOCUMENT_INDEX.md`: 文書と実装ファイルの役割、正本関係、配置を確認する。
3. `docs/ADLAIRE_CI_SPEC.md`: 方針、ポリシー、実装状態、実装可否、ロードマップを確認する。
4. `docs/ADLAIRE_CI_DETAIL_SPEC.md`: 詳細仕様の入口、読み順、共通固定値、対応表、ソース配置を確認する。
5. `docs/ADLAIRE_CI_DETAIL_*_SPEC.md`: 対象 owner component の入出力、状態、処理順序、異常系、検証条件を確認する。

## 主要ファイル

- `docs/ADLAIRE_CI_SPEC.md`: Adlaire CI のマスター仕様書。
- `docs/ADLAIRE_CI_DETAIL_SPEC.md`: Adlaire CI の Part 3 詳細仕様の入口。読み順、索引、共通固定値、実装前確認項目、検証マトリクス、Phase、§0i.1〜§0i.4 の詳細節対応表、§0j のリポジトリ内ソース配置、横断補足契約を持つ。
- `docs/ADLAIRE_CI_DETAIL_*_SPEC.md`: owner component 別の入出力、状態、処理順序、異常系、検証条件の詳細仕様本文。
- `main.go`: 起動入口。実行ファイル名に応じて builder / runner component を呼び出す。
- `components/builder.go`: Markdown を静的 Web サイトへ変換する Go 版ビルドスクリプト。Phase 1 実装済み。
- `components/builder_test.go`: `components/builder.go` の Phase 1 fixture テスト。
- `go.mod`: Go module 定義。外部 module は追加しない。
- `testdata/builder/`: Phase 1 受け入れ fixture。
- `components/runner.go`: GitHub API で対象 Markdown の変更を検出し、ビルドパイプラインを実行する Go 版 CI ランナー。Phase 2 完了判定パス実装済み。
- `components/runner_test.go`: `components/runner.go` の Phase 2 fixture、hardening、完了判定パステスト。
- `docs/DESIGN.md`: 生成静的 Web サイトのデザイン仕様を整理する補助文書。
- `docs/DOCUMENT_INDEX.md`: 文書・実装ファイルの役割を整理する索引。
- `AGENTS.md`: 本リポジトリにおけるエージェント作業ルール。

## 仕様化済みコンポーネント

Adlaire CI は、最初から Go を前提として仕様策定します。

標準ソース配置への移行は完了済みです。`components/builder.go` は Phase 1 実装済み、`components/runner.go` は Phase 2 完了判定パス実装済みです。いずれも `gofmt` と `go test` による検証対象です。

標準ソース配置は `docs/ADLAIRE_CI_SPEC.md` Part 1 §4.3 と `docs/ADLAIRE_CI_DETAIL_SPEC.md` §0j で定義し、`main.go` と `components/*.go` を正とします。

`components/api.go` は実装中・検証未完了のコンポーネントです。

`admin/adlaire-ci-sdk.js`、`admin/index.html` は仕様化済み・未実装のコンポーネントです。`components/mcp.go` は将来計画のコンポーネントです。

## リリース形式

Adlaire CI の標準リリース形式は、GitHub Releases に添付する OS/arch 別の Go 実行バイナリです。初期標準は Linux x86_64 とし、`adlaire-ci-build`、`adlaire-ci-runner`、管理 API 導入後の `adlaire-ci-api` を SHA-256 checksum 検証後に `/usr/local/bin/` へ配置します。

## 注意

仕様や挙動を変更する場合は、先に `docs/ADLAIRE_CI_SPEC.md`、`docs/ADLAIRE_CI_DETAIL_SPEC.md`、または該当する owner component 別の `docs/ADLAIRE_CI_DETAIL_*_SPEC.md` を改訂し、その内容に基づいて実装ファイルと関連文書を整合させます。
